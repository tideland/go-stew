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

// TestNil tests the Nil assertion.
func TestNil(t *testing.T) {
	stb := newSubTB()

	var is []int
	var ch chan int
	var f func() bool

	Assert(stb, Nil(nil), "nil")
	Assert(stb, Nil(is), "empty slice")
	Assert(stb, Nil(ch), "channel")
	Assert(stb, Nil(f), "function")

	Assert(stb, Nil(""), "empty string")
	Assert(stb, Nil(&struct{}{}), "referenced struct")
	Assert(stb, Nil([]int{1, 2, 3}), "filled slice")

	Assert(t, Equal(stb.Calls(), 7), "should be seven calls")
	Assert(t, Equal(stb.Len(), 3), "should be three fails")
}

// TestNotNil tests the NotNil assertion.
func TestNotNil(t *testing.T) {
	stb := newSubTB()

	var is []int
	var ch chan int
	var f func() bool

	Assert(stb, NotNil(&struct{}{}), "referenced struct")
	Assert(stb, NotNil([]int{1, 2, 3}), "filled slice")

	Assert(stb, NotNil(nil), "nil")
	Assert(stb, NotNil(is), "empty slice")
	Assert(stb, NotNil(ch), "channel")
	Assert(stb, NotNil(f), "function")

	Assert(t, Equal(stb.Calls(), 6), "should be seven calls")
	Assert(t, Equal(stb.Len(), 4), "should be five fails")
}

// TestZero tests the Zero assertion.
func TestZero(t *testing.T) {
	stb := newSubTB()

	var is []int
	var ch chan int
	var f func() bool

	Assert(stb, Zero(false), "false")
	Assert(stb, Zero(0), "0")
	Assert(stb, Zero(0.0), "0.0")
	Assert(stb, Zero(""), "empty string")
	Assert(stb, Zero(is), "empty slice")
	Assert(stb, Zero(ch), "channel")
	Assert(stb, Zero([]int{}), "empty slice")
	Assert(stb, Zero(f), "function")

	Assert(stb, Zero(true), "true")
	Assert(stb, Zero(1), "1")
	Assert(stb, Zero(1.0), "1.0")
	Assert(stb, Zero("foo"), "filled string")
	Assert(stb, Zero([]int{1, 2, 3}), "filled slice")

	Assert(t, Equal(stb.Calls(), 13), "should be thirteen calls")
	Assert(t, Equal(stb.Len(), 5), "should be five fails")
}

// TestOK tests the OK assertion.
func TestOK(t *testing.T) {
	stb := newSubTB()

	Assert(stb, OK(""), "empty string")
	Assert(stb, OK(func() bool { return true }), "function returning true")
	Assert(stb, OK(func() error { return nil }), "function returning nil error")

	Assert(stb, OK(1), "1")
	Assert(stb, OK(interfacor("ouch")), "error")

	Assert(t, Equal(stb.Calls(), 5), "should be five calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestNotOK tests the NotOK assertion.
func TestNotOK(t *testing.T) {
	stb := newSubTB()

	Assert(stb, NotOK(1), "1")
	Assert(stb, NotOK(fmt.Errorf("error")), "error")
	Assert(stb, NotOK(func() bool { return false }), "function returning false")

	Assert(stb, NotOK(nil), "nil")
	Assert(stb, NotOK(""), "empty string")
	Assert(stb, NotOK(func() error { return nil }), "function returning nil error")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 3), "should be three fails")
}

// TestAnyError tests the AnyError assertion.
func TestAnyError(t *testing.T) {
	stb := newSubTB()

	Assert(stb, AnyError(interfacor("ouch")), "error")
	Assert(stb, AnyError(func() error { return fmt.Errorf("error") }), "function returning error")

	Assert(stb, AnyError(nil), "nil")
	Assert(stb, AnyError(func() error { return nil }), "function returning nil error")

	Assert(t, Equal(stb.Calls(), 4), "should be four calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestNoError tests the NoError assertion.
func TestNoError(t *testing.T) {
	stb := newSubTB()

	Assert(stb, NoError(func() error { return nil }), "function returning nil error")

	err := interfacor("ouch")

	Assert(stb, NoError(err), "error")
	Assert(stb, NoError(fmt.Errorf("error")), "function returning error")

	Assert(t, Equal(stb.Calls(), 3), "should be three calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestTrue tests the True assertion.
func TestTrue(t *testing.T) {
	stb := newSubTB()
	one := 1
	foo := "foo"

	Assert(stb, True(true), "true")
	Assert(stb, True(one == 1), "positive int comparison")
	Assert(stb, True(foo == "foo"), "positive string comparison")

	Assert(stb, True(false), "false")
	Assert(stb, True(1 == 2), "negative int comparison")
	Assert(stb, True("foo" == "bar"), "negative string comparison")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 3), "should be three fails")
}

// TestFalse tests the False assertion.
func TestFalse(t *testing.T) {
	stb := newSubTB()
	one := 1
	foo := "foo"

	Assert(stb, False(false), "false")
	Assert(stb, False(1 == 2), "negative int comparison")
	Assert(stb, False("foo" == "bar"), "negative string comparison")

	Assert(stb, False(true), "true")
	Assert(stb, False(one == 1), "positive int comparison")
	Assert(stb, False(foo == "foo"), "positive string comparison")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 3), "should be three fails")
}

