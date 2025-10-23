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
	"math/bits"
	"strings"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

// DivInt is an integer type for holding generic division values, which can be larger than segment values
type DivInt = uint64

const DivIntSize = 64

type divderiver interface {
	// deriveNew produces a new division with the same bit count as the old,
	// but with the new values and prefix length
	deriveNew(val, upperVal DivInt, prefLen PrefixLen) divisionValues

	// derivePrefixed produces a new division with the same bit count and values as the old,
	// but with the new prefix length
	derivePrefixed(prefLen PrefixLen) divisionValues
}

type divIntVals interface {
	// getDivisionValue gets the lower value for a division
	getDivisionValue() DivInt

	// getUpperDivisionValue gets the upper value for a division
	getUpperDivisionValue() DivInt
}

func newDivValues(value, upperValue DivInt, prefLen PrefixLen, bitCount BitCount) *divIntValues {
	if value > upperValue {
		value, upperValue = upperValue, value
	}
	if bitCount <= 0 {
		value = 0
		upperValue = 0
	} else if (1 << uint(bitCount)) <= upperValue { // upperValue too big
		max := ^(^DivInt(0) << uint(bitCount))
		value &= max
		upperValue &= max
	}
	prefLen = checkPrefLen(prefLen, bitCount)
	return newDivValuesUnchecked(value, upperValue, prefLen, bitCount)
}

func newDivValuesUnchecked(value, upperValue DivInt, prefLen PrefixLen, bitCount BitCount) *divIntValues {
	return &divIntValues{
		value:      value,
		upperValue: upperValue,
		prefLen:    prefLen,
		bitCount:   bitCount,
	}
}

// divIntValues are used by AddressDivision
type divIntValues struct {
	bitCount          BitCount
	value, upperValue DivInt
	prefLen           PrefixLen
	cache             divCache
}

func (div *divIntValues) getBitCount() BitCount {
	return div.bitCount
}

func (div *divIntValues) getByteCount() int {
	return (int(div.getBitCount()) + 7) >> 3
}

func (div *divIntValues) getDivisionPrefixLength() PrefixLen {
	return div.prefLen
}

func (div *divIntValues) getValue() *BigDivInt {
	return big.NewInt(int64(div.value))
}

func (div *divIntValues) getUpperValue() *BigDivInt {
	return big.NewInt(int64(div.upperValue))
}

func (div *divIntValues) includesZero() bool {
	return div.value == 0
}

func (div *divIntValues) includesMax() bool {
	return div.upperValue == ^((^DivInt(0)) << uint(div.getBitCount()))
}

func (div *divIntValues) isMultiple() bool {
	return div.value != div.upperValue
}

func (div *divIntValues) getCount() *big.Int {
	res := bigZero()
	return res.SetUint64(uint64(div.upperValue-div.value)).Add(res, bigOneConst())
}

func (div *divIntValues) calcBytesInternal() (bytes, upperBytes []byte) {
	return calcBytesInternal(div.getByteCount(), div.getDivisionValue(), div.getUpperDivisionValue())
}

func (div *divIntValues) bytesInternal(upper bool) []byte {
	if upper {
		return calcSingleBytes(div.getByteCount(), div.getUpperDivisionValue())
	}
	return calcSingleBytes(div.getByteCount(), div.getDivisionValue())
}

func calcBytesInternal(byteCount int, val, upperVal DivInt) (bytes, upperBytes []byte) {
	byteIndex := byteCount - 1
	isMultiple := val != upperVal
	bytes = make([]byte, byteCount)
	if isMultiple {
		upperBytes = make([]byte, byteCount)
	} else {
		upperBytes = bytes
	}
	for {
		bytes[byteIndex] |= byte(val)
		val >>= 8
		if isMultiple {
			upperBytes[byteIndex] |= byte(upperVal)
			upperVal >>= 8
		}
		if byteIndex == 0 {
			return bytes, upperBytes
		}
		byteIndex--
	}
}

func calcSingleBytes(byteCount int, val DivInt) (bytes []byte) {
	byteIndex := byteCount - 1
	bytes = make([]byte, byteCount)
	for {
		bytes[byteIndex] |= byte(val)
		val >>= 8
		if byteIndex == 0 {
			return bytes
		}
		byteIndex--
	}
}

func (div *divIntValues) getCache() *divCache {
	return &div.cache
}

func (div *divIntValues) getAddrType() addrType {
	return zeroType
}

func (div *divIntValues) getDivisionValue() DivInt {
	return div.value
}

func (div *divIntValues) getUpperDivisionValue() DivInt {
	return div.upperValue
}

func (div *divIntValues) getSegmentValue() SegInt {
	return SegInt(div.value)
}

func (div *divIntValues) getUpperSegmentValue() SegInt {
	return SegInt(div.upperValue)
}

func (div *divIntValues) deriveNew(val, upperVal DivInt, prefLen PrefixLen) divisionValues {
	return newDivValuesUnchecked(val, upperVal, prefLen, div.bitCount)
}

func (div *divIntValues) derivePrefixed(prefLen PrefixLen) divisionValues {
	return newDivValuesUnchecked(div.value, div.upperValue, prefLen, div.bitCount)
}

func (div *divIntValues) deriveNewMultiSeg(val, upperVal SegInt, prefLen PrefixLen) divisionValues {
	return newDivValuesUnchecked(DivInt(val), DivInt(upperVal), prefLen, div.bitCount)
}

func (div *divIntValues) deriveNewSeg(val SegInt, prefLen PrefixLen) divisionValues {
	value := DivInt(val)
	return newDivValuesUnchecked(value, value, prefLen, div.bitCount)
}

var _ divisionValues = &divIntValues{}

func createAddressDivision(vals divisionValues) *AddressDivision {
	return &AddressDivision{
		addressDivisionInternal{
			addressDivisionBase: addressDivisionBase{vals},
		},
	}
}

