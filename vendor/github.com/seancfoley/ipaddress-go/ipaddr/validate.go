//
// Copyright 2020-2024 Sean C Foley
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ipaddr

import (
	"math"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

var chars, extendedChars = createChars()

func createChars() (chars [int('z') + 1]byte, extendedChars [int('~') + 1]byte) {
	i := byte(1)
	for c := '1'; i < 10; i, c = i+1, c+1 {
		chars[c] = i
	}
	for c, c2 := 'a', 'A'; i < 26; i, c, c2 = i+1, c+1, c2+1 {
		chars[c] = i
		chars[c2] = i
	}
	var extendedDigits = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B',
		'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
		'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z', '!', '#', '$', '%', '&', '(', ')', '*', '+', '-',
		';', '<', '=', '>', '?', '@', '^', '_', '`', '{', '|', '}',
		'~'}
	extLen := byte(len(extendedDigits))
	for i = 0; i < extLen; i++ {
		c := extendedDigits[i]
		extendedChars[c] = i
	}
	return
}

const (
	longSize                           = 64
	maxHostLength                      = 253
	maxHostSegments                    = 127
	maxLabelLength                     = 63
	macDoubleSegmentDigitCount         = 6
	macExtendedDoubleSegmentDigitCount = 10
	macSingleSegmentDigitCount         = 12
	macExtendedSingleSegmentDigitCount = 16
	ipv6SingleSegmentDigitCount        = 32
	ipv6BinarySingleSegmentDigitCount  = 128
	ipv4BinarySingleSegmentDigitCount  = 32
	ipv6Base85SingleSegmentDigitCount  = 20
	maxWildcards                       = ipv6Base85SingleSegmentDigitCount - 1 // 20 wildcards is equivalent to a base 85 address
	ipv4SingleSegmentOctalDigitCount   = 11
	longHexDigits                      = longSize >> 2
	longBinaryDigits                   = longSize
)

var (
	macMaxTriple uint64 = (MACMaxValuePerSegment << uint64(MACBitsPerSegment*2)) |
		(MACMaxValuePerSegment << uint64(MACBitsPerSegment)) | MACMaxValuePerSegment
	macMaxQuintuple = (macMaxTriple << uint64(MACBitsPerSegment*2)) | uint64(macMaxTriple>>uint64(MACBitsPerSegment))
)

func isSingleSegmentIPv6(
	str string,
	totalDigits int,
	isRange bool,
	frontTotalDigits int,
	ipv6SpecificOptions addrstrparam.IPv6AddressStringParams) (isSingle bool, err addrerr.AddressStringError) {
	backIsIpv6 := totalDigits == ipv6SingleSegmentDigitCount || // 32 hex chars with or without 0x
		(ipv6SpecificOptions.AllowsBinary() && totalDigits == ipv6BinarySingleSegmentDigitCount+2) || // 128 binary chars with 0b
		(isRange && totalDigits == 0 && (frontTotalDigits == ipv6SingleSegmentDigitCount ||
			(ipv6SpecificOptions.AllowsBinary() && frontTotalDigits == ipv6BinarySingleSegmentDigitCount+2)))
	if backIsIpv6 && isRange && totalDigits != 0 {
		frontIsIpv6 := frontTotalDigits == ipv6SingleSegmentDigitCount ||
			(ipv6SpecificOptions.AllowsBinary() && frontTotalDigits == ipv6BinarySingleSegmentDigitCount+2) ||
			frontTotalDigits == 0
		if !frontIsIpv6 {
			err = &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments.digit.count"}}
			return
		}
	}
	isSingle = backIsIpv6
	return
}

// When checking for binary single segment, we must check for the exact number of digits for IPv4.
// This is because of ambiguity between IPv6 hex 32 chars starting with 0b and 0b before 30 binary chars.
// So we must therefore avoid 0b before 30 binary chars for IPv4.  We must require 0b before 32 binary chars.
// This only applies to single-segment.
// For segmented IPv4, there is no ambiguity and we allow binary segments of varying lengths,
// just like we do for inet_aton.

func isSingleSegmentIPv4(
	str string,
	nonZeroDigits,
	totalDigits int,
	isRange bool,
	frontNonZeroDigits,
	frontTotalDigits int,
	ipv4SpecificOptions addrstrparam.IPv4AddressStringParams) (isSingle bool, err addrerr.AddressStringError) {
	backIsIpv4 := nonZeroDigits <= ipv4SingleSegmentOctalDigitCount ||
		(ipv4SpecificOptions.AllowsBinary() && totalDigits == ipv4BinarySingleSegmentDigitCount+2) ||
		(isRange && totalDigits == 0 && (frontTotalDigits <= ipv4SingleSegmentOctalDigitCount ||
			(ipv4SpecificOptions.AllowsBinary() && frontTotalDigits == ipv4BinarySingleSegmentDigitCount+2)))
	if backIsIpv4 && isRange && totalDigits != 0 {
		frontIsIpv4 := frontNonZeroDigits <= ipv4SingleSegmentOctalDigitCount ||
			(ipv4SpecificOptions.AllowsBinary() && frontTotalDigits == ipv4BinarySingleSegmentDigitCount+2) ||
			frontTotalDigits == 0
		if !frontIsIpv4 {
			err = &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments.digit.count"}}
			return
		}
	}
	isSingle = backIsIpv4
	return
}

func hasExtraneousDigitsIPv4(
	validationOptions addrstrparam.IPAddressStringParams,
	ipv4SpecificOptions addrstrparam.IPv4AddressStringParams,
	ipv6SpecificOptions addrstrparam.IPv6AddressStringParams,
	totalDigits int) bool {
	if !ipv4SpecificOptions.Allows_inet_aton_extraneous_digits() {
		return false
	} else if !validationOptions.AllowsIPv6() {
		return true // any number of digits is allowed when IPv6 is ruled out
	} else if !ipv6SpecificOptions.AllowsBase85() {
		return totalDigits < ipv6SingleSegmentDigitCount
	}
	return totalDigits < ipv6Base85SingleSegmentDigitCount
}

type strValidator struct{}

func (strValidator) validateIPAddressStr(fromString *IPAddressString, validationOptions addrstrparam.IPAddressStringParams) (prov ipAddressProvider, err addrerr.AddressStringError) {
	str := fromString.str
	pa := parsedIPAddress{
		originator:         fromString,
		options:            validationOptions,
		ipAddressParseData: ipAddressParseData{addressParseData: addressParseData{str: str}},
	}
	if err = validateIPAddress(validationOptions, str, 0, len(str), pa.getIPAddressParseData(), false); err == nil {
		if err = parseAddressQualifier(str, validationOptions, nil, pa.getIPAddressParseData(), len(str)); err == nil {
			prov, err = chooseIPAddressProvider(fromString, str, validationOptions, &pa)
		} else {
			prov = getInvalidProvider(validationOptions)
		}
	} else {
		prov = getInvalidProvider(validationOptions)
	}
	return
}

func getInvalidProvider(validationOptions addrstrparam.IPAddressStringParams) ipAddressProvider {
	if validationOptions == defaultIPAddrParameters {
		return invalidProvider
	}
	return &nullProvider{isInvalidVal: true, ipType: invalidType, params: validationOptions}
}

func (strValidator) validateMACAddressStr(fromString *MACAddressString, validationOptions addrstrparam.MACAddressStringParams) (prov macAddressProvider, err addrerr.AddressStringError) {
	str := fromString.str
	pa := parsedMACAddress{
		originator:          fromString,
		macAddressParseData: macAddressParseData{addressParseData: addressParseData{str: str}},
		params:              validationOptions,
		creationLock:        &sync.Mutex{},
	}
	if err = validateMACAddress(validationOptions, str, 0, len(str), pa.getMACAddressParseData()); err == nil {
		addressParseData := pa.getAddressParseData()
		prov, err = chooseMACAddressProvider(fromString, validationOptions, &pa, addressParseData)
	} else {
		prov = getInvalidMACProvider(validationOptions)
	}
	if err != nil && prov == nil {
		prov = getInvalidMACProvider(validationOptions)
	}
	return
}

func getInvalidMACProvider(validationOptions addrstrparam.MACAddressStringParams) macAddressProvider {
	if validationOptions == defaultMACAddrParameters {
		return invalidMACProvider
	}
	return macAddressNullProvider{validationOptions}
}

func validateIPAddress(
	validationOptions addrstrparam.IPAddressStringParams,
	str string,
	strStartIndex, strEndIndex int,
	parseData *ipAddressParseData,
	isEmbeddedIPv4 bool) addrerr.AddressStringError {
	return validateAddress(validationOptions, nil, str, strStartIndex, strEndIndex, parseData, nil, isEmbeddedIPv4)
}

func validateMACAddress(
	validationOptions addrstrparam.MACAddressStringParams,
	str string,
	strStartIndex, strEndIndex int,
	parseData *macAddressParseData) addrerr.AddressStringError {
	return validateAddress(nil, validationOptions, str, strStartIndex, strEndIndex, nil, parseData, false)
}

/**
* This method is the mega-parser.
* It is designed to go through the characters one-by-one as a big if/else.
* You have basically several cases: digits, segment separators (. : -), end characters like zone or prefix length,
* range characters denoting a range a-b, wildcard char *, and the 'x' character used to denote hex like 0xf.
*
* Most of the processing occurs in the segment characters, where each segment is analyzed based on what chars came before.
*
* We can parse all possible imaginable variations of mac, ipv4, and ipv6.
*
* This is not the clearest way to write such a parser, because the code for each possible variation is interspersed amongst the various cases,
* so you cannot easily see the code for a given variation clearly, but written this way it may be the fastest parser since we basically account
* for all possibilities simultaneously as we move through the characters just once.
*
 */
