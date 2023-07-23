// Tideland Go Stew - Semantic Versions
//
// Copyright (C) 2014-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package semver // import "tideland.dev/go/stew/semver"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strconv"
	"strings"
)

//--------------------
// CONST
//--------------------

// Precedence describes if a version is newer, equal, or older.
type Precedence int

// Level describes the level, on which a version differentiates
// from an other.
type Level string

// Separator, precedences, and part identifiers.
const (
	Metadata = "+"

	Newer Precedence = 1
	Equal Precedence = 0
	Older Precedence = -1

	Major      Level = "major"
	Minor      Level = "minor"
	Patch      Level = "patch"
	PreRelease Level = "pre-release"
	All        Level = "all"
)

//--------------------
// VERSION
//--------------------

// Version implements a semantic version.
type Version struct {
	major      int
	minor      int
	patch      int
	preRelease []string
	metadata   []string
}

// NewVersion returns a simple version instance. Parts of pre-release
// and metadata are passed as optional strings separated by
// version.Metadata ("+").
func NewVersion(major, minor, patch int, prmds ...string) *Version {
	if major < 0 {
		major = 0
	}
	if minor < 0 {
		minor = 0
	}
	if patch < 0 {
		patch = 0
	}
	vsn := &Version{
		major: major,
		minor: minor,
		patch: patch,
	}
	isPR := true
	for _, prmd := range prmds {
		if isPR {
			if prmd == Metadata {
				isPR = false
				continue
			}
			vsn.preRelease = append(vsn.preRelease, validID(prmd, true))
		} else {
			vsn.metadata = append(vsn.metadata, validID(prmd, false))
		}
	}
	return vsn
}

// Parse retrieves a version out of a string.
func Parse(vsnstr string) (*Version, error) {
	// Split version, pre-release, and metadata.
	npmstrs, err := splitVersionString(vsnstr)
	if err != nil {
		return nil, err
	}
	// Parse these parts.
	nums, err := parseNumberString(npmstrs[0])
	if err != nil {
		return nil, err
	}
	prmds := []string{}
	if npmstrs[1] != "" {
		prmds = strings.Split(npmstrs[1], ".")
	}
	if npmstrs[2] != "" {
		prmds = append(prmds, Metadata)
		prmds = append(prmds, strings.Split(npmstrs[2], ".")...)
	}
	// Done.
	return NewVersion(nums[0], nums[1], nums[2], prmds...), nil
}

// Major returns the major version number.
func (vsn *Version) Major() int {
	return vsn.major
}

// Minor returns the minor version number.
func (vsn *Version) Minor() int {
	return vsn.minor
}

// Patch returns the patch version number.
func (vsn *Version) Patch() int {
	return vsn.patch
}

// PreRelease returns the pre-release string.
func (vsn *Version) PreRelease() string {
	return strings.Join(vsn.preRelease, ".")
}

// Metadata returns the metadata string.
func (vsn *Version) Metadata() string {
	return strings.Join(vsn.metadata, ".")
}

// Compare implements the Version interface.
func (vsn *Version) Compare(cvsn *Version) (Precedence, Level) {
	// Standard version parts.
	switch {
	case vsn.major < cvsn.Major():
		return Older, Major
	case vsn.major > cvsn.Major():
		return Newer, Major
	case vsn.minor < cvsn.Minor():
		return Older, Minor
	case vsn.minor > cvsn.Minor():
		return Newer, Minor
	case vsn.patch < cvsn.Patch():
		return Older, Patch
	case vsn.patch > cvsn.Patch():
		return Newer, Patch
	}
	// Now the parts of the pre-release.
	cvsnpr := []string{}
	for _, cvsnprPart := range strings.Split(cvsn.PreRelease(), ".") {
		if cvsnprPart != "" {
			cvsnpr = append(cvsnpr, cvsnprPart)
		}
	}
	vsnlen := len(vsn.preRelease)
	cvsnlen := len(cvsnpr)
	count := vsnlen
	if cvsnlen < vsnlen {
		count = cvsnlen
	}
	for i := 0; i < count; i++ {
		vsnn, vsnerr := strconv.Atoi(vsn.preRelease[i])
		cvsnn, cvsnerr := strconv.Atoi(cvsnpr[i])
		if vsnerr == nil && cvsnerr == nil {
			// Numerical comparison.
			switch {
			case vsnn < cvsnn:
				return Older, PreRelease
			case vsnn > cvsnn:
				return Newer, PreRelease
			}
			continue
		}
		// Alphanumerical comparison.
		switch {
		case vsn.preRelease[i] < cvsnpr[i]:
			return Older, PreRelease
		case vsn.preRelease[i] > cvsnpr[i]:
			return Newer, PreRelease
		}
	}
	// Still no clean result, so the shorter
	// pre-relese is older.
	switch {
	case vsnlen < cvsnlen:
		return Newer, PreRelease
	case vsnlen > cvsnlen:
		return Older, PreRelease
	}
	// Last but not least: we are equal.
	return Equal, All
}

