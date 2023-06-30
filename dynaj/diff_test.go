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
	"testing"

	"tideland.dev/go/stew/asserts"
	"tideland.dev/go/stew/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestCompare tests comparing two documents.
func TestCompare(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	first, _ := createDocument(assert)
	second := createCompareDocument(assert)
	firstDoc, err := dynaj.Unmarshal(first)
	assert.NoError(err)
	secondDoc, err := dynaj.Unmarshal(second)
	assert.NoError(err)

	diff, err := dynaj.Compare(first, first)
	assert.NoError(err)
	assert.Length(diff.Differences(), 0)

	diff, err = dynaj.Compare(first, second)
	assert.NoError(err)
	assert.Length(diff.Differences(), 13)

	diff, err = dynaj.CompareDocuments(firstDoc, secondDoc)
	assert.NoError(err)
	assert.Length(diff.Differences(), 13)

	for _, path := range diff.Differences() {
		fv, sv := diff.DifferenceAt(path)
		fvs := fv.AsString("<first undefined>")
		svs := sv.AsString("<second undefined>")
		assert.Different(fvs, svs, path)
	}

	first, err = diff.FirstDocument().MarshalJSON()
	assert.NoError(err)
	second, err = diff.SecondDocument().MarshalJSON()
	assert.NoError(err)
	diff, err = dynaj.Compare(first, second)
	assert.NoError(err)
	assert.Length(diff.Differences(), 13)

	// Special case of empty arrays, objects, and null.
	first = []byte(`{}`)
	second = []byte(`{"a":[],"b":{},"c":null}`)

	sdocParsed, err := dynaj.Unmarshal(second)
	assert.NoError(err)
	sdocMarshalled, err := sdocParsed.MarshalJSON()
	assert.NoError(err)
	assert.Equal(string(sdocMarshalled), string(second))

	diff, err = dynaj.Compare(first, second)
	assert.NoError(err)
	assert.Length(diff.Differences(), 4)

	first = []byte(`[]`)
	diff, err = dynaj.Compare(first, second)
	assert.NoError(err)
	assert.Length(diff.Differences(), 4)

	first = []byte(`["A", "B", "C"]`)
	diff, err = dynaj.Compare(first, second)
	assert.NoError(err)
	assert.Length(diff.Differences(), 6)

	first = []byte(`"foo"`)
	diff, err = dynaj.Compare(first, second)
	assert.NoError(err)
	assert.Length(diff.Differences(), 4)
}

// EOF
