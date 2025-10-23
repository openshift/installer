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
	"strconv"
)

type boolSetting struct {
	isSet, val bool
}

var (
	falseVal = false
	trueVal  = true
)

// A PrefixLen indicates the length of the prefix for an address, section, division grouping, segment, or division.
// The zero value, which is nil, indicates that there is no prefix length.
type PrefixLen = *PrefixBitCount

// ToPrefixLen converts the given int to a prefix length
func ToPrefixLen(i int) PrefixLen {
	res := PrefixBitCount(i)
	return &res
}

// BitCount represents a count of bits in an address, section, grouping, segment, or division.
// Using signed integers allows for easier arithmetic, avoiding bugs.
// However, all methods adjust bit counts to match address size,
// so negative bit counts or bit counts larger than address size are meaningless.
type BitCount = int // using signed integers allows for easier arithmetic

const maxBitCountInternal, minBitCountInternal = math.MaxUint8, 0

// A PrefixBitCount is the count of bits in a non-nil PrefixLen.
// For arithmetic, you may wish to use the signed integer type BitCount instead, which you can get from a PrefixLen using the Len method.
type PrefixBitCount uint8

// Len returns the length of the prefix.  If the receiver is nil, representing the absence of a prefix length, returns 0.
// It will also return 0 if the receiver is a prefix with length of 0.  To distinguish the two, compare the receiver with nil.
func (prefixBitCount *PrefixBitCount) Len() BitCount {
	if prefixBitCount == nil {
		return 0
	}
	return prefixBitCount.bitCount()
}

// IsNil returns true if this is nil, meaning it represents having no prefix length, or the absence of a prefix length
func (prefixBitCount *PrefixBitCount) IsNil() bool {
	return prefixBitCount == nil
}

func (prefixBitCount *PrefixBitCount) bitCount() BitCount {
	return BitCount(*prefixBitCount)
}

func (prefixBitCount *PrefixBitCount) copy() PrefixLen {
	if prefixBitCount == nil {
		return nil
	}
	res := *prefixBitCount
	return &res
}

// Equal compares two PrefixLen values for equality.  This method is intended for the PrefixLen type.  BitCount values should be compared with the == operator.
func (prefixBitCount *PrefixBitCount) Equal(other PrefixLen) bool {
	if prefixBitCount == nil {
		return other == nil
	}
	return other != nil && prefixBitCount.bitCount() == other.bitCount()
}

// Matches compares a PrefixLen value with a bit count
func (prefixBitCount *PrefixBitCount) Matches(other BitCount) bool {
	return prefixBitCount != nil && prefixBitCount.bitCount() == other
}

// Compare compares PrefixLen values, returning -1, 0, or 1 if this prefix length is less than, equal to, or greater than the given prefix length.
// This method is intended for the PrefixLen type.  BitCount values should be compared with ==, >, <, >= and <= operators.
func (prefixBitCount *PrefixBitCount) Compare(other PrefixLen) int {
	if prefixBitCount == nil {
		if other == nil {
			return 0
		}
		return 1
	} else if other == nil {
		return -1
	}
	return prefixBitCount.bitCount() - other.bitCount()
}

// String returns the bit count as a base-10 positive integer string, or "<nil>" if the receiver is a nil pointer.
func (prefixBitCount *PrefixBitCount) String() string {
	if prefixBitCount == nil {
		return nilString()
	}
	return strconv.Itoa(prefixBitCount.bitCount())
}

// HostBitCount is the count of bits in a host.
// For arithmetic, you may wish to use the signed integer type BitCount instead, which you can get from a HostBitCount using the Len method.
type HostBitCount uint8

