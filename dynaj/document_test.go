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
	"bytes"
	"testing"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestDocumentFromJSON tests the creation of a document from JSON.
func TestDocumentFromJSON(t *testing.T) {
	doc, err := dynaj.Unmarshal(createJSON())
	Assert(t, NoError(err), "document must be unmarshalled w/o error")
	Assert(t, NotNil(doc), "document must exist")

	Assert(t, NoError(doc.Root()), "root must be accessible")
	Assert(t, NoError(doc.At("object")), "object element must be accessible")
	Assert(t, NoError(doc.At("nested", "1", "d")), "nested element must be accessible")

	r := bytes.NewReader(createJSON())
	doc, err = dynaj.Read(r)
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")
}

//--------------------
// HELPER
//--------------------

// createDocument creates a simple JSON test document.
func createDocument() *dynaj.Document {
	doc, err := dynaj.Unmarshal([]byte(createJSON()))
	if err != nil {
		panic(err)
	}
	return doc
}

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
