// Tideland Go Stew - Generic JSON
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
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

//--------------------
// JSON TYPES
//--------------------

// ID represents a key or index in a JSON object or array.
type ID = string

// Path represents a list of IDs.
type Path = []ID

// Element represents a JSON element, i.e. a simple value, an object or an array.
type Element = any

// Value represents a simple JSON value.
type Value = any

// Object represents a JSON object.
type Object = map[string]any

// Array represents a JSON array.
type Array = []any

// ValueConstraint descibes allowed types for JSON values.
type ValueConstraint interface {
	string | int | float64 | bool | time.Time | time.Duration
}

//--------------------
// JSON DOCUMENT
//--------------------

// Document represents a JSON document. The structure is hidden to
// avoid direct access to the elements.
type Document struct {
	root Element
}

// New creates a new empty document.
func New() *Document {
	return &Document{}
}

// Read reads a document from a reader.
func Read(r io.Reader) (*Document, error) {
	// Read the data.
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read document: %v", err)
	}
	// Unmarshal the data to JSON.
	var root Element
	err = json.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal document: %v", err)
	}
	return &Document{
		root: root,
	}, nil
}

// Write writes a document to a writer.
func Write(d *Document, w io.Writer) error {
	// Marshal the data to JSON.
	data, err := json.Marshal(d.root)
	if err != nil {
		return fmt.Errorf("cannot marshal document: %v", err)
	}
	// Write the data.
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write document: %v", err)
	}
	return nil
}

//--------------------
// FUNCTIONS WORKING ON DOCUMENTS
//--------------------

// Get returns the addressed element or an error if the path is invalid. In case of a
// found Value first the matching type will be checked, otherwise it will be tried to
// convert it if matches ValueConstraint. If this is not possible an error will be
// returned.
func Get[V ValueConstraint](d *Document, path ...ID) (V, error) {
	var v V
	// Walk the path.
	h, _, err := walk(d.root, path)
	if err != nil {
		return v, fmt.Errorf("cannot get element: %v", err)
	}
	// Check if wanted type.
	switch ht := h.(type) {
	case V:
		return ht, nil
	case Object, Array:
		return v, fmt.Errorf("path points to object or array")
	}
	// No match, try to convert.
	nv, ok := elementToValue(h, v)
	if !ok {
		return v, fmt.Errorf("element is not of type %T", v)
	}
	ev, ok := nv.(V)
	if !ok {
		return v, fmt.Errorf("element is not of type %T", v)
	}
	return ev, nil
}

// GetDefault returns the addressed element or the default value if the path is invalid
// or a conversion fails. Only in case of the path pointing to an Object or Array the
// given ValueConstraints standard value will be returned.
func GetDefault[V ValueConstraint](d *Document, def V, path ...ID) V {
	var v V
	// Walk the path.
	h, _, err := walk(d.root, path)
	if err != nil {
		return def
	}
	// Check if wanted type.
	switch ht := h.(type) {
	case V:
		return ht
	case Object, Array:
		return v
	}
	// Try to convert or take default.
	nv, ok := elementToValue(h, def)
	if !ok {
		return def
	}
	ev, ok := nv.(V)
	if !ok {
		return def
	}
	return ev
}

// Set sets the addressed element to the given value. The new value has to match the
// type of the current value. The path has to be valid.
func Set[V ValueConstraint](d *Document, v V, path ...ID) error {
	// Walk the path.
	h, t, err := walk(d.root, path)
	if err != nil {
		return fmt.Errorf("cannot set element: %v", err)
	}
	// Check if type of head is correct.
	nv, ok := valueToElement(h, v)
	if !ok {
		return fmt.Errorf("current element of type %T does not match to new type %T: %v", h, v, nv)
	}
	// Last of tail will be the parent. Set the new value.
	switch pt := t[len(t)-1].(type) {
	case Object:
		pt[path[len(path)-1]] = nv
	case Array:
		i, err := strconv.Atoi(path[len(path)-1])
		if err != nil {
			return fmt.Errorf("cannot set element: %v", err)
		}
		pt[i] = nv
	}
	return nil
}

// SetAny sets the addressed element to the given value. The path has to be valid. The type of
// the current value doesn't matter.
func SetAny(d *Document, v any, path ...ID) error {
	// Walk the path.
	h, t, err := walk(d.root, path)
	if err != nil {
		return fmt.Errorf("cannot set element: %v", err)
	}
	// Head is not allowed to be an object or array.
	switch h.(type) {
	case Object, Array:
		return fmt.Errorf("current element is not allowed to be an object or array")
	}
	// Last of tail will be the parent. Set the new value.
	switch pt := t[len(t)-1].(type) {
	case Object:
		pt[path[len(path)-1]] = v
	case Array:
		i, err := strconv.Atoi(path[len(path)-1])
		if err != nil {
			return fmt.Errorf("cannot set element: %v", err)
		}
		pt[i] = v
	}
	return nil
}

// EOF
