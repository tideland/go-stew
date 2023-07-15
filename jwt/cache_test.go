// Tideland Go Stew - JSON Web Token - Cache - Unit Tests
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
	"context"
	"fmt"
	"testing"
	"time"

	. "tideland.dev/go/stew/assert"

	"tideland.dev/go/stew/jwt"
)

//--------------------
// TESTS
//--------------------

// TestCachePutGet verifies the putting and getting of tokens
// to the cache.
func TestCachePutGet(t *testing.T) {
	ctx := context.Background()
	maxEntries := 10
	cache := jwt.NewCache(ctx, time.Minute, time.Minute, time.Minute, maxEntries)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	_, err = cache.Put(jwtIn)
	Assert(t, NoError(err), "putting of token failed")
	jwt := jwtIn.String()
	jwtOut, err := cache.Get(jwt)
	Assert(t, NoError(err), "getting of token failed")
	Assert(t, Equal(jwtIn, jwtOut), "token is correct")
	jwtOut, err = cache.Get("is.not.there")
	Assert(t, NoError(err), "getting of token failed")
	Assert(t, Nil(jwtOut), "token is not there")
}

// TestCacheAccessCleanup verifies the access based cleanup
// of the JWT cache.
func TestCacheAccessCleanup(t *testing.T) {
	ctx := context.Background()
	maxEntries := 10
	cache := jwt.NewCache(ctx, time.Second, time.Second, time.Second, maxEntries)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	_, err = cache.Put(jwtIn)
	Assert(t, NoError(err), "putting of token failed")
	jwt := jwtIn.String()
	jwtOut, err := cache.Get(jwt)
	Assert(t, NoError(err), "getting of token failed")
	Assert(t, Equal(jwtIn, jwtOut), "token is correct")
	// Now wait a bit an try again.
	time.Sleep(5 * time.Second)
	jwtOut, err = cache.Get(jwt)
	Assert(t, NoError(err), "getting of token failed")
	Assert(t, Nil(jwtOut), "token is not there")
}

// TestCacheValidityCleanup verifies the validity based cleanup
// of the JWT cache.
func TestCacheValidityCleanup(t *testing.T) {
	ctx := context.Background()
	maxEntries := 10
	cache := jwt.NewCache(ctx, time.Minute, time.Second, time.Second, maxEntries)
	key := []byte("secret")
	now := time.Now()
	nbf := now.Add(-2 * time.Second)
	exp := now.Add(2 * time.Second)
	claims := initClaims()
	claims.SetNotBefore(nbf)
	claims.SetExpiration(exp)
	jwtIn, err := jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	_, err = cache.Put(jwtIn)
	Assert(t, NoError(err), "putting of token failed")
	jwt := jwtIn.String()
	jwtOut, err := cache.Get(jwt)
	Assert(t, NoError(err), "getting of token failed")
	Assert(t, Equal(jwtIn, jwtOut), "token is correct")
	// Now access until it is invalid and not
	// available anymore.
	var i int
	for i = 0; i < 5; i++ {
		time.Sleep(time.Second)
		jwtOut, err = cache.Get(jwt)
		Assert(t, NoError(err), "getting of token failed")
		if jwtOut == nil {
			break
		}
		Assert(t, Equal(jwtIn, jwtOut), "token is correct")
	}
	Assert(t, Range(i, 1, 4), "token is not there")
}

// TestCacheLoad verifies the cache load based cleanup.
func TestCacheLoad(t *testing.T) {
	ctx := context.Background()
	cacheTime := 100 * time.Millisecond
	maxEntries := 5
	cache := jwt.NewCache(ctx, 2*cacheTime, cacheTime, cacheTime, maxEntries)
	claims := initClaims()
	// Now fill the cache and check that it doesn't
	// grow too high.
	var i int
	for i = 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		key := []byte(fmt.Sprintf("secret-%d", i))
		jwtIn, err := jwt.Encode(claims, key, jwt.HS512)
		Assert(t, NoError(err), "encoding of token failed")
		size, err := cache.Put(jwtIn)
		Assert(t, NoError(err), "putting of token failed")
		Assert(t, True(size <= maxEntries), "cache size is correct")
	}
}

// TestCacheContext verifies the cache stopping by context.
func TestCacheContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	maxEntries := 10
	cache := jwt.NewCache(ctx, time.Minute, time.Minute, time.Minute, maxEntries)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := jwt.Encode(claims, key, jwt.HS512)
	Assert(t, NoError(err), "encoding of token failed")
	_, err = cache.Put(jwtIn)
	Assert(t, NoError(err), "putting of token failed")
	// Now cancel and test to get jwt.
	cancel()
	time.Sleep(10 * time.Millisecond)
	jwt := jwtIn.String()
	jwtOut, err := cache.Get(jwt)
	Assert(t, ErrorContains(err, "cache action timeout"), "getting of token failed")
	Assert(t, Nil(jwtOut), "token is not there")
}

//--------------------
// HELPERS
//--------------------

// initClaims creates test claims.
func initClaims() jwt.Claims {
	c := jwt.NewClaims()
	c.SetSubject("1234567890")
	c.Set("name", "John Doe")
	c.Set("admin", true)
	return c
}

// EOF
