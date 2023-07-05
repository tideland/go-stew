// Tideland Go Stew - Assert - Unit Tests
//
// Copyright (C) 2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package assert_test

//------------------------------
// IMPORTS
//------------------------------

import (
	"fmt"
	"sync"
	"testing"

	. "tideland.dev/go/stew/assert"
)

//------------------------------
// TESTS
//------------------------------

// TestOK tests the OK assertion.
func TestOK(t *testing.T) {
	stb := newSubTB()

	Assert(stb, OK(nil), "should not fail")
	Assert(stb, OK(""), "should not fail")
	Assert(stb, OK(1), "should fail")
	Assert(stb, OK(fmt.Errorf("error")), "should fail")
	Assert(stb, OK(func() bool { return true }), "should not fail")
	Assert(stb, OK(func() error { return nil }), "should not fail")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 2), "should be two errors")
}

// TestNotOK tests the NotOK assertion.
func TestNotOK(t *testing.T) {
	stb := newSubTB()

	Assert(stb, NotOK(nil), "should fail")
	Assert(stb, NotOK(""), "should fail")
	Assert(stb, NotOK(1), "should not fail")
	Assert(stb, NotOK(fmt.Errorf("error")), "should not fail")
	Assert(stb, NotOK(func() bool { return false }), "should not fail")
	Assert(stb, NotOK(func() error { return nil }), "should fail")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 3), "should be three errors")
}

func TestTrue(t *testing.T) {
	stb := newSubTB()
	one := 1
	foo := "foo"

	Assert(stb, True(true), "should not fail")
	Assert(stb, True(false), "should fail")
	Assert(stb, True(one == 1), "should not fail")
	Assert(stb, True(1 == 2), "should fail")
	Assert(stb, True(foo == "foo"), "should not fail")
	Assert(stb, True("foo" == "bar"), "should fail")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 3), "should be three errors")
}

// TestEqual tests the Equal assertion.
func TestEqual(t *testing.T) {
	stb := newSubTB()
	one := 1
	foo := "foo"

	Assert(stb, Equal(1, 1), "should be equal")
	Assert(stb, Equal(1, 2), "should not be equal")
	Assert(stb, Equal(one, 1), "should be equal")
	Assert(stb, Equal("foo", "foo"), "should be equal")
	Assert(stb, Equal("foo", "bar"), "should not be equal")
	Assert(stb, Equal(foo, "foo"), "should be equal")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 2), "should be two errors")
}

// TestLength tests the Length assertion for different types.
func TestLength(t *testing.T) {
	stb := newSubTB()
	foo := []string{"foo", "bar"}
	bar := "foobarbaz"
	l := lenable(3)
	ch := make(chan string, 3)

	ch <- "foo"
	ch <- "bar"

	Assert(stb, Length(foo, 2), "should have length 2")
	Assert(stb, Length(foo, 3), "should have length 3")
	Assert(stb, Length(bar, 9), "should have length 9")
	Assert(stb, Length(bar, 1), "should have length 1")
	Assert(stb, Length(l, 3), "should have length 3")
	Assert(stb, Length(ch, 2), "should have length 2")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 2), "should be one error")
}

//------------------------------
// HELPER
//------------------------------

// subtb is a subset of testing.TB.
type subtb struct {
	mux    sync.Mutex
	calls  int
	errors []string
}

func newSubTB() *subtb {
	return &subtb{}
}

func (stb *subtb) Helper() {
	stb.mux.Lock()
	defer stb.mux.Unlock()
	stb.calls++
}

func (stb *subtb) Errorf(format string, args ...any) {
	stb.mux.Lock()
	defer stb.mux.Unlock()
	stb.errors = append(stb.errors, fmt.Sprintf(format, args...))
}

func (stb *subtb) Fatalf(format string, args ...any) {
	stb.mux.Lock()
	defer stb.mux.Unlock()
	stb.errors = append(stb.errors, fmt.Sprintf(format, args...))
}

func (stb *subtb) Calls() int {
	stb.mux.Lock()
	defer stb.mux.Unlock()

	return stb.calls
}

func (stb *subtb) Len() int {
	stb.mux.Lock()
	defer stb.mux.Unlock()

	return len(stb.errors)
}

type lenable int

func (l lenable) Len() int {
	return int(l)
}

// EOF
