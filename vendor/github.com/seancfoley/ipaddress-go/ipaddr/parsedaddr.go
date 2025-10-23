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
	"strconv"
	"sync"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

// How address sections and addresses and ranges can be created here:
// section (with no error) -> address -> sequential range
// non-nil hostSection -> hostAddress
// nil hostSection -> section (with no error) -> address -> hostAddress
// lower/upper boundary -> sequential range
// lower boundary -> mask (ie address used as mask)

type translatedResult struct {
	sections *sectionResult

	rng *SequentialRange[*IPAddress]

	mask *IPAddress

	//series IPAddressDivisionSeries // TODO LATER division grouping creation
}

type boundaryResult struct {
	lowerSection, upperSection *IPAddressSection
}

func (res *boundaryResult) createRange() *SequentialRange[*IPAddress] {
	//we need to add zone in order to reuse the lower and upper
	lowerSection := res.lowerSection
	creator := lowerSection.getAddrType().getIPNetwork().getIPAddressCreator()
	rangeLower := creator.createAddressInternalFromSection(lowerSection, NoZone, nil)
	var rangeUpper *IPAddress
	if res.upperSection == nil {
		rangeUpper = rangeLower
	} else {
		rangeUpper = creator.createAddressInternalFromSection(res.upperSection, NoZone, nil)
	}
	result := rangeLower.SpanWithRange(rangeUpper)
	return result
}

func (res *boundaryResult) createMask() *IPAddress {
	lowerSection := res.lowerSection
	creator := lowerSection.getAddrType().getIPNetwork().getIPAddressCreator()
	return creator.createAddressInternalFromSection(res.lowerSection, NoZone, nil)
}

type sectionResult struct {
	section, hostSection *IPAddressSection

	address, hostAddress *IPAddress

	joinHostError, joinAddressError /* inet_aton, single seg */, mixedError, maskError addrerr.IncompatibleAddressError
}

func (res *sectionResult) withoutAddressException() bool {
	return res.joinAddressError == nil && res.mixedError == nil && res.maskError == nil
}

type parsedIPAddress struct {
	ipAddressParseData

	ipAddrProvider // provides a few methods like isInvalid

	options               addrstrparam.IPAddressStringParams
	originator            HostIdentifierString
	vals                  translatedResult
	skipCntains           *bool
	maskers, mixedMaskers []Masker

	creationLock sync.Mutex
}

func (parseData *parsedIPAddress) values() *translatedResult {
	return &parseData.vals
}

func (parseData *parsedIPAddress) providerCompare(other ipAddressProvider) (int, addrerr.IncompatibleAddressError) {
	return providerCompare(parseData, other)
}

func (parseData *parsedIPAddress) providerEquals(other ipAddressProvider) (bool, addrerr.IncompatibleAddressError) {
	return providerEquals(parseData, other)
}

func (parseData *parsedIPAddress) isProvidingIPAddress() bool {
	return true
}

func (parseData *parsedIPAddress) getType() ipType {
	return fromVersion(parseData.getProviderIPVersion())
}

func (parseData *parsedIPAddress) getParameters() addrstrparam.IPAddressStringParams {
	return parseData.options
}

// Note: the following are needed because we have two anonymous fields and there are name clashes
// Instead of defaulting to the default methods in ipAddressProvider, we need to defer to our parsed data for these methods
//

func (parseData *parsedIPAddress) isProvidingMixedIPv6() bool {
	return parseData.ipAddressParseData.isProvidingMixedIPv6()
}

func (parseData *parsedIPAddress) isProvidingIPv6() bool {
	return parseData.ipAddressParseData.isProvidingIPv6()
}

func (parseData *parsedIPAddress) isProvidingIPv4() bool {
	return parseData.ipAddressParseData.isProvidingIPv4()
}

func (parseData *parsedIPAddress) isProvidingBase85IPv6() bool {
	return parseData.ipAddressParseData.isProvidingBase85IPv6()
}

func (parseData *parsedIPAddress) getProviderIPVersion() IPVersion {
	return parseData.ipAddressParseData.getProviderIPVersion()
}

func (parseData *parsedIPAddress) getIPAddressParseData() *ipAddressParseData {
	return &parseData.ipAddressParseData
}

// creation methods start here

func (parseData *parsedIPAddress) createSections(doSections, doRangeBoundaries, withUpper bool) (sections sectionResult, boundaries boundaryResult) {
	version := parseData.getProviderIPVersion()
	if version.IsIPv4() {
		return parseData.createIPv4Sections(doSections, doRangeBoundaries, withUpper)
	} else if version.IsIPv6() {
		return parseData.createIPv6Sections(doSections, doRangeBoundaries, withUpper)
	}
	return
}

func (parseData *parsedIPAddress) getProviderSeqRange() *SequentialRange[*IPAddress] {
	val := parseData.values()
	result := (*SequentialRange[*IPAddress])(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&val.rng))))
	if result == nil {
		parseData.creationLock.Lock()
		result = val.rng
		if result == nil {
			sections := val.sections
			if sections == nil {
				_, boundaries := parseData.createSections(false, true, true)
				// creates lower, upper, then range from the two
				result = boundaries.createRange()
			} else {
				if sections.withoutAddressException() {
					result = sections.address.ToSequentialRange()
				} else {
					_, boundaries := parseData.createSections(false, true, true)
					result = boundaries.createRange()
				}
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&val.rng))
			atomicStorePointer(dataLoc, unsafe.Pointer(result))
		}
		parseData.creationLock.Unlock()
	}
	return result
}

// this is for parsed addresses which are masks in and of themselves
// with masks, only the lower value matters
func (parseData *parsedIPAddress) getValForMask() *IPAddress {
	val := parseData.values()
	mask := (*IPAddress)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&val.mask))))
	if mask == nil {
		parseData.creationLock.Lock()
		mask = val.mask
		if mask == nil {
			_, boundaries := parseData.createSections(false, true, false)
			mask = boundaries.createMask()
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&val.mask))
			atomicStorePointer(dataLoc, unsafe.Pointer(mask))
		}
		parseData.creationLock.Unlock()
	}
	return mask
}

func (parseData *parsedIPAddress) getCachedAddresses(forHostAddr bool) *sectionResult {
	val := parseData.values()
	sections := (*sectionResult)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&val.sections))))
	if sections == nil {
		parseData.creationLock.Lock()
		sections = val.sections
		if sections == nil {
			sects, _ := parseData.createSections(true, false, false)
			sections = &sects
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&val.sections))
			atomicStorePointer(dataLoc, unsafe.Pointer(sections))
		}
		parseData.creationLock.Unlock()
	}
	if sections.withoutAddressException() {
		var addr *IPAddress
		if forHostAddr {
			addr = (*IPAddress)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&sections.hostAddress))))
		} else {
			addr = (*IPAddress)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&sections.address))))
		}
		if addr == nil {
			parseData.creationLock.Lock()
			if forHostAddr {
				addr = sections.hostAddress
			} else {
				addr = sections.address
			}
			if addr == nil {
				var section *IPAddressSection
				var originator HostIdentifierString
				if forHostAddr {
					section = sections.hostSection
					if section == nil {
						section = sections.section
					}
				} else {
					section = sections.section
					originator = parseData.originator
				}
				creator := section.getAddrType().getIPNetwork().getIPAddressCreator()
				addr = creator.createAddressInternalFromSection(section, parseData.getQualifier().getZone(), originator)
				var dataLoc *unsafe.Pointer
				if forHostAddr {
					dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&sections.hostAddress))
				} else {
					// if range created first, stick the lower and upper into the address cache,
					// but only if the address no prefix, because the range never has prefix lengths
					if rng := val.rng; rng != nil && !addr.IsPrefixed() {
						cache := addr.cache
						if cache != nil {
							cache.addrsCache = &addrsCache{
								lower: rng.lower.ToAddressBase(),
								upper: rng.upper.ToAddressBase(),
							}
						}
					}
					dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&sections.address))
				}
				atomicStorePointer(dataLoc, unsafe.Pointer(addr))
			}
			parseData.creationLock.Unlock()
		}
	}
	return sections
}

// this is for parsed addresses which have associated masks
func (parseData *parsedIPAddress) getProviderMask() *IPAddress {
	return parseData.getQualifier().getMaskLower()
}

