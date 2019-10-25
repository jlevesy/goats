package exec

import (
	"fmt"
	"reflect"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
)

// InstructionFunc is an helper for instructions.
type InstructionFunc func(prev *State) (*State, error)

func (f InstructionFunc) Exec(prev *State) (*State, error) {
	return f(prev)
}

var Symbols = map[string]map[string]reflect.Value{
	"github.com/jlevesy/batman/pkg/exec": map[string]reflect.Value{
		"State": reflect.ValueOf((*State)(nil)),
	},
}

const src = `
package assert

import (
	"errors"

	"github.com/jlevesy/batman/pkg/exec"
)

func OK(st *exec.State) (*exec.State, error) {
	if st.ExitCode != 0 {
		return nil, errors.New("status code different from 0")
	}

	return st, nil
}
`

func buildAssertOK() (Instruction, error) {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(Symbols)

	_, err := i.Eval(src)
	if err != nil {
		return nil, fmt.Errorf("unable to eval src: %w", err)
	}

	v, err := i.Eval("assert.OK")
	if err != nil {
		panic(err)
	}

	inst, ok := v.Interface().(func(*State) (*State, error))
	if !ok {
		return nil, fmt.Errorf("unable to get a reference to instruction")
	}

	return InstructionFunc(inst), nil
}
