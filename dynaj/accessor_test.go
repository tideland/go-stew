// Tideland Go Stew - Dynamic JSON - Unit Tests
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dynaj_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestAccessAsString verifies the access to values as string.
func TestAccessAsString(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	s, err := doc.At("string").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "value"), "accessor must return right string")

	s, err = doc.At("int").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "42"), "accessor must return right string")

	s, err = doc.At("float").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "3.1415"), "accessor must return right string")

	s, err = doc.At("bool").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "true"), "accessor must return right string")

	s, err = doc.At("time").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "2019-09-01T12:00:00Z"), "accessor must return right string")

	s, err = doc.At("duration").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "1h30m"), "accessor must return right string")

	s, err = doc.At("array").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "[...]"), "accessor must return right string")

	s, err = doc.At("object").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "{...}"), "accessor must return right string")

	// Negative tests.
	s, err = doc.At("not", "existing").AsString()
	Assert(t, ErrorContains(err, `invalid path [not]`), "accessor must return error")
	Assert(t, Equal(s, ""), "accessor must return no value")

	acc := doc.At("not", "existing")
	Assert(t, Equal(acc.Type(), dynaj.TypeError), "accessor must return error")
	acc = doc.At("not", "existing")
	Assert(t, ErrorContains(acc, `invalid path [not]`), "accessor must return error")
}

// TestAccessAsInt verifies the access to values as int.
func TestAccessAsInt(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	i, err := doc.At("int").AsInt()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(i, 42), "accessor must return right int")

	i, err = doc.At("intstring").AsInt()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(i, 4711), "accessor must return right  int")

	i, err = doc.At("float").AsInt()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(i, 3), "accessor must return right  int")

	i, err = doc.At("bool").AsInt()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(i, 1), "accessor must return right int")

	i, err = doc.At("array").AsInt()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(i, 3), "accessor must return length")

	i, err = doc.At("object").AsInt()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(i, 3), "accessor must return length")

	// Negative tests.
	i, err = doc.At("string").AsInt()
	Assert(t, ErrorContains(err, `cannot convert "value" to int`), "accessor must return error")
	Assert(t, Equal(i, 0), "accessor must return no value")

	i, err = doc.At("not", "existing").AsInt()
	Assert(t, ErrorContains(err, `invalid path [not]`), "accessor must return error")
	Assert(t, Equal(i, 0), "accessor must return no value")
}

// TestAccessAsFloat64 verifies the access to values as float.
func TestAccessAsFloat64(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	f, err := doc.At("int").AsFloat64()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(f, 42.0), "accessor must return right float")

	f, err = doc.At("float").AsFloat64()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(f, 3.1415), "accessor must return right float")

	f, err = doc.At("floatstring").AsFloat64()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(f, 2.7182), "accessor must return right float")

	f, err = doc.At("bool").AsFloat64()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(f, 1.0), "accessor must return right float")

	f, err = doc.At("array").AsFloat64()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(f, 3.0), "accessor must return length")

	f, err = doc.At("object").AsFloat64()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(f, 3.0), "accessor must return length")

	// Negative tests.
	f, err = doc.At("string").AsFloat64()
	Assert(t, ErrorContains(err, `cannot convert "value" to float64`), "accessor must return error")
	Assert(t, Equal(f, 0.0), "accessor must return no value")

	f, err = doc.At("not", "existing").AsFloat64()
	Assert(t, ErrorContains(err, `invalid path [not]`), "accessor must return error")
	Assert(t, Equal(f, 0.0), "accessor must return no value")
}

// TestAccessAsBool verifies the access to values as bool.
func TestAccessAsBool(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	b, err := doc.At("bool").AsBool()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, True(b), "accessor must return right bool")

	b, err = doc.At("boolstring").AsBool()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, True(b), "accessor must return right bool")

	b, err = doc.At("int").AsBool()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, True(b), "accessor must return right bool")

	b, err = doc.At("float").AsBool()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, True(b), "accessor must return right bool")

	// Negative tests.
	b, err = doc.At("string").AsBool()
	Assert(t, ErrorContains(err, `invalid syntax`), "accessor must return error")
	Assert(t, False(b), "accessor must return no value")

	b, err = doc.At("not").At("existing").AsBool()
	Assert(t, ErrorContains(err, `invalid path [not]`), "accessor must return error")
	Assert(t, False(b), "accessor must return no value")
}