func (parseData *parsedIPAddress) getProviderHostAddress() (*IPAddress, addrerr.IncompatibleAddressError) {
	addrs := parseData.getCachedAddresses(true)
	if addrs.mixedError != nil {
		return nil, addrs.mixedError
	} else if addrs.joinHostError != nil {
		return nil, addrs.joinHostError
	}
	return addrs.hostAddress, nil
}

func (parseData *parsedIPAddress) getProviderAddress() (*IPAddress, addrerr.IncompatibleAddressError) {
	addrs := parseData.getCachedAddresses(false)
	if addrs.mixedError != nil {
		return nil, addrs.mixedError
	} else if addrs.maskError != nil {
		return nil, addrs.maskError
	} else if addrs.joinAddressError != nil {
		return nil, addrs.joinAddressError
	}
	return addrs.address, nil
}

func (parseData *parsedIPAddress) getVersionedAddress(version IPVersion) (*IPAddress, addrerr.IncompatibleAddressError) {
	thisVersion := parseData.getProviderIPVersion()
	if version != thisVersion {
		return nil, nil
	}
	return parseData.getProviderAddress()
}

func (parseData *parsedIPAddress) getProviderNetworkPrefixLen() PrefixLen {
	return parseData.getQualifier().getEquivalentPrefixLen()
}

// TODO LATER getDivisionGrouping
//func (parseData *parsedIPAddress)   groupingIsSequential() bool {
//		try {
//			return getDivisionGrouping().isSequential();
//		} catch(IncompatibleAddressException e) {
//			// division groupings avoid all IncompatibleAddressException caused by regrouping the values into segments of different size
//			// that takes care of two of the sources of IncompatibleAddressException: joining mixed segs, and expanding inet_aton ipv4 or single-segment ipv6 into the standard number of ipv4 or ipv6 segments
//
//			// Those remaining are the IncompatibleAddressException caused by masks, which are the result of individual divisions becoming non-sequential
//			// So in such cases, you know we are not sequential.  So we return false.
//			// the usual caveat is that this cannot happen with standard network or host masks
//			return false;
//		}
//	}
//
//func (parseData *parsedIPAddress) IsSequential() bool {
//		TranslatedResult<?,?> val = values;
//		if(val != null) {
//			// check address first
//			if(!val.withoutSections()) {
//				// address already there, use it if we can
//				if(val.withoutAddressException()) {
//					return val.getAddress().isSequential();
//				}
//				return groupingIsSequential();
//			}
//			if(!val.withoutGrouping()) {
//				return groupingIsSequential();
//			}
//		}
//		// neither address nor grouping is there, create the address
//		val = getCachedAddresses(false);
//		if(val.withoutAddressException()) {
//			return val.getAddress().isSequential();
//		}
//		return groupingIsSequential();
//	}

// skips contains checking for addresses already parsed -
// so this is not a case of unusual string formatting, because this is not for comparing strings,
// but more a case of whether the parsing data structures are easy to use or not
func (parseData *parsedIPAddress) skipContains() bool {
	segmentCount := parseData.getAddressParseData().getSegmentCount()
	// first we must excluded cases where the segments line up differently than standard, although we do not exclude ipv6 compressed
	if parseData.isProvidingIPv4() {
		if segmentCount != IPv4SegmentCount { // accounts for is_inet_aton_joined, singleSegment and wildcard segments
			return true
		}
	} else {
		if parseData.isProvidingMixedIPv6() || (segmentCount != IPv6SegmentCount && !parseData.isCompressed()) { // accounts for single segment and wildcard segments
			return true
		}
	}
	// exclude non-standard masks which will modify segment values from their parsed values
	mask := parseData.getProviderMask()
	if mask != nil && mask.GetBlockMaskPrefixLen(true) == nil { // handles non-standard masks
		return true
	}
	return false
}

//we do not call this method with parse data from inet_aton or single segment strings, so the cast to int is fine.
//this is only for addresses with standard segment counts, although we do allow compressed.
func (parseData *parsedIPAddress) isPrefixSubnet(networkPrefixLength BitCount) bool {
	var bytesPerSegment int
	var max SegInt
	var bitsPerSegment BitCount
	if parseData.isProvidingIPv4() {
		bytesPerSegment = IPv4BytesPerSegment
		bitsPerSegment = IPv4BitsPerSegment
		max = IPv4MaxValuePerSegment
	} else {
		bytesPerSegment = IPv6BytesPerSegment
		bitsPerSegment = IPv6BitsPerSegment
		max = IPv6MaxValuePerSegment
	}
	addressParseData := parseData.getAddressParseData()
	segmentCount := addressParseData.getSegmentCount()
	if parseData.isCompressed() {
		compressedCount := IPv6SegmentCount - segmentCount
		compressedIndex := addressParseData.getConsecutiveSeparatorSegmentIndex()
		return isPrefixSubnet(
			func(segmentIndex int) SegInt {
				if segmentIndex >= compressedIndex {
					if segmentIndex-compressedIndex < compressedCount {
						return 0
					}
					segmentIndex -= compressedCount
				}
				return SegInt(parseData.getValue(segmentIndex, keyLower))
			},
			func(segmentIndex int) SegInt {
				if segmentIndex >= compressedIndex {
					if segmentIndex-compressedIndex < compressedCount {
						return 0
					}
					segmentIndex -= compressedCount
				}
				return SegInt(parseData.getValue(segmentIndex, keyUpper))
			},
			segmentCount+compressedCount,
			bytesPerSegment,
			bitsPerSegment,
			max,
			networkPrefixLength,
			zerosOrFullRange)
	}
	return isPrefixSubnet(
		func(segmentIndex int) SegInt {
			return SegInt(parseData.getValue(segmentIndex, keyLower))
		},
		func(segmentIndex int) SegInt {
			return SegInt(parseData.getValue(segmentIndex, keyUpper))
		},
		segmentCount,
		bytesPerSegment,
		bitsPerSegment,
		max,
		networkPrefixLength,
		zerosOrFullRange)
}

func (parseData *parsedIPAddress) containmentCheck(other ipAddressProvider, networkOnly, equals, checkZone bool) (res boolSetting) {
	if otherParsed, ok := other.(*parsedIPAddress); ok {
		sect := (*sectionResult)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&parseData.vals.sections))))
		otherSect := (*sectionResult)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&otherParsed.vals.sections))))
		if sect == nil || otherSect == nil {
			// one or the other value not yet created, so take the shortcut that provides an answer most (but not all) of the time
			// An answer is provided for all normalized, conventional or canonical addresses
			res = parseData.containsProv(otherParsed, networkOnly, equals)
			if checkZone && res.isSet && res.val {
				res.val = parseData.getQualifier().getZone() == otherParsed.getQualifier().getZone()
			}
		} // else we defer to the values-based containment check (in the caller), which is best since it is ready to go
	}
	return
}

func (parseData *parsedIPAddress) containsProvider(other ipAddressProvider) (res boolSetting) {
	return parseData.containmentCheck(other, false, false, true)
}

func (parseData *parsedIPAddress) parsedEquals(other ipAddressProvider) (res boolSetting) {
	return parseData.containmentCheck(other, false, true, true)
}

func (parseData *parsedIPAddress) prefixContainsProvider(other ipAddressProvider) boolSetting {
	return parseData.containmentCheck(other, true, false, false)
}

func (parseData *parsedIPAddress) prefixEqualsProvider(other ipAddressProvider) boolSetting {
	return parseData.containmentCheck(other, true, true, false)
}

