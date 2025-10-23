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
	"math/big"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

func createIPSection(segments []*AddressDivision, prefixLength PrefixLen, addrType addrType) *IPAddressSection {
	sect := &IPAddressSection{
		ipAddressSectionInternal{
			addressSectionInternal{
				addressDivisionGroupingInternal{
					addressDivisionGroupingBase: addressDivisionGroupingBase{
						divisions:    standardDivArray(segments),
						addrType:     addrType,
						cache:        &valueCache{},
						prefixLength: prefixLength,
					},
				},
			},
		},
	}
	assignStringCache(&sect.addressDivisionGroupingBase, addrType)
	return sect
}

func createIPSectionFromSegs(isIPv4 bool, orig []*IPAddressSegment, prefLen PrefixLen) (result *IPAddressSection) {
	segProvider := func(index int) *IPAddressSegment {
		return orig[index]
	}
	var divs []*AddressDivision
	var newPref PrefixLen
	var isMultiple bool
	if isIPv4 {
		divs, newPref, isMultiple = createDivisionsFromSegs(
			segProvider,
			len(orig),
			ipv4BitsToSegmentBitshift,
			IPv4BitsPerSegment,
			IPv4BytesPerSegment,
			IPv4MaxValuePerSegment,
			zeroIPv4Seg.ToIP(),
			zeroIPv4SegZeroPrefix.ToIP(),
			zeroIPv4SegPrefixBlock.ToIP(),
			prefLen)
		result = createIPv4Section(divs).ToIP()
	} else {
		divs, newPref, isMultiple = createDivisionsFromSegs(
			segProvider,
			len(orig),
			ipv6BitsToSegmentBitshift,
			IPv6BitsPerSegment,
			IPv6BytesPerSegment,
			IPv6MaxValuePerSegment,
			zeroIPv6Seg.ToIP(),
			zeroIPv6SegZeroPrefix.ToIP(),
			zeroIPv6SegPrefixBlock.ToIP(),
			prefLen)
		result = createIPv6Section(divs).ToIP()
	}
	result.prefixLength = newPref
	result.isMult = isMultiple
	return result
}

// Callers to this function have already initialized the segments to have consistent prefix lengths,
// but in here we need to determine what that prefix length might be.
func deriveIPAddressSection(from *IPAddressSection, segments []*AddressDivision) (res *IPAddressSection) {
	res = createIPSection(segments, nil, from.getAddrType())
	res.initMultAndPrefLen()
	return
}

// Callers to this function have already initialized the segments to have prefix lengths corresponding to the supplied argument prefixLength
// So we need only check if multiple and assign the prefix length.
func deriveIPAddressSectionPrefLen(from *IPAddressSection, segments []*AddressDivision, prefixLength PrefixLen) (res *IPAddressSection) {
	res = createIPSection(segments, prefixLength, from.getAddrType())
	res.initMultiple()
	return
}

type ipAddressSectionInternal struct {
	addressSectionInternal
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (section *ipAddressSectionInternal) GetSegment(index int) *IPAddressSegment {
	return section.getDivision(index).ToIP()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (section *ipAddressSectionInternal) ForEachSegment(consumer func(segmentIndex int, segment *IPAddressSegment) (stop bool)) int {
	divArray := section.getDivArray()
	if divArray != nil {
		for i, div := range divArray {
			if consumer(i, div.ToIP()) {
				return i + 1
			}
		}
	}
	return len(divArray)
}

// GetIPVersion returns the IP version of this IP address section.
func (section *ipAddressSectionInternal) GetIPVersion() IPVersion {
	addrType := section.getAddrType()
	if addrType.isIPv4() {
		return IPv4
	} else if addrType.isIPv6() {
		return IPv6
	}
	return IndeterminateIPVersion
}

func (section *ipAddressSectionInternal) getNetworkPrefixLen() PrefixLen {
	return section.prefixLength
}

// GetNetworkPrefixLen returns the prefix length, or nil if there is no prefix length.  It is equivalent to GetPrefixLen.
//
// A prefix length indicates the number of bits in the initial part of the address item that comprises the prefix.
//
// A prefix is a part of the address item that is not specific to that address but common amongst a group of such items, such as a CIDR prefix block subnet.
func (section *ipAddressSectionInternal) GetNetworkPrefixLen() PrefixLen {
	return section.getNetworkPrefixLen().copy()
}

// GetBlockMaskPrefixLen returns the prefix length if this address section is equivalent to the mask for a CIDR prefix block.
// Otherwise, it returns nil.
// A CIDR network mask is an address section with all ones in the network section and then all zeros in the host section.
// A CIDR host mask is an address section with all zeros in the network section and then all ones in the host section.
// The prefix length is the bit-length of the network section.
//
// Also, keep in mind that the prefix length returned by this method is not equivalent to the prefix length of this instance,
// indicating the network and host section of this address section.
// The prefix length returned here indicates the whether the value of this address can be used as a mask for the network and host
// section of any other address.  Therefore the two values can be different values, or one can be nil while the other is not.
//
// This method applies only to the lower value of the range if this section represents multiple values.
func (section *ipAddressSectionInternal) GetBlockMaskPrefixLen(network bool) PrefixLen {
	cache := section.cache
	if cache == nil {
		return nil // no prefix
	}
	cachedMaskLens := (*maskLenSetting)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.cachedMaskLens))))
	if cachedMaskLens == nil {
		networkMaskLen, hostMaskLen := section.checkForPrefixMask()
		cachedMaskLens = &maskLenSetting{networkMaskLen, hostMaskLen}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.cachedMaskLens))
		atomicStorePointer(dataLoc, unsafe.Pointer(cachedMaskLens))
	}
	if network {
		return cachedMaskLens.networkMaskLen
	}
	return cachedMaskLens.hostMaskLen
}

func (section *ipAddressSectionInternal) checkForPrefixMask() (networkMaskLen, hostMaskLen PrefixLen) {
	count := section.GetSegmentCount()
	if count == 0 {
		return
	}
	firstSeg := section.GetSegment(0)
	checkingNetworkFront, checkingHostFront := true, true
	var checkingNetworkBack, checkingHostBack bool
	var prefixedSeg int
	prefixedSegPrefixLen := BitCount(0)
	maxVal := firstSeg.GetMaxValue()
	for i := 0; i < count; i++ {
		seg := section.GetSegment(i)
		val := seg.GetSegmentValue()
		if val == 0 {
			if checkingNetworkFront {
				prefixedSeg = i
				checkingNetworkFront, checkingNetworkBack = false, true
			} else if !checkingHostFront && !checkingNetworkBack {
				return
			}
			checkingHostBack = false
		} else if val == maxVal {
			if checkingHostFront {
				prefixedSeg = i
				checkingHostFront, checkingHostBack = false, true
			} else if !checkingHostBack && !checkingNetworkFront {
				return
			}
			checkingNetworkBack = false
		} else {
			segNetworkMaskLen, segHostMaskLen := seg.checkForPrefixMask()
			if segNetworkMaskLen != nil {
				if checkingNetworkFront {
					prefixedSegPrefixLen = segNetworkMaskLen.bitCount()
					checkingNetworkBack = true
					checkingHostBack = false
					prefixedSeg = i
				} else {
					return
				}
			} else if segHostMaskLen != nil {
				if checkingHostFront {
					prefixedSegPrefixLen = segHostMaskLen.bitCount()
					checkingHostBack = true
					checkingNetworkBack = false
					prefixedSeg = i
				} else {
					return
				}
			} else {
				return
			}
			checkingNetworkFront, checkingHostFront = false, false
		}
	}
	if checkingNetworkFront {
		// all ones
		networkMaskLen = cacheBitCount(section.GetBitCount())
		hostMaskLen = cacheBitCount(0)
	} else if checkingHostFront {
		// all zeros
		hostMaskLen = cacheBitCount(section.GetBitCount())
		networkMaskLen = cacheBitCount(0)
	} else if checkingNetworkBack {
		// ending in zeros, network mask
		networkMaskLen = getNetworkPrefixLen(firstSeg.GetBitCount(), prefixedSegPrefixLen, prefixedSeg)
	} else if checkingHostBack {
		// ending in ones, host mask
		hostMaskLen = getNetworkPrefixLen(firstSeg.GetBitCount(), prefixedSegPrefixLen, prefixedSeg)
	}
	return
}

// IncludesZeroHost returns whether the address section contains an individual address section with a host of zero.  If the address section has no prefix length it returns false.
// If the prefix length matches the bit count, then it returns true.
//
// Otherwise, it checks whether it contains an individual address section for which all bits past the prefix are zero.
func (section *ipAddressSectionInternal) IncludesZeroHost() bool {
	networkPrefixLength := section.getPrefixLen()
	return networkPrefixLength != nil && section.IncludesZeroHostLen(networkPrefixLength.bitCount())
}

// IncludesZeroHostLen returns whether the address section contains an individual section with a host of zero, a section for which all bits past the given prefix length are zero.
func (section *ipAddressSectionInternal) IncludesZeroHostLen(networkPrefixLength BitCount) bool {
	networkPrefixLength = checkSubnet(section, networkPrefixLength)
	bitsPerSegment := section.GetBitsPerSegment()
	bytesPerSegment := section.GetBytesPerSegment()
	prefixedSegmentIndex := getHostSegmentIndex(networkPrefixLength, bytesPerSegment, bitsPerSegment)
	divCount := section.GetSegmentCount()
	for i := prefixedSegmentIndex; i < divCount; i++ {
		div := section.GetSegment(i)
		segmentPrefixLength := getPrefixedSegmentPrefixLength(bitsPerSegment, networkPrefixLength, i)
		mask := div.GetSegmentHostMask(segmentPrefixLength.bitCount())
		if (mask & div.GetSegmentValue()) != 0 {
			return false
		}
		for i++; i < divCount; i++ {
			div = section.GetSegment(i)
			if !div.includesZero() {
				return false
			}
		}
	}
	return true
}

// IncludesMaxHost returns whether the address section contains an individual address section with a host of all one-bits.  If the address section has no prefix length it returns false.
// If the prefix length matches the bit count, then it returns true.
//
// Otherwise, it checks whether it contains an individual address section for which all bits past the prefix are one.
func (section *ipAddressSectionInternal) IncludesMaxHost() bool {
	networkPrefixLength := section.getPrefixLen()
	return networkPrefixLength != nil && section.IncludesMaxHostLen(networkPrefixLength.bitCount())
}

// IncludesMaxHostLen returns whether the address section contains an individual address section with a host of all one-bits, an address section for which all bits past the given prefix length are all ones.
func (section *ipAddressSectionInternal) IncludesMaxHostLen(networkPrefixLength BitCount) bool {
	networkPrefixLength = checkSubnet(section, networkPrefixLength)
	bitsPerSegment := section.GetBitsPerSegment()
	bytesPerSegment := section.GetBytesPerSegment()
	prefixedSegmentIndex := getHostSegmentIndex(networkPrefixLength, bytesPerSegment, bitsPerSegment)
	divCount := section.GetSegmentCount()
	for i := prefixedSegmentIndex; i < divCount; i++ {
		div := section.GetSegment(i)
		segmentPrefixLength := getPrefixedSegmentPrefixLength(bitsPerSegment, networkPrefixLength, i)
		if segmentPrefixLength != nil {
			mask := div.GetSegmentHostMask(segmentPrefixLength.bitCount())
			if (mask & div.getUpperSegmentValue()) != mask {
				return false
			}
			for i++; i < divCount; i++ {
				div = section.GetSegment(i)
				if !div.includesMax() {
					return false
				}
			}
		}
	}
	return true
}

