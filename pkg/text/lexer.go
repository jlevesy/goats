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
	TypeCommentLine
	TypeTestDeclaration
	TypeDoubleQuote

	TypeOpenBlock  // {
	TypeCloseBlock // }

	TypeEOL
	TypeEOF

	TypeError
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Content string
}

const (
	whitespace     = ' '
	commentLine    = '#'
	functionMarker = '@'
	tab            = '\t'
	eol            = '\n'
	eof            = rune(0)
)

type stateFunc func(l *Lexer) stateFunc

// Lexer is a lexical scanner which returns lexical items based on given parsed text.
type Lexer struct {
	content *bufio.Reader
	state   stateFunc
	tokenCh chan *Token
}

// NewLexer returns a Lexer.
func NewLexer(content io.Reader) *Lexer {
	l := &Lexer{
		content: bufio.NewReader(content),
		state:   parseText,
		tokenCh: make(chan *Token),
	}

	go l.run()

	return l
}

func (l *Lexer) run() {
	for l.state != nil {
		l.state = l.state(l)
	}

	close(l.tokenCh)
}

func (l *Lexer) Next() (*Token, bool) {
	tok, ok := <-l.tokenCh
	return tok, ok
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

		if isWhitespace(r) {
			return string(word), nil
		}
	}
}

func (l *Lexer) skipWhitespace() error {
	for {
		r, err := l.peekRune()
		if err != nil {
			return err
		}

		if !isWhitespace(r) {
			return nil
		}

		if _, _, err = l.content.ReadRune(); err != nil {
			return err
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

func parseText(l *Lexer) stateFunc {
	for {
		r, _, err := l.content.ReadRune()
		// Standard termination.
		if err == io.EOF {
			l.emitToken(TypeEOF, "")
			return nil
		}
		if err != nil {
			l.errorf("unable to read rune %w", err)
			return nil
		}

		if isWhitespace(r) {
			continue
		}

		switch r {
		case commentLine:
			l.emitToken(TypeCommentLine, string(r))
			return parseComment
		case functionMarker:
			return parseFunction
		default:
			l.errorf("unexpected rune %q", r)
			return nil
		}
	}

	return nil
}

func parseComment(l *Lexer) stateFunc {
	if err := l.skipWhitespace(); err != nil {
		l.errorf("unable to pase comment: %w", err)
		return nil
	}

	for {
		word, err := l.readWord()
		if err != nil {
			l.errorf("unable to parse comment: %w", err)
			return nil
		}

		l.emitToken(TypeWord, word)

		sep, _, err := l.content.ReadRune()
		if err == io.EOF {
			l.emitToken(TypeEOF, "")
			return nil
		}
		if err != nil {
			l.errorf("unable to parse comment: %w", err)
			return nil
		}

		if sep == eol {
			l.emitToken(TypeEOL, "")
			return parseText
		}
	}
}

func parseFunction(l *Lexer) stateFunc {
	rawFuncType, err := l.readWord()
	if err != nil {
		l.errorf("unable to parse function: %w", err)
		return nil
	}

	funcType, err := parseFunctionType(rawFuncType)
	if err != nil {
		l.errorf("unable to parse function: %w", err)
		return nil
	}

	l.emitToken(funcType, "")

	// TODO

	return parseText
}

func isWhitespace(ch rune) bool {
	return ch == whitespace || ch == tab || ch == eol
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func parseFunctionType(t string) (TokenType, error) {
	switch t {
	case "test":
		return TypeTestDeclaration, nil
	default:
		return TypeUnknown, fmt.Errorf("unknown function type: %q", t)
	}
}
