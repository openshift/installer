//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 07/04/2021    CLH             Adapt for TKE SDK

package common

import (
	"errors"
	"strconv"
	"strings"

	"github.com/IBM/ibm-hpcs-tke-sdk/rest"
)

type Location struct {
	AvailZone   string
	HostSystem  string
	CMIndex     int
	DomainIndex int
}

/*----------------------------------------------------------------------------*/
/* Submits the POST /hsms request that sends an HTPRequest to a TKE catcher   */
/* program.                                                                   */
/*                                                                            */
/* Returns the HTPResponse string from the TKE catcher program.               */
/*----------------------------------------------------------------------------*/
func SubmitHTPRequest(req *rest.Request) (htpResponse string, err error) {

	var outmap map[string]string

	// Create client with default TLS handshake timeout of 10 seconds
	client := rest.NewClient()

	// Create client with no TLS handshake timeout
	//	tr := &http.Transport{
	//		TLSHandshakeTimeout: 0,
	// for 30 second timeout use 30 * 1000 * 1000 * 1000
	//	}
	//	dfcl := &http.Client{
	//		Transport: tr,
	//	}
	//	client := &rest.Client{
	//		HTTPClient: dfcl,
	//		DefaultHeader: make(http.Header),
	//	}

	_, err = client.Do(req, &outmap, nil)
	if err != nil {
		t1, ok := err.(*rest.ErrorResponse)
		if ok {
			return "", errors.New(
				"Error sending HTPRequest to target service instance." +
					"\nStatus code: " + strconv.Itoa(t1.StatusCode) +
					"\n" + getMessageText(t1))
		} else if strings.Contains(err.Error(), "no such host") {
			return "", errors.New(
				"Error sending HTPRequest to target service instance." +
					"\nMessage: " + err.Error() +
					"\n\nHyper Protect Crypto Services may not be available in " +
					"the region you have selected." +
					"\nCheck the REFERENCE section of the Hyper Protect Crypto " +
					"Services online documentation to determine the regions and " +
					"locations where the service is available.")
		} else {
			return "", errors.New(
				"Error sending HTPRequest to target service instance." +
					"\nMessage: " + err.Error())
		}
	}
	resp := outmap["response"]
	if resp == "" {
		return "", errors.New(
			"Error sending HTPRequest to target service instance." +
				"\nNo HPTResponse returned.")
	}
	return resp, nil
}

