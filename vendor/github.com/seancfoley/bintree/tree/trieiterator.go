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

package tree

type TrieKeyIterator[E TrieKey[E]] interface {
	HasNext

	Next() E

	// Remove removes the last iterated element from the underlying trie, and returns that element.
	// If there is no such element, it returns the zero value.
	Remove() E
}

type trieKeyIterator[E TrieKey[E]] struct {
	keyIterator[E]
}

func (iter trieKeyIterator[E]) Next() E {
	return iter.keyIterator.Next()
}

func (iter trieKeyIterator[E]) Remove() E {
	return iter.keyIterator.Remove()
}

type TrieNodeIterator[E TrieKey[E], V any] interface {
	HasNext

	Next() *BinTrieNode[E, V]
}

type TrieNodeIteratorRem[E TrieKey[E], V any] interface {
	TrieNodeIterator[E, V]

	// Remove removes the last iterated element from the underlying trie, and returns that element.
	// If there is no such element, it returns the zero value.
	Remove() *BinTrieNode[E, V]
}

type trieNodeIteratorRem[E TrieKey[E], V any] struct {
	nodeIteratorRem[E, V]
}

func (iter trieNodeIteratorRem[E, V]) Next() *BinTrieNode[E, V] {
	return toTrieNode(iter.nodeIteratorRem.Next())
}

func (iter trieNodeIteratorRem[E, V]) Remove() *BinTrieNode[E, V] {
	return toTrieNode(iter.nodeIteratorRem.Remove())
}

type trieNodeIterator[E TrieKey[E], V any] struct {
	nodeIterator[E, V]
}

func (iter trieNodeIterator[E, V]) Next() *BinTrieNode[E, V] {
	return toTrieNode(iter.nodeIterator.Next())
}

type CachingTrieNodeIterator[E TrieKey[E], V any] interface {
	TrieNodeIteratorRem[E, V]
	CachingIterator
}

type cachingTrieNodeIterator[E TrieKey[E], V any] struct {
	cachingNodeIterator[E, V] // an interface
}

func (iter *cachingTrieNodeIterator[E, V]) Next() *BinTrieNode[E, V] {
	return toTrieNode(iter.cachingNodeIterator.Next())
}

func (iter *cachingTrieNodeIterator[E, V]) Remove() *BinTrieNode[E, V] {
	return toTrieNode(iter.cachingNodeIterator.Remove())
}

type dualTrieNodeIterator struct {
	oneCurrentChange, twoCurrentChange change
	oneChangeTracker, twoChangeTracker *changeTracker
}

type dualIterator[E TrieKey[E], V any] struct {
	dualTrieNodeIterator

	one, two TrieNodeIteratorRem[E, V]

	onSecond bool
}

func (iter *dualIterator[E, V]) HasNext() bool {
	if iter.onSecond {
		return iter.two.HasNext()
	}
	return iter.one.HasNext() || iter.two.HasNext()
}

func (iter *dualIterator[E, V]) Next() *BinTrieNode[E, V] {
	if iter.onSecond {
		ct := iter.oneChangeTracker
		if ct != nil && ct.changedSince(iter.oneCurrentChange) {
			changePanic()
		}
		return iter.two.Next()
	} else if !iter.one.HasNext() {
		iter.onSecond = true
		ct := iter.oneChangeTracker
		if ct != nil && ct.changedSince(iter.oneCurrentChange) {
			changePanic()
		}
		return iter.two.Next()
	}
	ct := iter.twoChangeTracker
	if ct != nil && ct.changedSince(iter.twoCurrentChange) {
		changePanic()
	}
	return iter.one.Next()
}

func (iter *dualIterator[E, V]) Remove() (res *BinTrieNode[E, V]) {
	if iter.onSecond {
		// we check the change tracker of the iterator we are not using,
		// since the call to Remove will check the change tracker of the iterator we are using
		ct := iter.oneChangeTracker
		if ct != nil && ct.changedSince(iter.oneCurrentChange) {
			changePanic()
		}
		res = iter.two.Remove()
		// we update the change value of the iterator we used, since we just made a change
		ct = iter.twoChangeTracker
		if ct != nil {
			iter.twoCurrentChange = ct.getCurrent()
		}
	} else {
		// we check the change tracker of the iterator we are not using,
		// since the call to Remove will check the change tracker of the iterator we are using
		ct := iter.twoChangeTracker
		if ct != nil && ct.changedSince(iter.twoCurrentChange) {
			changePanic()
		}
		res = iter.one.Remove()
		// we update the change value of the iterator we used, since we just made a change
		ct = iter.oneChangeTracker
		if ct != nil {
			iter.oneCurrentChange = ct.getCurrent()
		}
	}
	return
}

// CombineSequentiallyRem will combine iterators from two tries into a single iterator.
// Each trie must have a non-nil root, which is provided as an argument to this function.
// The root is needed to monitor for changes to either trie during iteration.
func CombineSequentiallyRem[E TrieKey[E], V any](oneRoot, twoRoot *BinTrieNode[E, V], one, two TrieNodeIteratorRem[E, V], forward bool) TrieNodeIteratorRem[E, V] {
	if !forward {
		t := one
		one = two
		two = t
		r := oneRoot
		oneRoot = twoRoot
		twoRoot = r
	}

	oneTracker := oneRoot.cTracker
	var oneChange, twoChange change
	if oneTracker != nil {
		oneChange = oneTracker.getCurrent()
	}
	twoTracker := twoRoot.cTracker
	if twoTracker != nil {
		twoChange = twoTracker.getCurrent()
	}

	return &dualIterator[E, V]{
		dualTrieNodeIterator: dualTrieNodeIterator{
			oneCurrentChange: oneChange,
			twoCurrentChange: twoChange,
			oneChangeTracker: oneTracker,
			twoChangeTracker: twoTracker,
		},
		one: one,
		two: two,
	}
}

