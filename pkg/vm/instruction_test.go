package vm_test

import (
	"context"

	"github.com/jlevesy/goats/pkg/vm"
)

func reportOutput(out vm.InstructionOutput) vm.InstructionFunc {
	return func(_ context.Context, rt *vm.Runtime) error {
		rt.Report(out)
		return nil
	}
}

func fail(err error) vm.InstructionFunc {
	return func(_ context.Context, _ *vm.Runtime) error { return err }
}
