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
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/** Represents EP11 SignerInfo in xcpAdminRsp with EC OA signature */
type SignerInfo struct {
	Version              int
	SubjectKeyID         []byte // 32-byte SKI, without ASN.1 tags or lengths
	DigestAlgorithmID    []byte
	SignatureAlgorithmID []byte
	Signature            []byte // signature, ECC signatures are represented as an ASN.1 sequence of two INTEGERs
	SignatureR           []byte
	SignatureS           []byte
}

/** ASN.1 representation of the integer 3 (to designate version = 3) */
var VERSION_3 = []byte{
	0x02, 0x01, 0x03,
}

/** ASN.1 representation of NULL */
var ASN1_NULL = []byte{
	0x05, 0x00,
}

/** Object identifier for digest algorithm: SHA-256 */
var OID_sha256 = []byte{
	0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01,
}

//#B@T372621CLH
/** Object identifier for digest algorithm: SHA-512 */
var OID_sha512 = []byte{
	0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x03,
}

//#E@T372621CLH

/** Object identifier for signature algorithm: rsaEncryption */
var OID_rsaEncryption = []byte{
	0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x01, 0x01,
}

/** Object identifier for signature algorithm: ecdsaWithSHA512 */
var OID_ecdsaWithSHA512 = []byte{
	0x06, 0x08, 0x2A, 0x86, 0x48, 0xCE, 0x3D, 0x04, 0x03, 0x04,
}

/*----------------------------------------------------------------------------*/
/* Signs the input data using a set of signature keys.  A set of concatenated */
/* ASN.1 sequences, one for each signature key, is returned.                  */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte dataToSign -- the data to be signed                                 */
/* []string sigkeys -- identifies the signature keys to be used               */
/* []string sigkeySkis -- the Subject Key Identifiers for the signature keys  */
/* []string sigkeyTokens -- authentication tokens for the signature keys      */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- a set of concatenated ASN.1 structures, one for each signature   */
/* error -- reports any error encountered                                     */
/*----------------------------------------------------------------------------*/
func CreateSignerInfo(dataToSign []byte, sigkeys []string,
	sigkeySkis []string, sigkeyTokens []string) ([]byte, error) {

	// Check if the environment variable is set indicating a signing service
	// should be used
	ssURL := os.Getenv("TKE_SIGNSERV_URL")

	finalResult := make([]byte, 0)
	for i := 0; i < len(sigkeys); i++ {

		var signerInfoFields [][]byte
		var err error
		if ssURL != "" {
			// Only P521 EC keys are supported when signing service is used
			signerInfoFields, err = CreateP521ECSignerInfoFields(
				dataToSign, sigkeys[i], sigkeySkis[i], sigkeyTokens[i])
			if err != nil {
				return nil, err
			}
		} else {
			// Signature key files may contain either 2048-bit RSA or P521 EC
			// signature keys

			// Read the signature key file
			data, err := ioutil.ReadFile(sigkeys[i])
			if err != nil {
				return nil, err
			}
			var skfields map[string]string
			err = json.Unmarshal(data, &skfields)
			if err != nil {
				panic(err)
			}

			if skfields["keyType"] == "p521ec" {
				// Use a P521 EC signature key
				signerInfoFields, err = CreateP521ECSignerInfoFields(
					dataToSign, sigkeys[i], sigkeySkis[i], sigkeyTokens[i])
				if err != nil {
					return nil, err
				}
			} else {
				// Use a 2048-bit RSA signature key
				signerInfoFields, err = Create2048RSASignerInfoFields(
					dataToSign, sigkeys[i], sigkeySkis[i], sigkeyTokens[i])
				if err != nil {
					return nil, err
				}
			}
		}

		// Append this ASN.1 sequence to the final result
		finalResult = append(finalResult, common.Asn1FormSequence(signerInfoFields)...)
	}
	return finalResult, nil
}

