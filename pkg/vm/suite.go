package vm

import (
	"context"
	"fmt"

	"github.com/jlevesy/goats/pkg/testing"
	"golang.org/x/sync/errgroup"
)

// A Suite is a representation of a goats file. It carries a collection of tests.
// TODO: add SuiteSetup and SuiteTearDown to the suite. Wait for it: SuiteUp and SuiteDown ?
type Suite struct {
	Tests []*Test
}

// Exec executes a suite and returns a SuiteResult.
func (s Suite) Exec(ctx context.Context, workers int) (*SuiteResult, error) {
	var result SuiteResult

	testInput := make(chan *Test, len(s.Tests))
	testOutput := make(chan *testing.T, len(s.Tests))
	errg, routineCtx := errgroup.WithContext(ctx)

	for i := 0; i < workers; i++ {
		errg.Go(func() error {
			for {
				select {
				case <-routineCtx.Done():
					return routineCtx.Err()
				case t, ok := <-testInput:
					if !ok {
						return nil
					}

					testOutput <- t.Exec(routineCtx)
				}
			}
		})
	}

	for _, t := range s.Tests {
		testInput <- t
	}

	close(testInput)

	if err := errg.Wait(); err != nil {
		return nil, fmt.Errorf("unable to execute tests: %w", err)
	}

	close(testOutput)

	for tr := range testOutput {
		result.Tests = append(result.Tests, tr)
	}

	return &result, nil
}

// SuiteResult represents the result of a suite.
type SuiteResult struct {
	Tests []*testing.T
}

// ExecStatus return the suite status.
func (s *SuiteResult) Status() testing.Status {
	if len(s.Tests) == 0 {
		return testing.StatusUnknown
	}

	for _, tr := range s.Tests {
		if tr.Status() == testing.StatusUnknown {
			return testing.StatusUnknown
		}

		if tr.Status() != testing.StatusSuccess {
			return testing.StatusFailure
		}
	}

	return testing.StatusSuccess
}
