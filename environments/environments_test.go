// Tideland Go Stew - Environments - Unit Tests
//
// Copyright (C) 2012-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package environments_test

//--------------------
// IMPORTS
//--------------------

import (
	"os"
	"testing"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/environments"
)

//--------------------
// TESTS
//--------------------

// TestTempDirCreate tests the creation of temporary directories.
func TestTempDirCreate(t *testing.T) {
	testDir := func(dir string) {
		fi, err := os.Stat(dir)
		Assert(t, NoError(err), "os.Stat worked")
		Assert(t, True(fi.IsDir()), "file is a directory")
		Assert(t, Equal(fi.Mode().Perm(), os.FileMode(0700)), "directory has correct permissions")
	}

	td, err := environments.NewTempDir()
	Assert(t, NoError(err), "new temp dir created")
	defer td.Restore()

	tds := td.String()
	Assert(t, NotEmpty(tds), "temp dir has a name")
	testDir(tds)

	sda, err := td.Mkdir("subdir", "foo")
	Assert(t, NoError(err), "subdir created")
	Assert(t, NotEmpty(sda), "subdir has name foo")
	testDir(sda)

	sdb, err := td.Mkdir("subdir", "bar")
	Assert(t, NoError(err), "subdir created")
	Assert(t, NotEmpty(sdb), "subdir has name bar")
	testDir(sdb)
}

// TestTempDirRestore tests the restoring of temporary created
// directories.
func TestTempDirRestore(t *testing.T) {
	td, err := environments.NewTempDir()
	Assert(t, NoError(err), "new temp dir created")
	Assert(t, NotNil(td), "temp dir is not nil")

	tds := td.String()
	fi, err := os.Stat(tds)
	Assert(t, NoError(err), "temp dir exists")
	Assert(t, True(fi.IsDir()), "temp dir is a directory")

	td.Restore()

	_, err = os.Stat(tds)
	Assert(t, ErrorMatches(err, "stat .* no such file or directory"), "temp dir does not exist")
}

// TestEnvVarsSet tests the setting of temporary environment variables.
func TestEnvVarsSet(t *testing.T) {
	testEnv := func(key, value string) {
		v := os.Getenv(key)
		Assert(t, Equal(v, value), "environment variable has value")
	}

	ev := environments.NewVariables()
	Assert(t, NotNil(ev), "environment variables is not nil")
	defer ev.Restore()

	ev.Set("TESTING_ENV_A", "FOO")
	testEnv("TESTING_ENV_A", "FOO")
	ev.Set("TESTING_ENV_B", "BAR")
	testEnv("TESTING_ENV_B", "BAR")

	ev.Unset("TESTING_ENV_A")
	testEnv("TESTING_ENV_A", "")
}

// TestEnvVarsREstore tests the restoring of temporary set environment
// variables.
func TestEnvVarsRestore(t *testing.T) {
	testEnv := func(key, value string) {
		v := os.Getenv(key)
		Assert(t, Equal(v, value), "environment variable has value")
	}

	ev := environments.NewVariables()
	Assert(t, NotNil(ev), "environment variables is not nil")

	path := os.Getenv("PATH")
	Assert(t, NotEmpty(path), "PATH is not empty")

	ev.Set("PATH", "/foo:/bar/bin")
	testEnv("PATH", "/foo:/bar/bin")
	ev.Set("PATH", "/bar:/foo:/yadda/bin")
	testEnv("PATH", "/bar:/foo:/yadda/bin")

	ev.Restore()

	testEnv("PATH", path)
}

// EOF