//not used for invalid, or cases where parseData.isEmpty or parseData.isAll
func (parseData *parsedIPAddress) containsProv(other *parsedIPAddress, networkOnly, equals bool) (res boolSetting) {
	pd := parseData.getAddressParseData()
	otherParseData := other.getAddressParseData()
	segmentData := pd.getSegmentData()                  //grab this field for thread safety, other threads can make it disappear
	otherSegmentData := otherParseData.getSegmentData() //grab this field for thread safety, other threads can make it disappear
	if segmentData == nil || otherSegmentData == nil {
		return
	} else if parseData.skipContains() || other.skipContains() { // this excludes mixed addresses, amongst others
		return
	}
	ipVersion := parseData.getProviderIPVersion()
	if ipVersion != other.getProviderIPVersion() {
		return boolSetting{true, false}
	}
	segmentCount := pd.getSegmentCount()
	otherSegmentCount := otherParseData.getSegmentCount()
	var max SegInt
	var compressedAlready, otherCompressedAlready bool
	var expectedSegCount, bytesPerSegment int
	var bitsPerSegment BitCount
	if parseData.isProvidingIPv4() {
		max = IPv4MaxValuePerSegment
		expectedSegCount = IPv4SegmentCount
		bitsPerSegment = IPv4BitsPerSegment
		bytesPerSegment = IPv4BytesPerSegment
		compressedAlready = true
		otherCompressedAlready = true
	} else {
		max = IPv6MaxValuePerSegment
		expectedSegCount = IPv6SegmentCount
		bitsPerSegment = IPv6BitsPerSegment
		bytesPerSegment = IPv6BytesPerSegment
		compressedAlready = expectedSegCount == segmentCount
		otherCompressedAlready = expectedSegCount == otherSegmentCount
	}
	pref := parseData.getProviderNetworkPrefixLen()
	otherPref := other.getProviderNetworkPrefixLen()
	var networkSegIndex, hostSegIndex, endIndex, otherHostAllSegIndex, hostAllSegIndex int
	endIndex = segmentCount

	// determine what indexes to use for network, host, and prefix block adjustments (hostAllSegIndex and otherHostAllSegIndex)
	var adjustedOtherPref PrefixLen
	if pref == nil {
		networkOnly = false
		hostAllSegIndex = expectedSegCount
		otherHostAllSegIndex = expectedSegCount
		hostSegIndex = expectedSegCount
		networkSegIndex = hostSegIndex - 1
	} else {
		prefLen := pref.bitCount()
		if networkOnly {
			hostSegIndex = getHostSegmentIndex(prefLen, bytesPerSegment, bitsPerSegment)
			hostAllSegIndex = hostSegIndex
			otherHostAllSegIndex = hostSegIndex
			networkSegIndex = getNetworkSegmentIndex(prefLen, bytesPerSegment, bitsPerSegment)
			// we treat the other as if it were a prefix block of the same prefix length
			// this allows us to compare entire segments for prefixEquals, ignoring the host values
			adjustedOtherPref = pref
		} else {
			otherHostAllSegIndex = expectedSegCount
			hostSegIndex = getHostSegmentIndex(prefLen, bytesPerSegment, bitsPerSegment)
			networkSegIndex = getNetworkSegmentIndex(prefLen, bytesPerSegment, bitsPerSegment)
			if parseData.isPrefixSubnet(prefLen) {
				hostAllSegIndex = hostSegIndex
				if !equals {
					// no need to look at host for containment when a prefix subnet
					networkOnly = true
				}
			} else {
				hostAllSegIndex = expectedSegCount
			}
		}
	}
	// Now determine if the other is a prefix block subnet, and if so, adjust otherHostAllSegIndex
	if otherPref != nil {
		otherPrefLen := otherPref.bitCount()
		if adjustedOtherPref == nil || otherPrefLen < adjustedOtherPref.bitCount() {
			otherHostIndex := getHostSegmentIndex(otherPrefLen, bytesPerSegment, bitsPerSegment)
			if otherHostIndex < otherHostAllSegIndex &&
				other.isPrefixSubnet(otherPrefLen) {
				otherHostAllSegIndex = otherHostIndex
			}
		} else {
			otherPref = adjustedOtherPref
		}
	} else {
		otherPref = adjustedOtherPref
	}
	i, j, normalizedCount := 0, 0, 0
	var compressedCount, otherCompressedCount int
	for i < endIndex || compressedCount > 0 {
		if networkOnly && normalizedCount > networkSegIndex {
			break
		}
		var lower, upper SegInt
		if compressedCount <= 0 {
			lower = SegInt(parseData.getValue(i, keyLower))
			upper = SegInt(parseData.getValue(i, keyUpper))
		}
		if normalizedCount >= hostAllSegIndex { // we've reached the prefixed segment
			segPrefLength := getSegmentPrefixLength(bitsPerSegment, pref, normalizedCount)
			segPref := segPrefLength.bitCount()
			networkMask := ^SegInt(0) << uint(bitsPerSegment-segPref)
			hostMask := ^networkMask
			lower &= networkMask
			upper |= hostMask
		}
		var otherLower, otherUpper SegInt
		if normalizedCount > otherHostAllSegIndex {
			otherLower = 0
			otherUpper = max
		} else {
			if otherCompressedCount <= 0 {
				otherLower = SegInt(otherParseData.getValue(j, keyLower))
				otherUpper = SegInt(otherParseData.getValue(j, keyUpper))
			}
			if normalizedCount == otherHostAllSegIndex { // we've reached the prefixed segment
				segPrefLength := getSegmentPrefixLength(bitsPerSegment, otherPref, normalizedCount)
				segPref := segPrefLength.bitCount()
				networkMask := ^SegInt(0) << uint(bitsPerSegment-segPref)
				hostMask := ^networkMask
				otherLower &= networkMask
				otherUpper |= hostMask
			}
		}
		if equals {
			if lower != otherLower || upper != otherUpper {
				return boolSetting{true, false}
			}
		} else {
			if lower > otherLower || upper < otherUpper {
				return boolSetting{true, false}
			}
		}
		if !compressedAlready {
			if compressedCount > 0 {
				compressedCount--
				if compressedCount == 0 {
					compressedAlready = true
				}
			} else if parseData.segmentIsCompressed(i) {
				i++
				compressedCount = expectedSegCount - segmentCount
			} else {
				i++
			}
		} else {
			i++
		}
		if !otherCompressedAlready {
			if otherCompressedCount > 0 {
				otherCompressedCount--
				if otherCompressedCount == 0 {
					otherCompressedAlready = true
				}
			} else if other.segmentIsCompressed(j) {
				j++
				otherCompressedCount = expectedSegCount - otherSegmentCount
			} else {
				j++
			}
		} else {
			j++
		}
		normalizedCount++
	}
	return boolSetting{true, true}
}

func allocateSegments(
	segments,
	originalSegments []*AddressDivision,
	segmentCount,
	originalCount int) []*AddressDivision {
	if segments == nil {
		segments = createSegmentArray(segmentCount)
		if originalCount > 0 {
			copy(segments, originalSegments[:originalCount])
		}
	}
	return segments
}

