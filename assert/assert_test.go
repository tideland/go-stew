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
	"os"
	"path/filepath"
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

	Assert(stb, Nil(""), "empty string (fail)")
	Assert(stb, Nil(&struct{}{}), "referenced struct (fail)")
	Assert(stb, Nil([]int{1, 2, 3}), "filled slice (fail)")

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

	Assert(stb, NotNil(nil), "nil (fail)")
	Assert(stb, NotNil(is), "empty slice (fail)")
	Assert(stb, NotNil(ch), "channel (fail)")
	Assert(stb, NotNil(f), "function (fail)")

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

	Assert(stb, Zero(true), "true (fail)")
	Assert(stb, Zero(1), "1 (fail)")
	Assert(stb, Zero(1.0), "1.0 (fail)")
	Assert(stb, Zero("foo"), "filled string (fail)")
	Assert(stb, Zero([]int{1, 2, 3}), "filled slice (fail)")

	Assert(t, Equal(stb.Calls(), 13), "should be thirteen calls")
	Assert(t, Equal(stb.Len(), 5), "should be five fails")
}

// TestOK tests the OK assertion.
func TestOK(t *testing.T) {
	stb := newSubTB()

	Assert(stb, OK(""), "empty string")
	Assert(stb, OK(func() bool { return true }), "function returning true")
	Assert(stb, OK(func() error { return nil }), "function returning nil error")

	Assert(stb, OK(1), "1 (fail)")
	Assert(stb, OK(interfacor("ouch")), "error (fail)")

	Assert(t, Equal(stb.Calls(), 5), "should be five calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestNotOK tests the NotOK assertion.
func TestNotOK(t *testing.T) {
	stb := newSubTB()

	Assert(stb, NotOK(1), "1")
	Assert(stb, NotOK(fmt.Errorf("error")), "error")
	Assert(stb, NotOK(func() bool { return false }), "function returning false")

	Assert(stb, NotOK(nil), "nil (fail)")
	Assert(stb, NotOK(""), "empty string (fail)")
	Assert(stb, NotOK(func() error { return nil }), "function returning nil error (fail)")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 3), "should be three fails")
}

