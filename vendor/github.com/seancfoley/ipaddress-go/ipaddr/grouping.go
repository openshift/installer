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

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
)

func createGrouping(divs []*AddressDivision, prefixLength PrefixLen, addrType addrType) *AddressDivisionGrouping {
	grouping := &AddressDivisionGrouping{
		addressDivisionGroupingInternal{
			addressDivisionGroupingBase: addressDivisionGroupingBase{
				divisions:    standardDivArray(divs),
				prefixLength: prefixLength,
				addrType:     addrType,
				cache:        &valueCache{},
			},
		},
	}
	assignStringCache(&grouping.addressDivisionGroupingBase, addrType)
	return grouping
}

// callers to this function have segments/divisions with prefix length consistent with the supplied prefix length
func createGroupingMultiple(divs []*AddressDivision, prefixLength PrefixLen, isMultiple bool) *AddressDivisionGrouping {
	result := createGrouping(divs, prefixLength, zeroType)
	result.isMult = isMultiple
	return result
}

// callers to this function have segments/divisions with prefix length consistent with the supplied prefix length
func createInitializedGrouping(divs []*AddressDivision, prefixLength PrefixLen) *AddressDivisionGrouping {
	result := createGrouping(divs, prefixLength, zeroType)
	result.initMultiple() // assigns isMultiple
	return result
}

var (
	emptyBytes = make([]byte, 0, 0)
)

type addressDivisionGroupingInternal struct {
	addressDivisionGroupingBase

	// TODO LATER refactor to support infiniband, which will involve multiple types.
	// But that will be a joint effort with Java and will wait to later.
}

func createSegmentArray(length int) []*AddressDivision {
	return make([]*AddressDivision, length)
}

func (grouping *addressDivisionGroupingInternal) initMultiple() {
	divCount := grouping.getDivisionCount()
	for i := divCount - 1; i >= 0; i-- {
		div := grouping.getDivision(i)
		if div.isMultiple() {
			grouping.isMult = true
			return
		}
	}
	return
}

func (grouping *addressDivisionGroupingInternal) getDivArray() standardDivArray {
	if divsArray := grouping.divisions; divsArray != nil {
		return divsArray.(standardDivArray)
	}
	return nil
}

// getDivision returns the division or panics if the index is negative or too large
func (grouping *addressDivisionGroupingInternal) getDivision(index int) *AddressDivision {
	return grouping.getDivArray()[index]
}

// getDivisionsInternal returns the divisions slice, only to be used internally
func (grouping *addressDivisionGroupingInternal) getDivisionsInternal() []*AddressDivision {
	return grouping.getDivArray()
}

func (grouping *addressDivisionGroupingInternal) getDivisionCount() int {
	if divArray := grouping.getDivArray(); divArray != nil {
		return divArray.getDivisionCount()
	}
	return 0
}

func (grouping *addressDivisionGroupingInternal) forEachSubDivision(start, end int, target func(index int, div *AddressDivision), targetLen int) (count int) {
	divArray := grouping.getDivArray()
	if divArray != nil {
		// if not enough space in target, adjust
		if targetEnd := start + targetLen; end > targetEnd {
			end = targetEnd
		}
		divArray = divArray[start:end]
		for i, div := range divArray {
			target(i, div)
		}
	}
	return len(divArray)
}

func adjust1To1StartIndices(sourceStart, sourceEnd, sourceCount, targetCount int) (newSourceStart, newSourceEnd, newTargetStart int) {
	// both sourceCount and targetCount are lengths of their respective slices, so never negative
	targetStart := 0
	if sourceStart < 0 {
		targetStart -= sourceStart
		sourceStart = 0
		if targetStart > targetCount || targetStart < 0 /* overflow */ {
			targetStart = targetCount
		}
	} else if sourceStart > sourceCount {
		sourceStart = sourceCount
	}
	// how many to copy?
	if sourceEnd > sourceCount { // end index exceeds available
		sourceEnd = sourceCount
	} else if sourceEnd < sourceStart {
		sourceEnd = sourceStart
	}
	return sourceStart, sourceEnd, targetStart
}

func adjust1To1Indices(sourceStart, sourceEnd, sourceCount, targetCount int) (newSourceStart, newSourceEnd, newTargetStart int) {
	var targetStart int
	sourceStart, sourceEnd, targetStart = adjust1To1StartIndices(sourceStart, sourceEnd, sourceCount, targetCount)
	if limitEnd := sourceStart + (targetCount - targetStart); sourceEnd > limitEnd {
		sourceEnd = limitEnd
	}
	return sourceStart, sourceEnd, targetStart
}

func adjustIndices(
	startIndex, endIndex, sourceCount,
	replacementStartIndex, replacementEndIndex, replacementSegmentCount int) (int, int, int, int) {
	if startIndex < 0 {
		startIndex = 0
	} else if startIndex > sourceCount {
		startIndex = sourceCount
	}
	if endIndex < startIndex {
		endIndex = startIndex
	} else if endIndex > sourceCount {
		endIndex = sourceCount
	}
	if replacementStartIndex < 0 {
		replacementStartIndex = 0
	} else if replacementStartIndex > replacementSegmentCount {
		replacementStartIndex = replacementSegmentCount
	}
	if replacementEndIndex < replacementStartIndex {
		replacementEndIndex = replacementStartIndex
	} else if replacementEndIndex > replacementSegmentCount {
		replacementEndIndex = replacementSegmentCount
	}
	return startIndex, endIndex, replacementStartIndex, replacementEndIndex
}

// copySubDivisions copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (grouping *addressDivisionGroupingInternal) copySubDivisions(start, end int, divs []*AddressDivision) (count int) {
	if divArray := grouping.getDivArray(); divArray != nil {
		start, end, targetIndex := adjust1To1Indices(start, end, grouping.GetDivisionCount(), len(divs))
		return divArray.copySubDivisions(start, end, divs[targetIndex:])
	}
	return
}

