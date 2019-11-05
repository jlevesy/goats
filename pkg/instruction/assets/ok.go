package assert

import (
	"context"
	"errors"

	"github.com/jlevesy/goats/pkg/testing"
)

// @instruction{name=assert_ok,builder=NewOK}
// NewOK returns a new instance of the OK instruction.
func NewOK(_ []string) (func(ctx context.Context, t *testing.T), error) {
	return OK, nil
}

// OK asserts that the previous exec instructon has no error.
func OK(ctx context.Context, t *testing.T) {
	ok(ctx, t)
}

// @instruction{name=fail,builder=NewFail}
// NewFail fails.
func NewFail(_ []string) (func(ctx context.Context, t *testing.T), error) {
	// TODO: returning nil, errors.New("nope") makes Yaegi panic.
	return func(ctx context.Context, t *testing.T) {}, errors.New("nope")
}
