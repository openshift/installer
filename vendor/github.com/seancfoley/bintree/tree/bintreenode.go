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
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

var freezeRoot = true

func bigOne() *big.Int {
	return big.NewInt(1)
}

var one = bigOne()

type change struct {
	big   *big.Int
	small uint64
}

func (c change) Equal(c2 change) bool {
	if c.small == c2.small {
		if c.big == nil {
			return c2.big == nil
		} else if c2.big != nil {
			return c.big.Cmp(c2.big) == 0
		}
	}
	return false
}

func (c *change) increment() {
	val := c.small
	val++
	if val == 0 {
		if c.big == nil {
			c.big = bigOne()
		} else {
			c.big.Add(c.big, one)
		}
	}
	c.small = val
}

func (c change) String() string {
	return c.big.String() + " " + strconv.FormatUint(c.small, 10)
}

type changeTracker struct {
	currentChange change
	watched       bool
}

func (c *changeTracker) changed() {
	if c.watched {
		c.watched = false
		c.currentChange.increment()
	} // else nobody is watching the current change, so no need to do anything
}

func (c *changeTracker) changedSince(otherChange change) bool {
	return !c.currentChange.Equal(otherChange)
}

func (c *changeTracker) getCurrent() change {
	c.watched = true
	return c.currentChange
}

func (c *changeTracker) String() string {
	return "current change: " + c.currentChange.String()
}

type bounds[E Key] struct {
}

func (b *bounds[E]) isInBounds(_ E) bool {
	return true
}

func (b *bounds[E]) isWithinLowerBound(_ E) bool {
	return true
}

func (b *bounds[E]) isBelowLowerBound(_ E) bool {
	return true
}

func (b *bounds[E]) isWithinUpperBound(_ E) bool {
	return true
}

func (b *bounds[E]) isAboveUpperBound(_ E) bool {
	return true
}

type Key interface {
	comparable // needed by populateCacheItem
}

// C represents cached values in iterators
type C any

const sizeUnknown = -1

type binTreeNode[E Key, V any] struct {

	// the key for the node
	item E

	// only for associative trie nodes
	value V

	parent, lower, upper *binTreeNode[E, V]

	storedSize int

	cTracker *changeTracker

	// used to store opResult objects for search operations
	pool *sync.Pool

	// some nodes represent elements added to the tree and others are nodes generated internally when other nodes are added
	added bool

	self *binTreeNode[E, V]
}

// This hideptr trick is used in strings.Builder to trick escape analysis to ensure that this self-referential pointer does not cause automatic heap allocation
// cannot hurt to use it
// https://github.com/golang/go/issues/23382
// https://github.com/golang/go/issues/7921
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.6:src/strings/builder.go;l=28

//go:nosplit
//go:nocheckptr
func hideptr(p unsafe.Pointer) unsafe.Pointer {
	ptr := uintptr(p)
	return unsafe.Pointer(ptr ^ 0)
}

func (node *binTreeNode[E, V]) setAddr() {
	node.self = (*binTreeNode[E, V])(hideptr(unsafe.Pointer(node)))
}

func (node *binTreeNode[E, V]) checkCopy() {
	if node != nil && node.self != nil && node.self != node {
		panic("attempting to modify trie with a copied node")
	}
}

func (node *binTreeNode[E, V]) getChangeTracker() *changeTracker {
	if node == nil {
		return nil
	}
	return node.cTracker
}

func toTrieNode[E TrieKey[E], V any](node *binTreeNode[E, V]) *BinTrieNode[E, V] {
	return (*BinTrieNode[E, V])(unsafe.Pointer(node))
}

// when FREEZE_ROOT is true, this is never called (and FREEZE_ROOT is always true)
func (node *binTreeNode[E, V]) setKey(item E) {
	node.item = item
}

// Gets the key used for placing the node in the tree.
func (node *binTreeNode[E, V]) GetKey() (key E) {
	if node != nil {
		key = node.item
	}
	return
}

// SetValue assigns a value to the node, overwriting any previous value
func (node *binTreeNode[E, V]) SetValue(val V) {
	// new value assignment
	node.value = val
}

// GetValue returns the value assigned to the node
func (node *binTreeNode[E, V]) GetValue() (val V) {
	if node != nil {
		val = node.value
	}
	return
}

func (node *binTreeNode[E, V]) ClearValue() {
	var v V
	// new value assignment
	node.value = v
}

// IsRoot returns whether this is the root of the backing tree.
func (node *binTreeNode[E, V]) IsRoot() bool {
	return node != nil && node.parent == nil
}

