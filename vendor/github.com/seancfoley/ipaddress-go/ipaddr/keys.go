//
// Copyright 2024 Sean C Foley
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

package ipaddr

import "fmt"

func newSequentialRangeKey[T SequentialRangeConstraint[T]](rng *SequentialRange[T]) (key SequentialRangeKey[T]) {
	lower := rng.GetLower()
	upper := rng.GetUpper()
	lowerIp := lower.ToIP()
	upperIp := upper.ToIP()

	var t T
	anyt := any(t)
	_, isIP := anyt.(*IPAddress)
	if lowerIp.isIPv4() {
		section := lowerIp.GetSection()
		divs := section.getDivArray()
		for _, div := range divs {
			seg := div.ToIPv4()
			val := &key.vals[0]
			newLower := (val.lower << IPv4BitsPerSegment) | uint64(seg.GetIPv4SegmentValue())
			val.lower = newLower
		}
		section = upperIp.GetSection()
		divs = section.getDivArray()
		for _, div := range divs {
			seg := div.ToIPv4()
			val := &key.vals[0]
			newUpper := (val.upper << IPv4BitsPerSegment) | uint64(seg.GetIPv4SegmentValue())
			val.upper = newUpper
		}
		if isIP {
			key.addrType = ipv4Type
		}
	} else if lowerIp.isIPv6() {
		section := lowerIp.GetSection()
		divs := section.getDivArray()
		for i, div := range divs {
			seg := div.ToIPv6()
			val := &key.vals[i>>2]
			newLower := (val.lower << IPv6BitsPerSegment) | uint64(seg.GetIPv6SegmentValue())
			val.lower = newLower
		}
		section = upperIp.GetSection()
		divs = section.getDivArray()
		for i, div := range divs {
			seg := div.ToIPv6()
			val := &key.vals[i>>2]
			newUpper := (val.upper << IPv6BitsPerSegment) | uint64(seg.GetIPv6SegmentValue())
			val.upper = newUpper
		}
		if isIP {
			key.addrType = ipv6Type
		}
	} else { // nether IPv4 nor IPv6, the zero IP address
		// key.addrType is zeroType
	}
	return
}

// SequentialRangeKey is a representation of SequentialRange that is comparable as defined by the language specification.
// See https://go.dev/ref/spec#Comparison_operators
//
// It can be used as a map key.
// The zero value is a range from a zero-length address to itself.
type SequentialRangeKey[T SequentialRangeConstraint[T]] struct {
	vals [2]struct {
		lower,
		upper uint64
	}
	addrType addrType // only used when T is *IPAddress to indicate version for non-zero-valued address
}

// ToSeqRange converts back to a sequential range instance.
func (key SequentialRangeKey[T]) ToSeqRange() *SequentialRange[T] {
	var lower, upper T
	var isMult bool
	isIP, isIPv4, isIPv6 := false, false, false
	anyt := any(lower)
	if _, isIPv4 = anyt.(*IPv4Address); !isIPv4 {
		if _, isIPv6 = anyt.(*IPv6Address); !isIPv6 {
			if _, isIP = anyt.(*IPAddress); isIP {
				addressType := key.addrType
				if isIPv4 = addressType.isIPv4(); !isIPv4 {
					if isIPv6 = addressType.isIPv6(); !isIPv6 {
						if isNeither := addressType.isZeroSegments(); isNeither {
							lower = any(zeroIPAddr).(T)
							upper = lower
						} else {
							panic("supports only IP addresses")
						}
					}
				}
			} else {
				panic("supports only IP addresses")
			}
		}
	}
	if isIPv6 {
		lower6 := NewIPv6AddressFromVals(
			func(segmentIndex int) IPv6SegInt {
				valsIndex := segmentIndex >> 2
				segIndex := ((IPv6SegmentCount - 1) - segmentIndex) & 0x3
				return IPv6SegInt(key.vals[valsIndex].lower >> (segIndex << ipv6BitsToSegmentBitshift))
			})
		upper6 := NewIPv6AddressFromVals(
			func(segmentIndex int) IPv6SegInt {
				valsIndex := segmentIndex >> 2
				segIndex := ((IPv6SegmentCount - 1) - segmentIndex) & 0x3
				return IPv6SegInt(key.vals[valsIndex].upper >> (segIndex << ipv6BitsToSegmentBitshift))
			})
		isMult = key.vals[1].lower != key.vals[1].upper || key.vals[0].lower != key.vals[0].upper
		if isIP {
			lower = any(lower6.ToIP()).(T)
			upper = any(upper6.ToIP()).(T)
		} else {
			lower = any(lower6).(T)
			upper = any(upper6).(T)
		}
	} else if isIPv4 {
		l := uint32(key.vals[0].lower)
		u := uint32(key.vals[0].upper)
		lower4 := NewIPv4AddressFromUint32(l)
		upper4 := NewIPv4AddressFromUint32(u)
		isMult = l != u
		if isIP {
			lower = any(lower4.ToIP()).(T)
			upper = any(upper4.ToIP()).(T)
		} else {
			lower = any(lower4).(T)
			upper = any(upper4).(T)
		}
	}
	return newSequRangeUnchecked(lower, upper, isMult)
}