// TestAccessAsTime verifies the access to values as time.Time.
func TestAccessAsTime(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	tt, err := doc.At("time").AsTime(time.RFC3339Nano)
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(tt, time.Date(2019, 9, 1, 12, 0, 0, 0, time.UTC)), "accessor must return right time")

	tt, err = doc.At("int").AsTime(time.RFC3339Nano)
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(tt, time.Unix(42, 0)), "accessor must return right time")

	tt, err = doc.At("float").AsTime(time.RFC3339Nano)
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(tt, time.Unix(3, 0)), "accessor must return right time")

	// Negative tests.
	tt, err = doc.At("string").AsTime(time.RFC3339Nano)
	Assert(t, ErrorContains(err, `cannot convert "value" to time.Time`), "accessor must return error")
	Assert(t, Equal(tt, time.Time{}), "accessor must return no value")

	tt, err = doc.At("not").At("existing").AsTime(time.RFC3339Nano)
	Assert(t, ErrorContains(err, `invalid path [not]`), "accessor must return error")
	Assert(t, Equal(tt, time.Time{}), "accessor must return no value")
}

// TestAccessAsDuration verifies the access to values as time.Duration.
func TestAccessAsDuration(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	d, err := doc.At("duration").AsDuration()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(d, time.Hour+30*time.Minute), "accessor must return right duration")

	d, err = doc.At("int").AsDuration()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(d, 42*time.Second), "accessor must return right duration")

	d, err = doc.At("float").AsDuration()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(d, 3*time.Second), "accessor must return right duration")

	// Negative tests.
	d, err = doc.At("string").AsDuration()
	Assert(t, ErrorContains(err, `cannot convert "value" to time.Duration`), "accessor must return error")
	Assert(t, Equal(d, time.Duration(0)), "accessor must return no value")

	d, err = doc.At("not", "existing").AsDuration()
	Assert(t, ErrorContains(err, `invalid path [not]`), "accessor must return error")
	Assert(t, Equal(d, time.Duration(0)), "accessor must return no value")
}

// TestAccessUpdate verifies the update of values.
func TestAccessUpdate(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("string").Update("new value")
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, NotNil(acc), "accessor must be returned")
	s, err := acc.AsString()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(s, "new value"), "accessor must return right string")
	s, err = doc.At("string").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "new value"), "accessor must return right string")

	acc = doc.At("nested", "0", "a").Update(4711)
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, NotNil(acc), "accessor must be returned")
	i, err := acc.AsInt()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(i, 4711), "accessor must return right int")
	i, err = doc.At("nested", "0", "a").AsInt()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(i, 4711), "accessor must return right int")

	acc = doc.At("nested", "1", "d", "2").Update("yadda")
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, NotNil(acc), "accessor must be returned")
	s, err = acc.AsString()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(s, "yadda"), "accessor must return right string")
	s, err = doc.At("nested", "1", "d", "2").AsString()
	Assert(t, NoError(err), "accessor must be created and used w/o error")
	Assert(t, Equal(s, "yadda"), "accessor must return right string")

	acc = doc.At("string").Update(dynaj.Array{})
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, NotNil(acc), "accessor must be returned")

	// Negative tests.
	acc = doc.At("not", "existing").Update("yadda")
	Assert(t, ErrorContains(acc, `invalid path [not]`), "accessor must return error")

	acc = doc.At("nested", "0", "a").Update(dynaj.Object{"foo": 12345})
	Assert(t, ErrorContains(acc, `invalid type or non-empty container type for update`), "accessor must return error")

	// Final positive test.
	acc = doc.Root().Update("1-2-3")
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, Length(doc.Root(), 1), "document must have length 1")
	acc = doc.Root()
	s, err = acc.AsString()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(s, "1-2-3"), "accessor must return right string")
}

// TestAccessSet verifies the setting of values.
func TestAccessSet(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.Root().Set("new", "value")
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	s, err := acc.AsString()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(s, "value"), "accessor must return right string")
	s, err = doc.At("new").AsString()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(s, "value"), "accessor must return right string")

	acc = doc.At("nested", "0", "d").Set("0", 4711)
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	i, err := acc.AsInt()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(i, 4711), "accessor must return right int")

	// Negative tests.
	acc = doc.At("not", "existing").Set("yadda", "yadda")
	Assert(t, ErrorContains(acc, `invalid path [not]`), "accessor must return error")

	acc = doc.At("string").Set("foo", "bar")
	Assert(t, ErrorContains(acc, `cannot set element: not an array or object`), "accessor must return error")

	acc = doc.At("array").Set("foo", "bar")
	Assert(t, ErrorContains(acc, `cannot set element: illegal index`), "accessor must return error")

	// Complex positive scenario.
	acc = doc.At("object").Set("nested", dynaj.Object{})
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	sub := doc.At("object", "nested").Set("foo", 1)
	Assert(t, NoError(sub), "accessor must be created and used w/o error")
	sub = doc.At("object", "nested").Set("bar", 2)
	Assert(t, NoError(sub), "accessor must be created and used w/o error")
	sub = doc.At("object", "nested").Set("baz", 3)
	Assert(t, NoError(sub), "accessor must be created and used w/o error")
	Assert(t, Length(doc.At("object", "nested"), 3), "document must have length 3")
}

