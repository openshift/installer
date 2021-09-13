//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/29/2021    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"encoding/hex"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Removes an administrator                                                   */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain with the administrator to be removed  */
/* string -- the Subject Key Identifier of the administator to be removed     */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func RemoveDomainAdministrator(authToken string, urlStart string,
	de common.DomainEntry, ski string, sigkeys []string,
	sigkeySkis []string, sigkeyTokens []string) error {

	// Convert from hexadecimal string to []byte
	skibytes, err := hex.DecodeString(ski)
	if err != nil {
		return err
	}

	htpRequestString, err := RemoveDomainAdminReq(authToken, urlStart, de,
		skibytes, sigkeys, sigkeySkis, sigkeyTokens)
	if err != nil {
		return err
	}

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return err
	}

	_, err = buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return err
	}

	return nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for removing a domain administrator                 */
/*----------------------------------------------------------------------------*/
func RemoveDomainAdminReq(authToken string, urlStart string,
	de common.DomainEntry, ski []byte, sigkeys []string,
	sigkeySkis []string, sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_DOM_ADMIN_LOGOUT
	// DomainID, ModuleID, and TransactionCounter get filled in later when sending the request
	adminBlk.CmdInput = ski
	return CreateSignedHTPRequest(authToken, urlStart, de, adminBlk, sigkeys,
		sigkeySkis, sigkeyTokens)
}