type addressDivisionInternal struct {
	addressDivisionBase
}

var (
	// wildcards differ, here we use only range since div size not implicit
	octalParamsDiv   = new(addrstr.IPStringOptionsBuilder).SetRadix(8).SetSegmentStrPrefix(OctalPrefix).SetWildcards(rangeWildcard).ToOptions()
	hexParamsDiv     = new(addrstr.IPStringOptionsBuilder).SetRadix(16).SetSegmentStrPrefix(HexPrefix).SetWildcards(rangeWildcard).ToOptions()
	decimalParamsDiv = new(addrstr.IPStringOptionsBuilder).SetRadix(10).SetWildcards(rangeWildcard).ToOptions()
)

// String produces a string that is useful when a division string is provided with no context.
// If the division was originally constructed as an address segment, uses the default radix for that segment, which is decimal for IPv4 and hexadecimal for IPv6, MAC or other.
// It uses a string prefix for octal or hex ("0" or "0x"), and does not use the wildcard '*', because division size is variable, so '*' is ambiguous.
// GetWildcardString is more appropriate in context with other segments or divisions.  It does not use a string prefix and uses '*' for full-range segments.
// GetString is more appropriate in context with prefix lengths, it uses zeros instead of wildcards for prefix block ranges.
func (div *addressDivisionInternal) String() string {
	return div.toString()
}

// toString produces a string that is useful when a division string is provided with no context.
// It uses a string prefix for octal or hex ("0" or "0x"), and does not use the wildcard '*', because division size is variable, so '*' is ambiguous.
// GetWildcardString() is more appropriate in context with other segments or divisions.  It does not use a string prefix and uses '*' for full-range segments.
// GetString() is more appropriate in context with prefix lengths, it uses zeros instead of wildcards for prefix block ranges.
func (div *addressDivisionInternal) toString() string { // this can be moved to addressDivisionBase when we have ContainsPrefixBlock and similar methods implemented for big.Int in the base.
	return toString(div.toAddressDivision())
}

// Format implements [fmt.Formatter] interface. It accepts the formats
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
func (div addressDivisionInternal) Format(state fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		_, _ = state.Write([]byte(div.toString()))
		return
	}
	// we try to filter through the flags provided to the DivInt values, as if the fmt string were applied to the int(s) directly
	formatStr := flagsFromState(state, verb)
	if div.isMultiple() {
		formatStr = fmt.Sprintf("%s%c%s", formatStr, RangeSeparator, formatStr)
		_, _ = state.Write([]byte(fmt.Sprintf(formatStr, div.getDivisionValue(), div.getUpperDivisionValue())))
	} else {
		_, _ = state.Write([]byte(fmt.Sprintf(formatStr, div.getDivisionValue())))
	}
}

func (div *addressDivisionInternal) toStringOpts(opts addrstr.StringOptions) string {
	return toStringOpts(opts, div.toAddressDivision())
}

func (div *addressDivisionInternal) isPrefixed() bool {
	return div.getDivisionPrefixLength() != nil
}

// return whether the division range includes the block of values for the given prefix length.
func (div *addressDivisionInternal) containsPrefixBlock(divisionPrefixLen BitCount) bool {
	return div.isPrefixBlockVals(div.getDivisionValue(), div.getUpperDivisionValue(), divisionPrefixLen)
}

// Returns whether the division range includes the block of values for its prefix length.
func (div *addressDivisionInternal) isPrefixBlockVals(divisionValue, upperValue DivInt, divisionPrefixLen BitCount) bool {
	return isPrefixBlockVals(divisionValue, upperValue, divisionPrefixLen, div.GetBitCount())
}

func isPrefixBlockVals(divisionValue, upperValue DivInt, divisionPrefixLen, divisionBitCount BitCount) bool {
	if divisionPrefixLen <= 0 {
		if divisionValue != 0 {
			return false
		}
		maxValue := ^(^DivInt(0) << uint(divisionBitCount))
		return upperValue == maxValue
	}
	if divisionPrefixLen >= divisionBitCount {
		return true
	}
	var ones = ^DivInt(0)
	divisionBitMask := ^(ones << uint(divisionBitCount))
	divisionPrefixMask := ones << uint(divisionBitCount-divisionPrefixLen)
	var divisionNonPrefixMask = ^divisionPrefixMask
	return testRange(divisionValue,
		upperValue,
		upperValue,
		divisionPrefixMask&divisionBitMask,
		divisionNonPrefixMask)
}

// Returns whether the given range of divisionValue to upperValue is equivalent to the range of segmentValue with the prefix of divisionPrefixLen
func (div *addressDivisionInternal) isSinglePrefix(divisionValue, upperValue DivInt, divisionPrefixLen BitCount) bool {
	bitCount := div.GetBitCount()
	divisionPrefixLen = checkBitCount(divisionPrefixLen, bitCount)
	shift := uint(bitCount - divisionPrefixLen)
	return (divisionValue >> shift) == (upperValue >> shift)
}

// Returns whether the given range of divisionValue to upperValue is equivalent to the range of segmentValue with the prefix of divisionPrefixLen
func (div *addressDivisionInternal) isSinglePrefixBlock(divisionValue, upperValue DivInt, divisionPrefixLen BitCount) bool {
	if divisionPrefixLen == 0 {
		return divisionValue == 0 && upperValue == div.getMaxValue()
	}
	bitCount := div.GetBitCount()
	ones := ^DivInt(0)
	divisionBitMask := ^(ones << uint(bitCount))
	divisionPrefixMask := ones << uint(bitCount-divisionPrefixLen)
	divisionHostMask := ^divisionPrefixMask
	return testRange(divisionValue,
		divisionValue,
		upperValue,
		divisionPrefixMask&divisionBitMask,
		divisionHostMask)
}

