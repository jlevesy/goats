package text_test

import (
	"testing"

	"github.com/jlevesy/goats/pkg/text"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []*text.Token
		wantCmds [][][]string
		wantErr  bool
	}{
		{
			name: "handles test declaration",
			tokens: []*text.Token{
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
				{Type: text.TypeWord, Content: "echo"},
				{Type: text.TypeWord, Content: "coucou"},
				{Type: text.TypeEOL, Content: ""},
				{Type: text.TypeCloseFunctionBody, Content: ""},
				{Type: text.TypeEOF, Content: ""},
			},
			wantCmds: [][][]string{
				{
					{"ls", "/foo/bar"},
					{"echo", "coucou"},
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sc := &mockScanner{tokens: test.tokens}
			parser := text.NewParser(sc, mockResolver(spewerResolver))

			suite, err := parser.Parse()
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.wantCmds, instructionsFromSuite(suite))
		})
	}
}
