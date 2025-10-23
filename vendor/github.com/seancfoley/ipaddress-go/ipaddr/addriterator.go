//
// Copyright 2020-2022 Sean C Foley
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

// Iterator iterates collections, such as subnets and sequential address ranges.
// Use StdPushIterator or StdPullIterator to convert an Iterator to a standard library iterator.
type Iterator[T any] interface {
	// HasNext returns true if there is another item to iterate, false otherwise.
	HasNext() bool

	// Next returns the next item, or the zero value for T if there is none left.
	Next() T
}

// IteratorWithRemove is an iterator that provides a removal operation.
// Use NewPointIteratorWithRemove followed by StdPushIterator or StdPullIterator to convert an IteratorWithRemove to a standard library iterator.
type IteratorWithRemove[T any] interface {
	Iterator[T]

	// Remove removes the last iterated item from the underlying data structure or collection, and returns that element.
	// If there is no such element, it returns the zero value for T.
	Remove() T
}

type singleIterator[T any] struct {
	empty    bool
	original T
}

func (it *singleIterator[T]) HasNext() bool {
	return !it.empty
}

func (it *singleIterator[T]) Next() (res T) {
	if it.HasNext() {
		res = it.original
		it.empty = true
	}
	return
}

type multiAddrIterator struct {
	Iterator[*AddressSection]
	zone Zone
}

func (it multiAddrIterator) Next() (res *Address) {
	if it.HasNext() {
		sect := it.Iterator.Next()
		res = createAddress(sect, it.zone)
	}
	return
}

func nilAddrIterator() Iterator[*Address] {
	return &singleIterator[*Address]{}
}

func nilIterator[T any]() Iterator[T] {
	return &singleIterator[T]{}
}

func addrIterator(
	single bool,
	original *Address,
	prefixLen PrefixLen,
	valsAreMultiple bool,
	iterator Iterator[[]*AddressDivision]) Iterator[*Address] {
	if single {
		return &singleIterator[*Address]{original: original}
	}
	return multiAddrIterator{
		Iterator: &multiSectionIterator{
			original:        original.section,
			iterator:        iterator,
			valsAreMultiple: valsAreMultiple,
			prefixLen:       prefixLen,
		},
		zone: original.zone,
	}
}

func prefixAddrIterator(
	single bool,
	original *Address,
	prefixLen PrefixLen,
	iterator Iterator[[]*AddressDivision]) Iterator[*Address] {
	if single {
		return &singleIterator[*Address]{original: original}
	}
	var zone Zone
	if original != nil {
		zone = original.zone
	}
	return multiAddrIterator{
		Iterator: &prefixSectionIterator{
			original:  original.section,
			iterator:  iterator,
			prefixLen: prefixLen,
		},
		zone: zone,
	}
}

// this one is used by the sequential ranges
func rangeAddrIterator(
	single bool,
	original *Address,
	prefixLen PrefixLen,
	valsAreMultiple bool,
	iterator Iterator[[]*AddressDivision]) Iterator[*Address] {
	return addrIterator(single, original, prefixLen, valsAreMultiple, iterator)
}

type ipAddrIterator struct {
	Iterator[*Address]
}

func (iter ipAddrIterator) Next() *IPAddress {
	return iter.Iterator.Next().ToIP()
}

type sliceIterator[T any] struct {
	elements []T
}

func (iter *sliceIterator[T]) HasNext() bool {
	return len(iter.elements) > 0
}

func (iter *sliceIterator[T]) Next() (res T) {
	if iter.HasNext() {
		res = iter.elements[0]
		iter.elements = iter.elements[1:]
	}
	return
}

type ipv4AddressIterator struct {
	Iterator[*Address]
}

func (iter ipv4AddressIterator) Next() *IPv4Address {
	return iter.Iterator.Next().ToIPv4()
}

type ipv6AddressIterator struct {
	Iterator[*Address]
}

