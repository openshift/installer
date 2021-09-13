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
	"crypto/elliptic"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/Logicalis/asn1"
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

var delimiter byte = ';'

/** Represents an EP11 xcpAdminBlk */
type AdminBlk struct {
	CmdID              []byte
	DomainID           []byte
	ModuleID           []byte
	TransactionCounter []byte
	CmdInput           []byte
}

/** Represents an EP11 xcpAdminReq */
type AdminReq struct {
	CmdID      []byte
	DomainID   []byte
	AdminBlock []byte
	SignerInfo []byte
}

/** Represents an EP11 xcpAdminRsp (normal case) */
type AdminRsp struct {
	CmdID       []byte
	DomainID    []byte
	ReturnCode  []byte
	AdminRspBlk []byte
	SignerInfo  []byte
}

/** Represents an EP11 xcpAdminRsp (error case) */
type AdminRspError struct {
	CmdID      []byte
	DomainID   []byte
	ReturnCode []byte
}

/* Some errors on the EP11 crypto module are reported in the xcpAdminRsp
 * rather than in the xcpAdminRspBlk.  For these errors no xcpAdminRspBlk or
 * SignerInfo is returned in the xcpAdminRsp.  If we get an error decoding
 * into the first structure because of missing fields, decode into the second
 * structure to determine the return code. */

/** Represents an EP11 xcpAdminRspBlk */
type AdminRspBlk struct {
	CmdID              []byte
	DomainID           []byte
	ModuleID           []byte
	TransactionCounter []byte
	ReturnCode         []byte
	CmdOutput          []byte
}

/** Length of the public key in the OA certificate for a CEX6P crypto module */
const CEX6P_PUBLIC_KEY_LENGTH = 133

/** Function ID for get_xcp_info */
var FNID_GET_XCP_INFO = []byte{0x00, 0x01, 0x00, 0x26}

/** get_xcp_info subtype to retrieve module information */
var CK_IBM_XCPQ_MODULE = []byte{0x00, 0x00, 0x00, 0x01}

/** get_xcp_info subtype to retrieve domain information */
var CK_IBM_XCPQ_DOMAIN = []byte{0x00, 0x00, 0x00, 0x03}

/** XCPReq fields for get_xcp_info to retrieve module or domain information */
type XCPReq struct {
	CmdId      []byte
	DomainID   []byte
	CmdSubtype []byte
	Unused     []byte
}

/** XCPRsp fields for get_xcp_info to retrieve module or domain information */
type XCPRsp struct {
	CmdId      []byte
	DomainID   []byte
	ReturnCode []byte
	Payload    []byte
}

/**  Function ID for m_admin */
var FNID_ADMIN = []byte{0x00, 0x01, 0x00, 0x29}

var tkeAdministration = "TKE Administration"

type Certificate struct {
	TheBody   CertBody
	Algorithm AlgID
	Signature []byte `asn1:"universal,tag:3"`
}

type CertBody struct {
	TheVersion    CertVersion `asn1:"tag:0"`
	SerialNumber  int
	Algorithm     AlgID
	TheIssuer     Names
	Validity      TimeRange
	TheSubject    Names
	ThePublicKey  PublicKey
	TheExtensions Extensions `asn1:"tag:3"`
}

type CertVersion struct {
	Version int
}

type Names struct {
	OrgName    Name2 `asn1:"set"`
	CommonName Name2 `asn1:"set"`
}

type Name2 struct {
	TheName Name
}

type Name struct {
	OID             asn1.Oid
	PrintableString []byte `asn1:"universal,tag:19"`
}

type TimeRange struct {
	NotBefore []byte `asn1:"universal,tag:24"`
	NotAfter  []byte `asn1:"universal,tag:24"`
}

type PublicKey struct {
	Algorithm    AlgID
	ThePublicKey []byte `asn1:"universal,tag:3"`
}

//#B@T372621CLH
type PublicKeyECP521 struct { // used for EC P521 encode
	Algorithm    AlgID
	ThePublicKey []byte `asn1:"universal,tag:3"`
}
//#E@T372621CLH

