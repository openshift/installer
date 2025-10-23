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
	"fmt"
	"math/big"
	"unsafe"
)

type addressDivisionGroupingBase struct {
	// the non-cacheBitCount elements are assigned at creation and are immutable
	divisions divArray // either standard or large

	prefixLength PrefixLen // must align with the divisions if they store prefix lengths
	isMult       bool

	// When a top-level section is created, it is assigned an address type, IPv4, IPv6, or MAC,
	// and determines if an *AddressDivisionGrouping can be converted back to a section of the original type.
	//
	// Type-specific functions in IPAddressSection and lower levels, such as functions returning strings,
	// can rely on this field.
	addrType addrType

	// assigned on creation only; for zero-value groupings it is never assigned, but in that case it is not needed since there is nothing to cache
	cache *valueCache
}

// this is used by methods that are used by both mac and ipv4/6, even though the prefix length assignment does not apply to MAC.  It is also used by large division groupings.
func (grouping *addressDivisionGroupingBase) initMultAndPrefLen() {
	segCount := grouping.GetDivisionCount()
	bitsSoFar := 0
	if segCount != 0 {
		var previousSegmentPrefix PrefixLen
		isMultiple := false
		//bitsPerSegment := grouping.getBitsPerSegment()
		for i := 0; i < segCount; i++ {
			segment := grouping.getDivision(i)
			if !isMultiple && segment.isMultiple() {
				isMultiple = true
				grouping.isMult = true
				if grouping.prefixLength != nil { // nothing left to do
					break
				}
			}

			//Calculate the segment-level prefix
			//
			//Across an address prefixes are:
			//IPv6: (nil):...:(nil):(1 to 16):(0):...:(0)
			//or IPv4: ...(nil).(1 to 8).(0)...
			//For MAC, all segs have nil prefix since prefix is not segment-level
			segPrefix := segment.getDivisionPrefixLength()
			if previousSegmentPrefix == nil {
				if segPrefix != nil {
					newPref := bitsSoFar + segPrefix.bitCount()
					//newPref := getNetworkPrefixLen(bitsPerSegment, segPrefix.bitCount(), i)
					grouping.prefixLength = cacheBitCount(newPref)
					if isMultiple { // nothing left to do
						break
					}
				}
			}
			previousSegmentPrefix = segPrefix
			bitsSoFar += segment.GetBitCount()
		}
	}
	return
}

func (grouping *addressDivisionGroupingBase) getAddrType() addrType {
	return grouping.addrType
}

func (grouping *addressDivisionGroupingBase) isPrefixed() bool {
	return grouping.prefixLength != nil
}

func (grouping *addressDivisionGroupingBase) getPrefixLen() PrefixLen {
	return grouping.prefixLength
}

// GetPrefixLen returns the prefix length, or nil if there is no prefix length.
//
// A prefix length indicates the number of bits in the initial part of the address item that comprises the prefix.
//
// A prefix is a part of the address item that is not specific to that address but common amongst a group of such items, such as a CIDR prefix block subnet.
func (grouping *addressDivisionGroupingBase) GetPrefixLen() PrefixLen {
	return grouping.getPrefixLen().copy()
}

// isMultiple returns whether this address or grouping represents more than one address or grouping.
// Such addresses include CIDR/IP addresses (eg 1.2.3.4/11) or wildcard addresses (eg 1.2.*.4) or range addresses (eg 1.2.3-4.5)
func (grouping *addressDivisionGroupingBase) isMultiple() bool {
	return grouping.isMult
}

// hasNoDivisions() returns whether this grouping is the zero grouping,
// which is what you get when constructing a grouping or section with no divisions
func (grouping *addressDivisionGroupingBase) hasNoDivisions() bool {
	divisions := grouping.divisions
	return divisions == nil || divisions.getDivisionCount() == 0
}

// GetBitCount returns the total number of bits across all divisions
func (grouping *addressDivisionGroupingBase) GetBitCount() (res BitCount) {
	for i := 0; i < grouping.GetDivisionCount(); i++ {
		res += grouping.getDivision(i).GetBitCount()
	}
	return
}

