package testing

// T is the test handler for a test.
type T struct {
	Name string

	status  Status
	outputs map[interface{}]interface{}
}

// NewT returns a new test handler.
func NewT(name string) *T {
	return &T{
		Name:    name,
		status:  StatusSuccess,
		outputs: make(map[interface{}]interface{}),
	}
}

// Fail marks the test as failed.
func (t *T) Fail() {
	t.status = StatusFailure
}

// Fatal marks the test as fatal.
func (t *T) Fatal() {
	t.status = StatusFatal
}

// Status returns the current test Status.
func (t *T) Status() Status {
	return t.status
}

// StoreOutput stores an instruction output object into the test handler.
func (t *T) StoreOutput(k interface{}, v interface{}) {
	t.outputs[k] = v
}

// GetOutput retrieves an instruction output for the test handler.
func (t *T) GetOutput(k interface{}) interface{} {
	return t.outputs[k]
}
