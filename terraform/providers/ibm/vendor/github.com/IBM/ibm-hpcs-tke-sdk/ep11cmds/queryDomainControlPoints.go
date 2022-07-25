//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/07/2021    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Queries the domain control points                                          */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain to be queried                         */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the domain control points (16 bytes long)                        */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDomainControlPoints(authToken string, urlStart string,
	de common.DomainEntry) ([]byte, error) {

	htpRequestString := QueryDomainControlPointsReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex())

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return nil, err
	}

	adminRspBlk, err := buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return nil, err
	}

	return adminRspBlk.CmdOutput, nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for querying domain control points                  */
/*----------------------------------------------------------------------------*/
func QueryDomainControlPointsReq(cryptoModuleIndex int, domainIndex int) string {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADMQ_DOM_CTRLPOINTS
	adminBlk.DomainID = BuildAdminDomainIndex(domainIndex)
	// module ID not used for queries
	// transaction counter not used for queries
	// no input parameters
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminBlk)
}
