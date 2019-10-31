package vm_test

import (
	"context"
	"testing"

	gtesting "github.com/jlevesy/goats/pkg/testing"
	"github.com/jlevesy/goats/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestTest_Exec(t *testing.T) {
	tests := []struct {
		name       string
		t          vm.Test
		wantStatus gtesting.Status
	}{
		{
			name: "executes all instructions",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportStatus(gtesting.StatusSuccess),
					reportStatus(gtesting.StatusSuccess),
					reportStatus(gtesting.StatusSuccess),
				},
			},
			wantStatus: gtesting.StatusSuccess,
		},
		{
			name:       "repports success by default",
			t:          vm.Test{},
			wantStatus: gtesting.StatusSuccess,
		},
		{
			name: "continues if failure happens",
			t: vm.Test{
				Instructions: []vm.Instruction{
					reportStatus(gtesting.StatusSuccess),
					reportStatus(gtesting.StatusFailure),
					reportStatus(gtesting.StatusSuccess),
				},
			},
			wantStatus: gtesting.StatusFailure,
		},
		{
			name: "stops if fatal failure happens",
			t: vm.Test{
				Instructions: []vm.Instruction{
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
