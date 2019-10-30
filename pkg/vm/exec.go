package vm

import (
	"context"
	"errors"
)

type ExecInstruction struct {
	Cmd []string
}

func (e *ExecInstruction) Exec(ctx context.Context, r *Runtime) error {
	return errors.New("not impemented")
}
