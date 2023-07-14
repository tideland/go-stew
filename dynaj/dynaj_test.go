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
	"encoding/json"
	"testing"
	"time"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestBuildTypes tests the creation of documents with different
// root types.
func TestBuildTypes(t *testing.T) {
	// Just one value.
	doc := dynaj.NewDocument()
	err := doc.SetValueAt("", "foo")
	Assert(t, NoError(err), "value / set")

	sv := doc.NodeAt("").AsString("bar")
	Assert(t, Equal(sv, "foo"), "value / retrieved")

	// Now an object.
	doc = dynaj.NewDocument()
	err = doc.SetValueAt("/a", 1)
	Assert(t, NoError(err), "value /a set")

	nv := doc.NodeAt("")
	Assert(t, True(nv.IsObject()), "value / is object")
	iv := doc.NodeAt("a").AsInt(-1)
	Assert(t, Equal(iv, 1), "value /a retrieved")

	// And finally an array.
	doc = dynaj.NewDocument()
	err = doc.SetValueAt("/0", 1)
	Assert(t, NoError(err), "value /0 set")

	nv = doc.NodeAt("")
	Assert(t, True(nv.IsArray()), "value / is array")
	iv = doc.NodeAt("0").AsInt(-1)
	Assert(t, Equal(iv, 1), "value /0 retrieved")
}

// TestBuilding tests the creation of documents.
func TestBuilding(t *testing.T) {
	// Positive cases.
	doc := dynaj.NewDocument()
	err := doc.SetValueAt("/a/b/x", 1)
	Assert(t, NoError(err), "value /a/b/x set")
	err = doc.SetValueAt("/a/b/y", true)
	Assert(t, NoError(err), "value /a/b/y set")
	err = doc.SetValueAt("/a/b/0", "foo")
	Assert(t, ErrorContains(err, "cannot insert value"), "cannot insert array value in object")
	err = doc.SetValueAt("/a/c", "quick brown fox")
	Assert(t, NoError(err), "value /a/c set")
	err = doc.SetValueAt("/a/d/0/z", 47.11)
	Assert(t, NoError(err), "value /a/d/0/z set")
	err = doc.SetValueAt("/a/d/1/z", nil)
	Assert(t, NoError(err), "value /a/d/1/z set")
	err = doc.SetValueAt("/a/d/2", 2)
	Assert(t, NoError(err), "value /a/d/2 set")

	iv := doc.NodeAt("a/b/x").AsInt(0)
	Assert(t, Equal(iv, 1), "value aa/b/x retrieved")
	iv = doc.NodeAt("a").NodeAt("b").NodeAt("x").AsInt(0)
	Assert(t, Equal(iv, 1), "value x in b in a retrieved")
	nv := doc.NodeAt("a").NodeAt("b").NodeAt("x")
	Assert(t, Equal(nv.Path(), "/a/b/x"), "path to x correct")
	bv := doc.NodeAt("a/b/y").AsBool(false)
	Assert(t, True(bv), "value /a/b/y retrieved")
	sv := doc.NodeAt("a/c").AsString("")
	Assert(t, Equal(sv, "quick brown fox"), "value a/c retrieved")
	fv := doc.NodeAt("a/d/0/z").AsFloat64(8.15)
	Assert(t, Equal(fv, 47.11), "value a/d/0/z retrieved")
	nvt := doc.NodeAt("a/d/1/z").IsUndefined()
	Assert(t, True(nvt), "value a/d/1/z retrieved, it's undefined")

	nodes, err := doc.Root().Query("*x")
	Assert(t, NoError(err), "query for *x")
	Assert(t, Length(nodes, 1), "one node found")

	// Now provoke errors.
	err = doc.SetValueAt("/a/d", "stupid")
	Assert(t, ErrorContains(err, "cannot insert value"), "cannot insert 1st value")
	err = doc.SetValueAt("/a/d/0", "stupid")
	Assert(t, ErrorContains(err, "cannot insert value"), "cannot insert 2nd value")
	err = doc.SetValueAt("/a/d/2/z", "stupid")
	Assert(t, ErrorContains(err, "cannot insert value"), "cannot insert 3rd value")
	err = doc.SetValueAt("/a/b/y/z", "stupid")
	Assert(t, ErrorContains(err, "cannot insert value"), "cannot insert 4th value")

	err = doc.SetValueAt("a", "stupid")
	Assert(t, ErrorMatches(err, ".*corrupt.*"), "cannot currupt tree")
	err = doc.SetValueAt("a/b/x/y", "stupid")
	Assert(t, ErrorMatches(err, ".*corrupt.*"), "cannot currupt tree")
	err = doc.SetValueAt("/a/d/x", "stupid")
	Assert(t, ErrorMatches(err, ".*invalid index.*"), "cannot address invalid index")
	err = doc.SetValueAt("/a/d/-1", "stupid")
	Assert(t, ErrorMatches(err, ".*negative index.*"), "cannot address invalid index")

	// Legal change of values.
	err = doc.SetValueAt("/a/b/x", 2)
	Assert(t, NoError(err), "value /a/b/x changed")
	iv = doc.NodeAt("a/b/x").AsInt(0)
	Assert(t, Equal(iv, 2), "value /a/b/x retrieved")
}

