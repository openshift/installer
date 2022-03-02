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
	"crypto/ecdsa"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Exports the current wrapping key register using a single key part          */
/*                                                                            */
/* A more complicated form of this could be written to support multiple key   */
/* parts, but for our purposes it isn't necessary.  The more complicated      */
/* version would need to have an array of KPH certificates as input, and      */
/* specify the M policy (number of key parts needed to reconstruct the key).  */
/* The output would be a concatenated set of RecipientInfo structures.        */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose current wrapping key register   */
/*    is to be exported                                                       */
/* []byte -- parameter file with format described in section 5.3 ("Serialized */
/*    module state") of the EP11 wire formats document.  Contains inputs to   */
/*    the Export WK command, such as the M policy and KPH certificates to use.*/
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- single RecipientInfo structure containing the encrypted key part */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func ExportWK(authToken string, urlStart string, de common.DomainEntry,
	pfile []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) ([]byte, error) {

	htpRequestString, err := ExportWKReq(authToken, urlStart, de, pfile,
		sigkeys, sigkeySkis, sigkeyTokens)
	if err != nil {
		return nil, err
	}

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
/* Creates the HTPRequest for exporting the current wrapping key register     */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose current wrapping key register   */
/*    is to be exported                                                       */
/* []byte -- parameter file with format described in section 5.3 ("Serialized */
/*    module state") of the EP11 wire formats document.  Contains inputs to   */
/*    the Export WK command, such as the M policy and KPH certificates to use.*/
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the HTPRequest string with the signed CPRB for the command       */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func ExportWKReq(authToken string, urlStart string, de common.DomainEntry,
	pfile []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_EXPORT_WK
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	adminBlk.CmdInput = pfile
	return CreateSignedHTPRequest(authToken, urlStart, de, adminBlk, sigkeys,
		sigkeySkis, sigkeyTokens)
}

/*----------------------------------------------------------------------------*/
/* Exports the pending wrapping key register                                  */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose pending wrapping key register   */
/*    is to be exported                                                       */
/* []byte -- parameter file with format described in section 5.3 ("Serialized */
/*    module state") of the EP11 wire formats document.  Contains inputs to   */
/*    the Export WK command, such as the M policy and KPH certificates to use.*/
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- raw encrypted key parts                                          */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func ExportPendingWK(authToken string, urlStart string, de common.DomainEntry,
	pfile []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) ([]byte, error) {

	htpRequestString, err := ExportPendingWKReq(authToken, urlStart, de, pfile,
		sigkeys, sigkeySkis, sigkeyTokens)
	if err != nil {
		return nil, err
	}

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
/* Creates the HTPRequest for exporting the pending wrapping key register     */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain whose pending wrapping key register   */
/*    is to be exported                                                       */
/* []byte -- parameter file with format described in section 5.3 ("Serialized */
/*    module state") of the EP11 wire formats document.  Contains inputs to   */
/*    the Export WK command, such as the M policy and KPH certificates to use.*/
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the HTPRequest string with the signed CPRB for the command       */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func ExportPendingWKReq(authToken string, urlStart string, de common.DomainEntry,
	pfile []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_EXPORT_NEXT_WK
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	adminBlk.CmdInput = pfile
	return CreateSignedHTPRequest(authToken, urlStart, de, adminBlk,
		sigkeys, sigkeySkis, sigkeyTokens)
}

/*----------------------------------------------------------------------------*/
/* Construct and return a KPH certificate containing the given P521 EC public */
/* key using the proprietary TKE format.                                      */
/*----------------------------------------------------------------------------*/
func KPHCert(pubKey ecdsa.PublicKey) []byte {
	var cert [309]byte
	cert[36] = 'K'  // certificate type (KPH)
	cert[37] = 1    // M policy
	cert[38] = 1    // N policy
	cert[40] = 0    // Prime curve
	cert[41] = 0x85 // public key length (0x85 = 133)
	cert[42] = 0x02 // curve size (0x0209 = 521)
	cert[43] = 0x09
	cert[44] = 0x04 // compression byte of public key

	bytes := pubKey.X.Bytes()
	length := len(bytes)
	// copy the X coordinate into the certificate
	for i:=0; i<length; i++ {
		cert[45 + (66 - length) + i] = bytes[i]
	}

	bytes = pubKey.Y.Bytes()
	length = len(bytes)
	// copy the Y coordinate into the certificate
	for i:=0; i<length; i++ {
		cert[45 + 66 + (66 - length) + i] = bytes[i]
	}

	return cert[:]
}

/*----------------------------------------------------------------------------*/
/* Construct a parameter file for the Export WK request.                      */
/*                                                                            */
/* Parameter files are described in section 5.3 "Serialized module state" in  */
/* the EP11 wire formats document.  A parameter file is an ASN.1 sequence of  */
/* octet strings, where the data in each octet string is a two-byte type,     */
/* followed by a four-byte index, followed by variable length data.  The data */
/* may be omitted (have zero length).                                         */
/*                                                                            */
/* Input:                                                                     */
/* []byte -- KPH certificate to use to encrypt key parts.  Can use the same   */
/*    certificate for multiple key parts.                                     */
/*                                                                            */
/* Output:                                                                    */
/* []byte -- the parameter file                                               */
/*----------------------------------------------------------------------------*/
func ExportWKParameterFile(kphcert []byte) []byte {
	pMap := common.NewParameterMap()
	// Export using two key parts to work around EP11 attribute problem
	pMap.Put(common.PMTAG_M_POLICY, 1, make([]byte, 0))
	pMap.Put(common.PMTAG_KPH_CERTIFICATE, 0, kphcert)
	//	pMap.Put(common.PMTAG_KPH_CERTIFICATE, 1, kphcert)
	pMap.Put(common.PMTAG_STATE_SCOPE, 0x0000000C, make([]byte, 0))
		// Don't include KPH and OA certificates in output
		// If this parameter is omitted, the resulting output is larger than
		// what the TKE catcher program can handle.
	return pMap.GenerateBytes()
}
