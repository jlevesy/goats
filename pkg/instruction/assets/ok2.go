package assert

import (
	"context"

	"github.com/jlevesy/goats/pkg/testing"
)

// @instruction{name=assert_ok_2,builder=NewOK2}
// NewOK2 returns a new instance of the OK instruction.
func NewOK2(_ []string) (func(ctx context.Context, t *testing.T), error) {
	return OK2, nil
}

func OK2(ctx context.Context, t *testing.T) {
	ok(ctx, t)
}
