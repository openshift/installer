//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 11/13/2020    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"encoding/binary"
	"errors"
	"strconv"
)

/* Converted to Go from com.ibm.tke.model.pcix.OACertificateX.java */
/* Represents OA certificates for the CEX5C/CEX5P and earlier crypto modules */

type OACertificateX struct {
	Header      []byte
	Body        []byte
	PublicKey   []byte
	DescriptorA []byte
	DescriptorB []byte
	Signature   []byte

	HeaderIdName      byte
	HeaderIdVersion   byte
	HeaderTData       uint32
	HeaderVDataOffset uint32
	HeaderVDataLen    uint32
	HeaderVSigOffset  uint32
	HeaderVSigLen     uint32
	HeaderTSig        uint32
	HeaderCkoName     []byte
	HeaderCkoType     uint32
	HeaderCkoStatus   uint32
	HeaderParentName  []byte

	BodyIdName              byte
	BodyIdVersion           byte
	BodyTPublic             uint32
	BodyVPublicOffset       uint32
	BodyVPublicLen          uint32
	BodyVDescAOffset        uint32
	BodyVDescALen           uint32
	BodyVDescBOffset        uint32
	BodyVDescBLen           uint32
	BodyDeviceNameIdName    byte
	BodyDeviceNameIdVersion byte
	BodyDeviceNameAdapterId []byte
	BodyCkoName             []byte
	BodyCkoNameNameType     uint16
	BodyCkoNameIndex        uint16
	BodyCkoType             uint32
	BodyParentName          []byte
	BodyParentNameNameType  uint16
	BodyParentNameIndex     uint16

	// fields set when the certificate contains an RSA public key
	PublicKeyType        uint32
	PublicKeyTokenLength uint32
	PublicKeyNBitLength  uint32
	PublicKeyNLength     uint32
	PublicKeyELength     uint32
	PublicKeyNOffset     uint32
	PublicKeyEOffset     uint32
	PublicKeyTokenData   []byte
	PublicKeyN           []byte
	PublicKeyE           []byte

	// fields set when the certificate contains an EC public key
	ECCCurveType     uint32
	ECCCurveSize     uint16
	ECCPublicKeyQLen uint16
	ECCPublicKeyQ    []byte
	ECCPublicKeyX    []byte
	ECCPublicKeyY    []byte

	// fields set when the certificate contains an EC signature
	ECCSignatureR []byte
	ECCSignatureS []byte
}

const OA_NEW_CERT = 1
const OA_RSA = 0
const OA_ECC = 2

const ECC_PRIME = 0x00
const ECC_BRAINPOOL = 0x01

