// Tideland Go Stew - Etc - Unit Tests
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package etc_test

//--------------------
// IMPORTS
//--------------------

import (
	"bufio"
	"bytes"
	"context"
	"testing"
	"time"

	"tideland.dev/go/stew/qaenv"
	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/etc"
)

//--------------------
// TESTS
//--------------------

// TestRead tests reading a configuration out of a reader.
func TestRead(t *testing.T) {
	// Simple reader.
	cfg, err := etc.Read(bytes.NewReader(createEtcOne()))
	Assert(t, NoError(err), "no error expected")
	Assert(t, NotNil(cfg), "configuration expected")

	// Invalid string.
	cfg, err = etc.ReadString("[[foo}}")
	Assert(t, ErrorContains(err, `invalid source format`), "error expected")
	Assert(t, Nil(cfg), "no configuration expected")

	// File.
	td, err := qaenv.MkdirTemp("etc-test-")
	Assert(t, NoError(err), "no error expected")
	defer td.Restore()

	fn, err := td.WriteFile("etc.json", createEtcOne())
	Assert(t, NoError(err), "no error expected")

	cfg, err = etc.ReadFile(fn)
	Assert(t, NoError(err), "no error expected")
	Assert(t, NotNil(cfg), "configuration expected")
}

// TestWrite tests the writing of configurations.
func TestWrite(t *testing.T) {
	// Create and write first.
	cfgA, err := etc.Read(bytes.NewReader(createEtcOne()))
	Assert(t, NoError(err), "no error expected")
	Assert(t, NotNil(cfgA), "configuration expected")

	var buf bytes.Buffer

	err = cfgA.Write(&buf)
	Assert(t, NoError(err), "no error expected")

	// Read and compare, defaults are not needed.
	cfgB, err := etc.Read(bufio.NewReader(&buf))
	Assert(t, NoError(err), "no error expected")
	Assert(t, NotNil(cfgB), "configuration expected")

	ts := cfgB.At("global", "hostAddress").AsString("services.example.com:8080")
	Assert(t, Equal(ts, "localhost:8080"), "host address should be localhost:8080")

	ti := cfgB.At("global", "maxUsers").AsInt(0)
	Assert(t, Equal(ti, 50), "max users should be 50")

	tf := cfgB.At("services", "0", "bandwidth").AsFloat64(0.0)
	Assert(t, Equal(tf, 1.5), "bandwidth should be 1.5")

	tb := cfgB.At("services", "0", "active").AsBool(false)
	Assert(t, Equal(tb, true), "service-a should be active")

	tt := cfgB.At("global", "startTime").AsTime(time.TimeOnly, time.Now())
	h, m, s := tt.Clock()
	Assert(t, Equal(h, 6), "hour is 6")
	Assert(t, Equal(m, 0), "minute is 0")
	Assert(t, Equal(s, 0), "second is 0")

	td := cfgB.At("global", "working").AsDuration(time.Hour)
	Assert(t, Equal(td, 12*time.Hour), "working time is 12h")

	// Access non-existing value to get the default.
	ts = cfgB.At("unset", "string").AsString("expected")
	Assert(t, Equal(ts, "expected"), "default string expected")

	ti = cfgB.At("unset", "int").AsInt(42)
	Assert(t, Equal(ti, 42), "default int expected")

	tf = cfgB.At("unset", "float").AsFloat64(42.42)
	Assert(t, Equal(tf, 42.42), "default float expected")

	tb = cfgB.At("unset", "bool").AsBool(true)
	Assert(t, Equal(tb, true), "default bool expected")

	tt = cfgB.At("unset", "time").AsTime(time.TimeOnly, time.Now())
	h, m, s = tt.Clock()
	Assert(t, Equal(h, time.Now().Hour()), "default hour expected")
	Assert(t, Equal(m, time.Now().Minute()), "default minute expected")
	Assert(t, Equal(s, time.Now().Second()), "default second expected")

	td = cfgB.At("unset", "duration").AsDuration(time.Hour)
	Assert(t, Equal(td, time.Hour), "default duration expected")
}

// TestMacros tests the resolving of macros.
func TestMacros(t *testing.T) {
	cfg, err := etc.Read(bytes.NewReader(createEtcOne()))
	Assert(t, NoError(err), "no error expected")
	Assert(t, NotNil(cfg), "configuration expected")

	// Deep resolving.
	ts := cfg.At("services", "0", "url").AsString("")
	Assert(t, Equal(ts, "http://localhost:8040/service-a"), "URL should be http://localhost:8040/service-a")

	// Deep resolving with default.
	ts = cfg.At("services", "2", "url").AsString("http://localhost:8040/service-c")
	Assert(t, Equal(ts, "http://localhost:8040/service-c"), "URL should be http://localhost:8040/service-c")

	// Resolving with unset environment variable.
	ts = cfg.At("global", "baseDirectory").AsString("")
	Assert(t, Equal(ts, "/var/lib/my-server"), "base directory should be /var/lib/my-server")

	// Resolving with set environment variable.
	env := qaenv.NewEinvironment()
	defer env.Restore()

	ts = cfg.At("global", "title").AsString("Test in Progress")
	Assert(t, Equal(ts, "Fantastic Test Servicee"), "title should be Fantastic Test Servicee")
}

// TestContext tests adding a configuration to a context
// an retrieve it again.
func TestContext(t *testing.T) {
	cfgIn, err := etc.Read(bytes.NewReader(createEtcOne()))
	Assert(t, NoError(err), "no error expected")
	Assert(t, NotNil(cfgIn), "configuration expected")

	ctx := etc.NewContext(context.Background(), cfgIn)
	Assert(t, NotNil(ctx), "context expected")

	cfgOut, ok := etc.FromContext(ctx)
	Assert(t, True(ok), "configuration expected")
	Assert(t, Equal(cfgIn, cfgOut), "configurations should be equal")
}

//--------------------
// HELPER
//--------------------

// createEtcOne creates a simple configuration.
func createEtcOne() []byte {
	return []byte(`
	{
		"global": {
			"title": "[[$MYAPP_TITLE||My Server]]",
			"baseDirectory": "[[$MYAPP_BASEDIR||/var/lib/my-server]]",
			"hostAddress": "localhost:[[global::port||8080]]",
			"port": 8040,
			"maxUsers": 50,
			"startTime": "06:00:00",
			"working": "12h"
		},
		"services": [
			{
				"id": "service-a",
				"url": "http://[[global::hostAddress]]/service-a",
				"directory": "[[global::baseDirectory||.]]/service-a",
				"active": true,
				"bandwidth": 1.5
			}, {
				"id": "service-b",
				"url": "http://[[global::hostAddress]]/service-b",
				"directory": "[[global::baseDirectory||.]]/service-b",
				"active": false,
				"bandwidth": 0.5
			}
		]
	}
	`)
}

// EOF
