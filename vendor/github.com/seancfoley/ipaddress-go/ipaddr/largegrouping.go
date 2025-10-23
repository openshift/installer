//
// Copyright 2022-2024 Sean C Foley
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
)

func createLargeGrouping(divs []*IPAddressLargeDivision) *IPAddressLargeDivisionGrouping {
	addrType := zeroType
	grouping := &IPAddressLargeDivisionGrouping{
		largeDivisionGroupingInternal{
			addressDivisionGroupingBase: addressDivisionGroupingBase{
				divisions: largeDivArray(divs),
				addrType:  addrType,
				cache:     &valueCache{},
			},
		},
	}
	assignStringCache(&grouping.addressDivisionGroupingBase, addrType)
	return grouping
}

type largeDivisionGroupingInternal struct {
	addressDivisionGroupingBase
}

func (grouping *largeDivisionGroupingInternal) getDivArray() largeDivArray {
	if divsArray := grouping.divisions; divsArray != nil {
		return divsArray.(largeDivArray)
	}
	return nil
}

func (grouping *largeDivisionGroupingInternal) getDivisionCount() int {
	if divArray := grouping.getDivArray(); divArray != nil {
		return divArray.getDivisionCount()
	}
	return 0
}

// getDivision returns the division or panics if the index is negative or too large
func (grouping *largeDivisionGroupingInternal) getDivision(index int) *IPAddressLargeDivision {
	return grouping.getDivArray()[index]
}