func (parseData *parsedIPAddress) createIPv4Sections(doSections, doRangeBoundaries, withUpper bool) (sections sectionResult, boundaries boundaryResult) {
	qualifier := parseData.getQualifier()
	prefLen := getPrefixLength(qualifier)
	isMultiple := false
	isHostMultiple := false
	var segIsMult bool
	mask := parseData.getProviderMask()
	if mask != nil && mask.GetBlockMaskPrefixLen(true) != nil {
		mask = nil //we don't do any masking if the mask is a subnet mask, instead we just map it to the corresponding prefix length
	}
	hasMask := mask != nil
	addrParseData := parseData.getAddressParseData()
	segmentCount := addrParseData.getSegmentCount()
	if hasMask && parseData.maskers == nil {
		parseData.maskers = make([]Masker, segmentCount)
	}
	creator := ipv4Type.getIPNetwork().getIPAddressCreator()
	missingCount := IPv4SegmentCount - segmentCount

	var hostSegments, segments, lowerSegments, upperSegments []*AddressDivision
	if doSections {
		segments = createSegmentArray(IPv4SegmentCount)
	} else if doRangeBoundaries {
		lowerSegments = createSegmentArray(IPv4SegmentCount)
	} else {
		return
	}
	expandedSegments := missingCount <= 0
	expandedStart, expandedEnd := -1, -1
	addressString := parseData.str
	maskedIsDifferent := false
	for i, normalizedSegmentIndex := 0, 0; i < segmentCount; i++ {
		lower := addrParseData.getValue(i, keyLower)
		upper := addrParseData.getValue(i, keyUpper)
		if !expandedSegments {
			//check for any missing segments that we should account for here
			isLastSegment := i == segmentCount-1
			isWildcard := addrParseData.isWildcard(i)
			expandedSegments = isLastSegment
			if !expandedSegments {
				// if we are inet_aton, we must wait for last segment
				// otherwise, we check if we are wildcard and no other wildcard further down
				expandedSegments = !parseData.is_inet_aton_joined() && isWildcard
				if expandedSegments {
					for j := i + 1; j < segmentCount; j++ {
						if addrParseData.isWildcard(j) { //another wildcard further down
							expandedSegments = false
							break
						}
					}
				}
			}
			if expandedSegments {
				if isWildcard {
					upper = 0xffffffff >> uint((3-missingCount)<<3)
				} else {
					expandedStart = i
					expandedEnd = i + missingCount
				}
				bits := BitCount(missingCount+1) << ipv4BitsToSegmentBitshift // BitCount(missingCount+1) * IPv4BitsPerSegment
				var maskedLower, maskedUpper uint64
				if hasMask {
					var divMask uint64
					for k := 0; k <= missingCount; k++ {
						divMask = (divMask << uint(IPv4BitsPerSegment)) | uint64(mask.GetSegment(normalizedSegmentIndex+k).GetSegmentValue())
					}
					masker := parseData.maskers[i]
					if masker == nil {
						maxValue := ^(^uint64(0) << uint(bits))
						masker = MaskRange(lower, upper, divMask, maxValue)
						parseData.maskers[i] = masker
					}
					if !masker.IsSequential() && sections.maskError == nil {
						sections.maskError = &incompatibleAddressError{
							addressError: addressError{
								str: maskString(lower, upper, divMask),
								key: "ipaddress.error.maskMismatch",
							},
						}
					}
					maskedLower = masker.GetMaskedLower(lower, divMask)
					maskedUpper = masker.GetMaskedUpper(upper, divMask)
					maskedIsDifferent = maskedIsDifferent || maskedLower != lower || maskedUpper != upper
				} else {
					maskedLower = lower
					maskedUpper = upper
				}
				shift := bits
				count := missingCount
				for count >= 0 { //add the missing segments
					shift -= IPv4BitsPerSegment
					currentPrefix := getSegmentPrefixLength(IPv4BitsPerSegment, prefLen, normalizedSegmentIndex)
					//currentPrefix := getQualifierSegmentPrefixLength(normalizedSegmentIndex, , qualifier)
					hostSegLower := SegInt((lower >> uint(shift)) & IPv4MaxValuePerSegment)
					var hostSegUpper SegInt
					if lower == upper {
						hostSegUpper = hostSegLower
					} else {
						hostSegUpper = SegInt((upper >> uint(shift)) & IPv4MaxValuePerSegment)
					}
					var maskedSegLower, maskedSegUpper SegInt
					if hasMask {
						maskedSegLower = SegInt((maskedLower >> uint(shift)) & IPv4MaxValuePerSegment)
						if maskedLower == maskedUpper {
							maskedSegUpper = maskedSegLower
						} else {
							maskedSegUpper = SegInt((maskedUpper >> uint(shift)) & IPv4MaxValuePerSegment)
						}
					} else {
						maskedSegLower = hostSegLower
						maskedSegUpper = hostSegUpper
					}
					if doSections {
						if maskedIsDifferent || currentPrefix != nil {
							hostSegments = allocateSegments(hostSegments, segments, IPv4SegmentCount, normalizedSegmentIndex)
							hostSegments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
								addressString,
								IPv4,
								hostSegLower,
								hostSegUpper,
								false,
								i,
								nil,
								creator)
							isHostMultiple = isHostMultiple || segIsMult
						}
						segments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
							addressString,
							IPv4,
							maskedSegLower,
							maskedSegUpper,
							false,
							i,
							currentPrefix,
							creator)
						isMultiple = isMultiple || segIsMult
					}
					if doRangeBoundaries {
						isRange := maskedSegLower != maskedSegUpper
						if !doSections || isRange {
							if doSections {
								lowerSegments = allocateSegments(lowerSegments, segments, IPv4SegmentCount, normalizedSegmentIndex)
							} // else segments already allocated
							lowerSegments[normalizedSegmentIndex], _ = parseData.createSegment(
								addressString,
								IPv4,
								maskedSegLower,
								maskedSegLower,
								false,
								i,
								currentPrefix,
								creator)
						} else if lowerSegments != nil {
							lowerSegments[normalizedSegmentIndex] = segments[normalizedSegmentIndex]
						}
						if withUpper {
							if isRange {
								upperSegments = allocateSegments(upperSegments, lowerSegments, IPv4SegmentCount, normalizedSegmentIndex)
								upperSegments[normalizedSegmentIndex], _ = parseData.createSegment(
									addressString,
									IPv4,
									maskedSegUpper,
									maskedSegUpper,
									false,
									i,
									currentPrefix,
									creator)
							} else if upperSegments != nil {
								upperSegments[normalizedSegmentIndex] = lowerSegments[normalizedSegmentIndex]
							}
						}
					}
					normalizedSegmentIndex++
					count--
				}
				addrParseData.setBitLength(i, bits)
				continue
			} //end handle inet_aton joined segments
		}
		hostLower, hostUpper := lower, upper
		var masker Masker
		unmasked := true
		if hasMask {
			masker = parseData.maskers[i]
			maskInt := uint64(mask.GetSegment(normalizedSegmentIndex).GetSegmentValue())
			if masker == nil {
				masker = MaskRange(lower, upper, maskInt, uint64(creator.getMaxValuePerSegment()))
				parseData.maskers[i] = masker
			}
			if !masker.IsSequential() && sections.maskError == nil {
				sections.maskError = &incompatibleAddressError{
					addressError: addressError{
						str: maskString(lower, upper, maskInt),
						key: "ipaddress.error.maskMismatch",
					},
				}
			}
			lower = masker.GetMaskedLower(lower, maskInt)
			upper = masker.GetMaskedUpper(upper, maskInt)
			unmasked = hostLower == lower && hostUpper == upper
			maskedIsDifferent = maskedIsDifferent || !unmasked
		}
		segmentPrefixLength := getSegmentPrefixLength(IPv4BitsPerSegment, prefLen, normalizedSegmentIndex)
		if doSections {
			if maskedIsDifferent || segmentPrefixLength != nil {
				hostSegments = allocateSegments(hostSegments, segments, IPv4SegmentCount, normalizedSegmentIndex)
				hostSegments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
					addressString,
					IPv4,
					SegInt(hostLower),
					SegInt(hostUpper),
					true,
					i,
					nil,
					creator)
				isHostMultiple = isHostMultiple || segIsMult
			}
			segments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
				addressString,
				IPv4,
				SegInt(lower),
				SegInt(upper),
				unmasked,
				i,
				segmentPrefixLength,
				creator)
			isMultiple = isMultiple || segIsMult
		}
		if doRangeBoundaries {
			isRange := lower != upper
			if !doSections || isRange {
				if doSections {
					lowerSegments = allocateSegments(lowerSegments, segments, IPv4SegmentCount, normalizedSegmentIndex)
				} // else segments already allocated
				lowerSegments[normalizedSegmentIndex], _ = parseData.createSegment(
					addressString,
					IPv4,
					SegInt(lower),
					SegInt(lower),
					false,
					i,
					segmentPrefixLength,
					creator)
			} else if lowerSegments != nil {
				lowerSegments[normalizedSegmentIndex] = segments[normalizedSegmentIndex]
			}
			if withUpper {
				if isRange {
					upperSegments = allocateSegments(upperSegments, lowerSegments, IPv4SegmentCount, normalizedSegmentIndex)
					upperSegments[normalizedSegmentIndex], _ = parseData.createSegment(
						addressString,
						IPv4,
						SegInt(upper),
						SegInt(upper),
						false,
						i,
						segmentPrefixLength,
						creator)
				} else if upperSegments != nil {
					upperSegments[normalizedSegmentIndex] = lowerSegments[normalizedSegmentIndex]
				}
			}
		}
		normalizedSegmentIndex++
		addrParseData.setBitLength(i, IPv4BitsPerSegment)
	}
	prefLength := getPrefixLength(qualifier)
	var result, hostResult *IPAddressSection
	if doSections {
		result = creator.createPrefixedSectionInternal(segments, isMultiple, prefLength)
		sections.section = result
		if hostSegments != nil {
			hostResult = creator.createSectionInternal(hostSegments, isHostMultiple).ToIP()
			sections.hostSection = hostResult
			if checkExpandedValues(hostResult, expandedStart, expandedEnd) {
				sections.joinHostError = &incompatibleAddressError{
					addressError{
						str: addressString,
						key: "ipaddress.error.invalid.joined.ranges",
					},
				}
			}
		}
		if checkExpandedValues(result, expandedStart, expandedEnd) {
			sections.joinAddressError = &incompatibleAddressError{addressError{str: addressString, key: "ipaddress.error.invalid.joined.ranges"}}
			if hostResult == nil {
				sections.joinHostError = sections.joinAddressError
			}
		}
	}
	if doRangeBoundaries {
		// if we have a prefix subnet, it is possible our lower and upper boundaries exceed what appears in the parsed address
		prefixLength := getPrefixLength(qualifier)
		isPrefixSub := false
		if prefixLength != nil {
			var lowerSegs, upperSegs []*AddressDivision
			if doSections {
				upperSegs = segments
				lowerSegs = upperSegs
			} else {
				lowerSegs = lowerSegments
				if upperSegments == nil {
					upperSegs = lowerSegments
				} else {
					upperSegs = upperSegments
				}
			}
			isPrefixSub = isPrefixSubnet(
				func(index int) SegInt { return lowerSegs[index].ToSegmentBase().GetSegmentValue() },
				func(index int) SegInt { return upperSegs[index].ToSegmentBase().GetUpperSegmentValue() },
				len(lowerSegs),
				IPv4BytesPerSegment,
				IPv4BitsPerSegment,
				IPv4MaxValuePerSegment,
				prefixLength.bitCount(),
				zerosOnly)
			if isPrefixSub {
				if lowerSegments == nil {
					//allocate lower segments from address segments
					lowerSegments = allocateSegments(lowerSegments, segments, IPv4SegmentCount, IPv4SegmentCount)
				}
				if upperSegments == nil {
					//allocate upper segments from lower segments
					upperSegments = allocateSegments(upperSegments, lowerSegments, IPv4SegmentCount, IPv4SegmentCount)
				}
			}
		}
		if lowerSegments != nil {
			boundaries.lowerSection = creator.createPrefixedSectionInternalSingle(lowerSegments, false, prefLength)
		}
		if upperSegments != nil {
			section := creator.createPrefixedSectionInternal(upperSegments, false, prefLength)
			if isPrefixSub {
				section = section.ToPrefixBlock()
			}
			boundaries.upperSection = section.GetUpper()
		}
	}
	return
}

