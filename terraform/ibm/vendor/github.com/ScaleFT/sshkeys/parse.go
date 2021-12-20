// Portions of this file are based on https://github.com/golang/crypto/blob/master/ssh/keys.go
//
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sshkeys

import (
	"crypto/x509"

	"golang.org/x/crypto/ssh"
)

// ErrIncorrectPassword is returned when the supplied passphrase was not correct for an encrypted private key.
var ErrIncorrectPassword = x509.IncorrectPasswordError

// ParseEncryptedPrivateKey returns a Signer from an encrypted private key. It supports
// the same keys as ParseEncryptedRawPrivateKey.
func ParseEncryptedPrivateKey(data []byte, passphrase []byte) (ssh.Signer, error) {
	key, err := ParseEncryptedRawPrivateKey(data, passphrase)
	if err != nil {
		return nil, err
	}

	return ssh.NewSignerFromKey(key)
}

// ParseEncryptedRawPrivateKey returns a private key from an encrypted private key. It
// supports RSA (PKCS#1 or OpenSSH), DSA (OpenSSL), and ECDSA private keys.
//
// ErrIncorrectPassword will be returned if the supplied passphrase is wrong,
// but some formats like RSA in PKCS#1 detecting a wrong passphrase is difficult,
// and other parse errors may be returned.
func ParseEncryptedRawPrivateKey(data []byte, passphrase []byte) (interface{}, error) {
	if passphrase == nil {
		return ssh.ParseRawPrivateKey(data)
	}
	return ssh.ParseRawPrivateKeyWithPassphrase(data, passphrase)
}