/*----------------------------------------------------------------------------*/
/* Submits the GET /hsms request that queries the Cloud for the domains       */
/* associated with a crypto instance.                                         */
/*                                                                            */
/* Input:                                                                     */
/* *rest.Request -- the GET /hsms request to be sent to the cloud.            */
/*                                                                            */
/* Outputs:                                                                   */
/* []string -- hsm_ids of each domain in the crypto instance                  */
/* []string -- locations of the crypto modules for each domain                */
/* []string -- serial numbers of the crypto modules for each domain           */
/* []string -- hsm_types, "recovery" or "operational"                         */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func SubmitQueryDomainsRequest(req *rest.Request) ([]string, []string,
	[]string, []string, error) {

	/*
	 * The format of the response for a GET /hsms request is:
	 *
	 * {
	 *    "metadata": [
	 *       {
	 *          "collectionType": "string",
	 *          "collectionTotal": int
	 *       }
	 *    ]
	 *    "hsms": [
	 *       {
	 *          "hsm_id": "string",
	 *          "location": "string",
	 *          "state": "string",
	 *          "version": int,
	 *          "serial_number: "string"
	 *       }
	 *    ]
	 *    "source_hsms": [
	 *       {
	 *          "hsm_id": "string",
	 *          "location": "string",
	 *          "state": "string",
	 *          "version": int,
	 *          "serial_number: "string"
	 *       }
	 *    ]
	 * }
	 *
	 * The original version of the GET /hsms request did not include a source_hsms
	 * array in the response.  To allow a smooth transition, continue to accept
	 * responses that lack a source_hsms array.
	 */

	var outmap = make(map[string]interface{})

	// Create client with default TLS handshake timeout of 10 seconds
	client := rest.NewClient()

	// Create client with no TLS handshake timeout
	//	tr := &http.Transport{
	//		TLSHandshakeTimeout: 0,
	// for 30 second timeout use 30 * 1000 * 1000 * 1000
	//	}
	//	dfcl := &http.Client{
	//		Transport: tr,
	//	}
	//	client := &rest.Client{
	//		HTTPClient: dfcl,
	//		DefaultHeader: make(http.Header),
	//	}

	_, err := client.Do(req, &outmap, nil)
	if err != nil {
		t1, ok := err.(*rest.ErrorResponse)
		if ok {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nStatus code: " + strconv.Itoa(t1.StatusCode) +
					"\n" + getMessageText(t1))
		} else {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nMessage: " + err.Error())
		}
	}

	hsm_ids := make([]string, 0)
	locations := make([]string, 0)
	serial_nums := make([]string, 0)
	hsm_types := make([]string, 0)

	// Process operational HSMs
	hsms := outmap["hsms"]
	hsmsArray, ok := hsms.([]interface{})
	if !ok {
		return nil, nil, nil, nil, errors.New(
			"Error querying crypto units." +
				"\nUnexpected hsms data, not an array.")
	}

	for _, hsm := range hsmsArray {
		hsmEntry, ok := hsm.(map[string]interface{})
		if !ok {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nUnexpected hsms entry.")
		}

		hsm_id := hsmEntry["hsm_id"]
		if hsm_id == nil {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nhsm_id not found")
		}
		hsm_id_string, ok := hsm_id.(string)
		if !ok {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nhsm_id is not a string")
		}
		location := hsmEntry["location"]
		if location == nil {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nlocation not found")
		}
		location_string, ok := location.(string)
		if !ok {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nlocation is not a string")
		}
		serial_num := hsmEntry["serial_number"]
		if serial_num == nil {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nserial_number not found")
		}
		serial_num_string, ok := serial_num.(string)
		if !ok {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nserial_number is not a string")
		}

		hsm_ids = append(hsm_ids, hsm_id_string)
		locations = append(locations, location_string)
		serial_nums = append(serial_nums, serial_num_string)
		hsm_types = append(hsm_types, "operational")
	}

	// Process recovery HSMs
	// This is a little confusing.  Originally "recovery HSMs" were called
	// "source HSMs".  The query returns a "source_hsms" array containing
	// the recovery HSMs for the service instance.  In this code continue
	// to use the term "source HSM".
	source_hsms := outmap["source_hsms"]
	if source_hsms != nil {
		sourceHsmsArray, ok := source_hsms.([]interface{})
		if !ok {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nUnexpected source_hsms data, not an array.")
		}

		for _, hsm := range sourceHsmsArray {
			hsmEntry, ok := hsm.(map[string]interface{})
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nUnexpected source_hsms entry.")
			}

			hsm_id := hsmEntry["hsm_id"]
			if hsm_id == nil {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nhsm_id not found")
			}
			hsm_id_string, ok := hsm_id.(string)
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nhsm_id is not a string")
			}
			location := hsmEntry["location"]
			if location == nil {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nlocation not found")
			}
			location_string, ok := location.(string)
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nlocation is not a string")
			}
			serial_num := hsmEntry["serial_number"]
			if serial_num == nil {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nserial_number not found")
			}
			serial_num_string, ok := serial_num.(string)
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nserial_number is not a string")
			}

			hsm_ids = append(hsm_ids, hsm_id_string)
			locations = append(locations, location_string)
			serial_nums = append(serial_nums, serial_num_string)
			hsm_types = append(hsm_types, "recovery")
		}
	}

	// Process failover HSMs
	failover_hsms := outmap["failover_hsms"]
	if failover_hsms != nil {
		sourceHsmsArray, ok := failover_hsms.([]interface{})
		if !ok {
			return nil, nil, nil, nil, errors.New(
				"Error querying crypto units." +
					"\nUnexpected failover_hsms data, not an array.")
		}

		for _, hsm := range sourceHsmsArray {
			hsmEntry, ok := hsm.(map[string]interface{})
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nUnexpected failover_hsms entry.")
			}

			hsm_id := hsmEntry["hsm_id"]
			if hsm_id == nil {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nhsm_id not found")
			}
			hsm_id_string, ok := hsm_id.(string)
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nhsm_id is not a string")
			}
			location := hsmEntry["location"]
			if location == nil {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nlocation not found")
			}
			location_string, ok := location.(string)
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nlocation is not a string")
			}
			serial_num := hsmEntry["serial_number"]
			if serial_num == nil {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nserial_number not found")
			}
			serial_num_string, ok := serial_num.(string)
			if !ok {
				return nil, nil, nil, nil, errors.New(
					"Error querying crypto units." +
						"\nserial_number is not a string")
			}

			hsm_ids = append(hsm_ids, hsm_id_string)
			locations = append(locations, location_string)
			serial_nums = append(serial_nums, serial_num_string)
			hsm_types = append(hsm_types, "failover")
		}
	}

	return hsm_ids, locations, serial_nums, hsm_types, nil
}