// BitsForCount returns the number of bits required outside the prefix length
// for a single prefix block to span at least as many addresses as the given count.
// Mathematically, it is the ceiling of the base 2 logarithm of the given count.
// A count of zero returns nil.
func BitsForCount(count uint64) (result *HostBitCount) {
	if count != 0 {
		var res HostBitCount
		countMinusOne := count - 1
		if (countMinusOne & (0xfff0000000000000)) != 0 { // conversion to float64 will fail
			count = (countMinusOne >> 53) + 1
			res = 53
		}
		res += HostBitCount(math.Ilogb(float64((count << 1) - 1)))
		return &res
	}
	return nil
}

// BlockSize is the reverse of BitsForCount, giving the total number of values when ranging across the number of host bits.
// The nil *HostBitCount returns 0.
func (hostBitCount *HostBitCount) BlockSize() *big.Int {
	if hostBitCount == nil {
		return bigZero()
	}
	return bigZero().Lsh(bigOneConst(), uint(*hostBitCount))
}

// Len returns the length of the host.  If the receiver is nil, representing the absence of a host length, returns 0.
// It will also return 0 if the receiver has a host length of 0.  To distinguish the two, compare the receiver with nil.
func (hostBitCount *HostBitCount) Len() BitCount {
	if hostBitCount == nil {
		return 0
	}
	return BitCount(*hostBitCount)
}

// String returns the bit count as a base-10 positive integer string, or "<nil>" if the receiver is a nil pointer.
func (hostBitCount *HostBitCount) String() string {
	if hostBitCount == nil {
		return nilString()
	}
	return strconv.Itoa(hostBitCount.Len())
}

// IsNil returns true if this is nil, meaning it represents having no identified host length.
func (hostBitCount *HostBitCount) IsNil() bool {
	return hostBitCount == nil
}

var cachedPrefixBitCounts, cachedPrefixLens = initPrefLens()

func initPrefLens() ([]PrefixBitCount, []PrefixLen) {
	cachedPrefBitcounts := make([]PrefixBitCount, maxBitCountInternal)
	cachedPrefLens := make([]PrefixLen, maxBitCountInternal)
	for i := 0; i <= IPv6BitCount; i++ {
		cachedPrefBitcounts[i] = PrefixBitCount(i)
		cachedPrefLens[i] = &cachedPrefBitcounts[i]
	}
	return cachedPrefBitcounts, cachedPrefLens
}

func cacheBitCount(i BitCount) PrefixLen {
	if i < minBitCountInternal {
		i = minBitCountInternal
	}
	if i < len(cachedPrefixBitCounts) {
		return &cachedPrefixBitCounts[i]
	}
	if i > maxBitCountInternal {
		i = maxBitCountInternal
	}
	res := PrefixBitCount(i)
	return &res
}

func cachePrefix(i BitCount) *PrefixLen {
	if i < minBitCountInternal {
		i = minBitCountInternal
	}
	if i < len(cachedPrefixLens) {
		return &cachedPrefixLens[i]
	}
	if i > maxBitCountInternal {
		i = maxBitCountInternal
	}
	val := PrefixBitCount(i)
	res := &val
	return &res
}

func cachePrefixLen(external PrefixLen) PrefixLen {
	if external == nil {
		return nil
	}
	return cacheBitCount(external.bitCount())
}

var p PrefixLen

func cacheNilPrefix() *PrefixLen {
	return &p
}

const maxPortNumInternal, minPortNumInternal = math.MaxUint16, 0

// Port represents the port of a UDP or TCP address.  A nil value indicates no port.
type Port = *PortNum

type PortInt = int // using signed integers allows for easier arithmetic

// PortNum is the port number for a non-nil Port.  For arithmetic, you might wish to use the signed integer type PortInt instead.
type PortNum uint16

func (portNum *PortNum) portNum() PortInt {
	return PortInt(*portNum)
}

func (portNum *PortNum) copy() Port {
	if portNum == nil {
		return nil
	}
	res := *portNum
	return &res
}

// Num converts to a PortPortIntNum, returning 0 if the receiver is nil.
func (portNum *PortNum) Num() PortInt {
	if portNum == nil {
		return 0
	}
	return PortInt(*portNum)
}

