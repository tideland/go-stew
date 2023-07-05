// Tideland Go Stew - Callstack
//
// Copyright (C) 2017-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package callstack provides a way to retrieve informations about the current position
// in the call stack. This can be used to create more detailed error messages or
// to log the call stack. Offsets can be used to retrieve informations about the
// caller of the caller.
//
//	here := callstack.Here()
//	downunder := callstack.At(5)
//	deeperID := callstack.At(2).ID
//	deeperCode := callstack.At(2).Code("ERR")
//	stack := callstack.Dive(5)
//
// Internal caching fastens retrieval after first call.
package callstack // import "tideland.dev/go/stew/callstack"

// EOF
