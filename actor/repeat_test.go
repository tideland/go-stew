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
	"testing"
	"time"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/actor"
)

//--------------------
// TESTS
//--------------------

// TestRepeatStopActor verifies Repeat working and being
// stopped when the Actor is stopped.
func TestRepeatStopActor(t *testing.T) {
	finalized := make(chan struct{})
	counter := 0
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		defer close(finalized)

		counter = 0

		return err
	}))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	// Start the repeated action.
	stop, err := act.Repeat(10*time.Millisecond, func() {
		counter++
	})
	Assert(t, NoError(err), "action repeated")
	Assert(t, NotNil(stop), "stop not nil")

	time.Sleep(100 * time.Millisecond)
	Assert(t, True(counter >= 9), "possibly only 9 due to late interval start")

	// Stop the Actor and check the finalization.
	act.Stop()

	<-finalized

	Assert(t, NoError(act.Err()), "actor stopped via context")
	Assert(t, Equal(counter, 0), "counter is 0")

	// Check if the Interval is stopped too.
	time.Sleep(100 * time.Millisecond)
	Assert(t, Equal(counter, 0), "counter is still 0")
}

// TestPeriodicalStopInterval verifies Periodical working and being
// stopped when the periodical is stopped.
func TestIntervalStopInterval(t *testing.T) {
	counter := 0
	act, err := actor.Go()
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	// Start the repeated action.
	stop, err := act.Repeat(10*time.Millisecond, func() {
		counter++
	})
	Assert(t, NoError(err), "action repeated")
	Assert(t, NotNil(stop), "stop not nil")

	time.Sleep(100 * time.Millisecond)
	Assert(t, True(counter >= 9), "possibly only 9 due to late interval start")

	// Stop the periodical and check that it doesn't work anymore.
	counterNow := counter
	stop()

	time.Sleep(100 * time.Millisecond)
	Assert(t, Equal(counter, counterNow), "counter not increased")

	act.Stop()
}

// EOF