// ContainsPrefixBlock returns whether the division range includes the block of values for the given prefix length.
func (div *addressDivisionInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	return div.isPrefixBlockVals(div.getDivisionValue(), div.getUpperDivisionValue(), prefixLen)
}

// ContainsSinglePrefixBlock returns whether the division range matches exactly the block of values for the given prefix length and has just a single prefix for that prefix length.
func (div *addressDivisionInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	prefixLen = checkDiv(div.toAddressDivision(), prefixLen)
	return div.isSinglePrefixBlock(div.getDivisionValue(), div.getUpperDivisionValue(), prefixLen)
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this division includes the block of all values for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this division represents a single value, this returns the bit count.
func (div *addressDivisionInternal) GetMinPrefixLenForBlock() BitCount {
	return getMinPrefixLenForBlock(div.getDivisionValue(), div.getUpperDivisionValue(), div.GetBitCount())
}

// GetPrefixLenForSingleBlock returns a prefix length for which there is only one prefix in this division,
// and the range of values in this division matches the block of all values for that prefix.
//
// If the range of division values can be described this way, then this method returns the same value as GetMinPrefixLenForBlock.
//
// If no such prefix length exists, returns nil.
//
// If this division represents a single value, this returns the bit count of the segment.
func (div *addressDivisionInternal) GetPrefixLenForSingleBlock() PrefixLen {
	return getPrefixLenForSingleBlock(div.getDivisionValue(), div.getUpperDivisionValue(), div.GetBitCount())
}

// return whether the division range includes the block of values for the division prefix length,
// or false if the division has no prefix length.
func (div *addressDivisionInternal) isPrefixBlock() bool {
	prefLen := div.getDivisionPrefixLength()
	return prefLen != nil && div.containsPrefixBlock(prefLen.bitCount())
}

func (div *addressDivisionInternal) getMaxValue() DivInt {
	return ^(^DivInt(0) << uint(div.GetBitCount()))
}

func (div *addressDivisionInternal) getDivisionValue() DivInt {
	vals := div.divisionValues
	if vals == nil {
		return 0
	}
	return vals.getDivisionValue()
}

func (div *addressDivisionInternal) getUpperDivisionValue() DivInt {
	vals := div.divisionValues
	if vals == nil {
		return 0
	}
	return vals.getUpperDivisionValue()
}

func (div *addressDivisionInternal) matches(value DivInt) bool {
	return !div.isMultiple() && value == div.getDivisionValue()
}

func (div *addressDivisionInternal) matchesWithMask(value, mask DivInt) bool {
	if div.isMultiple() {
		//we want to ensure that any of the bits that can change from value to upperValue is masked out (zeroed) by the mask.
		//In other words, when masked we need all values represented by this segment to become just a single value
		diffBits := div.getDivisionValue() ^ div.getUpperDivisionValue()
		leadingZeros := bits.LeadingZeros64(diffBits)
		//the bits that can change are all bits following the first leadingZero bits
		//all the bits that follow must be zeroed out by the mask
		fullMask := ^DivInt(0) >> uint(leadingZeros)
		if (fullMask & mask) != 0 {
			return false
		} //else we know that the mask zeros out all the bits that can change from value to upperValue, so now we just compare with either one
	}
	return value == (div.getDivisionValue() & mask)
}

// matchesWithMask returns whether masking with the given mask results in a valid contiguous range for this segment,
// and if it does, if the result matches the range of lowerValue to upperValue.
func (div *addressDivisionInternal) matchesValsWithMask(lowerValue, upperValue, mask DivInt) bool {
	if lowerValue == upperValue {
		return div.matchesWithMask(lowerValue, mask)
	}
	if !div.isMultiple() {
		// the values to match, lowerValue and upperValue, are not the same, so impossible to match those two values with a single value from this segment
		return false
	}
	thisValue := div.getDivisionValue()
	thisUpperValue := div.getUpperDivisionValue()
	masker := MaskRange(thisValue, thisUpperValue, mask, div.getMaxValue())
	if !masker.IsSequential() {
		return false
	}
	return lowerValue == masker.GetMaskedLower(thisValue, mask) && upperValue == masker.GetMaskedUpper(thisUpperValue, mask)
}

func (div *addressDivisionInternal) toPrefixedNetworkDivision(divPrefixLength PrefixLen) *AddressDivision {
	return div.toNetworkDivision(divPrefixLength, true)
}

func (div *addressDivisionInternal) toNetworkDivision(divPrefixLength PrefixLen, withPrefixLength bool) *AddressDivision {
	vals := div.divisionValues
	if vals == nil {
		return div.toAddressDivision()
	}
	lower := div.getDivisionValue()
	upper := div.getUpperDivisionValue()
	var newLower, newUpper DivInt
	hasPrefLen := divPrefixLength != nil
	if hasPrefLen {
		prefBits := divPrefixLength.bitCount()
		bitCount := div.GetBitCount()
		prefBits = checkBitCount(prefBits, bitCount)
		mask := ^DivInt(0) << uint(bitCount-prefBits)
		newLower = lower & mask
		newUpper = upper | ^mask
		if !withPrefixLength {
			divPrefixLength = nil
		}
		if divsSame(divPrefixLength, div.getDivisionPrefixLength(), newLower, lower, newUpper, upper) {
			return div.toAddressDivision()
		}
	} else {
		divPrefixLength = nil
		if div.getDivisionPrefixLength() == nil {
			return div.toAddressDivision()
		}
	}
	newVals := div.deriveNew(newLower, newUpper, divPrefixLength)
	return createAddressDivision(newVals)
}

func (div *addressDivisionInternal) toPrefixedHostDivision(divPrefixLength PrefixLen) *AddressDivision {
	return div.toHostDivision(divPrefixLength, true)
}

func (div *addressDivisionInternal) toHostDivision(divPrefixLength PrefixLen, withPrefixLength bool) *AddressDivision {
	vals := div.divisionValues
	if vals == nil {
		return div.toAddressDivision()
	}
	lower := div.getDivisionValue()
	upper := div.getUpperDivisionValue()
	//var newLower, newUpper DivInt
	hasPrefLen := divPrefixLength != nil
	var mask SegInt
	if hasPrefLen {
		prefBits := divPrefixLength.bitCount()
		bitCount := div.GetBitCount()
		prefBits = checkBitCount(prefBits, bitCount)
		mask = ^(^SegInt(0) << uint(bitCount-prefBits))
	}
	divMask := uint64(mask)
	maxVal := uint64(^SegInt(0))
	masker := MaskRange(lower, upper, divMask, maxVal)
	newLower, newUpper := masker.GetMaskedLower(lower, divMask), masker.GetMaskedUpper(upper, divMask)
	if !withPrefixLength {
		divPrefixLength = nil
	}
	if divsSame(divPrefixLength, div.getDivisionPrefixLength(), newLower, lower, newUpper, upper) {
		return div.toAddressDivision()
	}
	newVals := div.deriveNew(newLower, newUpper, divPrefixLength)
	return createAddressDivision(newVals)
}

func (div *addressDivisionInternal) toPrefixedDivision(divPrefixLength PrefixLen) *AddressDivision {
	hasPrefLen := divPrefixLength != nil
	bitCount := div.GetBitCount()
	if hasPrefLen {
		prefBits := divPrefixLength.bitCount()
		prefBits = checkBitCount(prefBits, bitCount)
		if div.isPrefixed() && prefBits == div.getDivisionPrefixLength().bitCount() {
			return div.toAddressDivision()
		}
	} else {
		return div.toAddressDivision()
	}
	lower := div.getDivisionValue()
	upper := div.getUpperDivisionValue()
	newVals := div.deriveNew(lower, upper, divPrefixLength)
	return createAddressDivision(newVals)
}

func (div *addressDivisionInternal) getCount() *big.Int {
	if !div.isMultiple() {
		return bigOne()
	}
	if div.IsFullRange() {
		res := bigZero()
		return res.SetUint64(0xffffffffffffffff).Add(res, bigOneConst())
	}
	return bigZero().SetUint64((div.getUpperDivisionValue() - div.getDivisionValue()) + 1)
}

// IsSinglePrefix returns true if the division value range spans just a single prefix value for the given prefix length.
func (div *addressDivisionInternal) IsSinglePrefix(divisionPrefixLength BitCount) bool {
	return div.isSinglePrefix(div.getDivisionValue(), div.getUpperDivisionValue(), divisionPrefixLength)
}

// GetPrefixCountLen returns the number of distinct prefixes in the division value range for the given prefix length.
func (div *addressDivisionInternal) GetPrefixCountLen(divisionPrefixLength BitCount) *big.Int {
	if div.IsFullRange() {
		return bigZero().Add(bigOneConst(), bigZero().SetUint64(div.getMaxValue()))
	}
	bitCount := div.GetBitCount()
	divisionPrefixLength = checkBitCount(divisionPrefixLength, bitCount)
	shiftAdjustment := bitCount - divisionPrefixLength
	count := ((div.getUpperDivisionValue() >> uint(shiftAdjustment)) - (div.getDivisionValue() >> uint(shiftAdjustment))) + 1
	return bigZero().SetUint64(count)
}

func (div *addressDivisionInternal) matchesIPSegment() bool {
	return div.divisionValues == nil || div.getAddrType().isIP()
}

func (div *addressDivisionInternal) matchesIPv4Segment() bool {
	// the init() methods ensure even zero-IPv4 segments (IPv4Segment{}) have addr type IPv4
	return div.divisionValues != nil && div.getAddrType().isIPv4()
}

func (div *addressDivisionInternal) matchesIPv6Segment() bool {
	// the init() methods ensure even zero IPv6 segments (IPv6Segment{}) have addr type IPv6
	return div.divisionValues != nil && div.getAddrType().isIPv6()
}

func (div *addressDivisionInternal) matchesMACSegment() bool {
	// the init() methods ensure even zero MAC segments (MACSegment{}) have addr type MAC
	return div.divisionValues != nil && div.getAddrType().isMAC()
}

func (div *addressDivisionInternal) matchesSegment() bool {
	return div.GetBitCount() <= SegIntSize
}

func (div *addressDivisionInternal) toAddressDivision() *AddressDivision {
	return (*AddressDivision)(unsafe.Pointer(div))
}

func (div *addressDivisionInternal) toAddressSegment() *AddressSegment {
	if div.matchesSegment() {
		return (*AddressSegment)(unsafe.Pointer(div))
	}
	return nil
}

func (div *addressDivisionInternal) getStringAsLower() string {
	if seg := div.toAddressDivision().ToIP(); seg != nil {
		return seg.getStringAsLower()
	}
	return div.getStringFromStringer(div.getDefaultLowerString)
}

func (div *addressDivisionInternal) getDivString() string {
	if !div.isMultiple() {
		return div.getStringFromStringer(div.getDefaultLowerString)
	} else {
		return div.getStringFromStringer(div.getDefaultRangeString)
	}
}

func (div *addressDivisionInternal) getStringFromStringer(stringer func() string) string {
	if div.divisionValues != nil {
		if cache := div.getCache(); cache != nil {
			return cacheStr(&cache.cachedString, stringer)
		}
	}
	return stringer()
}

func (div *addressDivisionInternal) getString() string {
	if seg := div.toAddressDivision().ToIP(); seg != nil {
		return seg.GetString()
	}
	return div.getDivString()
}

func (div *addressDivisionInternal) getWildcardString() string {
	if seg := div.toAddressDivision().ToIP(); seg != nil {
		return seg.GetWildcardString()
	}
	return div.getDivString() // same string as GetString() when not an IP segment
}

func (div *addressDivisionInternal) getDefaultRangeStringVals(val1, val2 uint64, radix int) string {
	return getDefaultRangeStringVals(div, val1, val2, radix)
}

func (div *addressDivisionInternal) buildDefaultRangeString(radix int) string {
	return buildDefaultRangeString(div, radix)
}

func (div *addressDivisionInternal) getLowerStringLength(radix int) int {
	return toUnsignedStringLength(div.getDivisionValue(), radix)
}

func (div *addressDivisionInternal) getUpperStringLength(radix int) int {
	return toUnsignedStringLength(div.getUpperDivisionValue(), radix)
}

func (div *addressDivisionInternal) getLowerString(radix int, uppercase bool, appendable *strings.Builder) {
	toUnsignedStringCased(div.getDivisionValue(), radix, 0, uppercase, appendable)
}

func (div *addressDivisionInternal) getLowerStringChopped(radix int, choppedDigits int, uppercase bool, appendable *strings.Builder) {
	toUnsignedStringCased(div.getDivisionValue(), radix, choppedDigits, uppercase, appendable)
}

func (div *addressDivisionInternal) getUpperString(radix int, uppercase bool, appendable *strings.Builder) {
	toUnsignedStringCased(div.getUpperDivisionValue(), radix, 0, uppercase, appendable)
}

func (div *addressDivisionInternal) getUpperStringMasked(radix int, uppercase bool, appendable *strings.Builder) {
	if seg := div.toAddressDivision().ToIP(); seg != nil {
		seg.getUpperStringMasked(radix, uppercase, appendable)
	} else if div.isPrefixed() {
		upperValue := div.getUpperDivisionValue()
		mask := ^DivInt(0) << uint(div.GetBitCount()-div.getDivisionPrefixLength().bitCount())
		upperValue &= mask
		toUnsignedStringCased(upperValue, radix, 0, uppercase, appendable)
	} else {
		div.getUpperString(radix, uppercase, appendable)
	}
}

func (div *addressDivisionInternal) getSplitLowerString(radix int, choppedDigits int, uppercase bool,
	splitDigitSeparator byte, reverseSplitDigits bool, stringPrefix string, appendable *strings.Builder) {
	toSplitUnsignedString(div.getDivisionValue(), radix, choppedDigits, uppercase, splitDigitSeparator, reverseSplitDigits, stringPrefix, appendable)
}

func (div *addressDivisionInternal) getSplitRangeString(rangeSeparator string, wildcard string, radix int, uppercase bool,
	splitDigitSeparator byte, reverseSplitDigits bool, stringPrefix string, appendable *strings.Builder) addrerr.IncompatibleAddressError {
	return toUnsignedSplitRangeString(
		div.getDivisionValue(),
		div.getUpperDivisionValue(),
		rangeSeparator,
		wildcard,
		radix,
		uppercase,
		splitDigitSeparator,
		reverseSplitDigits,
		stringPrefix,
		appendable)
}

func (div *addressDivisionInternal) getSplitRangeStringLength(rangeSeparator string, wildcard string, leadingZeroCount int, radix int, uppercase bool,
	splitDigitSeparator byte, reverseSplitDigits bool, stringPrefix string) int {
	return toUnsignedSplitRangeStringLength(
		div.getDivisionValue(),
		div.getUpperDivisionValue(),
		rangeSeparator,
		wildcard,
		leadingZeroCount,
		radix,
		uppercase,
		splitDigitSeparator,
		reverseSplitDigits,
		stringPrefix)
}

func (div *addressDivisionInternal) getRangeDigitCount(radix int) int {
	if !div.isMultiple() {
		return 0
	}
	if radix == 16 {
		prefix := div.GetMinPrefixLenForBlock()
		bitCount := div.GetBitCount()
		if prefix < bitCount && div.ContainsSinglePrefixBlock(prefix) {
			bitsPerCharacter := BitCount(4)
			if prefix%bitsPerCharacter == 0 {
				return int((bitCount - prefix) / bitsPerCharacter)
			}
		}
		return 0
	}
	value := div.getDivisionValue()
	upperValue := div.getUpperDivisionValue()
	maxValue := div.getMaxValue()
	factorRadix := DivInt(radix)
	factor := factorRadix
	numDigits := 1
	for {
		lowerRemainder := value % factor
		if lowerRemainder == 0 {
			//Consider in ipv4 the segment 24_
			//what does this mean?  It means 240 to 249 (not 240 to 245)
			//Consider 25_.  It means 250-255.
			//so the last digit ranges between 0-5 or 0-9 depending on whether the front matches the max possible front of 25.
			//If the front matches, the back ranges from 0 to the highest value of 255.
			//if the front does not match, the back must range across all values for the radix (0-9)
			var max DivInt
			if maxValue/factor == upperValue/factor {
				max = maxValue % factor
			} else {
				max = factor - 1
			}
			upperRemainder := upperValue % factor
			if upperRemainder == max {
				//whatever range there is must be accounted entirely by range digits, otherwise the range digits is 0
				//so here we check if that is the case
				if upperValue-upperRemainder == value {
					return numDigits
				} else {
					numDigits++
					factor *= factorRadix
					continue
				}
			}
		}
		return 0
	}
}

// if leadingZeroCount is -1, returns the number of leading zeros for maximum width, based on the width of the value
func (div *addressDivisionInternal) adjustLowerLeadingZeroCount(leadingZeroCount int, radix int) int {
	return div.adjustLeadingZeroCount(leadingZeroCount, div.getDivisionValue(), radix)
}

// if leadingZeroCount is -1, returns the number of leading zeros for maximum width, based on the width of the value
func (div *addressDivisionInternal) adjustUpperLeadingZeroCount(leadingZeroCount int, radix int) int {
	return div.adjustLeadingZeroCount(leadingZeroCount, div.getUpperDivisionValue(), radix)
}

func (div *addressDivisionInternal) adjustLeadingZeroCount(leadingZeroCount int, value DivInt, radix int) int {
	if leadingZeroCount < 0 {
		width := getDigitCount(value, radix)
		num := div.getMaxDigitCountRadix(radix) - width
		if num < 0 {
			return 0
		}
		return num
	}
	return leadingZeroCount
}

func (div *addressDivisionInternal) getDigitCount(radix int) int {
	if !div.isMultiple() && radix == div.getDefaultTextualRadix() { //optimization - just get the string, which is cached, which speeds up further calls to this or getString()
		return len(div.getWildcardString())
	}
	return getDigitCount(div.getUpperDivisionValue(), radix)
}

func (div *addressDivisionInternal) getMaxDigitCountRadix(radix int) int {
	return getMaxDigitCount(radix, div.GetBitCount(), div.getMaxValue())
}

// returns the number of digits for the maximum possible value of the division when using the default radix
func (div *addressDivisionInternal) getMaxDigitCount() int {
	return div.getMaxDigitCountRadix(div.getDefaultTextualRadix())
}

// returns the default radix for textual representations of addresses (10 for IPv4, 16 for IPv6, MAC and other)
func (div *addressDivisionInternal) getDefaultTextualRadix() int {
	addrType := div.getAddrType()
	if addrType.isIPv4() {
		return IPv4DefaultTextualRadix
	}
	return 16
}

// A simple string using just the lower value and the default radix.
func (div *addressDivisionInternal) getDefaultLowerString() string {
	return toDefaultString(div.getDivisionValue(), div.getDefaultTextualRadix())
}

// A simple string using just the lower and upper values and the default radix, separated by the default range character.
func (div *addressDivisionInternal) getDefaultRangeString() string {
	return div.getDefaultRangeStringVals(div.getDivisionValue(), div.getUpperDivisionValue(), div.getDefaultTextualRadix())
}

// getDefaultSegmentWildcardString() is the wildcard string to be used when producing the default strings with getString() or getWildcardString()
//
// Since no parameters for the string are provided, default settings are used, but they must be consistent with the address.
//
// For instance, generally the '*' is used as a wildcard to denote all possible values for a given segment,
// but in some cases that character is used for a segment separator.
//
// Note that this only applies to "default" settings, there are additional string methods that allow you to specify these separator characters.
// Those methods must be aware of the defaults as well, to know when they can defer to the defaults and when they cannot.
func (div *addressDivisionInternal) getDefaultSegmentWildcardString() string {
	if seg := div.toAddressDivision().ToSegmentBase(); seg != nil {
		return seg.getDefaultSegmentWildcardString()
	}
	return "" // for divisions, the width is variable and max values can change, so using wildcards make no sense
}

// getDefaultRangeSeparatorString() is the wildcard string to be used when producing the default strings with getString() or getWildcardString()
//
// Since no parameters for the string are provided, default settings are used, but they must be consistent with the address.
//
// For instance, generally the '-' is used as a range separator, but in some cases that character is used for a segment separator.
//
// Note that this only applies to "default" settings, there are additional string methods that allow you to specify these separator characters.
// Those methods must be aware of the defaults as well, to know when they can defer to the defaults and when they cannot.
func (div *addressDivisionInternal) getDefaultRangeSeparatorString() string {
	return "-"
}

//// only needed for godoc / pkgsite

// GetBitCount returns the number of bits in each value comprising this address item.
func (div *addressDivisionInternal) GetBitCount() BitCount {
	return div.addressDivisionBase.GetBitCount()
}

// GetByteCount returns the number of bytes required for each value comprising this address item,
// rounding up if the bit count is not a multiple of 8.
func (div *addressDivisionInternal) GetByteCount() int {
	return div.addressDivisionBase.GetByteCount()
}

// GetValue returns the lowest value in the address division range as a big integer.
func (div *addressDivisionInternal) GetValue() *BigDivInt {
	return div.addressDivisionBase.GetValue()
}

// GetUpperValue returns the highest value in the address division range as a big integer.
func (div *addressDivisionInternal) GetUpperValue() *BigDivInt {
	return div.addressDivisionBase.GetUpperValue()
}

// Bytes returns the lowest value in the address division range as a byte slice.
func (div *addressDivisionInternal) Bytes() []byte {
	return div.addressDivisionBase.Bytes()
}

// UpperBytes returns the highest value in the address division range as a byte slice.
func (div *addressDivisionInternal) UpperBytes() []byte {
	return div.addressDivisionBase.UpperBytes()
}

// CopyBytes copies the lowest value in the address division range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (div *addressDivisionInternal) CopyBytes(bytes []byte) []byte {
	return div.addressDivisionBase.CopyBytes(bytes)
}

