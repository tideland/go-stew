// Tideland Go Stew - Assert
//
// Copyright (C) 2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package assert // import "tideland.dev/go/stew/assert"

//------------------------------
// IMPORTS
//------------------------------

import (
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"
)

//------------------------------
// ASSERTION TYPE AND FUNCTION
//------------------------------

// Assertion describes a function performing an assertion and returning false if
// it failed. Additionally it returns a message describing the assertion failure.
type Assertion func() (bool, string, error)

// SubTB is a sub interface of testung.TB.
type SubTB interface {
	Helper()

	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// Assert executes the given assertion and fails a test if it returns false.
func Assert(stb SubTB, assert Assertion, msg string) bool {
	stb.Helper()
	ok, info, err := assert()
	if err != nil {
		stb.Fatalf("error in assertion: %s (%s)", err.Error(), msg)
		return false
	}
	if !ok {
		stb.Errorf("assertion failed: %s (%s)", info, msg)
		return false
	}
	return true
}

//------------------------------
// ASSERTIONS
//------------------------------

// Nil asserts that a value is nil.
func Nil(v any) Assertion {
	return func() (bool, string, error) {
		isNil, err := inspectNil(v)
		if err != nil {
			return false, "", err
		}
		if isNil {
			return true, "", nil
		}
		return false, typedValue(v) + " is not nil", nil
	}
}

// NotNil asserts that a value is not nil.
func NotNil(v any) Assertion {
	return func() (bool, string, error) {
		isNil, err := inspectNil(v)
		if err != nil {
			return false, "", err
		}
		if isNil {
			return false, typedValue(v) + " is nil", nil
		}
		return true, "", nil
	}
}

// Zero asserts that a value has its zero value.
func Zero(v any) Assertion {
	return func() (bool, string, error) {
		var ok bool
		var err error
		ok, err = inspectZero(v)
		info := ""
		if err != nil {
			return false, "", err
		}
		if !ok {
			info = typedValue(v) + " is not zero"
		}
		return ok, info, err
	}
}

// OK asserts that a value is true, nil, 0, "", no error or a function returning true
// or a nil error.
func OK(v any) Assertion {
	return func() (bool, string, error) {
		ok, err := inspectOK(v)
		if err != nil {
			return false, "", err
		}
		return ok, "", err
	}
}

// NotOK asserts that a value is false, not nil, not 0, not "", an error or a function
// returning false or a non-nil error.
func NotOK(v any) Assertion {
	return func() (bool, string, error) {
		ok, err := inspectOK(v)
		if err != nil {
			return false, "", err
		}
		return !ok, "", err
	}
}

// AnyError asserts that a value is an error, a function returning an error or
// a type implementing an Err() error method and the error is not nil.
func AnyError(v any) Assertion {
	return func() (bool, string, error) {
		var ierr, err error
		ierr, err = inspectError(v)
		ok := ierr != nil
		info := ""
		if err != nil {
			return false, "", err
		}
		if !ok {
			info = typedValue(v) + " is or returns no error"
		}
		return ok, info, err
	}
}

// ErrorContains asserts that a value is an error like in AnyError and the error
// contains the given string.
func ErrorContains(v any, contains string) Assertion {
	return func() (bool, string, error) {
		var ierr, err error
		ierr, err = inspectError(v)
		ok := ierr != nil
		info := ""
		if err != nil {
			return false, "", err
		}
		if !ok {
			info = typedValue(v) + " is or returns no error"
		}
		if ok && !strings.Contains(ierr.Error(), contains) {
			info = typedValue(v) + " does not con	tain " + contains
			ok = false
		}
		return ok, info, err
	}
}

// ErrorMatches asserts that a value is an error like in AnyError and the error
// matches the given regular expression.
func ErrorMatches(v any, pattern string) Assertion {
	return func() (bool, string, error) {
		var ierr, err error
		ierr, err = inspectError(v)
		ok := ierr != nil
		info := ""
		if err != nil {
			return false, "", err
		}
		if !ok {
			info = typedValue(v) + " is or returns no error"
		}
		if ok && !regexp.MustCompile(pattern).MatchString(ierr.Error()) {
			info = typedValue(v) + " does not match " + pattern
			ok = false
		}
		return ok, info, err
	}
}

// NoError asserts that a value is an error like in AnyError and the error is nil.
func NoError(v any) Assertion {
	return func() (bool, string, error) {
		var ierr, err error
		ierr, err = inspectError(v)
		ok := ierr == nil
		info := ""
		if err != nil {
			return false, "", err
		}
		if !ok {
			info = typedValue(v) + " is or returns an error"
		}
		return ok, info, err
	}
}

// True asserts that a value is true.
func True(v any) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := false
		info := ""
		switch obtained := v.(type) {
		case bool:
			ok = obtained
		case func() bool:
			ok = obtained()
		default:
			info = typedValue(v) + " is not true"
		}
		return ok, info, err
	}
}

