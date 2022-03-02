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
	"github.com/Logicalis/asn1"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Loads the new wrapping key register.                                       */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose new wrapping key register is    */
/*    to be loaded.                                                           */
/* [][]byte -- array of recipient info, one entry per key part                */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*    Only one signature is needed for this command.                          */
/*                                                                            */
/* Output:                                                                    */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func ImportWK(authToken string, urlStart string, de common.DomainEntry,
	recipientInfo [][]byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) error {

	// Create a concatenated set of signed xcpAdminReq, one for each
	// key part.
	var adminRequests []byte

	// Issue Query Domain Attributes to get the administrative domain,
	// the module identifier, the transaction counter, and the
	// signature thresholds
	_, adminRspBlk, err := QueryDomainAttributes(authToken, urlStart, de)
	if err != nil {
		return err
	}

	// Increment the transaction counter value
	transactionCounter :=
		IncrementTransactionCounter(adminRspBlk.TransactionCounter)

	// Repeat for each key part
	for i := range recipientInfo {

		// Build the individual import wrapping key admin block
		adminBlk := ImportWKRequest(adminRspBlk.DomainID,
			adminRspBlk.ModuleID, transactionCounter, recipientInfo[i])
		if err != nil {
			return err
		}
		adminBlockSeq, err := asn1.Encode(adminBlk)
		if err != nil {
			return err
		}

		// Sign the admin block
		signerInfo, err := CreateSignerInfo(adminBlockSeq, sigkeys,
			sigkeySkis, sigkeyTokens)
		if err != nil {
			return err
		}

		// Create the xcpAdminReq sequence
		var adminReq AdminReq
		adminReq.CmdID = FNID_ADMIN
		adminReq.DomainID = adminRspBlk.DomainID[0:4]
		adminReq.AdminBlock = adminBlockSeq
		adminReq.SignerInfo = signerInfo

		adminReqSeq, err := asn1.Encode(adminReq)
		if err != nil {
			return err
		}
		adminRequests = append(adminRequests, adminReqSeq...)
	}

	// Build the enveloping xcpAdminReq
	var bigAdminBlk AdminBlk
	bigAdminBlk = BigImportWKRequest(adminRspBlk.DomainID,
		adminRspBlk.ModuleID, transactionCounter, adminRequests)

	bigAdminBlockSeq, err := asn1.Encode(bigAdminBlk)
	if err != nil {
		return err
	}
	var bigAdminReq AdminReq
	bigAdminReq.CmdID = FNID_ADMIN
	bigAdminReq.DomainID = adminRspBlk.DomainID[0:4]
	bigAdminReq.AdminBlock = bigAdminBlockSeq

	bigAdminReqSeq, err := asn1.Encode(bigAdminReq)
	if err != nil {
		return err
	}

	// Form and submit the final HTPRequest string
	xpNumRequest := NewXPNUMRequest(de.GetCryptoModuleIndex(),
		de.GetDomainIndex(), bigAdminReqSeq)

	xpNumReq := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, xpNumRequest)

	htpResponse, err := common.SubmitHTPRequest(xpNumReq)
	if err != nil {
		return err
	}

	_, err = ImportWKResponse(htpResponse, de)
	if err != nil {
		return err
	}

	return nil
}

/*----------------------------------------------------------------------------*/
/* Build an import wrapping key request for a single key part                 */
/*----------------------------------------------------------------------------*/
func ImportWKRequest(domainID []byte, moduleID []byte, transactionCounter []byte,
	recipientInfo []byte) AdminBlk {

	var adminBlock AdminBlk

	adminBlock.CmdID = XCP_ADM_IMPORT_WK
	adminBlock.DomainID = domainID
	adminBlock.ModuleID = moduleID
	adminBlock.TransactionCounter = transactionCounter
	adminBlock.CmdInput = recipientInfo
	return adminBlock
}

/*----------------------------------------------------------------------------*/
/* Build an import wrapping key request containing all of the inport wrapping */
/* requests for each individual key part                                      */
/*----------------------------------------------------------------------------*/
func BigImportWKRequest(domainID []byte, moduleID []byte, transactionCounter []byte,
	keyPartRequests []byte) AdminBlk {

	var adminBlock AdminBlk

	adminBlock.CmdID = XCP_ADM_IMPORT_WK
	adminBlock.DomainID = domainID
	adminBlock.ModuleID = moduleID
	adminBlock.TransactionCounter = transactionCounter
	adminBlock.CmdInput = keyPartRequests
	return adminBlock
}

/*----------------------------------------------------------------------------*/
/* Parse an import wrapping key response                                      */
/*----------------------------------------------------------------------------*/
func ImportWKResponse(htpResponse string, de common.DomainEntry) (AdminRspBlk, error) {
	rspBlk, err := buildAdminRspBlk(htpResponse, de)
	return rspBlk, err
}
