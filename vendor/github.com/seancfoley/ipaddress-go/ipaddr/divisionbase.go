//
// Copyright 2020-2023 Sean C Foley
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
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

// divisionValuesBase provides an interface for divisions of any bit-size.
// It is shared by standard and large divisions.
// All the methods can be called for any division.
type divisionValuesBase interface {
	getBitCount() BitCount

	getByteCount() int

	// getDivisionPrefixLength provides the prefix length
	// if is aligned is true and the prefix is non-nil, any divisions that follow in the same grouping have a zero-length prefix
	getDivisionPrefixLength() PrefixLen

	// getValue gets the lower value as a BigDivInt
	getValue() *BigDivInt

	// getValue gets the upper value as a BigDivInt
	getUpperValue() *BigDivInt

	includesZero() bool

	includesMax() bool

	isMultiple() bool

	getCount() *big.Int

	// convert lower and upper values to byte arrays
	calcBytesInternal() (bytes, upperBytes []byte)

	bytesInternal(upper bool) (bytes []byte)

	// getCache returns a divCache for those divisions which cache their values, or nil otherwise
	getCache() *divCache

	getAddrType() addrType
}

// divisionValues provides methods to provide the values from divisions,
// and to create new divisions from values.
// Values may be truncated if the stored values in the interface implementation
// have larger bit-size than the return values.
// Similarly, values may be truncated if the supplied values have greater bit-size
// than the returned types.
type divisionValues interface {
	divisionValuesBase

	divIntVals

	divderiver

	segderiver

	segmentValues
}

type divCache struct {
	cachedString, cachedWildcardString, cached0xHexString, cachedHexString, cachedNormalizedString *string

	isSinglePrefBlock *bool
}

// addressDivisionBase is a division of any bit-size.
// It is shared by standard and large divisions types.
// Large divisions must not use the methods of divisionValues and use only the methods in divisionValuesBase.
type addressDivisionBase struct {
	// I've looked into making this divisionValuesBase.
	// If you do that, then to get access to the methods in divisionValues, you can either do type assertions like divisionValuesBase.(divisionValiues),
	// or you can add a method getDivisionValues to divisionValuesBase.
	// But in the end, either way you are assuming you know that divisionValuesBase is a divisionValues.  So no point.
	// Instead, each division type like IPAddressSegment and LargeDivision will know which value methods apply to that type.
	divisionValues
	// The field could possibly be generic.  However, since we aggregate implementations of divisionValues, what we have may be better
}

func (div *addressDivisionBase) getDivisionPrefixLength() PrefixLen {
	vals := div.divisionValues
	if vals == nil {
		return nil
	}
	return vals.getDivisionPrefixLength()
}

// GetBitCount returns the number of bits in each value comprising this address item.
func (div *addressDivisionBase) GetBitCount() BitCount {
	vals := div.divisionValues
	if vals == nil {
		return 0
	}
	return vals.getBitCount()
}

// GetByteCount returns the number of bytes required for each value comprising this address item,
// rounding up if the bit count is not a multiple of 8.
func (div *addressDivisionBase) GetByteCount() int {
	vals := div.divisionValues
	if vals == nil {
		return 0
	}
	return vals.getByteCount()
}

// GetValue returns the lowest value in the address division range as a big integer.
func (div *addressDivisionBase) GetValue() *BigDivInt {
	vals := div.divisionValues
	if vals == nil {
		return bigZero()
	}
	return vals.getValue()
}

// GetUpperValue returns the highest value in the address division range as a big integer.
func (div *addressDivisionBase) GetUpperValue() *BigDivInt {
	vals := div.divisionValues
	if vals == nil {
		return bigZero()
	}
	return vals.getUpperValue()
}

// Bytes returns the lowest value in the address division range as a byte slice.
func (div *addressDivisionBase) Bytes() []byte {
	if div.divisionValues == nil {
		return emptyBytes
	}
	return div.getBytes()
}

// UpperBytes returns the highest value in the address division range as a byte slice.
func (div *addressDivisionBase) UpperBytes() []byte {
	if div.divisionValues == nil {
		return emptyBytes
	}
	return div.getUpperBytes()
}

// CopyBytes copies the lowest value in the address division range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (div *addressDivisionBase) CopyBytes(bytes []byte) []byte {
	if div.divisionValues == nil {
		if bytes != nil {
			return bytes
		}
		return emptyBytes
	}
	cached := div.getBytes()
	return getBytesCopy(bytes, cached)
}

// CopyUpperBytes copies the highest value in the address division range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (div *addressDivisionBase) CopyUpperBytes(bytes []byte) []byte {
	if div.divisionValues == nil {
		if bytes != nil {
			return bytes
		}
		return emptyBytes
	}
	cached := div.getUpperBytes()
	return getBytesCopy(bytes, cached)
}

