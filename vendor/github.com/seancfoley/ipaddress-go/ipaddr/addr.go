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
	"reflect"
	"unsafe"

	"github.com/seancfoley/bintree/tree"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

const (
	HexPrefix                    = "0x"
	OctalPrefix                  = "0"
	BinaryPrefix                 = "0b"
	RangeSeparator               = '-'
	RangeSeparatorStr            = "-"
	AlternativeRangeSeparator    = '\u00bb'
	AlternativeRangeSeparatorStr = "\u00bb" // '»'

	ExtendedDigitsRangeSeparatorStr = AlternativeRangeSeparatorStr

	SegmentWildcard    = '*'
	SegmentWildcardStr = "*"

	SegmentSqlWildcard          = '%'
	SegmentSqlWildcardStr       = "%"
	SegmentSqlSingleWildcard    = '_'
	SegmentSqlSingleWildcardStr = "_"

	//ExtendedDigitsRangeSeparator    = '\u00bb'
	//AlternativeSegmentWildcard  = '¿'
)

var segmentWildcardStr = SegmentWildcardStr

func createAddress(section *AddressSection, zone Zone) *Address {
	res := &Address{
		addressInternal{
			section: section,
			zone:    zone,
			cache:   &addressCache{},
		},
	}
	return res
}

// SegmentValueProvider provides values for segments.
// Values that fall outside the segment value type range will be truncated using standard golang integer type conversions https://golang.org/ref/spec#Conversions
type SegmentValueProvider func(segmentIndex int) SegInt

// AddressValueProvider provides values for addresses.
type AddressValueProvider interface {
	GetSegmentCount() int

	GetValues() SegmentValueProvider

	GetUpperValues() SegmentValueProvider
}

type addrsCache struct {
	lower, upper *Address
}

// identifierStr is a string representation of an address or host name.
type identifierStr struct {
	idStr HostIdentifierString // MACAddressString or IPAddressString or HostName
}

type addressCache struct {
	addrsCache *addrsCache

	stringCache *stringCache // only used by IPv6 when there is a zone

	identifierStr *identifierStr

	trieKeyCache *tree.TrieKeyData
}

type addressInternal struct {
	section *AddressSection
	zone    Zone
	cache   *addressCache
}

func (addr *addressInternal) assignTrieCache() {
	cache := addr.cache
	if cache != nil && cache.trieKeyCache != nil {
		cache.trieKeyCache = addr.constructTrieCache()
	}
}

func (addr *addressInternal) constructTrieCache() *tree.TrieKeyData {
	sect := addr.section
	prefLen := sect.getPrefixLen()
	var cache *tree.TrieKeyData
	if sectionIPv4 := sect.ToIPv4(); sectionIPv4 != nil {
		cache = &tree.TrieKeyData{
			Is32Bits:  true,
			PrefLen:   tree.PrefixLen(prefLen),
			Uint32Val: sectionIPv4.Uint32Value(),
		}
		if prefLen != nil {
			bits := prefLen.bitCount()
			cache.NextBitMask32Val = uint32(0x80000000) >> bits
			cache.Mask32Val = ipv4NetworkMasks[bits]
		}
	} else if sectionIPv6 := sect.ToIPv6(); sectionIPv6 != nil {
		cache = &tree.TrieKeyData{
			Is128Bits: true,
			PrefLen:   tree.PrefixLen(prefLen),
		}
		cache.Uint64HighVal, cache.Uint64LowVal = sectionIPv6.Uint64Values()
		if prefLen != nil {
			bits := prefLen.bitCount()
			mask := ipv6NetworkMasks[bits]
			cache.Mask64HighVal, cache.Mask64LowVal = mask[0], mask[1]
			if bits > 63 {
				cache.NextBitMask64Val = uint64(0x8000000000000000) >> (bits - 64)
			} else {
				cache.NextBitMask64Val = uint64(0x8000000000000000) >> bits
			}
		}
	} else {
		cache = &tree.TrieKeyData{}
	}
	return cache
}

func (addr *addressInternal) getTrieCache() *tree.TrieKeyData {
	cache := addr.cache
	if cache != nil {
		cached := (*tree.TrieKeyData)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.trieKeyCache))))
		if cached == nil {
			cached = addr.constructTrieCache()
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.trieKeyCache))
			atomicStorePointer(dataLoc, unsafe.Pointer(cached))
		}
		return cached
	}
	return &tree.TrieKeyData{}
}

// GetBitCount returns the number of bits comprising this address,
// or each address in the range if a subnet.
func (addr *addressInternal) GetBitCount() BitCount {
	section := addr.section
	if section == nil {
		return 0
	}
	return section.GetBitCount()
}

// GetByteCount returns the number of bytes required for this address,
// or each address in the range if a subnet.
func (addr *addressInternal) GetByteCount() int {
	section := addr.section
	if section == nil {
		return 0
	}
	return section.GetByteCount()
}

func (addr *addressInternal) getBytes() []byte {
	return addr.section.getBytes()
}

func (addr *addressInternal) getUpperBytes() []byte {
	return addr.section.getUpperBytes()
}

func (addr *addressInternal) getTrailingBitCount(ones bool) BitCount {
	return addr.section.GetTrailingBitCount(ones)
}

func (addr *addressInternal) getLeadingBitCount(ones bool) BitCount {
	return addr.section.GetLeadingBitCount(ones)
}

func (addr *addressInternal) getCount() *big.Int {
	section := addr.section
	if section == nil {
		return bigOne()
	}
	return section.GetCount()
}

// GetPrefixCount returns the count of prefixes in this address or subnet.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the count of the range of values in the prefix.
//
// If this has a nil prefix length, returns the same value as GetCount.
func (addr *addressInternal) GetPrefixCount() *big.Int {
	section := addr.section
	if section == nil {
		return bigOne()
	}
	return section.GetPrefixCount()
}

// GetPrefixCountLen returns the count of prefixes in this address or subnet for the given prefix length.
//
// If not a subnet of multiple addresses, or a subnet with just single prefix of the given length, returns 1.
func (addr *addressInternal) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	section := addr.section
	if section == nil {
		return bigOne()
	}
	return section.GetPrefixCountLen(prefixLen)
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) segments.
func (addr *addressInternal) GetBlockCount(segments int) *big.Int {
	section := addr.section
	if section == nil {
		return bigOne()
	}
	return section.GetBlockCount(segments)
}

// testBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this address.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (addr *addressInternal) testBit(n BitCount) bool {
	return addr.section.TestBit(n)
}

// isOneBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the most significant bit.
// isOneBit will panic if bitIndex is less than zero or larger than the bit count of this item.
func (addr *addressInternal) isOneBit(bitIndex BitCount) bool {
	return addr.section.IsOneBit(bitIndex)
}

// isMultiple returns true if this address represents more than a single individual address, whether it is a subnet of multiple addresses.
func (addr *addressInternal) isMultiple() bool {
	return addr.section != nil && addr.section.isMultiple()
}

// isPrefixed returns whether this address has an associated prefix length.
func (addr *addressInternal) isPrefixed() bool {
	return addr.section.IsPrefixed()
}

// GetPrefixLen returns the prefix length, or nil if there is no prefix length.
//
// A prefix length indicates the number of bits in the initial part (most significant bits) of the address that comprise the prefix.
//
// A prefix is a part of the address that is not specific to that address but common amongst a group of addresses, such as a CIDR prefix block subnet.
//
// For IP addresses, the prefix is explicitly defined when the address is created. For example, "1.2.0.0/16" has a prefix length of 16, while "1.2.*.*" has no prefix length,
// even though they both represent the same set of addresses and are considered equal.  Prefixes can be considered variable for a given IP address and can depend on routing.
//
// The methods GetMinPrefixLenForBlock and GetPrefixLenForSingleBlock can help you to obtain or define a prefix length if one does not exist already.
// The method ToPrefixBlockLen allows you to create the subnet consisting of the block of addresses for any given prefix length.
//
// For MAC addresses, the prefix is initially inferred from the range, so "1:2:3:*:*:*" has a prefix length of 24.
// MAC addresses derived from an address with a prefix length may retain the prefix length regardless of their own range of values.
func (addr *addressInternal) GetPrefixLen() PrefixLen {
	return addr.getPrefixLen().copy()
}

func (addr *addressInternal) getPrefixLen() PrefixLen {
	if addr.section == nil {
		return nil
	}
	return addr.section.getPrefixLen()
}

// IsSinglePrefixBlock returns whether the address range matches the block of values for a single prefix identified by the prefix length of this address.
// This is similar to IsPrefixBlock except that it returns false when the subnet has multiple prefixes.
//
// What distinguishes this method from ContainsSinglePrefixBlock is that this method returns
// false if the series does not have a prefix length assigned to it,
// or a prefix length that differs from the prefix length for which ContainsSinglePrefixBlock returns true.
//
// It is similar to IsPrefixBlock but returns false when there are multiple prefixes.
//
// For instance, "1.*.*.* /16" returns false from this method and returns true from IsPrefixBlock.
func (addr *addressInternal) IsSinglePrefixBlock() bool {
	prefLen := addr.getPrefixLen()
	return prefLen != nil && addr.section.IsSinglePrefixBlock()
}

