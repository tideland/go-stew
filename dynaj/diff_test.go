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
	"testing"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestCompare tests comparing two documents.
func TestCompare(t *testing.T) {
	first, _ := createDocument(t)
	second := createCompareDocument(t)
	firstDoc, err := dynaj.Unmarshal(first)
	Assert(t, NoError(err), "first document unmarshalled")
	secondDoc, err := dynaj.Unmarshal(second)
	Assert(t, NoError(err), "second document unmarshalled")

	diff, err := dynaj.Compare(first, first)
	Assert(t, Length(diff.Differences(), 0), "first document compared with itself has no differences")
	Assert(t, NoError(err), "first unmarshalled document compared with itself")
	Assert(t, Empty(diff.Differences()), "first unmarshalled document compared with itself has no differences")

	diff, err = dynaj.Compare(first, second)
	Assert(t, NoError(err), "first unmarshalled document compared with second")
	Assert(t, Length(diff.Differences(), 15), "first unmarshalled document compared with second has differences")

	diff, err = dynaj.CompareDocuments(firstDoc, secondDoc)
	Assert(t, NoError(err), "first document compared with second")
	Assert(t, Length(diff.Differences(), 15), "first document compared with second has differences")

	for _, path := range diff.Differences() {
		fv, sv := diff.DifferenceAt(path)
		fvs := fv.AsString("<first undefined>")
		svs := sv.AsString("<second undefined>")
		Assert(t, Different(fvs, svs), fmt.Sprintf("first and second pathes different at %q", path))
	}

	first, err = diff.FirstDocument().MarshalJSON()
	Assert(t, NoError(err), "first document marshalled")
	second, err = diff.SecondDocument().MarshalJSON()
	Assert(t, NoError(err), "second document marshalled")
	diff, err = dynaj.Compare(first, second)
	Assert(t, NoError(err), "first marshalled document compared with second")
	Assert(t, Length(diff.Differences(), 15), "first marshalled document compared with second has differences")

	// Special case of empty arrays, objects, and null.
	first = []byte(`{}`)
	second = []byte(`{"a":[],"b":{},"c":null}`)

	sdocParsed, err := dynaj.Unmarshal(second)
	Assert(t, NoError(err), "second document unmarshalled")
	sdocMarshalled, err := sdocParsed.MarshalJSON()
	Assert(t, NoError(err), "second document marshalled")
	Assert(t, Equal(string(sdocMarshalled), string(second)), "second document marshalled equals original")

	diff, err = dynaj.Compare(first, second)
	Assert(t, NoError(err), "first document compared with second")
	Assert(t, Length(diff.Differences(), 4), "first document compared with second has differences")

	first = []byte(`[]`)
	diff, err = dynaj.Compare(first, second)
	Assert(t, NoError(err), "first document compared with second")
	Assert(t, Length(diff.Differences(), 4), "first document compared with second has differences")

	first = []byte(`["A", "B", "C"]`)
	diff, err = dynaj.Compare(first, second)
	Assert(t, NoError(err), "first document compared with second")
	Assert(t, Length(diff.Differences(), 6), "first document compared with second has differences")

	first = []byte(`"foo"`)
	diff, err = dynaj.Compare(first, second)
	Assert(t, NoError(err), "first document compared with second")
	Assert(t, Length(diff.Differences(), 4), "first document compared with second has differences")
}

// EOF
