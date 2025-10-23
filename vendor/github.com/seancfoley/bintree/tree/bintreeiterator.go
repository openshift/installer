//
// Copyright 2022-2024 Sean C Foley
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

import (
	"container/heap"
)

type HasNext interface {
	HasNext() bool
}

type keyIterator[E Key] interface {
	HasNext

	Next() E

	// Remove removes the last iterated element from the underlying trie, and returns that element.
	// If there is no such element, it returns nil.
	Remove() E
}

type nodeIteratorRem[E Key, V any] interface {
	nodeIterator[E, V]

	// Remove removes the last iterated element from the underlying trie, and returns that element.
	// If there is no such element, it returns nil.
	Remove() *binTreeNode[E, V]
}

type nodeIterator[E Key, V any] interface {
	HasNext

	Next() *binTreeNode[E, V]
}

type binTreeKeyIterator[E Key, V any] struct {
	nodeIteratorRem[E, V]
}

func (iter binTreeKeyIterator[E, V]) Next() E {
	return iter.nodeIteratorRem.Next().GetKey()
}

func (iter binTreeKeyIterator[E, V]) Remove() E {
	return iter.nodeIteratorRem.Remove().GetKey()
}

func newNodeIterator[E Key, V any](forward, addedOnly bool, start, end *binTreeNode[E, V], ctracker *changeTracker) nodeIteratorRem[E, V] {
	var nextOperator func(current *binTreeNode[E, V], end *binTreeNode[E, V]) *binTreeNode[E, V]
	if forward {
		nextOperator = (*binTreeNode[E, V]).nextNodeBounded
	} else {
		nextOperator = (*binTreeNode[E, V]).previousNodeBounded
	}
	if addedOnly {
		wrappedOp := nextOperator
		nextOperator = func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
			return currentNode.nextAdded(endNode, wrappedOp)
		}
	}
	res := binTreeNodeIterator[E, V]{end: end}
	res.initChangeTracker(ctracker)
	res.operator = nextOperator
	res.next = res.getStart(start, end, nil, addedOnly)
	return &res
}

type binTreeNodeIterator[E Key, V any] struct {
	// takes current node and end as args
	operator      func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) (nextNode *binTreeNode[E, V])
	end           *binTreeNode[E, V] // a non-nil node that denotes the end, possibly parent of the starting node
	cTracker      *changeTracker
	currentChange change

	current, next *binTreeNode[E, V]
}

func (iter *binTreeNodeIterator[E, V]) getStart(
	start,
	end *binTreeNode[E, V],
	bounds *bounds[E],
	addedOnly bool) *binTreeNode[E, V] {
	if start == end || start == nil {
		return nil
	}
	if !addedOnly || start.IsAdded() {
		if bounds == nil || bounds.isInBounds(start.GetKey()) {
			return start
		}
	}
	return iter.toNext(start)
}

func (iter *binTreeNodeIterator[E, V]) initChangeTracker(ctracker *changeTracker) {
	if ctracker != nil {
		iter.cTracker, iter.currentChange = ctracker, ctracker.getCurrent()
	}
}

func (iter *binTreeNodeIterator[E, V]) HasNext() bool {
	return iter.next != nil
}

func (iter *binTreeNodeIterator[E, V]) Next() *binTreeNode[E, V] {
	if !iter.HasNext() {
		return nil
	}
	cTracker := iter.cTracker
	if cTracker != nil && cTracker.changedSince(iter.currentChange) {
		changePanic()
	}
	iter.current = iter.next
	iter.next = iter.toNext(iter.next)
	return iter.current
}

func changePanic() {
	panic("the tree has been modified since the iterator was created")
}

func (iter *binTreeNodeIterator[E, V]) toNext(current *binTreeNode[E, V]) *binTreeNode[E, V] {
	return iter.operator(current, iter.end)
}

func (iter *binTreeNodeIterator[E, V]) Remove() *binTreeNode[E, V] {
	if iter.current == nil {
		return nil
	}
	cTracker := iter.cTracker
	if cTracker != nil && cTracker.changedSince(iter.currentChange) {
		changePanic()
	}
	result := iter.current
	result.Remove()
	iter.current = nil
	if cTracker != nil {
		iter.currentChange = cTracker.getCurrent()
	}
	return result
}

