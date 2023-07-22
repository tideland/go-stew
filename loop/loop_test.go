// Tideland Go Stew - Loop - Unit Tests
//
// Copyright (C) 2017-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/loop"
)

//--------------------
// TESTS
//--------------------

// TestPureOK tests a loop without any options, stopping without an error.
func TestPureOK(t *testing.T) {
	started := make(chan struct{})
	stopped := make(chan struct{})
	beenThereDoneThat := false
	worker := func(ctx context.Context) error {
		defer close(stopped)
		close(started)
		for {
			select {
			case <-ctx.Done():
				beenThereDoneThat = true
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(worker)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	<-started
	l.Stop()
	<-stopped

	Assert(t, NoError(l.Err()), "stopped loop.Err() returned an error")
	Assert(t, True(beenThereDoneThat), "worker not called")
}

// TestPureError tests a loop without any options, error during stopping.
func TestPureError(t *testing.T) { // Init.
	started := make(chan struct{})
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		close(started)
		for {
			select {
			case <-ctx.Done():
				return errors.New("ouch")
			case <-time.Tick(50 * time.Millisecond):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(worker)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	<-started
	l.Stop()
	<-stopped

	Assert(t, ErrorContains(l.Err(), "ouch"), "stopped loop.Err() returned wrong error")
}

// TestContextCancelOK tests the stopping after a context cancel w/o error.
func TestContextCancelOK(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(
		worker,
		loop.WithContext(ctx),
	)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	cancel()
	<-stopped

	Assert(t, NoError(l.Err()), "stopped loop.Err() returned an error")
}

// TestContextCancelError tests the stopping after a context cancel w/ error.
func TestContextCancelError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		for {
			select {
			case <-ctx.Done():
				return errors.New("oh, no")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(
		worker,
		loop.WithContext(ctx),
	)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	cancel()
	<-stopped

	Assert(t, ErrorContains(l.Err(), "oh, no"), "stopped loop.Err() returned wrong error")
}

// TestFinalizerOK tests calling a finalizer returning an own error.
func TestFinalizerOK(t *testing.T) {
	stopped := make(chan struct{})
	finalized := false
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	finalizer := func(err error) error {
		defer close(stopped)
		Assert(t, NoError(err), "finalizer called with error")
		finalized = true
		return err
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
	)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	l.Stop()
	<-stopped

	Assert(t, NoError(l.Err()), "stopped loop.Err() returned an error")
	Assert(t, True(finalized), "finalizer not called")
}

// TestFinalizerError tests the stopping with an error but
// finalizer returns an own one.
func TestFinalizerError(t *testing.T) {
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return errors.New("don't want to stop")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	finalizer := func(err error) error {
		defer close(stopped)
		Assert(t, ErrorContains(err, "don't want to stop"), "finalizer called with wrong error")
		return errors.New("don't care")
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
	)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	l.Stop()
	<-stopped

	Assert(t, ErrorContains(l.Err(), "don't care"), "stopped loop.Err() returned wrong error")
}

// TestInternalError tests the stopping after an internal error.
func TestInternalError(t *testing.T) {
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return errors.New("time over")
			}
		}
	}
	l, err := loop.Go(worker)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	<-stopped
	l.Stop()

	Assert(t, ErrorContains(l.Err(), "time over"), "stopped loop.Err() returned wrong error")
}

// TestRepairerOK tests the stopping without an error if Loop has a repairer.
// Repairer must never been called.
func TestRepairerOK(t *testing.T) {
	repaired := make(chan struct{})
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	repairer := func(reason interface{}) error {
		defer close(repaired)
		return nil
	}
	l, err := loop.Go(
		worker,
		loop.WithRepairer(repairer),
	)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	// Test.
	l.Stop()

	select {
	case <-repaired:
		Assert(t, Fail("repaired has been closed"), "repairer called")
	case <-time.After(100 * time.Millisecond):
	}
}

// TestRepairerErrorOK tests the stopping with an error if Loop has a repairer.
// Repairer must never been called.
func TestRepairerErrorOK(t *testing.T) {
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				close(stopped)
				return errors.New("oh, no")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	repairer := func(reason interface{}) error {
		return fmt.Errorf("unexpected")
	}
	l, err := loop.Go(
		worker,
		loop.WithRepairer(repairer),
	)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	l.Stop()
	<-stopped

	Assert(t, ErrorContains(l.Err(), "oh, no"), "stopped loop.Err() returned wrong error")
}

// TestRecoverPanics tests the stopping handling and later stopping
// after panics.
func TestRecoverPanicsOK(t *testing.T) {
	panics := 0
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				panic("bam")
			}
		}
	}
	finalizer := func(err error) error {
		defer close(stopped)
		return err
	}
	repairer := func(reason interface{}) error {
		panics++
		if panics > 10 {
			return fmt.Errorf("too many panics: %v", reason)
		}
		return nil
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
		loop.WithRepairer(repairer),
	)
	Assert(t, NoError(err), "loop.Go() failed")
	Assert(t, NotNil(l), "loop is nil")
	Assert(t, NoError(l.Err()), "initial loop.Err() returned an error")

	<-stopped

	Assert(t, ErrorContains(l.Err(), "too many panics: bam"), "stopped loop.Err() returned wrong error")
}

//--------------------
// EXAMPLES
//--------------------

// ExampleWorker shows the usage of Loop with no repairer. The inner loop
// contains a select listening to the channel returned by Closer.Done().
// Other channels are for the standard communication with the Loop.
func ExampleWorker() {
	prints := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	// Sample loop worker.
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				// We shall stop.
				return ctx.Err()
			case str := <-prints:
				// Standard work of example loop.
				if str == "panic" {
					return errors.New("panic")
				}
				println(str)
			}
		}
	}
	l, err := loop.Go(worker, loop.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	prints <- "Hello"
	prints <- "World"

	// cancel() terminates the loop via the context.
	cancel()

	// Returned error must be nil in this example.
	if l.Err() != nil {
		panic(l.Err())
	}
}

// ExampleRepairer demonstrates the usage of a repairer.
func ExampleRepairer() {
	panics := make(chan string)
	// Sample loop worker.
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case str := <-panics:
				panic(str)
			}
		}
	}
	// Repairer function checks the reasion. "never mind" will
	// be repaired, all others lead to an error. The repairer
	// is also responsable for fixing the owners state crashed
	// during panic.
	repairer := func(reason interface{}) error {
		why := reason.(string)
		if why == "never mind" {
			return nil
		}
		return fmt.Errorf("worker panic: %v", why)
	}
	l, err := loop.Go(worker, loop.WithRepairer(repairer))
	if err != nil {
		panic(err)
	}
	l.Stop()
}

// EOF