/*----------------------------------------------------------------------------*/
/* Determines the message text to display for an error.                       */
/*----------------------------------------------------------------------------*/
func getMessageText(rsp *rest.ErrorResponse) string {
	message := "Message: " + rsp.Message
	if rsp.StatusCode == 401 {
		message = message +
			"\nYour access token is invalid, expired, or does not have the " +
			"necessary permissions to access this instance. You may need to log " +
			"out and log back in to the IBM Cloud CLI to refresh your token. If " +
			"the problem persists, note the input and output of this attempt and " +
			"contact IBM Cloud support."
	} else if rsp.StatusCode == 404 {
		message = message +
			"\nThe service instance or crypto unit could not be found. Check and " +
			"ensure that you have the correct resource group and region targeted " +
			"and that you are logged into the correct account. If the problem " +
			"persists, note the input and output of this attempt and contact " +
			"IBM Cloud support."
	} else if rsp.StatusCode == 500 {
		message = message +
			"\nIBM Hyper Protect Crypto Services is currently unavailable.  " +
			"Your request could not be processed.  Please try again " +
			"later.  If the problem persists, note the input and output " +
			"of this attempt and contact IBM Cloud support."
	}
	// Additional message text is not needed for the 400 status code.
	return message
}

func ParseLocation(location string) (result *Location, err error) {
	pieces := strings.Split(location, ".")
	if len(pieces) != 4 {
		return nil, errors.New(`String "` + location + `" does not contain 4 pieces`)
	}
	// remove the leading and trailing bracket from each field
	for i := range pieces {
		pieces[i] = strings.TrimSpace(pieces[i])
		pieceLen := len(pieces[i])
		if len(pieces[i]) >= 2 {
			if pieces[i][pieceLen-1] == ']' {
				pieces[i] = pieces[i][:pieceLen-1]
				pieceLen--
			}
			if pieces[i][0] == '[' {
				pieces[i] = pieces[i][1:pieceLen]
			}
		}
	}
	returnStruct := new(Location)
	returnStruct.AvailZone = pieces[0]
	returnStruct.HostSystem = pieces[1]

	returnStruct.CMIndex, err = convertStringToIndex(pieces[2])
	if err != nil {
		return nil, err
	}
	returnStruct.DomainIndex, err = convertStringToIndex(pieces[3])
	if err != nil {
		return nil, err
	}
	return returnStruct, nil
}

func convertStringToIndex(theString string) (int, error) {
	indexValue, err := strconv.Atoi(theString)
	if err != nil {
		return -1, err
	}
	if indexValue < 0 {
		return indexValue, errors.New("CM index value must be a positive number")
	}
	return indexValue, nil
}