// Gets the node from which this node is a direct child node, or nil if this is the root.
func (node *binTreeNode[E, V]) getParent() (parent *binTreeNode[E, V]) {
	if node != nil {
		parent = node.parent
	}
	return
}

func (node *binTreeNode[E, V]) setParent(parent *binTreeNode[E, V]) {
	node.parent = parent
}

// Gets the direct child node whose key is largest in value
func (node *binTreeNode[E, V]) getUpperSubNode() (upper *binTreeNode[E, V]) {
	if node != nil {
		upper = node.upper
	}
	return
}

// Gets the direct child node whose key is smallest in value
func (node *binTreeNode[E, V]) getLowerSubNode() (lower *binTreeNode[E, V]) {
	if node != nil {
		lower = node.lower
	}
	return
}

func (node *binTreeNode[E, V]) setUpper(upper *binTreeNode[E, V]) {
	node.upper = upper
	if upper != nil {
		upper.setParent(node)
	}
}

func (node *binTreeNode[E, V]) setLower(lower *binTreeNode[E, V]) {
	node.lower = lower
	if lower != nil {
		lower.setParent(node)
	}
}

// IsAdded returns whether the node was "added".
// Some binary tree nodes are considered "added" and others are not.
// Those nodes created for key elements added to the tree are "added" nodes.
// Those that are not added are those nodes created to serve as junctions for the added nodes.
// Only added elements contribute to the size of a tree.
// When removing nodes, non-added nodes are removed automatically whenever they are no longer needed,
// which is when an added node has less than two added sub-nodes.
func (node *binTreeNode[E, V]) IsAdded() bool {
	return node != nil && node.added
}

// SetAdded makes this node an added node, which is equivalent to adding the corresponding key to the tree.
// If the node is already an added node, this method has no effect.
// You cannot set an added node to non-added, for that you should Remove the node from the tree by calling Remove.
// A non-added node will only remain in the tree if it needs to in the tree.
func (node *binTreeNode[E, V]) SetAdded() {
	if !node.added {
		node.setNodeAdded(true)
		node.adjustCount(1)
	}
}

func (node *binTreeNode[E, V]) setNodeAdded(added bool) {
	node.added = added
}

// Size returns the count of nodes added to the sub-tree starting from this node as root and moving downwards to sub-nodes.
// This is a constant-time operation since the size is maintained in each node and adjusted with each add and Remove operation in the sub-tree.
func (node *binTreeNode[E, V]) Size() (storedSize int) {
	if node != nil {
		storedSize = node.storedSize
		if storedSize == sizeUnknown {
			iterator := node.containedFirstAllNodeIterator(true)
			for next := iterator.Next(); next != nil; next = iterator.Next() {
				var nodeSize int
				if next.IsAdded() {
					nodeSize = 1
				}
				lower := next.getLowerSubNode()
				if lower != nil {
					nodeSize += lower.storedSize
				}
				upper := next.getUpperSubNode()
				if upper != nil {
					nodeSize += upper.storedSize
				}
				next.storedSize = nodeSize
			}
			storedSize = node.storedSize
		}
	}
	return
}

// NodeSize returns the count of all nodes in the tree starting from this node and extending to all sub-nodes.
// Unlike for the Size method, this is not a constant-time operation and must visit all sub-nodes of this node.
func (node *binTreeNode[E, V]) NodeSize() int {
	totalCount := 0
	iterator := node.allNodeIterator(false)
	next := iterator.Next()
	for next != nil {
		totalCount++
		next = iterator.Next()
	}
	return totalCount
}

func (node *binTreeNode[E, V]) adjustCount(delta int) {
	if delta != 0 {
		thisNode := node
		for {
			thisNode.storedSize += delta
			thisNode = thisNode.getParent()
			if thisNode == nil {
				break
			}
		}
	}
}

// Remove removes this node from the collection of added nodes,
// and also removes from the tree if possible.
// If it has two sub-nodes, it cannot be removed from the tree, in which case it is marked as not "added",
// nor is it counted in the tree size.
// Only added nodes can be removed from the tree.  If this node is not added, this method does nothing.
func (node *binTreeNode[E, V]) Remove() {
	node.checkCopy()
	if !node.IsAdded() {
		return
	} else if freezeRoot && node.IsRoot() {
		node.removed()
	} else if node.getUpperSubNode() == nil {
		node.replaceThis(node.getLowerSubNode()) // also handles case of lower == nil
	} else if node.getLowerSubNode() == nil {
		node.replaceThis(node.getUpperSubNode())
	} else { // has two sub-nodes
		node.removed()
	}
}

