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
)

//------------------------------
// PRINT HELPERS
//------------------------------

// valueDescription returns a description of a value as string.
func valueDescription(value any) string {
	rvalue := reflect.ValueOf(value)
	kind := rvalue.Kind()
	switch kind {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return kind.String() + " of " + rvalue.Type().Elem().String()
	case reflect.Func:
		return kind.String() + " " + rvalue.Type().Name() + "()"
	case reflect.Interface, reflect.Struct:
		return kind.String() + " " + rvalue.Type().Name()
	case reflect.Ptr:
		return kind.String() + " to " + rvalue.Type().Elem().String()
	default:
		return kind.String()
	}
}

// typedValue returns a value including its type.
func typedValue(value any) string {
	kind := reflect.ValueOf(value).Kind()
	switch kind {
	case reflect.String:
		return fmt.Sprintf("%q (string)", value)
	default:
		return fmt.Sprintf("%v (%s)", value, kind.String())
	}
}

// EOF
