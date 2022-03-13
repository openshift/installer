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
/* Verifies an OA certificate returned by a CEX6P or CEX7P.                   */
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
func VerifyOA2Certificate(authToken string, urlStart string, de common.DomainEntry,
	certIndex uint32, aCertificate OA2CertificateX) error {

	if common.ByteSlicesAreEqual(aCertificate.MetaDataSubjectSKI, aCertificate.MetaDataSignerSKI) {
		return errors.New("Self-signed OA certificates are not allowed")
	}

	var parentCert OA2CertificateX
	parentExists := false
	var xbytes, ybytes []byte
	parentCertData, err := QueryDeviceCertificate(authToken, urlStart, de, certIndex+1)
	if err == nil {
		parentExists = true
		err = parentCert.Init(parentCertData)
		if err != nil {
			return err
		}
		if !common.ByteSlicesAreEqual(aCertificate.MetaDataSignerSKI,
			parentCert.MetaDataSubjectSKI) {
			return errors.New("Incorrect parent certificate fetched in OA certificate chain.")
		}
		// Use public key from parent certificate to verify signature
		xbytes = parentCert.SpkiXCoordinate
		ybytes = parentCert.SpkiYCoordinate
	} else {
		// Look for certificate not found error
		ve, ok := err.(VerbError)
		if !ok {
			return err
		} else if ve.ReturnCode() == 96 && ve.ReasonCode() == 60 {
			// Use hardcoded IBM root key to verify signature
			cex6RootSKI, err := hex.DecodeString(
				"5f4717480b75b8bcbc224c62dfbd7a3b8987e54193d40a1a664fa221e5387123")
			if err != nil {
				panic(err)
			}
			cex7RootSKI, err := hex.DecodeString(
				"b23884d4dea35eaabab109e0a6c400b2f5bea652f14ab23d2a47220696eb0e6b")
			if err != nil {
				panic(err)
			}
			if common.ByteSlicesAreEqual(aCertificate.MetaDataSignerSKI, cex6RootSKI) {
				// Use root key for 4768
				xbytes, err = hex.DecodeString(
					"018bd55521ae1d52a07630d04cd1adeb" +
						"82f22fea54da7ccf5d4e33e0c5ff0755" +
						"19f4d9d902a6800d47dc2c0e5fd4f6eb" +
						"e63096f37fadf9717b8a679d0fb1d595" +
						"7cdc")
				if err != nil {
					panic(err)
				}
				ybytes, err = hex.DecodeString(
					"0116f34a637240bf0abcef0528cb8eb4" +
						"cdde17154010d10692c753fdadebcfb4" +
						"52130fa664bc8546c1bd72211297ad00" +
						"38836a4e55182530d6ccb10ffa2004f7" +
						"f9cd")
				if err != nil {
					panic(err)
				}
			} else if common.ByteSlicesAreEqual(aCertificate.MetaDataSignerSKI, cex7RootSKI) {
				// Use root key for 4769
				xbytes, err = hex.DecodeString(
					"0112ec7b6c304f2ded0dcbc0d8ec3c43" +
						"ee2fb377052a46517747bfe5ad35e68f" +
						"50016644a438cc6c20c4d1b25085eea9" +
						"c83a99514d9b39b6e3d60f472237ac89" +
						"4d81")
				if err != nil {
					panic(err)
				}
				ybytes, err = hex.DecodeString(
					"002b53cdf6a39ca67e258a0e714d6609" +
						"d7ea3f687aa0b3bf8b69f92e7dcf5397" +
						"92abb8b063be8f9393124c10271d4c56" +
						"9626cb6e8ea79f502670c201848e3280" +
						"41a7")
				if err != nil {
					panic(err)
				}
			} else {
				return errors.New("Unrecognized IBM root key in OA certificate chain.")
			}
		} else {
			return err
		}
	}
	var x, y big.Int
	x.SetBytes(xbytes)
	y.SetBytes(ybytes)
	pubkey := ecdsa.PublicKey{Curve: elliptic.P521(), X: &x, Y: &y}

	// Calculate the SHA-512 hash of the certificate body
	hasher := sha512.New()
	hasher.Write(aCertificate.Body)
	sha512hash := hasher.Sum(nil)

	// Verify the signature
	var r, s big.Int
	r.SetBytes(aCertificate.SignerInfoR)
	s.SetBytes(aCertificate.SignerInfoS)
	if !ecdsa.Verify(&pubkey, sha512hash, &r, &s) {
		return errors.New("Invalid signature in OA certificate")
	}

	// If a parent certificate was found, recursively call this function to
	// verify the parent certificate
	if parentExists {
		return VerifyOA2Certificate(authToken, urlStart, de, certIndex+1, parentCert)
	} else {
		return nil
	}
}
