// Tideland Go Stew
//
// Copyright (C) 2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stew // import "tideland.dev/go/stew"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/stew/semver"
)

//--------------------
// VERSION
//--------------------

func Version() *semver.Version {
	return semver.NewVersion(0, 1, 0, "alpha", "2023-08-11")
}

// EOF