func validateAddress(
	validationOptions addrstrparam.IPAddressStringParams,
	macOptions addrstrparam.MACAddressStringParams,
	str string,
	strStartIndex, strEndIndex int,
	ipParseData *ipAddressParseData,
	macParseData *macAddressParseData,
	isEmbeddedIPv4 bool) addrerr.AddressStringError {

	isMac := macParseData != nil

	var parseData *addressParseData
	var stringFormatParams addrstrparam.AddressStringFormatParams
	var ipv6SpecificOptions addrstrparam.IPv6AddressStringParams
	var ipv4SpecificOptions addrstrparam.IPv4AddressStringParams
	var macSpecificOptions addrstrparam.MACAddressStringFormatParams
	var baseOptions addrstrparam.AddressStringParams
	var macFormat macFormat
	canBeBase85 := false
	if isMac {
		baseOptions = macOptions
		macSpecificOptions = macOptions.GetFormatParams()
		stringFormatParams = macSpecificOptions
		macParseData.init(str)
		parseData = macParseData.getAddressParseData()
	} else {
		baseOptions = validationOptions
		// we set stringFormatParams when we know what ip version we have
		ipParseData.init(str)
		parseData = ipParseData.getAddressParseData()
		ipv6SpecificOptions = validationOptions.GetIPv6Params()
		canBeBase85 = ipv6SpecificOptions.AllowsBase85() && validationOptions.AllowsIPv6()
		ipv4SpecificOptions = validationOptions.GetIPv4Params()
	}

	index := strStartIndex

	// per segment variables
	var frontDigitCount, frontLeadingZeroCount, frontSingleWildcardCount, leadingZeroCount,
		singleWildcardCount, wildcardCount, frontWildcardCount int

	var extendedCharacterIndex, extendedRangeWildcardIndex, rangeWildcardIndex,
		hexDelimiterIndex, frontHexDelimiterIndex, segmentStartIndex,
		segmentValueStartIndex = -1, -1, -1, -1, -1, index, index

	var isSegmented, leadingWithZero, hasDigits, frontIsStandardRangeChar, atEnd,
		firstSegmentDashedRange, frontUppercase, uppercase,
		isSingleIPv6, isSingleSegment, isDoubleSegment bool

	var err addrerr.AddressStringError

	checkCharCounts := true

	var version IPVersion
	var currentValueHex, currentFrontValueHex, extendedValue uint64

	charArray := chars

	var currentChar byte
	for {
		if index >= strEndIndex {
			// for ipv6 or ipv4 we set current char to be the right separator to parse the last segment,
			// so that is why we have the check here for index == strEndIndex and not index >= strEndIndex
			atEnd = index == strEndIndex
			if atEnd {
				parseData.setAddressEndIndex(index)
				if isSegmented {
					if isMac {
						currentChar = byte(*macFormat)
						isDoubleSegment = parseData.getSegmentCount() == 1 && currentChar == RangeSeparator
						macParseData.setDoubleSegment(isDoubleSegment)
						if isDoubleSegment {
							totalDigits := index - segmentValueStartIndex
							macParseData.setExtended(totalDigits == macExtendedDoubleSegmentDigitCount)
						}
					} else {
						// we are not base 85, so error if necessary
						if extendedCharacterIndex >= 0 {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
								extendedCharacterIndex}
						}
						//current char is either . or : to handle last segment, unless we have double :: in which case we already handled last segment
						if version.IsIPv4() {
							currentChar = IPv4SegmentSeparator
						} else { //ipv6
							if index == segmentStartIndex {
								if index == parseData.getConsecutiveSeparatorIndex()+2 {
									//ends with ::, we've already parsed the last segment
									break
								}
								return &addressStringError{addressError{str: str, key: "ipaddress.error.cannot.end.with.single.separator"}}
							} else if ipParseData.isProvidingMixedIPv6() {
								//no need to parse the last segment, since it is mixed we already have
								break
							} else {
								currentChar = IPv6SegmentSeparator
							}
						}
					}
				} else {
					// no segment separator so far and segmentCount is 0
					// it could be all addresses like "*", empty "", prefix-only ip address like /64, single segment like 12345, or single segment range like 12345-67890
					totalCharacterCount := index - strStartIndex
					if totalCharacterCount == 0 {
						//it is ""
						if !isMac && ipParseData.hasPrefixSeparator() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.prefix.only"}}
						} else if !baseOptions.AllowsEmpty() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.empty"}}
						}
						parseData.setEmpty(true)
						break
					} else if wildcardCount == totalCharacterCount && wildcardCount <= maxWildcards { //20 wildcards are base 85!
						if !baseOptions.AllowsAll() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.all"}}
						}
						parseData.setHasWildcard()
						parseData.setAll()
						break
					}
					// At this point it is single segment like 12345 or single segment range like 12345-67890
					totalDigits := index - segmentValueStartIndex
					frontTotalDigits := frontLeadingZeroCount + frontDigitCount
					if isMac {
						// we handle the double segment format abcdef-abcdef here
						isDoubleSeg := (totalDigits == macDoubleSegmentDigitCount || totalDigits == macExtendedDoubleSegmentDigitCount) &&
							(frontTotalDigits == macDoubleSegmentDigitCount || frontWildcardCount > 0)
						isDoubleSeg = isDoubleSeg || (frontTotalDigits == macDoubleSegmentDigitCount && wildcardCount > 0)
						isDoubleSeg = isDoubleSeg || (frontWildcardCount > 0 && wildcardCount > 0)
						if isDoubleSeg && !firstSegmentDashedRange { //checks for *-abcdef and abcdef-* and abcdef-abcdef and *-* two segment addresses
							// firstSegmentDashedRange means that the range character is '|'
							addressSize := macOptions.GetPreferredLen()
							if addressSize == addrstrparam.EUI64Len && totalDigits == macDoubleSegmentDigitCount {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments"}}
							} else if addressSize == addrstrparam.MAC48Len && totalDigits == macExtendedDoubleSegmentDigitCount {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.too.many.segments"}}
							}
							// we have aaaaaa-bbbbbb
							if !macOptions.AllowsSingleDashed() {
								return &addressStringError{addressError{str: str, key: "ipaddress.mac.error.format"}}
							}
							isDoubleSegment = true
							macParseData.setDoubleSegment(true)
							macParseData.setExtended(totalDigits == macExtendedDoubleSegmentDigitCount) //we have aaaaaa-bbbbbbbbbb
							currentChar = MACDashSegmentSeparator
							checkCharCounts = false //counted chars already
						} else if frontWildcardCount > 0 || wildcardCount > 0 {
							// either x-* or *-x, we treat these as if they can be expanded to x-*-*-*-*-* or *-*-*-*-*-x
							if !macOptions.AllowsSingleDashed() {
								return &addressStringError{addressError{str: str, key: "ipaddress.mac.error.format"}}
							}
							currentChar = MACDashSegmentSeparator
						} else {
							// a string of digits with no segment separator
							// here we handle abcdefabcdef or abcdefabcdef|abcdefabcdef or abcdefabcdef-abcdefabcdef
							if !baseOptions.AllowsSingleSegment() {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.single.segment"}}
							}
							is12Digits := totalDigits == macSingleSegmentDigitCount
							is16Digits := totalDigits == macExtendedSingleSegmentDigitCount
							isNoDigits := totalDigits == 0
							if is12Digits || is16Digits || isNoDigits {
								var frontIs12Digits, frontIs16Digits, frontIsNoDigits bool
								if rangeWildcardIndex >= 0 {
									frontIs12Digits = frontTotalDigits == macSingleSegmentDigitCount
									frontIs16Digits = frontTotalDigits == macExtendedSingleSegmentDigitCount
									frontIsNoDigits = frontTotalDigits == 0
									if is12Digits {
										if !frontIs12Digits && !frontIsNoDigits {
											return &addressStringError{addressError{str: str, key: "ipaddress.error.front.digit.count"}}
										}
									} else if is16Digits {
										if !frontIs16Digits && !frontIsNoDigits {
											return &addressStringError{addressError{str: str, key: "ipaddress.error.front.digit.count"}}
										}
									} else if isNoDigits {
										if !frontIs12Digits && !frontIs16Digits {
											return &addressStringError{addressError{str: str, key: "ipaddress.error.front.digit.count"}}
										}
									}
								} else if isNoDigits {
									return &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments.digit.count"}}
								}
								isSingleSegment = true
								parseData.setSingleSegment()
								macParseData.setExtended(is16Digits || frontIs16Digits)
								currentChar = MACColonSegmentSeparator
								checkCharCounts = false //counted chars already
							} else {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments.digit.count"}}
							}
						}
					} else {
						//a string of digits with no segment separator
						if !baseOptions.AllowsSingleSegment() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.single.segment"}}
						}

						if canBeBase85 {
							if canBeBase85, base85Err := parseBase85(
								validationOptions, str, strStartIndex, strEndIndex, ipParseData,
								extendedRangeWildcardIndex, totalCharacterCount, index); canBeBase85 {
								break
							} else if base85Err != nil {
								// base85Err can only be non-nil when returning false
								return base85Err
							}
						}
						// we are not base 85, so error if necessary
						if extendedCharacterIndex >= 0 {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
								extendedCharacterIndex,
							}
						}

						isRange := rangeWildcardIndex >= 0
						if isSingleSeg, serr := isSingleSegmentIPv6(
							str,
							totalDigits,
							isRange,
							frontTotalDigits,
							ipv6SpecificOptions); serr != nil {
							return serr
						} else if isSingleSeg {
							isSingleIPv6 = true
							currentChar = IPv6SegmentSeparator
						} else if validationOptions.AllowsIPv4() {

							leadingZeros := leadingZeroCount
							if leadingWithZero {
								leadingZeros++
							}

							if isSingleSeg, serr = isSingleSegmentIPv4(
								str,
								totalDigits-leadingZeros,
								totalDigits,
								isRange,
								frontDigitCount,
								frontTotalDigits,
								ipv4SpecificOptions); serr != nil {
								return serr
							} else if isSingleSeg {
								currentChar = IPv4SegmentSeparator
							} else if hasExtraneousDigitsIPv4(validationOptions, ipv4SpecificOptions, ipv6SpecificOptions, totalDigits) {
								if singleWildcardCount > 0 || wildcardCount > 0 {
									return &addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character"}}
								} else if isRange {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
										rangeWildcardIndex}
								} else if ipParseData.getQualifierIndex() >= 0 {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
										ipParseData.getQualifierIndex()}

								}
								var val, radix uint32
								if hexDelimiterIndex >= 0 {
									if !ipv4SpecificOptions.AllowsLeadingZeros() {
										// the '0' preceding the 'x' is not allowed
										return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
									} else if !ipv4SpecificOptions.Allows_inet_aton_hex() {
										return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv4.segment.hex"}}
									} else if leadingZeros > 1 && !ipv4SpecificOptions.Allows_inet_aton_leading_zeros() {
										// the '0' following the 'x' is not allowed
										return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
									}
									radix = 16
									val = parseInt16(str, strEndIndex-8, strEndIndex)
								} else if leadingZeros > 0 && ipv4SpecificOptions.Allows_inet_aton_octal() {
									if !ipv4SpecificOptions.AllowsLeadingZeros() {
										return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
									} else if leadingZeros > 1 && !ipv4SpecificOptions.Allows_inet_aton_leading_zeros() {
										return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
									}
									radix = 8
									if val, serr = parseInt8(str, strStartIndex, strEndIndex); serr != nil {
										return serr
									}
								} else {
									radix = 10
									if val, serr = parseInt10(str, strStartIndex, strEndIndex); serr != nil {
										return serr
									}
								}

								// We cannot allow normal parsing for decimal, because the remainder (mod of max uint32) does not work with the way we shift-store values hex values currentValueHex.
								// We cannot allow normal parsing for octal.  This is because, even though the way we shift-store hex digits in currentValueHex works for octal, giving us the correct mod of max uint32,
								// it does not allow us to check if overflow chars are octal.
								// For hex, maybe we could use normal parsing, but we'd need to check for range chars here, we need to disallow them, and the code as written does not.
								// So far all radices, we parse extraneous digits here.

								ipParseData.setVersion(IPv4)
								ipParseData.set_has_inet_aton_value(true)
								parseData.initSegmentData(1)
								parseData.incrementSegmentCount()
								val64 := uint64(val)
								assign3Attributes1Values1Flags(strStartIndex+leadingZeros, strEndIndex, strStartIndex, parseData, 0, val64, radix)
								parseData.setSingleSegment()
								break
							} else {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments.digit.count"}}
							}
						} else {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments.digit.count"}}
						}

						// single segment IPv4 or single-segment IPv6
						isSingleSegment = true
						parseData.setSingleSegment()
						checkCharCounts = false // counted chars already
					}
				}
			} else {
				break
			}
		} else {
			currentChar = str[index]
		}

		// evaluate the character
		if currentChar <= '9' && currentChar >= '0' {
			if hasDigits {
				currentValueHex = currentValueHex<<4 | uint64(charArray[currentChar])
			} else {
				if currentChar == '0' {
					if leadingWithZero {
						leadingZeroCount++
					} else {
						leadingWithZero = true
					}
				} else {
					hasDigits = true
					currentValueHex = currentValueHex<<4 | uint64(charArray[currentChar])
				}
			}
			index++
		} else if currentChar >= 'a' && currentChar <= 'f' {
			currentValueHex = currentValueHex<<4 | uint64(charArray[currentChar])
			hasDigits = true
			index++
		} else if currentChar == IPv4SegmentSeparator {
			segCount := parseData.getSegmentCount()
			// could be mac or ipv4, we handle either one
			if isMac {
				if segCount == 0 {
					if !macOptions.AllowsDotted() {
						return &addressStringError{addressError{str: str, key: "ipaddress.mac.error.format"}}
					}
					macFormat = dotted
					macParseData.setFormat(macFormat)
					parseData.initSegmentData(MediaAccessControlDotted64SegmentCount)
					isSegmented = true
				} else {
					if macFormat != dotted {
						return &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.mac.error.mix.format.characters.at.index"}},
							index}
					}
					var limit int
					if macOptions.GetPreferredLen() == addrstrparam.MAC48Len {
						limit = MediaAccessControlDottedSegmentCount
					} else {
						limit = MediaAccessControlDotted64SegmentCount
					}
					if segCount >= limit {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.too.many.segments"}}
					}
				}
			} else {
				//end of an ipv4 segment
				if segCount == 0 {
					if !validationOptions.AllowsIPv4() {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv4"}}
					}
					version = IPv4
					ipParseData.setVersion(version)
					stringFormatParams = ipv4SpecificOptions
					canBeBase85 = false
					parseData.initSegmentData(IPv4SegmentCount)
					isSegmented = true
				} else if ipParseData.getProviderIPVersion().IsIPv6() {
					//mixed IPv6 address like 1:2:3:4:5:6:1.2.3.4
					if !ipv6SpecificOptions.AllowsMixed() {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.no.mixed"}}
					}
					totalSegmentCount := segCount + IPv6MixedReplacedSegmentCount
					if totalSegmentCount > IPv6SegmentCount {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.too.many.segments"}}
					}
					if wildcardCount > 0 {
						if parseData.getConsecutiveSeparatorIndex() < 0 &&
							totalSegmentCount < IPv6SegmentCount &&
							ipv6SpecificOptions.AllowsWildcardedSeparator() {
							// the '*' is covering an additional ipv6 segment (eg 1:2:3:4:5:*.2.3.4, the * covers both an ipv4 and ipv6 segment)
							// we flag this IPv6 segment with keyMergedMixed
							parseData.setHasWildcard()
							assign6Attributes2Values1Flags(segmentStartIndex, index, segmentStartIndex, segmentStartIndex, index, segmentStartIndex,
								parseData, segCount, 0, IPv6MaxValuePerSegment, keyWildcard|keyMergedMixed)
							parseData.incrementSegmentCount()
						}
					}
					mixedOptions := ipv6SpecificOptions.GetMixedParams()
					pa := &parsedIPAddress{
						ipAddressParseData: ipAddressParseData{addressParseData: addressParseData{str: str}},
						options:            mixedOptions,
					}
					err = validateIPAddress(mixedOptions, str, segmentStartIndex, strEndIndex, &pa.ipAddressParseData, true)
					if err != nil {
						return err
					}
					pa.clearQualifier()
					err = checkSegments(str, mixedOptions, pa.getIPAddressParseData())
					if err != nil {
						return err
					}
					ipParseData.setMixedParsedAddress(pa)
					index = pa.getAddressParseData().getAddressEndIndex()
					continue
				} else if segCount >= IPv4SegmentCount {
					return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv4.too.many.segments"}}
				}
			}
			if wildcardCount > 0 {
				if !stringFormatParams.GetRangeParams().AllowsWildcard() {
					return &addressStringError{addressError{str: str, key: "ipaddress.error.no.wildcard"}}
				}
				//wildcards must appear alone
				totalDigits := index - segmentStartIndex
				if wildcardCount != totalDigits || hexDelimiterIndex >= 0 {
					return &addressStringIndexError{
						addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
						index}
				}
				parseData.setHasWildcard()
				startIndex := index - wildcardCount
				var max uint64
				if isMac {
					max = MACMaxValuePerDottedSegment
				} else {
					max = IPv4MaxValuePerSegment
				}
				assign6Attributes2Values1Flags(startIndex, index, startIndex, startIndex, index, startIndex,
					parseData, segCount, 0, max, keyWildcard)
				wildcardCount = 0
			} else {
				var flags, rangeFlags, radix uint32
				var value uint64
				digitStartIndex := segmentValueStartIndex + leadingZeroCount
				digitCount := index - digitStartIndex
				if leadingWithZero {
					if digitCount == 1 {
						if leadingZeroCount == 0 && rangeWildcardIndex < 0 && hexDelimiterIndex < 0 {
							// handles 0, but not 1-0 or 0x0
							assign4Attributes(digitStartIndex, index, parseData, segCount, 10, segmentValueStartIndex)
							parseData.incrementSegmentCount()
							index++
							segmentStartIndex = index
							segmentValueStartIndex = index
							leadingWithZero = false
							continue
						}
					} else {
						leadingZeroCount++
						digitStartIndex++
						digitCount--
					}
					leadingWithZero = false // reset this flag now that we've used it
				}
				noValuesToSet := false
				if digitCount == 0 {
					// we allow an empty range boundary to denote the max value
					if rangeWildcardIndex < 0 || hexDelimiterIndex >= 0 || !stringFormatParams.GetRangeParams().AllowsInferredBoundary() {
						// starts with '.', or has two consecutive '.'
						return &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
							index}
					} else if isMac {
						value = MACMaxValuePerDottedSegment
						radix = MACDefaultTextualRadix
					} else {
						value = IPv4MaxValuePerSegment // for inet-aton multi-segment, this will be adjusted later
						radix = IPv4DefaultTextualRadix
					}
					rangeFlags = keyInferredUpperBoundary
				} else { // digitCount > 0
					// Note: we cannot do max value check on ipv4 until after all segments have been read due to inet_aton joined segments,
					// although we can do a preliminary check here that is in fact needed to prevent overflow when calculating values later
					isBinary := false
					hasLeadingZeros := leadingZeroCount > 0
					isSingleWildcard := singleWildcardCount > 0
					if isMac || hexDelimiterIndex >= 0 {
						if isMac { // mac dotted segments aabb.ccdd.eeff
							maxMacChars := 4
							if digitCount > maxMacChars { //
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
									segmentValueStartIndex}
							}
							totalDigits := digitCount + leadingZeroCount
							if hexDelimiterIndex >= 0 {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
									hexDelimiterIndex}
							} else if leadingZeroCount > 0 && !stringFormatParams.AllowsLeadingZeros() {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
							} else if !stringFormatParams.AllowsUnlimitedLeadingZeros() && totalDigits > maxMacChars {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
									segmentValueStartIndex}
							} else if !macSpecificOptions.AllowsShortSegments() && totalDigits < maxMacChars {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.short.at.index"}},
									segmentValueStartIndex}
							}
						} else if !stringFormatParams.AllowsLeadingZeros() {
							// the '0' preceding the 'x' is not allowed
							return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
						} else if !ipv4SpecificOptions.Allows_inet_aton_hex() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv4.segment.hex"}}
						} else if hasLeadingZeros && !ipv4SpecificOptions.Allows_inet_aton_leading_zeros() {
							// the '0' following the 'x' is not allowed
							return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
						} else {
							if digitCount > 8 { // 0xffffffff
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
									segmentValueStartIndex}
							}
							ipParseData.set_has_inet_aton_value(true)
						}
						radix = 16
						if isSingleWildcard {
							if rangeWildcardIndex >= 0 {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
									index}
							}
							err = assignSingleWildcard16(currentValueHex, str, digitStartIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams)
							if err != nil {
								return err
							}
							value = 0
							noValuesToSet = true
							singleWildcardCount = 0
						} else {
							value = currentValueHex
						}
						hexDelimiterIndex = -1
					} else {
						isBinaryOrOctal := hasLeadingZeros
						if isBinaryOrOctal {
							isBinary = ipv4SpecificOptions.AllowsBinary() && isBinaryDelimiter(str, digitStartIndex)
							isBinaryOrOctal = isBinary || ipv4SpecificOptions.Allows_inet_aton_octal()
						}
						if isBinaryOrOctal {
							if !stringFormatParams.AllowsLeadingZeros() {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
							}
							if isBinary {
								if digitCount > 33 {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
										segmentValueStartIndex}
								}
								digitStartIndex++ // exclude the 'b' in 0b1100
								digitCount--      // exclude the 'b'
								radix = 2
								ipParseData.setHasBinaryDigits(true)
								if isSingleWildcard {
									if rangeWildcardIndex >= 0 {
										return &addressStringIndexError{
											addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
											index}
									}
									if digitCount > 16 {
										if parseErr := parseSingleSegmentSingleWildcard2(str, digitStartIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
											return parseErr
										}
									} else {
										if parseErr := switchSingleWildcard2(currentValueHex, str, digitStartIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
											return parseErr
										}
									}
									value = 0
									noValuesToSet = true
									singleWildcardCount = 0
								} else {
									if digitCount > 16 {
										value = parseLong2(str, digitStartIndex, index)
									} else {
										value, err = switchValue2(currentValueHex, str, digitCount)
										if err != nil {
											return err
										}
									}
								}
							} else {
								if leadingZeroCount > 1 && !ipv4SpecificOptions.Allows_inet_aton_leading_zeros() {
									return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
								} else if digitCount > 11 { //octal 037777777777
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
										segmentValueStartIndex}
								}
								ipParseData.set_has_inet_aton_value(true)
								radix = 8
								if isSingleWildcard {
									if rangeWildcardIndex >= 0 {
										return &addressStringIndexError{
											addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
											index}
									}
									if parseErr := switchSingleWildcard8(currentValueHex, str, digitStartIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
										return parseErr
									}
									value = 0
									noValuesToSet = true
									singleWildcardCount = 0
								} else {
									value, err = switchValue8(currentValueHex, str, digitCount)
									if err != nil {
										return err
									}
								}
							}
						} else {
							if hasLeadingZeros {
								if !stringFormatParams.AllowsLeadingZeros() {
									return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
								}
								ipParseData.setHasIPv4LeadingZeros(true)
							}
							if digitCount > 10 { // 4294967295
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
									segmentValueStartIndex}
							}
							radix = 10
							if isSingleWildcard {
								if rangeWildcardIndex >= 0 {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
										index}
								}
								if parseErr := switchSingleWildcard10(currentValueHex, str, digitStartIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, ipv4SpecificOptions); parseErr != nil {
									return parseErr
								}
								value = 0
								noValuesToSet = true
								singleWildcardCount = 0
							} else {
								value, err = switchValue10(currentValueHex, str, digitCount)
								if err != nil {
									return err
								}
								flags = keyStandardStr
							}
						}
					}
					hasDigits = false
					currentValueHex = 0
				}
				if rangeWildcardIndex >= 0 {
					var frontRadix uint32
					var front uint64
					frontStartIndex := rangeWildcardIndex - frontDigitCount
					frontEndIndex := rangeWildcardIndex
					frontLeadingZeroStartIndex := frontStartIndex - frontLeadingZeroCount
					if !stringFormatParams.GetRangeParams().AllowsRangeSeparator() {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.no.range"}}
					} else if frontSingleWildcardCount > 0 || frontWildcardCount > 0 { //no wildcards in ranges
						return &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
							rangeWildcardIndex}
					}
					frontEmpty := frontStartIndex == frontEndIndex
					isReversed := false
					hasFrontLeadingZeros := frontLeadingZeroCount > 0
					if isMac || frontHexDelimiterIndex >= 0 {
						if isMac {
							totalFrontDigits := frontDigitCount + frontLeadingZeroCount
							maxMacChars := 4
							if frontHexDelimiterIndex >= 0 {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
									frontHexDelimiterIndex}
							} else if hasFrontLeadingZeros && !stringFormatParams.AllowsLeadingZeros() {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
							} else if !stringFormatParams.AllowsUnlimitedLeadingZeros() && totalFrontDigits > maxMacChars {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
									frontLeadingZeroStartIndex}
							} else if !macSpecificOptions.AllowsShortSegments() && totalFrontDigits < maxMacChars {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.short.at.index"}},
									frontLeadingZeroStartIndex}
							} else if frontEmpty { //we allow the front of a range to be empty in which case it is 0
								if !stringFormatParams.GetRangeParams().AllowsInferredBoundary() {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
										index}
								}
								rangeFlags |= keyInferredLowerBoundary
								front = 0
							} else if frontDigitCount > maxMacChars { // mac dotted segments aaa.bbb.ccc.ddd
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
									frontLeadingZeroStartIndex}
							} else {
								front = currentFrontValueHex
								isReversed = front > value && digitCount != 0
							}
						} else if !stringFormatParams.AllowsLeadingZeros() {
							// the '0' preceding the 'x' is not allowed
							return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
						} else if !ipv4SpecificOptions.Allows_inet_aton_hex() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv4.segment.hex"}}
						} else if hasFrontLeadingZeros && !ipv4SpecificOptions.Allows_inet_aton_leading_zeros() {
							// the '0' following the 'x' is not allowed
							return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
						} else if frontEmpty {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
								index}
						} else if frontDigitCount > 8 { // 0xffffffff
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
								frontLeadingZeroStartIndex}
						} else {
							ipParseData.set_has_inet_aton_value(true)
							front = currentFrontValueHex
							isReversed = front > value && digitCount != 0
						}
						frontRadix = 16
					} else {
						if hasFrontLeadingZeros {
							if !stringFormatParams.AllowsLeadingZeros() {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
							}
							if ipv4SpecificOptions.AllowsBinary() && isBinaryDelimiter(str, frontStartIndex) {
								if frontDigitCount > 33 {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
										frontLeadingZeroStartIndex}
								}
								ipParseData.setHasBinaryDigits(true)
								frontStartIndex++
								frontDigitCount--
								if frontDigitCount > 16 {
									front = parseLong2(str, frontStartIndex, frontEndIndex)
								} else {
									front, err = switchValue2(currentFrontValueHex, str, frontDigitCount)
									if err != nil {
										return err
									}
								}
								frontRadix = 2
								isReversed = digitCount != 0 && front > value
							} else if ipv4SpecificOptions.Allows_inet_aton_octal() {
								if frontLeadingZeroCount > 1 && !ipv4SpecificOptions.Allows_inet_aton_leading_zeros() {
									return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
								} else if frontDigitCount > 11 { // 037777777777
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
										frontLeadingZeroStartIndex}
								}
								ipParseData.set_has_inet_aton_value(true)
								front, err = switchValue8(currentFrontValueHex, str, frontDigitCount)
								if err != nil {
									return err
								}
								frontRadix = 8
								isReversed = digitCount != 0 && front > value
							}
						}
						if frontRadix == 0 {
							frontRadix = 10
							if frontEmpty { //we allow the front of a range to be empty in which case it is 0
								if !stringFormatParams.GetRangeParams().AllowsInferredBoundary() {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
										index}
								}
								rangeFlags |= keyInferredLowerBoundary
							} else if frontDigitCount > 10 { // 4294967295
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
									frontLeadingZeroStartIndex}
							} else {
								front, err = switchValue10(currentFrontValueHex, str, frontDigitCount)
								if err != nil {
									return err
								}
								if hasFrontLeadingZeros {
									if !stringFormatParams.AllowsLeadingZeros() {
										return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
									}
									ipParseData.setHasIPv4LeadingZeros(true)
								}
								isReversed = digitCount != 0 && front > value
								if !isReversed {
									if leadingZeroCount == 0 && (flags&keyStandardStr) != 0 {
										rangeFlags |= keyStandardRangeStr | keyStandardStr
									} else {
										rangeFlags |= keyStandardStr
									}
								}
							}
						}
					}
					backEndIndex := index
					if isReversed {
						if !stringFormatParams.GetRangeParams().AllowsReverseRange() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.invalidRange"}}
						}
						// switcheroo
						frontStartIndex, digitStartIndex = digitStartIndex, frontStartIndex
						frontEndIndex, backEndIndex = backEndIndex, frontEndIndex
						frontLeadingZeroStartIndex, segmentValueStartIndex = segmentValueStartIndex, frontLeadingZeroStartIndex
						frontRadix, radix = radix, frontRadix
						front, value = value, front
					}
					assign6Attributes2Values2Flags(frontStartIndex, frontEndIndex, frontLeadingZeroStartIndex, digitStartIndex, backEndIndex, segmentValueStartIndex,
						parseData, segCount, front, value, rangeFlags|keyRangeWildcard|frontRadix, radix)
					rangeWildcardIndex = -1
				} else if !noValuesToSet {
					assign3Attributes1Values1Flags(digitStartIndex, index, segmentValueStartIndex, parseData, segCount, value, flags|radix)
				}
				leadingZeroCount = 0
			}
			parseData.incrementSegmentCount()
			index++
			segmentValueStartIndex = index
			segmentStartIndex = index
			// end of IPv4 segments and mac segments with '.' separators
		} else {
			//checking for all IPv6 and MAC48Len segments, as well as the front range of all segments IPv4, IPv6, and MAC48Len
			//the range character '-' is the same as one of the separators '-' for MAC48Len,
			//so further work is required to distinguish between the front of IPv6/IPv4/MAC48Len range and MAC48Len segment
			//we also handle IPv6 segment and MAC48Len segment in the same place to avoid code duplication
			var isSpace, isDashedRangeChar, isRangeChar bool
			if currentChar == IPv6SegmentSeparator {
				isRangeChar = false
				isSpace = false
			} else {
				isRangeChar = currentChar == RangeSeparator
				if isRangeChar || (isMac && (currentChar == MacDashedSegmentRangeSeparator)) {
					isSpace = false
					isDashedRangeChar = !isRangeChar

					/*
					 There are 3 cases here, A, B and C.
					 A - we have two MAC48Len segments a-b-
					 B - we have the front of a range segment, either a-b which is MAC48Len or ipv6AddrType,  or a|b or a<space>b which is MAC48Len
					 C - we have a single segment, either a MAC48Len segment a- or an IPv6 or MAC48Len segment a:
					*/

					/*
					 Here we have either a '-' or '|' character or a space ' '

					 If we have a '-' character:
					 For MAC48Len address, the cases are:
					 1. we did not previously set macFormat and we did not previously encounter '|'
					 		-if rangeWildcardIndex >= 0 we have dashed a-b- we treat as two segments, case A (we cannot have a|b because that would have set macFormat previously)
					 		-if rangeWildcardIndex < 0, we treat as front of range, case B, later we will know for sure if really front of range
					 2. we previously set macFormat or we previously encountered '|'
					 		if set to dashed we treat as one segment, may or may not be range segment, case C
					 		if we previously encountered '|' we treat as dashed range segment, case C
					 		if not we treat as front of range, case B

					 For IPv6, this is always front of range, case B

					 If we have a '|' character, we have front of range MAC48Len, case B
					*/
					// we know either isRangeChar or isDashedRangeChar is true at this point
					endOfHexSegment := false
					if isMac {
						if macFormat == nil {
							if rangeWildcardIndex >= 0 && !firstSegmentDashedRange {

								//case A, we have two segments a-b- or a-b|

								//we handle the first segment here, we handle the second segment in the usual place below
								if frontHexDelimiterIndex >= 0 {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
										frontHexDelimiterIndex}
								} else if hexDelimiterIndex >= 0 {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
										hexDelimiterIndex}
								} else if !macOptions.AllowsDashed() {
									return &addressStringError{addressError{str: str, key: "ipaddress.mac.error.format"}}
								}
								macFormat = dashed
								macParseData.setFormat(macFormat)
								checkCharCounts = false //counting chars later
								parseData.initSegmentData(ExtendedUniqueIdentifier64SegmentCount)
								isSegmented = true
								if frontWildcardCount > 0 {
									if !stringFormatParams.GetRangeParams().AllowsWildcard() {
										return &addressStringError{addressError{str: str, key: "ipaddress.error.no.wildcard"}}
									} else if frontSingleWildcardCount > 0 || frontLeadingZeroCount > 0 || frontDigitCount > 0 || frontHexDelimiterIndex >= 0 { //wildcards must appear alone
										return &addressStringIndexError{
											addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
											rangeWildcardIndex}
									}
									parseData.setHasWildcard()
									backDigits := index - segmentValueStartIndex
									var upperValue uint64
									if isDoubleSegment || backDigits == macDoubleSegmentDigitCount {
										//even when not already identified as a double segment address, which is something we can see
										//only when we reach the end of the address, we may have a-b| where a is * and b is a 6 digit value.
										//Here we are considering the max value of a.
										//If b is 6 digits, we need to consider the max value of * as if we know already it will be double segment.
										//We can do this because the max values will be checked after the address has been parsed,
										//so even if a-b| ends up being a full address a-b|c-d-e-f-a and not a-b|c,
										//the fact that we have 6 digits here will invalidate the first address,
										//so we can safely assume that this address must be a double segment a-b|c even before we have seen that.
										upperValue = macMaxTriple
									} else {
										upperValue = MACMaxValuePerSegment
									}
									startIndex := rangeWildcardIndex - frontWildcardCount
									assign6Attributes2Values1Flags(startIndex, rangeWildcardIndex, startIndex, startIndex, rangeWildcardIndex, startIndex,
										parseData, 0, 0, upperValue, keyWildcard)
								} else {
									if !stringFormatParams.AllowsLeadingZeros() && frontLeadingZeroCount > 0 {
										return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
									}
									startIndex := rangeWildcardIndex - frontDigitCount
									leadingZeroStartIndex := startIndex - frontLeadingZeroCount
									if frontSingleWildcardCount > 0 {
										if parseErr := assignSingleWildcard16(currentFrontValueHex, str, startIndex, rangeWildcardIndex, singleWildcardCount, parseData, 0, leadingZeroStartIndex, stringFormatParams); parseErr != nil {
											return parseErr
										}
									} else {
										var flags uint32
										if !frontUppercase {
											if frontDigitCount == 0 {
												return &addressStringIndexError{
													addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
													startIndex}
											}
											flags = keyStandardStr
										}
										assign3Attributes1Values1Flags(startIndex, rangeWildcardIndex, leadingZeroStartIndex, parseData, 0, currentFrontValueHex, flags)
									}
								}
								segmentValueStartIndex = rangeWildcardIndex + 1
								segmentStartIndex = segmentValueStartIndex
								rangeWildcardIndex = -1
								parseData.incrementSegmentCount()
								//end of handling the first segment a- in a-b-
								//below we handle b- by setting endOfSegment here
								endOfHexSegment = isRangeChar
							} else { //we will treat this as the front of a range
								if isDashedRangeChar {
									firstSegmentDashedRange = true
								} else {
									endOfHexSegment = firstSegmentDashedRange
								}
							}
						} else {
							if macFormat == dashed {
								endOfHexSegment = isRangeChar
							} else if isDashedRangeChar {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
									index}
							}
						}
					}
					if !endOfHexSegment {
						if extendedCharacterIndex < 0 {
							//case B
							if rangeWildcardIndex >= 0 {
								if canBeBase85 {
									index++
									extendedCharacterIndex = index
								} else {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
										index}
								}
							} else {
								//here is where we handle the front 'a' of a range like 'a-b'
								rangeWildcardIndex = index

								frontIsStandardRangeChar = isRangeChar
								frontDigitCount = ((index - segmentValueStartIndex) - leadingZeroCount) - wildcardCount
								frontLeadingZeroCount = leadingZeroCount
								if leadingWithZero {
									if frontDigitCount != 1 {
										frontLeadingZeroCount++
										frontDigitCount--
									}
								}
								frontUppercase = uppercase
								frontHexDelimiterIndex = hexDelimiterIndex
								frontWildcardCount = wildcardCount
								frontSingleWildcardCount = singleWildcardCount
								currentFrontValueHex = currentValueHex

								index++
								segmentValueStartIndex = index
								hasDigits = false
								uppercase = false
								leadingWithZero = false
								hexDelimiterIndex = -1
								leadingZeroCount = 0
								wildcardCount = 0
								singleWildcardCount = 0
								currentValueHex = 0
							}
						} else {
							index++
						}
						continue
					}
				} else if isMac && currentChar == space {
					isSpace = true
				} else {
					// other characters handled here
					isZoneChar := false
					if currentChar == PrefixLenSeparator {
						if isMac {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
								index}
						}
						strEndIndex = index
						ipParseData.setHasPrefixSeparator(true)
						ipParseData.setQualifierIndex(index + 1)
					} else if currentChar >= 'A' && currentChar <= 'F' { // this is not paired with 'a' to 'f' because these are not canonical and hence not part of the fast path
						index++
						currentValueHex = (currentValueHex << 4) | uint64(charArray[currentChar])
						hasDigits = true
						uppercase = true
					} else {
						isSegWildcard := currentChar == SegmentWildcard
						if !isSegWildcard {
							isZoneChar = currentChar == SegmentSqlWildcard
							isSegWildcard = isZoneChar
						}
						if isSegWildcard {
							//the character * is always treated as wildcard (but later can be determined to be a base 85 character)

							//the character % denotes a zone and is also a character for the SQL wildcard,
							//and it is also a base 85 character,
							//so we treat it as zone only if the options allow it and it is in the zone position.
							//Either we have seen an ipv6 segment separator, or we are at the end of the correct number of digits for ipv6 single segment (which rules out base 85 or ipv4 single segment),
							//or we are the '*' all wildcard so far which can represent everything including ipv6
							//
							//In all other cases, the character is treated as wildcard,
							//but as is the case of other characters we may later discover we are base 85 ipv6
							//For base 85, we decided that having the same character mean two different thing depending on position in the string, that is not reasonable.
							//In fact, if the zone character were allowed, can you tell if there is a zone here or not: %%%%%%%%%%%%%%%%%%%%%%

							canBeZone := isZoneChar &&
								!isMac &&
								ipv6SpecificOptions.AllowsZone()
							if canBeZone {
								isIPv6 := parseData.getSegmentCount() > 0 && (isEmbeddedIPv4 || ipParseData.getProviderIPVersion() == IPv6) /* at end of IPv6 regular or mixed */
								if !isIPv6 {
									isIPv6, _ = isSingleSegmentIPv6(str, index-segmentValueStartIndex, rangeWildcardIndex >= 0, frontLeadingZeroCount+frontDigitCount, ipv6SpecificOptions)
									if !isIPv6 {
										isIPv6 = wildcardCount == index && wildcardCount <= maxWildcards /* all wildcards so far */
									}
								}
								canBeZone = isIPv6
							}
							if canBeZone {
								//we are not base 85
								canBeBase85 = false
								strEndIndex = index
								ipParseData.setZoned(true)
								ipParseData.setQualifierIndex(index + 1)
							} else {
								wildcardCount++
								index++
							}
						} else if currentChar == SegmentSqlSingleWildcard {
							hasDigits = true
							index++
							singleWildcardCount++
						} else if isHexDelimiter(currentChar) {
							if hasDigits || !leadingWithZero || leadingZeroCount > 0 || hexDelimiterIndex >= 0 || singleWildcardCount > 0 {
								if canBeBase85 {
									if extendedCharacterIndex < 0 {
										extendedCharacterIndex = index
									}
									index++
								} else {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
										index}
								}
							} else {
								if isMac {
									if parseData.getSegmentCount() > 0 {
										return &addressStringIndexError{
											addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
											index}
									}
								} else if version.IsIPv6() {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
										index}
								}
								hexDelimiterIndex = index
								leadingWithZero = false
								index++
								segmentValueStartIndex = index
							}
							//the remaining possibilities are base85 only
						} else if currentChar == AlternativeRangeSeparatorStr[0] {
							if index+1 == strEndIndex {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
									index}
							}
							currentChar = str[index+1]
							if currentChar == AlternativeRangeSeparatorStr[1] {
								if canBeBase85 {
									if extendedCharacterIndex < 0 {
										extendedCharacterIndex = index
									} else if extendedRangeWildcardIndex >= 0 {
										return &addressStringIndexError{
											addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
											index}
									}
									extendedRangeWildcardIndex = index
								} else {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
										index}
								}
								index += 2
							} else if currentChar == IPv6AlternativeZoneSeparatorStr[1] {
								if canBeBase85 && !isMac && ipv6SpecificOptions.AllowsZone() {
									strEndIndex = index
									ipParseData.setZoned(true)
									ipParseData.setBase85Zoned(true)
									ipParseData.setQualifierIndex(index + 2)
								} else {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
										index}
								}
							} else {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
									index - 1}
							}
						} else {
							if canBeBase85 {
								if currentChar < 0 || int(currentChar) >= len(extendedChars) {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
										index}
								}
								val := extendedChars[currentChar]
								if val == 0 { //note that we already check for the currentChar '0' character at another else/if block, so any other character mapped to the value 0 is an invalid character
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
										index}
								} else if extendedCharacterIndex < 0 {
									extendedCharacterIndex = index
								}
							} else {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
									index}
							}
							index++
						}
					}
					continue
				}
			}
			// ipv6 and mac segments handled here
			segCount := parseData.getSegmentCount()
			var hexMaxChars int
			if isMac {
				if segCount == 0 {
					if isSingleSegment {
						parseData.initSegmentData(1)
					} else {
						if hexDelimiterIndex >= 0 {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
								hexDelimiterIndex}
						} else {
							var isNoExc bool
							if isRangeChar {
								isNoExc = macOptions.AllowsDashed()
							} else if isSpace {
								isNoExc = macOptions.AllowsSpaceDelimited()
							} else {
								isNoExc = macOptions.AllowsColonDelimited()
							}
							if !isNoExc {
								return &addressStringError{addressError{str: str, key: "ipaddress.mac.error.format"}}
							} else if isRangeChar {
								macFormat = dashed
								macParseData.setFormat(macFormat)
								checkCharCounts = false //counting chars later
							} else {
								if isSpace {
									macFormat = spaceDelimited
								} else {
									macFormat = colonDelimited
								}
								macParseData.setFormat(macFormat)
							}
						}
						parseData.initSegmentData(ExtendedUniqueIdentifier64SegmentCount)
						isSegmented = true
					}
				} else {
					isExc := false
					if isRangeChar {
						isExc = macFormat != dashed
					} else if isSpace {
						isExc = macFormat != spaceDelimited
					} else {
						isExc = macFormat != colonDelimited
					}
					if isExc {
						return &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.mac.error.mix.format.characters.at.index"}},
							index}
					}
					var segLimit int
					if macOptions.GetPreferredLen() == addrstrparam.MAC48Len {
						segLimit = MediaAccessControlSegmentCount
					} else {
						segLimit = ExtendedUniqueIdentifier64SegmentCount
					}
					if segCount >= segLimit {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.too.many.segments"}}
					}
				}
				hexMaxChars = MACSegmentMaxChars //will be ignored for single or double segments due to checkCharCounts booleans
			} else {
				if segCount == 0 {
					if !validationOptions.AllowsIPv6() {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv6"}}
					}
					canBeBase85 = false
					version = IPv6
					ipParseData.setVersion(version)
					stringFormatParams = ipv6SpecificOptions
					isSegmented = true

					if index == strStartIndex {
						firstIndex := index
						index++
						if index == strEndIndex {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.too.few.segments"}}
						} else if str[index] != IPv6SegmentSeparator {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv6.cannot.start.with.single.separator"}}
						}
						parseData.initSegmentData(IPv6SegmentCount)
						parseData.setConsecutiveSeparatorSegmentIndex(0)
						parseData.setConsecutiveSeparatorIndex(firstIndex)
						assign3Attributes(index, index, parseData, 0, index)
						parseData.incrementSegmentCount()
						index++
						segmentValueStartIndex = index
						segmentStartIndex = segmentValueStartIndex
						continue
					} else {
						if isSingleSegment {
							parseData.initSegmentData(1)
						} else {
							if hexDelimiterIndex >= 0 {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
									hexDelimiterIndex}
							}
							parseData.initSegmentData(IPv6SegmentCount)
						}
					}
				} else if ipParseData.getProviderIPVersion().IsIPv4() {
					return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv6.separator"}}
				} else if segCount >= IPv6SegmentCount {
					return &addressStringError{addressError{str: str, key: "ipaddress.error.too.many.segments"}}
				}
				hexMaxChars = IPv6SegmentMaxChars // will be ignored for single segment due to checkCharCounts boolean
			}
			if index == segmentStartIndex { // empty segment
				if isMac {
					return &addressStringIndexError{
						addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
						index}
				} else if parseData.getConsecutiveSeparatorIndex() >= 0 {
					return &addressStringError{addressError{str: str, key: "ipaddress.error.ipv6.ambiguous"}}
				}
				parseData.setConsecutiveSeparatorSegmentIndex(segCount)
				parseData.setConsecutiveSeparatorIndex(index - 1)
				assign3Attributes(index, index, parseData, segCount, index)
				parseData.incrementSegmentCount()
			} else if wildcardCount > 0 && !isSingleIPv6 {
				if !stringFormatParams.GetRangeParams().AllowsWildcard() {
					return &addressStringError{addressError{str: str, key: "ipaddress.error.no.wildcard"}}
				}
				totalDigits := index - segmentStartIndex
				if wildcardCount != totalDigits || hexDelimiterIndex >= 0 {
					return &addressStringIndexError{
						addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
						index}
				}
				parseData.setHasWildcard()
				startIndex := index - wildcardCount
				var maxVal uint64
				if isMac {
					if isDoubleSegment {
						maxVal = macMaxTriple
					} else {
						maxVal = MACMaxValuePerSegment
					}
				} else {
					maxVal = IPv6MaxValuePerSegment
				}
				assign6Attributes2Values1Flags(startIndex, index, startIndex, startIndex, index, startIndex,
					parseData, segCount, 0, maxVal, keyWildcard)
				parseData.incrementSegmentCount()
				wildcardCount = 0
			} else {
				startIndex := segmentValueStartIndex
				digitCount := index - startIndex
				noValuesToSet := false
				var value uint64
				var flags, rangeFlags uint32
				if leadingWithZero {
					if digitCount == 1 {
						// can only be a single 0
						if leadingZeroCount == 0 && rangeWildcardIndex < 0 {
							// handles 0 but not 1-0
							assign3Attributes(startIndex, index, parseData, segCount, segmentValueStartIndex)
							parseData.incrementSegmentCount()
							index++
							segmentValueStartIndex = index
							segmentStartIndex = segmentValueStartIndex
							leadingWithZero = false
							continue
						}
					} else {
						if hasDigits {
							leadingZeroCount++
						}
						startIndex += leadingZeroCount
						digitCount -= leadingZeroCount
					}
					leadingWithZero = false
				}
				if leadingZeroCount == 0 {
					if digitCount == 0 {
						// since we have already checked for an empty segment, this can only happen with a range, ie rangeWildcardIndex >= 0
						// we allow an empty range boundary to denote the max value
						if !stringFormatParams.GetRangeParams().AllowsInferredBoundary() {
							// starts with '.', or has two consecutive '.'
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
								index}
						} else if isMac {
							if isSingleSegment {
								if macParseData.isExtended() {
									value = 0xffffffffffffffff
								} else {
									value = 0xffffffffffff
								}
							} else {
								value = MACMaxValuePerSegment
							}
						} else {
							if isSingleIPv6 {
								value = 0xffffffffffffffff
								extendedValue = value
							} else {
								value = IPv6MaxValuePerSegment
							}
						}
						rangeFlags = keyInferredUpperBoundary
					} else {
						if digitCount > hexMaxChars && checkCharCounts {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
								segmentValueStartIndex}
						}
						if singleWildcardCount > 0 {
							noValuesToSet = true
							if rangeWildcardIndex >= 0 {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
									index}
							} else if isSingleIPv6 { //We need this special call here because single ipv6 hex is 128 bits and cannot fit into a long
								if parseErr := parseSingleSegmentSingleWildcard16(currentValueHex, str, startIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
									return parseErr
								}
							} else {
								if parseErr := assignSingleWildcard16(currentValueHex, str, startIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
									return parseErr
								}
							}
							uppercase = false
							singleWildcardCount = 0
						} else {
							if isSingleIPv6 { //We need this special branch here because single ipv6 hex is 128 bits and cannot fit into a long
								midIndex := index - 16
								if startIndex < midIndex {
									extendedValue = parseLong16(str, startIndex, midIndex)
									value = parseLong16(str, midIndex, index)
								} else {
									value = currentValueHex
								}
							} else {
								value = currentValueHex
								if uppercase {
									uppercase = false
								} else {
									flags = keyStandardStr
								}
							}
						}
						hasDigits = false
						currentValueHex = 0
					}
				} else {
					if leadingZeroCount == 1 && (digitCount == 17 || digitCount == 129) &&
						ipv6SpecificOptions.AllowsBinary() && isBinaryDelimiter(str, startIndex) {
						// IPv6 binary - to avoid ambiguity, all binary digits must be present, and preceded by 0b.
						// So for a single segment IPv6 segment:
						// 0b11 is a hex segment, 0b111 is hex, 0b1111 is invalid (too short for binary, too long for hex), 0b1100110011001100 is a binary segment
						if !stringFormatParams.AllowsLeadingZeros() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
						}
						startIndex++ // exclude the 'b' in 0b1100
						digitCount-- // exclude the 'b'
						if singleWildcardCount > 0 {
							if rangeWildcardIndex >= 0 {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
									index}
							} else if isSingleIPv6 {
								if parseErr := parseSingleSegmentSingleWildcard2(str, startIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
									return parseErr
								}
							} else {
								if parseErr := switchSingleWildcard2(currentValueHex, str, startIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
									return parseErr
								}
							}
							noValuesToSet = true
							singleWildcardCount = 0
						} else {
							if isSingleIPv6 { //We need this special branch here because single ipv6 hex is 128 bits and cannot fit into a long
								midIndex := index - 64
								extendedValue = parseLong2(str, startIndex, midIndex)
								value = parseLong2(str, midIndex, index)
							} else {
								value, err = switchValue2(currentValueHex, str, digitCount)
								if err != nil {
									return err
								}
							}
							flags = 2 // radix
						}
						ipParseData.setHasBinaryDigits(true)
						hasDigits = false
						currentValueHex = 0
					} else {
						if digitCount > hexMaxChars && checkCharCounts {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
								segmentValueStartIndex}
						} else if !stringFormatParams.AllowsLeadingZeros() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
						} else if !stringFormatParams.AllowsUnlimitedLeadingZeros() && checkCharCounts && (digitCount+leadingZeroCount) > hexMaxChars {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
								segmentValueStartIndex}
						}
						if singleWildcardCount > 0 {
							noValuesToSet = true
							if rangeWildcardIndex >= 0 {
								return &addressStringIndexError{
									addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
									index}
							} else if isSingleIPv6 { //We need this special call here because single ipv6 hex is 128 bits and cannot fit into a long
								if parseErr := parseSingleSegmentSingleWildcard16(currentValueHex, str, startIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
									return parseErr
								}
							} else {
								if parseErr := assignSingleWildcard16(currentValueHex, str, startIndex, index, singleWildcardCount, parseData, segCount, segmentValueStartIndex, stringFormatParams); parseErr != nil {
									return parseErr
								}
							}
							uppercase = false
							singleWildcardCount = 0
						} else {
							if isSingleIPv6 { //We need this special branch here because single ipv6 hex is 128 bits and cannot fit into a long
								midIndex := index - 16
								if startIndex < midIndex {
									extendedValue = parseLong16(str, startIndex, midIndex)
									value = parseLong16(str, midIndex, index)
								} else {
									value = currentValueHex
								}
							} else {
								value = currentValueHex
								if uppercase {
									uppercase = false
								} else {
									flags = keyStandardStr
								}
							}
						}
						hasDigits = false
						currentValueHex = 0
					}
				}
				if rangeWildcardIndex >= 0 {
					frontStartIndex := rangeWildcardIndex - frontDigitCount
					frontEndIndex := rangeWildcardIndex
					frontLeadingZeroStartIndex := frontStartIndex - frontLeadingZeroCount
					frontTotalDigitCount := frontDigitCount + frontLeadingZeroCount //the stuff that uses frontLeadingZeroCount needs to be sectioned off when singleIPv6
					if !stringFormatParams.GetRangeParams().AllowsRangeSeparator() {
						return &addressStringError{addressError{str: str, key: "ipaddress.error.no.range"}}
					} else if frontHexDelimiterIndex >= 0 && !isSingleSegment {
						return &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
							frontHexDelimiterIndex}
					} else if frontSingleWildcardCount > 0 || frontWildcardCount > 0 { // no wildcards in ranges
						return &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.combination.at.index"}},
							rangeWildcardIndex}
					} else if isMac && !macSpecificOptions.AllowsShortSegments() && frontTotalDigitCount < 2 {
						return &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.short.at.index"}},
							frontLeadingZeroStartIndex}
					}
					var upperRadix uint32
					frontIsBinary := false
					var front, extendedFront uint64
					frontEmpty := frontStartIndex == frontEndIndex
					isReversed := false
					if frontEmpty {
						if !stringFormatParams.GetRangeParams().AllowsInferredBoundary() {
							return &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
								index}
						}
						rangeFlags |= keyInferredLowerBoundary
					} else {
						if frontLeadingZeroCount == 1 && frontDigitCount == 17 && ipv6SpecificOptions.AllowsBinary() && isBinaryDelimiter(str, frontStartIndex) {
							// IPv6 binary - to avoid ambiguity, all binary digits must be present, and preceded by 0b.
							frontStartIndex++ // exclude the 'b' in 0b1100
							frontDigitCount-- // exclude the 'b'
							front, err = switchValue2(currentFrontValueHex, str, frontDigitCount)
							if err != nil {
								return err
							}
							upperRadix = 2          // radix
							rangeFlags = upperRadix // radix
							ipParseData.setHasBinaryDigits(true)
							isReversed = front > value
							frontIsBinary = true
						} else if isSingleIPv6 { //We need this special block here because single ipv6 hex is 128 bits and cannot fit into a long
							if frontDigitCount == 129 { // binary
								frontStartIndex++ // exclude the 'b' in 0b1100
								frontDigitCount-- // exclude the 'b'
								upperRadix = 2    // radix
								rangeFlags = 2    // radix
								ipParseData.setHasBinaryDigits(true)
								frontMidIndex := frontEndIndex - 64
								extendedFront = parseLong2(str, frontStartIndex, frontMidIndex)
								front = parseLong2(str, frontMidIndex, frontEndIndex)
							} else {
								frontMidIndex := frontEndIndex - 16
								if frontStartIndex < frontMidIndex {
									extendedFront = parseLong16(str, frontStartIndex, frontMidIndex)
									front = parseLong16(str, frontMidIndex, frontEndIndex)
								} else {
									front = currentFrontValueHex
								}
							}
							isReversed = (extendedFront > extendedValue) || (extendedFront == extendedValue && front > value)
						} else {
							if !stringFormatParams.AllowsLeadingZeros() && frontLeadingZeroCount > 0 {
								return &addressStringError{addressError{str: str, key: "ipaddress.error.segment.leading.zeros"}}
							} else if checkCharCounts {
								if frontDigitCount > hexMaxChars {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
										frontLeadingZeroStartIndex}
								} else if !stringFormatParams.AllowsUnlimitedLeadingZeros() && frontTotalDigitCount > hexMaxChars {
									return &addressStringIndexError{
										addressStringError{addressError{str: str, key: "ipaddress.error.segment.too.long.at.index"}},
										frontLeadingZeroStartIndex}
								}
							}
							front = currentFrontValueHex
							isReversed = front > value
							extendedFront = 0
						}
					}
					backEndIndex := index
					if isReversed {
						if !stringFormatParams.GetRangeParams().AllowsReverseRange() {
							return &addressStringError{addressError{str: str, key: "ipaddress.error.invalidRange"}}
						}
						// switcheroo
						frontStartIndex, startIndex = startIndex, frontStartIndex
						frontEndIndex, backEndIndex = backEndIndex, frontEndIndex
						frontLeadingZeroStartIndex, segmentValueStartIndex = segmentValueStartIndex, frontLeadingZeroStartIndex
						front, value = value, front
						extendedFront, extendedValue = extendedValue, extendedFront
					}
					if isSingleIPv6 {
						assign7Attributes4Values1Flags(frontStartIndex, frontEndIndex, frontLeadingZeroStartIndex, startIndex, backEndIndex, segmentValueStartIndex,
							parseData, segCount, front, extendedFront, value, extendedValue, rangeFlags|keyRangeWildcard, upperRadix)
					} else {
						if !frontUppercase && !frontEmpty && !isReversed && !frontIsBinary {
							if leadingZeroCount == 0 && (flags&keyStandardStr) != 0 && frontIsStandardRangeChar {
								rangeFlags |= keyStandardRangeStr | keyStandardStr
							} else {
								rangeFlags |= keyStandardStr
							}
						}
						assign6Attributes2Values2Flags(frontStartIndex, frontEndIndex, frontLeadingZeroStartIndex, startIndex, backEndIndex, segmentValueStartIndex,
							parseData, segCount, front, value, rangeFlags|keyRangeWildcard, upperRadix)
					}
					rangeWildcardIndex = -1
				} else if !noValuesToSet {
					if isSingleIPv6 {
						assign3Attributes2Values1Flags(startIndex, index, segmentValueStartIndex, parseData, segCount, value, extendedValue, flags)
					} else {
						assign3Attributes1Values1Flags(startIndex, index, segmentValueStartIndex, parseData, segCount, value, flags)
					}
				}
				parseData.incrementSegmentCount()
				leadingZeroCount = 0
			}
			index++
			segmentValueStartIndex = index
			segmentStartIndex = segmentValueStartIndex
			// end of IPv6 and MAC48Len segments
		} // end of all cases
	} // end of character loop
	return nil
}

