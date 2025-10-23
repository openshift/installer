//
// Copyright 2020-2022 Sean C Foley
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
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
	"sync"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
)

type parsedMACAddress struct {
	macAddressParseData

	originator   *MACAddressString
	address      *MACAddress
	params       addrstrparam.MACAddressStringParams
	creationLock *sync.Mutex
}

func (provider *parsedMACAddress) getParameters() addrstrparam.MACAddressStringParams {
	return provider.params
}

func (parseData *parsedMACAddress) getMACAddressParseData() *macAddressParseData {
	return &parseData.macAddressParseData
}

func (parseData *parsedMACAddress) getAddress() (*MACAddress, addrerr.IncompatibleAddressError) {
	var err addrerr.IncompatibleAddressError
	addr := (*MACAddress)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&parseData.address))))
	if addr == nil {
		parseData.creationLock.Lock()
		addr = parseData.address
		if addr == nil {
			addr, err = parseData.createAddress()
			if err == nil {
				parseData.segmentData = nil // no longer needed
				dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&parseData.address))
				atomicStorePointer(dataLoc, unsafe.Pointer(addr))
			}
		}
		parseData.creationLock.Unlock()
	}
	return addr, err
}

func (parseData *parsedMACAddress) createAddress() (*MACAddress, addrerr.IncompatibleAddressError) {
	creator := macType.getNetwork().getAddressCreator()
	sect, err := parseData.createSection()
	if err != nil {
		return nil, err
	}
	return creator.createAddressInternal(sect.ToSectionBase(), parseData.originator).ToMAC(), nil
}

func (parseData *parsedMACAddress) createSection() (*MACAddressSection, addrerr.IncompatibleAddressError) {
	addressString := parseData.str
	addressParseData := parseData.getAddressParseData()
	actualInitialSegmentCount := addressParseData.getSegmentCount()
	creator := macType.getNetwork().getAddressCreator()
	isMultiple := false
	var segIsMult bool
	format := parseData.getFormat()
	var finalSegmentCount, initialSegmentCount int
	if format == nil {
		if parseData.isExtended() {
			initialSegmentCount = ExtendedUniqueIdentifier64SegmentCount
		} else {
			initialSegmentCount = MediaAccessControlSegmentCount
		}
		finalSegmentCount = initialSegmentCount
	} else if format == dotted {
		if parseData.isExtended() {
			initialSegmentCount = MediaAccessControlDotted64SegmentCount
		} else {
			initialSegmentCount = MediaAccessControlDottedSegmentCount
		}
		if actualInitialSegmentCount <= MediaAccessControlDottedSegmentCount && !parseData.isExtended() {
			finalSegmentCount = MediaAccessControlSegmentCount
		} else {
			finalSegmentCount = ExtendedUniqueIdentifier64SegmentCount
		}
	} else {
		if addressParseData.isSingleSegment() || parseData.isDoubleSegment() {
			if parseData.isExtended() {
				finalSegmentCount = ExtendedUniqueIdentifier64SegmentCount
			} else {
				finalSegmentCount = MediaAccessControlSegmentCount
			}
		} else if actualInitialSegmentCount <= MediaAccessControlSegmentCount && !parseData.isExtended() {
			finalSegmentCount = MediaAccessControlSegmentCount
		} else {
			finalSegmentCount = ExtendedUniqueIdentifier64SegmentCount
		}
		initialSegmentCount = finalSegmentCount
	}
	missingCount := initialSegmentCount - actualInitialSegmentCount
	expandedSegments := missingCount <= 0
	segments := make([]*AddressDivision, finalSegmentCount)
	for i, normalizedSegmentIndex := 0, 0; i < actualInitialSegmentCount; i++ {
		lower := addressParseData.getValue(i, keyLower)
		upper := addressParseData.getValue(i, keyUpper)
		if format == dotted { //aaa.bbb.ccc.ddd
			//aabb is becoming aa.bb
			segLower := SegInt(lower)
			segUpper := SegInt(upper)
			lowerHalfLower := segLower >> 8
			lowerHalfUpper := segUpper >> 8
			adjustedLower2 := segLower & 0xff
			adjustedUpper2 := segUpper & 0xff
			if lowerHalfLower != lowerHalfUpper && adjustedUpper2-adjustedLower2 != 0xff {
				return nil, &incompatibleAddressError{addressError{str: addressString, key: "ipaddress.error.invalid.joined.ranges"}}
			}
			segments[normalizedSegmentIndex], segIsMult = createSegment(
				addressString,
				lowerHalfLower,
				lowerHalfUpper,
				false,
				addressParseData,
				i,
				creator)
			normalizedSegmentIndex++
			isMultiple = isMultiple || segIsMult
			segments[normalizedSegmentIndex], segIsMult = createSegment(
				addressString,
				adjustedLower2,
				adjustedUpper2,
				false,
				addressParseData,
				i,
				creator)
			isMultiple = isMultiple || segIsMult
		} else {
			if addressParseData.isSingleSegment() || parseData.isDoubleSegment() {
				useStringIndicators := true
				var count int
				if i == actualInitialSegmentCount-1 {
					count = missingCount
				} else {
					count = MACOrganizationalUniqueIdentifierSegmentCount - 1
				}
				missingCount -= count
				isRange := lower != upper
				previousAdjustedWasRange := false
				for count >= 0 { //add the missing segments
					var newLower, newUpper uint64
					if isRange {
						segmentMask := uint64(MACMaxValuePerSegment)
						shift := uint64(count) << macBitsToSegmentBitshift
						newLower = (lower >> shift) & segmentMask
						newUpper = (upper >> shift) & segmentMask
						if previousAdjustedWasRange && newUpper-newLower != MACMaxValuePerSegment {
							//any range extending into upper segments must have full range in lower segments
							//otherwise there is no way for us to represent the address
							//so we need to check whether the lower parts cover the full range
							//eg cannot represent 0.0.0x100-0x10f or 0.0.1-1ff, but can do 0.0.0x100-0x1ff or 0.0.0-1ff
							return nil, &incompatibleAddressError{addressError{str: addressString, key: "ipaddress.error.invalid.joined.ranges"}}
						}
						previousAdjustedWasRange = newLower != newUpper

						//we may be able to reuse our strings on the final segment
						//for previous segments, strings can be reused only when the value is 0, which we do not need to cacheBitCountx.  Any other value changes when shifted.
						if count == 0 && newLower == lower {
							if newUpper != upper {
								addressParseData.unsetFlag(i, keyStandardRangeStr)
							}
						} else {
							useStringIndicators = false
						}
					} else {
						newLower = (lower >> uint(count<<3)) & MACMaxValuePerSegment
						newUpper = newLower
						if count != 0 || newLower != lower {
							useStringIndicators = false
						}
					}
					segments[normalizedSegmentIndex], segIsMult = createSegment(
						addressString,
						SegInt(newLower),
						SegInt(newUpper),
						useStringIndicators,
						addressParseData,
						i,
						creator)
					isMultiple = isMultiple || segIsMult
					normalizedSegmentIndex++
					count--
				}
				continue
			} //end joined segments
			segments[normalizedSegmentIndex], segIsMult = createSegment(
				addressString,
				SegInt(lower),
				SegInt(upper),
				true,
				addressParseData,
				i,
				creator)
			isMultiple = isMultiple || segIsMult
		}
		if !expandedSegments {
			//check for any missing segments that we should account for here
			if addressParseData.isWildcard(i) {
				expandSegments := true
				for j := i + 1; j < actualInitialSegmentCount; j++ {
					if addressParseData.isWildcard(j) { //another wildcard further down
						expandSegments = false
						break
					}
				}
				if expandSegments {
					expandedSegments = true
					count := missingCount
					for ; count > 0; count-- { //add the missing segments
						if format == dotted {
							seg, _ := createSegment(
								addressString,
								0,
								MACMaxValuePerSegment,
								false,
								addressParseData,
								i,
								creator)
							normalizedSegmentIndex++
							segments[normalizedSegmentIndex] = seg
							normalizedSegmentIndex++
							segments[normalizedSegmentIndex] = seg
						} else {
							normalizedSegmentIndex++
							segments[normalizedSegmentIndex], _ = createSegment(
								addressString,
								0,
								MACMaxValuePerSegment,
								false,
								addressParseData,
								i,
								creator)
						}
						isMultiple = true
					}
				}
			}
		}
		normalizedSegmentIndex++
	}
	return creator.createSectionInternal(segments, isMultiple).ToMAC(), nil
}

