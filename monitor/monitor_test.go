// Tideland Go Stew - Monitor - Unit Tests
//
// Copyright (C) 2009-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package monitor_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/stew/qagen"
	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/monitor"
)

//--------------------
// TESTS
//--------------------

// TestSimpleMonitor test creating and stopping a monitor.
func TestSimpleMonitor(t *testing.T) {
	m := monitor.New()
	defer m.Stop()

	measures := m.StopWatch().Measure("simple", func() { time.Sleep(time.Millisecond) })
	Assert(t, True(measures > 0), "stop watch is not nil")

	mp, err := m.StopWatch().Read("simple")
	Assert(t, NoError(err), "stop watch read works")
	Assert(t, Equal(mp.ID, "simple"), "stop watch id is correct")
}

// TestStopWatch tests the stop watch.
func TestStopWatch(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())
	m := monitor.New()
	defer m.Stop()

	// Generate some measurings first.
	for i := 0; i < 500; i++ {
		m.StopWatch().Measure("watch", func() {
			gen.SleepOneOf(1*time.Millisecond, 3*time.Millisecond, 10*time.Millisecond)
		})
	}

	sw, err := m.StopWatch().Read("doesnotexist")
	Assert(t, ErrorMatches(err, `watch value 'doesnotexist' does not exist`), "error for non-existing watch")

	// Check access of one measuring point.
	sw, err = m.StopWatch().Read("watch")
	Assert(t, NoError(err), "no error for existing watch")
	Assert(t, Equal(sw.ID, "watch"), "watch id is correct")
	Assert(t, Equal(sw.Count, 500), "watch count is correct")
	Assert(t, Range(sw.Avg, sw.Min, sw.Max), "watch avg is in range")

	// Check iteration over all measuring points.
	wvs := monitor.WatchValues{}
	err = m.StopWatch().Do(func(wv monitor.WatchValue) error {
		wvs = append(wvs, wv)
		return nil
	})
	Assert(t, NoError(err), "no error for iteration")
	Assert(t, Length(wvs, 1), "one watch value")

	// Check resetting the measurings.
	m.Reset()

	wvs = monitor.WatchValues{}
	err = m.StopWatch().Do(func(wv monitor.WatchValue) error {
		wvs = append(wvs, wv)
		return nil
	})
	Assert(t, NoError(err), "no error for iteration")
	Assert(t, Empty(wvs), "no watch values")
}

// Test of the stay-set indicators  of the monitor.
func TestStaySetIndicators(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())
	m := monitor.New()
	defer m.Stop()

	// Generate some indicators first.
	for i := 0; i < 500; i++ {
		id := gen.OneStringOf("foo", "bar", "baz", "yadda", "deadbeef")
		if gen.FlipCoin(60) {
			m.StaySetIndicator().Increase(id)
		} else {
			m.StaySetIndicator().Decrease(id)
		}
	}

	iv, err := m.StaySetIndicator().Read("doesnotexist")
	Assert(t, ErrorMatches(err, `indicator value 'doesnotexist' does not exist`), "error for non-existing indicator")

	// Check access of one stay-set indicator.
	iv, err = m.StaySetIndicator().Read("foo")
	Assert(t, NoError(err), "no error for existing indicator")
	Assert(t, Equal(iv.ID, "foo"), "indicator id is correct")
	Assert(t, Equal(iv.Count, 99), "indicator count is correct")
	Assert(t, Range(iv.Current, iv.Min, iv.Max), "indicator current is in range")

	// Check iteration over all measuring points.
	ivs := monitor.IndicatorValues{}
	err = m.StaySetIndicator().Do(func(iv monitor.IndicatorValue) error {
		ivs = append(ivs, iv)
		return nil
	})
	Assert(t, NoError(err), "no error for iteration")
	Assert(t, Length(ivs, 5), "five indicator values")

	// Check resetting the measurings.
	m.Reset()

	ivs = monitor.IndicatorValues{}
	err = m.StaySetIndicator().Do(func(iv monitor.IndicatorValue) error {
		ivs = append(ivs, iv)
		return nil
	})
	Assert(t, NoError(err), "no error for iteration")
	Assert(t, Empty(ivs), "no indicator values")
}

//--------------------
// BENCHMARKS
//--------------------

// BenchmarkMonitor checks the performance of monitor.
func BenchmarkMonitor(b *testing.B) {
	gen := qagen.New(qagen.SimpleRand())
	m := monitor.New()
	defer m.Stop()

	for i := 0; i < b.N; i++ {
		m.StopWatch().Measure("bench", func() {
			gen.SleepOneOf(1*time.Millisecond, 3*time.Millisecond, 5*time.Millisecond)
		})
	}
}

// EOF
