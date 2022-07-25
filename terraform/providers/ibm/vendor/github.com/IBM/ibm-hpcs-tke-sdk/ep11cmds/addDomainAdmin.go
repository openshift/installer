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
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Adds a domain administrator                                                */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain where an administrator is to be added */
/* []byte -- certificate containing the public key for the administrator      */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func AddDomainAdmin(authToken string, urlStart string, de common.DomainEntry,
	cert []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) error {

	htpRequestString, err := AddDomainAdminReq(
		authToken, urlStart, de, cert, sigkeys, sigkeySkis, sigkeyTokens)
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
/* Creates the HTPRequest for adding a domain administrator                   */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose current wrapping key register   */
/*    is to be exported                                                       */
/* []byte -- certificate containing the public key for the administrator to   */
/*    be added                                                                */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the HTPRequest string with the signed CPRB for the command       */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func AddDomainAdminReq(authToken string, urlStart string, de common.DomainEntry,
	cert []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_DOM_ADMIN_LOGIN
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	// the certificate is the payload
	adminBlk.CmdInput = cert
	return CreateSignedHTPRequest(authToken, urlStart, de, adminBlk, sigkeys,
		sigkeySkis, sigkeyTokens)
}
