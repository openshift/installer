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
	"strings"
	"sync"
	"unsafe"
)

// BinTrie is a binary trie.
//
// To use BinTrie, your keys implement TrieKey.
//
// All keys are either fixed, in which the key value does not change,
// or comprising of a prefix in which an initial sequence of bits does not change, and the the remaining bits represent all bit values.
// The length of the initial fixed sequence of bits is the prefix length.
// The total bit length is the same for all keys.
//
// A key with a prefix is also known as a prefix block, and represents all bit sequences with the same prefix.
//
// The zero value for BinTrie is a binary trie ready for use.
//
// Each node can be associated with a value, making BinTrie an associative binary trie.
// If you do not wish to associate values to nodes, then use the type EmptyValueType,
// in which case the value will be ignored in functions that print node strings.
type BinTrie[E TrieKey[E], V any] struct {
	binTree[E, V]
}

type EmptyValueType struct{}

func (trie *BinTrie[E, V]) toBinTree() *binTree[E, V] {
	return (*binTree[E, V])(unsafe.Pointer(trie))
}

// GetRoot returns the root of this trie (in the case of bounded tries, this would be the bounded root)
func (trie *BinTrie[E, V]) GetRoot() (root *BinTrieNode[E, V]) {
	if trie != nil {
		root = toTrieNode(trie.root)
	}
	return
}

// Returns the root of this trie (in the case of bounded tries, the absolute root ignores the bounds)
func (trie *BinTrie[E, V]) absoluteRoot() (root *BinTrieNode[E, V]) {
	if trie != nil {
		root = toTrieNode(trie.root)
	}
	return
}

// Size returns the number of elements in the tree.
// Only nodes for which IsAdded() returns true are counted.
// When zero is returned, IsEmpty() returns true.
func (trie *BinTrie[E, V]) Size() int {
	return trie.toBinTree().Size()
}

// NodeSize returns the number of nodes in the tree, which is always more than the number of elements.
func (trie *BinTrie[E, V]) NodeSize() int {
	return trie.toBinTree().NodeSize()
}

// EnsureRoot ensures this trie has a root, whose type will match the type and version of the given key.
// The given key is not added to the trie when calling this method,
// instead it is used to construct an address of the given address type and version with a prefix length of 0 and with a host that spans all addresses, which will become the root key.
// Once the root is added, the version and type of the trie is fixed.
// So, for instance, if a trie has key type of IPAddress, then the root could be version IPv6 or IPv4, one of either ::/0 or 0.0.0.0/0,
// Once the root exists, all keys to follow must match that same version (ie the same bit-size).
// Another consequence of fixing the root is that any copy of the BinTrie instance will reference the same underlying trie.
// If the trie has no root, then any copy of the trie will result in an entirely new and undependent trie.
// Calling this method will have no effect if the trie root already exists.
// Calling this method is not necessary, adding any key to the trie will also cause the tree root to be created, matching the type and version of the added key.
// It only makes sense to call this method on empty tries and only if you wish to fix the root to a given type or version,
// or if you wish copies of the empty trie instance to reference the same underlying trie structure.
func (trie *BinTrie[E, V]) EnsureRoot(key E) *BinTrieNode[E, V] {
	root := trie.root
	if root == nil {
		root = trie.setRoot(key.ToPrefixBlockLen(0))
	}
	return toTrieNode(root)
}

func (trie *BinTrie[E, V]) setRoot(key E) *binTreeNode[E, V] {
	root := &binTreeNode[E, V]{
		item:     key,
		cTracker: &changeTracker{},
		pool: &sync.Pool{
			New: func() any { return &opResult[E, V]{} },
		},
	}
	root.setAddr()
	trie.root = root
	return root
}

// Iterator returns an iterator that iterates through the elements of the trie.
// The iteration is in sorted element order.
func (trie *BinTrie[E, V]) Iterator() TrieKeyIterator[E] {
	return trie.GetRoot().Iterator()
}

// DescendingIterator returns an iterator that iterates through the elements of the trie.
// The iteration is in reverse sorted element order.
func (trie *BinTrie[E, V]) DescendingIterator() TrieKeyIterator[E] {
	return trie.GetRoot().DescendingIterator()
}

// String returns a visual representation of the tree with one node per line.
// It is equivalent to calling TreeString(true)
func (trie *BinTrie[E, V]) String() string {
	if trie == nil {
		return nilString()
	}
	return trie.binTree.String()
}

