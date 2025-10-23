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
	"math/bits"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

// SegInt is an integer type for holding generic address segment values.  It is at least as large as all address segment values: [IPv6SegInt], [IPv4SegInt], [MACSegInt].
type SegInt = uint32 // must be at least uint16 to handle IPv6, at least 32 to handle single segment IPv4, and no larger than 64 because we use bits.TrailingZeros64.  IP address segment code uses bits.TrailingZeros32 and bits.LeadingZeros32, so it cannot be larger than 32.

const SegIntSize = 32 // must match the bit count of SegInt

type SegIntCount = uint64 // must be able to hold: (max value of SegInt) + 1

type segderiver interface {
	// deriveNew produces a new segment with the same bit count as the old
	deriveNewMultiSeg(val, upperVal SegInt, prefLen PrefixLen) divisionValues

	// deriveNew produces a new segment with the same bit count as the old
	deriveNewSeg(val SegInt, prefLen PrefixLen) divisionValues
}

type segmentValues interface {
	// getSegmentValue gets the lower value for a segment
	getSegmentValue() SegInt

	// getUpperSegmentValue gets the upper value for a segment
	getUpperSegmentValue() SegInt
}

// compareSegInt returns a negative number, 0 or a positive number if integer one is less than, equal to, or greater than integer two.
func compareSegInt(one, two SegInt) int {
	if one < two {
		return -1
	} else if one > two {
		return 1
	}
	return 0
}

func createAddressSegment(vals divisionValues) *AddressSegment {
	return &AddressSegment{
		addressSegmentInternal{
			addressDivisionInternal{
				addressDivisionBase{
					vals,
				},
			},
		},
	}
}

type addressSegmentInternal struct {
	addressDivisionInternal
}

func (seg *addressSegmentInternal) sameTypeContains(otherSeg *AddressSegment) bool {
	return otherSeg.GetSegmentValue() >= seg.GetSegmentValue() &&
		otherSeg.GetUpperSegmentValue() <= seg.GetUpperSegmentValue()
}

func (seg *addressSegmentInternal) contains(other AddressSegmentType) bool {
	if other == nil {
		return true
	}
	otherSeg := other.ToSegmentBase()
	if seg.toAddressSegment() == otherSeg || otherSeg == nil {
		return true
	} else if matchesStructure, _ := seg.matchesStructure(other); matchesStructure {
		return seg.sameTypeContains(otherSeg)
	}
	return false
}

func (seg *addressSegmentInternal) sameTypeOverlaps(otherSeg *AddressSegment) bool {
	return otherSeg.GetSegmentValue() <= seg.GetUpperSegmentValue() &&
		otherSeg.GetUpperSegmentValue() >= seg.GetSegmentValue()
}

func (seg *addressSegmentInternal) overlaps(other AddressSegmentType) bool {
	if other == nil {
		return true
	}
	otherSeg := other.ToSegmentBase()
	if seg.toAddressSegment() == otherSeg || otherSeg == nil {
		return true
	} else if matchesStructure, _ := seg.matchesStructure(other); matchesStructure {
		return seg.sameTypeOverlaps(otherSeg)
	}
	return false
}

func (seg *addressSegmentInternal) equal(other AddressSegmentType) bool {
	if other == nil || other.ToSegmentBase() == nil {
		return false
	}
	if seg.isMultiple() {
		if other.IsMultiple() {
			matches, _ := seg.matchesStructure(other)
			otherDivision := other.ToSegmentBase()
			return matches && segValsSame(seg.getSegmentValue(), otherDivision.getSegmentValue(),
				seg.getUpperSegmentValue(), otherDivision.getUpperSegmentValue())
		} else {
			return false
		}
	} else if other.IsMultiple() {
		return false
	}
	matches, _ := seg.matchesStructure(other)
	otherDivision := other.ToSegmentBase()
	return matches && segValSame(seg.GetSegmentValue(), otherDivision.GetSegmentValue())
}

func (seg *addressSegmentInternal) equalsSegment(other *AddressSegment) bool {
	matchesStructure, _ := seg.matchesStructure(other)
	return matchesStructure && seg.sameTypeEquals(other)
}

func (seg *addressSegmentInternal) sameTypeEquals(other *AddressSegment) bool {
	if seg.isMultiple() {
		return other.isMultiple() && segValsSame(seg.getSegmentValue(), other.getSegmentValue(),
			seg.getUpperSegmentValue(), other.getUpperSegmentValue())
	}
	return !other.isMultiple() && seg.getSegmentValue() == other.getSegmentValue()
}

// PrefixContains returns whether the prefix values in the prefix of the given segment are also prefix values in this segment.
// It returns whether the prefix of this segment contains the prefix of the given segment.
func (seg *addressSegmentInternal) PrefixContains(other AddressSegmentType, prefixLength BitCount) bool {
	prefixLength = checkBitCount(prefixLength, seg.GetBitCount())
	shift := seg.GetBitCount() - prefixLength
	if shift <= 0 {
		return seg.contains(other)
	}
	return (other.GetSegmentValue()>>uint(shift)) >= (seg.GetSegmentValue()>>uint(shift)) &&
		(other.GetUpperSegmentValue()>>uint(shift)) <= (seg.GetUpperSegmentValue()>>uint(shift))
}