// copyDivisions copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (grouping *addressDivisionGroupingInternal) copyDivisions(divs []*AddressDivision) (count int) {
	if divArray := grouping.getDivArray(); divArray != nil {
		return divArray.copyDivisions(divs)
	}
	return
}

func (grouping *addressDivisionGroupingInternal) getSubDivisions(start, end int) []*AddressDivision {
	divArray := grouping.getDivArray()
	if divArray != nil {
		return divArray.getSubDivisions(start, end)
	} else if start != 0 || end != 0 {
		panic("invalid subslice")
	}
	return make([]*AddressDivision, 0)
}

func (grouping *addressDivisionGroupingInternal) isAddressSection() bool {
	return grouping != nil && grouping.matchesAddrSectionType()
}

func (grouping *addressDivisionGroupingInternal) compareSize(other AddressItem) int { // the getCount() is optimized which is why we do not defer to the method in addressDivisionGroupingBase
	return compareCount(grouping.toAddressDivisionGrouping(), other)
}

func (grouping *addressDivisionGroupingInternal) getCount() *big.Int {
	if !grouping.isMultiple() {
		return bigOne()
	} else {
		g := grouping.toAddressDivisionGrouping()
		if sect := g.ToIPv4(); sect != nil {
			return sect.GetCount()
		} else if sect := g.ToIPv6(); sect != nil {
			return sect.GetCount()
		} else if sect := g.ToMAC(); sect != nil {
			return sect.GetCount()
		}
	}
	return grouping.addressDivisionGroupingBase.getCount()
}

func (grouping *addressDivisionGroupingInternal) getCachedCount() *big.Int {
	if !grouping.isMultiple() {
		return bigOne()
	} else {
		g := grouping.toAddressDivisionGrouping()
		if sect := g.ToIPv4(); sect != nil {
			return sect.getCachedCount()
		} else if sect := g.ToIPv6(); sect != nil {
			return sect.getCachedCount()
		} else if sect := g.ToMAC(); sect != nil {
			return sect.getCachedCount()
		}
	}
	return grouping.addressDivisionGroupingBase.getCachedCount()
}

// GetPrefixCount returns the number of distinct prefix values in this item.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetCount.
func (grouping *addressDivisionGroupingInternal) GetPrefixCount() *big.Int {
	if section := grouping.toAddressSection(); section != nil {
		return section.GetPrefixCount()
	}
	return grouping.addressDivisionGroupingBase.GetPrefixCount()
}

// GetPrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
func (grouping *addressDivisionGroupingInternal) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	if section := grouping.toAddressSection(); section != nil {
		return section.GetPrefixCountLen(prefixLen)
	}
	return grouping.addressDivisionGroupingBase.GetPrefixCountLen(prefixLen)
}

func (grouping *addressDivisionGroupingInternal) getDivisionStrings() []string {
	if grouping.hasNoDivisions() {
		return []string{}
	}
	result := make([]string, grouping.GetDivisionCount())
	for i := range result {
		result[i] = grouping.getDivision(i).String()
	}
	return result
}

func (grouping *addressDivisionGroupingInternal) getSegmentStrings() []string {
	if grouping.hasNoDivisions() {
		return []string{}
	}
	result := make([]string, grouping.GetDivisionCount())
	for i := range result {
		result[i] = grouping.getDivision(i).GetWildcardString()
	}
	return result
}

func (grouping *addressDivisionGroupingInternal) toAddressDivisionGrouping() *AddressDivisionGrouping {
	return (*AddressDivisionGrouping)(unsafe.Pointer(grouping))
}

func (grouping *addressDivisionGroupingInternal) toAddressSection() *AddressSection {
	return grouping.toAddressDivisionGrouping().ToSectionBase()
}

func (grouping *addressDivisionGroupingInternal) matchesIPv6AddressType() bool {
	return grouping.getAddrType().isIPv6() // no need to check segment count because addresses cannot be constructed with incorrect segment count
}

func (grouping *addressDivisionGroupingInternal) matchesIPv4AddressType() bool {
	return grouping.getAddrType().isIPv4() // no need to check segment count because addresses cannot be constructed with incorrect segment count
}

func (grouping *addressDivisionGroupingInternal) matchesIPAddressType() bool {
	return grouping.matchesIPSectionType() // no need to check segment count because addresses cannot be constructed with incorrect segment count (note the zero IPAddress has zero-segments)
}

func (grouping *addressDivisionGroupingInternal) matchesMACAddressType() bool {
	return grouping.getAddrType().isMAC()
}

// The adaptive zero grouping, produced by zero sections like IPv4AddressSection{} or AddressDivisionGrouping{}, can represent a zero-length section of any address type,
// It is not considered equal to constructions of specific zero length sections of groupings like NewIPv4Section(nil) which can only represent a zero-length section of a single address type.
func (grouping *addressDivisionGroupingInternal) matchesZeroGrouping() bool {
	addrType := grouping.getAddrType()
	return addrType.isZeroSegments() && grouping.hasNoDivisions()
}

func (grouping *addressDivisionGroupingInternal) matchesAddrSectionType() bool {
	addrType := grouping.getAddrType()
	// because there are no init() conversions for IPv6/IPV4/MAC sections, an implicitly zero-valued IPv6/IPV4/MAC or zero IP section has addr type nil
	return addrType.isIP() || addrType.isMAC() || grouping.matchesZeroGrouping()
}

func (grouping *addressDivisionGroupingInternal) matchesIPv6SectionType() bool {
	// because there are no init() conversions for IPv6 sections, an implicitly zero-valued IPV6 section has addr type nil
	return grouping.getAddrType().isIPv6() || grouping.matchesZeroGrouping()
}