// TreeString returns a visual representation of the tree with one node per line, with or without the non-added keys.
func (trie *BinTrie[E, V]) TreeString(withNonAddedKeys bool) string {
	if trie == nil {
		return "\n" + nilString()
	}
	return trie.binTree.TreeString(withNonAddedKeys)
}

// Add adds the given key to the trie, returning true if not there already.
func (trie *BinTrie[E, V]) Add(key E) bool {
	root := trie.EnsureRoot(key)
	result := &opResult[E, V]{
		key: key,
		op:  insert,
	}
	root.matchBits(result)
	return !result.exists
}

// AddNode is similar to Add but returns the new or existing node.
func (trie *BinTrie[E, V]) AddNode(key E) *BinTrieNode[E, V] {
	root := trie.EnsureRoot(key)
	result := &opResult[E, V]{
		key: key,
		op:  insert,
	}
	root.matchBits(result)
	node := result.existingNode
	if node == nil {
		node = result.inserted
	}
	return node
}

func (trie *BinTrie[E, V]) addNode(result *opResult[E, V], fromNode *BinTrieNode[E, V]) *BinTrieNode[E, V] {
	fromNode.matchBitsFromIndex(fromNode.GetKey().GetPrefixLen().Len(), result)
	node := result.existingNode
	if node == nil {
		return result.inserted
	}
	return node
}

func (trie *BinTrie[E, V]) addTrie(addedTreeNode *BinTrieNode[E, V], withValues bool) *BinTrieNode[E, V] {
	var valueMap func(V) V
	if withValues {
		valueMap = func(v V) V { return v }
	}
	return addTrie(trie, addedTreeNode, func(e E) E { return e }, valueMap)
}

func addTrie[ED TrieKey[ED], ES TrieKey[ES], VD, VS any](
	targetTrie *BinTrie[ED, VD],
	sourceNode *BinTrieNode[ES, VS],
	keyMap func(ES) ED,
	valueMap func(VS) VD) *BinTrieNode[ED, VD] {
	if sourceNode == nil { // addedTreeNode can be nil when the root of a zero-valued trie
		return nil
	}
	sourceIterator := sourceNode.ContainingFirstAllNodeIterator(true)
	toAdd := sourceIterator.Next()
	firstKey := keyMap(toAdd.GetKey())
	result := &opResult[ED, VD]{
		key: firstKey,
		op:  insert,
	}
	var firstNode *BinTrieNode[ED, VD]
	root := targetTrie.EnsureRoot(firstKey)
	firstAdded := toAdd.IsAdded()
	if firstAdded {
		if valueMap != nil {
			result.newValue = valueMap(toAdd.GetValue())
			// new value assignment
		}
		firstNode = targetTrie.addNode(result, root)
	} else {
		firstNode = root
	}
	lastAddedNode := firstNode
	for sourceIterator.HasNext() {
		sourceIterator.CacheWithLowerSubNode(lastAddedNode)
		sourceIterator.CacheWithUpperSubNode(lastAddedNode)
		toAdd = sourceIterator.Next()
		cachedNode := sourceIterator.GetCached().(*BinTrieNode[ED, VD])
		if toAdd.IsAdded() {
			addrNext := toAdd.GetKey()
			result.key = keyMap(addrNext)
			result.existingNode = nil
			result.inserted = nil
			if valueMap != nil {
				result.newValue = valueMap(toAdd.GetValue())
				// new value assignment
			}
			lastAddedNode = targetTrie.addNode(result, cachedNode)
		} else {
			lastAddedNode = cachedNode
		}
	}
	if !firstAdded {
		firstNode = targetTrie.GetNode(keyMap(sourceNode.GetKey()))
	}
	return firstNode
}

