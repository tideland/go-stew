// Tideland Go Stew - QA Generators - Unit Tests
//
// Copyright (C) 2013-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the New BSD license.

package qagen_test

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"
	"time"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/qagen"
)

//--------------------
// TESTS
//--------------------

// TestBuildDate tests the generation of dates.
func TestBuildDate(t *testing.T) {
	layouts := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	for _, layout := range layouts {
		ts, tim := qagen.BuildTime(layout, 0)
		tsp, err := time.Parse(layout, ts)
		Assert(t, NoError(err), "time parsed")
		Assert(t, Equal(tim, tsp), "time equal")

		ts, tim = qagen.BuildTime(layout, -30*time.Minute)
		tsp, err = time.Parse(layout, ts)
		Assert(t, NoError(err), "time parsed")
		Assert(t, Equal(tim, tsp), "time equal")

		ts, tim = qagen.BuildTime(layout, time.Hour)
		tsp, err = time.Parse(layout, ts)
		Assert(t, NoError(err), "time parsed")
		Assert(t, Equal(tim, tsp), "time equal")
	}
}

// TestBytes tests the generation of bytes.
func TestBytes(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	// Test individual bytes.
	for i := 0; i < 10000; i++ {
		lo := gen.Byte(0, 255)
		hi := gen.Byte(0, 255)
		n := gen.Byte(lo, hi)
		if hi < lo {
			lo, hi = hi, lo
		}
		Assert(t, Range(n, lo, hi), "byte range")
	}

	// Test byte slices.
	ns := gen.Bytes(1, 200, 1000)
	Assert(t, Length(ns, 1000), "length")
	for _, n := range ns {
		Assert(t, Range(n, 1, 200), "byte slice")
	}

	// Test UUIDs.
	for i := 0; i < 10000; i++ {
		uuid := gen.UUID()
		Assert(t, Length(uuid, 16), "UUID length")
	}
}

// TestInts tests the generation of ints.
func TestInts(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	// Test individual ints.
	for i := 0; i < 10000; i++ {
		lo := gen.Int(-100, 100)
		hi := gen.Int(-100, 100)
		n := gen.Int(lo, hi)
		if hi < lo {
			lo, hi = hi, lo
		}
		Assert(t, Range(n, lo, hi), "int range")
	}

	// Test int slices.
	ns := gen.Ints(0, 500, 10000)
	Assert(t, Length(ns, 10000), "length")
	for _, n := range ns {
		Assert(t, Range(n, 0, 500), "int slice")
	}

	// Test the generation of percent.
	for i := 0; i < 10000; i++ {
		p := gen.Percent()
		Assert(t, Range(p, 0, 100), "percent")
	}

	// Test the flipping of coins.
	ct := 0
	cf := 0
	for i := 0; i < 10000; i++ {
		c := gen.FlipCoin(50)
		if c {
			ct++
		} else {
			cf++
		}
	}
	Assert(t, About(ct, cf, 500), "coin flipping")
}

// TestOneOf tests the generation of selections.
func TestOneOf(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())
	stuff := []any{1, true, "three", 47.11, []byte{'A', 'B', 'C'}}

	for i := 0; i < 10000; i++ {
		m := gen.OneOf(stuff...)
		Assert(t, DeepContains(stuff, m), "contains")

		b := gen.OneByteOf(1, 2, 3, 4, 5)
		Assert(t, Range(b, 1, 5), "byte")

		r := gen.OneRuneOf("abcdef")
		Assert(t, Range(r, 'a', 'f'), "rune")

		n := gen.OneIntOf(1, 2, 3, 4, 5)
		Assert(t, Range(n, 1, 5), "int")

		ss := []string{"one", "two", "three", "four", "five"}
		s := gen.OneStringOf(ss...)
		find := func() bool {
			for _, os := range ss {
				if os == s {
					return true
				}
			}
			return false
		}
		Assert(t, True(find()), "string")

		d := gen.OneDurationOf(1*time.Second, 2*time.Second, 3*time.Second)
		Assert(t, Range(d, 1*time.Second, 3*time.Second), "duration")
	}
}

