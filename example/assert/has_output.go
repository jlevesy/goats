package assert

// @instruction{name=assert_has_output,builder=NewHasOutput}

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/jlevesy/goats/pkg/testing"
)

// NewHasOutput returns a new instance of the hasOutput instruction.
func NewHasOutput(cmd []string) (func(ctx context.Context, t *testing.T), error) {
	if len(cmd) != 2 {
		return nil, errors.New("expected one argument")
	}

	expected := cmd[1]

	return func(ctx context.Context, t *testing.T) { hasOutput(ctx, t, expected) }, nil
}

func hasOutput(ctx context.Context, t *testing.T, expected string) {
	execResult, err := instruction.GetExecOutput(t)
	if err != nil {
		t.Fatal()
		return
	}

	var found bool

	scanner := bufio.NewScanner(bytes.NewBuffer(execResult.Stdout))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), expected) {
			found = true
			break
		}
	}

	if err = scanner.Err(); err != nil {
		t.Fatal(err)
		return
	}

	if !found {
		t.Failf("output %q does not contains substring %q", string(execResult.Stdout), expected)
	}
}
