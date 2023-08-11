// Tideland Go Stew - QA One-Liners
//
// Copyright (C) 2012-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package qaone // import "tideland.dev/go/stew/qaone"

//------------------------------
// IMPORTS
//------------------------------

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

//------------------------------
// INSPECTORS
//------------------------------

// inspectNil checks if obtained is nil in a safe way.
func inspectNil(obtained any) (bool, error) {
	if obtained == nil {
		// Standard test.
		return true, nil
	}
	// Some types have to be tested via reflection.
	value := reflect.ValueOf(obtained)
	kind := value.Kind()
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil(), nil
	}
	return false, fmt.Errorf("obtained %s cannot be nil", valueDescription(obtained))
}

// inspectZero checks if obtained is the zero value of its type.
func inspectZero(obtained any) (bool, error) {
	value := reflect.ValueOf(obtained)
	kind := value.Kind()
	switch kind {
	case reflect.Func, reflect.Interface, reflect.Ptr:
		return value.IsNil(), nil
	case reflect.Chan, reflect.Map, reflect.Slice:
		if value.IsNil() {
			return true, nil
		}
		l := value.Len()
		return l == 0, nil
	case reflect.Array:
		l := value.Len()
		return l == 0, nil
	case reflect.Bool:
		return !value.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0, nil
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0.0, nil
	case reflect.Complex64, reflect.Complex128:
		return value.Complex() == 0.0, nil
	case reflect.String:
		return value.String() == "", nil
	}
	return false, fmt.Errorf("obtained %s cannot be zero", valueDescription(obtained))
}

// errable describes a type able to return an error state
// with the method Err().
type errable interface {
	Err() error
}

// inspectOK checks if obtained is ok in a safe way.
func inspectOK(obtained any) (bool, error) {
	var ok bool
	var err error
	switch value := obtained.(type) {
	case bool:
		ok = value
	case int:
		ok = value == 0
	case string:
		ok = value == ""
	case error:
		ok = value == nil
	case func() bool:
		ok = value()
	case func() error:
		ok = value() == nil
	default:
		var oerr error
		oerr, err = inspectError(obtained)
		if err != nil {
			return false, err
		}
		ok = oerr == nil
	}
	return ok, nil
}

// inspectError checks if obtained is an error in a safe way.
func inspectError(obtained any) (error, error) {
	switch v := obtained.(type) {
	case error:
		if v == nil {
			return nil, nil
		}
		return v, nil
	case func() error:
		if v == nil {
			return nil, nil
		}
		return v(), nil
	case errable:
		if v == nil {
			return nil, nil
		}
		return v.Err(), nil
	default:
		return nil, fmt.Errorf("no error or errable: %v", valueDescription(obtained))
	}
}

// lenable is used to check if a type has a Len() method.
type lenable interface {
	Len() int
}

// inspctLength checks the length of the obtained string, array, slice, map or channel.
func inspctLength(obtained any) (int, error) {
	// Check using the lenable interface.
	if ol, ok := obtained.(lenable); ok {
		return ol.Len(), nil
	}
	// Check for sting due to UTF-8 rune handling.
	if s, ok := obtained.(string); ok {
		l := utf8.RuneCountInString(s)
		return l, nil
	}
	// Check the standard types.
	ov := reflect.ValueOf(obtained)
	ok := ov.Kind()
	switch ok {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		l := ov.Len()
		return l, nil
	default:
		descr := valueDescription(obtained)
		return 0, fmt.Errorf("obtained %s is no array, chan, map, slice, string and does not understand Len()", descr)
	}
}

// EOF
