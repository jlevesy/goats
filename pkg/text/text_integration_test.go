package text_test

import (
	"os"
	"testing"

	"github.com/jlevesy/goats/pkg/text"
	"github.com/jlevesy/goats/pkg/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsing(t *testing.T) {
	file, err := os.Open("./assets/test.goats")
	require.NoError(t, err)
	defer file.Close()

	want := &vm.Suite{
		Tests: []*vm.Test{
			{
				Name: "test parsing integration",
				Instructions: []vm.Instruction{
					&vm.ExecInstruction{
						Cmd: []string{"ls", "/foo/bar"},
					},
					&vm.ExecInstruction{
						Cmd: []string{"assert_ok"},
					},
				},
			},
			{
				Name: "test parsing integration 2",
				Instructions: []vm.Instruction{
					&vm.ExecInstruction{
						Cmd: []string{"ls", "/bar/biz"},
					},
					&vm.ExecInstruction{
						Cmd: []string{"assert_ok"},
					},
				},
			},
		},
	}

	parser := text.NewParser(text.NewLexer(file))
	suite, err := parser.Parse()
	assert.NoError(t, err)
	assert.Equal(t, want, suite)
}