// CopyUpperBytes copies the highest value in the address division range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (div *addressDivisionInternal) CopyUpperBytes(bytes []byte) []byte {
	return div.addressDivisionBase.CopyUpperBytes(bytes)
}

// IsZero returns whether this division matches exactly the value of zero.
func (div *addressDivisionInternal) IsZero() bool {
	return div.addressDivisionBase.IsZero()
}

// IncludesZero returns whether this item includes the value of zero within its range.
func (div *addressDivisionInternal) IncludesZero() bool {
	return div.addressDivisionBase.IncludesZero()
}

// IsMax returns whether this division matches exactly the maximum possible value, the value whose bits are all ones.
func (div *addressDivisionInternal) IsMax() bool {
	return div.addressDivisionBase.IsMax()
}

// IncludesMax returns whether this division includes the max value, the value whose bits are all ones, within its range.
func (div *addressDivisionInternal) IncludesMax() bool {
	return div.addressDivisionBase.IncludesMax()
}

// IsFullRange returns whether the division range includes all possible values for its bit length.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (div *addressDivisionInternal) IsFullRange() bool {
	return div.addressDivisionBase.IsFullRange()
}

func (div *addressDivisionInternal) compareSize(other AddressItem) int {
	return compareCount(div.toAddressDivision(), other)
}