type AlgID struct {
	ObjID  asn1.Oid
	Null   asn1.Null `asn1:"optional"`
	ObjID2 asn1.Oid  `asn1:"optional"`
}

type Extensions struct {
	TheSeq1 Seq1
}

type Seq1 struct {
	TheSeq2 Seq2
}

type Seq2 struct {
	ObjID asn1.Oid
	SKI   []byte
}

type RecipientInfo2048 struct {
	Version      int
	SKI          []byte `asn1:"tag:0"`
	AlgID        AlgIDRSA2048
	EncryptedKey []byte
}

/** Administrative function IDs */

/** Add administrator certificate */
var XCP_ADM_ADMIN_LOGIN = []byte{0x00, 0x00, 0x00, 0x01}

/** Add domain administrator certificate */
var XCP_ADM_DOM_ADMIN_LOGIN = []byte{0x00, 0x00, 0x00, 0x02}

/** Revoke administrator certificate */
var XCP_ADM_ADMIN_LOGOUT = []byte{0x00, 0x00, 0x00, 0x03}

/** Revoke domain administrator certificate */
var XCP_ADM_DOM_ADMIN_LOGOUT = []byte{0x00, 0x00, 0x00, 0x04}

/** Transition administrator certificate */
var XCP_ADM_ADMIN_REPLACE = []byte{0x00, 0x00, 0x00, 0x05}

/** Transition domain administrator certificate */
var XCP_ADM_DOM_ADMIN_REPLACE = []byte{0x00, 0x00, 0x00, 0x06}

/** Set card attribute/s */
var XCP_ADM_SET_ATTR = []byte{0x00, 0x00, 0x00, 0x07}

/** Set domain attribute/s */
var XCP_ADM_DOM_SET_ATTR = []byte{0x00, 0x00, 0x00, 0x08}

/** Generate new importer (PK) key */
var XCP_ADM_GEN_IMPORTER = []byte{0x00, 0x00, 0x00, 0x09}

/** Create random domain WK */
var XCP_ADM_GEN_WK = []byte{0x00, 0x00, 0x00, 0x0A}

/** Wrap+output WK or parts */
var XCP_ADM_EXPORT_WK = []byte{0x00, 0x00, 0x00, 0x0B}

/** Set (set of) WK (parts) to pending */
var XCP_ADM_IMPORT_WK = []byte{0x00, 0x00, 0x00, 0x0C}

/** Activate pending WK */
var XCP_ADM_COMMIT_WK = []byte{0x00, 0x00, 0x00, 0x0D}

/** Remove previous WK's */
var XCP_ADM_FINALIZE_WK = []byte{0x00, 0x00, 0x00, 0x0E}

/** Release CSPs from entire module */
var XCP_ADM_ZEROIZE = []byte{0x00, 0x00, 0x00, 0x0F}

/** Release CSPs from domain/s */
var XCP_ADM_DOM_ZEROIZE = []byte{0x00, 0x00, 0x00, 0x10}

/** Fix domain control points */
var XCP_ADM_DOM_CONTROLPOINT_SET = []byte{0x00, 0x00, 0x00, 0x11}

/** Enable domain control points */
var XCP_ADM_DOM_CONTROLPOINT_ADD = []byte{0x00, 0x00, 0x00, 0x12}

/** Disable domain control points */
var XCP_ADM_DOM_CONTROLPOINT_DEL = []byte{0x00, 0x00, 0x00, 0x13}

/** Set module-internal UTC time */
var XCP_ADM_SET_CLOCK = []byte{0x00, 0x00, 0x00, 0x14}

/** Set function-control vector */
var XCP_ADM_SET_FCV = []byte{0x00, 0x00, 0x00, 0x15}

/** Fix control points -- not currently used */
var XCP_ADM_CONTROLPOINT_SET = []byte{0x00, 0x00, 0x00, 0x16}

/** Enable control points -- not currently used */
var XCP_ADM_CONTROLPOINT_ADD = []byte{0x00, 0x00, 0x00, 0x17}