/*----------------------------------------------------------------------------*/
/* Initializes an OACertificateX structure using data read from a crypto      */
/* module                                                                     */
/*----------------------------------------------------------------------------*/
func (cert *OACertificateX) Init(data []byte) error {

	/* This and other length checks will detect some problems with the input
	   data, but not all.  If an array index is out of range, Go panics and
	   displays a call stack. */
	if len(data) < 52 {
		return errors.New("OACertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing header.")
	}

	cert.Header            = data[0:52]
	cert.HeaderIdName      = data[0]
	cert.HeaderIdVersion   = data[1]
	cert.HeaderTData       = binary.BigEndian.Uint32(data[4:8])
	cert.HeaderVDataOffset = binary.BigEndian.Uint32(data[8:12])
	cert.HeaderVDataLen    = binary.BigEndian.Uint32(data[12:16])
	cert.HeaderVSigOffset  = binary.BigEndian.Uint32(data[16:20])
	cert.HeaderVSigLen     = binary.BigEndian.Uint32(data[20:24])
	cert.HeaderTSig        = binary.BigEndian.Uint32(data[24:28])
	cert.HeaderCkoName     = data[28:36]
	cert.HeaderCkoType     = binary.BigEndian.Uint32(data[36:40])
	cert.HeaderCkoStatus   = binary.BigEndian.Uint32(data[40:44])
	cert.HeaderParentName  = data[44:52]

	checklen := int(8 + cert.HeaderVDataOffset + cert.HeaderVDataLen)
	if len(data) < checklen {
		return errors.New("OACertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing body.")
	}

	cert.Body = data[8 + cert.HeaderVDataOffset :
	                 8 + cert.HeaderVDataOffset + cert.HeaderVDataLen]
	cert.BodyIdName              = cert.Body[0]
	cert.BodyIdVersion           = cert.Body[1]
	cert.BodyTPublic             = binary.BigEndian.Uint32(cert.Body[4:8])
	cert.BodyVPublicOffset       = binary.BigEndian.Uint32(cert.Body[8:12])
	cert.BodyVPublicLen          = binary.BigEndian.Uint32(cert.Body[12:16])
	cert.BodyVDescAOffset        = binary.BigEndian.Uint32(cert.Body[16:20])
	cert.BodyVDescALen           = binary.BigEndian.Uint32(cert.Body[20:24])
	cert.BodyVDescBOffset        = binary.BigEndian.Uint32(cert.Body[24:28])
	cert.BodyVDescBLen           = binary.BigEndian.Uint32(cert.Body[28:32])
	cert.BodyDeviceNameIdName    = cert.Body[32]
	cert.BodyDeviceNameIdVersion = cert.Body[33]
	cert.BodyDeviceNameAdapterId = cert.Body[36:44]
	cert.BodyCkoName             = cert.Body[52:60]
	cert.BodyCkoNameNameType     = binary.BigEndian.Uint16(cert.Body[56:58])
	cert.BodyCkoNameIndex        = binary.BigEndian.Uint16(cert.Body[58:60])
	cert.BodyCkoType             = binary.BigEndian.Uint32(cert.Body[60:64])
	cert.BodyParentName          = cert.Body[64:72]
	cert.BodyParentNameNameType  = binary.BigEndian.Uint16(cert.Body[68:70])
	cert.BodyParentNameIndex     = binary.BigEndian.Uint16(cert.Body[70:72])

	checklen = int(8 + cert.BodyVPublicOffset + cert.BodyVPublicLen)
	if len(cert.Body) < checklen {
		return errors.New("OACertificateX body length too short, body length = " +
			strconv.Itoa(len(cert.Body)) + ", processing public key.")
	}

	cert.PublicKey = cert.Body[8 + cert.BodyVPublicOffset :
	                           8 + cert.BodyVPublicOffset + cert.BodyVPublicLen]  

	if cert.BodyTPublic == OA_RSA {
		cert.PublicKeyType        = binary.BigEndian.Uint32(cert.PublicKey[0:4])
		cert.PublicKeyTokenLength = binary.BigEndian.Uint32(cert.PublicKey[4:8])
		cert.PublicKeyNBitLength  = binary.BigEndian.Uint32(cert.PublicKey[8:12])
		cert.PublicKeyNLength     = binary.BigEndian.Uint32(cert.PublicKey[12:16])
		cert.PublicKeyELength     = binary.BigEndian.Uint32(cert.PublicKey[16:20])
		cert.PublicKeyNOffset     = binary.BigEndian.Uint32(cert.PublicKey[52:56])
		cert.PublicKeyEOffset     = binary.BigEndian.Uint32(cert.PublicKey[56:60])
		cert.PublicKeyTokenData   = cert.PublicKey[92 : cert.PublicKeyTokenLength]
		cert.PublicKeyN = cert.PublicKey[cert.PublicKeyNOffset :
		                                 cert.PublicKeyNOffset + cert.PublicKeyNLength]
		cert.PublicKeyE = cert.PublicKey[cert.PublicKeyEOffset :
		                                 cert.PublicKeyEOffset + cert.PublicKeyELength]
	} else {
		cert.ECCCurveType     = uint32(cert.PublicKey[20])
		cert.ECCCurveSize     = binary.BigEndian.Uint16(cert.PublicKey[22:24])
		cert.ECCPublicKeyQLen = binary.BigEndian.Uint16(cert.PublicKey[24:26])
		cert.ECCPublicKeyQ    = cert.PublicKey[26 : 26 + cert.ECCPublicKeyQLen]
		
		if (cert.ECCPublicKeyQLen > 1) && (cert.ECCPublicKeyQLen % 2 == 1) {
			pointLen := (cert.ECCPublicKeyQLen - 1) / 2
			cert.ECCPublicKeyX = cert.ECCPublicKeyQ[1 : 1 + pointLen]
			cert.ECCPublicKeyY = cert.ECCPublicKeyQ[1 + pointLen :]	
		}
	}

	checklen = int(16 + cert.BodyVDescAOffset + cert.BodyVDescALen)
	if len(cert.Body) < checklen {
		return errors.New("OACertificateX body length too short, body length = " +
			strconv.Itoa(len(cert.Body)) + ", processing Descriptor A.")
	}
	cert.DescriptorA = cert.Body[16 + cert.BodyVDescAOffset :
	                             16 + cert.BodyVDescAOffset + cert.BodyVDescALen]	                             

	checklen = int(24 + cert.BodyVDescBOffset + cert.BodyVDescBLen)
	if len(cert.Body) < checklen {
		return errors.New("OACertificateX body length too short, body length = " +
			strconv.Itoa(len(cert.Body)) + ", processing Descriptor B.")
	}
	cert.DescriptorB = cert.Body[24 + cert.BodyVDescBOffset :
	                             24 + cert.BodyVDescBOffset + cert.BodyVDescBLen]

	checklen = int(16 + cert.HeaderVSigOffset + cert.HeaderVSigLen)
	if len(data) < checklen {
		return errors.New("OACertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing signature.")
	}
	cert.Signature = data[16 + cert.HeaderVSigOffset :
	                      16 + cert.HeaderVSigOffset + cert.HeaderVSigLen]
	if cert.BodyTPublic == OA_ECC {
		sublen := (len(cert.Signature) - 8) / 2
		cert.ECCSignatureR = cert.Signature[8 : 8 + sublen] 
		cert.ECCSignatureS = cert.Signature[8 + sublen :]
	}

	return nil
}
