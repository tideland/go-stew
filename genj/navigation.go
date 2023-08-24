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
		return nil, fmt.Errorf("empty path")
	}
	// Walk the path.
	for _, id := range path {
		switch et := e.(type) {
		case Object:
			v, ok := et[id]
			if !ok {
				return nil, fmt.Errorf("element %q not found", id)
			}
			e = v
		case Array:
			i, err := strconv.Atoi(id)
			if err != nil {
				return nil, fmt.Errorf("element %q not found", id)
			}
			if i < 0 || i >= len(et) {
				return nil, fmt.Errorf("element %q not found", id)
			}
			e = et[i]
		default:
			return nil, fmt.Errorf("element %q not found", id)
		}
	}
	return e, nil
}

// EOF