func (section *ipAddressSectionInternal) toZeroHost(boundariesOnly bool) (res *IPAddressSection, err addrerr.IncompatibleAddressError) {
	segmentCount := section.GetSegmentCount()
	if segmentCount == 0 {
		return section.toIPAddressSection(), nil
	}
	var prefLen BitCount
	if section.isPrefixed() {
		prefLen = section.getPrefixLen().bitCount()
	}
	if section.IsZeroHostLen(prefLen) {
		return section.toIPAddressSection(), nil
	}
	if section.IncludesZeroHost() && section.IsSingleNetwork() {
		res = section.getLower().ToIP() //cached
		return
	}
	if !section.isPrefixed() {
		mask := section.addrType.getIPNetwork().GetPrefixedNetworkMask(0)
		res = mask.GetSubSection(0, segmentCount)
		return
	}
	return section.createZeroHost(prefLen, boundariesOnly)
}

// boundariesOnly: whether we care if the masking works for all values in a range.
// For instance, 1.2.3.2-4/31 cannot be zero-hosted, because applyng to the boundaries results in 1.2.3.2-4/31,
// and that includes 1.2.3.3/31 which does not have host of zero.
// So in that case, we'd normally haveaddrerr.IncompatibleAddressError.  boundariesOnly as true avoids the exception,
// if we are really just interested in getting the zero-host boundaries,
// and we don't care about the remaining values in-between.
func (section *ipAddressSectionInternal) createZeroHost(prefLen BitCount, boundariesOnly bool) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	mask := section.addrType.getIPNetwork().GetNetworkMask(prefLen)
	return section.getSubnetSegments(
		getNetworkSegmentIndex(prefLen, section.GetBytesPerSegment(), section.GetBitsPerSegment()),
		cacheBitCount(prefLen),
		!boundariesOnly, //verifyMask
		section.getDivision,
		func(i int) SegInt { return mask.GetSegment(i).GetSegmentValue() })
}

func (section *ipAddressSectionInternal) toZeroHostLen(prefixLength BitCount) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	var minIndex int
	if section.isPrefixed() {
		existingPrefLen := section.getNetworkPrefixLen().bitCount()
		if prefixLength == existingPrefLen {
			return section.toZeroHost(false)
		}
		if prefixLength < existingPrefLen {
			minIndex = getNetworkSegmentIndex(prefixLength, section.GetBytesPerSegment(), section.GetBitsPerSegment())
		} else {
			minIndex = getNetworkSegmentIndex(existingPrefLen, section.GetBytesPerSegment(), section.GetBitsPerSegment())
		}
	} else {
		minIndex = getNetworkSegmentIndex(prefixLength, section.GetBytesPerSegment(), section.GetBitsPerSegment())
	}
	mask := section.addrType.getIPNetwork().GetNetworkMask(prefixLength)
	return section.getSubnetSegments(
		minIndex,
		nil, // intentionally no prefix length
		true,
		section.getDivision,
		func(i int) SegInt { return mask.GetSegment(i).GetSegmentValue() })
}

func (section *ipAddressSectionInternal) toZeroNetwork() *IPAddressSection {
	segmentCount := section.GetSegmentCount()
	if segmentCount == 0 {
		return section.toIPAddressSection()
	}
	if !section.isPrefixed() {
		mask := section.addrType.getIPNetwork().GetHostMask(section.GetBitCount())
		return mask.GetSubSection(0, segmentCount)
	}
	return section.createZeroNetwork()
}

func (section *ipAddressSectionInternal) createZeroNetwork() *IPAddressSection {
	prefixLength := section.getNetworkPrefixLen() // we know it is prefixed here so no panic on the derefence
	mask := section.addrType.getIPNetwork().GetHostMask(prefixLength.bitCount())
	res, _ := section.getSubnetSegments(
		0,
		prefixLength,
		false,
		section.getDivision,
		func(i int) SegInt { return mask.GetSegment(i).GetSegmentValue() })
	return res
}

func (section *ipAddressSectionInternal) toMaxHost() (res *IPAddressSection, err addrerr.IncompatibleAddressError) {
	segmentCount := section.GetSegmentCount()
	if segmentCount == 0 {
		return section.toIPAddressSection(), nil
	}
	if !section.isPrefixed() {
		mask := section.addrType.getIPNetwork().GetPrefixedHostMask(0)
		res = mask.GetSubSection(0, segmentCount)
		return
	}
	if section.IsMaxHostLen(section.getPrefixLen().bitCount()) {
		return section.toIPAddressSection(), nil
	}
	if section.IncludesMaxHost() && section.IsSingleNetwork() {
		return section.getUpper().ToIP(), nil // cached
	}
	return section.createMaxHost()
}

func (section *ipAddressSectionInternal) createMaxHost() (*IPAddressSection, addrerr.IncompatibleAddressError) {
	prefixLength := section.getNetworkPrefixLen() // we know it is prefixed here so no panic on the derefence
	mask := section.addrType.getIPNetwork().GetHostMask(prefixLength.bitCount())
	return section.getOredSegments(
		prefixLength,
		true,
		section.getDivision,
		func(i int) SegInt { return mask.GetSegment(i).GetSegmentValue() })
}

func (section *ipAddressSectionInternal) toMaxHostLen(prefixLength BitCount) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	if section.isPrefixed() && prefixLength == section.getNetworkPrefixLen().bitCount() {
		return section.toMaxHost()
	}
	mask := section.addrType.getIPNetwork().GetHostMask(prefixLength)
	return section.getOredSegments(
		nil,
		true,
		section.getDivision,
		func(i int) SegInt { return mask.GetSegment(i).GetSegmentValue() })
}

// IsSingleNetwork returns whether the network section of the address, the prefix, consists of a single value.
//
// If it has no prefix length, it returns true if not multiple, if it contains only a single individual address section.
func (section *ipAddressSectionInternal) IsSingleNetwork() bool {
	networkPrefixLength := section.getNetworkPrefixLen()
	if networkPrefixLength == nil {
		return !section.isMultiple()
	}
	prefLen := networkPrefixLength.bitCount()
	if prefLen >= section.GetBitCount() {
		return !section.isMultiple()
	}
	bitsPerSegment := section.GetBitsPerSegment()
	prefixedSegmentIndex := getNetworkSegmentIndex(prefLen, section.GetBytesPerSegment(), bitsPerSegment)
	if prefixedSegmentIndex < 0 {
		return true
	}
	for i := 0; i < prefixedSegmentIndex; i++ {
		if section.getDivision(i).isMultiple() {
			return false
		}
	}
	div := section.GetSegment(prefixedSegmentIndex)
	divPrefLen := getPrefixedSegmentPrefixLength(bitsPerSegment, prefLen, prefixedSegmentIndex)
	shift := bitsPerSegment - divPrefLen.bitCount()
	return (div.GetSegmentValue() >> uint(shift)) == (div.GetUpperSegmentValue() >> uint(shift))
}

// IsMaxHost returns whether this section has a prefix length and if so,
// whether the host is all all one-bits, the max value, for all individual sections in this address section.
//
// If the host section is zero length (there are zero host bits), IsMaxHost returns true.
func (section *ipAddressSectionInternal) IsMaxHost() bool {
	if !section.isPrefixed() {
		return false
	}
	return section.IsMaxHostLen(section.getNetworkPrefixLen().bitCount())
}

// IsMaxHostLen returns whether the host host is all one-bits, the max value, for all individual sections in this address section,
// for the given prefix length, the host being the bits following the prefix.
//
// If the host section is zero length (there are zero host bits), IsMaxHostLen returns true.
func (section *ipAddressSectionInternal) IsMaxHostLen(prefLen BitCount) bool {
	divCount := section.GetSegmentCount()
	if divCount == 0 {
		return true
	} else if prefLen < 0 {
		prefLen = 0
	}
	bytesPerSegment := section.GetBytesPerSegment()
	bitsPerSegment := section.GetBitsPerSegment()
	// Note: 1.2.3.4/32 has a max host
	prefixedSegmentIndex := getHostSegmentIndex(prefLen, bytesPerSegment, bitsPerSegment)
	if prefixedSegmentIndex < divCount {
		segmentPrefixLength := getPrefixedSegmentPrefixLength(bitsPerSegment, prefLen, prefixedSegmentIndex)
		i := prefixedSegmentIndex
		div := section.GetSegment(i)
		mask := div.GetSegmentHostMask(segmentPrefixLength.bitCount())
		if div.isMultiple() || (mask&div.getSegmentValue()) != mask {
			return false
		}
		i++
		for ; i < divCount; i++ {
			div = section.GetSegment(i)
			if !div.IsMax() {
				return false
			}
		}
	}
	return true
}

// IsZeroHost returns whether this section has a prefix length and if so,
// whether the host section is always zero for all individual sections in this address section.
//
// If the host section is zero length (there are zero host bits), IsZeroHost returns true.
func (section *ipAddressSectionInternal) IsZeroHost() bool {
	if !section.isPrefixed() {
		return false
	}
	return section.IsZeroHostLen(section.getNetworkPrefixLen().bitCount())
}

// IsZeroHostLen returns whether the host section is always zero for all individual sections in this address section,
// for the given prefix length.
//
// If the host section is zero length (there are zero host bits), IsZeroHostLen returns true.
func (section *ipAddressSectionInternal) IsZeroHostLen(prefLen BitCount) bool {
	segmentCount := section.GetSegmentCount()
	if segmentCount == 0 {
		return true
	} else if prefLen < 0 {
		prefLen = 0
	}
	bitsPerSegment := section.GetBitsPerSegment()
	// Note: 1.2.3.4/32 has a zero host
	prefixedSegmentIndex := getHostSegmentIndex(prefLen, section.GetBytesPerSegment(), bitsPerSegment)
	if prefixedSegmentIndex < segmentCount {
		segmentPrefixLength := getPrefixedSegmentPrefixLength(bitsPerSegment, prefLen, prefixedSegmentIndex)
		i := prefixedSegmentIndex
		div := section.GetSegment(i)
		if div.isMultiple() || (div.GetSegmentHostMask(segmentPrefixLength.bitCount())&div.getSegmentValue()) != 0 {
			return false
		}
		for i++; i < segmentCount; i++ {
			div := section.GetSegment(i)
			if !div.IsZero() {
				return false
			}
		}
	}
	return true
}

func (section *ipAddressSectionInternal) adjustPrefixLength(adjustment BitCount, withZeros bool) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	if adjustment == 0 && section.isPrefixed() {
		return section.toIPAddressSection(), nil
	}
	prefix := section.getAdjustedPrefix(adjustment)
	sec, err := section.setPrefixLength(prefix, withZeros)
	return sec.ToIP(), err
}

func (section *ipAddressSectionInternal) adjustPrefixLen(adjustment BitCount) *IPAddressSection {
	// no zeroing
	res, _ := section.adjustPrefixLength(adjustment, false)
	return res
}

func (section *ipAddressSectionInternal) adjustPrefixLenZeroed(adjustment BitCount) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	return section.adjustPrefixLength(adjustment, true)
}

func (section *ipAddressSectionInternal) withoutPrefixLen() *IPAddressSection {
	if !section.isPrefixed() {
		return section.toIPAddressSection()
	}
	if section.hasNoDivisions() {
		return createIPSection(section.getDivisionsInternal(), nil, section.getAddrType())
	}
	existingPrefixLength := section.getPrefixLen().bitCount()
	maxVal := section.GetMaxSegmentValue()
	var startIndex int
	if existingPrefixLength > 0 {
		bitsPerSegment := section.GetBitsPerSegment()
		bytesPerSegment := section.GetBytesPerSegment()
		startIndex = getNetworkSegmentIndex(existingPrefixLength, bytesPerSegment, bitsPerSegment)
	}
	res, _ := section.getSubnetSegments(
		startIndex,
		nil,
		false,
		func(i int) *AddressDivision {
			return section.getDivision(i)
		},
		func(i int) SegInt {
			return maxVal
		},
	)
	return res
}

func (section *ipAddressSectionInternal) checkSegmentCount(other *IPAddressSection) addrerr.SizeMismatchError {
	if other.GetSegmentCount() != section.GetSegmentCount() {
		return &sizeMismatchError{incompatibleAddressError{addressError{key: "ipaddress.error.sizeMismatch"}}}
	}
	return nil
}