// Less checks if the version is older than the passed one.
func (vsn *Version) Less(cvsn *Version) bool {
	precedence, _ := vsn.Compare(cvsn)
	return precedence == Older
}

// String implements the fmt.Stringer interface.
func (vsn *Version) String() string {
	vsns := fmt.Sprintf("%d.%d.%d", vsn.major, vsn.minor, vsn.patch)
	if len(vsn.preRelease) > 0 {
		vsns += "-" + vsn.PreRelease()
	}
	if len(vsn.metadata) > 0 {
		vsns += Metadata + vsn.Metadata()
	}
	return vsns
}

//--------------------
// TOOLS
//--------------------

// validID reduces the passed identifier to a valid one. If we care
// for numeric identifiers leading zeros will be removed.
func validID(id string, numeric bool) string {
	out := []rune{}
	letter := false
	digit := false
	hyphen := false
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z':
			letter = true
			out = append(out, r)
		case r >= 'A' && r <= 'Z':
			letter = true
			out = append(out, r)
		case r >= '0' && r <= '9':
			digit = true
			out = append(out, r)
		case r == '-':
			hyphen = true
			out = append(out, r)
		}
	}
	if numeric && digit && !letter && !hyphen {
		// Digits only, and we care for it.
		// Remove leading zeros.
		for len(out) > 0 && out[0] == '0' {
			out = out[1:]
		}
		if len(out) == 0 {
			out = []rune{'0'}
		}
	}
	return string(out)
}

// splitVersionString separates the version string into numbers,
// pre-release, and metadata strings.
func splitVersionString(vsnstr string) ([]string, error) {
	npXm := strings.SplitN(vsnstr, Metadata, 2)
	switch len(npXm) {
	case 1:
		nXp := strings.SplitN(npXm[0], "-", 2)
		switch len(nXp) {
		case 1:
			return []string{nXp[0], "", ""}, nil
		case 2:
			return []string{nXp[0], nXp[1], ""}, nil
		}
	case 2:
		nXp := strings.SplitN(npXm[0], "-", 2)
		switch len(nXp) {
		case 1:
			return []string{nXp[0], "", npXm[1]}, nil
		case 2:
			return []string{nXp[0], nXp[1], npXm[1]}, nil
		}
	}
	return nil, fmt.Errorf("illegal version format: %q", vsnstr)
}

// parseNumberString retrieves major, minor, and patch number
// of the passed string.
func parseNumberString(nstr string) ([]int, error) {
	nstrs := strings.Split(nstr, ".")
	if len(nstrs) < 1 || len(nstrs) > 3 {
		return nil, fmt.Errorf("illegal version format: %q", nstr)
	}
	vsn := []int{1, 0, 0}
	for i, nstr := range nstrs {
		num, err := strconv.Atoi(nstr)
		if err != nil {
			return nil, fmt.Errorf("illegal version format: %v", err)
		}
		if num < 0 {
			return nil, fmt.Errorf("illegal version format: %q", nstr)
		}
		vsn[i] = num
	}
	return vsn, nil
}

// EOF
