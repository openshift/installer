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

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

// ExtendedSegmentSeries wraps either an Address or AddressSection.
// ExtendedSegmentSeries can be used to write code that works with both addresses and address sections,
// going further than AddressSegmentSeries to offer additional methods with the series types in their signature.
type ExtendedSegmentSeries interface {
	AddressSegmentSeries

	// Unwrap returns the wrapped address or address section as an interface, AddressSegmentSeries.
	Unwrap() AddressSegmentSeries

	// Equal returns whether the given address series is equal to this address series.
	// Two address series are equal if they represent the same set of series.
	// Both must be equal addresses or both must be equal sections.
	Equal(ExtendedSegmentSeries) bool

	// Contains returns whether this is same type and version as the given address series and whether it contains all values in the given series.
	//
	// Series must also have the same number of segments to be comparable, otherwise false is returned.
	Contains(ExtendedSegmentSeries) bool

	// GetSection returns the backing section for this series, comprising all segments.
	GetSection() *AddressSection

	// GetTrailingSection returns an ending subsection of the full address section.
	GetTrailingSection(index int) *AddressSection

	// GetSubSection returns a subsection of the full address section.
	GetSubSection(index, endIndex int) *AddressSection

	// GetSegment returns the segment at the given index.
	// The first segment is at index 0.
	// GetSegment will panic given a negative index or an index matching or larger than the segment count.
	GetSegment(index int) *AddressSegment

	// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
	GetSegments() []*AddressSegment

	// CopySegments copies the existing segments into the given slice,
	// as much as can be fit into the slice, returning the number of segments copied.
	CopySegments(segs []*AddressSegment) (count int)

	// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
	// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
	CopySubSegments(start, end int, segs []*AddressSegment) (count int)

	// IsIP returns true if this series originated as an IPv4 or IPv6 series, or a zero-length IP series.  If so, use ToIP to convert back to the IP-specific type.
	IsIP() bool

	// IsIPv4 returns true if this series originated as an IPv4 series.  If so, use ToIPv4 to convert back to the IPv4-specific type.
	IsIPv4() bool

	// IsIPv6 returns true if this series originated as an IPv6 series.  If so, use ToIPv6 to convert back to the IPv6-specific type.
	IsIPv6() bool

	// IsMAC returns true if this series originated as a MAC series.  If so, use ToMAC to convert back to the MAC-specific type.
	IsMAC() bool

	// ToIP converts to an IPAddressSegmentSeries if this series originated as IPv4 or IPv6, or an implicitly zero-valued IP.
	// If not, ToIP returns nil.
	ToIP() IPAddressSegmentSeries

	// ToIPv4 converts to an IPv4AddressSegmentSeries if this series originated as an IPv4 series.
	// If not, ToIPv4 returns nil.
	//
	// ToIPv4 implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToIPv4() IPv4AddressSegmentSeries

	// ToIPv6 converts to an IPv4AddressSegmentSeries if this series originated as an IPv6 series.
	// If not, ToIPv6 returns nil.
	//
	// ToIPv6 implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToIPv6() IPv6AddressSegmentSeries

	// ToMAC converts to a MACAddressSegmentSeries if this series originated as a MAC series.
	// If not, ToMAC returns nil.
	//
	// ToMAC implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToMAC() MACAddressSegmentSeries

	// ToBlock creates a new series block by changing the segment at the given index to have the given lower and upper value,
	// and changing the following segments to be full-range.
	ToBlock(segmentIndex int, lower, upper SegInt) ExtendedSegmentSeries

	// ToPrefixBlock returns the series with the same prefix as this series while the remaining bits span all values.
	// The series will be the block of all series with the same prefix.
	//
	// If this series has no prefix, this series is returned.
	ToPrefixBlock() ExtendedSegmentSeries

	// ToPrefixBlockLen returns the series with the same prefix of the given length as this series while the remaining bits span all values.
	// The returned series will be the block of all series with the same prefix.
	ToPrefixBlockLen(prefLen BitCount) ExtendedSegmentSeries

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
	Increment(int64) ExtendedSegmentSeries

	// IncrementBoundary returns the item that is the given increment from the range boundaries of this item.
	//
	// If the given increment is positive, adds the value to the highest (GetUpper) in the range to produce a new item.
	// If the given increment is negative, adds the value to the lowest (GetLower) in the range to produce a new item.
	// If the increment is zero, returns this.
	//
	// If this represents just a single value, this item is simply incremented by the given increment value, positive or negative.
	//
	// On overflow or underflow, IncrementBoundary returns nil.
	IncrementBoundary(int64) ExtendedSegmentSeries

	// Enumerate indicates where an address series sits relative to the range ordering.
	//
	// Determines how many address series elements of a range precede the given address series element, if the address series is in the range.
	// If above the range, it is the distance to the upper boundary added to the range count less one, and if below the range, the distance to the lower boundary.
	//
	// In other words, if the given address series is not in the range but above it, returns the number of address series preceding the address series from the upper range boundary,
	// added to one less than the total number of range address series.  If the given address series is not in the subnet but below it, returns the number of address series following the address to the lower subnet boundary.
	//
	// Returns nil when the argument is multi-valued. The argument must be an individual address series.
	//
	// When this is also an individual address series, the returned value is the distance (difference) between the two address series values.
	//
	// If the given address does not have the same version or type, then nil is returned.
	Enumerate(ExtendedSegmentSeries) *big.Int

	// GetLower returns the series in the range with the lowest numeric value,
	// which will be the same series if it represents a single value.
	GetLower() ExtendedSegmentSeries

	// GetUpper returns the series in the range with the highest numeric value,
	// which will be the same series if it represents a single value.
	GetUpper() ExtendedSegmentSeries

	// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this series.
	// The returned block will have an assigned prefix length indicating the prefix length for the block.
	//
	// There may be no such series - it is required that the range of values match the range of a prefix block.
	// If there is no such series, then nil is returned.
	AssignPrefixForSingleBlock() ExtendedSegmentSeries

	// AssignMinPrefixForBlock returns an equivalent series, assigned the smallest prefix length possible,
	// such that the prefix block for that prefix length is in this series.
	//
	// In other words, this method assigns a prefix length to this series matching the largest prefix block in this series.
	AssignMinPrefixForBlock() ExtendedSegmentSeries

	// Iterator provides an iterator to iterate through the individual series of this series.
	//
	// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual series.
	//
	// Call IsMultiple to determine if this instance represents multiple series, or GetCount for the count.
	Iterator() Iterator[ExtendedSegmentSeries]

	// PrefixIterator provides an iterator to iterate through the individual prefixes of this series,
	// each iterated element spanning the range of values for its prefix.
	//
	// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
	// instead constraining themselves to values from this series.
	//
	// If the series has no prefix length, then this is equivalent to Iterator.
	PrefixIterator() Iterator[ExtendedSegmentSeries]

	// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this series.
	// Each iterated series will be a prefix block with the same prefix length as this series.
	//
	// If this series has no prefix length, then this is equivalent to Iterator.
	PrefixBlockIterator() Iterator[ExtendedSegmentSeries]

	// AdjustPrefixLen increases or decreases the prefix length by the given increment.
	//
	// A prefix length will not be adjusted lower than zero or beyond the bit length of the series.
	//
	// If this series has no prefix length, then the prefix length will be set to the adjustment if positive,
	// or it will be set to the adjustment added to the bit count if negative.
	AdjustPrefixLen(BitCount) ExtendedSegmentSeries

	// AdjustPrefixLenZeroed increases or decreases the prefix length by the given increment while zeroing out the bits that have moved into or outside the prefix.
	//
	// A prefix length will not be adjusted lower than zero or beyond the bit length of the series.
	//
	// If this series has no prefix length, then the prefix length will be set to the adjustment if positive,
	// or it will be set to the adjustment added to the bit count if negative.
	//
	// When prefix length is increased, the bits moved within the prefix become zero.
	// When a prefix length is decreased, the bits moved outside the prefix become zero.
	AdjustPrefixLenZeroed(BitCount) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError)

	// SetPrefixLen sets the prefix length.
	//
	// A prefix length will not be set to a value lower than zero or beyond the bit length of the series.
	// The provided prefix length will be adjusted to these boundaries if necessary.
	SetPrefixLen(BitCount) ExtendedSegmentSeries

	// SetPrefixLenZeroed sets the prefix length.
	//
	// A prefix length will not be set to a value lower than zero or beyond the bit length of the series.
	// The provided prefix length will be adjusted to these boundaries if necessary.
	//
	// If this series has a prefix length, and the prefix length is increased when setting the new prefix length, the bits moved within the prefix become zero.
	// If this series has a prefix length, and the prefix length is decreased when setting the new prefix length, the bits moved outside the prefix become zero.
	//
	// In other words, bits that move from one side of the prefix length to the other (bits moved into the prefix or outside the prefix) are zeroed.
	//
	// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
	SetPrefixLenZeroed(BitCount) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError)

	// WithoutPrefixLen provides the same address series but with no prefix length.  The values remain unchanged.
	WithoutPrefixLen() ExtendedSegmentSeries

	// ReverseBytes returns a new segment series with the bytes reversed.  Any prefix length is dropped.
	//
	// If each segment is more than 1 byte long, and the bytes within a single segment cannot be reversed because the segment represents a range,
	// and reversing the segment values results in a range that is not contiguous, then this returns an error.
	//
	// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
	ReverseBytes() (ExtendedSegmentSeries, addrerr.IncompatibleAddressError)

	// ReverseBits returns a new segment series with the bits reversed.  Any prefix length is dropped.
	//
	// If the bits within a single segment cannot be reversed because the segment represents a range,
	// and reversing the segment values results in a range that is not contiguous, this returns an error.
	//
	// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
	ReverseBits(perByte bool) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError)

	// ReverseSegments returns a new series with the segments reversed.
	ReverseSegments() ExtendedSegmentSeries

	// ToCustomString creates a customized string from this series according to the given string option parameters.
	ToCustomString(stringOptions addrstr.StringOptions) string
}