func (grouping *addressDivisionGroupingInternal) matchesIPv6v4MixedGroupingType() bool {
	// because there are no init() conversions for IPv6v4MixedGrouping groupings, an implicitly zero-valued IPv6v4MixedGrouping has addr type nil
	return grouping.getAddrType().isIPv6v4Mixed() || grouping.matchesZeroGrouping()
}

func (grouping *addressDivisionGroupingInternal) matchesIPv4SectionType() bool {
	// because there are no init() conversions for IPV4 sections, an implicitly zero-valued IPV4 section has addr type nil
	return grouping.getAddrType().isIPv4() || grouping.matchesZeroGrouping()
}

func (grouping *addressDivisionGroupingInternal) matchesIPSectionType() bool {
	// because there are no init() conversions for IPv6 or IPV4 sections, an implicitly zero-valued IPv4, IPv6 or IP section has addr type nil
	return grouping.getAddrType().isIP() || grouping.matchesZeroGrouping()
}

func (grouping *addressDivisionGroupingInternal) matchesMACSectionType() bool {
	// because there are no init() conversions for MAC sections, an implicitly zero-valued MAC section has addr type nil
	return grouping.getAddrType().isMAC() || grouping.matchesZeroGrouping()
}

// Format implements the [fmt.Formatter] interface. It accepts the formats
//   - 'v' for the default address and section format (either the normalized or canonical string),
//   - 's' (string) for the same,
//   - 'b' (binary), 'o' (octal with 0 prefix), 'O' (octal with 0o prefix),
//   - 'd' (decimal), 'x' (lowercase hexadecimal), and
//   - 'X' (uppercase hexadecimal).
//
// Also supported are some of fmt's format flags for integral types.
// Sign control is not supported since addresses and sections are never negative.
// '#' for an alternate format is supported, which adds a leading zero for octal, and for hexadecimal it adds
// a leading "0x" or "0X" for "%#x" and "%#X" respectively.
// Also supported is specification of minimum digits precision, output field width,
// space or zero padding, and '-' for left or right justification.
func (grouping addressDivisionGroupingInternal) Format(state fmt.State, verb rune) {
	if sect := grouping.toAddressSection(); sect != nil {
		sect.Format(state, verb)
		return
	} else if mixed := grouping.toAddressDivisionGrouping().ToMixedIPv6v4(); mixed != nil {
		mixed.Format(state, verb)
		return
	}
	// divisions are printed like slices of *AddressDivision (which are Stringers) with division separated by spaces and enclosed in square brackets,
	// sections are printed like addresses with segments separated by segment separators
	grouping.defaultFormat(state, verb)
}

func (grouping addressDivisionGroupingInternal) defaultFormat(state fmt.State, verb rune) {
	s := flagsFromState(state, verb)
	_, _ = state.Write([]byte(fmt.Sprintf(s, grouping.initDivs().getDivArray())))
}

func (grouping *addressDivisionGroupingInternal) toString() string {
	if sect := grouping.toAddressSection(); sect != nil {
		return sect.ToNormalizedString()
	}
	return fmt.Sprint(grouping.initDivs().getDivArray())
}

func (grouping *addressDivisionGroupingInternal) initDivs() *addressDivisionGroupingInternal {
	if grouping.divisions == nil {
		return &zeroSection.addressDivisionGroupingInternal
	}
	return grouping
}

// ContainsPrefixBlock returns whether the values of this item contains the block of values for the given prefix length.
//
// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values in this item for the given prefix length makes no difference.
//
// Use GetMinPrefixLenForBlock to determine the smallest prefix length for which this method returns true.
func (grouping *addressDivisionGroupingInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	if section := grouping.toAddressSection(); section != nil {
		return section.ContainsPrefixBlock(prefixLen)
	}
	prefixLen = checkSubnet(grouping, prefixLen)
	divisionCount := grouping.GetDivisionCount()
	var prevBitCount BitCount
	for i := 0; i < divisionCount; i++ {
		division := grouping.getDivision(i)
		bitCount := division.GetBitCount()
		totalBitCount := bitCount + prevBitCount
		if prefixLen < totalBitCount {
			divPrefixLen := prefixLen - prevBitCount
			if !division.containsPrefixBlock(divPrefixLen) {
				return false
			}
			for i++; i < divisionCount; i++ {
				division = grouping.getDivision(i)
				if !division.IsFullRange() {
					return false
				}
			}
			return true
		}
		prevBitCount = totalBitCount
	}
	return true
}

// ContainsSinglePrefixBlock returns whether the values of this grouping contains a single prefix block for the given prefix length.
//
// This means there is only one prefix of the given length in this item, and this item contains the prefix block for that given prefix.
//
// Use GetPrefixLenForSingleBlock to determine whether there is a prefix length for which this method returns true.
func (grouping *addressDivisionGroupingInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	prefixLen = checkSubnet(grouping, prefixLen)
	divisionCount := grouping.GetDivisionCount()
	var prevBitCount BitCount
	for i := 0; i < divisionCount; i++ {
		division := grouping.getDivision(i)
		bitCount := division.getBitCount()
		totalBitCount := bitCount + prevBitCount
		if prefixLen >= totalBitCount {
			if division.isMultiple() {
				return false
			}
		} else {
			divPrefixLen := prefixLen - prevBitCount
			if !division.ContainsSinglePrefixBlock(divPrefixLen) {
				return false
			}
			for i++; i < divisionCount; i++ {
				division = grouping.getDivision(i)
				if !division.IsFullRange() {
					return false
				}
			}
			return true
		}
		prevBitCount = totalBitCount
	}
	return true
}