func (strValidator) validatePrefixLenStr(fullAddr string, version IPVersion) (prefixLen PrefixLen, err addrerr.AddressStringError) {
	var qualifier parsedHostIdentifierStringQualifier
	isPrefix, err := validatePrefix(fullAddr, nil, defaultIPAddrParameters, nil,
		&qualifier, 0, len(fullAddr), version)
	if !isPrefix {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalidCIDRPrefix"}}
	} else {
		prefixLen = qualifier.getNetworkPrefixLen()
	}
	return
}

func parsePortOrService(
	fullAddr string,
	zone *Zone,
	validationOptions addrstrparam.HostNameParams,
	res *parsedHostIdentifierStringQualifier,
	index,
	endIndex int) (err addrerr.AddressStringError) {
	isPort := true
	var hasLetter, hasDigits, isAll bool
	var charCount, digitCount int
	var port int
	lastHyphen := -1
	charArray := chars
	for i := index; i < endIndex; i++ {
		c := fullAddr[i]
		if c >= '1' && c <= '9' {
			if isPort {
				digitCount++
				if digitCount > 5 { // 65535 is max
					isPort = false
				} else {
					hasDigits = true
					port = port*10 + int(charArray[c])
				}
			}
			charCount++
		} else if c == '0' {
			if isPort && hasDigits {
				digitCount++
				if digitCount > 5 { // 65535 is max
					isPort = false
				} else {
					port *= 10
				}
			}
			charCount++
		} else {
			//http://www.iana.org/assignments/port-numbers
			//valid service name chars:
			//https://tools.ietf.org/html/rfc6335#section-5.1
			//https://tools.ietf.org/html/rfc6335#section-10.1
			isPort = false
			isHyphen := c == '-'
			isAll = c == SegmentWildcard
			if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || isHyphen || isAll {
				if isHyphen {
					if i == index {
						err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalid.service.hyphen.start"}}
						return
					} else if i-1 == lastHyphen {
						err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalid.service.hyphen.consecutive"}}
						return
					} else if i == endIndex-1 {
						err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalid.service.hyphen.end"}}
						return
					}
					lastHyphen = i
				} else if isAll {
					if i > index {
						err = &addressStringIndexError{
							addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.character.combination.at.index"}},
							i}
						return
					} else if i+1 < endIndex {
						err = &addressStringIndexError{
							addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.character.combination.at.index"}},
							i + 1}
						return
					}
					hasLetter = true
					charCount++
					break
				} else {
					hasLetter = true
				}
				charCount++
			} else {
				err = &addressStringIndexError{
					addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalid.port.service"}},
					i}
				return
			}
		}
	}
	if isPort {
		if !validationOptions.AllowsPort() {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.port"}}
			return
		} else if port == 0 {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalidPort.no.digits"}}
			return
		} else if port > 65535 {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalidPort.too.large"}}
			return
		}
		res.setZone(zone)
		res.port = cachePorts(PortInt(port))
		return
	} else if !validationOptions.AllowsService() {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.service"}}
		return
	} else if charCount == 0 {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalidService.no.chars"}}
		return
	} else if charCount > 15 {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalidService.too.long"}}
		return
	} else if !hasLetter {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.host.error.invalidService.no.letter"}}
		return
	}
	res.setZone(zone)
	res.service = fullAddr[index:endIndex]
	return
}

