//
// Copyright contributors to the ibm-hpcs-tke-sdk project
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 01/09/2025    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
	"math/big"
)

/*----------------------------------------------------------------------------*/
/* Verifies an OA certificate returned by a CEX8P.                            */
/*                                                                            */
/* Adapted from a method of the same name in                                  */
/* com.ibm.tke.model.xcp.XCPCryptoModuleClass                                 */
/*                                                                            */
/* OA certificates for the 4770 contain both ECC and Dilithium public keys    */
/* and signatures.  Without Go language support for Dilithium cryptographic   */
/* algorithms, only the ECC signature in a certificate can be verified.       */
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
func VerifyOA3Certificate(authToken string, urlStart string, de common.DomainEntry,
	certIndex uint32, aCertificate OA3CertificateX) error {

	if common.ByteSlicesAreEqual(aCertificate.EccKeyMetaDataSubjectSKI,
		aCertificate.SignerInfoEccSignerSKI) {
		return errors.New("Self-signed OA certificates are not allowed")
	}

	var parentCert OA3CertificateX
	parentExists := false
	var xbytes, ybytes []byte
	parentCertData, err := QueryDeviceCertificate(authToken, urlStart, de, certIndex+1)
	if err == nil {
		parentExists = true
		err = parentCert.Init(parentCertData)
		if err != nil {
			return err
		}
		if !common.ByteSlicesAreEqual(aCertificate.SignerInfoEccSignerSKI,
		                              parentCert.EccKeyMetaDataSubjectSKI) {
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
		} else if ve.ReturnCode() == 96  &&  ve.ReasonCode() == 60 {
			// Use hardcoded IBM root key to verify signature
			cex8RootECCSKI, err := hex.DecodeString(
				"96c87ef655f67918b4f38c612ed9e492fe152555a63cd5fe2b719ee77b3277f6")
			if err != nil {
				panic(err)
			}
			if common.ByteSlicesAreEqual(aCertificate.SignerInfoEccSignerSKI, cex8RootECCSKI) {
				// Use root key for 4770
				xbytes, err = hex.DecodeString(
					"01b6de79e40847a31b10ed5e2458d70b"+
					"5acb8f2e3c1a78875a99cdee110cdd27"+
					"020174485b0e71deb3668e574181443e"+
					"b3801d0d2060b90ce8f2c09bf0fb346c"+
					"3657")
				if err != nil {
					panic(err)
				}
				ybytes, err = hex.DecodeString(
					"018f9caef6a395b32b03f3327993f55b"+
					"96e0a3c5fab1dc640e033ceaede78b41"+
					"d27400f611494f227c5695abd40d8fac"+
					"dbb30a0a4a63e62bb2867805e85768cb"+
					"cd33")
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
	pubkey := ecdsa.PublicKey{ elliptic.P521(), &x, &y }

	// Calculate the SHA-512 hash of the part of the certificate body
	// signed by the ECC key
	hasher := sha512.New()
	hasher.Write(aCertificate.EccBody)
	sha512hash := hasher.Sum(nil)

	// Verify the ECC signature
	var r, s big.Int
	r.SetBytes(aCertificate.SignerInfoR)
	s.SetBytes(aCertificate.SignerInfoS)
	if !ecdsa.Verify(&pubkey, sha512hash, &r, &s) {
		return errors.New("Invalid signature in OA certificate")
	}

	// If a parent certificate was found, recursively call this function to
	// verify the parent certificate
	if parentExists {
		return VerifyOA3Certificate(authToken, urlStart, de, certIndex+1, parentCert)
	} else {
		return nil
	}
}