/*----------------------------------------------------------------------------*/
/* Submits a GET /keys request to a signing service to retrieve the public    */
/* part of a signature key.                                                   */
/*                                                                            */
/* Input:                                                                     */
/* *rest.Request -- the GET /keys request to be sent to the signing service   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the base64 encoded public key                                    */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func SubmitQueryPublicKeyRequest(req *rest.Request) (string, error) {

	/*
	 * The format of the response for a GET /hsms request is:
	 *
	 * {
	 *     "publickey": <base64 encoded string of public key (ASN.1 DER encoded struct containing X and Y)>
	 * }
	 *
	 */

	var outmap = make(map[string]interface{})

	// Create client with default TLS handshake timeout of 10 seconds
	client := rest.NewClient()

	// Create client with no TLS handshake timeout
	//	tr := &http.Transport{
	//		TLSHandshakeTimeout: 0,
	// for 30 second timeout use 30 * 1000 * 1000 * 1000
	//	}
	//	dfcl := &http.Client{
	//		Transport: tr,
	//	}
	//	client := &rest.Client{
	//		HTTPClient: dfcl,
	//		DefaultHeader: make(http.Header),
	//	}

	_, err := client.Do(req, &outmap, nil)
	if err != nil {
		t1, ok := err.(*rest.ErrorResponse)
		if ok {
			return "", errors.New(
				"Error requesting public part of signature key." +
					"\nStatus code: " + strconv.Itoa(t1.StatusCode) +
					"\nMessage: " + t1.Message)
		} else {
			return "", errors.New(
				"Error requesting public part of signature key." +
					"\nMessage: " + err.Error())
		}
	}

	pubkey := outmap["publickey"]
	if pubkey == nil {
		return "", errors.New(
			"Error requesting public part of signature key." +
				"\npublickey not found")
	}
	pubkey_string, ok := pubkey.(string)
	if !ok {
		return "", errors.New(
			"Error requesting public part of signature key." +
				"\npublickey is not a string")
	}

	return pubkey_string, nil
}

/*----------------------------------------------------------------------------*/
/* Submits a POST /sign request to a signing service to sign the supplied     */
/* data.                                                                      */
/*                                                                            */
/* Input:                                                                     */
/* *rest.Request -- the POST /sign request to be sent to the signing service  */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the base64 encoded signature                                     */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func SubmitSignDataRequest(req *rest.Request) (string, error) {

	/*
	 * The format of the response for a POST /sign request is:
	 *
	 * {
	 *     "signature": "<base64 encoded string of binary data (signature: ASN.1 DER encoded struct of integers R and S)>"
	 * }
	 *
	 */

	var outmap = make(map[string]interface{})

	// Create client with default TLS handshake timeout of 10 seconds
	client := rest.NewClient()

	// Create client with no TLS handshake timeout
	//	tr := &http.Transport{
	//		TLSHandshakeTimeout: 0,
	// for 30 second timeout use 30 * 1000 * 1000 * 1000
	//	}
	//	dfcl := &http.Client{
	//		Transport: tr,
	//	}
	//	client := &rest.Client{
	//		HTTPClient: dfcl,
	//		DefaultHeader: make(http.Header),
	//	}

	_, err := client.Do(req, &outmap, nil)
	if err != nil {
		t1, ok := err.(*rest.ErrorResponse)
		if ok {
			return "", errors.New(
				"Error requesting a signature over supplied data." +
					"\nStatus code: " + strconv.Itoa(t1.StatusCode) +
					"\nMessage: " + t1.Message)
		} else {
			return "", errors.New(
				"Error requesting a signature over supplied data." +
					"\nMessage: " + err.Error())
		}
	}

	signature := outmap["signature"]
	if signature == nil {
		return "", errors.New(
			"Error requesting a signature over supplied data." +
				"\nsignature not found")
	}
	signature_string, ok := signature.(string)
	if !ok {
		return "", errors.New(
			"Error requesting a signature over supplied data." +
				"\nsignature is not a string")
	}

	return signature_string, nil
}

