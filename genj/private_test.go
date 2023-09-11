// Tideland Go Stew - Generic JSON - Private Unit Tests
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
	"bytes"
	"testing"

	. "tideland.dev/go/stew/qaone"
)

//--------------------
// TESTS
//--------------------

// TestMoonwalk tests the new moonwalk function.
func TestMoonwalk(t *testing.T) {
	doc, err := Read(bytes.NewReader(deepNestedJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Exact walks.
	id, elem, path, err := moonwalk(doc.root, Path{"l1a", "l2a", "a"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "a"), "path ID must be 'a'")
	Assert(t, Equal(elem, 1.0), "element must be 1")
	Assert(t, Nil(path), "path must be empty")

	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2b", "z", "l4a", "color"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "color"), "path ID must be 'color'")
	Assert(t, Equal(elem, "red"), "element must be 'red'")
	Assert(t, Nil(path), "path must be empty")

	// Short paths.
	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2b"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "l2b"), "path ID must be 'l2b'")
	object, ok := elem.(Object)
	Assert(t, OK(ok), "element must be an object")
	Assert(t, Length(object, 3), "element must be not nil")
	Assert(t, Nil(path), "path must be empty")

	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2c"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "l2c"), "path ID must be 'l2c'")
	array, ok := elem.(Array)
	Assert(t, OK(ok), "element must be an array")
	Assert(t, Length(array, 5), "element must be not nil")
	Assert(t, Nil(path), "path must be empty")

	// Long paths.
	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2b", "z", "l4a", "production", "count"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "count"), "path ID must be 'count'")
	object, ok = elem.(Object)
	Assert(t, OK(ok), "element must be an object")
	Assert(t, Length(object, 2), "element must have length 2")
	Assert(t, Length(path, 2), "path must contain production and count")

	// Little element test using contains.
	value, err := contains(elem, "color")
	Assert(t, NoError(err), "contains must be successful")
	Assert(t, Equal(value, "red"), "value must be 'red'")

	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2b", "x", "999"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "999"), "path ID must be '999'")
	array, ok = elem.(Array)
	Assert(t, OK(ok), "element must be an array")
	Assert(t, Length(array, 3), "element have length 3")
	Assert(t, Length(path, 1), "path must have length 1")

	// Illegal paths.
	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2b", "x", "-99", "myid"})
	Assert(t, ErrorContains(err, `array index -99 out of bounds`), "path must not be walked")
	Assert(t, Equal(id, "myid"), "path ID must be 'myid'")
	Assert(t, NotNil(elem), "element must be not nil")
	Assert(t, Length(path, 2), "path contains a rest")

	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2a", "a", "not", "existing"})
	Assert(t, ErrorMatches(err, `path .* is illegal`), "path must not be walkable")
	Assert(t, Equal(id, "existing"), "path ID must be 'existing'")
	Assert(t, Nil(elem), "element must be nil")
	Assert(t, Length(path, 2), "path contains a rest")
}

// TestContains tests the contains function.
func TestContains(t *testing.T) {
	doc, err := Read(bytes.NewReader(deepNestedJSON()))
	Assert(t, NoError(err), "document must be read w/o error")
	Assert(t, NotNil(doc), "document must exist")

	// Elements in Object.
	id, elem, path, err := moonwalk(doc.root, Path{"l1a", "l2b"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "l2b"), "path ID must be 'l2b'")
	Assert(t, Length(path, 0), "path must be empty")

	value, err := contains(elem, "x")
	Assert(t, NoError(err), "contains must be successful")
	Assert(t, NotNil(value), "value must be not nil")
	array, ok := value.(Array)
	Assert(t, OK(ok), "value must be an array")
	Assert(t, Length(array, 3), "value must have length 3")

	value, err = contains(elem, "z")
	Assert(t, NoError(err), "contains must be successful")
	Assert(t, NotNil(value), "value must be not nil")
	object, ok := value.(Object)
	Assert(t, OK(ok), "value must be an object")
	Assert(t, Length(object, 2), "value must have length 2")

	value, err = contains(elem, "not")
	Assert(t, ErrorContains(err, `element "not" not found`), "contains must fail")
	Assert(t, Nil(value), "value must be nil")

	// Elements in Array.
	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2c"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "l2c"), "path ID must be 'l2c'")
	Assert(t, Length(path, 0), "path must be empty")

	value, err = contains(elem, "0")
	Assert(t, NoError(err), "contains must be successful")
	Assert(t, Equal(value, "one"), "value must be 'one'")

	value, err = contains(elem, "999")
	Assert(t, ErrorContains(err, `index 999 out of bounds`), "contains must fail")
	Assert(t, Nil(value), "value must be nil")

	// Neither Object nor Array.
	id, elem, path, err = moonwalk(doc.root, Path{"l1a", "l2a", "a"})
	Assert(t, NoError(err), "path must be walked w/o error")
	Assert(t, Equal(id, "a"), "path ID must be 'a'")
	Assert(t, Length(path, 0), "path must be empty")

	value, err = contains(elem, "not")
	Assert(t, ErrorContains(err, `element is no Object or Array`), "cannot test element type")
	Assert(t, Nil(value), "value must be nil")
}

//--------------------
// HELPER
//--------------------

// deepNestedJSON creates a deep nested tree of elements.
func deepNestedJSON() []byte {
	return []byte(`{
		"l1a": {
			"l2a": {
				"a": 1,
				"b": 2
			},
			"l2b": {
				"x": [1, 2, 3],
				"y": [4, 5, 6],
				"z": {
					"l4a": {
						"color": "red",
						"size":  42
					},
					"l4b": {
						"color": "blue",
						"size":  23
					}
				}
			},
			"l2c": [
				"one",
				"two",
				"three",
				"four",
				{
					"l3a": ["a", "b", "c"],
					"l3b": ["d", "e", "f"]
				}
			]
		}
	}`)
}

// EOF