type CachingIterator interface {
	// GetCached returns an object previously cached with the current iterated node.
	// After Next has returned a node,
	// if an object was cached by a call to CacheWithLowerSubNode or CacheWithUpperSubNode
	// was called when that node's parent was previously returned by Next,
	// then this returns that cached object.
	GetCached() C

	// CacheWithLowerSubNode caches an object with the lower sub-node of the current iterated node.
	// After Next has returned a node,
	// calling this method caches the provided object with the lower sub-node so that it can
	// be retrieved with GetCached when the lower sub-node is visited later.
	//
	// Returns false if it could not be cached, either because the node has since been removed with a call to Remove,
	// because Next has not been called yet, or because there is no lower sub node for the node previously returned by  Next.
	//
	// The caching and retrieval is done in constant time.
	CacheWithLowerSubNode(C) bool

	// CacheWithUpperSubNode caches an object with the upper sub-node of the current iterated node.
	// After Next has returned a node,
	// calling this method caches the provided object with the upper sub-node so that it can
	// be retrieved with GetCached when the upper sub-node is visited later.
	//
	// Returns false if it could not be cached, either because the node has since been removed with a call to Remove,
	// because Next has not been called yet, or because there is no upper sub node for the node previously returned by Next.
	//
	// The caching and retrieval is done in constant time.
	CacheWithUpperSubNode(C) bool
}

type cachingNodeIterator[E Key, V any] interface {
	nodeIteratorRem[E, V]

	CachingIterator
}

type queueType = any

// see https://pkg.go.dev/container/heap
type nodePriorityQueue struct {
	queue      []queueType
	comparator func(one, two queueType) int // -1, 0 or 1 if one is <, == or > two
}

func (prioQueue nodePriorityQueue) Len() int {
	return len(prioQueue.queue)
}

func (prioQueue nodePriorityQueue) Less(i, j int) bool {
	queue := prioQueue.queue
	return prioQueue.comparator(queue[i], queue[j]) < 0
}

func (prioQueue nodePriorityQueue) Swap(i, j int) {
	queue := prioQueue.queue
	queue[i], queue[j] = queue[j], queue[i]
}

func (prioQueue *nodePriorityQueue) Push(x queueType) {
	prioQueue.queue = append(prioQueue.queue, x)
}

func (prioQueue *nodePriorityQueue) Pop() queueType {
	current := prioQueue.queue
	queueLen := len(current)
	topNode := current[queueLen-1]
	current[queueLen-1] = nil
	prioQueue.queue = current[:queueLen-1]
	return topNode
}

func newPriorityNodeIterator[E Key, V any](
	treeSize int,
	addedOnly bool,
	start *binTreeNode[E, V],
	comparator func(E, E) int,
) binTreeNodeIterator[E, V] {
	return newPriorityNodeIteratorBounded(
		nil,
		treeSize,
		addedOnly,
		start,
		comparator)
}

func newPriorityNodeIteratorBounded[E Key, V any](
	bnds *bounds[E],
	treeSize int,
	addedOnly bool,
	start *binTreeNode[E, V],
	comparator func(E, E) int) binTreeNodeIterator[E, V] {

	comp := func(one, two queueType) int {
		node1, node2 := one.(*binTreeNode[E, V]), two.(*binTreeNode[E, V])
		addr1, addr2 := node1.GetKey(), node2.GetKey()
		return comparator(addr1, addr2)
	}
	queue := &nodePriorityQueue{comparator: comp}
	if treeSize > 0 {
		queue.queue = make([]queueType, 0, (treeSize+2)>>1)
	}
	op := func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
		lower := currentNode.getLowerSubNode()
		if lower != nil {
			heap.Push(queue, lower)
		}
		upper := currentNode.getUpperSubNode()
		if upper != nil {
			heap.Push(queue, upper)
		}
		var node *binTreeNode[E, V]
		if queue.Len() > 0 {
			node = heap.Pop(queue).(*binTreeNode[E, V])
		}
		if node == endNode {
			return nil
		}
		return node
	}
	if addedOnly {
		wrappedOp := op
		op = func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
			return currentNode.nextAdded(endNode, wrappedOp)
		}
	}
	if bnds != nil {
		wrappedOp := op
		op = func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
			return currentNode.nextInBounds(endNode, wrappedOp, bnds)
		}
	}
	res := binTreeNodeIterator[E, V]{operator: op}
	start = res.getStart(start, nil, bnds, addedOnly)
	if start != nil {
		res.next = start
		res.initChangeTracker(start.cTracker)
	}
	return res
}

