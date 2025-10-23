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

type singleSegmentsIterator struct {
	original []*AddressDivision
}

func (it *singleSegmentsIterator) HasNext() bool {
	return it.original != nil
}

func (it *singleSegmentsIterator) Next() (res []*AddressDivision) {
	if it.HasNext() {
		res = it.original
		it.original = nil
	}
	return
}

type multiSegmentsIterator struct {
	done       bool
	variations []Iterator[*AddressSegment]
	nextSet    []*AddressDivision

	segIteratorProducer,
	hostSegIteratorProducer func(int) Iterator[*AddressSegment]

	networkSegmentIndex,
	hostSegmentIndex int

	excludeFunc func([]*AddressDivision) bool
}

func (it *multiSegmentsIterator) HasNext() bool {
	return !it.done
}

func (it *multiSegmentsIterator) updateVariations(start int) {
	i := start
	nextSet := it.nextSet
	variations := it.variations
	segIteratorProducer := it.segIteratorProducer
	for ; i < it.hostSegmentIndex; i++ {
		variations[i] = segIteratorProducer(i)
		nextSet[i] = variations[i].Next().ToDiv()
	}
	if i == it.networkSegmentIndex {
		variations[i] = it.hostSegIteratorProducer(i)
		nextSet[i] = variations[i].Next().ToDiv()
	}
}

func (it *multiSegmentsIterator) init() {
	it.updateVariations(0)
	nextSet := it.nextSet
	variations := it.variations
	divCount := len(variations)
	hostSegIteratorProducer := it.hostSegIteratorProducer
	// for regular iterators (not prefix block), networkSegmentIndex is last segment (count - 1)
	for i := it.networkSegmentIndex + 1; i < divCount; i++ {
		variations[i] = hostSegIteratorProducer(i)
		nextSet[i] = variations[i].Next().ToDiv()
	}
	excludeFunc := it.excludeFunc
	if excludeFunc != nil && excludeFunc(it.nextSet) {
		it.increment()
	}
}

func (it *multiSegmentsIterator) Next() (res []*AddressDivision) {
	if it.HasNext() {
		res = it.increment()
	}
	return
}

func (it *multiSegmentsIterator) increment() (res []*AddressDivision) {
	var previousSegs []*AddressDivision
	// the current set of segments already holds the next iteration,
	// this searches for the set of segments to follow.
	variations := it.variations
	nextSet := it.nextSet
	for j := it.networkSegmentIndex; j >= 0; j-- { //for regular iterators (not prefix block), networkSegmentIndex is last segment (count - 1)
		for variations[j].HasNext() {
			if previousSegs == nil {
				previousSegs = clone(nextSet)
			}
			nextSet[j] = variations[j].Next().ToDiv()
			it.updateVariations(j + 1)
			excludeFunc := it.excludeFunc
			if excludeFunc != nil && excludeFunc(nextSet) {
				// try again, starting over
				j = it.networkSegmentIndex
			} else {
				return previousSegs
			}
		}
	}
	it.done = true
	if previousSegs == nil {
		// never found set of candidate segments
		return nextSet
	}
	// found a candidate to follow, but was rejected.
	// nextSet has that rejected candidate,
	// so we must return the set that was created prior to that.
	return previousSegs
}

// this iterator function used by addresses and segment arrays, for iterators that are not prefix or prefix block iterators
func allSegmentsIterator(
	divCount int,
	segSupplier func() []*AddressDivision, // only useful for a segment iterator.  Address/section iterators use address/section for single valued iterator.
	segIteratorProducer func(int) Iterator[*AddressSegment],
	excludeFunc func([]*AddressDivision) bool /* can be nil */) Iterator[[]*AddressDivision] {
	return segmentsIterator(divCount, segSupplier, segIteratorProducer, excludeFunc, divCount-1, divCount, nil)
}

