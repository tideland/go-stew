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
	"fmt"
	"strconv"
	"time"
)

//--------------------
// TREE FUNCTIONS
//--------------------

// last retrieves the last ID from a Path.
func last(path Path) ID {
	if len(path) == 0 {
		return ""
	}
	return path[len(path)-1]
}

// headTail retrieves the head and the tail from a Path.
func headTail(path Path) (ID, Path) {
	switch len(path) {
	case 0:
		return "", Path{}
	case 1:
		return path[0], Path{}
	default:
		return path[0], path[1:]
	}
}

// initLast retrieves the initials and the last from a Path.
func initLast(path Path) (Path, ID) {
	switch len(path) {
	case 0:
		return Path{}, ""
	case 1:
		return Path{}, path[0]
	default:
		return path[:len(path)-1], path[len(path)-1]
	}
}

// asIndex converts the given key into an index.
func asIndex(id ID) (int, bool) {
	index, err := strconv.Atoi(id)
	if err != nil {
		return -1, false
	}
	return index, true
}

// elementAt returns the element at the given path recursively
// starting at the given start element.
func elementAt(start Element, stack, path Path) (Element, error) {
	if len(path) == 0 {
		// End of the path.
		return start, nil
	}
	// Further access depends on part content and type.
	h, t := headTail(path)
	current := append(stack, h)
	if h == "" {
		return start, nil
	}
	switch typed := start.(type) {
	case Object:
		// JSON object.
		field, ok := typed[h]
		if !ok {
			return nil, fmt.Errorf("invalid path %v", current)
		}
		return elementAt(field, current, t)
	case Array:
		// JSON array.
		index, ok := asIndex(h)
		if !ok {
			return nil, fmt.Errorf("invalid path %v: no index", current)
		}
		if index < 0 || index >= len(typed) {
			return nil, fmt.Errorf("invalid path %v: index out of range", current)
		}
		return elementAt(typed[index], current, t)
	}
	// Path is longer than existing node structure.
	return nil, fmt.Errorf("key or index not found")
}

// replaceAt replaces the element at the end of the given path with the
// given value.
func replaceAt(start Element, stack, path Path, value Element) error {
	if len(path) == 0 {
		// End of the path.
		return nil
	}
	// Further access depends on part content and type.
	h, t := headTail(path)
	if h == "" {
		return nil
	}
	current := append(stack, h)
	switch typed := start.(type) {
	case Object:
		// JSON object.
		field, ok := typed[h]
		if !ok {
			return fmt.Errorf("invalid path %v", current)
		}
		if len(t) > 0 {
			return replaceAt(field, current, t, value)
		}
		typed[h] = value
		return nil
	case Array:
		// JSON array.
		index, ok := asIndex(h)
		if !ok {
			return fmt.Errorf("invalid path %v: no index", current)
		}
		if index < 0 || index >= len(typed) {
			return fmt.Errorf("invalid path %v: index out of range", current)
		}
		if len(t) > 0 {
			return replaceAt(typed[index], current, t, value)
		}
		typed[index] = value
		return nil
	}
	// Path is longer than existing node structure.
	return fmt.Errorf("key or index not found")
}

// isValidElement checks if the element is valid. These are strings, ints, floats,
// bools, time.Time, time.Duration, and empty Array and Object.
func isValidElement(element Element) bool {
	switch typed := element.(type) {
	case string, int, float64, bool, time.Time, time.Duration:
		return true
	case Array:
		return len(typed) == 0
	case Object:
		return len(typed) == 0
	}
	return false
}

// EOF