// PrefixEqual returns whether the prefix bits of this segment match the same bits of the given segment.
// It returns whether the two segments share the same range of prefix values using the given prefix length.
func (seg *addressSegmentInternal) PrefixEqual(other AddressSegmentType, prefixLength BitCount) bool {
	prefixLength = checkBitCount(prefixLength, seg.GetBitCount())
	shift := seg.GetBitCount() - prefixLength
	if shift <= 0 {
		return seg.GetSegmentValue() == other.GetSegmentValue() && seg.GetUpperSegmentValue() == other.GetUpperSegmentValue()
	}
	return (other.GetSegmentValue()>>uint(shift)) == (seg.GetSegmentValue()>>uint(shift)) &&
		(other.GetUpperSegmentValue()>>uint(shift)) == (seg.GetUpperSegmentValue()>>uint(shift))
}

func (seg *addressSegmentInternal) toAddressSegment() *AddressSegment {
	return (*AddressSegment)(unsafe.Pointer(seg))
}

// GetSegmentValue returns the lower value of the segment value range.
func (seg *addressSegmentInternal) GetSegmentValue() SegInt {
	vals := seg.divisionValues
	if vals == nil {
		return 0
	}
	return vals.getSegmentValue()
}

// GetUpperSegmentValue returns the upper value of the segment value range.
func (seg *addressSegmentInternal) GetUpperSegmentValue() SegInt {
	vals := seg.divisionValues
	if vals == nil {
		return 0
	}
	return vals.getUpperSegmentValue()
}

// Matches returns true if the segment range matches the given single value.
func (seg *addressSegmentInternal) Matches(value SegInt) bool {
	return seg.matches(DivInt(value))
}

// MatchesWithMask applies the mask to this segment and then compares the result with the given value,
// returning true if the range of the resulting segment matches that single value.
func (seg *addressSegmentInternal) MatchesWithMask(value, mask SegInt) bool {
	return seg.matchesWithMask(DivInt(value), DivInt(mask))
}

// MatchesValsWithMask applies the mask to this segment and then compares the result with the given values,
// returning true if the range of the resulting segment matches the given range.
func (seg *addressSegmentInternal) MatchesValsWithMask(lowerValue, upperValue, mask SegInt) bool {
	return seg.matchesValsWithMask(DivInt(lowerValue), DivInt(upperValue), DivInt(mask))
}

// GetPrefixCountLen returns the count of the number of distinct prefix values for the given prefix length in the range of values of this segment.
func (seg *addressSegmentInternal) GetPrefixCountLen(segmentPrefixLength BitCount) *big.Int {
	return bigZero().SetUint64(seg.GetPrefixValueCountLen(segmentPrefixLength))
}

// GetPrefixValueCountLen returns the same value as GetPrefixCountLen as an integer.
func (seg *addressSegmentInternal) GetPrefixValueCountLen(segmentPrefixLength BitCount) SegIntCount {
	return getPrefixValueCount(seg.toAddressSegment(), segmentPrefixLength)
}

// GetValueCount returns the same value as GetCount as an integer.
func (seg *addressSegmentInternal) GetValueCount() SegIntCount {
	return uint64(seg.GetUpperSegmentValue()-seg.GetSegmentValue()) + 1
}

// GetMaxValue gets the maximum possible value for this type or version of segment, determined by the number of bits.
//
// For the highest range value of this particular segment, use GetUpperSegmentValue.
func (seg *addressSegmentInternal) GetMaxValue() SegInt {
	return ^(^SegInt(0) << uint(seg.GetBitCount()))
}

// TestBit returns true if the bit in the lower value of this segment at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this section.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (seg *addressSegmentInternal) TestBit(n BitCount) bool {
	value := seg.GetSegmentValue()
	if n < 0 || n >= seg.GetBitCount() {
		panic("invalid bit index")
	}
	return (value & (1 << uint(n))) != 0
}

// IsOneBit returns true if the bit in the lower value of this segment at the given index is 1, where index 0 refers to the most significant bit.
// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this item.
func (seg *addressSegmentInternal) IsOneBit(segmentBitIndex BitCount) bool {
	value := seg.GetSegmentValue()
	bitCount := seg.GetBitCount()
	if segmentBitIndex < 0 || segmentBitIndex >= seg.GetBitCount() {
		panic("invalid bit index")
	}
	return (value & (1 << uint(bitCount-(segmentBitIndex+1)))) != 0
}

func (seg *addressSegmentInternal) getLower() *AddressSegment {
	if !seg.isMultiple() {
		return seg.toAddressSegment()
	}
	vals := seg.divisionValues
	var newVals divisionValues
	if vals != nil {
		newVals = seg.deriveNewMultiSeg(seg.GetSegmentValue(), seg.GetSegmentValue(), seg.getDivisionPrefixLength())
	}
	return createAddressSegment(newVals)
}