// AddTrieKeys copies the trie node structure of addedTreeNode into trie, but does not copy node mapped values
func AddTrieKeys[E TrieKey[E], V1 any, V2 any](trie *BinTrie[E, V1], addedTreeNode *BinTrieNode[E, V2]) *BinTrieNode[E, V1] {
	if addedTreeNode == nil { // addedTreeNode can be nil when the root of a zero-valued trie
		return nil
	}
	iterator := addedTreeNode.ContainingFirstAllNodeIterator(true)
	toAdd := iterator.Next()
	firstKey := toAdd.GetKey()
	result := &opResult[E, V1]{
		key: firstKey,
		op:  insert,
	}
	var firstNode *BinTrieNode[E, V1]
	root := trie.EnsureRoot(firstKey)
	firstAdded := toAdd.IsAdded()
	if firstAdded {
		firstNode = trie.addNode(result, root)
	} else {
		firstNode = root
	}
	lastAddedNode := firstNode
	for iterator.HasNext() {
		iterator.CacheWithLowerSubNode(lastAddedNode)
		iterator.CacheWithUpperSubNode(lastAddedNode)
		toAdd = iterator.Next()
		cachedNode := iterator.GetCached().(*BinTrieNode[E, V1])
		if toAdd.IsAdded() {
			result.key = toAdd.GetKey()
			result.existingNode = nil
			result.inserted = nil
			lastAddedNode = trie.addNode(result, cachedNode)
		} else {
			lastAddedNode = cachedNode
		}
	}
	if !firstAdded {
		firstNode = trie.GetNode(addedTreeNode.GetKey())
	}
	return firstNode
}

func (trie *BinTrie[E, V]) AddTrie(trieNode *BinTrieNode[E, V]) *BinTrieNode[E, V] {
	if trieNode == nil {
		return nil
	}
	trie.EnsureRoot(trieNode.GetKey())
	return trie.addTrie(trieNode, false)
}

// AddConvertibleTrie maps from one trie type to another.
// The keys must follow the same containment structure, meaning any sub-node to parent-node relationship must be maintained across mappings.
// If the containment structure is not maintained, then do not use this function, instead simply iterate through the source trie and add each to the target trie.
// If the valueMap is non-nil, then the values will be mapped by that function, otherwise the new values will all be zero values.
func AddConvertibleTrie[ED TrieKey[ED], ES TrieKey[ES], VD, VS any](
	target *BinTrie[ED, VD],
	sourceNode *BinTrieNode[ES, VS],
	keyMap func(ES) ED,
	valueMap func(VS) VD) *BinTrieNode[ED, VD] {
	if sourceNode == nil {
		return nil
	}
	target.EnsureRoot(keyMap(sourceNode.GetKey()))
	return addTrie(target, sourceNode, keyMap, valueMap)
}

func (trie *BinTrie[E, V]) Contains(key E) bool {
	return trie.absoluteRoot().Contains(key)
}

func (trie *BinTrie[E, V]) Remove(key E) bool {
	return trie.absoluteRoot().RemoveNode(key)
}

func (trie *BinTrie[E, V]) RemoveElementsContainedBy(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().RemoveElementsContainedBy(key)
}

func (trie *BinTrie[E, V]) ElementsContainedBy(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().ElementsContainedBy(key)
}

func (trie *BinTrie[E, V]) ElementsContaining(key E) *Path[E, V] {
	return trie.absoluteRoot().ElementsContaining(key)
}

// LongestPrefixMatch finds the added key with the longest matching prefix.
func (trie *BinTrie[E, V]) LongestPrefixMatch(key E) (E, bool) {
	return trie.absoluteRoot().LongestPrefixMatch(key)
}

// LongestPrefixMatchNode finds the added node with the longest matching prefix.
func (trie *BinTrie[E, V]) LongestPrefixMatchNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().LongestPrefixMatchNode(key)
}

func (trie *BinTrie[E, V]) ElementContains(key E) bool {
	return trie.absoluteRoot().ElementContains(key)
}

// ShortestPrefixMatch finds the added key with the shortest matching prefix.
func (trie *BinTrie[E, V]) ShortestPrefixMatch(key E) (E, bool) {
	return trie.absoluteRoot().ShortestPrefixMatch(key)
}

// ShortestPrefixMatchNode finds the added node whose key has the shortest matching prefix.
func (trie *BinTrie[E, V]) ShortestPrefixMatchNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().ShortestPrefixMatchNode(key)
}

// GetNode gets the node in the sub-trie corresponding to the given address,
// or returns nil if not such element exists.
//
// It returns any node, whether added or not,
// including any prefix block node that was not added.
func (trie *BinTrie[E, V]) GetNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().GetNode(key)
}

// GetAddedNode gets trie nodes representing added elements.
//
// Use Contains to check for the existence of a given address in the trie,
// as well as GetNode to search for all nodes including those not-added but also auto-generated nodes for subnet blocks.
func (trie *BinTrie[E, V]) GetAddedNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().GetAddedNode(key)
}

