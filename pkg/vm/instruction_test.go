package vm_test

import (
	"context"

	"github.com/jlevesy/goats/pkg/testing"
	"github.com/jlevesy/goats/pkg/vm"
)

func reportStatus(st testing.Status) vm.InstructionFunc {
	return func(_ context.Context, t *testing.T) {
		switch st {
		case testing.StatusFatal:
			t.Fatal()
		case testing.StatusFailure:
			t.Fail()
		}
	}
}