// TestWords tests the generation of words.
func TestWords(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	// Test single words.
	for i := 0; i < 10000; i++ {
		w := gen.Word()
		for _, r := range w {
			Assert(t, Range(r, 'a', 'z'), "rune")
		}
	}

	// Test limited words.
	for i := 0; i < 10000; i++ {
		lo := gen.Int(qagen.MinWordLen, qagen.MaxWordLen)
		hi := gen.Int(qagen.MinWordLen, qagen.MaxWordLen)
		w := gen.LimitedWord(lo, hi)
		wl := len(w)
		if hi < lo {
			lo, hi = hi, lo
		}
		Assert(t, Range(wl, lo, hi), "limited word length")
	}
}

// TestPattern tests the generation based on patterns.
func TestPattern(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())
	assertPattern := func(pattern, runes string) {
		set := make(map[rune]bool)
		for _, r := range runes {
			set[r] = true
		}
		for i := 0; i < 10; i++ {
			result := gen.Pattern(pattern)
			for _, r := range result {
				Assert(t, True(set[r]), "pattern %q contains %q", pattern, r)
			}
		}
	}

	assertPattern("^^", "^")
	assertPattern("^0^0^0^0^0", "0123456789")
	assertPattern("^1^1^1^1^1", "123456789")
	assertPattern("^o^o^o^o^o", "01234567")
	assertPattern("^h^h^h^h^h", "0123456789abcdef")
	assertPattern("^H^H^H^H^H", "0123456789ABCDEF")
	assertPattern("^a^a^a^a^a", "abcdefghijklmnopqrstuvwxyz")
	assertPattern("^A^A^A^A^A", "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	assertPattern("^c^c^c^c^c", "bcdfghjklmnpqrstvwxyz")
	assertPattern("^C^C^C^C^C", "BCDFGHJKLMNPQRSTVWXYZ")
	assertPattern("^v^v^v^v^v", "aeiou")
	assertPattern("^V^V^V^V^V", "AEIOU")
	assertPattern("^z^z^z^z^z", "abcdefghijklmnopqrstuvwxyz0123456789")
	assertPattern("^Z^Z^Z^Z^Z", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	assertPattern("^1^0.^0^0^0,^0^0 €", "0123456789 .,€")
}

// TestText tests the generation of text.
func TestText(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())
	names := gen.Names(4)

	for i := 0; i < 10000; i++ {
		s := gen.Sentence()
		ws := strings.Split(s, " ")
		lws := len(ws)
		Assert(t, Range(lws, 2, 15), "S: %v SL: %d", s, lws)
		Assert(t, Range(s[0], 'A', 'Z'), "SUC: %v", s[0])
	}

	for i := 0; i < 10; i++ {
		s := gen.SentenceWithNames(names)
		Assert(t, NotEmpty(s), "sentence with names")
	}

	for i := 0; i < 10000; i++ {
		p := gen.Paragraph()
		ss := strings.Split(p, ". ")
		lss := len(ss)
		Assert(t, Range(lss, 2, 10), "PL: %d", lss)
		for _, s := range ss {
			ws := strings.Split(s, " ")
			lws := len(ws)
			Assert(t, Range(lws, 2, 15), "S: %v PSL: %d", s, lws)
			Assert(t, Range(s[0], 'A', 'Z'), "PSUC: %v", s[0])
		}
	}

	for i := 0; i < 10; i++ {
		s := gen.ParagraphWithNames(names)
		Assert(t, NotEmpty(s), "paragraph with names")
	}
}

