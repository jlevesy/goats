package vm

import (
	"context"

	"github.com/jlevesy/goats/pkg/testing"
)

// Instruction is an executable statement.
type Instruction interface {
	Exec(ctx context.Context, t *testing.T)
}

// InstructionFunc is an helper to create instructions as func.
type InstructionFunc func(ctx context.Context, t *testing.T)

// Exec calls the inner function.
func (f InstructionFunc) Exec(ctx context.Context, t *testing.T) {
	f(ctx, t)
}