// Put associates the specified value with the specified key in this map.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The Partition type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// If this map previously contained a mapping for a key,
// the old value is replaced by the specified value, and false is returned along with the old value, which may be the zero value.
// If this map did not previously contain a mapping for the key, true is returned along with the zero value.
func (trie *BinTrie[E, V]) Put(key E, value V) (V, bool) {
	root := trie.EnsureRoot(key)
	result := &opResult[E, V]{
		key:      key,
		op:       insert,
		newValue: value,
		// new value assignment
	}
	root.matchBits(result)
	return result.existingValue, !result.exists
}

func (trie *BinTrie[E, V]) PutTrie(trieNode *BinTrieNode[E, V]) *BinTrieNode[E, V] {
	if trieNode == nil {
		return nil
	}
	trie.EnsureRoot(trieNode.GetKey())
	return trie.addTrie(trieNode, true)
}

func (trie *BinTrie[E, V]) PutNode(key E, value V) *BinTrieNode[E, V] {
	root := trie.EnsureRoot(key)
	result := &opResult[E, V]{
		key:      key,
		op:       insert,
		newValue: value,
		// new value assignment
	}
	root.matchBits(result)
	resultNode := result.existingNode
	if resultNode == nil {
		resultNode = result.inserted
	}
	return resultNode
}

func (trie *BinTrie[E, V]) Remap(key E, remapper func(existing V, found bool) (mapped V, mapIt bool)) *BinTrieNode[E, V] {
	return trie.remapImpl(key,
		func(existingVal V, exists bool) (V, remapAction) {
			result, mapIt := remapper(existingVal, exists)
			if mapIt {
				return result, remapValue
			}
			var v V
			return v, removeNode
		})
}

func (trie *BinTrie[E, V]) RemapIfAbsent(key E, supplier func() V) *BinTrieNode[E, V] {
	return trie.remapImpl(key,
		func(existingVal V, exists bool) (V, remapAction) {
			if !exists {
				return supplier(), remapValue
			}
			var v V
			return v, doNothing
		})
}

func (trie *BinTrie[E, V]) remapImpl(key E, remapper func(val V, exists bool) (V, remapAction)) *BinTrieNode[E, V] {
	root := trie.EnsureRoot(key)
	result := &opResult[E, V]{
		key:      key,
		op:       remap,
		remapper: remapper,
	}
	root.matchBits(result)
	resultNode := result.existingNode
	if resultNode == nil {
		resultNode = result.inserted
	}
	return resultNode
}

func (trie *BinTrie[E, V]) Get(key E) (V, bool) {
	return trie.absoluteRoot().Get(key)
}

// NodeIterator returns an iterator that iterates through the added nodes of the trie in forward or reverse tree order.
func (trie *BinTrie[E, V]) NodeIterator(forward bool) TrieNodeIteratorRem[E, V] {
	return trie.absoluteRoot().NodeIterator(forward)
}

// AllNodeIterator returns an iterator that iterates through all the nodes of the trie in forward or reverse tree order.
func (trie *BinTrie[E, V]) AllNodeIterator(forward bool) TrieNodeIteratorRem[E, V] {
	return trie.absoluteRoot().AllNodeIterator(forward)
}

// BlockSizeNodeIterator returns an iterator that iterates the added nodes in the trie, ordered by keys from largest prefix blocks to smallest, and then to individual addresses.
//
// If lowerSubNodeFirst is true, for blocks of equal size the lower is first, otherwise the reverse order
func (trie *BinTrie[E, V]) BlockSizeNodeIterator(lowerSubNodeFirst bool) TrieNodeIteratorRem[E, V] {
	return trie.absoluteRoot().BlockSizeNodeIterator(lowerSubNodeFirst)
}

// BlockSizeAllNodeIterator returns an iterator that iterates all nodes in the trie, ordered by keys from largest prefix blocks to smallest, and then to individual addresses.
//
// If lowerSubNodeFirst is true, for blocks of equal size the lower is first, otherwise the reverse order
func (trie *BinTrie[E, V]) BlockSizeAllNodeIterator(lowerSubNodeFirst bool) TrieNodeIteratorRem[E, V] {
	return trie.absoluteRoot().BlockSizeAllNodeIterator(lowerSubNodeFirst)
}