func parseValidatedPrefix(
	result BitCount,
	fullAddr string,
	zone *Zone,
	validationOptions addrstrparam.IPAddressStringParams,
	res *parsedHostIdentifierStringQualifier,
	digitCount,
	leadingZeros int,
	ipVersion IPVersion) (err addrerr.AddressStringError) {
	if digitCount == 0 {
		//we know leadingZeroCount is > 0 since we have checked already if there were no characters at all
		leadingZeros--
		// digitCount++ digitCount is unused after this, no need for it to be accurate
	}
	asIPv4 := ipVersion.IsIPv4()
	if asIPv4 {
		if leadingZeros > 0 && !validationOptions.GetIPv4Params().AllowsPrefixLenLeadingZeros() {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipv4.prefix.leading.zeros"}}
			return
		}
		if result > IPv4BitCount {
			allowPrefixesBeyondAddressSize := validationOptions.GetIPv4Params().AllowsPrefixesBeyondAddressSize()
			if !allowPrefixesBeyondAddressSize {
				err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.prefixSize"}}
				return
			}
			result = IPv4BitCount
		}
	} else {
		if leadingZeros > 0 && !validationOptions.GetIPv6Params().AllowsPrefixLenLeadingZeros() {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipv6.prefix.leading.zeros"}}
			return
		}
		if result > IPv6BitCount {
			allowPrefixesBeyondAddressSize := validationOptions.GetIPv6Params().AllowsPrefixesBeyondAddressSize()
			if !allowPrefixesBeyondAddressSize {
				err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.prefixSize"}}
				return
			}
			result = IPv6BitCount
		}
	}
	res.networkPrefixLength = cacheBitCount(result)
	res.setZone(zone)
	return
}

func validatePrefix(
	fullAddr string,
	zone *Zone,
	validationOptions addrstrparam.IPAddressStringParams,
	hostValidationOptions addrstrparam.HostNameParams,
	res *parsedHostIdentifierStringQualifier,
	index,
	endIndex int,
	ipVersion IPVersion) (isPrefix bool, err addrerr.AddressStringError) {
	if index == len(fullAddr) {
		return
	}
	isPrefix = true
	prefixEndIndex := endIndex
	hasDigits := false
	var result BitCount
	var leadingZeros int
	charArray := chars
	for i := index; i < endIndex; i++ {
		c := fullAddr[i]
		if c >= '1' && c <= '9' {
			hasDigits = true
			result = result*10 + BitCount(charArray[c])
		} else if c == '0' {
			if hasDigits {
				result *= 10
			} else {
				leadingZeros++
			}
		} else if c == PortSeparator && hostValidationOptions != nil &&
			(hostValidationOptions.AllowsPort() || hostValidationOptions.AllowsService()) && i > index {
			// check if we have a port or service.  If not, possibly an IPv6 mask.
			// Also, parsing for port first (rather than prefix) allows us to call
			// parseValidatedPrefix with the knowledge that whatever is supplied can only be a prefix.
			portErr := parsePortOrService(fullAddr, zone, hostValidationOptions, res, i+1, endIndex)
			//portQualifier, err = parsePortOrService(fullAddr, zone, hostValidationOptions, res, i+1, endIndex)
			if portErr != nil {
				return
			}
			prefixEndIndex = i
			break
		} else {
			isPrefix = false
			break
		}
	}
	//we treat as a prefix if all the characters were digits, even if there were too many, unless the mask options allow for inet_aton single segment
	if isPrefix {
		err = parseValidatedPrefix(result, fullAddr,
			zone, validationOptions, res, prefixEndIndex-index /* digitCount */, leadingZeros, ipVersion)
	}
	return
}

func parseAddressQualifier(
	fullAddr string,
	validationOptions addrstrparam.IPAddressStringParams,
	hostValidationOptions addrstrparam.HostNameParams,
	ipAddressParseData *ipAddressParseData,
	endIndex int) (err addrerr.AddressStringError) {
	qualifierIndex := ipAddressParseData.getQualifierIndex()
	addressIsEmpty := ipAddressParseData.getAddressParseData().isProvidingEmpty()
	ipVersion := ipAddressParseData.getProviderIPVersion()
	res := ipAddressParseData.getQualifier()
	if ipAddressParseData.hasPrefixSeparator() {
		return parsePrefix(fullAddr, nil, validationOptions, hostValidationOptions,
			res, addressIsEmpty, qualifierIndex, endIndex, ipVersion)
	} else if ipAddressParseData.isZoned() {
		if ipAddressParseData.isBase85Zoned() && !ipAddressParseData.isProvidingBase85IPv6() {
			err = &addressStringIndexError{
				addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.character.at.index"}},
				qualifierIndex - 1}
			return
		}
		if addressIsEmpty {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.only.zone"}}
			return
		}
		return parseZone(fullAddr, validationOptions, res, addressIsEmpty, qualifierIndex, endIndex, ipVersion)
	}
	return
}

func parseHostAddressQualifier(
	fullAddr string,
	validationOptions addrstrparam.IPAddressStringParams,
	hostValidationOptions addrstrparam.HostNameParams,
	isPrefixed,
	hasPort bool,
	ipAddressParseData *ipAddressParseData,
	qualifierIndex,
	endIndex int) (err addrerr.AddressStringError) {
	res := ipAddressParseData.getQualifier()
	addressIsEmpty := ipAddressParseData.getAddressParseData().isProvidingEmpty()
	ipVersion := ipAddressParseData.getProviderIPVersion()
	if isPrefixed {
		return parsePrefix(fullAddr, nil, validationOptions, hostValidationOptions,
			res, addressIsEmpty, qualifierIndex, endIndex, ipVersion)
	} else if ipAddressParseData.isZoned() {
		if addressIsEmpty {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.only.zone"}}
			return
		}
		return parseEncodedZone(fullAddr, validationOptions, res, addressIsEmpty, qualifierIndex, endIndex, ipVersion)
	} else if hasPort { //isPort is always false when validating an address
		return parsePortOrService(fullAddr, nil, hostValidationOptions, res, qualifierIndex, endIndex)
	}
	return
}

func parsePrefix(
	fullAddr string,
	zone *Zone,
	validationOptions addrstrparam.IPAddressStringParams,
	hostValidationOptions addrstrparam.HostNameParams,
	res *parsedHostIdentifierStringQualifier,
	addressIsEmpty bool,
	index,
	endIndex int,
	ipVersion IPVersion) (err addrerr.AddressStringError) {
	if validationOptions.AllowsPrefix() {
		var isPrefix bool
		isPrefix, err = validatePrefix(fullAddr, zone, validationOptions, hostValidationOptions,
			res, index, endIndex, ipVersion)
		if err != nil || isPrefix {
			return
		}
	}
	if addressIsEmpty {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.mask.address.empty"}}
	} else if validationOptions.AllowsMask() {
		//check for a mask
		//check if we need a new validation options for the mask
		maskOptions := toMaskOptions(validationOptions, ipVersion)
		pa := &parsedIPAddress{
			ipAddressParseData: ipAddressParseData{addressParseData: addressParseData{str: fullAddr}},
			options:            maskOptions,
		}
		err = validateIPAddress(maskOptions, fullAddr, index, endIndex, pa.getIPAddressParseData(), false)
		if err != nil {
			err = &addressStringNestedError{
				addressStringError: addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalidCIDRPrefixOrMask"}},
				nested:             err,
			}
			return
		}
		maskParseData := pa.getAddressParseData()
		if maskParseData.isProvidingEmpty() {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.mask.empty"}}
			return
		} else if maskParseData.isAll() {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.mask.wildcard"}}
			return
		}
		err = checkSegments(fullAddr, maskOptions, pa.getIPAddressParseData())
		if err != nil {
			err = &addressStringNestedError{
				addressStringError: addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalidCIDRPrefixOrMask"}},
				nested:             err,
			}
			return
		}
		maskEndIndex := maskParseData.getAddressEndIndex()
		if maskEndIndex != endIndex { // 1.2.3.4/ or 1.2.3.4// or 1.2.3.4/%
			err = &addressStringIndexError{
				addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.mask.extra.chars"}},
				maskEndIndex + 1}
			return
		}
		maskVersion := pa.getProviderIPVersion()
		if maskVersion.IsIPv4() && maskParseData.getSegmentCount() == 1 && !maskParseData.hasWildcard() &&
			!validationOptions.GetIPv4Params().Allows_inet_aton_single_segment_mask() { //1.2.3.4/33 where 33 is an aton_inet single segment address and not a prefix length
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.mask.single.segment"}}
			return
		} else if !ipVersion.IsIndeterminate() && (maskVersion.IsIPv4() != ipVersion.IsIPv4() || maskVersion.IsIPv6() != ipVersion.IsIPv6()) {
			//note that this also covers the cases of non-standard addresses in the mask, ie mask neither ipv4 or ipv6
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipMismatch"}}
			return
		}
		res.mask = pa
		res.setZone(zone)
	} else if validationOptions.AllowsPrefix() {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalidCIDRPrefixOrMask"}}
	} else {
		err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.CIDRNotAllowed"}}
	}
	return
}

func parseHostNameQualifier(
	fullAddr string,
	validationOptions addrstrparam.IPAddressStringParams,
	hostValidationOptions addrstrparam.HostNameParams,
	res *parsedHostIdentifierStringQualifier,
	isPrefixed,
	isPort, // always false for address
	addressIsEmpty bool,
	index,
	endIndex int,
	ipVersion IPVersion) (err addrerr.AddressStringError) {
	if isPrefixed {
		return parsePrefix(fullAddr, nil, validationOptions, hostValidationOptions,
			res, addressIsEmpty, index, endIndex, ipVersion)
	} else if isPort { // isPort is always false when validating an address
		return parsePortOrService(fullAddr, nil, hostValidationOptions, res, index, endIndex)
	}
	//res = noQualifier
	return
}

// ValidateZoneStr returns an error if the given zone is invalid
func ValidateZoneStr(zoneStr string) (zone Zone, err addrerr.AddressStringError) {
	for i := 0; i < len(zoneStr); i++ {
		c := zone[i]
		if c == PrefixLenSeparator {
			err = &addressStringIndexError{addressStringError{addressError{str: zoneStr, key: "ipaddress.error.invalid.zone"}}, i}
			return
		}
		if c == IPv6SegmentSeparator {
			err = &addressStringIndexError{addressStringError{addressError{str: zoneStr, key: "ipaddress.error.invalid.zone"}}, i}
			return
		}
	}
	return Zone(zoneStr), nil
}

func isReserved(c byte) bool {
	isUnreserved :=
		(c >= '0' && c <= '9') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			c == RangeSeparator ||
			c == LabelSeparator ||
			c == '_' ||
			c == '~'
	return !isUnreserved
}

