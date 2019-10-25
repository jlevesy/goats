package exec

import (
	"fmt"
	"io"
	"os"
)

func Exec(insts []Instruction) error {
	st := &State{}

	for id, inst := range insts {
		var err error

		st, err = inst.Exec(st)
		if err != nil {
			return fmt.Errorf("instruction %d rose an error: %w", id, err)
		}

		io.Copy(os.Stdout, st.StdOut)
		io.Copy(os.Stdout, st.StdErr)
	}

	return nil
}
