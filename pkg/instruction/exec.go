package instruction

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/jlevesy/goats/pkg/vm"
)

// Exec is an instruction is an instructions that executes a command on the system.
type Exec struct {
	cmd []string
}

// NewExec returns a new instance of Exc
func NewExec(cmd []string) *Exec {
	return &Exec{cmd: cmd}
}

// Exec executes given command.
func (e Exec) Exec(ctx context.Context, t *vm.TestResult) error {
	cmd := exec.CommandContext(ctx, e.cmd[0], e.cmd[1:]...)

	var buf = bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = &buf, &buf

	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		// TODO collect more output.
		// TODO this is shitty, we need to be able to run multiple assertions on the same output.
		// At the moment we can't.
		// We need to have a collection of outputs and a collection of results.
		t.Report(vm.ExecStatusFailure)
		return nil
	}
	if err != nil {
		return err
	}

	t.Report(vm.ExecStatusSuccess)

	return nil
}
