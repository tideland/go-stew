// Tideland Go Stew - QA Environments
//
// Copyright (C) 2012-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package qaenv // import "tideland.dev/go/stew/qaenv"

//--------------------
// IMPORTS
//--------------------

import (
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
//	td, err := qaenv.MkdirTemp()
//	if err != nil { ... }
//	defer td.Restore()
//
//	tdName := td.String()
//	subName:= td.Mkdir("my", "sub", "directory")
//
// The deferred Restore() removes the temporary directory with all
// contents.
type TempDir struct {
	dirname string
}

// MkdirTemp creates a new temporary directory usable for direct
// usage or further subdirectories. The pattern is used as individual
// name part of the directory name.
func MkdirTemp(pattern string) (*TempDir, error) {
	dirname, err := os.MkdirTemp("", pattern)
	if err != nil {
		return nil, fmt.Errorf("cannot create temporary directory: %v", err)
	}
	return &TempDir{dirname: dirname}, nil
}

// Restore deletes the temporary directory and all contents.
func (td *TempDir) Restore() error {
	err := os.RemoveAll(td.dirname)
	if err != nil {
		return fmt.Errorf("cannot remove temporary directory %q: %v", td.dirname, err)
	}
	return nil
}

// Mkdir creates a potentially nested directory inside the
// temporary directory.
func (td *TempDir) Mkdir(name ...string) (string, error) {
	innerName := filepath.Join(name...)
	fullName := filepath.Join(td.dirname, innerName)
	if err := os.MkdirAll(fullName, 0700); err != nil {
		return "", fmt.Errorf("cannot create nested temporary directory %q: %v", fullName, err)
	}
	return fullName, nil
}

// WriteFile writes a file with the passed data into the temporary directory.
func (td *TempDir) WriteFile(filename string, data []byte) (string, error) {
	fullName := filepath.Join(td.dirname, filename)
	if err := os.WriteFile(fullName, data, 0600); err != nil {
		return "", fmt.Errorf("cannot write file %q: %v", fullName, err)
	}
	return fullName, nil
}

// OpenFile opens a file inside the temporary directory.
func (td *TempDir) OpenFile(filename string) (*os.File, error) {
	fullName := filepath.Join(td.dirname, filename)
	file, err := os.Open(fullName)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %v", fullName, err)
	}
	return file, nil
}

// RemoveFile removes a file inside the temporary directory.
func (td *TempDir) RemoveFile(filename string) error {
	fullName := filepath.Join(td.dirname, filename)
	if err := os.Remove(fullName); err != nil {
		return fmt.Errorf("cannot remove file %q: %v", fullName, err)
	}
	return nil
}

// String returns the temporary directory.
func (td *TempDir) String() string {
	return td.dirname
}

//--------------------
// VARIABLES
//--------------------

// Env allows to change and restore environment variables. The
// same variable can be set multiple times. Simply do
//
//	env := qaenv.NewEnvironment()
//	defer env.Restore()
//
//	ev.Set("MY_VAR", myValue)
//
//	...
//
//	ev.Set("MY_VAR", anotherValue)
//
// The deferred Restore() resets to the original values.
type Env struct {
	vars map[string]string
}

// NewEinvironment create a new changer for environment variables.
func NewEinvironment() *Env {
	env := &Env{
		vars: make(map[string]string),
	}
	return env
}

// Restore resets all changed environment variables
func (env *Env) Restore() error {
	for key, value := range env.vars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("cannot reset environment variable %q: %v", key, err)
		}
	}
	return nil
}

// Set sets an environment variable to a new value.
func (env *Env) Set(key, value string) error {
	ov := os.Getenv(key)
	_, ok := env.vars[key]
	if !ok {
		env.vars[key] = ov
	}
	if err := os.Setenv(key, value); err != nil {
		return fmt.Errorf("cannot set environment variable %q: %v", key, err)
	}
	return nil
}

// Get gets an environment variable.
func (env *Env) Get(key string) string {
	return os.Getenv(key)
}

// Unset unsets an environment variable.
func (env *Env) Unset(key string) error {
	ov := os.Getenv(key)
	_, ok := env.vars[key]
	if !ok {
		env.vars[key] = ov
	}
	if err := os.Unsetenv(key); err != nil {
		return fmt.Errorf("cannot unset environment variable %q: %v", key, err)
	}
	return nil
}

// EOF
