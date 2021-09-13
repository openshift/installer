//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/26/2020    CLH             T372621 - Support P521 EC signature keys

package ep11cmds

import (
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
	"github.com/Logicalis/asn1"
)


type CertificateRSA2048 struct {
	TheBody   CertBodyRSA2048
	Algorithm AlgID
	Signature []byte `asn1:"universal,tag:3"`
}

type CertBodyRSA2048 struct {
	TheVersion    CertVersion `asn1:"tag:0"`
	SerialNumber  int
	Algorithm     AlgID
	TheIssuer     Names
	Validity      TimeRange
	TheSubject    Names
	ThePublicKey  PublicKeyRSA2048
	TheExtensions Extensions `asn1:"tag:3"`
}

type PublicKeyRSA2048 struct { // used for RSA 2048 encode
	Algorithm    AlgIDRSA2048
	ThePublicKey []byte `asn1:"universal,tag:3"`
}

type AlgIDRSA2048 struct { // used for RSA 2048 encode
	ObjID  asn1.Oid
	Null   asn1.Null // not optional in this case
	ObjID2 asn1.Oid  `asn1:"optional"` // will always be nil for encode, here so we can copy
}

type ModAndExp struct {
	Modulus  []byte `asn1:"universal,tag:2"`
	Exponent int
}

/*----------------------------------------------------------------------------*/
/* Initialize an administrator certificate that will contain a 2048-bit RSA   */
/* public key.                                                                */
/*                                                                            */
/* Operates on CertificateRSA2048 (both input and output).                    */
/*----------------------------------------------------------------------------*/
func (cert *CertificateRSA2048) Initialize() {
	cert.TheBody.TheVersion.Version = 2
	cert.TheBody.SerialNumber = 0

	cert.TheBody.Algorithm.ObjID = SHA256WithRSAEncryption

	orgUnitName := []uint{2, 5, 4, 11} // 2.5.4.11
	cert.TheBody.TheIssuer.OrgName.TheName.OID = orgUnitName
	cert.TheBody.TheIssuer.OrgName.TheName.PrintableString = []byte(tkeAdministration)

	commonName := []uint{2, 5, 4, 3} // 2.5.4.3
	cert.TheBody.TheIssuer.CommonName.TheName.OID = commonName
	cert.TheBody.TheIssuer.CommonName.TheName.PrintableString = setNewSliceToValue(30, 0x20)

	notBefore := []byte{0x32, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x5A} // -- 20000101 000000 Z -- not valid before 2000/1/1
	notAfter := []byte{0x32, 0x31, 0x30, 0x30, 0x31, 0x32, 0x33, 0x31, 0x32, 0x33, 0x35, 0x39, 0x35, 0x39, 0x5A}  // -- 21001231 235959 Z  -- not valid after 2100/12/31
	cert.TheBody.Validity.NotBefore = notBefore
	cert.TheBody.Validity.NotAfter = notAfter

	cert.TheBody.TheSubject.OrgName.TheName.OID = orgUnitName
	cert.TheBody.TheSubject.OrgName.TheName.PrintableString = []byte(tkeAdministration)

	cert.TheBody.TheSubject.CommonName.TheName.OID = commonName
	cert.TheBody.TheSubject.CommonName.TheName.PrintableString = setNewSliceToValue(30, 0x20)

	cert.TheBody.ThePublicKey.Algorithm.ObjID = RSAEncryption
	pubKey := []byte{0x00, 0x30, 0x82, 0x01, 0x0a, 0x02, 0x82, 0x01, 0x01, 00}
	pubExponent := []byte{0x02, 0x03, 0x01, 0x00, 0x01}
	publicModulus := setNewSliceToValue(256, 0xff)
	pubKey = append(pubKey, publicModulus...)
	pubKey = append(pubKey, pubExponent...)
	cert.TheBody.ThePublicKey.ThePublicKey = pubKey

	subjectKeyID := []uint{2, 5, 29, 14}
	cert.TheBody.TheExtensions.TheSeq1.TheSeq2.ObjID = subjectKeyID
	theSKI := []byte{0x04, 0x20}
	skiValue := setNewSliceToValue(32, 0xff)
	theSKI = append(theSKI, skiValue...)
	cert.TheBody.TheExtensions.TheSeq1.TheSeq2.SKI = theSKI

	cert.Algorithm.ObjID = SHA256WithRSAEncryption

	// initialize signature
	cert.Signature = setNewSliceToValue(257, 0xff)
	cert.Signature[0] = 0
}

