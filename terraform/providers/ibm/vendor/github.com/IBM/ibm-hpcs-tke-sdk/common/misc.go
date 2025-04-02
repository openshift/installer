//
// Copyright contributors to the ibm-hpcs-tke-sdk project
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/30/2021    CLH             Modify for TKE SDK
// 01/09/2025    CLH             Set last four bytes of VP to zero

package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var DomainsFileName = "DOMAINS"
var CryptoModulesFileName = "CRYPTOMODULES"

/** Entry in the DOMAINS file describing a single domain.  This version
  lacks a "type" field. */
type DomainEntryNoType struct {
	Domain_num int    `json:"domain_num"`
	Hsm_id     string `json:"hsm_id"`
	// UUID for this particular domain
	Crypto_instance_id string `json:"crypto_instance_id"`
	// UUID for the crypto instance containing this domain
	Location string `json:"location"`
	// Describes the location of the domain
	// Format is [Availability zone].[Host].[Crypto module index].[domain index]
	Serial_num string `json:"serial_num"`
	Public_key string `json:"public_key"`
	Selected   bool   `json:"selected"`
	// Indicates whether the user has selected this domain to work with
}

/** Entry in the DOMAINS file describing a single domain */
type DomainEntry struct {
	Domain_num int    `json:"domain_num"`
	Hsm_id     string `json:"hsm_id"`
	// UUID for this particular domain
	Crypto_instance_id string `json:"crypto_instance_id"`
	// UUID for the crypto instance containing this domain
	Location string `json:"location"`
	// Describes the location of the domain
	// Format is [Availability zone].[Host].[Crypto module index].[domain index]
	Serial_num string `json:"serial_num"`
	Public_key string `json:"public_key"`
	Type       string `json:"type"` //@T390301CLH
	// "operational" or "recovery"  //@T407032CLH
	Selected bool `json:"selected"`
	// Indicates whether the user has selected this domain to work with
}

/** Entry in the CRYPTOMODULES file */
type CryptoModuleEntry struct {
	Serial_num string `json:"serial_num"`
	Public_key string `json:"public_key"`
}

type RotateStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

/** Used to work with an ASN.1 sequence representing an EC public key */
type ECPublicKey struct {
	X *big.Int
	Y *big.Int
}

/*----------------------------------------------------------------------------*/
/* Returns the crypto module index from the Location field of a DomainEntry   */
/*----------------------------------------------------------------------------*/
func (de DomainEntry) GetCryptoModuleIndex() int {
	parts := strings.Split(de.Location, ".")
	if len(parts) != 4 {
		panic("Invalid Location entry in DOMAINS file.")
	}
	// remove the surrounding brackets
	indexString := parts[2][1 : len(parts[2])-1]
	index, err := strconv.Atoi(indexString)
	if err != nil {
		panic("Invalid crypto module index in Location entry: '" + parts[2] + "'")
	}
	return index
}

/*----------------------------------------------------------------------------*/
/* Returns the domain index from the Location field of a DomainEntry          */
/*----------------------------------------------------------------------------*/
func (de DomainEntry) GetDomainIndex() int {
	return GetDomainIndexFromLocation(de.Location)
}

/*----------------------------------------------------------------------------*/
/* Returns the domain index from a location string                            */
/*----------------------------------------------------------------------------*/
func GetDomainIndexFromLocation(location string) int {
	parts := strings.Split(location, ".")
	if len(parts) != 4 {
		panic("Error extracting domain index from location: invalid location format.")
	}
	// remove the surrounding brackets
	indexString := parts[3][1 : len(parts[3])-1]
	index, err := strconv.Atoi(indexString)
	if err != nil {
		panic("Invalid domain index in location data: '" + parts[3] + "'")
	}
	return index
}

/*----------------------------------------------------------------------------*/
/* Returns the part of the location string that identifies a crypto module.   */
/* That is, everything except the domain index at the end.                    */
/*----------------------------------------------------------------------------*/
func GetPartialLocation(location string) string {
	parts := strings.Split(location, ".")
	if len(parts) != 4 {
		panic("Invalid location format: '" + location + "'")
	}
	return parts[0] + "." + parts[1] + "." + parts[2]
}

type AdminInfo struct {
	Domain DomainEntry
	Name   string
	Ski    string
}

