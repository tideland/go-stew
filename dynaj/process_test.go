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
	"errors"
	"fmt"
	"testing"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestProcess tests the processing of documents.
func TestProcess(t *testing.T) {
	bs, _ := createDocument(t)

	values := []string{}
	processor := func(node *dynaj.Node) error {
		value := fmt.Sprintf("%q = %q", node.Path(), node.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")

	// Verify iteration of all nodes.
	err = doc.Root().Process(processor)
	Assert(t, NoError(err), "document processed")
	Assert(t, Length(values, 29), "document processed all nodes")
	Assert(t, Contains(values, `"/A" = "Level One"`), "document processed all nodes")
	Assert(t, Contains(values, `"/B/0/D/A" = "Level Three - 0"`), "document processed all nodes")
	Assert(t, Contains(values, `"/B/0/D/B" = "10.1"`), "document processed all nodes")
	Assert(t, Contains(values, `"/B/1/S/2" = "white"`), "document processed all nodes")

	// Verifiy processing error.
	processor = func(node *dynaj.Node) error {
		return errors.New("ouch")
	}
	err = doc.Root().Process(processor)
	Assert(t, ErrorContains(err, "ouch"), "document processed with error")
}

// TestValueAtProcess tests the processing of documents starting at a
// deeper node.
func TestValueAtProcess(t *testing.T) {
	bs, _ := createDocument(t)

	values := []string{}
	processor := func(node *dynaj.Node) error {
		value := fmt.Sprintf("%q = %q", node.Path(), node.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")

	// Verify iteration of all nodes.
	err = doc.NodeAt("/B/0/D").Process(processor)
	Assert(t, NoError(err), "document processed")
	Assert(t, Length(values, 2), "document processed all nodes")
	Assert(t, Contains(values, `"/B/0/D/A" = "Level Three - 0"`), "document processed all nodes")
	Assert(t, Contains(values, `"/B/0/D/B" = "10.1"`), "document processed all nodes")

	values = []string{}
	err = doc.NodeAt("/B/1").Process(processor)
	Assert(t, NoError(err), "document processed")
	Assert(t, Length(values, 8), "document processed all nodes")
	Assert(t, Contains(values, `"/B/1/S/2" = "white"`), "document processed all nodes")
	Assert(t, Contains(values, `"/B/1/B" = "200"`), "document processed all nodes")

	// Verifiy iteration of non-existing path.
	err = doc.NodeAt("/B/3").Process(processor)
	Assert(t, ErrorContains(err, "invalid path"), "document processed with error")

	// Verify procesing error.
	processor = func(node *dynaj.Node) error {
		return errors.New("ouch")
	}
	err = doc.NodeAt("/A").Process(processor)
	Assert(t, ErrorContains(err, "ouch"), "document processed with error")
}

// TestRange tests the range processing of documents.
func TestRange(t *testing.T) {
	bs, _ := createDocument(t)

	values := []string{}
	processor := func(node *dynaj.Node) error {
		value := fmt.Sprintf("%q = %q", node.Path(), node.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")

	// Verify range of object.
	values = []string{}
	err = doc.NodeAt("/B/0/D").Range(processor)
	Assert(t, NoError(err), "node range processed")
	Assert(t, Length(values, 2), "note range processed all nodes")
	Assert(t, Contains(values, `"/B/0/D/A" = "Level Three - 0"`), "note range processed all nodes")
	Assert(t, Contains(values, `"/B/0/D/B" = "10.1"`), "note range processed all nodes")

	// Verify range of array.
	values = []string{}
	err = doc.NodeAt("/B/1/S").Range(processor)
	Assert(t, NoError(err), "note range processed")
	Assert(t, Length(values, 3), "note range processed all nodes")
	Assert(t, Contains(values, `"/B/1/S/0" = "orange"`), "note range processed all nodes")
	Assert(t, Contains(values, `"/B/1/S/1" = "blue"`), "note range processed all nodes")
	Assert(t, Contains(values, `"/B/1/S/2" = "white"`), "note range processed all nodes")

	// Verify range of value.
	values = []string{}
	err = doc.NodeAt("/A").Range(processor)
	Assert(t, NoError(err), "node range processed")
	Assert(t, Length(values, 1), "note range processed all nodes")
	Assert(t, Contains(values, `"/A" = "Level One"`), "note range processed all nodes")

	// Verify range of non-existing path.
	err = doc.NodeAt("/B/0/D/X").Range(processor)
	Assert(t, ErrorContains(err, "invalid path"), "node range processed with error")

	// Verify range of mixed types.
	err = doc.NodeAt("/B/0").Range(processor)
	Assert(t, ErrorContains(err, "is object or array"), "node range processed with error")

	// Verify procesing error.
	processor = func(node *dynaj.Node) error {
		return errors.New("ouch")
	}
	err = doc.NodeAt("/A").Range(processor)
	Assert(t, ErrorContains(err, "ouch"), "node range processed with error")
}

// TestRootQuery tests querying a document.
func TestRootQuery(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	nodes, err := doc.Root().Query("Z/*")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 0), "document queried invalid nodes")
	nodes, err = doc.Root().Query("*")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 29), "document queried all nodes")
	nodes, err = doc.Root().Query("/A")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 1), "document queried node /A")
	nodes, err = doc.Root().Query("/B/*")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 24), "document queried node /B/*")
	nodes, err = doc.Root().Query("/B/[01]/*")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 18), "document queried node /B/[01]/*")
	nodes, err = doc.Root().Query("/B/[01]/*A")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 4), "document queried node /B/[01]/*A")
	nodes, err = doc.Root().Query("*/S/*")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 8), "document queried node */S/*")
	nodes, err = doc.Root().Query("*/S/3")
	Assert(t, NoError(err), "document queried")
	Assert(t, Length(nodes, 1), "document queried node */S/3")

	// Verify the content
	nodes, err = doc.Root().Query("/A")
	Assert(t, NoError(err), "document queried")
	Assert(t, Equal(nodes[0].Path(), "/A"), "document queried node /A")
	Assert(t, Equal(nodes[0].AsString(""), "Level One"), "document queried node /A")
}

