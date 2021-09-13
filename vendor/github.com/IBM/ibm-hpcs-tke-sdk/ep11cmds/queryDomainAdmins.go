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
	"encoding/hex"
	"errors"
	"strings"

	"github.com/Logicalis/asn1"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Queries the domain administrators.                                         */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain to be queried                         */
/*                                                                            */
/* Outputs:                                                                   */
/* [][]byte -- array of Subject Key Identifiers (SKIs), one for each          */
/*    administrator installed in the domain                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDomainAdmins(authToken string, urlStart string,
	de common.DomainEntry) ([][]byte, error) {

	htpRequestString := QueryDomainAdminsReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex(), nil)

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString) //@TxxxxxxCLH

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return nil, err
	}

	_, skis, err := QueryDomainAdminsListRsp(htpResponseString, de)
	return skis, err
}

/*----------------------------------------------------------------------------*/
/* Retrieves the name of a domain administrator.                              */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain to be queried                         */
/* []byte -- Subject Key Identifier of the domain administrator of interest   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- name of the domain administrator                                 */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDomainAdminName(authToken string, urlStart string,
	de common.DomainEntry, ski []byte) (string, error) {

	htpRequestString := QueryDomainAdminsReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex(), ski)

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return "", err
	}

	_, cert := QueryDomainAdminRsp(htpResponseString, de)
	return strings.TrimSpace(string(cert.TheBody.TheIssuer.CommonName.TheName.PrintableString)), nil
}

/**
Create a query domain administrators request.  The aSKI parameter may contain the
SKI for an administrator or be nil.
*/
func QueryDomainAdminsReq(cryptoModuleIndex int, domainIndex int, aSKI []byte) string {
	var adminCmd AdminBlk

	//TODO
	/*
		if domainIndex == 21 {
			domainIndex = 0
			fmt.Println("changing domain index to 0...")
		}

		if cryptoModuleIndex == 06 {
			cryptoModuleIndex = 07
			fmt.Println("changing crypto module index to 07...")
		}
	*/

	adminCmd.CmdID = XCP_ADMQ_DOMADMIN
	adminCmd.DomainID = BuildAdminDomainIndex(domainIndex)
	if aSKI != nil {
		adminCmd.CmdInput = aSKI
	}
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminCmd)
}

/**
Process the response to a query domain administrators list request.
This function will return the admin response block plus an
array containing a SKI value for each domain administrator.
*/
func QueryDomainAdminsListRsp(theRsp string, de common.DomainEntry) (AdminRspBlk, [][]byte, error) {
	var SKIs [][]byte
	adminRspBlk, err := buildAdminRspBlk(theRsp, de)

	if err != nil {
		return adminRspBlk, nil, err
	}
	outputLength := len(adminRspBlk.CmdOutput)
	if (outputLength % 32) != 0 {
		panic(errors.New("Command output length is not a multiple of 32"))
	}
	numberOfSKIs := outputLength / 32
	SKIs = make([][]byte, numberOfSKIs)
	for i := 0; i < outputLength; i += 32 {
		index := i / 32
		SKIs[index] = adminRspBlk.CmdOutput[i : i+32]
	}
	return adminRspBlk, SKIs, nil
}

/*
Process the response to a query domain administrator request.
This function will return the admin response block plus a
certificate containing information about the domain administrator.
*/
func QueryDomainAdminRsp(theRsp string, de common.DomainEntry) (AdminRspBlk, Certificate) {
	adminRspBlk, err := buildAdminRspBlk(theRsp, de)

	if err != nil {
		panic(err)
	}

	var cert Certificate
	_, err2 := asn1.Decode(adminRspBlk.CmdOutput, &cert)

	if err2 != nil {
		panic(err2)
	}
	return adminRspBlk, cert
}

/*
Get the administrator name out of the certificate
*/
func (cert Certificate) GetAdminName() []byte {
	return cert.TheBody.TheIssuer.CommonName.TheName.PrintableString
}

func (cert Certificate) GetPublicKey() []byte {
	return cert.TheBody.ThePublicKey.ThePublicKey[1:]
}

func (cert Certificate) GetSKI() []byte {
	return cert.TheBody.TheExtensions.TheSeq1.TheSeq2.SKI[2:]
}

/*
Returns true if the certificate is one associated with TKE operations
*/
func (cert Certificate) IsTKECertificate() bool {

	tkeAdminBytes := []byte(tkeAdministration)
	tkeAdminHex := hex.EncodeToString(tkeAdminBytes)
	orgNameHex := hex.EncodeToString(cert.TheBody.TheIssuer.OrgName.TheName.PrintableString)

	return strings.Contains(orgNameHex, tkeAdminHex)
}

func setNewSliceToValue(length int, value byte) []byte {
	slice := make([]byte, length)
	for i := range slice {
		slice[i] = value
	}
	return slice
}
