// Tideland Go Stew - JSON Web Token - Crypto
//
// Copyright (C) 2016-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt // import "tideland.dev/go/stew/jwt"

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

//--------------------
// KEY
//--------------------

// Key is the used key to sign a token. The real implementation
// controls signing and verification.
type Key any

// ReadECPrivateKey reads a PEM formated ECDSA private key
// from the passed reader.
func ReadECPrivateKey(r io.Reader) (Key, error) {
	var pemkey bytes.Buffer
	_, err := pemkey.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey.Bytes()); block == nil {
		return nil, fmt.Errorf("cannot decode the PEM")
	}
	var parsed *ecdsa.PrivateKey
	if parsed, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		return nil, fmt.Errorf("cannot parse the ECDSA: %v", err)
	}
	return parsed, nil
}

// ReadECPublicKey reads a PEM encoded ECDSA public key
// from the passed reader.
func ReadECPublicKey(r io.Reader) (Key, error) {
	var pemkey bytes.Buffer
	_, err := pemkey.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey.Bytes()); block == nil {
		return nil, fmt.Errorf("cannot decode the PEM")
	}
	var parsed any
	parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse the ECDSA: %v", err)
		}
		parsed = certificate.PublicKey
	}
	publicKey, ok := parsed.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("passed key is no ECDSA key")
	}
	return publicKey, nil
}

// ReadRSAPrivateKey reads a PEM encoded PKCS1 or PKCS8 private key
// from the passed reader.
func ReadRSAPrivateKey(r io.Reader) (Key, error) {
	var pemkey bytes.Buffer
	_, err := pemkey.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey.Bytes()); block == nil {
		return nil, fmt.Errorf("cannot decode the PEM")
	}
	var parsed any
	parsed, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse the RSA: %v", err)
		}
	}
	privateKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("passed key is no RSA key")
	}
	return privateKey, nil
}

// ReadRSAPublicKey reads a PEM encoded PKCS1 or PKCS8 public key
// from the passed reader.
func ReadRSAPublicKey(r io.Reader) (Key, error) {
	var pemkey bytes.Buffer
	_, err := pemkey.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey.Bytes()); block == nil {
		return nil, fmt.Errorf("cannot decode the PEM")
	}
	var parsed any
	parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse the RSA: %v", err)
		}
		parsed = certificate.PublicKey
	}
	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("passed key is no RSA key")
	}
	return publicKey, nil
}

// EOF
