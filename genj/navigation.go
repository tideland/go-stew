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

// walk walks the path and returns the addressed element.
func walk(e Element, path Path) (Element, error) {
	// Check the path.
	if len(path) == 0 {
		return e, nil
	}
	// Walk the path.
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
		switch et := e.(type) {
		case Object:
			v, ok := et[id]
			if !ok {
				return nil, notFoundError()
			}
			e = v
		case Array:
			i, err := strconv.Atoi(id)
			if err != nil {
				return nil, notFoundError()
			}
			if i < 0 || i >= len(et) {
				return nil, notFoundError()
			}
			e = et[i]
		default:
			return nil, notFoundError()
		}
	}
	return e, nil
}

// EOF