// IsPrefixBlock returns whether the address has a prefix length and the address range includes the block of values for that prefix length.
// If the prefix length matches the bit count, this returns true.
//
// To create a prefix block from any address, use ToPrefixBlock.
//
// This is different from ContainsPrefixBlock in that this method returns
// false if the series has no prefix length, or a prefix length that differs from a prefix length for which ContainsPrefixBlock returns true.
func (addr *addressInternal) IsPrefixBlock() bool {
	prefLen := addr.getPrefixLen()
	return prefLen != nil && addr.section.ContainsPrefixBlock(prefLen.bitCount())
}

// ContainsPrefixBlock returns whether the range of this address or subnet contains the block of addresses for the given prefix length.
//
// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values in this item for the given prefix length makes no difference.
//
// Use GetMinPrefixLenForBlock to determine the smallest prefix length for which this method returns true.
func (addr *addressInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	return addr.section == nil || addr.section.ContainsPrefixBlock(prefixLen)
}

// ContainsSinglePrefixBlock returns whether this address contains a single prefix block for the given prefix length.
//
// This means there is only one prefix value for the given prefix length, and it also contains the full prefix block for that prefix, all addresses with that prefix.
//
// Use GetPrefixLenForSingleBlock to determine whether there is a prefix length for which this method returns true.
func (addr *addressInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	return addr.section == nil || addr.section.ContainsSinglePrefixBlock(prefixLen)
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this includes the block of addresses for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this represents just a single address, returns the bit length of this address.
func (addr *addressInternal) GetMinPrefixLenForBlock() BitCount {
	section := addr.section
	if section == nil {
		return 0
	}
	return section.GetMinPrefixLenForBlock()
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this address subnet matches exactly the block of addresses for that prefix.
//
// If the range can be described this way, then this method returns the same value as GetMinPrefixLenForBlock.
//
// If no such prefix exists, returns nil.
//
// If this segment grouping represents a single value, returns the bit length of this address.
//
// IP address examples:
//   - 1.2.3.4 returns 32
//   - 1.2.3.4/16 returns 32
//   - 1.2.*.* returns 16
//   - 1.2.*.0/24 returns 16
//   - 1.2.0.0/16 returns 16
//   - 1.2.*.4 returns nil
//   - 1.2.252-255.* returns 22
func (addr *addressInternal) GetPrefixLenForSingleBlock() PrefixLen {
	section := addr.section
	if section == nil {
		return cacheBitCount(0)
	}
	return section.GetPrefixLenForSingleBlock()
}

// In callers, we always need to ensure init is called, otherwise a nil section will be zero-size instead of having size one.
func (addr *addressInternal) compareSize(other AddressItem) int {
	return addr.section.compareSize(other)
}

func (addr *addressInternal) trieCompare(other *Address) int {
	if addr.toAddress() == other {
		return 0
	}
	segmentCount := addr.getDivisionCount()
	bitsPerSegment := addr.GetBitsPerSegment()
	o1Pref := addr.GetPrefixLen()
	o2Pref := other.GetPrefixLen()
	bitsMatchedSoFar := 0
	i := 0
	for {
		segment1 := addr.getSegment(i)
		segment2 := other.getSegment(i)
		pref1 := getSegmentPrefLen(addr.toAddress(), o1Pref, bitsPerSegment, bitsMatchedSoFar, segment1)
		pref2 := getSegmentPrefLen(other, o2Pref, bitsPerSegment, bitsMatchedSoFar, segment2)
		if pref1 != nil {
			segmentPref1 := pref1.Len()
			segmentPref2 := pref2.Len()
			if pref2 != nil && segmentPref2 <= segmentPref1 {
				matchingBits := getMatchingBits(segment1, segment2, segmentPref2, bitsPerSegment)
				if matchingBits >= segmentPref2 {
					if segmentPref2 == segmentPref1 {
						// same prefix block
						return 0
					}
					// segmentPref2 is shorter prefix, prefix bits match, so depends on bit at index segmentPref2
					if segment1.IsOneBit(segmentPref2) {
						return 1
					}
					return -1
				}
				return compareSegInt(segment1.GetSegmentValue(), segment2.GetSegmentValue())
			} else {
				matchingBits := getMatchingBits(segment1, segment2, segmentPref1, bitsPerSegment)
				if matchingBits >= segmentPref1 {
					if segmentPref1 < bitsPerSegment {
						if segment2.IsOneBit(segmentPref1) {
							return -1
						}
						return 1
					} else {
						i++
						if i == segmentCount {
							return 1 // o1 with prefix length matching bit count is the bigger
						} // else must check the next segment
					}
				} else {
					return compareSegInt(segment1.GetSegmentValue(), segment2.GetSegmentValue())
				}
			}
		} else if pref2 != nil {
			segmentPref2 := pref2.Len()
			matchingBits := getMatchingBits(segment1, segment2, segmentPref2, bitsPerSegment)
			if matchingBits >= segmentPref2 {
				if segmentPref2 < bitsPerSegment {
					if segment1.IsOneBit(segmentPref2) {
						return 1
					}
					return -1
				} else {
					i++
					if i == segmentCount {
						return -1 // o2 with prefix length matching bit count is the bigger
					} // else must check the next segment
				}
			} else {
				return compareSegInt(segment1.GetSegmentValue(), segment2.GetSegmentValue())
			}
		} else {
			matchingBits := getMatchingBits(segment1, segment2, bitsPerSegment, bitsPerSegment)
			if matchingBits < bitsPerSegment { // no match - the current subnet/address is not here
				return compareSegInt(segment1.GetSegmentValue(), segment2.GetSegmentValue())
			} else {
				i++
				if i == segmentCount {
					// same address
					return 0
				} // else must check the next segment
			}
		}
		bitsMatchedSoFar += bitsPerSegment
	}
}

func trieIncrement[T TrieKeyConstraint[T]](addr T) (t T, ok bool) {
	if res, ok := tree.TrieIncrement(createKey(addr)); ok {
		return res.address, true
	}
	return
}

func trieDecrement[T TrieKeyConstraint[T]](addr T) (t T, ok bool) {
	if res, ok := tree.TrieDecrement(createKey(addr)); ok {
		return res.address, true
	}
	return
}

func (addr *addressInternal) toString() string {
	section := addr.section
	if section == nil {
		return nilSection() // note no zone possible since a zero-address like Address{} or IPAddress{} cannot have a zone
	} else if addr.isMAC() {
		return addr.toNormalizedString()
	}
	return addr.toCanonicalString()
}

// IsSequential returns whether the address or subnet represents a range of addresses that are sequential.
//
// Generally, for a subnet this means that any segment covering a range of values must be followed by segments that are full range, covering all values.
//
// Individual addresses are sequential and CIDR prefix blocks are sequential.
// The subnet "1.2.3-4.5" is not sequential, since the two addresses it represents, "1.2.3.5" and "1.2.4.5", are not ("1.2.3.6" is in-between the two but not in the subnet).
//
// With any IP address subnet, you can use SequentialBlockIterator to convert any subnet to a collection of sequential subnets.
func (addr *addressInternal) IsSequential() bool {
	section := addr.section
	if section == nil {
		return true
	}
	return section.IsSequential()
}

func (addr *addressInternal) getSegment(index int) *AddressSegment {
	return addr.section.GetSegment(index)
}

// GetBitsPerSegment returns the number of bits comprising each segment in this address or subnet.  Segments in the same address are equal length.
func (addr *addressInternal) GetBitsPerSegment() BitCount {
	section := addr.section
	if section == nil {
		return 0
	}
	return section.GetBitsPerSegment()
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this address or subnet.  Segments in the same address are equal length.
func (addr *addressInternal) GetBytesPerSegment() int {
	section := addr.section
	if section == nil {
		return 0
	}
	return section.GetBytesPerSegment()
}

func (addr *addressInternal) getMaxSegmentValue() SegInt {
	return addr.section.GetMaxSegmentValue()
}

func (addr *addressInternal) checkIdentity(section *AddressSection) *Address {
	if section == nil {
		return nil
	} else if section == addr.section {
		return addr.toAddress()
	}
	return createAddress(section, addr.zone)
}

func (addr *addressInternal) getLower() *Address {
	lower, _ := addr.getLowestHighestAddrs()
	return lower
}

func (addr *addressInternal) getUpper() *Address {
	_, upper := addr.getLowestHighestAddrs()
	return upper
}

func (addr *addressInternal) getLowestHighestAddrs() (lower, upper *Address) {
	if !addr.isMultiple() {
		lower = addr.toAddress()
		upper = lower
		return
	}
	cache := addr.cache
	if cache == nil {
		return addr.createLowestHighestAddrs()
	}
	cached := (*addrsCache)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.addrsCache))))
	if cached == nil {
		cached = &addrsCache{}
		cached.lower, cached.upper = addr.createLowestHighestAddrs()
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.addrsCache))
		atomicStorePointer(dataLoc, unsafe.Pointer(cached))
	}
	lower, upper = cached.lower, cached.upper
	return
}

