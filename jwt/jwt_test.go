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

	"tideland.dev/go/stew/asserts"
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
	assert := asserts.NewTesting(t, asserts.FailStop)
	// Decode.
	token, err := jwt.Decode(rawToken)
	assert.Nil(err)
	assert.Equal(token.Algorithm(), jwt.HS512)
	key, err := token.Key()
	assert.Nil(key)
	assert.ErrorMatch(err, ".*no key available, only after encoding or verifying.*")
	assert.Length(token.Claims(), 4)

	sub, ok := token.Claims().GetString("sub")
	assert.True(ok)
	assert.Equal(sub, subClaim)
	name, ok := token.Claims().GetString("name")
	assert.True(ok)
	assert.Equal(name, nameClaim)
	admin, ok := token.Claims().GetBool("admin")
	assert.True(ok)
	assert.Equal(admin, adminClaim)
	iat, ok := token.Claims().IssuedAt()
	assert.True(ok)
	assert.Equal(iat, time.Unix(iatClaim, 0))
	exp, ok := token.Claims().Expiration()
	assert.False(ok)
	assert.Equal(exp, time.Time{})
}

// TestIsValid verifies the time validation of a token.
func TestIsValid(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing time validation")
	now := time.Now()
	leeway := time.Minute
	key := []byte("secret")
	// Create token with no times set, encode, decode, validate ok.
	claims := jwt.NewClaims()
	tokenEnc, err := jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	tokenDec, err := jwt.Decode(tokenEnc.String())
	assert.Nil(err)
	ok := tokenDec.IsValid(leeway)
	assert.True(ok)
	// Now a token with a long timespan, still valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(-time.Hour))
	claims.SetExpiration(now.Add(time.Hour))
	tokenEnc, err = jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	tokenDec, err = jwt.Decode(tokenEnc.String())
	assert.Nil(err)
	ok = tokenDec.IsValid(leeway)
	assert.True(ok)
	// Now a token with a long timespan in the past, not valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(-2 * time.Hour))
	claims.SetExpiration(now.Add(-time.Hour))
	tokenEnc, err = jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	tokenDec, err = jwt.Decode(tokenEnc.String())
	assert.Nil(err)
	ok = tokenDec.IsValid(leeway)
	assert.False(ok)
	// And at last a token with a long timespan in the future, not valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(time.Hour))
	claims.SetExpiration(now.Add(2 * time.Hour))
	tokenEnc, err = jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	tokenDec, err = jwt.Decode(tokenEnc.String())
	assert.Nil(err)
	ok = tokenDec.IsValid(leeway)
	assert.False(ok)
}

// EOF
