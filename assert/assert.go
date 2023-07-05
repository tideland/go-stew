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
		var ok bool
		var err error
		ok, err = inspectNil(v)
		info := ""
		if err != nil {
			return false, "", err
		}
		if !ok {
			info = typedValue(v) + " is not nil"
		}
		return ok, info, err
	}
}

// NotNil asserts that a value is not nil.
func NotNil(v any) Assertion {
	return func() (bool, string, error) {
		ok, info, err := Nil(v)()
		if err != nil {
			return false, "", err
		}
		return !ok, info, err
	}
}

// OK asserts that a value is true, nil, 0, "", or no error.
func OK(v any) Assertion {
	return func() (bool, string, error) {
		var err error
		ok := false
		info := ""
		switch obtained := v.(type) {
		case bool:
			ok = obtained
		case int:
			ok = obtained == 0
		case string:
			ok = obtained == ""
		case error:
			ok = obtained == nil
		case func() bool:
			ok = obtained()
		case func() error:
			ok = obtained() == nil
		default:
			var oerr error
			oerr, err = inspctError(obtained)
			if err != nil {
				return false, "", err
			}
			ok = oerr == nil
		}
		return ok, info, err
	}
}

// NotOK asserts that a value is false, not nil, not 0, not "", or an error.
func NotOK(v any) Assertion {
	return func() (bool, string, error) {
		ok, info, err := OK(v)()
		if err != nil {
			return false, "", err
		}
		return !ok, info, err
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
		ok, info, err := True(v)()
		if err != nil {
			return false, "", err
		}
		return !ok, info, err
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
		return va == vb, info, err
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

// EOF
