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

// newValues creates a new value.
func newValue(etc *Etc, path Path) *Accessor {
	acc := &Accessor{
		etc:  etc,
		path: path,
		acc:  etc.data.At(path...),
	}
	return acc
}

// AsString returns the value as string. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsString(def string) string {
	acc.checkMacro(def)
	s, err := acc.acc.AsString()
	if err != nil {
		return def
	}
	return s
}

// AsInt returns the value as int. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsInt(def int) int {
	acc.checkMacro(def)
	i, err := acc.acc.AsInt()
	if err != nil {
		return def
	}
	return i
}

// AsFloat64 returns the value as float. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsFloat64(def float64) float64 {
	acc.checkMacro(def)
	f, err := acc.acc.AsFloat64()
	if err != nil {
		return def
	}
	return f
}

// AsBool returns the value as bool. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsBool(def bool) bool {
	acc.checkMacro(def)
	b, err := acc.acc.AsBool()
	if err != nil {
		return def
	}
	return b
}

// AsTime returns the value as time using the given layout. In case
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

// AsDuration returns the value as duration. In case of an error
// of the accessor the default value is returned.
func (acc *Accessor) AsDuration(def time.Duration) time.Duration {
	acc.checkMacro(def)
	d, err := acc.acc.AsDuration()
	if err != nil {
		return def
	}
	return d
}

// checkMacro checks if a macro is inside the value
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
