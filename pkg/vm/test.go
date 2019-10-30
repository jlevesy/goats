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

		if tr.Status() == ExecStatusFatal {
			// Abort execution if there is a fatal error.
			return &tr, nil
		}
	}

	return &tr, nil
}

// TestResult is the result
type TestResult struct {
	Name     string
	Statuses []ExecStatus
}

func (r *TestResult) Report(st ExecStatus) {
	r.Statuses = append(r.Statuses, st)
}

func (r *TestResult) LastStatus() ExecStatus {
	if len(r.Statuses) == 0 {
		return ExecStatusUnkown
	}

	return r.Statuses[len(r.Statuses)-1]
}

func (r *TestResult) Status() ExecStatus {
	if len(r.Statuses) == 0 {
		return ExecStatusUnkown
	}

	for _, st := range r.Statuses {
		if st == ExecStatusUnkown {
			return ExecStatusUnkown
		}

		if st == ExecStatusFatal {
			return ExecStatusFatal
		}

		if st != ExecStatusSuccess {
			return ExecStatusFailure
		}
	}

	return ExecStatusSuccess
}