func (seg *addressSegmentInternal) getUpper() *AddressSegment {
	if !seg.isMultiple() {
		return seg.toAddressSegment()
	}
	vals := seg.divisionValues
	var newVals divisionValues
	if vals != nil {
		newVals = seg.deriveNewMultiSeg(seg.GetUpperSegmentValue(), seg.GetUpperSegmentValue(), seg.getDivisionPrefixLength())
	}
	return createAddressSegment(newVals)
}

func (seg *addressSegmentInternal) withoutPrefixLen() *AddressSegment {
	if seg.isPrefixed() {
		return createAddressDivision(seg.derivePrefixed(nil)).ToSegmentBase()
	}
	return seg.toAddressSegment()
}

func (seg *addressSegmentInternal) getDefaultSegmentWildcardString() string {
	return SegmentWildcardStr
}

func (seg *addressSegmentInternal) iterator() Iterator[*AddressSegment] {
	return seg.segmentIterator(seg.getDivisionPrefixLength(), false, false)
}

func (seg *addressSegmentInternal) identityIterator() Iterator[*AddressSegment] {
	return &singleSegmentIterator{original: seg.toAddressSegment()}
}

func (seg *addressSegmentInternal) prefixBlockIterator() Iterator[*AddressSegment] {
	return seg.segmentIterator(seg.getDivisionPrefixLength(), true, true)
}

func (seg *addressSegmentInternal) prefixedBlockIterator(segPrefLen BitCount) Iterator[*AddressSegment] {
	return seg.segmentIterator(cacheBitCount(segPrefLen), true, true)
}

func (seg *addressSegmentInternal) prefixIterator() Iterator[*AddressSegment] {
	return seg.segmentIterator(seg.getDivisionPrefixLength(), true, false)
}

func (seg *addressSegmentInternal) prefixedIterator(segPrefLen BitCount) Iterator[*AddressSegment] {
	return seg.segmentIterator(cacheBitCount(segPrefLen), true, false)
}

func (seg *addressSegmentInternal) segmentIterator(segPrefLen PrefixLen, isPrefixIterator, isBlockIterator bool) Iterator[*AddressSegment] {
	vals := seg.divisionValues
	if vals == nil {
		return segIterator(seg,
			0,
			0,
			0,
			nil,
			nil,
			false,
			false,
		)
	}
	return segIterator(seg,
		seg.getSegmentValue(),
		seg.getUpperSegmentValue(),
		seg.getBitCount(),
		vals,
		segPrefLen,
		isPrefixIterator,
		isBlockIterator,
	)
}

// GetLeadingBitCount returns the number of consecutive leading one or zero bits.
// If ones is true, returns the number of consecutive leading one bits.
// Otherwise, returns the number of consecutive leading zero bits.
//
// This method applies only to the lower value of the range if this segment represents multiple values.
func (seg *addressSegmentInternal) GetLeadingBitCount(ones bool) BitCount {
	extraLeading := 32 - seg.GetBitCount()
	val := seg.GetSegmentValue()
	if ones {
		//leading ones
		return BitCount(bits.LeadingZeros32(uint32(^val&seg.GetMaxValue()))) - extraLeading
	}
	// leading zeros
	return BitCount(bits.LeadingZeros32(uint32(val))) - extraLeading
}

// GetTrailingBitCount returns the number of consecutive trailing one or zero bits.
// If ones is true, returns the number of consecutive trailing zero bits.
// Otherwise, returns the number of consecutive trailing one bits.
//
// This method applies only to the lower value of the range if this segment represents multiple values.
func (seg *addressSegmentInternal) GetTrailingBitCount(ones bool) BitCount {
	val := seg.GetSegmentValue()
	if ones {
		// trailing ones
		return BitCount(bits.TrailingZeros32(uint32(^val)))
	}
	//trailing zeros
	bitCount := uint(seg.GetBitCount())
	return BitCount(bits.TrailingZeros32(uint32(val | (1 << bitCount))))
}

// GetSegmentNetworkMask returns a value comprising the same number of total bits as the bit-length of this segment,
// the value that is all one-bits for the given number of bits followed by all zero-bits.
func (seg *addressSegmentInternal) GetSegmentNetworkMask(networkBits BitCount) SegInt {
	bitCount := seg.GetBitCount()
	networkBits = checkBitCount(networkBits, bitCount)
	return seg.GetMaxValue() & (^SegInt(0) << uint(bitCount-networkBits))
}

// GetSegmentHostMask returns a value comprising the same number of total bits as the bit-length of this segment,
// the value that is all zero-bits for the given number of bits followed by all one-bits.
func (seg *addressSegmentInternal) GetSegmentHostMask(networkBits BitCount) SegInt {
	bitCount := seg.GetBitCount()
	networkBits = checkBitCount(networkBits, bitCount)
	return ^(^SegInt(0) << uint(bitCount-networkBits))
}

var (
	// wildcards differ, for divs we use only range since div size not implicit, here we use both range and *
	hexParamsSeg     = new(addrstr.IPStringOptionsBuilder).SetRadix(16).SetSegmentStrPrefix(HexPrefix).ToOptions()
	decimalParamsSeg = new(addrstr.IPStringOptionsBuilder).SetRadix(10).ToOptions()
)

