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
	"encoding/binary"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

type DomainAttributes struct {
	SignatureThreshold           uint32
	RevocationSignatureThreshold uint32
	Permissions                  uint32
	OperationalMode              uint32
	StandardsCompliance          uint32
}

/** Domain signature threshold */
const XCP_ADMINT_SIGN_THR = 1

/** Domain revocation signature threshold */
const XCP_ADMINT_REVOKE_THR = 2

/** Domain permissions */
const XCP_ADMINT_PERMITS = 3

/** Domain operational mode */
const XCP_ADMINT_MODE = 4

/** Domain security standards compliance */
const XCP_ADMINT_STD = 5

/*----------------------------------------------------------------------------*/
/* Queries the domain attributes                                              */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain to be queried                         */
/*                                                                            */
/* Outputs:                                                                   */
/* DomainAttributes -- structure with the domain attributes                   */
/* AdminRspBlk -- the xcpAdminRspBlk from the query.  The Query Domain        */
/*    Attributes command is issued before each signed command to determine    */
/*    the administrative domain,  module ID, and transaction counter fields   */
/*    that should be used for the signed command.                             */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDomainAttributes(authToken string, urlStart string,
	de common.DomainEntry) (DomainAttributes, AdminRspBlk, error) {

	var domainAttributes DomainAttributes
	var adminRspBlk AdminRspBlk

	htpRequestString := QueryDomainAttributesReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex())

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return domainAttributes, adminRspBlk, err
	}

	adminRspBlk, err = buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return domainAttributes, adminRspBlk, err
	}

	// Assemble the domain attributes in the output structure
	attributes := make(map[int]uint32)
	payload := adminRspBlk.CmdOutput

	for i := 0; i+8 <= len(payload); i += 8 {
		attributes[int(binary.BigEndian.Uint32(payload[i:i+4]))] =
			binary.BigEndian.Uint32(payload[i+4 : i+8])
	}

	signatureThreshold, ok := attributes[XCP_ADMINT_SIGN_THR]
	if ok {
		domainAttributes.SignatureThreshold = signatureThreshold
	} else {
		domainAttributes.SignatureThreshold = 999
		// an impossible value, checked elsewhere
	}

	revocationSignatureThreshold, ok := attributes[XCP_ADMINT_REVOKE_THR]
	if ok {
		domainAttributes.RevocationSignatureThreshold = revocationSignatureThreshold
	} else {
		domainAttributes.RevocationSignatureThreshold = 999
		// an impossible value, checked elsewhere
	}

	domainAttributes.Permissions = uint32(attributes[XCP_ADMINT_PERMITS])
	domainAttributes.OperationalMode = uint32(attributes[XCP_ADMINT_MODE])
	domainAttributes.StandardsCompliance = uint32(attributes[XCP_ADMINT_STD])
	// these set to 0 if not found

	return domainAttributes, adminRspBlk, nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for querying domain attributes                      */
/*----------------------------------------------------------------------------*/
func QueryDomainAttributesReq(cryptoModuleIndex int, domainIndex int) string {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADMQ_DOM_ATTRS
	adminBlk.DomainID = BuildAdminDomainIndex(domainIndex)
	// module ID not used for queries
	// transaction counter not used for queries
	// no input parameters
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminBlk)
}