// String calls the String method in the corresponding sequential range.
func (key SequentialRangeKey[T]) String() string {
	return key.ToSeqRange().String()
}

// IPv4AddressKey is a representation of an IPv4 address that is comparable as defined by the language specification.
// See https://go.dev/ref/spec#Comparison_operators
//
// It can be used as a map key.  It can be obtained from its originating address instances.
// The zero value corresponds to the zero-value for IPv4Address.
// Keys do not incorporate prefix length to ensure that all equal addresses have equal keys.
// To create a key that has prefix length, combine into a struct with the PrefixKey obtained by passing the address into PrefixKeyFrom.
// IPv4Address can be compared using the Compare or Equal methods, or using an AddressComparator.
type IPv4AddressKey struct {
	vals uint64 // upper and lower combined into one uint64
}

// ToAddress converts back to an address instance.
func (key IPv4AddressKey) ToAddress() *IPv4Address {
	return fromIPv4Key(key)
}

// String calls the String method in the corresponding address.
func (key IPv4AddressKey) String() string {
	return key.ToAddress().String()
}

type testComparableConstraint[T comparable] struct{}

var (
	// ensure our 5 key types are indeed comparable
	_ testComparableConstraint[IPv4AddressKey]
	_ testComparableConstraint[IPv6AddressKey]
	_ testComparableConstraint[MACAddressKey]
	_ testComparableConstraint[Key[*IPAddress]]
	_ testComparableConstraint[Key[*Address]]
	//_ testComparableConstraint[RangeBoundaryKey[*IPv4Address]] // does not compile, as expected, because it has an interface field.  But it is still go-comparable.
)

// IPv6AddressKey is a representation of an IPv6 address that is comparable as defined by the language specification.
// See https://go.dev/ref/spec#Comparison_operators
//
// It can be used as a map key.  It can be obtained from its originating address instances.
// The zero value corresponds to the zero-value for IPv6Address.
// Keys do not incorporate prefix length to ensure that all equal addresses have equal keys.
// To create a key that has prefix length, combine into a struct with the PrefixKey obtained by passing the address into PrefixKeyFrom.
// IPv6Address can be compared using the Compare or Equal methods, or using an AddressComparator.
type IPv6AddressKey struct {
	keyContents
}

// ToAddress converts back to an address instance.
func (key IPv6AddressKey) ToAddress() *IPv6Address {
	return fromIPv6Key(key)
}

// String calls the String method in the corresponding address.
func (key IPv6AddressKey) String() string {
	return key.ToAddress().String()
}

// MACAddressKey is a representation of a MAC address that is comparable as defined by the language specification.
// See https://go.dev/ref/spec#Comparison_operators
//
// It can be used as a map key.  It can be obtained from its originating address instances.
// The zero value corresponds to the zero-value for MACAddress.
// Keys do not incorporate prefix length to ensure that all equal addresses have equal keys.
// To create a key that has prefix length, combine into a struct with the PrefixKey obtained by passing the address into PrefixKeyFrom.
// MACAddress can be compared using the Compare or Equal methods, or using an AddressComparator.
type MACAddressKey struct {
	vals struct {
		lower,
		upper uint64
	}
	additionalByteCount uint8 // 0 for MediaAccessControlSegmentCount or 2 for ExtendedUniqueIdentifier64SegmentCount
}

// ToAddress converts back to an address instance.
func (key MACAddressKey) ToAddress() *MACAddress {
	return fromMACKey(key)
}

// String calls the String method in the corresponding address.
func (key MACAddressKey) String() string {
	return key.ToAddress().String()
}

// KeyConstraint is the generic type constraint for an address type that can be generated from a generic address key.
type KeyConstraint[T any] interface {
	fmt.Stringer
	fromKey(addressScheme, *keyContents) T // implemented by IPAddress and Address
}

type addressScheme byte

const (
	adaptiveZeroScheme addressScheme = 0 // adaptiveZeroScheme needs to be zero, to coincide with the zero value for Address and IPAddress, which is a zero-length address
	ipv4Scheme         addressScheme = 1
	ipv6Scheme         addressScheme = 2
	mac48Scheme        addressScheme = 3
	eui64Scheme        addressScheme = 4
)

