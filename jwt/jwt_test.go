// Tideland Go Stew - JSON Web Token - Unit Tests
//
// Copyright (C) 2016-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/jwt"
)

//--------------------
// CONSTANTS
//--------------------

const (
	subClaim   = "1234567890"
	nameClaim  = "John Doe"
	adminClaim = true
	iatClaim   = 1600000000
	rawToken   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYwMDAwMDAwMH0." +
		"P50peTbENKIPw0tjuHLgosFmJRYGTh_kNA9IcyWIoJ39uYMa4JfKYhnQw5mkgSLB2WYVT68QaDeWWErn4lU69g"
)

//--------------------
// TESTS
//--------------------

// TestDecode verifies a token decoding without internal verification the signature.
func TestDecode(t *testing.T) {
	// Decode.
	token, err := jwt.Decode(rawToken)
	Assert(t, NoError(err), "decoding of token failed")
	Assert(t, Equal(token.Algorithm(), jwt.HS512), "algorithm is correct")
	key, err := token.Key()
	Assert(t, Nil(key), "no key available")
	Assert(t, ErrorMatches(err, ".*no key available, only after encoding or verifying.*"), "error is correct")
	Assert(t, Length(token.Claims(), 4), "length of token claims is correct")

	sub, ok := token.Claims().GetString("sub")
	Assert(t, OK(ok), "sub claim is available")
	Assert(t, Equal(sub, subClaim), "sub claim is correct")
	name, ok := token.Claims().GetString("name")
	Assert(t, OK(ok), "name claim is available")
	Assert(t, Equal(name, nameClaim), "name claim is correct")
	admin, ok := token.Claims().GetBool("admin")
	Assert(t, OK(ok), "admin claim is available")
	Assert(t, Equal(admin, adminClaim), "admin claim is correct")
	iat, ok := token.Claims().IssuedAt()
	Assert(t, OK(ok), "iat claim is available")
	Assert(t, Equal(iat, time.Unix(iatClaim, 0)), "iat claim is correct")
	exp, ok := token.Claims().Expiration()
	Assert(t, NotOK(ok), "exp claim is not available")
	Assert(t, Equal(exp, time.Time{}), "exp claim is correct set empty")
}

// TestIsValid verifies the time validation of a token.
func TestIsValid(t *testing.T) {
	now := time.Now()
	leeway := time.Minute
	key := []byte("secret")
	// Create token with no times set, encode, decode, validate ok.
	claims := jwt.NewClaims()
	tokenEnc, err := jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	tokenDec, err := jwt.Decode(tokenEnc.String())
	Assert(t, NoError(err), "decoding of token failed")
	ok := tokenDec.IsValid(leeway)
	Assert(t, OK(ok), "token is valid")
	// Now a token with a long timespan, still valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(-time.Hour))
	claims.SetExpiration(now.Add(time.Hour))
	tokenEnc, err = jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	tokenDec, err = jwt.Decode(tokenEnc.String())
	Assert(t, NoError(err), "decoding of token failed")
	ok = tokenDec.IsValid(leeway)
	Assert(t, OK(ok), "token is valid")
	// Now a token with a long timespan in the past, not valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(-2 * time.Hour))
	claims.SetExpiration(now.Add(-time.Hour))
	tokenEnc, err = jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	tokenDec, err = jwt.Decode(tokenEnc.String())
	Assert(t, NoError(err), "decoding of token failed")
	ok = tokenDec.IsValid(leeway)
	Assert(t, NotOK(ok), "token is not valid")
	// And at last a token with a long timespan in the future, not valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(time.Hour))
	claims.SetExpiration(now.Add(2 * time.Hour))
	tokenEnc, err = jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	tokenDec, err = jwt.Decode(tokenEnc.String())
	Assert(t, NoError(err), "decoding of token failed")
	ok = tokenDec.IsValid(leeway)
	Assert(t, NotOK(ok), "token is not valid")
}

// EOF
