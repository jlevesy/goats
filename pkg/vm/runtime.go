package vm

// Runtime represents some kind of state carried across instructions during executions.
type Runtime struct {
	Outputs []InstructionOutput
}

func (r *Runtime) Report(out InstructionOutput) {
	r.Outputs = append(r.Outputs, out)
}

func (r *Runtime) LastOutput() InstructionOutput {
	if len(r.Outputs) == 0 {
		return InstructionOutput{}
	}

	return r.Outputs[len(r.Outputs)-1]
}

func (r *Runtime) Status() ExecStatus {
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
