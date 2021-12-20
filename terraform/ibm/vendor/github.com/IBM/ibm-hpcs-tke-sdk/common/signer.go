//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/19/2021    CLH             Adapt for TKE SDK

package common

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
)


/** Used to create an ASN.1 sequence representing an EC signature */
type ECSignature struct {
	R *big.Int
	S *big.Int
}

/*----------------------------------------------------------------------------*/
/* Signs the input data.  Checks the TKE_SIGNSERV_URL environment variable.   */
/* If set, uses a signing service provided by the user to sign the data.      */
/* Otherwise, assumes signature keys are in files on the local workstation.   */
/*                                                                            */
/* Inputs:                                                                    */
/* dataToSign []byte -- the data to be signed                                 */
/* sigkey string -- identifies the signature key to use                       */
/* sigkeyToken string -- authentication token for the signature key           */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the calculated signature                                         */
/* error -- any error encountered                                             */
/*----------------------------------------------------------------------------*/
func SignWithSignatureKey(dataToSign []byte, sigkey string, sigkeyToken string) ([]byte, error) {

	// Check if the environment variable is set indicating a signing service
	// should be used
	ssURL := os.Getenv("TKE_SIGNSERV_URL")
	if ssURL != "" {
		// Use the signing service to create the signature
		encodedData := base64.StdEncoding.EncodeToString(dataToSign)
		req := CreateSignDataRequest(sigkeyToken, ssURL, sigkey, encodedData)
		sig, err := SubmitSignDataRequest(req)
		if err != nil {
			return nil, err
		}
		decodedSig, err := base64.StdEncoding.DecodeString(sig)
		if err != nil {
			return nil, err
		}
		return decodedSig, nil
	} else {
		// Sign the data using the key in a signature key file
		return SignWithSignatureKeyFile(dataToSign, sigkey, sigkeyToken)
	}
}

/*----------------------------------------------------------------------------*/
/* Signs the input data using the private key in a signature key file.        */
/* The signature key could be either a 2048-bit RSA key or a P521 EC key.     */
/*                                                                            */
/* Inputs:                                                                    */
/* dataToSign []byte -- the data to be signed                                 */
/* sigkey string -- identifies the signature key to use                       */
/* sigkeyToken string -- authentication token for the signature key           */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the calculated signature                                         */
/* error -- any error encountered                                             */
/*----------------------------------------------------------------------------*/
func SignWithSignatureKeyFile(dataToSign []byte, sigkey string, sigkeyToken string) ([]byte, error) {

	// Read the signature key file
	data, err := ioutil.ReadFile(sigkey)
	if err != nil {
		return nil, err
	}
	var skfields map[string]string
	err = json.Unmarshal(data, &skfields)
	if err != nil {
		panic(err)
	}

	// Derive the encryption key from the password and salt field
	aeskey, _ := Derive_aes_key(sigkeyToken, skfields["seaSalt"])

	// Decrypt the signature key
	enckey, err := hex.DecodeString(skfields["enckey"])
	if err != nil {
		panic(err)
	}
	pemBytes, err := Decrypt(enckey, aeskey)
	if err != nil {
		return nil, errors.New("Invalid password.")
	}

	if skfields["keyType"] == "p521ec" {
		// Sign using a P521 EC signature key
		return SignWithP521ECKey(dataToSign, pemBytes, skfields["ski"])
	} else {
		// Sign using a 2048-bit RSA signature key
		return SignWithRSA2048Key(dataToSign, pemBytes, skfields["ski"])
	}
}

/*----------------------------------------------------------------------------*/
/* Signs the input data using a 2048-bit RSA key.                             */
/*                                                                            */
/* A SHA-256 hash is calculated over dataToSign.  This is padded using the    */
/* ANSI X9.31 method, and the result is enciphered using the RSA private key. */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte dataToSign -- the data to be signed                                 */
/* []byte pemBytes -- PEM encoded representation of RSA private key           */
/* string savedSKI -- subject key identifier for the RSA public key from the  */
/*     signature key file, represented as a hexadecimal string                */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the RSA signature.  Padded with leading zeroes if needed to make */
/*     it 256 bytes long.                                                     */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func SignWithRSA2048Key(dataToSign []byte, pemBytes []byte, savedSKI string) ([]byte, error) {

	// Recover the private RSA key
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
		return nil, errors.New("PEM decode of signature key failed.")
	}
	rsaKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}

	// Calculate the Subject Key Identifier	(SKI)
	publicBytes, err := asn1.Marshal(rsaKey.PublicKey)
	// This produces an ASN.1 sequence containing two integers:
	// the modulus and the public exponent.  The Subject Key
	// Identifier is the SHA-256 hash of this.
	ski := sha256.Sum256(publicBytes)

	// Compare the calculated and saved SKIs
	if hex.EncodeToString(ski[:]) != savedSKI {
		return nil, errors.New("Miscompare on saved and calculated Subject Key Identifier.")
	}

	return Signature256(dataToSign, rsaKey), nil
}

