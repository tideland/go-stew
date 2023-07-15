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
	"encoding/json"
	"testing"
	"time"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/jwt"
)

//--------------------
// TESTS
//--------------------

// TestClaimsMarshalling verifies the marshalling of claims to JSON and back.
func TestClaimsMarshalling(t *testing.T) {
	// First with uninitialised or empty jwt.
	var c jwt.Claims
	jsonValue, err := json.Marshal(c)
	Assert(t, Equal(string(jsonValue), "{}"), "empty claims marshalled")
	Assert(t, NoError(err), "no error")
	c = jwt.NewClaims()
	jsonValue, err = json.Marshal(c)
	Assert(t, Equal(string(jsonValue), "{}"), "empty claims marshalled")
	Assert(t, NoError(err), "no error")
	// Now fill it.
	c.Set("foo", "yadda")
	c.Set("bar", 12345)
	Assert(t, Length(c, 2), "length of claims")
	jsonValue, err = json.Marshal(c)
	Assert(t, NotNil(jsonValue), "claims marshalled")
	Assert(t, NoError(err), "no error")
	var uc jwt.Claims
	err = json.Unmarshal(jsonValue, &uc)
	Assert(t, NoError(err), "no error")
	Assert(t, Length(uc, 2), "length of claims")
	foo, ok := uc.Get("foo")
	Assert(t, Equal(foo, "yadda"), "foo claim")
	Assert(t, OK(ok), "foo claim")
	bar, ok := uc.GetInt("bar")
	Assert(t, Equal(bar, 12345), "bar claim")
	Assert(t, OK(ok), "bar claim")
}

// TestClaimsBasic verifies the low level operations on claims.
func TestClaimsBasic(t *testing.T) {
	// First with uninitialised jwt.
	var c jwt.Claims
	ok := c.Contains("foo")
	Assert(t, NotOK(ok), "initial foo not contained")
	nothing, ok := c.Get("foo")
	Assert(t, Nil(nothing), "foo not contained")
	Assert(t, NotOK(ok), "foo not contained")
	old := c.Set("foo", "bar")
	Assert(t, Nil(old), "foo not contained")
	old = c.Delete("foo")
	Assert(t, Nil(old), "foo not contained")
	// Now initialise it.
	c = jwt.NewClaims()
	ok = c.Contains("foo")
	Assert(t, NotOK(ok), "initial foo not contained")
	nothing, ok = c.Get("foo")
	Assert(t, Nil(nothing), "foo not contained")
	Assert(t, NotOK(ok), "foo not contained")
	old = c.Set("foo", "bar")
	Assert(t, Nil(old), "foo not contained so far")
	ok = c.Contains("foo")
	Assert(t, OK(ok), "foo contained now")
	foo, ok := c.Get("foo")
	Assert(t, Equal(foo, "bar"), "foo contained")
	Assert(t, OK(ok), "foo contained")
	old = c.Set("foo", "yadda")
	Assert(t, Equal(old, "bar"), "foo contained and replaced")
	// Finally delete it.
	old = c.Delete("foo")
	Assert(t, Equal(old, "yadda"), "foo contained and deleted")
	old = c.Delete("foo")
	Assert(t, Nil(old), "foo not contained anymore")
	ok = c.Contains("foo")
	Assert(t, NotOK(ok), "foo not contained anymore")
}

// TestClaimsString verifies the string operations on claims.
func TestClaimsString(t *testing.T) {
	c := jwt.NewClaims()
	nothing := c.Set("foo", "bar")
	Assert(t, Nil(nothing), "foo not contained so far")
	var foo string
	foo, ok := c.GetString("foo")
	Assert(t, Equal(foo, "bar"), "foo contained")
	Assert(t, OK(ok), "foo contained")
	c.Set("foo", 4711)
	foo, ok = c.GetString("foo")
	Assert(t, Equal(foo, "4711"), "foo contained")
	Assert(t, OK(ok), "foo contained")
}

