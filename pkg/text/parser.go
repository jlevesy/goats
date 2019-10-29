package text

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jlevesy/goats/pkg/vm"
)

// Scanner scans for token.
type Scanner interface {
	Next() *Token
}

type parserState func(*Parser) (parserState, error)

// Parser transforms lexer tokens into an actual vm TestSuite.
type Parser struct {
	lexer Scanner
	state parserState

	suite  *vm.Suite
	testID int32
}

func NewParser(l Scanner) *Parser {
	return &Parser{
		lexer:  l,
		state:  parseSuite,
		suite:  &vm.Suite{},
		testID: -1, // This is not good !
	}
}

// Parse returns the *vm.Suite
func (p *Parser) Parse() (*vm.Suite, error) {
	var err error

	for {
		p.state, err = p.state(p)
		if err != nil {
			return nil, fmt.Errorf("unable to parse suite: %w", err)
		}

		if p.state == nil {
			return p.suite, nil
		}
	}
}

func (p *Parser) nextToken() (*Token, error) {
	tok := p.lexer.Next()
	if tok == nil {
		return nil, nil
	}

	if tok.Type == TypeError {
		return nil, errors.New(tok.Content)
	}

	return tok, nil
}

func (p *Parser) requireNextToken() (*Token, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if tok == nil {
		return nil, errors.New("received nil token")
	}

	return tok, nil
}

func (p *Parser) requireNextTokenWithType(allowedTypes ...TokenType) (*Token, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	for _, t := range allowedTypes {
		if tok.Type == t {
			return tok, nil
		}
	}

	return nil, fmt.Errorf("received unexpected token %q, expected %q", tok.Type, allowedTypes)
}

func parseSuite(p *Parser) (parserState, error) {
	tok, err := p.requireNextToken()
	if err != nil {
		return nil, fmt.Errorf("unable to parse suite: %w", err)
	}

	switch tok.Type {
	case TypeTestDeclaration:
		return parseTestName, nil
	case TypeEOF:
		// Normal termination.
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token type %q", tok.Type)
	}
}

func parseTestName(p *Parser) (parserState, error) {
	// Consume the first double quote.
	tok, err := p.requireNextTokenWithType(TypeDoubleQuote)
	if err != nil {
		return nil, fmt.Errorf("unable to parse test name: %w", err)
	}

	var words []string

	for {
		tok, err = p.requireNextTokenWithType(TypeWord, TypeDoubleQuote)
		if err != nil {
			return nil, fmt.Errorf("unable to parse test name: %w", err)
		}

		if tok.Type == TypeDoubleQuote {
			break
		}

		words = append(words, tok.Content)
	}

	t := vm.Test{
		Name: strings.Join(words, string(spaceRune)),
	}

	p.suite.Tests = append(p.suite.Tests, &t)
	p.testID++

	return parseTestBody, nil
}

func parseTestBody(p *Parser) (parserState, error) {
	tok, err := p.requireNextTokenWithType(TypeOpenFunctionBody)
	if err != nil {
		return nil, fmt.Errorf("unable to parse test body: %w", err)
	}

	var (
		instructions       []vm.Instruction
		currentInstruction []string
	)
	for {
		tok, err = p.requireNextTokenWithType(TypeWord, TypeEOL, TypeCloseFunctionBody)
		if err != nil {
			return nil, fmt.Errorf("unable to parse test body: %w", err)
		}

		switch tok.Type {
		case TypeWord:
			currentInstruction = append(currentInstruction, tok.Content)
		case TypeEOL:
			// TODO => resolve instruction here
			instructions = append(instructions, &vm.ExecInstruction{Cmd: currentInstruction})
			currentInstruction = nil
		case TypeCloseFunctionBody:
			p.suite.Tests[p.testID].Instructions = instructions
			return parseSuite, nil
		}
	}
}
