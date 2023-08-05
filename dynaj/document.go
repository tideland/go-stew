// Tideland Go Stew - Dynamic JSON
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dynaj // import "tideland.dev/go/stew/dynaj"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"fmt"
	"io"
)

//--------------------
// DOCUMENT
//--------------------

// Document represents one JSON document.
type Document struct {
	root Element
}

// NewDocument creates a new empty document.
func NewDocument() *Document {
	return &Document{}
}

// Read reads a JSON document from the reader.
func Read(r io.Reader) (*Document, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read document: %v", err)
	}
	return Unmarshal(data)
}

// Unmarshal parses a JSON-encoded document.
func Unmarshal(data []byte) (*Document, error) {
	var root any
	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal document: %v", err)
	}
	return &Document{
		root: root,
	}, nil
}

// At returns the addressed Accessor containing the addressed
// element or an error if the path is invalid.
func (doc *Document) At(path ...ID) *Accessor {
	elem, err := elementAt(doc.root, Path{}, path)
	acc := newAccessor(doc, path, elem, err)
	return acc
}

// Root returns the root Accesor of the document.
func (d *Document) Root() *Accessor {
	return newAccessor(d, Path{}, d.root, nil)
}

// Clone returns a clone of the document.
func (d *Document) Clone() (*Document, error) {
	var raw []byte
	raw, err := json.Marshal(d.root)
	if err != nil {
		return nil, fmt.Errorf("cannot clone document: %v", err)
	}
	dc := &Document{}
	err = json.Unmarshal(raw, &dc.root)
	if err != nil {
		return nil, fmt.Errorf("cannot clone document: %v", err)
	}
	return dc, nil
}

// MarshalJSON returns the JSON encoding of the document.
func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.root)
}

// MarshalJSONIndent returns the indented JSON encoding of the document.
func (d *Document) MarshalJSONIndent() ([]byte, error) {
	return json.MarshalIndent(d.root, "", "  ")
}

// String returns the string representation of the document.
func (d *Document) String() string {
	data, err := d.MarshalJSONIndent()
	if err != nil {
		return fmt.Sprintf("cannot marshal document: %v", err)
	}
	return string(data)
}

// EOF