func (addr *addressInternal) createLowestHighestAddrs() (lower, upper *Address) {
	lower = addr.checkIdentity(addr.section.GetLower())
	upper = addr.checkIdentity(addr.section.GetUpper())
	return
}

func (addr *addressInternal) toMaxLower() *Address {
	section := addr.section
	if section == nil {
		return addr.toAddress()
	}
	return addr.checkIdentity(addr.section.toMaxLower())
}

func (addr *addressInternal) toMinUpper() *Address {
	section := addr.section
	if section == nil {
		return addr.toAddress()
	}
	return addr.checkIdentity(addr.section.toMinUpper())
}

// IsZero returns whether this address matches exactly the value of zero.
func (addr *addressInternal) IsZero() bool {
	section := addr.section
	if section == nil {
		return true
	}
	return section.IsZero()
}

// IncludesZero returns whether this address includes the zero address within its range.
func (addr *addressInternal) IncludesZero() bool {
	section := addr.section
	if section == nil {
		return true
	}
	return section.IncludesZero()
}

// IsFullRange returns whether this address covers the entire address space of this address version or type.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (addr *addressInternal) IsFullRange() bool {
	section := addr.section
	if section == nil {
		// when no bits, the only value 0 is the max value too
		return true
	}
	return section.IsFullRange()
}

func (addr *addressInternal) toAddress() *Address {
	return (*Address)(unsafe.Pointer(addr))
}

func (addr *addressInternal) getDivision(index int) *AddressDivision {
	return addr.section.getDivision(index)
}

func (addr *addressInternal) getDivisionCount() int {
	if addr.section == nil {
		return 0
	}
	return addr.section.GetDivisionCount()
}

func (addr *addressInternal) getDivisionsInternal() []*AddressDivision {
	return addr.section.getDivisionsInternal()
}

func (addr *addressInternal) toPrefixBlock() *Address {
	return addr.checkIdentity(addr.section.toPrefixBlock())
}

func (addr *addressInternal) toPrefixBlockLen(prefLen BitCount) *Address {
	return addr.checkIdentity(addr.section.toPrefixBlockLen(prefLen))
}

func (addr *addressInternal) toBlock(segmentIndex int, lower, upper SegInt) *Address {
	return addr.checkIdentity(addr.section.toBlock(segmentIndex, lower, upper))
}

func (addr *addressInternal) reverseBytes() (*Address, addrerr.IncompatibleAddressError) {
	sect, err := addr.section.ReverseBytes()
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(sect), nil
}

func (addr *addressInternal) reverseBits(perByte bool) (*Address, addrerr.IncompatibleAddressError) {
	sect, err := addr.section.ReverseBits(perByte)
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(sect), nil
}

// reverseSegments returns a new address with the segments reversed.
func (addr *addressInternal) reverseSegments() *Address {
	return addr.checkIdentity(addr.section.ReverseSegments())
}

// isIPv4 returns whether this matches an IPv4 address.
// we allow nil receivers to allow this to be called following a failed conversion like ToIP()
func (addr *addressInternal) isIPv4() bool {
	return addr.section != nil && addr.section.matchesIPv4AddressType()
}

// isIPv6 returns whether this matches an IPv6 address.
// we allow nil receivers to allow this to be called following a failed conversion like ToIP()
func (addr *addressInternal) isIPv6() bool {
	return addr.section != nil && addr.section.matchesIPv6AddressType()
}

// isIPv6 returns whether this matches an IPv6 address.
// we allow nil receivers to allow this to be called following a failed conversion like ToIP()
func (addr *addressInternal) isMAC() bool {
	return addr.section != nil && addr.section.matchesMACAddressType()
}

// isIP returns whether this matches an IP address.
// It must be IPv4, IPv6, or the zero IPAddress which has no segments
// we allow nil receivers to allow this to be called following a failed conversion like ToIP()
func (addr *addressInternal) isIP() bool {
	return addr.section == nil /* zero addr */ || addr.section.matchesIPAddressType()
}

func (addr *addressInternal) prefixEquals(other AddressType) bool {
	otherAddr := other.ToAddressBase()
	if addr.toAddress() == otherAddr {
		return true
	}
	otherSection := otherAddr.GetSection()
	if addr.section == nil {
		return otherSection.GetSegmentCount() == 0
	}
	return addr.section.PrefixEqual(otherSection) &&
		// if it is IPv6 and has a zone, then it does not contain addresses from other zones
		addr.isSameZone(otherAddr)
}

func (addr *addressInternal) prefixContains(other AddressType) bool {
	otherAddr := other.ToAddressBase()
	if addr.toAddress() == otherAddr {
		return true
	}
	otherSection := otherAddr.GetSection()
	if addr.section == nil {
		return otherSection.GetSegmentCount() == 0
	}
	return addr.section.PrefixContains(otherSection) &&
		// if it is IPv6 and has a zone, then it does not contain addresses from other zones
		addr.isSameZone(otherAddr)
}

func (addr *addressInternal) contains(other AddressType) bool {
	if other == nil {
		return true
	}
	otherAddr := other.ToAddressBase()
	if addr.toAddress() == otherAddr || otherAddr == nil {
		return true
	}
	otherSection := otherAddr.GetSection()
	return addr.section.Contains(otherSection) &&
		// if it is IPv6 and has a zone, then it does not contain addresses from other zones
		addr.isSameZone(otherAddr)
}

// Returns whether this is same type and version of the given address and whether it overlaps with the values in the given address or subnet
func (addr *addressInternal) overlaps(other AddressType) bool {
	if other == nil {
		return true
	}
	otherAddr := other.ToAddressBase()
	if addr.toAddress() == otherAddr || otherAddr == nil {
		return true
	}
	otherSection := otherAddr.GetSection()
	if addr.section == nil {
		return otherSection.GetSegmentCount() == 0
	}
	return addr.section.Overlaps(otherSection) &&
		// if it is IPv6 and has a zone, then it does not overlap addresses from other zones
		addr.isSameZone(otherAddr)
}

func (addr *addressInternal) equals(other AddressType) bool {
	if other == nil {
		return false
	}
	otherAddr := other.ToAddressBase()
	if addr.toAddress() == otherAddr {
		return true
	} else if otherAddr == nil {
		return false
	}
	otherSection := otherAddr.GetSection()
	if addr.section == nil {
		return otherSection.GetSegmentCount() == 0
	}
	return addr.section.Equal(otherSection) &&
		// if it it is IPv6 and has a zone, then it does not equal addresses from other zones
		addr.isSameZone(otherAddr)
}

// returns whether two addresses, already known to be the same version and address type, are equal
func (addr *addressInternal) equalsSameVersion(other AddressType) bool {
	otherAddr := other.ToAddressBase()
	if addr.toAddress() == otherAddr {
		return true
	} else if otherAddr == nil {
		return false
	}
	otherSection := otherAddr.GetSection()
	return addr.section.sameCountTypeEquals(otherSection) &&
		// if it it is IPv6 and has a zone, then it does not equal addresses from other zones
		addr.isSameZone(otherAddr)
}

// withoutPrefixLen returns the same address but with no associated prefix length.
func (addr *addressInternal) withoutPrefixLen() *Address {
	return addr.checkIdentity(addr.section.withoutPrefixLen())
}

func (addr *addressInternal) adjustPrefixLen(prefixLen BitCount) *Address {
	return addr.checkIdentity(addr.section.adjustPrefixLen(prefixLen))
}

func (addr *addressInternal) adjustPrefixLenZeroed(prefixLen BitCount) (res *Address, err addrerr.IncompatibleAddressError) {
	section, err := addr.section.adjustPrefixLenZeroed(prefixLen)
	if err == nil {
		res = addr.checkIdentity(section)
	}
	return
}

func (addr *addressInternal) setPrefixLen(prefixLen BitCount) *Address {
	return addr.checkIdentity(addr.section.setPrefixLen(prefixLen))
}

func (addr *addressInternal) setPrefixLenZeroed(prefixLen BitCount) (res *Address, err addrerr.IncompatibleAddressError) {
	section, err := addr.section.setPrefixLenZeroed(prefixLen)
	if err == nil {
		res = addr.checkIdentity(section)
	}
	return
}

func (addr *addressInternal) assignPrefixForSingleBlock() *Address {
	newPrefix := addr.GetPrefixLenForSingleBlock()
	if newPrefix == nil {
		return nil
	}
	return addr.checkIdentity(addr.section.setPrefixLen(newPrefix.bitCount()))
}

// assignMinPrefixForBlock constructs an equivalent address section with the smallest CIDR prefix possible (largest network),
// such that the range of values are a set of subnet blocks for that prefix.
func (addr *addressInternal) assignMinPrefixForBlock() *Address {
	return addr.setPrefixLen(addr.GetMinPrefixLenForBlock())
}