// WrappedAddress is the implementation of ExtendedSegmentSeries for addresses.
type WrappedAddress struct {
	*Address
}

// Unwrap returns the wrapped address as an interface, AddressSegmentSeries.
func (addr WrappedAddress) Unwrap() AddressSegmentSeries {
	res := addr.Address
	if res == nil {
		return nil
	}
	return res
}

// ToIPv4 converts to an IPv4AddressSegmentSeries if this series originated as an IPv4 series.
// If not, ToIPv4 returns nil.
//
// ToIPv4 implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr WrappedAddress) ToIPv4() IPv4AddressSegmentSeries {
	return addr.Address.ToIPv4()
}

// ToIPv6 converts to an IPv4AddressSegmentSeries if this series originated as an IPv6 series.
// If not, ToIPv6 returns nil.
//
// ToIPv6 implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr WrappedAddress) ToIPv6() IPv6AddressSegmentSeries {
	return addr.Address.ToIPv6()
}

// ToIP converts to an IP address if this originated as IPv4 or IPv6, or an implicitly zero-valued IP.
// If not, ToIP returns nil.
func (addr WrappedAddress) ToIP() IPAddressSegmentSeries {
	return addr.Address.ToIP()
}

// ToMAC converts to a MACAddressSegmentSeries if this series originated as a MAC series.
// If not, ToMAC returns nil.
//
// ToMAC implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr WrappedAddress) ToMAC() MACAddressSegmentSeries {
	return addr.Address.ToMAC()
}

