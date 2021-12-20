//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/27/2020    CLH             T372621 - Support P521 EC signature keys
// 12/08/2020    CLH             T390301 - Add minimal touch functions

package common

import (
	"errors"
)

/* Converted to Go from com.ibm.tke.util.ASN1.java */

/** ASN.1 tag for INTEGER */
var ASN1_INTEGER_TAG byte = 0x02

/** ASN.1 tag for BIT STRING */
var ASN1_BIT_STRING_TAG byte = 0x03

/** ASN.1 tag for OCTET STRING */
var ASN1_OCTET_STRING_TAG byte = 0x04

/** ASN.1 tag for an object identifier (OID) */
var ASN1_OID_TAG byte = 0x06

/** ASN.1 tag for SEQUENCE */
var ASN1_SEQUENCE_TAG byte = 0x30

/** ASN.1 tag for context-specific entry */
var ASN1_CONTEXT_SPECIFIC_TAG byte = 0x80

/*----------------------------------------------------------------------------*/
/* Returns a DER-encoded length for a non-negative four-byte integer.         */
/*----------------------------------------------------------------------------*/
func Asn1EncodeLength(length int) []byte {
	if length < 0 {
		panic("Negative length cannot be DER encoded.")
	}
	result := make([]byte, 0)
	if length < 128 {
		result = append(result, byte(length))
	} else if length < 256 {
		result = append(result, 0x81)
		result = append(result, byte(length))
	} else if length < 256*256 {
		result = append(result, 0x82)
		result = append(result, byte(length/256))
		result = append(result, byte(length%256))
	} else if length < 256*256*256 {
		result = append(result, 0x83)
		result = append(result, byte(length/(256*256)))
		result = append(result, byte((length/256)%256))
		result = append(result, byte(length%256))
	} else {
		result = append(result, 0x84)
		result = append(result, byte(length/(256*256*256)))
		result = append(result, byte((length/(256*256))%256))
		result = append(result, byte((length/256)%256))
		result = append(result, byte(length%256))
	}
	return result
}

/*----------------------------------------------------------------------------*/
/* Forms an ASN.1 OCTET STRING containing the input byte stream.              */
/*----------------------------------------------------------------------------*/
func Asn1FormOctetString(source []byte) []byte {
	result := make([]byte, 0)
	result = append(result, ASN1_OCTET_STRING_TAG)
	result = append(result, Asn1EncodeLength(len(source))...)
	result = append(result, source...)
	return result
}

//#B@T372621CLH
/*----------------------------------------------------------------------------*/
/* Forms an ASN.1 BIT STRING containing the input byte stream.  Adds a 0x00   */
/* byte to the beginning of the source data.                                  */
/*----------------------------------------------------------------------------*/
func Asn1FormBitString(source []byte) []byte {
	newSource := make([]byte, 1)
	newSource = append(newSource, source...)
	result := make([]byte, 0)
	result = append(result, ASN1_BIT_STRING_TAG)
	result = append(result, Asn1EncodeLength(len(newSource))...)
	result = append(result, newSource...)
	return result
}

//#E@T372621CLH

/*----------------------------------------------------------------------------*/
/* Forms an ASN.1 SEQUENCE containing the set of input elements.              */
/*                                                                            */
/* The input argument is an array of ASN.1 encoded elements that are          */
/* concatenated together to form the final SEQUENCE.                          */
/*----------------------------------------------------------------------------*/
func Asn1FormSequence(elements [][]byte) []byte {
	seqlen := 0
	for i := range elements {
		seqlen += len(elements[i])
	}
	result := make([]byte, 0)
	result = append(result, ASN1_SEQUENCE_TAG)
	result = append(result, Asn1EncodeLength(seqlen)...)
	for i := range elements {
		result = append(result, elements[i]...)
	}
	return result
}

/*----------------------------------------------------------------------------*/
/* Interprets the byte stream at the specified offset as a BER-encoded        */
/* length, and returns the length value.                                      */
/*----------------------------------------------------------------------------*/
func Asn1GetLength(source []byte, offset int) (int, error) {
	if (source[offset] & 0x80) == 0x00 {
		return int(source[offset]), nil
	}
	size := int(source[offset] & 0x7F)
	if (size == 0) || (size > 4) {
		return 0, errors.New("Invalid length of ASN.1 length field")
	}
	length := 0
	for i := 0; i < size; i++ {
		length = (length * 256) + int(source[offset+1+i])
	}
	if length < 0 {
		return 0, errors.New("Overflow in ASN.1 length")
	}
	return length, nil
}