func newCachingPriorityNodeIterator[E Key, V any](
	start *binTreeNode[E, V],
	comparator func(E, E) int,
) cachingPriorityNodeIterator[E, V] {
	return newCachingPriorityNodeIteratorSized(
		0,
		start,
		comparator)
}

func newCachingPriorityNodeIteratorSized[E Key, V any](
	treeSize int,
	start *binTreeNode[E, V],
	comparator func(E, E) int) cachingPriorityNodeIterator[E, V] {

	comp := func(one, two queueType) int {
		cached1, cached2 := one.(*cached[E, V]), two.(*cached[E, V])
		node1, node2 := cached1.node, cached2.node
		addr1, addr2 := node1.GetKey(), node2.GetKey()
		return comparator(addr1, addr2)
	}
	queue := &nodePriorityQueue{comparator: comp}
	if treeSize > 0 {
		queue.queue = make([]queueType, 0, (treeSize+2)>>1)
	}
	res := cachingPriorityNodeIterator[E, V]{cached: &cachedObjs[E, V]{}}
	res.operator = res.getNextOperation(queue)
	start = res.getStart(start, nil, nil, false)
	if start != nil {
		res.next = start
		res.initChangeTracker(start.cTracker)
	}
	return res
}

type cachedObjs[E Key, V any] struct {
	cacheItem                    C
	nextCachedItem               *cached[E, V]
	lowerCacheObj, upperCacheObj *cached[E, V]
}

type cachingPriorityNodeIterator[E Key, V any] struct {
	binTreeNodeIterator[E, V]
	cached *cachedObjs[E, V]
}

func (iter *cachingPriorityNodeIterator[E, V]) getNextOperation(queue *nodePriorityQueue) func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
	return func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
		lower := currentNode.getLowerSubNode()
		cacheObjs := iter.cached
		if lower != nil {
			cachd := &cached[E, V]{
				node: lower,
			}
			cacheObjs.lowerCacheObj = cachd
			heap.Push(queue, cachd)
		} else {
			cacheObjs.lowerCacheObj = nil
		}
		upper := currentNode.getUpperSubNode()
		if upper != nil {
			cachd := &cached[E, V]{
				node: upper,
			}
			cacheObjs.upperCacheObj = cachd
			heap.Push(queue, cachd)
		} else {
			cacheObjs.upperCacheObj = nil
		}
		if cacheObjs.nextCachedItem != nil {
			cacheObjs.cacheItem = cacheObjs.nextCachedItem.cached
		}
		var item queueType
		if queue.Len() > 0 {
			item = heap.Pop(queue)
		}
		if item != nil {
			cachd := item.(*cached[E, V])
			node := cachd.node
			if node != endNode {
				cacheObjs.nextCachedItem = cachd
				return node
			}
		}
		cacheObjs.nextCachedItem = nil
		return nil
	}
}

func (iter *cachingPriorityNodeIterator[E, V]) GetCached() C {
	return iter.cached.cacheItem
}

func (iter *cachingPriorityNodeIterator[E, V]) CacheWithLowerSubNode(object C) bool {
	cached := iter.cached
	if cached.lowerCacheObj != nil {
		cached.lowerCacheObj.cached = object
		return true
	}
	return false
}

func (iter *cachingPriorityNodeIterator[E, V]) CacheWithUpperSubNode(object C) bool {
	cached := iter.cached
	if cached.upperCacheObj != nil {
		cached.upperCacheObj.cached = object
		return true
	}
	return false
}

type cached[E Key, V any] struct {
	node   *binTreeNode[E, V]
	cached C
}

// The caching only useful when in reverse order, since you have to visit parent nodes first for it to be useful.
func newPostOrderNodeIterator[E Key, V any](
	forward, addedOnly bool,
	start, end *binTreeNode[E, V],
	ctracker *changeTracker,
) subNodeCachingIterator[E, V] {
	return newPostOrderNodeIteratorBounded(
		nil,
		forward, addedOnly,
		start, end,
		ctracker)
}

