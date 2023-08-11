// Tideland Go Trace - callstack - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package callstack_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/callstack"
)

//--------------------
// TESTS
//--------------------

// TestHere tests retrieving the callstack in a detailed way at the call position.
func TestHere(t *testing.T) {
	l := callstack.Here()

	Assert(t, Equal(l.Package(), "tideland.dev/go/stew/callstack_test"), "location package")
	Assert(t, Equal(l.File(), "callstack_test.go"), "location file")
	Assert(t, Equal(l.Func(), "TestHere"), "location function")
	Assert(t, Equal(l.Line(), 28), "location line")

	s := callstack.Here().String()

	Assert(t, Equal(s, "(tideland.dev/go/stew/callstack_test:callstack_test.go:TestHere:35)"), "location string")
}

// TestAt tests retrieving the callstack in a detailed way at the depth zero.
func TestAt(t *testing.T) {
	l := callstack.At(0)

	Assert(t, Equal(l.Package(), "tideland.dev/go/stew/callstack_test"), "location package")
	Assert(t, Equal(l.File(), "callstack_test.go"), "location file")
	Assert(t, Equal(l.Func(), "TestAt"), "location function")
	Assert(t, Equal(l.Line(), 42), "location line")

	s := callstack.At(0).String()

	Assert(t, Equal(s, "(tideland.dev/go/stew/callstack_test:callstack_test.go:TestAt:49)"), "location string")
}

// TestStack tests retrieving a call stack.
func TestStack(t *testing.T) {
	st := stackOne()
	s := []string{
		st[0].String(),
		st[1].String(),
		st[2].String(),
		st[3].String(),
		st[4].String(),
	}

	Assert(t, Length(st, 5), "stack length")
	Assert(t, Equal(s[0], "(tideland.dev/go/stew/callstack_test:callstack_test.go:stackFive:148)"), "stack depth 0")
	Assert(t, Equal(s[1], "(tideland.dev/go/stew/callstack_test:callstack_test.go:stackFour:143)"), "stack depth 1")
	Assert(t, Equal(s[2], "(tideland.dev/go/stew/callstack_test:callstack_test.go:stackThree:138)"), "stack depth 2")
	Assert(t, Equal(s[3], "(tideland.dev/go/stew/callstack_test:callstack_test.go:stackTwo:133)"), "stack depth 3")
	Assert(t, Equal(s[4], "(tideland.dev/go/stew/callstack_test:callstack_test.go:stackOne:128)"), "stack depth 4")
}

// TestOffset tests retrieving the callstack with an offset.
func TestOffset(t *testing.T) {
	s := there()

	Assert(t, Equal(s, "(tideland.dev/go/stew/callstack_test:callstack_test.go:TestOffset:75)"), "stack there")

	s = nestedThere()

	Assert(t, Equal(s, "(tideland.dev/go/stew/callstack_test:callstack_test.go:TestOffset:79)"), "stack nested there")

	s = nameless()

	Assert(t, Equal(s, "(tideland.dev/go/stew/callstack_test:callstack_test.go:nameless.func1:121)"), "stack nameless")

	s = callstack.At(-5).String()

	Assert(t, Equal(s, "(tideland.dev/go/stew/callstack_test:callstack_test.go:TestOffset:87)"), "stack offset -5")
}

// TestCache tests the caching of callstacks.
func TestCache(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := nameless()

		Assert(t, Equal(s, "(tideland.dev/go/stew/callstack_test:callstack_test.go:nameless.func1:121)"), "cached stack nameless")
	}
}

//--------------------
// HELPER
//--------------------

// there returns the id at the callstack of the caller.
func there() string {
	return callstack.At(1).String()
}

// nestedThere returns the id at the callstack of the caller but inside a local func.
func nestedThere() string {
	where := func() string {
		return callstack.At(2).String()
	}
	return where()
}

// nameless returns the id from calling a nested nameless function w/o an offset.
func nameless() string {
	noname := func() string {
		return callstack.Here().String()
	}
	return noname()
}

// stackOne is the first one of a stack calling function set.
func stackOne() callstack.Stack {
	return stackTwo()
}

// stackTwo is the second one of a stack calling function set.
func stackTwo() callstack.Stack {
	return stackThree()
}

// stackThree is the third one of a stack calling function set.
func stackThree() callstack.Stack {
	return stackFour()
}

// stackFour is the fourth one of a stack calling function set.
func stackFour() callstack.Stack {
	return stackFive()
}

// stackFive is the fifth one of a stack calling function set.
func stackFive() callstack.Stack {
	return callstack.Dive(5)
}

// EOF