/*----------------------------------------------------------------------------*/
/* Returns an updated offset into a source byte stream when the stream is     */
/* interpreted as a BER-encoded length and we want to skip to the next field  */
/* in the stream.                                                             */
/*----------------------------------------------------------------------------*/
func Asn1SkipLength(source []byte, offset int) (int, error) {
	if (source[offset] & 0x80) == 0x00 {
		return offset + 1, nil
	}
	size := int(source[offset] & 0x7F)
	if (size == 0) || (size > 4) {
		return 0, errors.New("Invalid length of ASN.1 length field")
	}
	return offset + 1 + size, nil
}

/*----------------------------------------------------------------------------*/
/* Interprets the byte stream at the specified offset as a SEQUENCE, and      */
/* returns the payload bytes without the SEQUENCE tag and length.             */
/*----------------------------------------------------------------------------*/
func Asn1GetSequenceBytes(source []byte, offset int) ([]byte, error) {
	if source[offset] != ASN1_SEQUENCE_TAG {
		return nil, errors.New("Expected SEQUENCE tag not found")
	}
	length, err := Asn1GetLength(source, offset+1)
	if err != nil {
		return nil, err
	}
	newOffset, err := Asn1SkipLength(source, offset+1)
	if err != nil {
		return nil, err
	}
	return source[newOffset : newOffset+length], nil
}

/*----------------------------------------------------------------------------*/
/* Returns an updated offset into a source byte stream when the stream is     */
/* interpreted as a SEQUENCE and we want to skip to the next field in the     */
/* stream.                                                                    */
/*----------------------------------------------------------------------------*/
func Asn1SkipSequence(source []byte, offset int) (int, error) {
	sequence, err := Asn1GetSequenceBytes(source, offset)
	if err != nil {
		return 0, err
	}
	newOffset, err := Asn1SkipLength(source, offset+1)
	if err != nil {
		return 0, err
	}
	return newOffset + len(sequence), nil
}

/*----------------------------------------------------------------------------*/
/* Interprets the byte stream at the specified offset as an OCTET STRING,     */
/* and returns the payload bytes without the OCTET STRING tag and length.     */
/*----------------------------------------------------------------------------*/
func Asn1GetOctetStringBytes(source []byte, offset int) ([]byte, error) {
	if source[offset] != ASN1_OCTET_STRING_TAG {
		return nil, errors.New("Expected OCTET STRING tag not found")
	}
	length, err := Asn1GetLength(source, offset+1)
	if err != nil {
		return nil, err
	}
	newOffset, err := Asn1SkipLength(source, offset+1)
	if err != nil {
		return nil, err
	}
	return source[newOffset : newOffset+length], nil
}

/*----------------------------------------------------------------------------*/
/* Interprets the byte stream at the specified offset as an INTEGER, and      */
/* returns the payload bytes without the INTEGER tag and length.              */
/*----------------------------------------------------------------------------*/
func Asn1GetIntegerBytes(source []byte, offset int) ([]byte, error) {
	if source[offset] != ASN1_INTEGER_TAG {
		return nil, errors.New("Expected INTEGER tag not found")
	}
	length, err := Asn1GetLength(source, offset+1)
	if err != nil {
		return nil, err
	}
	newOffset, err := Asn1SkipLength(source, offset+1)
	if err != nil {
		return nil, err
	}
	return source[newOffset : newOffset+length], nil
}

/*----------------------------------------------------------------------------*/
/* Returns an updated offset into a source byte stream when the stream is     */
/* interpreted as an INTEGER and we want to skip to the next field in the     */
/* stream.                                                                    */
/*----------------------------------------------------------------------------*/
func Asn1SkipInteger(source []byte, offset int) (int, error) {
	intBytes, err := Asn1GetIntegerBytes(source, offset)
	if err != nil {
		return 0, err
	}
	newOffset, err := Asn1SkipLength(source, offset+1)
	if err != nil {
		return 0, err
	}
	return newOffset + len(intBytes), nil
}

//#B@T390301CLH
/*----------------------------------------------------------------------------*/
/* Returns an updated offset into a source byte stream when the stream is     */
/* interpreted as an OCTET STRING and we want to skip to the next field in    */
/* the stream.                                                                */
/*----------------------------------------------------------------------------*/
func Asn1SkipOctetString(source []byte, offset int) (int, error) {
	octetStringBytes, err := Asn1GetOctetStringBytes(source, offset)
	if err != nil {
		return 0, err
	}
	newOffset, err := Asn1SkipLength(source, offset+1)
	if err != nil {
		return 0, err
	}
	return newOffset + len(octetStringBytes), nil
}

//#E@T390301CLH
