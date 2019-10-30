package vm

// TestResult represents the current state of the execution.
type TestResult struct {
	Name    string
	Outputs []InstructionOutput
}

func (r *TestResult) Report(out InstructionOutput) {
	r.Outputs = append(r.Outputs, out)
}

func (r *TestResult) LastOutput() InstructionOutput {
	if len(r.Outputs) == 0 {
		return InstructionOutput{}
	}

	return r.Outputs[len(r.Outputs)-1]
}

func (r *TestResult) Status() ExecStatus {
	if len(r.Outputs) == 0 {
		return ExecStatusUnkown
	}

	for _, out := range r.Outputs {
		if out.Status == ExecStatusUnkown {
			return ExecStatusUnkown
		}

		if out.Status == ExecStatusFatal {
			return ExecStatusFatal
		}

		if out.Status != ExecStatusSuccess {
			return ExecStatusFailure
		}
	}

	return ExecStatusSuccess
}