// TestAccessAppend verifies the appending of values.
func TestAccessAppend(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("array").Append("new value")
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, Length(doc.At("array"), 4), "document must have length 4")
	s, err := doc.At("array", "3").AsString()
	Assert(t, NoError(err), "value must be returned w/o error")
	Assert(t, Equal(s, "new value"), "accessor must return right string")

	// Negative tests.
	acc = doc.At("not", "existing").Append("yadda")
	Assert(t, ErrorContains(acc, `invalid path [not]`), "accessor must return error")

	acc = doc.At("string").Append("yadda")
	Assert(t, ErrorContains(acc, `cannot append element: not an array`), "accessor must return error")
}

// TestAccessDelete verifies the deletion of values.
func TestAccessDelete(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("string").Delete()
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, Length(doc.Root(), 11), "document must have length 11")
	acc = doc.At("string")
	Assert(t, ErrorContains(acc, `invalid path [string]`), "accessor must return error")

	acc = doc.At("array", "0").Delete()
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, Length(doc.At("array"), 2), "document at location must have length 2")

	acc = doc.At("nested", "0", "d").Delete()
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, Length(doc.At("nested", "0"), 3), "document at location must have length 3")
	acc = doc.At("nested", "0", "d")
	Assert(t, ErrorContains(acc, `invalid path [nested 0 d]`), "accessor must return error")

	acc = doc.At("nested", "1").Delete()
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, Length(doc.At("nested"), 1), "document at location must have length 1")

	// Negative tests.
	acc = doc.At("not", "existing").Delete()
	Assert(t, ErrorContains(acc, `invalid path [not]`), "accessor must return error")

	// Hard positive test.
	acc = doc.Root().Delete()
	Assert(t, NoError(acc), "accessor must be created and used w/o error")
	Assert(t, Length(doc.Root(), 0), "document must have length 0")
}

// TestAccessAt verifies the access to values at a location based on
// another accessor.
func TestAccessAt(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("nested").At("0").At("d").At("1")
	Assert(t, NotNil(acc), "accessor must be created")
	Assert(t, NoError(acc), "accessor must be used w/o error")

	s, err := acc.AsString()
	Assert(t, NoError(err), "accessor must be used w/o error")
	Assert(t, Equal(s, "bar"), "accessor must return right string")

	s, err = doc.At("nested").At("1").At("a").AsString()
	Assert(t, NoError(err), "accessor must be used w/o error")
	Assert(t, Equal(s, "9"), "accessor must return right string")

	s, err = doc.At("nested").At("0", "d", "0").AsString()
	Assert(t, NoError(err), "accessor must be used w/o error")
	Assert(t, Equal(s, "foo"), "accessor must return right string")

	// Negative tests.
	acc = doc.At("nested").At("2").At("z").At("99")
	Assert(t, NotNil(acc), "accessor must be created")
	Assert(t, ErrorContains(acc, `invalid path [nested 2]`), "accessor must return error")
}

// TestAccessDo verifies the looping on elements and
// trees of a JSON document.
func TestAccessDo(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("string")
	acc.Do(func(acc *dynaj.Accessor) error {
		s, err := acc.AsString()
		if err != nil {
			return err
		}
		return acc.Update(strings.ToUpper(s)).Err()
	})
	s, err := doc.At("string").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "VALUE"), "string should be 'VALUE'")

	acc = doc.At("array")
	acc.Do(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err = doc.At("array", "0").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "first element should be 'new'")

	acc = doc.At("object")
	acc.Do(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err = doc.At("object", "one").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "first element should be 'new'")

	// Negetive tests.
	acc = doc.At("array")
	err = acc.Do(func(acc *dynaj.Accessor) error {
		return fmt.Errorf("ouch")
	}).Err()
	Assert(t, ErrorContains(err, "ouch"), "error expected")
}

// TestAccessDeepDo verifies the loop on elements and
// trees of a JSON document.
func TestAccessDeepDo(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("string")
	acc.DeepDo(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err := doc.At("string").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "the element should be 'new'")

	acc = doc.At("nested")
	acc.DeepDo(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err = doc.At("nested", "1", "d", "0").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "deep nested elements should be 'new'")

	// Negetive tests.
	acc = doc.At("array")
	err = acc.DeepDo(func(acc *dynaj.Accessor) error {
		return fmt.Errorf("ouch")
	}).Err()
	Assert(t, ErrorContains(err, "ouch"), "error expected")
}

// EOF