/*----------------------------------------------------------------------------*/
/* Creates a set of ASN.1 OCTET STRINGS to be made into an ASN.1 sequence     */
/* that forms the SignerInfo for data signed by a P521 EC signature key.      */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte dataToSign -- the data to be signed                                 */
/* string sigkey -- identifies the signature key to use                       */
/* string sigkeySki -- Subject Key Identifier for the signature key           */
/* string sigkeyToken -- authentication token for the signature key           */
/*                                                                            */
/* Outputs:                                                                   */
/* [][]byte -- a set of ASN.1 OCTET STRINGS that will form SignerInfo         */
/*     containing an EC signature                                             */
/* error -- reports any error encountered                                     */
/*----------------------------------------------------------------------------*/
func CreateP521ECSignerInfoFields(dataToSign []byte, sigkey string,
	sigkeySki string, sigkeyToken string) ([][]byte, error) {

	signerInfoFields := make([][]byte, 5)
	signerInfoFields[0] = VERSION_3

	ski, err := hex.DecodeString(sigkeySki)
	if err != nil {
		return nil, err
	}
	signerInfoFields[1] = common.Asn1FormOctetString(ski)
	signerInfoFields[1][0] = common.ASN1_CONTEXT_SPECIFIC_TAG //hack

	algIdFields := make([][]byte, 2)
	algIdFields[0] = OID_sha512
	algIdFields[1] = ASN1_NULL
	signerInfoFields[2] = common.Asn1FormSequence(algIdFields)

	algIdFields[0] = OID_ecdsaWithSHA512
	// algIdFields[1] is still ASN1_NULL
	signerInfoFields[3] = common.Asn1FormSequence(algIdFields)

	signature, err := common.SignWithSignatureKey(
		dataToSign, sigkey, sigkeyToken)
	if err != nil {
		return nil, err
	}
	signerInfoFields[4] = common.Asn1FormOctetString(signature)

	return signerInfoFields, nil
}

/*----------------------------------------------------------------------------*/
/* Creates a set of ASN.1 OCTET STRINGS to be made into an ASN.1 sequence     */
/* that forms the SignerInfo for data signed by a 2048-bit RSA signature key. */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte dataToSign -- the data to be signed                                 */
/* string sigkey -- identifies the signature key to use                       */
/* string sigkeySki -- Subject Key Identifier for the signature key           */
/* string sigkeyToken -- authentication token for the signature key           */
/*                                                                            */
/* Outputs:                                                                   */
/* [][]byte -- a set of ASN.1 OCTET STRINGS that will form SignerInfo         */
/*     containing an RSA signature                                            */
/* error -- reports any error encountered                                     */
/*----------------------------------------------------------------------------*/
func Create2048RSASignerInfoFields(dataToSign []byte, sigkey string,
	sigkeySki string, sigkeyToken string) ([][]byte, error) {

	signerInfoFields := make([][]byte, 5)
	signerInfoFields[0] = VERSION_3

	ski, err := hex.DecodeString(sigkeySki)
	if err != nil {
		return nil, err
	}
	signerInfoFields[1] = common.Asn1FormOctetString(ski)
	signerInfoFields[1][0] = common.ASN1_CONTEXT_SPECIFIC_TAG //hack

	algIdFields := make([][]byte, 2)
	algIdFields[0] = OID_sha256
	algIdFields[1] = ASN1_NULL
	signerInfoFields[2] = common.Asn1FormSequence(algIdFields)

	algIdFields[0] = OID_rsaEncryption
	// algIdFields[1] is still ASN1_NULL
	signerInfoFields[3] = common.Asn1FormSequence(algIdFields)

	signature, err := common.SignWithSignatureKey(
		dataToSign, sigkey, sigkeyToken)
	if err != nil {
		return nil, err
	}
	signerInfoFields[4] = common.Asn1FormOctetString(signature)

	return signerInfoFields, nil
}