// TestDeleteValueAt tests the deletion of values.
func TestDeleteValueAt(t *testing.T) {
	// Create a document.
	doc := dynaj.NewDocument()
	err := doc.SetValueAt("/obj/a", 1)
	Assert(t, NoError(err), "value /obj/a set")
	err = doc.SetValueAt("/obj/b", 2)
	Assert(t, NoError(err), "value /obj/b set")
	err = doc.SetValueAt("/obj/c", 3)
	Assert(t, NoError(err), "value /obj/c set")
	err = doc.SetValueAt("/obj/d/x", "x")
	Assert(t, NoError(err), "value /obj/d/x set")
	err = doc.SetValueAt("/obj/d/y", "y")
	Assert(t, NoError(err), "value /obj/d/y set")
	err = doc.SetValueAt("/obj/d/z", "z")
	Assert(t, NoError(err), "value /obj/d/z set")

	err = doc.SetValueAt("/arr/0", "foo")
	Assert(t, NoError(err), "value /arr/0 set")
	err = doc.SetValueAt("/arr/1", "bar")
	Assert(t, NoError(err), "value /arr/1 set")
	err = doc.SetValueAt("/arr/2", "baz")
	Assert(t, NoError(err), "value /arr/2 set")

	err = doc.SetValueAt("/val", true)
	Assert(t, NoError(err), "value /val set")

	// Delete values, object removes key, array shifts.
	err = doc.DeleteValueAt("/obj/b")
	Assert(t, NoError(err), "value /obj/b deleted")
	node := doc.NodeAt("/obj/b")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /obj/b not found")
	err = doc.DeleteValueAt("/obj/d/z")
	Assert(t, NoError(err), "value /obj/d/z deleted")

	err = doc.DeleteValueAt("/arr/1")
	Assert(t, NoError(err), "value /arr/1 deleted")
	node = doc.NodeAt("/arr/1")
	Assert(t, Equal(node.AsString("ouch"), "baz"), "value /arr/1 is baz")
	node = doc.NodeAt("/arr/2")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /arr/2 not found")

	err = doc.DeleteValueAt("/val")
	Assert(t, NoError(err), "value /val deleted")
	node = doc.NodeAt("/val")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /val not found")
	err = doc.DeleteValueAt("/not_found")
	Assert(t, NoError(err), "don't care about not found")

	// Provoke errors.
	err = doc.DeleteValueAt("/obj/a/not_found")
	Assert(t, ErrorContains(err, "path too long"), "path too long")
	err = doc.DeleteValueAt("/obj/d")
	Assert(t, ErrorContains(err, "is no value"), "is no value")

	err = doc.DeleteValueAt("/deep/not_found")
	Assert(t, ErrorContains(err, "invalid path"), "invalid path")
}

