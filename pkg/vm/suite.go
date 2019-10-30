package vm

import (
	"context"
	"fmt"

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
	testOutput := make(chan *TestResult, len(s.Tests))
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

					tr, err := t.Exec(routineCtx)
					if err != nil {
						return err
					}

					testOutput <- tr
				}
			}
		})
	}

	for _, t := range s.Tests {
		testInput <- t
	}

	close(testInput)

	if err := errg.Wait(); err != nil {
		return nil, fmt.Errorf("one or more tests reported a technical failure: %w", err)
	}

	close(testOutput)

	for tr := range testOutput {
		result.Tests = append(result.Tests, tr)
	}

	return &result, nil
}

// SuiteResult represents the result of a suite.
type SuiteResult struct {
	Tests []*TestResult
}

// ExecStatus return the suite status.
func (s *SuiteResult) Status() ExecStatus {
	if len(s.Tests) == 0 {
		return ExecStatusUnkown
	}

	for _, tr := range s.Tests {
		if tr.Status() == ExecStatusUnkown {
			return ExecStatusUnkown
		}

		if tr.Status() != ExecStatusSuccess {
			return ExecStatusFailure
		}
	}

	return ExecStatusSuccess
}
