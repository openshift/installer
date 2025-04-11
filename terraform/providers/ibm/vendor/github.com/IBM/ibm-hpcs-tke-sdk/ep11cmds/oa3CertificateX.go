//
// Copyright contributors to the ibm-hpcs-tke-sdk project
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/11/2024    CLH             Fix issue with EccBody and DilithiumBody
// 01/09/2025    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"encoding/binary"
	"errors"
	"strconv"
)

/* Converted to Go from com.ibm.tke.model.pcix.OA3CertificateX.java */
/* Represents OA certificates for CEX8 crypto modules */

type OA3CertificateX struct {
	Der1 []byte
	MetaData []byte
	AuxData []byte
	SignerInfo []byte

	MetaDataDer1 []byte
	MetaDataObjVersion uint32
	MetaDataDer2 []byte
	MetaDataTrustBits byte
	MetaDataIssuingSeg byte
	MetaDataSubjectSeg byte
	MetaDataDer3 []byte
	MetaDataAdapterType []byte
	MetaDataDer4 []byte
	MetaDataAdapterID []byte
	MetaDataDer5 []byte
	MetaDataSourceConfig []byte
	MetaDataDer6 []byte
	MetaDataCreationBoot uint32
	MetaDataDer7 []byte
	MetaDataSigningTime []byte

	AuxDataDer1 []byte
	AuxDataVersion byte
	AuxDataDer2 []byte
	AuxDataSegHash []byte
	AuxDataDer3 []byte
	AuxDataOwner2 uint32
	AuxDataDer4 []byte
	AuxDataOwner3 uint32

	SignerInfoDer1 []byte
	SignerInfoVersion byte
	SignerInfoDer2 []byte
	SignerInfoEccSignerSKI []byte
	SignerInfoDer3 []byte
	SignerInfoEccDigestAlgID []byte
	SignerInfoDer4 []byte
	SignerInfoEccKeyInfo []byte   // For subfields, see below
	SignerInfoEccKeyValue []byte   // For subfields, see below
	SignerInfoDer5 []byte
	SignerInfoEccSigAlgID []byte
	SignerInfoDer6 []byte
	SignerInfoR []byte
	SignerInfoDer7 []byte
	SignerInfoS []byte
	SignerInfoDer8 []byte
	SignerInfoDilVersion byte
	SignerInfoDer9 []byte
	SignerInfoDilSignerSKI []byte
	SignerInfoDer10 []byte
	SignerInfoDilDigestAlgID []byte
	SignerInfoDer11 []byte
	SignerInfoDilKeyInfo []byte   // For subfields, see below
	SignerInfoDilKeyValue []byte   // For subfields, see below
	SignerInfoDer12 []byte
	SignerInfoDilSigAlgID []byte
	SignerInfoDer13 []byte
	SignerInfoDilSig []byte   // For subfields, see below

	// Subfields of SignerInfoEccKeyInfo
	EccKeyMetaDataDer1 []byte
	EccKeyMetaDataObjectVersion uint32
	EccKeyMetaDataDer2 []byte
	EccKeyMetaDataRootCA byte
	EccKeyMetaDataDer3 []byte
	EccKeyMetaDataKeyType uint32
	EccKeyMetaDataDer4 []byte
	EccKeyMetaDataAlgorithmID uint32
	EccKeyMetaDataDer5 []byte
	EccKeyMetaDataSubjectSKI []byte
	EccKeyMetaDataDer6 []byte
	EccKeyMetaDataOAInfo []byte

	// Subfields of SignerInfoEccKeyValue
	SpkiDer1 []byte
	SpkiOID1 []byte
	SpkiDer2 []byte
	SpkiOID2 []byte
	SpkiDer3 []byte
	SpkiCompressed byte
	SpkiXCoordinate []byte
	SpkiYCoordinate []byte

	SpkiPublicKey []byte

	// Subfields of SignerInfoDilKeyInfo
	DilKeyMetaDataDer1 []byte
	DilKeyMetaDataObjectVersion uint32
	DilKeyMetaDataDer2 []byte
	DilKeyMetaDataRootCA byte
	DilKeyMetaDataDer3 []byte
	DilKeyMetaDataKeyType uint32
	DilKeyMetaDataDer4 []byte
	DilKeyMetaDataAlgorithmID uint32
	DilKeyMetaDataDer5 []byte
	DilKeyMetaDataSubjectSKI []byte
	DilKeyMetaDataDer6 []byte
	DilKeyMetaDataOAInfo []byte

	// Subfields of SignerInfoDilKeyValue
	DilithiumPublicDer1 []byte
	DilithiumPublicOID1 []byte
	DilithiumPublicDer2 []byte
	DilithiumPublicDer3 []byte
	DilithiumPublicNonce []byte
	DilithiumPublicDer4 []byte
	DilithiumPublicT1 []byte

	// Subfields of SignerInfoDilSig
	DilithiumSigDer1 []byte
	DilithiumSigOID1 []byte
	DilithiumSigDer2 []byte
	DilithiumSigSig []byte

	EccBody []byte  // the certificate bytes signed by the ECC key of the parent
	DilithiumBody []byte  // the certificate bytes signed by the Dilithium key of the parent
}

