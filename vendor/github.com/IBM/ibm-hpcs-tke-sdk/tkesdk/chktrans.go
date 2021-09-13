//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/03/2021    CLH             Initial version
// 07/23/2021    CLH             Change message when a signature key cannot be used

package tkesdk

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Checks for invalid inputs and checks that the transition from initial      */
/* state to final state is possible.                                          */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, and an IBM Cloud authentication token.           */
/* HsmConfig -- A structure containing information from the hsm_config        */
/*      section of the resource block for the HPCS service instance.  This    */
/*      provides access to signature keys for signing commands to crypto      */
/*      units.                                                                */
/*                                                                            */
/* Outputs:                                                                   */
/* []string -- set of messages identifying either an invalid input or a       */
/*      reason the transition from initial state to desired final state is    */
/*      not possible                                                          */
/* error -- identifies any error encountered when running the function        */
/*----------------------------------------------------------------------------*/
func CheckTransition(ci CommonInputs, hc HsmConfig) ([]string, error) {

	// Check inputs in the resource block
	problems, err := checkInputs(hc)
	if err != nil {
		return make([]string, 0), err
	}
	if len(problems) > 0 {
		return problems, nil
	}

	// Read the initial configuration
	hsminfo, _, _, err := internalQuery(ci)
	if err != nil {
		return make([]string, 0), err
	}

	// Check for invalid transitions
	problems, err, _, _, _ = internalCheckTransition(ci, hc, hsminfo)
	if err != nil {
		return make([]string, 0), err
	}

	return problems, nil
}

/*----------------------------------------------------------------------------*/
/* Check for problems with the inputs specified by the user.                  */
/*----------------------------------------------------------------------------*/
func checkInputs(hc HsmConfig) ([]string, error) {

	problems := make([]string, 0)
	if hc.SignatureThreshold < 1 || hc.SignatureThreshold > 8 {
		problems = append(problems, "The signature threshold must be an integer between 1 and 8.")
	}
	if hc.RevocationThreshold < 1 || hc.RevocationThreshold > 8 {
		problems = append(problems, "The revocation threshold must be an integer between 1 and 8.")
	}
	if len(hc.Admins) < hc.SignatureThreshold {
		problems = append(problems, "Not enough administrators are specified to meet the signature threshold value.")
	}
	if len(hc.Admins) < hc.RevocationThreshold {
		problems = append(problems, "Not enough administrators are specified to meet the revocation threshold value.")
	}
	if len(hc.Admins) > 8 {
		problems = append(problems, "No more than 8 administrators can be specified.")
	}

	allKeysValid := true
	for _, admin := range hc.Admins {
		if len(admin.Name) > 30 {
			problems = append(problems, "An administrator name is too long.  Names must be 30 characters or less.")
		}
		if !validKey(admin) {
			ssURL := os.Getenv("TKE_SIGNSERV_URL")
			if ssURL != "" {
				problems = append(problems, "The signature key associated with " +
					admin.Name + " could not be accessed.  An attempt was made " +
					"to use a signing service.  The signing service may not be " +
					"running at the specified URL and port.")
			} else {
				problems = append(problems, "The signature key associated with " +
					admin.Name + " could not be accessed.")
			}
			allKeysValid = false
		}
	}

	if allKeysValid {
		uniqueKeys, err := keysAreUnique(hc.Admins)
		if err != nil {
			return problems, err
		}
		if !uniqueKeys {
			problems = append(problems, "Signature keys are not unique.  The same signature key is specified for more than one administrator.")
		}
	}

	return problems, nil
}

