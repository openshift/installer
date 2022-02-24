//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/09/2021    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"math/big"
	"strconv"

	"github.com/Logicalis/asn1"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

// Request to generate 2048-bit RSA importer key
const XCP_IMPRKEY_RSA_2048 = 0 //@T390301CLH
// Request to generate P521 EC importer key
const XCP_IMPRKEY_EC_P521 = 3 //@T390301CLH

/*----------------------------------------------------------------------------*/
/* Generates a 2048-bit RSA importer key.                                     */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose attributes are to be set        */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* rsa.PublicKey -- the public part of the generated 2048-bit RSA key         */
/* []byte -- the Subject Key Identifier of the RSA public key                 */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func Generate2048RSAImporterKey(authToken string, urlStart string,
	de common.DomainEntry, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (rsa.PublicKey, []byte, error) {

	var pubKey rsa.PublicKey
	var ski []byte

	htpRequestString, err := GenerateImporterKeyRequest(
		authToken, urlStart, de, XCP_IMPRKEY_RSA_2048, sigkeys, sigkeySkis,
		sigkeyTokens)
	if err != nil {
		return pubKey, ski, err
	}

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return pubKey, ski, err
	}

	return Generate2048RSAImporterKeyResponse(htpResponseString, de)
}

/*----------------------------------------------------------------------------*/
/* Parse a generate importer key response for an RSA 2048 importer key.       */
/*                                                                            */
/* Returns the RSA public importer key and its SKI, and any error             */
/*----------------------------------------------------------------------------*/
func Generate2048RSAImporterKeyResponse(htpResponse string, de common.DomainEntry) (rsa.PublicKey, []byte, error) {

	var pubKey rsa.PublicKey
	var ski []byte

	rspBlk, err := buildAdminRspBlk(htpResponse, de)
	if err != nil {
		return pubKey, ski, err
	}

	var pubKeyASN1 PublicKeyRSA2048
	_, err = asn1.Decode(rspBlk.CmdOutput, &pubKeyASN1)
	if err != nil {
		return pubKey, ski, err
	}

	var innerPubKey []byte
	innerPubKey = append(innerPubKey, pubKeyASN1.ThePublicKey[1:]...)
	tempSKI := sha256.Sum256(innerPubKey)
	ski = tempSKI[:]

	var modAndExp ModAndExp
	_, err = asn1.Decode(innerPubKey, &modAndExp)
	if err != nil {
		return pubKey, ski, err
	}

	var pubModulus big.Int
	pubModulus.SetBytes(modAndExp.Modulus)
	pubKey.N = &pubModulus
	pubKey.E = modAndExp.Exponent
	return pubKey, ski, nil
}

/*----------------------------------------------------------------------------*/
/* Initializes an RSA 2048 recipient info structure, including the encrypting */
/* key's SKI and the encrypted key part.                                      */
/*----------------------------------------------------------------------------*/
func (info *RecipientInfo2048) Initialize(ski []byte, encrKey []byte) {
	info.Version = 2
	info.SKI = ski
	info.AlgID.ObjID = RSAEncryption
	info.EncryptedKey = encrKey
}

/*----------------------------------------------------------------------------*/
/* Generates a P521 EC importer key.                                          */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose attributes are to be set        */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* ecdsa.PublicKey -- the public part of the generated P521 EC key            */
/* []byte -- the Subject Key Identifier of the EC public key                  */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func GenerateP521ECImporterKey(authToken string, urlStart string,
	de common.DomainEntry, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (ecdsa.PublicKey, []byte, error) {

	var pubKey ecdsa.PublicKey
	var ski []byte

	htpRequestString, err := GenerateImporterKeyRequest(
		authToken, urlStart, de, XCP_IMPRKEY_EC_P521, sigkeys, sigkeySkis,
		sigkeyTokens)
	if err != nil {
		return pubKey, ski, err
	}

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return pubKey, ski, err
	}

	return GenerateP521ECImporterKeyResponse(htpResponseString, de)
}

/*----------------------------------------------------------------------------*/
/* Parse the response from a generate importer key request when a P521 EC     */
/* key was requested.                                                         */
/*                                                                            */
/* Returns the EC public importer key, its SKI, and any error                 */
/*----------------------------------------------------------------------------*/
func GenerateP521ECImporterKeyResponse(htpResponse string, de common.DomainEntry) (ecdsa.PublicKey, []byte, error) {

	var pubKey ecdsa.PublicKey
	var ski []byte

	rspBlk, err := buildAdminRspBlk(htpResponse, de)
	if err != nil {
		return pubKey, ski, err
	}

	var pubKeyASN1 PublicKeyECP521
	_, err = asn1.Decode(rspBlk.CmdOutput, &pubKeyASN1)
	if err != nil {
		return pubKey, ski, err
	}

	length := len(pubKeyASN1.ThePublicKey)
	if length != 134 {
		panic("Unexpected length for P521 EC public key, length = " + strconv.Itoa(length))
	}

	pubKey.Curve = elliptic.P521()
	var X, Y big.Int
	X.SetBytes(pubKeyASN1.ThePublicKey[2:68])
	Y.SetBytes(pubKeyASN1.ThePublicKey[68:])
	pubKey.X = &X
	pubKey.Y = &Y

	tempSKI := sha256.Sum256(pubKeyASN1.ThePublicKey[2:])
		// skip the leading 0x00 and 0x04 bytes when calculating the SKI
	ski = tempSKI[:]

	return pubKey, ski, nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for generating a domain importer key                */
/*----------------------------------------------------------------------------*/
func GenerateImporterKeyRequest(authToken string, urlStart string, de common.DomainEntry,
	importerKeyType uint32, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_GEN_IMPORTER
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	adminBlk.CmdInput = common.Uint32To4ByteSlice(importerKeyType)

	return CreateSignedHTPRequest(authToken, urlStart, de, adminBlk, sigkeys,
		sigkeySkis, sigkeyTokens)
}
