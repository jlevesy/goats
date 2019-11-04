package instruction_test

import (
	"context"
	"testing"

	"github.com/jlevesy/goats/pkg/instruction"
	gtesting "github.com/jlevesy/goats/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestExec_Exec(t *testing.T) {
	tests := []struct {
		name         string
		cmd          []string
		wantStatus   gtesting.Status
		wantStdout   []byte
		wantExitCode int
	}{
		{
			name:       "it executes the command and captures stdout",
			cmd:        []string{"echo", "-n", "coucou"},
			wantStatus: gtesting.StatusSuccess,
			wantStdout: []byte("coucou"),
		},
		{
			name:         "it captures the error",
			cmd:          []string{"curl", "lalalalalalala"},
			wantStatus:   gtesting.StatusSuccess,
			wantStdout:   []byte{},
			wantExitCode: 6,
		},
		{
			name:       "it reports fatal on technical error",
			cmd:        []string{"dhdhdhdh"},
			wantStatus: gtesting.StatusFatal,
			wantStdout: []byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			exec := instruction.NewExec(test.cmd)
			gt := gtesting.NewT("test")

			exec.Exec(ctx, gt)

			assert.Equal(t, test.wantStatus, gt.Status())

			if test.wantStatus == gtesting.StatusFatal {
				// Should get an output here !
				return
			}

			out, _ := instruction.GetExecOutput(gt)
			assert.NotNil(t, out)
			assert.Equal(t, test.wantStdout, out.Stdout)

			if out.Err != nil {
				assert.Equal(t, test.wantExitCode, out.Err.ExitCode())
			}
		})
	}
}