// KeyGeneratorConstraint is the generic type constraint for an address type that can generate a generic address key.
type KeyGeneratorConstraint[T KeyConstraint[T]] interface {
	ToGenericKey() Key[T]
}

// GenericKeyConstraint is the generic type constraint for an address type that can generate and be generated from a generic address key.
type GenericKeyConstraint[T KeyConstraint[T]] interface {
	KeyGeneratorConstraint[T]
	KeyConstraint[T]
}

// Key is a representation of an address that is comparable as defined by the language specification.
// See https://go.dev/ref/spec#Comparison_operators
//
// It can be used as a map key.  It can be obtained from its originating address instances.
// The zero value corresponds to the zero-value for its generic address type.
// Keys do not incorporate prefix length to ensure that all equal addresses have equal keys.
// To create a key that has prefix length, combine into a struct with the PrefixKey obtained by passing the address into PrefixKeyFrom.
type Key[T KeyConstraint[T]] struct {
	scheme addressScheme
	keyContents
}

// ToAddress converts back to an address instance.
func (key Key[T]) ToAddress() T {
	var t T
	return t.fromKey(key.scheme, &key.keyContents)
}

// String calls the String method in the corresponding address.
func (key Key[T]) String() string {
	return key.ToAddress().String()
}

type keyContents struct {
	vals [2]struct {
		lower,
		upper uint64
	}
	zone Zone
}

type (
	AddressKey             = Key[*Address]
	IPAddressKey           = Key[*IPAddress]
	IPAddressSeqRangeKey   = SequentialRangeKey[*IPAddress]
	IPv4AddressSeqRangeKey = SequentialRangeKey[*IPv4Address]
	IPv6AddressSeqRangeKey = SequentialRangeKey[*IPv6Address]
)

var (
	_ Key[*IPv4Address]
	_ Key[*IPv6Address]
	_ Key[*MACAddress]

	_ AddressKey   // Key[*Address]
	_ IPAddressKey // Key[*IPAddress]
	_ IPv4AddressKey
	_ IPv6AddressKey
	_ MACAddressKey

	_ IPAddressSeqRangeKey   // SequentialRangeKey[*IPAddress]
	_ IPv4AddressSeqRangeKey // SequentialRangeKey[*IPv4Address]
	_ IPv6AddressSeqRangeKey // SequentialRangeKey[*IPv6Address]
)

// PrefixKey is a representation of a prefix length that is comparable as defined by the language specification.
// See https://go.dev/ref/spec#Comparison_operators
//
// It can be used as a map key.
// The zero value is the absence of a prefix length.
type PrefixKey struct {
	// If true, the prefix length is indicated by PrefixLen.
	// If false, this indicates no prefix length for the associated address or subnet.
	IsPrefixed bool

	// If IsPrefixed is true, this holds the prefix length.
	// Otherwise, this should be zero if you wish that each address has a unique key.
	PrefixLen PrefixBitCount
}

// ToPrefixLen converts this key to its corresponding prefix length.
func (pref PrefixKey) ToPrefixLen() PrefixLen {
	if pref.IsPrefixed {
		return &pref.PrefixLen
	}
	return nil
}

func PrefixKeyFrom(addr AddressType) PrefixKey {
	if addr.IsPrefixed() {
		return PrefixKey{
			IsPrefixed: true,
			PrefixLen:  *addr.ToAddressBase().getPrefixLen(), // doing this instead of calling GetPrefixLen() on AddressType avoids the prefix len copy
		}
	}
	return PrefixKey{}
}

// TODO LATER serialization of addresses using https://pkg.go.dev/encoding/gob#GobDecoder - need to use custom encoder and decoder
// You may need new types that have both the prefix key and the address key, due to the MarshalBinary and UnmarshalBinary signatures.
// Don't encode the prefix unless there is one.  So you may need to encode a boolean for that.
// Of course, using strings is an alternative that is faily effective, but not quite as effective.
// For instance, IPv4AddressKey is a uint64, 8 bytes, while an ascii string is 15 bytes.
// IPv6AddressKey is 32 bytes.  In fact, to make this worthwhile, you should check for multiple, check for prefix, check for zone, use a bit for each.
// Otherwise, there is no point, using a string is just as good.
// https://stackoverflow.com/questions/28020070/golang-serialize-and-deserialize-back
// https://stackoverflow.com/questions/12854125/how-do-i-dump-the-struct-into-the-byte-array-without-reflection/12854659#12854659