// TestClaimsBool verifies the bool operations on claims.
func TestClaimsBool(t *testing.T) {
	c := jwt.NewClaims()
	c.Set("foo", true)
	c.Set("bar", false)
	c.Set("baz", "T")
	c.Set("bingo", "0")
	c.Set("yadda", "nope")
	foo, ok := c.GetBool("foo")
	Assert(t, True(foo), "foo contained")
	Assert(t, OK(ok), "foo contained")
	bar, ok := c.GetBool("bar")
	Assert(t, False(bar), "bar contained")
	Assert(t, OK(ok), "bar contained")
	baz, ok := c.GetBool("baz")
	Assert(t, True(baz), "baz contained")
	Assert(t, OK(ok), "baz contained")
	bingo, ok := c.GetBool("bingo")
	Assert(t, False(bingo), "bingo contained")
	Assert(t, OK(ok), "bingo contained")
	yadda, ok := c.GetBool("yadda")
	Assert(t, False(yadda), "yadda no bool")
	Assert(t, NotOK(ok), "yadda no bool")
}

// TestClaimsInt verifies the int operations on claims.
func TestClaimsInt(t *testing.T) {
	c := jwt.NewClaims()
	c.Set("foo", 4711)
	c.Set("bar", "4712")
	c.Set("baz", 4713.0)
	c.Set("yadda", "nope")
	foo, ok := c.GetInt("foo")
	Assert(t, Equal(foo, 4711), "foo contained")
	Assert(t, OK(ok), "foo contained")
	bar, ok := c.GetInt("bar")
	Assert(t, Equal(bar, 4712), "bar contained")
	Assert(t, OK(ok), "bar contained")
	baz, ok := c.GetInt("baz")
	Assert(t, Equal(baz, 4713), "baz contained")
	Assert(t, OK(ok), "baz contained")
	yadda, ok := c.GetInt("yadda")
	Assert(t, Equal(yadda, 0), "yadda no int")
	Assert(t, NotOK(ok), "yadda no int")
}

// TestClaimsFloat64 verifies the float64 operations on claims.
func TestClaimsFloat64(t *testing.T) {
	c := jwt.NewClaims()
	c.Set("foo", 4711)
	c.Set("bar", "4712")
	c.Set("baz", 4713.0)
	c.Set("yadda", "nope")
	foo, ok := c.GetFloat64("foo")
	Assert(t, Equal(foo, 4711.0), "foo contained")
	Assert(t, OK(ok), "foo contained")
	bar, ok := c.GetFloat64("bar")
	Assert(t, Equal(bar, 4712.0), "bar contained")
	Assert(t, OK(ok), "bar contained")
	baz, ok := c.GetFloat64("baz")
	Assert(t, Equal(baz, 4713.0), "baz contained")
	Assert(t, OK(ok), "baz contained")
	yadda, ok := c.GetFloat64("yadda")
	Assert(t, Equal(yadda, 0.0), "yadda no float64")
	Assert(t, NotOK(ok), "yadda no float64")
}

// TestClaimsTime verifies the time operations on claims.
func TestClaimsTime(t *testing.T) {
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := jwt.NewClaims()
	c.SetTime("foo", goLaunch)
	c.Set("bar", goLaunch.Unix())
	c.Set("baz", goLaunch.Format(time.RFC3339))
	c.Set("yadda", "nope")
	foo, ok := c.GetTime("foo")
	Assert(t, Equal(foo.Unix(), goLaunch.Unix()), "foo contained")
	Assert(t, OK(ok), "foo contained")
	bar, ok := c.GetTime("bar")
	Assert(t, Equal(bar.Unix(), goLaunch.Unix()), "bar contained")
	Assert(t, OK(ok), "bar contained")
	baz, ok := c.GetTime("baz")
	Assert(t, Equal(baz.Unix(), goLaunch.Unix()), "baz contained")
	Assert(t, OK(ok), "baz contained")
	yadda, ok := c.GetTime("yadda")
	Assert(t, Equal(yadda, time.Time{}), "yadda no time")
	Assert(t, NotOK(ok), "yadda no time")
}

