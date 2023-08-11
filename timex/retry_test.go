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
	"errors"
	"testing"
	"time"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/timex"
)

//--------------------
// TESTS
//--------------------

// TestRetrySuccess verifies a successful retry.
func TestRetrySuccess(t *testing.T) {
	count := 0
	err := timex.Retry(func() (bool, error) {
		count++
		return count == 5, nil
	}, timex.ShortAttempt())
	Assert(t, NoError(err), "no error")
	Assert(t, Equal(count, 5), "retry executed five times")
}

// TestRetryFuncError verifies a retry with a function error.
func TestRetryFuncError(t *testing.T) {
	err := timex.Retry(func() (bool, error) {
		return false, errors.New("ouch")
	}, timex.ShortAttempt())
	Assert(t, ErrorContains(err, "ouch"), "error matches")
}

// TestRetryTooLong verifies a retry timeout error.
func TestRetryTooLong(t *testing.T) {
	rs := timex.RetryStrategy{
		Count:          10,
		Break:          5 * time.Millisecond,
		BreakIncrement: 5 * time.Millisecond,
		Timeout:        50 * time.Millisecond,
	}
	err := timex.Retry(func() (bool, error) {
		return false, nil
	}, rs)
	Assert(t, ErrorContains(err, "retried longer than"), "error matches")
}

// TestRetryTooOften verifies a retry count error.
func TestRetryTooOften(t *testing.T) {
	rs := timex.RetryStrategy{
		Count:          5,
		Break:          5 * time.Millisecond,
		BreakIncrement: 5 * time.Millisecond,
		Timeout:        time.Second,
	}
	err := timex.Retry(func() (bool, error) {
		return false, nil
	}, rs)
	Assert(t, ErrorContains(err, "retried more than"), "error matches")
}

// EOF