//// end needed for godoc / pkgsite

// NewDivision creates a division of the given bit length, assigning it the given value.
// If the value's bit length exceeds the given bit length, it is truncated.
func NewDivision(val DivInt, bitCount BitCount) *AddressDivision {
	return NewRangePrefixDivision(val, val, nil, bitCount)
}

// NewRangeDivision creates a division of the given bit length, assigning it the given value range.
// If a value's bit length exceeds the given bit length, it is truncated.
func NewRangeDivision(val, upperVal DivInt, bitCount BitCount) *AddressDivision {
	return NewRangePrefixDivision(val, upperVal, nil, bitCount)
}

// NewPrefixDivision creates a division of the given bit length, assigning it the given value and prefix length.
// If the value's bit length exceeds the given bit length, it is truncated.
// If the prefix length exceeds the bit length, it is adjusted to the bit length.  If the prefix length is negative, it is adjusted to zero.
func NewPrefixDivision(val DivInt, prefixLen PrefixLen, bitCount BitCount) *AddressDivision {
	return NewRangePrefixDivision(val, val, prefixLen, bitCount)
}

// NewRangePrefixDivision creates a division of the given bit length, assigning it the given value range and prefix length.
// If a value's bit length exceeds the given bit length, it is truncated.
// If the prefix length exceeds the bit length, it is adjusted to the bit length.  If the prefix length is negative, it is adjusted to zero.
func NewRangePrefixDivision(val, upperVal DivInt, prefixLen PrefixLen, bitCount BitCount) *AddressDivision {
	return createAddressDivision(newDivValues(val, upperVal, prefixLen, bitCount))
}

