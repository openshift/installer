//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/09/2021    CLH             Initial version

package tkesdk

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/IBM/ibm-hpcs-tke-sdk/ep11cmds"
)

// Permission bit to zeroize with one signature
const XCP_ADMP_ZERO_1SIGN = 0x00000040

/*----------------------------------------------------------------------------*/
/* Zeroizes the crypto units assigned to a service instance, or returns an    */
/* error if that is not possible.                                             */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, and an IBM Cloud authentication token.           */
/* HsmConfig -- A structure containing information from the hsm_config        */
/*      section of the resource block for the HPCS service instance.  This    */
/*      provides access to signature keys for signing commands to crypto      */
/*      units.                                                                */
/*----------------------------------------------------------------------------*/
func Zeroize(ci CommonInputs, hc HsmConfig) error {

	// Query the initial configuration of the crypto units
	hsminfo, urlStart, domains, err := internalQuery(ci)
	if err != nil {
		return err
	}

	// Check if any crypto unit has left imprint mode
	imprintModeOnly := true
	for i := 0; i < len(hsminfo); i++ {
		if hsminfo[i].SignatureThreshold > 0 {
			imprintModeOnly = false
			break
		}
	}

	// Handle special case of all crypto units in imprint mode
	if imprintModeOnly {
		sigkeys := make([]string, 0)
		sigkeySkis := make([]string, 0)
		sigkeyTokens := make([]string, 0)
		for i := 0; i < len(hsminfo); i++ {
			err := ep11cmds.ZeroizeDomain(ci.AuthToken, urlStart, domains[i],
				sigkeys, sigkeySkis, sigkeyTokens)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Check that all signature keys specified in the resource block can be
	// accessed
	for _, adminInfo := range hc.Admins {
		if !validKey(adminInfo) {
			return errors.New("One or more signature keys cannot be accessed.")
		}
	}

	// Determine what signature keys are available
	suppliedSkis, sigKeyMap, sigKeyTokenMap, _, err :=
		GetSignatureKeysFromResourceBlock(hc)
	if err != nil {
		return err
	}

	// Determine the zeroize with one signature attribute for all crypto units
	zeroizeWithOne := make([]bool, 0)
	for i := 0; i < len(domains); i++ {
		attr, _, err := ep11cmds.QueryDomainAttributes(ci.AuthToken, urlStart, domains[i])
		if err != nil {
			return err
		}
		zeroizeWithOne = append(zeroizeWithOne, ((attr.Permissions & XCP_ADMP_ZERO_1SIGN) == XCP_ADMP_ZERO_1SIGN))
	}

	// Read the installed administrators for all crypto units
	installedAdminSkis := make([][]string, 0)
	for i := 0; i < len(domains); i++ {
		skiBytes, err := ep11cmds.QueryDomainAdmins(ci.AuthToken, urlStart, domains[i])
		if err != nil {
			return err
		}
		iasThisHsm := make([]string, 0)
		for j := 0; j < len(skiBytes); j++ {
			iasThisHsm = append(iasThisHsm, strings.ToLower(hex.EncodeToString(skiBytes[j])))
		}
		installedAdminSkis = append(installedAdminSkis, iasThisHsm)
	}

	// Check whether the supplied administrator signature keys are sufficient
	failingDomainLocations := make([]string, 0)
	for i := 0; i < len(domains); i++ {
		// Count the number of valid supplied signature keys for this crypto unit
		count := 0
		for j := 0; j < len(installedAdminSkis[i]); j++ {
			if suppliedSkis[installedAdminSkis[i][j]] {
				count++
			}
		}
		if (count < hsminfo[i].SignatureThreshold) && !zeroizeWithOne[i] {
			failingDomainLocations = append(failingDomainLocations, hsminfo[i].HsmLocation)
		} else if zeroizeWithOne[i] && (hsminfo[i].SignatureThreshold > 0) && (count == 0) {
			failingDomainLocations = append(failingDomainLocations, hsminfo[i].HsmLocation)
		}
	}
	if len(failingDomainLocations) > 0 {
		errorMsg := "The administrator signature keys specified in the " +
			"resource block do not allow the operation to be performed " +
			"in the following crypto units: "
		for i := 0; i < len(failingDomainLocations); i++ {
			if i < len(failingDomainLocations)-1 {
				errorMsg = errorMsg + failingDomainLocations[i] + ", "
			} else {
				errorMsg = errorMsg + failingDomainLocations[i]
			}
		}
		return errors.New(errorMsg)
	}

	// Zeroize all crypto units
	for i := 0; i < len(domains); i++ {

		sigkeys := make([]string, 0)
		sigkeySkis := make([]string, 0)
		sigkeyTokens := make([]string, 0)

		var signaturesNeeded int
		if hsminfo[i].SignatureThreshold == 0 {
			signaturesNeeded = 0
		} else if zeroizeWithOne[i] {
			signaturesNeeded = 1
		} else {
			signaturesNeeded = hsminfo[i].SignatureThreshold
		}

		// Select the signature keys to be used for this crypto unit
		for j := 0; j < len(installedAdminSkis[i]); j++ {
			if len(sigkeys) == signaturesNeeded {
				break
			}
			if suppliedSkis[installedAdminSkis[i][j]] {
				sigkeys = append(sigkeys, sigKeyMap[installedAdminSkis[i][j]])
				sigkeySkis = append(sigkeySkis, installedAdminSkis[i][j])
				sigkeyTokens = append(sigkeyTokens, sigKeyTokenMap[installedAdminSkis[i][j]])
			}
		}
		if len(sigkeys) != signaturesNeeded {
			return errors.New("Error selecting signature keys to sign zeroize command")
			// Previous checks should prevent this from ever being reported
		}

		// Zeroize the crypto unit
		err := ep11cmds.ZeroizeDomain(ci.AuthToken, urlStart, domains[i],
			sigkeys, sigkeySkis, sigkeyTokens)
		if err != nil {
			return err
		}
	}

	return nil
}
