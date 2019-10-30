package vm

import (
	"context"
)

// Test is a collection of instructions (+ some metadata) to execute in order to validate that an application is working.
// TODO: add Setup and Teardown.
type Test struct {
	Name         string
	Instructions []Instruction
}

func (t *Test) Exec(ctx context.Context) (*TestResult, error) {
	tr := TestResult{
		Name: t.Name,
	}

	for _, inst := range t.Instructions {
		if err := inst.Exec(ctx, &tr); err != nil {
			return nil, err
		}

		if tr.LastOutput().Status == ExecStatusFatal {
			// Abort execution if there is a fatal error.
			return &tr, nil
		}
	}

	return &tr, nil
}