func (section *ipAddressSectionInternal) checkMaskSegmentCount(other *IPAddressSection) addrerr.SizeMismatchError {
	if other.GetSegmentCount() < section.GetSegmentCount() {
		return &sizeMismatchError{incompatibleAddressError{addressError{key: "ipaddress.error.sizeMismatch"}}}
	}
	return nil
}

// error can be addrerr.IncompatibleAddressError or addrerr.SizeMismatchError
func (section *ipAddressSectionInternal) mask(msk *IPAddressSection, retainPrefix bool) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	if err := section.checkMaskSegmentCount(msk); err != nil {
		return nil, err
	}
	var prefLen PrefixLen
	if retainPrefix {
		prefLen = section.getPrefixLen()
	}
	return section.getSubnetSegments(
		0,
		prefLen,
		true,
		section.getDivision,
		func(i int) SegInt { return msk.GetSegment(i).GetSegmentValue() })
}

// error can be addrerr.IncompatibleAddressError or addrerr.SizeMismatchError
func (section *ipAddressSectionInternal) bitwiseOr(msk *IPAddressSection, retainPrefix bool) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	if err := section.checkMaskSegmentCount(msk); err != nil {
		return nil, err
	}
	var prefLen PrefixLen
	if retainPrefix {
		prefLen = section.getPrefixLen()
	}
	return section.getOredSegments(
		prefLen,
		true,
		section.getDivision,
		func(i int) SegInt { return msk.GetSegment(i).GetSegmentValue() })
}

func (section *ipAddressSectionInternal) matchesWithMask(other *IPAddressSection, mask *IPAddressSection) bool {
	if err := section.checkMaskSegmentCount(other); err != nil {
		return false
	} else if err := section.checkMaskSegmentCount(mask); err != nil {
		return false
	}
	divCount := section.GetSegmentCount()
	for i := 0; i < divCount; i++ {
		seg := section.GetSegment(i)
		maskSegment := mask.GetSegment(i)
		otherSegment := other.GetSegment(i)
		if !seg.MatchesValsWithMask(
			otherSegment.getSegmentValue(),
			otherSegment.getUpperSegmentValue(),
			maskSegment.getSegmentValue()) {
			return false
		}
	}
	return true
}

func (section *ipAddressSectionInternal) intersect(other *IPAddressSection) (res *IPAddressSection, err addrerr.SizeMismatchError) {
	//check if they are comparable section.  We only check segment count, we do not care about start index.
	err = section.checkSegmentCount(other)
	if err != nil {
		return
	}
	//larger prefix length should prevail?    hmmmmm... I would say that is true, choose the larger prefix
	pref := section.getNetworkPrefixLen()
	otherPref := other.getNetworkPrefixLen()
	if pref != nil {
		if otherPref != nil {
			if otherPref.bitCount() > pref.bitCount() {
				pref = otherPref
			}
		} else {
			pref = nil
		}
	}

	if other.Contains(section.toIPAddressSection()) {
		if pref.Equal(section.getNetworkPrefixLen()) {
			res = section.toIPAddressSection()
			return
		}
	} else if !section.isMultiple() {
		// no intersection, for single valued section, any intersection would have to be containment
		return
	}
	if section.contains(other) {
		if pref.Equal(other.getNetworkPrefixLen()) {
			res = other.toIPAddressSection()
			return
		}
	} else if !other.isMultiple() {
		// no intersection, for single valued section, any intersection would have to be containment
		return
	}

	segCount := section.GetSegmentCount()
	for i := 0; i < segCount; i++ {
		seg := section.GetSegment(i)
		otherSeg := other.GetSegment(i)
		lower := seg.GetSegmentValue()
		higher := seg.getUpperSegmentValue()
		otherLower := otherSeg.GetSegmentValue()
		otherHigher := otherSeg.getUpperSegmentValue()
		if otherLower > higher || lower > otherHigher {
			//no overlap in this segment means no overlap at all
			return
		}
	}

	// all segments have overlap
	segs := createSegmentArray(segCount)
	for i := 0; i < segCount; i++ {
		seg := section.GetSegment(i)
		otherSeg := other.GetSegment(i)
		segPref := getSegmentPrefixLength(seg.getBitCount(), pref, i)
		if seg.Contains(otherSeg) {
			if segPref.Equal(otherSeg.GetSegmentPrefixLen()) {
				segs[i] = otherSeg.ToDiv()
				continue
			}
		}
		if otherSeg.Contains(seg) {
			if segPref.Equal(seg.GetSegmentPrefixLen()) {
				segs[i] = seg.ToDiv()
				continue
			}
		}
		lower := seg.GetSegmentValue()
		higher := seg.getUpperSegmentValue()
		otherLower := otherSeg.GetSegmentValue()
		otherHigher := otherSeg.getUpperSegmentValue()
		lower = max(lower, otherLower)
		higher = min(higher, otherHigher)
		segs[i] = createAddressDivision(seg.deriveNewMultiSeg(lower, higher, segPref))
	}
	res = deriveIPAddressSectionPrefLen(section.toIPAddressSection(), segs, pref)
	return
}

func (section *ipAddressSectionInternal) subtract(other *IPAddressSection) (res []*IPAddressSection, err addrerr.SizeMismatchError) {
	//check if they are comparable section
	err = section.checkSegmentCount(other)
	if err != nil {
		return
	}
	//Since this is only called from IPv4 and IPv6, we need not check section versions or types here
	if !section.isMultiple() {
		if other.Contains(section.toIPAddressSection()) {
			return
		}
		res = []*IPAddressSection{section.toIPAddressSection()}
		return
	}
	//getDifference: same as removing the intersection
	//   section you confirm there is an intersection in each segment.
	// Then you remove each intersection, one at a time, leaving the other segments the same, since only one segment needs to differ.
	// To prevent adding the same section twice, use only the intersection (ie the relative complement of the diff)
	// of segments already handled and not the whole segment.

	// For example: 0-3.0-3.2.4 subtracting 1-4.1-3.2.4, the intersection is 1-3.1-3.2.4
	// The diff of the section segment is just 0, giving 0.0-3.2.4 (subtract the section segment, leave the others the same)
	// The diff of the second segment is also 0, but for the section segment we use the intersection since we handled the section already, giving 1-3.0.2.4
	// 	(take the intersection of the section segment, subtract the second segment, leave remaining segments the same)

	segCount := section.GetSegmentCount()
	for i := 0; i < segCount; i++ {
		seg := section.GetSegment(i)
		otherSeg := other.GetSegment(i)
		lower := seg.GetSegmentValue()
		higher := seg.getUpperSegmentValue()
		otherLower := otherSeg.GetSegmentValue()
		otherHigher := otherSeg.getUpperSegmentValue()
		if otherLower > higher || lower > otherHigher {
			//no overlap in this segment means no overlap at all
			res = []*IPAddressSection{section.toIPAddressSection()}
			return
		}
	}

	// As we create each section, the initial segments created by us will have no prefix,
	// the trailing segments come from the original section, so the resulting section will have prefix-consistent segments.
	// We do not care what that prefix length is, because at the end of this method we assign a new prefix.
	// We just need to be sure that we create a valid section, one with prefix-consistent segments.
	intersections := createSegmentArray(segCount)
	sections := make([]*IPAddressSection, 0, segCount<<1)
	for i := 0; i < segCount; i++ {
		seg := section.GetSegment(i)
		otherSeg := other.GetSegment(i)
		lower := seg.GetSegmentValue()
		higher := seg.getUpperSegmentValue()
		otherLower := otherSeg.GetSegmentValue()
		otherHigher := otherSeg.getUpperSegmentValue()
		if lower >= otherLower {
			if higher <= otherHigher {
				//this segment is contained in the other
				if seg.isPrefixed() {
					intersections[i] = createAddressDivision(seg.deriveNewMultiSeg(lower, higher, nil))
				} else {
					intersections[i] = seg.ToDiv()
				}
				continue
			}
			//otherLower <= lower <= otherHigher < higher
			intersections[i] = createAddressDivision(seg.deriveNewMultiSeg(lower, otherHigher, nil))
			section := section.createDiffSection(seg, otherHigher+1, higher, i, intersections)
			sections = append(sections, section)
		} else {
			//lower < otherLower <= otherHigher
			section := section.createDiffSection(seg, lower, otherLower-1, i, intersections)
			sections = append(sections, section)
			if higher <= otherHigher {
				intersections[i] = createAddressDivision(seg.deriveNewMultiSeg(otherLower, higher, nil))
			} else {
				//lower < otherLower <= otherHigher < higher
				intersections[i] = createAddressDivision(seg.deriveNewMultiSeg(otherLower, otherHigher, nil))
				section = section.createDiffSection(seg, otherHigher+1, higher, i, intersections)
				sections = append(sections, section)
			}
		}
	}
	if len(sections) == 0 {
		return
	}

	//apply the prefix to the sections
	//for each section, we figure out what each prefix length should be
	if section.isPrefixed() {
		thisPrefix := section.getNetworkPrefixLen().bitCount()
		for i := 0; i < len(sections); i++ {
			section := sections[i]
			bitCount := section.GetBitCount()
			totalPrefix := bitCount
			for j := section.GetSegmentCount() - 1; j >= 0; j-- {
				seg := section.GetSegment(j)
				segBitCount := seg.GetBitCount()
				segPrefix := seg.GetMinPrefixLenForBlock()
				if segPrefix == segBitCount {
					break
				} else {
					totalPrefix -= segBitCount
					if segPrefix != 0 {
						totalPrefix += segPrefix
						break
					}
				}
			}
			if totalPrefix != bitCount {
				if totalPrefix < thisPrefix {
					totalPrefix = thisPrefix
				}
				section = section.SetPrefixLen(totalPrefix)
				sections[i] = section
			}
		}
	}
	res = sections
	return
}

func (section *ipAddressSectionInternal) createDiffSection(
	seg *IPAddressSegment,
	lower,
	upper SegInt,
	diffIndex int,
	intersectingValues []*AddressDivision) *IPAddressSection {
	segCount := section.GetSegmentCount()
	segments := createSegmentArray(segCount)
	for j := 0; j < diffIndex; j++ {
		segments[j] = intersectingValues[j]
	}
	diff := createAddressDivision(seg.deriveNewMultiSeg(lower, upper, nil))
	segments[diffIndex] = diff
	for j := diffIndex + 1; j < segCount; j++ {
		segments[j] = section.getDivision(j)
	}
	return deriveIPAddressSection(section.toIPAddressSection(), segments)
}

func (section *ipAddressSectionInternal) spanWithPrefixBlocks() []ExtendedIPSegmentSeries {
	if section.IsSequential() {
		if section.IsSinglePrefixBlock() {
			wrapped := wrapIPSection(section.toIPAddressSection())
			return []ExtendedIPSegmentSeries{wrapped}
		}
		return cloneIPSections(nil, getSpanningPrefixBlocks(section.toIPAddressSection(), section.toIPAddressSection()))
	}
	return cloneIPSections(nil, spanWithPrefixBlocks(section.toIPAddressSection()))
}

func (section *ipAddressSectionInternal) spanWithSequentialBlocks() []ExtendedIPSegmentSeries {
	if section.IsSequential() {
		wrapped := wrapIPSection(section.toIPAddressSection())
		return []ExtendedIPSegmentSeries{wrapped}
	}
	return cloneIPSections(nil, spanWithSequentialBlocks(section.toIPAddressSection()))
}

func (section *ipAddressSectionInternal) coverSeriesWithPrefixBlock() ExtendedIPSegmentSeries {
	if section.IsSinglePrefixBlock() {
		return wrapIPSection(section.toIPAddressSection())
	}
	return coverWithPrefixBlock(
		section.getLower().ToIP(),
		section.getUpper().ToIP()).Wrap()
}