func (div *addressDivisionBase) getBytes() (bytes []byte) {
	return div.bytesInternal(false)
}

func (div *addressDivisionBase) getUpperBytes() (bytes []byte) {
	return div.bytesInternal(true)
}

func (div *addressDivisionBase) getCount() *big.Int {
	if !div.isMultiple() {
		return bigOne()
	}
	return div.divisionValues.getCount()
}

func (div *addressDivisionBase) isMultiple() bool {
	vals := div.divisionValues
	if vals == nil {
		return false
	}
	return vals.isMultiple()
}

// GetPrefixCountLen returns the count of the number of distinct values within the prefix part of the address item, the bits that appear within the prefix length.
func (div *addressDivisionBase) GetPrefixCountLen(prefixLength BitCount) *big.Int {
	if prefixLength < 0 {
		return bigOne()
	}
	bitCount := div.GetBitCount()
	if prefixLength >= bitCount {
		return div.getCount()
	}
	ushiftAdjustment := uint(bitCount - prefixLength)
	lower := div.GetValue()
	upper := div.GetUpperValue()
	upper.Rsh(upper, ushiftAdjustment)
	lower.Rsh(lower, ushiftAdjustment)
	upper.Sub(upper, lower).Add(upper, bigOneConst())
	return upper
}

// IsZero returns whether this division matches exactly the value of zero.
func (div *addressDivisionBase) IsZero() bool {
	return !div.isMultiple() && div.IncludesZero()
}

// IncludesZero returns whether this item includes the value of zero within its range.
func (div *addressDivisionBase) IncludesZero() bool {
	vals := div.divisionValues
	if vals == nil {
		return true
	}
	return vals.includesZero()
}

// IsMax returns whether this address matches exactly the maximum possible value, the value whose bits are all ones.
func (div *addressDivisionBase) IsMax() bool {
	return !div.isMultiple() && div.includesMax()
}

// IncludesMax returns whether this division includes the max value, the value whose bits are all ones, within its range.
func (div *addressDivisionBase) IncludesMax() bool {
	vals := div.divisionValues
	if vals == nil {
		return false
	}
	return vals.includesMax()
}

// IsFullRange returns whether the division range includes all possible values for its bit length.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (div *addressDivisionBase) IsFullRange() bool {
	return div.includesZero() && div.includesMax()
}

func (div *addressDivisionBase) getAddrType() addrType {
	vals := div.divisionValues
	if vals == nil {
		return zeroType
	}
	return vals.getAddrType()
}

func (div *addressDivisionBase) matchesStructure(other DivisionType) (res bool, addrType addrType) {
	addrType = div.getAddrType()
	if addrType != other.getAddrType() || (addrType.isZeroSegments() && (div.GetBitCount() != other.GetBitCount())) {
		return
	}
	res = true
	return
}

// toString produces a string that is useful when a division string is provided with no context.
// It uses a string prefix for octal or hex ("0" or "0x"), and does not use the wildcard '*', because division size is variable, so '*' is ambiguous.
// GetWildcardString() is more appropriate in context with other segments or divisions.  It does not use a string prefix and uses '*' for full-range segments.
// GetString() is more appropriate in context with prefix lengths, it uses zeros instead of wildcards for prefix block ranges.
func toString(div DivisionType) string { // this can be moved to addressDivisionBase when we have ContainsPrefixBlock and similar methods implemented for big.Int in the base.
	radix := div.getDefaultTextualRadix()
	var opts addrstr.IPStringOptions
	switch radix {
	case 16:
		opts = hexParamsDiv
	case 10:
		opts = decimalParamsDiv
	case 8:
		opts = octalParamsDiv
	default:
		opts = new(addrstr.IPStringOptionsBuilder).SetRadix(radix).SetWildcards(rangeWildcard).ToOptions()
	}
	return toStringOpts(opts, div)
}

func toStringOpts(opts addrstr.StringOptions, div DivisionType) string {
	builder := strings.Builder{}
	params := toParams(opts)
	builder.Grow(params.getDivisionStringLength(div))
	params.appendDivision(&builder, div)
	return builder.String()
}

func bigDivsSame(onePref, twoPref PrefixLen, oneVal, twoVal, oneUpperVal, twoUpperVal *BigDivInt) bool {
	return onePref.Equal(twoPref) &&
		oneVal.CmpAbs(twoVal) == 0 && oneUpperVal.CmpAbs(twoUpperVal) == 0
}

func bigDivValsSame(oneVal, twoVal, oneUpperVal, twoUpperVal *BigDivInt) bool {
	return oneVal.CmpAbs(twoVal) == 0 && oneUpperVal.CmpAbs(twoUpperVal) == 0
}

func bigDivValSame(oneVal, twoVal *big.Int) bool {
	return oneVal.CmpAbs(twoVal) == 0
}
