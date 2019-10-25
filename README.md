## Batman

This project is a PoC of a script interpretor which allows mixing command instructions with instructions written in go and evaluated at runtime using [yaegi](https://github.com/containous/yaegi).

The main use case behind this is to create a command line test framework heavilly inspired from [bats](https://github.com/bats-core/bats-core) but without any kind of bash knowledge required to use it.

## Current status

This is a wild PoC which only shows that it is possible to do something like this.

Right now it applies the following algorithm.
- It reads statements line per line
- Then maps statements to instruction:
  - To an interpreted go function if the statement is `assert_ok`
  - Otherwise, to a command call.

## Using it

With go modules enabled

`go run cmd/main.go assets/test.batman`
