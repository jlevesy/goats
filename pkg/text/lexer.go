package text

import (
	"bufio"
	"fmt"
	"io"
)

// TokenType is a type for a token.
type TokenType int32

// Token types.
const (
	TypeUnknown TokenType = iota

	TypeWord
	TypeLineComment
	TypeTestDeclaration
	TypeDoubleQuote

	TypeOpenFunctionBody  // {
	TypeCloseFunctionBody // }

	TypeEOL
	TypeEOF

	TypeError
)

func (t TokenType) String() string {
	switch t {
	case TypeWord:
		return "word"
	case TypeLineComment:
		return "line-comment"
	case TypeTestDeclaration:
		return "test-declaration"
	case TypeDoubleQuote:
		return "double-quote"
	case TypeOpenFunctionBody:
		return "open-function-body"
	case TypeCloseFunctionBody:
		return "close-function-body"
	case TypeEOL:
		return "eol"
	case TypeEOF:
		return "eof"
	case TypeError:
		return "error"
	default:
		return "unknown"
	}
}

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Content string
}

const (
	spaceRune               = ' '
	tabRune                 = '\t'
	eolRune                 = '\n'
	commentLineRune         = '#'
	functionDeclarationRune = '@'
	doubleQuoteRune         = '"'
	openBlockRune           = '{'
	closeBlockRune          = '}'
	escapeNextRune          = '\\'
)

// Useful runesets.
var (
	tabsAndSpaces          = []rune{spaceRune, tabRune}
	tabsSpacesAndComments  = append(tabsAndSpaces, commentLineRune)
	whitespaces            = append(tabsAndSpaces, eolRune)
	whitespacesAndComments = append(whitespaces, commentLineRune)
	endOfWord              = append(whitespaces, doubleQuoteRune, escapeNextRune)
)

type lexerState func(l *Lexer) lexerState

// Lexer is a lexical scanner which produces a sequence of tokens describing the given content.
type Lexer struct {
	content *bufio.Reader
	state   lexerState
	tokenCh chan *Token
}

// NewLexer returns a Lexer.
func NewLexer(content io.Reader) *Lexer {
	l := &Lexer{
		content: bufio.NewReader(content),
		state:   scanText,
		tokenCh: make(chan *Token),
	}

	go l.run()

	return l
}

func (l *Lexer) run() {
	defer close(l.tokenCh)

	for l.state != nil {
		l.state = l.state(l)
	}
}

func (l *Lexer) Next() *Token {
	return <-l.tokenCh
}

func (l *Lexer) emitToken(t TokenType, content string) {
	l.tokenCh <- &Token{Type: t, Content: content}
}

func (l *Lexer) errorf(pattern string, args ...interface{}) {
	l.tokenCh <- &Token{Type: TypeError, Content: fmt.Sprintf(pattern, args...)}
}

func (l *Lexer) readWord() (string, error) {
	var word []rune
	for {
		r, _, err := l.content.ReadRune()
		if err != nil {
			return "", err
		}

		word = append(word, r)

		r, err = l.peekRune()
		if err != nil {
			return "", err
		}

		if contains(endOfWord, r) {
			return string(word), nil
		}
	}
}

func (l *Lexer) skip(skippedSet []rune) error {
	for {
		r, err := l.peekRune()
		if err != nil {
			return err
		}

		if !contains(skippedSet, r) {
			return nil
		}

		if _, _, err = l.content.ReadRune(); err != nil {
			return err
		}

		if r == commentLineRune {
			if err = l.skipComment(); err != nil {
				return err
			}
		}
	}
}

func (l *Lexer) skipComment() error {
	for {
		if err := l.skip(whitespaces); err != nil {
			return err
		}

		_, err := l.readWord()
		if err != nil {
			return err
		}

		next, err := l.peekRune()
		if err != nil {
			return err
		}

		if next == eolRune {
			return nil
		}
	}
}

func (l *Lexer) peekRune() (rune, error) {
	r, _, err := l.content.ReadRune()
	if err != nil {
		return 0, err
	}

	if err = l.content.UnreadRune(); err != nil {
		return 0, err
	}

	return r, nil
}

func scanText(l *Lexer) lexerState {
	for {
		err := l.skip(whitespacesAndComments)
		if err == io.EOF {
			l.emitToken(TypeEOF, "")
			return nil
		}
		if err != nil {
			l.errorf("unable to scan file: %w", err)
			return nil
		}

		r, _, err := l.content.ReadRune()
		if err != nil {
			l.errorf("unable to scan file: %w", err)
			return nil
		}

		switch r {
		case functionDeclarationRune:
			return scanFunctionDeclaration
		default:
			l.errorf("unexpected rune %q", r)
			return nil
		}
	}

	return nil
}

