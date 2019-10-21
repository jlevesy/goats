package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Context is the current context of execution.
type ExecState struct {
	ExitCode int
	StdOut   io.Reader
	StdErr   io.Reader
}

// Instruction is something you can execute.
type Instruction interface {
	Exec(prev *ExecState) (*ExecState, error)
}

type Expression string

func parse(content io.Reader) ([]Expression, error) {
	scanner := bufio.NewScanner(content)
	var exps []Expression
	for scanner.Scan() {
		exps = append(exps, Expression(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return exps, nil
}

type execInstruction struct {
	content string
}

func (e execInstruction) Exec(prev *ExecState) (*ExecState, error) {
	rawCmd := strings.Split(e.content, " ")
	if len(rawCmd) == 0 {
		return nil, errors.New("empty command")
	}

	cmd := exec.Command(rawCmd[0])

	if len(rawCmd) > 1 {
		cmd.Args = rawCmd
	}
	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Println("executing", e.content, *cmd)
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		return &ExecState{
			StdOut:   &stdout,
			StdErr:   &stderr,
			ExitCode: exitErr.ExitCode(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &ExecState{
		StdOut:   &stdout,
		StdErr:   &stderr,
		ExitCode: 0,
	}, nil

}

func resolve(exps []Expression) ([]Instruction, error) {
	var insts []Instruction

	for _, exp := range exps {
		insts = append(insts, &execInstruction{content: string(exp)})
	}

	return insts, nil
}

func execute(insts []Instruction) error {
	st := &ExecState{}

	for _, instruction := range insts {
		var err error

		st, err = instruction.Exec(st)
		if err != nil {
			return fmt.Errorf("unable to execute instruction %d: %w", err)
		}

		io.Copy(os.Stdout, st.StdOut)
		io.Copy(os.Stdout, st.StdErr)
	}

	return nil
}

func main() {
	file, err := os.Open("./test.batman")
	if err != nil {
		fmt.Printf("unable to open file: %w\n", err)
		os.Exit(1)
	}
	// useless, but not to forget after.
	defer file.Close()

	exps, err := parse(file)
	if err != nil {
		fmt.Printf("unable to parse: %w\n", err)
		os.Exit(1)
	}

	insts, err := resolve(exps)
	if err != nil {
		fmt.Printf("unable to resolve expressions to known instructions: %w\n", err)
		os.Exit(1)
	}

	if err := execute(insts); err != nil {
		fmt.Printf("unable to execute instructions: %w\n", err)
		os.Exit(1)
	}
}