func (node *binTreeNode[E, V]) removed() {
	node.adjustCount(-1)
	node.setNodeAdded(false)
	node.cTracker.changed()
	node.ClearValue()
}

// Makes the parent of this point to something else, thus removing this and all sub-nodes from the tree
func (node *binTreeNode[E, V]) replaceThis(replacement *binTreeNode[E, V]) {
	node.replaceThisRecursive(replacement, 0)
	node.cTracker.changed()
}

func (node *binTreeNode[E, V]) replaceThisRecursive(replacement *binTreeNode[E, V], additionalSizeAdjustment int) {
	if node.IsRoot() {
		node.replaceThisRoot(replacement)
		return
	}
	parent := node.getParent()
	if parent.getUpperSubNode() == node {
		// we adjust parents first, using the size and other characteristics of ourselves,
		// before the parent severs the link to ourselves with the call to setUpper,
		// since the setUpper call is allowed to change the characteristics of the child,
		// and in some cases this does adjust the size of the child.
		node.adjustTree(parent, replacement, additionalSizeAdjustment, true)
		parent.setUpper(replacement)
	} else if parent.getLowerSubNode() == node {
		node.adjustTree(parent, replacement, additionalSizeAdjustment, false)
		parent.setLower(replacement)
	} else {
		panic("corrupted trie") // will never reach here
	}
}

func (node *binTreeNode[E, V]) adjustTree(parent, replacement *binTreeNode[E, V], additionalSizeAdjustment int, replacedUpper bool) {
	sizeAdjustment := -node.storedSize
	if replacement == nil {
		if !parent.IsAdded() && (!freezeRoot || !parent.IsRoot()) {
			parent.storedSize += sizeAdjustment
			var parentReplacement *binTreeNode[E, V]
			if replacedUpper {
				parentReplacement = parent.getLowerSubNode()
			} else {
				parentReplacement = parent.getUpperSubNode()
			}
			parent.replaceThisRecursive(parentReplacement, sizeAdjustment)
		} else {
			parent.adjustCount(sizeAdjustment + additionalSizeAdjustment)
		}
	} else {
		parent.adjustCount(replacement.storedSize + sizeAdjustment + additionalSizeAdjustment)
	}
	node.setParent(nil)
}

func (node *binTreeNode[E, V]) replaceThisRoot(replacement *binTreeNode[E, V]) {
	if replacement == nil {
		node.setNodeAdded(false)
		node.setUpper(nil)
		node.setLower(nil)
		if !freezeRoot {
			var e E
			node.setKey(e)
			//node.setKey(nil)
			// here we'd need to replace with the default root (ie call setKey with key of 0.0.0.0/0 or ::/0 or 0:0:0:0:0:0)
		}
		node.storedSize = 0
		node.ClearValue()
	} else {
		// We never go here when FREEZE_ROOT is true
		node.setNodeAdded(replacement.IsAdded())
		node.setUpper(replacement.getUpperSubNode())
		node.setLower(replacement.getLowerSubNode())
		node.setKey(replacement.GetKey())
		node.storedSize = replacement.storedSize
		node.SetValue(replacement.GetValue())
	}
}

// Clear removes this node and all sub-nodes from the sub-tree with this node as the root, after which isEmpty() will return true.
func (node *binTreeNode[E, V]) Clear() {
	node.checkCopy()
	if node != nil {
		node.replaceThis(nil)
	}
}

// IsEmpty returns where there are not any elements in the sub-tree with this node as the root.
func (node *binTreeNode[E, V]) IsEmpty() bool {
	return !node.IsAdded() && node.getUpperSubNode() == nil && node.getLowerSubNode() == nil
}

// IsLeaf returns whether this node is in the tree (a node for which IsAdded() is true)
// and there are no elements in the sub-tree with this node as the root.
func (node *binTreeNode[E, V]) IsLeaf() bool {
	return node.IsAdded() && node.getUpperSubNode() == nil && node.getLowerSubNode() == nil
}

// Returns the first (lowest valued) node in the sub-tree originating from this node.
func (node *binTreeNode[E, V]) firstNode() *binTreeNode[E, V] {
	first := node
	for {
		lower := first.getLowerSubNode()
		if lower == nil {
			return first
		}
		first = lower
	}
}