// toSingleBlockOrAddress converts to a single prefix block or address.
// If the given address is a single prefix block, it is returned.
// If it can be converted to a single prefix block by assigning a prefix length, the converted block is returned.
// If it is a single address, any prefix length is removed and the address is returned.
// Otherwise, nil is returned.
func (addr *addressInternal) toSinglePrefixBlockOrAddr() *Address {
	if !addr.isMultiple() {
		if !addr.isPrefixed() {
			return addr.toAddress()
		}
		return addr.withoutPrefixLen()
		//} else if addr.IsSinglePrefixBlock() {
		//	return addr.toAddress()
	} else {
		series := addr.assignPrefixForSingleBlock()
		if series != nil {
			return series
		}
	}
	return nil
}

func (addr *addressInternal) isSameZone(other *Address) bool {
	return addr.zone == other.ToAddressBase().zone
}

// Be careful when calling this, because for IPv6Address{} and IPv4Address{} it gives the wrong answer without the init method called first.
// An alternative is to call GetIPVersion, or to be sure you have called the init method.
func (addr *addressInternal) getAddrType() addrType {
	if addr.section == nil {
		return zeroType
	}
	return addr.section.addrType
}

// equivalent to section.sectionIterator
func (addr *addressInternal) addrIterator(excludeFunc func([]*AddressDivision) bool) Iterator[*Address] {
	useOriginal := !addr.isMultiple()
	original := addr.toAddress()
	var iterator Iterator[[]*AddressDivision]
	if useOriginal {
		if excludeFunc != nil && excludeFunc(addr.getDivisionsInternal()) {
			original = nil // the single-valued iterator starts out empty
		}
	} else {
		address := addr.toAddress()
		iterator = allSegmentsIterator(
			addr.getDivisionCount(),
			nil,
			func(index int) Iterator[*AddressSegment] { return address.getSegment(index).iterator() },
			excludeFunc)
	}
	return addrIterator(
		useOriginal,
		original,
		original.getPrefixLen(),
		false,
		iterator)
}

func (addr *addressInternal) prefixIterator(isBlockIterator bool) Iterator[*Address] {
	prefLen := addr.getPrefixLen()
	if prefLen == nil {
		return addr.addrIterator(nil)
	}
	var useOriginal bool
	if isBlockIterator {
		useOriginal = addr.IsSinglePrefixBlock()
	} else {
		useOriginal = bigIsOne(addr.GetPrefixCount())
	}
	prefLength := prefLen.bitCount()
	bitsPerSeg := addr.GetBitsPerSegment()
	bytesPerSeg := addr.GetBytesPerSegment()
	networkSegIndex := getNetworkSegmentIndex(prefLength, bytesPerSeg, bitsPerSeg)
	hostSegIndex := getHostSegmentIndex(prefLength, bytesPerSeg, bitsPerSeg)
	segCount := addr.getDivisionCount()
	var iterator Iterator[[]*AddressDivision]
	address := addr.toAddress()
	if !useOriginal {
		var hostSegIteratorProducer func(index int) Iterator[*AddressSegment]
		if isBlockIterator {
			hostSegIteratorProducer = func(index int) Iterator[*AddressSegment] {
				seg := address.getSegment(index)
				if seg.isPrefixed() { // IP address segments know their own prefix, MAC segments do not
					return seg.prefixBlockIterator()
				}
				segPref := getPrefixedSegmentPrefixLength(bitsPerSeg, prefLength, index)
				return seg.prefixedBlockIterator(segPref.bitCount())
			}
		} else {
			hostSegIteratorProducer = func(index int) Iterator[*AddressSegment] {
				seg := address.getSegment(index)
				if seg.isPrefixed() { // IP address segments know their own prefix, MACS segments do not
					return seg.prefixIterator()
				}
				segPref := getPrefixedSegmentPrefixLength(bitsPerSeg, prefLength, index)
				return seg.prefixedIterator(segPref.bitCount())
			}
		}
		iterator = segmentsIterator(
			segCount,
			nil, //when no prefix we defer to other iterator, when there is one we use the whole original section in the encompassing iterator and not just the original segments
			func(index int) Iterator[*AddressSegment] { return address.getSegment(index).iterator() },
			nil,
			networkSegIndex,
			hostSegIndex,
			hostSegIteratorProducer)
	}
	if isBlockIterator {
		return addrIterator(
			useOriginal,
			address,
			address.getPrefixLen(),
			prefLength < addr.GetBitCount(),
			iterator)
	}
	return prefixAddrIterator(
		useOriginal,
		address,
		address.getPrefixLen(),
		iterator)
}

func (addr *addressInternal) blockIterator(segmentCount int) Iterator[*Address] {
	if segmentCount < 0 {
		segmentCount = 0
	}
	allSegsCount := addr.getDivisionCount()
	if segmentCount >= allSegsCount {
		return addr.addrIterator(nil)
	}
	useOriginal := !addr.section.isMultipleTo(segmentCount)
	address := addr.toAddress()
	var iterator Iterator[[]*AddressDivision]
	if !useOriginal {
		var hostSegIteratorProducer func(index int) Iterator[*AddressSegment]
		hostSegIteratorProducer = func(index int) Iterator[*AddressSegment] {
			return address.getSegment(index).identityIterator()
		}
		segIteratorProducer := func(index int) Iterator[*AddressSegment] {
			return address.getSegment(index).iterator()
		}
		iterator = segmentsIterator(
			allSegsCount,
			nil, //when no prefix we defer to other iterator, when there is one we use the whole original section in the encompassing iterator and not just the original segments
			segIteratorProducer,
			nil,
			segmentCount-1,
			segmentCount,
			hostSegIteratorProducer)
	}
	return addrIterator(
		useOriginal,
		address,
		address.getPrefixLen(),
		addr.section.isMultipleFrom(segmentCount),
		iterator)
}

// sequentialBlockIterator iterates through the minimal number of maximum-sized blocks comprising this subnet
// a block is sequential if given any two addresses in the block, any intervening address between the two is also in the block
func (addr *addressInternal) sequentialBlockIterator() Iterator[*Address] {
	return addr.blockIterator(addr.getSequentialBlockIndex())
}

func (addr *addressInternal) getSequentialBlockIndex() int {
	if addr.section == nil {
		return 0
	}
	return addr.section.GetSequentialBlockIndex()
}

func (addr *addressInternal) getSequentialBlockCount() *big.Int {
	if addr.section == nil {
		return bigOne()
	}
	return addr.section.GetSequentialBlockCount()
}

func (addr *addressInternal) hasZone() bool {
	return addr.zone != NoZone
}

func (addr *addressInternal) increment(increment int64) *Address {
	return addr.checkIdentity(addr.section.increment(increment))
}

func (addr *addressInternal) incrementBoundary(increment int64) *Address {
	return addr.checkIdentity(addr.section.incrementBoundary(increment))
}

func (addr *addressInternal) enumerate(other AddressType) *big.Int {
	if other == nil {
		return nil
	}
	otherAddr := other.ToAddressBase()
	if otherAddr == nil {
		return nil
	}
	otherSection := otherAddr.GetSection()
	return addr.section.Enumerate(otherSection)
}

func (addr *addressInternal) getStringCache() *stringCache {
	cache := addr.cache
	if cache == nil {
		return nil
	}
	return addr.cache.stringCache
}

func (addr *addressInternal) getSegmentStrings() []string {
	return addr.section.getSegmentStrings()
}

func (addr *addressInternal) toCanonicalString() string {
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.section.ToIPv6().toCanonicalString(addr.zone)
		}
		return cacheStr(&cache.canonicalString,
			func() string { return addr.section.ToIPv6().toCanonicalString(addr.zone) })
	}
	return addr.section.ToCanonicalString()
}

func (addr *addressInternal) toNormalizedString() string {
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.section.ToIPv6().toNormalizedString(addr.zone)
		}
		return cacheStr(&cache.normalizedIPv6String,
			func() string { return addr.section.ToIPv6().toNormalizedString(addr.zone) })
	}
	return addr.section.ToNormalizedString()
}

func (addr *addressInternal) toNormalizedWildcardString() string {
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.section.ToIPv6().toNormalizedWildcardStringZoned(addr.zone)
		}
		return cacheStr(&cache.normalizedIPv6String,
			func() string { return addr.section.ToIPv6().toNormalizedWildcardStringZoned(addr.zone) })
	}
	return addr.section.ToNormalizedWildcardString()
}

func (addr *addressInternal) toCompressedString() string {
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.section.ToIPv6().toCompressedString(addr.zone)
		}
		return cacheStr(&cache.compressedIPv6String,
			func() string { return addr.section.ToIPv6().toCompressedString(addr.zone) })
	}
	return addr.section.ToCompressedString()
}

func (addr *addressInternal) toOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.section.toOctalStringZoned(with0Prefix, addr.zone)
		}
		var cacheField **string
		if with0Prefix {
			cacheField = &cache.octalStringPrefixed
		} else {
			cacheField = &cache.octalString
		}
		return cacheStrErr(cacheField,
			func() (string, addrerr.IncompatibleAddressError) {
				return addr.section.toOctalStringZoned(with0Prefix, addr.zone)
			})
	}
	return addr.section.ToOctalString(with0Prefix)
}