// IsPrefixBlock returns whether this address division series has a prefix length and includes the block associated with its prefix length.
// If the prefix length matches the bit count, this returns true.
//
// This is different from ContainsPrefixBlock in that this method returns
// false if the series has no prefix length, or a prefix length that differs from a prefix length for which ContainsPrefixBlock returns true.
func (grouping *addressDivisionGroupingInternal) IsPrefixBlock() bool { //Note for any given prefix length you can compare with GetMinPrefixLenForBlock
	prefLen := grouping.getPrefixLen()
	return prefLen != nil && grouping.ContainsPrefixBlock(prefLen.bitCount())
}

// IsSinglePrefixBlock returns whether the range of values matches a single subnet block for the prefix length.
//
// What distinguishes this method with ContainsSinglePrefixBlock is that this method returns
// false if the series does not have a prefix length assigned to it,
// or a prefix length that differs from the prefix length for which ContainsSinglePrefixBlock returns true.
//
// It is similar to IsPrefixBlock but returns false when there are multiple prefixes.
func (grouping *addressDivisionGroupingInternal) IsSinglePrefixBlock() bool { //Note for any given prefix length you can compare with GetPrefixLenForSingleBlock
	calc := func() bool {
		prefLen := grouping.getPrefixLen()
		return prefLen != nil && grouping.ContainsSinglePrefixBlock(prefLen.bitCount())
	}
	return cacheIsSinglePrefixBlock(grouping.cache, grouping.getPrefixLen(), calc)
}

func cacheIsSinglePrefixBlock(cache *valueCache, prefLen PrefixLen, calc func() bool) bool {
	if cache == nil {
		return calc()
	}
	res := (*bool)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.isSinglePrefixBlock))))
	if res == nil {
		if calc() {
			res = &trueVal

			// we can also set related cache fields
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.equivalentPrefix))
			equivPref := cachePrefix(prefLen.bitCount())
			atomicStorePointer(dataLoc, unsafe.Pointer(equivPref))

			dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&cache.minPrefix))
			atomicStorePointer(dataLoc, unsafe.Pointer(prefLen))
		} else {
			res = &falseVal
		}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.isSinglePrefixBlock))
		atomicStorePointer(dataLoc, unsafe.Pointer(res))
	}
	return *res
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this grouping includes the block of all values for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this grouping represents a single value, this returns the bit count.
func (grouping *addressDivisionGroupingInternal) GetMinPrefixLenForBlock() BitCount {
	calc := func() BitCount {
		count := grouping.GetDivisionCount()
		totalPrefix := grouping.GetBitCount()
		for i := count - 1; i >= 0; i-- {
			div := grouping.getDivision(i)
			segBitCount := div.getBitCount()
			segPrefix := div.GetMinPrefixLenForBlock()
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
		return totalPrefix
	}
	return cacheMinPrefix(grouping.cache, calc)
}

func cacheMinPrefix(cache *valueCache, calc func() BitCount) BitCount {
	if cache == nil {
		return calc()
	}
	res := (PrefixLen)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.minPrefix))))
	if res == nil {
		val := calc()
		res = cacheBitCount(val)
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.minPrefix))
		atomicStorePointer(dataLoc, unsafe.Pointer(res))
	}
	return res.bitCount()
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this division grouping matches the block of addresses for that prefix.
//
// If no such prefix exists, GetPrefixLenForSingleBlock returns nil.
//
// If this division grouping represents a single value, returns the bit length.
func (grouping *addressDivisionGroupingInternal) GetPrefixLenForSingleBlock() PrefixLen {
	calc := func() *PrefixLen {
		count := grouping.GetDivisionCount()
		var totalPrefix BitCount
		for i := 0; i < count; i++ {
			div := grouping.getDivision(i)
			divPrefix := div.GetPrefixLenForSingleBlock()
			if divPrefix == nil {
				return cacheNilPrefix()
			}
			divPrefLen := divPrefix.bitCount()
			totalPrefix += divPrefLen
			if divPrefLen < div.GetBitCount() {
				//remaining segments must be full range or we return nil
				for i++; i < count; i++ {
					laterDiv := grouping.getDivision(i)
					if !laterDiv.IsFullRange() {
						return cacheNilPrefix()
					}
				}
			}
		}
		return cachePrefix(totalPrefix)
	}
	return cachePrefLenSingleBlock(grouping.cache, grouping.getPrefixLen(), calc)
}

func cachePrefLenSingleBlock(cache *valueCache, prefLen PrefixLen, calc func() *PrefixLen) PrefixLen {
	if cache == nil {
		return *calc()
	}
	res := (*PrefixLen)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.equivalentPrefix))))
	if res == nil {
		res = calc()
		if *res == nil {
			// we can also set related cache fields
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.isSinglePrefixBlock))
			atomicStorePointer(dataLoc, unsafe.Pointer(&falseVal))
		} else {
			// we can also set related cache fields
			var isSingleBlock *bool
			if prefLen != nil && (*res).Equal(prefLen) {
				isSingleBlock = &trueVal
			} else {
				isSingleBlock = &falseVal
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.isSinglePrefixBlock))
			atomicStorePointer(dataLoc, unsafe.Pointer(isSingleBlock))

			dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&cache.minPrefix))
			atomicStorePointer(dataLoc, unsafe.Pointer(*res))
		}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.equivalentPrefix))
		atomicStorePointer(dataLoc, unsafe.Pointer(res))
	}
	return *res
}

// GetValue returns the lowest individual address division grouping in this address division grouping as an integer value.
func (grouping *addressDivisionGroupingInternal) GetValue() *big.Int {
	if grouping.hasNoDivisions() {
		return bigZero()
	}
	return bigZero().SetBytes(grouping.getBytes())
}