// TestName tests the generation of names.
func TestName(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	Assert(t, Equal(qagen.ToUpperFirst("yadda"), "Yadda"), "to upper first")

	for i := 0; i < 10000; i++ {
		first, middle, last := gen.Name()

		Assert(t, Matches(first, `[A-Z][a-z]+(-[A-Z][a-z]+)?`), "first name")
		Assert(t, Matches(middle, `[A-Z][a-z]+(-[A-Z][a-z]+)?`), "middle name")
		Assert(t, Matches(last, `[A-Z]['a-zA-Z]+`), "last name")

		first, middle, last = gen.MaleName()

		Assert(t, Matches(first, `[A-Z][a-z]+(-[A-Z][a-z]+)?`), "first name")
		Assert(t, Matches(middle, `[A-Z][a-z]+(-[A-Z][a-z]+)?`), "middle name")
		Assert(t, Matches(last, `[A-Z]['a-zA-Z]+`), "last name")

		first, middle, last = gen.FemaleName()

		Assert(t, Matches(first, `[A-Z][a-z]+(-[A-Z][a-z]+)?`), "first name")
		Assert(t, Matches(middle, `[A-Z][a-z]+(-[A-Z][a-z]+)?`), "middle name")
		Assert(t, Matches(last, `[A-Z]['a-zA-Z]+`), "last name")

		count := gen.Int(0, 5)
		names := gen.Names(count)

		Assert(t, Length(names, count), "length")

		for _, name := range names {
			Assert(t, Matches(name, `[A-Z][a-z]+(-[A-Z][a-z]+)?\s([A-Z]\.\s)?[A-Z]['a-zA-Z]+`), "name")
		}
	}
}

// TestDomain tests the generation of domains.
func TestDomain(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	for i := 0; i < 00100; i++ {
		domain := gen.Domain()

		Assert(t, Matches(domain, `[a-z0-9.-]+\.[a-z]{2,4}`), "domain")
	}
}

// TestURL tests the generation of URLs.
func TestURL(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	for i := 0; i < 10000; i++ {
		url := gen.URL()

		Assert(t, Matches(url, `(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`), "URL")
	}
}

// TestEMail tests the generation of e-mail addresses.
func TestEMail(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	for i := 0; i < 10000; i++ {
		addr := gen.EMail()

		Assert(t, Matches(addr, `[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}`), "e-mail address")
	}
}

// TestTimes tests the generation of durations and times.
func TestTimes(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	for i := 0; i < 10000; i++ {
		// Test durations.
		lo := gen.Duration(time.Second, time.Minute)
		hi := gen.Duration(time.Second, time.Minute)
		d := gen.Duration(lo, hi)
		if hi < lo {
			lo, hi = hi, lo
		}
		Assert(t, Range(d, lo, hi), "duration range")

		// Test times.
		loc := time.Local
		now := time.Now()
		dur := gen.Duration(24*time.Hour, 30*24*time.Hour)
		tim := gen.Time(loc, now, dur)
		Assert(t, True(tim.Equal(now) || tim.After(now)), "equal or after now")
		Assert(t, True(tim.Before(now.Add(dur)) || tim.Equal(now.Add(dur))), "before or equal now plus duration")
	}

	sleeps := map[int]time.Duration{
		1: 1 * time.Millisecond,
		2: 2 * time.Millisecond,
		3: 3 * time.Millisecond,
		4: 4 * time.Millisecond,
		5: 5 * time.Millisecond,
	}
	for i := 0; i < 1000; i++ {
		sleep := gen.SleepOneOf(sleeps[1], sleeps[2], sleeps[3], sleeps[4], sleeps[5])
		s := int(sleep) / 1000000
		_, ok := sleeps[s]
		Assert(t, OK(ok), "chosen sleep is one the arguments")
	}
}

// TestConcurrency simply produces a number of concurrent calls simply to let
// the race detection do its work.
func TestConcurrency(t *testing.T) {
	gen := qagen.New(qagen.FixedRand())

	run := func() {
		go gen.Byte(0, 255)
		go gen.Int(0, 9999)
		go gen.Duration(25*time.Millisecond, 75*time.Millisecond)
	}
	for i := 0; i < 100; i++ {
		go run()
	}

	time.Sleep(3 * time.Second)
}

// EOF
