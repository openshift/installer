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

/* Converted to Go from com.ibm.tke.model.pcix.OA2CertificateX.java */
/* Represents OA certificates for the CEX6C/CEX6P and later crypto modules */

type OA2CertificateX struct {
	Header     []byte
	Info       []byte
	MetaData   []byte
	AuxData    []byte
	Spki       []byte
	SignerInfo []byte

	HdrMagic           byte
	HdrSpare           []byte
	HdrSectionCount    byte
	HdrByteCount       uint32
	HdrSectionLengths0 uint32
	HdrSectionLengths1 uint32
	HdrSectionLengths2 uint32
	HdrSectionLengths3 uint32
	HdrSectionLengths4 uint32
	HdrPadBytes        []byte

	InfoAlgID    uint32
	InfoPadBytes []byte

	MetaDataDer1         []byte
	MetaDataObjVersion   uint32
	MetaDataDer2         []byte
	MetaDataRootCA       byte
	MetaDataTrustBits    byte
	MetaDataIssuingSeg   byte
	MetaDataSubjectSeg   byte
	MetaDataDer3         []byte
	MetaDataAdapterType  []byte
	MetaDataDer4         []byte
	MetaDataAdapterID    []byte
	MetaDataDer5         []byte
	MetaDataSourceConfig []byte
	MetaDataDer6         []byte
	MetaDataCreationBoot uint32
	MetaDataDer7         []byte
	MetaDataKeyType      uint32
	MetaDataDer8         []byte
	MetaDataSubjectSKI   []byte
	MetaDataDer9         []byte
	MetaDataSignerSKI    []byte
	MetaDataDer10        []byte
	MetaDataSigningTime  []byte
	MetaDataDer11        []byte
	MetaDataOAInfo       []byte
	MetaDataPadBytes     []byte

	AuxDataDer1     []byte
	AuxDataVersion  byte
	AuxDataDer2     []byte
	AuxDataSegHash  []byte
	AuxDataDer3     []byte
	AuxDataOwner2   []byte
	AuxDataDer4     []byte
	AuxDataOwner3   []byte
	AuxDataPadBytes []byte

	SpkiDer1        []byte
	SpkiOID1        []byte
	SpkiDer2        []byte
	SpkiOID2        []byte
	SpkiDer3        []byte
	SpkiCompressed  byte
	SpkiXCoordinate []byte
	SpkiYCoordinate []byte
	SpkiPadBytes    []byte

	SpkiPublicKey []byte

	Body []byte // everything except the SignerInfo

	SignerInfoDer1        []byte
	SignerInfoVersion     byte
	SignerInfoDer2        []byte
	SignerInfoSKI         []byte
	SignerInfoDer3        []byte
	SignerInfoDigestAlgID []byte
	SignerInfoDer4        []byte
	SignerInfoSigAlgID    []byte
	SignerInfoDer5        []byte
	SignerInfoR           []byte
	SignerInfoDer6        []byte
	SignerInfoS           []byte
	SignerInfoPadByte     byte
}

