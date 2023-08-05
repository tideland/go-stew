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
// ACCESSOR
//--------------------

// Accessor provides access to a JSON document.
type Accessor struct {
	doc     *Document
	path    Path
	element Element
	err     error
}

// newAccessor quickly creates a new Accessor.
func newAccessor(doc *Document, path Path, element Element, err error) *Accessor {
	return &Accessor{
		doc:     doc,
		path:    path,
		element: element,
		err:     err,
	}
}

// newError creates a new Accessor with an error.
func newError(acc *Accessor, format string, a ...any) *Accessor {
	return newAccessor(acc.doc, acc.path, nil, fmt.Errorf(format, a...))
}

// Path returns the path of the Accessor.
func (acc *Accessor) Path() Path {
	path := make(Path, len(acc.path))
	copy(path, acc.path)
	return path
}

// ID returns the ID of the Accessor. It is the last element
// of the path.
func (acc *Accessor) ID() ID {
	return last(acc.path)
}

// IsError returns true if the Accessor has an error.
func (acc *Accessor) IsError() bool {
	return acc.err != nil
}

// Err returns a possible error of the Accessor.
func (acc *Accessor) Err() error {
	return acc.err
}

// Len returns the length of the element. For Arrays and Objects
// this is the number of elements, for all others it is 1. If the
// Accessor has an error it returns 0.
func (acc *Accessor) Len() int {
	if acc.element == nil || acc.err != nil {
		return 0
	}
	switch typed := acc.element.(type) {
	case Array:
		return len(typed)
	case Object:
		return len(typed)
	}
	return 1
}

// AsString returns the value as string.
func (acc *Accessor) AsString() (string, error) {
	if acc.err != nil {
		return "", acc.err
	}
	switch telem := acc.element.(type) {
	case string:
		return telem, nil
	case int:
		return strconv.Itoa(telem), nil
	case float64:
		return strconv.FormatFloat(telem, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(telem), nil
	case time.Time:
		return telem.Format(time.RFC3339Nano), nil
	case time.Duration:
		return telem.String(), nil
	case Array:
		return "[...]", nil
	case Object:
		return "{...}", nil
	}
	return "", fmt.Errorf("cannot convert %T to string", acc.element)
}

// AsInt returns the value as int.
func (acc *Accessor) AsInt() (int, error) {
	if acc.err != nil {
		return 0, acc.err
	}
	switch telem := acc.element.(type) {
	case string:
		i, ok := asIndex(telem)
		if !ok {
			return 0, fmt.Errorf("cannot convert %q to int", telem)
		}
		return i, nil
	case int:
		return telem, nil
	case float64:
		return int(telem), nil
	case bool:
		if telem {
			return 1, nil
		}
		return 0, nil
	case Array:
		return len(telem), nil
	case Object:
		return len(telem), nil
	}
	return 0, fmt.Errorf("cannot convert %T to int", acc.element)
}

// AsFloat64 returns the value as float64.
func (acc *Accessor) AsFloat64() (float64, error) {
	if acc.err != nil {
		return 0.0, acc.err
	}
	switch telem := acc.element.(type) {
	case string:
		f, err := strconv.ParseFloat(telem, 64)
		if err != nil {
			return 0.0, fmt.Errorf("cannot convert %q to float64: %v", telem, err)
		}
		return f, nil
	case int:
		return float64(telem), nil
	case float64:
		return telem, nil
	case bool:
		if telem {
			return 1.0, nil
		}
		return 0.0, nil
	case Array:
		return float64(len(telem)), nil
	case Object:
		return float64(len(telem)), nil
	}
	return 0.0, fmt.Errorf("cannot convert %T to float64", acc.element)
}

// AsBool returns the value as bool.
func (acc *Accessor) AsBool() (bool, error) {
	if acc.err != nil {
		return false, acc.err
	}
	switch telem := acc.element.(type) {
	case string:
		b, err := strconv.ParseBool(telem)
		if err != nil {
			return false, fmt.Errorf("cannot convert %q to bool: %v", telem, err)
		}
		return b, nil
	case int:
		return telem != 0, nil
	case float64:
		return telem != 0.0, nil
	case bool:
		return telem, nil
	}
	return false, fmt.Errorf("cannot convert %T to bool", acc.element)
}

// AsTime returns the value as time.Time.
func (acc *Accessor) AsTime(format string) (time.Time, error) {
	if acc.err != nil {
		return time.Time{}, acc.err
	}
	switch telem := acc.element.(type) {
	case time.Time:
		return telem, nil
	case string:
		t, err := time.Parse(format, telem)
		if err != nil {
			return time.Time{}, fmt.Errorf("cannot convert %q to time.Time: %v", telem, err)
		}
		return t, nil
	case int:
		return time.Unix(int64(telem), 0), nil
	case float64:
		return time.Unix(int64(telem), 0), nil
	}
	return time.Time{}, fmt.Errorf("cannot convert %T to time.Time", acc.element)
}

// AsDuration returns the value as time.Duration.
func (acc *Accessor) AsDuration() (time.Duration, error) {
	if acc.err != nil {
		return 0, acc.err
	}
	switch telem := acc.element.(type) {
	case time.Duration:
		return telem, nil
	case string:
		d, err := time.ParseDuration(telem)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %q to time.Duration: %v", telem, err)
		}
		return d, nil
	case int:
		return time.Duration(telem) * time.Second, nil
	case float64:
		return time.Duration(telem) * time.Second, nil
	}
	return 0, fmt.Errorf("cannot convert %T to time.Duration", acc.element)
}

