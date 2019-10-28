package vm

import "errors"

type ExecInstruction struct {
	Cmd []string
}

func (e *ExecInstruction) Exec(r *Runtime) error {
	return errors.New("not impemented")
}