const CEX8_OA_CERTIFICATE_LENGTH = 7973
const CEX8_MB_CERTIFICATE_LENGTH = 7810

/*----------------------------------------------------------------------------*/
/* Initializes an OA3CertificateX structure using data read from a crypto     */
/* module                                                                     */
/*----------------------------------------------------------------------------*/
func (cert *OA3CertificateX) Init(data []byte) error {

	// There are two forms of OA certificates on the CEX8.  Certificates
	// created by segment 3 have 80-byte OA_info fields in the key metadata
	// fields of the SignerInfo.  For certificates created by miniboot this
	// field is empty.  The difference changes length values in the ASN.1
	// structure and this changes the length of some "der" fields.
	// ("der" fields contain ASN.1 tags and lengths.)  The two forms can be
	// distinguished by the overall length of the certificate.

	var isMBCert bool
	if len(data) == 7810 {
		isMBCert = true
	} else if len(data) == 7973 {
		isMBCert = false
	} else {
		return errors.New("Invalid data for OA3CertificateX, length = " +
			strconv.Itoa(len(data)) + ".")
	}

	cert.Der1       = data[0:4]
	cert.MetaData   = data[4:99]
	cert.AuxData    = data[99:182]
	cert.SignerInfo = data[182:]

	cert.MetaDataDer1         = cert.MetaData[0:4]
	cert.MetaDataObjVersion   = binary.BigEndian.Uint32(cert.MetaData[4:8])
	cert.MetaDataDer2         = cert.MetaData[8:10]
	cert.MetaDataTrustBits    = cert.MetaData[10]
	cert.MetaDataIssuingSeg   = cert.MetaData[11]
	cert.MetaDataSubjectSeg   = cert.MetaData[12]
	cert.MetaDataDer3         = cert.MetaData[13:15]
	cert.MetaDataAdapterType  = cert.MetaData[15:23]
	cert.MetaDataDer4         = cert.MetaData[23:25]
	cert.MetaDataAdapterID    = cert.MetaData[25:37]
	cert.MetaDataDer5         = cert.MetaData[37:39]
	cert.MetaDataSourceConfig = cert.MetaData[39:71]
	cert.MetaDataDer6         = cert.MetaData[71:73]
	cert.MetaDataCreationBoot = binary.BigEndian.Uint32(cert.MetaData[73:77])
	cert.MetaDataDer7         = cert.MetaData[77:79]
	cert.MetaDataSigningTime  = cert.MetaData[79:95]

	cert.AuxDataDer1    = cert.AuxData[0:4]
	cert.AuxDataVersion = cert.AuxData[4]
	cert.AuxDataDer2    = cert.AuxData[5:7]
	cert.AuxDataSegHash = cert.AuxData[7:71]
	cert.AuxDataDer3    = cert.AuxData[71:73]
	cert.AuxDataOwner2  = binary.BigEndian.Uint32(cert.AuxData[73:77])
	cert.AuxDataDer4    = cert.AuxData[77:79]
	cert.AuxDataOwner3  = binary.BigEndian.Uint32(cert.AuxData[79:83])

	cert.SignerInfoDer1               = cert.SignerInfo[0:6]
	cert.SignerInfoVersion            = cert.SignerInfo[6]
	cert.SignerInfoDer2               = cert.SignerInfo[7:9]
	cert.SignerInfoEccSignerSKI       = cert.SignerInfo[9:41]
	cert.SignerInfoDer3               = cert.SignerInfo[41:45]
	cert.SignerInfoEccDigestAlgID     = cert.SignerInfo[45:54]
	if isMBCert {
		cert.SignerInfoDer4           = cert.SignerInfo[54:59]
		cert.SignerInfoEccKeyInfo     = cert.SignerInfo[59:116]
		cert.SignerInfoEccKeyValue    = cert.SignerInfo[116:274]
		cert.SignerInfoDer5           = cert.SignerInfo[274:278]
		cert.SignerInfoEccSigAlgID    = cert.SignerInfo[278:286]
		cert.SignerInfoDer6           = cert.SignerInfo[286:294]
		cert.SignerInfoR              = cert.SignerInfo[294:360]
		cert.SignerInfoDer7           = cert.SignerInfo[360:362]
		cert.SignerInfoS              = cert.SignerInfo[362:428]
		cert.SignerInfoDer8           = cert.SignerInfo[428:434]
		cert.SignerInfoDilVersion     = cert.SignerInfo[434]
		cert.SignerInfoDer9           = cert.SignerInfo[435:437]
		cert.SignerInfoDilSignerSKI   = cert.SignerInfo[437:469]
		cert.SignerInfoDer10          = cert.SignerInfo[469:473]
		cert.SignerInfoDilDigestAlgID = cert.SignerInfo[473:482]
		cert.SignerInfoDer11          = cert.SignerInfo[482:488]
		cert.SignerInfoDilKeyInfo     = cert.SignerInfo[488:545]
		cert.SignerInfoDilKeyValue    = cert.SignerInfo[545:2919]
		cert.SignerInfoDer12          = cert.SignerInfo[2919:2923]
		cert.SignerInfoDilSigAlgID    = cert.SignerInfo[2923:2934]
		cert.SignerInfoDer13          = cert.SignerInfo[2934:2938]
		cert.SignerInfoDilSig         = cert.SignerInfo[2938:7628]
	} else {
		cert.SignerInfoDer4           = cert.SignerInfo[54:61]
		cert.SignerInfoEccKeyInfo     = cert.SignerInfo[61:198]
		cert.SignerInfoEccKeyValue    = cert.SignerInfo[198:356]
		cert.SignerInfoDer5           = cert.SignerInfo[356:360]
		cert.SignerInfoEccSigAlgID    = cert.SignerInfo[360:368]
		cert.SignerInfoDer6           = cert.SignerInfo[368:376]
		cert.SignerInfoR              = cert.SignerInfo[376:442]
		cert.SignerInfoDer7           = cert.SignerInfo[442:444]
		cert.SignerInfoS              = cert.SignerInfo[444:510]
		cert.SignerInfoDer8           = cert.SignerInfo[510:516]
		cert.SignerInfoDilVersion     = cert.SignerInfo[516]
		cert.SignerInfoDer9           = cert.SignerInfo[517:519]
		cert.SignerInfoDilSignerSKI   = cert.SignerInfo[519:551]
		cert.SignerInfoDer10          = cert.SignerInfo[551:555]
		cert.SignerInfoDilDigestAlgID = cert.SignerInfo[555:564]
		cert.SignerInfoDer11          = cert.SignerInfo[564:571]
		cert.SignerInfoDilKeyInfo     = cert.SignerInfo[571:708]
		cert.SignerInfoDilKeyValue    = cert.SignerInfo[708:3082]
		cert.SignerInfoDer12          = cert.SignerInfo[3082:3086]
		cert.SignerInfoDilSigAlgID    = cert.SignerInfo[3086:3097]
		cert.SignerInfoDer13          = cert.SignerInfo[3097:3101]
		cert.SignerInfoDilSig         = cert.SignerInfo[3101:7791]
	}

	// Subfields of SignerInfoEccKeyInfo
	cert.EccKeyMetaDataDer1          = cert.SignerInfoEccKeyInfo[0:2]
	cert.EccKeyMetaDataObjectVersion = binary.BigEndian.Uint32(cert.SignerInfoEccKeyInfo[2:6])
	cert.EccKeyMetaDataDer2          = cert.SignerInfoEccKeyInfo[6:8]
	cert.EccKeyMetaDataRootCA        = cert.SignerInfoEccKeyInfo[8]
	cert.EccKeyMetaDataDer3          = cert.SignerInfoEccKeyInfo[9:11]
	cert.EccKeyMetaDataKeyType       = binary.BigEndian.Uint32(cert.SignerInfoEccKeyInfo[11:15])
	cert.EccKeyMetaDataDer4          = cert.SignerInfoEccKeyInfo[15:17]
	cert.EccKeyMetaDataAlgorithmID   = binary.BigEndian.Uint32(cert.SignerInfoEccKeyInfo[17:21])
	cert.EccKeyMetaDataDer5          = cert.SignerInfoEccKeyInfo[21:23]
	cert.EccKeyMetaDataSubjectSKI    = cert.SignerInfoEccKeyInfo[23:55]
	cert.EccKeyMetaDataDer6          = cert.SignerInfoEccKeyInfo[55:57]
	cert.EccKeyMetaDataOAInfo        = cert.SignerInfoEccKeyInfo[57:]  // empty for MB cert

	// Subfields of SignerInfoEccKeyValue
	cert.SpkiDer1        = cert.SignerInfoEccKeyValue[0:7]
	cert.SpkiOID1        = cert.SignerInfoEccKeyValue[7:14]
	cert.SpkiDer2        = cert.SignerInfoEccKeyValue[14:16]
	cert.SpkiOID2        = cert.SignerInfoEccKeyValue[16:21]
	cert.SpkiDer3        = cert.SignerInfoEccKeyValue[21:25]
	cert.SpkiCompressed  = cert.SignerInfoEccKeyValue[25]
	cert.SpkiXCoordinate = cert.SignerInfoEccKeyValue[26:92]
	cert.SpkiYCoordinate = cert.SignerInfoEccKeyValue[92:158]

	cert.SpkiPublicKey   = cert.SignerInfoEccKeyValue[25:158]

	// Subfields of SignerInfoDilKeyInfo
	cert.DilKeyMetaDataDer1          = cert.SignerInfoDilKeyInfo[0:2]
	cert.DilKeyMetaDataObjectVersion = binary.BigEndian.Uint32(cert.SignerInfoDilKeyInfo[2:6])
	cert.DilKeyMetaDataDer2          = cert.SignerInfoDilKeyInfo[6:8]
	cert.DilKeyMetaDataRootCA        = cert.SignerInfoDilKeyInfo[8]
	cert.DilKeyMetaDataDer3          = cert.SignerInfoDilKeyInfo[9:11]
	cert.DilKeyMetaDataKeyType       = binary.BigEndian.Uint32(cert.SignerInfoDilKeyInfo[11:15])
	cert.DilKeyMetaDataDer4          = cert.SignerInfoDilKeyInfo[15:17]
	cert.DilKeyMetaDataAlgorithmID   = binary.BigEndian.Uint32(cert.SignerInfoDilKeyInfo[17:21])
	cert.DilKeyMetaDataDer5          = cert.SignerInfoDilKeyInfo[21:23]
	cert.DilKeyMetaDataSubjectSKI    = cert.SignerInfoDilKeyInfo[23:55]
	cert.DilKeyMetaDataDer6          = cert.SignerInfoDilKeyInfo[55:57]
	cert.DilKeyMetaDataOAInfo        = cert.SignerInfoDilKeyInfo[57:]  // empty for MB cert

	// Subfields of SignerInfoDilKeyValue
	cert.DilithiumPublicDer1  = cert.SignerInfoDilKeyValue[0:8]
	cert.DilithiumPublicOID1  = cert.SignerInfoDilKeyValue[8:19]
	cert.DilithiumPublicDer2  = cert.SignerInfoDilKeyValue[19:26]
	cert.DilithiumPublicDer3  = cert.SignerInfoDilKeyValue[26:33]
	cert.DilithiumPublicNonce = cert.SignerInfoDilKeyValue[33:65]
	cert.DilithiumPublicDer4  = cert.SignerInfoDilKeyValue[65:70]
	cert.DilithiumPublicT1    = cert.SignerInfoDilKeyValue[70:2374]

	// Subfields of SignerInfoDilSig
	cert.DilithiumSigDer1 = cert.SignerInfoDilSig[0:4]
	cert.DilithiumSigOID1 = cert.SignerInfoDilSig[4:15]
	cert.DilithiumSigDer2 = cert.SignerInfoDilSig[15:22]
	cert.DilithiumSigSig  = cert.SignerInfoDilSig[22:4690]

	// Make a copy of the data signed by the ECC OA signature key
	// to avoid overwriting the base certificate bytes
	eccBodyData := make([]byte, len(cert.MetaData) + len(cert.AuxData) +
		len(cert.SignerInfoDer4) + len(cert.SignerInfoEccKeyInfo) +
		len(cert.SignerInfoEccKeyValue))
	for i := 0; i < len(cert.MetaData); i++ {
        eccBodyData[i] = cert.MetaData[i]
	}
	for i := 0; i < len(cert.AuxData); i++ {
		eccBodyData[i + len(cert.MetaData)] = cert.AuxData[i]
	}
	for i := 0; i < len(cert.SignerInfoDer4); i++ {
		eccBodyData[i + len(cert.MetaData) + len(cert.AuxData)] =
			cert.SignerInfoDer4[i]
	}
	for i := 0; i < len(cert.SignerInfoEccKeyInfo); i++ {
		eccBodyData[i + len(cert.MetaData) + len(cert.AuxData) +
			len(cert.SignerInfoDer4)] = cert.SignerInfoEccKeyInfo[i]
	}
	for i := 0; i< len(cert.SignerInfoEccKeyValue); i++ {
		eccBodyData[i + len(cert.MetaData) + len(cert.AuxData) +
			len(cert.SignerInfoDer4) + len(cert.SignerInfoEccKeyInfo)] =
				cert.SignerInfoEccKeyValue[i]
	}
	cert.EccBody = eccBodyData[:]

	// Make a copy of the data signed by the Dilithium OA signature key
	// to avoid overwriting the base certificate bytes
	dilithiumBodyData := make([]byte, len(cert.MetaData) + len(cert.AuxData) +
		len(cert.SignerInfoDer11) + len(cert.SignerInfoDilKeyInfo) +
		len(cert.SignerInfoDilKeyValue))
	for i := 0; i < len(cert.MetaData); i++ {
        dilithiumBodyData[i] = cert.MetaData[i]
	}
	for i := 0; i < len(cert.AuxData); i++ {
		dilithiumBodyData[i + len(cert.MetaData)] = cert.AuxData[i]
	}
	for i := 0; i < len(cert.SignerInfoDer11); i++ {
		dilithiumBodyData[i + len(cert.MetaData) + len(cert.AuxData)] =
			cert.SignerInfoDer11[i]
	}
	for i := 0; i < len(cert.SignerInfoDilKeyInfo); i++ {
		dilithiumBodyData[i + len(cert.MetaData) + len(cert.AuxData) +
			len(cert.SignerInfoDer11)] = cert.SignerInfoDilKeyInfo[i]
	}
	for i := 0; i< len(cert.SignerInfoDilKeyValue); i++ {
		dilithiumBodyData[i + len(cert.MetaData) + len(cert.AuxData) +
			len(cert.SignerInfoDer11) + len(cert.SignerInfoDilKeyInfo)] =
				cert.SignerInfoDilKeyValue[i]
	}
	cert.DilithiumBody = dilithiumBodyData[:]

	return nil
}