// Iterator provides an iterator to iterate through the individual series of this series.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual series.
//
// Call IsMultiple to determine if this instance represents multiple series, or GetCount for the count.
func (addr WrappedAddress) Iterator() Iterator[ExtendedSegmentSeries] {
	return addressSeriesIterator{addr.Address.Iterator()}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this series,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this series.
//
// If the series has no prefix length, then this is equivalent to Iterator.
func (addr WrappedAddress) PrefixIterator() Iterator[ExtendedSegmentSeries] {
	return addressSeriesIterator{addr.Address.PrefixIterator()}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this series.
// Each iterated series will be a prefix block with the same prefix length as this series.
//
// If this series has no prefix length, then this is equivalent to Iterator.
func (addr WrappedAddress) PrefixBlockIterator() Iterator[ExtendedSegmentSeries] {
	return addressSeriesIterator{addr.Address.PrefixBlockIterator()}
}

// ToBlock creates a new series block by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (addr WrappedAddress) ToBlock(segmentIndex int, lower, upper SegInt) ExtendedSegmentSeries {
	return wrapAddress(addr.Address.ToBlock(segmentIndex, lower, upper))
}

// ToPrefixBlock returns the series with the same prefix as this series while the remaining bits span all values.
// The series will be the block of all series with the same prefix.
//
// If this series has no prefix, this series is returned.
func (addr WrappedAddress) ToPrefixBlock() ExtendedSegmentSeries {
	return wrapAddress(addr.Address.ToPrefixBlock())
}

// ToPrefixBlockLen returns the series with the same prefix of the given length as this series while the remaining bits span all values.
// The returned series will be the block of all series with the same prefix.
func (addr WrappedAddress) ToPrefixBlockLen(prefLen BitCount) ExtendedSegmentSeries {
	return wrapAddress(addr.Address.ToPrefixBlockLen(prefLen))
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
func (addr WrappedAddress) Increment(i int64) ExtendedSegmentSeries {
	return convAddrToIntf(addr.Address.Increment(i))
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
func (addr WrappedAddress) IncrementBoundary(i int64) ExtendedSegmentSeries {
	return convAddrToIntf(addr.Address.IncrementBoundary(i))
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
// If the given argument is not an address or does not have the same address version or type, then nil is returned.
func (addr WrappedAddress) Enumerate(other ExtendedSegmentSeries) *big.Int {
	if a, ok := other.Unwrap().(AddressType); ok {
		return addr.Address.Enumerate(a)
	}
	return nil
}

// GetLower returns the series in the range with the lowest numeric value,
// which will be the same series if it represents a single value.
func (addr WrappedAddress) GetLower() ExtendedSegmentSeries {
	return wrapAddress(addr.Address.GetLower())
}

// GetUpper returns the series in the range with the highest numeric value,
// which will be the same series if it represents a single value.
func (addr WrappedAddress) GetUpper() ExtendedSegmentSeries {
	return wrapAddress(addr.Address.GetUpper())
}

// GetSection returns the backing section for this series, comprising all segments.
func (addr WrappedAddress) GetSection() *AddressSection {
	return addr.Address.GetSection()
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this series.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such series - it is required that the range of values match the range of a prefix block.
// If there is no such series, then nil is returned.
func (addr WrappedAddress) AssignPrefixForSingleBlock() ExtendedSegmentSeries {
	return convAddrToIntf(addr.Address.AssignPrefixForSingleBlock())
}

// AssignMinPrefixForBlock returns an equivalent series, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this series.
//
// In other words, this method assigns a prefix length to this series matching the largest prefix block in this series.
func (addr WrappedAddress) AssignMinPrefixForBlock() ExtendedSegmentSeries {
	return wrapAddress(addr.Address.AssignMinPrefixForBlock())
}

// WithoutPrefixLen provides the same address series but with no prefix length.  The values remain unchanged.
func (addr WrappedAddress) WithoutPrefixLen() ExtendedSegmentSeries {
	return wrapAddress(addr.Address.WithoutPrefixLen())
}

// Contains returns whether this is same type and version as the given address series and whether it contains all values in the given series.
//
// Series must also have the same number of segments to be comparable, otherwise false is returned.
func (addr WrappedAddress) Contains(other ExtendedSegmentSeries) bool {
	a, ok := other.Unwrap().(AddressType)
	return ok && addr.Address.Contains(a)
}

// Equal returns whether the given address series is equal to this address series.
// Two address series are equal if they represent the same set of series.
// Both must be equal addresses.
func (addr WrappedAddress) Equal(other ExtendedSegmentSeries) bool {
	a, ok := other.Unwrap().(AddressType)
	return ok && addr.Address.Equal(a)
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the series.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (addr WrappedAddress) SetPrefixLen(prefixLen BitCount) ExtendedSegmentSeries {
	return wrapAddress(addr.Address.SetPrefixLen(prefixLen))
}

// SetPrefixLenZeroed sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the series.
// The provided prefix length will be adjusted to these boundaries if necessary.
//
// If this series has a prefix length, and the prefix length is increased when setting the new prefix length, the bits moved within the prefix become zero.
// If this series has a prefix length, and the prefix length is decreased when setting the new prefix length, the bits moved outside the prefix become zero.
//
// In other words, bits that move from one side of the prefix length to the other (bits moved into the prefix or outside the prefix) are zeroed.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (addr WrappedAddress) SetPrefixLenZeroed(prefixLen BitCount) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	return wrapAddrWithErr(addr.Address.SetPrefixLenZeroed(prefixLen))
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the series.
//
// If this series has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (addr WrappedAddress) AdjustPrefixLen(prefixLen BitCount) ExtendedSegmentSeries {
	return wrapAddress(addr.Address.AdjustPrefixLen(prefixLen))
}

// AdjustPrefixLenZeroed increases or decreases the prefix length by the given increment while zeroing out the bits that have moved into or outside the prefix.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the series.
//
// If this series has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
//
// When prefix length is increased, the bits moved within the prefix become zero.
// When a prefix length is decreased, the bits moved outside the prefix become zero.
func (addr WrappedAddress) AdjustPrefixLenZeroed(prefixLen BitCount) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	return wrapAddrWithErr(addr.Address.AdjustPrefixLenZeroed(prefixLen))
}

// ReverseBytes returns a new segment series with the bytes reversed.  Any prefix length is dropped.
//
// If each segment is more than 1 byte long, and the bytes within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, then this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
func (addr WrappedAddress) ReverseBytes() (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	return wrapAddrWithErr(addr.Address.ReverseBytes())
}

// ReverseBits returns a new segment series with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (addr WrappedAddress) ReverseBits(perByte bool) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	a, err := addr.Address.ReverseBits(perByte)
	if err != nil {
		return nil, err
	}
	return wrapAddress(a), nil
}

// ReverseSegments returns a new series with the segments reversed.
func (addr WrappedAddress) ReverseSegments() ExtendedSegmentSeries {
	return wrapAddress(addr.Address.ReverseSegments())
}

// WrappedAddressSection is the implementation of ExtendedSegmentSeries for address sections.
type WrappedAddressSection struct {
	*AddressSection
}

// Unwrap returns the wrapped address section as an interface, AddressSegmentSeries.
func (section WrappedAddressSection) Unwrap() AddressSegmentSeries {
	res := section.AddressSection
	if res == nil {
		return nil
	}
	return res
}

// ToIPv4 converts to an IPv4AddressSegmentSeries if this series originated as an IPv4 series.
// If not, ToIPv4 returns nil.
//
// ToIPv4 implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section WrappedAddressSection) ToIPv4() IPv4AddressSegmentSeries {
	return section.AddressSection.ToIPv4()
}

// ToIPv6 converts to an IPv4AddressSegmentSeries if this series originated as an IPv6 series.
// If not, ToIPv6 returns nil.
//
// ToIPv6 implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section WrappedAddressSection) ToIPv6() IPv6AddressSegmentSeries {
	return section.AddressSection.ToIPv6()
}

// ToIP converts to an IP address section if this originated as IPv4 or IPv6, or an implicitly zero-valued IP.
// If not, ToIP returns nil.
func (section WrappedAddressSection) ToIP() IPAddressSegmentSeries {
	return section.AddressSection.ToIP()
}

// ToMAC converts to a MACAddressSegmentSeries if this series originated as a MAC series.
// If not, ToMAC returns nil.
//
// ToMAC implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section WrappedAddressSection) ToMAC() MACAddressSegmentSeries {
	return section.AddressSection.ToMAC()
}

