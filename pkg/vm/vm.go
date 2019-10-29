package vm

// A Suite is a representation of a goats file. It carries a collection of tests.
// TODO: add SuiteSetup and SuiteTearDown to the suite. Wait for it: SuiteUp and SuiteDown ?
type Suite struct {
	Tests []*Test
}

// Test is a collection of instructions (+ some metadata) to execute in order to validate that an application is working.
// TODO: add Setup and Teardown.
type Test struct {
	Name         string
	Instructions []Instruction
}

// Runtime represents some kind of state carried across instructions during executions.
type Runtime struct {
}

// Instruction is an executable statement.
type Instruction interface {
	Exec(prev *Runtime) error
}
