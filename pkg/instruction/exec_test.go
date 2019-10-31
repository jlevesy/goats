package instruction_test

import (
	"testing"

	"github.com/jlevesy/goats/pkg/goats"
	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/stretchr/testify/assert"
)

func TestExecImplements(t *testing.T) {
	assert.Implements(t, (*goats.Instruction)(nil), &instruction.Exec{})
}