func (grouping *largeDivisionGroupingInternal) initMultiple() {
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

// divisions are printed like slices of *IPAddressLargeDivision (which are Stringers) with division separated by spaces and enclosed in square brackets,
// sections are printed like addresses with segments separated by segment separators
func (grouping largeDivisionGroupingInternal) Format(state fmt.State, verb rune) {
	arr := grouping.initDivs().getDivArray()
	if len(arr) == 0 {
		return
	}
	s := flagsFromState(state, verb)
	_, _ = state.Write([]byte(fmt.Sprintf(s, arr)))
}

func (grouping *largeDivisionGroupingInternal) toString() string {
	return fmt.Sprint(grouping.initDivs().getDivArray())
}

var zeroLargeGrouping = createLargeGrouping(zeroLargeDivs)

func (grouping *largeDivisionGroupingInternal) initDivs() *largeDivisionGroupingInternal {
	if grouping.divisions == nil {
		return &zeroLargeGrouping.largeDivisionGroupingInternal
	}
	return grouping
}

func (grouping *largeDivisionGroupingInternal) getBytes() (bytes []byte) {
	bytes, _ = grouping.getCachedBytes(grouping.calcBytes)
	return
}

func (grouping *largeDivisionGroupingInternal) getUpperBytes() (bytes []byte) {
	_, bytes = grouping.getCachedBytes(grouping.calcBytes)
	return
}

func (grouping *largeDivisionGroupingInternal) calcBytes() (bytes, upperBytes []byte) {
	divisionCount := grouping.GetDivisionCount()
	isMultiple := grouping.isMultiple()
	byteCount := grouping.GetByteCount()
	bytes = make([]byte, byteCount)
	if isMultiple {
		upperBytes = make([]byte, byteCount)
	} else {
		upperBytes = bytes
	}
	// for each division in reverse order
	for k, byteIndex, bitIndex := divisionCount-1, byteCount-1, BitCount(8); k >= 0; k-- {
		div := grouping.getDivision(k)
		bigBytes := div.getValue().Bytes()
		var bigUpperBytes []byte
		if isMultiple {
			bigUpperBytes = div.getUpperValue().Bytes()
		}

		// for each 64 bits of the division in reverse order
		for totalDivBits := div.GetBitCount(); totalDivBits > 0; totalDivBits -= 64 {

			// grab those 64 bits (from bigBytes and bigUpperBytes) and put them in val and upperVal
			divBits := min(totalDivBits, 64)
			var divBytes []byte
			var val, upperVal uint64
			if len(bigBytes) > 8 {
				byteLen := len(bigBytes) - 8
				divBytes = bigBytes[byteLen:]
				bigBytes = bigBytes[:byteLen]
			} else {
				divBytes = bigBytes
				bigBytes = nil
			}
			for _, b := range divBytes {
				val = (val << 8) | uint64(b)
			}
			if isMultiple {
				var divUpperBytes []byte
				if len(upperBytes) > 8 {
					byteLen := len(bigUpperBytes) - 8
					divUpperBytes = bigBytes[byteLen:]
					bigUpperBytes = bigBytes[:byteLen]
				} else {
					divUpperBytes = bigUpperBytes
					bigUpperBytes = nil
				}
				for _, b := range divUpperBytes {
					upperVal = (upperVal << 8) | uint64(b)
				}
			}

			// insert the 64 bits into the  bytes slice
			for divBits > 0 {
				rbi := 8 - bitIndex
				bytes[byteIndex] |= byte(val << uint(rbi))
				val >>= uint(bitIndex)
				if isMultiple {
					upperBytes[byteIndex] |= byte(upperVal << uint(rbi))
					upperVal >>= uint(bitIndex)
				}
				if divBits < bitIndex {
					// bitIndex is the index into the last copied byte that was already occupied previously
					// so here we were able to copy all the bits and there was still space left over
					bitIndex -= divBits
					break
				} else {
					// we used up all the space available
					// if we also copied all the bits, then divBits will be assigned zero
					// otherwise it will have the number of bits still left to copy
					divBits -= bitIndex
					bitIndex = 8
					byteIndex--
				}
			}
		}
	}
	return
}

// CopyBytes copies the value of the lowest division grouping in the range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
//
// You can use GetByteCount to determine the required array length for the bytes.
func (grouping *largeDivisionGroupingInternal) CopyBytes(bytes []byte) []byte {
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
func (grouping *largeDivisionGroupingInternal) CopyUpperBytes(bytes []byte) []byte {
	if grouping.hasNoDivisions() {
		if bytes != nil {
			return bytes[:0]
		}
		return emptyBytes
	}
	return getBytesCopy(bytes, grouping.getUpperBytes())
}

// Bytes returns the lowest individual division grouping in this grouping as a byte slice.
func (grouping *largeDivisionGroupingInternal) Bytes() []byte {
	if grouping.hasNoDivisions() {
		return emptyBytes
	}
	return clone(grouping.getBytes())
}

// UpperBytes returns the highest individual division grouping in this grouping as a byte slice.
func (grouping *largeDivisionGroupingInternal) UpperBytes() []byte {
	if grouping.hasNoDivisions() {
		return emptyBytes
	}
	return clone(grouping.getUpperBytes())
}

// GetValue returns the lowest individual address division grouping in this address division grouping as an integer value.
func (grouping *largeDivisionGroupingInternal) GetValue() *big.Int {
	res := big.Int{}
	if grouping.hasNoDivisions() {
		return &res
	}
	return res.SetBytes(grouping.getBytes())
}

// GetUpperValue returns the highest individual address division grouping in this address division grouping as an integer value.
func (grouping *largeDivisionGroupingInternal) GetUpperValue() *big.Int {
	res := big.Int{}
	if grouping.hasNoDivisions() {
		return &res
	}
	return res.SetBytes(grouping.getUpperBytes())
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this division grouping matches the block of addresses for that prefix.
//
// If no such prefix exists, GetPrefixLenForSingleBlock returns nil.
//
// If this division grouping represents a single value, returns the bit length.
func (grouping *largeDivisionGroupingInternal) GetPrefixLenForSingleBlock() PrefixLen {
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

// GetMinPrefixLenForBlock returns the smallest prefix length such that this grouping includes the block of all values for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this grouping represents a single value, this returns the bit count.
func (grouping *largeDivisionGroupingInternal) GetMinPrefixLenForBlock() BitCount {
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

// IsPrefixBlock returns whether this division grouping has a prefix length and includes the block associated with its prefix length.
// If the prefix length matches the bit count, this returns true.
//
// This is different from ContainsPrefixBlock in that this method returns
// false if the series has no prefix length, or a prefix length that differs from a prefix length for which ContainsPrefixBlock returns true.
func (grouping *largeDivisionGroupingInternal) IsPrefixBlock() bool {
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
func (grouping *largeDivisionGroupingInternal) IsSinglePrefixBlock() bool {
	calc := func() bool {
		prefLen := grouping.getPrefixLen()
		return prefLen != nil && grouping.ContainsSinglePrefixBlock(prefLen.bitCount())
	}
	return cacheIsSinglePrefixBlock(grouping.cache, grouping.getPrefixLen(), calc)
}

// ContainsPrefixBlock returns whether the values of this item contains the block of values for the given prefix length.
//
// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values in this item for the given prefix length makes no difference.
//
// Use GetMinPrefixLenForBlock to determine the smallest prefix length for which this method returns true.
func (grouping *largeDivisionGroupingInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	prefixLen = checkSubnet(grouping, prefixLen)
	divisionCount := grouping.GetDivisionCount()
	var prevBitCount BitCount
	for i := 0; i < divisionCount; i++ {
		division := grouping.getDivision(i)
		bitCount := division.GetBitCount()
		totalBitCount := bitCount + prevBitCount
		if prefixLen < totalBitCount {
			divPrefixLen := prefixLen - prevBitCount
			if !division.ContainsPrefixBlock(divPrefixLen) {
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
func (grouping *largeDivisionGroupingInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
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

// copySubDivisions copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (grouping *largeDivisionGroupingInternal) copySubDivisions(start, end int, divs []*IPAddressLargeDivision) (count int) {
	if divArray := grouping.getDivArray(); divArray != nil {
		start, end, targetIndex := adjust1To1Indices(start, end, grouping.GetDivisionCount(), len(divs))
		return divArray.copySubDivisions(start, end, divs[targetIndex:])
	}
	return
}

// copyDivisions copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (grouping *largeDivisionGroupingInternal) copyDivisions(divs []*IPAddressLargeDivision) (count int) {
	if divArray := grouping.getDivArray(); divArray != nil {
		return divArray.copyDivisions(divs)
	}
	return
}

// NewIPAddressLargeDivGrouping creates an arbitrary grouping of divisions of arbitrary size, each division can have an arbitrarily large bit-length.
// To create address sections or addresses, use the constructors that are specific to the address version or type.
// The IPAddressLargeDivision instances can be created with the NewLargeIPDivision, NewLargeIPRangeDivision, NewLargeIPPrefixDivision, NewLargeIPRangePrefixDivision functions.
func NewIPAddressLargeDivGrouping(divs []*IPAddressLargeDivision) *IPAddressLargeDivisionGrouping {
	// We do not check for prefix subnet because an explicit prefix length must be supplied for that
	newDivs, newPref, isMult := normalizeLargeDivisions(divs)
	result := createLargeGrouping(newDivs)
	result.isMult = isMult
	result.prefixLength = newPref
	return result
}

func normalizeLargeDivisions(divs []*IPAddressLargeDivision) (newDivs []*IPAddressLargeDivision, newPref PrefixLen, isMultiple bool) {
	divCount := len(divs)
	newDivs = make([]*IPAddressLargeDivision, 0, divCount)
	var previousDivPrefixed bool
	var bits BitCount
	for _, div := range divs {
		if div == nil || div.GetBitCount() == 0 {
			// nil divisions are divisions with zero bit-length, which we ignore
			continue
		}
		var newDiv *IPAddressLargeDivision
		// The final prefix length is the minimum amongst the divisions' own prefixes
		divPrefix := div.getDivisionPrefixLength()
		divIsPrefixed := divPrefix != nil
		if previousDivPrefixed {
			if !divIsPrefixed || divPrefix.bitCount() != 0 {
				newDiv = createLargeAddressDiv(div.derivePrefixed(cacheBitCount(0)), div.getDefaultRadix()) // change prefix to 0
			} else {
				newDiv = div // div prefix is already 0
			}
		} else {
			if divIsPrefixed {
				if divPrefix.bitCount() == 0 && len(newDivs) > 0 {
					// normalize boundaries by looking back
					lastDiv := newDivs[len(newDivs)-1]
					if !lastDiv.IsPrefixed() {
						newDivs[len(newDivs)-1] = createLargeAddressDiv(
							lastDiv.derivePrefixed(cacheBitCount(lastDiv.GetBitCount())), div.getDefaultRadix())
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

type IPAddressLargeDivisionGrouping struct {
	largeDivisionGroupingInternal
}

// GetCount returns the count of possible distinct values for this division grouping.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (grouping *IPAddressLargeDivisionGrouping) GetCount() *big.Int {
	if !grouping.isMultiple() {
		return bigOne()
	}
	return grouping.addressDivisionGroupingBase.getCount()
}

// IsMultiple returns whether this grouping represents multiple values rather than a single value.
func (grouping *IPAddressLargeDivisionGrouping) IsMultiple() bool {
	return grouping != nil && grouping.isMultiple()
}

// Compare returns a negative integer, zero, or a positive integer if this address division grouping is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (grouping *IPAddressLargeDivisionGrouping) Compare(item AddressItem) int {
	return CountComparator.Compare(grouping, item)
}

// CompareSize compares the counts of two items, the number of individual values within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one represents more individual values than another.
//
// CompareSize returns a positive integer if this division has a larger count than the item given, zero if they are the same, or a negative integer if the other has a larger count.
func (grouping *IPAddressLargeDivisionGrouping) CompareSize(other AddressItem) int {
	if grouping == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return compareCount(grouping, other)
	//return grouping.compareSize(other)
}

// String implements the [fmt.Stringer] interface.
// It returns "<nil>" if the receiver is a nil pointer.
// Otherwise, the string is printed like a slice, with each division converted to a string by its own String method (like "[ div0 div1 ... ]").
func (grouping *IPAddressLargeDivisionGrouping) String() string {
	if grouping == nil {
		return nilString()
	}
	return grouping.toString()
}

// IsPrefixed returns whether this division grouping has an associated prefix length.
// If so, the prefix length is given by GetPrefixLen.
func (grouping *IPAddressLargeDivisionGrouping) IsPrefixed() bool {
	if grouping == nil {
		return false
	}
	return grouping.isPrefixed()
}

// GetDivision returns the division at the given index.
func (grouping *IPAddressLargeDivisionGrouping) GetDivision(index int) *IPAddressLargeDivision {
	return grouping.getDivision(index)
}

// ForEachDivision visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// ForEachDivision returns the number of visited segments.
func (grouping *IPAddressLargeDivisionGrouping) ForEachDivision(consumer func(divisionIndex int, division *IPAddressLargeDivision) (stop bool)) int {
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

func (grouping *IPAddressLargeDivisionGrouping) isNil() bool {
	return grouping == nil
}

// CopySubDivisions copies the existing divisions from the given start index until but not including the division at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of divisions copied.
func (grouping *IPAddressLargeDivisionGrouping) CopySubDivisions(start, end int, divs []*IPAddressLargeDivision) (count int) {
	return grouping.copySubDivisions(start, end, divs)
}

// CopyDivisions copies the existing divisions from the given start index until but not including the division at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of divisions copied.
func (grouping *IPAddressLargeDivisionGrouping) CopyDivisions(divs []*IPAddressLargeDivision) (count int) {
	return grouping.copyDivisions(divs)
}
