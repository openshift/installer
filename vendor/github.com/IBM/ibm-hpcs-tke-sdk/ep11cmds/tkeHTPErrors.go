//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 12/08/2020    CLH             T390301 - Add minimal touch functions

package ep11cmds

import (
	"strings"
)

var msgMap map[string]string
var msgInit bool = false

func GetErrorMsg(errorType string, returnCode string, reasonCode string) string {
	if !msgInit {
		initializeMsgs()
	}
	errorType = strings.TrimSpace(errorType)
	returnCode = strings.TrimSpace(returnCode)
	reasonCode = strings.TrimSpace(reasonCode)
	key := errorType + "-" + returnCode + "-" + reasonCode
	msgResult, found := msgMap[key]
	if found == false {
		msgResult = "Error type: " + errorType +
			", return code: " + returnCode +
			", reason code: " + reasonCode
	}
	return msgResult
}

/*----------------------------------------------------------------------------*/
/* Returns message text for an error reported by an EP11 crypto module.       */
/*                                                                            */
/* On TKE, the message type for these errors is "3".                          */
/* The EP11 crypto module does not always supply a reason code.               */
/* When both a return code and a reason code are supplied, look for the       */
/* specific pair first.  If not found, look for just the return code.         */
/* If no message is found, returns an empty string.                           */
/*----------------------------------------------------------------------------*/
func GetEP11ErrorMsg(returnCode string, reasonCode string) string {
	if !msgInit {
		initializeMsgs()
	}
	returnCode = strings.TrimSpace(returnCode)
	reasonCode = strings.TrimSpace(reasonCode)
	if reasonCode == "" {
		return msgMap["3-" + returnCode + "--"]
	} else {
		key := "3-" + returnCode + "-" + reasonCode
		msgResult, found := msgMap[key]
		if found {
			return msgResult
		} else {
			return msgMap["3-" + returnCode + "--"]
		}
	}
}