/*----------------------------------------------------------------------------*/
/* Isolates the individual fields of the SignerInfo returned in an            */
/* xcpAdminRsp.                                                               */
/*                                                                            */
/* SignerInfo is an ASN.1 sequence but the sequence elements are not all      */
/* octet strings.  asn1.Decode returns with an error when given SignerInfo    */
/* to process.  This function can be called instead of asn1.Decode.           */
/*                                                                            */
/* Inputs:                                                                    */
/* signerInfoBytes -- the ASN.1 SignerInfo sequence to process                */
/*                                                                            */
/* Outputs:                                                                   */
/* SignerInfo -- structure with entries for each of the SignerInfo fields     */
/* error -- identifies any errors found parsing the input sequence            */
/*----------------------------------------------------------------------------*/
func DecodeSignerInfo(signerInfoBytes []byte) (SignerInfo, error) {

	var signerInfo SignerInfo
	if signerInfoBytes[0] != common.ASN1_SEQUENCE_TAG {
		return signerInfo, errors.New("Invalid SignerInfo, missing SEQUENCE tag")
	}
	seqlen, err := common.Asn1GetLength(signerInfoBytes, 1)
	if err != nil {
		return signerInfo, err
	}
	offset, err := common.Asn1SkipLength(signerInfoBytes, 1)
	if len(signerInfoBytes) != offset+seqlen {
		return signerInfo, errors.New("Invalid SignerInfo, invalid SEQUENCE length")
	}

	// Process the Version field
	if signerInfoBytes[offset] != common.ASN1_INTEGER_TAG {
		return signerInfo, errors.New("Invalid SignerInfo, missing INTEGER tag")
	}
	offset++
	if signerInfoBytes[offset] != 0x01 {
		return signerInfo, errors.New("Invalid SignerInfo, unexpected Version length")
	}
	offset++
	signerInfo.Version = int(signerInfoBytes[offset])

	// Process the Subject Key Identifier field
	offset++
	if signerInfoBytes[offset] != common.ASN1_CONTEXT_SPECIFIC_TAG {
		return signerInfo, errors.New("Invalid SignerInfo, missing CONTEXT_SPECIFIC tag")
	}
	offset++
	if signerInfoBytes[offset] != 0x20 {
		return signerInfo, errors.New("Invalid SignerInfo, unexpected SubjectKeyID length")
	}
	offset++
	signerInfo.SubjectKeyID = signerInfoBytes[offset : offset+32]
	offset = offset + 32

	// Process the DigestAlgorithmID field
	oidBytes, err := common.Asn1GetSequenceBytes(signerInfoBytes, offset)
	if oidBytes[0] != common.ASN1_OID_TAG {
		return signerInfo, errors.New("Expected OID tag not found")
	}
	oidLength, err := common.Asn1GetLength(oidBytes, 1)
	if err != nil {
		return signerInfo, err
	}
	newOffset, err := common.Asn1SkipLength(oidBytes, 1)
	signerInfo.DigestAlgorithmID = oidBytes[0 : newOffset+oidLength]
	offset, err = common.Asn1SkipSequence(signerInfoBytes, offset)
	if err != nil {
		return signerInfo, err
	}

	// Process the SignatureAlgorithmID field
	oidBytes, err = common.Asn1GetSequenceBytes(signerInfoBytes, offset)
	if oidBytes[0] != common.ASN1_OID_TAG {
		return signerInfo, errors.New("Expected OID tag not found")
	}
	oidLength, err = common.Asn1GetLength(oidBytes, 1)
	if err != nil {
		return signerInfo, err
	}
	newOffset, err = common.Asn1SkipLength(oidBytes, 1)
	signerInfo.SignatureAlgorithmID = oidBytes[0 : newOffset+oidLength]
	offset, err = common.Asn1SkipSequence(signerInfoBytes, offset)
	if err != nil {
		return signerInfo, err
	}

	// Process the Signature field
	signatureBytes, err := common.Asn1GetOctetStringBytes(signerInfoBytes, offset)
	if err != nil {
		return signerInfo, err
	}
	signerInfo.Signature = signatureBytes

	// The signature returned by an EP11 crypto module is an ASN.1 sequence
	// of two INTEGERs, R and S.
	bytes, err := common.Asn1GetSequenceBytes(signatureBytes, 0)
	if err != nil {
		return signerInfo, err
	}
	rBytes, err := common.Asn1GetIntegerBytes(bytes, 0)
	if err != nil {
		return signerInfo, err
	}
	signerInfo.SignatureR = rBytes
	offset, err = common.Asn1SkipInteger(bytes, 0)
	if err != nil {
		return signerInfo, err
	}
	sBytes, err := common.Asn1GetIntegerBytes(bytes, offset)
	if err != nil {
		return signerInfo, err
	}
	signerInfo.SignatureS = sBytes

	return signerInfo, nil
}