func (section *ipAddressSectionInternal) coverWithPrefixBlock() *IPAddressSection {
	if section.IsSinglePrefixBlock() {
		return section.toIPAddressSection()
	}
	return coverWithPrefixBlock(section.getLower().ToIP(), section.getUpper().ToIP())
}

func (section *ipAddressSectionInternal) coverWithPrefixBlockTo(other *IPAddressSection) (*IPAddressSection, addrerr.SizeMismatchError) {
	if err := section.checkSegmentCount(other); err != nil {
		return nil, err
	}
	return getCoveringPrefixBlock(
		section.toIPAddressSection(),
		other), nil
}

func (section *ipAddressSectionInternal) getNetworkSection() *IPAddressSection {
	var prefLen BitCount
	if section.isPrefixed() {
		prefLen = section.getPrefixLen().bitCount()
	} else {
		prefLen = section.GetBitCount()
	}
	return section.getNetworkSectionLen(prefLen)
}

func (section *ipAddressSectionInternal) getNetworkSectionLen(networkPrefixLength BitCount) *IPAddressSection {
	segmentCount := section.GetSegmentCount()
	if segmentCount == 0 {
		return section.toIPAddressSection()
	}
	networkPrefixLength = checkBitCount(networkPrefixLength, section.GetBitCount())
	bitsPerSegment := section.GetBitsPerSegment()
	prefixedSegmentIndex := getNetworkSegmentIndex(networkPrefixLength, section.GetBytesPerSegment(), bitsPerSegment)
	var newSegments []*AddressDivision
	if prefixedSegmentIndex >= 0 {
		segPrefLength := getPrefixedSegmentPrefixLength(bitsPerSegment, networkPrefixLength, prefixedSegmentIndex) // prefixedSegmentIndex of -1 already handled
		lastSeg := section.GetSegment(prefixedSegmentIndex)
		prefBits := segPrefLength.bitCount()
		mask := ^SegInt(0) << uint(bitsPerSegment-prefBits)
		lower, upper := lastSeg.getSegmentValue()&mask, lastSeg.getUpperSegmentValue()|^mask
		networkSegmentCount := prefixedSegmentIndex + 1
		if networkSegmentCount == segmentCount && segsSame(segPrefLength, lastSeg.GetSegmentPrefixLen(), lower, lastSeg.getSegmentValue(), upper, lastSeg.getUpperSegmentValue()) {
			// the segment count and prefixed segment matches
			return section.toIPAddressSection()
		}
		newSegments = createSegmentArray(networkSegmentCount)
		section.copySubDivisions(0, prefixedSegmentIndex, newSegments)
		newSegments[prefixedSegmentIndex] = createAddressDivision(lastSeg.deriveNewMultiSeg(lower, upper, segPrefLength))
	} else {
		newSegments = createSegmentArray(0)
	}
	return deriveIPAddressSectionPrefLen(section.toIPAddressSection(), newSegments, cacheBitCount(networkPrefixLength))
}

func (section *ipAddressSectionInternal) getHostSection() *IPAddressSection {
	var prefLen BitCount
	if section.isPrefixed() {
		prefLen = section.getPrefixLen().bitCount()
	}
	return section.getHostSectionLen(prefLen)
}

func (section *ipAddressSectionInternal) getHostSectionLen(networkPrefixLength BitCount) *IPAddressSection {
	segmentCount := section.GetSegmentCount()
	if segmentCount == 0 {
		return section.toIPAddressSection()
	}
	networkPrefixLength = checkBitCount(networkPrefixLength, section.GetBitCount())
	bitsPerSegment := section.GetBitsPerSegment()
	bytesPerSegment := section.GetBytesPerSegment()
	prefixedSegmentIndex := getHostSegmentIndex(networkPrefixLength, bytesPerSegment, bitsPerSegment)
	var prefLen PrefixLen
	var newSegments []*AddressDivision
	if prefixedSegmentIndex < segmentCount {
		firstSeg := section.GetSegment(prefixedSegmentIndex)
		segPrefLength := getPrefixedSegmentPrefixLength(bitsPerSegment, networkPrefixLength, prefixedSegmentIndex)
		prefLen = segPrefLength
		prefBits := segPrefLength.bitCount()
		//mask the boundary segment
		mask := ^(^SegInt(0) << uint(bitsPerSegment-prefBits))
		divLower := uint64(firstSeg.getDivisionValue())
		divUpper := uint64(firstSeg.getUpperDivisionValue())
		divMask := uint64(mask)
		maxVal := uint64(^SegInt(0))
		masker := MaskRange(divLower, divUpper, divMask, maxVal)
		lower, upper := masker.GetMaskedLower(divLower, divMask), masker.GetMaskedUpper(divUpper, divMask)
		segLower, segUpper := SegInt(lower), SegInt(upper)
		if prefixedSegmentIndex == 0 && segsSame(segPrefLength, firstSeg.GetSegmentPrefixLen(), segLower, firstSeg.getSegmentValue(), segUpper, firstSeg.getUpperSegmentValue()) {
			// the segment count and prefixed segment matches
			return section.toIPAddressSection()
		}
		hostSegmentCount := segmentCount - prefixedSegmentIndex
		newSegments = createSegmentArray(hostSegmentCount)
		newSegments[0] = createAddressDivision(firstSeg.deriveNewMultiSeg(segLower, segUpper, segPrefLength))

		// the remaining segments each must have zero-segment prefix length
		var zeroPrefixIndex int
		if section.isPrefixed() {
			zeroPrefixIndex = getNetworkSegmentIndex(section.GetPrefixLen().bitCount(), bytesPerSegment, bitsPerSegment) + 1
		} else {
			zeroPrefixIndex = segmentCount
		}
		zeroPrefixIndex -= prefixedSegmentIndex
		zeroPrefixIndex = max(zeroPrefixIndex, 1)
		for i := 1; i < zeroPrefixIndex; i++ {
			seg := section.GetSegment(prefixedSegmentIndex + i)
			newSegments[i] = createAddressDivision(seg.derivePrefixed(cacheBitCount(0)))
		}
		// the rest already have zero-segment prefix length, just copy them
		section.copySubDivisions(prefixedSegmentIndex+zeroPrefixIndex, prefixedSegmentIndex+hostSegmentCount, newSegments[zeroPrefixIndex:])
	} else {
		prefLen = cacheBitCount(0)
		newSegments = createSegmentArray(0)
	}
	return deriveIPAddressSectionPrefLen(section.toIPAddressSection(), newSegments, prefLen)
}

func (section *ipAddressSectionInternal) getSubnetSegments( // called by methods to adjust/remove/set prefix length, masking methods, zero host and zero network methods
	startIndex int,
	networkPrefixLength PrefixLen,
	verifyMask bool,
	segProducer func(int) *AddressDivision,
	segmentMaskProducer func(int) SegInt,
) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	newSect, err := section.addressSectionInternal.getSubnetSegments(startIndex, networkPrefixLength, verifyMask, segProducer, segmentMaskProducer)
	return newSect.ToIP(), err
}

func (section *ipAddressSectionInternal) getOredSegments(
	networkPrefixLength PrefixLen,
	verifyMask bool,
	segProducer func(int) *AddressDivision,
	segmentMaskProducer func(int) SegInt) (res *IPAddressSection, err addrerr.IncompatibleAddressError) {
	networkPrefixLength = checkPrefLen(networkPrefixLength, section.GetBitCount())
	bitsPerSegment := section.GetBitsPerSegment()
	count := section.GetSegmentCount()
	for i := 0; i < count; i++ {
		segmentPrefixLength := getSegmentPrefixLength(bitsPerSegment, networkPrefixLength, i)
		seg := segProducer(i)
		//note that the mask can represent a range (for example a CIDR mask),
		//but we use the lowest value (maskSegment.value) in the range when masking (ie we discard the range)
		maskValue := segmentMaskProducer(i)
		origValue, origUpperValue := seg.getSegmentValue(), seg.getUpperSegmentValue()
		value, upperValue := origValue, origUpperValue
		if verifyMask {
			mask64 := uint64(maskValue)
			val64 := uint64(value)
			upperVal64 := uint64(upperValue)
			masker := bitwiseOrRange(val64, upperVal64, mask64, seg.GetMaxValue())
			if !masker.IsSequential() {
				err = &incompatibleAddressError{addressError{key: "ipaddress.error.maskMismatch"}}
				return
			}
			value = SegInt(masker.GetOredLower(val64, mask64))
			upperValue = SegInt(masker.GetOredUpper(upperVal64, mask64))
		} else {
			value |= maskValue
			upperValue |= maskValue
		}
		if !segsSame(segmentPrefixLength, seg.getDivisionPrefixLength(), value, origValue, upperValue, origUpperValue) {
			newSegments := createSegmentArray(count)
			section.copySubDivisions(0, i, newSegments)
			newSegments[i] = createAddressDivision(seg.deriveNewMultiSeg(value, upperValue, segmentPrefixLength))
			for i++; i < count; i++ {
				segmentPrefixLength = getSegmentPrefixLength(bitsPerSegment, networkPrefixLength, i)
				seg = segProducer(i)
				maskValue = segmentMaskProducer(i)
				origValue, origUpperValue = seg.getSegmentValue(), seg.getUpperSegmentValue()
				value, upperValue = origValue, origUpperValue
				if verifyMask {
					mask64 := uint64(maskValue)
					val64 := uint64(value)
					upperVal64 := uint64(upperValue)
					masker := bitwiseOrRange(val64, upperVal64, mask64, seg.GetMaxValue())
					if !masker.IsSequential() {
						err = &incompatibleAddressError{addressError{key: "ipaddress.error.maskMismatch"}}
						return
					}
					value = SegInt(masker.GetOredLower(val64, mask64))
					upperValue = SegInt(masker.GetOredUpper(upperVal64, mask64))
				} else {
					value |= maskValue
					upperValue |= maskValue
				}
				if !segsSame(segmentPrefixLength, seg.getDivisionPrefixLength(), value, origValue, upperValue, origUpperValue) {
					newSegments[i] = createAddressDivision(seg.deriveNewMultiSeg(value, upperValue, segmentPrefixLength))
				} else {
					newSegments[i] = seg
				}
			}
			res = deriveIPAddressSectionPrefLen(section.toIPAddressSection(), newSegments, networkPrefixLength)
			return
		}
	}
	res = section.toIPAddressSection()
	return
}

func (section *ipAddressSectionInternal) getNetwork() IPAddressNetwork {
	if addrType := section.getAddrType(); addrType.isIPv4() {
		return ipv4Network
	} else if addrType.isIPv6() {
		return ipv6Network
	}
	return nil
}

func (section *ipAddressSectionInternal) getNetworkMask(network IPAddressNetwork) *IPAddressSection {
	var prefLen BitCount
	if section.isPrefixed() {
		prefLen = section.getNetworkPrefixLen().bitCount()
	} else {
		prefLen = section.GetBitCount()
	}
	return network.GetNetworkMask(prefLen).GetSubSection(0, section.GetSegmentCount())
}

func (section *ipAddressSectionInternal) getHostMask(network IPAddressNetwork) *IPAddressSection {
	var prefLen BitCount
	if section.isPrefixed() {
		prefLen = section.getNetworkPrefixLen().bitCount()
	}
	return network.GetHostMask(prefLen).GetSubSection(0, section.GetSegmentCount())
}

func (section *ipAddressSectionInternal) insert(index int, other *IPAddressSection, segmentToBitsShift uint) *IPAddressSection {
	return section.replaceLen(index, index, other, 0, other.GetSegmentCount(), segmentToBitsShift)
}

