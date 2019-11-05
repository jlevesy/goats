package instruction

import (
	"bytes"
	"context"
	"errors"
	"os/exec"

	"github.com/jlevesy/goats/pkg/testing"
)

// Exec is an instruction is an instructions that executes a command on the system.
type Exec struct {
	cmd []string
}

// NewExec returns a new instance of an Exec instruction.
func NewExec(cmd []string) *Exec {
	return &Exec{cmd: cmd}
}

// Exec executes given command.
func (e Exec) Exec(ctx context.Context, t *testing.T) {
	cmd := exec.CommandContext(ctx, e.cmd[0], e.cmd[1:]...)

	var stdout = bytes.Buffer{}
	cmd.Stdout = &stdout

	err := cmd.Run()
	if exit, ok := err.(*exec.ExitError); ok {
		t.StoreOutput(execOutputStoreKey, &ExecOutput{Stdout: stdout.Bytes(), Err: exit})
		return
	}
	if err != nil {
		t.Fatal()
		return
	}

	t.StoreOutput(execOutputStoreKey, &ExecOutput{Stdout: stdout.Bytes()})
}

type execOutputKey string

const (
	execOutputStoreKey execOutputKey = "k"
)

// ExecOutput is an output.
type ExecOutput struct {
	Stdout []byte
	Err    *exec.ExitError
}

// GetExecOutput returns the exec output.
func GetExecOutput(t *testing.T) (*ExecOutput, error) {
	v := t.GetOutput(execOutputStoreKey)
	if v == nil {
		return nil, errors.New("no registered ExecOutput")
	}

	out, ok := v.(*ExecOutput)
	if !ok {
		return nil, errors.New("unable to convert stored value to ExecOutput")
	}

	return out, nil
}
