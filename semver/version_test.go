// Tideland Go Stew - Semantic Versions - Unit Tests
//
// Copyright (C) 2014-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package semver_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/semver"
)

//--------------------
// TESTS
//--------------------

// TestNewVersion tests the creation of new versions and their
// accessor methods.
func TestNewVersion(t *testing.T) {
	tests := []struct {
		id         string
		vsn        *semver.Version
		major      int
		minor      int
		patch      int
		preRelease string
		metadata   string
	}{
		{
			id:         "1.2.3",
			vsn:        semver.NewVersion(1, 2, 3),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "1.0.3",
			vsn:        semver.NewVersion(1, -2, 3),
			major:      1,
			minor:      0,
			patch:      3,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "1.2.3-alpha.2014-08-03",
			vsn:        semver.NewVersion(1, 2, 3, "alpha", "2014-08-03"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alpha.2014-08-03",
			metadata:   "",
		}, {
			id:         "1.2.3-alphabeta.7.11",
			vsn:        semver.NewVersion(1, 2, 3, "alpha beta", "007", "1+1"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alphabeta.7.11",
			metadata:   "",
		}, {
			id:         "1.2.3+007.a",
			vsn:        semver.NewVersion(1, 2, 3, semver.Metadata, "007", "a"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "",
			metadata:   "007.a",
		}, {
			id:         "1.2.3-alpha+007.a",
			vsn:        semver.NewVersion(1, 2, 3, "alpha", semver.Metadata, "007", "a"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alpha",
			metadata:   "007.a",
		}, {
			id:         "1.2.3-ALPHA+007.a",
			vsn:        semver.NewVersion(1, 2, 3, "ALPHA", semver.Metadata, "007", "a"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "ALPHA",
			metadata:   "007.a",
		},
	}
	// Perform tests.
	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			Assert(t, Equal(test.vsn.Major(), test.major), "major semver matches")
			Assert(t, Equal(test.vsn.Minor(), test.minor), "minor semver matches")
			Assert(t, Equal(test.vsn.Patch(), test.patch), "patch semver matches")
			Assert(t, Equal(test.vsn.PreRelease(), test.preRelease), "pre-release semver matches")
			Assert(t, Equal(test.vsn.Metadata(), test.metadata), "metadata semver matches")
			Assert(t, Equal(test.vsn.String(), test.id), "semver string matches")
		})
	}
}

// TestParse tests the creation of new versions and their
// accessor methods by parsing strings.
func TestParse(t *testing.T) {
	tests := []struct {
		id         string
		vsn        string
		err        string
		major      int
		minor      int
		patch      int
		preRelease string
		metadata   string
	}{
		{
			id:         "1",
			vsn:        "1.0.0",
			major:      1,
			minor:      0,
			patch:      0,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "1.1",
			vsn:        "1.1.0",
			major:      1,
			minor:      1,
			patch:      0,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "1.2.3",
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "1.0.3",
			major:      1,
			minor:      0,
			patch:      3,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "1.2.3-alpha.2016-11-14",
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alpha.2016-11-14",
			metadata:   "",
		}, {
			id:         "1.2.3-alphabeta.7.11",
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alphabeta.7.11",
			metadata:   "",
		}, {
			id:         "1.2.3+007.a",
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "",
			metadata:   "007.a",
		}, {
			id:         "1.2.3-alpha+007.a",
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alpha",
			metadata:   "007.a",
		}, {
			id:         "1.2.3-ALPHA+007.a",
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "ALPHA",
			metadata:   "007.a",
		}, {
			id:  "",
			err: "illegal version format: strconv.Atoi: parsing",
		}, {
			id:  "foobar",
			err: "illegal version format: strconv.Atoi: parsing",
		}, {
			id:  "1.a",
			err: "illegal version format: strconv.Atoi: parsing",
		}, {
			id:  "1,1",
			err: "illegal version format: strconv.Atoi: parsing",
		}, {
			id:  "-1",
			err: "illegal version format: strconv.Atoi: parsing",
		}, {
			id:  "1.-1",
			err: "illegal version format: strconv.Atoi: parsing",
		}, {
			id:  "+",
			err: "illegal version format: strconv.Atoi: parsing",
		}, {
			id:  "1.1.1.1",
			err: "illegal version format: \"1.1.1.1\"",
		},
	}
	// Perform tests.
	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			vsn, err := semver.Parse(test.id)
			if test.err != "" {
				Assert(t, ErrorContains(err, test.err), "error matches")
				return
			}
			Assert(t, NoError(err), "no error")
			Assert(t, NotNil(vsn), "semver created")
			Assert(t, Equal(vsn.Major(), test.major), "major semver matches")
			Assert(t, Equal(vsn.Minor(), test.minor), "minor semver matches")
			Assert(t, Equal(vsn.Patch(), test.patch), "patch semver matches")
			Assert(t, Equal(vsn.PreRelease(), test.preRelease), "pre-release semver matches")
			Assert(t, Equal(vsn.Metadata(), test.metadata), "metadata semver matches")
			if test.vsn != "" {
				Assert(t, Equal(vsn.String(), test.vsn), "semver string matches")
			} else {
				Assert(t, Equal(vsn.String(), test.id), "semver string matches")
			}
		})
	}
}