func (addr *addressInternal) toBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.section.toBinaryStringZoned(with0bPrefix, addr.zone)
		}
		var cacheField **string
		if with0bPrefix {
			cacheField = &cache.binaryStringPrefixed
		} else {
			cacheField = &cache.binaryString
		}
		return cacheStrErr(cacheField,
			func() (string, addrerr.IncompatibleAddressError) {
				return addr.section.toBinaryStringZoned(with0bPrefix, addr.zone)
			})
	}
	return addr.section.ToBinaryString(with0bPrefix)
}

func (addr *addressInternal) toHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.section.toHexStringZoned(with0xPrefix, addr.zone)
		}
		var cacheField **string
		if with0xPrefix {
			cacheField = &cache.hexStringPrefixed
		} else {
			cacheField = &cache.hexString
		}
		return cacheStrErr(cacheField,
			func() (string, addrerr.IncompatibleAddressError) {
				return addr.section.toHexStringZoned(with0xPrefix, addr.zone)
			})
	}
	return addr.section.ToHexString(with0xPrefix)
}

func (addr *addressInternal) format(state fmt.State, verb rune) {
	section := addr.section
	section.format(state, verb, addr.zone, addr.isIP())
}

var zeroAddr = createAddress(zeroSection, NoZone)

// Address represents a single address, or a collection of multiple addresses, such as with an IP subnet or a set of MAC addresses.
//
// Addresses consist of a sequence of segments, each of equal bit-size.
// The number of such segments and the bit-size are determined by the underlying version or type of the address, whether IPv4, IPv6, MAC, or other.
// Each segment can represent a single value or a sequential range of values.  Addresses can also have an associated prefix length,
// which is the number of consecutive bits comprising the prefix, the most significant bits of an address.
//
// To construct one from a string, use
// NewIPAddressString or NewMACAddressString,
// then use the ToAddress or GetAddress methods to get an [IPAddress] or [MACAddress],
// and then you can convert to this type using the ToAddressBase method.
//
// Any given specific address types can be converted to Address with the ToAddressBase method,
// and then back again to their original types with methods like ToIPv6, ToIP, ToIPv4, and ToMAC.
// When calling such a method on a given address, if the address was not originally constructed as the type returned from the method,
// then the method will return nil.  Conversion methods work with nil pointers (returning nil) so that they can be chained together safely.
//
// This allows for polymorphic code that works with all addresses, such as with the address trie code in this library,
// while still allowing for methods and code specific to each address version or type.
//
// You can also use the methods IsIPv6, IsIP, IsIPv4, and IsMAC,
// which will return true if and only if the corresponding method ToIPv6, ToIP, ToIPv4, and ToMAC returns non-nil, respectively.
//
// The zero value for an Address is an address with no segments and no associated address version or type, also known as the adaptive zero.
type Address struct {
	addressInternal
}

func (addr *Address) init() *Address {
	if addr.section == nil {
		return zeroAddr // this has a zero section rather than a nil section
	}
	return addr
}

// GetCount returns the count of addresses that this address or subnet represents.
//
// If just a single address, not a collection nor subnet of multiple addresses, returns 1.
//
// For instance, the IP address subnet "2001:db8::/64" has the count of 2 to the power of 64.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (addr *Address) GetCount() *big.Int {
	if addr == nil {
		return bigZero()
	}
	return addr.getCount()
}

// IsMultiple returns true if this represents more than a single individual address, whether it is a collection or subnet of multiple addresses.
func (addr *Address) IsMultiple() bool {
	return addr != nil && addr.isMultiple()
}

// IsPrefixed returns whether this address has an associated prefix length.
func (addr *Address) IsPrefixed() bool {
	return addr != nil && addr.isPrefixed()
}

// PrefixEqual determines if the given address matches this address up to the prefix length of this address.
// It returns whether the two addresses share the same range of prefix values.
func (addr *Address) PrefixEqual(other AddressType) bool {
	return addr.init().prefixEquals(other)
}

// PrefixContains returns whether the prefix values in the given address or subnet
// are prefix values in this address or subnet, using the prefix length of this address or subnet.
// If this address has no prefix length, the entire address is compared.
//
// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
func (addr *Address) PrefixContains(other AddressType) bool {
	return addr.init().prefixContains(other)
}

// containsSame returns whether this address contains all addresses in the given address or subnet of the same type.
func (addr *Address) containsSame(other *Address) bool {
	return addr.Contains(other)
}

// Contains returns whether this is the same type and version as the given address or subnet and whether it contains all addresses in the given address or subnet.
func (addr *Address) Contains(other AddressType) bool {
	if addr == nil {
		return other == nil || other.ToAddressBase() == nil
	}
	return addr.init().contains(other)
}

// Overlaps returns true if this address overlaps the given address or subnet
func (addr *Address) Overlaps(other AddressType) bool {
	if addr == nil {
		return true
	}
	return addr.init().overlaps(other)
}

// Compare returns a negative integer, zero, or a positive integer if this address or subnet is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (addr *Address) Compare(item AddressItem) int {
	return CountComparator.Compare(addr, item)
}

// Equal returns whether the given address or subnet is equal to this address or subnet.
// Two address instances are equal if they represent the same set of addresses.
func (addr *Address) Equal(other AddressType) bool {
	if addr == nil {
		return other == nil || other.ToAddressBase() == nil
	} else if other.ToAddressBase() == nil {
		return false
	}
	return addr.init().equals(other)
}

// CompareSize compares the counts of two subnets or addresses or other address items, the number of individual items within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one subnet or collection represents more individual items than another.
//
// CompareSize returns a positive integer if this address or subnet has a larger count than the item given, zero if they are the same, or a negative integer if the other has a larger count.
func (addr *Address) CompareSize(other AddressItem) int {
	if addr == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return addr.init().compareSize(other)
}

// TrieCompare compares two addresses according to address trie ordering.
// It returns a number less than zero, zero, or a number greater than zero if the first address argument is less than, equal to, or greater than the second.
//
// The comparison is intended for individual addresses and CIDR prefix blocks.
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *Address) TrieCompare(other *Address) (int, addrerr.IncompatibleAddressError) {
	if thisAddr := addr.ToIPv4(); thisAddr != nil {
		if oth := other.ToIPv4(); oth != nil {
			return thisAddr.TrieCompare(oth), nil
		}
	} else if thisAddr := addr.ToIPv6(); thisAddr != nil {
		if oth := other.ToIPv6(); oth != nil {
			return thisAddr.TrieCompare(oth), nil
		}
	} else if thisAddr := addr.ToMAC(); thisAddr != nil {
		if oth := other.ToMAC(); oth != nil {
			return thisAddr.TrieCompare(oth)
		}
	}
	if segmentCount, otherSegmentCount := addr.getDivisionCount(), other.getDivisionCount(); segmentCount == otherSegmentCount {
		if bitsPerSegment, otherBitsPerSegment := addr.GetBitsPerSegment(), other.GetBitsPerSegment(); bitsPerSegment == otherBitsPerSegment {
			return addr.trieCompare(other), nil
		}
	}
	return 0, &incompatibleAddressError{addressError{key: "ipaddress.error.mismatched.bit.size"}}
}

// TrieIncrement returns the next address or block according to address trie ordering.
//
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *Address) TrieIncrement() *Address {
	if res, ok := trieIncrement(addr); ok {
		return res
	}
	return nil
}

// TrieDecrement returns the previous or block address according to address trie ordering.
//
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *Address) TrieDecrement() *Address {
	if res, ok := trieDecrement(addr); ok {
		return res
	}
	return nil
}

