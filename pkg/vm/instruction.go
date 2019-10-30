package vm

import (
	"context"
	"io"
)

// Execution satuses.
const (
	ExecStatusUnkown ExecStatus = iota
	ExecStatusSuccess
	ExecStatusFailure
	ExecStatusFatal
)

type ExecStatus int

func (e ExecStatus) String() string {
	switch e {
	case ExecStatusSuccess:
		return "success"
	case ExecStatusFailure:
		return "failure"
	case ExecStatusFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// Instruction is an executable statement.
type Instruction interface {
	Exec(ctx context.Context, rt *Runtime) error
}

// InstructionFunc is an helper to create instructions as func.
type InstructionFunc func(ctx context.Context, rt *Runtime) error

// Exec calls the inner function.
func (f InstructionFunc) Exec(ctx context.Context, rt *Runtime) error {
	return f(ctx, rt)
}

// InstructionOutput is the output for an instruction.
type InstructionOutput struct {
	Status ExecStatus
	Output io.ReadSeeker
}
