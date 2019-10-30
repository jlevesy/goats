package exec

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type execInstruction struct {
	content string
}

func (e execInstruction) Exec(prev *State) (*State, error) {
	rawCmd := strings.Split(e.content, " ")
	if len(rawCmd) == 0 {
		return nil, errors.New("empty command")
	}

	cmd := exec.Command(rawCmd[0])
	cmd.Args = rawCmd

	var buf bytes.Buffer

	cmd.Stdout, cmd.Stderr = &buf, &buf

	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		return &State{
			StdOut:   &stdout,
			StdErr:   &stderr,
			ExitCode: exitErr.ExitCode(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &State{
		StdOut:   &stdout,
		StdErr:   &stderr,
		ExitCode: 0,
	}, nil
}