// We do not need to "override" ToNormalizedString() and ToHexString(bool) because neither prints leading zeros according to bit count, so zero-segments of type IPv4/IPv6/MAC are printed consistently

// ToNormalizedString produces a string that is consistent for all address segments of the same type and version.
// IPv4 segments use base 10, while other segment types use base 16.
func (seg *addressSegmentInternal) ToNormalizedString() string {
	stringer := func() string {
		switch seg.getDefaultTextualRadix() {
		case 10:
			return seg.toStringOpts(decimalParamsSeg)
		default:
			return seg.toStringOpts(macCompressedParams)
		}
	}
	if seg.divisionValues != nil {
		if cache := seg.getCache(); cache != nil {
			return cacheStr(&cache.cachedNormalizedString, stringer)
		}
	}
	return stringer()
}

// ToHexString writes this address segment as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// For segments, the error is always nil.
func (seg *addressSegmentInternal) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	var stringer func() string
	if with0xPrefix {
		stringer = func() string {
			return seg.toStringOpts(hexParamsSeg)
		}
	} else {
		stringer = func() string {
			return seg.toStringOpts(macCompressedParams)
		}
	}
	if seg.divisionValues != nil {
		if cache := seg.getCache(); cache != nil {
			if with0xPrefix {
				return cacheStr(&cache.cached0xHexString, stringer), nil
			}
			return cacheStr(&cache.cachedHexString, stringer), nil
		}
	}
	return stringer(), nil
}

func (seg *addressSegmentInternal) reverseMultiValSeg(perByte bool) (res *AddressSegment, err addrerr.IncompatibleAddressError) {
	if isReversible := seg.isReversibleRange(perByte); isReversible {
		// all reversible multi-valued segs reverse to the same segment
		res = seg.withoutPrefixLen()
		return
	}
	err = &incompatibleAddressError{addressError{key: "ipaddress.error.reverseRange"}}
	return
}

// ReverseBits returns a segment with the bits reversed.
//
// If this segment represents a range of values that cannot be reversed, then this returns an error.
//
// To be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
// Otherwise the result is not contiguous and thus cannot be represented by a sequential range of values.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (seg *addressSegmentInternal) ReverseBits(perByte bool) (res *AddressSegment, err addrerr.IncompatibleAddressError) {
	if seg.divisionValues == nil {
		res = seg.toAddressSegment()
		return
	}
	if seg.isMultiple() {
		return seg.reverseMultiValSeg(perByte)
	}
	byteCount := seg.GetByteCount()
	oldVal := seg.GetSegmentValue()
	var val SegInt
	switch byteCount {
	case 1:
		val = SegInt(reverseUint8(uint8(oldVal)))
	case 2:
		val = SegInt(reverseUint16(uint16(oldVal)))
		if perByte {
			val = ((val & 0xff) << 8) | (val >> 8)
		}
	case 3:
		val = reverseUint32(uint32(oldVal)) >> 8
		if perByte {
			val = ((val & 0xff) << 16) | (val & 0xff00) | (val >> 16)
		}
	case 4:
		val = reverseUint32(uint32(oldVal))
		if perByte {
			val = ((val & 0xff) << 24) | (val&0xff00)<<8 | (val&0xff0000)>>8 | (val >> 24)
		}
	default: // SegInt is at most 32 bits so this default case is not possible
		err = &incompatibleAddressError{addressError{key: "ipaddress.error.reverseRange"}}
		return
	}
	if oldVal == val && !seg.isPrefixed() {
		res = seg.toAddressSegment()
	} else {
		res = createAddressSegment(seg.deriveNewSeg(val, nil))
	}
	return
}

// ReverseBytes returns a segment with the bytes reversed.
//
// If this segment represents a range of values that cannot be reversed, then this returns an error.
//
// To be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
// Otherwise the result is not contiguous and thus cannot be represented by a sequential range of values.
func (seg *addressSegmentInternal) ReverseBytes() (res *AddressSegment, err addrerr.IncompatibleAddressError) {
	byteCount := seg.GetByteCount()
	if byteCount <= 1 {
		res = seg.toAddressSegment()
		return
	}
	if seg.isMultiple() {
		return seg.reverseMultiValSeg(false)
	}
	oldVal := seg.GetSegmentValue()
	var val SegInt
	switch byteCount {
	case 2:
		val = ((oldVal & 0xff) << 8) | (oldVal >> 8)
	case 3:
		val = ((oldVal & 0xff) << 16) | (oldVal & 0xff00) | (oldVal >> 16)
	case 4:
		val = ((oldVal & 0xff) << 24) | (oldVal&0xff00)<<8 | (oldVal&0xff0000)>>8 | (oldVal >> 24)
	default: // SegInt is at most 32 bits so this default case is not possible
		err = &incompatibleAddressError{addressError{key: "ipaddress.error.reverseRange"}}
		return
	}
	if oldVal == val && !seg.isPrefixed() {
		res = seg.toAddressSegment()
	} else {
		res = createAddressSegment(seg.deriveNewSeg(val, nil))
	}
	return
}