// TestCompare tests the comparing of two versions.
func TestCompare(t *testing.T) {
	tests := []struct {
		vsnA       *semver.Version
		vsnB       *semver.Version
		precedence semver.Precedence
		level      semver.Level
	}{
		{
			vsnA:       semver.NewVersion(1, 2, 3),
			vsnB:       semver.NewVersion(1, 2, 3),
			precedence: semver.Equal,
			level:      semver.All,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3),
			vsnB:       semver.NewVersion(1, 2, 4),
			precedence: semver.Older,
			level:      semver.Patch,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3),
			vsnB:       semver.NewVersion(1, 3, 3),
			precedence: semver.Older,
			level:      semver.Minor,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3),
			vsnB:       semver.NewVersion(2, 2, 3),
			precedence: semver.Older,
			level:      semver.Major,
		}, {
			vsnA:       semver.NewVersion(3, 2, 1),
			vsnB:       semver.NewVersion(1, 2, 3),
			precedence: semver.Newer,
			level:      semver.Major,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha"),
			vsnB:       semver.NewVersion(1, 2, 3),
			precedence: semver.Older,
			level:      semver.PreRelease,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha", "1"),
			vsnB:       semver.NewVersion(1, 2, 3, "alpha"),
			precedence: semver.Older,
			level:      semver.PreRelease,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha", "1"),
			vsnB:       semver.NewVersion(1, 2, 3, "alpha", "2"),
			precedence: semver.Older,
			level:      semver.PreRelease,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha", "4711"),
			vsnB:       semver.NewVersion(1, 2, 3, "alpha", "471"),
			precedence: semver.Newer,
			level:      semver.PreRelease,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha", "48"),
			vsnB:       semver.NewVersion(1, 2, 3, "alpha", "4711"),
			precedence: semver.Older,
			level:      semver.PreRelease,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "1"),
			vsnB:       semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "2"),
			precedence: semver.Equal,
			level:      semver.All,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "2"),
			vsnB:       semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "1"),
			precedence: semver.Equal,
			level:      semver.All,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha", semver.Metadata, "alpha", "2"),
			vsnB:       semver.NewVersion(1, 2, 3, "alpha", semver.Metadata, "alpha", "1"),
			precedence: semver.Equal,
			level:      semver.All,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha", "48", semver.Metadata, "alpha", "2"),
			vsnB:       semver.NewVersion(1, 2, 3, "alpha", "4711", semver.Metadata, "alpha", "1"),
			precedence: semver.Older,
			level:      semver.PreRelease,
		}, {
			vsnA:       semver.NewVersion(1, 2, 3, "alpha", "2"),
			vsnB:       semver.NewVersion(1, 2, 3, "alpha", "1b"),
			precedence: semver.Newer,
			level:      semver.PreRelease,
		},
	}
	// Perform tests.
	for i, test := range tests {
		id := fmt.Sprintf("compare test #%d: %q <> %q -> %d / %s", i, test.vsnA, test.vsnB, test.precedence, test.level)
		t.Run(id, func(t *testing.T) {
			precedence, level := test.vsnA.Compare(test.vsnB)
			Assert(t, Equal(precedence, test.precedence), "precedence matches")
			Assert(t, Equal(level, test.level), "level matches")
		})
	}
}

// TestLess tests if a semver is less (older) than another.
func TestLess(t *testing.T) {
	tests := []struct {
		vsnA *semver.Version
		vsnB *semver.Version
		less bool
	}{
		{
			vsnA: semver.NewVersion(1, 2, 3),
			vsnB: semver.NewVersion(1, 2, 3),
			less: false,
		}, {
			vsnA: semver.NewVersion(1, 2, 3),
			vsnB: semver.NewVersion(1, 2, 4),
			less: true,
		}, {
			vsnA: semver.NewVersion(1, 2, 3),
			vsnB: semver.NewVersion(1, 3, 3),
			less: true,
		}, {
			vsnA: semver.NewVersion(1, 2, 3),
			vsnB: semver.NewVersion(2, 2, 3),
			less: true,
		}, {
			vsnA: semver.NewVersion(3, 2, 1),
			vsnB: semver.NewVersion(1, 2, 3),
			less: false,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha"),
			vsnB: semver.NewVersion(1, 2, 3),
			less: true,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha", "1"),
			vsnB: semver.NewVersion(1, 2, 3, "alpha"),
			less: true,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha", "1"),
			vsnB: semver.NewVersion(1, 2, 3, "alpha", "2"),
			less: true,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha", "4711"),
			vsnB: semver.NewVersion(1, 2, 3, "alpha", "471"),
			less: false,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha", "48"),
			vsnB: semver.NewVersion(1, 2, 3, "alpha", "4711"),
			less: true,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "1"),
			vsnB: semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "2"),
			less: false,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "2"),
			vsnB: semver.NewVersion(1, 2, 3, semver.Metadata, "alpha", "1"),
			less: false,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha", semver.Metadata, "alpha", "2"),
			vsnB: semver.NewVersion(1, 2, 3, "alpha", semver.Metadata, "alpha", "1"),
			less: false,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha", "48", semver.Metadata, "alpha", "2"),
			vsnB: semver.NewVersion(1, 2, 3, "alpha", "4711", semver.Metadata, "alpha", "1"),
			less: true,
		}, {
			vsnA: semver.NewVersion(1, 2, 3, "alpha", "2"),
			vsnB: semver.NewVersion(1, 2, 3, "alpha", "1b"),
			less: false,
		},
	}
	// Perform tests.
	for i, test := range tests {
		id := fmt.Sprintf("less test #%d: %q <> %q -> %v", i, test.vsnA, test.vsnB, test.less)
		t.Run(id, func(t *testing.T) {
			Assert(t, Equal(test.vsnA.Less(test.vsnB), test.less), "less matches")
		})
	}
}

// EOF
