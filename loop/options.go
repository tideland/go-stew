// Tideland Go Stew - Loop
//
// Copyright (C) 2017-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop // import "tideland.dev/go/stew/loop"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
)

//--------------------
// OPTIONS
//--------------------

// Option defines the signature of an option setting function.
type Option func(loop *Loop) error

// WithContext allows to pass a context for cancellation or timeout.
func WithContext(ctx context.Context) Option {
	return func(l *Loop) error {
		if ctx == nil {
			return fmt.Errorf("invalid loop option: context is nil")
		}
		l.ctx = ctx
		return nil
	}
}

// WithRepairer defines the panic handler of a loop.
func WithRepairer(repairer Repairer) Option {
	return func(l *Loop) error {
		if repairer == nil {
			return fmt.Errorf("invalid loop option: repairer is nil")
		}
		l.repairer = repairer
		return nil
	}
}

// WithFinalizer sets a function for finalizing the
// work of a Loop.
func WithFinalizer(finalizer Finalizer) Option {
	return func(l *Loop) error {
		if finalizer == nil {
			return fmt.Errorf("invalid loop option: finalizer is nil")
		}
		l.finalizer = finalizer
		return nil
	}
}

// EOF