func (seg *addressSegmentInternal) isReversibleRange(perByte bool) (isReversible bool) {
	// Consider the case of reversing the bits of a range
	// Any range that can be successfully reversed must span all bits (otherwise after flipping you'd have a range in which the lower bit is constant, which is impossible in any contiguous range)
	// So that means at least one value has 0xxxx and another has 1xxxx (using 5 bits for our example). This means you must have the values 01111 and 10000 since the range is contiguous.
	// But reversing a range twice results in the original again, meaning the reversed must also be reversible, so the reversed also has 01111 and 10000.
	// So this means both the original and the reversed also have those two patterns flipped, which are 00001 and 11110.
	// So this means both ranges must span from at most 1 to at least 11110.
	// However, the two remaining values, 0 and 11111, are optional, as they are boundary value and remain themselves when reversed, and hence have no effect on whether the reversed range is contiguous.
	// So the only reversible ranges are 0-11111, 0-11110, 1-11110, and 1-11111.

	//-----------------------
	// Consider the case of reversing each of the bytes of a range.
	//
	// You can apply the same argument to the top multiple byte,
	// which means it is 0 or 1 to 254 or 255.
	// Suppose there is another byte to follow.
	// If you take the upper byte range, and you hold it constant, then reversing the next byte applies the same argument to that byte.
	// And so the lower byte must span from at most 1 to at least 11111110.
	// This argument holds when holding the upper byte constant at any value.
	// So the lower byte must span from at most 1 to at least 111111110 for any value.
	// So you have x 00000001-x 111111110 and y 00000001-y 111111110 and so on.

	// But all the bytes form a range, so you must also have the values in-between.
	// So that means you have 1 00000001 to 1 111111110 to 10 111111110 to 11 111111110 all the way to x 11111110, where x is at least 11111110.
	// In all cases, the upper byte lower value is at most 1, and 1 < 10000000.
	// That means you always have 10000000 00000000.
	// So you have the reverse as well (as argued above, for any value we also have the reverse).
	// So you always have 00000001 00000000.
	//
	// In other words, if the upper byte has lower 0, then the full bytes lower must be at most 0 00000001
	// Otherwise, when the upper byte has lower 1, the the full bytes lower is at most 1 00000000.
	//
	// In other words, if any upper byte has lower value 1, then all lower values to follow are 0.
	// If all upper bytes have lower value 0, then the next byte is permitted to have lower value 1.
	//
	// In summary, any upper byte having lower of 1 forces the remaining lower values to be 0.
	//
	// WHen the upper bytes are all zero, and thus the lower is at most 0 0 0 0 1,
	// then the only remaining lower value is 0 0 0 0 0.  This reverses to itself, so it is optional.
	//
	// The same argument applies to upper boundaries.
	//

	//-----------------------
	// Consider the case of reversing the bytes of a range.
	// Any range that can be successfully reversed must span all bits
	// (otherwise after flipping you'd have a range in which a lower bit is constant, which is impossible in any contiguous range)
	// So that means at least one value has 0xxxxx and another has 1xxxxx (we use 6 bits for our example, and we assume each byte has 3 bits).
	// This means you must have the values 011111 and 100000 since the range is contiguous.
	// But reversing a range twice results in the original again, meaning the reversed must also be reversible, so the reversed also has 011111 and 100000.

	// So this means both the original and the reversed also have those two bytes in each flipped, which are 111011 and 000100.
	// So the range must have 000100, 011111, 100000, 111011, so it must be at least 000100 to 111011.
	// So what if the range does not have 000001?  then the reversed range cannot have 001000, the byte-reversed address.
	// But we know it spans 000100 to 111011. So the original must have 000001.
	// What if it does not have 111110?  Then the reversed cannot have 110111, the byte-reversed address.
	// But we know it ranges from 000100 to 111011.  So the original must have 111110.
	// So it must range from 000001 to 111110.  The only remaining values in question are 000000 and 111111.
	// But once again, the two remaining values are optional, because the byte-reverse to themselves.
	// So for the byte-reverse case, we have the same potential ranges as in the bit-reverse case: 0-111111, 0-111110, 1-111110, and 1-111111
	if perByte {
		byteCount := seg.GetByteCount()
		bitCount := seg.GetBitCount()
		val := seg.GetSegmentValue()
		upperVal := seg.GetUpperSegmentValue()
		for i := 1; i <= byteCount; i++ {
			bitShift := i << 3
			shift := bitCount - BitCount(bitShift)
			byteVal := val >> uint(shift)
			upperByteVal := upperVal >> uint(shift)
			if byteVal != upperByteVal {
				if byteVal > 1 || upperByteVal < 254 {
					return false
				}
				i++
				if i <= byteCount {
					lowerIsZero := byteVal == 1
					upperIsMax := upperByteVal == 254
					for {
						bitShift = i << 3
						shift = bitCount - BitCount(bitShift)
						byteVal = val >> uint(shift)
						upperByteVal = upperVal >> uint(shift)
						if lowerIsZero {
							if byteVal != 0 {
								return
							}
						} else {
							if byteVal > 1 {
								return
							}
							lowerIsZero = byteVal == 1
						}
						if upperIsMax {
							if upperByteVal != 255 {
								return
							}
						} else {
							if upperByteVal < 254 {
								return
							}
							upperIsMax = upperByteVal == 254
						}
						i++
						if i > byteCount {
							break
						}
					}
				}
				return true
			}
		}
		return true
	}
	isReversible = seg.GetSegmentValue() <= 1 && seg.GetUpperSegmentValue() >= seg.GetMaxValue()-1
	return
}

