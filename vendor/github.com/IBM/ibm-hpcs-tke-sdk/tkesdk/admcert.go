//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/30/2021    CLH             Initial version

package tkesdk

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/Logicalis/asn1"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
	"github.com/IBM/ibm-hpcs-tke-sdk/ep11cmds"
)

/*----------------------------------------------------------------------------*/
/* Creates an administrator certificate using the signature key in a file on  */
/* the local workstation.  The file can contain either a 2048-bit RSA key or  */
/* a P521 EC key.                                                             */
/*                                                                            */
/* Inputs:                                                                    */
/* string sigkey -- the full path and name of the signature key file          */
/* string ski -- the Subject Key Identifier of the signature key,             */
/*     represented as a hexadecimal string                                    */
/* string sigkeyToken -- the file password                                    */
/* string adminName -- administrator name                                     */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- an administrator certificate containing the public key for the   */
/*     signature key                                                          */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func CreateAdminCertFromFile(sigkey string, ski string,
	sigkeyToken string, adminName string) ([]byte, error) {

	// Initialize an output variable
	certBytes := make([]byte, 0)

	// Read the signature key file
	data, err := ioutil.ReadFile(sigkey)
	if err != nil {
		return certBytes, err
	}
	var fields map[string]string
	err = json.Unmarshal(data, &fields)
	if err != nil {
		return certBytes, err
	}

	// Derive the encryption key from the file password and salt field
	aeskey, _ := common.Derive_aes_key(sigkeyToken, fields["seaSalt"])

	// Decrypt the signature key
	enckey, err := hex.DecodeString(fields["enckey"])
	if err != nil {
		return certBytes, err
	}
	pemBytes, err := common.Decrypt(enckey, aeskey)
	if err != nil {
		return certBytes, errors.New("Invalid password for a signature key file")
	}

	// Process the signature key
	if fields["keyType"] == "p521ec" {
		// Process a P521 EC key
		certBytes, err = CreateAdminCertForECKey(pemBytes, ski, adminName)
		if err != nil {
			return certBytes, err
		}
	} else {
		// Process a 2048-bit RSA key
		certBytes, err = CreateAdminCertForRSAKey(pemBytes, ski, adminName)
		if err != nil {
			return certBytes, err
		}
	}
	return certBytes, nil
}

/*----------------------------------------------------------------------------*/
/* Creates an administrator certificate containing a 2048-bit RSA public key  */
/* using the PEM representation of an RSA private key.                        */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte pemBytes -- PEM encoded representation of RSA private key           */
/* string savedSKI -- subject key identifier for the RSA public key from the  */
/*     signature key file, represented as a hexadecimal string                */
/* string adminName -- administrator name                                     */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- an administrator certificate containing the RSA public key       */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func CreateAdminCertForRSAKey(pemBytes []byte, savedSKI string, adminName string) ([]byte, error) {

	// Recover the private RSA key
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
		return nil, errors.New("PEM decode of signature key failed.")
	}

	rsaKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Calculate the Subject Key Identifier	(SKI)
	ski, err := common.CalculateRSAKeyHash(rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}

	// Compare the calculated and saved SKIs
	if hex.EncodeToString(ski) != savedSKI {
		return nil, errors.New("Miscompare on saved and calculated Subject Key Identifier.")
	}

	// Create an administrator certificate with the RSA public key
	var cert ep11cmds.CertificateRSA2048
	cert.Initialize()
	cert.SetAdminName([]byte(adminName))
	pubKey := []byte{0}
	pubKey = append(pubKey, rsaKey.PublicKey.N.Bytes()...)
	cert.SetPublicKey(pubKey)
	cert.SetSignature(rsaKey)

	return asn1.Encode(cert)
}

/*----------------------------------------------------------------------------*/
/* Creates an administrator certificate containing a P521 EC public key       */
/* using the PEM representation of an EC private key.                         */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte pemBytes -- PEM encoded representation of EC private key            */
/* string savedSKI -- subject key identifier for the EC public key from the   */
/*     signature key file, represented as a hexadecimal string                */
/* string adminName -- administrator name                                     */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- an administrator certificate containing the EC public key        */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func CreateAdminCertForECKey(pemBytes []byte, savedSKI string, adminName string) ([]byte, error) {

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
	ski := common.CalculateECKeyHash(ecKey.PublicKey)

	// Compare the calculated and saved SKIs
	if hex.EncodeToString(ski[:]) != savedSKI {
		return nil, errors.New("Miscompare on saved and calculated Subject Key Identifier.")
	}

	// Create an administrator certificate with the EC public key
	return ep11cmds.CreateAdminCertP521EC(*ecKey, adminName)
}
