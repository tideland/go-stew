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
	"testing"
	"time"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/timex"
)

//--------------------
// TESTS
//--------------------

// TestTimeContainments verifies the containment of a time in a
// list or range of times.
func TestTimeContainments(t *testing.T) {
	ts := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	years := []int{2008, 2009, 2010}
	months := []time.Month{10, 11, 12}
	days := []int{10, 11, 12, 13, 14}
	hours := []int{20, 21, 22, 23}
	minutes := []int{0, 5, 10, 15, 20, 25}
	seconds := []int{0, 15, 30, 45}
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday}

	Assert(t, True(timex.YearInList(ts, years)), "time in year list")
	Assert(t, True(timex.YearInRange(ts, 2005, 2015)), "time in year range")
	Assert(t, True(timex.MonthInList(ts, months)), "time in month list")
	Assert(t, True(timex.MonthInRange(ts, 7, 12)), "time in month range")
	Assert(t, True(timex.DayInList(ts, days)), "time in day list")
	Assert(t, True(timex.DayInRange(ts, 5, 15)), "time in day range")
	Assert(t, True(timex.HourInList(ts, hours)), "time in hour list")
	Assert(t, True(timex.HourInRange(ts, 20, 31)), "time in hour range")
	Assert(t, True(timex.MinuteInList(ts, minutes)), "time in minute list")
	Assert(t, True(timex.MinuteInRange(ts, 0, 5)), "time in minute range")
	Assert(t, True(timex.SecondInList(ts, seconds)), "time in second list")
	Assert(t, True(timex.SecondInRange(ts, 0, 5)), "time in second range")
	Assert(t, True(timex.WeekdayInList(ts, weekdays)), "time in weekday list")
	Assert(t, True(timex.WeekdayInRange(ts, time.Monday, time.Friday)), "time in weekday range")
}

// TestBeginOf verifies the calculation of a beginning of a unit of time.
func TestBeginOf(t *testing.T) {
	ts := time.Date(2015, time.August, 2, 15, 10, 45, 12345, time.UTC)

	Assert(t, Equal(timex.BeginOf(ts, timex.Second), time.Date(2015, time.August, 2, 15, 10, 45, 0, time.UTC)), "beginning of second")
	Assert(t, Equal(timex.BeginOf(ts, timex.Minute), time.Date(2015, time.August, 2, 15, 10, 0, 0, time.UTC)), "beginning of minute")
	Assert(t, Equal(timex.BeginOf(ts, timex.Hour), time.Date(2015, time.August, 2, 15, 0, 0, 0, time.UTC)), "beginning of hour")
	Assert(t, Equal(timex.BeginOf(ts, timex.Day), time.Date(2015, time.August, 2, 0, 0, 0, 0, time.UTC)), "beginning of day")
	Assert(t, Equal(timex.BeginOf(ts, timex.Month), time.Date(2015, time.August, 1, 0, 0, 0, 0, time.UTC)), "beginning of month")
	Assert(t, Equal(timex.BeginOf(ts, timex.Year), time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)), "beginning of year")
}

// TestEndOf verifies the calculation of a ending of a unit of time.
func TestEndOf(t *testing.T) {
	ts := time.Date(2012, time.February, 2, 15, 10, 45, 12345, time.UTC)

	Assert(t, Equal(timex.EndOf(ts, timex.Second), time.Date(2012, time.February, 2, 15, 10, 45, 999999999, time.UTC)), "end of second")
	Assert(t, Equal(timex.EndOf(ts, timex.Minute), time.Date(2012, time.February, 2, 15, 10, 59, 999999999, time.UTC)), "end of minute")
	Assert(t, Equal(timex.EndOf(ts, timex.Hour), time.Date(2012, time.February, 2, 15, 59, 59, 999999999, time.UTC)), "end of hour")
	Assert(t, Equal(timex.EndOf(ts, timex.Day), time.Date(2012, time.February, 2, 23, 59, 59, 999999999, time.UTC)), "end of day")
	Assert(t, Equal(timex.EndOf(ts, timex.Month), time.Date(2012, time.February, 29, 23, 59, 59, 999999999, time.UTC)), "end of month")
	Assert(t, Equal(timex.EndOf(ts, timex.Year), time.Date(2012, time.December, 31, 23, 59, 59, 999999999, time.UTC)), "end of year")
}

// EOF
