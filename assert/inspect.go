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

// errable describes a type able to return an error state
// with the method Err().
type errable interface {
	Err() error
}

// inspctError converts an any variable into an error.
func inspctError(obtained any) (error, error) {
	if obtained == nil {
		return nil, nil
	}
	err, ok := obtained.(error)
	if ok {
		return err, nil
	}
	able, ok := obtained.(errable)
	if ok {
		if able == nil {
			return nil, nil
		}
		return able.Err(), nil
	}
	// No error and not errable, so return error about it.
	return nil, fmt.Errorf("no error or errable: %T", valueDescription(obtained))
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
