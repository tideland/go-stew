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

	f, err := genj.Get[float64](doc, "float")
	Assert(t, NoError(err), "float must be accessible")
	Assert(t, Equal(f, 3.1415), "float must be correct")

	b, err := genj.Get[bool](doc, "bool")
	Assert(t, NoError(err), "bool must be accessible")
	Assert(t, Equal(b, true), "bool must be correct")

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
	Assert(t, ErrorContains(err, "path points to object or array"), "array must not be accessible")
	Assert(t, Equal(a, 0), "int must be default value")

	s, err = genj.Get[string](doc, "this", "does", "not", "exist")
	Assert(t, ErrorContains(err, "cannot get element"), "path does not exist")
	Assert(t, Equal(s, ""), "string must be standard value")

	s, err = genj.Get[string](doc, "nested", "1", "d")
	Assert(t, ErrorContains(err, "path points to object or array"), "path is array")
	Assert(t, Equal(s, ""), "string must be standard value")
}

// TestGetShort tests the from a JSON document with only one element.
func TestGetShort(t *testing.T) {
	// Small, but standard.
	doc, err := genj.Read(bytes.NewReader([]byte(`{"foo": "bar"}`)))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	s, err := genj.Get[string](doc, "foo")
	Assert(t, NoError(err), "string must be accessible")
	Assert(t, Equal(s, "bar"), "string must be correct")

	// Now even smaller.
	doc, err = genj.Read(bytes.NewReader([]byte(`"bar"`)))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	s, err = genj.Get[string](doc)
	Assert(t, NoError(err), "string must be accessible")
	Assert(t, Equal(s, "bar"), "string must be correct")
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

	i = genj.GetDefault[int](doc, 1234, "string")
	Assert(t, Equal(i, 1234), "int string must be default value as string is not an int")

	f := genj.GetDefault[float64](doc, 2.7182, "float")
	Assert(t, Equal(f, 3.1415), "float must be correct")

	b := genj.GetDefault[bool](doc, false, "bool")
	Assert(t, Equal(b, true), "bool must be correct")

	tm := genj.GetDefault[time.Time](doc, time.Date(2000, 9, 1, 1, 0, 0, 0, time.UTC), "time")
	Assert(t, Equal(tm, time.Date(2019, 9, 1, 12, 0, 0, 0, time.UTC)), "time must be correct")

	d := genj.GetDefault[time.Duration](doc, 5*time.Second, "duration")
	Assert(t, Equal(d, 90*time.Minute), "duration ust be correct")

	s = genj.GetDefault[string](doc, "default", "nested", "0", "d", "1")
	Assert(t, Equal(s, "bar"), "string must be correct")

	s = genj.GetDefault[string](doc, "default", "nested", "0", "d", "3")
	Assert(t, Equal(s, "default"), "string must be default value")

	s = genj.GetDefault[string](doc, "default", "this", "does", "not", "exist")
	Assert(t, Equal(s, "default"), "string must be default value")

	// Access with returning of default..
	s = genj.GetDefault[string](doc, "don' care", "nested", "0", "d")
	Assert(t, Equal(s, ""), "int must be standard value")

	i = genj.GetDefault[int](doc, 4711, "array")
	Assert(t, Equal(i, 0), "int must be standard value")
}

