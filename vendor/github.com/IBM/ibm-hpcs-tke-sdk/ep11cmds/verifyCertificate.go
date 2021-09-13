//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/12/2021    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Verifies an OA certificate returned by a CEX5P.                            */
/*                                                                            */
/* Adapted from a method of the same name in                                  */
/* com.ibm.tke.model.xcp.XCPCryptoModuleClass                                 */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies a domain assigned to the user.  The OA           */
/*    certificate chain for the crypto module containing that domain is to    */
/*    be verified.                                                            */
/* certIndex -- index of the OA certificate to start with.  0 = currently     */
/*    active epoch key, 1 = its parent, etc.  This method calls itself        */
/*    recursively to read and verify the entire OA certificate chain.         */
/* aCertificate -- the OA certificate to be verified                          */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func VerifyCertificate(authToken string, urlStart string, de common.DomainEntry,
	certIndex uint32, aCertificate OACertificateX) error {

	if aCertificate.HeaderTData != OA_NEW_CERT {
		return errors.New("Invalid OA certificate")
	}

	if aCertificate.BodyTPublic == OA_RSA {
		return errors.New("The plug-in does not support OA certificates with an RSA public key")
	} else if aCertificate.BodyTPublic == OA_ECC {
		return verifyECCCertificate(authToken, urlStart, de, certIndex, aCertificate)
	} else {
		return errors.New("Unrecognized OA certificate key type")
	}
}

/*----------------------------------------------------------------------------*/
/* Verifies an OA certificate containing an ECC public key.                   */
/*                                                                            */
/* Has the same inputs and outputs as the previous function.                  */
/*----------------------------------------------------------------------------*/
func verifyECCCertificate(authToken string, urlStart string, de common.DomainEntry,
	certIndex uint32, aCertificate OACertificateX) error {

	if common.ByteSlicesAreEqual(aCertificate.BodyCkoName, aCertificate.BodyParentName) {
		return errors.New("Self-signed OA certificates are not allowed")
	}

	var parentCert OACertificateX
	var parentExists bool = false
	var xbytes, ybytes []byte
	parentCertData, err := QueryDeviceCertificate(authToken, urlStart, de, certIndex+1)
	if err == nil {
		parentExists = true
		err = parentCert.Init(parentCertData)
		if err != nil {
			return err
		}
		// Use public key from parent certificate to verify signature
		if (parentCert.ECCCurveType != ECC_PRIME) || (parentCert.ECCCurveSize != 521) {
			return errors.New("Unsupported ECC curve type or size in OA certificate")
		}
		xbytes = parentCert.ECCPublicKeyX
		ybytes = parentCert.ECCPublicKeyY
	} else {
		// Look for certificate not found error
		ve, ok := err.(VerbError)
		if !ok {
			return err
		} else if ve.ReturnCode() == 96 &&
			aCertificate.BodyCkoNameNameType == 0x9000 &&
			aCertificate.BodyCkoNameIndex > 0 {
			// Use hardcoded IBM root key to verify signature
			xbytes, err = hex.DecodeString(
				"01CF9238348503B59D5FE207467DDA6E" +
					"D13E6593BE42376523BF2CD65FC37729" +
					"66FD90F210D4B96046927418037C8534" +
					"0D6E98D97551656E89F8650F9DEC54D5" +
					"7C2D")
			if err != nil {
				panic(err)
			}
			ybytes, err = hex.DecodeString(
				"00DDFC3D77761BD913C9165534534904" +
					"C36835DCDADA6DD8C8FD6A226E7EE363" +
					"98DF81EC4C534996993EBD6EF421410E" +
					"FABB09286BB2212CEAD62CFA1D0E5D6F" +
					"26E7")
			if err != nil {
				panic(err)
			}
		} else {
			return err
		}
	}
	var x, y big.Int
	x.SetBytes(xbytes)
	y.SetBytes(ybytes)
	pubkey := ecdsa.PublicKey{elliptic.P521(), &x, &y}

	// Calculate the SHA-512 hash of the certificate body
	hasher := sha512.New()
	hasher.Write(aCertificate.Body)
	sha512hash := hasher.Sum(nil)

	// Verify the signature
	var r, s big.Int
	r.SetBytes(aCertificate.ECCSignatureR)
	s.SetBytes(aCertificate.ECCSignatureS)
	if !ecdsa.Verify(&pubkey, sha512hash, &r, &s) {
		return errors.New("Invalid signature in OA certificate")
	}

	// If a parent certificate was found, recursively call this function to
	// verify the parent certificate
	if parentExists {
		return VerifyCertificate(authToken, urlStart, de, certIndex+1, parentCert)
	} else {
		return nil
	}
}
