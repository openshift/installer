//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/04/2021    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

var baseTemplate = 
  "30820195" +                                       //   (0)
    "A003" +                                         //   (4) -- [0]
      "020102" +                                     //   (6)
    "020100" +                                       //   (9)
    "300A" +                                         //  (12)  *** SIGNATURE ALGORITHM ***
      "06082A8648CE3D040304" +                       //  (14) -- OID ecdsaWithSHA512
    "3046" +                                         //  (24)  *** ISSUER ***
      "311B" +                                       //  (26)
        "3019" +                                     //  (28)
          "060355040B" +                             //  (30) -- OID organizationUnitName
          "1312" +                                   //  (35)
            "544B452041646D696E697374726174696F6E" + //       -- "TKE Administration"
      "3127" +                                       //  (55)
        "3025" +                                     //  (57)
          "0603550403" +                             //  (59) -- OID commonName
          "131E" +                                   //  (64)
            "202020202020202020202020202020" +       //       -- fill in administrator name
            "202020202020202020202020202020" +       // 
    "3022" +                                         //  (96)  *** VALIDITY ***
      "180F" +                                       //  (98)
        "32303030303130313030303030305A" +           //       -- not valid before 2000/1/1
      "180F" +                                       // (115)
        "32313030313233313233353935395A" +           //       -- not valid after 2100/12/31
    "3046" +                                         // (132) *** SUBJECT ***
      "311B" +                                       // (134)
        "3019" +                                     // (136)
          "060355040B" +                             // (138) -- OID organizationUnitName
          "1312" +                                   // (143)
            "544B452041646D696E697374726174696F6E" + //       -- "TKE Administration"
      "3127" +                                       // (163)
        "3025" +                                     // (165)
          "0603550403" +                             // (167) -- OID commonName
          "131E" +                                   // (172)
            "202020202020202020202020202020" +       //       -- fill in administrator name
            "202020202020202020202020202020" +       //
    "30819B" +                                       // (204) *** SUBJECT PUBLIC KEY INFO ***
      "3010" +                                       // (207)
        "06072A8648CE3D0201" +                       // (209) -- OID ecPublicKey
        "06052B81040023" +                           // (218) -- OID secp521r1
      "038186" +                                     // (225)
        "00" +                                       // (228)
        "00000000000000000000000000000000" +         //       -- fill in the EC public key
        "00000000000000000000000000000000" +         //
        "00000000000000000000000000000000" +         //
        "00000000000000000000000000000000" +         //
        "00000000000000000000000000000000" +         //
        "00000000000000000000000000000000" +         //       -- fill in the EC public key
        "00000000000000000000000000000000" +         //
        "00000000000000000000000000000000" +         //
        "0000000000"                       +         //
    "A32D" +                                         // (362) -- [3]
      "302B" +                                       // (364)
        "3029" +                                     // (366)
          "0603551D0E" +                             // (368) -- OID subjectKeyIdentifier
          "0422" +                                   // (373)
            "0420" +                                 // (375)
              "00000000000000000000000000000000" +   // (377) -- fill in the SKI
              "00000000000000000000000000000000"

const admin_name_1_offset = 66
const admin_name_2_offset = 174
const public_key_offset = 228
const ski_offset = 377

/** Used to create an ASN.1 sequence representing an EC signature */
type ECSignature struct {
	R *big.Int
	S *big.Int
}