// False asserts that a value is false.
func False(v any) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := false
		info := ""
		switch obtained := v.(type) {
		case bool:
			ok = !obtained
		case func() bool:
			ok = !obtained()
		default:
			info = typedValue(v) + " is not false"
		}
		return ok, info, err
	}
}

// Equal asserts that two comparable values are equal.
func Equal[T comparable](va, vb T) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := va == vb
		info := ""
		if !ok {
			info = typedValue(va) + " is not equal to " + typedValue(vb)
		}
		return ok, info, err
	}
}

// Different asserts that two comparable values are different.
func Different[T comparable](va, vb T) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := va != vb
		info := ""
		if !ok {
			info = typedValue(va) + " is equal to " + typedValue(vb)
		}
		return ok, info, err
	}
}

// Length asserts that a value has a specific size.
func Length(v any, l int) Assertion {
	return func() (bool, string, error) {
		var err error
		obtained, err := inspctLength(v)
		ok := obtained == l
		info := ""
		if !ok {
			info = typedValue(v) + " has not the expected length"
		}
		return ok, info, err
	}
}

// Empty asserts that a value is empty.
func Empty(v any) Assertion {
	return func() (bool, string, error) {
		var err error
		obtained, err := inspctLength(v)
		ok := obtained == 0
		info := ""
		if !ok {
			info = typedValue(v) + " is not empty"
		}
		return ok, info, err
	}
}

// NotEmpty asserts that a value is not empty.
func NotEmpty(v any) Assertion {
	return func() (bool, string, error) {
		var err error
		obtained, err := inspctLength(v)
		ok := obtained != 0
		info := ""
		if !ok {
			info = typedValue(v) + " is empty"
		}
		return ok, info, err
	}
}

// Contains asserts that a slice contains a specific value.
func Contains[S ~[]T, T comparable](vs S, content T) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := false
		info := ""
		for _, v := range vs {
			if v == content {
				ok = true
				break
			}
		}
		if !ok {
			info = typedValue(vs) + " does not contain " + typedValue(content)
		}
		return ok, info, err
	}
}

// ContainsNot asserts that a slice does not contain a specific value.
func ContainsNot[S ~[]T, T comparable](vs S, content T) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := true
		info := ""
		for _, v := range vs {
			if v == content {
				ok = false
				break
			}
		}
		if !ok {
			info = typedValue(vs) + " contains " + typedValue(content)
		}
		return ok, info, err
	}
}

// About asserts that two numbers are about equal within a delta.
func About[T constraints.Integer | constraints.Float](va, vb, delta T) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := va >= vb-delta && va <= vb+delta
		info := ""
		if !ok {
			info = typedValue(va) + " is not about equal to " + typedValue(vb)
		}
		return ok, info, err
	}
}

// Range asserts that a number is in a range.
func Range[T constraints.Integer | constraints.Float](v, min, max T) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := v >= min && v <= max
		info := ""
		if !ok {
			info = typedValue(v) + " is not in range [" + typedValue(min) + ", " + typedValue(max) + "]"
		}
		return ok, info, err
	}
}

// OneCase asserts that a string is in one case.
func OneCase(v string) Assertion {
	return func() (bool, string, error) {
		upper := strings.ToUpper(v)
		lower := strings.ToLower(v)
		if v == upper || v == lower {
			return true, "", nil
		}
		return false, "not one case", nil
	}
}

// EOF