// used to produce regular iterators with or without zero-host values, and prefix block iterators
func segmentsIterator(
	divCount int,
	segSupplier func() []*AddressDivision,
	segIteratorProducer func(int) Iterator[*AddressSegment], // unused at this time, since we do not have a public segments iterator
	excludeFunc func([]*AddressDivision) bool, // can be nil
	networkSegmentIndex,
	hostSegmentIndex int,
	hostSegIteratorProducer func(int) Iterator[*AddressSegment]) Iterator[[]*AddressDivision] { // returns Iterator<S[]>
	if segSupplier != nil {
		return &singleSegmentsIterator{segSupplier()}
	}
	iterator := &multiSegmentsIterator{
		variations:              make([]Iterator[*AddressSegment], divCount),
		nextSet:                 make([]*AddressDivision, divCount),
		segIteratorProducer:     segIteratorProducer,
		hostSegIteratorProducer: hostSegIteratorProducer,
		networkSegmentIndex:     networkSegmentIndex,
		hostSegmentIndex:        hostSegmentIndex,
		excludeFunc:             excludeFunc,
	}
	iterator.init()
	return iterator
}

// this iterator function used by sequential ranges
func rangeSegmentsIterator(
	divCount int,
	segIteratorProducer func(int) Iterator[*AddressSegment],
	networkSegmentIndex,
	hostSegmentIndex int,
	prefixedSegIteratorProducer func(int) Iterator[*AddressSegment]) Iterator[[]*AddressDivision] {
	return segmentsIterator(
		divCount,
		nil,
		segIteratorProducer,
		nil,
		networkSegmentIndex,
		hostSegmentIndex,
		prefixedSegIteratorProducer)
}

type singleSectionIterator struct {
	original *AddressSection
}

func (it *singleSectionIterator) HasNext() bool {
	return it.original != nil
}

func (it *singleSectionIterator) Next() (res *AddressSection) {
	if it.HasNext() {
		res = it.original
		it.original = nil
	}
	return
}

type multiSectionIterator struct {
	original        *AddressSection
	iterator        Iterator[[]*AddressDivision]
	valsAreMultiple bool
	prefixLen       PrefixLen
}

func (it *multiSectionIterator) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *multiSectionIterator) Next() (res *AddressSection) {
	if it.HasNext() {
		segs := it.iterator.Next()
		original := it.original
		res = createSection(segs, it.prefixLen, original.addrType)
		res.isMult = it.valsAreMultiple
	}
	return
}

func nilSectIterator() Iterator[*AddressSection] {
	return &singleSectionIterator{}
}

func sectIterator(
	useOriginal bool,
	original *AddressSection,
	valsAreMultiple bool,
	iterator Iterator[[]*AddressDivision],
) Iterator[*AddressSection] {
	if useOriginal {
		return &singleSectionIterator{original: original}
	}
	return &multiSectionIterator{
		original:        original,
		iterator:        iterator,
		valsAreMultiple: valsAreMultiple,
		prefixLen:       original.getPrefixLen(),
	}
}

type prefixSectionIterator struct {
	original   *AddressSection
	iterator   Iterator[[]*AddressDivision]
	isNotFirst bool
	prefixLen  PrefixLen
}

func (it *prefixSectionIterator) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *prefixSectionIterator) Next() (res *AddressSection) {
	if it.HasNext() {
		segs := it.iterator.Next()
		original := it.original
		res = createSection(segs, it.prefixLen, original.addrType)
		if !it.isNotFirst {
			res.initMultiple() // sets isMultiple
			it.isNotFirst = true
		} else if !it.HasNext() {
			res.initMultiple() // sets isMultiple
		} else {
			res.isMult = true
		}
	}
	return
}

func prefixSectIterator(
	useOriginal bool,
	original *AddressSection,
	iterator Iterator[[]*AddressDivision],
) Iterator[*AddressSection] {
	if useOriginal {
		return &singleSectionIterator{original: original}
	}
	return &prefixSectionIterator{
		original:  original,
		iterator:  iterator,
		prefixLen: original.getPrefixLen(),
	}
}

type ipSectionIterator struct {
	Iterator[*AddressSection]
}

func (iter ipSectionIterator) Next() *IPAddressSection {
	return iter.Iterator.Next().ToIP()
}

type ipv4SectionIterator struct {
	Iterator[*AddressSection]
}

func (iter ipv4SectionIterator) Next() *IPv4AddressSection {
	return iter.Iterator.Next().ToIPv4()
}

type ipv6SectionIterator struct {
	Iterator[*AddressSection]
}

func (iter ipv6SectionIterator) Next() *IPv6AddressSection {
	return iter.Iterator.Next().ToIPv6()
}

type macSectionIterator struct {
	Iterator[*AddressSection]
}

func (iter macSectionIterator) Next() *MACAddressSection {
	return iter.Iterator.Next().ToMAC()
}
