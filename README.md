# Tideland Go Stew

[![GitHub release](https://img.shields.io/github/release/tideland/go-stew.svg)](https://github.com/tideland/go-stew)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-stew/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-stew)](https://github.com/tideland/go-stew/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/stew?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/stew?tab=packages)
[![Workflow](https://github.com/tideland/go-stew/actions/workflows/go.yml/badge.svg)](https://github.com/tideland/go-stew/actions/)

### Description

**Tideland Go Stew** provides a good collection of useful Go packages, like a delicious stew. These are usefult for many purposes. It's a kind of toolbox for the daily work of Go developers.

* `actor` supports easier synchronous and asynchronous concurrent programming based on the actor model
* `callstack` helps diving into the call stack and provides information about the current function
* `capture` allows to capture the output on STDOUT and STDERR; useful for testing
* `dynaj` helps to work with JSON documents without definition of structs first
* `etc` provides a simple configuration management based on JSON files
* `genj` helps working with JSON documents by flexibly using generics
* `jwt` implements the JSON Web Tokens
* `loop` supports the management of channel selection loops in goroutines
* `matcher` matches strings against patterns, like regular expressions, only simpler
* `monitor` supports monitoring variables as well as the runtime of functions
* `qaenv` lets you set environment variables or create temporary directories and files for tests
* `qagen` provides a generator for random data, typially used in tests
* `qaone` enebales convenient one-liner assertions in unit tests based on the standard `testing` package
* `semver` implements the semantic versioning
* `slices` provides a powerful set of functions for slices based on generics
* `timex` provides a set of helpful functions for the work with times
* `uuid` creates and parses UUIDs in the versions 1, 3, 4, and 5
* `wait` helps you to wait for certain conditions by polling; additionally it conaints a throttle

I hope you like it. ;)

### Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