/** Disable control points -- not currently used */
var XCP_ADM_CONTROLPOINT_DEL = []byte{0x00, 0x00, 0x00, 0x18}

/** Transform blobs to next WK -- not used by TKE */
var XCP_ADM_REENCRYPT = []byte{0x00, 0x00, 0x00, 0x19}

/** Remove (semi-) retained key -- not used by TKE */
var XCP_ADM_RK_REMOVE = []byte{0x00, 0x00, 0x00, 0x1A}

/** Erase current WK */
var XCP_ADM_CLEAR_WK = []byte{0x00, 0x00, 0x00, 0x1B}

/** Erase pending WK */
var XCP_ADM_CLEAR_NEXT_WK = []byte{0x00, 0x00, 0x00, 0x1C}

/** Card zeroize, preserving system key, if it is present (i.e. retaining
 * the Support Element administrator) */
var XCP_ADM_SYSTEM_ZEROIZE = []byte{0x00, 0x00, 0x00, 0x1D}

/** Export state */
var XCP_ADM_EXPORT_STATE = []byte{0x00, 0x00, 0x00, 0x1E}

/** Import state (part) */
var XCP_ADM_IMPORT_STATE = []byte{0x00, 0x00, 0x00, 0x1F}

/** Commit imported state */
var XCP_ADM_COMMIT_STATE = []byte{0x00, 0x00, 0x00, 0x20}

/** Remove cloning state */
var XCP_ADM_REMOVE_STATE = []byte{0x00, 0x00, 0x00, 0x21}

/** Generate module imported key */
var XCP_ADM_GEN_MODULE_IMPORTER = []byte{0x00, 0x00, 0x00, 0x22}

/** Set trusted attribute on blob/SPKI */
var XCP_ADM_SET_TRUSTED = []byte{0x00, 0x00, 0x00, 0x23}

/** Multi-domain zeroize */
var XCP_ADM_DOMAINS_ZEROIZE = []byte{0x00, 0x00, 0x00, 0x24}

/** wrap+output next WK or parts */
var XCP_ADM_EXPORT_NEXT_WK = []byte{0x00, 0x00, 0x00, 0x26} //@T390301CLH

/** Query administrative SKI/certificate */
var XCP_ADMQ_ADMIN = []byte{0x00, 0x01, 0x00, 0x01}

/** Query domain administrative SKI/certificate */
var XCP_ADMQ_DOMADMIN = []byte{0x00, 0x01, 0x00, 0x02}

/** Query module CA (OA) certificate */
var XCP_ADMQ_DEVICE_CERT = []byte{0x00, 0x01, 0x00, 0x03}

/** Query card control points or profile */
var XCP_ADMQ_CTRLPOINTS = []byte{0x00, 0x01, 0x00, 0x05}

/** Query domain control points or profile */
var XCP_ADMQ_DOM_CTRLPOINTS = []byte{0x00, 0x01, 0x00, 0x06}

/** Query current wrapping key */
var XCP_ADMQ_WK = []byte{0x00, 0x01, 0x00, 0x07}

/** Query pending wrapping key */
var XCP_ADMQ_NEXT_WK = []byte{0x00, 0x01, 0x00, 0x08}

/** Query card attributes */
var XCP_ADMQ_ATTRS = []byte{0x00, 0x01, 0x00, 0x09}

/** Query domain attributes */
var XCP_ADMQ_DOM_ATTRS = []byte{0x00, 0x01, 0x00, 0x0A}

/** Query Function Control Vector */
var XCP_ADMQ_FCV = []byte{0x00, 0x01, 0x00, 0x0B}

/** Query information on original WK */
var XCP_ADMQ_WK_ORIGINS = []byte{0x00, 0x01, 0x00, 0x0C}

/** Query retained keys */
var XCP_ADMQ_RKLIST = []byte{0x00, 0x01, 0x00, 0x0D}

/** Query cloning state */
var XCP_ADMQ_INTERNAL_STATE = []byte{0x00, 0x01, 0x00, 0x0E}

/** Query current WK importer */
var XCP_ADMQ_IMPORTER_CERT = []byte{0x00, 0x01, 0x00, 0x0F}