// GetSection returns the backing section for this address or subnet, comprising all segments.
func (addr *Address) GetSection() *AddressSection {
	return addr.init().section
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (addr *Address) GetTrailingSection(index int) *AddressSection {
	return addr.GetSection().GetTrailingSection(index)
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (addr *Address) GetSubSection(index, endIndex int) *AddressSection {
	return addr.GetSection().GetSubSection(index, endIndex)
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (addr *Address) CopySubSegments(start, end int, segs []*AddressSegment) (count int) {
	return addr.GetSection().CopySubSegments(start, end, segs)
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (addr *Address) CopySegments(segs []*AddressSegment) (count int) {
	return addr.GetSection().CopySegments(segs)
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
func (addr *Address) GetSegments() []*AddressSegment {
	return addr.GetSection().GetSegments()
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *Address) GetSegment(index int) *AddressSegment {
	return addr.getSegment(index)
}

// GetSegmentCount returns the segment count, the number of segments in this address.
// For example, IPv4 addresses have 4, IPv6 addresses have 8.
func (addr *Address) GetSegmentCount() int {
	return addr.getDivisionCount()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (addr *Address) ForEachSegment(consumer func(segmentIndex int, segment *AddressSegment) (stop bool)) int {
	return addr.GetSection().ForEachSegment(consumer)
}

// GetGenericDivision returns the segment at the given index as a DivisionType.
// The first segment is at index 0.
// GetGenericDivision will panic given a negative index or index larger than the division count.
func (addr *Address) GetGenericDivision(index int) DivisionType {
	return addr.getDivision(index)
}

// GetGenericSegment returns the segment at the given index as an AddressSegmentType.
// The first segment is at index 0.
// GetGenericSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *Address) GetGenericSegment(index int) AddressSegmentType {
	return addr.getSegment(index)
}

// GetDivisionCount returns the division count, which is the same as the segment count, since the divisions of an address are the segments.
func (addr *Address) GetDivisionCount() int {
	return addr.getDivisionCount()
}

// TestBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this address.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (addr *Address) TestBit(n BitCount) bool {
	return addr.init().testBit(n)
}

// IsOneBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the most significant bit.
// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this item.
func (addr *Address) IsOneBit(bitIndex BitCount) bool {
	return addr.init().isOneBit(bitIndex)
}

// GetLower returns the address in the subnet or address collection with the lowest numeric value,
// which will be the receiver if it represents a single address.
// For example, for "1.2-3.4.5-6", the series "1.2.4.5" is returned.
func (addr *Address) GetLower() *Address {
	return addr.init().getLower()
}

// GetUpper returns the address in the subnet or address collection with the highest numeric value,
// which will be the receiver if it represents a single address.
// For example, for the subnet "1.2-3.4.5-6", the address "1.3.4.6" is returned.
func (addr *Address) GetUpper() *Address {
	return addr.init().getUpper()
}

// GetValue returns the lowest address in this subnet or address collection as an integer value.
func (addr *Address) GetValue() *big.Int {
	return addr.init().section.GetValue()
}

// GetUpperValue returns the highest address in this subnet or address collection as an integer value.
func (addr *Address) GetUpperValue() *big.Int {
	return addr.init().section.GetUpperValue()
}

// Bytes returns the lowest address in this subnet or address collection as a byte slice.
func (addr *Address) Bytes() []byte {
	return addr.init().section.Bytes()
}

// UpperBytes returns the highest address in this subnet or address collection as a byte slice.
func (addr *Address) UpperBytes() []byte {
	return addr.init().section.UpperBytes()
}

// CopyBytes copies the value of the lowest individual address in the subnet into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *Address) CopyBytes(bytes []byte) []byte {
	return addr.init().section.CopyBytes(bytes)
}

// CopyUpperBytes copies the value of the highest individual address in the subnet into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *Address) CopyUpperBytes(bytes []byte) []byte {
	return addr.init().section.CopyUpperBytes(bytes)
}

// IsMax returns whether this address matches exactly the maximum possible value, the address whose bits are all ones.
func (addr *Address) IsMax() bool {
	return addr.init().section.IsMax()
}

// IncludesMax returns whether this address includes the max address, the address whose bits are all ones, within its range.
func (addr *Address) IncludesMax() bool {
	return addr.init().section.IncludesMax()
}

// ToPrefixBlock returns the address collection associated with the prefix of this address or address collection,
// the address whose prefix matches the prefix of this address, and the remaining bits span all values.
// If this address has no prefix length, this address is returned.
//
// The returned address collection will include all addresses with the same prefix as this one, the prefix "block".
func (addr *Address) ToPrefixBlock() *Address {
	return addr.init().toPrefixBlock()
}

// ToPrefixBlockLen returns the address associated with the prefix length provided,
// the address collection whose prefix of that length matches the prefix of this address, and the remaining bits span all values.
//
// The returned address will include all addresses with the same prefix as this one, the prefix "block".
func (addr *Address) ToPrefixBlockLen(prefLen BitCount) *Address {
	return addr.init().toPrefixBlockLen(prefLen)
}

// ToBlock creates a new block of addresses by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (addr *Address) ToBlock(segmentIndex int, lower, upper SegInt) *Address {
	return addr.init().toBlock(segmentIndex, lower, upper)
}

// WithoutPrefixLen provides the same address but with no prefix length.  The values remain unchanged.
func (addr *Address) WithoutPrefixLen() *Address {
	if !addr.IsPrefixed() {
		return addr
	}
	return addr.init().withoutPrefixLen()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (addr *Address) SetPrefixLen(prefixLen BitCount) *Address {
	return addr.init().setPrefixLen(prefixLen)
}

// SetPrefixLenZeroed sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address.
// The provided prefix length will be adjusted to these boundaries if necessary.
//
// If this address has a prefix length, and the prefix length is increased when setting the new prefix length, the bits moved within the prefix become zero.
// If this address has a prefix length, and the prefix length is decreased when setting the new prefix length, the bits moved outside the prefix become zero.
//
// In other words, bits that move from one side of the prefix length to the other (bits moved into the prefix or outside the prefix) are zeroed.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (addr *Address) SetPrefixLenZeroed(prefixLen BitCount) (*Address, addrerr.IncompatibleAddressError) {
	return addr.init().setPrefixLenZeroed(prefixLen)
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address.
//
// If this address has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (addr *Address) AdjustPrefixLen(prefixLen BitCount) *Address {
	return addr.adjustPrefixLen(prefixLen).ToAddressBase()
}

// AdjustPrefixLenZeroed increases or decreases the prefix length by the given increment while zeroing out the bits that have moved into or outside the prefix.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address.
//
// If this address has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
//
// When prefix length is increased, the bits moved within the prefix become zero.
// When a prefix length is decreased, the bits moved outside the prefix become zero.
//
// For example, "1.2.0.0/16" adjusted by -8 becomes "1.0.0.0/8".
// "1.2.0.0/16" adjusted by 8 becomes "1.2.0.0/24".
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (addr *Address) AdjustPrefixLenZeroed(prefixLen BitCount) (*Address, addrerr.IncompatibleAddressError) {
	res, err := addr.adjustPrefixLenZeroed(prefixLen)
	return res.ToAddressBase(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address - it is required that the range of values match the range of a prefix block.
// If there is no such address, then nil is returned.
//
// Examples:
//   - 1.2.3.4 returns 1.2.3.4/32
//   - 1.2.*.* returns 1.2.0.0/16
//   - 1.2.*.0/24 returns 1.2.0.0/16
//   - 1.2.*.4 returns nil
//   - 1.2.0-1.* returns 1.2.0.0/23
//   - 1.2.1-2.* returns nil
//   - 1.2.252-255.* returns 1.2.252.0/22
//   - 1.2.3.4/16 returns 1.2.3.4/32
func (addr *Address) AssignPrefixForSingleBlock() *Address {
	return addr.init().assignPrefixForSingleBlock()
}

// AssignMinPrefixForBlock returns an equivalent subnet, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this subnet.
//
// In other words, this method assigns a prefix length to this subnet matching the largest prefix block in this subnet.
//
// Examples:
//   - 1.2.3.4 returns 1.2.3.4/32
//   - 1.2.*.* returns 1.2.0.0/16
//   - 1.2.*.0/24 returns 1.2.0.0/16
//   - 1.2.*.4 returns 1.2.*.4/32
//   - 1.2.0-1.* returns 1.2.0.0/23
//   - 1.2.1-2.* returns 1.2.1-2.0/24
//   - 1.2.252-255.* returns 1.2.252.0/22
//   - 1.2.3.4/16 returns 1.2.3.4/32
func (addr *Address) AssignMinPrefixForBlock() *Address {
	return addr.init().assignMinPrefixForBlock()
}

// ToSinglePrefixBlockOrAddress converts to a single prefix block or address.
// If the given address is a single prefix block, it is returned.
// If it can be converted to a single prefix block by assigning a prefix length, the converted block is returned.
// If it is a single address, any prefix length is removed and the address is returned.
// Otherwise, nil is returned.
// This method provides the address formats used by tries.
// ToSinglePrefixBlockOrAddress is quite similar to AssignPrefixForSingleBlock, which always returns prefixed addresses, while this does not.
func (addr *Address) ToSinglePrefixBlockOrAddress() *Address {
	return addr.init().toSinglePrefixBlockOrAddr()
}

func (addr *Address) toSinglePrefixBlockOrAddress() (*Address, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.address.not.block", str: addr.String()}}
	}
	res := addr.ToSinglePrefixBlockOrAddress()
	if res == nil {
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.address.not.block", str: addr.String()}}
	}
	return res, nil
}

// GetMaxSegmentValue returns the maximum possible segment value for this type of address.
//
// Note this is not the maximum of the range of segment values in this specific address,
// this is the maximum value of any segment for this address type and version, determined by the number of bits per segment.
func (addr *Address) GetMaxSegmentValue() SegInt {
	return addr.init().getMaxSegmentValue()
}

// Iterator provides an iterator to iterate through the individual addresses of this address or subnet.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual addresses.
//
// Call IsMultiple to determine if this instance represents multiple addresses, or GetCount for the count.
func (addr *Address) Iterator() Iterator[*Address] {
	if addr == nil {
		return nilAddrIterator()
	}
	return addr.addrIterator(nil)
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this subnet,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this subnet.
//
// If the subnet has no prefix length, then this is equivalent to Iterator.
func (addr *Address) PrefixIterator() Iterator[*Address] {
	return addr.prefixIterator(false)
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address or subnet.
// Each iterated address or subnet will be a prefix block with the same prefix length as this address or subnet.
//
// If this address has no prefix length, then this is equivalent to Iterator.
func (addr *Address) PrefixBlockIterator() Iterator[*Address] {
	return addr.prefixIterator(true)
}

// BlockIterator iterates through the addresses that can be obtained by iterating through all the upper segments up to the given segment count.
// The segments following remain the same in all iterated addresses.
//
// For instance, given the IPv4 subnet "1-2.3-4.5-6.7" and the count argument 2,
// BlockIterator will iterate through "1.3.5-6.7", "1.4.5-6.7", "2.3.5-6.7" and "2.4.5-6.7".
func (addr *Address) BlockIterator(segmentCount int) Iterator[*Address] {
	return addr.init().blockIterator(segmentCount)
}

// SequentialBlockIterator iterates through the sequential subnets or addresses that make up this address or subnet.
//
// Practically, this means finding the count of segments for which the segments that follow are not full range, and then using BlockIterator with that segment count.
//
// For instance, given the IPv4 subnet "1-2.3-4.5-6.7-8", it will iterate through "1.3.5.7-8", "1.3.6.7-8", "1.4.5.7-8", "1.4.6.7-8", "2.3.5.7-8", "2.3.6.7-8", "2.4.6.7-8" and "2.4.6.7-8".
//
// Use GetSequentialBlockCount to get the number of iterated elements.
func (addr *Address) SequentialBlockIterator() Iterator[*Address] {
	return addr.init().sequentialBlockIterator()
}

// GetSequentialBlockIndex gets the minimal segment index for which all following segments are full-range blocks.
//
// The segment at this index is not a full-range block itself, unless all segments are full-range.
// The segment at this index and all following segments form a sequential range.
// For the full subnet to be sequential, the preceding segments must be single-valued.
func (addr *Address) GetSequentialBlockIndex() int {
	return addr.getSequentialBlockIndex()
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential subnets that comprise this subnet.
func (addr *Address) GetSequentialBlockCount() *big.Int {
	return addr.getSequentialBlockCount()
}

// IncrementBoundary returns the address that is the given increment from the range boundaries of this subnet or address collection.
//
// If the given increment is positive, adds the value to the upper address (GetUpper) in the range to produce a new address.
// If the given increment is negative, adds the value to the lower address (GetLower) in the range to produce a new address.
// If the increment is zero, returns this address.
//
// If this is a single address value, that address is simply incremented by the given increment value, positive or negative.
//
// On address overflow or underflow, IncrementBoundary returns nil.
func (addr *Address) IncrementBoundary(increment int64) *Address {
	return addr.init().IncrementBoundary(increment)
}

// Increment returns the address from the subnet that is the given increment upwards into the subnet range,
// with the increment of 0 returning the first address in the range.
//
// If the increment i matches or exceeds the subnet size count c, then i - c + 1
// is added to the upper address of the range.
// An increment matching the subnet count gives you the address just above the highest address in the subnet.
//
// If the increment is negative, it is added to the lower address of the range.
// To get the address just below the lowest address of the subnet, use the increment -1.
//
// If this is just a single address value, the address is simply incremented by the given increment, positive or negative.
//
// If this is a subnet with multiple values, a positive increment i is equivalent i + 1 values from the subnet iterator and beyond.
// For instance, a increment of 0 is the first value from the iterator, an increment of 1 is the second value from the iterator, and so on.
// An increment of a negative value added to the subnet count is equivalent to the same number of iterator values preceding the upper bound of the iterator.
// For instance, an increment of count - 1 is the last value from the iterator, an increment of count - 2 is the second last value, and so on.
//
// On address overflow or underflow, Increment returns nil.
func (addr *Address) Increment(increment int64) *Address {
	return addr.init().increment(increment)
}

// Enumerate indicates where an address sits relative to the subnet ordering.
//
// Determines how many address elements of the subnet precede the given address element, if the address is in the subnet.
// If above the subnet range, it is the distance to the upper boundary added to the subnet count less one, and if below the subnet range, the distance to the lower boundary.
//
// In other words, if the given address is not in the subnet but above it, returns the number of addresses preceding the address from the upper range boundary,
// added to one less than the total number of subnet addresses.  If the given address is not in the subnet but below it, returns the number of addresses following the address to the lower subnet boundary.
//
// Returns nil when the argument is multi-valued. The argument must be an individual address.
//
// When this is also an individual address, the returned value is the distance (difference) between the two addresses.
//
// Enumerate is the inverse of the increment method:
//   - subnet.Enumerate(subnet.Increment(inc)) = inc
//   - subnet.Increment(subnet.Enumerate(newAddr)) = newAddr
//
// If the given address does not have the same version or type, then nil is returned.
func (addr *Address) Enumerate(other AddressType) *big.Int {
	return addr.init().enumerate(other)
}

// ReverseBytes returns a new address with the bytes reversed.  Any prefix length is dropped.
//
// If each segment is more than 1 byte long, and the bytes within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, then this returns an error.
//
// In practice this means that to be reversible, a segment range must include all values except possibly the largest and/or smallest, which reverse to themselves.
func (addr *Address) ReverseBytes() (*Address, addrerr.IncompatibleAddressError) {
	return addr.init().reverseBytes()
}

// ReverseBits returns a new address with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a segment range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (addr *Address) ReverseBits(perByte bool) (*Address, addrerr.IncompatibleAddressError) {
	return addr.init().reverseBits(perByte)
}

// ReverseSegments returns a new address with the segments reversed.
func (addr *Address) ReverseSegments() *Address {
	return addr.init().reverseSegments()
}

// IsMulticast returns whether this address is multicast.
func (addr *Address) IsMulticast() bool {
	if thisAddr := addr.ToIPv4(); thisAddr != nil {
		return thisAddr.IsMulticast()
	} else if thisAddr := addr.ToIPv6(); thisAddr != nil {
		return thisAddr.IsMulticast()
	} else if thisAddr := addr.ToMAC(); thisAddr != nil {
		return thisAddr.IsMulticast()
	}
	return false
}

// IsLocal returns whether the address can be considered a local address (as opposed to a global one).
func (addr *Address) IsLocal() bool {
	if thisAddr := addr.ToIPv4(); thisAddr != nil {
		return thisAddr.IsLocal()
	} else if thisAddr := addr.ToIPv6(); thisAddr != nil {
		return thisAddr.IsLocal()
	} else if thisAddr := addr.ToMAC(); thisAddr != nil {
		return thisAddr.IsLocal()
	}
	return false
}

// GetLeadingBitCount returns the number of consecutive leading one or zero bits.
// If ones is true, returns the number of consecutive leading one bits.
// Otherwise, returns the number of consecutive leading zero bits.
//
// This method applies to the lower address of the range if this is a subnet representing multiple values.
func (addr *Address) GetLeadingBitCount(ones bool) BitCount {
	return addr.init().getLeadingBitCount(ones)
}

// GetTrailingBitCount returns the number of consecutive trailing one or zero bits.
// If ones is true, returns the number of consecutive trailing zero bits.
// Otherwise, returns the number of consecutive trailing one bits.
//
// This method applies to the lower value of the range if this is a subnet representing multiple values.
func (addr *Address) GetTrailingBitCount(ones bool) BitCount {
	return addr.init().getTrailingBitCount(ones)
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
func (addr Address) Format(state fmt.State, verb rune) {
	addr.init().format(state, verb)
}

// String implements the [fmt.Stringer] interface, returning the canonical string provided by ToCanonicalString, or "<nil>" if the receiver is a nil pointer.
func (addr *Address) String() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toString()
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (addr *Address) GetSegmentStrings() []string {
	if addr == nil {
		return nil
	}
	return addr.init().getSegmentStrings()
}

// ToCanonicalString produces a canonical string for the address.
//
// For IPv4, dotted octet format, also known as dotted decimal format, is used.
// https://datatracker.ietf.org/doc/html/draft-main-ipaddr-text-rep-00#section-2.1
//
// For IPv6, RFC 5952 describes canonical string representation.
// https://en.wikipedia.org/wiki/IPv6_address#Representation
// http://tools.ietf.org/html/rfc5952
//
// For MAC, it uses the canonical standardized IEEE 802 MAC address representation of xx-xx-xx-xx-xx-xx.  An example is "01-23-45-67-89-ab".
// For range segments, '|' is used: "11-22-33|44-55-66".
//
// Each address has a unique canonical string, not counting the prefix length.
// With IP addresses, the prefix length is included in the string, and the prefix length can cause two equal addresses to have different strings, for example "1.2.3.4/16" and "1.2.3.4".
// It can also cause two different addresses to have the same string, such as "1.2.0.0/16" for the individual address "1.2.0.0" and also the prefix block "1.2.*.*".
// Use the IPAddress method ToCanonicalWildcardString for a unique string for each IP address and subnet.
func (addr *Address) ToCanonicalString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCanonicalString()
}

// ToNormalizedString produces a normalized string for the address.
//
// For IPv4, it is the same as the canonical string.
//
// For IPv6, it differs from the canonical string.  Zero-segments are not compressed.
//
// For MAC, it differs from the canonical string.  It uses the most common representation of MAC addresses: "xx:xx:xx:xx:xx:xx".  An example is "01:23:45:67:89:ab".
// For range segments, '-' is used: "11:22:33-44:55:66".
//
// Each address has a unique normalized string, not counting the prefix length.
// With IP addresses, the prefix length can cause two equal addresses to have different strings, for example "1.2.3.4/16" and "1.2.3.4".
// It can also cause two different addresses to have the same string, such as "1.2.0.0/16" for the individual address "1.2.0.0" and also the prefix block "1.2.*.*".
// Use the IPAddress method ToNormalizedWildcardString for a unique string for each IP address and subnet.
func (addr *Address) ToNormalizedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toNormalizedString()
}

// ToNormalizedWildcardString produces a string similar to the normalized string but avoids the CIDR prefix length in IP addresses.
// Multi-valued segments will be shown with wildcards and ranges (denoted by '*' and '-').
func (addr *Address) ToNormalizedWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toNormalizedWildcardString()
}

// ToCompressedString produces a short representation of this address while remaining within the confines of standard representation(s) of the address.
//
// For IPv4, it is the same as the canonical string.
//
// For IPv6, it differs from the canonical string.  It compresses the maximum number of zeros and/or host segments with the IPv6 compression notation '::'.
//
// For MAC, it differs from the canonical string.  It produces a shorter string for the address that has no leading zeros.
func (addr *Address) ToCompressedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCompressedString()
}

// ToHexString writes this address as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If an address collection cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *Address) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toHexString(with0xPrefix)
}