/*----------------------------------------------------------------------------*/
/* Initializes an OA2CertificateX structure using data read from a crypto     */
/* module                                                                     */
/*----------------------------------------------------------------------------*/
func (cert *OA2CertificateX) Init(data []byte) error {

	// Process the header
	if len(data) < 32 {
		/* This and other length checks will detect some problems with the input
		   data, but not all.  If an array index is out of range, Go panics and
		   displays a call stack. */
		return errors.New("OA2CertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing header.")
	}
	cert.Header             = data[0:32]
	
	cert.HdrMagic           = data[0]
	cert.HdrSpare           = data[1:3]
	cert.HdrSectionCount    = data[3]
	cert.HdrByteCount       = binary.BigEndian.Uint32(data[4:8])
	cert.HdrSectionLengths0 = binary.BigEndian.Uint32(data[8:12])
	cert.HdrSectionLengths1 = binary.BigEndian.Uint32(data[12:16])
	cert.HdrSectionLengths2 = binary.BigEndian.Uint32(data[16:20])
	cert.HdrSectionLengths3 = binary.BigEndian.Uint32(data[20:24])
	cert.HdrSectionLengths4 = binary.BigEndian.Uint32(data[24:28])
	cert.HdrPadBytes        = data[28:32]

	var info_pad_bytes_len, metadata_pad_bytes_len, auxdata_pad_bytes_len,
		spki_pad_bytes_len, signerinfo_pad_bytes_len uint32 = 0, 0, 0, 0, 0

	if cert.HdrSectionLengths0 > 0 {
		info_pad_bytes_len = 4
	}
	if cert.HdrSectionLengths1 > 0 {
		metadata_pad_bytes_len = 3
	}
	if cert.HdrSectionLengths2 > 0 {
		auxdata_pad_bytes_len = 5
	}
	if cert.HdrSectionLengths3 > 0 {
		spki_pad_bytes_len = 2
	}
	if cert.HdrSectionLengths4 > 0 {
		signerinfo_pad_bytes_len = 1
	}

	// Process the info bytes
	checklen := 32 + cert.HdrSectionLengths0 + info_pad_bytes_len
	if len(data) < int(checklen) {
		return errors.New("OA2CertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing info section.")
	}
	cert.Info = data[32 : 32+cert.HdrSectionLengths0+info_pad_bytes_len]

	cert.InfoAlgID = binary.BigEndian.Uint32(cert.Info[0:4])
	cert.InfoPadBytes = cert.Info[4:8]

	// Process the meta data
	checklen = 32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
		cert.HdrSectionLengths1 + metadata_pad_bytes_len
	if len(data) < int(checklen) {
		return errors.New("OA2CertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing meta data.")
	}

	cert.MetaData = data[ 32 + cert.HdrSectionLengths0 + info_pad_bytes_len :
	                      32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
		                       cert.HdrSectionLengths1 + metadata_pad_bytes_len ]

	cert.MetaDataDer1         = cert.MetaData[0:5]
	cert.MetaDataObjVersion   = binary.BigEndian.Uint32(cert.MetaData[5:9])
	cert.MetaDataDer2         = cert.MetaData[9:11]
	cert.MetaDataRootCA       = cert.MetaData[11]
	cert.MetaDataTrustBits    = cert.MetaData[12]
	cert.MetaDataIssuingSeg   = cert.MetaData[13]
	cert.MetaDataSubjectSeg   = cert.MetaData[14]
	cert.MetaDataDer3         = cert.MetaData[15:17]
	cert.MetaDataAdapterType  = cert.MetaData[17:25]
	cert.MetaDataDer4         = cert.MetaData[25:27]
	cert.MetaDataAdapterID    = cert.MetaData[27:39]
	cert.MetaDataDer5         = cert.MetaData[39:41]
	cert.MetaDataSourceConfig = cert.MetaData[41:73]
	cert.MetaDataDer6         = cert.MetaData[73:75]
	cert.MetaDataCreationBoot = binary.BigEndian.Uint32(cert.MetaData[75:79])
	cert.MetaDataDer7         = cert.MetaData[79:81]
	cert.MetaDataKeyType      = binary.BigEndian.Uint32(cert.MetaData[81:85])
	cert.MetaDataDer8         = cert.MetaData[85:87]
	cert.MetaDataSubjectSKI   = cert.MetaData[87:119]
	cert.MetaDataDer9         = cert.MetaData[119:121]
	cert.MetaDataSignerSKI    = cert.MetaData[121:153]
	cert.MetaDataDer10        = cert.MetaData[153:155]
	cert.MetaDataSigningTime  = cert.MetaData[155:171]
	cert.MetaDataDer11        = cert.MetaData[171:173]        

	metaDataOAInfo_len := int(cert.MetaDataDer11[1])
	if metaDataOAInfo_len > 0 {
		cert.MetaDataOAInfo = cert.MetaData[173 : 173+metaDataOAInfo_len]
	} else {
		cert.MetaDataOAInfo = nil
	}

	cert.MetaDataPadBytes = cert.MetaData[173 + metaDataOAInfo_len :
	                                      176 + metaDataOAInfo_len ]

	// Process the auxiliary data, if present
	if cert.HdrSectionLengths2 > 0 {
		checklen = 32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
			cert.HdrSectionLengths1 + metadata_pad_bytes_len +
			cert.HdrSectionLengths2 + auxdata_pad_bytes_len
		if len(data) < int(checklen) {
			return errors.New("OA2CertificateX data length too short, length = " +
				strconv.Itoa(len(data)) + ", processing auxiliary data.")
		}

		cert.AuxData =
			data[32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
			          cert.HdrSectionLengths1 + metadata_pad_bytes_len :
			     32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
			          cert.HdrSectionLengths1 + metadata_pad_bytes_len +
			          cert.HdrSectionLengths2 + auxdata_pad_bytes_len]
					
		cert.AuxDataDer1     = cert.AuxData[0:4]
		cert.AuxDataVersion  = cert.AuxData[4]
		cert.AuxDataDer2     = cert.AuxData[5:7]
		cert.AuxDataSegHash  = cert.AuxData[7:71]
		cert.AuxDataDer3     = cert.AuxData[71:73]
		cert.AuxDataOwner2   = cert.AuxData[73:77]
		cert.AuxDataDer4     = cert.AuxData[77:79]
		cert.AuxDataOwner3   = cert.AuxData[79:83]
		cert.AuxDataPadBytes = cert.AuxData[83:88]
	}

	// Process the SPKI (Subject Public Key Identifier)
	checklen = 32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
		cert.HdrSectionLengths1 + metadata_pad_bytes_len +
		cert.HdrSectionLengths2 + auxdata_pad_bytes_len +
		cert.HdrSectionLengths3 + spki_pad_bytes_len
	if len(data) < int(checklen) {
		return errors.New("OA2CertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing SPKI.")
	}

	cert.Spki =
		data[32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
		          cert.HdrSectionLengths1 + metadata_pad_bytes_len +
		          cert.HdrSectionLengths2 + auxdata_pad_bytes_len :
		     32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
		          cert.HdrSectionLengths1 + metadata_pad_bytes_len +
		          cert.HdrSectionLengths2 + auxdata_pad_bytes_len +
		          cert.HdrSectionLengths3 + spki_pad_bytes_len]
				
	cert.SpkiDer1        = cert.Spki[0:7]
	cert.SpkiOID1        = cert.Spki[7:14]
	cert.SpkiDer2        = cert.Spki[14:16]
	cert.SpkiOID2        = cert.Spki[16:21]
	cert.SpkiDer3        = cert.Spki[21:25]
	cert.SpkiCompressed  = cert.Spki[25]
	cert.SpkiXCoordinate = cert.Spki[26:92]
	cert.SpkiYCoordinate = cert.Spki[92:158]
	cert.SpkiPadBytes    = cert.Spki[158:160]
	cert.SpkiPublicKey   = cert.Spki[25:158]

	cert.Body = data[0 :
	                 32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
	                      cert.HdrSectionLengths1 + metadata_pad_bytes_len +
	                      cert.HdrSectionLengths2 + auxdata_pad_bytes_len +
	                      cert.HdrSectionLengths3 + spki_pad_bytes_len]

	// Process the signer info
	checklen = 32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
	                cert.HdrSectionLengths1 + metadata_pad_bytes_len +
	                cert.HdrSectionLengths2 + auxdata_pad_bytes_len +
	                cert.HdrSectionLengths3 + spki_pad_bytes_len +
	                cert.HdrSectionLengths4 + signerinfo_pad_bytes_len
	if len(data) < int(checklen) {
		return errors.New("OA2CertificateX data length too short, length = " +
			strconv.Itoa(len(data)) + ", processing signer info.")
	}

	cert.SignerInfo =
		data[32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
		          cert.HdrSectionLengths1 + metadata_pad_bytes_len +
		          cert.HdrSectionLengths2 + auxdata_pad_bytes_len +
		          cert.HdrSectionLengths3 + spki_pad_bytes_len :
		     32 + cert.HdrSectionLengths0 + info_pad_bytes_len +
		          cert.HdrSectionLengths1 + metadata_pad_bytes_len +
		          cert.HdrSectionLengths2 + auxdata_pad_bytes_len +
		          cert.HdrSectionLengths3 + spki_pad_bytes_len +
		          cert.HdrSectionLengths4 + signerinfo_pad_bytes_len]
	        
	cert.SignerInfoDer1        = cert.SignerInfo[0:5]
	cert.SignerInfoVersion     = cert.SignerInfo[5]
	cert.SignerInfoDer2        = cert.SignerInfo[6:8]
	cert.SignerInfoSKI         = cert.SignerInfo[8:40]
	cert.SignerInfoDer3        = cert.SignerInfo[40:44]
	cert.SignerInfoDigestAlgID = cert.SignerInfo[44:53]
	cert.SignerInfoDer4        = cert.SignerInfo[53:57]
	cert.SignerInfoSigAlgID    = cert.SignerInfo[57:65]
	cert.SignerInfoDer5        = cert.SignerInfo[65:73]
	cert.SignerInfoR           = cert.SignerInfo[73:139]
	cert.SignerInfoDer6        = cert.SignerInfo[139:141]
	cert.SignerInfoS           = cert.SignerInfo[141:207]
	cert.SignerInfoPadByte     = cert.SignerInfo[207]
	
	return nil
}
