package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/jlevesy/goats/pkg/text"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("expected a path, got nothing")
		os.Exit(1)
	}

	builders := make(instruction.Builders)
	if err := instruction.LoadDynamic([]string{"./example/assert"}, builders); err != nil {
		fmt.Printf("unable to load dynamic instructons: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("unable to open file: %v\n", err)
		os.Exit(1)
	}
	// useless, but not to forget after.
	defer file.Close()

	suite, err := text.NewParser(text.NewLexer(file), builders).Parse()
	if err != nil {
		fmt.Printf("unable to parse suite: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	results, err := suite.Exec(ctx, 1)
	if err != nil {
		fmt.Printf("unable to execute suite: %v\n", err)
		os.Exit(1)
	}

	for _, res := range results.Tests {
		fmt.Printf("Tests %q has status %q\n", res.Name(), res.Status())

		for _, err := range res.Errors() {
			fmt.Println(err)
		}
	}
}
