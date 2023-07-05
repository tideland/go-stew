// Tideland Go Stew - Callstack
//
// Copyright (C) 2017-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package callstack // import "tideland.dev/go/stew/callstack"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
)

//--------------------
// LOCATION
//--------------------

// Cached locations.
var (
	mu        sync.Mutex
	locations = make(map[uintptr]Location)
)

// Location contains the details of one location.
type Location struct {
	pkg  string
	file string
	fun  string
	line int
}

// Here return the current location.
func Here() Location {
	return At(1)
}

// At returns the location at the given offset.
func At(offset int) Location {
	mu.Lock()
	defer mu.Unlock()
	// Fix the offset.
	offset += 2
	if offset < 2 {
		offset = 2
	}
	// Retrieve program counters.
	pcs := make([]uintptr, 1)
	n := runtime.Callers(offset, pcs)
	if n == 0 {
		return Location{}
	}
	pcs = pcs[:n]
	// Check cache.
	pc := pcs[0]
	l, ok := locations[pc]
	if ok {
		return l
	}
	// Build ID based on program counters.
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		pkg, fun := path.Split(frame.Function)
		parts := strings.Split(fun, ".")
		pkg = path.Join(pkg, parts[0])
		fun = strings.Join(parts[1:], ".")
		_, file := path.Split(frame.File)
		if !more {
			l := Location{
				pkg:  pkg,
				file: file,
				fun:  fun,
				line: frame.Line,
			}
			locations[pc] = l
			return l
		}
	}
}

// Package returns the package of the location.
func (l Location) Package() string {
	return l.pkg
}

// File returns the file of the location.
func (l Location) File() string {
	return l.file
}

// Func returns the function of the location.
func (l Location) Func() string {
	return l.fun
}

// Line returns the line of the location.
func (l Location) Line() int {
	return l.line
}

// String implements fmt.Stringer.
func (l Location) String() string {
	return fmt.Sprintf("(%s:%s:%s:%d)", l.pkg, l.file, l.fun, l.line)
}

//--------------------
// STACK
//--------------------

// Stack contains a number of locations of a call stack.
type Stack []Location

// String returns a string representation of the stack.
func (s Stack) String() string {
	var ids []string
	for _, l := range s {
		ids = append(ids, l.String())
	}
	return strings.Join(ids, " :: ")
}

// Dive returns the current callstack until the given depth.
func Dive(depth int) Stack {
	var stack Stack
	var end = depth + 1
	if end < 1 {
		end = 1
	}
	for i := 1; i < end; i++ {
		stack = append(stack, At(i))
	}
	return stack
}

// EOF
