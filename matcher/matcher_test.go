// Tideland Go Stew - Matcher - Unit Tests
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package matcher_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/stew/asserts"
	"tideland.dev/go/stew/matcher"
)

//--------------------
// TESTS
//--------------------

// TestMatches tests matching string.
func TestMatches(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	tests := []struct {
		name       string
		pattern    string
		value      string
		ignoreCase bool
		out        bool
	}{
		{
			"equal pattern and string without wildcards",
			"quick brown fox",
			"quick brown fox",
			matcher.IgnoreCase,
			true,
		}, {
			"unequal pattern and string without wildcards",
			"quick brown fox",
			"lazy dog",
			matcher.IgnoreCase,
			false,
		}, {
			"matching pattern with one question mark",
			"quick brown f?x",
			"quick brown fox",
			matcher.IgnoreCase,
			true,
		}, {
			"matching pattern with one asterisk",
			"quick*fox",
			"quick brown fox",
			matcher.IgnoreCase,
			true,
		}, {
			"matching pattern with char group",
			"quick brown f[ao]x",
			"quick brown fox",
			matcher.IgnoreCase,
			true,
		}, {
			"not-matching pattern with char group",
			"quick brown f[eiu]x",
			"quick brown fox",
			matcher.IgnoreCase,
			false,
		}, {
			"matching pattern with char range",
			"quick brown f[a-u]x",
			"quick brown fox",
			matcher.IgnoreCase,
			true,
		}, {
			"not-matching pattern with char range",
			"quick brown f[^a-u]x",
			"quick brown fox",
			matcher.IgnoreCase,
			false,
		}, {
			"matching pattern with char group not ignoring care",
			"quick * F[aeiou]x",
			"quick * Fox",
			matcher.ValidateCase,
			true,
		}, {
			"not-matching pattern with char group not ignoring care",
			"quick * F[aeiou]x",
			"quick * fox",
			matcher.ValidateCase,
			false,
		}, {
			"matching pattern with escape",
			"quick \\* f\\[o\\]x",
			"quick * f[o]x",
			matcher.IgnoreCase,
			true,
		}, {
			"not-matching pattern with escape",
			"quick \\* f\\[o\\]x",
			"quick brown fox",
			matcher.IgnoreCase,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer assert.SetFailable(t)()
			out := matcher.Matches(test.pattern, test.value, test.ignoreCase)
			assert.Equal(out, test.out)
		})
	}
}

// EOF