// Returns the first (lowest valued) added node in the sub-tree originating from this node,
// or nil if there are no added entries in this tree or sub-tree
func (node *binTreeNode[E, V]) firstAddedNode() *binTreeNode[E, V] {
	first := node.firstNode()
	if first.IsAdded() {
		return first
	}
	return first.nextAddedNode()
}

// Returns the last (highest valued) node in the sub-tree originating from this node.
func (node *binTreeNode[E, V]) lastNode() *binTreeNode[E, V] {
	last := node
	for {
		upper := last.getUpperSubNode()
		if upper == nil {
			return last
		}
		last = upper
	}
}

// Returns the last (highest valued) added node in the sub-tree originating from this node,
// or nil if there are no added entries in this tree or sub-tree
func (node *binTreeNode[E, V]) lastAddedNode() *binTreeNode[E, V] {
	last := node.lastNode()
	if last.IsAdded() {
		return last
	}
	return last.previousAddedNode()
}

func (node *binTreeNode[E, V]) firstPostOrderNode() *binTreeNode[E, V] {
	next := node
	var nextNext *binTreeNode[E, V]
	for {
		nextNext = next.getLowerSubNode()
		if nextNext == nil {
			nextNext = next.getUpperSubNode()
			if nextNext == nil {
				return next
			}
		}
		next = nextNext
	}
}

func (node *binTreeNode[E, V]) lastPreOrderNode() *binTreeNode[E, V] {
	next := node
	var nextNext *binTreeNode[E, V]
	for {
		nextNext = next.getUpperSubNode()
		if nextNext == nil {
			nextNext = next.getLowerSubNode()
			if nextNext == nil {
				return next
			}
		}
		next = nextNext
	}
}

// Returns the node that follows this node following the tree order
func (node *binTreeNode[E, V]) nextNode() *binTreeNode[E, V] {
	return node.nextNodeBounded(nil)
}

//	in-order
//
//				8x
//		4x					12x
//	2x		6x			10x		14x
//
// 1x 3x		5x 7x		9x 11x	13x 15x
func (node *binTreeNode[E, V]) nextNodeBounded(bound *binTreeNode[E, V]) *binTreeNode[E, V] {
	next := node.getUpperSubNode()
	if next != nil {
		for {
			nextLower := next.getLowerSubNode()
			if nextLower == nil {
				return next
			}
			next = nextLower
		}
	} else {
		next = node.getParent()
		if next == bound {
			return nil
		}
		current := node
		for next != nil && current == next.getUpperSubNode() {
			current = next
			next = next.getParent()
			if next == bound {
				return nil
			}
		}
	}
	return next
}

// Returns the node that precedes this node following the tree order.
func (node *binTreeNode[E, V]) previousNode() *binTreeNode[E, V] {
	return node.previousNodeBounded(nil)
}

//	reverse order
//
//				8x
//		12x					4x
//	14x		10x			6x		2x
//
// 15x 13x	11x 9x		7x 5x	3x 1x
func (node *binTreeNode[E, V]) previousNodeBounded(bound *binTreeNode[E, V]) *binTreeNode[E, V] {
	previous := node.getLowerSubNode()
	if previous != nil {
		for {
			previousUpper := previous.getUpperSubNode()
			if previousUpper == nil {
				break
			}
			previous = previousUpper
		}
	} else {
		previous = node.getParent()
		if previous == bound {
			return nil
		}
		current := node
		for previous != nil && current == previous.getLowerSubNode() {
			current = previous
			previous = previous.getParent()
			if previous == bound {
				return nil
			}
		}
	}
	return previous
}

//	pre order
//				1x
//		2x						9x
//
// 3x		6x				10x		13x
// 4x 5x		7x 8x		11x 12x		14x 15x
// this one starts from root, ends at last node, all the way right
func (node *binTreeNode[E, V]) nextPreOrderNode(end *binTreeNode[E, V]) *binTreeNode[E, V] {
	next := node.getLowerSubNode()
	if next == nil {
		// cannot go left/lower
		next = node.getUpperSubNode()
		if next == nil {
			// cannot go right/upper
			current := node
			next = node.getParent()
			// so instead, keep going up until we can go right
			for next != nil {
				if next == end {
					return nil
				}
				if current == next.getLowerSubNode() {
					// parent is higher
					nextNext := next.getUpperSubNode()
					if nextNext != nil {
						return nextNext
					}
				}
				current = next
				next = next.getParent()
			}
		}
	}
	return next
}