// GetByteCount returns the total number of bytes across all divisions (rounded up)
func (grouping *addressDivisionGroupingBase) GetByteCount() int {
	return (int(grouping.GetBitCount()) + 7) >> 3
}

// getDivision returns the division or panics if the index is negative or it is too large
func (grouping *addressDivisionGroupingBase) getDivision(index int) *addressDivisionBase {
	return grouping.divisions.getDivision(index)
}

// GetGenericDivision returns the division as a DivisionType,
// allowing all division types and aggregated division types to be represented by a single type,
// useful for comparisons and other common uses.
func (grouping *addressDivisionGroupingBase) GetGenericDivision(index int) DivisionType {
	return grouping.divisions.getGenericDivision(index)
}

// GetDivisionCount returns the number of divisions in this grouping.
func (grouping *addressDivisionGroupingBase) GetDivisionCount() int {
	divisions := grouping.divisions
	if divisions != nil {
		return divisions.getDivisionCount()
	}
	return 0
}

// IsZero returns whether this grouping matches exactly the value of zero.
func (grouping *addressDivisionGroupingBase) IsZero() bool {
	divCount := grouping.GetDivisionCount()
	for i := 0; i < divCount; i++ {
		if !grouping.getDivision(i).IsZero() {
			return false
		}
	}
	return true
}

// IncludesZero returns whether this grouping includes the value of zero within its range.
func (grouping *addressDivisionGroupingBase) IncludesZero() bool {
	divCount := grouping.GetDivisionCount()
	for i := 0; i < divCount; i++ {
		if !grouping.getDivision(i).IncludesZero() {
			return false
		}
	}
	return true
}

// IsMax returns whether this grouping matches exactly the maximum possible value, the value whose bits are all ones.
func (grouping *addressDivisionGroupingBase) IsMax() bool {
	divCount := grouping.GetDivisionCount()
	for i := 0; i < divCount; i++ {
		if !grouping.getDivision(i).IsMax() {
			return false
		}
	}
	return true
}

// IncludesMax returns whether this grouping includes the max value, the value whose bits are all ones, within its range.
func (grouping *addressDivisionGroupingBase) IncludesMax() bool {
	divCount := grouping.GetDivisionCount()
	for i := 0; i < divCount; i++ {
		if !grouping.getDivision(i).IncludesMax() {
			return false
		}
	}
	return true
}

// IsFullRange returns whether this address item represents all possible values attainable by an address item of this type.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (grouping *addressDivisionGroupingBase) IsFullRange() bool {
	divCount := grouping.GetDivisionCount()
	for i := 0; i < divCount; i++ {
		if !grouping.getDivision(i).IsFullRange() {
			return false
		}
	}
	return true
}