/*----------------------------------------------------------------------------*/
/* Creates an administrator certificate containing a P521 EC public key.      */
/*                                                                            */
/* Inputs:                                                                    */
/* ecdsa.PrivateKey ecKey -- the EC private key                               */
/* string adminName -- the administrator name                                 */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the administrator certificate                                    */
/* error -- reports any error                                                 */
/*----------------------------------------------------------------------------*/
func CreateAdminCertP521EC(ecKey ecdsa.PrivateKey, adminName string) ([]byte, error) {

	// Copy the template
	certBase, err := hex.DecodeString(baseTemplate)
	if err != nil {
		return nil, err
	}

	// Add the administrator name in two places
	if len(adminName) > 30 {
		return nil, errors.New("Administrator name is too long.")
	}
	nameBuffer := []byte(adminName)
	nameLen := len(nameBuffer)
	if nameLen < 30 {
		nameBuffer = append(nameBuffer, setNewSliceToValue(30-nameLen, 0x20)...)
	}
	copy(certBase[admin_name_1_offset:admin_name_1_offset+30], nameBuffer[:])
	copy(certBase[admin_name_2_offset:admin_name_2_offset+30], nameBuffer[:])

	// Add the EC public key
	certBase[public_key_offset + 1] = 0x04  // compression byte
	bytes := ecKey.PublicKey.X.Bytes()
	length := len(bytes)
	// copy the X coordinate
	for i:=0; i<length; i++ {
		certBase[public_key_offset + 2 + (66 - length) + i] = bytes[i]
	}
	bytes = ecKey.PublicKey.Y.Bytes()
	length = len(bytes)
	// copy the Y coordinate
	for i:=0; i<length; i++ {
		certBase[public_key_offset + 2 + 66 + (66 - length) + i] = bytes[i]
	}

	// Add the subject key identifier
	ski := common.CalculateECKeyHash(ecKey.PublicKey)
	copy(certBase[ski_offset:ski_offset+32], ski[:])

	// Calculate the signature
	hash := sha512.Sum512(certBase)
	r, s, err := ecdsa.Sign(rand.Reader, &ecKey, hash[:])
	if err != nil {
		return nil, err
	}

	var ecSig ECSignature
	ecSig.R = r
	ecSig.S = s
	// Represent the signature as an ASN.1 sequence
	signature, err := asn1.Marshal(ecSig)
	if err != nil {
		return nil, err
	}

	// Assemble the final certificate
	elements := make([][]byte, 3)
	elements[0] = certBase
	signatureOid := make([][]byte, 1)
	signatureOid[0] = OID_ecdsaWithSHA512
	elements[1] = common.Asn1FormSequence(signatureOid)
	elements[2] = common.Asn1FormBitString(signature)
	return common.Asn1FormSequence(elements), nil
}

/*----------------------------------------------------------------------------*/
/* Creates an administrator certificate containing a P521 EC public key       */
/* using a signing service.                                                   */
/*                                                                            */
/* Inputs:                                                                    */
/* string -- base URL for the signing service                                 */
/* string -- identifies the signature key to be used                          */
/* string -- authentication token for the signature key                       */
/* string -- administrator name to be placed in the certificate               */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- an administrator certificate containing the EC public key        */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func CreateAdminCertUsingSigningService(ssURL string, sigkey string,
	sigkeyToken string, adminName string) ([]byte, error) {

	// Copy the template
	certBase, err := hex.DecodeString(baseTemplate)
	if err != nil {
		return nil, err
	}

	// Add the administrator name in two places
	if len(adminName) > 30 {
		return nil, errors.New("Administrator name is too long.")
	}
	nameBuffer := []byte(adminName)
	nameLen := len(nameBuffer)
	if nameLen < 30 {
		nameBuffer = append(nameBuffer, setNewSliceToValue(30-nameLen, 0x20)...)
	}
	copy(certBase[admin_name_1_offset:admin_name_1_offset+30], nameBuffer[:])
	copy(certBase[admin_name_2_offset:admin_name_2_offset+30], nameBuffer[:])

	// Add the EC public key

	// Use the signing service to get the EC public key
	pubkey, err := common.GetPublicKeyFromSigningService(
		ssURL, sigkey, sigkeyToken)
	if err != nil {
		return make([]byte, 0), err
	}

	// Add the EC public key
	copy(certBase[public_key_offset + 1 : public_key_offset + 1 + 133], pubkey)

	// Add the subject key identifier
	ski := sha256.Sum256(pubkey[:])
	copy(certBase[ski_offset:ski_offset+32], ski[:])

	// Calculate the signature
	hash := sha512.Sum512(certBase)
	signature, err := common.SignWithSignatureKey(hash[:], sigkey, sigkeyToken)
	if err != nil {
		return make([]byte, 0), err
	}

	// Assemble the final certificate
	elements := make([][]byte, 3)
	elements[0] = certBase
	signatureOid := make([][]byte, 1)
	signatureOid[0] = OID_ecdsaWithSHA512
	elements[1] = common.Asn1FormSequence(signatureOid)
	elements[2] = common.Asn1FormBitString(signature)
	return common.Asn1FormSequence(elements), nil
}