//	reverse post order
//				1x
//		9x					2x
//	13x		10x			6x		3x
//
// 15x 14x	12x 11x		8x 7x	5x 4x
// this one starts from root, ends at first node, all the way left
// this is the mirror image of nextPreOrderNode, so no comments
func (node *binTreeNode[E, V]) previousPostOrderNode(end *binTreeNode[E, V]) *binTreeNode[E, V] {
	next := node.getUpperSubNode()
	if next == nil {
		next = node.getLowerSubNode()
		if next == nil {
			current := node
			next = node.getParent()
			for next != nil {
				if next == end {
					return nil
				}
				if current == next.getUpperSubNode() {
					nextNext := next.getLowerSubNode()
					if nextNext != nil {
						next = nextNext
						break
					}
				}
				current = next
				next = next.getParent()
			}
		}
	}
	return next
}

//	reverse pre order
//
//				15x
//		14x					7x
//	13x		10x			6x		3x
//12x 11x	9x 8x		5x 4x	2x 1x

// this one starts from last node, all the way right, ends at root
// this is the mirror image of nextPostOrderNode, so no comments
func (node *binTreeNode[E, V]) previousPreOrderNode(end *binTreeNode[E, V]) *binTreeNode[E, V] {
	next := node.getParent()
	if next == nil || next == end {
		return nil
	}
	if next.getLowerSubNode() == node {
		return next
	}
	nextNext := next.getLowerSubNode()
	if nextNext == nil {
		return next
	}
	next = nextNext
	for {
		nextNext = next.getUpperSubNode()
		if nextNext == nil {
			nextNext = next.getLowerSubNode()
			if nextNext == nil {
				return next
			}
		}
		next = nextNext
	}
}

//	post order
//				15x
//		7x					14x
//	3x		6x			10x		13x
//
// 1x 2x		4x 5x		8x 9x	11x 12x
// this one starts from first node, all the way left, ends at root
func (node *binTreeNode[E, V]) nextPostOrderNode(end *binTreeNode[E, V]) *binTreeNode[E, V] {
	next := node.getParent()
	if next == nil || next == end {
		return nil
	}
	if next.getUpperSubNode() == node {
		// we are the upper sub-node, so parent is next
		return next
	}
	// we are the lower sub-node
	nextNext := next.getUpperSubNode()
	if nextNext == nil {
		// parent has no upper sub-node, so parent is next
		return next
	}
	// go to parent's upper sub-node
	next = nextNext
	// now go all the way down until we can go no further, favoring left/lower turns over right/upper
	for {
		nextNext = next.getLowerSubNode()
		if nextNext == nil {
			nextNext = next.getUpperSubNode()
			if nextNext == nil {
				return next
			}
		}
		next = nextNext
	}
}

// Returns the next node in the tree that is an added node, following the tree order,
// or nil if there is no such node.
func (node *binTreeNode[E, V]) nextAddedNode() *binTreeNode[E, V] {
	return node.nextAdded(nil, (*binTreeNode[E, V]).nextNodeBounded)
}

// Returns the previous node in the tree that is an added node, following the tree order in reverse,
// or nil if there is no such node.
func (node *binTreeNode[E, V]) previousAddedNode() *binTreeNode[E, V] {
	return node.nextAdded(nil, (*binTreeNode[E, V]).previousNodeBounded)
}

// The generic method pointers are fine.  The parser errors are just a Goland problem.  Try it out in playground: https://go.dev/play/p/lf8zJtGCKYI

func nextTest[E Key, V any](current, end *binTreeNode[E, V], nextOperator func(current *binTreeNode[E, V], end *binTreeNode[E, V]) *binTreeNode[E, V], tester func(current *binTreeNode[E, V]) bool) *binTreeNode[E, V] {
	for {
		current = nextOperator(current, end)
		if current == end || current == nil {
			return nil
		}
		if tester(current) {
			break
		}
	}
	return current
}

func (node *binTreeNode[E, V]) nextAdded(end *binTreeNode[E, V], nextOperator func(current *binTreeNode[E, V], end *binTreeNode[E, V]) *binTreeNode[E, V]) *binTreeNode[E, V] {
	return nextTest(node, end, nextOperator, (*binTreeNode[E, V]).IsAdded)
}

func (node *binTreeNode[E, V]) nextInBounds(end *binTreeNode[E, V], nextOperator func(current *binTreeNode[E, V], end *binTreeNode[E, V]) *binTreeNode[E, V], bnds *bounds[E]) *binTreeNode[E, V] {
	return nextTest(node, end, nextOperator, func(current *binTreeNode[E, V]) bool {
		return bnds.isInBounds(current.GetKey())
	})
}