func scanFunctionDeclaration(l *Lexer) lexerState {
	rawFuncType, err := l.readWord()
	if err != nil {
		l.errorf("unable to scan function: %w", err)
		return nil
	}

	funcType, err := functionType(rawFuncType)
	if err != nil {
		l.errorf("unable to scan function: %w", err)
		return nil
	}

	l.emitToken(funcType, "")

	if err = l.skip(whitespaces); err != nil {
		l.errorf("unable to scan function declaration: %w", err)
		return nil
	}

	return scanFunctionName
}

func scanFunctionName(l *Lexer) lexerState {
	if err := l.skip(whitespaces); err != nil {
		l.errorf("unable to scan function declaration: %w", err)
		return nil
	}

	openQuote, _, err := l.content.ReadRune()
	if err != nil {
		l.errorf("unable to scan function: %w", err)
		return nil
	}

	if openQuote != doubleQuoteRune {
		l.errorf("unexpected rune %q, expected %q", openQuote, doubleQuoteRune)
		return nil
	}

	l.emitToken(TypeDoubleQuote, "")

	// Read all the words for function name body.
	for {
		if err := l.skip(whitespaces); err != nil {
			l.errorf("unable to scan function declaration: %w", err)
			return nil
		}

		word, err := l.readWord()
		if err != nil {
			l.errorf("unable to scan function declaration: %w", err)
			return nil
		}

		l.emitToken(TypeWord, word)

		next, err := l.peekRune()
		if err != nil {
			l.errorf("unable to scan function declaration: %w", err)
			return nil
		}

		if next != doubleQuoteRune {
			continue
		}

		// Consuming the closing double quote.
		_, _, err = l.content.ReadRune()
		if err != nil {
			l.errorf("unable to scan function declaration: %w", err)
			return nil
		}

		l.emitToken(TypeDoubleQuote, "")

		return scanFunctionBody
	}
}

func scanFunctionBody(l *Lexer) lexerState {
	if err := l.skip(whitespaces); err != nil {
		l.errorf("unable to scan function declaration: %w", err)
		return nil
	}

	openBlock, _, err := l.content.ReadRune()
	if err != nil {
		l.errorf("unable to scan function: %w", err)
		return nil
	}

	if openBlock != openBlockRune {
		l.errorf("unexpected rune %q, expected %q", openBlock, openBlockRune)
		return nil
	}

	l.emitToken(TypeOpenFunctionBody, "")

	return scanInstruction
}

func scanInstruction(l *Lexer) lexerState {
	if err := l.skip(whitespacesAndComments); err != nil {
		l.errorf("unable to scan instruction: %w", err)
		return nil
	}

	for {
		word, err := l.readWord()
		if err != nil {
			l.errorf("unable to scan instruction: %w", err)
			return nil
		}

		l.emitToken(TypeWord, word)

		if err := l.skip(tabsSpacesAndComments); err != nil {
			l.errorf("unable to scan instruction: %w", err)
			return nil
		}

		next, err := l.peekRune()
		if err != nil {
			l.errorf("unable to scan instruction: %w", err)
			return nil
		}

		// TODO handle escaping of instructions....
		if next != eolRune {
			continue
		}

		break
	}

	l.emitToken(TypeEOL, "")

	if err := l.skip(whitespaces); err != nil {
		l.errorf("unable to scan instruction: %w", err)
		return nil
	}

	next, err := l.peekRune()
	if err != nil {
		l.errorf("unable to scan instruction: %w", err)
		return nil
	}

	if next != closeBlockRune {
		return scanInstruction
	}

	// consuming the closing bracket
	_, _, err = l.content.ReadRune()
	if err != nil {
		l.errorf("unable to scan instruction: %w", err)
		return nil
	}

	l.emitToken(TypeCloseFunctionBody, "")

	return scanText
}

func contains(set []rune, c rune) bool {
	for _, v := range set {
		if c == v {
			return true
		}
	}

	return false
}

func functionType(t string) (TokenType, error) {
	switch t {
	case "test":
		return TypeTestDeclaration, nil
	default:
		return TypeUnknown, fmt.Errorf("unknown function type: %q", t)
	}
}