func parseZone(
	fullAddr string,
	validationOptions addrstrparam.IPAddressStringParams,
	res *parsedHostIdentifierStringQualifier,
	addressIsEmpty bool,
	index,
	endIndex int,
	ipVersion IPVersion) (err addrerr.AddressStringError) {
	if index == endIndex && !validationOptions.GetIPv6Params().AllowsEmptyZone() {
		err = &addressStringIndexError{addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.zone"}}, index}
		return
	}
	for i := index; i < endIndex; i++ {
		c := fullAddr[i]
		if c == PrefixLenSeparator {
			if i == index && !validationOptions.GetIPv6Params().AllowsEmptyZone() {
				err = &addressStringIndexError{addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.zone"}}, index}
				return
			}
			zone := Zone(fullAddr[index:i])
			return parsePrefix(fullAddr, &zone, validationOptions, nil, res, addressIsEmpty, i+1, endIndex, ipVersion)
		} else if c == IPv6SegmentSeparator {
			err = &addressStringIndexError{
				addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.zone"}},
				i}
			return
		}
	}
	z := Zone(fullAddr[index:endIndex])
	res.setZone(&z)
	return
}

func parseEncodedZone(
	fullAddr string,
	validationOptions addrstrparam.IPAddressStringParams,
	res *parsedHostIdentifierStringQualifier,
	addressIsEmpty bool,
	index,
	endIndex int,
	ipVersion IPVersion) (err addrerr.AddressStringError) {
	if index == endIndex && !validationOptions.GetIPv6Params().AllowsEmptyZone() {
		err = &addressStringIndexError{addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.zone"}}, index}
		return
	}
	var result strings.Builder
	var zone string
	for i := index; i < endIndex; i++ {
		c := fullAddr[i]

		//we are in here when we have a square bracketed host like [::1]
		//not if we have a HostName with no brackets

		//https://tools.ietf.org/html/rfc6874
		//https://tools.ietf.org/html/rfc4007#section-11.7
		if c == IPv6ZoneSeparator {
			if i+2 >= endIndex {
				err = &addressStringIndexError{
					addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.zone.encoding"}},
					i}
				return
			}
			//percent encoded
			if result.Cap() == 0 {
				result.Grow(endIndex - index)
				result.WriteString(fullAddr[index:i])
			}
			charArray := chars
			i++
			c = charArray[fullAddr[i]] << 4
			i++
			c |= charArray[fullAddr[i]]
		} else if c == PrefixLenSeparator {
			if i == index && !validationOptions.GetIPv6Params().AllowsEmptyZone() {
				err = &addressStringIndexError{addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.zone"}}, index}
				return
			}
			if result.Cap() > 0 {
				zone = result.String()
			} else {
				zone = fullAddr[index:i]
			}
			z := Zone(zone)
			return parsePrefix(fullAddr, &z, validationOptions, nil, res, addressIsEmpty, i+1, endIndex, ipVersion)
		} else if isReserved(c) {
			err = &addressStringIndexError{
				addressStringError{addressError{str: fullAddr, key: "ipaddress.error.invalid.zone"}},
				i}
			return
		}
		if result.Cap() > 0 {
			result.WriteByte(c)
		}
	}
	if result.Len() == 0 {
		zone = fullAddr[index:endIndex]
	} else {
		zone = result.String()
	}
	z := Zone(zone)
	res.setZone(&z)
	return
}

// whether no wildcards or range characters allowed
func isNoRange(rp addrstrparam.RangeParams) bool {
	return !rp.AllowsWildcard() && !rp.AllowsRangeSeparator() && !rp.AllowsSingleWildcard()
}

/**
 * Some options are not supported in masks (prefix, wildcards, etc)
 * So we eliminate those options while preserving the others from the address options.
 * @param validationOptions
 * @param ipVersion
 * @return
 */
func toMaskOptions(validationOptions addrstrparam.IPAddressStringParams,
	ipVersion IPVersion) (res addrstrparam.IPAddressStringParams) {
	//We must provide options that do not allow a mask with wildcards or ranges
	var builder *addrstrparam.IPAddressStringParamsBuilder
	if ipVersion.IsIndeterminate() || ipVersion.IsIPv6() {
		ipv6Options := validationOptions.GetIPv6Params()
		if !isNoRange(ipv6Options.GetRangeParams()) {
			builder = new(addrstrparam.IPAddressStringParamsBuilder).Set(validationOptions)
			builder.GetIPv6AddressParamsBuilder().SetRangeParams(addrstrparam.NoRange)
		}
		if ipv6Options.AllowsMixed() && !isNoRange(ipv6Options.GetMixedParams().GetIPv4Params().GetRangeParams()) {
			if builder == nil {
				builder = new(addrstrparam.IPAddressStringParamsBuilder).Set(validationOptions)
			}
			builder.GetIPv6AddressParamsBuilder().SetRangeParams(addrstrparam.NoRange)
		}
	}
	if ipVersion.IsIndeterminate() || ipVersion.IsIPv4() {
		ipv4Options := validationOptions.GetIPv4Params()
		if !isNoRange(ipv4Options.GetRangeParams()) {
			if builder == nil {
				builder = new(addrstrparam.IPAddressStringParamsBuilder).Set(validationOptions)
			}
			builder.GetIPv4AddressParamsBuilder().SetRangeParams(addrstrparam.NoRange)
		}
	}
	if validationOptions.AllowsAll() {
		if builder == nil {
			builder = new(addrstrparam.IPAddressStringParamsBuilder).Set(validationOptions)
		}
		builder.AllowAll(false)
	}
	if builder == nil {
		res = validationOptions
	} else {
		res = builder.ToParams()
	}
	return
}

func assign3Attributes(start, end int, parseData *addressParseData, parsedSegIndex, leadingZeroStartIndex int) {
	ustart := uint32(start)
	uend := uint32(end)
	uleadingZeroStart := uint32(leadingZeroStartIndex)
	parseData.setIndex(parsedSegIndex,
		keyLowerStrDigitsIndex, uleadingZeroStart,
		keyLowerStrStartIndex, ustart,
		keyLowerStrEndIndex, uend,
		keyUpperStrDigitsIndex, uleadingZeroStart,
		keyUpperStrStartIndex, ustart,
		keyUpperStrEndIndex, uend)
}

func assign4Attributes(start, end int, parseData *addressParseData, parsedSegIndex, radix, leadingZeroStartIndex int) {
	ustart := uint32(start)
	uend := uint32(end)
	uleadingZeroStart := uint32(leadingZeroStartIndex)
	parseData.set7IndexFlags(parsedSegIndex,
		keyLowerRadixIndex, uint32(radix),
		keyLowerStrDigitsIndex, uleadingZeroStart,
		keyLowerStrStartIndex, ustart,
		keyLowerStrEndIndex, uend,
		keyUpperStrDigitsIndex, uleadingZeroStart,
		keyUpperStrStartIndex, ustart,
		keyUpperStrEndIndex, uend)
}

func assign7Attributes4Values1Flags(frontStart, frontEnd, frontLeadingZeroStartIndex, start, end, leadingZeroStartIndex int,
	parseData *addressParseData, parsedSegIndex int, frontValue, frontExtendedValue, value, extendedValue uint64, flags uint32, upperRadix uint32) {
	parseData.set8Index4ValuesFlags(parsedSegIndex,
		flagsIndex, flags,
		keyLowerStrDigitsIndex, uint32(frontLeadingZeroStartIndex),
		keyLowerStrStartIndex, uint32(frontStart),
		keyLowerStrEndIndex, uint32(frontEnd),
		keyUpperRadixIndex, uint32(upperRadix),
		keyUpperStrDigitsIndex, uint32(leadingZeroStartIndex),
		keyUpperStrStartIndex, uint32(start),
		keyUpperStrEndIndex, uint32(end),
		keyLower, frontValue,
		keyExtendedLower, frontExtendedValue,
		keyUpper, value,
		keyExtendedUpper, extendedValue)
}

func assign6Attributes4Values1Flags(frontStart, frontEnd, frontLeadingZeroStartIndex, start, end, leadingZeroStartIndex int,
	parseData *addressParseData, parsedSegIndex int, frontValue, frontExtendedValue, value, extendedValue uint64, flags uint32) {
	parseData.set7Index4ValuesFlags(parsedSegIndex,
		flagsIndex, flags,
		keyLowerStrDigitsIndex, uint32(frontLeadingZeroStartIndex),
		keyLowerStrStartIndex, uint32(frontStart),
		keyLowerStrEndIndex, uint32(frontEnd),
		keyUpperStrDigitsIndex, uint32(leadingZeroStartIndex),
		keyUpperStrStartIndex, uint32(start),
		keyUpperStrEndIndex, uint32(end),
		keyLower, frontValue,
		keyExtendedLower, frontExtendedValue,
		keyUpper, value,
		keyExtendedUpper, extendedValue)
}

func assign6Attributes2Values1Flags(frontStart, frontEnd, frontLeadingZeroStartIndex, start, end, leadingZeroStartIndex int,
	parseData *addressParseData, parsedSegIndex int, frontValue, value uint64, flags uint32) {
	parseData.set7Index2ValuesFlags(parsedSegIndex,
		flagsIndex, flags,
		keyLowerStrDigitsIndex, uint32(frontLeadingZeroStartIndex),
		keyLowerStrStartIndex, uint32(frontStart),
		keyLowerStrEndIndex, uint32(frontEnd),
		keyUpperStrDigitsIndex, uint32(leadingZeroStartIndex),
		keyUpperStrStartIndex, uint32(start),
		keyUpperStrEndIndex, uint32(end),
		keyLower, frontValue,
		keyUpper, value)
}

func assign6Attributes2Values2Flags(frontStart, frontEnd, frontLeadingZeroStartIndex, start, end, leadingZeroStartIndex int,
	parseData *addressParseData, parsedSegIndex int, frontValue, value uint64, flags /* includes lower radix */ uint32, upperRadix uint32) {
	parseData.set8Index2ValuesFlags(parsedSegIndex,
		flagsIndex, flags,
		keyLowerStrDigitsIndex, uint32(frontLeadingZeroStartIndex),
		keyLowerStrStartIndex, uint32(frontStart),
		keyLowerStrEndIndex, uint32(frontEnd),
		keyUpperRadixIndex, uint32(upperRadix),
		keyUpperStrDigitsIndex, uint32(leadingZeroStartIndex),
		keyUpperStrStartIndex, uint32(start),
		keyUpperStrEndIndex, uint32(end),
		keyLower, frontValue,
		keyUpper, value)
}

func assign3Attributes2Values1Flags(start, end, leadingZeroStart int,
	parseData *addressParseData, parsedSegIndex int, value, extendedValue uint64, flags uint32) {
	ustart := uint32(start)
	uend := uint32(end)
	uleadingZeroStart := uint32(leadingZeroStart)
	parseData.set7Index4ValuesFlags(parsedSegIndex,
		flagsIndex, flags,
		keyLowerStrDigitsIndex, uleadingZeroStart,
		keyLowerStrStartIndex, ustart,
		keyLowerStrEndIndex, uend,
		keyUpperStrDigitsIndex, uleadingZeroStart,
		keyUpperStrStartIndex, ustart,
		keyUpperStrEndIndex, uend,
		keyLower, value,
		keyExtendedLower, extendedValue,
		keyUpper, value,
		keyExtendedUpper, extendedValue)
}

func assign3Attributes1Values1Flags(start, end, leadingZeroStart int,
	parseData *addressParseData, parsedSegIndex int, value uint64, flags uint32) {
	ustart := uint32(start)
	uend := uint32(end)
	uleadingZeroStart := uint32(leadingZeroStart)
	parseData.set7Index2ValuesFlags(parsedSegIndex,
		flagsIndex, flags,
		keyUpperStrDigitsIndex, uleadingZeroStart,
		keyLowerStrDigitsIndex, uleadingZeroStart,
		keyUpperStrStartIndex, ustart,
		keyLowerStrStartIndex, ustart,
		keyUpperStrEndIndex, uend,
		keyLowerStrEndIndex, uend,
		keyLower, value,
		keyUpper, value)
}

func isBinaryDelimiter(str string, index int) bool {
	c := str[index]
	return c == 'b' || c == 'B'
}

func isHexDelimiter(c byte) bool {
	return c == 'x' || c == 'X'
}

var lowBitsVal uint64 = 0xffffffffffffffff
var lowBitsMask = bigZero().SetUint64(lowBitsVal)

func parseBase85(
	validationOptions addrstrparam.IPAddressStringParams,
	str string,
	strStartIndex,
	strEndIndex int,
	ipAddressParseData *ipAddressParseData,
	extendedRangeWildcardIndex,
	totalCharacterCount,
	index int) (bool, addrerr.AddressStringError) {

	parseData := ipAddressParseData.getAddressParseData()
	if extendedRangeWildcardIndex < 0 {
		if totalCharacterCount == ipv6Base85SingleSegmentDigitCount {
			ipAddressParseData.setVersion(IPv6)
			val := parse85(str, strStartIndex, strEndIndex)
			var lowBits, highBits big.Int
			lowBits.And(val, lowBitsMask)
			val.Rsh(val, 64)
			highBits.And(val, lowBitsMask)
			val.Rsh(val, 64)
			//note that even with the correct number of base-85 digits, we can have a value too large
			if !bigIsZero(val) {
				return false, &addressStringError{addressError{str: str, key: "ipaddress.error.address.too.large"}}
			}
			value, extendedValue := lowBits.Uint64(), highBits.Uint64()
			parseData.initSegmentData(1)
			parseData.incrementSegmentCount()
			assign3Attributes2Values1Flags(strStartIndex, strEndIndex, strStartIndex, parseData, 0, value, extendedValue, 85)
			ipAddressParseData.setBase85(true)
			return true, nil
		}
	} else {
		if totalCharacterCount == (ipv6Base85SingleSegmentDigitCount<<1)+len(IPv6AlternativeRangeSeparatorStr) /* two base 85 addresses */ ||
			(totalCharacterCount == ipv6Base85SingleSegmentDigitCount+len(IPv6AlternativeRangeSeparatorStr) &&
				(extendedRangeWildcardIndex == 0 || extendedRangeWildcardIndex+len(IPv6AlternativeRangeSeparatorStr) == strEndIndex)) { /* note that we already check that extendedRangeWildcardIndex is at index 20 */
			ipv6SpecificOptions := validationOptions.GetIPv6Params()
			if !ipv6SpecificOptions.GetRangeParams().AllowsRangeSeparator() {
				return false, &addressStringError{addressError{str: str, key: "ipaddress.error.no.range"}}
			}
			ipAddressParseData.setVersion(IPv6)
			frontEndIndex, flags := extendedRangeWildcardIndex, uint32(0)
			var value, value2, extendedValue, extendedValue2 uint64
			var lowerStart, lowerEnd, upperStart, upperEnd int
			//
			if frontEndIndex == strStartIndex+ipv6Base85SingleSegmentDigitCount {
				val := parse85(str, strStartIndex, frontEndIndex)
				var lowBits, highBits big.Int
				lowBits.And(val, lowBitsMask)
				val.Rsh(val, 64)
				highBits.And(val, lowBitsMask)
				value, extendedValue = lowBits.Uint64(), highBits.Uint64()
				if frontEndIndex+len(IPv6AlternativeRangeSeparatorStr) < strEndIndex {
					val2 := parse85(str, frontEndIndex+len(IPv6AlternativeRangeSeparatorStr), strEndIndex)
					var lowBits, highBits big.Int
					lowBits.And(val2, lowBitsMask)
					val2.Rsh(val2, 64)
					highBits.And(val2, lowBitsMask)
					value2, extendedValue2 = lowBits.Uint64(), highBits.Uint64()
					if val.Cmp(val2) > 0 {
						if !ipv6SpecificOptions.GetRangeParams().AllowsReverseRange() {
							return false, &addressStringError{addressError{str: str, key: "ipaddress.error.invalidRange"}}
						}
						//note that even with the correct number of base-85 digits, we can have a value too large
						val.Rsh(val, 64)
						if !bigIsZero(val) {
							return false, &addressStringError{addressError{str: str, key: "ipaddress.error.address.too.large"}}
						}
						lowerStart = frontEndIndex + len(IPv6AlternativeRangeSeparatorStr)
						lowerEnd = strEndIndex
						upperStart = strStartIndex
						upperEnd = frontEndIndex
					} else {
						//note that even with the correct number of base-85 digits, we can have a value too large
						val2.Rsh(val2, 64)
						if !bigIsZero(val2) {
							return false, &addressStringError{addressError{str: str, key: "ipaddress.error.address.too.large"}}
						}
						lowerStart = strStartIndex
						lowerEnd = frontEndIndex
						upperStart = frontEndIndex + len(IPv6AlternativeRangeSeparatorStr)
						upperEnd = strEndIndex
					}
				} else {
					if !ipv6SpecificOptions.GetRangeParams().AllowsInferredBoundary() {
						return false, &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
							extendedRangeWildcardIndex}
					}
					lowerStart = strStartIndex
					lowerEnd = frontEndIndex
					upperStart = strEndIndex
					upperEnd = strEndIndex
					value2 = lowBitsVal
					extendedValue2 = lowBitsVal
					flags = keyInferredUpperBoundary
				}
			} else if frontEndIndex == 0 {
				if !ipv6SpecificOptions.GetRangeParams().AllowsInferredBoundary() {
					return false, &addressStringIndexError{
						addressStringError{addressError{str: str, key: "ipaddress.error.empty.segment.at.index"}},
						index}
				}
				flags = keyInferredLowerBoundary
				val2 := parse85(str, frontEndIndex+len(IPv6AlternativeRangeSeparatorStr), strEndIndex)
				var lowBits, highBits big.Int
				lowBits.And(val2, lowBitsMask)
				val2.Rsh(val2, 64)
				highBits.And(val2, lowBitsMask)
				val2.Rsh(val2, 64)
				//note that even with the correct number of base-85 digits, we can have a value too large
				if !bigIsZero(val2) {
					return false, &addressStringError{addressError{str: str, key: "ipaddress.error.address.too.large"}}
				}
				value2, extendedValue2 = lowBits.Uint64(), highBits.Uint64()
				upperStart = len(IPv6AlternativeRangeSeparatorStr)
				upperEnd = strEndIndex
			} else {
				return false, &addressStringIndexError{
					addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
					extendedRangeWildcardIndex}
			}
			parseData.incrementSegmentCount()
			parseData.initSegmentData(1)
			//parseData.setHasRange();
			assign7Attributes4Values1Flags(lowerStart, lowerEnd, lowerStart, upperStart, upperEnd, upperStart,
				parseData, 0, value, extendedValue, value2, extendedValue2,
				keyRangeWildcard|85|flags, 85)
			ipAddressParseData.setBase85(true)
			return true, nil
		}
	}
	return false, nil
}

func chooseMACAddressProvider(fromString *MACAddressString,
	validationOptions addrstrparam.MACAddressStringParams, pa *parsedMACAddress,
	addressParseData *addressParseData) (res macAddressProvider, err addrerr.AddressStringError) {
	if addressParseData.isProvidingEmpty() {
		if validationOptions == defaultMACAddrParameters {
			res = defaultMACAddressEmptyProvider
		} else {
			res = macAddressEmptyProvider{macAddressNullProvider{validationOptions}}
		}
	} else if addressParseData.isAll() {
		if validationOptions == defaultMACAddrParameters {
			res = macAddressDefaultAllProvider
		} else {
			res = &macAddressAllProvider{validationOptions: validationOptions, creationLock: &sync.Mutex{}}
		}
	} else {
		if err = checkMACSegments(fromString.str, validationOptions, pa); err == nil {
			res = pa
		}
	}
	return
}

var maskCache = [3][IPv6BitCount + 1]*maskCreator{}

var loopbackCache = newEmptyAddrCreator(defaultIPAddrParameters, NoZone)

func chooseIPAddressProvider(
	originator HostIdentifierString,
	fullAddr string,
	validationOptions addrstrparam.IPAddressStringParams,
	parseData *parsedIPAddress) (res ipAddressProvider, err addrerr.AddressStringError) {
	qualifier := parseData.getQualifier()
	version := parseData.getProviderIPVersion()
	if version.IsIndeterminate() {
		version = qualifier.inferVersion(validationOptions) // checks whether a mask, prefix length, or zone makes the version clear
		optionsVersion := inferVersion(validationOptions)   // checks whether IPv4 or IPv6 is disallowed
		if version.IsIndeterminate() {
			version = optionsVersion
			parseData.setVersion(version)
		} else if !optionsVersion.IsIndeterminate() && version != optionsVersion {
			var key string
			if version.IsIPv6() {
				key = "ipaddress.error.ipv6"
			} else {
				key = "ipaddress.error.ipv4"
			}
			err = &addressStringError{addressError{str: fullAddr, key: key}}
			return
		}
		addressParseData := parseData.getAddressParseData()
		if addressParseData.isProvidingEmpty() {
			networkPrefixLength := qualifier.getNetworkPrefixLen()
			if networkPrefixLength != nil {
				if version.IsIndeterminate() {
					version = IPVersion(validationOptions.GetPreferredVersion())
				}
				prefLen := networkPrefixLength.bitCount()
				if validationOptions == defaultIPAddrParameters && prefLen <= IPv6BitCount {
					index := 0
					if version.IsIPv4() {
						index = 1
					} else if version.IsIPv6() {
						index = 2
					}
					creator := (*maskCreator)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&maskCache[index][prefLen]))))
					if creator == nil {
						creator = newMaskCreator(defaultIPAddrParameters, version, networkPrefixLength)
						dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&maskCache[index][prefLen]))
						atomic.StorePointer(dataLoc, unsafe.Pointer(creator))
					}
					res = creator
					return
				}
				res = newMaskCreator(validationOptions, version, networkPrefixLength)
				return
			} else {
				emptyOpt := validationOptions.EmptyStrParsedAs()
				if emptyOpt == addrstrparam.LoopbackOption || emptyOpt == addrstrparam.ZeroAddressOption {
					result := getLoopbackCreator(validationOptions, qualifier)
					if result != nil {
						res = result
					}
					return
				}

				if validationOptions == defaultIPAddrParameters {
					res = emptyProvider
				} else {
					res = &nullProvider{isEmpty: true, ipType: emptyType, params: validationOptions}
				}
				return
			}
		} else { //isAll
			// Before reaching here, we already checked whether we allow "all", "allowAll".
			if version.IsIndeterminate() && // version not inferred, nor is a particular version disallowed
				validationOptions.AllStrParsedAs() == addrstrparam.AllPreferredIPVersion {
				preferredVersion := IPVersion(validationOptions.GetPreferredVersion())
				if !preferredVersion.IsIndeterminate() {
					var formatParams addrstrparam.IPAddressStringFormatParams
					if preferredVersion.IsIPv6() {
						formatParams = validationOptions.GetIPv6Params()
					} else {
						formatParams = validationOptions.GetIPv4Params()
					}
					if formatParams.AllowsWildcardedSeparator() {
						version = preferredVersion
					}
				}
			}
			res = newAllCreator(qualifier, version, originator, validationOptions)
			return
		}
	} else {
		if parseData.isZoned() && version.IsIPv4() {
			err = &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.only.ipv6.has.zone"}}
			return
		}
		if err = checkSegments(fullAddr, validationOptions, parseData.getIPAddressParseData()); err == nil {
			res = parseData
		}
	}
	return
}