// Replaces segments starting from startIndex and ending before endIndex with the segments starting at replacementStartIndex and
// ending before replacementEndIndex from the replacement section.
func (section *ipAddressSectionInternal) replaceLen(
	startIndex, endIndex int, replacement *IPAddressSection, replacementStartIndex, replacementEndIndex int, segmentToBitsShift uint) *IPAddressSection {

	segmentCount := section.GetSegmentCount()
	startIndex, endIndex, replacementStartIndex, replacementEndIndex =
		adjustIndices(startIndex, endIndex, segmentCount, replacementStartIndex, replacementEndIndex, replacement.GetSegmentCount())
	replacedCount := endIndex - startIndex
	replacementCount := replacementEndIndex - replacementStartIndex
	thizz := section.toAddressSection()
	if replacementCount == 0 && replacedCount == 0 { //keep in mind for ipvx, empty sections cannot have prefix lengths
		return section.toIPAddressSection()
	} else if segmentCount == replacedCount { //keep in mind for ipvx, empty sections cannot have prefix lengths
		return replacement
	}
	var newPrefixLen PrefixLen
	prefixLength := section.getPrefixLen()
	startBits := BitCount(startIndex << segmentToBitsShift)
	if prefixLength != nil && prefixLength.bitCount() <= startBits {
		newPrefixLen = prefixLength
		replacement = replacement.SetPrefixLen(0)
	} else {
		replacementEndBits := BitCount(replacementEndIndex << segmentToBitsShift)
		replacementPrefLen := replacement.getPrefixLen()
		endIndexBits := BitCount(endIndex << segmentToBitsShift)
		if replacementPrefLen != nil && replacementPrefLen.bitCount() <= replacementEndBits {
			var replacementPrefixLen BitCount
			replacementStartBits := BitCount(replacementStartIndex << segmentToBitsShift)
			replacementPrefLenIsZero := replacementPrefLen.bitCount() <= replacementStartBits
			if !replacementPrefLenIsZero {
				replacementPrefixLen = replacementPrefLen.bitCount() - replacementStartBits
			}
			newPrefixLen = cacheBitCount(startBits + replacementPrefixLen)
			if endIndex < segmentCount && (prefixLength == nil || prefixLength.bitCount() > endIndexBits) {
				if replacedCount > 0 || replacementPrefLenIsZero {
					thizz = section.setPrefixLen(endIndexBits)
				} else {
					// this covers the case of a:5:6:7:8 is getting b:c:d/47 at index 1 to 1
					// We need "a" to have no prefix, and "5" to get prefix len 0
					// But setting "5" to have prefix len 0 gives "a" the prefix len 16
					// This is not a problem if any segments are getting replaced or the replacement segments have prefix length 0
					//
					// we move the non-replaced host segments from the end of this to the end of the replacement segments
					// and we also remove the prefix length from this
					additionalSegs := segmentCount - endIndex
					thizz = section.getSubSection(0, startIndex)
					replacement = replacement.insert(
						replacementEndIndex, section.getSubSection(endIndex, segmentCount).ToIP(), segmentToBitsShift)
					replacementEndIndex += additionalSegs
				}
			}
		} else if prefixLength != nil {
			replacementBits := BitCount(replacementCount << segmentToBitsShift)
			var endPrefixBits BitCount
			if prefixLength.bitCount() > endIndexBits {
				endPrefixBits = prefixLength.bitCount() - endIndexBits
			}
			newPrefixLen = cacheBitCount(startBits + replacementBits + endPrefixBits)
		} // else newPrefixLen is nil
	}
	return thizz.replace(startIndex, endIndex, replacement.ToSectionBase(),
		replacementStartIndex, replacementEndIndex, newPrefixLen).ToIP()
}

func (section *ipAddressSectionInternal) toNormalizedWildcardString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToNormalizedWildcardString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToNormalizedWildcardString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toCanonicalWildcardString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToCanonicalWildcardString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToCanonicalWildcardString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toSegmentedBinaryString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToSegmentedBinaryString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToSegmentedBinaryString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toSQLWildcardString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToSQLWildcardString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToSQLWildcardString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toFullString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToFullString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToFullString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toReverseDNSString() (string, addrerr.IncompatibleAddressError) {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToReverseDNSString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToReverseDNSString()
	}
	return nilSection(), nil
}

func (section *ipAddressSectionInternal) toPrefixLenString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToPrefixLenString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToPrefixLenString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toSubnetString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToNormalizedWildcardString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToPrefixLenString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toCompressedWildcardString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToCompressedWildcardString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToCompressedWildcardString()
	}
	return nilSection()
}

func (section *ipAddressSectionInternal) toCustomString(stringOptions addrstr.IPStringOptions) string {
	return toNormalizedIPZonedString(stringOptions, section.toIPAddressSection(), NoZone)
}

func (section *ipAddressSectionInternal) toCustomZonedString(stringOptions addrstr.IPStringOptions, zone Zone) string {
	return toNormalizedIPZonedString(stringOptions, section.toIPAddressSection(), zone)
}

// Wrap wraps this IP address section, returning a WrappedIPAddressSection, an implementation of ExtendedIPSegmentSeries,
// which can be used to write code that works with both IP addresses and IP address sections.
// Wrap can be called with a nil receiver, wrapping a nil address section.
func (section *ipAddressSectionInternal) Wrap() WrappedIPAddressSection {
	return wrapIPSection(section.toIPAddressSection())
}

// WrapSection wraps this IP address section, returning a WrappedAddressSection, an implementation of ExtendedSegmentSeries,
// which can be used to write code that works with both addresses and address sections.
// WrapSection can be called with a nil receiver, wrapping a nil address section.
func (section *ipAddressSectionInternal) WrapSection() WrappedAddressSection {
	return wrapSection(section.toAddressSection())
}

func (section *ipAddressSectionInternal) toIPAddressSection() *IPAddressSection {
	return (*IPAddressSection)(unsafe.Pointer(section))
}

//// only needed for godoc / pkgsite

// GetBitCount returns the number of bits in each value comprising this address item.
func (section *ipAddressSectionInternal) GetBitCount() BitCount {
	return section.addressSectionInternal.GetBitCount()
}

// GetByteCount returns the number of bytes required for each value comprising this address item.
func (section *ipAddressSectionInternal) GetByteCount() int {
	return section.addressSectionInternal.GetByteCount()
}

//IPv6v4, Div,  Not needed Addr because of GetGenericSegment
//func (grouping *addressDivisionGroupingBase) GetGenericDivision(index int) DivisionType {
//
//IPv6v4, Div, Not needed Addr
//func (grouping *addressDivisionGroupingBase) GetDivisionCount() int {

// IsZero returns whether this section matches exactly the value of zero.
func (section *ipAddressSectionInternal) IsZero() bool {
	return section.addressSectionInternal.IsZero()
}

// IncludesZero returns whether this section includes the value of zero within its range.
func (section *ipAddressSectionInternal) IncludesZero() bool {
	return section.addressSectionInternal.IncludesZero()
}

// IsMax returns whether this section matches exactly the maximum possible value, the value whose bits are all ones.
func (section *ipAddressSectionInternal) IsMax() bool {
	return section.addressSectionInternal.IsMax()
}

// IncludesMax returns whether this section includes the max value, the value whose bits are all ones, within its range.
func (section *ipAddressSectionInternal) IncludesMax() bool {
	return section.addressSectionInternal.IncludesMax()
}

// IsFullRange returns whether this address item represents all possible values attainable by an address item of this type.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (section *ipAddressSectionInternal) IsFullRange() bool {
	return section.addressSectionInternal.IsFullRange()
}

// GetSequentialBlockIndex gets the minimal segment index for which all following segments are full-range blocks.
//
// The segment at this index is not a full-range block itself, unless all segments are full-range.
// The segment at this index and all following segments form a sequential range.
// For the full address section to be sequential, the preceding segments must be single-valued.
func (section *ipAddressSectionInternal) GetSequentialBlockIndex() int {
	return section.addressSectionInternal.GetSequentialBlockIndex()
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential address sections that comprise this address section.
func (section *ipAddressSectionInternal) GetSequentialBlockCount() *big.Int {
	return section.addressSectionInternal.GetSequentialBlockCount()
}

// ContainsPrefixBlock returns whether the values of this item contains the block of values for the given prefix length.
//
// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values in this item for the given prefix length makes no difference.
//
// Use GetMinPrefixLenForBlock to determine the smallest prefix length for which this method returns true.
func (section *ipAddressSectionInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	return section.addressSectionInternal.ContainsPrefixBlock(prefixLen)
}

// ContainsSinglePrefixBlock returns whether the values of this section contains a single prefix block for the given prefix length.
//
// This means there is only one prefix of the given length in this item, and this item contains the prefix block for that given prefix.
//
// Use GetPrefixLenForSingleBlock to determine whether there is a prefix length for which this method returns true.
func (section *ipAddressSectionInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	return section.addressSectionInternal.ContainsSinglePrefixBlock(prefixLen)
}

// IsPrefixBlock returns whether this address segment series has a prefix length and includes the block associated with its prefix length.
// If the prefix length matches the bit count, this returns true.
//
// This is different from ContainsPrefixBlock in that this method returns
// false if the series has no prefix length, or a prefix length that differs from a prefix length for which ContainsPrefixBlock returns true.
func (section *ipAddressSectionInternal) IsPrefixBlock() bool {
	return section.addressSectionInternal.IsPrefixBlock()
}

// IsSinglePrefixBlock returns whether the range matches the block of values for a single prefix identified by the prefix length of this address.
// This is similar to IsPrefixBlock except that it returns false when the subnet has multiple prefixes.
//
// What distinguishes this method from ContainsSinglePrefixBlock is that this method returns
// false if the series does not have a prefix length assigned to it,
// or a prefix length that differs from the prefix length for which ContainsSinglePrefixBlock returns true.
//
// It is similar to IsPrefixBlock but returns false when there are multiple prefixes.
func (section *ipAddressSectionInternal) IsSinglePrefixBlock() bool {
	return section.addressSectionInternal.IsSinglePrefixBlock()
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this section includes the block of all values for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this section represents a single value, this returns the bit count.
func (section *ipAddressSectionInternal) GetMinPrefixLenForBlock() BitCount {
	return section.addressSectionInternal.GetMinPrefixLenForBlock()
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this address section matches the block of addresses for that prefix.
//
// If no such prefix exists, GetPrefixLenForSingleBlock returns nil.
//
// If this address section represents a single value, returns the bit length.
func (section *ipAddressSectionInternal) GetPrefixLenForSingleBlock() PrefixLen {
	return section.addressSectionInternal.GetPrefixLenForSingleBlock()
}

// GetValue returns the lowest individual address section in this address section as an integer value.
func (section *ipAddressSectionInternal) GetValue() *big.Int {
	return section.addressSectionInternal.GetValue()
}

// GetUpperValue returns the highest individual address section in this address section as an integer value.
func (section *ipAddressSectionInternal) GetUpperValue() *big.Int {
	return section.addressSectionInternal.GetUpperValue()
}

// Bytes returns the lowest individual address section in this address section as a byte slice.
func (section *ipAddressSectionInternal) Bytes() []byte {
	return section.addressSectionInternal.Bytes()
}

// UpperBytes returns the highest individual address section in this address section as a byte slice.
func (section *ipAddressSectionInternal) UpperBytes() []byte {
	return section.addressSectionInternal.UpperBytes()
}

// CopyBytes copies the value of the lowest individual address section in the section into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (section *ipAddressSectionInternal) CopyBytes(bytes []byte) []byte {
	return section.addressSectionInternal.CopyBytes(bytes)
}

// CopyUpperBytes copies the value of the highest individual address section in the section into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (section *ipAddressSectionInternal) CopyUpperBytes(bytes []byte) []byte {
	return section.addressSectionInternal.CopyUpperBytes(bytes)
}

// IsSequential returns  whether the section represents a range of values that are sequential.
//
// Generally, this means that any segment covering a range of values must be followed by segment that are full range, covering all values.
func (section *ipAddressSectionInternal) IsSequential() bool {
	return section.addressSectionInternal.IsSequential()
}

// GetBitsPerSegment returns the number of bits comprising each segment in this section.  Segments in the same address section are equal length.
func (section *ipAddressSectionInternal) GetBitsPerSegment() BitCount {
	return section.addressSectionInternal.GetBitsPerSegment()
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this section.  Segments in the same address section are equal length.
func (section *ipAddressSectionInternal) GetBytesPerSegment() int {
	return section.addressSectionInternal.GetBytesPerSegment()
}

// GetGenericSegment returns the segment at the given index as an AddressSegmentType.
// The first segment is at index 0.
// GetGenericSegment will panic given a negative index or an index matching or larger than the segment count.
func (section *ipAddressSectionInternal) GetGenericSegment(index int) AddressSegmentType {
	return section.addressSectionInternal.GetGenericSegment(index)
}

// GetSegmentCount returns the segment/division count.
func (section *ipAddressSectionInternal) GetSegmentCount() int {
	return section.addressSectionInternal.GetSegmentCount()
}

// GetMaxSegmentValue returns the maximum possible segment value for this type of address.
//
// Note this is not the maximum of the range of segment values in this specific address,
// this is the maximum value of any segment for this address type and version, determined by the number of bits per segment.
func (section *ipAddressSectionInternal) GetMaxSegmentValue() SegInt {
	return section.addressSectionInternal.GetMaxSegmentValue()
}

// TestBit returns true if the bit in the lower value of this section at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this section.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (section *ipAddressSectionInternal) TestBit(n BitCount) bool {
	return section.addressSectionInternal.TestBit(n)
}

// IsOneBit returns true if the bit in the lower value of this section at the given index is 1, where index 0 refers to the most significant bit.
// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this item.
func (section *ipAddressSectionInternal) IsOneBit(prefixBitIndex BitCount) bool {
	return section.addressSectionInternal.IsOneBit(prefixBitIndex)
}

// PrefixEqual determines if the given section matches this section up to the prefix length of this section.
// It returns whether the argument section has the same address section prefix values as this.
//
// All prefix bits of this section must be present in the other section to be comparable, otherwise false is returned.
func (section *ipAddressSectionInternal) PrefixEqual(other AddressSectionType) bool {
	return section.addressSectionInternal.PrefixEqual(other)
}

// PrefixContains returns whether the prefix values in the given address section
// are prefix values in this address section, using the prefix length of this section.
// If this address section has no prefix length, the entire address is compared.
//
// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
//
// All prefix bits of this section must be present in the other section to be comparable.
func (section *ipAddressSectionInternal) PrefixContains(other AddressSectionType) bool {
	return section.addressSectionInternal.PrefixContains(other)
}

//// end needed for godoc / pkgsite

// An IPAddressSection is an address section of an IP address, containing a certain number of consecutive segments of an IP address.
//
// It is a series of individual address segments.  Each segment has equal bit-length.  Each address is backed by an address section that contains all the segments of the address.
//
// IPAddressSection objects are immutable.  This also makes them concurrency-safe.
//
// Most operations that can be performed on IPAddress instances can also be performed on IPAddressSection instances and vice-versa.
type IPAddressSection struct {
	ipAddressSectionInternal
}

// containsSame returns whether this address section contains all address sections in the given address section collection of the same type.
func (addr *IPAddressSection) containsSame(other *IPAddressSection) bool {
	return addr.Contains(other)
}

// Contains returns whether this is same type and version as the given address section and whether it contains all values in the given section.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *IPAddressSection) Contains(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.contains(other)
}

