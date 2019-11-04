package goats

import (
	"context"

	"github.com/jlevesy/goats/pkg/testing"
)

// Test is a collection of instructions (+ some metadata) to execute in order to validate that an application is working.
// TODO: add Setup and Teardown.
type Test struct {
	Name         string
	Instructions []Instruction
}

func (t *Test) Exec(ctx context.Context) *testing.T {
	tr := testing.NewT(t.Name)

	for _, inst := range t.Instructions {
		inst(ctx, tr)

		if tr.Status() == testing.StatusFatal {
			// Abort execution if there is a fatal error.
			return tr
		}
	}

	return tr
}
