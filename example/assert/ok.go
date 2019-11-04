package assert

// @instruction{name=assert_ok,builder=NewOK}

import (
	"context"

	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/jlevesy/goats/pkg/testing"
)

// NewOK returns a new instance of the OK instruction.
func NewOK(_ []string) (func(ctx context.Context, t *testing.T), error) {
	return OK, nil
}

// OK asserts that the previous exec instructon has no error.
func OK(ctx context.Context, t *testing.T) {
	execResult, err := instruction.GetExecOutput(t)
	if err != nil {
		t.Fatal()
		return
	}

	if execResult.Err != nil {
		t.Fail()
	}
}