/*----------------------------------------------------------------------------*/
/* Checks whether the transition is allowed.                                  */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, and an IBM Cloud authentication token.           */
/* HsmConfig -- A structure containing information from the hsm_config        */
/*      section of the resource block for the HPCS service instance.  This    */
/*      provides access to signature keys for signing commands to crypto      */
/*      units.                                                                */
/* []HsmInfo -- Contains information on the initial configuration of each     */
/*      crypto unit.                                                          */
/*                                                                            */
/* Outputs:                                                                   */
/* []string -- set of messages identifying either an invalid input or a       */
/*      reason the transition from initial state to desired final state is    */
/*      not possible                                                          */
/* error -- identifies any error encountered when running the function        */
/*                                                                            */
/* The next three outputs are available for possible use later.  The first    */
/* index references different crypto units.                                   */
/* [][]string -- the Subject Key Identifiers for existing administrators to   */
/*      remain on the crypto unit                                             */
/* [][]string -- the Subject Key Identifiers of new administrators to be      */
/*      added to the crypto unit                                              */
/* [][]string -- the Subject Key Identifiers of the existing administrators   */
/*      to be removed from the crypto unit                                    */
/*----------------------------------------------------------------------------*/
func internalCheckTransition(ci CommonInputs, hc HsmConfig,
	hsminfo []HsmInfo) ([]string, error, [][]string, [][]string, [][]string) {

	// Initialize the output variables
	problems := make([]string, 0)
	allKeepSKIs := make([][]string, 0)
	allAddSKIs := make([][]string, 0)
	allRmvSKIs := make([][]string, 0)

	// The service instance must have at least one recovery crypto unit
	foundRecovery := false
	for i := 0; i < len(hsminfo); i++ {
		if hsminfo[i].HsmType == "recovery" {
			foundRecovery = true
			break
		}
	}
	if !foundRecovery {
		problems = append(problems, "The service instance does not contain any recovery crypto units.")
	}

	// Determine the desired final set of administrator SKIs for all crypto units
	finalSKIs, _, _, adminNameMap, err := GetSignatureKeysFromResourceBlock(hc)
	if err != nil {
		return problems, err, allKeepSKIs, allAddSKIs, allRmvSKIs
	}

	// For each crypto unit, figure out what administrators we want to keep,
	// what administrators we want to add, and what administrators we want
	// to remove.  Determine whether the changes are possible.
	for i := 0; i < len(hsminfo); i++ {

		keepSKIs := make([]string, 0)
		addSKIs  := make([]string, 0)
		rmvSKIs  := make([]string, 0)

		// What administrators do we want to keep or remove?
		for j := 0; j < len(hsminfo[i].Admins); j++ {
			if finalSKIs[hsminfo[i].Admins[j].AdminSKI] {
				keepSKIs = append(keepSKIs, hsminfo[i].Admins[j].AdminSKI)
			} else {
				rmvSKIs = append(rmvSKIs, hsminfo[i].Admins[j].AdminSKI)
			}
		}

		// What administrators do we want to add?
		for ski := range finalSKIs {
			foundIt := false
			for _, admin := range hsminfo[i].Admins {
				if admin.AdminSKI == ski {
					foundIt = true
					break
				}
			}
			if !foundIt {
				addSKIs = append(addSKIs, ski)
			}
		}

		allKeepSKIs = append(allKeepSKIs, keepSKIs)
		allAddSKIs  = append(allAddSKIs, addSKIs)
		allRmvSKIs  = append(allRmvSKIs, rmvSKIs)

		// Check whether the changes are possible
		if len(keepSKIs) < hsminfo[i].SignatureThreshold {
			problems = append(problems, "Not enough signature keys for "+
				"installed administrators are provided in the resource "+
				"block to meet the current signature threshold.")
		} else if hsminfo[i].SignatureThreshold >= hsminfo[i].RevocationThreshold {
			// This case can be handled by removing administrators first, then
			// adding administrators, then changing the signature thresholds.
			// The first check ensures we can do this.
		} else if hc.RevocationThreshold <= len(keepSKIs) {
			// This can can be handled by changing the revocation threshold
			// first, then removing administrators, then adding administrators,
			// then changing the signature threshold.
		} else if len(keepSKIs)+len(addSKIs)+len(rmvSKIs) <= 8 {
			// This can be handled by adding administrators, then removing
			// administrators, then changing the signature thresholds.
			// The first check ensures we can add up to
			// 8 - (len(keepSKIs) + len(rmvSKIs)) administrators.
		} else {
			problems = append(problems, "Not enough signature keys for "+
				"installed administrators are provided in the resource "+
				"block to meet the current revocation threshold.")
		}
	}

	// Check whether new names are specified for any administrators
	// being kept
	for i := 0; i < len(hsminfo); i++ {
		// Only need to check administrator names for crypto units not in
		// imprint mode.  For Update(), a pre-emptive zeroize will be done
		// for crypto units still in imprint mode to work around undesired
		// behavior from an EP11 firmware update.
		if hsminfo[i].SignatureThreshold > 0 {
			for _, ski := range allKeepSKIs[i] {
				foundIt := false
				for j := 0; j < len(hsminfo[i].Admins); j++ {
					if hsminfo[i].Admins[j].AdminSKI == ski {
						foundIt = true
						if strings.TrimSpace(adminNameMap[ski]) !=
							strings.TrimSpace(hsminfo[i].Admins[j].AdminName) {
							problems = append(problems, "You are not allowed to change the name of an existing adminstrator.")
						}
						break
					}
				}
				if !foundIt {
					return problems, errors.New("Error checking for change in existing administrator name"),
						allKeepSKIs, allAddSKIs, allRmvSKIs
				}
			}
		}
	}

	// Check the current master key registers
	// Two cases can be handled:
	// 1) All current master key registers are empty.
	// 2) The current master key register in at least one recovery crypto
	//    unit is not empty and all other current master key registers are
	//    either empty or have the same verification pattern as the recovery
	//    crypto unit.

	allEmpty := true
	for i := 0; i < len(hsminfo); i++ {
		if hsminfo[i].CurrentMKStatus != "Empty" {
			allEmpty = false
			break
		}
	}

	if !allEmpty {
		vp := ""
		for i := 0; i < len(hsminfo); i++ {
			if (hsminfo[i].HsmType == "recovery") &&
				(hsminfo[i].CurrentMKStatus != "Empty") {

				vp = hsminfo[i].CurrentMKVP
				break
			}
		}
		if vp == "" {
			problems = append(problems, "The current master key register "+
				"is set in one or more operational crypto units but is not "+
				"set in any recovery crypto units.")
		} else {
			for i := 0; i < len(hsminfo); i++ {
				if (hsminfo[i].CurrentMKStatus != "Empty") &&
					(hsminfo[i].CurrentMKVP != vp) {

					problems = append(problems, "Current master key "+
						"registers are set in multiple crypto units but "+
						"are not set to the same value.")
					break
				}
			}
		}
	}

	return problems, nil, allKeepSKIs, allAddSKIs, allRmvSKIs
}

