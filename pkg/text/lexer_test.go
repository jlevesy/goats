package text_test

import (
	"bytes"
	"testing"

	"github.com/jlevesy/goats/pkg/text"
	"github.com/stretchr/testify/assert"
)

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []*text.Token
	}{
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
				{Type: text.TypeOpenFunctionBody, Content: ""},
				{Type: text.TypeWord, Content: "ls"},
				{Type: text.TypeWord, Content: "/foo/bar"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeWord, Content: "assert_ok"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeCloseFunctionBody, Content: ""},
				{Type: text.TypeEOF, Content: ""},
			},
		},
		{
			name: "handles very commented test declaration",
			text: `
# Hello this is a comment
@test "this is a random test" { # Hey I'm commenting here too
	ls /foo/bar # ouh funny commenting here
	assert_ok
} # WHY COMMENTING HERE ?
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
				{Type: text.TypeOpenFunctionBody, Content: ""},
				{Type: text.TypeWord, Content: "ls"},
				{Type: text.TypeWord, Content: "/foo/bar"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeWord, Content: "assert_ok"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeCloseFunctionBody, Content: ""},
				{Type: text.TypeEOF, Content: ""},
			},
		},
		{
			name: "handles multiline instructions",
			text: `
@test "this is a random test" {
	ls /foo/bar \
		biz \
		buz \
		bar
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
				{Type: text.TypeOpenFunctionBody, Content: ""},
				{Type: text.TypeWord, Content: "ls"},
				{Type: text.TypeWord, Content: "/foo/bar"},
				{Type: text.TypeWord, Content: "biz"},
				{Type: text.TypeWord, Content: "buz"},
				{Type: text.TypeWord, Content: "bar"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeWord, Content: "assert_ok"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeCloseFunctionBody, Content: ""},
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
				tok := lexer.Next()
				if tok == nil {
					break
				}
				tokens = append(tokens, tok)
			}

			assert.Equal(t, test.want, tokens)
		})
	}
}
