// Tideland Go Stew - Dynamic JSON
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dynaj // import "tideland.dev/go/stew/dynaj"

//--------------------
// TYPES
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

// EOF
