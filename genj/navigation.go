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
	"fmt"
	"strconv"

	"tideland.dev/go/stew/slices"
)

//--------------------
// DOCUMENT TREE NAVIGATION
//--------------------

// moonwalk walks down the path starting at the given element and
// returns the last ID of the whole path as well as the last element found.
// Next value is the rest of the path if it is longer than one element. In
// case of a non-integer path ID for an array the error is returned.
func moonwalk(start Element, path Path) (ID, Element, Path, error) {
	var id ID
	// Check basic path length.
	switch len(path) {
	case 0:
		return "", start, nil, fmt.Errorf("path is empty")
	case 1:
		path, id = slices.InitLast(path)
		return id, start, path, nil
	}
	// Start walking.
	id = slices.Last(path)
	current := start
	for {
		// Get head and tail of path.
		head, tail := slices.HeadTail(path)
		// Walk the path.
		switch ct := current.(type) {
		case Object:
			next, ok := ct[head]
			if !ok {
				// End of walk but not of path.
				return id, ct, path, nil
			}
			current = next
		case Array:
			i, err := strconv.Atoi(head)
			if err != nil {
				return id, current, path, fmt.Errorf("invalid array index %q", head)
			}
			if i < 0 {
				return id, current, path, fmt.Errorf("array index %d out of bounds", i)
			}
			if i >= len(ct) {
				// End of walk but not of path.
				return id, ct, path, nil
			}
			current = ct[i]
		default:
			return id, nil, path, fmt.Errorf("path %q is illegal", path)
		}
		// Check if we are done.
		if len(tail) == 0 {
			return id, current, nil, nil
		}
		// Continue with the tail.
		path = tail
	}
}

// contains takes an Element and an ID and checks if the element
// is an Object or Array and contains the ID. In this case the found
// Element and nil are returned, otherwise nil and an error.
func contains(elem Element, id ID) (Element, error) {
	switch et := elem.(type) {
	case Object:
		value, ok := et[id]
		if !ok {
			return nil, fmt.Errorf("element %q not found", id)
		}
		return value, nil
	case Array:
		i, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		if i < 0 || i >= len(et) {
			return nil, fmt.Errorf("index %d out of bounds", i)
		}
		value := et[i]
		return value, nil
	}
	return nil, fmt.Errorf("element is no Object or Array")
}

// walk walks down the path starting at the given element and
// returns the found element as head and all elements on the path
// as tail.
func walk(start Element, path Path) (Element, []Element, error) {
	// Check the path.
	if len(path) == 0 {
		return start, nil, nil
	}
	// Walk the path.
	head := start
	tail := []Element{}
	sofar := []string{}
	notFoundError := func() error {
		msg := "path ["
		for i, id := range sofar {
			if i > 0 {
				msg += ", "
			}
			msg += fmt.Sprintf("%q", id)
		}
		msg += "] not found"
		return fmt.Errorf(msg)
	}
	for _, id := range path {
		sofar = append(sofar, id)
		switch et := head.(type) {
		case Object:
			v, ok := et[id]
			if !ok {
				return nil, tail, notFoundError()
			}
			tail = append(tail, head)
			head = v
		case Array:
			i, err := strconv.Atoi(id)
			if err != nil {
				return nil, tail, notFoundError()
			}
			if i < 0 || i >= len(et) {
				return nil, tail, notFoundError()
			}
			tail = append(tail, head)
			head = et[i]
		default:
			return nil, tail, notFoundError()
		}
	}
	return head, tail, nil
}

// create walks down the path starting at the given element and ends with
// the last one or creates the missing elements. It return the last ID and
// the last element.
func create(start Element, path Path) (ID, Element, error) {
	// Check the path.
	switch len(path) {
	case 0:
		return "", start, fmt.Errorf("path is empty")
	case 1:
		return path[0], start, nil
	}
	// Fetch end of path as ID and dive into the tree until the
	// path ended or the document tree ended.
	begin, _ := slices.InitLast(path)
	ph, pt := slices.HeadTail(begin)
	head := start
	for {
		current := head
		switch et := current.(type) {
		case Object:
			v, ok := et[ph]
			if !ok {
				break
			}
			head = v
		case Array:
			i, err := strconv.Atoi(ph)
			if err != nil {
				return "", nil, err
			}
			if i < 0 {
				return "", nil, fmt.Errorf("negative array index %d", i)
			}
			if i >= len(et) {
				// Enlarge array and break.
				for j := len(et); j <= i; j++ {
					et = append(et, nil)
				}
				head = et
				break
			}
		}
		ph, pt = slices.HeadTail(pt)
	}
	// Now create the missing elements.

}

// EOF
