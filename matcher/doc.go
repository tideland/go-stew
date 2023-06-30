// Tideland Go Stew - Matcher
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package matcher provides the matching of simple patterns against strings.
// It matches the following pattterns:
//
// - ? matches one char
// - * matches a group of chars
// - [abc] matches any of the chars inside the brackets
// - [a-z] matches any of the chars of the range
// - [^abc] matches any but the chars inside the brackets
// - \ escapes any of the pattern chars
package matcher // import "tideland.dev/go/stew/matcher"

// EOF