// TestClaimsMarshalledValue verifies the marshalling and
// unmarshalling of structures as values.
func TestClaimsMarshalledValue(t *testing.T) {
	type nestedValue struct {
		Name  string
		Value int
	}

	baz := []*nestedValue{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}
	c := jwt.NewClaims()
	c.Set("foo", "bar")
	c.Set("baz", baz)
	// Now marshal and unmarshal the claim.
	jsonValue, err := json.Marshal(c)
	Assert(t, NoError(err), "no error")
	Assert(t, NotNil(jsonValue), "json marshalled")
	var uc jwt.Claims
	err = json.Unmarshal(jsonValue, &uc)
	Assert(t, NoError(err), "no error")
	Assert(t, Length(uc, 2), "length of claims")
	foo, ok := uc.Get("foo")
	Assert(t, OK(ok), "foo contained")
	Assert(t, Equal(foo, "bar"), "foo contained")
	var ubaz []*nestedValue
	ok, err = uc.GetMarshalled("baz", &ubaz)
	Assert(t, OK(ok), "baz contained")
	Assert(t, Nil(err), "no error")
	Assert(t, Length(ubaz, 3), "length of baz")
	Assert(t, Equal(ubaz[0].Name, "one"), "baz[0] name")
	Assert(t, Equal(ubaz[2].Value, 3), "baz[2] value")
}

// TestClaimsAudience verifies the setting, getting, and
// deleting of the audience claim.
func TestClaimsAudience(t *testing.T) {
	audience := []string{"foo", "bar", "baz"}
	c := jwt.NewClaims()
	aud, ok := c.Audience()
	Assert(t, NotOK(ok), "no audience")
	Assert(t, Nil(aud), "no audience")
	none := c.SetAudience(audience...)
	Assert(t, Length(none, 0), "no old audience")
	aud, ok = c.Audience()
	Assert(t, DeepEqual(aud, audience), "audience")
	Assert(t, OK(ok), "audience")
	old := c.DeleteAudience()
	Assert(t, DeepEqual(old, aud), "audience")
	_, ok = c.Audience()
	Assert(t, NotOK(ok), "no audience")
}

// TestClaimsExpiration verifies the setting, getting, and
// deleting of the expiration claim.
func TestClaimsExpiration(t *testing.T) {
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := jwt.NewClaims()
	exp, ok := c.Expiration()
	Assert(t, NotOK(ok), "no expiration")
	Assert(t, Equal(exp, time.Time{}), "no expiration")
	none := c.SetExpiration(goLaunch)
	Assert(t, Equal(none, time.Time{}), "no old expiration")
	exp, ok = c.Expiration()
	Assert(t, OK(ok), "expiration exists")
	Assert(t, Equal(exp.Unix(), goLaunch.Unix()), "expiration exists")
	old := c.DeleteExpiration()
	Assert(t, Equal(old.Unix(), exp.Unix()), "expiration exists, now deleted")
	exp, ok = c.Expiration()
	Assert(t, NotOK(ok), "no expiration anymore")
	Assert(t, Equal(exp, time.Time{}), "no expiration anymore")
}

// TestClaimsIdentifier verifies the setting, getting, and
// deleting of the identifier claim.
func TestClaimsIdentifier(t *testing.T) {
	identifier := "foo"
	c := jwt.NewClaims()
	jti, ok := c.Identifier()
	Assert(t, NotOK(ok), "no identifier")
	Assert(t, Empty(jti), "no identifier")
	none := c.SetIdentifier(identifier)
	Assert(t, Equal(none, ""), "no old identifier, set new one")
	jti, ok = c.Identifier()
	Assert(t, OK(ok), "identifier")
	Assert(t, Equal(jti, identifier), "identifier")
	old := c.DeleteIdentifier()
	Assert(t, Equal(old, jti), "deleted old identifier")
	_, ok = c.Identifier()
	Assert(t, NotOK(ok), "no identifier anymore")
}

// TestClaimsIssuedAt verifies the setting, getting, and
// deleting of the issued at claim.
func TestClaimsIssuedAt(t *testing.T) {
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := jwt.NewClaims()
	iat, ok := c.IssuedAt()
	Assert(t, NotOK(ok), "no issued at")
	Assert(t, Equal(iat, time.Time{}), "no issued at")
	none := c.SetIssuedAt(goLaunch)
	Assert(t, Equal(none, time.Time{}), "no old issued at")
	iat, ok = c.IssuedAt()
	Assert(t, OK(ok), "issued at exists")
	Assert(t, Equal(iat.Unix(), goLaunch.Unix()), "issued at exists")
	old := c.DeleteIssuedAt()
	Assert(t, Equal(old.Unix(), iat.Unix()), "issued at exists, now deleted")
	iat, ok = c.IssuedAt()
	Assert(t, NotOK(ok), "no issued at anymore")
	Assert(t, Equal(iat, time.Time{}), "no issued at anymore")
}

