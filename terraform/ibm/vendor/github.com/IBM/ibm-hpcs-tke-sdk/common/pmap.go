//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 12/08/2020    CLH             T390301 - Add minimal touch functions

package common

import (
	"encoding/binary"
	"encoding/hex"
)

/*----------------------------------------------------------------------------*/
/* Type for working with ASN.1 sequences of the form defined in section 5.3   */
/* ("Serialized module state") of the EP11 wire formats document.             */
/*                                                                            */
/* Export WK and Export pending WK use ASN.1 sequences of this form for their */
/* input and output parameters.                                               */
/*----------------------------------------------------------------------------*/
type ParameterMap struct {
	pMap map[string][]byte
}

/*----------------------------------------------------------------------------*/
/* Tag values used in parameter maps                                          */
/*----------------------------------------------------------------------------*/

/** Domain administrator subject key identifiers (SKIs) */
var PMTAG_DOMAIN_ADMIN_SKIS string = "0x000A"

/** Domain administrator certificates */
var PMTAG_DOMAIN_ADMIN_CERTS string = "0x000B"

/** Domain level information */
var PMTAG_DOMAIN_QUERY_INFO string = "0x000C"

/** Domain attributes */
var PMTAG_DOMAIN_ATTRIBUTES string = "0x000F"

/** Domain transaction counter */
var PMTAG_DOMAIN_TRANSACTION_COUNTER string = "0x0011"

/** OA certificate */
var PMTAG_OA_CERTIFICATE string = "0x0015"

/** Domain control points */
var PMTAG_DOMAIN_CONTROL_POINTS string = "0x0017"

/** ASN.1 structure containing an encrypted key part */
var PMTAG_ENCR_KEY_PART string = "0x0019"

/** ASN.1 structure containing an OA signature over encrypted key part */
var PMTAG_SIGNATURE_ENCR_KEY_PART string = "0x001A"

/** M policy (number of key parts required to reconstruct the key) */
var PMTAG_M_POLICY string = "0x001C"

/** KPH certificate containing a public key */
var PMTAG_KPH_CERTIFICATE string = "0x001D"

/** Scope restrictions on a state export request */
var PMTAG_STATE_SCOPE = "0x001F"

// See section 5.3 of the EP11 wire formats document for other defined tag
// values.

/*----------------------------------------------------------------------------*/
/* Creates a new parameter map and initializes it to empty.                   */
/*----------------------------------------------------------------------------*/
func NewParameterMap() ParameterMap {
	var pm ParameterMap
	pm.pMap = make(map[string][]byte, 0)
	return pm
}

/*----------------------------------------------------------------------------*/
/* Initializes a parameter map using an input ASN.1 sequence of the form      */
/* described in section 5.3 of the EP11 wire formats document.                */
/*                                                                            */
/* Input:                                                                     */
/* []byte -- input ASN.1 sequence                                             */
/*                                                                            */
/* Outputs:                                                                   */
/* ParameterMap -- the updated parameter map                                  */
/* error -- reports invalid ASN.1 input sequence                              */
/*----------------------------------------------------------------------------*/
func (pm ParameterMap) Load(data []byte) (ParameterMap, error) {

	// Make sure the map is initialized and discard any prior contents
	pm.pMap = make(map[string][]byte, 0)

	seqbytes, err := Asn1GetSequenceBytes(data, 0)
	if err != nil {
		return pm, err
	}
	next := 0
	for next < len(seqbytes) {
		octbytes, err := Asn1GetOctetStringBytes(seqbytes, next)
		if err != nil {
			return pm, err
		}
		next, err = Asn1SkipOctetString(seqbytes, next)
		if err != nil {
			return pm, err
		}

		tag := "0x" + hex.EncodeToString(octbytes[0:2])
		auxInt := hex.EncodeToString(octbytes[2:6])
		value := octbytes[2:]

		// For entries that can repeat, add the index to the tag.
		switch tag {
		// Entries that can repeat
		case PMTAG_DOMAIN_ADMIN_SKIS,
			PMTAG_DOMAIN_ADMIN_CERTS,
			PMTAG_DOMAIN_QUERY_INFO,
			PMTAG_DOMAIN_ATTRIBUTES,
			PMTAG_DOMAIN_TRANSACTION_COUNTER,
			PMTAG_OA_CERTIFICATE,
			PMTAG_DOMAIN_CONTROL_POINTS,
			PMTAG_ENCR_KEY_PART,
			PMTAG_SIGNATURE_ENCR_KEY_PART,
			PMTAG_KPH_CERTIFICATE:

			pm.pMap[tag+"+"+auxInt] = value

		// Entries that cannot repeat
		default:
			pm.pMap[tag] = value
		}
	}
	return pm, nil
}