//// only needed for godoc / pkgsite

// GetBitCount returns the number of bits in each value comprising this address item.
func (seg *addressSegmentInternal) GetBitCount() BitCount {
	return seg.addressDivisionInternal.GetBitCount()
}

// GetByteCount returns the number of bytes required for each value comprising this address item.
func (seg *addressSegmentInternal) GetByteCount() int {
	return seg.addressDivisionInternal.GetByteCount()
}

// GetValue returns the lowest value in the address segment range as a big integer.
func (seg *addressSegmentInternal) GetValue() *BigDivInt {
	return seg.addressDivisionInternal.GetValue()
}

// GetUpperValue returns the highest value in the address segment range as a big integer.
func (seg *addressSegmentInternal) GetUpperValue() *BigDivInt {
	return seg.addressDivisionInternal.GetUpperValue()
}

// Bytes returns the lowest value in the address segment range as a byte slice.
func (seg *addressSegmentInternal) Bytes() []byte {
	return seg.addressDivisionInternal.Bytes()
}

// UpperBytes returns the highest value in the address segment range as a byte slice.
func (seg *addressSegmentInternal) UpperBytes() []byte {
	return seg.addressDivisionInternal.UpperBytes()
}

// CopyBytes copies the lowest value in the address segment range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (seg *addressSegmentInternal) CopyBytes(bytes []byte) []byte {
	return seg.addressDivisionInternal.CopyBytes(bytes)
}

// CopyUpperBytes copies the highest value in the address segment range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (seg *addressSegmentInternal) CopyUpperBytes(bytes []byte) []byte {
	return seg.addressDivisionInternal.CopyUpperBytes(bytes)
}

// IsZero returns whether this segment matches exactly the value of zero.
func (seg *addressSegmentInternal) IsZero() bool {
	return seg.addressDivisionInternal.IsZero()
}

// IncludesZero returns whether this segment includes the value of zero within its range.
func (seg *addressSegmentInternal) IncludesZero() bool {
	return seg.addressDivisionInternal.IncludesZero()
}

// IsMax returns whether this segment matches exactly the maximum possible value, the value whose bits are all ones.
func (seg *addressSegmentInternal) IsMax() bool {
	return seg.addressDivisionInternal.IsMax()
}

// IncludesMax returns whether this segment includes the max value, the value whose bits are all ones, within its range.
func (seg *addressSegmentInternal) IncludesMax() bool {
	return seg.addressDivisionInternal.IncludesMax()
}

// IsFullRange returns whether the segment range includes all possible values for its bit length.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (seg *addressSegmentInternal) IsFullRange() bool {
	return seg.addressDivisionInternal.IsFullRange()
}

// ContainsPrefixBlock returns whether the segment range includes the block of values for the given prefix length.
func (seg *addressSegmentInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	return seg.addressDivisionInternal.ContainsPrefixBlock(prefixLen)
}

// ContainsSinglePrefixBlock returns whether the segment range matches exactly the block of values for the given prefix length and has just a single prefix for that prefix length.
func (seg *addressSegmentInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	return seg.addressDivisionInternal.ContainsSinglePrefixBlock(prefixLen)
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this segment includes the block of all values for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this segment represents a single value, this returns the bit count.
func (seg *addressSegmentInternal) GetMinPrefixLenForBlock() BitCount {
	return seg.addressDivisionInternal.GetMinPrefixLenForBlock()
}

// GetPrefixLenForSingleBlock returns a prefix length for which there is only one prefix in this segment,
// and the range of values in this segment matches the block of all values for that prefix.
//
// If the range of segment values can be described this way, then this method returns the same value as GetMinPrefixLenForBlock.
//
// If no such prefix length exists, returns nil.
//
// If this segment represents a single value, this returns the bit count of the segment.
func (seg *addressSegmentInternal) GetPrefixLenForSingleBlock() PrefixLen {
	return seg.addressDivisionInternal.GetPrefixLenForSingleBlock()
}

// IsSinglePrefix determines if the segment has a single prefix value for the given prefix length.  You can call GetPrefixCountLen to get the count of prefixes.
func (seg *addressSegmentInternal) IsSinglePrefix(divisionPrefixLength BitCount) bool {
	return seg.addressDivisionInternal.IsSinglePrefix(divisionPrefixLength)
}

//// end needed for godoc / pkgsite

//

// AddressSegment represents a single segment of an address.  A segment contains a single value or a range of sequential values and it has an assigned bit length.
//
// The current implementations of this type are the most common representations of IPv4, IPv6 and MAC;
// segments are 1 byte for Ipv4, they are two bytes for Ipv6, and they are 1 byte for MAC addresses.
//
// There are alternative forms of dividing addresses into divisions, such as the dotted representation for MAC like "1111.2222.3333",
// the embedded IPv4 representation for IPv6 like "f:f:f:f:f:f:1.2.3.4", the inet_aton formats like "1.2" for IPv4, and so on.
//
// The general rules are that segments have a whole number of bytes, and in a given address all segments have the same length.
//
// When alternatives forms do not follow the general rules for segments, you can use [AddressDivision] instead.
// Divisions do not have the restriction that divisions of an address are equal length and a whole number of bytes.
// Divisions can be grouped using [AddressDivisionGrouping].
//
// AddressSegment objects are immutable and thus are also concurrency-safe.
type AddressSegment struct {
	addressSegmentInternal
}

// Contains returns whether this is same type and version as the given segment and whether it contains all values in the given segment.
func (seg *AddressSegment) Contains(other AddressSegmentType) bool {
	if seg == nil {
		return other == nil || other.ToSegmentBase() == nil
	}
	return seg.contains(other)
}

// Overlaps returns whether this is same type and version as the given segment and whether it overlaps with the values in the given segment.
func (seg *AddressSegment) Overlaps(other AddressSegmentType) bool {
	if seg == nil {
		return other == nil || other.ToSegmentBase() == nil
	}
	return seg.overlaps(other)
}

// Equal returns whether the given segment is equal to this segment.
// Two segments are equal if they match:
//   - type/version (IPv4, IPv6, MAC)
//   - value range
//
// Prefix lengths are ignored.
func (seg *AddressSegment) Equal(other AddressSegmentType) bool {
	if seg == nil {
		return other == nil || other.ToDiv() == nil
	}
	return seg.equal(other)
}

// Compare returns a negative integer, zero, or a positive integer if this address segment is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (seg *AddressSegment) Compare(item AddressItem) int {
	return CountComparator.Compare(seg, item)
}