/*----------------------------------------------------------------------------*/
/* Checks that the CLOUDTKEFILES environment variable is set and points to    */
/* a usable subdirectory on the local workstation.                            */
/*                                                                            */
/* Also changes the current working subdirectory to the directory identified  */
/* by CLOUDTKEFILES.                                                          */
/*----------------------------------------------------------------------------*/
func CheckSubdir() error {

	// Check that the environment variable is set
	subdir := os.Getenv("CLOUDTKEFILES")
	if subdir == "" {
		return errors.New("The CLOUDTKEFILES environment variable is not " +
			"defined on this workstation.\n\nSet the CLOUDTKEFILES " +
			"environment variable to indicate the subdirectory to hold " +
			"files used by the Cloud TKE plug-in.") //@T390301CLH
	}

	// Check that the environment variable points to a valid subdirectory
	err := os.Chdir(subdir)
	if err != nil {
		return errors.New("Error accessing the subdirectory defined by the " +
			"CLOUDTKEFILES environment variable.\n\nCheck that the " +
			"subdirectory exists and is specified correctly by the " +
			"environment variable.\n\nCLOUDTKEFILES=" + subdir)
	}

	return nil
}

/*----------------------------------------------------------------------------*/
/* Derives an AES key from a password.                                        */
/*                                                                            */
/* This hashes the password 4096 times to get the AES key.  A previously      */
/* used salt value may be supplied for existing files, otherwise a random     */
/* salt value will be generated and used.                                     */
/*----------------------------------------------------------------------------*/
func Derive_aes_key(passwd string, salt string) ([]byte, []byte) {
	pwbytes := []byte(passwd)
	var saltbytes []byte = nil
	var err error

	if salt == "" {
		saltbytes = make([]byte, 32)
		_, err = io.ReadFull(rand.Reader, saltbytes)
		if err != nil {
			panic(err)
		}
	} else {
		saltbytes, err = hex.DecodeString(salt)
		if err != nil {
			panic(err)
		}
	}
	return pbkdf2.Key(pwbytes, saltbytes, 4096, 32, sha256.New), saltbytes
}

