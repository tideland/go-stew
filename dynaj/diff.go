// Tideland Go Stew - Dynamic JSON
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dynaj // import "tideland.dev/go/stew/dynaj"

//--------------------
// DIFFERENCE
//--------------------

// Diff manages the two parsed documents and their differences.
type Diff struct {
	first  *Document
	second *Document
	paths  []string
}

// Compare parses and compares the documents and returns their differences.
func Compare(first, second []byte) (*Diff, error) {
	fd, err := Unmarshal(first)
	if err != nil {
		return nil, err
	}
	sd, err := Unmarshal(second)
	if err != nil {
		return nil, err
	}
	d := &Diff{
		first:  fd,
		second: sd,
	}
	err = d.compare()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CompareDocuments compares the documents and returns their differences.
func CompareDocuments(first, second *Document) (*Diff, error) {
	d := &Diff{
		first:  first,
		second: second,
	}
	err := d.compare()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// FirstDocument returns the first document passed to Diff().
func (d *Diff) FirstDocument() *Document {
	return d.first
}

// SecondDocument returns the second document passed to Diff().
func (d *Diff) SecondDocument() *Document {
	return d.second
}

// Differences returns a list of paths where the documents
// have different content.
func (d *Diff) Differences() []string {
	return d.paths
}

// DifferenceAt returns the differences at the given path by
// returning the first and the second value.
func (d *Diff) DifferenceAt(path string) (*Node, *Node) {
	fstNode := d.first.NodeAt(path)
	sndNode := d.second.NodeAt(path)
	return fstNode, sndNode
}

// compare iterates over the both documents looking for different
// values or even paths.
func (d *Diff) compare() error {
	firstPaths := map[string]struct{}{}
	firstProcessor := func(node *Node) error {
		firstPaths[node.path] = struct{}{}
		if !node.Equals(d.second.NodeAt(node.path)) {
			d.paths = append(d.paths, node.path)
		}
		return nil
	}
	err := d.first.Root().Process(firstProcessor)
	if err != nil {
		return err
	}
	secondProcessor := func(node *Node) error {
		_, ok := firstPaths[node.path]
		if ok {
			// Been there, done that.
			return nil
		}
		d.paths = append(d.paths, node.path)
		return nil
	}
	return d.second.Root().Process(secondProcessor)
}

// EOF
