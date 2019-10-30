package vm

import (
	"context"
	"errors"
)

type ExecInstruction struct {
	Cmd []string
}

func (e *ExecInstruction) Exec(ctx context.Context, r *TestResult) error {
	return errors.New("not implemented")
}
