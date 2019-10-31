package text_test

import (
	"os"
	"testing"

	"github.com/jlevesy/goats/pkg/goats"
	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/jlevesy/goats/pkg/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsing(t *testing.T) {
	file, err := os.Open("./assets/test.goats")
	require.NoError(t, err)
	defer file.Close()

	want := &goats.Suite{
		Tests: []*goats.Test{
			{
				Name: "test parsing integration",
				Instructions: []goats.Instruction{
					instruction.NewExec([]string{"ls", "/foo/bar", "is", "a", "long", "instruction"}),
					instruction.NewExec([]string{"assert_ok"}),
				},
			},
			{
				Name: "test parsing integration 2",
				Instructions: []goats.Instruction{
					instruction.NewExec([]string{"ls", "/bar/biz"}),
					instruction.NewExec([]string{"assert_ok"}),
				},
			},
		},
	}

	parser := text.NewParser(text.NewLexer(file))
	suite, err := parser.Parse()
	assert.NoError(t, err)
	assert.Equal(t, want, suite)
}