// GetUpperValue returns the highest individual address division grouping in this address division grouping as an integer value.
func (grouping *addressDivisionGroupingInternal) GetUpperValue() *big.Int {
	if grouping.hasNoDivisions() {
		return bigZero()
	}
	return bigZero().SetBytes(grouping.getUpperBytes())
}

// Bytes returns the lowest individual division grouping in this grouping as a byte slice.
func (grouping *addressDivisionGroupingInternal) Bytes() []byte {
	if grouping.hasNoDivisions() {
		return emptyBytes
	}
	return clone(grouping.getBytes())
}

// UpperBytes returns the highest individual division grouping in this grouping as a byte slice.
func (grouping *addressDivisionGroupingInternal) UpperBytes() []byte {
	if grouping.hasNoDivisions() {
		return emptyBytes
	}
	return clone(grouping.getUpperBytes())
}

// CopyBytes copies the value of the lowest division grouping in the range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
//
// You can use GetByteCount to determine the required array length for the bytes.
func (grouping *addressDivisionGroupingInternal) CopyBytes(bytes []byte) []byte {
	if grouping.hasNoDivisions() {
		if bytes != nil {
			return bytes[:0]
		}
		return emptyBytes
	}
	return getBytesCopy(bytes, grouping.getBytes())
}

// CopyUpperBytes copies the value of the highest division grouping in the range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
//
// You can use GetByteCount to determine the required array length for the bytes.
func (grouping *addressDivisionGroupingInternal) CopyUpperBytes(bytes []byte) []byte {
	if grouping.hasNoDivisions() {
		if bytes != nil {
			return bytes[:0]
		}
		return emptyBytes
	}
	return getBytesCopy(bytes, grouping.getUpperBytes())
}

func (grouping *addressDivisionGroupingInternal) getBytes() (bytes []byte) {
	bytes, _ = grouping.getCachedBytes(grouping.calcBytes)
	return
}

func (grouping *addressDivisionGroupingInternal) getUpperBytes() (bytes []byte) {
	_, bytes = grouping.getCachedBytes(grouping.calcBytes)
	return
}

func (grouping *addressDivisionGroupingInternal) calcBytes() (bytes, upperBytes []byte) {
	addrType := grouping.getAddrType()
	divisionCount := grouping.GetDivisionCount()
	isMultiple := grouping.isMultiple()
	if addrType.isIPv4() || addrType.isMAC() {
		bytes = make([]byte, divisionCount)
		if isMultiple {
			upperBytes = make([]byte, divisionCount)
		} else {
			upperBytes = bytes
		}
		for i := 0; i < divisionCount; i++ {
			seg := grouping.getDivision(i).ToSegmentBase()
			bytes[i] = byte(seg.GetSegmentValue())
			if isMultiple {
				upperBytes[i] = byte(seg.GetUpperSegmentValue())
			}
		}
	} else if addrType.isIPv6() {
		byteCount := divisionCount << 1
		bytes = make([]byte, byteCount)
		if isMultiple {
			upperBytes = make([]byte, byteCount)
		} else {
			upperBytes = bytes
		}
		for i := 0; i < divisionCount; i++ {
			seg := grouping.getDivision(i).ToSegmentBase()
			byteIndex := i << 1
			val := seg.GetSegmentValue()
			bytes[byteIndex] = byte(val >> 8)
			var upperVal SegInt
			if isMultiple {
				upperVal = seg.GetUpperSegmentValue()
				upperBytes[byteIndex] = byte(upperVal >> 8)
			}
			nextByteIndex := byteIndex + 1
			bytes[nextByteIndex] = byte(val)
			if isMultiple {
				upperBytes[nextByteIndex] = byte(upperVal)
			}
		}
	} else {
		byteCount := grouping.GetByteCount()
		bytes = make([]byte, byteCount)
		if isMultiple {
			upperBytes = make([]byte, byteCount)
		} else {
			upperBytes = bytes
		}
		for k, byteIndex, bitIndex := divisionCount-1, byteCount-1, BitCount(8); k >= 0; k-- {
			div := grouping.getDivision(k)
			val := div.GetDivisionValue()
			var upperVal DivInt
			if isMultiple {
				upperVal = div.GetUpperDivisionValue()
			}
			divBits := div.GetBitCount()
			for divBits > 0 {
				rbi := 8 - bitIndex
				bytes[byteIndex] |= byte(val << uint(rbi))
				val >>= uint(bitIndex)
				if isMultiple {
					upperBytes[byteIndex] |= byte(upperVal << uint(rbi))
					upperVal >>= uint(bitIndex)
				}
				if divBits < bitIndex {
					bitIndex -= divBits
					break
				} else {
					divBits -= bitIndex
					bitIndex = 8
					byteIndex--
				}
			}
		}
	}
	return
}

func (grouping *addressDivisionGroupingInternal) createNewDivisions(bitsPerDigit BitCount) ([]*AddressDivision, addrerr.IncompatibleAddressError) {
	return grouping.createNewPrefixedDivisions(bitsPerDigit, nil)
}