// TestValueAtQuery tests querying a document starting at a deeper node.
func TestValueAtQuery(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	nodes, err := doc.NodeAt("/B/0/D").Query("Z/*")
	Assert(t, NoError(err), "document queried deep")
	Assert(t, Length(nodes, 0), "document queried invalid nodes")
	nodes, err = doc.NodeAt("/B/0/D").Query("*")
	Assert(t, NoError(err), "document queried /B/0/D/*")
	Assert(t, Length(nodes, 2), "document queried node /B/0/D/*")
	nodes, err = doc.NodeAt("/B/0/D").Query("A")
	Assert(t, NoError(err), "document queried /B/0/D/A")
	Assert(t, Length(nodes, 1), "document queried node /B/0/D/A")
	nodes, err = doc.NodeAt("/B/0/D").Query("B")
	Assert(t, NoError(err), "document queried /B/0/D/B")
	Assert(t, Length(nodes, 1), "document queried node /B/0/D/B")
	nodes, err = doc.NodeAt("/B/0/D").Query("C")
	Assert(t, NoError(err), "document queried /B/0/D/C")
	Assert(t, Length(nodes, 0), "document queried node /B/0/D/C")
	nodes, err = doc.NodeAt("/B/1").Query("S/*")
	Assert(t, NoError(err), "document queried /B/1/S/*")
	Assert(t, Length(nodes, 3), "document queried node /B/1/S/*")
	nodes, err = doc.NodeAt("/B/1").Query("S/2")
	Assert(t, NoError(err), "document queried /B/1/S/2")
	Assert(t, Length(nodes, 1), "document queried node /B/1/S/2")

	// Verify non-existing path.
	nodes, err = doc.NodeAt("Z/Z/Z").Query("/A")
	Assert(t, ErrorContains(err, "invalid path"), "document queried with error")
	Assert(t, Length(nodes, 0), "document queried invalid nodes")
}

// EOF
