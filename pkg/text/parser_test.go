package text_test

import (
	"testing"

	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/jlevesy/goats/pkg/text"
	"github.com/jlevesy/goats/pkg/vm"
	"github.com/stretchr/testify/assert"
)

type mockScanner struct {
	tokens []*text.Token
}

func (m *mockScanner) Next() *text.Token {
	if len(m.tokens) == 0 {
		return nil
	}

	var next *text.Token
	next, m.tokens = m.tokens[0], m.tokens[1:]

	return next
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []*text.Token
		want    *vm.Suite
		wantErr bool
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
			want: &vm.Suite{
				Tests: []*vm.Test{
					{
						Name: "this is a random test",
						Instructions: []vm.Instruction{
							instruction.NewExec([]string{"ls", "/foo/bar"}),
							instruction.NewExec([]string{"echo", "coucou"}),
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sc := &mockScanner{tokens: test.tokens}
			parser := text.NewParser(sc)
			suite, err := parser.Parse()
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.want, suite)
		})
	}
}