// GetSequentialBlockIndex gets the minimal division index for which all following divisions are full-range blocks.
//
// The division at this index is not a full-range block unless all divisions are full-range.
// The division at this index and all following divisions form a sequential range.
// For the full grouping to be sequential, the preceding divisions must be single-valued.
func (grouping *addressDivisionGroupingBase) GetSequentialBlockIndex() int {
	divCount := grouping.GetDivisionCount()
	if divCount > 0 {
		for divCount--; divCount > 0 && grouping.getDivision(divCount).IsFullRange(); divCount-- {
		}
	}
	return divCount
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential address division groupings that comprise this address division grouping.
func (grouping *addressDivisionGroupingBase) GetSequentialBlockCount() *big.Int {
	sequentialSegCount := grouping.GetSequentialBlockIndex()
	prefixLen := BitCount(0)
	for i := 0; i < sequentialSegCount; i++ {
		prefixLen += grouping.getDivision(i).GetBitCount()
	}
	return grouping.GetPrefixCountLen(prefixLen) // 0-1.0-1.*.* gives 1 as seq block index, and then you count only previous segments
}

func (grouping *addressDivisionGroupingBase) getCountBig() *big.Int {
	res := bigOne()
	count := grouping.GetDivisionCount()
	if count > 0 {
		for i := 0; i < count; i++ {
			div := grouping.getDivision(i)
			if div.isMultiple() {
				res.Mul(res, div.getCount())
			}
		}
	}
	return res
}

func (grouping *addressDivisionGroupingBase) getPrefixCountBig() *big.Int {
	prefixLen := grouping.prefixLength
	if prefixLen == nil {
		return grouping.getCountBig()
	}
	return grouping.getPrefixCountLenBig(prefixLen.bitCount())
}

func (grouping *addressDivisionGroupingBase) getPrefixCountLenBig(prefixLen BitCount) *big.Int {
	if prefixLen <= 0 {
		return bigOne()
	} else if prefixLen >= grouping.GetBitCount() {
		return grouping.getCountBig()
	}
	res := bigOne()
	if grouping.isMultiple() {
		divisionCount := grouping.GetDivisionCount()
		divPrefixLength := prefixLen
		for i := 0; i < divisionCount; i++ {
			div := grouping.getDivision(i)
			divBitCount := div.getBitCount()
			if div.isMultiple() {
				var divCount *big.Int
				if divPrefixLength < divBitCount {
					divCount = div.GetPrefixCountLen(divPrefixLength)
				} else {
					divCount = div.getCount()
				}
				res.Mul(res, divCount)
			}
			if divPrefixLength <= divBitCount {
				break
			}
			divPrefixLength -= divBitCount
		}
	}
	return res
}

func (grouping *addressDivisionGroupingBase) getBlockCountBig(segmentCount int) *big.Int {
	if segmentCount <= 0 {
		return bigOne()
	}
	divCount := grouping.GetDivisionCount()
	if segmentCount >= divCount {
		return grouping.getCountBig()
	}
	res := bigOne()
	if grouping.isMultiple() {
		for i := 0; i < segmentCount; i++ {
			division := grouping.getDivision(i)
			if division.isMultiple() {
				res.Mul(res, division.getCount())
			}
		}
	}
	return res
}

func (grouping *addressDivisionGroupingBase) getCount() *big.Int {
	return grouping.cacheCount(grouping.getCountBig)
}

func (grouping *addressDivisionGroupingBase) getCachedCount() *big.Int {
	return grouping.cachedCount(grouping.getCountBig)
}

// GetPrefixCount returns the number of distinct prefix values in this item.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetCount.
func (grouping *addressDivisionGroupingBase) GetPrefixCount() *big.Int {
	return grouping.cachePrefixCount(grouping.getPrefixCountBig)
}

// GetPrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
func (grouping *addressDivisionGroupingBase) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	return grouping.calcCount(func() *big.Int { return grouping.getPrefixCountLenBig(prefixLen) })
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) divisions.
func (grouping *addressDivisionGroupingBase) GetBlockCount(divisionCount int) *big.Int {
	return grouping.calcCount(func() *big.Int { return grouping.getBlockCountBig(divisionCount) })
}

func (grouping *addressDivisionGroupingBase) cacheCount(counter func() *big.Int) *big.Int {
	cache := grouping.cache
	if cache == nil {
		return grouping.calcCount(counter)
	}
	count := (*big.Int)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.cachedCount))))
	if count == nil {
		count = grouping.calcCount(counter)
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.cachedCount))
		atomicStorePointer(dataLoc, unsafe.Pointer(count))
	}
	return bigZero().Set(count)
}

// cachedCount returns the cached count value, not a duplicate
func (grouping *addressDivisionGroupingBase) cachedCount(counter func() *big.Int) *big.Int {
	cache := grouping.cache
	if cache == nil {
		return grouping.calcCount(counter)
	}
	count := (*big.Int)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.cachedCount))))
	if count == nil {
		count = grouping.calcCount(counter)
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.cachedCount))
		atomicStorePointer(dataLoc, unsafe.Pointer(count))
	}
	return count
}

func (grouping *addressDivisionGroupingBase) calcCount(counter func() *big.Int) *big.Int {
	if grouping != nil && !grouping.isMultiple() {
		return bigOne()
	}
	return counter()
}

