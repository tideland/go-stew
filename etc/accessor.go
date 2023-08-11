// Tideland Go Stew - Etc
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package etc // import "tideland.dev/go/stew/etc"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"os"
	"strings"
	"time"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// Value
//--------------------

// Accessor provides access to a configuration value.
type Accessor struct {
	etc  *Etc
	path Path
	acc  *dynaj.Accessor
}

// newAccessor creates a new Accessor for a value.
func newAccessor(etc *Etc, path Path) *Accessor {
	acc := &Accessor{
		etc:  etc,
		path: path,
		acc:  etc.data.At(path...),
	}
	return acc
}

// Path returns the path of the accessor.
func (acc *Accessor) Path() Path {
	return acc.acc.Path()
}

// ID returns the ID of the Accessor.
func (acc *Accessor) ID() string {
	return acc.acc.ID()
}

// AsString returns the value of the Accessor as string. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsString(def string) string {
	acc.checkMacro(def)
	s, err := acc.acc.AsString()
	if err != nil {
		return def
	}
	return s
}

// AsInt returns the value of the Accessor as int. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsInt(def int) int {
	acc.checkMacro(def)
	i, err := acc.acc.AsInt()
	if err != nil {
		return def
	}
	return i
}

// Err returns the error of the Accessor.
func (acc *Accessor) Err() error {
	return acc.acc.Err()
}

// Len returns the length of the value of the Accessor.
func (acc *Accessor) Len() int {
	return acc.acc.Len()
}

// AsFloat64 returns the value of the Accessor as float. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsFloat64(def float64) float64 {
	acc.checkMacro(def)
	f, err := acc.acc.AsFloat64()
	if err != nil {
		return def
	}
	return f
}

// AsBool returns the value of the Accessor as bool. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsBool(def bool) bool {
	acc.checkMacro(def)
	b, err := acc.acc.AsBool()
	if err != nil {
		return def
	}
	return b
}

// AsTime returns the value of the Accessor as time using the given layout. In case
// of an empty string time.RFC3339 will be taken. An error during reading
// or parsing will return the default value.
func (acc *Accessor) AsTime(layout string, def time.Time) time.Time {
	acc.checkMacro(def)
	if layout == "" {
		layout = time.RFC3339
	}
	t, err := acc.acc.AsTime(layout)
	if err != nil {
		return def
	}
	return t
}

// AsDuration returns the value of the Accessor as duration. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsDuration(def time.Duration) time.Duration {
	acc.checkMacro(def)
	d, err := acc.acc.AsDuration()
	if err != nil {
		return def
	}
	return d
}

// Update updates the configuration value.
func (acc *Accessor) Update(value Value) *Accessor {
	acc.acc.Update(value)
	return newAccessor(acc.etc, acc.path)
}

// Set sets a value at a given key in kase of an object or
// in case of an Array if the key contains an integer
// in range of the array.
func (acc *Accessor) Set(key string, value Value) *Accessor {
	acc.acc.Set(key, value)
	return newAccessor(acc.etc, acc.path)
}

// Append appends a value to a configuration array.
func (acc *Accessor) Append(value Value) *Accessor {
	acc.acc.Append(value)
	return newAccessor(acc.etc, acc.path)
}

// Delete deletes a value at a given key in kase of an object or
// in case of an Array if the key contains an integer
// in range of the array.
func (acc *Accessor) Delete() *Accessor {
	newAcc := acc.acc.Delete()
	return newAccessor(acc.etc, newAcc.Path())
}

// At returns a new Accessor for a sub value.
func (acc *Accessor) At(path ...ID) *Accessor {
	return newAccessor(acc.etc, append(acc.path, path...))
}

// checkMacro checks if a macro is inside the value of the Accessor
// and replaces it.
func (acc *Accessor) checkMacro(def any) {
	s, err := acc.acc.AsString()
	if err != nil {
		return
	}
	sidx := strings.Index(s, "[[")
	if sidx == -1 {
		return
	}
	eidx := strings.Index(s[sidx:], "]]")
	if eidx == -1 {
		return
	}
	// Macro found, now look for default value..
	prefix := s[:sidx]
	suffix := s[sidx+eidx+2:]
	macro := s[sidx+2 : sidx+eidx]
	macroDef := ""
	value := ""
	didx := strings.Index(macro, "||")
	if didx != -1 {
		macroDef = macro[didx+2:]
		macro = macro[:didx]
	}
	fmt.Printf("macro: %q, macroDef: %q\n", macro, macroDef)
	// Check if macro is an environment variable of a path.
	if strings.HasPrefix(macro, "$") {
		// Environment variable.
		value = os.Getenv(macro[1:])
	} else {
		// Path.
		value = acc.etc.At(strings.Split(macro, "::")...).AsString(macroDef)
	}
	fmt.Printf("value: %q\n", value)
	// Check if value is empty.
	if value == "" {
		value = macroDef
	}
	if value == "" {
		value = fmt.Sprintf("%v", def)
	}
	fmt.Printf("value: %q\n", value)
	// Replace macro.
	acc.acc.Update(prefix + value + suffix)
}

// EOF
