package instruction_test

import (
	"testing"

	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/jlevesy/goats/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestExecImplements(t *testing.T) {
	assert.Implements(t, (*vm.Instruction)(nil), &instruction.Exec{})
}