func (parseData *parsedIPAddress) createIPv6Sections(doSections, doRangeBoundaries, withUpper bool) (sections sectionResult, boundaries boundaryResult) {
	qualifier := parseData.getQualifier()
	prefLen := getPrefixLength(qualifier)
	mask := parseData.getProviderMask()
	if mask != nil && mask.GetBlockMaskPrefixLen(true) != nil {
		mask = nil //we don't do any masking if the mask is a subnet mask, instead we just map it to the corresponding prefix length
	}
	hasMask := mask != nil
	isMultiple := false
	isHostMultiple := false
	var segIsMult bool
	addressParseData := parseData.getAddressParseData()
	segmentCount := addressParseData.getSegmentCount()
	if hasMask && parseData.maskers == nil {
		parseData.maskers = make([]Masker, segmentCount)
	}
	creator := ipv6Type.getIPNetwork().getIPAddressCreator()
	ipv6SegmentCount := IPv6SegmentCount
	var hostSegments, segments, lowerSegments, upperSegments []*AddressDivision
	if doSections {
		segments = createSegmentArray(IPv6SegmentCount)
	} else if doRangeBoundaries {
		lowerSegments = createSegmentArray(IPv6SegmentCount)
	} else {
		return
	}
	mixed := parseData.isProvidingMixedIPv6()

	normalizedSegmentIndex := 0
	var missingSegmentCount int
	if mixed {
		missingSegmentCount = IPv6MixedOriginalSegmentCount
	} else {
		missingSegmentCount = IPv6SegmentCount
	}
	missingSegmentCount -= segmentCount

	expandedSegments := missingSegmentCount <= 0
	expandedStart, expandedEnd := -1, -1
	addressString := parseData.str
	maskedIsDifferent := false

	//get the segments for IPv6
	for i := 0; i < segmentCount; i++ {
		lower := addressParseData.getValue(i, keyLower)
		upper := addressParseData.getValue(i, keyUpper)

		if !expandedSegments {
			isLastSegment := i == segmentCount-1
			isWildcard := addressParseData.isWildcard(i)
			isCompressed := parseData.segmentIsCompressed(i)

			// figure out if this segment should be expanded
			expandedSegments = isLastSegment || isCompressed
			if !expandedSegments {
				// we check if we are wildcard and no other wildcard or compressed segment further down
				expandedSegments = isWildcard
				if expandedSegments {
					for j := i + 1; j < segmentCount; j++ {
						if addressParseData.isWildcard(j) || parseData.segmentIsCompressed(j) {
							expandedSegments = false
							break
						}
					}
				}
			}
			if expandedSegments {
				var lowerHighBytes, upperHighBytes uint64
				hostIsRange := false
				if !isCompressed {
					if isWildcard {
						if missingSegmentCount > 3 {
							upperHighBytes = 0xffffffffffffffff >> uint((7-missingSegmentCount)<<4)
							upper = 0xffffffffffffffff
						} else {
							upperHighBytes = 0
							upper = 0xffffffffffffffff >> uint((3-missingSegmentCount)<<4)
						}
						lower = 0
						hostIsRange = true
					} else {
						if missingSegmentCount > 3 {
							lowerHighBytes = addressParseData.getValue(i, keyExtendedLower) //the high half of the lower value
							upperHighBytes = addressParseData.getValue(i, keyExtendedUpper) //the high half of the upper value
							hostIsRange = (lower != upper) || (lowerHighBytes != upperHighBytes)
						} else {
							hostIsRange = lower != upper
						}
						expandedStart = i
						expandedEnd = i + missingSegmentCount
					}
				}
				bits := BitCount(missingSegmentCount+1) << ipv6BitsToSegmentBitshift // BitCount(missingSegmentCount+1) * IPv6BitsPerSegment
				var maskedLower, maskedUpper, maskedLowerHighBytes, maskedUpperHighBytes uint64
				maskedIsRange := false
				if hasMask {
					// line up the mask segments into two longs
					if isCompressed {
						parseData.maskers[i] = defaultMasker
					} else {
						bitsPerSegment := IPv6BitsPerSegment
						var maskVal uint64 = 0
						if missingSegmentCount >= 4 {
							cachedMasker := parseData.maskers[i]
							var extendedMaskVal uint64
							extendedCount := missingSegmentCount - 3
							for k := 0; k < extendedCount; k++ {
								extendedMaskVal = (extendedMaskVal << uint(bitsPerSegment)) | mask.GetSegment(normalizedSegmentIndex+k).getDivisionValue()
							}
							for k := extendedCount; k <= missingSegmentCount; k++ {
								maskVal = (maskVal << uint(bitsPerSegment)) | mask.GetSegment(normalizedSegmentIndex+k).getDivisionValue()
							}
							if cachedMasker == nil {
								// shift must be 6 bits at most for this shift to work per the java spec (so it must be less than 2^6 = 64)
								extendedMaxValue := ^(^DivInt(0) << uint(bits-DivIntSize))
								cachedMasker = MaskExtendedRange(
									lower, lowerHighBytes,
									upper, upperHighBytes,
									maskVal, extendedMaskVal,
									0xffffffffffffffff, extendedMaxValue)
								parseData.maskers[i] = cachedMasker
							}
							if !cachedMasker.IsSequential() && sections.maskError == nil {
								sections.maskError = &incompatibleAddressError{
									addressError: addressError{
										str: addressString,
										key: "ipaddress.error.maskMismatch",
									},
								}
							}
							masker := cachedMasker.(ExtendedMasker)
							maskedLowerHighBytes = masker.GetExtendedMaskedLower(lowerHighBytes, extendedMaskVal)
							maskedUpperHighBytes = masker.GetExtendedMaskedUpper(upperHighBytes, extendedMaskVal)
							maskedLower = masker.GetMaskedLower(lower, maskVal)
							maskedUpper = masker.GetMaskedUpper(upper, maskVal)
							maskedIsRange = (maskedLower != maskedUpper) || (maskedLowerHighBytes != maskedUpperHighBytes)
							maskedIsDifferent = maskedIsDifferent || maskedLower != lower || maskedUpper != upper || maskedLowerHighBytes != lowerHighBytes || maskedUpperHighBytes != upperHighBytes
						} else {
							masker := parseData.maskers[i]
							for k := 0; k <= missingSegmentCount; k++ {
								maskVal = (maskVal << uint(bitsPerSegment)) | mask.GetSegment(normalizedSegmentIndex+k).getDivisionValue()
							}
							if masker == nil {
								// shift must be 6 bits at most for this shift to work per the java spec (so it must be less than 2^6 = 64)
								maxValue := ^(^DivInt(0) << uint(bits))
								masker = MaskRange(lower, upper, maskVal, maxValue)
								parseData.maskers[i] = masker
							}
							if !masker.IsSequential() && sections.maskError == nil {
								sections.maskError = &incompatibleAddressError{
									addressError: addressError{
										str: maskString(lower, upper, maskVal),
										key: "ipaddress.error.maskMismatch",
									},
								}
							}
							maskedLower = masker.GetMaskedLower(lower, maskVal)
							maskedUpper = masker.GetMaskedUpper(upper, maskVal)
							maskedIsRange = maskedLower != maskedUpper
							maskedIsDifferent = maskedIsDifferent || maskedLower != lower || maskedUpper != upper
						}
					}
				} else {
					maskedLowerHighBytes = lowerHighBytes
					maskedUpperHighBytes = upperHighBytes
					maskedLower = lower
					maskedUpper = upper
					maskedIsRange = hostIsRange
				}
				shift := bits
				count := missingSegmentCount
				for count >= 0 { // add the missing segments
					currentPrefix := getSegmentPrefixLength(IPv6BitsPerSegment, prefLen, normalizedSegmentIndex)
					var hostSegLower, hostSegUpper, maskedSegLower, maskedSegUpper uint64
					if !isCompressed {
						shift -= IPv6BitsPerSegment
						if count >= 4 {
							shorterShift := shift - (IPv6BitsPerSegment << 2)
							hostSegLower = (lowerHighBytes >> uint(shorterShift)) & IPv6MaxValuePerSegment
							if hostIsRange {
								hostSegUpper = (upperHighBytes >> uint(shorterShift)) & IPv6MaxValuePerSegment
							} else {
								hostSegUpper = hostSegLower
							}
							if hasMask {
								maskedSegLower = (maskedLowerHighBytes >> uint(shorterShift)) & IPv6MaxValuePerSegment
								if maskedIsRange {
									maskedSegUpper = (maskedUpperHighBytes >> uint(shorterShift)) & IPv6MaxValuePerSegment
								} else {
									maskedSegUpper = maskedSegLower
								}
							} else {
								maskedSegLower = hostSegLower
								maskedSegUpper = hostSegUpper
							}
						} else {
							hostSegLower = (lower >> uint(shift)) & IPv6MaxValuePerSegment
							if hostIsRange {
								hostSegUpper = (upper >> uint(shift)) & IPv6MaxValuePerSegment
							} else {
								hostSegUpper = hostSegLower
							}
							if hasMask {
								maskedSegLower = (maskedLower >> uint(shift)) & IPv6MaxValuePerSegment
								if maskedIsRange {
									maskedSegUpper = (maskedUpper >> uint(shift)) & IPv6MaxValuePerSegment
								} else {
									maskedSegUpper = maskedSegLower
								}
							} else {
								maskedSegLower = hostSegLower
								maskedSegUpper = hostSegUpper
							}
						}
					}
					if doSections {
						if maskedIsDifferent || currentPrefix != nil {
							hostSegments = allocateSegments(hostSegments, segments, ipv6SegmentCount, normalizedSegmentIndex)
							hostSegments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
								addressString,
								IPv6,
								SegInt(hostSegLower),
								SegInt(hostSegUpper),
								false,
								i,
								nil,
								creator)
							isHostMultiple = isHostMultiple || segIsMult
						}
						segments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
							addressString,
							IPv6,
							SegInt(maskedSegLower),
							SegInt(maskedSegUpper),
							false,
							i,
							currentPrefix,
							creator)
						isMultiple = isMultiple || segIsMult
					}
					if doRangeBoundaries {
						isSegRange := maskedSegLower != maskedSegUpper
						if !doSections || isSegRange {
							if doSections {
								lowerSegments = allocateSegments(lowerSegments, segments, ipv6SegmentCount, normalizedSegmentIndex)
							} // else segments already allocated
							lowerSegments[normalizedSegmentIndex], _ = parseData.createSegment(
								addressString,
								IPv6,
								SegInt(maskedSegLower),
								SegInt(maskedSegLower),
								false,
								i,
								currentPrefix,
								creator)

						} else if lowerSegments != nil {
							lowerSegments[normalizedSegmentIndex] = segments[normalizedSegmentIndex]
						}
						if withUpper {
							if isSegRange {
								upperSegments = allocateSegments(upperSegments, lowerSegments, ipv6SegmentCount, normalizedSegmentIndex)
								upperSegments[normalizedSegmentIndex], _ = parseData.createSegment(
									addressString,
									IPv6,
									SegInt(maskedSegUpper),
									SegInt(maskedSegUpper),
									false,
									i,
									currentPrefix,
									creator)
							} else if upperSegments != nil {
								upperSegments[normalizedSegmentIndex] = lowerSegments[normalizedSegmentIndex]
							}
						}
					}
					normalizedSegmentIndex++
					count--
				}
				addressParseData.setBitLength(i, bits)
				continue
			} //end handle joined segments
		}

		hostLower, hostUpper := lower, upper
		var masker Masker
		unmasked := true
		if hasMask {
			masker = parseData.maskers[i]
			maskInt := uint64(mask.GetSegment(normalizedSegmentIndex).GetSegmentValue())
			if masker == nil {
				masker = MaskRange(lower, upper, maskInt, uint64(creator.getMaxValuePerSegment()))
				parseData.maskers[i] = masker
			}
			if !masker.IsSequential() && sections.maskError == nil {
				sections.maskError = &incompatibleAddressError{
					addressError: addressError{
						str: maskString(lower, upper, maskInt),
						key: "ipaddress.error.maskMismatch",
					},
				}
			}
			lower = masker.GetMaskedLower(lower, maskInt)
			upper = masker.GetMaskedUpper(upper, maskInt)
			unmasked = hostLower == lower && hostUpper == upper
			maskedIsDifferent = maskedIsDifferent || !unmasked
		}
		segmentPrefixLength := getSegmentPrefixLength(IPv6BitsPerSegment, prefLen, normalizedSegmentIndex)
		if doSections {
			if maskedIsDifferent || segmentPrefixLength != nil {
				hostSegments = allocateSegments(hostSegments, segments, ipv6SegmentCount, normalizedSegmentIndex)
				hostSegments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
					addressString,
					IPv6,
					SegInt(hostLower),
					SegInt(hostUpper),
					true,
					i,
					nil,
					creator)
				isHostMultiple = isHostMultiple || segIsMult
			}
			segments[normalizedSegmentIndex], segIsMult = parseData.createSegment(
				addressString,
				IPv6,
				SegInt(lower),
				SegInt(upper),
				unmasked,
				i,
				segmentPrefixLength,
				creator)
			isMultiple = isMultiple || segIsMult
		}
		if doRangeBoundaries {
			isRange := lower != upper
			if !doSections || isRange {
				if doSections {
					lowerSegments = allocateSegments(lowerSegments, segments, ipv6SegmentCount, normalizedSegmentIndex)
				} // else segments already allocated
				lowerSegments[normalizedSegmentIndex], _ = parseData.createSegment(
					addressString,
					IPv6,
					SegInt(lower),
					SegInt(lower),
					false,
					i,
					segmentPrefixLength,
					creator)

			} else if lowerSegments != nil {
				lowerSegments[normalizedSegmentIndex] = segments[normalizedSegmentIndex]
			}
			if withUpper {
				if isRange {
					upperSegments = allocateSegments(upperSegments, lowerSegments, ipv6SegmentCount, normalizedSegmentIndex)
					upperSegments[normalizedSegmentIndex], _ = parseData.createSegment(
						addressString,
						IPv6,
						SegInt(upper),
						SegInt(upper),
						false,
						i,
						segmentPrefixLength,
						creator)
				} else if upperSegments != nil {
					upperSegments[normalizedSegmentIndex] = lowerSegments[normalizedSegmentIndex]
				}
			}
		}
		normalizedSegmentIndex++
		addressParseData.setBitLength(i, IPv6BitsPerSegment)
	}
	prefLength := getPrefixLength(qualifier)
	if mixed {
		ipv4Range := parseData.mixedParsedAddress.getProviderSeqRange().ToIPv4()
		if hasMask && parseData.mixedMaskers == nil {
			parseData.mixedMaskers = make([]Masker, IPv4SegmentCount)
		}
		for n := 0; n < 2; n++ {
			m := n << 1
			segmentPrefixLength := getSegmentPrefixLength(IPv6BitsPerSegment, prefLen, normalizedSegmentIndex)
			//segmentPrefixLength := getQualifierSegmentPrefixLength(normalizedSegmentIndex, IPv6BitsPerSegment, qualifier)
			o := m + 1
			oneLow := ipv4Range.GetLower().GetSegment(m)
			twoLow := ipv4Range.GetLower().GetSegment(o)
			oneUp := ipv4Range.GetUpper().GetSegment(m)
			twoUp := ipv4Range.GetUpper().GetSegment(o)
			oneLower := oneLow.GetSegmentValue()
			twoLower := twoLow.GetSegmentValue()
			oneUpper := oneUp.GetSegmentValue()
			twoUpper := twoUp.GetSegmentValue()

			originalOneLower := oneLower
			originalTwoLower := twoLower
			originalOneUpper := oneUpper
			originalTwoUpper := twoUpper

			if hasMask {
				maskInt := uint64(mask.GetSegment(normalizedSegmentIndex).GetSegmentValue())
				shift := IPv4BitsPerSegment
				shiftedMask := maskInt >> uint(shift)
				masker := parseData.mixedMaskers[m]
				lstringLower := uint64(oneLower)
				lstringUpper := uint64(oneUpper)
				if masker == nil {
					masker = MaskRange(lstringLower, lstringUpper, shiftedMask, IPv4MaxValuePerSegment)
					parseData.mixedMaskers[m] = masker
				}
				if !masker.IsSequential() && sections.maskError == nil {
					sections.maskError = &incompatibleAddressError{
						addressError: addressError{
							str: maskString(lstringLower, lstringUpper, shiftedMask),
							key: "ipaddress.error.maskMismatch",
						},
					}
				}
				oneLower = SegInt(masker.GetMaskedLower(lstringLower, shiftedMask))
				oneUpper = SegInt(masker.GetMaskedUpper(lstringUpper, shiftedMask))
				lstringLower = uint64(twoLower)
				lstringUpper = uint64(twoUpper)
				masker = parseData.mixedMaskers[m+1]
				if masker == nil {
					masker = MaskRange(lstringLower, lstringUpper, maskInt, IPv4MaxValuePerSegment)
					parseData.mixedMaskers[m+1] = masker
				}
				if !masker.IsSequential() && sections.maskError == nil {
					sections.maskError = &incompatibleAddressError{
						addressError: addressError{
							str: maskString(lstringLower, lstringUpper, maskInt),
							key: "ipaddress.error.maskMismatch",
						},
					}
				}
				twoLower = SegInt(masker.GetMaskedLower(lstringLower, maskInt))
				twoUpper = SegInt(masker.GetMaskedUpper(lstringUpper, maskInt))
				maskedIsDifferent = maskedIsDifferent || oneLower != originalOneLower || oneUpper != originalOneUpper ||
					twoLower != originalTwoLower || twoUpper != originalTwoUpper
			}
			isRange := oneLower != oneUpper || twoLower != twoUpper
			if doSections {
				doHostSegment := maskedIsDifferent || segmentPrefixLength != nil
				if doHostSegment {
					hostSegments = allocateSegments(hostSegments, segments, ipv6SegmentCount, normalizedSegmentIndex)
				}
				if !isRange {
					if doHostSegment {
						hostSegments[normalizedSegmentIndex] = createIPv6Segment(originalOneLower, originalTwoLower, nil, creator)
					}
					segments[normalizedSegmentIndex] = createIPv6Segment(
						oneLower,
						twoLower,
						segmentPrefixLength,
						creator)
				} else {
					if doHostSegment {
						hostSegments[normalizedSegmentIndex] = createIPv6RangeSegment(
							&sections,
							ipv4Range,
							originalOneLower,
							originalOneUpper,
							originalTwoLower,
							originalTwoUpper,
							nil,
							creator)
					}
					segments[normalizedSegmentIndex] = createIPv6RangeSegment(
						&sections,
						ipv4Range,
						oneLower,
						oneUpper,
						twoLower,
						twoUpper,
						segmentPrefixLength,
						creator)
					isMultiple = true
				}
			}
			if doRangeBoundaries {
				if !doSections || isRange {
					if doSections {
						lowerSegments = allocateSegments(lowerSegments, segments, ipv6SegmentCount, normalizedSegmentIndex)
					} // else segments already allocated
					lowerSegments[normalizedSegmentIndex] = createIPv6Segment(
						oneLower,
						twoLower,
						segmentPrefixLength,
						creator)
				} else if lowerSegments != nil {
					lowerSegments[normalizedSegmentIndex] = segments[normalizedSegmentIndex]
				}
				if withUpper {
					if isRange {
						upperSegments = allocateSegments(upperSegments, lowerSegments, ipv6SegmentCount, normalizedSegmentIndex)
						upperSegments[normalizedSegmentIndex] = createIPv6Segment(
							oneUpper,
							twoUpper,
							segmentPrefixLength, // we must keep prefix length for upper to get prefix subnet creation
							creator)
					} else if upperSegments != nil {
						upperSegments[normalizedSegmentIndex] = lowerSegments[normalizedSegmentIndex]
					}
				}
			}
			normalizedSegmentIndex++
		}
	}
	var result, hostResult *IPAddressSection
	if doSections {
		if hostSegments != nil {
			hostResult = creator.createSectionInternal(hostSegments, isHostMultiple).ToIP()
			sections.hostSection = hostResult
			if checkExpandedValues(hostResult, expandedStart, expandedEnd) {
				sections.joinHostError = &incompatibleAddressError{addressError{str: addressString, key: "ipaddress.error.invalid.joined.ranges"}}
			}
		}
		result = creator.createPrefixedSectionInternal(segments, isMultiple, prefLength)
		sections.section = result
		if checkExpandedValues(result, expandedStart, expandedEnd) {
			sections.joinAddressError = &incompatibleAddressError{addressError{str: addressString, key: "ipaddress.error.invalid.joined.ranges"}}
			if hostResult == nil {
				sections.joinHostError = sections.joinAddressError
			}
		}
	}
	if doRangeBoundaries {
		prefixLength := getPrefixLength(qualifier)
		isPrefixSub := false
		if prefixLength != nil {
			var lowerSegs, upperSegs []*AddressDivision
			if doSections {
				lowerSegs = segments
				upperSegs = segments
			} else {
				lowerSegs = lowerSegments
				if upperSegments == nil {
					upperSegs = lowerSegments
				} else {
					upperSegs = upperSegments
				}
			}
			isPrefixSub = isPrefixSubnet(
				func(index int) SegInt { return lowerSegs[index].ToSegmentBase().GetSegmentValue() },
				func(index int) SegInt { return upperSegs[index].ToSegmentBase().GetUpperSegmentValue() },
				len(lowerSegs),
				IPv6BytesPerSegment,
				IPv6BitsPerSegment,
				IPv6MaxValuePerSegment,
				prefixLength.bitCount(),
				zerosOnly)
			if isPrefixSub {
				if lowerSegments == nil {
					//allocate lower segments from address segments
					lowerSegments = allocateSegments(lowerSegments, segments, ipv6SegmentCount, ipv6SegmentCount)
				}
				if upperSegments == nil {
					//allocate upper segments from lower segments
					upperSegments = allocateSegments(upperSegments, lowerSegments, ipv6SegmentCount, ipv6SegmentCount)
				}
			}
		}
		if lowerSegments != nil {
			boundaries.lowerSection = creator.createPrefixedSectionInternalSingle(lowerSegments, false, prefLength)
		}
		if upperSegments != nil {
			section := creator.createPrefixedSectionInternal(upperSegments, false, prefLength)
			if isPrefixSub {
				section = section.ToPrefixBlock()
			}
			boundaries.upperSection = section.GetUpper()
		}
	}
	return
}