// ToOctalString writes this address as a single octal value (possibly two values if a range),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If an address collection cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *Address) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toOctalString(with0Prefix)
}

// ToBinaryString writes this address as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *Address) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toBinaryString(with0bPrefix)
}

// ToCustomString creates a customized string from this address or subnet according to the given string option parameters.
func (addr *Address) ToCustomString(stringOptions addrstr.StringOptions) string {
	if addr == nil {
		return nilString()
	}
	return addr.GetSection().toCustomStringZoned(stringOptions, addr.zone)
}

// ToAddressString retrieves or generates a HostIdentifierString instance for this Address object.
//
// This same Address instance can be retrieved from the resulting HostIdentifierString object using the GetAddress method.
//
// In general, users create Address instances from IPAddressString or MACAddressString instances,
// while the reverse direction is generally not common and not useful.
//
// However, the reverse direction can be useful under certain circumstances, such as when maintaining a collection of HostIdentifierString instances.
func (addr *Address) ToAddressString() HostIdentifierString {
	if addr.isIP() {
		return addr.ToIP().ToAddressString()
	} else if addr.isMAC() {
		return addr.ToMAC().ToAddressString()
	}
	return nil
}

// IsIPv4 returns true if this address or subnet originated as an IPv4 address or subnet.  If so, use ToIPv4 to convert back to the IPv4-specific type.
func (addr *Address) IsIPv4() bool {
	return addr != nil && addr.isIPv4()
}