// TestDeleteElementAt tests the deletion of elements.
func TestDeleteElementAt(t *testing.T) {
	// Create a document.
	doc := dynaj.NewDocument()
	err := doc.SetValueAt("/obj/a", 1)
	Assert(t, NoError(err), "value /obj/a set")
	err = doc.SetValueAt("/obj/b", 2)
	Assert(t, NoError(err), "value /obj/b set")
	err = doc.SetValueAt("/obj/c", 3)
	Assert(t, NoError(err), "value /obj/c set")
	err = doc.SetValueAt("/obj/d/x", "x")
	Assert(t, NoError(err), "value /obj/d/x set")
	err = doc.SetValueAt("/obj/d/y", "y")
	Assert(t, NoError(err), "value /obj/d/y set")
	err = doc.SetValueAt("/obj/d/z", "z")
	Assert(t, NoError(err), "value /obj/d/z set")

	err = doc.SetValueAt("/arr/0", "foo")
	Assert(t, NoError(err), "value /arr/0 set")
	err = doc.SetValueAt("/arr/1", "bar")
	Assert(t, NoError(err), "value /arr/1 set")
	err = doc.SetValueAt("/arr/2", "baz")
	Assert(t, NoError(err), "value /arr/2 set")

	err = doc.SetValueAt("/val", true)
	Assert(t, NoError(err), "value /val set")

	// Delete elements.
	err = doc.DeleteElementAt("/obj/d/x")
	Assert(t, NoError(err), "value /obj/d/x deleted")
	node := doc.NodeAt("/obj/d/x")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /obj/d/x not found")
	err = doc.DeleteElementAt("/obj/d")
	Assert(t, NoError(err), "value /obj/d deleted")
	node = doc.NodeAt("/obj/d/y")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /obj/d/y not found")
	err = doc.DeleteElementAt("/obj")
	Assert(t, NoError(err), "value /obj deleted")
	node = doc.NodeAt("/obj")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /obj not found")

	err = doc.DeleteElementAt("/arr")
	Assert(t, NoError(err), "value /arr deleted")
	node = doc.NodeAt("/arr")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /arr not found")

	err = doc.DeleteElementAt("/val")
	Assert(t, NoError(err), "value /val deleted")
	node = doc.NodeAt("/val")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "value /val not found")

	// Provoke errors.
	err = doc.SetValueAt("/obj/a", 1)
	Assert(t, NoError(err), "value /obj/a set")
	err = doc.DeleteElementAt("/obj/a/not_found")
	Assert(t, ErrorContains(err, "path too long"), "path too long")

	err = doc.DeleteElementAt("/deep/not_found")
	Assert(t, ErrorContains(err, "invalid path"), "invalid path")
}

// TestParseError tests the returned error in case of
// an invalid document.
func TestParseError(t *testing.T) {
	bs := []byte(`abc{def`)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, Nil(doc), "document is nil")
	Assert(t, ErrorContains(err, "cannot unmarshal document"), "invalid document")
}

// TestClear tests to clear a document.
func TestClear(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	doc.Clear()
	err = doc.SetValueAt("/", "foo")
	Assert(t, NoError(err), "value / set")
	foo := doc.NodeAt("/").AsString("<undefined>")
	Assert(t, Equal(foo, "foo"), "value / retrieved")
}

// TestLength tests retrieving values as strings.
func TestLength(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	l := doc.Length("X")
	Assert(t, Equal(l, -1), "length of undefined X")
	l = doc.Length("")
	Assert(t, Equal(l, 4), "length of root")
	l = doc.Length("B")
	Assert(t, Equal(l, 3), "length of B")
	l = doc.Length("B/2")
	Assert(t, Equal(l, 5), "length of B/2")
	l = doc.Length("/B/2/D")
	Assert(t, Equal(l, 2), "length of /B/2/D")
	l = doc.Length("/B/1/S")
	Assert(t, Equal(l, 3), "length of /B/1/S")
	l = doc.Length("/B/1/S/0")
	Assert(t, Equal(l, 1), "length of /B/1/S/0")
}

// TestNotFound tests the handling of not found values.
func TestNotFound(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")

	// Check if is undefined.
	node := doc.NodeAt("you-wont-find-me")
	Assert(t, True(node.IsUndefined()), "node is undefined")
	Assert(t, AnyError(node.Err()), "node in error mode")
	Assert(t, ErrorContains(node.Err(), "invalid path"), "node has invalid path error")
}