func maskString(lower, upper, maskInt uint64) string {
	return strconv.FormatUint(lower, 10) + "-" + strconv.FormatUint(upper, 10) + " /" + strconv.FormatUint(maskInt, 10)
}

// When expanding a set of segments into multiple, it is possible that the new segments do not accurately
// cover the same ranges of values.  This occurs when there is a range in the upper segments and the lower
// segments do not cover the full range (as is the case in the original unexpanded segment).
//
// This does not include compressed 0 segments or compressed '*' segments, as neither can have the issue.
//
// Returns true if the expansion was invalid.
func checkExpandedValues(section *IPAddressSection, start, end int) bool {
	if section != nil && start < end {
		seg := section.GetSegment(start)
		lastWasRange := seg.isMultiple()
		for {
			start++
			seg = section.GetSegment(start)
			if lastWasRange {
				if !seg.IsFullRange() {
					return true
				}
			} else {
				lastWasRange = seg.isMultiple()
			}
			if start >= end {
				break
			}
		}
	}
	return false
}

func (parseData *parsedIPAddress) createSegment(
	addressString string,
	version IPVersion,
	val,
	upperVal SegInt,
	useFlags bool,
	parsedSegIndex int,
	segmentPrefixLength PrefixLen,
	creator parsedAddressCreator) (div *AddressDivision, isMultiple bool) {
	parsed := parseData.getAddressParseData()
	if val != upperVal {
		return createRangeSeg(addressString, version, val, upperVal,
			useFlags, parsed, parsedSegIndex,
			segmentPrefixLength, creator), true
	}
	var result *AddressDivision
	if !useFlags {
		result = creator.createSegment(val, val, segmentPrefixLength)
	} else {
		result = creator.createSegmentInternal(
			val,
			segmentPrefixLength,
			addressString,
			val,
			parsed.getFlag(parsedSegIndex, keyStandardStr),
			parsed.getIndex(parsedSegIndex, keyLowerStrStartIndex),
			parsed.getIndex(parsedSegIndex, keyLowerStrEndIndex))
	}
	return result, false
}

