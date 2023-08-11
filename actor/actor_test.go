// Tideland Go Stew - Actor - Unit Tests
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/actor"
)

//--------------------
// TESTS
//--------------------

// TestPureOK verifies the starting and stopping an Actor.
func TestPureOK(t *testing.T) {
	finalized := make(chan struct{})
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		defer close(finalized)
		return err
	}))

	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	act.Stop()

	<-finalized

	Assert(t, NoError(act.Err()), "actor stopped")
}

// TestPureDoubleStop verifies stopping an Actor twice.
func TestPureDoubleStop(t *testing.T) {
	finalized := make(chan struct{})
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		defer close(finalized)
		return err
	}))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	act.Stop()
	act.Stop()

	<-finalized

	Assert(t, NoError(act.Err()), "actor stopped")
}

// TestPureError verifies starting and stopping an Actor.
// Returning the stop error.
func TestPureError(t *testing.T) {
	finalized := make(chan struct{})
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		defer close(finalized)
		return errors.New("damn")
	}))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	act.Stop()

	<-finalized

	Assert(t, ErrorMatches(act.Err(), "damn"), "actor stopped with error")
}

// TestContext verifies starting and stopping an Actor
// with an external context.
func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	act, err := actor.Go(actor.WithContext(ctx))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	cancel()

	Assert(t, NoError(act.Err()), "actor stopped via context")
}

// TestSync verifies synchronous calls.
func TestSync(t *testing.T) {
	finalized := make(chan struct{})
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		defer close(finalized)
		return err
	}))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	counter := 0

	for i := 0; i < 5; i++ {
		err = act.DoSync(func() {
			counter++
		})
		Assert(t, NoError(err), "action done")
	}

	Assert(t, Equal(counter, 5), "counter is 5")

	act.Stop()

	<-finalized

	err = act.DoSync(func() {
		counter++
	})
	Assert(t, ErrorMatches(err, "actor is done"), "actor stopped, cannot use it anymore")
}

// TestTimeout verifies timout error of a synchronous Action.
func TestTimeout(t *testing.T) {
	act, err := actor.Go()
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	// Scenario: Timeout is shorter than needed time, so error
	// is returned.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	err = act.DoSyncWithContext(ctx, func() {
		time.Sleep(25 * time.Millisecond)
	})
	Assert(t, NoError(err), "action synchronous with context done")

	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	err = act.DoSyncWithContext(ctx, func() {
		time.Sleep(100 * time.Millisecond)
	})
	Assert(t, ErrorMatches(err, "action.*context deadline exceeded.*"), "action synchronous with context timeout")

	cancel()

	time.Sleep(150 * time.Millisecond)
	act.Stop()
}

// TestWithTimeoutContext verifies timout error of a synchronous Action
// when the Actor is configured with a context timeout.
func TestWithTimeoutContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	act, err := actor.Go(actor.WithContext(ctx))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	// Scenario: Configured timeout is shorter than needed
	// time, so error is returned.
	err = act.DoSync(func() {
		time.Sleep(10 * time.Millisecond)
	})
	Assert(t, NoError(err), "action synchronous done")

	err = act.DoSync(func() {
		time.Sleep(100 * time.Millisecond)
	})
	Assert(t, ErrorMatches(err, "actor.*context deadline exceeded.*"), "action synchronous timeout")

	act.Stop()
	cancel()
}

// TestAsyncWithQueueCap tests running multiple calls asynchronously.
func TestAsyncWithQueueCap(t *testing.T) {
	act, err := actor.Go(actor.WithQueueCap(128))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	sigs := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	// Start background func waiting for the signals of
	// the asynchrounous calls.
	go func() {
		count := 0
		for range sigs {
			count++
			if count == 128 {
				break
			}
		}
		close(done)
	}()

	// Now start asynchrounous calls.
	now := time.Now()
	for i := 0; i < 128; i++ {
		err = act.DoAsync(func() {
			time.Sleep(2 * time.Millisecond)
			sigs <- struct{}{}
		})
		Assert(t, NoError(err), "action asynchronous done")
	}
	enqueued := time.Since(now)

	// Expect signal done to be sent about one second later.
	<-done
	duration := time.Since(now)

	Assert(t, OK((duration-250*time.Millisecond) > enqueued), "duration is less than 2 seconds")

	act.Stop()
}

// TestRecovererOK tests successful handling of panic recoveries.
func TestNotifierOK(t *testing.T) {
	counter := 0
	recovered := false
	done := make(chan struct{})
	recoverer := func(reason any) error {
		defer close(done)
		recovered = true
		return nil
	}
	act, err := actor.Go(actor.WithRecoverer(recoverer))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	act.DoSync(func() {
		counter++
		// Will crash on first call.
		fmt.Printf("%v", counter/(counter-1))
	})
	<-done
	Assert(t, True(recovered), "recovered")

	err = act.DoSync(func() {
		counter++
	})
	Assert(t, NoError(err), "action synchronous done")
	Assert(t, Equal(counter, 2), "counter is 2")

	act.Stop()
}

// TestRecovererFail tests failing handling of panic recoveries.
func TestNotifierFail(t *testing.T) {
	counter := 0
	recovered := false
	done := make(chan struct{})
	recoverer := func(reason any) error {
		defer close(done)
		recovered = true
		return fmt.Errorf("ouch: %v", reason)
	}
	act, err := actor.Go(actor.WithRecoverer(recoverer))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	act.DoSync(func() {
		counter++
		// Will crash on first call.
		fmt.Printf("%v", counter/(counter-1))
	})
	<-done
	Assert(t, True(recovered), "recovered")

	Assert(t, True(act.IsDone()), "actor is done")
	Assert(t, ErrorMatches(act.Err(), "ouch:.*"), "actor stopped with error")
}

// EOF