// The following avoid the prefix length checks, value to BitCount checks, and low to high check inside newDivValues

func newDivision(val DivInt, bitCount BitCount) *AddressDivision {
	return newRangePrefixDivision(val, val, nil, bitCount)
}

func newRangeDivision(val, upperVal DivInt, bitCount BitCount) *AddressDivision {
	return newRangePrefixDivision(val, upperVal, nil, bitCount)
}

func newPrefixDivision(val DivInt, prefixLen PrefixLen, bitCount BitCount) *AddressDivision {
	return newRangePrefixDivision(val, val, prefixLen, bitCount)
}

func newRangePrefixDivision(val, upperVal DivInt, prefixLen PrefixLen, bitCount BitCount) *AddressDivision {
	return createAddressDivision(newDivValuesUnchecked(val, upperVal, prefixLen, bitCount))
}

// AddressDivision represents an arbitrary division in an address or address division grouping.
// It can contain a single value or a range of sequential values and it has an assigned bit length.
// Like all address components, it is immutable.
// Divisions that were converted from IPv4, IPv6 or MAC segments can be converted back to the same segment type and version.
// Divisions that were not converted from IPv4, IPv6 or MAC cannot be converted to segments.
type AddressDivision struct {
	addressDivisionInternal
}

//Note: many of the methods below are not public to addressDivisionInternal because segments have corresponding methods using segment values