// create an IPv6 segment by joining two IPv4 segments
func createIPv6Segment(value1, value2 SegInt, segmentPrefixLength PrefixLen, creator parsedAddressCreator) *AddressDivision {
	value := (value1 << uint(IPv4BitsPerSegment)) | value2
	result := creator.createPrefixSegment(value, segmentPrefixLength)
	return result
}

// create an IPv6 segment by joining two IPv4 segments
func createIPv6RangeSegment(
	//finalResult *translatedResult,
	sections *sectionResult,
	_ *SequentialRange[*IPv4Address], // this was only used to be put into any exceptions
	upperRangeLower,
	upperRangeUpper,
	lowerRangeLower,
	lowerRangeUpper SegInt,
	segmentPrefixLength PrefixLen,
	creator ipAddressCreator) *AddressDivision {
	shift := IPv4BitsPerSegment
	if upperRangeLower != upperRangeUpper {
		//if the high segment has a range, the low segment must match the full range,
		//otherwise it is not possible to create an equivalent IPv6 range when joining two IPv4 ranges
		if sections.mixedError == nil && lowerRangeLower != 0 || lowerRangeUpper != IPv4MaxValuePerSegment {
			sections.mixedError = &incompatibleAddressError{
				addressError: addressError{
					key: "ipaddress.error.invalidMixedRange",
				},
			}
		}
	}
	return creator.createSegment(
		(upperRangeLower<<uint(shift))|lowerRangeLower,
		(upperRangeUpper<<uint(shift))|lowerRangeUpper,
		segmentPrefixLength)
}