func CombineByBlockSize[E TrieKey[E], V any](oneRoot, twoRoot *BinTrieNode[E, V], one, two TrieNodeIteratorRem[E, V], lowerSubNodeFirst bool) TrieNodeIteratorRem[E, V] {
	oneTracker := oneRoot.cTracker
	var oneChange, twoChange change
	if oneTracker != nil {
		oneChange = oneTracker.getCurrent()
	}
	twoTracker := twoRoot.cTracker
	if twoTracker != nil {
		twoChange = twoTracker.getCurrent()
	}

	return &dualBlockSizeIterator[E, V]{
		dualTrieNodeIterator: dualTrieNodeIterator{
			oneCurrentChange: oneChange,
			twoCurrentChange: twoChange,
			oneChangeTracker: oneTracker,
			twoChangeTracker: twoTracker,
		},
		one:               one,
		two:               two,
		lowerSubNodeFirst: lowerSubNodeFirst,
	}
}

type dualBlockSizeIterator[E TrieKey[E], V any] struct {
	dualTrieNodeIterator

	oneItem, twoItem, lastItem *BinTrieNode[E, V]

	one, two TrieNodeIteratorRem[E, V]

	lastItemIsOne, lowerSubNodeFirst bool
}

func (iter *dualBlockSizeIterator[E, V]) HasNext() bool {
	return iter.oneItem != nil || iter.twoItem != nil || iter.one.HasNext() || iter.two.HasNext()
}

func (iter *dualBlockSizeIterator[E, V]) Next() (result *BinTrieNode[E, V]) {
	if !iter.HasNext() {
		return
	}

	var accessedOne, accessedTwo bool

	// replace whatever was returned previously
	if iter.oneItem == nil && iter.one.HasNext() {
		accessedOne = true
		iter.oneItem = iter.one.Next()
	}
	if iter.twoItem == nil && iter.two.HasNext() {
		accessedTwo = true
		iter.twoItem = iter.two.Next()
	}
	if !accessedOne {
		ct := iter.oneChangeTracker
		if ct != nil && ct.changedSince(iter.oneCurrentChange) {
			changePanic()
		}
	}
	if !accessedTwo {
		ct := iter.twoChangeTracker
		if ct != nil && ct.changedSince(iter.twoCurrentChange) {
			changePanic()
		}
	}

	// now return the lowest of the two
	if iter.oneItem == nil {
		result = iter.twoItem
		iter.twoItem = nil
		iter.lastItemIsOne = false
	} else if iter.twoItem == nil {
		result = iter.oneItem
		iter.oneItem = nil
		iter.lastItemIsOne = true
	} else {
		cmp := blockSizeCompare(iter.oneItem.GetKey(), iter.twoItem.GetKey(), !iter.lowerSubNodeFirst)
		if cmp < 0 {
			result = iter.oneItem
			iter.oneItem = nil
			iter.lastItemIsOne = true
		} else {
			result = iter.twoItem
			iter.twoItem = nil
			iter.lastItemIsOne = false
		}
	}
	iter.lastItem = result
	return
}

// blockSizeCompare compares keys by block size and then by prefix value if block sizes are equal
func blockSizeCompare[E TrieKey[E]](key1, key2 E, reverseBlocksEqualSize bool) int {
	if key2 == key1 {
		return 0
	}
	pref2 := key2.GetPrefixLen()
	pref1 := key1.GetPrefixLen()
	if pref1 != nil {
		if pref2 != nil {
			val := (key2.GetBitCount() - pref2.Len()) - (key1.GetBitCount() - pref1.Len())
			if val == 0 {
				compVal := key2.Compare(key1)
				if reverseBlocksEqualSize {
					compVal = -compVal
				}
				return compVal
			}
			return val
		}
		return -1
	}
	if pref2 != nil {
		return 1
	}
	compVal := key2.Compare(key1)
	if reverseBlocksEqualSize {
		compVal = -compVal
	}
	return compVal
}

func (iter *dualBlockSizeIterator[E, V]) Remove() (result *BinTrieNode[E, V]) {
	if iter.lastItem == nil {
		return
	}
	if iter.lastItemIsOne {
		ct := iter.twoChangeTracker
		if ct != nil && ct.changedSince(iter.twoCurrentChange) {
			changePanic()
		}
		result = iter.one.Remove()
		// we update the change value of the iterator we used, since we just made a change
		ct = iter.oneChangeTracker
		if ct != nil {
			iter.oneCurrentChange = ct.getCurrent()
		}
	} else {
		ct := iter.oneChangeTracker
		if ct != nil && ct.changedSince(iter.oneCurrentChange) {
			changePanic()
		}
		result = iter.two.Remove()
		ct = iter.twoChangeTracker
		if ct != nil {
			iter.twoCurrentChange = ct.getCurrent()
		}
	}
	iter.lastItem = nil
	return
}