func initializeMsgs() {
	if msgInit {
		return
	}
	var locMsgMap map[string]string = make(map[string]string)

	/************************************************************************
	# 9     Return / reason codes returned from  TKEHTP - syntax errors     *
	#************************************************************************/
	locMsgMap["9-1-0"] = "Crypto Module index is not an integer."
	locMsgMap["9-2-0"] = "Crypto Module index is not within the allowed range."
	locMsgMap["9-3-0"] = "Crypto module information length is not valid."
	locMsgMap["9-4-0"] = "Domain index is not an integer."
	locMsgMap["9-5-0"] = "Domain index is not within the allowed range."
	locMsgMap["9-6-0"] = "Invalid length of Domain information."
	locMsgMap["9-7-0"] = "Invalid length of Authority information."
	locMsgMap["9-8-0"] = "Authority index is not an integer."
	locMsgMap["9-9-0"] = "Authority index is not within the allowed range."
	locMsgMap["9-10-0"] = "Crypto module information already present when attempting to Add."
	locMsgMap["9-11-0"] = "Crypto module information not present when attempting to Update."
	locMsgMap["9-12-0"] = "Invalid hex data in PKSC command."
	locMsgMap["9-13-0"] = "Length of PKSC command is not even."
	locMsgMap["9-14-0"] = "Invalid length of ENQ pending information."
	locMsgMap["9-16-0"] = "Invalid value of Set PCR value."
	locMsgMap["9-17-0"] = "Invalid length of Set PCR information."
	locMsgMap["9-18-0"] = "Invalid length of ENQ file information."
	locMsgMap["9-19-0"] = "Length of PKA request not valid (not 3000 or 6000 bytes)."
	locMsgMap["9-20-0"] = "Invalid hexadecimal data in PKA request."
	locMsgMap["9-21-0"] = "Dataset name not specified."
	locMsgMap["9-22-0"] = "Version number not 01 in PKA request."
	locMsgMap["9-23-0"] = "Length field in PKA request is not 1500."
	locMsgMap["9-24-0"] = "Length of RSA token field is not 1024 bytes"
	locMsgMap["9-25-0"] = "RSA token type is not external."
	locMsgMap["9-26-0"] = "Invalid length RSA external token(greater than 1024 bytes)."
	locMsgMap["9-27-0"] = "Invalid TCP/IP message length."
	locMsgMap["9-28-0"] = "User info missing in TCP/IP message."
	locMsgMap["9-29-0"] = "Invalid hexadecimal data in TCP/IP user info."
	locMsgMap["9-30-0"] = "Length of user info is not even."
	locMsgMap["9-31-0"] = "Different signature keys (Cipher TCP/IP user info is greater than modulus)."
	locMsgMap["9-32-0"] = "Different signature keys (Cipher TCP/IP user info is to long)."
	locMsgMap["9-33-0"] = "Begin session. Seq.no missing or incorrect."
	locMsgMap["9-34-0"] = "No session established for user."
	locMsgMap["9-35-0"] = "Invalid RACF request - not Env or UserChk. Software error on host."
	locMsgMap["9-36-0"] = "Invalid parameter at RACF request UserChk. Software error on host. (RACF parameter is not 'Change or NoChange')"
	locMsgMap["9-37-0"] = "Invalid parameter at RACF request UserChk. Software error on host. (RACF method is not 'Clear' or 'Cipher')"
	locMsgMap["9-38-0"] = "Invalid parameter at RACF request Userchk. Software error on host. (Length of encrypted User Id is not 16 bytes)"
	locMsgMap["9-39-0"] = "User authentication data contains an invalid hex digit."
	locMsgMap["9-40-0"] = "Crypto module index is not a whole number."
	locMsgMap["9-41-0"] = "Crypto module index is out of range."
	locMsgMap["9-42-0"] = "Crypto module serial number length is not valid."
	locMsgMap["9-43-0"] = "Invalid hex data in request parameter block."
	locMsgMap["9-44-0"] = "Invalid hex data in request parameter data block."
	locMsgMap["9-47-0"] = "Crypto module letter part is not C or P."
	locMsgMap["9-48-0"] = "Authorization index is missing."
	locMsgMap["9-49-0"] = "Domain index is missing."
	locMsgMap["9-50-0"] = "Crypto module index is not a whole number."
	locMsgMap["9-51-0"] = "Crypto module index is out of range."
	locMsgMap["9-52-0"] = "Invalid RACF user ID format."
	locMsgMap["9-53-0"] = "User ID is not 8 bytes long."
	locMsgMap["9-54-0"] = "User authentication data should be blank."
	locMsgMap["9-55-0"] = "User authentication data is not 8 bytes long."
	locMsgMap["9-56-0"] = "Version number not 02 in PKA request."
	locMsgMap["9-57-0"] = "Length field in PKA request is not 3000."
	locMsgMap["9-58-0"] = "Length of RSA token field is not 2524 bytes."
	locMsgMap["9-59-0"] = "Invalid length RSA external token (greater than 2524 bytes)."

	/*************************************************************************
	# 10     Return / reason codes returned from TSO                        *
	#************************************************************************/
	locMsgMap["10-1--"] = "TSO LISTDSI error. The reason code is the LISTDSI return code."
	locMsgMap["10-2--"] = "TSO ALLOC error. The reason code is the ALLOC return code."
	locMsgMap["10-3--"] = "TSO FREE error. The reason code is the FREE return code."
	locMsgMap["10-4--"] = "TSO EXECIO error. The reason code is the EXECIO return code."
	locMsgMap["10-4-1"] = "TSO EXECIO error. Data was truncated on write to disk. The target data set attributes may be incorrect."

	/*************************************************************************
	# 11     Return / reason codes returned from TKEHTP - internal error    *
	#************************************************************************/
	locMsgMap["11-1-0"] = "Invalid number of records in the Crypto module data set."
	locMsgMap["11-2-0"] = "Invalid state of the Crypto Module data set."
	locMsgMap["11-3-0"] = "Invalid message type."
	locMsgMap["11-4-0"] = "Invalid number of records in the Flag data set."
	locMsgMap["11-5-0"] = "Invalid state of the Flag data set."
	locMsgMap["11-6-0"] = "Internal crypto module number is not a whole number."
	locMsgMap["11-7-0"] = "Internal crypto module number is out of range."

	/************************************************************************
	# 13     Return / reason codes returned from TKEHTP - external errors   *
	#************************************************************************/
	locMsgMap["13-1-0"] = "Error in access to the host dataset"
	locMsgMap["13-2-0"] = "Host dataset and member name not specified"
	locMsgMap["13-3-0"] = "The member contains no records"
	locMsgMap["13-4-0"] = "The member contains more than one record"
	locMsgMap["13-5-0"] = "The record length is not 1500"
	locMsgMap["13-6-0"] = "Format error in dataset name and or member name."
	locMsgMap["13-7-0"] = "Member name already exists"
	locMsgMap["13-8-0"] = "Dataset does not exist"
	locMsgMap["13-9-0"] = "Dataset name contains a blank."
	locMsgMap["13-10-0"] = "A semicolon was specified in the host user ID or password."

	/*************************************************************************
	# 15     Return / reason codes returned from TKEHTP - TCP/IP errors     *
	#************************************************************************/
	locMsgMap["15-1--"] = "TCP/IP verb error."
	locMsgMap["15-2--"] = "TCP/IP read time out occurred on the host.\nThe reason code is the time out value in seconds."

	/************************************************************************
	# 18     Return / reason codes returned from RACF                       *
	#************************************************************************/
	locMsgMap["18-4-0"] = "Security violation. User Id invalid."
	locMsgMap["18-8-0"] = "Security violation. Password invalid."
	locMsgMap["18-12-0"] = "The password has expired. Change password from TSO."
	locMsgMap["18-14-0"] = "The user is not defined to the group."
	locMsgMap["18-28-0"] = "The user's access has been revoked."
	locMsgMap["18-34-0"] = "The user is not authorized to use the application."

	/************************************************************************
	# 23     Return / reason codes returned from OS/1 TCP/IP                *
	#************************************************************************/
	locMsgMap["23-10035-0"] = "Operation would block"
	locMsgMap["23-10036-0"] = "Operation now in progress"
	locMsgMap["23-10037-0"] = "Operation already in progress"
	locMsgMap["23-10038-0"] = "Socket operation on non-socket"
	locMsgMap["23-10039-0"] = "Destination address required"
	locMsgMap["23-10040-0"] = "Message too long"
	locMsgMap["23-10041-0"] = "Protocol wrong type for socket"
	locMsgMap["23-10042-0"] = "Protocol not available"
	locMsgMap["23-10043-0"] = "Protocol not supported"
	locMsgMap["23-10044-0"] = "Socket type not supported"
	locMsgMap["23-10045-0"] = "Operation not supported on socket"
	locMsgMap["23-10046-0"] = "Protocol family not supported"
	locMsgMap["23-10047-0"] = "Address family not supported by protocol family"
	locMsgMap["23-10048-0"] = "Address already in use"
	locMsgMap["23-10049-0"] = "Can't assign requested address"
	locMsgMap["23-10050-0"] = "Network is down"
	locMsgMap["23-10051-0"] = "Network is unreachable"
	locMsgMap["23-10052-0"] = "Network dropped connection on reset"
	locMsgMap["23-10053-0"] = "Software caused connection abort"
	locMsgMap["23-10054-0"] = "Connection reset by peer"
	locMsgMap["23-10055-0"] = "No buffer space available"
	locMsgMap["23-10056-0"] = "Socket is already connected"
	locMsgMap["23-10057-0"] = "Socket is not connected"
	locMsgMap["23-10058-0"] = "Can't send after socket shutdown"
	locMsgMap["23-10059-0"] = "Too many references: can't splice"
	locMsgMap["23-10060-0"] = "Connection timed out"
	locMsgMap["23-10061-0"] = "Connection refused"
	locMsgMap["23-10062-0"] = "Too many levels of symbolic links"
	locMsgMap["23-10063-0"] = "File name too long"
	locMsgMap["23-10064-0"] = "Host is down"
	locMsgMap["23-10065-0"] = "No route to host"
	locMsgMap["23-10066-0"] = "Directory not empty"

	/************************************************************************
	# 3     Return / reason codes returned from xCP card                    *
	#************************************************************************/
	//CKR_OK
	locMsgMap["3-0-0"] = "The verb completed processing successfully."
	//CKR_SLOT_ID_INVALID
	locMsgMap["3-3--"] = "The domain specified for the command is invalid."
	//#B@T390301CLH
	locMsgMap["3-3-102"] = "An invalid KPH certificate was supplied for Export WK."
	//#E@T390301CLH
	//CKR_ARGUMENTS_BAD
	locMsgMap["3-7--"] = "Bad arguments were provided."
	//CKR_DATA_INVALID
	locMsgMap["3-32--"] = "Invalid data was provided for the command.  The format or value of the data is incorrect."
	locMsgMap["3-32-32"] = "Administrator not found.  The administrator does not exist on the target domain or crypto module, or the supplied Subject Key Identifier (SKI) has incorrect length."
	locMsgMap["3-32-61"] = "The operation could not be performed because the Function Control Vector has not been loaded."
	locMsgMap["3-32-63"] = "The operation could not be performed because the Function Control Vector on the target crypto module is more restrictive than the Function Control Vector from the source crypto module."
	locMsgMap["3-32-103"] = "Invalid data.  Saved configuration data from the source crypto module is missing or has an invalid format."
	locMsgMap["3-32-105"] = "Invalid data.  The rewrapped transport key parts needed to apply configuration data are missing or have an invalid format."
	//CKR_DATA_LENGTH_RANGE
	locMsgMap["3-33--"] = "Command arguments were not provided, or were provided when none were expected."
	//CKR_DEVICE_MEMORY
	locMsgMap["3-49--"] = "Insufficient transient memory is available for the command to be executed."
	//CKR_FUNCTION_CANCELED
	locMsgMap["3-80--"] = "The attribute or control point settings do not allow the command to be executed."
	locMsgMap["3-80-46"] = "Change not allowed.  The control point management bits for this domain prohibit setting control points."
	locMsgMap["3-80-47"] = "Change not allowed.  The control point management bits for this domain prohibit resetting control points."
	locMsgMap["3-80-52"] = "Operation not allowed.  The permissions for this domain prohibit loading random master key values."
	locMsgMap["3-80-80"] = "Operation not allowed.  The permissions for this domain prohibit loading master keys."
//#B@T390301CLH
//	locMsgMap["3-80-81"] = "Operation not allowed.  The permissions for this domain prohibit loading a master key using a single key part."
//	locMsgMap["3-80-110"] = " Operation not allowed.  The module permissions do not allow you to collect or apply configuration data using a migration zone with M = 1."
// In spite of our objections, the EP11 team added cases for the above two
// reason codes so they no longer report a single condition.  Fall back to the
// generic error message (3-80--) for these reason codes.
	locMsgMap["3-80-131"] = "Operation not allowed.  The permissions for this domain prohibit exporting master master keys."
	locMsgMap["3-80-132"] = "Operation not allowed.  The module permissions do not allow you to collect configuration data."
	locMsgMap["3-80-133"] = "Operation not allowed.  The module permissions do not allow you to apply configuration data."
//#E@T390301CLH

	//CKR_FUNCTION_NOT_SUPPORTED
	locMsgMap["3-84--"] = "The requested function is not supported on the host crypto module."
	//CKR_KEY_HANDLE_INVALID
	locMsgMap["3-96--"] = "A requested key, certificate, or audit log record was not found."
	//CKR_KEY_CHANGED
	// intended -- no generic message for 101 return code.  PKCS #11 meaning is that the key being used now doesn't match the key being used earlier.
	locMsgMap["3-101-48"] = "Commit failed.  The supplied verification pattern of the new master key does not match the verification pattern of the key in the register.  This error will occur if the master keys being committed in a domain group are not all the same."
	//CKR_KEY_NEEDED
	locMsgMap["3-102--"] = "A key needed by the operation is missing.  Concurrent Patch Apply during configuration migration can cause this condition.  To recover, restart the Apply Configuration Data task."
	//CKR_OPERATION_NOT_INITIALIZED
	// no good generic message for 145 return code
	locMsgMap["3-145-61"] = "The operation could not be performed because the Function Control Vector has not been loaded."
	//CKR_PIN_EXPIRED
	locMsgMap["3-163--"] = "Invalid transaction counter."
	//CKR_SESSION_HANDLE_INVALID
	locMsgMap["3-179--"] = "Operation not allowed in imprint mode."
	//CKR_SIGNATURE_INVALID
	locMsgMap["3-192--"] = "One or more signatures provided with the command is invalid."
	//CKR_SIGNATURE_LEN_RANGE
	locMsgMap["3-193--"] = "A signature was malformed."
	//CKR_TEMPLATE_INCOMPLETE
	locMsgMap["3-208--"] = "Not enough signatures were provided on the command."
	//CKR_TEMPLATE_INCONSISTENT
	locMsgMap["3-209--"] = "Inconsistent template.  The data associated with a command is invalid or inconsistent."
	locMsgMap["3-209-36"] = "Invalid value.  You are not allowed to exit imprint mode with the revocation signature threshold set to 0."
	locMsgMap["3-209-56"] = "Invalid value.  You are not allowed to exit imprint mode with the revocation signature threshold set to 0."
	locMsgMap["3-209-57"] = "Invalid value.  You may not set the signature threshold to a number greater than the number of installed administrators."
	locMsgMap["3-209-58"] = "Invalid value.  You may not set the revocation signature threshold to a number greater than the number of installed administrators."
	locMsgMap["3-209-66"] = "Change not allowed.  You may not remove an administrator if that would drop the number of installed administrators below one of the threshold values."
	locMsgMap["3-209-69"] = "Change not allowed.  You are not allowed to set an attribute control bit after it has been reset."
	locMsgMap["3-209-70"] = "Change not allowed.  You are not allowed to change an attribute if the corresponding attribute control bit is reset."
	locMsgMap["3-209-71"] = "Change not allowed.  You are not allowed to change a threshold value if the corresponding attribute control bit is reset."
	//3-209-74 = Change not allowed.  You may not remove the "Threshold values of 1 allowed" permission if a threshold value is 1.
	locMsgMap["3-209-74"] = "Change not allowed.  You may not remove the Threshold values of 1 allowed permission if a threshold value is 1."
	locMsgMap["3-209-76"] = "Invalid value.  You are not allowed to set the signature threshold to 0."
	locMsgMap["3-209-77"] = "Invalid value.  You are not allowed to set the revocation signature threshold to 0."
	locMsgMap["3-209-78"] = "Invalid value.  The permissions do not allow you to set the signature threshold to 1."
	locMsgMap["3-209-79"] = "Invalid value.  The permissions do not allow you to set the revocation signature threshold to 1."
	locMsgMap["3-209-85"] = "Change not allowed.  You are not allowed to exit imprint mode at the domain level when the crypto module is still in imprint mode.  You must exit imprint mode at the crypto module level first."
	//CKR_USER_ALREADY_LOGGED_IN
	locMsgMap["3-256--"] = "An administrator could not be added because it already exists on the target."
	locMsgMap["3-256-38"] = "An administrator could not be added because it already exists on the target."
	//CKR_USER_NOT_LOGGED_IN
	locMsgMap["3-257--"] = "An unauthorized administrator signed the command."
	//CKR_USER_TOO_MANY_TYPES
	locMsgMap["3-261--"] = "An administrator could not be added because the maximum supported number of administrators already exists on the target."
	locMsgMap["3-261-39"] = "An administrator could not be added because the maximum supported number of administrators already exists on the target."
	//CKR_IBM_INTERNAL_ERROR (0x80010002)
	locMsgMap["3--2147418110--"] = "An internal error occurred in an IBM Enterprise PKCS #11 Crypto Adapter.  For assistance, contact your service representative."

	msgMap = locMsgMap
	msgInit = true
}
