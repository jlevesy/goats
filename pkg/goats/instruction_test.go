package goats_test

import (
	"context"

	"github.com/jlevesy/goats/pkg/goats"
	"github.com/jlevesy/goats/pkg/testing"
)

func reportStatus(st testing.Status) goats.Instruction {
	return func(_ context.Context, t *testing.T) {
		switch st {
		case testing.StatusFatal:
			t.Fatal()
		case testing.StatusFailure:
			t.Fail()
		}
	}
}