func inferVersion(params addrstrparam.IPAddressStringParams) IPVersion {
	if params.AllowsIPv6() {
		if !params.AllowsIPv4() {
			return IPv6
		}
	} else if params.AllowsIPv4() {
		return IPv4
	}
	return IndeterminateIPVersion
}

func getLoopbackCreator(validationOptions addrstrparam.IPAddressStringParams, qualifier *parsedHostIdentifierStringQualifier) (res *emptyAddrCreator) {
	zone := qualifier.getZone()
	defaultParams := defaultIPAddrParameters
	if validationOptions == defaultParams && zone == NoZone {
		res = loopbackCache
		if res == nil {
			res = newEmptyAddrCreator(defaultIPAddrParameters, NoZone)
		}
		return
	}
	res = newEmptyAddrCreator(validationOptions, zone)
	return
}

func checkSegmentMaxValues(
	fullAddr string,
	parseData *addressParseData,
	segmentIndex int,
	params addrstrparam.AddressStringFormatParams,
	maxValue uint64,
	maxDigitCount,
	maxUpperDigitCount int) addrerr.AddressStringError {
	if parseData.getFlag(segmentIndex, keySingleWildcard) {
		value := parseData.getValue(segmentIndex, keyLower)
		if value > maxValue {
			return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipv4.segment.too.large"}}
		}
		if parseData.getValue(segmentIndex, keyUpper) > maxValue {
			parseData.setValue(segmentIndex, keyUpper, maxValue)
		}
		if !params.AllowsUnlimitedLeadingZeros() {
			lowerRadix := parseData.getRadix(segmentIndex, keyLowerRadixIndex)
			if parseData.getIndex(segmentIndex, keyLowerStrEndIndex)-parseData.getIndex(segmentIndex, keyLowerStrDigitsIndex)-getStringPrefixCharCount(lowerRadix) > maxDigitCount {
				return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.segment.too.long"}}
			}
		}
	} else {
		value := parseData.getValue(segmentIndex, keyUpper)
		if value > maxValue {
			return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipv4.segment.too.large"}}
		}
		if !params.AllowsUnlimitedLeadingZeros() {
			lowerRadix := parseData.getRadix(segmentIndex, keyLowerRadixIndex)
			lowerEndIndex := parseData.getIndex(segmentIndex, keyLowerStrEndIndex)
			upperEndIndex := parseData.getIndex(segmentIndex, keyUpperStrEndIndex)
			if lowerEndIndex-parseData.getIndex(segmentIndex, keyLowerStrDigitsIndex)-getStringPrefixCharCount(lowerRadix) > maxDigitCount {
				return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.segment.too.long"}}
			}
			if lowerEndIndex != upperEndIndex {
				upperRadix := parseData.getRadix(segmentIndex, keyUpperRadixIndex)
				if upperEndIndex-parseData.getIndex(segmentIndex, keyUpperStrDigitsIndex)-getStringPrefixCharCount(upperRadix) > maxUpperDigitCount {
					return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.segment.too.long"}}
				}
			}
		}
	}
	return nil
}

func checkMACSegments(
	fullAddr string,
	validationOptions addrstrparam.MACAddressStringParams,
	parseData *parsedMACAddress) addrerr.AddressStringError {
	var err addrerr.AddressStringError
	format := parseData.getFormat()
	if format != unknownFormat {
		addressParseData := parseData.getAddressParseData()
		hasWildcardSeparator := addressParseData.hasWildcard() && validationOptions.GetFormatParams().AllowsWildcardedSeparator()
		//note that too many segments is checked inside the general parsing method
		segCount := addressParseData.getSegmentCount()
		if format == dotted {
			if segCount <= MediaAccessControlDottedSegmentCount && validationOptions.GetPreferredLen() != addrstrparam.EUI64Len {
				if !hasWildcardSeparator && segCount != MediaAccessControlDottedSegmentCount {
					return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.too.few.segments"}}
				}
			} else if !hasWildcardSeparator && segCount < MediaAccessControlDotted64SegmentCount {
				return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.too.few.segments"}}
			} else {
				parseData.setExtended(true)
			}
		} else if segCount > 2 {
			if segCount <= MediaAccessControlSegmentCount && validationOptions.GetPreferredLen() != addrstrparam.EUI64Len {
				if !hasWildcardSeparator && segCount != MediaAccessControlSegmentCount {
					return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.too.few.segments"}}
				}
			} else if !hasWildcardSeparator && segCount < ExtendedUniqueIdentifier64SegmentCount {
				return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.too.few.segments"}}
			} else {
				parseData.setExtended(true)
			}
			// we do not check char counts in the main parsing code for dashed, because we allow both
			// aabbcc-ddeeff and aa-bb-cc-dd-ee-ff, so we defer to the check until here
			if parseData.getFormat() == dashed {
				for i := 0; i < segCount; i++ {
					err = checkSegmentMaxValues(
						fullAddr,
						addressParseData,
						i,
						validationOptions.GetFormatParams(),
						MACMaxValuePerSegment,
						MACSegmentMaxChars,
						MACSegmentMaxChars)
					if err != nil {
						return err
					}
				}
			}
		} else {
			if parseData.getFormat() == dashed {
				//for single segment, we have already counted the exact number of hex digits
				//for double segment, we have already counted the exact number of hex digits in some cases and not others.
				//Basically, for address like a-b we have already counted the exact number of hex digits,
				//for addresses starting with a|b- or a-b| we have not,
				//but rather than figure out which are checked and which are not it's just as quick to check them all here
				if parseData.isDoubleSegment() {
					params := validationOptions.GetFormatParams()
					err = checkSegmentMaxValues(fullAddr, addressParseData, 0, params, macMaxTriple, macDoubleSegmentDigitCount, macDoubleSegmentDigitCount)
					if err != nil {
						return err
					}
					if parseData.isExtended() {
						err = checkSegmentMaxValues(fullAddr, addressParseData, 1, params, macMaxQuintuple, macExtendedDoubleSegmentDigitCount, macExtendedDoubleSegmentDigitCount)
					} else {
						err = checkSegmentMaxValues(fullAddr, addressParseData, 1, params, macMaxTriple, macDoubleSegmentDigitCount, macDoubleSegmentDigitCount)
					}
					if err != nil {
						return err
					}
				}
			} else if !hasWildcardSeparator {
				return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.too.few.segments"}}
			}
			if validationOptions.GetPreferredLen() == addrstrparam.EUI64Len {
				parseData.setExtended(true)
			}
		}
	} //else single segment
	return nil
}

func checkSegments(
	fullAddr string,
	validationOptions addrstrparam.IPAddressStringParams,
	parseData *ipAddressParseData) addrerr.AddressStringError {
	addressParseData := parseData.getAddressParseData()
	segCount := addressParseData.getSegmentCount()
	version := parseData.getProviderIPVersion()
	if version.IsIPv4() {
		missingCount := IPv4SegmentCount - segCount
		ipv4Options := validationOptions.GetIPv4Params()
		hasWildcardSeparator := addressParseData.hasWildcard() && ipv4Options.AllowsWildcardedSeparator()

		//single segments are handled in the parsing code with the allowSingleSegment setting
		hasMissingSegs := false
		if missingCount > 0 {
			if segCount > 1 {
				if ipv4Options.Allows_inet_aton_joinedSegments() {
					hasMissingSegs = true
					parseData.set_inet_aton_joined(true)
				} else if !hasWildcardSeparator {
					return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipv4.too.few.segments"}}
				}
			} else {
				hasMissingSegs = ipv4Options.Allows_inet_aton_joinedSegments()
			}
		}

		//here we check whether values are too large
		notUnlimitedLength := !ipv4Options.AllowsUnlimitedLeadingZeros()
		for i := 0; i < segCount; i++ {
			var max uint64
			if hasMissingSegs && i == segCount-1 {
				max = getMaxIPv4Value(missingCount + 1)
				if addressParseData.isInferredUpperBoundary(i) {
					parseData.setValue(i, keyUpper, max)
					continue
				}
			} else {
				max = IPv4MaxValuePerSegment
			}
			if parseData.getFlag(i, keySingleWildcard) {
				value := parseData.getValue(i, keyLower)
				if value > max {
					return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipv4.segment.too.large"}}
				}
				if parseData.getValue(i, keyUpper) > max {
					parseData.setValue(i, keyUpper, max)
				}
				if notUnlimitedLength {
					lowerRadix := addressParseData.getRadix(i, keyLowerRadixIndex)
					maxDigitCount := getMaxIPv4StringLength(missingCount, lowerRadix)
					if parseData.getIndex(i, keyLowerStrEndIndex)-parseData.getIndex(i, keyLowerStrDigitsIndex)-getStringPrefixCharCount(lowerRadix) > maxDigitCount {
						return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.segment.too.long"}}
					}
				}
			} else {
				value := parseData.getValue(i, keyUpper)
				if value > max {
					return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.ipv4.segment.too.large"}}
				}
				if notUnlimitedLength {
					lowerRadix := addressParseData.getRadix(i, keyLowerRadixIndex)
					maxDigitCount := getMaxIPv4StringLength(missingCount, lowerRadix)
					lowerEndIndex := parseData.getIndex(i, keyLowerStrEndIndex)
					upperEndIndex := parseData.getIndex(i, keyUpperStrEndIndex)
					if lowerEndIndex-parseData.getIndex(i, keyLowerStrDigitsIndex)-getStringPrefixCharCount(lowerRadix) > maxDigitCount {
						return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.segment.too.long"}}
					}
					if lowerEndIndex != upperEndIndex {
						upperRadix := parseData.getRadix(i, keyUpperRadixIndex)
						maxUpperDigitCount := getMaxIPv4StringLength(missingCount, upperRadix)
						if upperEndIndex-parseData.getIndex(i, keyUpperStrDigitsIndex)-getStringPrefixCharCount(upperRadix) > maxUpperDigitCount {
							return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.segment.too.long"}}
						}
					}
				}
			}
		}
	} else {
		totalSegmentCount := segCount
		if parseData.isProvidingMixedIPv6() {
			totalSegmentCount += IPv6MixedReplacedSegmentCount
		}
		hasWildcardSeparator := addressParseData.hasWildcard() && validationOptions.GetIPv6Params().AllowsWildcardedSeparator()
		if !hasWildcardSeparator && totalSegmentCount != 1 && totalSegmentCount < IPv6SegmentCount && !parseData.isCompressed() {
			return &addressStringError{addressError{str: fullAddr, key: "ipaddress.error.too.few.segments"}}
		}
	}
	return nil
}

func checkSingleWildcard(str string, start, end, digitsEnd int, options addrstrparam.AddressStringFormatParams) addrerr.AddressStringError {
	_ = start
	if !options.GetRangeParams().AllowsSingleWildcard() {
		return &addressStringError{addressError{str: str, key: "ipaddress.error.no.single.wildcard"}}
	}
	for k := digitsEnd; k < end; k++ {
		if str[k] != SegmentSqlSingleWildcard {
			return &addressStringError{addressError{str: str, key: "ipaddress.error.single.wildcard.order"}}
		}
	}
	return nil
}

func switchSingleWildcard10(currentValueHex uint64, s string, start, end, numSingleWildcards int, parseData *addressParseData, parsedSegIndex, leadingZeroStartIndex int, options addrstrparam.AddressStringFormatParams) (err addrerr.AddressStringError) {
	digitsEnd := end - numSingleWildcards
	err = checkSingleWildcard(s, start, end, digitsEnd, options)
	if err != nil {
		return
	}
	var lower uint64
	if start < digitsEnd {
		lower, err = switchValue10(currentValueHex, s, digitsEnd-start)
		if err != nil {
			return
		}
	}
	var upper uint64

	switch numSingleWildcards {
	case 1:
		lower *= 10
		upper = lower + 9
	case 2:
		lower *= 100
		upper = lower + 99
	case 3:
		lower *= 1000
		upper = lower + 999
	default:
		power := uint64(math.Pow10(numSingleWildcards))
		lower *= power
		upper = lower + power - 1
	}
	var radix uint32 = 10
	assign6Attributes2Values2Flags(start, end, leadingZeroStartIndex, start, end, leadingZeroStartIndex,
		parseData, parsedSegIndex, lower, upper, keySingleWildcard|radix, radix)
	return
}

func switchSingleWildcard2(currentValueHex uint64, s string, start, end, numSingleWildcards int, parseData *addressParseData, parsedSegIndex, leadingZeroStartIndex int, options addrstrparam.AddressStringFormatParams) (err addrerr.AddressStringError) {
	digitsEnd := end - numSingleWildcards
	err = checkSingleWildcard(s, start, end, digitsEnd, options)
	if err != nil {
		return
	}
	var lower, upper uint64
	if start < digitsEnd {
		lower, err = switchValue2(currentValueHex, s, digitsEnd-start)
		if err != nil {
			return
		}
	} else {
		lower = 0
	}
	lower <<= uint(numSingleWildcards)
	switch numSingleWildcards {
	case 1:
		upper = lower | 0x1
	case 2:
		upper = lower | 0x3
	case 3:
		upper = lower | 0x7
	case 4:
		upper = lower | 0xf
	case 5:
		upper = lower | 0x1f
	case 6:
		upper = lower | 0x3f
	case 7:
		upper = lower | 0x7f
	case 8:
		upper = lower | 0xff
	case 9:
		upper = lower | 0x1ff
	case 10:
		upper = lower | 0x3ff
	case 11:
		upper = lower | 0x7ff
	case 12:
		upper = lower | 0xfff
	case 13:
		upper = lower | 0x1fff
	case 14:
		upper = lower | 0x3fff
	case 15:
		upper = lower | 0x7fff
	case 16:
		upper = lower | 0xffff
	default:
		upper = lower | ^(^uint64(0) << uint(numSingleWildcards))
	}
	var radix uint32 = 2
	assign6Attributes2Values2Flags(start, end, leadingZeroStartIndex, start, end, leadingZeroStartIndex,
		parseData, parsedSegIndex, lower, upper, keySingleWildcard|radix, radix)
	return
}

func switchSingleWildcard8(currentValueHex uint64, s string, start, end, numSingleWildcards int, parseData *addressParseData, parsedSegIndex, leadingZeroStartIndex int, options addrstrparam.AddressStringFormatParams) (err addrerr.AddressStringError) {
	digitsEnd := end - numSingleWildcards
	err = checkSingleWildcard(s, start, end, digitsEnd, options)
	if err != nil {
		return
	}
	var lower, upper uint64
	if start < digitsEnd {
		lower, err = switchValue8(currentValueHex, s, digitsEnd-start)
		if err != nil {
			return
		}
	}
	switch numSingleWildcards {
	case 1:
		lower <<= 3
		upper = lower | 07
	case 2:
		lower <<= 6
		upper = lower | 077
	case 3:
		lower <<= 9
		upper = lower | 0777
	default:
		shift := numSingleWildcards * 3
		lower <<= uint(shift)
		upper = lower | ^(^uint64(0) << uint(shift))
	}
	var radix uint32 = 8
	assign6Attributes2Values2Flags(start, end, leadingZeroStartIndex, start, end, leadingZeroStartIndex,
		parseData, parsedSegIndex, lower, upper, keySingleWildcard|radix, radix)
	return
}

func assignSingleWildcard16(lower uint64, s string, start, end, numSingleWildcards int, parseData *addressParseData, parsedSegIndex, leadingZeroStartIndex int, options addrstrparam.AddressStringFormatParams) (err addrerr.AddressStringError) {
	digitsEnd := end - numSingleWildcards
	err = checkSingleWildcard(s, start, end, digitsEnd, options)
	if err != nil {
		return
	}
	shift := numSingleWildcards << 2
	lower <<= uint(shift)
	upper := lower | ^(^uint64(0) << uint(shift))
	assign6Attributes2Values1Flags(start, end, leadingZeroStartIndex, start, end, leadingZeroStartIndex,
		parseData, parsedSegIndex, lower, upper, keySingleWildcard)
	return
}

func parseSingleSegmentSingleWildcard16(currentValueHex uint64, s string, start, end, numSingleWildcards int, parseData *addressParseData, parsedSegIndex, leadingZeroStartIndex int, options addrstrparam.AddressStringFormatParams) (err addrerr.AddressStringError) {
	digitsEnd := end - numSingleWildcards
	err = checkSingleWildcard(s, start, end, digitsEnd, options)
	if err != nil {
		return
	}
	var upper, lower, extendedLower, extendedUpper uint64
	if numSingleWildcards < longHexDigits {
		midIndex := end - longHexDigits
		lower = parseLong16(s, midIndex, digitsEnd)
		shift := numSingleWildcards << 2
		lower <<= uint(shift)
		upper = lower | ^(^uint64(0) << uint(shift))
		extendedLower = parseLong16(s, start, midIndex)
		extendedUpper = extendedLower
	} else if numSingleWildcards == longHexDigits {
		lower = 0
		upper = 0xffffffffffffffff
		extendedUpper = currentValueHex
		extendedLower = currentValueHex
	} else {
		lower = 0
		upper = 0xffffffffffffffff
		extendedLower = currentValueHex
		shift := (numSingleWildcards - longHexDigits) << 2
		extendedLower <<= uint(shift)
		extendedUpper = extendedLower | ^(^uint64(0) << uint(shift))
	}
	assign6Attributes4Values1Flags(start, end, leadingZeroStartIndex, start, end, leadingZeroStartIndex,
		parseData, parsedSegIndex, lower, extendedLower, upper, extendedUpper, keySingleWildcard)
	return
}

func parseSingleSegmentSingleWildcard2(s string, start, end, numSingleWildcards int, parseData *addressParseData, parsedSegIndex, leadingZeroStartIndex int, options addrstrparam.AddressStringFormatParams) (err addrerr.AddressStringError) {
	digitsEnd := end - numSingleWildcards
	err = checkSingleWildcard(s, start, end, digitsEnd, options)
	if err != nil {
		return
	}
	var upper, lower, extendedLower, extendedUpper uint64
	midIndex := end - longBinaryDigits
	if numSingleWildcards < longBinaryDigits {
		lower = parseLong2(s, midIndex, digitsEnd)
		shift := numSingleWildcards
		lower <<= uint(shift)
		upper = lower | ^(^uint64(0) << uint(shift))
		extendedLower = parseLong2(s, start, midIndex)
		extendedUpper = extendedLower
	} else if numSingleWildcards == longBinaryDigits {
		upper = 0xffffffffffffffff
		extendedLower = parseLong2(s, start, midIndex)
		extendedUpper = extendedLower
	} else {
		upper = 0xffffffffffffffff
		shift := numSingleWildcards - longBinaryDigits
		extendedLower = parseLong2(s, start, midIndex-shift)
		extendedLower <<= uint(shift)
		extendedUpper = extendedLower | ^(^uint64(0) << uint(shift))
	}
	assign6Attributes4Values1Flags(start, end, leadingZeroStartIndex, start, end, leadingZeroStartIndex,
		parseData, parsedSegIndex, lower, extendedLower, upper, extendedUpper, keySingleWildcard)
	return
}

////////////////////////

var maxValues = [5]uint64{0, IPv4MaxValuePerSegment, 0xffff, 0xffffff, 0xffffffff}

func getMaxIPv4Value(segmentCount int) uint64 {
	return maxValues[segmentCount]
}

func getStringPrefixCharCount(radix uint32) int {
	switch radix {
	case 10:
		return 0
	case 16:
	case 2:
		return 2
	default:
	}
	return 1
}

var maxIPv4StringLen = [9][]int{ //indices are [radix / 2][additionalSegments], and we handle radices 8, 10, 16
	{3, 6, 8, 11},   //no radix supplied we treat as octal, the longest
	{8, 16, 24, 32}, // binary
	{}, {},
	{3, 6, 8, 11},                   //octal: 0377, 0177777, 077777777, 037777777777
	{IPv4SegmentMaxChars, 5, 8, 10}, //decimal: 255, 65535, 16777215, 4294967295
	{}, {},
	{2, 4, 6, 8}, //hex: 0xff, 0xffff, 0xffffff, 0xffffffff
}