func (grouping *addressDivisionGroupingBase) calcUint64Count(counter func() uint64) uint64 {
	if grouping != nil && !grouping.isMultiple() {
		return 1
	}
	return counter()
}

func (grouping *addressDivisionGroupingBase) cacheUint64PrefixCount(counter func() uint64) uint64 {
	cache := grouping.cache // isMultiple checks prior to this ensures cache not nil here
	if cache == nil {
		return grouping.calcUint64PrefixCount(counter)
	}
	count := (*big.Int)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.cachedPrefixCount))))
	if count == nil {
		count64 := grouping.calcUint64PrefixCount(counter)
		count = bigZero().SetUint64(count64)
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.cachedPrefixCount))
		atomicStorePointer(dataLoc, unsafe.Pointer(count))
		return count64
	}
	return count.Uint64()
}

func (grouping *addressDivisionGroupingBase) cachePrefixCount(counter func() *big.Int) *big.Int {
	cache := grouping.cache // isMultiple checks prior to this ensures cache not nil here
	if cache == nil {
		return grouping.calcPrefixCount(counter)
	}
	count := (*big.Int)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.cachedPrefixCount))))
	if count == nil {
		count = grouping.calcPrefixCount(counter)
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.cachedPrefixCount))
		atomicStorePointer(dataLoc, unsafe.Pointer(count))
	}
	return bigZero().Set(count)
}

func (grouping *addressDivisionGroupingBase) calcPrefixCount(counter func() *big.Int) *big.Int {
	if !grouping.isMultiple() {
		return bigOne()
	}
	prefixLen := grouping.prefixLength
	if prefixLen == nil || prefixLen.bitCount() >= grouping.GetBitCount() {
		return grouping.getCount()
	}
	return counter()
}

func (grouping *addressDivisionGroupingBase) calcUint64PrefixCount(counter func() uint64) uint64 {
	if !grouping.isMultiple() {
		return 1
	}
	//prefixLen := grouping.prefixLength
	//if prefixLen == nil || prefixLen.bitCount() >= grouping.GetBitCount() {
	//	return grouping.getCount()
	//}
	return counter()
}

func (grouping *addressDivisionGroupingBase) getCachedBytes(calcBytes func() (bytes, upperBytes []byte)) (bytes, upperBytes []byte) {
	cache := grouping.cache
	if cache == nil {
		return calcBytes()
	}
	cached := (*bytesCache)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.bytesCache))))
	if cached == nil {
		bytes, upperBytes = calcBytes()
		cached = &bytesCache{
			lowerBytes: bytes,
			upperBytes: upperBytes,
		}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.bytesCache))
		atomicStorePointer(dataLoc, unsafe.Pointer(cached))
	}
	bytes = cached.lowerBytes
	upperBytes = cached.upperBytes
	return
}

// IsSequential returns whether the grouping represents a range of values that are sequential.
//
// Generally, this means that any division covering a range of values must be followed by divisions that are full range, covering all values.
func (grouping *addressDivisionGroupingBase) IsSequential() bool {
	if grouping.isMultiple() {
		count := grouping.GetDivisionCount()
		if count > 1 {
			for i := 0; i < count; i++ {
				if grouping.getDivision(i).isMultiple() {
					for i++; i < count; i++ {
						if !grouping.getDivision(i).IsFullRange() {
							return false
						}
					}
					return true
				}
			}
		}
	}
	return true
}

type bytesCache struct {
	lowerBytes, upperBytes []byte
}

type mixedCache struct {
	defaultMixedAddressSection *IPv6v4MixedAddressGrouping
	embeddedIPv4Section        *IPv4AddressSection
	embeddedIPv6Section        *EmbeddedIPv6AddressSection
}

type uint128Cache struct {
	high, low uint64
}

type valueCache struct {
	cachedCount, cachedPrefixCount *big.Int

	cachedMaskLens *maskLenSetting

	bytesCache *bytesCache

	uint128Cache *uint128Cache

	uint32Cache *uint32

	stringCache stringCache

	sectionCache *groupingCache

	mixed *mixedCache

	minPrefix PrefixLen

	equivalentPrefix *PrefixLen

	isSinglePrefixBlock *bool
}

