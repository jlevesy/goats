package exec

import (
	"fmt"
)

// Instruction is an executable statement.
type Instruction interface {
	Exec(prev *State) (*State, error)
}

// Resolve resolves statements into instructions.
func Resolve(stmts []Statement) ([]Instruction, error) {
	var insts []Instruction

	for _, exp := range stmts {
		if exp == "assert_ok" {
			assertOK, err := buildAssertOK()
			if err != nil {
				return nil, fmt.Errorf("unable to build assert OK inst: %w", err)
			}
			insts = append(insts, assertOK)
			continue
		}

		insts = append(insts, &execInstruction{content: string(exp)})
	}

	return insts, nil
}