// TestAnyError tests the AnyError assertion.
func TestAnyError(t *testing.T) {
	stb := newSubTB()

	Assert(stb, AnyError(interfacor("ouch")), "error")
	Assert(stb, AnyError(func() error { return fmt.Errorf("error") }), "function returning error")

	Assert(stb, AnyError(nil), "nil (fail)")
	Assert(stb, AnyError(func() error { return nil }), "function returning nil error (fail)")

	Assert(t, Equal(stb.Calls(), 4), "should be four calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestErrorContains tests the ErrorContains assertion.
func TestErrorContains(t *testing.T) {
	stb := newSubTB()

	Assert(stb, ErrorContains(interfacor("ouch"), "ouch"), "error")
	Assert(stb, ErrorContains(func() error { return fmt.Errorf("ouch") }, "ouch"), "function returning error")

	Assert(stb, ErrorContains(nil, "ouch"), "nil (fail)")
	Assert(stb, ErrorContains(interfacor("nope"), "ouch"), "error (fail)")
	Assert(stb, ErrorContains(func() error { return nil }, "ouch"), "function returning nil error (fail)")

	Assert(t, Equal(stb.Calls(), 5), "should be four calls")
	Assert(t, Equal(stb.Len(), 3), "should be two fails")
}

// TestErrorMatches tests the ErrorMatches assertion.
func TestErrorMatches(t *testing.T) {
	stb := newSubTB()

	Assert(stb, ErrorMatches(interfacor("ERR ouch 123"), ".*ouch.*"), "error")
	Assert(stb, ErrorMatches(func() error { return fmt.Errorf("ERR ouch 123") }, ".*ouch.*"), "function returning error")

	Assert(stb, ErrorMatches(nil, "ouch"), "nil (fail)")
	Assert(stb, ErrorMatches(interfacor("nope"), "ouch"), "error (fail)")
	Assert(stb, ErrorMatches(func() error { return nil }, "ouch"), "function returning nil error (fail)")

	Assert(t, Equal(stb.Calls(), 5), "should be four calls")
	Assert(t, Equal(stb.Len(), 3), "should be two fails")
}

// TestNoError tests the NoError assertion.
func TestNoError(t *testing.T) {
	stb := newSubTB()

	Assert(stb, NoError(func() error { return nil }), "function returning nil error")

	err := interfacor("ouch")

	Assert(stb, NoError(err), "error (fail)")
	Assert(stb, NoError(fmt.Errorf("error")), "function returning error (fail)")

	Assert(t, Equal(stb.Calls(), 3), "should be three calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestPanics tests the Panics assertion.
func TestPanics(t *testing.T) {
	stb := newSubTB()

	Assert(stb, Panics(func() { panic("ouch") }), "function panics")

	Assert(stb, Panics(func() {}), "function does not panic (fail)")

	Assert(t, Equal(stb.Calls(), 2), "should be two calls")
	Assert(t, Equal(stb.Len(), 1), "should be one fail")
}

// TestPanicsNot tests the PanicsNot assertion.
func TestPanicsNot(t *testing.T) {
	stb := newSubTB()

	Assert(stb, PanicsNot(func() {}), "function does not panic")

	Assert(stb, PanicsNot(func() { panic("ouch") }), "function panics (fail)")

	Assert(t, Equal(stb.Calls(), 2), "should be two calls")
	Assert(t, Equal(stb.Len(), 1), "should be one fail")
}

// TestPanicsWith tests the PanicsWith assertion.
func TestPanicsWith(t *testing.T) {
	stb := newSubTB()

	Assert(stb, PanicsWith(func() { panic("ouch") }, "ouch"), "function panics with expected value")

	Assert(stb, PanicsWith(func() { panic("ouch") }, "nope"), "function panics with unexpected value (fail)")
	Assert(stb, PanicsWith(func() {}, "ouch"), "function does not panic (fail)")

	Assert(t, Equal(stb.Calls(), 3), "should be four calls")
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

	Assert(stb, True(false), "false (fail)")
	Assert(stb, True(1 == 2), "negative int comparison (fail)")
	Assert(stb, True("foo" == "bar"), "negative string comparison (fail)")

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

	Assert(stb, False(true), "true (fail)")
	Assert(stb, False(one == 1), "positive int comparison (fail)")
	Assert(stb, False(foo == "foo"), "positive string comparison (fail)")

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

	Assert(stb, Equal(1, 2), "different ints (fail)")
	Assert(stb, Equal("foo", "bar"), "different strings (fail)")

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

	Assert(stb, Different(1, 1), "equal ints (fail)")
	Assert(stb, Different(one, 1), "equal ints with variable (fail)")
	Assert(stb, Different("foo", "foo"), "equal strings (fail)")
	Assert(stb, Different(foo, "foo"), "equal strings with variable (fail)")

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

	Assert(stb, Length(foo, 3), "wrong length for slice (fail))")
	Assert(stb, Length(bar, 1), "wrong length for string (fail)")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestEmpty tests the Empty assertion for different types.
func TestEmpty(t *testing.T) {
	stb := newSubTB()
	foo := []string{}
	bar := ""
	l := interfacor("")
	ch := make(chan string, 3)

	Assert(stb, Empty(foo), "empty slice")
	Assert(stb, Empty(bar), "empty string")
	Assert(stb, Empty(l), "empty lenable")
	Assert(stb, Empty(ch), "empty channel")

	ch <- "foo"
	ch <- "bar"

	Assert(stb, Empty([]string{"foo"}), "non-empty slice (fail)")
	Assert(stb, Empty("foo"), "non-empty string (fail)")
	Assert(stb, Empty(interfacor("foo")), "non-empty lenable (fail)")
	Assert(stb, Empty(ch), "non-empty channel (fail)")

	Assert(t, Equal(stb.Calls(), 8), "should be eight calls")
	Assert(t, Equal(stb.Len(), 4), "should be four fails")
}

// TestNotEmpty tests the NotEmpty assertion for different types.
func TestNotEmpty(t *testing.T) {
	stb := newSubTB()
	foo := []string{}
	bar := ""
	l := interfacor("")
	ch := make(chan string, 3)

	ch <- "foo"
	ch <- "bar"

	Assert(stb, NotEmpty([]string{"foo"}), "non-empty slice")
	Assert(stb, NotEmpty("foo"), "non-empty string")
	Assert(stb, NotEmpty(interfacor("foo")), "non-empty lenable")
	Assert(stb, NotEmpty(ch), "non-empty channel")

	ch = make(chan string, 3)

	Assert(stb, NotEmpty(foo), "empty slice	(fail)")
	Assert(stb, NotEmpty(bar), "empty string (fail)")
	Assert(stb, NotEmpty(l), "empty lenable (fail)")
	Assert(stb, NotEmpty(ch), "empty channel (fail)")

	Assert(t, Equal(stb.Calls(), 8), "should be eight calls")
	Assert(t, Equal(stb.Len(), 4), "should be four fails")
}

// TestContains tests the Contains assertion for different types.
func TestContains(t *testing.T) {
	stb := newSubTB()
	fbb := []string{"foo", "bar", "baz"}
	ott := []int{1, 2, 3}

	Assert(stb, Contains(fbb, "foo"), "slice contains foo")
	Assert(stb, Contains(fbb, "bar"), "slice contains bar")
	Assert(stb, Contains(fbb, "baz"), "slice contains baz")
	Assert(stb, Contains(ott, 1), "slice contains 1")
	Assert(stb, Contains(ott, 2), "slice contains 2")
	Assert(stb, Contains(ott, 3), "slice contains 3")

	Assert(stb, Contains(fbb, "qux"), "slice does not contain string (fail)")
	Assert(stb, Contains(ott, 4), "slice does not contain 4	(fail)")

	Assert(t, Equal(stb.Calls(), 8), "should be eight calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestContainsNot tests the ContainsNot assertion for different types.
func TestContainsNot(t *testing.T) {
	stb := newSubTB()
	fbb := []string{"foo", "bar", "baz"}
	ott := []int{1, 2, 3}

	Assert(stb, ContainsNot(fbb, "qux"), "slice does not contain string")
	Assert(stb, ContainsNot(ott, 4), "slice does not contain 4")

	Assert(stb, ContainsNot(fbb, "foo"), "slice contains foo (fail)")
	Assert(stb, ContainsNot(fbb, "bar"), "slice contains bar (fail)")
	Assert(stb, ContainsNot(fbb, "baz"), "slice contains baz (fail)")
	Assert(stb, ContainsNot(ott, 1), "slice contains 1 (fail)")
	Assert(stb, ContainsNot(ott, 2), "slice contains 2 (fail)")
	Assert(stb, ContainsNot(ott, 3), "slice contains 3 (fail)")

	Assert(t, Equal(stb.Calls(), 8), "should be eight calls")
	Assert(t, Equal(stb.Len(), 6), "should be six fails")
}

// TestAbout tests the About assertion for different types.
func TestAbout(t *testing.T) {
	stb := newSubTB()

	Assert(stb, About(1.0, 1.0, 0.0), "1.0 ≈ 1.0")
	Assert(stb, About(1.0, 1.0, 0.1), "1.0 ≈ 1.0 ± 0.1")
	Assert(stb, About(1.0, 1.1, 0.1), "1.0 ≈ 1.1 ± 0.1")
	Assert(stb, About(100, 90, 10), "100 ≈ 90 ± 10")

	Assert(stb, About(1.0, 1.2, 0.1), "1.0 ≈ 1.2 ± 0.1 (fail)")
	Assert(stb, About(100, 80, 10), "100 ≈ 80 ± 10 (fail)")

	Assert(t, Equal(stb.Calls(), 6), "should be six calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestRange tests the Range assertion for different types.
func TestRange(t *testing.T) {
	stb := newSubTB()

	Assert(stb, Range(10, 10, 10), "10 <= 10 <= 10")
	Assert(stb, Range(10, 9, 11), "9 <= 10 <= 11")
	Assert(stb, Range(1.0, 0.9, 1.1), "ß.9 <= 1.0 <= 1.1")

	Assert(stb, Range(10, 15, 20), "15 <= 10 <= 20 (fail)")
	Assert(stb, Range(1.0, 5.0, 10.0), "5.0 <= 1.0 <= 10.0 (fail)")

	Assert(t, Equal(stb.Calls(), 5), "should be five calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestOneCase tests the OneCase assertion.
func TestOneCase(t *testing.T) {
	stb := newSubTB()

	Assert(stb, OneCase("lower"), "lower")
	Assert(stb, OneCase("UPPER"), "UPPER")

	Assert(stb, OneCase("Mixed"), "Mixed (fail)")

	Assert(t, Equal(stb.Calls(), 3), "should be three calls")
	Assert(t, Equal(stb.Len(), 1), "should be one fail")
}

// TestMatches tests the Matches assertion.
func TestMatches(t *testing.T) {
	stb := newSubTB()

	Assert(stb, Matches("foo", "foo"), "foo matches foo")
	Assert(stb, Matches("foo", "foo|bar"), "foo matches foo|bar")
	Assert(stb, Matches("bar", "foo|bar"), "bar matches foo|bar")

	Assert(stb, Matches("foo", "bar"), "foo matches bar (fail)")
	Assert(stb, Matches("bar", "foo"), "bar matches foo (fail)")

	Assert(t, Equal(stb.Calls(), 5), "should be five calls")
	Assert(t, Equal(stb.Len(), 2), "should be two fails")
}

// TestPathExists tests the PathExists assertion.
func TestPathExists(t *testing.T) {
	stb := newSubTB()

	dir := filepath.Join(os.TempDir(), "assert-path-exists")
	err := os.Mkdir(dir, 0700)
	Assert(stb, NoError(err), "create temporary test dir")
	defer func() {
		os.RemoveAll(dir)
	}()

	Assert(stb, PathExists(dir), "temporary test dir exists")

	Assert(stb, PathExists("/this/path/will/hopefully/not/exist"), "path not exists (fail)")

	Assert(t, Equal(stb.Calls(), 3), "should be four calls")
	Assert(t, Equal(stb.Len(), 1), "should be two fails")

	stb.LogFails(t)
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