func (grouping *addressDivisionGroupingInternal) createNewPrefixedDivisions(bitsPerDigit BitCount, networkPrefixLength PrefixLen) ([]*AddressDivision, addrerr.IncompatibleAddressError) {
	bitCount := grouping.GetBitCount()
	var bitDivs []BitCount

	// here we divide into divisions, each with an exact number of digits.
	// Each digit takes 3 bits.  So the division bit-sizes are a multiple of 3 until the last one.

	//ipv6 octal:
	//seg bit counts: 63, 63, 2
	//ipv4 octal:
	//seg bit counts: 30, 2

	largestBitCount := BitCount(64) // uint64, size of DivInt

	largestBitCount -= largestBitCount % bitsPerDigit // round off to a multiple of 3 bits
	for {
		if bitCount <= largestBitCount {
			mod := bitCount % bitsPerDigit
			secondLast := bitCount - mod
			if secondLast > 0 {
				bitDivs = append(bitDivs, secondLast)
			}
			if mod > 0 {
				bitDivs = append(bitDivs, mod)
			}
			break
		} else {
			bitCount -= largestBitCount
			bitDivs = append(bitDivs, largestBitCount)
		}
	}

	// at this point bitDivs has our division sizes

	divCount := len(bitDivs)
	divs := make([]*AddressDivision, divCount)
	if divCount > 0 {
		//S divs[] = groupingArrayCreator.apply(divCount);
		currentSegmentIndex := 0
		seg := grouping.getDivision(currentSegmentIndex)
		segLowerVal := seg.GetDivisionValue()
		segUpperVal := seg.GetUpperDivisionValue()
		segBits := seg.GetBitCount()
		bitsSoFar := BitCount(0)

		// 2 to the x is all ones shift left x, then not, then add 1
		// so, for x == 1, 1111111 -> 1111110 -> 0000001 -> 0000010
		//radix := ^(^(0) << uint(bitsPerDigit)) + 1

		//fill up our new divisions, one by one
		for i := divCount - 1; i >= 0; i-- {

			divBitSize := bitDivs[i]
			originalDivBitSize := divBitSize
			var divLowerValue, divUpperValue uint64
			for {
				if segBits >= divBitSize { // this segment fills the remainder of this division
					diff := uint(segBits - divBitSize)
					segBits = BitCount(diff)
					segL := segLowerVal >> diff
					segU := segUpperVal >> diff

					// if the division upper bits are multiple, then the lower bits inserted must be full range
					if divLowerValue != divUpperValue {
						if segL != 0 || segU != ^(^uint64(0)<<uint(divBitSize)) {
							return nil, &incompatibleAddressError{addressError: addressError{key: "ipaddress.error.invalid.joined.ranges"}}
						}
					}

					divLowerValue |= segL
					divUpperValue |= segU

					shift := ^(^uint64(0) << diff)
					segLowerVal &= shift
					segUpperVal &= shift

					// if a segment's bits are split into two divisions, and the bits going into the first division are multi-valued,
					// then the bits going into the second division must be full range
					if segL != segU {
						if segLowerVal != 0 || segUpperVal != ^(^uint64(0)<<uint(segBits)) {
							return nil, &incompatibleAddressError{addressError: addressError{key: "ipaddress.error.invalid.joined.ranges"}}
						}
					}

					var segPrefixBits PrefixLen
					if networkPrefixLength != nil {
						segPrefixBits = getDivisionPrefixLength(originalDivBitSize, networkPrefixLength.bitCount()-bitsSoFar)
					}
					div := newRangePrefixDivision(divLowerValue, divUpperValue, segPrefixBits, originalDivBitSize)
					divs[divCount-i-1] = div
					if segBits == 0 && i > 0 {
						//get next seg
						currentSegmentIndex++
						seg = grouping.getDivision(currentSegmentIndex)
						segLowerVal = seg.getDivisionValue()
						segUpperVal = seg.getUpperDivisionValue()
						segBits = seg.getBitCount()
					}
					break
				} else {
					// if the division upper bits are multiple, then the lower bits inserted must be full range
					if divLowerValue != divUpperValue {
						if segLowerVal != 0 || segUpperVal != ^(^uint64(0)<<uint(segBits)) {
							return nil, &incompatibleAddressError{addressError: addressError{key: "ipaddress.error.invalid.joined.ranges"}}
						}
					}
					diff := uint(divBitSize - segBits)
					divLowerValue |= segLowerVal << diff
					divUpperValue |= segUpperVal << diff
					divBitSize = BitCount(diff)

					//get next seg
					currentSegmentIndex++
					seg = grouping.getDivision(currentSegmentIndex)
					segLowerVal = seg.getDivisionValue()
					segUpperVal = seg.getUpperDivisionValue()
					segBits = seg.getBitCount()
				}
			}
			bitsSoFar += originalDivBitSize
		}
	}
	return divs, nil
}

//// only needed for godoc / pkgsite

// GetBitCount returns the number of bits in each value comprising this address item.
func (grouping *addressDivisionGroupingInternal) GetBitCount() BitCount {
	return grouping.addressDivisionGroupingBase.GetBitCount()
}

// GetByteCount returns the number of bytes required for each value comprising this address item,
// rounding up if the bit count is not a multiple of 8.
func (grouping *addressDivisionGroupingInternal) GetByteCount() int {
	return grouping.addressDivisionGroupingBase.GetByteCount()
}

// GetGenericDivision returns the division at the given index as a DivisionType implementation.
func (grouping *addressDivisionGroupingInternal) GetGenericDivision(index int) DivisionType {
	return grouping.addressDivisionGroupingBase.GetGenericDivision(index)
}

// GetDivisionCount returns the number of divisions in this grouping.
func (grouping *addressDivisionGroupingInternal) GetDivisionCount() int {
	return grouping.addressDivisionGroupingBase.GetDivisionCount()
}

// IsZero returns whether this grouping matches exactly the value of zero.
func (grouping *addressDivisionGroupingInternal) IsZero() bool {
	return grouping.addressDivisionGroupingBase.IsZero()
}

// IncludesZero returns whether this grouping includes the value of zero within its range.
func (grouping *addressDivisionGroupingInternal) IncludesZero() bool {
	return grouping.addressDivisionGroupingBase.IncludesZero()
}

// IsMax returns whether this grouping matches exactly the maximum possible value, the value whose bits are all ones.
func (grouping *addressDivisionGroupingInternal) IsMax() bool {
	return grouping.addressDivisionGroupingBase.IsMax()
}

// IncludesMax returns whether this grouping includes the max value, the value whose bits are all ones, within its range.
func (grouping *addressDivisionGroupingInternal) IncludesMax() bool {
	return grouping.addressDivisionGroupingBase.IncludesMax()
}

