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
	"sync"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

// timeout defines the time to wait for signals.
const timeout = 5 * time.Second

// Status defines the status of a loop.
type Status int

const (
	StatusStopped Status = iota
	StatusStarting
	StatusWorking
	StatusRepairing
	StatusStopping
	StatusFinalizing
	StatusError
)

//--------------------
// FUNCTION TYPES
//--------------------

// Worker discribes the backend function running the loop.
type Worker func(ctx context.Context) error

// Repairer allows the loop goroutine to react on a panic
// during its work. Returning a nil the loop will be continued
// by calling the worker again.
type Repairer func(reason interface{}) error

// Finalizer is called with the final error if the backend
// loop terminates.
type Finalizer func(err error) error

//--------------------
// LOOP
//--------------------

// Loop manages running for-select-loops in the background as goroutines
// in a controlled way. Users can get information about status and possible
// failure as well as control how to stop, restart, or recover via
// options.
type Loop struct {
	mu        sync.RWMutex
	ctx       context.Context
	cancel    func()
	worker    Worker
	repairer  Repairer
	finalizer Finalizer
	status    Status
	err       error
}

// Go starts a loop running the given worker with the
// given options.
func Go(worker Worker, options ...Option) (*Loop, error) {
	// Init with default values.
	l := &Loop{
		worker: worker,
		status: StatusStarting,
	}
	for _, option := range options {
		if err := option(l); err != nil {
			return nil, err
		}
	}
	// Ensure default settings for context.
	if l.ctx == nil {
		l.ctx, l.cancel = context.WithCancel(context.Background())
	} else {
		l.ctx, l.cancel = context.WithCancel(l.ctx)
	}
	// Start backend.
	started := make(chan struct{})
	go l.backend(started)
	select {
	case <-started:
		l.status = StatusWorking
		return l, nil
	case <-time.After(timeout):
		l.status = StatusError
		return nil, fmt.Errorf("loop starting timeout after %.1f seconds", timeout.Seconds())
	}
}

// Status returns the current status of the Loop.
func (l *Loop) Status() Status {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.status
}

// Err returns information about a possible error of the Loop.
func (l *Loop) Err() error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.err
}

// Stop terminates the Loop backend. It works asynchronous as
// the goroutine may need time for cleanup. Anyone wanting to
// be notified on state has to handle it in a Finalizer.
func (l *Loop) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.status == StatusWorking {
		l.status = StatusStopping
		l.cancel()
	}
}

// backend runs the loop worker as goroutine as long as
// the it isn't terminated or recovery returned false.
func (l *Loop) backend(started chan struct{}) {
	defer l.finalize()
	l.status = StatusWorking
	close(started)
	for l.status == StatusWorking {
		l.work()
	}
}

// work wraps the worker and handles possible panics.
func (l *Loop) work() {
	defer func() {
		// Check and handle panics!
		reason := recover()
		switch {
		case reason != nil && l.repairer != nil:
			// Try to repair.
			l.mu.Lock()
			l.status = StatusRepairing
			err := l.repairer(reason)
			l.err = err
			if l.err == nil {
				// Success, continue.
				l.status = StatusWorking
			} else {
				// Failure, stop.
				l.status = StatusError
			}
			l.mu.Unlock()
		case reason != nil && l.repairer == nil:
			// Accept panic.
			l.mu.Lock()
			l.err = fmt.Errorf("loop panic: %v", reason)
			l.status = StatusError
			l.mu.Unlock()
		}
	}()
	// Work without panic.
	err := l.worker(l.ctx)
	l.mu.Lock()
	l.err = err
	if l.err == nil {
		l.status = StatusStopped
	} else {
		l.status = StatusError
	}
	l.mu.Unlock()
}

// finalize takes care for a clean loop finalization.
func (l *Loop) finalize() {
	l.mu.Lock()
	defer l.mu.Unlock()
	status := l.status
	l.status = StatusFinalizing
	if l.finalizer != nil {
		l.err = l.finalizer(l.err)
	}
	l.status = status
}

// EOF
