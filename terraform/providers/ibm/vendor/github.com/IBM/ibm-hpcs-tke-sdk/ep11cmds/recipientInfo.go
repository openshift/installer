//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/27/2020    CLH             T390301 - Add minimal touch functions

package ep11cmds

import (
	"errors"
	"strconv"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

// VERSION_3 is declared in signerInfo.go
// ASN1_NULL is declared in signerInfo.go

/** Object identifier for ecPublicKey ( 1 2 840 10045 2 1 ) */
var OID_ecPublicKey = []byte{
	0x06, 0x07, 0x2A, 0x86, 0x48, 0xCE, 0x3D, 0x02, 0x01,
}

/** Object identifier for stdDH-sha156kdf ( 1 3 132 1 11 1 ) */
var OID_stdDH_sha256kdf = []byte{
	0x06, 0x06, 0x2B, 0x81, 0x04, 0x01, 0x0B, 0x01,
}

/** Object identifier for aes256-wrap ( 2 16 840 1 101 3 4 1 45 ) */
var OID_aes256_wrap = []byte{
	0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x01, 0x2D,
}

var A0_TAG byte = 0xA0
var A1_TAG byte = 0xA1

/*
Sample RecipientInfo for P521 EC importer key:
(taken from com.ibm.tke.message.xcp.XCPRecipientInfo)

30820135                           (0)
  020103                           (4)  -- version = 3
  A08198                           (7)  -- [0] originator tag
    A18195                         (10) -- [1] originatorKey choice
      308192                       (13)
        06072A8648CE3D0201         (16) -- OID ecPublicKey ( 1 2 840 10045 2 1 )
        038186                     (25)
          00                       (28)
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF -- EC public key from smart card for key import
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
          FFFFFFFFFF
  A12A                             (162) -- [1] ukm tag (UserKeyingMaterial)
    0428                           (164)
      FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF -- 40-byte random number UKM) used by ECDH
      FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
      FFFFFFFFFFFFFFFF
  3017                             (206) -- KeyEncryptionAlgorithm
    06062B8104010B01               (208) -- OID stdDH-sha256kdf ( 1 3 132 1 11 1 )
    300D                           (216)
      0609                         (218)
      60864801650304012D           (220) -- OID aes256-wrap ( 2 16 840 1 101 3 4 1 45 )
      0500                         (229) -- NULL parameter
  3050                             (231) -- RecipientEncryptedKeys
    304E                           (233) -- RecipientEncryptedKey
      A022                         (235) -- [0] RecipientKeyIdentifier tag
        0420                       (237)
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF -- SKI of importer key
          FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
      0428                         (271)
        FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF -- EncryptedKey (aes-wrap)
        FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
        FFFFFFFFFFFFFFFF
*/

/*----------------------------------------------------------------------------*/
/* Creates the RecipientInfo needed to import a master key part when a        */
/* P521 EC importer key is used.                                              */
/*                                                                            */
/* Inputs:                                                                    */
/* []byte publicKey -- the P521 EC public key (133 bytes)                     */
/* []byte ukm -- user key material, (40 random bytes)                         */
/* []byte ski -- subject key identifier for the P521 EC public key (32 bytes) */
/* []byte encryptedKeyPart -- the encrypted key part (40 bytes)               */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- ASN.1 sequence for the RecipientInfo                             */
/* error -- reports any error encountered                                     */
/*----------------------------------------------------------------------------*/
func CreateRecipientInfoP521EC(publicKey []byte, ukm []byte, ski []byte,
	encryptedKeyPart []byte) ([]byte, error) {

	finalResult := make([]byte, 0)

	if len(publicKey) != 133 {
		return finalResult, errors.New(
			"Error creating recipient info.  Invalid public key length, length = " +
				strconv.Itoa(len(publicKey)) + ".")
	}
	if len(ukm) != 40 {
		return finalResult, errors.New(
			"Error creating recipient info.  Invalid user key material length, length = " +
				strconv.Itoa(len(ukm)) + ".")
	}
	if len(ski) != 32 {
		return finalResult, errors.New(
			"Error creating recipient info.  Invalid subject key identifier length, length = " +
				strconv.Itoa(len(ski)) + ".")
	}
	if len(encryptedKeyPart) != 40 {
		return finalResult, errors.New(
			"Error creating recipient info.  Invalid encrypted key part length, length = " +
				strconv.Itoa(len(encryptedKeyPart)) + ".")
	}

	recipientInfoFields := make([][]byte, 5)

	recipientInfoFields[0] = VERSION_3

	publicKeyFields := make([][]byte, 2)
	publicKeyFields[0] = OID_ecPublicKey
	publicKeyFields[1] = common.Asn1FormBitString(publicKey)
	originatorKeyChoice := common.Asn1FormOctetString(
		common.Asn1FormSequence(publicKeyFields))
	originatorKeyChoice[0] = A1_TAG
	originator := common.Asn1FormOctetString(originatorKeyChoice)
	originator[0] = A0_TAG
	recipientInfoFields[1] = originator

	ukmField := common.Asn1FormOctetString(common.Asn1FormOctetString(ukm))
	ukmField[0] = A1_TAG
	recipientInfoFields[2] = ukmField

	keyEncrAlgorithmFields := make([][]byte, 2)
	keyEncrAlgorithmFields[0] = OID_stdDH_sha256kdf
	aes256_wrapFields := make([][]byte, 2)
	aes256_wrapFields[0] = OID_aes256_wrap
	aes256_wrapFields[1] = ASN1_NULL
	keyEncrAlgorithmFields[1] =
		common.Asn1FormSequence(aes256_wrapFields)
	recipientInfoFields[3] =
		common.Asn1FormSequence(keyEncrAlgorithmFields)

	recipientEncrKeysFields := make([][]byte, 1)
	recipientEncrKeyFields := make([][]byte, 2)
	recipientKeyId := common.Asn1FormOctetString(common.Asn1FormOctetString(ski))
	recipientKeyId[0] = A0_TAG
	recipientEncrKeyFields[0] = recipientKeyId
	recipientEncrKeyFields[1] = common.Asn1FormOctetString(encryptedKeyPart)
	recipientEncrKeysFields[0] = common.Asn1FormSequence(recipientEncrKeyFields)
	recipientInfoFields[4] = common.Asn1FormSequence(recipientEncrKeysFields)

	finalResult = append(finalResult, common.Asn1FormSequence(recipientInfoFields)...)
	return finalResult, nil
}
