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
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Commits the pending wrapping key register                                  */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose pending wrapping key register   */
/*    is to be committed                                                      */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func CommitPendingWK(authToken string, urlStart string, de common.DomainEntry,
	sigkeys []string, sigkeySkis []string, sigkeyTokens []string) error {

	// Get the verification pattern for the pending wrapping key register
	domainInfo, err := QueryDomainInfo(authToken, urlStart, de)
	if err != nil {
		return err
	}

	htpRequestString, err := CommitPendingWKReq(authToken, urlStart, de,
		domainInfo.NewMKVP, sigkeys, sigkeySkis, sigkeyTokens)
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
/* Creates the HTPRequest for committing the pending wrapping key register    */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose pending wrapping key register   */
/*    is to be committed                                                      */
/* []byte -- the verification pattern of the pending wrapping key register    */
/*    to be committed                                                         */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the HTPRequest string with the signed CPRB for the command       */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func CommitPendingWKReq(authToken string, urlStart string,
	de common.DomainEntry, vp []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_COMMIT_WK
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	adminBlk.CmdInput = vp
	return CreateSignedHTPRequest(authToken, urlStart, de, adminBlk, sigkeys,
		sigkeySkis, sigkeyTokens)
}