// TestClaimsIssuer verifies the setting, getting, and
// deleting of the issuer claim.
func TestClaimsIssuer(t *testing.T) {
	issuer := "foo"
	c := jwt.NewClaims()
	iss, ok := c.Issuer()
	Assert(t, NotOK(ok), "no issuer")
	Assert(t, Empty(iss), "no issuer")
	none := c.SetIssuer(issuer)
	Assert(t, Equal(none, ""), "no old issuer, set new one")
	iss, ok = c.Issuer()
	Assert(t, OK(ok), "issuer")
	Assert(t, Equal(iss, issuer), "issuer")
	old := c.DeleteIssuer()
	Assert(t, Equal(old, iss), "deleted old issuer")
	_, ok = c.Issuer()
	Assert(t, NotOK(ok), "no issuer anymore")
}

// TestClaimsNotBefore verifies the setting, getting, and
// deleting of the not before claim.
func TestClaimsNotBefore(t *testing.T) {
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := jwt.NewClaims()
	nbf, ok := c.NotBefore()
	Assert(t, NotOK(ok), "no not before")
	Assert(t, Equal(nbf, time.Time{}), "no not before")
	none := c.SetNotBefore(goLaunch)
	Assert(t, Equal(none, time.Time{}), "no old not before")
	nbf, ok = c.NotBefore()
	Assert(t, OK(ok), "not before exists")
	Assert(t, Equal(nbf.Unix(), goLaunch.Unix()), "not before exists")
	old := c.DeleteNotBefore()
	Assert(t, Equal(old.Unix(), nbf.Unix()), "not before exists, now deleted")
	_, ok = c.NotBefore()
	Assert(t, NotOK(ok), "no not before anymore")
}

// TestClaimsSubject verifies the setting, getting, and
// deleting of the subject claim.
func TestClaimsSubject(t *testing.T) {
	subject := "foo"
	c := jwt.NewClaims()
	sub, ok := c.Subject()
	Assert(t, NotOK(ok), "no subject")
	Assert(t, Empty(sub), "no subject")
	none := c.SetSubject(subject)
	Assert(t, Equal(none, ""), "no old subject, set new one")
	sub, ok = c.Subject()
	Assert(t, OK(ok), "subject")
	Assert(t, Equal(sub, subject), "subject")
	old := c.DeleteSubject()
	Assert(t, Equal(old, sub), "deleted old subject")
	_, ok = c.Subject()
	Assert(t, NotOK(ok), "no subject anymore")
}

// TestClaimsValidity verifies the validation of the not before
// and the expiring time.
func TestClaimsValidity(t *testing.T) {
	// Fresh jwt.
	now := time.Now()
	leeway := time.Minute
	c := jwt.NewClaims()
	valid := c.IsAlreadyValid(leeway)
	Assert(t, OK(valid), "is already valid with leeway")
	valid = c.IsStillValid(leeway)
	Assert(t, OK(valid), "is still valid with leeway")
	valid = c.IsValid(leeway)
	Assert(t, OK(valid), "is  valid with leeway")
	// Set times.
	nbf := now.Add(-time.Hour)
	exp := now.Add(time.Hour)
	c.SetNotBefore(nbf)
	valid = c.IsAlreadyValid(leeway)
	Assert(t, OK(valid), "is already valid with leeway")
	c.SetExpiration(exp)
	valid = c.IsStillValid(leeway)
	Assert(t, OK(valid), "is still valid with leeway")
	valid = c.IsValid(leeway)
	Assert(t, OK(valid), "is valid with leeway")
	// Invalid token.
	nbf = now.Add(time.Hour)
	exp = now.Add(-time.Hour)
	c.SetNotBefore(nbf)
	c.DeleteExpiration()
	valid = c.IsAlreadyValid(leeway)
	Assert(t, NotOK(valid), "is not already valid with leeway")
	valid = c.IsValid(leeway)
	Assert(t, NotOK(valid), "is not valid with leeway")
	c.DeleteNotBefore()
	c.SetExpiration(exp)
	valid = c.IsStillValid(leeway)
	Assert(t, NotOK(valid), "is not still valid with leeway")
	valid = c.IsValid(leeway)
	Assert(t, NotOK(valid), "is not valid with leeway")
	c.SetNotBefore(nbf)
	valid = c.IsValid(leeway)
	Assert(t, NotOK(valid), "is not valid with leeway")
}

// EOF