// Overlaps returns whether this is same type and version as the given address section and whether it overlaps the given section, both sections containing at least one individual section in common.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *IPAddressSection) Overlaps(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.overlaps(other)
}

// Equal returns whether the given address section is equal to this address section.
// Two address sections are equal if they represent the same set of sections.
// They must match:
//   - type/version: IPv4, IPv6
//   - segment counts
//   - segment value ranges
//
// Prefix lengths are ignored.
func (section *IPAddressSection) Equal(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.equal(other)
}

// Compare returns a negative integer, zero, or a positive integer if this address section is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (section *IPAddressSection) Compare(item AddressItem) int {
	return CountComparator.Compare(section, item)
}

// CompareSize compares the counts of two address sections or other items, the number of individual items represented.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether this section represents more individual address sections than another item.
//
// CompareSize returns a positive integer if this address section has a larger count than the item given, zero if they are the same, or a negative integer if the other has a larger count.
func (section *IPAddressSection) CompareSize(other AddressItem) int {
	if section == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return section.compareSize(other)
}

// GetCount returns the count of possible distinct values for this item.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (section *IPAddressSection) GetCount() *big.Int {
	if section == nil {
		return bigZero()
	} else if sect := section.ToIPv4(); sect != nil {
		return sect.GetCount()
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetCount()
	}
	return section.addressDivisionGroupingBase.getCount()
}

// IsMultiple returns whether this section represents multiple values.
func (section *IPAddressSection) IsMultiple() bool {
	return section != nil && section.isMultiple()
}

// IsPrefixed returns whether this section has an associated prefix length.
func (section *IPAddressSection) IsPrefixed() bool {
	return section != nil && section.isPrefixed()
}

// GetPrefixCount returns the number of distinct prefix values in this item.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, it returns the same value as GetCount.
func (section *IPAddressSection) GetPrefixCount() *big.Int {
	if sect := section.ToIPv4(); sect != nil {
		return sect.GetPrefixCount()
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetPrefixCount()
	}
	return section.addressDivisionGroupingBase.GetPrefixCount()
}

// GetPrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
func (section *IPAddressSection) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	if sect := section.ToIPv4(); sect != nil {
		return sect.GetPrefixCountLen(prefixLen)
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetPrefixCountLen(prefixLen)
	}
	return section.addressDivisionGroupingBase.GetPrefixCountLen(prefixLen)
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) segments.
func (section *IPAddressSection) GetBlockCount(segments int) *big.Int {
	if sect := section.ToIPv4(); sect != nil {
		return sect.GetBlockCount(segments)
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetBlockCount(segments)
	}
	return section.addressDivisionGroupingBase.GetBlockCount(segments)
}

// IsAdaptiveZero returns true if the division grouping was originally created as an implicitly zero-valued section or grouping (e.g. IPv4AddressSection{}),
// meaning it was not constructed using a constructor function.
// Such a grouping, which has no divisions or segments, is convertible to an implicitly zero-valued grouping of any type or version, whether IPv6, IPv4, MAC, or other.
// In other words, when a section or grouping is the zero-value, then it is equivalent and convertible to the zero value of any other section or grouping type.
func (section *IPAddressSection) IsAdaptiveZero() bool {
	return section != nil && section.matchesZeroGrouping()
}

// ToDivGrouping converts to an AddressDivisionGrouping, a polymorphic type usable with all address sections and division groupings.
// Afterwards, you can convert back with ToIP.
//
// ToDivGrouping can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPAddressSection) ToDivGrouping() *AddressDivisionGrouping {
	return section.ToSectionBase().ToDivGrouping()
}

// ToSectionBase converts to an AddressSection, a polymorphic type usable with all address sections.
// Afterwards, you can convert back with ToIP.
//
// ToSectionBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPAddressSection) ToSectionBase() *AddressSection {
	return (*AddressSection)(unsafe.Pointer(section))
}

// ToIPv6 converts to an IPv6AddressSection if this section originated as an IPv6 section.
// If not, ToIPv6 returns nil.
//
// ToIPv6 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPAddressSection) ToIPv6() *IPv6AddressSection {
	if section.IsIPv6() {
		return (*IPv6AddressSection)(section)
	}
	return nil
}

// ToIPv4 converts to an IPv4AddressSection if this section originated as an IPv4 section.
// If not, ToIPv4 returns nil.
//
// ToIPv4 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPAddressSection) ToIPv4() *IPv4AddressSection {
	if section.IsIPv4() {
		return (*IPv4AddressSection)(section)
	}
	return nil
}

// IsIPv4 returns true if this address section originated as an IPv4 section.  If so, use ToIPv4 to convert back to the IPv4-specific type.
func (section *IPAddressSection) IsIPv4() bool { // we allow nil receivers to allow this to be called following a failed converion like ToIP()
	return section != nil && section.matchesIPv4SectionType()
}

