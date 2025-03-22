//
// Copyright contributors to the ibm-hpcs-tke-sdk project
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/12/2021    CLH             Initial version
// 07/23/2021    CLH             Report original error when verifying OA cert chain
// 11/11/2022    CLH             T444610 - Support 4770 crypto modules

package tkesdk

import (
	"encoding/hex"
	"errors"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
	"github.com/IBM/ibm-hpcs-tke-sdk/ep11cmds"
)

/*----------------------------------------------------------------------------*/
/* Returns an array of information describing the crypto units allocated to   */
/* an HPCS service instance                                                   */
/*                                                                            */
/* Inputs:                                                                    */
/* authToken -- IBM Cloud authority token to use for requests                 */
/* urlStart -- base URL to use for requests to the IBM Cloud                  */
/* cryptoInstance -- identifies the HPCS service instance to work with        */
/*                                                                            */
/* Outputs:                                                                   */
/* []common.DomainEntry -- describes the crypto units assigned to the         */
/*     service instance                                                       */
/* error -- reports any error found during processing                         */
/*----------------------------------------------------------------------------*/
func getDomains(authToken string, urlStart string, cryptoInstance string) ([]common.DomainEntry, error) {

	// This function is based on code in tkefuncs/dlist.go.

	// Initialize the output variable
	domains := make([]common.DomainEntry, 0)

	// Determine what crypto units are assigned to the service instance
	req := common.CreateGetHsmsRequest(authToken, urlStart, cryptoInstance)
	hsm_ids, locations, serial_nums, hsm_types, err := common.SubmitQueryDomainsRequest(req)
	if err != nil {
		return domains, err
	}

	// Use maps to avoid processing a crypto module more than once
	mapLocationPublicKey := make(map[string]string)
	mapLocationOACert := make(map[string][]byte)
	mapLocationSerialNum := make(map[string]string)
	// The key for these maps is the first part of the location string
	// (everything except the domain index at the end).  Two locations
	// can map to the same crypto module (two paths to the same place),
	// but this is the best we can do.

	// For each crypto unit, read the OA certificate
	for i := 0; i < len(hsm_ids); i++ {
		partialLocation := common.GetPartialLocation(locations[i])
		_, found := mapLocationSerialNum[partialLocation]
		if !found {
			// Create a DomainEntry for querying the crypto module
			de := common.DomainEntry{
				0, // Domain_num -- don't care
				hsm_ids[i],
				cryptoInstance,
				locations[i],
				"",              // Serial_num -- don't care
				"not available", // Public_key -- not available for initial read
				hsm_types[i],
				false} // Selected -- don't care
			// Read the epoch OA certificate with the current OA
			// signature key
			certbytes, err := ep11cmds.QueryDeviceCertificate(
				authToken, urlStart, de, 0)
			if err != nil {
				return domains, err
			}
			mapLocationOACert[partialLocation] = certbytes
			//#B@T444610CLH
			if len(certbytes) == ep11cmds.CEX8_OA_CERTIFICATE_LENGTH ||
			   len(certbytes) == ep11cmds.CEX8_MB_CERTIFICATE_LENGTH {
				// Handle OA certificate for the CEX8P
				var cert ep11cmds.OA3CertificateX
				err = cert.Init(certbytes)
				if err != nil {
					return domains, err
				}
				mapLocationPublicKey[partialLocation] =
					hex.EncodeToString(cert.SpkiPublicKey)
			} else if certbytes[0] == 0x45 {
			//#E@T444610CLH	
				// Handle OA certificate for the CEX6P or CEX7P
				var cert ep11cmds.OA2CertificateX
				err = cert.Init(certbytes)
				if err != nil {
					return domains, err
				}
				mapLocationPublicKey[partialLocation] =
					hex.EncodeToString(cert.SpkiPublicKey)
			} else {
				// Handle OA certificate for the CEX5P
				var cert ep11cmds.OACertificateX
				err = cert.Init(certbytes)
				if err != nil {
					return domains, err
				}
				mapLocationPublicKey[partialLocation] =
					hex.EncodeToString(cert.PublicKey)
			}
			// Read the actual serial number for the crypto module
			de.Public_key = mapLocationPublicKey[partialLocation]
			// We want to check the OA signature on this query
			_, resp, err := ep11cmds.QueryDomainAttributes(authToken, urlStart, de)
			if err != nil {
				return domains, err
			}
			mapLocationSerialNum[partialLocation] = resp.GetSerialNumber()
		}
	}

	// Check that actual and reported serial numbers match
	for i := 0; i < len(hsm_ids); i++ {
		partialLocation := common.GetPartialLocation(locations[i])
		if serial_nums[i] !=
			mapLocationSerialNum[partialLocation] {
			return domains, errors.New("Serial number mismatch detected")
		}
	}

	// Verify the OA certificate chain for all crypto modules accessed

	// Create map of serial_num --> cert_chain_checked
	certChainChecked := make(map[string]bool)
	// The map allows us to avoid processing a crypto module more than
	// once.  True means the OA certificate chain has already been
	// verified.

	for i := 0; i < len(hsm_ids); i++ {
		partialLocation := common.GetPartialLocation(locations[i])
		serialNum := mapLocationSerialNum[partialLocation]
		newkey := mapLocationPublicKey[partialLocation]

		// Check if we have already tried to verify the OA certificate
		// chain in this loop
		if !certChainChecked[serialNum] {
			// Create a DomainEntry for verifying the OA certificate chain
			de := common.DomainEntry{
				0,              // Domain_num -- don't care
				hsm_ids[i],     // accurate, but don't care
				cryptoInstance, // accurate, but don't care
				locations[i],
				serialNum, // accurate, but don't care
				newkey,    // Public_key
				hsm_types[i],
				false} // Selected -- don't care
			certbytes := mapLocationOACert[partialLocation]
			//#B@T444610CLH
			if len(certbytes) == ep11cmds.CEX8_OA_CERTIFICATE_LENGTH ||
			   len(certbytes) == ep11cmds.CEX8_MB_CERTIFICATE_LENGTH {
				// Handle OA certificate for the CEX8P
				var cert ep11cmds.OA3CertificateX
				err = cert.Init(certbytes)
				if err != nil {
					return domains, err
				}
				err = ep11cmds.VerifyOA3Certificate(authToken, urlStart, de, 0, cert)
				if err != nil {
					return domains, err
				} else {
					certChainChecked[serialNum] = true
				}
			} else if certbytes[0] == 0x45 {
			//#E@T444610CLH
				// Handle OA certificate for the CEX6P or CEX7P
				var cert ep11cmds.OA2CertificateX
				err = cert.Init(certbytes)
				if err != nil {
					return domains, err
				}
				err = ep11cmds.VerifyOA2Certificate(authToken, urlStart, de, 0, cert)
				if err != nil {
					return domains, err
				} else {
					certChainChecked[serialNum] = true
				}
			} else {
				// Handle OA certificate for the CEX5P
				var cert ep11cmds.OACertificateX
				err = cert.Init(certbytes)
				if err != nil {
					return domains, err
				}
				err = ep11cmds.VerifyCertificate(authToken, urlStart, de, 0, cert)
				if err != nil {
					return domains, err
				} else {
					certChainChecked[serialNum] = true
				}
			}
		}
	}

	// Assemble the information to be returned
	domainNum := 1
	for i := 0; i < len(hsm_ids); i++ {
		partialLocation := common.GetPartialLocation(locations[i])
		serialNum := mapLocationSerialNum[partialLocation]
		domains = append(
			domains,
			common.DomainEntry{
				domainNum,
				hsm_ids[i],
				cryptoInstance,
				locations[i],
				serialNum,
				mapLocationPublicKey[partialLocation],
				hsm_types[i],
				true})
		domainNum++
	}

	return domains, nil
}