// BlockSizeCachingAllNodeIterator returns an iterator that iterates all nodes, ordered by keys from largest prefix blocks to smallest, and then to individual addresses.
func (trie *BinTrie[E, V]) BlockSizeCachingAllNodeIterator() CachingTrieNodeIterator[E, V] {
	return trie.absoluteRoot().BlockSizeCachingAllNodeIterator()
}

func (trie *BinTrie[E, V]) ContainingFirstIterator(forwardSubNodeOrder bool) CachingTrieNodeIterator[E, V] {
	return trie.absoluteRoot().ContainingFirstIterator(forwardSubNodeOrder)
}

func (trie *BinTrie[E, V]) ContainingFirstAllNodeIterator(forwardSubNodeOrder bool) CachingTrieNodeIterator[E, V] {
	return trie.absoluteRoot().ContainingFirstAllNodeIterator(forwardSubNodeOrder)
}

func (trie *BinTrie[E, V]) ContainedFirstIterator(forwardSubNodeOrder bool) TrieNodeIteratorRem[E, V] {
	return trie.absoluteRoot().ContainedFirstIterator(forwardSubNodeOrder)
}

func (trie *BinTrie[E, V]) ContainedFirstAllNodeIterator(forwardSubNodeOrder bool) TrieNodeIterator[E, V] {
	return trie.absoluteRoot().ContainedFirstAllNodeIterator(forwardSubNodeOrder)
}

func (trie *BinTrie[E, V]) FirstNode() *BinTrieNode[E, V] {
	return trie.absoluteRoot().FirstNode()
}

func (trie *BinTrie[E, V]) FirstAddedNode() *BinTrieNode[E, V] {
	return trie.absoluteRoot().FirstAddedNode()
}

func (trie *BinTrie[E, V]) LastNode() *BinTrieNode[E, V] {
	return trie.absoluteRoot().LastNode()
}

func (trie *BinTrie[E, V]) LastAddedNode() *BinTrieNode[E, V] {
	return trie.absoluteRoot().LastAddedNode()
}

func (trie *BinTrie[E, V]) LowerAddedNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().LowerAddedNode(key)
}

func (trie *BinTrie[E, V]) FloorAddedNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().FloorAddedNode(key)
}

func (trie *BinTrie[E, V]) HigherAddedNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().HigherAddedNode(key)
}

func (trie *BinTrie[E, V]) CeilingAddedNode(key E) *BinTrieNode[E, V] {
	return trie.absoluteRoot().CeilingAddedNode(key)
}

func (trie *BinTrie[E, V]) Clone() *BinTrie[E, V] {
	if trie == nil {
		return nil
	}
	return &BinTrie[E, V]{binTree[E, V]{root: trie.absoluteRoot().CloneTree().toBinTreeNode()}}
}

// DeepEqual returns whether the given argument is a trie with a set of nodes with the same keys as in this trie according to the Compare method,
// and the same values according to the reflect.DeepEqual method
func (trie *BinTrie[E, V]) DeepEqual(other *BinTrie[E, V]) bool {
	return trie.absoluteRoot().TreeDeepEqual(other.absoluteRoot())
}

// Equal returns whether the given argument is a trie with a set of nodes with the same keys as in this trie according to the Compare method
func (trie *BinTrie[E, V]) Equal(other *BinTrie[E, V]) bool {
	return trie.absoluteRoot().TreeEqual(other.absoluteRoot())
}

// For some reason Format must be here and not in addressTrieNode for nil node.
// It panics in fmt code either way, but if in here then it is handled by a recover() call in fmt properly.
// Seems to be a problem only in the debugger.

// Format implements the fmt.Formatter interface
func (trie BinTrie[E, V]) Format(state fmt.State, verb rune) {
	trie.format(state, verb)
}

// NewBinTrie creates a new trie with root key.ToPrefixBlockLen(0).
// If the key argument is not Equal to its zero-length prefix block, then the key will be added as well.
func NewBinTrie[E TrieKey[E], V any](key E) BinTrie[E, V] {
	trie := BinTrie[E, V]{binTree[E, V]{}}
	root := key.ToPrefixBlockLen(0)
	trie.setRoot(root)
	if key.Compare(root) != 0 {
		trie.Add(key)
	}
	return trie
}