func (iter ipv6AddressIterator) Next() *IPv6Address {
	return iter.Iterator.Next().ToIPv6()
}

type macAddressIterator struct {
	Iterator[*Address]
}

func (iter macAddressIterator) Next() *MACAddress {
	return iter.Iterator.Next().ToMAC()
}

type addressSeriesIterator struct {
	Iterator[*Address]
}

func (iter addressSeriesIterator) Next() ExtendedSegmentSeries {
	if !iter.HasNext() {
		return nil
	}
	return wrapAddress(iter.Iterator.Next())
}

type ipaddressSeriesIterator struct {
	Iterator[*IPAddress]
}

func (iter ipaddressSeriesIterator) Next() ExtendedIPSegmentSeries {
	if !iter.HasNext() {
		return nil
	}
	return iter.Iterator.Next().Wrap()
}

type sectionSeriesIterator struct {
	Iterator[*AddressSection]
}

func (iter sectionSeriesIterator) Next() ExtendedSegmentSeries {
	if !iter.HasNext() {
		return nil
	}
	return wrapSection(iter.Iterator.Next())
}

type ipSectionSeriesIterator struct {
	Iterator[*IPAddressSection]
}

func (iter ipSectionSeriesIterator) Next() ExtendedIPSegmentSeries {
	if !iter.HasNext() {
		return nil
	}
	return wrapIPSection(iter.Iterator.Next())
}

// StdPushIterator converts a "pull" iterator in this libary to a "push" iterator assignable to the type iter.Seq in the standard library.
//
// The returned iterator is a single-use iterator.
//
// This function does not return iter.Seq directly, instead it returns a func(yield func(V) bool) assignable to a variable of type iter.Seq[V].
// This avoids adding a dependency of this libary on Go version 1.23 while still integrating with the iter package introduced with Go 1.23.
//
// To convert an instance of IteratorWithRemove, wrap it by calling NewPointIteratorWithRemove first, then pass the returned iterator this function.
// To convert an instance of CachingTrieIterator, wrap it by calling NewPointCachingTrieIterator first, then pass the returned iterator this function.
//
// You should avoid doing a double conversion on an iterator from this library,
// first to a "push" iterator with StdPushIterator and then to a "pull" iterator using iter.Pull in the standard libary.
// The result is an iterator less efficient than the original that also requires a call to the "stop" function to release resources.
// Instead, use StdPullIterator to get a pull iterator with an API similar to that provided by iter.Pull.
func StdPushIterator[V any](iterator Iterator[V]) func(yield func(V) bool) {
	return func(yield func(V) bool) {
		for iterator.HasNext() && yield(iterator.Next()) {
		}
	}
}

// StdPullIterator converts a "pull" iterator from this library to a standard library "pull" iterator
// consisting of a single function for iterating, and a second function for stopping.
//
// Note that the stop function is a no-op for all iterators in this library.
// It does nothing, and can be ignored.  It is provided only to match the returned values of iter.Pull.
// To convert an instance of IteratorWithRemove, wrap it by calling NewPointIteratorWithRemove first, then pass the returned iterator this function.
// To convert an instance of CachingTrieIterator, wrap it by calling NewPointCachingTrieIterator first, then pass the returned iterator this function.
//
// This function produces an iterator equivalent to the original.  It is a single-use iterator, like the original.
//
// Use this function rather than a double conversion: firstly to a "push" iterator with StdPushIterator,
// and secondly to a "pull" iterator using iter.Pull from the standard libary.
// The double-conversion produces a final iterator less efficient than the original.
func StdPullIterator[V any](iterator Iterator[V]) (next func() (V, bool), stop func()) {
	stop, next = func() {}, func() (V, bool) {
		hasNext := iterator.HasNext()
		return iterator.Next(), hasNext
	}
	return
}

