package text

import (
	"errors"
	"fmt"

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

	suite *vm.Suite
}

func NewParser(l Scanner) *Parser {
	return &Parser{
		lexer: l,
		state: parseSuite,

		suite: &vm.Suite{},
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

func parseSuite(p *Parser) (parserState, error) {
	return nil, errors.New("not implemented")
}
