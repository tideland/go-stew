// Tideland Go Stew - Dynamic JSON
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dynaj // import "tideland.dev/go/stew/dynaj"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
)

//--------------------
// PROCESSOR
//--------------------

// Handler defines the signature of function for processing
// an Accessor. It will be called by the Processor and here
// for the one value or for all values of an Array or Object.
// In case of given Arrays or Objects a Handler can operate
// recursively.
type Handler func(acc *Accessor) error

// Processor is a function processing a JSON document starting at a
// given Path.
type Processor struct {
	acc *Accessor
	err error
}

// newProcessor creates a new processor for the given Accessor.
func newProcessor(acc *Accessor) *Processor {
	return &Processor{
		acc: acc,
	}
}

// IsError returns true if the Processor has an error.
func (p *Processor) IsError() bool {
	return p.err != nil
}

// Err returns a possible error of the Processor.
func (p *Processor) Err() error {
	return p.err
}

// Do calls the given handler for the current Accessor.
func (p *Processor) Do(handle Handler) *Processor {
	if p.acc.element == nil || p.err != nil {
		return p
	}
	switch typed := p.acc.element.(type) {
	case Array:
		for i, elem := range typed {
			path := append(p.acc.path, fmt.Sprintf("%d", i))
			acc := newAccessor(p.acc.doc, path, elem, nil)
			err := handle(acc)
			if err != nil {
				p.acc = nil
				p.err = err
				return p
			}
		}
	case Object:
		for id, elem := range typed {
			path := append(p.acc.path, id)
			acc := newAccessor(p.acc.doc, path, elem, nil)
			err := handle(acc)
			if err != nil {
				p.acc = nil
				p.err = err
				return p
			}
		}
	default:
		err := handle(p.acc)
		if err != nil {
			p.acc = nil
			p.err = err
			return p
		}
	}
	return p
}

// DeepDo calls the given handler for the current Accessor and all
// Accessors of the tree below.
func (p *Processor) DeepDo(handle Handler) *Processor {
	if p.acc.element == nil || p.err != nil {
		return p
	}
	diver := func(acc *Accessor) error {
		switch typed := acc.element.(type) {
		case Array:
			for i, elem := range typed {
				path := append(acc.path, fmt.Sprintf("%d", i))
				acc := newAccessor(acc.doc, path, elem, nil)
				if err := acc.Processor().DeepDo(handle).Err(); err != nil {
					return err
				}
			}
			return nil
		case Object:
			for id, elem := range typed {
				path := append(acc.path, id)
				acc := newAccessor(acc.doc, path, elem, nil)
				if err := acc.Processor().DeepDo(handle).Err(); err != nil {
					return err
				}
			}
			return nil
		default:
			return handle(acc)
		}
	}
	return p.Do(diver)
}

// EOF