// GetDivisionValue returns the lower division value in the range.
func (div *AddressDivision) GetDivisionValue() DivInt {
	return div.getDivisionValue()
}

// GetUpperDivisionValue returns the upper division value in the range.
func (div *AddressDivision) GetUpperDivisionValue() DivInt {
	return div.getUpperDivisionValue()
}

// IsMultiple returns  whether this division represents a sequential range of values, vs a single value.
func (div *AddressDivision) IsMultiple() bool {
	return div != nil && div.isMultiple()
}

// GetCount returns the count of possible distinct values for this division.
// If not representing multiple values, the count is 1.
//
// For instance, a division with the value range of 3-7 has count 5.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (div *AddressDivision) GetCount() *big.Int {
	if div == nil {
		return bigZero()
	}
	return div.getCount()
}

// Compare returns a negative integer, zero, or a positive integer if this address division is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (div *AddressDivision) Compare(item AddressItem) int {
	return CountComparator.Compare(div, item)
}

// CompareSize compares the counts of two items, the number of individual values within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one represents more individual values than another.
//
// CompareSize returns a positive integer if this division has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (div *AddressDivision) CompareSize(other AddressItem) int {
	if div == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return div.compareSize(other)
}

// Matches returns true if the division range matches the given single value.
func (div *AddressDivision) Matches(value DivInt) bool {
	return div.matches(value)
}

// MatchesWithMask applies the mask to this division and then compares the result with the given value,
// returning true if the range of the resulting division matches that single value.
func (div *AddressDivision) MatchesWithMask(value, mask DivInt) bool {
	return div.matchesWithMask(value, mask)
}

// MatchesValsWithMask applies the mask to this division and then compares the result with the given values,
// returning true if the range of the resulting division matches the given range.
func (div *AddressDivision) MatchesValsWithMask(lowerValue, upperValue, mask DivInt) bool {
	return div.matchesValsWithMask(lowerValue, upperValue, mask)
}

// GetMaxValue gets the maximum possible value for this type of division, determined by the number of bits.
//
// For the highest range value of this particular segment, use GetUpperDivisionValue.
func (div *AddressDivision) GetMaxValue() DivInt {
	return div.getMaxValue()
}

// IsSegmentBase returns true if this division originated as an address segment, and this can be converted back with ToSegmentBase.
func (div *AddressDivision) IsSegmentBase() bool {
	return div != nil && div.matchesSegment()
}

// IsIP returns true if this division originated as an IPv4 or IPv6 segment, or an implicitly zero-valued IP segment.  If so, use ToIP to convert back to the IP-specific type.
func (div *AddressDivision) IsIP() bool {
	return div != nil && div.matchesIPSegment()
}

// IsIPv4 returns true if this division originated as an IPv4 segment.  If so, use ToIPv4 to convert back to the IPv4-specific type.
func (div *AddressDivision) IsIPv4() bool {
	return div != nil && div.matchesIPv4Segment()
}

