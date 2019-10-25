package exec

import (
	"bufio"
	"io"
)

// Statement is a statement in a program
type Statement string

// Parse parses the content of a given reader into a collection of statements.
func Parse(content io.Reader) ([]Statement, error) {
	scanner := bufio.NewScanner(content)
	var stmts []Statement
	for scanner.Scan() {
		stmts = append(stmts, Statement(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stmts, nil
}
