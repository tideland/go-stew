// Tideland Go Stew - Environments
//
// Copyright (C) 2012-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package environments // import "tideland.dev/go/stew/environments"

//--------------------
// IMPORTS
//--------------------

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

//--------------------
// TEMPDIR
//--------------------

// TempDir represents a temporary directory and possible subdirectories
// for testing purposes. It simply is created with
//
//	td, err := environments.NewTempDir()
//	if err != nil { ... }
//	defer td.Restore()
//
//	tdName := td.String()
//	subName:= td.Mkdir("my", "sub", "directory")
//
// The deferred Restore() removes the temporary directory with all
// contents.
type TempDir struct {
	dir string
}

// NewTempDir creates a new temporary directory usable for direct
// usage or further subdirectories.
func NewTempDir() (*TempDir, error) {
	id := make([]byte, 8)
	td := &TempDir{}
	for i := 0; i < 256; i++ {
		_, err := rand.Read(id[:])
		if err != nil {
			return nil, fmt.Errorf("cannot create random id: %v", err)
		}
		dir := filepath.Join(os.TempDir(), fmt.Sprintf("stew-environment-%x", id))
		if err = os.Mkdir(dir, 0700); err == nil {
			td.dir = dir
			break
		}
		if td.dir == "" {
			return nil, fmt.Errorf("cannot create temporary directory %q: %v", td.dir, err)
		}
	}
	return td, nil
}

// Restore deletes the temporary directory and all contents.
func (td *TempDir) Restore() {
	err := os.RemoveAll(td.dir)
	if err != nil {
		panic(fmt.Errorf("cannot remove temporary directory %q: %v", td.dir, err))
	}
}

// Mkdir creates a potentially nested directory inside the
// temporary directory.
func (td *TempDir) Mkdir(name ...string) (string, error) {
	innerName := filepath.Join(name...)
	fullName := filepath.Join(td.dir, innerName)
	if err := os.MkdirAll(fullName, 0700); err != nil {
		return "", fmt.Errorf("cannot create nested temporary directory %q: %v", fullName, err)
	}
	return fullName, nil
}

// String returns the temporary directory.
func (td *TempDir) String() string {
	return td.dir
}

//--------------------
// VARIABLES
//--------------------

// Variables allows to change and restore environment variables. The
// same variable can be set multiple times. Simply do
//
//	ev := environments.NewVariables()
//	defer ev.Restore()
//
//	ev.Set("MY_VAR", myValue)
//
//	...
//
//	ev.Set("MY_VAR", anotherValue)
//
// The deferred Restore() resets to the original values.
type Variables struct {
	vars map[string]string
}

// NewVariables create a new changer for environment variables.
func NewVariables() *Variables {
	v := &Variables{
		vars: make(map[string]string),
	}
	return v
}

// Restore resets all changed environment variables
func (v *Variables) Restore() error {
	for key, value := range v.vars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("cannot reset environment variable %q: %v", key, err)
		}
	}
	return nil
}

// Set sets an environment variable to a new value.
func (v *Variables) Set(key, value string) error {
	ov := os.Getenv(key)
	_, ok := v.vars[key]
	if !ok {
		v.vars[key] = ov
	}
	if err := os.Setenv(key, value); err != nil {
		return fmt.Errorf("cannot set environment variable %q: %v", key, err)
	}
	return nil
}

// Unset unsets an environment variable.
func (v *Variables) Unset(key string) error {
	ov := os.Getenv(key)
	_, ok := v.vars[key]
	if !ok {
		v.vars[key] = ov
	}
	if err := os.Unsetenv(key); err != nil {
		return fmt.Errorf("cannot unset environment variable %q: %v", key, err)
	}
	return nil
}

// EOF