/*----------------------------------------------------------------------------*/
/* Returns the auxiuliary integer associated with a parameter entry.          */
/*                                                                            */
/* Input:                                                                     */
/* string -- tag identifying the parameter to retrieve                        */
/*                                                                            */
/* Output:                                                                    */
/* uint32 -- integer value associated with the parameter                      */
/*----------------------------------------------------------------------------*/
func (pm ParameterMap) GetAuxInt(tag string) uint32 {
	data := pm.pMap[tag]
	if data == nil {
		panic("Map entry not found")
	}
	return binary.BigEndian.Uint32(data[0:4])
}

/*----------------------------------------------------------------------------*/
/* Returns data from a parameter map when an index value is used.             */
/*                                                                            */
/* Inputs:                                                                    */
/* string -- tag identifying the parameter to retrieve                        */
/* uint32 -- index value to combine with the tag                              */
/*                                                                            */
/* Output:                                                                    */
/* []byte -- the parameter from the map, nil if no map entry exists           */
/*----------------------------------------------------------------------------*/
func (pm ParameterMap) GetDataUsingIndex(tag string, index uint32) []byte {

	auxInt := Uint32To4ByteSlice(index)
	data := pm.pMap[tag+"+"+hex.EncodeToString(auxInt)]

	if data == nil {
		panic("Map entry not found")
	}
	return data[4:]
}

/*----------------------------------------------------------------------------*/
/* Adds a value to a parameter map.                                           */
/*                                                                            */
/* Inputs:                                                                    */
/* string -- tag identifying the parameter to add                             */
/* uint32 -- index or associated integer                                      */
/* []byte -- additional data                                                  */
/*----------------------------------------------------------------------------*/
func (pm ParameterMap) Put(tag string, index uint32, data []byte) {

	auxInt := Uint32To4ByteSlice(index)

	value := make([]byte, 0)
	value = append(value, auxInt...)
	value = append(value, data...)

	// For entries that can repeat, add the index to the tag.
	switch tag {
	// Entries that can repeat
	case PMTAG_DOMAIN_ADMIN_SKIS,
		PMTAG_DOMAIN_ADMIN_CERTS,
		PMTAG_DOMAIN_QUERY_INFO,
		PMTAG_DOMAIN_ATTRIBUTES,
		PMTAG_DOMAIN_TRANSACTION_COUNTER,
		PMTAG_OA_CERTIFICATE,
		PMTAG_DOMAIN_CONTROL_POINTS,
		PMTAG_ENCR_KEY_PART,
		PMTAG_SIGNATURE_ENCR_KEY_PART,
		PMTAG_KPH_CERTIFICATE:

		pm.pMap[tag+"+"+hex.EncodeToString(auxInt)] = value

	// Entries that cannot repeat
	default:
		pm.pMap[tag] = value
	}
}

/*----------------------------------------------------------------------------*/
/* Returns an ASN.1 sequence of octet strings for the parameters in the map.  */
/*                                                                            */
/* Output:                                                                    */
/* []byte -- ASN.1 sequence                                                   */
/*----------------------------------------------------------------------------*/
func (pm ParameterMap) GenerateBytes() []byte {
	elements := make([][]byte, 0)
	for key, value := range pm.pMap {
		tag, err := hex.DecodeString(key[2:6])
		if err != nil {
			panic("Invalid tag string in parameter map")
		}
		data := make([]byte, 0)
		data = append(data, tag...)
		data = append(data, value...)
		elements = append(elements, Asn1FormOctetString(data))
	}
	return Asn1FormSequence(elements)
}