// Iterator provides an iterator to iterate through the individual series of this series.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual series.
//
// Call IsMultiple to determine if this instance represents multiple series, or GetCount for the count.
func (section WrappedAddressSection) Iterator() Iterator[ExtendedSegmentSeries] {
	return sectionSeriesIterator{section.AddressSection.Iterator()}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this series,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this series.
//
// If the series has no prefix length, then this is equivalent to Iterator.
func (section WrappedAddressSection) PrefixIterator() Iterator[ExtendedSegmentSeries] {
	return sectionSeriesIterator{section.AddressSection.PrefixIterator()}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this series.
// Each iterated series will be a prefix block with the same prefix length as this series.
//
// If this series has no prefix length, then this is equivalent to Iterator.
func (section WrappedAddressSection) PrefixBlockIterator() Iterator[ExtendedSegmentSeries] {
	return sectionSeriesIterator{section.AddressSection.PrefixBlockIterator()}
}

// ToBlock creates a new series block by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (section WrappedAddressSection) ToBlock(segmentIndex int, lower, upper SegInt) ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.ToBlock(segmentIndex, lower, upper))
}

// ToPrefixBlock returns the series with the same prefix as this series while the remaining bits span all values.
// The series will be the block of all series with the same prefix.
//
// If this series has no prefix, this series is returned.
func (section WrappedAddressSection) ToPrefixBlock() ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.ToPrefixBlock())
}

