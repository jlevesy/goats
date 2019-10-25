package main

import (
	"fmt"
	"os"

	"github.com/jlevesy/batman/pkg/exec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("expected a path, got nothing")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("unable to open file: %w\n", err)
		os.Exit(1)
	}
	// useless, but not to forget after.
	defer file.Close()

	stmts, err := exec.Parse(file)
	if err != nil {
		fmt.Printf("unable to parse: %w\n", err)
		os.Exit(1)
	}

	insts, err := exec.Resolve(stmts)
	if err != nil {
		fmt.Printf("unable to resolve statements: %w\n", err)
		os.Exit(1)
	}

	if err := exec.Exec(insts); err != nil {
		fmt.Printf("unable to execute instructions: %w\n", err)
		os.Exit(1)
	}
}