func createRangeSeg(
	addressString string,
	_ IPVersion,
	stringLower,
	stringUpper SegInt,
	useFlags bool,
	parseData *addressParseData,
	parsedSegIndex int,
	segmentPrefixLength PrefixLen,
	creator parsedAddressCreator) *AddressDivision {
	var lower, upper = stringLower, stringUpper
	var result *AddressDivision
	if !useFlags {
		result = creator.createSegment(lower, upper, segmentPrefixLength)
	} else {
		result = creator.createRangeSegmentInternal(
			lower,
			upper,
			segmentPrefixLength,
			addressString,
			stringLower,
			stringUpper,
			parseData.getFlag(parsedSegIndex, keyStandardStr),
			parseData.getFlag(parsedSegIndex, keyStandardRangeStr),
			parseData.getIndex(parsedSegIndex, keyLowerStrStartIndex),
			parseData.getIndex(parsedSegIndex, keyLowerStrEndIndex),
			parseData.getIndex(parsedSegIndex, keyUpperStrEndIndex))
	}
	return result
}

func createFullRangeSegment(
	version IPVersion,
	stringLower,
	stringUpper SegInt,
	parsedSegIndex int,
	segmentPrefixLength PrefixLen,
	mask *SegInt,
	creator parsedAddressCreator) (result, hostResult, lower, upper *AddressDivision, err addrerr.IncompatibleAddressError) {
	var maskedLower, maskedUpper SegInt
	hasMask := mask != nil
	maskedIsDifferent := false
	if hasMask {
		maskInt := DivInt(*mask)
		lstringLower := uint64(stringLower)
		lstringUpper := uint64(stringUpper)
		masker := MaskRange(lstringLower, lstringUpper, maskInt, uint64(creator.getMaxValuePerSegment()))
		if !masker.IsSequential() {
			err = &incompatibleAddressError{
				addressError{
					str: maskString(lstringLower, lstringUpper, maskInt),
					key: "ipaddress.error.maskMismatch",
				},
			}
		}
		maskedLower = SegInt(masker.GetMaskedLower(lstringLower, maskInt))
		maskedUpper = SegInt(masker.GetMaskedUpper(lstringUpper, maskInt))
		maskedIsDifferent = maskedLower != stringLower || maskedUpper != stringUpper
	} else {
		maskedLower = stringLower
		maskedUpper = stringUpper
	}
	result = createRangeSeg("", version, maskedLower, maskedUpper,
		false, nil, parsedSegIndex, segmentPrefixLength, creator)
	if maskedIsDifferent || segmentPrefixLength != nil {
		hostResult = createRangeSeg("", version, stringLower, stringUpper,
			false, nil, parsedSegIndex, nil, creator)
	} else {
		hostResult = result
	}
	if maskedLower == maskedUpper {
		lower = result
		upper = result
	} else {
		lower = createRangeSeg("", version, maskedLower, maskedLower,
			false, nil, parsedSegIndex, segmentPrefixLength, creator)
		upper = createRangeSeg("", version, maskedUpper, maskedUpper,
			false, nil, parsedSegIndex, segmentPrefixLength, creator)
	}
	return
}

func createAllAddress(
	version IPVersion,
	qualifier *parsedHostIdentifierStringQualifier,
	originator HostIdentifierString) (res, hostAddr, lower, upper *IPAddress, err addrerr.IncompatibleAddressError) {
	creator := version.toType().getIPNetwork().getIPAddressCreator()
	mask := qualifier.getMaskLower()
	if mask != nil && mask.GetBlockMaskPrefixLen(true) != nil {
		mask = nil //we don't do any masking if the mask is a subnet mask, instead we just map it to the corresponding prefix length
	}
	segmentCount := version.GetSegmentCount()
	segments := make([]*AddressDivision, segmentCount)
	hostSegments := make([]*AddressDivision, segmentCount)
	lowerSegments := make([]*AddressDivision, segmentCount)
	upperSegments := make([]*AddressDivision, segmentCount)
	segMaxVal := creator.getMaxValuePerSegment()
	hasMask := mask != nil
	prefLen := getPrefixLength(qualifier)
	bitsPerSegment := version.GetBitsPerSegment()
	for i := 0; i < segmentCount; i++ {
		var segmentMask *SegInt
		if hasMask {
			maskVal := mask.getSegment(i).getSegmentValue()
			segmentMask = &maskVal
		}
		newSeg, hostSeg, lowSeg, upperSeg, rngErr := createFullRangeSegment(
			version,
			0,
			segMaxVal,
			i,
			getSegmentPrefixLength(bitsPerSegment, prefLen, i),
			segmentMask,
			creator)
		if rngErr != nil && err == nil {
			err = rngErr
		}
		segments[i] = newSeg
		hostSegments[i] = hostSeg
		lowerSegments[i] = lowSeg
		upperSegments[i] = upperSeg
	}
	if err == nil {
		section := creator.createPrefixedSectionInternal(segments, true, prefLen)
		res = creator.createAddressInternalFromSection(section, qualifier.getZone(), originator).ToIP()
	}
	hostSection := creator.createSectionInternal(hostSegments, true)
	hostAddr = creator.createAddressInternal(hostSection.ToSectionBase(), nil).ToIP()
	lowerSection := creator.createPrefixedSectionInternal(lowerSegments, false, prefLen)
	lower = creator.createAddressInternal(lowerSection.ToSectionBase(), nil).ToIP()
	upperSection := creator.createPrefixedSectionInternal(upperSegments, false, prefLen)
	upper = creator.createAddressInternal(upperSection.ToSectionBase(), nil).ToIP()
	return
}

func getPrefixLength(qualifier *parsedHostIdentifierStringQualifier) PrefixLen {
	return qualifier.getEquivalentPrefixLen()
}
