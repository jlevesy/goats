package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jlevesy/goats/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestTest_Exec(t *testing.T) {
	tests := []struct {
		name            string
		t               vm.Test
		wantErr         bool
		wantStatusesLen int
		wantStatus      vm.ExecStatus
	}{
		{
			name: "executes all instructions",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportStatus(vm.ExecStatusSuccess),
					reportStatus(vm.ExecStatusSuccess),
					reportStatus(vm.ExecStatusSuccess),
				},
			},
			wantErr:         false,
			wantStatusesLen: 3,
			wantStatus:      vm.ExecStatusSuccess,
		},
		{
			name:            "repports unknown if no insructions",
			t:               vm.Test{},
			wantErr:         false,
			wantStatusesLen: 0,
			wantStatus:      vm.ExecStatusUnkown,
		},
		{
			name: "continues if failure happens",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportStatus(vm.ExecStatusSuccess),
					reportStatus(vm.ExecStatusFailure),
					reportStatus(vm.ExecStatusSuccess),
				},
			},
			wantErr:         false,
			wantStatusesLen: 3,
			wantStatus:      vm.ExecStatusFailure,
		},
		{
			name: "stops if fatal failure happens",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportStatus(vm.ExecStatusSuccess),
					reportStatus(vm.ExecStatusFatal),
					reportStatus(vm.ExecStatusSuccess),
				},
			},
			wantErr:         false,
			wantStatusesLen: 2,
			wantStatus:      vm.ExecStatusFatal,
		},
		{
			name: "stops if instruction reports a technical error",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportStatus(vm.ExecStatusSuccess),
					reportStatus(vm.ExecStatusSuccess),
					fail(errors.New("nope")),
				},
			},
			wantErr:    true,
			wantStatus: vm.ExecStatusFatal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := test.t.Exec(ctx)

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Len(t, result.Statuses, test.wantStatusesLen)
			assert.NoError(t, err)
			assert.Equal(t, test.wantStatus, result.Status())
		})
	}
}