// Returns an iterator that iterates through the elements of the sub-tree with this node as the root.
// The iteration is in sorted element order.
func (node *binTreeNode[E, V]) iterator() keyIterator[E] {
	return binTreeKeyIterator[E, V]{node.nodeIterator(true)}
}

// Returns an iterator that iterates through the elements of the subtrie with this node as the root.
// The iteration is in reverse sorted element order.
func (node *binTreeNode[E, V]) descendingIterator() keyIterator[E] {
	return binTreeKeyIterator[E, V]{node.nodeIterator(false)}
}

// Iterates through the added nodes of the sub-tree with this node as the root, in forward or reverse tree order.
func (node *binTreeNode[E, V]) nodeIterator(forward bool) nodeIteratorRem[E, V] {
	return node.configuredNodeIterator(forward, true)
}

// Iterates through all the nodes of the sub-tree with this node as the root, in forward or reverse tree order.
func (node *binTreeNode[E, V]) allNodeIterator(forward bool) nodeIteratorRem[E, V] {
	return node.configuredNodeIterator(forward, false)
}

func (node *binTreeNode[E, V]) containingFirstIterator(forwardSubNodeOrder bool) cachingNodeIterator[E, V] {
	return node.containingFirstNodeIterator(forwardSubNodeOrder, true)
}

func (node *binTreeNode[E, V]) containingFirstAllNodeIterator(forwardSubNodeOrder bool) cachingNodeIterator[E, V] {
	return node.containingFirstNodeIterator(forwardSubNodeOrder, false)
}

func (node *binTreeNode[E, V]) containingFirstNodeIterator(forwardSubNodeOrder, addedNodesOnly bool) cachingNodeIterator[E, V] {
	var iter subNodeCachingIterator[E, V]
	if forwardSubNodeOrder {
		iter = newPreOrderNodeIterator[E, V]( // remove is allowed
			true,           // forward
			addedNodesOnly, // added only
			node,
			node.getParent(),
			node.getChangeTracker())
	} else {
		iter = newPostOrderNodeIterator[E, V]( // remove is allowed
			false,          // forward
			addedNodesOnly, // added only
			node,
			node.getParent(),
			node.getChangeTracker())
	}
	return &iter
}

func (node *binTreeNode[E, V]) containedFirstIterator(forwardSubNodeOrder bool) nodeIteratorRem[E, V] {
	return node.containedFirstNodeIterator(forwardSubNodeOrder, true)
}

func (node *binTreeNode[E, V]) containedFirstAllNodeIterator(forwardSubNodeOrder bool) nodeIterator[E, V] {
	return node.containedFirstNodeIterator(forwardSubNodeOrder, false)
}

func (node *binTreeNode[E, V]) containedFirstNodeIterator(forwardSubNodeOrder, addedNodesOnly bool) nodeIteratorRem[E, V] {
	var iter subNodeCachingIterator[E, V]
	if forwardSubNodeOrder {
		iter = newPostOrderNodeIterator[E, V]( // Remove is allowed if and only if added only
			true,
			addedNodesOnly, // added only
			node.firstPostOrderNode(),
			node.getParent(),
			node.getChangeTracker())
	} else {
		iter = newPreOrderNodeIterator[E, V]( // Remove is allowed if and only if added only
			false,
			addedNodesOnly, // added only
			node.lastPreOrderNode(),
			node.getParent(),
			node.getChangeTracker())
	}
	return &iter
}

func (node *binTreeNode[E, V]) configuredNodeIterator(forward, addedOnly bool) nodeIteratorRem[E, V] {
	var startNode *binTreeNode[E, V]
	if forward {
		startNode = node.firstNode()
	} else {
		startNode = node.lastNode()
	}
	return newNodeIterator[E, V](
		forward,
		addedOnly,
		startNode,
		node.getParent(),
		node.getChangeTracker())
}

// https://jrgraphix.net/r/Unicode/2500-257F
// https://jrgraphix.net/r/Unicode/25A0-25FF
const (
	nonAddedNodeCircle = "\u25cb"
	addedNodeCircle    = "\u25cf"

	leftElbow       = "\u251C\u2500" // |-
	inBetweenElbows = "\u2502 "      // |
	rightElbow      = "\u2514\u2500" // --
	belowElbows     = "  "
)

type nodePrinter[E Key, V any] interface {
	GetKey() E
	GetValue() V
	IsAdded() bool
}