// TestString verifies the string representation of a document.
func TestString(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	s := doc.String()
	Assert(t, Equal(s, string(bs)), "document stringified correctly")
}

// TestAsString tests retrieving values as strings.
func TestAsString(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	sv := doc.NodeAt("A").AsString("default")
	Assert(t, Equal(sv, "Level One"), "value A retrieved")
	sv = doc.NodeAt("B/0/B").AsString("default")
	Assert(t, Equal(sv, "100"), "value B/0/B retrieved")
	sv = doc.NodeAt("B/0/C").AsString("default")
	Assert(t, Equal(sv, "true"), "value B/0/C retrieved")
	sv = doc.NodeAt("B/0/D/B").AsString("default")
	Assert(t, Equal(sv, "10.1"), "value B/0/D/B retrieved")
	sv = doc.NodeAt("Z/Z/Z").AsString("default")
	Assert(t, Equal(sv, "default"), "value Z/Z/Z retrieved")

	sv = doc.NodeAt("A").String()
	Assert(t, Equal(sv, "Level One"), "value A retrieved")
	sv = doc.NodeAt("Z/Z/Z").String()
	Assert(t, Equal(sv, "null"), "value Z/Z/Z retrieved")
	// assert.Contains("invalid path", sv)

	// Difference between invalid path and nil value.
	doc.SetValueAt("Z/Z/Z", nil)
	sv = doc.NodeAt("Z/Z/Z").String()
	Assert(t, Equal(sv, "null"), "value Z/Z/Z retrieved")
}

// TestAsInt tests retrieving values as ints.
func TestAsInt(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	iv := doc.NodeAt("A").AsInt(-1)
	Assert(t, Equal(iv, -1), "value A retrieved")
	iv = doc.NodeAt("B/0/B").AsInt(-1)
	Assert(t, Equal(iv, 100), "value B/0/B retrieved")
	iv = doc.NodeAt("B/0/C").AsInt(-1)
	Assert(t, Equal(iv, 1), "value B/0/C retrieved")
	iv = doc.NodeAt("B/0/S/2").AsInt(-1)
	Assert(t, Equal(iv, 1), "value B/0/S/2 retrieved")
	iv = doc.NodeAt("B/0/D/B").AsInt(-1)
	Assert(t, Equal(iv, 10), "value B/0/D/B retrieved")
	iv = doc.NodeAt("Z/Z/Z").AsInt(-1)
	Assert(t, Equal(iv, -1), "value Z/Z/Z retrieved")
}

// TestAsFloat64 tests retrieving values as float64.
func TestAsFloat64(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	fv := doc.NodeAt("A").AsFloat64(-1.0)
	Assert(t, Equal(fv, -1.0), "value A retrieved")
	fv = doc.NodeAt("B/0/B").AsFloat64(-1.0)
	Assert(t, Equal(fv, 100.0), "value B/0/B retrieved")
	fv = doc.NodeAt("B/1/B").AsFloat64(-1.0)
	Assert(t, Equal(fv, 200.0), "value B/1/B retrieved")
	fv = doc.NodeAt("B/0/C").AsFloat64(-99)
	Assert(t, Equal(fv, 1.0), "value B/0/C retrieved")
	fv = doc.NodeAt("B/0/S/3").AsFloat64(-1.0)
	Assert(t, Equal(fv, 2.2), "value B/0/S/3 retrieved")
	fv = doc.NodeAt("B/1/D/B").AsFloat64(-1.0)
	Assert(t, Equal(fv, 99.9), "value B/1/D/B retrieved")
	fv = doc.NodeAt("Z/Z/Z").AsFloat64(-1.0)
	Assert(t, Equal(fv, -1.0), "value Z/Z/Z retrieved")
}

