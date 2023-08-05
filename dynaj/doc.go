// Tideland Go Stew - Dynamic JSON
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package dynaj provides a simple dynamic handling of JSON documents.
// Values can be retrieved, set and added by paths like "foo/bar/3".
// Methods provide typesafe access to the values as well as flat and
// deep processing.
//
//	doc, err := dynaj.Unmarshal(myDoc)
//	if err != nil {
//	    ...
//	}
//	name, err := doc.At("name").AsString()
//	street, err := doc.At("address", "street").AsString()
//
// Another way is to create an empty document with
//
//	doc := dynaj.NewDocument()
//
// Here as well as in parsed documents values can be set with
//
//	err := doc.At("a", "b", "3").Set(4711)
//
// Additionally values of the document can be processed recursively
// using
//
//	err := doc.Root().Processor().DeepDo(func(acc *dynaj.Accessor) error {
//	    ...
//	})
//
// or from deeper nodes with doc.At("a", "b").Processor().Do(...).
//
// To retrieve the differences between two documents the function
// dynaj.Compare() can be used:
//
//	diff, err := dynaj.Compare(firstDoc, secondDoc)
//
// privides a dynaj.Diff instance which helps to compare individual
// paths of the two document.
package dynaj // import "tideland.dev/go/stew/dynaj"

// EOF
