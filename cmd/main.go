package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jlevesy/goats/pkg/text"
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

	suite, err := text.NewParser(text.NewLexer(file)).Parse()
	if err != nil {
		fmt.Printf("unable to parse suite: %w\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	results, err := suite.Exec(ctx, 1)
	if err != nil {
		fmt.Printf("unable to execute suite: %w\n", err)
		os.Exit(1)
	}

	for _, res := range results.Tests {
		fmt.Printf("Tests %q has status %q", res.Name, res.Status())
	}
}
