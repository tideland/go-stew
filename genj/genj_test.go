// Tideland Go Stew - Generic JSON - Unit Tests
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package genj_test

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"testing"
	"time"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/genj"
)

//--------------------
// TESTS
//--------------------

// TestRead tests the reading and writing of a JSON document.
func TestRead(t *testing.T) {
	doc, err := genj.Read(bytes.NewReader(createJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")
}

// TestGet tests the getting of values from a JSON document.
func TestGet(t *testing.T) {
	doc, err := genj.Read(bytes.NewReader(createJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Valid access.
	s, err := genj.Get[string](doc, "string")
	Assert(t, NoError(err), "string must be accessible")
	Assert(t, Equal(s, "value"), "string must be correct")

	i, err := genj.Get[int](doc, "int")
	Assert(t, NoError(err), "int must be accessible")
	Assert(t, Equal(i, 42), "int must be correct")

	i, err = genj.Get[int](doc, "intstring")
	Assert(t, NoError(err), "int string must be accessible")
	Assert(t, Equal(i, 4711), "int string must be correct")

	f, err := genj.Get[float64](doc, "float")
	Assert(t, NoError(err), "float must be accessible")
	Assert(t, Equal(f, 3.1415), "float must be correct")

	f, err = genj.Get[float64](doc, "floatstring")
	Assert(t, NoError(err), "float string must be accessible")
	Assert(t, Equal(f, 2.7182), "float string must be correct")

	b, err := genj.Get[bool](doc, "bool")
	Assert(t, NoError(err), "bool must be accessible")
	Assert(t, Equal(b, true), "bool must be correct")

	b, err = genj.Get[bool](doc, "boolstring")
	Assert(t, NoError(err), "bool string must be accessible")
	Assert(t, Equal(b, true), "bool string must be correct")

	tm, err := genj.Get[time.Time](doc, "time")
	Assert(t, NoError(err), "time must be accessible")
	Assert(t, Equal(tm, time.Date(2019, 9, 1, 12, 0, 0, 0, time.UTC)), "time must be correct")

	d, err := genj.Get[time.Duration](doc, "duration")
	Assert(t, NoError(err), "duration must be accessible")
	Assert(t, Equal(d, 90*time.Minute), "duration must be correct")

	s, err = genj.Get[string](doc, "nested", "0", "d", "1")
	Assert(t, NoError(err), "string must be accessible")
	Assert(t, Equal(s, "bar"), "string must be correct")

	// Invalid access.
	a, err := genj.Get[int](doc, "array")
	Assert(t, ErrorContains(err, "element is not of type int"), "array must not be accessible")
	Assert(t, Equal(a, 0), "int must be default value")

	s, err = genj.Get[string](doc, "this", "does", "not", "exist")
	Assert(t, ErrorContains(err, "cannot get element"), "string must not be accessible")
	Assert(t, Equal(s, ""), "string must be default value")
}

// TestGetDefault tests the getting of values from a JSON document.
func TestGetDefault(t *testing.T) {
	doc, err := genj.Read(bytes.NewReader(createJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Valid access.
	s := genj.GetDefault[string](doc, "default", "string")
	Assert(t, Equal(s, "value"), "string must be correct")

	i := genj.GetDefault[int](doc, 4711, "int")
	Assert(t, Equal(i, 42), "int must be correct")

	i = genj.GetDefault[int](doc, 1174, "intstring")
	Assert(t, Equal(i, 4711), "int string must be correct")

	f := genj.GetDefault[float64](doc, 2.7182, "float")
	Assert(t, Equal(f, 3.1415), "float must be correct")

	f = genj.GetDefault[float64](doc, 3.1415, "floatstring")
	Assert(t, Equal(f, 2.7182), "float string must be correct")

	b := genj.GetDefault[bool](doc, false, "bool")
	Assert(t, Equal(b, true), "bool must be correct")

	b = genj.GetDefault[bool](doc, false, "boolstring")
	Assert(t, Equal(b, true), "bool string must be correct")

	tm := genj.GetDefault[time.Time](doc, time.Date(2000, 9, 1, 1, 0, 0, 0, time.UTC), "time")
	Assert(t, Equal(tm, time.Date(2019, 9, 1, 12, 0, 0, 0, time.UTC)), "time must be correct")

	d := genj.GetDefault[time.Duration](doc, 5*time.Second, "duration")
	Assert(t, Equal(d, 90*time.Minute), "duration ust be correct")

	s = genj.GetDefault[string](doc, "default", "nested", "0", "d", "1")
	Assert(t, Equal(s, "bar"), "string must be correct")

	// Access with returning of default..
	i = genj.GetDefault[int](doc, 4711, "array")
	Assert(t, Equal(i, 4711), "int must be default value")

	s = genj.GetDefault[string](doc, "default", "this", "does", "not", "exist")
	Assert(t, Equal(s, "default"), "string must be default value")
}

//--------------------
// TESTS
//--------------------

// createJSON creates a simple JSON test document as bytes.
func createJSON() []byte {
	return []byte(`{
		"string": "value",
		"int": 42,
		"intstring": "4711",
		"float": 3.1415,
		"floatstring": "2.7182",
		"bool": true,
		"boolstring": "true",
		"time": "2019-09-01T12:00:00Z",
		"duration": "1h30m",
		"array": [
			"one",
			"two",
			"three"
		],
		"object": {
			"one": 1,
			"two": 2,
			"three": 3
		},
		"nested": [
			{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": ["foo", "bar", "baz"]
			},
			{
				"a": 9,
				"b": 8,
				"c": 7,
				"d": ["baz", "bar", "foo"]
			}
		]
	}`)
}

// EOF
