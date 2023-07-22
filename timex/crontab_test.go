// Tideland Go Library - Time Extensions - Unit Tests
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package timex_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/timex"
)

//--------------------
// TESTS
//--------------------

// TestCrontabEndlessJob verifies the endless execution of a crontab job.
func TestCrontabEndlessJob(t *testing.T) {
	interval := 100 * time.Millisecond
	counter := 0
	job := func() (bool, error) {
		counter++
		return true, nil
	}
	ctx := context.Background()
	ct, err := timex.NewCrontab(ctx, interval)
	Assert(t, NoError(err), "no error creating crontab")
	Assert(t, NotNil(ct), "crontab created")
	Assert(t, True(ct.Status() == timex.StatusWorking), "crontab is working")
	defer ct.Stop()

	ct.Add("job", interval, job)
	time.Sleep(5 * interval)
	Assert(t, Range(counter, 3, 6), "job executed at three to six times")
}

// TestCrontabStoppingJob verifies the stopping of a crontab job internally.
func TestCrontabStoppingJob(t *testing.T) {
	interval := 100 * time.Millisecond
	counter := 0
	job := func() (bool, error) {
		counter++
		if counter == 2 {
			return false, nil
		}
		return true, nil
	}
	ctx := context.Background()
	ct, err := timex.NewCrontab(ctx, interval)
	Assert(t, NoError(err), "no error creating crontab")
	Assert(t, NotNil(ct), "crontab created")
	Assert(t, True(ct.Status() == timex.StatusWorking), "crontab is working")
	defer ct.Stop()

	ct.Add("job", interval, job)
	time.Sleep(5 * interval)
	Assert(t, Equal(counter, 2), "job executed twice")

	ok, err := ct.JobStatus("job")

	Assert(t, False(ok), "job not found")
	Assert(t, NoError(err), "no error getting job status")
}

// TestCrontabRemovedJob verifies the stopping of a crontab job externally.
func TestCrontabRemovedJob(t *testing.T) {
	interval := 100 * time.Millisecond
	active := false
	job := func() (bool, error) {
		active = true
		return true, nil
	}
	ctx := context.Background()
	ct, err := timex.NewCrontab(ctx, interval)
	Assert(t, NoError(err), "no error creating crontab")
	Assert(t, NotNil(ct), "crontab created")
	Assert(t, True(ct.Status() == timex.StatusWorking), "crontab is working")
	defer ct.Stop()

	ct.Add("job", interval, job)
	time.Sleep(5 * interval)
	Assert(t, True(active), "job had been active")
	ct.Remove("job")
	time.Sleep(1 * interval)
	ok, err := ct.JobStatus("job")

	Assert(t, False(ok), "job not found")
	Assert(t, NoError(err), "no error getting job status")
}

// TestCrontabErrorJob verifies the stopping of a crontab job due to an error.
func TestCrontabErrorJob(t *testing.T) {
	interval := 100 * time.Millisecond
	counter := 0
	job := func() (bool, error) {
		counter++
		if counter == 2 {
			return false, fmt.Errorf("ouch")
		}
		return true, nil
	}
	ctx := context.Background()
	ct, err := timex.NewCrontab(ctx, interval)
	Assert(t, NoError(err), "no error creating crontab")
	Assert(t, NotNil(ct), "crontab created")
	Assert(t, True(ct.Status() == timex.StatusWorking), "crontab is working")
	defer ct.Stop()

	ct.Add("job", interval, job)
	time.Sleep(5 * interval)
	Assert(t, Equal(counter, 2), "job executed twice")

	ok, err := ct.JobStatus("job")

	Assert(t, False(ok), "job not found")
	Assert(t, ErrorContains(err, "ouch"), "error matches	")
}

// TestCrontabDifferentFrequencies verifies the execution of multiple jobs
// with different frequencies.
func TestCronjobDifferentFrequencies(t *testing.T) {
	interval := 100 * time.Millisecond
	counters := make([]int, 3)
	mkjob := func(i int) func() (bool, error) {
		return func() (bool, error) {
			counters[i]++
			return true, nil
		}
	}
	ctx := context.Background()
	ct, err := timex.NewCrontab(ctx, interval)
	Assert(t, NoError(err), "no error creating crontab")
	Assert(t, NotNil(ct), "crontab created")
	Assert(t, True(ct.Status() == timex.StatusWorking), "crontab is working")
	defer ct.Stop()

	ct.Add("job-0", interval/2, mkjob(0))
	ct.Add("job-1", interval, mkjob(1))
	ct.Add("job-2", interval*2, mkjob(2))
	time.Sleep(5 * interval)
	Assert(t, Equal(counters[0], counters[1]), "too low frequency of job-0 has been increased")
	Assert(t, True(counters[2] < counters[1]), "high frequency of job-2 has fewer executions")
}

// EOF