/*----------------------------------------------------------------------------*/
/* Checks whether a signature key can be used.                                */
/*----------------------------------------------------------------------------*/
func validKey(ai AdminInfo) bool {

	// Tries to sign some data.  If successful, the signature key can be used.

	// SignWithSignatureKey handles both signature keys stored in files and
	// signature keys accessed by a signing service.

	dataToSign := make([]byte, 100)
	_, err := common.SignWithSignatureKey(dataToSign, ai.Key, ai.Token)
	return err == nil
}

/*----------------------------------------------------------------------------*/
/* Checks that a unique key is specified for each administrator.              */
/*                                                                            */
/* Input:                                                                     */
/* []AdminInfo -- administrator signature key information from the Terraform  */
/*     resource block                                                         */
/*                                                                            */
/* Outputs:                                                                   */
/* bool -- true if unique signature keys are specified, false if a signature  */
/*     key is specified more than once                                        */
/* error -- reports any error found during processing                         */
/*----------------------------------------------------------------------------*/
func keysAreUnique(admins []AdminInfo) (bool, error) {
	skis := make(map[string]bool)
	for _, admin := range admins {
		ski, err := GetSigKeySKI(admin.Key, admin.Token)
		if err != nil {
			return false, err
		}
		skis[ski] = true
	}
	return len(skis) == len(admins), nil
}
