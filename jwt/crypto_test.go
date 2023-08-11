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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	. "tideland.dev/go/stew/qaone"

	"tideland.dev/go/stew/jwt"
)

//--------------------
// TESTS
//--------------------

var (
	esTests = []jwt.Algorithm{jwt.ES256, jwt.ES384, jwt.ES512}
	hsTests = []jwt.Algorithm{jwt.HS256, jwt.HS384, jwt.HS512}
	psTests = []jwt.Algorithm{jwt.PS256, jwt.PS384, jwt.PS512}
	rsTests = []jwt.Algorithm{jwt.RS256, jwt.RS384, jwt.RS512}
	data    = []byte("the quick brown fox jumps over the lazy dog")
)

// TestESAlgorithms verifies the ECDSA algorithms.
func TestESAlgorithms(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Assert(t, NoError(err), "generation of private key worked")
	for _, algo := range esTests {
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		Assert(t, NoError(err), "signing wilth algo %q worked", algo)
		Assert(t, NotEmpty(signature), "signature is not empty")
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		Assert(t, NoError(err), "verification with algo %q worked", algo)
	}
}

// TestHSAlgorithms verifies the HMAC algorithms.
func TestHSAlgorithms(t *testing.T) {
	key := []byte("secret")
	for _, algo := range hsTests {
		// Sign.
		signature, err := algo.Sign(data, key)
		Assert(t, NoError(err), "signing wilth algo %q worked", algo)
		Assert(t, NotEmpty(signature), "signature is not empty")
		// Verify.
		err = algo.Verify(data, signature, key)
		Assert(t, NoError(err), "verification with algo %q worked", algo)
	}
}

// TestPSAlgorithms verifies the RSAPSS algorithms.
func TestPSAlgorithms(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	Assert(t, NoError(err), "generation of private key worked")
	for _, algo := range psTests {
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		Assert(t, NoError(err), "signing wilth algo %q worked", algo)
		Assert(t, NotEmpty(signature), "signature is not empty")
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		Assert(t, NoError(err), "verification with algo %q worked", algo)
	}
}

// TestRSAlgorithms verifies the RSA algorithms.
func TestRSAlgorithms(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	Assert(t, NoError(err), "generation of private key worked")
	for _, algo := range rsTests {
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		Assert(t, NoError(err), "signing wilth algo %q worked", algo)
		Assert(t, NotEmpty(signature), "signature is not empty")
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		Assert(t, NoError(err), "verification with algo %q worked", algo)
	}
}

// TestNoneAlgorithm verifies the none algorithm.
func TestNoneAlgorithm(t *testing.T) {
	// Sign.
	signature, err := jwt.NONE.Sign(data, "")
	Assert(t, NoError(err), "signing without key worked")
	Assert(t, Empty(signature), "signature is empty")
	// Verify.
	err = jwt.NONE.Verify(data, signature, "")
	Assert(t, NoError(err), "verification without key worked")
}

// TestNotMatchingAlgorithm checks when algorithms of
// signing and verifying don't match.'
func TestNotMatchingAlgorithm(t *testing.T) {
	esPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	esPublicKey := esPrivateKey.Public()
	Assert(t, NoError(err), "generation of private key worked")
	hsKey := []byte("secret")
	rsPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	rsPublicKey := rsPrivateKey.Public()
	Assert(t, NoError(err), "generation of private key worked")
	noneKey := ""
	errorMatch := ".* combination of algorithm .* and key type .*"
	tests := []struct {
		description string
		algorithm   jwt.Algorithm
		key         jwt.Key
		signKeys    []jwt.Key
		verifyKeys  []jwt.Key
	}{
		{"ECDSA", jwt.ES512, esPrivateKey,
			[]jwt.Key{hsKey, rsPrivateKey, noneKey}, []jwt.Key{hsKey, rsPublicKey, noneKey}},
		{"HMAC", jwt.HS512, hsKey,
			[]jwt.Key{esPrivateKey, rsPrivateKey, noneKey}, []jwt.Key{esPublicKey, rsPublicKey, noneKey}},
		{"RSA", jwt.RS512, rsPrivateKey,
			[]jwt.Key{esPrivateKey, hsKey, noneKey}, []jwt.Key{esPublicKey, hsKey, noneKey}},
		{"RSAPSS", jwt.PS512, rsPrivateKey,
			[]jwt.Key{esPrivateKey, hsKey, noneKey}, []jwt.Key{esPublicKey, hsKey, noneKey}},
		{"none", jwt.NONE, noneKey,
			[]jwt.Key{esPrivateKey, hsKey, rsPrivateKey}, []jwt.Key{esPublicKey, hsKey, rsPublicKey}},
	}
	// Run the tests.
	for _, test := range tests {
		for _, key := range test.signKeys {
			_, err := test.algorithm.Sign(data, key)
			Assert(t, ErrorMatches(err, errorMatch), "signing with %s algorithm and %T key failed", test.algorithm, key)
		}
		signature, err := test.algorithm.Sign(data, test.key)
		Assert(t, NoError(err), "signing with %s algorithm and %T key worked", test.algorithm, test.key)
		for _, key := range test.verifyKeys {
			err = test.algorithm.Verify(data, signature, key)
			Assert(t, ErrorMatches(err, errorMatch), "verification with %s algorithm and %T key failed", test.algorithm, key)
		}
	}
}

// TestESTools verifies the tools for the reading of PEM encoded
// ECDSA keys.
func TestESTools(t *testing.T) {
	// Generate keys and PEMs.
	privateKeyIn, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Assert(t, NoError(err), "generation of private key worked")
	privateBytes, err := x509.MarshalECPrivateKey(privateKeyIn)
	Assert(t, NoError(err), "marshaling of private key worked")
	privateBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	Assert(t, NoError(err), "marshaling of public key worked")
	publicBlock := pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	Assert(t, NotEmpty(publicPEM), "public PEM is not empty")
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := jwt.ReadECPrivateKey(buf)
	Assert(t, NoError(err), "reading of private key worked")
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := jwt.ReadECPublicKey(buf)
	Assert(t, NoError(err), "reading of public key worked")
	// And as a last step check if they are correctly usable.
	signature, err := jwt.ES512.Sign(data, privateKeyOut)
	Assert(t, NoError(err), "signing with ES512 algorithm and ECDSA key worked")
	err = jwt.ES512.Verify(data, signature, publicKeyOut)
	Assert(t, NoError(err), "verification with ES512 algorithm and ECDSA key worked")
}

// TestRSTools verifies the tools for the reading of PEM encoded
// RSA keys.
func TestRSTools(t *testing.T) {
	// Generate keys and PEMs.
	privateKeyIn, err := rsa.GenerateKey(rand.Reader, 2048)
	Assert(t, NoError(err), "generation of private key worked")
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKeyIn)
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	Assert(t, NoError(err), "marshaling of public key worked")
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	Assert(t, NotNil(publicPEM), "public PEM is not nil")
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := jwt.ReadRSAPrivateKey(buf)
	Assert(t, NoError(err), "reading of private key worked")
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := jwt.ReadRSAPublicKey(buf)
	Assert(t, NoError(err), "reading of public key worked")
	// And as a last step check if they are correctly usable.
	signature, err := jwt.RS512.Sign(data, privateKeyOut)
	Assert(t, NoError(err), "signing with RS512 algorithm and RSA key worked")
	err = jwt.RS512.Verify(data, signature, publicKeyOut)
	Assert(t, NoError(err), "verification with RS512 algorithm and RSA key worked")
}

// EOF
