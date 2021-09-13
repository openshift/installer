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

/*----------------------------------------------------------------------------*/
/* Reads an OA certificate                                                    */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the crypto module and domain to be queried       */
/* certificateIndex -- index into the certificate chain                       */
/*    0 = currently active epoch key, 1 = its parent, etc.                    */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the returned OA certificate                                      */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDeviceCertificate(authToken string, urlStart string,
	de common.DomainEntry, certificateIndex uint32) ([]byte, error) {

	htpRequestString := QueryDeviceCertificateReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex(), certificateIndex)

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
/* Creates the HTPRequest to return a specific OA certificate                 */
/*                                                                            */
/* Inputs:                                                                    */
/* cryptoModuleIndex -- identifies the crypto module to be queried            */
/* domainIndex -- the domain assigned to the user.  We know this is a control */
/*    domain, and it will get past the domain check when the cloud processes  */
/*    the POST /hsms request.                                                 */
/* certificateIndex -- index into the certificate chain                       */
/*    0 = currently active epoch key, 1 = its parent, etc.                    */
/*                                                                            */
/* Output:                                                                    */
/* string -- hexadecimal string representing the HTPRequest                   */
/*----------------------------------------------------------------------------*/
func QueryDeviceCertificateReq(cryptoModuleIndex int, domainIndex int,
	certificateIndex uint32) string {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADMQ_DEVICE_CERT
	adminBlk.DomainID = XCP_DOMAIN_0
	// module ID not used for queries
	// transaction counter not used for queries
	adminBlk.CmdInput = common.Uint32To4ByteSlice(certificateIndex)
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminBlk)
}

/*----------------------------------------------------------------------------*/
/* Returns the number of OA certificates in the OA certificate chain          */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the crypto module and domain to be queried       */
/*                                                                            */
/* Outputs:                                                                   */
/* uint32 -- the number of certificates in the OA certificate chain           */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryNumberDeviceCertificates(authToken string, urlStart string,
	de common.DomainEntry) (uint32, error) {

	htpRequestString := QueryNumberDeviceCertificatesReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex())

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return 0, err
	}

	adminRspBlk, err := buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(adminRspBlk.CmdOutput), nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest to return the number of OA certificates in the OA   */
/* certificate chain                                                          */
/*                                                                            */
/* Inputs:                                                                    */
/* cryptoModuleIndex -- identifies the crypto module to be queried            */
/* domainIndex -- the domain assigned to the user.  We know this is a control */
/*    domain, and it will get past the domain check when the cloud processes  */
/*    the POST /hsms request.                                                 */
/*                                                                            */
/* Output:                                                                    */
/* string -- hexadecimal string representing the HTPRequest                   */
/*----------------------------------------------------------------------------*/
func QueryNumberDeviceCertificatesReq(cryptoModuleIndex int,
	domainIndex int) string {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADMQ_DEVICE_CERT
	adminBlk.DomainID = XCP_DOMAIN_0
	// module ID not used for queries
	// transaction counter not used for queries
	// empty payload to get the number of certificates
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminBlk)
}