/*----------------------------------------------------------------------------*/
/* Calculate a 256 byte RSA signature                                         */
/*----------------------------------------------------------------------------*/
func Signature256(dataToSign []byte, rsaKey *rsa.PrivateKey) []byte {
	hash := sha256.Sum256(dataToSign)
	//------------------------------------------------------------------------//
	// The Go crypto/rsa package doesn't support encryption, decryption,
	// or signing functions without some type of padding.
	//
	// https://stackoverflow.com/questions/40870178/golang-rsa-decrypt-no-padding
	//
	// This Stack Overflow post suggests using the following code to sign
	// something without padding:
	//
	// c := new(big.Int).SetBytes(cipherText)
	// plainText := c.Exp(c, privateKey.D, privateKey.N).Bytes()
	//------------------------------------------------------------------------//
	c := new(big.Int).SetBytes(PadANSIX931(hash[:], 0, len(hash), 2048))
	signature := c.Exp(c, rsaKey.D, rsaKey.N).Bytes()
	if len(signature) < 256 {
		// Add leading zeroes
		newsig := make([]byte, 0)
		for i:=0; i<(256-len(signature)); i++ {
			newsig = append(newsig, 0x00)
		}
		newsig = append(newsig, signature...)
		return newsig
	}
	return signature
}


const X931_PAD_BYTE byte = 0xBB
const X931_SIG_HASH_ID_SHA256 byte = 0x34

/*----------------------------------------------------------------------------*/
/* Adds ANSI X9.31 formatting to the input hash.                              */
/*                                                                            */
/* The goal is to create a result that looks like:                            */
/* 0x6B BB ... BB BA || hash || 0x34 CC                                       */
/*                                                                            */
/* Inputs:                                                                    */
/*    data -- contains the hash to be padded                                  */
/*    offset -- starting offset of the hash to be padded                      */
/*    length -- length in bytes of the hash to be padded                      */
/*    sigbits -- length in bits of the signature                              */
/*                                                                            */
/* Returns the padded hash.                                                   */
/*----------------------------------------------------------------------------*/
func PadANSIX931(data []byte, offset int, length int, sigbits int) []byte {

	// Converted to Go from padAnsiX931 in
	// com.ibm.ccc.smartcard.tke.implementation.Utilities.java

	pad_len := (sigbits - length*8 - 16) / 8
	x931_format := make([]byte, pad_len + length + 2)

	// Add the pad bytes
	x931_format[0] = 0x6B
	for i := 1; i < pad_len-1; i++ {
		x931_format[i] = X931_PAD_BYTE
	}
	x931_format[pad_len-1] = 0xBA

	// Add the hash
	for i := pad_len; i < pad_len+length; i++ {
		x931_format[i] = data[offset+i-pad_len]
	}

	// Add the hash format
	x931_format[pad_len+length] = X931_SIG_HASH_ID_SHA256

	// Add the trailer
	x931_format[pad_len+length+1] = 0xCC

	return x931_format
}

/*----------------------------------------------------------------------------*/
/* Signs the input data using a P521 EC key.                                  */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte dataToSign -- the data to be signed                                 */
/* []byte pemBytes -- PEM encoded representation of EC private key            */
/* string savedSKI -- subject key identifier for the EC public key from the   */
/*     signature key file, represented as a hexadecimal string                */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the EC signature.  This is an ASN.1 sequence containing two      */
/*     integers.                                                              */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func SignWithP521ECKey(dataToSign []byte, pemBytes []byte, savedSKI string) ([]byte, error) {

	// Recover the private EC key
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
		return nil, errors.New("PEM decode of signature key failed.")
	}

	ecKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Calculate the Subject Key Identifier	(SKI)
	ski := CalculateECKeyHash(ecKey.PublicKey)

	// Compare the calculated and saved SKIs
	if hex.EncodeToString(ski[:]) != savedSKI {
		return nil, errors.New("Miscompare on saved and calculated Subject Key Identifier.")
	}

	// Calculate the signature
	hash := sha512.Sum512(dataToSign)
	r, s, err := ecdsa.Sign(rand.Reader, ecKey, hash[:])
	if err != nil {
		return nil, err
	}

	var ecSig ECSignature
	ecSig.R = r
	ecSig.S = s
	// Represent the signature as an ASN.1 sequence
	return asn1.Marshal(ecSig)
}