func isNil[V any](v V) bool {
	valueType := reflect.ValueOf(&v).Elem()
	switch valueType.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return valueType.IsNil()
	}
	return false
}

// NodeString returns a visual representation of the given node including the key, with an open circle indicating this node is not an added node,
// a closed circle indicating this node is an added node.
func NodeString[E Key, V any](node nodePrinter[E, V]) string {
	if node == nil {
		return nilString()
	}
	key := node.GetKey()
	val := node.GetValue()

	if !node.IsAdded() {
		return fmt.Sprint(nonAddedNodeCircle, " ", key)
	} else if _, ok := any(val).(EmptyValueType); ok || isNil(val) {
		return fmt.Sprint(addedNodeCircle, " ", key)
	}
	return fmt.Sprint(addedNodeCircle, " ", key, " = ", val)
}

type indents struct {
	nodeIndent, subNodeInd string
}

// TreeString returns a visual representation of the sub-tree with this node as root, with one node per line.
//
// withNonAddedKeys: whether to show nodes that are not added nodes
// withSizes: whether to include the counts of added nodes in each sub-tree
func (node *binTreeNode[E, V]) TreeString(withNonAddedKeys, withSizes bool) string {
	builder := strings.Builder{}
	builder.WriteByte('\n')
	node.printTree(&builder, indents{}, withNonAddedKeys, withSizes)
	return builder.String()
}

func (node *binTreeNode[E, V]) printTree(
	builder *strings.Builder,
	initialIndents indents,
	withNonAdded,
	withSizes bool) {
	if node == nil {
		builder.WriteString(initialIndents.nodeIndent)
		builder.WriteString(nilString())
		builder.WriteByte('\n')
		return
	}
	iterator := node.containingFirstAllNodeIterator(true)
	next := iterator.Next()
	for next != nil {
		cached := iterator.GetCached()
		var nodeIndent, subNodeIndent string
		if cached == nil {
			nodeIndent = initialIndents.nodeIndent
			subNodeIndent = initialIndents.subNodeInd
		} else {
			cachedi := cached.(indents)
			nodeIndent = cachedi.nodeIndent
			subNodeIndent = cachedi.subNodeInd
		}
		if withNonAdded || next.IsAdded() {
			builder.WriteString(nodeIndent)
			builder.WriteString(next.String())
			if withSizes {
				builder.WriteString(" (")
				builder.WriteString(strconv.Itoa(next.Size()))
				builder.WriteByte(')')
			}
			builder.WriteByte('\n')
		} else {
			builder.WriteString(nodeIndent)
			builder.WriteString(nonAddedNodeCircle)
			builder.WriteByte('\n')
		}
		upper, lower := next.getUpperSubNode(), next.getLowerSubNode()
		if upper != nil {
			if lower != nil {
				lowerIndents := indents{
					nodeIndent: subNodeIndent + leftElbow,
					subNodeInd: subNodeIndent + inBetweenElbows,
				}
				iterator.CacheWithLowerSubNode(lowerIndents)
			}
			upperIndents := indents{
				nodeIndent: subNodeIndent + rightElbow,
				subNodeInd: subNodeIndent + belowElbows,
			}
			iterator.CacheWithUpperSubNode(upperIndents)
		} else if lower != nil {
			lowerIndents := indents{
				nodeIndent: subNodeIndent + rightElbow,
				subNodeInd: subNodeIndent + belowElbows,
			}
			iterator.CacheWithLowerSubNode(lowerIndents)
		}
		next = iterator.Next()
	}
}

func nilString() string {
	return "<nil>"
}

// Returns a visual representation of this node including the key, with an open circle indicating this node is not an added node,
// a closed circle indicating this node is an added node.
func (node *binTreeNode[E, V]) String() string {
	if node == nil {
		return NodeString[E, V](nil)
	}
	return NodeString[E, V](node)
}

func (node binTreeNode[E, V]) format(state fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		_, _ = state.Write([]byte(node.String()))
		return
	}
	s := flagsFromState(state, verb)
	_, _ = state.Write([]byte(fmt.Sprintf(s, binTreeNodePtr[E, V](node.self))))
}

// only used to eliminate the method set of *binTreeNode
type binTreeNodePtr[E Key, V any] *binTreeNode[E, V]

