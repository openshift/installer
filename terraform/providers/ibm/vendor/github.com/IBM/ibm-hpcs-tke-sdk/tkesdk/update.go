//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 06/21/2021    CLH             Initial version

package tkesdk

import (
	"errors"
	"os"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
	"github.com/IBM/ibm-hpcs-tke-sdk/ep11cmds"
)

/*----------------------------------------------------------------------------*/
/* Updates the crypto units in an HPCS service instance to match the desired  */
/* final configuration.                                                       */
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
func Update(ci CommonInputs, hc HsmConfig) ([]string, error) {

	// Check inputs in the resource block
	problems, err := checkInputs(hc)
	if err != nil {
		return make([]string, 0), err
	}
	if len(problems) > 0 {
		return problems, nil
	}

	// Read the initial configuration
	hsminfo, urlStart, domains, err := internalQuery(ci)

	// Check for invalid transitions
	problems, err, keepSKIs, addSKIs, rmvSKIs := internalCheckTransition(ci, hc, hsminfo)
	if err != nil {
		return make([]string, 0), err
	}
	if len(problems) > 0 {
		return problems, nil
	}

	// Do a pre-emptive zeroize to work around an undesired consequence of
	// an EP11 firmware update.
	anyAdminsRemoved := false
	for i := 0; i < len(hsminfo); i++ {
		// Only zeroize crypto units in imprint mode
		if hsminfo[i].SignatureThreshold == 0 {
			// Will the pre-emptive zeroize remove administrators?
			if len(hsminfo[i].Admins) > 0 {
				anyAdminsRemoved = true
			}
			// Do a pre-emptive zeroize
			sigkeys := make([]string, 0)
			sigkeySkis := make([]string, 0)
			sigkeyTokens := make([]string, 0)
			err := ep11cmds.ZeroizeDomain(ci.AuthToken, urlStart, domains[i],
				sigkeys, sigkeySkis, sigkeyTokens)
			if err != nil {
				return problems, err
			}
		}
	}
	// If administrators were removed by the pre-emptive zeroize, need to
	// refetch the initial configuration and redetermine what administrators
	// to keep, add, and remove.
	if anyAdminsRemoved {
		hsminfo, urlStart, domains, err = internalQuery(ci)
		if err != nil {
			return problems, err
		}
		problems, err, keepSKIs, addSKIs, rmvSKIs = internalCheckTransition(ci, hc, hsminfo)
		if err != nil {
			return make([]string, 0), err
		}
		if len(problems) > 0 {
			return problems, nil
		}
	}

	// Identify what signature keys are in the resource block
	suppliedSKIs, sigKeyMap, sigKeyTokenMap, adminNameMap, err :=
		GetSignatureKeysFromResourceBlock(hc)
	if err != nil {
		return make([]string, 0), err
	}

	// Create certificates for each signature key
	certMap := make(map[string][]byte, 0)
	// Maps SKI --> administrator certificate
	for ski := range suppliedSKIs {
		cert, err := createAdminCert(ski, sigKeyMap[ski],
			sigKeyTokenMap[ski], adminNameMap[ski])
		if err != nil {
			return make([]string, 0), err
		}
		certMap[ski] = cert
	}

	//--------------------------------------------------------------------------
	// Update administrators and signature thresholds
	//--------------------------------------------------------------------------

	for i, domain := range domains {

		if hsminfo[i].SignatureThreshold >= hsminfo[i].RevocationThreshold {

			//------------------------------------------------------------------
			// This case can be handled by removing administrators first,
			// then adding administrators, then changing the thresholds.
			//------------------------------------------------------------------

			// Assemble the set of signature keys to use to sign commands to
			// remove administrators
			sigkeys, sigkeySkis, sigkeyTokens :=
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hsminfo[i].RevocationThreshold)

			// Remove administrators
			for _, ski := range rmvSKIs[i] {
				err = ep11cmds.RemoveDomainAdministrator(ci.AuthToken,
					urlStart, domain, ski, sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}
			}

			// Assemble the set of signature keys to use to sign commands to
			// add administrators
			sigkeys, sigkeySkis, sigkeyTokens =
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hsminfo[i].SignatureThreshold)

			// Add administrators
			for _, ski := range addSKIs[i] {
				err = ep11cmds.AddDomainAdmin(ci.AuthToken, urlStart, domain,
					certMap[ski], sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}
				// Make this administrator available to sign subsequent commands
				keepSKIs[i] = append(keepSKIs[i], ski)
			}

			// Assemble the set of signature keys to use to sign the command
			// to change the signature thresholds
			if hsminfo[i].SignatureThreshold == 0 {
				// Leaving imprint mode is a special case.
				// The number of required signatures is the new signature
				// threshold value.
				sigkeys, sigkeySkis, sigkeyTokens =
					collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
						hc.SignatureThreshold)
			} else {
				// Not leaving imprint mode
				// Can use the set of signature keys already assembled
			}

			// Change the signature thresholds and other domain attributes
			err = SetDomainAttributes(ci.AuthToken, urlStart, domain,
				hc.SignatureThreshold, hc.RevocationThreshold,
				sigkeys, sigkeySkis, sigkeyTokens)
			if err != nil {
				return make([]string, 0), err
			}

		} else if hc.RevocationThreshold <= len(keepSKIs[i]) {

			//------------------------------------------------------------------
			// This case can be handled by changing the revocation threshold
			// first, then removing administrators, then adding administrators,
			// then changing the signature threshold.
			//------------------------------------------------------------------

			// Assemble the set of signature keys to use to sign the command to
			// change the revocation threshold
			sigkeys, sigkeySkis, sigkeyTokens :=
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hsminfo[i].SignatureThreshold)

			// Keep current signature threshold but change the revocation
			// threshold
			err = SetDomainAttributes(ci.AuthToken, urlStart, domain,
				hsminfo[i].SignatureThreshold, hc.RevocationThreshold,
				sigkeys, sigkeySkis, sigkeyTokens)
			if err != nil {
				return make([]string, 0), err
			}

			// Assemble the set of signature keys to use to sign commands to
			// remove administrators
			sigkeys, sigkeySkis, sigkeyTokens =
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hc.RevocationThreshold)

			// Remove administrators
			for _, ski := range rmvSKIs[i] {
				err = ep11cmds.RemoveDomainAdministrator(ci.AuthToken,
					urlStart, domain, ski, sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}
			}

			// Assemble the set of signature keys to use to sign commands to
			// add administrators
			sigkeys, sigkeySkis, sigkeyTokens =
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hsminfo[i].SignatureThreshold)

			// Add administrators
			for _, ski := range addSKIs[i] {
				err = ep11cmds.AddDomainAdmin(ci.AuthToken, urlStart, domain,
					certMap[ski], sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}
			}

			// Change the signature threshold
			// Can use the same signature keys as the previous operation
			err = SetDomainAttributes(ci.AuthToken, urlStart, domain,
				hc.SignatureThreshold, hc.RevocationThreshold,
				sigkeys, sigkeySkis, sigkeyTokens)
			if err != nil {
				return make([]string, 0), err
			}

		} else if len(keepSKIs[i])+len(addSKIs[i])+len(rmvSKIs[i]) <= 8 {

			//------------------------------------------------------------------
			// This case can be handled by adding administrators first, then
			// removing administrators, then changing the thresholds.
			//------------------------------------------------------------------

			// Assemble the set of signature keys to use to sign commands to
			// add administrators
			sigkeys, sigkeySkis, sigkeyTokens :=
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hsminfo[i].SignatureThreshold)

			// Add administrators
			for _, ski := range addSKIs[i] {
				err = ep11cmds.AddDomainAdmin(ci.AuthToken, urlStart, domain,
					certMap[ski], sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}
				// Make this administrator available to sign subsequent commands
				keepSKIs[i] = append(keepSKIs[i], ski)
			}

			// Assemble the set of signature keys to use to sign commands to
			// remove administrators
			sigkeys, sigkeySkis, sigkeyTokens =
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hsminfo[i].RevocationThreshold)

			// Remove administrators
			for _, ski := range rmvSKIs[i] {
				err = ep11cmds.RemoveDomainAdministrator(ci.AuthToken,
					urlStart, domain, ski, sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}
			}

			// Assemble the set of signature keys to use to sign the command
			// to change the signature thresholds
			sigkeys, sigkeySkis, sigkeyTokens =
				collectSigKeys(keepSKIs[i], sigKeyMap, sigKeyTokenMap,
					hsminfo[i].SignatureThreshold)

			// Change the signature thresholds
			err = SetDomainAttributes(ci.AuthToken, urlStart, domain,
				hc.SignatureThreshold, hc.RevocationThreshold,
				sigkeys, sigkeySkis, sigkeyTokens)
			if err != nil {
				return make([]string, 0), err
			}

		} else {
			// Previous checks should prevent us from ever getting here
			return make([]string, 0), errors.New("Unsupported state transition")
		}
	}

	//--------------------------------------------------------------------------
	// Update the current master key registers
	//--------------------------------------------------------------------------

	// Two cases are handled:
	// 1. All current master key registers are initially empty.
	// 2. The current master key register in at least one recovery crypto unit
	//    is set, and all other crypto units either have the same master key
	//    value or they are empty.
	//
	// We do not attempt to handle any other initial condition.  The call to
	// internalCheckTransition only allows these two cases, and the code below
	// relies on that check being done.

	// Don't have something we need in quite the right form
	availableSKIs := make([]string, 0)
	for ski := range suppliedSKIs {
		availableSKIs = append(availableSKIs, ski)
	}

	// Only need one signature for some commands
	singleSigkey, singleSigkeySki, singleSigkeyToken :=
		collectSigKeys(availableSKIs, sigKeyMap, sigKeyTokenMap, 1)

	// Other commands require the signature threshold number of signatures
	sigkeys, sigkeySkis, sigkeyTokens :=
		collectSigKeys(availableSKIs, sigKeyMap, sigKeyTokenMap, hc.SignatureThreshold)

	// Check if all master key registers are initially empty
	allEmpty := true
	for i := 0; i < len(hsminfo); i++ {
		if hsminfo[i].CurrentMKStatus != "Empty" {
			allEmpty = false
			break
		}
	}

	var recoveryHSM common.DomainEntry
	var recoveryHSMindex int

	if allEmpty {
		// Look for a recovery crypto unit
		foundIt := false
		for i, domain := range domains {
			if domain.Type == "recovery" {
				recoveryHSM = domain
				recoveryHSMindex = i
				foundIt = true
				break
			}
		}
		if !foundIt {
			return make([]string, 0), errors.New("No recovery crypto unit found when setting master key registers")
		}

		// Create a random WK in the recovery crypto unit
		err, _ := ep11cmds.CreateRandomWK(ci.AuthToken, urlStart, recoveryHSM,
			singleSigkey, singleSigkeySki, singleSigkeyToken)
		if err != nil {
			return make([]string, 0), err
		}

	} else {
		// Look for a recovery crypto unit whose current master key register
		// is set
		foundIt := false
		for i, domain := range domains {
			if domain.Type == "recovery" &&
				hsminfo[i].CurrentMKStatus != "Empty" {

				recoveryHSM = domain
				recoveryHSMindex = i
				foundIt = true
				break
			}
		}
		if !foundIt {
			return make([]string, 0), errors.New("No recovery crypto unit found whose current master key register is set")
		}
	}

	// Transfer the master key value to the other crypto units
	for i, domain := range domains {
		if i != recoveryHSMindex {
			if hsminfo[i].CurrentMKStatus == "Empty" {

				// Generate an importer key in the target domain
				pubKey, _, err := ep11cmds.GenerateP521ECImporterKey(
					ci.AuthToken, urlStart, domain, singleSigkey,
					singleSigkeySki, singleSigkeyToken)
				if err != nil {
					return make([]string, 0), err
				}

				// Export the master key value from the recovery crypto unit using
				// the importer key
				kphcert := ep11cmds.KPHCert(pubKey)
				pfile := ep11cmds.ExportWKParameterFile(kphcert)
				pdata, err := ep11cmds.ExportWK(ci.AuthToken, urlStart,
					recoveryHSM, pfile, sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}

				var pMap common.ParameterMap
				pMap, err = pMap.Load(pdata)
				if err != nil {
					return make([]string, 0), err
				}

				// Get the recipient info for the single key part
				recipientInfo := make([][]byte, 0)
				recipientInfo = append(recipientInfo, pMap.GetDataUsingIndex(common.PMTAG_ENCR_KEY_PART, 0))

				// Import the master key to the target domain
				err = ep11cmds.ImportWK(ci.AuthToken, urlStart, domain, recipientInfo,
					singleSigkey, singleSigkeySki, singleSigkeyToken)
				if err != nil {
					return make([]string, 0), err
				}

				// Commit the imported master key
				err = ep11cmds.CommitPendingWK(ci.AuthToken, urlStart, domain,
					sigkeys, sigkeySkis, sigkeyTokens)
				if err != nil {
					return make([]string, 0), err
				}

				// Finalize the imported master key
				err = ep11cmds.FinalizeWK(ci.AuthToken, urlStart, domain,
					singleSigkey, singleSigkeySki, singleSigkeyToken)
				if err != nil {
					return make([]string, 0), err
				}
			}
		}
	}

	return make([]string, 0), nil
}