func getMaxIPv4StringLength(additionalSegmentsCovered int, radix uint32) int {
	radixHalved := radix >> 1
	if radixHalved < uint32(len(maxIPv4StringLen)) {
		sl := maxIPv4StringLen[radixHalved]
		if additionalSegmentsCovered >= 0 && additionalSegmentsCovered < len(sl) {
			return sl[additionalSegmentsCovered]
		}
	}
	return 0
}

func switchValue2(currentHexValue uint64, s string, digitCount int) (result uint64, err addrerr.AddressStringError) {
	result = 0xf & currentHexValue
	if result > 1 {
		err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.binary.digit"}}
		return
	}
	shift := 0

	for digitCount--; digitCount > 0; digitCount-- {
		shift++
		currentHexValue >>= 4
		next := 0xf & currentHexValue
		if next >= 1 {
			if next == 1 {
				result |= 1 << uint(shift)
			} else {
				err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.binary.digit"}}
				return
			}
		}
	}
	return
}

/**
 * The digits were stored as a hex value, this switches them to an octal value.
 *
 * @param currentHexValue
 * @param digitCount
 * @return
 */
func switchValue8(currentHexValue uint64, s string, digitCount int) (result uint64, err addrerr.AddressStringError) {
	result = 0xf & currentHexValue
	if result >= 8 {
		err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.octal.digit"}}
		return
	}
	shift := 0
	for digitCount--; digitCount > 0; digitCount-- {
		shift += 3
		currentHexValue >>= 4
		next := 0xf & currentHexValue
		if next >= 8 {
			err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.octal.digit"}}
			return
		}
		result |= next << uint(shift)
	}
	return
}

func switchValue10(currentHexValue uint64, s string, digitCount int) (result uint64, err addrerr.AddressStringError) {
	result = 0xf & currentHexValue
	if result >= 10 {
		err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.decimal.digit"}}
		return
	}
	digitCount--
	if digitCount > 0 {
		factor := uint64(10)
		for {
			currentHexValue >>= 4
			next := 0xf & currentHexValue
			if next >= 10 {
				err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.decimal.digit"}}
				return
			}
			result += next * factor
			digitCount--
			if digitCount == 0 {
				break
			}
			if factor == 10 {
				factor = 100
			} else if factor == 100 {
				factor = 1000
			} else {
				factor *= 10
			}
		}
	}
	return
}

func parseLong2(s string, start, end int) uint64 {
	charArray := chars
	result := uint64(charArray[s[start]])
	for start++; start < end; start++ {
		c := s[start]
		if c == '1' {
			result = (result << 1) | 1
		} else {
			result <<= 1
		}
	}
	return result
}

func parseInt8(s string, start, end int) (result uint32, err addrerr.AddressStringError) {
	charArray := chars
	result = uint32(charArray[s[start]])
	if result >= 8 {
		err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.octal.digit"}}
		return
	}
	for start++; start < end; start++ {
		next := uint32(charArray[s[start]])
		if next >= 8 {
			err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.octal.digit"}}
			return
		}
		result = (result << 3) | next
	}
	return
}

func parseLong8(s string, start, end int) uint64 {
	charArray := chars
	result := uint64(charArray[s[start]])
	for start++; start < end; start++ {
		result = (result << 3) | uint64(charArray[s[start]])
	}
	return result
}

func parseInt10(s string, start, end int) (result uint32, err addrerr.AddressStringError) {
	charArray := chars
	result = uint32(charArray[s[start]])
	if result >= 10 {
		err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.decimal.digit"}}
		return
	}
	for start++; start < end; start++ {
		next := uint32(charArray[s[start]])
		if next >= 10 {
			err = &addressStringError{addressError{str: s, key: "ipaddress.error.ipv4.invalid.decimal.digit"}}
			return
		}
		result = (result * 10) + next
	}
	return
}

func parseLong10(s string, start, end int) uint64 {
	charArray := chars
	result := uint64(charArray[s[start]])
	for start++; start < end; start++ {
		result = (result * 10) + uint64(charArray[s[start]])
	}
	return result
}

func parseInt16(s string, start, end int) uint32 {
	charArray := chars
	result := uint32(charArray[s[start]])
	for start++; start < end; start++ {
		result = (result << 4) | uint32(charArray[s[start]])
	}
	return result
}

func parseLong16(s string, start, end int) uint64 {
	charArray := chars
	result := uint64(charArray[s[start]])
	for start++; start < end; start++ {
		result = (result << 4) | uint64(charArray[s[start]])
	}
	return result
}

var base85Powers = createBase85Powers()

func createBase85Powers() []big.Int {
	res := make([]big.Int, 10)
	eightyFive := big.NewInt(85)
	res[0].SetUint64(1)
	for i := 1; i < len(res); i++ {
		res[i].Mul(&res[i-1], eightyFive)
	}
	return res
}

func parse85(s string, start, end int) *big.Int {
	charArray := extendedChars
	var result big.Int
	var last bool
	for {
		var partialEnd, power int
		left := end - start
		if last = left <= 9; last {
			partialEnd = end
			power = left
		} else {
			partialEnd = start + 9
			power = 9
		}
		var partialResult = uint64(charArray[s[start]])
		for start++; start < partialEnd; start++ {
			next := charArray[s[start]]
			partialResult = (partialResult * 85) + uint64(next)
		}
		result.Mul(&result, &base85Powers[power]).Add(&result, bigZero().SetUint64(partialResult))
		if last {
			break
		}
	}
	return &result
}

//according to rfc 1035 or 952, a label must start with a letter, must end with a letter or digit, and must have in the middle a letter or digit or -
//rfc 1123 relaxed that to allow labels to start with a digit, section 2.1 has a discussion on this.  It states that the highest level component name must be alphabetic - referring to .com or .net or whatever.
//furthermore, the underscore has become generally acceptable, as indicated in rfc 2181
//there is actually a distinction between host names and domain names.  a host name is a specific type of domain name identifying hosts.
//hosts are not supposed to have the underscore.

//en.wikipedia.org/wiki/Domain_Name_System#Domain_name_syntax
//en.wikipedia.org/wiki/Hostname#Restrictions_on_valid_host_names

//max length is 63, cannot start or end with hyphen
//strictly speaking, the underscore is not allowed anywhere, but it seems that rule is sometimes broken
//also, underscores seem to be a part of dns names that are not part of host names, so we allow it here to be safe

//networkadminkb.com/KB/a156/windows-2003-dns-and-the-underscore.aspx

//It's a little confusing.  rfc 2181 https://www.ietf.org/rfc/rfc2181.txt in section 11 on name syntax says that any chars are allowed in dns.
//However, it also says internet host names might have restrictions of their own, and this was defined in rfc 1035.
//rfc 1035 defines the restrictions on internet host names, in section 2.3.1 http://www.ietf.org/rfc/rfc1035.txt

//So we will follow rfc 1035 and in addition allow the underscore.

var (
	ipvFutureUppercase = byte(unicode.ToUpper(rune(IPvFuture)))
	defaultEmptyHost   = &parsedHost{}
)

func (strValidator) validateHostName(fromHost *HostName, validationOptions addrstrparam.HostNameParams) (psdHost *parsedHost, err addrerr.HostNameError) {
	str := fromHost.str
	addrLen := len(str)
	if addrLen > maxHostLength {
		if addrLen > maxHostLength+1 || str[maxHostLength] != LabelSeparator {
			err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.length"}}
			return
		}
	}
	var segmentUppercase, isNotNormalized, squareBracketed,
		tryIPv6, tryIPv4,
		isPrefixed, hasPortOrService, hostIsEmpty bool
	isAllDigits, isPossiblyIPv6, isPossiblyIPv4 := true, true, true
	isSpecialOnlyIndex, qualifierIndex, index, lastSeparatorIndex := -1, -1, -1, -1
	labelCount := 0
	maxLocalLabels := 6 //should be at least 4 to avoid the array for ipv4 addresses
	var separatorIndices []int
	var normalizedFlags []bool

	sep0, sep1, sep2, sep3, sep4, sep5 := -1, -1, -1, -1, -1, -1
	var isUpper0, isUpper1, isUpper2, isUpper3, isUpper4, isUpper5 bool

	var currentChar byte
	for index++; index <= addrLen; index++ {

		//grab the character to evaluate
		if index == addrLen {
			if index == 0 {
				hostIsEmpty = true
				break
			}
			segmentCountMatchesIPv4 :=
				isPossiblyIPv4 &&
					(labelCount+1 == IPv4SegmentCount) ||
					(labelCount+1 < IPv4SegmentCount && isSpecialOnlyIndex >= 0) ||
					(labelCount+1 < IPv4SegmentCount && validationOptions.GetIPAddressParams().GetIPv4Params().Allows_inet_aton_joinedSegments()) ||
					labelCount == 0 && validationOptions.GetIPAddressParams().AllowsSingleSegment()
			if isAllDigits {
				if isPossiblyIPv4 && segmentCountMatchesIPv4 {
					tryIPv4 = true
					break
				}
				isPossiblyIPv4 = false
				if hasPortOrService && isPossiblyIPv6 { //isPossiblyIPv6 is already false if labelCount > 0
					//since it is all digits, it cannot be host, so we set tryIPv6 rather than just isPossiblyIPv6
					tryIPv6 = true
					break
				}
				err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid"}}
				return
			}
			isPossiblyIPv4 = isPossiblyIPv4 && segmentCountMatchesIPv4
			currentChar = LabelSeparator
		} else {
			currentChar = str[index]
		}

		//check that character
		//we break out of the loop if we hit '[', '*', '%' (as zone or wildcard), or ':' that is not interpreted as port (and this is ipv6)
		//we exit the loop prematurely if we hit '/' or ':' interpreted as port
		if currentChar >= 'a' && currentChar <= 'z' {
			if currentChar > 'f' {
				isPossiblyIPv6 = false
				isPossiblyIPv4 = isPossiblyIPv4 && (currentChar == 'x' && validationOptions.GetIPAddressParams().GetIPv4Params().Allows_inet_aton_hex())
			} else if currentChar == 'b' {
				isPossiblyIPv4 = isPossiblyIPv4 && validationOptions.GetIPAddressParams().GetIPv4Params().AllowsBinary()
			}
			isAllDigits = false
		} else if currentChar >= '0' && currentChar <= '9' {
			//nothing to do
			continue
		} else if currentChar >= 'A' && currentChar <= 'Z' {
			if currentChar > 'F' {
				isPossiblyIPv6 = false
				isPossiblyIPv4 = isPossiblyIPv4 && (currentChar == 'X' && validationOptions.GetIPAddressParams().GetIPv4Params().Allows_inet_aton_hex())
			} else if currentChar == 'B' {
				isPossiblyIPv4 = isPossiblyIPv4 && validationOptions.GetIPAddressParams().GetIPv4Params().AllowsBinary()
			}
			segmentUppercase = true
			isAllDigits = false
		} else if currentChar == LabelSeparator {
			length := index - lastSeparatorIndex - 1
			if length > maxLabelLength {
				err = &hostNameError{addressError{str: str, key: "ipaddress.error.segment.too.long"}}
				return
			} else if length == 0 {
				if index < addrLen {
					err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.segment.too.short"}}
					return
				}
				isPossiblyIPv4 = false
				isNotNormalized = true
			} else {
				if labelCount < maxLocalLabels {
					if labelCount < 3 {
						if labelCount == 0 {
							sep0 = index
							isUpper0 = segmentUppercase
						} else if labelCount == 1 {
							sep1 = index
							isUpper1 = segmentUppercase
						} else {
							sep2 = index
							isUpper2 = segmentUppercase
						}
					} else {
						if labelCount == 3 {
							sep3 = index
							isUpper3 = segmentUppercase
						} else if labelCount == 4 {
							sep4 = index
							isUpper4 = segmentUppercase
						} else {
							sep5 = index
							isUpper5 = segmentUppercase
						}
					}
					labelCount++
				} else if labelCount == maxLocalLabels {
					separatorIndices = make([]int, maxHostSegments+1)
					separatorIndices[labelCount] = index
					if validationOptions.NormalizesToLowercase() {
						normalizedFlags = make([]bool, maxHostSegments+1)
						normalizedFlags[labelCount] = !segmentUppercase
						isNotNormalized = isNotNormalized || segmentUppercase
					}
					labelCount++
				} else {
					separatorIndices[labelCount] = index
					if normalizedFlags != nil {
						normalizedFlags[labelCount] = !segmentUppercase
						isNotNormalized = isNotNormalized || segmentUppercase
					}
					labelCount++
					if labelCount > maxHostSegments {
						err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.too.many.segments"}}
						return
					}
				}
				segmentUppercase = false //this is per segment so reset it
			}
			lastSeparatorIndex = index
			isPossiblyIPv6 = isPossiblyIPv6 && (index == addrLen) //A '.' means not ipv6 (if we see ':' we jump out of loop so mixed address not possible), but for single segment we end up here even without a '.' character in the string
		} else if currentChar == '_' { //this is not supported in host names but is supported in domain names, see discussion in HostName
			isAllDigits = false
		} else if currentChar == '-' {
			//host name segments cannot end with '-'
			if index == lastSeparatorIndex+1 || index == addrLen-1 || str[index+1] == LabelSeparator {
				err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
				return
			}
			isAllDigits = false
		} else if currentChar == IPv6StartBracket {
			if index == 0 && labelCount == 0 && addrLen > 2 {
				squareBracketed = true
				break
			}
			err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
			return
		} else if currentChar == PrefixLenSeparator {
			isPrefixed = true
			qualifierIndex = index + 1
			addrLen = index
			isNotNormalized = true
			index--
		} else {
			a := currentChar == SegmentWildcard
			if a || currentChar == SegmentSqlSingleWildcard {
				b := !a
				addressOptions := validationOptions.GetIPAddressParams()
				if b && addressOptions.GetIPv6Params().AllowsZone() { //if we allow zones, we treat '%' as a zone and not as a wildcard
					if isPossiblyIPv6 && labelCount < IPv6SegmentCount {
						tryIPv6 = true
						isPossiblyIPv4 = false
						break
					}
					err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
					return
				} else {
					if isPossiblyIPv4 {
						if addressOptions.GetIPv4Params().GetRangeParams().AllowsWildcard() {
							if isSpecialOnlyIndex < 0 {
								isSpecialOnlyIndex = index
							}
						} else {
							isPossiblyIPv4 = false
						}
					}
					if isPossiblyIPv6 && addressOptions.GetIPv6Params().GetRangeParams().AllowsWildcard() {
						if isSpecialOnlyIndex < 0 {
							isSpecialOnlyIndex = index
						}
					} else {
						if !isPossiblyIPv4 {
							//needs to be either ipv4 or ipv6
							err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
							return
						}
						isPossiblyIPv6 = false
					}
				}
				isAllDigits = false
			} else if currentChar == IPv6SegmentSeparator { //also might denote a port
				if validationOptions.AllowsPort() || validationOptions.AllowsService() {
					hasPortOrService = true
					qualifierIndex = index + 1
					addrLen = index //causes loop to terminate, but only after handling the last segment
					isNotNormalized = true
					index--
				} else {
					isPossiblyIPv4 = false
					if isPossiblyIPv6 {
						tryIPv6 = true
						break
					}
					err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
					return
				}
			} else if currentChar == AlternativeRangeSeparatorStr[0] {
				//} else if currentChar == AlternativeRangeSeparator {
				if index+1 == addrLen {
					err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
					return
				}
				currentChar = str[index+1]
				if currentChar == AlternativeRangeSeparatorStr[1] {
					isAllDigits = false
					isPossiblyIPv4 = false
					isPossiblyIPv6 = false
					if isSpecialOnlyIndex < 0 {
						isSpecialOnlyIndex = index
					}
					index++
				} else {
					err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
					return
				}
			} else {
				err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, index}
				return
			}
		}
	}

	/*
		1. squareBracketed: [ addr ]
		2. tryIPv4 || tryIPv6: this is a string with characters that invalidate it as a host but it still may in fact be an address
			This includes ipv6 strings, ipv4/ipv6 strings with '*', or all dot/digit strings like 1.2.3.4 that are 4 segments
		3. isPossiblyIPv4: this is a string with digits, - and _ characters and the number of separators matches ipv4.  Such strings can also be valid hosts.
		  The range options flag (controlling whether we allow '-' or '_' in addresses) for ipv4 can control whether it is treated as host or address.
		  It also includes "" empty addresses.
		  isPossiblyIPv6: something like f:: or f:1, the former IPv6 and the latter a host "f" with port 1.  Such strings can be valid addresses or hosts.
		  If it parses as an address, we do not treat as host.
	*/
	psdHost = &parsedHost{originalStr: str, params: validationOptions}
	addressOptions := validationOptions.GetIPAddressParams()
	isIPAddress := squareBracketed || tryIPv4 || tryIPv6
	if !validationOptions.AllowsIPAddress() {
		if isIPAddress {
			err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.ipaddress"}}
			return
		}
	} else if isIPAddress || isPossiblyIPv4 || isPossiblyIPv6 {
		provider, addrErr, hostErr := func() (provider ipAddressProvider, addrErr addrerr.AddressError, hostErr addrerr.HostNameError) {
			pa := parsedIPAddress{
				ipAddressParseData: ipAddressParseData{addressParseData: addressParseData{str: str}},
				options:            addressOptions,
				originator:         fromHost,
			}
			hostQualifier := psdHost.getQualifier()
			if squareBracketed {
				//Note:
				//Firstly, we need to find the address end which is denoted by the end bracket
				//Secondly, while zones appear inside bracket, prefix or port appears outside, according to rfc 4038
				//So we keep track of the boolean endsWithPrefix to differentiate.
				endIndex := addrLen - 1
				endsWithQualifier := str[endIndex] != IPv6EndBracket
				if endsWithQualifier {
					for endIndex--; str[endIndex] != IPv6EndBracket; endIndex-- {
						if endIndex == 1 {
							err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.bracketed.missing.end"}}
							return
						}
					}
				}
				startIndex := 1
				if strings.HasPrefix(str[1:], SmtpIPv6Identifier) {
					//SMTP rfc 2821 allows [IPv6:ipv6address]
					startIndex = 6
				} else {
					/* RFC 3986 section 3.2.2
						host = IP-literal / IPv4address / reg-name
						IP-literal = "[" ( IPv6address / IPvFuture  ) "]"
						IPvFuture  = "v" 1*HEXDIG "." 1*( unreserved / sub-delims / ":" )
					If a URI containing an IP-literal that starts with "v" (case-insensitive),
					indicating that the version flag is present, is dereferenced by an application that does not know the meaning of that version flag,
					then the application should return an appropriate error for "address mechanism not supported".
					*/
					firstChar := str[1]
					if firstChar == IPvFuture || firstChar == ipvFutureUppercase {
						err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.mechanism"}}
						return
					}
				}
				addrErr = validateIPAddress(addressOptions, str, startIndex, endIndex, pa.getIPAddressParseData(), false)
				if addrErr != nil {
					return
				}
				if endsWithQualifier {
					//here we check what is in the qualifier that follows the bracket: prefix/mask or port?
					//if prefix/mask, we supply the qualifier to the address, otherwise we supply it to the host
					prefixIndex := endIndex + 1
					prefixChar := str[prefixIndex]
					if prefixChar == PrefixLenSeparator {
						isPrefixed = true
					} else if prefixChar == PortSeparator {
						hasPortOrService = true
					} else {
						hostErr = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, prefixIndex}
						return
					}
					qualifierIndex = prefixIndex + 1 //skip the ']/'
					endIndex = len(str)
					addressParseData := pa.getAddressParseData()
					addrErr = parseHostNameQualifier(
						str,
						addressOptions,
						validationOptions,
						hostQualifier,
						isPrefixed,
						hasPortOrService,
						addressParseData.isProvidingEmpty(),
						qualifierIndex,
						endIndex,
						pa.getProviderIPVersion())
					if addrErr != nil {
						return
					}
					insideBracketsQualifierIndex := pa.getQualifierIndex()
					if pa.isZoned() && str[insideBracketsQualifierIndex] == '2' &&
						str[insideBracketsQualifierIndex+1] == '5' {
						//handle %25 from rfc 6874
						insideBracketsQualifierIndex += 2
					}
					addrErr = parseHostAddressQualifier(
						str,
						addressOptions,
						nil,
						pa.hasPrefixSeparator(),
						false,
						pa.getIPAddressParseData(),
						insideBracketsQualifierIndex,
						prefixIndex-1)
					if addrErr != nil {
						return
					}
					if isPrefixed {
						// since we have an address, we apply the prefix to the address rather than to the host
						// rather than use the prefix as a host qualifier, we treat it as an address qualifier and leave the host qualifier as noQualifier
						// also, keep in mind you can combine prefix with zone like fe80::%2/64, see https://tools.ietf.org/html/rfc4007#section-11.7

						// if there are two prefix lengths, we choose the smaller (larger network)
						// if two masks, we combine them (if both network masks, this is the same as choosing smaller prefix)
						addrQualifier := pa.getIPAddressParseData().getQualifier()
						addrErr = addrQualifier.merge(hostQualifier)
						if addrErr != nil {
							return
						}
						hostQualifier.clearPrefixOrMask()
						// note it makes no sense to indicate a port or service with a prefix
					}
				} else {
					qualifierIndex = pa.getQualifierIndex()
					isPrefixed = pa.hasPrefixSeparator()
					hasPortOrService = false
					if pa.isZoned() && str[qualifierIndex] == '2' &&
						str[qualifierIndex+1] == '5' {
						//handle %25 from rfc 6874
						qualifierIndex += 2
					}
					addrErr = parseHostAddressQualifier(str, addressOptions, validationOptions, isPrefixed, hasPortOrService, pa.getIPAddressParseData(), qualifierIndex, endIndex)
					if addrErr != nil {
						return
					}
				}
				//SMTP rfc 2821 allows [ipv4address]
				version := pa.getProviderIPVersion()
				if !version.IsIPv6() && !validationOptions.AllowsBracketedIPv4() {
					err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.bracketed.not.ipv6"}}
					return
				}
			} else { //not square-bracketed
				/*
					there are cases where it can be ipv4 or ipv6, but not many
					any address with a '.' in it cannot be ipv6 at this point (if we hit a ':' first we would have jumped out of the loop)
					any address with a ':' has gone through tests to see if up until that point it could match an ipv4 address or an ipv6 address
					it can only be ipv4 if it has right number of segments, and only decimal digits.
					it can only be ipv6 if it has only hex digits.
					so when can it be both?  if it looks like *: at the start, so that it has the right number of segments for ipv4 but does not have a '.' invalidating ipv6
					so in that case we might have either something like *:1 for it to be ipv4 (ambiguous is treated as ipv4) or *:f:: to be ipv6
					So we validate the potential port (or ipv6 segment) to determine which one and then go from there
					Also, if it is single segment address that is all decimal digits.
				*/

				// We start by checking if there is potentially a port or service
				// if IPv6, we may need to try a :x as a port or service and as a trailing segment
				firstTrySucceeded := false
				hasAddressPortOrService := false
				addressQualifierIndex := -1
				isPotentiallyIPv6 := isPossiblyIPv6 || tryIPv6
				if isPotentiallyIPv6 {
					//find the last port separator, currently we point to the first one with qualifierIndex
					//note that the service we find here could be the ipv4 part of either an ipv6 address or ipv6 mask like this 1:2:3:4:5:6:1.2.3.4 or 1:2:3:4:5:6:1.2.3.4/1:2:3:4:5:6:1.2.3.4
					if !isPrefixed && (validationOptions.AllowsPort() || validationOptions.AllowsService()) {
						for j := len(str) - 1; j >= 0; j-- {
							c := str[j]
							if c == IPv6SegmentSeparator {
								hasAddressPortOrService = true
								addressQualifierIndex = j + 1
							} else if (c >= '0' && c <= '9') ||
								(c >= 'A' && c <= 'Z') ||
								(c >= 'a' && c <= 'z') ||
								(c == '-') ||
								(c == SegmentWildcard) {
								//see validateHostNamePort for more details on valid ports and service names
								continue
							}
							break
						}
					}
				} else {
					hasAddressPortOrService = hasPortOrService
					addressQualifierIndex = qualifierIndex
				}
				var endIndex int
				if hasAddressPortOrService {
					//validate the potential port
					addrErr = parsePortOrService(str, nil, validationOptions, hostQualifier, addressQualifierIndex, len(str))
					if addrErr != nil {
						//certainly not IPv4 since it doesn't qualify as port (see comment above)
						if !isPotentiallyIPv6 {
							//not IPv6 either, so we're done with checking for address
							return
						}
						// no need to call hostQualifier.clear() since parsePortOrService does not populate qualifier on error
						endIndex = len(str)
					} else if isPotentiallyIPv6 {
						//here it can be either a port or part of an IPv6 address, like this: fe80::6a05:caff:fe3:123
						expectPort := validationOptions.ExpectsPort()
						if expectPort {
							//try with port first, then try as IPv6 no port
							endIndex = addressQualifierIndex - 1
						} else {
							//try as IPv6 with no port first, try with port second
							endIndex = len(str)
						}
						//first try
						addrErr = validateIPAddress(addressOptions, str, 0, endIndex, pa.getIPAddressParseData(), false)
						if addrErr == nil {
							// Since no square brackets, we parse as an address (this can affect how zones are parsed).
							// Also, an address cannot end with a single ':' like a port, so we cannot take a shortcut here and parse for port, we must strip it off first (hence no host parameters passed)
							addrErr = parseAddressQualifier(str, addressOptions, nil, pa.getIPAddressParseData(), endIndex)
						}
						if firstTrySucceeded = addrErr == nil; !firstTrySucceeded {
							pa = parsedIPAddress{
								ipAddressParseData: ipAddressParseData{addressParseData: addressParseData{str: str}},
								options:            addressOptions,
								originator:         fromHost,
							}
							if expectPort {
								// we tried with port first, now we try as IPv6 no port
								hostQualifier.clearPortOrService()
								endIndex = len(str)
							} else {
								// we tried as IPv6 with no port first, now we try with port second
								endIndex = addressQualifierIndex - 1
							}
						} else if !expectPort {
							// it is an address
							// we tried with no port and succeeded, so clear the port, it was not a port
							hostQualifier.clearPortOrService()
						}
					} else {
						endIndex = addressQualifierIndex - 1
					}
				} else {
					endIndex = len(str)
				}
				if !firstTrySucceeded {
					if addrErr = validateIPAddress(addressOptions, str, 0, endIndex, pa.getIPAddressParseData(), false); addrErr == nil {
						//since no square brackets, we parse as an address (this can affect how zones are parsed)
						//Also, an address cannot end with a single ':' like a port, so we cannot take a shortcut here and parse for port, we must strip it off first (hence no host parameters passed)
						addrErr = parseAddressQualifier(str, addressOptions, nil, pa.getIPAddressParseData(), endIndex)
					}
					if addrErr != nil {
						return
					}
				}
			}
			// we successfully parsed an IP address
			provider, addrErr = chooseIPAddressProvider(fromHost, str, addressOptions, &pa)
			return
		}()
		if hostErr != nil {
			err = hostErr
			return
		}
		if addrErr != nil {
			if isIPAddress {
				err = &hostAddressNestedError{nested: addrErr}
				return
			}
			psdHost.labelsQualifier.clearPortOrService()
			//fall though and evaluate as a host
		} else {
			psdHost.embeddedAddress.addressProvider = provider
			return
		}
	}

	hostQualifier := psdHost.getQualifier()
	addrErr := parseHostNameQualifier(
		str,
		addressOptions,
		validationOptions,
		hostQualifier,
		isPrefixed,
		hasPortOrService,
		hostIsEmpty,
		qualifierIndex,
		len(str),
		IndeterminateIPVersion)
	if addrErr != nil {
		err = &hostAddressNestedError{nested: addrErr}
		return
	}
	if hostIsEmpty {
		if !validationOptions.AllowsEmpty() {
			err = &hostNameError{addressError{str: str, key: "ipaddress.host.error.empty"}}
			return
		}
		if *hostQualifier == defaultEmptyHost.labelsQualifier {
			psdHost = defaultEmptyHost
		}
	} else {
		if labelCount <= maxLocalLabels {
			maxLocalLabels = labelCount
			separatorIndices = make([]int, maxLocalLabels)
			if validationOptions.NormalizesToLowercase() {
				normalizedFlags = make([]bool, maxLocalLabels)
			}
		} else if labelCount != len(separatorIndices) {
			trimmedSeparatorIndices := make([]int, labelCount)
			copy(trimmedSeparatorIndices[maxLocalLabels:], separatorIndices[maxLocalLabels:labelCount])
			separatorIndices = trimmedSeparatorIndices
			if normalizedFlags != nil {
				trimmedNormalizedFlags := make([]bool, labelCount)
				copy(trimmedNormalizedFlags[maxLocalLabels:], normalizedFlags[maxLocalLabels:labelCount])
				normalizedFlags = trimmedNormalizedFlags
			}
		}
		for i := 0; i < maxLocalLabels; i++ {
			var nextSep int
			var isUpper bool
			if i < 2 {
				if i == 0 {
					nextSep = sep0
					isUpper = isUpper0
				} else {
					nextSep = sep1
					isUpper = isUpper1
				}
			} else if i < 4 {
				if i == 2 {
					nextSep = sep2
					isUpper = isUpper2
				} else {
					nextSep = sep3
					isUpper = isUpper3
				}
			} else if i == 4 {
				nextSep = sep4
				isUpper = isUpper4
			} else {
				nextSep = sep5
				isUpper = isUpper5
			}
			separatorIndices[i] = nextSep
			if normalizedFlags != nil {
				normalizedFlags[i] = !isUpper
				isNotNormalized = isNotNormalized || isUpper
			}
		}
		//We support a.b.com/24:80 (prefix and port combo)
		//or just port, or a service where-ever a port can appear
		//A prefix with port can mean a subnet of addresses using the same port everywhere (the subnet being the prefix block of the resolved address),
		//or just denote the prefix length of the resolved address along with a port

		//here we check what is in the qualifier that follows the bracket: prefix/mask or port?
		//if prefix/mask, we supply the qualifier to the address, otherwise we supply it to the host
		//also, it is possible the address has a zone
		var addrQualifier *parsedHostIdentifierStringQualifier
		if isPrefixed {
			addrQualifier = hostQualifier
		} else {
			addrQualifier = noQualifier
		}
		embeddedAddr := checkSpecialHosts(str, addrLen, addrQualifier)
		hasEmbeddedAddr := embeddedAddr.addressProvider != nil
		if isSpecialOnlyIndex >= 0 && (!hasEmbeddedAddr || embeddedAddr.addressStringError != nil) {
			if embeddedAddr.addressStringError != nil {
				err = &hostAddressNestedError{
					hostNameIndexError: hostNameIndexError{
						hostNameError: hostNameError{
							addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"},
						},
						index: isSpecialOnlyIndex,
					},
					nested: embeddedAddr.addressStringError,
				}
				return
			}
			err = &hostNameIndexError{hostNameError{addressError{str: str, key: "ipaddress.host.error.invalid.character.at.index"}}, isSpecialOnlyIndex}
			return
		}
		psdHost.separatorIndices = separatorIndices
		psdHost.normalizedFlags = normalizedFlags
		if !hasEmbeddedAddr {
			if !isNotNormalized {
				normalizedLabels := make([]string, len(separatorIndices))
				for i, lastSep := 0, -1; i < len(normalizedLabels); i++ {
					index := separatorIndices[i]
					normalizedLabels[i] = str[lastSep+1 : index]
					lastSep = index
				}
				psdHost.parsedHostCache = &parsedHostCache{
					host:             str,
					normalizedLabels: normalizedLabels,
				}
			}
		} else {
			if isPrefixed {
				psdHost.labelsQualifier.clearPrefixOrMask()
			}
			psdHost.embeddedAddress = embeddedAddr
		}
	}
	return
}