// IsIPv6 returns true if this division originated as an IPv6 segment.  If so, use ToIPv6 to convert back to the IPv6-specific type.
func (div *AddressDivision) IsIPv6() bool {
	return div != nil && div.matchesIPv6Segment()
}

// IsMAC returns true if this division originated as a MAC segment.  If so, use ToMAC to convert back to the MAC-specific type.
func (div *AddressDivision) IsMAC() bool {
	return div != nil && div.matchesMACSegment()
}

// ToIP converts to an IPAddressSegment if this division originated as an IPv4 or IPv6 segment, or an implicitly zero-valued IP segment.
// If not, ToIP returns nil.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (div *AddressDivision) ToIP() *IPAddressSegment {
	if div.IsIP() {
		return (*IPAddressSegment)(unsafe.Pointer(div))
	}
	return nil
}

// ToIPv4 converts to an IPv4AddressSegment if this division originated as an IPv4 segment.
// If not, ToIPv4 returns nil.
//
// ToIPv4 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (div *AddressDivision) ToIPv4() *IPv4AddressSegment {
	if div.IsIPv4() {
		return (*IPv4AddressSegment)(unsafe.Pointer(div))
	}
	return nil
}

// ToIPv6 converts to an IPv6AddressSegment if this division originated as an IPv6 segment.
// If not, ToIPv6 returns nil.
//
// ToIPv6 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (div *AddressDivision) ToIPv6() *IPv6AddressSegment {
	if div.IsIPv6() {
		return (*IPv6AddressSegment)(unsafe.Pointer(div))
	}
	return nil
}

// ToMAC converts to a MACAddressSegment if this division originated as a MAC segment.
// If not, ToMAC returns nil.
//
// ToMAC can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (div *AddressDivision) ToMAC() *MACAddressSegment {
	if div.IsMAC() {
		return (*MACAddressSegment)(unsafe.Pointer(div))
	}
	return nil
}

// ToSegmentBase converts to an AddressSegment if this division originated as a segment.
// If not, ToSegmentBase returns nil.
//
// ToSegmentBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (div *AddressDivision) ToSegmentBase() *AddressSegment {
	if div.IsSegmentBase() {
		return (*AddressSegment)(unsafe.Pointer(div))
	}
	return nil
}

// ToDiv is an identity method.
//
// ToDiv can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (div *AddressDivision) ToDiv() *AddressDivision {
	return div
}

// GetString produces a normalized string to represent the segment.
// If the segment is an IP segment string with CIDR network prefix block for its prefix length, then the string contains only the lower value of the block range.
// Otherwise, the explicit range will be printed.
// If the segment is not an IP segment, then the string is the same as that produced by GetWildcardString.
//
// The string returned is useful in the context of creating strings for address sections or full addresses,
// in which case the radix and bit-length can be deduced from the context.
// The String method produces strings more appropriate when no context is provided.
func (div *AddressDivision) GetString() string {
	if div == nil {
		return nilString()
	}
	return div.getString()
}

// GetWildcardString produces a normalized string to represent the segment, favouring wildcards and range characters regardless of any network prefix length.
// The explicit range of a range-valued segment will be printed.
//
// The string returned is useful in the context of creating strings for address sections or full addresses,
// in which case the radix and the bit-length can be deduced from the context.
// The String method produces strings more appropriate when no context is provided.
func (div *AddressDivision) GetWildcardString() string {
	if div == nil {
		return nilString()
	}
	return div.getWildcardString()
}

// String produces a string that is useful when a division string is provided with no context.
// It uses a string prefix for octal or hex ("0" or "0x"), and does not use the wildcard '*', because division size is variable, and so '*' is ambiguous.
// GetWildcardString is more appropriate in context with other segments or divisions.  It does not use a string prefix and uses '*' for full-range segments.
// GetString is more appropriate in context with prefix lengths, it uses zeros instead of wildcards for prefix block ranges.
func (div *AddressDivision) String() string {
	if div == nil {
		return nilString()
	}
	return div.toString()
}

func testRange(lowerValue, upperValue, finalUpperValue, networkMask, hostMask DivInt) bool {
	return lowerValue == (lowerValue&networkMask) && finalUpperValue == (upperValue|hostMask)
}

func divsSame(onePref, twoPref PrefixLen, oneVal, twoVal, oneUpperVal, twoUpperVal DivInt) bool {
	return onePref.Equal(twoPref) &&
		oneVal == twoVal && oneUpperVal == twoUpperVal
}

func divValsSame(oneVal, twoVal, oneUpperVal, twoUpperVal DivInt) bool {
	return oneVal == twoVal && oneUpperVal == twoUpperVal
}

func divValSame(oneVal, twoVal DivInt) bool {
	return oneVal == twoVal
}

func cacheStrPtr(cachedString **string, strPtr *string) {
	cachedVal := (*string)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(cachedString))))
	if cachedVal == nil {
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(cachedString))
		atomicStorePointer(dataLoc, unsafe.Pointer(strPtr))
	}
	return
}

func cacheStr(cachedString **string, stringer func() string) (str string) {
	cachedVal := (*string)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(cachedString))))
	if cachedVal == nil {
		str = stringer()
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(cachedString))
		atomicStorePointer(dataLoc, unsafe.Pointer(&str))
	} else {
		str = *cachedVal
	}
	return
}

func cacheStrErr(cachedString **string, stringer func() (string, addrerr.IncompatibleAddressError)) (str string, err addrerr.IncompatibleAddressError) {
	cachedVal := (*string)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(cachedString))))
	if cachedVal == nil {
		str, err = stringer()
		if err == nil {
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(cachedString))
			atomicStorePointer(dataLoc, unsafe.Pointer(&str))
		}
	} else {
		str = *cachedVal
	}
	return
}
