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
	"runtime"
	"testing"

	"tideland.dev/go/stew/qagen"
	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/slices"
)

//--------------------
// TESTS
//--------------------

// TestHeadTailInitLast verifies the splitting of slices.
func TestHeadTailInitLast(t *testing.T) {
	tests := []struct {
		descr  string
		values []int
		head   int
		tail   []int
		init   []int
		last   int
	}{
		{
			descr:  "Simple slice",
			values: []int{1, 2, 3, 4, 5},
			head:   1,
			tail:   []int{2, 3, 4, 5},
			init:   []int{1, 2, 3, 4},
			last:   5,
		}, {
			descr:  "Single value slice",
			values: []int{1},
			head:   1,
			tail:   []int{},
			init:   []int{},
			last:   1,
		}, {
			descr:  "Empty slice",
			values: []int{},
			head:   0,
			tail:   []int{},
			init:   []int{},
			last:   0,
		}, {
			descr:  "Nil slice",
			values: nil,
			head:   0,
			tail:   nil,
			init:   nil,
			last:   0,
		},
	}

	for _, test := range tests {
		t.Logf("test: %s", test.descr)

		head, tail := slices.HeadTail(test.values)
		init, last := slices.InitLast(test.values)

		Assert(t, Equal(head, test.head), test.descr)
		Assert(t, DeepEqual(tail, test.tail), test.descr)
		Assert(t, DeepEqual(init, test.init), test.descr)
		Assert(t, Equal(last, test.last), test.descr)
	}
}

// TestSort verifies the standard sorting of slices.
func TestSort(t *testing.T) {
	tests := []struct {
		descr  string
		values []int
		out    []int
	}{
		{
			descr:  "Simple unordered slice",
			values: []int{5, 7, 1, 3, 4, 2, 8, 6, 9},
			out:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}, {
			descr:  "Unordered double value slice",
			values: []int{9, 5, 7, 3, 1, 3, 5, 4, 2, 8, 5, 6, 9},
			out:    []int{1, 2, 3, 3, 4, 5, 5, 5, 6, 7, 8, 9, 9},
		}, {
			descr:  "Already ordered slice",
			values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			out:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}, {
			descr:  "Reverse ordered slice",
			values: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
			out:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}, {
			descr:  "Single value slice",
			values: []int{1, 1, 1, 1, 1},
			out:    []int{1, 1, 1, 1, 1},
		}, {
			descr:  "Empty slice",
			values: []int{},
			out:    []int{},
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    nil,
		},
	}

	for _, test := range tests {
		Assert(t, DeepEqual(slices.Sort(test.values), test.out), test.descr)
	}
}

// TestSortWith verifies the sorting of slices with a less function.
func TestSortWith(t *testing.T) {
	less := func(vs []string, i, j int) bool { return len(vs[i]) < len(vs[j]) }
	tests := []struct {
		descr  string
		values []string
		out    []string
	}{
		{
			descr:  "Simple unordered slice",
			values: []string{"alpha", "beta", "phi", "epsilon", "lambda", "pi"},
			out:    []string{"pi", "phi", "beta", "alpha", "lambda", "epsilon"},
		}, {
			descr:  "Unordered double value slice",
			values: []string{"phi", "alpha", "beta", "phi", "epsilon", "beta", "lambda", "pi"},
			out:    []string{"pi", "phi", "phi", "beta", "beta", "alpha", "lambda", "epsilon"},
		}, {
			descr:  "Already ordered slice",
			values: []string{"pi", "phi", "beta", "alpha", "lambda", "epsilon"},
			out:    []string{"pi", "phi", "beta", "alpha", "lambda", "epsilon"},
		}, {
			descr:  "Reverse ordered slice",
			values: []string{"epsilon", "lambda", "alpha", "beta", "phi", "pi"},
			out:    []string{"pi", "phi", "beta", "alpha", "lambda", "epsilon"},
		}, {
			descr:  "Single value slice",
			values: []string{"alpha", "alpha", "alpha", "alpha", "alpha"},
			out:    []string{"alpha", "alpha", "alpha", "alpha", "alpha"},
		}, {
			descr:  "Empty slice",
			values: []string{},
			out:    []string{},
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    nil,
		},
	}

	for _, test := range tests {
		Assert(t, DeepEqual(slices.SortWith(test.values, less), test.out), test.descr)
	}
}