func createSegment(
	addressString string,
	val,
	upperVal SegInt,
	useFlags bool,
	parseData *addressParseData,
	parsedSegIndex int,
	creator parsedAddressCreator) (div *AddressDivision, isMultiple bool) {
	if val != upperVal {
		return createRangeSegment(addressString, val, upperVal, useFlags, parseData, parsedSegIndex, creator), true
	}
	var result *AddressDivision
	if !useFlags {
		result = creator.createSegment(val, val, nil)
	} else {
		result = creator.createSegmentInternal(
			val,
			nil, //prefix length
			addressString,
			val,
			parseData.getFlag(parsedSegIndex, keyStandardStr),
			parseData.getIndex(parsedSegIndex, keyLowerStrStartIndex),
			parseData.getIndex(parsedSegIndex, keyLowerStrEndIndex))
	}
	return result, false
}

func createRangeSegment(
	addressString string,
	lower,
	upper SegInt,
	useFlags bool,
	parseData *addressParseData,
	parsedSegIndex int,
	creator parsedAddressCreator) *AddressDivision {
	var result *AddressDivision
	if !useFlags {
		result = creator.createSegment(lower, upper, nil)
	} else {
		result = creator.createRangeSegmentInternal(
			lower,
			upper,
			nil,
			addressString,
			lower,
			upper,
			parseData.getFlag(parsedSegIndex, keyStandardStr),
			parseData.getFlag(parsedSegIndex, keyStandardRangeStr),
			parseData.getIndex(parsedSegIndex, keyLowerStrStartIndex),
			parseData.getIndex(parsedSegIndex, keyLowerStrEndIndex),
			parseData.getIndex(parsedSegIndex, keyUpperStrEndIndex))
	}
	return result
}