/** Query cloning state */
var XCP_ADMQ_AUDIT_STATE = []byte{0x00, 0x01, 0x00, 0x10}

/** OCTET STRING representing domain 0 with all zero "instance identifier" */
var XCP_DOMAIN_0 = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

var RSAEncryption = []uint{1, 2, 840, 113549, 1, 1, 1}            // 1.2.840.113549.1.1.1
var SHA256WithRSAEncryption = []uint{1, 2, 840, 113549, 1, 1, 11} // 1.2.840.113549.1.1.11

/** Contains error information from the second section of an HTPResponse */
type RspInfo struct {
	ErrorType     string
	ReturnCode    string
	ReasonCode    string
	ProgramID     string
	ErrorLocation string
	ErrorText     string
}

/** Implements the first 32 bytes of an EP11 CPRB */
type CPRB struct {
	length        [2]byte
	version       byte
	reserved1     [2]byte
	flags         byte
	subtype       [2]byte
	partitionID   [4]byte
	domainID      [4]byte
	returnCode    [4]byte
	reserved2     [8]byte
	payloadLength [4]byte
}

/*----------------------------------------------------------------------------*/
/* Create a CPRB to send to a host system                                     */
/*----------------------------------------------------------------------------*/
func NewCPRB(domainIndex int, payloadParm []byte) []byte {
	newSlice := make([]byte, 0, 512)
	newSlice = append(newSlice, 0, 32) // CPRB length
	newSlice = append(newSlice, 4) // version
	newSlice = append(newSlice, 0, 0) // reserved field 1
	newSlice = append(newSlice, 0x80) // flags byte
	newSlice = append(newSlice, 'T', '4') // CPRB subtype
	newSlice = append(newSlice, 0, 0, 0, 0) // partition ID
	newSlice = append(newSlice, common.Uint32To4ByteSlice(uint32(domainIndex))...) // domain index
	newSlice = append(newSlice, 0, 0, 0, 0) // return code
	newSlice = append(newSlice, 0, 0, 0, 0, 0, 0, 0, 0) // reserved field 2
	newSlice = append(newSlice, common.Uint32To4ByteSlice(uint32(len(payloadParm)))...) // payload length
	newSlice = append(newSlice, payloadParm...)

	return newSlice
}