func TreesString[E TrieKey[E], V any](withNonAddedKeys bool, tries ...*BinTrie[E, V]) string {
	binTrees := make([]*binTree[E, V], 0, len(tries))
	for _, trie := range tries {
		binTrees = append(binTrees, tobinTree(trie))
	}
	return treesString(withNonAddedKeys, true, (*binTree[E, V]).printTree, binTrees...)
}

func tobinTree[E TrieKey[E], V any](trie *BinTrie[E, V]) *binTree[E, V] {
	return (*binTree[E, V])(unsafe.Pointer(trie))
}

// ConstructAddedNodesTree provides an associative trie in which the root and each added node of this trie are mapped to a list of their respective direct added sub-nodes.
// This trie provides an alternative non-binary tree structure of the added nodes.
// It is used by ToAddedNodesTreeString to produce a string showing the alternative structure.
// If there are no non-added nodes in this trie,
// then the alternative tree structure provided by this method is the same as the original trie.
// The trie values of this trie are of type []*BinTrieNode
func (trie *BinTrie[E, V]) ConstructAddedNodesTree() BinTrie[E, AddedSubnodeMapping] {
	newTrie := BinTrie[E, AddedSubnodeMapping]{binTree[E, AddedSubnodeMapping]{}}
	existingRoot := trie.GetRoot()
	if existingRoot == nil {
		return newTrie
	}

	newTrie.setRoot(existingRoot.GetKey())

	// populate the keys from the original trie into the new trie
	AddTrieKeys(&newTrie, existingRoot)

	// now, as we iterate,
	// we find our parent and add ourselves to that parent's list of subnodes

	var cachingIterator CachingTrieNodeIterator[E, AddedSubnodeMapping]
	cachingIterator = newTrie.ContainingFirstAllNodeIterator(true)
	thisIterator := trie.ContainingFirstAllNodeIterator(true)
	var newNext *BinTrieNode[E, AddedSubnodeMapping]
	var thisNext *BinTrieNode[E, V]
	for newNext, thisNext = cachingIterator.Next(), thisIterator.Next(); newNext != nil; newNext, thisNext = cachingIterator.Next(), thisIterator.Next() {

		// populate the values from the original trie into the new trie
		newNext.SetValue(SubNodesMapping[E, V]{Value: thisNext.GetValue()})

		cachingIterator.CacheWithLowerSubNode(newNext)
		cachingIterator.CacheWithUpperSubNode(newNext)

		// the cached object is our parent
		if newNext.IsAdded() {
			var parent *BinTrieNode[E, AddedSubnodeMapping]
			cached := cachingIterator.GetCached()
			if cached != nil {
				parent = cached.(*BinTrieNode[E, AddedSubnodeMapping])

				if parent != nil {
					// find added parent, or the root if no added parent
					for !parent.IsAdded() {
						parentParent := parent.GetParent()
						if parentParent == nil {
							break
						}
						parent = parentParent
					}
					// store ourselves with that added parent or root
					var val SubNodesMapping[E, V]
					val = parent.GetValue().(SubNodesMapping[E, V])
					var list []*BinTrieNode[E, AddedSubnodeMapping]
					if val.SubNodes == nil {
						list = make([]*BinTrieNode[E, AddedSubnodeMapping], 0, 3)
					} else {
						list = val.SubNodes
					}
					val.SubNodes = append(list, newNext)
					parent.SetValue(val)
				} // else root
			}
		}
	}
	return newTrie
}

// AddedNodesTreeString provides a flattened version of the trie showing only the contained added nodes and their containment structure, which is non-binary.
// The root node is included, which may or may not be added.
func (trie *BinTrie[E, V]) AddedNodesTreeString() string {
	if trie == nil {
		return "\n" + nilString()
	}
	addedTree := trie.ConstructAddedNodesTree()
	return AddedNodesTreeString[E, V](addedTree.GetRoot())
}