// IsIPv6 returns true if this address or subnet originated as an IPv6 address or subnet.  If so, use ToIPv6 to convert back to the IPv6-specific type.
func (addr *Address) IsIPv6() bool {
	return addr != nil && addr.isIPv6()
}

// IsIP returns true if this address or subnet originated as an IPv4 or IPv6 address or subnet, or an implicitly zero-valued IP.  If so, use ToIP to convert back to the IP-specific type.
func (addr *Address) IsIP() bool {
	return addr != nil && addr.isIP()
}

// IsMAC returns true if this address or address collection originated as a MAC address or address collection.  If so, use ToMAC to convert back to the MAC-specific type.
func (addr *Address) IsMAC() bool {
	return addr != nil && addr.isMAC()
}

// ToAddressBase is an identity method.
//
// ToAddressBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *Address) ToAddressBase() *Address {
	return addr
}

// ToIP converts to an IPAddress if this address or subnet originated as an IPv4 or IPv6 address or subnet, or an implicitly zero-valued IP.
// If not, ToIP returns nil.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *Address) ToIP() *IPAddress {
	if addr.IsIP() {
		return (*IPAddress)(unsafe.Pointer(addr))
	}
	return nil
}

// ToIPv6 converts to an IPv6Address if this address or subnet originated as an IPv6 address or subnet.
// If not, ToIPv6 returns nil.
//
// ToIPv6 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *Address) ToIPv6() *IPv6Address {
	if addr.IsIPv6() {
		return (*IPv6Address)(unsafe.Pointer(addr))
	}
	return nil
}

// ToIPv4 converts to an IPv4Address if this address or subnet originated as an IPv4 address or subnet.
// If not, ToIPv4 returns nil.
//
// ToIPv4 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *Address) ToIPv4() *IPv4Address {
	if addr.IsIPv4() {
		return (*IPv4Address)(unsafe.Pointer(addr))
	}
	return nil
}

// ToMAC converts to a MACAddress if this address or subnet originated as a MAC address or subnet.
// If not, ToMAC returns nil.
//
// ToMAC can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *Address) ToMAC() *MACAddress {
	if addr.IsMAC() {
		return (*MACAddress)(addr)
	}
	return nil
}

// Wrap wraps this address, returning a WrappedAddress, an implementation of ExtendedSegmentSeries,
// which can be used to write code that works with both addresses and address sections.
func (addr *Address) Wrap() WrappedAddress {
	return wrapAddress(addr.init())
}

// ToKey creates the associated address key.
// While addresses can be compared with the Compare, TrieCompare or Equal methods as well as various provided instances of AddressComparator,
// they are not comparable with Go operators.
// However, AddressKey instances are comparable with Go operators, and thus can be used as map keys.
func (addr *Address) ToKey() Key[*Address] {
	key := Key[*Address]{}
	contents := &key.keyContents
	if thisAddr := addr.ToIPv4(); thisAddr != nil {
		key.scheme = ipv4Scheme
		thisAddr.toIPv4Key(contents)
	} else if thisAddr := addr.ToIPv6(); thisAddr != nil {
		key.scheme = ipv6Scheme
		thisAddr.toIPv6Key(contents)
	} else if thisAddr := addr.ToMAC(); thisAddr != nil {
		if addr.GetSegmentCount() == ExtendedUniqueIdentifier64SegmentCount {
			key.scheme = eui64Scheme
		} else {
			key.scheme = mac48Scheme
		}
		thisAddr.toMACKey(contents)
	} // else key.scheme == adaptiveZeroScheme
	return key
}

// ToGenericKey produces a generic Key[*Address] that can be used with generic code working with [Address], [IPAddress], [IPv4Address], [IPv6Address] and [MACAddress].
func (addr *Address) ToGenericKey() Key[*Address] {
	return addr.ToKey()
}

func (addr *Address) fromKey(scheme addressScheme, key *keyContents) *Address {
	if scheme == ipv4Scheme {
		ipv4Addr := fromIPv4IPKey(key)
		return ipv4Addr.ToAddressBase()
	} else if scheme == ipv6Scheme {
		ipv6Addr := fromIPv6IPKey(key)
		return ipv6Addr.ToAddressBase()
	} else if scheme == eui64Scheme || scheme == mac48Scheme {
		macAddr := fromMACAddrKey(scheme, key)
		return macAddr.ToAddressBase()
	}
	// scheme == adaptiveZeroScheme
	zeroAddr := Address{}
	return zeroAddr.init()
}

// AddrsMatchUnordered checks if the two slices share the same list of addresses, subnets, or address collections, in any order, using address equality.
// The function can handle duplicates and nil addresses.
func AddrsMatchUnordered[T, U AddressType](addrs1 []T, addrs2 []U) (result bool) {
	len1 := len(addrs1)
	len2 := len(addrs2)
	sameLen := len1 == len2
	if len1 == 0 || len2 == 0 {
		result = sameLen
	} else if len1 == 1 && sameLen {
		result = addrs1[0].Equal(addrs2[0])
	} else if len1 == 2 && sameLen {
		if addrs1[0].Equal(addrs2[0]) {
			result = addrs1[1].Equal(addrs2[1])
		} else if result = addrs1[0].Equal(addrs2[1]); result {
			result = addrs1[1].Equal(addrs2[0])
		}
	} else {
		result = reflect.DeepEqual(asMap(addrs1), asMap(addrs2))
	}
	return
}

// AddrsMatchOrdered checks if the two slices share the same ordered list of addresses, subnets, or address collections, using address equality.
// Duplicates and nil addresses are allowed.
func AddrsMatchOrdered[T, U AddressType](addrs1 []T, addrs2 []U) (result bool) {
	len1 := len(addrs1)
	len2 := len(addrs2)
	if len1 != len2 {
		return
	}
	for i, addr := range addrs1 {
		if !addr.Equal(addrs2[i]) {
			return
		}
	}
	return true
}

func asMap[T AddressType](addrs []T) (result map[string]struct{}) {
	if addrLen := len(addrs); addrLen > 0 {
		result = make(map[string]struct{})
		for _, addr := range addrs {
			result[addr.ToAddressBase().ToNormalizedWildcardString()] = struct{}{}
		}
	}
	return
}