// ToPrefixBlockLen returns the series with the same prefix of the given length as this series while the remaining bits span all values.
// The returned series will be the block of all series with the same prefix.
func (section WrappedAddressSection) ToPrefixBlockLen(prefLen BitCount) ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.ToPrefixBlockLen(prefLen))
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
func (section WrappedAddressSection) Increment(i int64) ExtendedSegmentSeries {
	return convSectToIntf(section.AddressSection.Increment(i))
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
func (section WrappedAddressSection) IncrementBoundary(i int64) ExtendedSegmentSeries {
	return convSectToIntf(section.AddressSection.IncrementBoundary(i))
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
// If the given argument is not an address section, or does not have the same version or type, then nil is returned.
//
// Sections must also have the same number of segments to be comparable, otherwise nil is returned.
func (section WrappedAddressSection) Enumerate(other ExtendedSegmentSeries) *big.Int {
	if a, ok := other.Unwrap().(AddressSectionType); ok {
		return section.AddressSection.Enumerate(a)
	}
	return nil
}

// GetLower returns the series in the range with the lowest numeric value,
// which will be the same series if it represents a single value.
func (section WrappedAddressSection) GetLower() ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.GetLower())
}

// GetUpper returns the series in the range with the highest numeric value,
// which will be the same series if it represents a single value.
func (section WrappedAddressSection) GetUpper() ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.GetUpper())
}

