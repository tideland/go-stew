// Tideland Go Stew - QA Environments - Unit Tests
//
// Copyright (C) 2012-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package qaenv_test

//--------------------
// IMPORTS
//--------------------

import (
	"os"
	"testing"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/qaenv"
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

	td, err := qaenv.MkdirTemp("stew")
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
	td, err := qaenv.MkdirTemp("test")
	Assert(t, NoError(err), "new temp dir created")
	Assert(t, NotNil(td), "temp dir is not nil")

	tds := td.String()
	fi, err := os.Stat(tds)
	Assert(t, NoError(err), "temp dir exists")
	Assert(t, True(fi.IsDir()), "temp dir is a directory")

	err = td.Restore()
	Assert(t, NoError(err), "temp dir restored")

	_, err = os.Stat(tds)
	Assert(t, ErrorMatches(err, "stat .* no such file or directory"), "temp dir does not exist")
}

// TestTempDirWrite tests the writing of files into temporary
// directories.
func TestTempDirWrite(t *testing.T) {
	td, err := qaenv.MkdirTemp("test")
	Assert(t, NoError(err), "new temp dir created")
	defer td.Restore()

	fn, err := td.WriteFile("foo.txt", []byte("foo"))
	Assert(t, NoError(err), "file written")
	Assert(t, NotEmpty(fn), "file has name")

	file, err := td.OpenFile("foo.txt")
	Assert(t, NoError(err), "file opened")
	defer file.Close()

	fi, err := file.Stat()
	Assert(t, NoError(err), "file stat")
	Assert(t, Equal(fi.Size(), int64(3)), "file size")
}

// TestTempDirRemove tests the removing of files from temporary.
func TestTempDirRemove(t *testing.T) {
	td, err := qaenv.MkdirTemp("test")
	Assert(t, NoError(err), "new temp dir created")
	defer td.Restore()

	fn, err := td.WriteFile("foo.txt", []byte("foo"))
	Assert(t, NoError(err), "file written")
	Assert(t, NotEmpty(fn), "file has name")

	err = td.RemoveFile("foo.txt")
	Assert(t, NoError(err), "file removed")

	_, err = td.OpenFile("foo.txt")
	Assert(t, ErrorMatches(err, "open .* no such file or directory"), "file does not exist")
}

// TestEnvVarsSet tests the setting of temporary environment variables.
func TestEnvVarsSet(t *testing.T) {
	testEnv := func(key, value string) {
		v := os.Getenv(key)
		Assert(t, Equal(v, value), "environment variable has value")
	}

	env := qaenv.NewEinvironment()
	Assert(t, NotNil(env), "environment variables is not nil")
	defer env.Restore()

	env.Set("TESTING_ENV_A", "FOO")
	testEnv("TESTING_ENV_A", "FOO")
	env.Set("TESTING_ENV_B", "BAR")
	testEnv("TESTING_ENV_B", "BAR")

	env.Unset("TESTING_ENV_A")
	testEnv("TESTING_ENV_A", "")
}

// TestEnvVarsREstore tests the restoring of temporary set environment
// variables.
func TestEnvVarsRestore(t *testing.T) {
	testEnv := func(key, value string) {
		v := os.Getenv(key)
		Assert(t, Equal(v, value), "environment variable has value")
	}

	env := qaenv.NewEinvironment()
	Assert(t, NotNil(env), "environment variables is not nil")

	path := os.Getenv("PATH")
	Assert(t, NotEmpty(path), "PATH is not empty")

	env.Set("PATH", "/foo:/bar/bin")
	testEnv("PATH", "/foo:/bar/bin")
	env.Set("PATH", "/bar:/foo:/yadda/bin")
	testEnv("PATH", "/bar:/foo:/yadda/bin")

	env.Restore()

	testEnv("PATH", path)
}

// EOF