/*----------------------------------------------------------------------------*/
/* Decrypts ciphertext using an AES key                                       */
/*                                                                            */
/* The input data is a nonce followed by the ciphertext.                      */
/* Returns the plaintext.                                                     */
/*----------------------------------------------------------------------------*/
func Decrypt(data []byte, key []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("decrypt: input data shorter than nonce size")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func Uint32To4ByteSlice(theInt uint32) []byte {
	tempSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(tempSlice, theInt)

	prependSize := 4 - len(tempSlice)
	returnSlice := make([]byte, 0, 4)
	for i := 0; i < prependSize; i++ {
		returnSlice = append(returnSlice, 0)
	}
	returnSlice = append(returnSlice, tempSlice...)
	return returnSlice
}

func FourByteSliceToInt(theSlice []byte) int {
	result := binary.BigEndian.Uint32(theSlice)
	return int(result)
}

/*----------------------------------------------------------------------------*/
/* Checks if two []byte are equal                                             */
/*----------------------------------------------------------------------------*/
func ByteSlicesAreEqual(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil && b != nil {
		return false
	} else if a != nil && b == nil {
		return false
	} else if len(a) != len(b) {
		return false
	} else {
		// Both slices are not nil and have the same length
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}

/*----------------------------------------------------------------------------*/
/* Calculates the verification pattern of an AES key part.                    */
/*                                                                            */
/* The verification pattern of a symmetric key is defined in section 6.7 of   */
/* the EP11 wire formats document as:                                         */
/* SHA_256( 01 || <raw_key> ), with the last four bytes set to zero.          */
/*                                                                            */
/* Old EP11 crypto modules did not zeroize the last four bytes of the         */
/* verification pattern.                                                      */
/*----------------------------------------------------------------------------*/
func Calc_vp(rawkey []byte) []byte {
	var temp []byte
	temp = append(temp, 0x01)
	temp = append(temp, rawkey...)
	hasher := sha256.New()
	hasher.Write(temp)
	out := hasher.Sum(nil)
	for i := 28; i <= 31; i++ {
		out[i] = 0
	}
	return out
}

/*----------------------------------------------------------------------------*/
/* Calculates the Subject Key Identifier for an RSA key.                      */
/*                                                                            */
/* Inputs:                                                                    */
/* rsa.PublicKey pubKey -- the RSA public key                                 */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the calculated subject key identifier                            */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func CalculateRSAKeyHash(pubKey rsa.PublicKey) ([]byte, error) {
	bytes, err := asn1.Marshal(pubKey)
	// This produces an ASN.1 sequence containing two integers:
	// the modulus and the public exponent.  The Subject Key
	// Identifier is the SHA-256 hash of this.
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(bytes)
	return hash[:], nil
}

/*----------------------------------------------------------------------------*/
/* Calculates the Subject Key Identifier of an EC key.                        */
/*                                                                            */
/* Inputs:                                                                    */
/* ecdsa.PublicKey pubKey -- the EC public key                                */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the calculated subject key identifier                            */
/*----------------------------------------------------------------------------*/
func CalculateECKeyHash(pubKey ecdsa.PublicKey) []byte {

	var publicBytes [133]byte
	publicBytes[0] = 0x04
	bytes := pubKey.X.Bytes()
	length := len(bytes)
	// copy the X coordinate
	for i := 0; i < length; i++ {
		publicBytes[1+(66-length)+i] = bytes[i]
	}
	bytes = pubKey.Y.Bytes()
	length = len(bytes)
	// copy the Y coordinate
	for i := 0; i < length; i++ {
		publicBytes[1+66+(66-length)+i] = bytes[i]
	}
	hash := sha256.Sum256(publicBytes[:])
	return hash[:]
}

/*----------------------------------------------------------------------------*/
/* Encrypts plaintext using an AES key                                        */
/*                                                                            */
/* Returns a nonce followed by the ciphertext                                 */
/*----------------------------------------------------------------------------*/
func Encrypt(plaintext []byte, aeskey []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

/*----------------------------------------------------------------------------*/
/* Checks if a bit in a []byte is set.                                        */
/*----------------------------------------------------------------------------*/
func IsBitSet(data []byte, bitnum int) bool {
	if bitnum < 0 {
		return false
	} else if bitnum >= len(data)*8 {
		return false
	} else if (data[bitnum/8] & (0x80 >> uint(bitnum%8))) == 0 {
		return false
	} else {
		return true
	}
}

/*----------------------------------------------------------------------------*/
/* Gets the public key from a signing service                                 */
/*                                                                            */
/* Inputs:                                                                    */
/* string -- base URL for the signing service                                 */
/* string -- identifies the signature key to be accessed                      */
/* string -- authentication token for the signature key to be accessed        */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the public key.  Only P521 EC signature keys are supported.      */
/*     This will be a compression byte (0x04) followed by a 66 byte X-value   */
/*     and a 66-byte Y value.                                                 */
/* error -- reports any error encountered during processing                   */
/*----------------------------------------------------------------------------*/
func GetPublicKeyFromSigningService(ssURL string, sigkey string, sigkeyToken string) ([]byte, error) {

	// Create dummy public key for error return
	rtnkey := make([]byte, 0)

	// Get the public key from the signing service
	req := CreateGetPublicKeyRequest(sigkeyToken, ssURL, sigkey)
	pubkey, err := SubmitQueryPublicKeyRequest(req)
	if err != nil {
		return rtnkey, err
	}

	// Decode the public key
	decodedKey, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		return rtnkey, err
	}

	// When a signing service is used, only P521 EC signature keys are
	// supported.

	// Check that the decoded key is an ASN.1 sequence with expected fields.
	// Reformat if needed so X and Y coordinates are each 66 bytes long.
	var ecpubkey ECPublicKey
	_, err = asn1.Unmarshal(decodedKey, &ecpubkey)
	if err != nil {
		return rtnkey, err // invalid ASN.1 sequence
	}

	var publicKey [133]byte
	publicKey[0] = 0x04
	bytes := ecpubkey.X.Bytes()
	length := len(bytes)
	if length > 66 {
		return rtnkey, errors.New("Invalid public key returned by signing service.  Only P521 EC keys are supported.")
	}
	// copy the X coordinate
	for i:=0; i<length; i++ {
		publicKey[1 + (66 - length) + i] = bytes[i]
	}
	bytes = ecpubkey.Y.Bytes()
	length = len(bytes)
	if length > 66 {
		return rtnkey, errors.New("Invalid public key returned by signing service.  Only P521 EC keys are supported.")
	}
	// copy the Y coordinate
	for i:=0; i<length; i++ {
		publicKey[1 + 66 + (66 - length) + i] = bytes[i]
	}

	return publicKey[:], nil
}

