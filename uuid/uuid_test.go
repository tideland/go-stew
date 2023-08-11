// Tideland Go Stew - UUID - Unit Tests
//
// Copyright (C) 2021-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package uuid_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/uuid"
)

//--------------------
// TESTS
//--------------------

// TestStandard tests the standard UUID.
func TestStandard(t *testing.T) {
	uuidA := uuid.New()
	Assert(t, Equal(uuidA.Version(), uuid.V4), "wrong UUID version")
	uuidAShortStr := uuidA.ShortString()
	uuidAStr := uuidA.String()
	Assert(t, Length(uuidA, 16), "wrong UUID short string length")
	Assert(t, Length(uuidAShortStr, 32), "wrong UUID string length")
	Assert(t, Matches(uuidAShortStr, "[0-9a-f]{32}"), "wrong UUID short string format")
	Assert(t, Matches(uuidAStr, "[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}"), "wrong UUID string format")
	// Check for copy.
	uuidB := uuid.New()
	uuidC := uuidB.Copy()
	for i := 0; i < len(uuidB); i++ {
		uuidB[i] = 0
	}
	Assert(t, Different(uuidB, uuidC), "UUID copy not independent")
}

// TestNamespaces tests the creation of the different standard namespaces.
func TestNamespaces(t *testing.T) {
	Assert(t, Equal(uuid.NamespaceDNS().String(), "6ba7b810-9dad-11d1-80b4-00c04fd430c8"), "wrong UUID for DNS namespace")
	Assert(t, Equal(uuid.NamespaceURL().String(), "6ba7b811-9dad-11d1-80b4-00c04fd430c8"), "wrong UUID for URL namespace")
	Assert(t, Equal(uuid.NamespaceOID().String(), "6ba7b812-9dad-11d1-80b4-00c04fd430c8"), "wrong UUID for OID namespace")
	Assert(t, Equal(uuid.NamespaceX500().String(), "6ba7b814-9dad-11d1-80b4-00c04fd430c8"), "wrong UUID for X.500 namespace")
}

// TestRepeatability tests the repeatability of UUIDs v3 and v5 with same namespace and name.
func TestRepeatability(t *testing.T) {
	ns := uuid.NamespaceDNS()
	name := []byte("tideland.dev")
	uuidV3A, err := uuid.NewV3(ns, name)
	Assert(t, NoError(err), "error creating UUID v3")
	uuidV3B, err := uuid.NewV3(ns, name)
	Assert(t, NoError(err), "error creating UUID v3")
	Assert(t, Equal(uuidV3A, uuidV3B), "UUID v3 not repeatable")
	uuidV5A, err := uuid.NewV5(ns, name)
	Assert(t, NoError(err), "error creating UUID v5")
	uuidV5B, err := uuid.NewV5(ns, name)
	Assert(t, NoError(err), "error creating UUID v5")
	Assert(t, Equal(uuidV5A, uuidV5B), "UUID v5 not repeatable")
}

// TestVersions tests the creation of different UUID versions.
func TestVersions(t *testing.T) {
	ns := uuid.NamespaceOID()
	name := []byte{1, 3, 3, 7}
	uuidV1, err := uuid.NewV1()
	Assert(t, NoError(err), "error creating UUID V1")
	Assert(t, Equal(uuidV1.Version(), uuid.V1), "wrong UUID version")
	Assert(t, Equal(uuidV1.Variant(), uuid.VariantRFC4122), "wrong UUID variant")
	uuidV3, err := uuid.NewV3(ns, name)
	Assert(t, NoError(err), "error creating UUID V3")
	Assert(t, Equal(uuidV3.Version(), uuid.V3), "wrong UUID version")
	Assert(t, Equal(uuidV3.Variant(), uuid.VariantRFC4122), "wrong UUID variant")
	uuidV4, err := uuid.NewV4()
	Assert(t, NoError(err), "error creating UUID V4")
	Assert(t, Equal(uuidV4.Version(), uuid.V4), "wrong UUID version")
	Assert(t, Equal(uuidV4.Variant(), uuid.VariantRFC4122), "wrong UUID variant")
	uuidV5, err := uuid.NewV5(ns, name)
	Assert(t, NoError(err), "error creating UUID V5")
	Assert(t, Equal(uuidV5.Version(), uuid.V5), "wrong UUID version")
	Assert(t, Equal(uuidV5.Variant(), uuid.VariantRFC4122), "wrong UUID variant")
}

// TestNil tests the nil UUID.
func TestNil(t *testing.T) {
	var uuidNil uuid.UUID
	Assert(t, Equal(uuidNil.Version(), 0), "UUID version")
	Assert(t, Equal(uuidNil.Variant(), uuid.VariantNCS), "UUID variant")
	Assert(t, Equal(uuidNil.String(), "00000000-0000-0000-0000-000000000000"), "UUID string")
	Assert(t, Equal(uuidNil.ShortString(), "00000000000000000000000000000000"), "UUID short string")
}

// TestParse tests creating UUIDs from different string representations.
func TestParse(t *testing.T) {
	ns := uuid.NamespaceOID()
	name := []byte{1, 3, 3, 7}
	tests := []struct {
		source  func() string
		version uuid.Version
		variant uuid.Variant
		err     string
	}{
		{func() string { u, _ := uuid.NewV1(); return u.String() }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return u.String() }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return u.String() }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return u.String() }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV1(); return "urn:uuid:" + u.String() }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return "urn:uuid:" + u.String() }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return "urn:uuid:" + u.String() }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return "urn:uuid:" + u.String() }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV1(); return "{" + u.String() + "}" }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return "{" + u.String() + "}" }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return "{" + u.String() + "}" }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return "{" + u.String() + "}" }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV1(); return u.ShortString() }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return u.ShortString() }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return u.ShortString() }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return u.ShortString() }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return u.String() + "-ffaabb" }, 0, 0, "invalid source format"},
		{func() string { u, _ := uuid.NewV4(); return u.String() + "-ffxxyy" }, 0, 0, "invalid source format"},
		{func() string { u, _ := uuid.NewV4(); return "uuid:" + u.String() }, 0, 0, "invalid source format"},
		{func() string { u, _ := uuid.NewV4(); return "{" + u.ShortString() + "}" }, 0, 0, "invalid source format"},
		{func() string { return "ababababababababab" }, 0, 0, "invalid source format"},
		{func() string { return "abcdefabcdefZZZZefabcdefabcdefab" }, 0, 0, "source char 12 is no hex char"},
		{func() string { return "[abcdefabcdefabcdefabcdefabcdefab]" }, 0, 0, "invalid source format"},
		{func() string { return "abcdefab=cdef=abcd=efab=cdefabcdefab" }, 0, 0, "source char 8 does not match pattern"},
	}
	for _, test := range tests {
		source := test.source()
		uuidT, err := uuid.Parse(source)
		if test.err == "" {
			Assert(t, NoError(err), "error parsing UUID %q", source)
			Assert(t, Equal(uuidT.Version(), test.version), "wrong UUID version for %q", source)
			Assert(t, Equal(uuidT.Variant(), test.variant), "wrong UUID variant for %q", source)
		} else {
			Assert(t, ErrorContains(err, test.err), "wrong error parsing UUID %q", source)
		}
	}
}

// EOF
