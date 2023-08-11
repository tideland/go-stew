// Tideland Go Stew - Etc
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package etc // import "tideland.dev/go/stew/etc"

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"tideland.dev/go/stew/dynaj"
)

//--------------------
// ALIASES
//--------------------

type ID = dynaj.ID
type Path = dynaj.Path
type Value = dynaj.Value
type Object = dynaj.Object
type Array = dynaj.Array

//--------------------
// ETC
//--------------------

// Etc contains the read etc configuration and provides access to
// it. The syntax is JSON but extended by templates. These are
// formatted as [[reference||default]]. The reference can be an
// environment variable or a path inside the configuration. If
// the reference cannot be found the default value is used.
type Etc struct {
	mu   sync.RWMutex
	data *dynaj.Document
	orig *dynaj.Document
}

// Read reads the SML source of the configuration from a
// reader, parses it, and returns the etc instance.
func Read(source io.Reader) (*Etc, error) {
	// Read and parse the source.
	var buf bytes.Buffer
	_, err := buf.ReadFrom(source)
	if err != nil {
		return nil, fmt.Errorf("cannot read source: %v", err)
	}
	data, err := dynaj.Unmarshal(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("invalid source format: %v", err)
	}
	orig, _ := dynaj.Unmarshal(buf.Bytes())
	etc := &Etc{
		data: data,
		orig: orig,
	}
	return etc, nil
}

// ReadString reads the SML source of the configuration from a
// string, parses it, and returns the etc instance.
func ReadString(source string) (*Etc, error) {
	return Read(strings.NewReader(source))
}

// ReadFile reads the SML source of a configuration file,
// parses it, and returns the etc instance.
func ReadFile(filename string) (*Etc, error) {
	source, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file '%s': %v", filename, err)
	}
	return ReadString(string(source))
}

// At returns the Value at a given path.
func (e *Etc) At(path ...ID) *Accessor {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return newValue(e, path)
}

// Write writes the configuration as indented JSON to the passed writer. All
// macros will stay as long as the aren't explicitly overwritten.
func (e *Etc) Write(target io.Writer) error {
	buf, err := e.orig.MarshalJSONIndent()
	if err != nil {
		return fmt.Errorf("cannot write configuration: %v", err)
	}
	_, err = target.Write(buf)
	if err != nil {
		return fmt.Errorf("cannot write configuration: %v", err)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e *Etc) String() string {
	var sw strings.Builder
	err := e.Write(&sw)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return sw.String()
}

//--------------------
// CONTEXT
//--------------------

// etcKey is the key for the configuration in the context.
type etcKey string

const etcID etcKey = "etc"

// NewContext returns a new context that carries a configuration.
func NewContext(ctx context.Context, etc *Etc) context.Context {
	return context.WithValue(ctx, etcID, etc)
}

// FromContext returns the configuration stored in ctx, if any.
func FromContext(ctx context.Context) (*Etc, bool) {
	cfg, ok := ctx.Value(etcID).(*Etc)
	return cfg, ok
}

// EOF