var defaultUncOpts = new(addrstrparam.IPAddressStringParamsBuilder).
	AllowIPv4(false).AllowEmpty(false).AllowMask(false).AllowPrefix(false).ToParams()

var reverseDNSIPv4Opts = new(addrstrparam.IPAddressStringParamsBuilder).
	AllowIPv6(false).AllowEmpty(false).AllowMask(false).AllowPrefix(false).
	GetIPv4AddressParamsBuilder().Allow_inet_aton(false).GetParentBuilder().ToParams()

var reverseDNSIPv6Opts = new(addrstrparam.IPAddressStringParamsBuilder).
	AllowIPv4(false).AllowEmpty(false).AllowMask(false).AllowPrefix(false).
	GetIPv6AddressParamsBuilder().AllowMixed(false).AllowZone(false).GetParentBuilder().ToParams()

func checkSpecialHosts(str string, addrLen int, hostQualifier *parsedHostIdentifierStringQualifier) (emb embeddedAddress) {
	suffix := IPv6UncSuffix
	//note that by using addrLen we are omitting any terminating prefix
	if addrLen > len(suffix) {
		suffixStartIndex := addrLen - len(suffix)
		//get the address for the UNC IPv6 host
		if strings.EqualFold(str[suffixStartIndex:suffixStartIndex+len(suffix)], suffix) {
			var builder strings.Builder
			beginStr := str[:suffixStartIndex]
			foundZone := false
			for i := 0; i < len(beginStr); i++ {
				c := beginStr[i]
				if c == IPv6UncSegmentSeparator {
					c = IPv6SegmentSeparator
				} else if c == IPv6UncRangeSeparatorStr[0] {
					if i+1 < len(beginStr) {
						c = beginStr[i+1]
						if c == IPv6UncRangeSeparatorStr[1] {
							c = RangeSeparator
							i++
						}
					}
				} else if c == IPv6UncZoneSeparator && !foundZone {
					foundZone = true
					c = IPv6ZoneSeparator
				}
				builder.WriteByte(c)
			}
			emb = embeddedAddress{
				isUNCIPv6Literal: true,
			}
			pa := parsedIPAddress{
				options:            defaultUncOpts,
				ipAddressParseData: ipAddressParseData{addressParseData: addressParseData{str: str}},
			}
			var err addrerr.AddressStringError
			addrStr := builder.String()
			addrStrLen := len(addrStr)
			if err = validateIPAddress(defaultUncOpts, addrStr, 0, addrStrLen, pa.getIPAddressParseData(), false); err == nil {
				if err = parseAddressQualifier(addrStr, defaultUncOpts, nil, pa.getIPAddressParseData(), addrStrLen); err == nil {
					//if err = parseAddressQualifier(str, defaultUncOpts, nil, pa.getIPAddressParseData(), addrStrLen); err == nil {
					if *pa.getQualifier() == *noQualifier {
						*pa.getQualifier() = *hostQualifier
					} else if *hostQualifier != *noQualifier {
						_ = pa.getQualifier().merge(hostQualifier)
					}
					emb.addressProvider, err = chooseIPAddressProvider(nil, str, defaultUncOpts, &pa)
				}
			}
			emb.addressStringError = err
			return
		}
	}

	//Note: could support bitstring labels and support subnets in them, however they appear to be generally unused in the real world
	//RFC 2673
	//arpa: https://www.ibm.com/support/knowledgecenter/SSLTBW_1.13.0/com.ibm.zos.r13.halz002/f1a1b3b1220.htm
	//Also, support partial dns lookups and map then to the associated subnet with prefix length, which I think we may
	//already do for ipv4 but not for ipv6, ipv4 uses the prefix notation d.c.b.a/x but ipv6 uses fewer nibbles
	//on the ipv6 side, would just need to add the proper number of zeros and the prefix length
	suffix3 := IPv6ReverseDnsSuffixDeprecated
	if addrLen > len(suffix3) {
		suffix = IPv4ReverseDnsSuffix
		suffix2 := IPv6ReverseDnsSuffix
		var isIpv4, isMatch bool
		suffixStartIndex := addrLen - len(suffix)
		if isMatch = suffixStartIndex > 0 && strings.EqualFold(str[suffixStartIndex:suffixStartIndex+len(suffix)], suffix); !isMatch {
			suffixStartIndex = addrLen - len(suffix2)
			if isMatch = suffixStartIndex > 0 && strings.EqualFold(str[suffixStartIndex:suffixStartIndex+len(suffix2)], suffix2); !isMatch {
				suffixStartIndex = addrLen - len(suffix3)
				isMatch = suffixStartIndex > 0 && strings.EqualFold(str[suffixStartIndex:suffixStartIndex+len(suffix3)], suffix3)
			}
		} else {
			isIpv4 = true
		}
		if isMatch {
			emb = embeddedAddress{
				isReverseDNS: true,
			}
			var err addrerr.AddressStringError
			var sequence string
			var params addrstrparam.IPAddressStringParams
			if isIpv4 {
				sequence, err = convertReverseDNSIPv4(str, suffixStartIndex)
				params = reverseDNSIPv4Opts
			} else {
				sequence, err = convertReverseDNSIPv6(str, suffixStartIndex)
				params = reverseDNSIPv6Opts
			}
			if err == nil {
				pa := parsedIPAddress{
					options:            params,
					ipAddressParseData: ipAddressParseData{addressParseData: addressParseData{str: sequence}},
				}
				if err = validateIPAddress(params, sequence, 0, len(sequence), pa.getIPAddressParseData(), false); err == nil {
					if err = parseAddressQualifier(str, params, nil, pa.getIPAddressParseData(), len(str)); err == nil {
						pa.qualifier = *hostQualifier
						emb.addressProvider, err = chooseIPAddressProvider(nil, sequence, params, &pa)
					}
				}
			}
			emb.addressStringError = err
		}
	}
	//			//handle TLD host https://tools.ietf.org/html/draft-osamu-v6ops-ipv4-literal-in-url-02
	//			//https://www.ietf.org/proceedings/87/slides/slides-87-v6ops-6.pdf
	//			suffix = ".v4";
	//			if(addrLen > suffix.length() &&
	//					str.regionMatches(true, suffixStartIndex = addrLen - suffix.length(), suffix, 0, suffix.length())) {
	//				//not an rfc, so let's leave it
	//			}
	return
}

//123.2.3.4 is 4.3.2.123.in-addr.arpa.

func convertReverseDNSIPv4(str string, suffixStartIndex int) (string, addrerr.AddressStringError) {
	var builder strings.Builder
	builder.Grow(suffixStartIndex)
	segCount := 0
	j := suffixStartIndex
	for i := suffixStartIndex - 1; i > 0; i-- {
		c1 := str[i]
		if c1 == IPv4SegmentSeparator {
			if j-i <= 1 {
				return "", &addressStringIndexError{
					addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
					i}
			}
			for k := i + 1; k < j; k++ {
				builder.WriteByte(str[k])
			}
			builder.WriteByte(c1)
			j = i
			segCount++
		}
	}
	for k := 0; k < j; k++ {
		builder.WriteByte(str[k])
	}
	if segCount+1 != IPv4SegmentCount {
		return "", &addressStringIndexError{
			addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
			0}
	}
	return builder.String(), nil
}

//4321:0:1:2:3:4:567:89ab would be b.a.9.8.7.6.5.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.0.0.0.0.1.2.3.4.IP6.ARPA

func convertReverseDNSIPv6(str string, suffixStartIndex int) (string, addrerr.AddressStringError) {
	var builder strings.Builder
	builder.Grow(suffixStartIndex)
	segCount := 0
	for i := suffixStartIndex - 1; i >= 0; {
		var low, high strings.Builder
		isRange := false
		for j := 0; j < 4; j++ {
			c1 := str[i]
			i--
			if i >= 0 {
				c2 := str[i]
				i--
				if c2 == IPv4SegmentSeparator {
					if c1 == SegmentWildcard {
						isRange = true
						low.WriteByte('0')
						high.WriteByte('f')
					} else {
						if isRange {
							return "", &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
								i + 1}
						}
						low.WriteByte(c1)
						high.WriteByte(c1)
					}
				} else if c2 == RangeSeparator {
					high.WriteByte(c1)
					if i >= 1 {
						c2 = str[i]
						i--
						low.WriteByte(c2)
						isFullRange := (c2 == '0' && c1 == 'f')
						if isRange && !isFullRange {
							return "", &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
								i + 1}
						}
						c2 = str[i]
						i--
						if c2 != IPv4SegmentSeparator {
							return "", &addressStringIndexError{
								addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
								i + 1}
						}
					} else {
						return "", &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
							i}
					}
					isRange = true
				} else {
					return "", &addressStringIndexError{
						addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
						i + 1}
				}
			} else if j < 3 {
				return "", &addressStringIndexError{
					addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
					i + 1}
			} else {
				if c1 == SegmentWildcard {
					isRange = true
					low.WriteByte('0')
					high.WriteByte('f')
				} else {
					if isRange {
						return "", &addressStringIndexError{
							addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
							0}
					}
					low.WriteByte(c1)
					high.WriteByte(c1)
				}
			}
		}
		segCount++
		if builder.Len() > 0 {
			builder.WriteByte(IPv6SegmentSeparator)
		}
		builder.WriteString(low.String())
		if isRange {
			builder.WriteByte(RangeSeparator)
			builder.WriteString(high.String())
		}
	}
	if segCount != IPv6SegmentCount {
		return "", &addressStringIndexError{
			addressStringError{addressError{str: str, key: "ipaddress.error.invalid.character.at.index"}},
			0}
	}
	return builder.String(), nil
}

//
// we need to initialize parsing package variables first before using them, so we put these at the bottom of this file

var zeroIPAddressString = NewIPAddressString("")

var ipv4MappedPrefix = NewIPAddressString("::ffff:0:0/96")