// Update changes the value of the node. It is allowed for strings,
// ints, floats, bools, time.Time, time.Duration, empty Objects or
// Arrays.
func (acc *Accessor) Update(value Value) *Accessor {
	if acc.err != nil {
		return acc
	}
	if !isValidElement(value) {
		return newError(acc, "invalid type or non-empty container type for update")
	}
	if len(acc.path) == 0 {
		acc.doc.root = value
		acc.element = value
		return acc
	}
	if err := replaceAt(acc.doc.root, Path{}, acc.path, value); err != nil {
		return newError(acc, "cannot update element: %v", err)
	}
	acc.element = value
	return acc
}

// Set sets the value at a given key. It only works in case of
// an Object or in case of an Array if the key contains an integer
// in range of the array. The values are only allowed for strings,
// ints, floats, bools, time.Time, time.Duration, empty Objects or
// Arrays.
func (acc *Accessor) Set(key ID, value Value) *Accessor {
	if acc.err != nil {
		return acc
	}
	if !isValidElement(value) {
		return newError(acc, "invalid type or non-empty container type for set")
	}
	// Check element type.
	switch typed := acc.element.(type) {
	case Array:
		idx, ok := asIndex(key)
		if !ok || (idx < 0 || idx >= len(typed)) {
			return newError(acc, "cannot set element: illegal index")
		}
		typed[idx] = value
		path := append(acc.path, key)
		return newAccessor(acc.doc, path, value, nil)
	case Object:
		typed[key] = value
		path := append(acc.path, key)
		return newAccessor(acc.doc, path, value, nil)
	}
	return newError(acc, "cannot set element: not an array or object")
}

// Append appends a value to an array. The value is only allowed
// for strings, ints, floats, bools, time.Time, time.Duration,
// empty Objects or Arrays.
func (acc *Accessor) Append(value Value) *Accessor {
	if acc.err != nil {
		return acc
	}
	if !isValidElement(value) {
		return newError(acc, "invalid type or non-empty container type for append")
	}
	// Check element type.
	switch typed := acc.element.(type) {
	case Array:
		typed = append(typed, value)
		path := append(acc.path, strconv.Itoa(len(typed)-1))
		if err := replaceAt(acc.doc.root, Path{}, acc.path, typed); err != nil {
			return newError(acc, "cannot append element: %v", err)
		}
		return newAccessor(acc.doc, path, value, nil)
	}
	return newError(acc, "cannot append element: not an array")
}

// Delete deletes the element at the location and returns the
// parent element.
func (acc *Accessor) Delete() *Accessor {
	if acc.err != nil {
		return acc
	}
	// Start delete with check for root.
	if len(acc.path) == 0 {
		acc.doc.root = nil
		acc.element = nil
		acc.err = nil
		return acc
	}
	// No, so find parent and delete element.
	i, l := initLast(acc.path)
	parent, err := elementAt(acc.doc.root, Path{}, i)
	if err != nil {
		return newError(acc, "cannot delete element: %v", err)
	}
	switch typed := parent.(type) {
	case Array:
		idx, ok := asIndex(l)
		if !ok || (idx < 0 || idx >= len(typed)) {
			return newError(acc, "invalid index for delete")
		}
		typed = append(typed[:idx], typed[idx+1:]...)
		if err := replaceAt(acc.doc.root, Path{}, i, typed); err != nil {
			return newError(acc, "cannot delete element: %v", err)
		}
		return newAccessor(acc.doc, i, typed, nil)
	case Object:
		delete(typed, l)
		return newAccessor(acc.doc, i, typed, nil)
	}
	return newError(acc, "cannot delete element: parent not an array or object")
}

// At returns the addressed Accessor containing the element
// or an error if the path is invalid.
func (acc *Accessor) At(path ...ID) *Accessor {
	if acc.err != nil {
		return newAccessor(acc.doc, acc.path, nil, acc.err)
	}
	extendedPath := append(acc.path, path...)
	elem, err := elementAt(acc.doc.root, Path{}, extendedPath)
	if err != nil {
		return newAccessor(acc.doc, extendedPath, nil, err)
	}
	return newAccessor(acc.doc, extendedPath, elem, nil)
}

// Processor returns a processor starting it's work at the Accessors location.
func (acc *Accessor) Processor() *Processor {
	return newProcessor(acc)
}

// EOF