// CompareSize compares the counts of two items, the number of individual values within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one represents more individual values than another.
//
// CompareSize returns a positive integer if this segment has a larger count than the item given, zero if they are the same, or a negative integer if the other has a larger count.
func (seg *AddressSegment) CompareSize(other AddressItem) int {
	if seg == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return seg.compareSize(other)
}

// GetLower returns a segment representing just the lowest value in the range, which will be the same segment if it represents a single value.
func (seg *addressSegmentInternal) GetLower() *AddressSegment {
	return seg.getLower()
}

// GetUpper returns a segment representing just the highest value in the range, which will be the same segment if it represents a single value.
func (seg *addressSegmentInternal) GetUpper() *AddressSegment {
	return seg.getUpper()
}

// IsMultiple returns whether this segment represents multiple values.
func (seg *AddressSegment) IsMultiple() bool {
	return seg != nil && seg.isMultiple()
}

// GetCount returns the count of possible distinct values for this item.
// If not representing multiple values, the count is 1.
//
// For instance, a segment with the value range of 3-7 has count 5.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (seg *AddressSegment) GetCount() *big.Int {
	if seg == nil {
		return bigZero()
	}
	return seg.getCount()
}

// IsIP returns true if this segment originated as an IPv4 or IPv6 segment, or an implicitly zero-valued IP segment.  If so, use ToIP to convert back to the IP-specific type.
func (seg *AddressSegment) IsIP() bool {
	return seg != nil && seg.matchesIPSegment()
}

// IsIPv4 returns true if this segment originated as an IPv4 segment.  If so, use ToIPv4 to convert back to the IPv4-specific type.
func (seg *AddressSegment) IsIPv4() bool {
	return seg != nil && seg.matchesIPv4Segment()
}

// IsIPv6 returns true if this segment originated as an IPv6 segment.  If so, use ToIPv6 to convert back to the IPv6-specific type.
func (seg *AddressSegment) IsIPv6() bool {
	return seg != nil && seg.matchesIPv6Segment()
}

// IsMAC returns true if this segment originated as a MAC segment.  If so, use ToMAC to convert back to the MAC-specific type.
func (seg *AddressSegment) IsMAC() bool {
	return seg != nil && seg.matchesMACSegment()
}

// Iterator provides an iterator to iterate through the individual address segments of this address segment.
//
// Call IsMultiple to determine if this instance represents multiple address segments, or GetValueCount for the count.
func (seg *AddressSegment) Iterator() Iterator[*AddressSegment] {
	if seg == nil {
		return nilSegIterator()
	}
	return seg.iterator()
}

// ToIP converts to an IPAddressSegment if this division originated as an IPv4 or IPv6 segment, or an implicitly zero-valued IP segment.
// If not, ToIP returns nil.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *AddressSegment) ToIP() *IPAddressSegment {
	if seg.IsIP() {
		return (*IPAddressSegment)(unsafe.Pointer(seg))
	}
	return nil
}

// ToIPv4 converts to an IPv4AddressSegment if this segment originated as an IPv4 segment.
// If not, ToIPv4 returns nil.
//
// ToIPv4 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *AddressSegment) ToIPv4() *IPv4AddressSegment {
	if seg.IsIPv4() {
		return (*IPv4AddressSegment)(unsafe.Pointer(seg))
	}
	return nil
}

