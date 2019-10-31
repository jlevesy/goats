package goats_test

import (
	"context"
	"testing"

	"github.com/jlevesy/goats/pkg/goats"
	gtesting "github.com/jlevesy/goats/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestSuite_Exec(t *testing.T) {
	tests := []struct {
		name           string
		suite          goats.Suite
		workers        int
		wantStatus     gtesting.Status
		wantResultsLen int
	}{
		{
			name:    "executes all tests in parallel",
			workers: 2,
			suite: goats.Suite{
				Tests: []*goats.Test{
					{
						Name: "test-1",
						Instructions: []goats.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
					{
						Name: "test-2",
						Instructions: []goats.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
				},
			},
			wantStatus:     gtesting.StatusSuccess,
			wantResultsLen: 2,
		},
		{
			name:    "reports unknown if there is no workers registered",
			workers: 0,
			suite: goats.Suite{
				Tests: []*goats.Test{
					{
						Name: "test-1",
						Instructions: []goats.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
					{
						Name: "test-2",
						Instructions: []goats.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
				},
			},
			wantStatus:     gtesting.StatusUnknown,
			wantResultsLen: 0,
		},
		{
			name:    "reports failure if one test is failed",
			workers: 1,
			suite: goats.Suite{
				Tests: []*goats.Test{
					{
						Name: "test-1",
						Instructions: []goats.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusFailure),
						},
					},
					{
						Name: "test-2",
						Instructions: []goats.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
				},
			},
			wantStatus:     gtesting.StatusFailure,
			wantResultsLen: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			result, _ := test.suite.Exec(ctx, test.workers)
			assert.Equal(t, test.wantStatus, result.Status())
		})
	}
}
