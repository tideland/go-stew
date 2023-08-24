// Tideland Go Stew - Generic JSON
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package genj // import "tideland.dev/go/stew/genj"

//--------------------
// IMPORTS
//--------------------

import (
	"reflect"
	"strconv"
	"time"
)

//--------------------
// ELEMENT CONVERSION
//--------------------

// elementToValue converts an element into a value returned as any.
func elementToValue[V ValueConstraint](e Element, def V) (any, bool) {
	deft := reflect.TypeOf(def)
	switch deft.Kind() {
	case reflect.String:
		return elementToString(e, def)
	case reflect.Int:
		return elementToInt(e, def)
	case reflect.Float64:
		return elementToFloat(e, def)
	case reflect.Bool:
		return elementToBool(e, def)
	case reflect.Struct:
		return elementToTime(e, def)
	case reflect.Int64:
		return elementToDuration(e, def)
	}
	return def, true
}

// elementToString converts an element into a string.
func elementToString[V ValueConstraint](e Element, def V) (any, bool) {
	switch et := e.(type) {
	case string:
		return et, false
	case int:
		return strconv.Itoa(et), false
	case float64:
		return strconv.FormatFloat(et, 'f', -1, 64), false
	case bool:
		return strconv.FormatBool(et), false
	case time.Time:
		return et.Format(time.RFC3339), false
	case time.Duration:
		return et.String(), false
	default:
		return def, true
	}
}

// elementToInt converts an element into an int.
func elementToInt[V ValueConstraint](e Element, def V) (any, bool) {
	switch et := e.(type) {
	case string:
		i, err := strconv.Atoi(et)
		if err != nil {
			return def, true
		}
		return i, false
	case int:
		return et, false
	case float64:
		return int(et), false
	case bool:
		if et {
			return 1, false
		}
		return 0, false
	case time.Time:
		return int(et.Unix()), false
	case time.Duration:
		return int(et.Seconds()), false
	default:
		return def, true
	}
}

// elementToFloat converts an element into a float.
func elementToFloat[V ValueConstraint](e Element, def V) (any, bool) {
	switch et := e.(type) {
	case string:
		f, err := strconv.ParseFloat(et, 64)
		if err != nil {
			return def, true
		}
		return f, false
	case int:
		return float64(et), false
	case float64:
		return et, false
	case bool:
		if et {
			return 1.0, false
		}
		return 0.0, false
	case time.Time:
		return float64(et.Unix()), false
	case time.Duration:
		return et.Seconds(), false
	default:
		return def, true
	}
}

// elementToBool converts an element into a bool.
func elementToBool[V ValueConstraint](e Element, def V) (any, bool) {
	switch et := e.(type) {
	case string:
		b, err := strconv.ParseBool(et)
		if err != nil {
			return def, true
		}
		return b, false
	case int:
		return et != 0, false
	case float64:
		return et != 0.0, false
	case bool:
		return et, false
	case time.Time:
		return !et.IsZero(), false
	case time.Duration:
		return et != 0, false
	default:
		return def, true
	}
}

// elementToTime converts an element into a time.
func elementToTime[V ValueConstraint](e Element, def V) (any, bool) {
	switch et := e.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, et)
		if err != nil {
			return def, true
		}
		return t, false
	case int:
		return time.Unix(int64(et), 0), false
	case float64:
		return time.Unix(int64(et), 0), false
	case bool:
		if et {
			return time.Now(), false
		}
		return time.Time{}, false
	case time.Time:
		return et, false
	case time.Duration:
		return time.Now().Add(et), false
	default:
		return def, true
	}
}

// elementToDuration converts an element into a duration.
func elementToDuration[V ValueConstraint](e Element, def V) (any, bool) {
	switch et := e.(type) {
	case string:
		d, err := time.ParseDuration(et)
		if err != nil {
			return def, true
		}
		return d, false
	case int:
		return time.Duration(et) * time.Second, false
	case float64:
		return time.Duration(et) * time.Second, false
	case bool:
		if et {
			return time.Second, false
		}
		return 0, false
	case time.Time:
		return time.Since(et), false
	case time.Duration:
		return et, false
	default:
		return def, true
	}
}

// EOF