// ToIPv6 converts to an IPv6AddressSegment if this segment originated as an IPv6 segment.
// If not, ToIPv6 returns nil.
//
// ToIPv6 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *AddressSegment) ToIPv6() *IPv6AddressSegment {
	if seg.IsIPv6() {
		return (*IPv6AddressSegment)(unsafe.Pointer(seg))
	}
	return nil
}

// ToMAC converts to a MACAddressSegment if this segment originated as a MAC segment.
// If not, ToMAC returns nil.
//
// ToMAC can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *AddressSegment) ToMAC() *MACAddressSegment {
	if seg.IsMAC() {
		return (*MACAddressSegment)(seg)
	}
	return nil
}

// ToSegmentBase is an identity method.
//
// ToSegmentBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *AddressSegment) ToSegmentBase() *AddressSegment {
	return seg
}

// ToDiv converts to an AddressDivision, a polymorphic type usable with all address segments and divisions.
// Afterwards, you can convert back with ToSegmentBase.
//
// ToDiv can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *AddressSegment) ToDiv() *AddressDivision {
	return (*AddressDivision)(unsafe.Pointer(seg))
}

// GetString produces a normalized string to represent the segment.
// If the segment is an IP segment string with CIDR network prefix block for its prefix length, then the string contains only the lower value of the block range.
// Otherwise, the explicit range will be printed.
// If the segment is not an IP segment, then the string is the same as that produced by GetWildcardString.
//
// The string returned is useful in the context of creating strings for address sections or full addresses,
// in which case the radix and bit-length can be deduced from the context.
// The String method produces strings more appropriate when no context is provided.
func (seg *AddressSegment) GetString() string {
	if seg == nil {
		return nilString()
	}
	return seg.getString()
}

// GetWildcardString produces a normalized string to represent the segment, favouring wildcards and range characters while ignoring any network prefix length.
// The explicit range of a range-valued segment will be printed.
//
// The string returned is useful in the context of creating strings for address sections or full addresses,
// in which case the radix and the bit-length can be deduced from the context.
// The String method produces strings more appropriate when no context is provided.
func (seg *AddressSegment) GetWildcardString() string {
	if seg == nil {
		return nilString()
	}
	return seg.getWildcardString()
}

// String produces a string that is useful when a segment string is provided with no context.
// If the segment was originally constructed as an IPv4 address segment it uses decimal, otherwise hexadecimal.
// It uses a string prefix for hex ("0x"), and does not use the wildcard '*', because division size is variable, so '*' is ambiguous.
// GetWildcardString is more appropriate in context with other segments or divisions.  It does not use a string prefix and uses '*' for full-range segments.
// GetString is more appropriate in context with prefix lengths, it uses zeros instead of wildcards with full prefix block ranges alongside prefix lengths.
func (seg *AddressSegment) String() string {
	if seg == nil {
		return nilString()
	}
	return seg.toString()
}

func segsSame(onePref, twoPref PrefixLen, oneVal, twoVal, oneUpperVal, twoUpperVal SegInt) bool {
	return onePref.Equal(twoPref) &&
		oneVal == twoVal && oneUpperVal == twoUpperVal
}

func segValsSame(oneVal, twoVal, oneUpperVal, twoUpperVal SegInt) bool {
	return oneVal == twoVal && oneUpperVal == twoUpperVal
}

func segValSame(oneVal, twoVal SegInt) bool {
	return oneVal == twoVal
}

func getPrefixValueCount(segment *AddressSegment, segmentPrefixLength BitCount) SegIntCount {
	shiftAdjustment := segment.GetBitCount() - segmentPrefixLength
	if shiftAdjustment <= 0 {
		return SegIntCount(segment.GetUpperSegmentValue()) - SegIntCount(segment.GetSegmentValue()) + 1
	}
	return SegIntCount(segment.GetUpperSegmentValue()>>uint(shiftAdjustment)) - SegIntCount(segment.GetSegmentValue()>>uint(shiftAdjustment)) + 1
}

func getSegmentPrefLen(
	_ AddressSegmentSeries,
	prefLen PrefixLen,
	bitsPerSegment,
	bitsMatchedSoFar BitCount,
	segment *AddressSegment) PrefixLen {
	if ipSeg := segment.ToIP(); ipSeg != nil {
		return ipSeg.GetSegmentPrefixLen()
	} else if prefLen != nil {
		result := prefLen.Len() - bitsMatchedSoFar
		if result <= bitsPerSegment {
			if result < 0 {
				result = 0
			}
			return cacheBitCount(result)
		}
	}
	return nil
}

func getMatchingBits(segment1, segment2 *AddressSegment, maxBits, bitsPerSegment BitCount) BitCount {
	if maxBits == 0 {
		return 0
	}
	val1 := segment1.getSegmentValue()
	val2 := segment2.getSegmentValue()
	xor := val1 ^ val2
	switch bitsPerSegment {
	case IPv4BitsPerSegment:
		return BitCount(bits.LeadingZeros8(uint8(xor)))
	case IPv6BitsPerSegment:
		return BitCount(bits.LeadingZeros16(uint16(xor)))
	default:
		return BitCount(bits.LeadingZeros32(xor)) - 32 + bitsPerSegment
	}
}