/*----------------------------------------------------------------------------*/
/* Assembles a set of signature keys that can be used to sign a command.      */
/*                                                                            */
/* Inputs:                                                                    */
/* []string -- set of the Subject Key Identifiers for signature keys that     */
/*     can be used to sign the command.  These SKIs must be for signature     */
/*     keys that are specified in the resource block and that are already     */
/*     installed as administrators on the target crypto unit.                 */
/* map[string]string -- maps SKI --> signature key                            */
/* map[string]string -- maps SKI --> signature key token                      */
/* int -- number of signatures needed                                         */
/*                                                                            */
/* Outputs:                                                                   */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*----------------------------------------------------------------------------*/
func collectSigKeys(allowedSKIs []string, sigKeyMap map[string]string,
	sigKeyTokenMap map[string]string, needed int) ([]string, []string,
	[]string) {

	sigkeys := make([]string, 0)
	sigkeySkis := make([]string, 0)
	sigkeyTokens := make([]string, 0)
	for _, ski := range allowedSKIs {
		if len(sigkeys) < needed {
			sigkeys = append(sigkeys, sigKeyMap[ski])
			sigkeySkis = append(sigkeySkis, ski)
			sigkeyTokens = append(sigkeyTokens, sigKeyTokenMap[ski])
		} else {
			break
		}
	}
	if len(sigkeys) < needed {
		panic("Internal error: not enough administrators to meet threshold value")
	}
	return sigkeys, sigkeySkis, sigkeyTokens
}

