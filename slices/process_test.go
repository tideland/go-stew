// Tideland Go Stew - Slices - Unit Tests
//
// Copyright (C) 2022-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// MatchesAll rights reserved. Use of this source code is governed
// by the new BSD license.

package slices_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/slices"
)

//--------------------
// TESTS
//--------------------

// TestFoldL verifies the left folding of a slice.
func TestFoldL(t *testing.T) {
	acc := "0"
	stringer := func(v int, acc string) string { return fmt.Sprintf("%s%d", acc, v) }
	tests := []struct {
		descr  string
		values []int
		out    string
	}{
		{
			descr:  "Many value slice",
			values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			out:    "0123456789",
		}, {
			descr:  "Single value slice",
			values: []int{1},
			out:    "01",
		}, {
			descr:  "Empty slice",
			values: []int{},
			out:    "0",
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    "0",
		},
	}

	for _, test := range tests {
		Assert(t, Equal(slices.FoldL(test.values, acc, stringer), test.out), test.descr)
	}
}

// TestFoldLFirst verifies the left folding of a slice with first as
// accumulator.
func TestFoldLFirst(t *testing.T) {
	potentiator := func(v, acc int) int { return acc*10 + v }
	tests := []struct {
		descr  string
		values []int
		out    int
	}{
		{
			descr:  "Many values slice",
			values: []int{1, 2, 3, 4, 5},
			out:    12345,
		}, {
			descr:  "Single value slice",
			values: []int{1},
			out:    1,
		}, {
			descr:  "Empty slice",
			values: []int{},
			out:    0,
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    0,
		},
	}

	for _, test := range tests {
		Assert(t, Equal(slices.FoldLFirst(test.values, potentiator), test.out), test.descr)
	}
}

// TestFoldR verifies the right folding of a slice.
func TestFoldR(t *testing.T) {
	acc := "0"
	stringer := func(v int, acc string) string { return fmt.Sprintf("%s%d", acc, v) }
	tests := []struct {
		descr  string
		values []int
		out    string
	}{
		{
			descr:  "Many value slice",
			values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			out:    "0987654321",
		}, {
			descr:  "Single value slice",
			values: []int{1},
			out:    "01",
		}, {
			descr:  "Empty slice",
			values: []int{},
			out:    "0",
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    "0",
		},
	}

	for _, test := range tests {
		Assert(t, Equal(slices.FoldR(test.values, acc, stringer), test.out), test.descr)
	}
}

// TestFoldRLast verifies the right folding of a slice.
func TestFoldRLast(t *testing.T) {
	potentiator := func(v, acc int) int { return acc*10 + v }
	tests := []struct {
		descr  string
		values []int
		out    int
	}{
		{
			descr:  "Many values slice",
			values: []int{1, 2, 3, 4, 5},
			out:    54321,
		}, {
			descr:  "Single value slice",
			values: []int{1},
			out:    1,
		}, {
			descr:  "Empty slice",
			values: []int{},
			out:    0,
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    0,
		},
	}

	for _, test := range tests {
		Assert(t, Equal(slices.FoldRLast(test.values, potentiator), test.out), test.descr)
	}
}

// TestMapFoldL verifies the left combined mapping and folding.
func TestMapFoldL(t *testing.T) {
	in := "0"
	plusStringer := func(v int, acc string) (int, string) { return v + 1, fmt.Sprintf("%s%d", acc, v) }
	tests := []struct {
		descr  string
		values []int
		mapped []int
		out    string
	}{
		{
			descr:  "Many values slice",
			values: []int{1, 2, 3, 4, 5},
			mapped: []int{2, 3, 4, 5, 6},
			out:    "012345",
		}, {
			descr:  "Single value slice",
			values: []int{1},
			mapped: []int{2},
			out:    "01",
		}, {
			descr:  "Empty slice",
			values: []int{},
			mapped: []int{},
			out:    "0",
		}, {
			descr:  "Nil slice",
			values: nil,
			mapped: nil,
			out:    "0",
		},
	}

	for _, test := range tests {
		mapped, out := slices.MapFoldL(test.values, in, plusStringer)
		Assert(t, DeepEqual(mapped, test.mapped), test.descr)
		Assert(t, Equal(out, test.out), test.descr)
	}
}

// TestMapFoldR verifies the right combined mapping and folding.
func TestMapFoldR(t *testing.T) {
	in := "0"
	plusStringer := func(v int, acc string) (int, string) { return v + 1, fmt.Sprintf("%s%d", acc, v) }
	tests := []struct {
		descr  string
		values []int
		mapped []int
		out    string
	}{
		{
			descr:  "Many values slice",
			values: []int{1, 2, 3, 4, 5},
			mapped: []int{2, 3, 4, 5, 6},
			out:    "054321",
		}, {
			descr:  "Single value slice",
			values: []int{1},
			mapped: []int{2},
			out:    "01",
		}, {
			descr:  "Empty slice",
			values: []int{},
			mapped: []int{},
			out:    "0",
		}, {
			descr:  "Nil slice",
			values: nil,
			mapped: nil,
			out:    "0",
		},
	}

	for _, test := range tests {
		mapped, out := slices.MapFoldR(test.values, in, plusStringer)
		Assert(t, DeepEqual(mapped, test.mapped), test.descr)
		Assert(t, Equal(out, test.out), test.descr)
	}
}

// TestPartition verifies the partitioning of slices.
func TestPartition(t *testing.T) {
	isMod := func(v int) bool { return v%2 == 0 }
	tests := []struct {
		descr         string
		values        []int
		satisfying    []int
		notSatisfying []int
	}{
		{
			descr:         "Many values slice containing partiion matching values",
			values:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			satisfying:    []int{2, 4, 6, 8},
			notSatisfying: []int{1, 3, 5, 7, 9},
		}, {
			descr:         "Many values slice containing no partiion matching values",
			values:        []int{2, 4, 6, 8, 10},
			satisfying:    []int{2, 4, 6, 8, 10},
			notSatisfying: nil,
		}, {
			descr:         "Many values slice containing only partiion matching values",
			values:        []int{1, 3, 5, 7, 9},
			satisfying:    nil,
			notSatisfying: []int{1, 3, 5, 7, 9},
		}, {
			descr:         "Single value slice",
			values:        []int{1},
			satisfying:    nil,
			notSatisfying: []int{1},
		}, {
			descr:         "Empty slice",
			values:        []int{},
			satisfying:    nil,
			notSatisfying: nil,
		}, {
			descr:         "Nil slice",
			values:        nil,
			satisfying:    nil,
			notSatisfying: nil,
		},
	}

	for _, test := range tests {
		satisfying, notSatisfying := slices.Partition(test.values, isMod)
		Assert(t, DeepEqual(satisfying, test.satisfying), test.descr)
		Assert(t, DeepEqual(notSatisfying, test.notSatisfying), test.descr)
	}
}

// EOF
