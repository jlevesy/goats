package vm_test

import (
	"context"

	"github.com/jlevesy/goats/pkg/vm"
)

func reportStatus(st vm.ExecStatus) vm.InstructionFunc {
	return func(_ context.Context, tr *vm.TestResult) error {
		tr.Report(st)
		return nil
	}
}

func fail(err error) vm.InstructionFunc {
	return func(_ context.Context, _ *vm.TestResult) error { return err }
}