/*----------------------------------------------------------------------------*/
/* Creates an administrator certificate                                       */
/*                                                                            */
/* Inputs:                                                                    */
/* string -- Subject Key Identifier for the administrator to be created,      */
/*     represented as a hexadecimal string                                    */
/* string -- Key parameter for the signature key to be used                   */
/* string -- Token parameter for the signature key to be used                 */
/* string -- the name to be placed in the certificate                         */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the administrator certificate                                    */
/* error -- identifies any error encountered during processing                */
/*----------------------------------------------------------------------------*/
func createAdminCert(ski string, sigkey string, sigkeyToken string,
	adminName string) ([]byte, error) {

	ssURL := os.Getenv("TKE_SIGNSERV_URL")
	if ssURL != "" {
		return ep11cmds.CreateAdminCertUsingSigningService(ssURL, sigkey, sigkeyToken, adminName)
	} else {
		return CreateAdminCertFromFile(sigkey, ski, sigkeyToken, adminName)
	}
}

/*----------------------------------------------------------------------------*/
/* Sets the domain attributes.  Different attributes are set for recovery     */
/* HSMs and operational HSMs.                                                 */
/*                                                                            */
/* Inputs:                                                                    */
/* PluginContext -- contains the IAM access token and parameters identifying  */
/*    what resource group the user is working with                            */
/* DomainEntry -- identifies the domain whose attributes are to be set        */
/* int -- new signature threshold value to set                                */
/* int -- new revocation signature threshold value to set                     */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Output:                                                                    */
/* error -- reports any errors accessing the domain                           */
/*----------------------------------------------------------------------------*/
func SetDomainAttributes(authToken string, urlStart string,
	domain common.DomainEntry, newSigThr int, newRevThr int,
	sigkeys []string, sigkeySkis []string, sigkeyTokens []string) error {

	// Get the current domain attributes
	domainAttributes, _, err := ep11cmds.QueryDomainAttributes(
		authToken, urlStart, domain)
	if err != nil {
		return err
	}

	// Allow domains to be zeroized with a single signature
	domainAttributes.Permissions |= 0x00000040
		// May remove this in the future.  Bank of America was told the
		// signature threshold value applies.  Should change TKE plug-in
		// and HPCS Management Utilities at the same time.

	// Allow master key import
	domainAttributes.Permissions |= 0x00000001

	// Allow importing using a single key part
	domainAttributes.Permissions |= 0x00000004

	// Set some attributes differently for recovery HSMs and operational HSMs
	if domain.Type == "recovery" {
		// Recovery HSM

		// Allow randomly generated wrapping key
		domainAttributes.Permissions |= 0x00000008

		// Allow master key export
		domainAttributes.Permissions |= 0x00000002

		// Set "do not disturb"
		domainAttributes.Permissions |= 0x00002000

		// Prevent "do not disturb" from being changed
		domainAttributes.Permissions &= 0x7FFFFFFF

	} else {
		// Operational HSM

		// Disable master key export
		domainAttributes.Permissions &= 0xFFFFFFFD

		// Conditionally set "do not disturb" and prevent it from being
		// changed
		if (domainAttributes.Permissions & 0x80000000) == 0x80000000 {

			// Set "do not disturb"
			domainAttributes.Permissions |= 0x00002000

			// Prevent "do not disturb" from being changed
			domainAttributes.Permissions &= 0x7FFFFFFF

			// We would like to always turn on "do not disturb" and
			// prevent future changes, but that may not be possible.
			// If a domain was initialized before the EP11 firmware update
			// that implements "do not disturb" was applied, both
			// "do not disturb" and the corresponding "may change" control
			// are off.  For that case we cannot change the attribute.
			// Recovery HSMs were added to customer configurations after the
			// EP11 firmware update that implements "do not disturb" was
			// applied.  So the conditional handling of "do not disturb"
			// needs to be added only for operational HSMs.
		}
	}

	// Set the new signature thresholds
	domainAttributes.SignatureThreshold = uint32(newSigThr)
	domainAttributes.RevocationSignatureThreshold = uint32(newRevThr)

	err = ep11cmds.SetDomainAttributes(
		authToken, urlStart, domain, domainAttributes, sigkeys, sigkeySkis,
		sigkeyTokens)
	if err != nil {
		return err
	}

	return nil
}
