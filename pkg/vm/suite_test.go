package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jlevesy/goats/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestSuite_Exec(t *testing.T) {
	tests := []struct {
		name           string
		suite          vm.Suite
		workers        int
		wantErr        bool
		wantStatus     vm.ExecStatus
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
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusSuccess),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusSuccess),
						},
					},
				},
			},
			wantStatus:     vm.ExecStatusSuccess,
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
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusSuccess),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusSuccess),
						},
					},
				},
			},
			wantStatus:     vm.ExecStatusUnkown,
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
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusUnkown),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusSuccess),
						},
					},
				},
			},
			wantStatus:     vm.ExecStatusUnkown,
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
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusFailure),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusSuccess),
						},
					},
				},
			},
			wantStatus:     vm.ExecStatusFailure,
			wantResultsLen: 2,
		},
		{
			name:    "raise an error if a test reports an error",
			workers: 1,
			suite: vm.Suite{
				Tests: []*vm.Test{
					{
						Name: "test-1",
						Instructions: []vm.Instruction{
							reportStatus(vm.ExecStatusSuccess),
						},
					},
					{
						Name: "test-2",
						Instructions: []vm.Instruction{
							reportStatus(vm.ExecStatusSuccess),
							reportStatus(vm.ExecStatusSuccess),
							fail(errors.New("nope")),
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := test.suite.Exec(ctx, test.workers)

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Len(t, result.Tests, test.wantResultsLen)
			assert.Equal(t, test.wantStatus, result.Status())
		})
	}
}