// TestSet tests the setting of values in a JSON document.
func TestSet(t *testing.T) {
	doc, err := genj.Read(bytes.NewReader(createJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Unnested valid and invalid.
	err = genj.Set(doc, "new value", "string")
	Assert(t, NoError(err), "string must be set")

	s, err := genj.Get[string](doc, "string")
	Assert(t, NoError(err), "string must be accessible")
	Assert(t, Equal(s, "new value"), "string must be correct")

	err = genj.Set[int](doc, 12345, "int")
	Assert(t, NoError(err), "new int must be set")

	i, err := genj.Get[int](doc, "int")
	Assert(t, NoError(err), "int must be accessible")
	Assert(t, Equal(i, 12345), "int must be correct")

	now := time.Now()
	err = genj.Set[time.Time](doc, now, "time")
	Assert(t, NoError(err), "time must be set")

	tm, err := genj.Get[time.Time](doc, "time")
	Assert(t, NoError(err), "time must be accessible")
	Assert(t, True(tm.Equal(now)), "time must be correct")

	dur := 5 * time.Second
	err = genj.Set[time.Duration](doc, dur, "duration")
	Assert(t, NoError(err), "duration must be set")

	d, err := genj.Get[time.Duration](doc, "duration")
	Assert(t, NoError(err), "duration must be accessible")
	Assert(t, Equal(d, dur), "duration must be correct")

	err = genj.Set(doc, 4711, "string")
	Assert(t, ErrorContains(err, "current element of type string does not match to new type int"), "string must not be set as int")

	// Nested valid and invalid.
	err = genj.Set(doc, "new value", "nested", "0", "d", "1")
	Assert(t, NoError(err), "string must be set")

	s, err = genj.Get[string](doc, "nested", "0", "d", "1")
	Assert(t, NoError(err), "string must be accessible")
	Assert(t, Equal(s, "new value"), "string must be correct")

	err = genj.Set(doc, 4711, "nested", "0", "d", "1")
	Assert(t, ErrorContains(err, "current element of type string does not match to new type int"), "string must not be set as int")
}

// TestSetWriteRead test the writing and reading of a JSON document with set values.
func TestSetWriteRead(t *testing.T) {
	doc, err := genj.Read(bytes.NewReader(createJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Set some values.
	err = genj.Set(doc, "new value", "string")
	Assert(t, NoError(err), "string must be set")

	err = genj.Set[int](doc, 12345, "int")
	Assert(t, NoError(err), "new int must be set")

	now := time.Now()
	err = genj.Set[time.Time](doc, now, "time")
	Assert(t, NoError(err), "time must be set")

	dur := 5 * time.Second
	err = genj.Set[time.Duration](doc, dur, "duration")
	Assert(t, NoError(err), "duration must be set")

	// Write and read.
	buf := &bytes.Buffer{}
	err = genj.Write(doc, buf)
	Assert(t, NoError(err), "document must be written w/o error")

	doc, err = genj.Read(buf)
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Test set values.
	s, err := genj.Get[string](doc, "string")
	Assert(t, NoError(err), "string must be accessible")
	Assert(t, Equal(s, "new value"), "string must be correct")

	i, err := genj.Get[int](doc, "int")
	Assert(t, NoError(err), "int must be accessible")
	Assert(t, Equal(i, 12345), "int must be correct")

	tm, err := genj.Get[time.Time](doc, "time")
	Assert(t, NoError(err), "time must be accessible")
	Assert(t, True(tm.Equal(now)), "time must be correct")

	d, err := genj.Get[time.Duration](doc, "duration")
	Assert(t, NoError(err), "duration must be accessible")
	Assert(t, Equal(d, dur), "duration must be correct")
}

// TestSetAny tests the setting of values in a JSON document.
func TestSetAny(t *testing.T) {
	doc, err := genj.Read(bytes.NewReader(createJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Valid setting.
	err = genj.SetAny(doc, 4711, "string")
	Assert(t, NoError(err), "string must be set to int value")

	i, err := genj.Get[int](doc, "string")
	Assert(t, NoError(err), "int must be accessible")
	Assert(t, Equal(i, 4711), "int must be correct")

	// Invalid setting.
	err = genj.SetAny(doc, 4711, "nested", "0", "x")
	Assert(t, ErrorContains(err, "not found"), "not existing path")

	err = genj.SetAny(doc, 4711, "nested", "0", "d")
	Assert(t, ErrorContains(err, "current element is not allowed to be an object or array"), "array to int")
}

//--------------------
// TESTS
//--------------------

// createJSON creates a simple JSON test document as bytes.
func createJSON() []byte {
	return []byte(`{
		"string": "value",
		"int": 42,
		"float": 3.1415,
		"bool": true,
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
