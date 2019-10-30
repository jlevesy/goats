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
		name           string
		t              vm.Test
		wantErr        bool
		wantResultsLen int
		wantStatus     vm.ExecStatus
	}{
		{
			name: "executes all instructions",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
				},
			},
			wantErr:        false,
			wantResultsLen: 3,
			wantStatus:     vm.ExecStatusSuccess,
		},
		{
			name:           "repports unknown if no insructions",
			t:              vm.Test{},
			wantErr:        false,
			wantResultsLen: 0,
			wantStatus:     vm.ExecStatusUnkown,
		},
		{
			name: "continues if failure happens",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusFailure}),
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
				},
			},
			wantErr:        false,
			wantResultsLen: 3,
			wantStatus:     vm.ExecStatusFailure,
		},
		{
			name: "stops if fatal failure happens",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusFatal}),
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
				},
			},
			wantErr:        false,
			wantResultsLen: 2,
			wantStatus:     vm.ExecStatusFatal,
		},
		{
			name: "stops if instruction reports a technical error",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
					reportOutput(vm.InstructionOutput{Status: vm.ExecStatusSuccess}),
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

			assert.Len(t, result.Outputs, test.wantResultsLen)
			assert.NoError(t, err)
			assert.Equal(t, test.wantStatus, result.Status())
		})
	}
}