// TestAsBool tests retrieving values as bool.
func TestAsBool(t *testing.T) {
	bs, _ := createDocument(t)

	doc, err := dynaj.Unmarshal(bs)
	Assert(t, NoError(err), "document unmarshalled")
	bv := doc.NodeAt("A").AsBool(false)
	Assert(t, Equal(bv, false), "value A retrieved")
	bv = doc.NodeAt("B/0/C").AsBool(false)
	Assert(t, Equal(bv, true), "value B/0/C retrieved")
	bv = doc.NodeAt("B/0/S/0").AsBool(false)
	Assert(t, Equal(bv, false), "value B/0/S/0 retrieved")
	bv = doc.NodeAt("B/0/S/2").AsBool(false)
	Assert(t, Equal(bv, true), "value B/0/S/2 retrieved")
	bv = doc.NodeAt("B/0/S/4").AsBool(false)
	Assert(t, Equal(bv, true), "value B/0/S/4 retrieved")
	bv = doc.NodeAt("Z/Z/Z").AsBool(false)
	Assert(t, Equal(bv, false), "value Z/Z/Z retrieved")
}

// TestMarshalJSON tests building a JSON document again.
func TestMarshalJSON(t *testing.T) {
	// Compare input and output.
	bsIn, _ := createDocument(t)
	parsedDoc, err := dynaj.Unmarshal(bsIn)
	Assert(t, NoError(err), "document unmarshalled")
	bsOut, err := parsedDoc.MarshalJSON()
	Assert(t, NoError(err), "document marshalled")
	Assert(t, DeepEqual(bsOut, bsIn), "document marshalled correctly")

	// Now doc one.
	doc := dynaj.NewDocument()
	err = doc.SetValueAt("/a/2/x", 1)
	Assert(t, NoError(err), "value /a/2/x set")
	err = doc.SetValueAt("/a/4/y", true)
	Assert(t, NoError(err), "value /a/4/y set")
	bsIn = []byte(`{"a":[null,null,{"x":1},null,{"y":true}]}`)
	bsOut, err = doc.MarshalJSON()
	Assert(t, NoError(err), "document marshalled")
	Assert(t, DeepEqual(bsOut, bsIn), "document marshalled correctly")
}

//--------------------
// HELPERS
//--------------------

type levelThree struct {
	A string
	B float64
}

type levelTwo struct {
	A string
	B int
	C bool
	D *levelThree
	S []string
}

type levelOne struct {
	A string
	B []*levelTwo
	D time.Duration
	T time.Time
}

func createDocument(t *testing.T) ([]byte, *levelOne) {
	lo := &levelOne{
		A: "Level One",
		B: []*levelTwo{
			{
				A: "Level Two - 0",
				B: 100,
				C: true,
				D: &levelThree{
					A: "Level Three - 0",
					B: 10.1,
				},
				S: []string{
					"red",
					"green",
					"1",
					"2.2",
					"true",
				},
			},
			{
				A: "Level Two - 1",
				B: 200,
				C: false,
				D: &levelThree{
					A: "Level Three - 1",
					B: 20.2,
				},
				S: []string{
					"orange",
					"blue",
					"white",
				},
			},
			{
				A: "Level Two - 2",
				B: 300,
				C: true,
				D: &levelThree{
					A: "Level Three - 2",
					B: 30.3,
				},
			},
		},
		D: 5 * time.Second,
		T: time.Date(2018, time.April, 29, 20, 30, 0, 0, time.UTC),
	}
	bs, err := json.Marshal(lo)
	Assert(t, NoError(err), "document marshalled")
	return bs, lo
}

func createCompareDocument(t *testing.T) []byte {
	lo := &levelOne{
		A: "Level One",
		B: []*levelTwo{
			{
				A: "Level Two - 0",
				B: 100,
				C: true,
				D: &levelThree{
					A: "Level Three - 0",
					B: 10.1,
				},
				S: []string{
					"red",
					"green",
					"0",
					"2.2",
					"false",
				},
			},
			{
				A: "Level Two - 1",
				B: 300,
				C: false,
				D: &levelThree{
					A: "Level Three - 1",
					B: 99.9,
				},
				S: []string{
					"orange",
					"blue",
					"white",
					"red",
				},
			},
		},
		D: 10 * time.Second,
		T: time.Date(2018, time.April, 29, 20, 59, 0, 0, time.UTC),
	}
	bs, err := json.Marshal(lo)
	Assert(t, NoError(err), "document marshalled")
	return bs
}

// EOF
