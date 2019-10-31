## Goats

Bats + Go = Goats

This project is a PoC of a script interpretor which allows mixing command instructions with instructions written in go and evaluated at runtime using [yaegi](https://github.com/containous/yaegi).

The main use case behind this is to create a command line test framework heavily inspired from [bats](https://github.com/bats-core/bats-core) but without any kind of bash knowledge required to use it.

## Current status

This is a wild PoC which only shows that it is possible to do something like this.

At the moment it is able to run a set of commands on the system and capture their output.

## Using it

With go modules enabled

`go run cmd/main.go assets/test.goats`
