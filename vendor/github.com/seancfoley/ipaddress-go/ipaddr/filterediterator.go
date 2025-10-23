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

type filteredAddrIterator struct {
	skip func(*Address) bool
	iter Iterator[*Address]
	next *Address
}

func (it *filteredAddrIterator) Next() (res *Address) {
	res = it.next
	for {
		next := it.iter.Next()
		if next == nil || !it.skip(next) {
			it.next = next
			break
		}
	}
	return res
}

func (it *filteredAddrIterator) HasNext() bool {
	return it.next != nil
}

// NewFilteredAddrIterator modifies an address iterator to skip certain addresses,
// skipping those elements for which the "skip" function returns true
func NewFilteredAddrIterator(iter Iterator[*Address], skip func(*Address) bool) Iterator[*Address] {
	res := &filteredAddrIterator{skip: skip, iter: iter}
	res.Next()
	return res
}

type filteredIPAddrIterator struct {
	skip func(*IPAddress) bool
	iter Iterator[*IPAddress]
	next *IPAddress
}

func (it *filteredIPAddrIterator) Next() (res *IPAddress) {
	res = it.next
	for {
		next := it.iter.Next()
		if next == nil || !it.skip(next) {
			it.next = next
			break
		}
	}
	return res
}

func (it *filteredIPAddrIterator) HasNext() bool {
	return it.next != nil
}

// NewFilteredIPAddrIterator returns an iterator similar to the passed in iterator,
// but skipping those elements for which the "skip" function returns true
func NewFilteredIPAddrIterator(iter Iterator[*IPAddress], skip func(*IPAddress) bool) Iterator[*IPAddress] {
	res := &filteredIPAddrIterator{skip: skip, iter: iter}
	res.Next()
	return res
}
