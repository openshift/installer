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

var MK_STATUS_EMPTY int = 0
var CMK_STATUS_VALID int = 1
var NMK_STATUS_FULL_UNCOMMITTED int = 2
var NMK_STATUS_FULL_COMMITTED int = 3

type DomainInfoRspInfo struct {
	CurrentMKVP     []byte
	NewMKVP         []byte
	CurrentMKStatus int
	NewMKStatus     int
}

/*----------------------------------------------------------------------------*/
/* Queries the domain master key register status and verification pattern     */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- the authority token to use for the request                    */
/* urlStart -- the base URL to use for the request                            */
/* DomainEntry -- identifies the domain to be queried                         */
/*                                                                            */
/* Outputs:                                                                   */
/* DomainInfoRspInfo -- contains the status and verification patterns of the  */
/*    new and current master key registers                                    */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDomainInfo(authToken string, urlStart string,
	de common.DomainEntry) (DomainInfoRspInfo, error) {

	htpRequestString := QueryDomainInfoRequest(
		de.GetCryptoModuleIndex(), de.GetDomainIndex())

	req := common.CreatePostHsmsRequest(
		authToken, urlStart, de.Crypto_instance_id, de.Hsm_id, htpRequestString)

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		var dummy DomainInfoRspInfo
		return dummy, err
	}

	return QueryDomainInfoRsp(htpResponseString)
}

func QueryDomainInfoRequest(cryptoModuleIndex int, domainIndex int) string {

	var req XCPReq
	req.CmdId = FNID_GET_XCP_INFO
	req.DomainID = common.Uint32To4ByteSlice(uint32(domainIndex))
	req.CmdSubtype = CK_IBM_XCPQ_DOMAIN
	req.Unused = common.Uint32To4ByteSlice(uint32(0)) // set to 4 bytes of zeroes

	// Form an ASN.1 sequence of octet strings
	seq1, anError := asn1.Encode(req)
	if anError != nil {
		panic(anError)
	}
	return NewXPNUMRequest(cryptoModuleIndex, domainIndex, seq1)
}

/*
Process the response from get_xcp_info to retrieve domain information.
*/
func QueryDomainInfoRsp(theRsp string) (DomainInfoRspInfo, error) {
	var returnInfo DomainInfoRspInfo

	asn1Data, err := GetRspPayload(theRsp)

	if err != nil {
		return returnInfo, err
	}

	var resp XCPRsp

	_, err2 := asn1.Decode(asn1Data, &resp)

	if err2 != nil {
		return returnInfo, err2
	}
	if len(resp.Payload) < 80 {
		return returnInfo, errors.New("Query domain information response is too short.")
	}
	/*
			30820066
		  040400010026
		  040400000000
		  040400000000
		  04820050
		    00000000

		    2058C870E9D3194F4200EAAFD5BDDD40 current VP 32 bytes
		    2BAE0E6620E1B4BD1CB2E8CAB95211F7

		    00000000000000000000000000000000
		    00000000000000000000000000000000 pending VP 32 bytes

		    80000003  flags     00000002 -- current MK present if on, 00000004 -- next wrapping key present, can't be committed if on,
		                        00000008 -- next wrapping key present and committed

		    00000000  operational mode

		    00000001
	*/
	returnInfo.CurrentMKVP = make([]byte, 32)
	returnInfo.NewMKVP = make([]byte, 32)
	copy(returnInfo.CurrentMKVP[0:32], resp.Payload[4:36])
	copy(returnInfo.NewMKVP[0:32], resp.Payload[36:68])
	// current MK status -- empty or valid
	// new MK status -- empty, full uncommitted, full committed
	mkStatus := resp.Payload[71]
	returnInfo.CurrentMKStatus = MK_STATUS_EMPTY
	returnInfo.NewMKStatus = NMK_STATUS_FULL_COMMITTED
	if (mkStatus & 0x02) > 0 {
		returnInfo.CurrentMKStatus = CMK_STATUS_VALID
	}
	if (mkStatus & 0x04) == 0 {
		returnInfo.NewMKStatus = MK_STATUS_EMPTY
	} else {
		if (mkStatus & 0x08) == 0 {
			returnInfo.NewMKStatus = NMK_STATUS_FULL_UNCOMMITTED
		}
	}
	return returnInfo, nil
}
