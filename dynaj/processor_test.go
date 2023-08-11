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

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestProcessDo verifies the processing on elements and
// trees of a JSON document.
func TestProcessDo(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("string")
	acc.Processor().Do(func(acc *dynaj.Accessor) error {
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
	acc.Processor().Do(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err = doc.At("array", "0").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "first element should be 'new'")

	acc = doc.At("object")
	acc.Processor().Do(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err = doc.At("object", "one").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "first element should be 'new'")

	// Negetive tests.
	acc = doc.At("array")
	err = acc.Processor().Do(func(acc *dynaj.Accessor) error {
		return fmt.Errorf("ouch")
	}).Err()
	Assert(t, ErrorContains(err, "ouch"), "error expected")
}

// TestProcessDeepDo verifies the processing on elements and
// trees of a JSON document.
func TestProcessDeepDo(t *testing.T) {
	doc := createDocument()

	// Positive tests.
	acc := doc.At("string")
	acc.Processor().DeepDo(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err := doc.At("string").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "the element should be 'new'")

	acc = doc.At("nested")
	acc.Processor().DeepDo(func(acc *dynaj.Accessor) error {
		return acc.Update("new").Err()
	})
	s, err = doc.At("nested", "1", "d", "0").AsString()
	Assert(t, NoError(err), "no error expected")
	Assert(t, Equal(s, "new"), "deep nested elements should be 'new'")

	// Negetive tests.
	acc = doc.At("array")
	err = acc.Processor().DeepDo(func(acc *dynaj.Accessor) error {
		return fmt.Errorf("ouch")
	}).Err()
	Assert(t, ErrorContains(err, "ouch"), "error expected")
}

// EOF
