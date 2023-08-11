// Tideland Go Stew - Capture - Unit Tests
//
// Copyright (C) 2017-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package capture_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"os"
	"testing"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/capture"
)

//--------------------
// TESTS
//--------------------

// TestStdout tests the capturing of writings to stdout.
func TestStdout(t *testing.T) {
	hello := "Hello, World!"
	cptrd := capture.Stdout(func() {
		fmt.Print(hello)
	})
	Assert(t, Equal(cptrd.String(), hello), "captured stdout is not equal to expected one")
	Assert(t, Length(cptrd, len(hello)), "captured stdout has wrong length")
}

// TestStderr tests the capturing of writings to stderr.
func TestStderr(t *testing.T) {
	ouch := "ouch"
	cptrd := capture.Stderr(func() {
		fmt.Fprint(os.Stderr, ouch)
	})
	Assert(t, Equal(cptrd.String(), ouch), "captured stderr is not equal to expected one")
	Assert(t, Length(cptrd, len(ouch)), "captured stderr has wrong length")
}

// TestBoth tests the capturing of writings to stdout
// and stderr.
func TestBoth(t *testing.T) {
	hello := "Hello, World!"
	ouch := "ouch"
	cout, cerr := capture.Both(func() {
		fmt.Fprint(os.Stdout, hello)
		fmt.Fprint(os.Stderr, ouch)
	})
	Assert(t, Equal(cout.String(), hello), "captured stdout is not equal to expected one")
	Assert(t, Length(cout, len(hello)), "captured stdout has wrong length")
	Assert(t, Equal(cerr.String(), ouch), "captured stderr is not equal to expected one")
	Assert(t, Length(cerr, len(ouch)), "captured stderr has wrong length")
}

// TestBytes tests the retrieving of captures as bytes.
func TestBytes(t *testing.T) {
	foo := "foo"
	cout, cerr := capture.Both(func() {
		fmt.Fprint(os.Stdout, foo)
		fmt.Fprint(os.Stderr, foo)
	})
	Assert(t, Equal(cout.String(), foo), "captured stdout is not equal to expected one")
	Assert(t, Equal(cerr.String(), foo), "captured stderr is not equal to expected one")
}

// TestRestore tests the restoring of os.Stdout
// and os.Stderr after capturing.
func TestRestore(t *testing.T) {
	foo := "foo"
	oldOut := os.Stdout
	oldErr := os.Stderr
	cout, cerr := capture.Both(func() {
		fmt.Fprint(os.Stdout, foo)
		fmt.Fprint(os.Stderr, foo)
	})
	Assert(t, Equal(cout.String(), foo), "captured stdout is not equal to expected one")
	Assert(t, Length(cout, len(foo)), "captured stdout has wrong length")
	Assert(t, Equal(cerr.String(), foo), "captured stderr is not equal to expected one")
	Assert(t, Length(cerr, len(foo)), "captured stderr has wrong length")
	Assert(t, Equal(os.Stdout, oldOut), "os.Stdout not restored")
	Assert(t, Equal(os.Stderr, oldErr), "os.Stderr not restored")
}

// EOF