func newPostOrderNodeIteratorBounded[E Key, V any](
	bnds *bounds[E],
	forward, addedOnly bool,
	start, end *binTreeNode[E, V],
	ctracker *changeTracker) subNodeCachingIterator[E, V] {
	var op func(current *binTreeNode[E, V], end *binTreeNode[E, V]) *binTreeNode[E, V]
	if forward {
		op = (*binTreeNode[E, V]).nextPostOrderNode
	} else {
		op = (*(binTreeNode[E, V])).previousPostOrderNode
	}
	// do the added-only filter first, because it is simpler
	if addedOnly {
		wrappedOp := op
		op = func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
			return currentNode.nextAdded(endNode, wrappedOp)
		}
	}
	if bnds != nil {
		wrappedOp := op
		op = func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
			return currentNode.nextInBounds(endNode, wrappedOp, bnds)
		}
	}
	return newSubNodeCachingIterator[E, V](
		bnds,
		forward, addedOnly,
		start, end,
		ctracker,
		op,
		!forward,
		!forward || addedOnly)
}

// The caching only useful when in forward order, since you have to visit parent nodes first for it to be useful.
func newPreOrderNodeIterator[E Key, V any](
	forward, addedOnly bool,
	start, end *binTreeNode[E, V],
	ctracker *changeTracker) subNodeCachingIterator[E, V] {
	return newPreOrderNodeIteratorBounded(
		nil,
		forward, addedOnly,
		start, end,
		ctracker)
}

func newPreOrderNodeIteratorBounded[E Key, V any](
	bnds *bounds[E],
	forward, addedOnly bool,
	start, end *binTreeNode[E, V],
	ctracker *changeTracker) subNodeCachingIterator[E, V] {
	var op func(current *binTreeNode[E, V], end *binTreeNode[E, V]) *binTreeNode[E, V]
	if forward {
		op = (*binTreeNode[E, V]).nextPreOrderNode
	} else {
		op = (*binTreeNode[E, V]).previousPreOrderNode
	}
	// do the added-only filter first, because it is simpler
	if addedOnly {
		wrappedOp := op
		op = func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
			return currentNode.nextAdded(endNode, wrappedOp)
		}
	}
	if bnds != nil {
		wrappedOp := op
		op = func(currentNode *binTreeNode[E, V], endNode *binTreeNode[E, V]) *binTreeNode[E, V] {
			return currentNode.nextInBounds(endNode, wrappedOp, bnds)
		}
	}
	return newSubNodeCachingIterator(
		bnds,
		forward, addedOnly,
		start, end,
		ctracker,
		op,
		forward,
		forward || addedOnly)
}

func newSubNodeCachingIterator[E Key, V any](
	bnds *bounds[E],
	forward, addedOnly bool,
	start, end *binTreeNode[E, V],
	ctracker *changeTracker,
	nextOperator func(current *binTreeNode[E, V], end *binTreeNode[E, V]) *binTreeNode[E, V],
	allowCaching,
	allowRemove bool,
) subNodeCachingIterator[E, V] {
	res := subNodeCachingIterator[E, V]{
		allowCaching:        allowCaching,
		allowRemove:         allowRemove,
		stackIndex:          -1,
		bnds:                bnds,
		isForward:           forward,
		addedOnly:           addedOnly,
		binTreeNodeIterator: binTreeNodeIterator[E, V]{end: end},
	}
	res.initChangeTracker(ctracker)
	res.operator = nextOperator
	res.next = res.getStart(start, end, bnds, addedOnly)
	return res
}

const ipv6BitCount = 128
const stackSize = ipv6BitCount + 2 // 129 for prefixes /0 to /128 and also 1 more for non-prefixed

type subNodeCachingIterator[E Key, V any] struct {
	binTreeNodeIterator[E, V]

	cacheItem  C
	nextKey    E
	nextCached C
	stack      []C
	stackIndex int

	bnds                 *bounds[E]
	addedOnly, isForward bool

	// Both these fields are not really necessary because
	// the caching and removal functionality should not be exposed when it is not usable.
	// The interfaces will not include the caching and Remove() methods in the cases where they are not usable.
	// So these fields are both runtime checks for coding errors.
	allowCaching, allowRemove bool
}

func (iter *subNodeCachingIterator[E, V]) Next() *binTreeNode[E, V] {
	result := iter.binTreeNodeIterator.Next()
	if result != nil && iter.allowCaching {
		iter.populateCacheItem(result)
	}
	return result
}

func (iter *subNodeCachingIterator[E, V]) GetCached() C {
	if !iter.allowCaching {
		panic("no caching allowed, this code path should not be accessible")
	}
	return iter.cacheItem
}

