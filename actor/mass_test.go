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
	"math/rand"
	"testing"
	"time"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/actor"
)

//--------------------
// TESTS
//--------------------

// TestMass verifies the starting and stopping an Actor.
func TestMass(t *testing.T) {
	pps := make([]*PingPong, 1000)
	for i := 0; i < len(pps); i++ {
		pps[i] = NewPingPong(pps)
	}
	// Let's start the ping pong party.
	for i := 0; i < 5; i++ {
		n := rand.Intn(len(pps))
		pps[n].Ping()
		n = rand.Intn(len(pps))
		pps[n].Pong()
	}
	// Let's wait one seconds before stopping.
	time.Sleep(1 * time.Second)
	// Let's check some random ping pong pairs.
	for _, pp := range pps {
		pings, pongs := pp.PingPongs()
		Assert(t, True(pings > 0), "pings > 0")
		Assert(t, True(pongs > 0), "pongs > 0")
		pp.Stop()
	}
}

// TestPerformance verifies the starting and stopping an Actor.
func TestPerformance(t *testing.T) {
	finalized := make(chan struct{})
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		defer close(finalized)
		return err
	}))
	Assert(t, NoError(err), "actor started")
	Assert(t, NotNil(act), "actor not nil")

	now := time.Now()
	for i := 0; i < 10000; i++ {
		act.DoAsync(func() {})
	}
	duration := time.Since(now)
	Assert(t, True(duration < 100*time.Millisecond), "duration < 100ms")

	act.Stop()

	<-finalized

	Assert(t, NoError(act.Err()), "actor stopped via context")
}

//--------------------
// TEST ACTOR
//--------------------

type PingPong struct {
	pps   []*PingPong
	pings int
	pongs int

	act *actor.Actor
}

func NewPingPong(pps []*PingPong) *PingPong {
	pp := &PingPong{
		pps:   pps,
		pings: 0,
		pongs: 0,
	}
	act, err := actor.Go(actor.WithQueueCap(256))
	if err != nil {
		panic(err)
	}
	pp.act = act
	return pp
}

func (pp *PingPong) Ping() {
	pp.act.DoAsync(func() {
		pp.pings++
		n := rand.Intn(len(pp.pps))
		pp.pps[n].Pong()
	})
}

func (pp *PingPong) Pong() {
	pp.act.DoAsync(func() {
		pp.pongs++
		n := rand.Intn(len(pp.pps))
		pp.pps[n].Ping()
	})
}

func (pp *PingPong) PingPongs() (int, int) {
	var pings int
	var pongs int
	pp.act.DoSync(func() {
		pings = pp.pings
		pongs = pp.pongs
	})
	return pings, pongs
}

func (pp *PingPong) Err() error {
	return pp.act.Err()
}

func (pp *PingPong) Stop() {
	pp.act.Stop()
}

// EOF