// IsIPv6 returns true if this address section originated as an IPv6 section.  If so, use ToIPv6 to convert back to the IPv6-specific type.
func (section *IPAddressSection) IsIPv6() bool {
	return section != nil && section.matchesIPv6SectionType()
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (section *IPAddressSection) GetTrailingSection(index int) *IPAddressSection {
	return section.GetSubSection(index, section.GetSegmentCount())
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (section *IPAddressSection) GetSubSection(index, endIndex int) *IPAddressSection {
	return section.getSubSection(index, endIndex).ToIP()
}

// GetNetworkSection returns a subsection containing the segments with the network bits of the address section.
// The returned section will have only as many segments as needed as determined by the existing CIDR network prefix length.
//
// If this series has no CIDR prefix length, the returned network section will
// be the entire series as a prefixed section with prefix length matching the address bit length.
func (section *IPAddressSection) GetNetworkSection() *IPAddressSection {
	return section.getNetworkSection()
}

// GetNetworkSectionLen returns a subsection containing the segments with the network of the address section, the prefix bits according to the given prefix length.
// The returned section will have only as many segments as needed to contain the network.
//
// The new section will be assigned the given prefix length,
// unless the existing prefix length is smaller, in which case the existing prefix length will be retained.
func (section *IPAddressSection) GetNetworkSectionLen(prefLen BitCount) *IPAddressSection {
	return section.getNetworkSectionLen(prefLen)
}

// GetHostSection returns a subsection containing the segments with the host of the address section, the bits beyond the CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
//
// If this series has no prefix length, the returned host section will be the full section.
func (section *IPAddressSection) GetHostSection() *IPAddressSection {
	return section.getHostSection()
}

// GetHostSectionLen returns a subsection containing the segments with the host of the address section, the bits beyond the given CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
// The returned section will have an assigned prefix length indicating the beginning of the host.
func (section *IPAddressSection) GetHostSectionLen(prefLen BitCount) *IPAddressSection {
	return section.getHostSectionLen(prefLen)
}

// GetNetworkMask returns the network mask associated with the CIDR network prefix length of this address section.
// If this section has no prefix length, then the all-ones mask is returned.
func (section *IPAddressSection) GetNetworkMask() *IPAddressSection {
	return section.getNetworkMask(section.getNetwork())
}

// GetHostMask returns the host mask associated with the CIDR network prefix length of this address section.
// If this section has no prefix length, then the all-ones mask is returned.
func (section *IPAddressSection) GetHostMask() *IPAddressSection {
	return section.getHostMask(section.getNetwork())
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (section *IPAddressSection) CopySubSegments(start, end int, segs []*IPAddressSegment) (count int) {
	start, end, targetStart := adjust1To1StartIndices(start, end, section.GetDivisionCount(), len(segs))
	segs = segs[targetStart:]
	return section.forEachSubDivision(start, end, func(index int, div *AddressDivision) {
		segs[index] = div.ToIP()
	}, len(segs))
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (section *IPAddressSection) CopySegments(segs []*IPAddressSegment) (count int) {
	return section.ForEachSegment(func(index int, seg *IPAddressSegment) (stop bool) {
		if stop = index >= len(segs); !stop {
			segs[index] = seg
		}
		return
	})
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
func (section *IPAddressSection) GetSegments() (res []*IPAddressSegment) {
	res = make([]*IPAddressSegment, section.GetSegmentCount())
	section.CopySegments(res)
	return
}

// GetLower returns the section in the range with the lowest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1.2-3.4.5-6", the section "1.2.4.5" is returned.
func (section *IPAddressSection) GetLower() *IPAddressSection {
	return section.getLower().ToIP()
}

// GetUpper returns the section in the range with the highest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1.2-3.4.5-6", the section "1.3.4.6" is returned.
func (section *IPAddressSection) GetUpper() *IPAddressSection {
	return section.getUpper().ToIP()
}

// ToZeroHost converts the address section to one in which all individual address sections have a host of zero,
// the host being the bits following the prefix length.
// If the address section has no prefix length, then it returns an all-zero address section.
//
// The returned section will have the same prefix and prefix length.
//
// This returns an error if the section is a range of address sections which cannot be converted to a range in which all sections have zero hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPAddressSection) ToZeroHost() (res *IPAddressSection, err addrerr.IncompatibleAddressError) {
	return section.toZeroHost(false)
}

// ToZeroHostLen converts the address section to one in which all individual sections have a host of zero,
// the host being the bits following the given prefix length.
// If this address section has the same prefix length, then the returned one will too, otherwise the returned section will have no prefix length.
//
// This returns an error if the section is a range of which cannot be converted to a range in which all sections have zero hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPAddressSection) ToZeroHostLen(prefixLength BitCount) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	return section.ToZeroHostLen(prefixLength)
}

// ToZeroNetwork converts the address section to one in which all individual address sections have a network of zero,
// the network being the bits within the prefix length.
// If the address section has no prefix length, then it returns an all-zero address section.
//
// The returned address section will have the same prefix length.
func (section *IPAddressSection) ToZeroNetwork() *IPAddressSection {
	return section.toZeroNetwork()
}

// ToMaxHost converts the address section to one in which all individual address sections have a host of all one-bits, the max value,
// the host being the bits following the prefix length.
// If the address section has no prefix length, then it returns an all-ones section, the max address section.
//
// The returned address section will have the same prefix and prefix length.
//
// This returns an error if the address section is a range of address sections which cannot be converted to a range in which all sections have max hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPAddressSection) ToMaxHost() (res *IPAddressSection, err addrerr.IncompatibleAddressError) {
	return section.toMaxHost()
}

// ToMaxHostLen converts the address section to one in which all individual address sections have a host of all one-bits, the max host,
// the host being the bits following the given prefix length.
// If this section has the same prefix length, then the resulting section will too, otherwise the resulting section will have no prefix length.
//
// This returns an error if the section is a range of address sections which cannot be converted to a range in which all address sections have max hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPAddressSection) ToMaxHostLen(prefixLength BitCount) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	return section.toMaxHostLen(prefixLength)
}

// WithoutPrefixLen provides the same address section but with no prefix length.  The values remain unchanged.
func (section *IPAddressSection) WithoutPrefixLen() *IPAddressSection {
	if !section.IsPrefixed() {
		return section
	}
	return section.withoutPrefixLen()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address section.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (section *IPAddressSection) SetPrefixLen(prefixLen BitCount) *IPAddressSection {
	return section.setPrefixLen(prefixLen).ToIP()
}

// SetPrefixLenZeroed sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address section.
// The provided prefix length will be adjusted to these boundaries if necessary.
//
// If this address section has a prefix length, and the prefix length is increased when setting the new prefix length, the bits moved within the prefix become zero.
// If this address section has a prefix length, and the prefix length is decreased when setting the new prefix length, the bits moved outside the prefix become zero.
//
// In other words, bits that move from one side of the prefix length to the other (bits moved into the prefix or outside the prefix) are zeroed.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (section *IPAddressSection) SetPrefixLenZeroed(prefixLen BitCount) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.setPrefixLenZeroed(prefixLen)
	return res.ToIP(), err
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address section.
//
// If this address section has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (section *IPAddressSection) AdjustPrefixLen(prefixLen BitCount) *IPAddressSection {
	return section.adjustPrefixLen(prefixLen)
}

// AdjustPrefixLenZeroed increases or decreases the prefix length by the given increment while zeroing out the bits that have moved into or outside the prefix.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address section.
//
// If this address section has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
//
// When prefix length is increased, the bits moved within the prefix become zero.
// When a prefix length is decreased, the bits moved outside the prefix become zero.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (section *IPAddressSection) AdjustPrefixLenZeroed(prefixLen BitCount) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	return section.adjustPrefixLenZeroed(prefixLen)
}

// ToPrefixBlock returns the section with the same prefix as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
//
// If this section has no prefix, this section is returned.
func (section *IPAddressSection) ToPrefixBlock() *IPAddressSection {
	return section.toPrefixBlock().ToIP()
}

// ToPrefixBlockLen returns the section with the same prefix of the given length as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
func (section *IPAddressSection) ToPrefixBlockLen(prefLen BitCount) *IPAddressSection {
	return section.toPrefixBlockLen(prefLen).ToIP()
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address section.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address section - it is required that the range of values match the range of a prefix block.
// If there is no such address section, then nil is returned.
func (section *IPAddressSection) AssignPrefixForSingleBlock() *IPAddressSection {
	return section.assignPrefixForSingleBlock().ToIP()
}

// AssignMinPrefixForBlock returns an equivalent address section, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this address section.
//
// In other words, this method assigns a prefix length to this address section matching the largest prefix block in this address section.
func (section *IPAddressSection) AssignMinPrefixForBlock() *IPAddressSection {
	return section.assignMinPrefixForBlock().ToIP()
}

// ToBlock creates a new block of address sections by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (section *IPAddressSection) ToBlock(segmentIndex int, lower, upper SegInt) *IPAddressSection {
	return section.toBlock(segmentIndex, lower, upper).ToIP()
}

// Iterator provides an iterator to iterate through the individual address sections of this address section.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual address sections.
//
// Call IsMultiple to determine if this instance represents multiple address sections, or GetCount for the count.
func (section *IPAddressSection) Iterator() Iterator[*IPAddressSection] {
	if section == nil {
		return ipSectionIterator{nilSectIterator()}
	}
	return ipSectionIterator{section.sectionIterator(nil)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this address section,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this address section.
//
// If the series has no prefix length, then this is equivalent to Iterator.
func (section *IPAddressSection) PrefixIterator() Iterator[*IPAddressSection] {
	return ipSectionIterator{section.prefixIterator(false)}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address section.
// Each iterated address section will be a prefix block with the same prefix length as this address section.
//
// If this address section has no prefix length, then this is equivalent to Iterator.
func (section *IPAddressSection) PrefixBlockIterator() Iterator[*IPAddressSection] {
	return ipSectionIterator{section.prefixIterator(true)}
}

// BlockIterator Iterates through the address sections that can be obtained by iterating through all the upper segments up to the given segment count.
// The segments following remain the same in all iterated sections.
func (section *IPAddressSection) BlockIterator(segmentCount int) Iterator[*IPAddressSection] {
	return ipSectionIterator{section.blockIterator(segmentCount)}
}

// SequentialBlockIterator iterates through the sequential address sections that make up this address section.
//
// Practically, this means finding the count of segments for which the segments that follow are not full range, and then using BlockIterator with that segment count.
//
// Use GetSequentialBlockCount to get the number of iterated elements.
func (section *IPAddressSection) SequentialBlockIterator() Iterator[*IPAddressSection] {
	return ipSectionIterator{section.sequentialBlockIterator()}
}

// IncrementBoundary returns the item that is the given increment from the range boundaries of this item.
//
// If the given increment is positive, adds the value to the highest (GetUpper) in the range to produce a new item.
// If the given increment is negative, adds the value to the lowest (GetLower) in the range to produce a new item.
// If the increment is zero, returns this.
//
// If this represents just a single value, this item is simply incremented by the given increment value, positive or negative.
//
// On overflow or underflow, IncrementBoundary returns nil.
func (section *IPAddressSection) IncrementBoundary(increment int64) *IPAddressSection {
	return section.incrementBoundary(increment).ToIP()
}

// Increment returns the item that is the given increment upwards into the range,
// with the increment of 0 returning the first in the range.
//
// If the increment i matches or exceeds the range count c, then i - c + 1
// is added to the upper item of the range.
// An increment matching the count gives you the item just above the highest in the range.
//
// If the increment is negative, it is added to the lowest of the range.
// To get the item just below the lowest of the range, use the increment -1.
//
// If this represents just a single value, the item is simply incremented by the given increment, positive or negative.
//
// If this item represents multiple values, a positive increment i is equivalent i + 1 values from the iterator and beyond.
// For instance, a increment of 0 is the first value from the iterator, an increment of 1 is the second value from the iterator, and so on.
// An increment of a negative value added to the count is equivalent to the same number of iterator values preceding the last value of the iterator.
// For instance, an increment of count - 1 is the last value from the iterator, an increment of count - 2 is the second last value, and so on.
//
// On overflow or underflow, Increment returns nil.
func (section *IPAddressSection) Increment(increment int64) *IPAddressSection {
	return section.increment(increment).ToIP()
}

// Enumerate indicates where an individual address section sits relative to the address section range ordering.
//
// Determines how many address section elements of a range precede the given address section element, if the address section is in the range.
// If above the range, it is the distance to the upper boundary added to the range count less one, and if below the range, the distance to the lower boundary.
//
// In other words, if the given address section is not in the range but above it, returns the number of address sections preceding the address from the upper range boundary,
// added to one less than the total number of range address sections.  If the given address section is not in the subnet but below it, returns the number of address sections following the address section to the lower subnet boundary.
//
// If the argument is not in the range, but neither above nor below the range, then nil is returned.
//
// Enumerate returns nil when the argument is multi-valued. The argument must be an individual address section.
//
// When this is also an individual address section, the returned value is the distance (difference) between the two address section values.
//
// If the given address section does not have the same version or type, then nil is returned.
//
// Sections must also have the same number of segments to be comparable, otherwise nil is returned.
func (section *IPAddressSection) Enumerate(other AddressSectionType) *big.Int {
	return section.enumerate(other)
}

// SpanWithPrefixBlocks returns an array of prefix blocks that spans the same set of individual address sections as this section.
//
// Unlike SpanWithPrefixBlocksTo, the result only includes blocks that are a part of this section.
func (section *IPAddressSection) SpanWithPrefixBlocks() []*IPAddressSection {
	if section.IsSequential() {
		if section.IsSinglePrefixBlock() {
			return []*IPAddressSection{section}
		}
		return getSpanningPrefixBlocks(section, section)
	}
	return spanWithPrefixBlocks(section)
}

// SpanWithSequentialBlocks produces the smallest slice of sequential blocks that cover the same set of sections as this.
//
// This slice can be shorter than that produced by SpanWithPrefixBlocks and is never longer.
func (section *IPAddressSection) SpanWithSequentialBlocks() []*IPAddressSection {
	if section.IsSequential() {
		return []*IPAddressSection{section}
	}
	return spanWithSequentialBlocks(section)
}

// CoverWithPrefixBlock returns the minimal-size prefix block that covers all the individual address sections in this section.
// The resulting block will have a larger count than this, unless this section is already a prefix block.
func (section *IPAddressSection) CoverWithPrefixBlock() *IPAddressSection {
	return section.coverWithPrefixBlock()
}

// ReverseBits returns a new section with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (section *IPAddressSection) ReverseBits(perByte bool) (*IPAddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.reverseBits(perByte)
	return res.ToIP(), err
}

// ReverseBytes returns a new section with the bytes reversed.  Any prefix length is dropped.
//
// If each segment is more than 1 byte long, and the bytes within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, then this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
func (section *IPAddressSection) ReverseBytes() (*IPAddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.reverseBytes(false)
	return res.ToIP(), err
}

// ReverseSegments returns a new section with the segments reversed.
func (section *IPAddressSection) ReverseSegments() *IPAddressSection {
	if section.GetSegmentCount() <= 1 {
		if section.IsPrefixed() {
			return section.WithoutPrefixLen()
		}
		return section
	}
	res, _ := section.reverseSegments(
		func(i int) (*AddressSegment, addrerr.IncompatibleAddressError) {
			return section.GetSegment(i).withoutPrefixLen().ToSegmentBase(), nil
		},
	)
	return res.ToIP()
}

// String implements the [fmt.Stringer] interface, returning the normalized string provided by ToNormalizedString, or "<nil>" if the receiver is a nil pointer.
func (section *IPAddressSection) String() string {
	if section == nil {
		return nilString()
	}
	return section.toString()
}

// ToCanonicalString produces a canonical string for the address section.
//
// For IPv4, dotted octet format, also known as dotted decimal format, is used.
// https://datatracker.ietf.org/doc/html/draft-main-ipaddr-text-rep-00#section-2.1
//
// For IPv6, RFC 5952 describes canonical string representation.
// https://en.wikipedia.org/wiki/IPv6_address#Representation
// http://tools.ietf.org/html/rfc5952
//
// With IP addresses, the prefix length is included in the string.
func (section *IPAddressSection) ToCanonicalString() string {
	if section == nil {
		return nilString()
	}
	return section.toCanonicalString()
}

// ToNormalizedString produces a normalized string for the address section.
//
// For IPv4, it is the same as the canonical string.
//
// For IPv6, it differs from the canonical string.  Zero-segments are not compressed.
//
// With IP addresses, the prefix length is included in the string.
func (section *IPAddressSection) ToNormalizedString() string {
	if section == nil {
		return nilString()
	}
	return section.toNormalizedString()
}

// ToCompressedString produces a short representation of this address section while remaining within the confines of standard representation(s) of the address.
//
// For IPv4, it is the same as the canonical string.
//
// For IPv6, it differs from the canonical string.  It compresses the maximum number of zeros and/or host segments with the IPv6 compression notation '::'.
func (section *IPAddressSection) ToCompressedString() string {
	if section == nil {
		return nilString()
	}
	return section.toCompressedString()
}

// ToHexString writes this address section as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPAddressSection) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toHexString(with0xPrefix)
}

// ToOctalString writes this address section as a single octal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPAddressSection) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toOctalString(with0Prefix)
}