// GetSection returns the backing section for this series, comprising all segments.
func (section WrappedAddressSection) GetSection() *AddressSection {
	return section.AddressSection
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this series.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such series - it is required that the range of values match the range of a prefix block.
// If there is no such series, then nil is returned.
func (section WrappedAddressSection) AssignPrefixForSingleBlock() ExtendedSegmentSeries {
	return convSectToIntf(section.AddressSection.AssignPrefixForSingleBlock())
}

// AssignMinPrefixForBlock returns an equivalent series, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this series.
//
// In other words, this method assigns a prefix length to this series matching the largest prefix block in this series.
func (section WrappedAddressSection) AssignMinPrefixForBlock() ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.AssignMinPrefixForBlock())
}

// WithoutPrefixLen provides the same address series but with no prefix length.  The values remain unchanged.
func (section WrappedAddressSection) WithoutPrefixLen() ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.WithoutPrefixLen())
}

// Contains returns whether this is same type and version as the given address series and whether it contains all values in the given series.
//
// Series must also have the same number of segments to be comparable, otherwise false is returned.
func (section WrappedAddressSection) Contains(other ExtendedSegmentSeries) bool {
	s, ok := other.Unwrap().(AddressSectionType)
	return ok && section.AddressSection.Contains(s)
}

// Equal returns whether the given address series is equal to this address series.
// Two address series are equal if they represent the same set of series.
// Both must be equal sections.
func (section WrappedAddressSection) Equal(other ExtendedSegmentSeries) bool {
	s, ok := other.Unwrap().(AddressSectionType)
	return ok && section.AddressSection.Equal(s)
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the series.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (section WrappedAddressSection) SetPrefixLen(prefixLen BitCount) ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.SetPrefixLen(prefixLen))
}