/*----------------------------------------------------------------------------*/
/* Create a HTPRequest specifying the XPNUM rule                              */
/*----------------------------------------------------------------------------*/
func NewXPNUMRequest(cryptoModuleIndex int, domainIndex int, sequence []byte) string {
	var stringBuilder strings.Builder
	reqHeader1 := ";PCI request;XPNUM   ;"
	cmiString := strconv.Itoa(cryptoModuleIndex)
	reqHeader2 := ";        ;"
	cprb := strings.ToUpper(hex.EncodeToString(NewCPRB(domainIndex, sequence)))
	reqLength := len(reqHeader1) + len(cmiString) + len(reqHeader2) + len(cprb)
	reqLengthString := strconv.Itoa(reqLength)
	stringBuilder.WriteString(reqLengthString)
	stringBuilder.WriteString(reqHeader1)
	stringBuilder.WriteString(cmiString)
	stringBuilder.WriteString(reqHeader2)
	stringBuilder.WriteString(cprb)
	return stringBuilder.String()
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest string for an unsigned query.                       */
/*                                                                            */
/* For unsigned queries, the administrative domain, the module identifier,    */
/* and the transaction counter in the xcpAdminReq do not need to be filled    */
/* in.                                                                        */
/*----------------------------------------------------------------------------*/
func CreateQueryHTPRequest(cryptoModuleIndex int, domainIndex int,
	adminBlock AdminBlk) string {

	var adminReq AdminReq
	adminReq.CmdID = FNID_ADMIN
	adminReq.DomainID = common.Uint32To4ByteSlice(uint32(domainIndex))
	// now form a sequence from the admin block
	seq1, anError := asn1.Encode(adminBlock)
	if anError != nil {
		panic(anError)
	}
	adminReq.AdminBlock = seq1
	adminReq.SignerInfo = nil // no signature required
	// now form a sequence from the admin request
	seq2, anError2 := asn1.Encode(adminReq)
	if anError2 != nil {
		panic(anError2)
	}
	return NewXPNUMRequest(cryptoModuleIndex, domainIndex, seq2)
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest string for a signed command.                        */
/*                                                                            */
/* For a signed command, the administrative domain, the module identifier,    */
/* and the transaction counter are determined by issuing a Query Domain       */
/* Attributes command.                                                        */
/*                                                                            */
/* The number of signature keys provided indicates the number of signatures   */
/* that need to be collected for the command.                                 */
/*----------------------------------------------------------------------------*/
func CreateSignedHTPRequest(authToken string, urlStart string, de common.DomainEntry,
	adminBlock AdminBlk, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	// Issue Query Domain Attributes to get the administrative domain, the
	// module identifier, and the transaction counter.
	_, adminRspBlk, err := QueryDomainAttributes(authToken, urlStart, de)
	if err != nil {
		return "", err
	}

	// Transfer returned fields to the input xcpAdminBlk
	adminBlock.DomainID = adminRspBlk.DomainID
	adminBlock.ModuleID = adminRspBlk.ModuleID
	adminBlock.TransactionCounter =
		IncrementTransactionCounter(adminRspBlk.TransactionCounter)

	// Create the xcpAdminBlk sequence
	adminBlockSeq, err := asn1.Encode(adminBlock)
	if err != nil {
		panic(err)
	}

	signerInfo, err := CreateSignerInfo(adminBlockSeq, sigkeys, sigkeySkis,
		sigkeyTokens)
	if err != nil {
		return "", err
	}

	// Create the xcpAdminReq sequence
	var adminReq AdminReq
	adminReq.CmdID = FNID_ADMIN
	adminReq.DomainID = adminRspBlk.DomainID[0:4]
	adminReq.AdminBlock = adminBlockSeq
	adminReq.SignerInfo = signerInfo

	adminReqSeq, err := asn1.Encode(adminReq)
	if err != nil {
		panic(err)
	}

	// Form and return the final HTPRequest string
	return NewXPNUMRequest(de.GetCryptoModuleIndex(), de.GetDomainIndex(), adminReqSeq), nil
}

/*----------------------------------------------------------------------------*/
/* Increments the transaction counter                                         */
/*----------------------------------------------------------------------------*/
func IncrementTransactionCounter(original []byte) []byte {

	// Converted to Go from increment in com.ibm.tke.util.ASN1.java

	// The behavior of the Go big.Int type is different from that of the Java
	// BigInteger class.  Some cases to consider:
	// { 0x00, 0x00 } increments to { 0x01 }; need to pad on left to preserve
	//      correct length.
	// { 0xFF, 0xFF } increments to { 0x01, 0x00, 0x00 }; need to strip off
	//      the leading byte.
	// { 0xFF, 0xFD } increments to { 0xFF, 0xFE }; no special handling needed.
	// { 0x7F, 0xFF } increments to { 0x80, 0x00 }; no special handling needed.

	tc := new(big.Int).SetBytes(original)
	tc.Add(tc, big.NewInt(1))

	incremented := tc.Bytes()
	if len(incremented) == len(original) {
		return incremented
	} else if len(incremented) < len(original) {
		rtnbytes := make([]byte, len(original))
		copy(rtnbytes[len(original)-len(incremented):], incremented)
		// leftmost bytes are 0x00
		return rtnbytes
	} else {
		return incremented[1:] // remove first byte
	}
}

/*----------------------------------------------------------------------------*/
/* Initializes the administrative domain field of an xcpAdminBlk.             */
/*                                                                            */
/* The first four bytes contains the domain index.  The last four bytes       */
/* contains the domain instance identifier.  For unsigned commands and        */
/* module-level commands, the domain instance identifier can be all zeros.    */
/* For signed commands to a domain, the domain instance identifier is filled  */
/* in later using data returned by the Query Domain Attributes command.       */
/*----------------------------------------------------------------------------*/
func BuildAdminDomainIndex(domainIndex int) []byte {
	domainIndexBytes := common.Uint32To4ByteSlice(uint32(domainIndex))
	domainIndexBytes = append(domainIndexBytes, 0, 0, 0, 0)
	return domainIndexBytes
}

/*----------------------------------------------------------------------------*/
/* Processes an HTPResponse containing an xcpAdminRsp.                        */
/*                                                                            */
/* Checks the return value in the xcpAdminRsp.                                */
/* Isolates the xcpAdminRspBlk and SignerInfo and verifies the OA signature.  */
/* Checks the return value in the xcpAdminRspBlk.                             */
/* Returns the xcpAdminRspBlk if no errors are found.                         */
/*                                                                            */
/* Only valid for HTPResponses containing an EP11 response CPRB.              */
/*                                                                            */
/* Inputs:                                                                    */
/* htpResponse string -- HTPResponse string from POST /hsms                   */
/* DomainEntry -- identifies the target domain for the request.  The          */
/*    DomainEntry contains the public key for verifying the OA signature.     */
/*                                                                            */
/* Outputs:                                                                   */
/* AdminRspBlk -- structure containing fields from the xcpAdminRspBlk         */
/* error -- identifies a terminating condition the caller should report       */
/*----------------------------------------------------------------------------*/
func buildAdminRspBlk(htpResponse string, de common.DomainEntry) (AdminRspBlk, error) {
	var adminRspBlk AdminRspBlk

	asn1Data, err := GetRspPayload(htpResponse)
	if err != nil {
		return adminRspBlk, err
	}

	// Decode the xcpAdminRsp sequence
	var adminRsp AdminRsp
	_, err = asn1.Decode(asn1Data, &adminRsp)
	if err != nil {
		// The decode will fail if the return code in the xcpAdminRsp is
		// nonzero.  In that case no xcpAdminRspBlk or SignerInfo is
		// returned and there aren't enough returned fields to fill the
		// AdminRsp structure.
		//
		// Decode into an alternate structure to determine the return code.
		var adminRspErr AdminRspError
		_, err = asn1.Decode(asn1Data, &adminRspErr)
		if err != nil {
			return adminRspBlk, err
		}
		retval := int(binary.BigEndian.Uint32(adminRspErr.ReturnCode))
		if retval != 0 {
			errorMsg := GetEP11ErrorMsg(strconv.Itoa(retval), "")
			if errorMsg == "" {
				return adminRspBlk, NewVerbError(
					"Error reported in xcpAdminRsp." +
					"\nReturn value: " + strconv.Itoa(retval), retval, -1)
			} else {
				return adminRspBlk, NewVerbError(
					"Error reported in xcpAdminRsp." +
					"\nReturn value: " + strconv.Itoa(retval) +
					"\nError message: " + errorMsg, retval, -1)
			}
		} else {
			return adminRspBlk, errors.New(
				"xcpAdminRsp returnValue is zero, but following fields are missing.")
		}
	}

	// Check the return value in the xcpAdminRsp, even when all five fields
	// are present in the xcpAdminRsp.
	retval := int(binary.BigEndian.Uint32(adminRsp.ReturnCode))
	if retval != 0 {
		errorMsg := GetEP11ErrorMsg(strconv.Itoa(retval), "")
		if errorMsg == "" {
			return adminRspBlk, NewVerbError(
				"Error reported in xcpAdminRsp." +
				"\nReturn value: " + strconv.Itoa(retval), retval, -1)
		} else {
			return adminRspBlk, NewVerbError(
				"Error reported in xcpAdminRsp." +
				"\nReturn value: " + strconv.Itoa(retval) +
				"\nError message: " + errorMsg, retval, -1)
		}
	}

	// When querying the epoch key certificate, we don't have the public key
	// used to verify the OA signature.  For this special case the Public_key
	// in the DomainEntry is set to "not available".
	if de.Public_key != "not available" {
		// For all other cases, verify the OA signature in the xcpAdminRsp

		// Decode the SignerInfo sequence
		signerInfo, err := DecodeSignerInfo(adminRsp.SignerInfo)
		if err != nil {
			return adminRspBlk, err
		}

		// Check for expected signature algorithm type
		if !common.ByteSlicesAreEqual(
				signerInfo.SignatureAlgorithmID, OID_ecdsaWithSHA512) {
			return adminRspBlk, errors.New("Unsupported signature algorithm in xcpAdminRsp")	
		}

		// Calculate the SHA-512 hash of the xcpAdminRspBlk
		hasher := sha512.New()
		hasher.Write(adminRsp.AdminRspBlk)
		sha512hash := hasher.Sum(nil)

		// Get the public key from the DomainEntry
		publicKey, err := hex.DecodeString(de.Public_key)
		if err != nil {
			panic(err)
		}
		var x, y big.Int
		qlen := binary.BigEndian.Uint16(publicKey[24:26])
		if len(publicKey) == CEX6P_PUBLIC_KEY_LENGTH {
			// Handle public key from CEX6P
			x.SetBytes(publicKey[1 : 67])
			y.SetBytes(publicKey[67 : 133])
		} else if len(publicKey) == int(26 + qlen) {
			// Handle public key from CEX5P
			pointLen := (qlen - 1) / 2
			x.SetBytes(publicKey[27 : 27 + pointLen])
			y.SetBytes(publicKey[27 + pointLen :])
		} else {
			return adminRspBlk, errors.New("Unrecognized format of crypto module public key")
		}
		pubkey := ecdsa.PublicKey{elliptic.P521(), &x, &y}

		// Verify the signature
		var r, s big.Int
		r.SetBytes(signerInfo.SignatureR)
		s.SetBytes(signerInfo.SignatureS)
		if !ecdsa.Verify(&pubkey, sha512hash, &r, &s) {
			return adminRspBlk, errors.New("Invalid signature in xcpAdminRsp")
		}
	}

	// Decode the xcpAdminRspBlk sequence
	_, err = asn1.Decode(adminRsp.AdminRspBlk, &adminRspBlk)
	if err != nil {
		return adminRspBlk, err
	}

	// Exit if the return value in the xcpAdminRspBlk is not 0.
	returnCode := int(binary.BigEndian.Uint32(adminRspBlk.ReturnCode))
	if returnCode != 0 {
		returnCodeString := strconv.Itoa(returnCode)
		// For some errors, a reason code is returned in the first four
		// bytes of the payload.
		reasonCode := -1
		var reasonCodeString string
		if len(adminRspBlk.CmdOutput) >= 4 {
			reasonCode = int(binary.BigEndian.Uint32(adminRspBlk.CmdOutput[0:4]))
			reasonCodeString = strconv.Itoa(reasonCode)
		}
		errorMsg := GetEP11ErrorMsg(returnCodeString, reasonCodeString)
		if errorMsg == "" {
			return adminRspBlk, NewVerbError(
				"Error reported by EP11 crypto module."+
					"\nReturn code: "+returnCodeString+
					"\nReason code: "+reasonCodeString,
				returnCode, reasonCode)
		} else {
			return adminRspBlk, NewVerbError(
				"Error reported by EP11 crypto module."+
					"\nReturn code: "+returnCodeString+
					"\nReason code: "+reasonCodeString+
					"\nError message: "+errorMsg,
				returnCode, reasonCode)
		}
	}
	return adminRspBlk, nil
}

// rsp error type
// rsp return code
// rsp reason code
/*----------------------------------------------------------------------------*/
/* Extracts the error information and EP11 response CPRB from an HTPResponse. */
/*                                                                            */
/* Only valid for HTPResponses containing an EP11 response CPRB.              */
/*                                                                            */
/* The HTPResponse string contains three fields, separated by ";".            */
/* 1. A length field (ignored here)                                           */
/* 2. Error information                                                       */
/* 3. The EP11 response CPRB                                                  */
/*                                                                            */
/* Inputs:                                                                    */
/* theRsp string -- HTPResponse string from POST /hsms                        */
/*                                                                            */
/* Outputs:                                                                   */
/* RspInfo -- structure containing error information from the HTPResponse     */
/* string -- the EP11 response CPRB                                           */
/* error -- identifies any errors found parsing the HTPResponse               */
/*----------------------------------------------------------------------------*/
func ParseResponse(theRsp string) (RspInfo, string, error) {
	var rspInfo RspInfo

	stringParts := strings.Split(theRsp, ";")
	// ignore stringParts[0]
	if len(stringParts) < 2 {
		return rspInfo, "", errors.New("Invalid HTPResponse, error information not found.")
	}
	if len(stringParts[1]) < 156 {
		return rspInfo, "", errors.New("Invalid HTPResponse, invalid error information.")
	}
	rspInfo.ErrorType = stringParts[1][0:2]
	rspInfo.ReturnCode = stringParts[1][2:10]
	rspInfo.ReasonCode = stringParts[1][10:18]
	rspInfo.ProgramID = stringParts[1][18:26]
	rspInfo.ErrorLocation = stringParts[1][26:56]
	rspInfo.ErrorText = stringParts[1][56:156]
	if rspInfo.ErrorType != "00" || len(stringParts) < 3 {
		return rspInfo, "", nil
	}
	return rspInfo, stringParts[2], nil
}

/*
Remove the front part of the HTPResponse (including the CPRB) and just
return the payload
*/
func GetRspPayload(theRsp string) ([]byte, error) {
	var asn1Data []byte
	// Extract error information and the EP11 response CPRB
	rspInfo, cprbString, err := ParseResponse(theRsp)
	if err != nil {
		return asn1Data, err
	}
	// Exit if the TKE catcher program reports an error
	if rspInfo.ErrorType != "00" {
		return asn1Data, errors.New(
			"HTPResponse error." +
				"\nError message:  " + GetErrorMsg(rspInfo.ErrorType, rspInfo.ReturnCode, rspInfo.ReasonCode) +
				"\nProgram ID:     " + rspInfo.ProgramID +
				"\nError location: " + rspInfo.ErrorLocation +
				"\nError text:     " + rspInfo.ErrorText)
	}
	cprbData, err2 := hex.DecodeString(cprbString)
	if err2 != nil {
		return asn1Data, err2
	}
	// Parse the CPRB
	cprb := cprbData[0:32]
	var returnCPRB = new(CPRB)
	copy(returnCPRB.length[:], cprb[0:2])
	returnCPRB.version = cprb[2]
	copy(returnCPRB.reserved1[:], cprb[3:5])
	returnCPRB.flags = cprb[5]
	copy(returnCPRB.subtype[:], cprb[6:8])
	copy(returnCPRB.partitionID[:], cprb[8:12])
	copy(returnCPRB.domainID[:], cprb[12:16])
	copy(returnCPRB.returnCode[:], cprb[16:20])
	// Exit if return code in the CPRB is not zero.
	returnCode := int(binary.BigEndian.Uint16(cprb[16:18]))
	reasonCode := int(binary.BigEndian.Uint16(cprb[18:20]))
	// TKEHTP.send breaks the four-byte field in the CPRB into a
	// return code and a reason code.
	if returnCode != 0 || reasonCode != 0 {
		return asn1Data, errors.New(
			"Error reported in EP11 CPRB." +
				"\nReturn code: " + strconv.Itoa(returnCode) +
				"\nReason code: " + strconv.Itoa(reasonCode))
	}
	copy(returnCPRB.reserved2[:], cprb[20:28])
	copy(returnCPRB.payloadLength[:], cprb[28:32])
	asn1DataLength := common.FourByteSliceToInt(cprb[28:32])
	asn1DataEndOffset := asn1DataLength + 32

	asn1Data = cprbData[32:asn1DataEndOffset]
	return asn1Data, nil
}

/*----------------------------------------------------------------------------*/
/* Extracts the crypto module serial number from the moduleIdentifier field   */
/* of an xcpAdminRspBlk.                                                      */
/*----------------------------------------------------------------------------*/
func (resp *AdminRspBlk) GetSerialNumber() string {
	return string(resp.ModuleID[0:8])
}
