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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/** Used to work with an ASN.1 sequence representing an EC public key */
type ECPublicKey struct {
	X *big.Int
	Y *big.Int
}

/*----------------------------------------------------------------------------*/
/* Assembles information on the signature keys identified in the Terraform    */
/* resource block.                                                            */
/*                                                                            */
/* Handles both signature key files on the local workstation and a            */
/* user-provided signing service.                                             */
/*                                                                            */
/* Inputs:                                                                    */
/* HsmConfig -- A structure containing information from the hsm_config        */
/*     section of the resource block for the HPCS service instance.  This     */
/*     provides access to the signature keys for signing commands.            */
/*                                                                            */
/* Outputs:                                                                   */
/* map[string]bool -- set of the Subject Key Identifiers for the signature    */
/*     keys identified in the resource block.  maps SKI --> true.             */
/* map[string]string -- maps SKI --> signature key                            */
/* map[string]string -- maps SKI --> signature key token                      */
/* map[string]string -- maps SKI --> administrator name                       */
/* error -- reports any error during processing                               */
/*----------------------------------------------------------------------------*/
func GetSignatureKeysFromResourceBlock(hc HsmConfig) (map[string]bool,
	map[string]string, map[string]string, map[string]string, error) {

	// Set of Subject Key Identifiers
	suppliedSKIs := make(map[string]bool)
		// Use a map to check if a signature key is specified more than once
	// Maps SKIs to signature keys
	sigKeyMap := make(map[string]string)
	// Maps SKIs to signature key tokens
	sigKeyTokenMap := make(map[string]string)
	// Maps SKIs to administrator name
	adminNameMap := make(map[string]string)

	for i := 0; i < len(hc.Admins); i++ {
		ski, err := GetSigKeySKI(hc.Admins[i].Key, hc.Admins[i].Token)
		if err != nil {
			return suppliedSKIs, sigKeyMap, sigKeyTokenMap, adminNameMap, err
		}
		if suppliedSKIs[ski] {
			return suppliedSKIs, sigKeyMap, sigKeyTokenMap, adminNameMap,
				errors.New("A signature key has been specified more than once in the resource block")
		}
		suppliedSKIs[ski] = true
		sigKeyMap[ski] = hc.Admins[i].Key
		sigKeyTokenMap[ski] = hc.Admins[i].Token
		adminNameMap[ski] = hc.Admins[i].Name
	}
	return suppliedSKIs, sigKeyMap, sigKeyTokenMap, adminNameMap, nil
}

/*----------------------------------------------------------------------------*/
/* Returns the Subject Key Identifier (SKI) for a signature key.  Checks an   */
/* environment variable to determine whether a signing service should be used */
/* or whether the signature key is in a signature key file on the local       */
/* workstation.                                                               */
/*                                                                            */
/* Inputs:                                                                    */
/* sigkey string -- a string identifying which signature key to access        */
/* sigkeyToken string -- associated authentication token for the signature    */
/*     key                                                                    */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- Subject Key Identifier for the signature key, represented as a   */
/*     hexadecimal string.                                                    */
/* error -- reports any error during processing                               */
/*----------------------------------------------------------------------------*/
func GetSigKeySKI(sigkey string, sigkeyToken string) (string, error) {

	// Check if the environment variable is set indicating a signing service
	// should be used
	ssURL := os.Getenv("TKE_SIGNSERV_URL")
	if ssURL != "" {

		// Use the signing service to get the public key
		pubkey, err := common.GetPublicKeyFromSigningService(
			ssURL, sigkey, sigkeyToken)
		if err != nil {
			return "", err
		}
		hash := sha256.Sum256(pubkey[:])
		return strings.ToLower(hex.EncodeToString(hash[:])), nil

	} else {

		// When a signing service is not used, assume signature keys are in
		// files on the local workstation

		// Read the signature key file
		data, err := ioutil.ReadFile(sigkey)
		if err != nil {
			return "", err
		}
		var skfields map[string]string
		err = json.Unmarshal(data, &skfields)
		if err != nil {
			return "", err
		}
		ski := skfields["ski"]
		if len(ski) != 64 {
			return "", errors.New("Invalid Subject Key Identifier, length = " + strconv.Itoa(len(ski)))
		}

		return strings.ToLower(ski), nil
	}
}
