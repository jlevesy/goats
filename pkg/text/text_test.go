package text_test

import (
	"context"

	"github.com/jlevesy/goats/pkg/goats"
	gtesting "github.com/jlevesy/goats/pkg/testing"
	"github.com/jlevesy/goats/pkg/text"
)

type mockResolver func(cmd []string) (goats.Instruction, error)

func (m mockResolver) Resolve(cmd []string) (goats.Instruction, error) {
	return m(cmd)
}

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

func spewerResolver(cmd []string) (goats.Instruction, error) {
	return (&cmdSpewer{cmd: cmd}).Exec, nil
}

type cmdSpewer struct {
	cmd []string
}

func (c *cmdSpewer) Exec(ctx context.Context, t *gtesting.T) {
	var cmds [][]string

	if v := t.GetOutput(commandsKey); v != nil {
		cmds = v.([][]string)
	}

	t.StoreOutput(commandsKey, append(cmds, c.cmd))
}

func instructionsFromSuite(suite *goats.Suite) [][][]string {
	ctx := context.Background()
	var gotCmds [][][]string

	for _, tt := range suite.Tests {
		tr := gtesting.NewT("test")

		for _, inst := range tt.Instructions {
			inst(ctx, tr)
		}

		gotCmds = append(gotCmds, tr.GetOutput(commandsKey).([][]string))
	}

	return gotCmds
}

func testNamesFromSuite(suite *goats.Suite) []string {
	var testNames []string

	for _, t := range suite.Tests {
		testNames = append(testNames, t.Name)
	}

	return testNames
}

const commandsKey = "commands"
