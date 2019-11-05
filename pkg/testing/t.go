package testing

import (
	"errors"
	"fmt"
)

// T is the test handler for a test.
type T struct {
	name string

	status  Status
	outputs map[interface{}]interface{}
	errs    []error
}

// NewT returns a new test handler.
func NewT(name string) *T {
	return &T{
		name:    name,
		status:  StatusSuccess,
		outputs: make(map[interface{}]interface{}),
	}
}

// Fail marks the test as failed.
func (t *T) Fail(args ...interface{}) {
	t.status = StatusFailure
	t.errs = append(t.errs, errors.New(fmt.Sprintln(args...)))
}

// Failf marks the test as fail.
func (t *T) Failf(format string, args ...interface{}) {
	t.status = StatusFailure
	t.errs = append(t.errs, fmt.Errorf(format, args...))
}

// Fatal marks the test as fatal and stops test execution.
func (t *T) Fatal(args ...interface{}) {
	t.status = StatusFatal
	t.errs = append(t.errs, errors.New(fmt.Sprintln(args...)))
}

// Failf marks the test as fatal and stops the execution.
func (t *T) Fatalf(format string, args ...interface{}) {
	t.status = StatusFatal
	t.errs = append(t.errs, fmt.Errorf(format, args...))
}

func (t *T) Name() string {
	return t.name
}

// Status returns the current test Status.
func (t *T) Status() Status {
	return t.status
}

// Errors return all the errors collected during execution.
func (t *T) Errors() []error {
	return t.errs
}

// StoreOutput stores an instruction output object into the test handler.
func (t *T) StoreOutput(k interface{}, v interface{}) {
	t.outputs[k] = v
}

// GetOutput retrieves an instruction output for the test handler.
func (t *T) GetOutput(k interface{}) interface{} {
	return t.outputs[k]
}
