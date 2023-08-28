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
	"time"
)

//--------------------
// ELEMENT CONVERSION
//--------------------

// elementToValue converts an element into a typed Element. The second
// return value is true if the element could not be converted.
func elementToValue[V ValueConstraint](e Element, def V) (Element, bool) {
	deft := reflect.TypeOf(def)
	switch deft.Kind() {
	case reflect.String:
		se, ok := e.(string)
		if ok {
			return se, false
		}
		return def, true
	case reflect.Int:
		ie, ok := e.(int)
		if ok {
			return ie, false
		}
		fe, ok := e.(float64)
		if ok {
			return int(fe), false
		}
		return def, true
	case reflect.Float64:
		fe, ok := e.(float64)
		if ok {
			return fe, false
		}
		ie, ok := e.(int)
		if ok {
			return float64(ie), false
		}
		return def, true
	case reflect.Bool:
		be, ok := e.(bool)
		if ok {
			return be, false
		}
		return def, true
	case reflect.Struct:
		// Time is a JSON string.
		t, ok := elementToTime(e)
		if !ok {
			return def, true
		}
		return t, false
	case reflect.Int64:
		// Duration is a JSON string.
		d, ok := elementToDuration(e)
		if !ok {
			return def, true
		}
		return d, false
	}
	return def, true
}

// valueToElement converts a value into an element if it is allowed. It's the
// reverse of elementToValue.
func valueToElement[V ValueConstraint](cv Element, nv V) (Element, bool) {
	nvt := reflect.TypeOf(nv)
	switch nvt.Kind() {
	case reflect.String:
		_, ok := cv.(string)
		if ok {
			return nv, true
		}
		return cv, false
	case reflect.Int:
		_, ok := cv.(int)
		if ok {
			return nv, true
		}
		_, ok = cv.(float64)
		if ok {
			return nv, true
		}
		return cv, false
	case reflect.Float64:
		_, ok := cv.(float64)
		if ok {
			return nv, true
		}
		_, ok = cv.(int)
		if ok {
			return nv, true
		}
		return cv, false
	case reflect.Bool:
		_, ok := cv.(bool)
		if ok {
			return nv, true
		}
		return cv, false
	case reflect.Struct:
		// Time is a JSON string.
		_, ok := elementToTime(cv)
		if !ok {
			return cv, false
		}
		return timeToElement(nv)
	case reflect.Int64:
		// Duration is a JSON string.
		_, ok := elementToDuration(cv)
		if !ok {
			return cv, false
		}
		return durationToElement(nv)
	}
	return cv, false
}

// timeToElement converts a time.Time into an element.
func timeToElement(a any) (Element, bool) {
	switch t := a.(type) {
	case time.Time:
		return t.Format(time.RFC3339Nano), true
	default:
		return nil, false
	}
}

// elementToTime converts an element into a time.Time.
func elementToTime(e Element) (time.Time, bool) {
	s, ok := e.(string)
	if !ok {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// durationToElement converts a time.Duration into an element.
func durationToElement(a any) (Element, bool) {
	switch d := a.(type) {
	case time.Duration:
		return d.String(), true
	default:
		return nil, false
	}
}

// elementToDuration converts an element into a time.Duration.
func elementToDuration(e Element) (time.Duration, bool) {
	s, ok := e.(string)
	if !ok {
		return time.Duration(0), false
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return time.Duration(0), false
	}
	return d, true
}

// EOF
