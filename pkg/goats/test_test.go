package goats_test

import (
	"context"
	"testing"

	"github.com/jlevesy/goats/pkg/goats"
	gtesting "github.com/jlevesy/goats/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestTest_Exec(t *testing.T) {
	tests := []struct {
		name       string
		t          goats.Test
		wantStatus gtesting.Status
	}{
		{
			name: "executes all instructions",
			t: goats.Test{
				Instructions: []goats.Instruction{
					reportStatus(gtesting.StatusSuccess),
					reportStatus(gtesting.StatusSuccess),
					reportStatus(gtesting.StatusSuccess),
				},
			},
			wantStatus: gtesting.StatusSuccess,
		},
		{
			name:       "repports success by default",
			t:          goats.Test{},
			wantStatus: gtesting.StatusSuccess,
		},
		{
			name: "continues if failure happens",
			t: goats.Test{
				Instructions: []goats.Instruction{
					reportStatus(gtesting.StatusSuccess),
					reportStatus(gtesting.StatusFailure),
					reportStatus(gtesting.StatusSuccess),
				},
			},
			wantStatus: gtesting.StatusFailure,
		},
		{
			name: "stops if fatal failure happens",
			t: goats.Test{
				Instructions: []goats.Instruction{
					reportStatus(gtesting.StatusSuccess),
					reportStatus(gtesting.StatusFatal),
					reportStatus(gtesting.StatusSuccess),
				},
			},
			wantStatus: gtesting.StatusFatal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			result := test.t.Exec(ctx)
			assert.Equal(t, test.wantStatus, result.Status())
		})
	}
}