func AddedNodesTreesString[E TrieKey[E], V any](addedTrees ...*BinTrie[E, AddedSubnodeMapping]) string {
	// one complication here is that you may wish to omit the usually non-added ::/0 and 0.0.0.0/0 nodes.
	// In such cases, you actually want to drop the root, and supply to treesString each of the nodes in AddedSubnodeMapping
	// See below how nextNodes is populated.
	var binTrees []*binTree[E, AddedSubnodeMapping]
	for _, trie := range addedTrees {
		if trie == nil {
			binTrees = append(binTrees, nil)
		} else {
			trieNode := trie.GetRoot()
			if trieNode.IsAdded() {
				binTrees = append(binTrees, &trie.binTree)
			} else {
				nextVal := trieNode.GetValue()
				if nextVal != nil {
					mapping := nextVal.(SubNodesMapping[E, V])
					if mapping.SubNodes != nil {
						for _, trieNode = range mapping.SubNodes {
							binTrees = append(binTrees, &binTree[E, AddedSubnodeMapping]{
								root: trieNode.toBinTreeNode(),
							})
						}
					}
				}
			}
		}
	}
	return treesString(
		true,
		false,
		func(tree *binTree[E, AddedSubnodeMapping], builder *strings.Builder, inds indents, withNonAddedKeys bool) {
			if tree == nil {
				addedNodesTreeString[E, V](builder, inds, nil)
			} else {
				addedNodesTreeString[E, V](builder, inds, toTrieNode(tree.GetRoot()))
			}
		},
		binTrees...)
}

func AddedNodesTreeString[E TrieKey[E], V any](addedTreeRoot *BinTrieNode[E, AddedSubnodeMapping]) string {
	builder := strings.Builder{}
	builder.WriteByte('\n')
	addedNodesTreeString[E, V](&builder, indents{}, addedTreeRoot)
	return builder.String()
}

func addedNodesTreeString[E TrieKey[E], V any](builder *strings.Builder, inds indents, addedTreeRoot *BinTrieNode[E, AddedSubnodeMapping]) {
	var stack []indentsNode[E]
	nextNode := addedTreeRoot
	for {
		builder.WriteString(inds.nodeIndent)
		builder.WriteString(NodeString[E, V](printWrapper[E, V]{nextNode}))
		builder.WriteByte('\n')

		var nextVal AddedSubnodeMapping // SubNodesMapping[E, V]
		nextVal = nextNode.GetValue()
		var nextNodes []*BinTrieNode[E, AddedSubnodeMapping]
		if nextVal != nil {
			mapping := nextVal.(SubNodesMapping[E, V])
			if mapping.SubNodes != nil {
				nextNodes = mapping.SubNodes
			}
		}
		if len(nextNodes) > 0 {
			i := len(nextNodes) - 1
			subNodeIndent := inds.subNodeInd
			lastIndents := indents{
				nodeIndent: subNodeIndent + rightElbow,
				subNodeInd: subNodeIndent + belowElbows,
			}

			var nNode *BinTrieNode[E, AddedSubnodeMapping] // SubNodesMapping[E, V]
			nNode = nextNodes[i]
			if stack == nil {
				stack = make([]indentsNode[E], 0, addedTreeRoot.Size())
			}
			stack = append(stack, indentsNode[E]{lastIndents, nNode})
			if len(nextNodes) > 1 {
				firstIndents := indents{
					nodeIndent: subNodeIndent + leftElbow,
					subNodeInd: subNodeIndent + inBetweenElbows,
				}
				for i--; i >= 0; i-- {
					nNode = nextNodes[i]
					stack = append(stack, indentsNode[E]{firstIndents, nNode})
				}
			}
		}
		stackLen := len(stack)
		if stackLen == 0 {
			break
		}
		newLen := stackLen - 1
		nextItem := stack[newLen]
		stack = stack[:newLen]
		nextNode = nextItem.node
		inds = nextItem.inds
	}
}

type SubNodesMapping[E TrieKey[E], V any] struct {
	Value V

	// subNodes is the list of direct and indirect added subnodes in the original trie
	SubNodes []*BinTrieNode[E, AddedSubnodeMapping]
}

type AddedSubnodeMapping any // AddedSubnodeMapping / any is always SubNodesMapping[E,V]

type printWrapper[E TrieKey[E], V any] struct {
	*BinTrieNode[E, AddedSubnodeMapping]
}

func (p printWrapper[E, V]) GetValue() V {
	var nodeValue AddedSubnodeMapping = p.BinTrieNode.GetValue()
	if nodeValue == nil {
		var v V
		return v
	}
	return nodeValue.(SubNodesMapping[E, V]).Value
}

type indentsNode[E TrieKey[E]] struct {
	inds indents
	node *BinTrieNode[E, AddedSubnodeMapping]
}
