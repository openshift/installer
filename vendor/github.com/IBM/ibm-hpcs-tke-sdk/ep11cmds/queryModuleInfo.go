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
	"errors"

	"github.com/Logicalis/asn1"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/** Output fields from get_xcp_info to retrieve module information */
type ModuleInfoRspInfo struct {
	APIOrdinalNumber      []byte
	FirmwareIdentifier    []byte
	APIVersionMajor       byte
	APIVersionMinor       byte
	CSPVersionMajor       byte
	CSPVersionMinor       byte
	FirmwareConfiguration []byte
	XCPConfiguration      []byte
	SerialNumber          []byte
	// other fields omitted...see section 5.1.1 in XCP wire formats document
	// for full list
	SerialNumberString string
}

/*----------------------------------------------------------------------------*/
/* Issues get_xcp_info to retrieve module information                         */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the crypto module and domain to be queried       */
/*                                                                            */
/* Outputs:                                                                   */
/* ModuleInfoRspInfo -- returned data from the query                          */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryModuleInfo(authToken string, urlStart string,
	de common.DomainEntry) (ModuleInfoRspInfo, error) {

	htpRequestString := QueryModuleInfoRequest(
		de.GetCryptoModuleIndex(), de.GetDomainIndex())

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		var dummy ModuleInfoRspInfo
		return dummy, err
	}
	return QueryModuleInfoRsp(htpResponseString)
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for get_xcp_info with a request for module          */
/* information                                                                */
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
func QueryModuleInfoRequest(cryptoModuleIndex int, domainIndex int) string {

	var req XCPReq
	req.CmdId = FNID_GET_XCP_INFO
	req.DomainID = common.Uint32To4ByteSlice(uint32(domainIndex))
	req.CmdSubtype = CK_IBM_XCPQ_MODULE
	req.Unused = common.Uint32To4ByteSlice(uint32(0)) // set to 4 bytes of zeroes

	// Form an ASN.1 sequence of octet strings
	seq, err := asn1.Encode(req)
	if err != nil {
		panic(err)
	}
	return NewXPNUMRequest(cryptoModuleIndex, domainIndex, seq)
}

/*----------------------------------------------------------------------------*/
/* Processes the HTP response from get_xcp_data to retrieve module            */
/* information.                                                               */
/*                                                                            */
/* Input:                                                                     */
/* htpResponseString -- HTPResponse string from POST /hsms                    */
/*                                                                            */
/* ModuleInfoRspInfo -- structure containing output fields from the query     */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryModuleInfoRsp(htpResponseString string) (ModuleInfoRspInfo, error) {

	var rtnData ModuleInfoRspInfo

	payload, err := GetRspPayload(htpResponseString)
	if err != nil {
		return rtnData, err
	}

	var resp XCPRsp
	_, err = asn1.Decode(payload, &resp)
	if err != nil {
		return rtnData, err
	}

	if len(resp.Payload) < 196 {
		return rtnData, errors.New("Query module information response is too short.")
	}

	rtnData.APIOrdinalNumber      = resp.Payload[0:4]
	rtnData.FirmwareIdentifier    = resp.Payload[4:8]
	rtnData.APIVersionMajor       = resp.Payload[8]
	rtnData.APIVersionMinor       = resp.Payload[9]
	rtnData.CSPVersionMajor       = resp.Payload[10]
	rtnData.CSPVersionMinor       = resp.Payload[11]
	rtnData.FirmwareConfiguration = resp.Payload[12:44]
	rtnData.XCPConfiguration      = resp.Payload[44:76]
	rtnData.SerialNumber          = resp.Payload[108:124]
	rtnData.SerialNumberString    = string(rtnData.SerialNumber[0:8])

	return rtnData, nil
}