type ipStringCache struct {
	normalizedWildcardString,
	fullString,
	sqlWildcardString,

	reverseDNSString,

	segmentedBinaryString *string
}

type ipv4StringCache struct {
	inetAtonOctalString,
	inetAtonHexString *string
}

type ipv6StringCache struct {
	normalizedIPv6String,
	compressedIPv6String,
	mixedString,
	compressedWildcardString,
	canonicalWildcardString,
	networkPrefixLengthString,
	base85String,
	uncString *string
}

type macStringCache struct {
	normalizedMACString,
	compressedMACString,
	dottedString,
	spaceDelimitedString *string
}

type stringCache struct {
	canonicalString *string

	octalString, octalStringPrefixed,
	binaryString, binaryStringPrefixed,
	hexString, hexStringPrefixed *string

	*ipv6StringCache

	*ipv4StringCache

	*ipStringCache

	*macStringCache
}

var zeroStringCache = stringCache{
	ipv6StringCache: &ipv6StringCache{},
	ipv4StringCache: &ipv4StringCache{},
	ipStringCache:   &ipStringCache{},
	macStringCache:  &macStringCache{},
}

type groupingCache struct {
	lower, upper *AddressSection
}

type maskLenSetting struct {
	networkMaskLen, hostMaskLen PrefixLen
}

type divArray interface {
	getDivision(index int) *addressDivisionBase

	getGenericDivision(index int) DivisionType

	getDivisionCount() int

	fmt.Stringer
}

var zeroDivs = make([]*AddressDivision, 0)
var zeroStandardDivArray = standardDivArray(zeroDivs)

type standardDivArray []*AddressDivision

func (grouping standardDivArray) getDivisionCount() int {
	return len(grouping)
}

func (grouping standardDivArray) getDivision(index int) *addressDivisionBase {
	return (*addressDivisionBase)(unsafe.Pointer(grouping[index]))
}

func (grouping standardDivArray) getGenericDivision(index int) DivisionType {
	return grouping[index]
}

func (grouping standardDivArray) copyDivisions(divs []*AddressDivision) (count int) {
	return copy(divs, grouping)
}

func (grouping standardDivArray) copySubDivisions(start, end int, divs []*AddressDivision) (count int) {
	return copy(divs, grouping[start:end])
}

func (grouping standardDivArray) getSubDivisions(index, endIndex int) (divs []*AddressDivision) {
	return grouping[index:endIndex]
}

func (grouping standardDivArray) init() standardDivArray {
	if grouping == nil {
		return zeroStandardDivArray
	}
	return grouping
}

func (grouping standardDivArray) String() string {
	return fmt.Sprint([]*AddressDivision(grouping.init()))
}

var zeroLargeDivs = make([]*IPAddressLargeDivision, 0)
var zeroLargeDivArray = largeDivArray(zeroLargeDivs)

type largeDivArray []*IPAddressLargeDivision

func (grouping largeDivArray) getDivisionCount() int {
	return len(grouping)
}

func (grouping largeDivArray) getDivision(index int) *addressDivisionBase {
	return (*addressDivisionBase)(unsafe.Pointer(grouping[index]))
}

func (grouping largeDivArray) getGenericDivision(index int) DivisionType {
	return grouping[index]
}

func (grouping largeDivArray) copyDivisions(divs []*IPAddressLargeDivision) (count int) {
	return copy(divs, grouping)
}

func (grouping largeDivArray) copySubDivisions(start, end int, divs []*IPAddressLargeDivision) (count int) {
	return copy(divs, grouping[start:end])
}

func (grouping largeDivArray) getSubDivisions(index, endIndex int) (divs []*IPAddressLargeDivision) {
	return grouping[index:endIndex]
}

func (grouping largeDivArray) init() largeDivArray {
	if grouping == nil {
		return zeroLargeDivArray
	}
	return grouping
}

func (grouping largeDivArray) String() string {
	return fmt.Sprint([]*IPAddressLargeDivision(grouping.init()))
}

var _, _ divArray = standardDivArray{}, largeDivArray{}