/*----------------------------------------------------------------------------*/
/* Sets the administrator name in an administrator certificate containing a   */
/* 2048-bit RSA public key.                                                   */
/*                                                                            */
/* Operates on CertificateRSA2048 (both input and output).                    */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte newName -- the administrator name to be set in the certificate.     */
/*     Can be up to 30 bytes long.                                            */
/*----------------------------------------------------------------------------*/
func (cert *CertificateRSA2048) SetAdminName(newName []byte) {
	if len(newName) > 30 {
		panic("Administrator name is too long")
	}
	nameBuffer := newName
	nameLen := len(nameBuffer)
	if nameLen < 30 {
		nameBuffer = append(nameBuffer, setNewSliceToValue(30-nameLen, 0x20)...)
	}
	cert.TheBody.TheIssuer.CommonName.TheName.PrintableString = nameBuffer
	cert.TheBody.TheSubject.CommonName.TheName.PrintableString = nameBuffer
}

/*----------------------------------------------------------------------------*/
/* Sets the 2048-bit RSA public key and associated subject key identifier in  */
/* an administrator certificate.                                              */
/*                                                                            */
/* Operates on CertificateRSA2048 (both input and output).                    */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte publicKey -- the RSA public key to be set in the certificate.       */
/*     Actually, just the modulus.  The public exponent is assumed to be      */
/*     65537.  Must be 257 bytes long.                                        */
/*----------------------------------------------------------------------------*/
func (cert *CertificateRSA2048) SetPublicKey(publicKey []byte) {
	if len(publicKey) != 257 {
		panic("An RSA public key must be 257 bytes in length")
	}
	// offset 9 into PublicKeyRSA2048.ThePublicKey
	copy(cert.TheBody.ThePublicKey.ThePublicKey[9:], publicKey)
	// calculate and set the SKI based on the public modulus and exponent
	ski := calculateRSAKeyHash(publicKey, 65537)
	encodedSKI, err := asn1.Encode(ski)
	if err != nil {
		panic(err)
	}
	cert.TheBody.TheExtensions.TheSeq1.TheSeq2.SKI = encodedSKI
}

/*----------------------------------------------------------------------------*/
/* Sets the signature field in an administrator certificate containing a      */
/* 2048-bit RSA public key.                                                   */
/*                                                                            */
/* Operates on CertificateRSA2048 (both input and output).                    */
/*                                                                            */
/* Inputs:                                                                    */
/* rsa.PrivateKey rsaKey -- the RSA private key to use to create the          */
/*     signature.                                                             */
/*----------------------------------------------------------------------------*/
func (cert *CertificateRSA2048) SetSignature(rsaKey *rsa.PrivateKey) {
	encoded, err := asn1.Encode(*cert)
	if err != nil {
		panic(err)
	}
	var bytesToSign []byte
	copy(bytesToSign, encoded[4:550])

	signature := common.Signature256(bytesToSign, rsaKey)
	if len(signature) == 256 {
		cert.Signature = setNewSliceToValue(257, 0)
		copy(cert.Signature[1:], signature)
	} else if len(signature) == 257 {
		cert.Signature = signature
	} else {
		panic(errors.New("Signature length is not valid"))
	}
}

/*----------------------------------------------------------------------------*/
/* Calculates the RSA key hash for a given modulus and public exponent.       */
/* 2048-bit RSA public key.                                                   */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte modulus -- the modulus of the RSA key                               */
/* []byte publicExponent -- the public exponent of the RSA key                */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the RSA key hash.  This is the SHA-256 hash of the ASN.1 encoded */
/*     modulus and public exponent.                                           */
/*----------------------------------------------------------------------------*/
func calculateRSAKeyHash(modulus []byte, publicExponent int) []byte {
	var me ModAndExp

	me.Modulus = modulus
	me.Exponent = publicExponent
	encoded, err := asn1.Encode(me)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(encoded)
	return hash[0:]
}
