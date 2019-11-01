package assert

import (
	"context"

	"github.com/jlevesy/goats/pkg/goats"
	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/jlevesy/goats/pkg/testing"
)

// @instruction{name=assert_ok,builder=NewOK}

// OK asserts that the previous exec call is ok.
type OK struct{}

// NewOK returns a new instance of the OK instruction.
func NewOK(_ []string) (goats.Instruction, error) {
	return &OK{}, nil
}

func (a *OK) Exec(ctx context.Context, t *testing.T) {
	execResult, err := instruction.GetExecOutput(t)
	if err != nil {
		t.Fatal()
		return
	}

	if execResult.Err != nil {
		t.Fail()
	}
}
