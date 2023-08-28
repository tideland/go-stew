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
)

//--------------------
// DOCUMENT TREE NAVIGATION
//--------------------

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

// EOF
