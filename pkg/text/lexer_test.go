package text_test

import (
	"bytes"
	"testing"

	"github.com/jlevesy/goats/pkg/text"
	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []*text.Token
	}{
		{
			name: "handles comment",
			text: "# hello this is a comment\n",
			want: []*text.Token{
				{Type: text.TypeLineComment, Content: ""},
				{Type: text.TypeWord, Content: "hello"},
				{Type: text.TypeWord, Content: "this"},
				{Type: text.TypeWord, Content: "is"},
				{Type: text.TypeWord, Content: "a"},
				{Type: text.TypeWord, Content: "comment"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeEOF, Content: ""},
			},
		},
		{
			name: "handles test declaration",
			text: `
@test "this is a random test" {
	ls /foo/bar
	assert_ok
}
			`,
			want: []*text.Token{
				{Type: text.TypeTestDeclaration, Content: ""},
				{Type: text.TypeDoubleQuote, Content: ""},
				{Type: text.TypeWord, Content: "this"},
				{Type: text.TypeWord, Content: "is"},
				{Type: text.TypeWord, Content: "a"},
				{Type: text.TypeWord, Content: "random"},
				{Type: text.TypeWord, Content: "test"},
				{Type: text.TypeDoubleQuote, Content: ""},
				{Type: text.TypeOpenBlock, Content: ""},
				{Type: text.TypeWord, Content: "ls"},
				{Type: text.TypeWord, Content: "/foo/bar"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeWord, Content: "assert_ok"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeCloseBlock, Content: ""},
				{Type: text.TypeEOF, Content: ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := bytes.NewBuffer([]byte(test.text))
			lexer := text.NewLexer(in)

			var tokens []*text.Token

			for {
				tok, more := lexer.Next()
				if !more {
					break
				}
				tokens = append(tokens, tok)
			}

			assert.Equal(t, test.want, tokens)
		})
	}
}
