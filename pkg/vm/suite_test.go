package vm_test

import (
	"context"
	"testing"

	gtesting "github.com/jlevesy/goats/pkg/testing"
	"github.com/jlevesy/goats/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestSuite_Exec(t *testing.T) {
	tests := []struct {
		name           string
		suite          vm.Suite
		workers        int
		wantStatus     gtesting.Status
		wantResultsLen int
	}{
		{
			name:    "executes all tests in parallel",
			workers: 2,
			suite: vm.Suite{
				Tests: []*vm.Test{
					{
						Name: "test-1",
						Instructions: []vm.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
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
			suite: vm.Suite{
				Tests: []*vm.Test{
					{
						Name: "test-1",
						Instructions: []vm.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
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
			name:    "reports unknown if one test is unknown",
			workers: 1,
			suite: vm.Suite{
				Tests: []*vm.Test{
					{
						Name: "test-1",
						Instructions: []vm.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusUnknown),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusSuccess),
						},
					},
				},
			},
			wantStatus:     gtesting.StatusUnknown,
			wantResultsLen: 2,
		},
		{
			name:    "reports failure if one test is failed",
			workers: 1,
			suite: vm.Suite{
				Tests: []*vm.Test{
					{
						Name: "test-1",
						Instructions: []vm.Instruction{
							reportStatus(gtesting.StatusSuccess),
							reportStatus(gtesting.StatusFailure),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
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