// TestIsSorted verifies the check of sorted slices.
func TestIsSorted(t *testing.T) {
	tests := []struct {
		descr  string
		values []int
		out    bool
	}{
		{
			descr:  "Unordered slice",
			values: []int{5, 7, 1, 3, 4, 2, 8, 6, 9},
			out:    false,
		}, {
			descr:  "Ordered slice",
			values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			out:    true,
		}, {
			descr:  "Reverse ordered slice",
			values: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
			out:    false,
		}, {
			descr:  "Single value slice",
			values: []int{1, 1, 1, 1, 1},
			out:    true,
		}, {
			descr:  "Empty slice",
			values: []int{},
			out:    true,
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    true,
		},
	}

	for _, test := range tests {
		Assert(t, Equal(slices.IsSorted(test.values), test.out), test.descr)
	}
}

// TestLargeSort verifies the sorting of large slices with a parallel QuickSort.
func TestLargeSort(t *testing.T) {
	size := runtime.NumCPU()*2048 + 1
	gen := qagen.New(qagen.FixedRand())
	ivs := gen.Ints(0, 10000, size)

	Assert(t, False(slices.IsSorted(ivs)), "unsorted slice")
	ovs := slices.Sort(ivs)
	Assert(t, True(slices.IsSorted(ovs)), "sorted slice")
}

// TestIsSortedWith verifies the check of sorted slices.
func TestIsSortedWith(t *testing.T) {
	less := func(a, b string) bool { return len(a) < len(b) }
	tests := []struct {
		descr  string
		values []string
		out    bool
	}{
		{
			descr:  "Unordered slice",
			values: []string{"alpha", "beta", "phi", "epsilon", "lambda", "pi"},
			out:    false,
		}, {
			descr:  "Ordered slice",
			values: []string{"pi", "phi", "beta", "alpha", "lambda", "epsilon"},
			out:    true,
		}, {
			descr:  "Reverse ordered slice",
			values: []string{"epsilon", "lambda", "alpha", "beta", "phi", "pi"},
			out:    false,
		}, {
			descr:  "Single value slice",
			values: []string{"alpha", "alpha", "alpha", "alpha", "alpha"},
			out:    true,
		}, {
			descr:  "Empty slice",
			values: []string{},
			out:    true,
		}, {
			descr:  "Nil slice",
			values: nil,
			out:    true,
		},
	}

	for _, test := range tests {
		Assert(t, Equal(slices.IsSortedWith(test.values, less), test.out), test.descr)
	}
}

// TestShuffle verifies the random shuffling of slices.
func TestShuffle(t *testing.T) {
	tests := []struct {
		descr  string
		values []int
		sorted bool
	}{
		{
			descr:  "Unordered slice",
			values: []int{5, 7, 1, 3, 4, 2, 8, 6, 9},
			sorted: false,
		}, {
			descr:  "Ordered slice",
			values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			sorted: false,
		}, {
			descr:  "Single value slice",
			values: []int{1, 1, 1, 1, 1},
			sorted: true,
		}, {
			descr:  "Empty slice",
			values: []int{},
			sorted: true,
		}, {
			descr:  "Nil slice",
			values: nil,
			sorted: true,
		},
	}

	for _, test := range tests {
		shuffled := slices.Shuffle(test.values)
		sorted := slices.Sort(shuffled)

		Assert(t, Equal(len(shuffled), len(test.values)), test.descr)
		Assert(t, Equal(slices.IsSorted(shuffled), test.sorted), test.descr)
		Assert(t, Equal(slices.IsSorted(sorted), true), test.descr)
	}
}

//--------------------
// BENCHMARKS AND FUZZ TESTS
//--------------------

// BenchmarkSort runs a performance test on standard sorting.
func BenchmarkSort(b *testing.B) {
	gen := qagen.New(qagen.FixedRand())
	vs := gen.Ints(0, 1000, 10000)

	slices.Sort(vs)
}

// BenchmarkSortWith runs a performance test on sorting with comparator.
func BenchmarkSortWith(b *testing.B) {
	gen := qagen.New(qagen.FixedRand())
	vs := gen.Words(10000)
	less := func(vs []string, i, j int) bool { return len(vs[i]) < len(vs[j]) }

	slices.SortWith(vs, less)
}

// FuzzSort runs a fuzz test on the standard sorting.
func FuzzSort(f *testing.F) {
	gen := qagen.New(qagen.FixedRand())

	f.Add(5)
	f.Fuzz(func(t *testing.T, i int) {
		vs := gen.Ints(0, 1000, 10000)

		slices.Sort(vs)
	})
}

// EOF