// SetPrefixLenZeroed sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the series.
// The provided prefix length will be adjusted to these boundaries if necessary.
//
// If this series has a prefix length, and the prefix length is increased when setting the new prefix length, the bits moved within the prefix become zero.
// If this series has a prefix length, and the prefix length is decreased when setting the new prefix length, the bits moved outside the prefix become zero.
//
// In other words, bits that move from one side of the prefix length to the other (bits moved into the prefix or outside the prefix) are zeroed.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (section WrappedAddressSection) SetPrefixLenZeroed(prefixLen BitCount) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	return wrapSectWithErr(section.AddressSection.SetPrefixLenZeroed(prefixLen))
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the series.
//
// If this series has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (section WrappedAddressSection) AdjustPrefixLen(adjustment BitCount) ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.AdjustPrefixLen(adjustment))
}

// AdjustPrefixLenZeroed increases or decreases the prefix length by the given increment while zeroing out the bits that have moved into or outside the prefix.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the series.
//
// If this series has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
//
// When prefix length is increased, the bits moved within the prefix become zero.
// When a prefix length is decreased, the bits moved outside the prefix become zero.
func (section WrappedAddressSection) AdjustPrefixLenZeroed(adjustment BitCount) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	return wrapSectWithErr(section.AddressSection.AdjustPrefixLenZeroed(adjustment))
}

// ReverseBytes returns a new segment series with the bytes reversed.  Any prefix length is dropped.
//
// If each segment is more than 1 byte long, and the bytes within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, then this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
func (section WrappedAddressSection) ReverseBytes() (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	return wrapSectWithErr(section.AddressSection.ReverseBytes())
}

// ReverseBits returns a new segment series with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (section WrappedAddressSection) ReverseBits(perByte bool) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	return wrapSectWithErr(section.AddressSection.ReverseBits(perByte))
}

// ReverseSegments returns a new series with the segments reversed.
func (section WrappedAddressSection) ReverseSegments() ExtendedSegmentSeries {
	return wrapSection(section.AddressSection.ReverseSegments())
}

var _ ExtendedSegmentSeries = WrappedAddress{}
var _ ExtendedSegmentSeries = WrappedAddressSection{}

// In go, a nil value is not converted to a nil interface, it is converted to a non-nil interface instance with underlying value nil.
func convAddrToIntf(addr *Address) ExtendedSegmentSeries {
	if addr == nil {
		return nil
	}
	return wrapAddress(addr)
}

func convSectToIntf(sect *AddressSection) ExtendedSegmentSeries {
	if sect == nil {
		return nil
	}
	return wrapSection(sect)
}

func wrapSectWithErr(section *AddressSection, err addrerr.IncompatibleAddressError) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	if err == nil {
		return wrapSection(section), nil
	}
	return nil, err
}

func wrapAddrWithErr(addr *Address, err addrerr.IncompatibleAddressError) (ExtendedSegmentSeries, addrerr.IncompatibleAddressError) {
	if err == nil {
		return wrapAddress(addr), nil
	}
	return nil, err
}

func wrapAddress(addr *Address) WrappedAddress {
	return WrappedAddress{addr}
}

func wrapSection(section *AddressSection) WrappedAddressSection {
	return WrappedAddressSection{section}
}