func (iter *subNodeCachingIterator[E, V]) populateCacheItem(current *binTreeNode[E, V]) {
	nextKey := iter.nextKey
	if current.GetKey() == nextKey {
		iter.cacheItem = iter.nextCached
		iter.nextCached = nil
	} else {
		stack := iter.stack
		if stack != nil {
			stackIndex := iter.stackIndex
			if stackIndex >= 0 && stack[stackIndex] == current.GetKey() {
				iter.cacheItem = stack[stackIndex+stackSize].(C)
				stack[stackIndex+stackSize] = nil
				stack[stackIndex] = nil
				iter.stackIndex--
			} else {
				iter.cacheItem = nil
			}
		} else {
			iter.cacheItem = nil
		}
	}
}

func (iter *subNodeCachingIterator[E, V]) Remove() *binTreeNode[E, V] {
	if !iter.allowRemove {
		// Example:
		// Suppose we are at right sub-node, just visited left.  Next node is parent, but not added.
		// When right is removed, so is the parent, so that the left takes its place.
		// But parent is our next node.  Now our next node is invalid.  So we are lost.
		// This is avoided for iterators that are "added" only.
		panic("no removal allowed, this code path should not be accessible")
	}
	return iter.binTreeNodeIterator.Remove()
}

func (iter *subNodeCachingIterator[E, V]) checkCaching() {
	if !iter.allowCaching {
		panic("no caching allowed, this code path should not be accessible")
	}
}

func (iter *subNodeCachingIterator[E, V]) CacheWithLowerSubNode(object C) bool {
	iter.checkCaching()
	if iter.isForward {
		return iter.cacheWithFirstSubNode(object)
	}
	return iter.cacheWithSecondSubNode(object)

}

func (iter *subNodeCachingIterator[E, V]) CacheWithUpperSubNode(object C) bool {
	iter.checkCaching()
	if iter.isForward {
		return iter.cacheWithSecondSubNode(object)
	}
	return iter.cacheWithFirstSubNode(object)
}

// the sub-node will be the next visited node
func (iter *subNodeCachingIterator[E, V]) cacheWithFirstSubNode(object C) bool {
	iter.checkCaching()
	if iter.current != nil {
		var firstNode *binTreeNode[E, V]
		if iter.isForward {
			firstNode = iter.current.getLowerSubNode()
		} else {
			firstNode = iter.current.getUpperSubNode()
		}
		if firstNode != nil {
			if (iter.addedOnly && !firstNode.IsAdded()) || (iter.bnds != nil && !iter.bnds.isInBounds(firstNode.GetKey())) {
				firstNode = iter.operator(firstNode, iter.current)
			}
			if firstNode != nil {
				// the lower sub-node is always next if it exists
				iter.nextKey = firstNode.GetKey()
				iter.nextCached = object
				return true
			}
		}
	}
	return false
}

// the sub-node will only be the next visited node if there is no other sub-node,
// otherwise it might not be visited for a while
func (iter *subNodeCachingIterator[E, V]) cacheWithSecondSubNode(object C) bool {
	iter.checkCaching()
	if iter.current != nil {
		var secondNode *binTreeNode[E, V]
		if iter.isForward {
			secondNode = iter.current.getUpperSubNode()
		} else {
			secondNode = iter.current.getLowerSubNode()
		}
		if secondNode != nil {
			if (iter.addedOnly && !secondNode.IsAdded()) || (iter.bnds != nil && !iter.bnds.isInBounds(secondNode.GetKey())) {
				secondNode = iter.operator(secondNode, iter.current)
			}
			if secondNode != nil {
				// if there is no lower node, we can use the nextCached field since upper is next when no lower sub-node
				var firstNode *binTreeNode[E, V]
				if iter.isForward {
					firstNode = iter.current.getLowerSubNode()
				} else {
					firstNode = iter.current.getUpperSubNode()
				}
				if firstNode == nil {
					iter.nextKey = secondNode.GetKey()
					iter.nextCached = object
				} else {
					if iter.stack == nil {
						iter.stack = make([]C, stackSize<<1)
					}
					iter.stackIndex++
					iter.stack[iter.stackIndex] = secondNode.GetKey()
					iter.stack[iter.stackIndex+stackSize] = object
				}
				return true
			}
		}
	}
	return false
}