// ToBinaryString writes this address section as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPAddressSection) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toBinaryString(with0bPrefix)
}

// ToNormalizedWildcardString produces a string similar to the normalized string but avoids the CIDR prefix length.
// CIDR addresses will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix notation.
func (section *IPAddressSection) ToNormalizedWildcardString() string {
	if section == nil {
		return nilString()
	}
	return section.toNormalizedWildcardString()
}

// ToCanonicalWildcardString produces a string similar to the canonical string but avoids the CIDR prefix length.
// Address sections with a network prefix length will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix length notation.
// IPv6 sections will be compressed according to the canonical representation.
// For IPv4 it is the same as ToNormalizedWildcardString.
func (section *IPAddressSection) ToCanonicalWildcardString() string {
	if section == nil {
		return nilString()
	}
	return section.toCanonicalWildcardString()
}

// ToSegmentedBinaryString writes this IP address section as segments of binary values preceded by the "0b" prefix.
func (section *IPAddressSection) ToSegmentedBinaryString() string {
	if section == nil {
		return nilString()
	}
	return section.toSegmentedBinaryString()
}

// ToSQLWildcardString create a string similar to that from toNormalizedWildcardString except that
// it uses SQL wildcards.  It uses '%' instead of '*' and also uses the wildcard '_'.
func (section *IPAddressSection) ToSQLWildcardString() string {
	if section == nil {
		return nilString()
	}
	return section.toSQLWildcardString()
}

// ToFullString produces a string with no compressed segments and all segments of full length with leading zeros,
// which is 4 characters for IPv6 segments and 3 characters for IPv4 segments.
func (section *IPAddressSection) ToFullString() string {
	if section == nil {
		return nilString()
	}
	return section.toFullString()
}

// ToReverseDNSString generates the reverse-DNS lookup string,
// returning an error if this address section is an IPv6 multiple-valued section for which the range cannot be represented.
// For "8.255.4.4" it is "4.4.255.8.in-addr.arpa".
// For "2001:db8::567:89ab" it is "b.a.9.8.7.6.5.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa".
func (section *IPAddressSection) ToReverseDNSString() (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toReverseDNSString()
}

// ToPrefixLenString returns a string with a CIDR network prefix length if this address has a network prefix length.
// For IPv6, a zero host section will be compressed with "::". For IPv4 the string is equivalent to the canonical string.
func (section *IPAddressSection) ToPrefixLenString() string {
	if section == nil {
		return nilString()
	}
	return section.toPrefixLenString()
}

// ToSubnetString produces a string with specific formats for subnets.
// The subnet string looks like "1.2.*.*" or "1:2::/16".
//
// In the case of IPv4, this means that wildcards are used instead of a network prefix when a network prefix has been supplied.
// In the case of IPv6, when a network prefix has been supplied, the prefix will be shown and the host section will be compressed with "::".
func (section *IPAddressSection) ToSubnetString() string {
	if section == nil {
		return nilString()
	}
	return section.toSubnetString()
}

// ToCompressedWildcardString produces a string similar to ToNormalizedWildcardString, avoiding the CIDR prefix, but with full IPv6 segment compression as well, including single zero-segments.
// For IPv4 it is the same as ToNormalizedWildcardString.
func (section *IPAddressSection) ToCompressedWildcardString() string {
	if section == nil {
		return nilString()
	}
	return section.toCompressedWildcardString()
}

// ToCustomString creates a customized string from this address section according to the given string option parameters.
func (section *IPAddressSection) ToCustomString(stringOptions addrstr.IPStringOptions) string {
	if section == nil {
		return nilString()
	}
	return section.toCustomString(stringOptions)
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (section *IPAddressSection) GetSegmentStrings() []string {
	if section == nil {
		return nil
	}
	return section.getSegmentStrings()
}

var (
	rangeWildcard                 = new(addrstr.WildcardsBuilder).ToWildcards()
	allWildcards                  = new(addrstr.WildcardOptionsBuilder).SetWildcardOptions(addrstr.WildcardsAll).ToOptions()
	wildcardsRangeOnlyNetworkOnly = new(addrstr.WildcardOptionsBuilder).SetWildcards(rangeWildcard).ToOptions()
	allSQLWildcards               = new(addrstr.WildcardOptionsBuilder).SetWildcardOptions(addrstr.WildcardsAll).SetWildcards(
		new(addrstr.WildcardsBuilder).SetWildcard(SegmentSqlWildcardStr).SetSingleWildcard(SegmentSqlSingleWildcardStr).ToWildcards()).ToOptions()
)

// handles prefix block subnets, and ensures segment prefixes match the section prefix
func assignPrefix(prefixLength PrefixLen, segments []*AddressDivision, res *IPAddressSection, singleOnly, checkPrefixes bool, boundaryBits BitCount) {
	prefLen := prefixLength.bitCount()
	if prefLen < 0 {
		prefLen = 0
	} else if prefLen > boundaryBits {
		prefLen = boundaryBits
		prefixLength = cacheBitCount(boundaryBits)
	} else {
		prefixLength = cachePrefixLen(prefixLength) // use our own cache of prefix lengths so callers cannot overwrite a section's prefix length
	}
	segLen := len(segments)
	if segLen > 0 {
		var segProducer func(*AddressDivision, PrefixLen) *AddressDivision
		applyPrefixSubnet := !singleOnly && isPrefixSubnetDivs(segments, prefLen)
		if applyPrefixSubnet || checkPrefixes {
			if applyPrefixSubnet {
				segProducer = (*AddressDivision).toPrefixedNetworkDivision
			} else {
				segProducer = (*AddressDivision).toPrefixedDivision
			}
			applyPrefixToSegments(
				prefLen,
				segments,
				res.GetBitsPerSegment(),
				res.GetBytesPerSegment(),
				segProducer)
			if applyPrefixSubnet {
				res.isMult = res.isMult || res.GetSegment(segLen-1).isMultiple()
			}
		}
	}
	res.prefixLength = prefixLength
	return
}

// Starting from the first host bit according to the prefix, if the section is a sequence of zeros in both low and high values,
// followed by a sequence where low values are zero and high values are 1, then the section is a subnet prefix.
//
// Note that this includes sections where hosts are all zeros, or sections where hosts are full range of values,
// so the sequence of zeros can be empty and the sequence of where low values are zero and high values are 1 can be empty as well.
// However, if they are both empty, then this returns false, there must be at least one bit in the sequence.
func isPrefixSubnetDivs(sectionSegments []*AddressDivision, networkPrefixLength BitCount) bool {
	segmentCount := len(sectionSegments)
	if segmentCount == 0 {
		return false
	}
	seg := sectionSegments[0]
	return isPrefixSubnet(
		func(segmentIndex int) SegInt {
			return sectionSegments[segmentIndex].ToSegmentBase().GetSegmentValue()
		},
		func(segmentIndex int) SegInt {
			return sectionSegments[segmentIndex].ToSegmentBase().GetUpperSegmentValue()
		},
		segmentCount,
		seg.GetByteCount(),
		seg.GetBitCount(),
		seg.ToSegmentBase().GetMaxValue(),
		networkPrefixLength,
		zerosOnly)
}

func applyPrefixToSegments(
	sectionPrefixBits BitCount,
	segments []*AddressDivision,
	segmentBitCount BitCount,
	segmentByteCount int,
	segProducer func(*AddressDivision, PrefixLen) *AddressDivision) {
	var i int
	if sectionPrefixBits != 0 {
		i = getNetworkSegmentIndex(sectionPrefixBits, segmentByteCount, segmentBitCount)
	}
	for ; i < len(segments); i++ {
		pref := getPrefixedSegmentPrefixLength(segmentBitCount, sectionPrefixBits, i)
		if pref != nil {
			segments[i] = segProducer(segments[i], pref)
		}
	}
}

// called prior to check for prefix subnet.  The segments must be created first before that can happen.
func createSegmentsUint64(
	segLen int,
	highBytes,
	lowBytes uint64,
	bytesPerSegment int,
	bitsPerSegment BitCount,
	creator addressSegmentCreator,
	assignedPrefixLength PrefixLen) []*AddressDivision {
	segmentMask := ^(^SegInt(0) << uint(bitsPerSegment))
	lowSegCount := getHostSegmentIndex(64, bytesPerSegment, bitsPerSegment)
	newSegs := make([]*AddressDivision, segLen)
	lowIndex := segLen - lowSegCount
	if lowIndex < 0 {
		lowIndex = 0
	}
	segmentIndex := segLen - 1
	bytes := lowBytes
	for {
		for {
			segmentPrefixLength := getSegmentPrefixLength(bitsPerSegment, assignedPrefixLength, segmentIndex)
			value := segmentMask & SegInt(bytes)
			seg := creator.createSegment(value, value, segmentPrefixLength)
			newSegs[segmentIndex] = seg
			segmentIndex--
			if segmentIndex < lowIndex {
				break
			}
			bytes >>= uint(bitsPerSegment)
		}
		if lowIndex == 0 {
			break
		}
		lowIndex = 0
		bytes = highBytes
	}
	return newSegs
}

// called prior to check for prefix subnet.  The segments must be created first before that can happen.
func createSegments(
	lowerValueProvider,
	upperValueProvider SegmentValueProvider,
	segmentCount int,
	bitsPerSegment BitCount,
	creator addressSegmentCreator,
	prefixLength PrefixLen) (segments []*AddressDivision, isMultiple bool) {
	segments = createSegmentArray(segmentCount)
	for segmentIndex := 0; segmentIndex < segmentCount; segmentIndex++ {
		segmentPrefixLength := getSegmentPrefixLength(bitsPerSegment, prefixLength, segmentIndex)
		var value, value2 SegInt
		if lowerValueProvider == nil {
			value = upperValueProvider(segmentIndex)
			value2 = value
		} else {
			value = lowerValueProvider(segmentIndex)
			if upperValueProvider != nil {
				value2 = upperValueProvider(segmentIndex)
				if !isMultiple && value2 != value {
					isMultiple = true

				}
			} else {
				value2 = value
			}
		}
		seg := creator.createSegment(value, value2, segmentPrefixLength)
		segments[segmentIndex] = seg
	}
	return
}