// IsFullRange returns whether this address item represents all possible values attainable by an address item of this type.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (grouping *addressDivisionGroupingInternal) IsFullRange() bool {
	return grouping.addressDivisionGroupingBase.IsFullRange()
}

// GetSequentialBlockIndex gets the minimal division index for which all following divisions are full-range blocks.
//
// The division at this index is not a full-range block unless all divisions are full-range.
// The division at this index and all following divisions form a sequential range.
// For the full grouping to be sequential, the preceding divisions must be single-valued.
func (grouping *addressDivisionGroupingInternal) GetSequentialBlockIndex() int {
	return grouping.addressDivisionGroupingBase.GetSequentialBlockIndex()
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential address division groupings that comprise this address division grouping.
func (grouping *addressDivisionGroupingInternal) GetSequentialBlockCount() *big.Int {
	return grouping.addressDivisionGroupingBase.GetSequentialBlockCount()
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) divisions.
func (grouping *addressDivisionGroupingInternal) GetBlockCount(divisionCount int) *big.Int {
	return grouping.addressDivisionGroupingBase.GetBlockCount(divisionCount)
}

//// end needed for godoc / pkgsite

// NewDivisionGrouping creates an arbitrary grouping of divisions.
// To create address sections or addresses, use the constructors that are specific to the address version or type.
// The AddressDivision instances can be created with the NewDivision, NewRangeDivision, NewPrefixDivision or NewRangePrefixDivision functions.
func NewDivisionGrouping(divs []*AddressDivision) *AddressDivisionGrouping {
	newDivs, newPref, isMult := normalizeDivisions(divs)
	result := createGrouping(newDivs, newPref, zeroType)
	result.isMult = isMult
	return result
}

func normalizeDivisions(divs []*AddressDivision) (newDivs []*AddressDivision, newPref PrefixLen, isMultiple bool) {
	divCount := len(divs)
	newDivs = make([]*AddressDivision, 0, divCount)
	var previousDivPrefixed bool
	var bits BitCount
	for _, div := range divs {
		if div == nil || div.GetBitCount() == 0 {
			// nil divisions are divisions with zero bit-length, which we ignore
			continue
		}
		var newDiv *AddressDivision
		// The final prefix length is the minimum amongst the divisions' own prefixes
		divPrefix := div.getDivisionPrefixLength()
		divIsPrefixed := divPrefix != nil
		if previousDivPrefixed {
			if !divIsPrefixed || divPrefix.bitCount() != 0 {
				newDiv = createAddressDivision(
					div.derivePrefixed(cacheBitCount(0))) // change prefix to 0
			} else {
				newDiv = div // div prefix is already 0
			}
		} else {
			if divIsPrefixed {
				if divPrefix.bitCount() == 0 && len(newDivs) > 0 {
					// normalize boundaries by looking back
					lastDiv := newDivs[len(newDivs)-1]
					if !lastDiv.isPrefixed() {
						newDivs[len(newDivs)-1] = createAddressDivision(
							lastDiv.derivePrefixed(cacheBitCount(lastDiv.GetBitCount())))
					}
				}
				newPref = cacheBitCount(bits + divPrefix.bitCount())
				previousDivPrefixed = true
			}
			newDiv = div
		}
		newDivs = append(newDivs, newDiv)
		bits += newDiv.GetBitCount()
		isMultiple = isMultiple || newDiv.isMultiple()
	}
	return
}

// AddressDivisionGrouping objects consist of a series of AddressDivision objects, each division containing a sequential range of values.
//
// AddressDivisionGrouping objects are immutable.  This also makes them concurrency-safe.
//
// AddressDivision objects use uint64 to represent their values, so this places a cap on the size of the divisions in AddressDivisionGrouping.
//
// AddressDivisionGrouping objects are similar to address sections and addresses, except that groupings can have divisions of differing bit-length,
// including divisions that are not an exact number of bytes, whereas all segments in an address or address section must be equal bit size and an exact number of bytes.
type AddressDivisionGrouping struct {
	addressDivisionGroupingInternal
}

// Compare returns a negative integer, zero, or a positive integer if this address division grouping is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (grouping *AddressDivisionGrouping) Compare(item AddressItem) int {
	return CountComparator.Compare(grouping, item)
}

// CompareSize compares the counts of two items, the number of individual items represented in each.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether this grouping represents more individual items than another.
//
// CompareSize returns a positive integer if this address division grouping has a larger count than the item given, zero if they are the same, or a negative integer if the other has a larger count.
func (grouping *AddressDivisionGrouping) CompareSize(other AddressItem) int {
	if grouping == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return grouping.compareSize(other)
}

// GetCount returns the count of possible distinct values for this division grouping.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (grouping *AddressDivisionGrouping) GetCount() *big.Int {
	if grouping == nil {
		return bigZero()
	}
	return grouping.getCount()
}

// IsMultiple returns whether this grouping represents multiple values rather than a single value
func (grouping *AddressDivisionGrouping) IsMultiple() bool {
	return grouping != nil && grouping.isMultiple()
}

// IsPrefixed returns whether this grouping has an associated prefix length.
func (grouping *AddressDivisionGrouping) IsPrefixed() bool {
	if grouping == nil {
		return false
	}
	return grouping.isPrefixed()
}

// CopySubDivisions copies the existing divisions from the given start index until but not including the division at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of divisions copied.
func (grouping *AddressDivisionGrouping) CopySubDivisions(start, end int, divs []*AddressDivision) (count int) {
	return grouping.copySubDivisions(start, end, divs)
}

// CopyDivisions copies the existing divisions from the given start index until but not including the division at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of divisions copied.
func (grouping *AddressDivisionGrouping) CopyDivisions(divs []*AddressDivision) (count int) {
	return grouping.copyDivisions(divs)
}