// Port dereferences this PortNum, while returning 0 if the receiver is nil.
func (portNum *PortNum) Port() PortNum {
	if portNum == nil {
		return 0
	}
	return *portNum
}

// Equal compares two Port values for equality.
func (portNum *PortNum) Equal(other Port) bool {
	if portNum == nil {
		return other == nil
	}
	return other != nil && portNum.portNum() == other.portNum()
}

// Matches compares a Port value with a port number.
func (portNum *PortNum) Matches(other PortInt) bool {
	return portNum != nil && portNum.portNum() == other
}

// Compare compares PrefixLen values, returning -1, 0, or 1 if the receiver is less than, equal to, or greater than the argument.
func (portNum *PortNum) Compare(other Port) int {
	if portNum == nil {
		if other == nil {
			return 0
		}
		return -1
	} else if other == nil {
		return 1
	}
	return portNum.portNum() - other.portNum()
}

// String returns the bit count as a base-10 positive integer string, or "<nil>" if the receiver is a nil pointer.
func (portNum *PortNum) String() string {
	if portNum == nil {
		return nilString()
	}
	return strconv.Itoa(portNum.portNum())
}

func cachePorts(i PortInt) Port {
	if i < minPortNumInternal {
		i = minPortNumInternal
	} else if i > maxPortNumInternal {
		i = maxPortNumInternal
	}
	res := PortNum(i)
	return &res
}

func bigOne() *big.Int {
	return big.NewInt(1)
}

var one = bigOne()

func bigOneConst() *big.Int {
	return one
}

func bigZero() *big.Int {
	return new(big.Int)
}

var zero = bigZero()

func bigZeroConst() *big.Int {
	return zero
}

var minusOne = big.NewInt(-1)

func bigMinusOneConst() *big.Int {
	return minusOne
}

func bigSixteen() *big.Int {
	return big.NewInt(16)
}

func bigIsZero(val *big.Int) bool {
	return len(val.Bits()) == 0 // slightly faster than div.value.BitLen() == 0
}

func bigIsOne(val *big.Int) bool {
	return bigAbsIsOne(val) && val.Sign() > 0
}

func bigAbsIsOne(val *big.Int) bool {
	bits := val.Bits()
	return len(bits) == 1 && bits[0] == 1
}

func bigIsNegative(val *big.Int) bool {
	return val.Sign() < 0
}

func bigIsNonPositive(val *big.Int) bool {
	return val.Sign() <= 0
}

func checkSubnet(item BitItem, prefixLength BitCount) BitCount {
	return checkBitCount(prefixLength, item.GetBitCount())
}

func checkDiv(div DivisionType, prefixLength BitCount) BitCount {
	return checkBitCount(prefixLength, div.GetBitCount())
}

func checkBitCount(prefixLength, max BitCount) BitCount {
	if prefixLength > max {
		return max
	} else if prefixLength < 0 {
		return 0
	}
	return prefixLength
}

func checkPrefLen(prefixLength PrefixLen, max BitCount) PrefixLen {
	if prefixLength != nil {
		prefLen := prefixLength.bitCount()
		if prefLen > max {
			return cacheBitCount(max)
		} else if prefLen < 0 {
			return cacheBitCount(0)
		}
	}
	return prefixLength

}

// wrapperIterator notifies the iterator to the right when wrapperIterator reaches its final value
type wrappedIterator struct {
	iterator   Iterator[*IPAddressSegment]
	finalValue []bool
	indexi     int
}

func (wrapped *wrappedIterator) HasNext() bool {
	return wrapped.iterator.HasNext()
}

func (wrapped *wrappedIterator) Next() *IPAddressSegment {
	iter := wrapped.iterator
	next := iter.Next()
	if !iter.HasNext() {
		wrapped.finalValue[wrapped.indexi+1] = true
	}
	return next
}
