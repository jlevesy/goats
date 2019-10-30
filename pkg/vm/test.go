package vm

import (
	"context"
	"errors"
)

var ErrFatal = errors.New("fatal instruction error")

// Test is a collection of instructions (+ some metadata) to execute in order to validate that an application is working.
// TODO: add Setup and Teardown.
type Test struct {
	Name         string
	Instructions []Instruction
}

func (t *Test) Exec(ctx context.Context) (*Runtime, error) {
	var rt Runtime

	for _, inst := range t.Instructions {
		if err := inst.Exec(ctx, &rt); err != nil {
			return nil, err
		}

		if rt.LastOutput().Status == ExecStatusFatal {
			// Abort execution if there is a fatal error.
			return &rt, nil
		}
	}

	return &rt, nil
}