// GetDivisionStrings returns a slice containing each string returned from the String method of each division in the grouping.
func (grouping *AddressDivisionGrouping) GetDivisionStrings() []string {
	if grouping == nil {
		return nil
	}
	return grouping.getDivisionStrings()
}

// IsAdaptiveZero returns true if this is an adaptive zero grouping.
// The adaptive zero grouping, produced by zero sections like IPv4AddressSection{} or AddressDivisionGrouping{}, can represent a zero-length section of any address type.
// It is not considered equal to constructions of specific zero length sections or groupings like NewIPv4Section(nil) which can only represent a zero-length section of a single address type.
func (grouping *AddressDivisionGrouping) IsAdaptiveZero() bool {
	return grouping != nil && grouping.matchesZeroGrouping()
}

// IsSectionBase returns true if this address division grouping originated as an address section.  If so, use ToSectionBase to convert back to the section type.
func (grouping *AddressDivisionGrouping) IsSectionBase() bool {
	return grouping != nil && grouping.isAddressSection()
}

// IsIP returns true if this address division grouping originated as an IPv4 or IPv6 section, or a zero-length IP section.  If so, use ToIP to convert back to the IP-specific type.
func (grouping *AddressDivisionGrouping) IsIP() bool {
	return grouping.ToSectionBase().IsIP()
}

// IsIPv4 returns true if this grouping originated as an IPv4 section.  If so, use ToIPv4 to convert back to the IPv4-specific type.
func (grouping *AddressDivisionGrouping) IsIPv4() bool {
	return grouping.ToSectionBase().IsIPv4()
}

// IsIPv6 returns true if this grouping originated as an IPv6 section.  If so, use ToIPv6 to convert back to the IPv6-specific type.
func (grouping *AddressDivisionGrouping) IsIPv6() bool {
	return grouping.ToSectionBase().IsIPv6()
}

// IsMixedIPv6v4 returns true if this grouping originated as a mixed IPv6-IPv4 grouping.  If so, use ToMixedIPv6v4 to convert back to the more specific grouping type.
func (grouping *AddressDivisionGrouping) IsMixedIPv6v4() bool {
	return grouping != nil && grouping.matchesIPv6v4MixedGroupingType()
}

// IsMAC returns true if this grouping originated as a MAC section.  If so, use ToMAC to convert back to the MAC-specific type.
func (grouping *AddressDivisionGrouping) IsMAC() bool {
	return grouping.ToSectionBase().IsMAC()
}

// ToSectionBase converts to an address section if this grouping originated as an address section.
// Otherwise, the result will be nil.
//
// ToSectionBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *AddressDivisionGrouping) ToSectionBase() *AddressSection {
	if grouping == nil || !grouping.isAddressSection() {
		return nil
	}
	return (*AddressSection)(unsafe.Pointer(grouping))
}

// ToMixedIPv6v4 converts to a mixed IPv6/4 address section if this grouping originated as a mixed IPv6/4 address section.
// Otherwise, the result will be nil.
//
// ToMixedIPv6v4 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *AddressDivisionGrouping) ToMixedIPv6v4() *IPv6v4MixedAddressGrouping {
	if grouping.matchesIPv6v4MixedGroupingType() {
		return (*IPv6v4MixedAddressGrouping)(grouping)
	}
	return nil
}

// ToIP converts to an IPAddressSection if this grouping originated as an IPv4 or IPv6 section, or an implicitly zero-valued IP section.
// If not, ToIP returns nil.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *AddressDivisionGrouping) ToIP() *IPAddressSection {
	return grouping.ToSectionBase().ToIP()
}

// ToIPv6 converts to an IPv6AddressSection if this grouping originated as an IPv6 section.
// If not, ToIPv6 returns nil.
//
// ToIPv6 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *AddressDivisionGrouping) ToIPv6() *IPv6AddressSection {
	return grouping.ToSectionBase().ToIPv6()
}

// ToIPv4 converts to an IPv4AddressSection if this grouping originated as an IPv4 section.
// If not, ToIPv4 returns nil.
//
// ToIPv4 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *AddressDivisionGrouping) ToIPv4() *IPv4AddressSection {
	return grouping.ToSectionBase().ToIPv4()
}

// ToMAC converts to a MACAddressSection if this grouping originated as a MAC section.
// If not, ToMAC returns nil.
//
// ToMAC can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *AddressDivisionGrouping) ToMAC() *MACAddressSection {
	return grouping.ToSectionBase().ToMAC()
}

// ToDivGrouping is an identity method.
//
// ToDivGrouping can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *AddressDivisionGrouping) ToDivGrouping() *AddressDivisionGrouping {
	return grouping
}

// GetDivision returns the division at the given index.
func (grouping *AddressDivisionGrouping) GetDivision(index int) *AddressDivision {
	return grouping.getDivision(index)
}

// ForEachDivision visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// ForEachDivision returns the number of visited segments.
func (grouping *AddressDivisionGrouping) ForEachDivision(consumer func(divisionIndex int, division *AddressDivision) (stop bool)) int {
	divArray := grouping.getDivArray()
	if divArray != nil {
		for i, div := range divArray {
			if consumer(i, div) {
				return i + 1
			}
		}
	}
	return len(divArray)
}

// String implements the [fmt.Stringer] interface.
// It returns "<nil>" if the receiver is a nil pointer.
// It returns the normalized string provided by ToNormalizedString if this grouping originated as an address section.
// Otherwise, the string is printed like a slice, with each division converted to a string by its own String method (like "[ div0 div1 ... ]").
func (grouping *AddressDivisionGrouping) String() string {
	if grouping == nil {
		return nilString()
	}
	return grouping.toString()
}