func flagsFromState(state fmt.State, verb rune) string {
	flags := "# +-0"
	vals := make([]rune, 0, len(flags)+5) // %, flags, width, '.', precision, verb
	vals = append(vals, '%')
	for i := 0; i < len(flags); i++ {
		b := flags[i]
		if state.Flag(int(b)) {
			vals = append(vals, rune(b))
		}
	}
	width, widthOK := state.Width()
	precision, precisionOK := state.Precision()
	if widthOK || precisionOK {
		var wpv string
		if widthOK && precisionOK {
			wpv = fmt.Sprintf("%d.%d%c", width, precision, verb)
		} else if widthOK {
			wpv = fmt.Sprintf("%d%c", width, verb)
		} else {
			wpv = fmt.Sprintf(".%d%c", precision, verb)
		}
		return string(vals) + wpv
	}
	vals = append(vals, verb)
	return string(vals)
}

// Clones the node.
// Keys remain the same, but the parent node and the lower and upper sub-nodes are all set to nil.
func (node *binTreeNode[E, V]) clone() *binTreeNode[E, V] {
	if node == nil {
		return nil
	}
	result := *node // maintains same key and value which are not copied
	result.setParent(nil)
	result.setLower(nil)
	result.setUpper(nil)
	if node.IsAdded() {
		result.storedSize = 1
	} else {
		result.storedSize = 0
	}
	// it is ok to have no change tracker, because the parent, lower and upper are nil
	// so any attempt to remove or clear will do nothing,
	// and you cannot add to nodes, you can only add to tries,
	// so no calls to the change tracker
	result.cTracker = nil
	// no need to make use of the shared pool
	result.pool = nil
	result.setAddr()
	return &result
}

func (node *binTreeNode[E, V]) cloneTreeNode(cTracker *changeTracker, pool *sync.Pool) *binTreeNode[E, V] {
	if node == nil {
		return nil
	}
	result := *node // maintains same key and value which are not copied
	result.setParent(nil)
	result.cTracker = cTracker
	result.pool = pool
	result.setAddr()
	return &result
}

func (node *binTreeNode[E, V]) cloneTreeTrackerBounds(ctracker *changeTracker, pool *sync.Pool, bnds *bounds[E]) *binTreeNode[E, V] {
	if node == nil {
		return nil
	}
	rootClone := node.cloneTreeNode(ctracker, pool)
	clonedNode := rootClone
	iterator := clonedNode.containingFirstAllNodeIterator(true).(*subNodeCachingIterator[E, V])
	recalculateSize := false
	for {
		lower := clonedNode.getLowerSubNode()
		if bnds != nil {
			for {
				if lower == nil {
					break
				} else if bnds.isWithinLowerBound(lower.GetKey()) {
					if !lower.IsAdded() {
						next := lower.getLowerSubNode()
						for bnds.isBelowLowerBound(next.GetKey()) {
							next = next.getUpperSubNode()
							if next == nil {
								lower = lower.getUpperSubNode()
								recalculateSize = true
								break
							}
						}
					}
					break
				}
				recalculateSize = true
				// outside bounds, try again
				lower = lower.getUpperSubNode()
			}
		}
		if lower != nil {
			clonedNode.setLower(lower.cloneTreeNode(ctracker, pool))
		} else {
			clonedNode.setLower(nil)
		}
		upper := clonedNode.getUpperSubNode()
		if bnds != nil {
			for {
				if upper == nil {
					break
				} else if bnds.isWithinUpperBound(upper.GetKey()) {
					if !upper.IsAdded() {
						next := upper.getUpperSubNode()
						for bnds.isAboveUpperBound(next.GetKey()) {
							next = next.getLowerSubNode()
							if next == nil {
								upper = upper.getLowerSubNode()
								recalculateSize = true
								break
							}
						}
					}

					break
				}
				recalculateSize = true
				// outside bounds, try again
				upper = upper.getLowerSubNode()
			}
		}
		if upper != nil {
			clonedNode.setUpper(upper.cloneTreeNode(ctracker, pool))
		} else {
			clonedNode.setUpper(nil)
		}
		iterator.Next() // returns current clonedNode
		clonedNode = iterator.next
		if !iterator.HasNext() { /* basically this checks clonedNode != nil */
			break
		}
	}
	if !rootClone.IsAdded() && !node.IsRoot() {
		lower := rootClone.getLowerSubNode()
		if lower == nil {
			rootClone = rootClone.getUpperSubNode()
		} else if rootClone.getUpperSubNode() == nil {
			rootClone = lower
		}
	}
	if recalculateSize && rootClone != nil {
		rootClone.storedSize = sizeUnknown
		rootClone.Size()
	}
	return rootClone
}
