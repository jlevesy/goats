package exec

import "io"

// Context is the current context of execution.
type State struct {
	ExitCode int
	StdOut   io.Reader
	StdErr   io.Reader
}
