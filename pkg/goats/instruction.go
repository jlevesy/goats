package goats

import (
	"context"

	"github.com/jlevesy/goats/pkg/testing"
)

// Instruction is an executable statement.
type Instruction func(ctx context.Context, t *testing.T)