// NewPointIteratorWithRemove can be used to convert an IteratorWithRemove to standard libary iterators.
// Call this function, and then pass the returned iterator to either StdPushIterator or StdPullIterator.
// This preserves access to the removal operation of the original iterator, now available through each element, IteratorWithRemovePosition.
func NewPointIteratorWithRemove[V any](iterator IteratorWithRemove[V]) Iterator[IteratorWithRemovePosition[V]] {
	wrappedIter := pointIteratorWithRemove[V]{IteratorWithRemove: iterator}
	castIter := Iterator[IteratorWithRemovePosition[V]](&wrappedIter)
	return castIter
}

type pointIteratorWithRemove[V any] struct {
	IteratorWithRemove[V]
}

func (iter pointIteratorWithRemove[V]) Next() IteratorWithRemovePosition[V] {
	return IteratorWithRemovePosition[V]{
		v:    iter.IteratorWithRemove.Next(),
		iter: iter.IteratorWithRemove,
	}
}

func (iter pointIteratorWithRemove[V]) HasNext() bool {
	return iter.IteratorWithRemove.HasNext()
}

// IteratorWithRemovePosition is an element returned from a PointIteratorWithRemove.
type IteratorWithRemovePosition[V any] struct {
	v    V
	iter IteratorWithRemove[V]
}

// Value returns the iterator value associated with this iterator position.
func (iterPosition IteratorWithRemovePosition[V]) Value() V {
	return iterPosition.v
}

// Remove removes the current iterated value from the underlying data structure or collection, and returns that element.
// If there is no such element, because it has been removed already or there are no more iterated elements, it returns the zero value for T.
func (iterPosition IteratorWithRemovePosition[V]) Remove() V {
	return iterPosition.iter.Remove()
}

// NewPointCachingTrieIterator can be used to convert a CachingTrieIterator to standard libary iterators.
// Call this function, and then pass the returned iterator to either StdPushIterator or StdPullIterator.
// This preserves access to the removal and caching operations of the original iterator, now available through each element of type CachingTrieIteratorPosition.
func NewPointCachingTrieIterator[V any](iterator CachingTrieIterator[V]) Iterator[CachingTrieIteratorPosition[V]] {
	wrappedIter := pointCachingTrieIterator[V]{CachingTrieIterator: iterator}
	castIter := Iterator[CachingTrieIteratorPosition[V]](&wrappedIter)
	return castIter
}

type pointCachingTrieIterator[V any] struct {
	CachingTrieIterator[V]
}

func (iter pointCachingTrieIterator[V]) Next() CachingTrieIteratorPosition[V] {
	return CachingTrieIteratorPosition[V]{
		v:    iter.CachingTrieIterator.Next(),
		iter: iter.CachingTrieIterator,
	}
}

func (iter pointCachingTrieIterator[V]) HasNext() bool {
	return iter.CachingTrieIterator.HasNext()
}

// CachingTrieIteratorPosition is an element returned from an iterator created with NewPointCachingTrieIterator.
type CachingTrieIteratorPosition[V any] struct {
	v    V
	iter CachingTrieIterator[V]
}

// Value returns the iterator value associated with this iterator position.
func (iterPosition CachingTrieIteratorPosition[V]) Value() V {
	return iterPosition.v
}

// Remove removes the current iterated value from the underlying data structure or collection, and returns that element.
// If there is no such element, because it has been removed already or there are no more iterated elements, it returns the zero value for T.
func (iterPosition CachingTrieIteratorPosition[V]) Remove() V {
	return iterPosition.iter.Remove()
}

func (iterPosition CachingTrieIteratorPosition[V]) GetCached() Cached {
	return iterPosition.iter.GetCached()
}

func (iterPosition CachingTrieIteratorPosition[V]) CacheWithLowerSubNode(cached Cached) bool {
	return iterPosition.iter.CacheWithLowerSubNode(cached)
}

func (iterPosition CachingTrieIteratorPosition[V]) CacheWithUpperSubNode(cached Cached) bool {
	return iterPosition.iter.CacheWithUpperSubNode(cached)
}