// TestEqual tests the Equal assertion.
func TestEqual(t *testing.T) {
	stb := newSubTB()
	one := 1
	foo := "foo"

	Assert(stb, Equal(1, 1), "equal ints")
	Assert(stb, Equal(one, 1), "equal ints with variable")
	Assert(stb, Equal("foo", "foo"), "equal strings")
	Assert(stb, Equal(foo, "foo"), "equal strings with variable")

	Assert(stb, Equal(1, 2), "different ints")
	Assert(stb, Equal("foo", "bar"), "different strings")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestDifferent tests the Different assertion.
func TestDifferent(t *testing.T) {
	stb := newSubTB()
	one := 1
	foo := "foo"

	Assert(stb, Different(1, 2), "different ints")
	Assert(stb, Different("foo", "bar"), "different strings")

	Assert(stb, Different(1, 1), "equal ints")
	Assert(stb, Different(one, 1), "equal ints with variable")
	Assert(stb, Different("foo", "foo"), "equal strings")
	Assert(stb, Different(foo, "foo"), "equal strings with variable")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 4), "should be four fails")
}

// TestLength tests the Length assertion for different types.
func TestLength(t *testing.T) {
	stb := newSubTB()
	foo := []string{"foo", "bar"}
	bar := "foobarbaz"
	l := interfacor("foo")
	ch := make(chan string, 3)

	ch <- "foo"
	ch <- "bar"

	Assert(stb, Length(foo, 2), "slice of two strings")
	Assert(stb, Length(bar, 9), "string of nine characters")
	Assert(stb, Length(l, 3), "lenable with value 3")
	Assert(stb, Length(ch, 2), "channel with two values")

	Assert(stb, Length(foo, 3), "wrong length for slice")
	Assert(stb, Length(bar, 1), "wrong length for string")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

//------------------------------
// TEST HELPER
//------------------------------

// subtb is a subset of testing.TB.
type subtb struct {
	mux   sync.Mutex
	calls int
	fails map[string]string
}

func newSubTB() *subtb {
	return &subtb{
		fails: make(map[string]string),
	}
}

func (stb *subtb) Helper() {
	stb.mux.Lock()
	defer stb.mux.Unlock()
	stb.calls++
}

func (stb *subtb) Errorf(format string, args ...any) {
	stb.mux.Lock()
	defer stb.mux.Unlock()
	key := fmt.Sprintf("E:%d", stb.calls)
	stb.fails[key] = fmt.Sprintf(format, args...)
}

func (stb *subtb) Fatalf(format string, args ...any) {
	stb.mux.Lock()
	defer stb.mux.Unlock()
	key := fmt.Sprintf("F:%d", stb.calls)
	stb.fails[key] = fmt.Sprintf(format, args...)
}

func (stb *subtb) Calls() int {
	stb.mux.Lock()
	defer stb.mux.Unlock()

	return stb.calls
}

func (stb *subtb) Len() int {
	stb.mux.Lock()
	defer stb.mux.Unlock()

	return len(stb.fails)
}

func (stb *subtb) Fails() string {
	stb.mux.Lock()
	defer stb.mux.Unlock()

	s := ""
	for k, v := range stb.fails {
		s += fmt.Sprintf("[%s] %s\n", k, v)
	}

	return s
}

func (stb *subtb) LogFails(t *testing.T) {
	t.Logf("Fails:\n%v", stb.Fails())
}

// interfacor is a type that implements different interface methods needed for
// testing.
type interfacor string

func (i interfacor) Len() int {
	return len(i)
}

func (i interfacor) Err() error {
	return fmt.Errorf(string(i))
}

func (i interfacor) Error() string {
	return string(i)
}

func (i interfacor) String() string {
	return string(i)
}

// EOF
