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
//

package ipaddr

import (
	"fmt"

	"github.com/seancfoley/bintree/tree"
)

type baseDualIPv4v6Tries[V any] struct {
	ipv4Trie, ipv6Trie trieBase[*IPAddress, V]
}

// ensureRoots ensures both the IPv4 and IPv6 tries each have a root.
// The roots will be 0.0.0.0/0 and ::/0 for the IPv4 and IPv6 tries, respectively.
// Once the roots are created, any copy of the this instance will reference the same underlying tries.
// Any copy of either trie, such as through GetIPv4Trie or GetIPv6Trie, will reference the same underlying trie.
// Calling this method will have no effect if both trie roots already exist.
// Calling this method is not necessary, adding any key to a trie will also cause its trie root to be created.
// It only makes sense to call this method if at least one trie is empty, and you wish to ensure that copies of this instance or either of its tries will reference the same underlying trie structures.
func (tries *baseDualIPv4v6Tries[V]) ensureRoots() {
	if tries.ipv4Trie.getRoot() == nil {
		key := IPv4Address{}
		tries.ipv4Trie.trie.EnsureRoot(trieKey[*IPAddress]{key.ToIP()})
	}
	if tries.ipv6Trie.getRoot() == nil {
		key := IPv6Address{}
		tries.ipv6Trie.trie.EnsureRoot(trieKey[*IPAddress]{key.ToIP()})
	}
}

// equal returns whether the two given pair of tries is equal to this pair of tries
func (tries *baseDualIPv4v6Tries[V]) equal(other *baseDualIPv4v6Tries[V]) bool {
	return tries.ipv4Trie.equal(&other.ipv4Trie) && tries.ipv6Trie.equal(&other.ipv6Trie)
}

// Clone clones this pair of tries.
func (tries *baseDualIPv4v6Tries[V]) clone() baseDualIPv4v6Tries[V] {
	return baseDualIPv4v6Tries[V]{
		ipv4Trie: *toBase(tries.ipv4Trie.clone()),
		ipv6Trie: *toBase(tries.ipv6Trie.clone()),
	}
}

// Clear removes all added nodes from the tries, after which IsEmpty will return true.
func (tries *baseDualIPv4v6Tries[V]) Clear() {
	tries.ipv4Trie.clear()
	tries.ipv6Trie.clear()
}

// Size returns the number of elements in the tries.
// Only added nodes are counted.
// When zero is returned, IsEmpty() returns true.
func (tries *baseDualIPv4v6Tries[V]) Size() int {
	return tries.ipv4Trie.size() + tries.ipv6Trie.size()
}

// iIsEmpty returns true if there are no added nodes within the two tries
func (tries *baseDualIPv4v6Tries[V]) IsEmpty() bool {
	return tries.ipv4Trie.size() == 0 && tries.ipv6Trie.size() == 0
}

// Add adds the given single address or prefix block subnet to one of the two tries.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Add returns true if the prefix block or address was inserted, false if already in one of the two tries.
func (tries *baseDualIPv4v6Tries[V]) Add(addr *IPAddress) bool {
	return addressFuncOp(addr, tries.ipv4Trie.add, tries.ipv6Trie.add)
}

// Contains returns whether the given address or prefix block subnet is in one of the two tries (as an added element).
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Contains returns true if the prefix block or address address exists already in one the two tries, false otherwise.
//
// Use GetAddedNode to get the node for the address rather than just checking for its existence.
func (tries *baseDualIPv4v6Tries[V]) Contains(addr *IPAddress) bool {
	return addressFuncOp(addr, tries.ipv4Trie.contains, tries.ipv6Trie.contains)
}

// Remove Removes the given single address or prefix block subnet from the tries.
//
// Removing an element will not remove contained elements (nodes for contained blocks and addresses).
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Returns true if the prefix block or address was removed, false if not already in one of the two tries.
//
// You can also remove by calling GetAddedNode to get the node and then calling Remove on the node.
//
// When an address is removed, the corresponding node may remain in the trie if it remains a subnet block for two sub-nodes.
// If the corresponding node can be removed from the trie, it will be.
func (tries *baseDualIPv4v6Tries[V]) Remove(addr *IPAddress) bool {
	return addressFuncOp(addr, tries.ipv4Trie.remove, tries.ipv6Trie.remove)
}

// ElementContains checks if a prefix block subnet or address in ones of the two tries contains the given subnet or address.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// ElementContains returns true if the subnet or address is contained by a trie element, false otherwise.
//
// To get all the containing addresses, use ElementsContaining.
func (tries *baseDualIPv4v6Tries[V]) ElementContains(addr *IPAddress) bool {
	return addressFuncOp(addr, tries.ipv4Trie.elementContains, tries.ipv6Trie.elementContains)
}

func (tries *baseDualIPv4v6Tries[V]) elementsContaining(addr *IPAddress) *containmentPath[*IPAddress, V] {
	return addressFuncOp(addr, tries.ipv4Trie.elementsContaining, tries.ipv6Trie.elementsContaining)
}

func (tries *baseDualIPv4v6Tries[V]) elementsContainedBy(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.elementsContainedBy, tries.ipv6Trie.elementsContainedBy)
}

func (tries *baseDualIPv4v6Tries[V]) removeElementsContainedBy(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.removeElementsContainedBy, tries.ipv6Trie.removeElementsContainedBy)
}

func (tries *baseDualIPv4v6Tries[V]) getAddedNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.getAddedNode, tries.ipv6Trie.getAddedNode)
}

func (tries *baseDualIPv4v6Tries[V]) longestPrefixMatchNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.longestPrefixMatchNode, tries.ipv6Trie.longestPrefixMatchNode)
}

func (tries *baseDualIPv4v6Tries[V]) LongestPrefixMatch(addr *IPAddress) *IPAddress {
	return addressFuncOp(addr, tries.ipv4Trie.longestPrefixMatch, tries.ipv6Trie.longestPrefixMatch)
}

func (tries *baseDualIPv4v6Tries[V]) shortestPrefixMatchNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.shortestPrefixMatchNode, tries.ipv6Trie.shortestPrefixMatchNode)
}

// ShortestPrefixMatch returns the address added to the trie
// with the same IP version amd the shortest matching prefix compared to the provided address, or nil if no matching address.
func (tries *baseDualIPv4v6Tries[V]) ShortestPrefixMatch(addr *IPAddress) *IPAddress {
	return addressFuncOp(addr, tries.ipv4Trie.shortestPrefixMatch, tries.ipv6Trie.shortestPrefixMatch)
}

func (tries *baseDualIPv4v6Tries[V]) addNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.addNode, tries.ipv6Trie.addNode)
}

func (tries *baseDualIPv4v6Tries[V]) addTrie(trie *trieNode[*IPAddress, V]) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return switchOp(trie.getKey(), trie, tries.ipv4Trie.addTrie, tries.ipv6Trie.addTrie)
}

func (tries *baseDualIPv4v6Tries[V]) floorAddedNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.floorAddedNode, tries.ipv6Trie.floorAddedNode)
}

func (tries *baseDualIPv4v6Tries[V]) lowerAddedNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.lowerAddedNode, tries.ipv6Trie.lowerAddedNode)
}

func (tries *baseDualIPv4v6Tries[V]) ceilingAddedNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.ceilingAddedNode, tries.ipv6Trie.ceilingAddedNode)
}

func (tries *baseDualIPv4v6Tries[V]) higherAddedNode(addr *IPAddress) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	return addressFuncOp(addr, tries.ipv4Trie.higherAddedNode, tries.ipv6Trie.higherAddedNode)
}

// Floor returns the highest address less than or equal to the given address, and with the same version as the given address.
func (tries *baseDualIPv4v6Tries[V]) Floor(addr *IPAddress) *IPAddress {
	return addressFuncOp(addr, tries.ipv4Trie.floor, tries.ipv6Trie.floor)
}

// Lower returns the highest address strictly less than the given address, and with the same version as the given address.
func (tries *baseDualIPv4v6Tries[V]) Lower(addr *IPAddress) *IPAddress {
	return addressFuncOp(addr, tries.ipv4Trie.lower, tries.ipv6Trie.lower)
}

// Ceiling returns the lowest address greater than or equal to the given address, and with the same version as the given address.
func (tries *baseDualIPv4v6Tries[V]) Ceiling(addr *IPAddress) *IPAddress {
	return addressFuncOp(addr, tries.ipv4Trie.ceiling, tries.ipv6Trie.ceiling)
}

// Higher returns the lowest address strictly greater than the given address, and with the same version as the given address.
func (tries *baseDualIPv4v6Tries[V]) Higher(addr *IPAddress) *IPAddress {
	return addressFuncOp(addr, tries.ipv4Trie.higher, tries.ipv6Trie.higher)
}

// Iterator returns an iterator that iterates through the added addresses and prefix blocks in both tries.
// The iteration is in sorted element order, with IPv4 first.
func (tries *baseDualIPv4v6Tries[V]) Iterator() IteratorWithRemove[*IPAddress] {
	return addrTrieIteratorRem[*IPAddress, V]{tries.nodeIterator(true)}
}

// DescendingIterator returns an iterator that iterates through the added addresses and prefix blocks in both tries.
// The iteration is in reverse sorted element order, with IPv6 first.
func (tries *baseDualIPv4v6Tries[V]) DescendingIterator() IteratorWithRemove[*IPAddress] {
	return addrTrieIteratorRem[*IPAddress, V]{tries.nodeIterator(false)}
}

func (tries *baseDualIPv4v6Tries[V]) nodeIterator(forward bool) tree.TrieNodeIteratorRem[trieKey[*IPAddress], V] {
	tries.ensureRoots() // we need a change tracker from each trie to monitor for changes, and each change tracker is created by the trie root, so we need the roots
	ipv4Iterator := tries.ipv4Trie.nodeIterator(forward)
	ipv6Iterator := tries.ipv6Trie.nodeIterator(forward)
	return tree.CombineSequentiallyRem(tries.ipv4Trie.getRoot(), tries.ipv6Trie.getRoot(), ipv4Iterator, ipv6Iterator, forward)
}

func (tries *baseDualIPv4v6Tries[V]) containingFirstIterator(forwardSubNodeOrder bool) tree.TrieNodeIteratorRem[trieKey[*IPAddress], V] {
	tries.ensureRoots() // we need a change tracker from each trie to monitor for changes, and each change tracker is created by the trie root, so we need the roots
	ipv4Iterator := tries.ipv4Trie.containingFirstIterator(forwardSubNodeOrder)
	ipv6Iterator := tries.ipv6Trie.containingFirstIterator(forwardSubNodeOrder)
	return tree.CombineSequentiallyRem[trieKey[*IPAddress], V](tries.ipv4Trie.getRoot(), tries.ipv6Trie.getRoot(), ipv4Iterator, ipv6Iterator, forwardSubNodeOrder)
}

func (tries *baseDualIPv4v6Tries[V]) containedFirstIterator(forwardSubNodeOrder bool) tree.TrieNodeIteratorRem[trieKey[*IPAddress], V] {
	tries.ensureRoots()
	ipv4Iterator := tries.ipv4Trie.containedFirstIterator(forwardSubNodeOrder)
	ipv6Iterator := tries.ipv6Trie.containedFirstIterator(forwardSubNodeOrder)
	return tree.CombineSequentiallyRem(tries.ipv4Trie.getRoot(), tries.ipv6Trie.getRoot(), ipv4Iterator, ipv6Iterator, forwardSubNodeOrder)
}

func (tries *baseDualIPv4v6Tries[V]) blockSizeNodeIterator(lowerSubNodeFirst bool) tree.TrieNodeIteratorRem[trieKey[*IPAddress], V] {
	tries.ensureRoots()
	ipv4Iterator := tries.ipv4Trie.blockSizeNodeIterator(lowerSubNodeFirst)
	ipv6Iterator := tries.ipv6Trie.blockSizeNodeIterator(lowerSubNodeFirst)
	return tree.CombineByBlockSize(tries.ipv4Trie.getRoot(), tries.ipv6Trie.getRoot(), ipv4Iterator, ipv6Iterator, lowerSubNodeFirst)
}

// Format implements the [fmt.Formatter] interface.
func (tries baseDualIPv4v6Tries[V]) Format(state fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		// same call as String()
		str := tree.TreesString(true, tries.ipv4Trie.toTrie(), tries.ipv6Trie.toTrie())
		_, _ = state.Write([]byte(str))
		return
	}
	// We follow the same pattern as for single tries
	s := flagsFromState(state, verb)
	ipv4Str := fmt.Sprintf(s, (*trieBase[*IPAddress, V])(&tries.ipv4Trie))
	ipv6Str := fmt.Sprintf(s, (*trieBase[*IPAddress, V])(&tries.ipv6Trie))
	totalLen := len(ipv4Str) + len(ipv6Str) + 1
	bytes := make([]byte, totalLen+2)
	bytes[0] = '{'
	shifted := bytes[1:]
	copy(shifted, ipv4Str)
	shifted[len(ipv4Str)] = ' '
	shifted = shifted[len(ipv4Str)+1:]
	copy(shifted, ipv6Str)
	shifted[len(ipv6Str)] = '}'
	_, _ = state.Write(bytes)
}

// DualIPv4v6Tries maintains a pair of tries to store both IPv4 and IPv6 addresses and subnets.
type DualIPv4v6Tries struct {
	baseDualIPv4v6Tries[emptyValue]
}

// GetIPv4Trie provides direct access to the IPv4 trie.
func (tries *DualIPv4v6Tries) GetIPv4Trie() *Trie[*IPAddress] {
	tries.ensureRoots() // Since we are making a copy of ipv4Trie, we need to ensure the root exists, otherwise the returned trie will not share the same root
	return &Trie[*IPAddress]{tries.ipv4Trie}
}

// GetIPv6Trie provides direct access to the IPv6 trie.
func (tries *DualIPv4v6Tries) GetIPv6Trie() *Trie[*IPAddress] {
	tries.ensureRoots() // Since we are making a copy of ipv6Trie, we need to ensure the root exists, otherwise the returned trie will not share the same root
	return &Trie[*IPAddress]{tries.ipv6Trie}
}

// Equal returns whether the two given pair of tries is equal to this pair of tries
func (tries *DualIPv4v6Tries) Equal(other *DualIPv4v6Tries) bool {
	return tries.equal(&other.baseDualIPv4v6Tries)
}

// String returns a string representation of the pair of tries.
func (tries *DualIPv4v6Tries) String() string {
	return tries.TreeString(true)
}

// TreeString merges the trie strings of the two tries into a single merged trie string.
func (tries *DualIPv4v6Tries) TreeString(withNonAddedKeys bool) string {
	if tries == nil {
		return nilString()
	}
	return tree.TreesString(withNonAddedKeys, tries.ipv4Trie.toTrie(), tries.ipv6Trie.toTrie())
}

// AddedNodesTreeString provides a string showing a flattened version of the two tries showing only the added nodes and their containment structure, which is non-binary.
func (tries *DualIPv4v6Tries) AddedNodesTreeString() string {
	ipv4AddedTrie := tries.ipv4Trie.constructAddedNodesTree()
	ipv6AddedTrie := tries.ipv6Trie.constructAddedNodesTree()
	return tree.AddedNodesTreesString[trieKey[*IPAddress], emptyValue](&ipv4AddedTrie.trie, &ipv6AddedTrie.trie)
}

// Clone clones this pair of tries.
func (tries *DualIPv4v6Tries) Clone() *DualIPv4v6Tries {
	return &DualIPv4v6Tries{tries.clone()}
}

// ElementsContaining finds the trie nodes in one of the two tries containing the given key and returns them as a linked list.
// Only added nodes are added to the linked list.
//
// If the argument is not a single address nor prefix block, this method will panic.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
func (tries *DualIPv4v6Tries) ElementsContaining(addr *IPAddress) *ContainmentPath[*IPAddress] {
	return &ContainmentPath[*IPAddress]{*tries.elementsContaining(addr)} // *containmentPath[*IPAddress, tree.EmptyValueType]
}

// ElementsContainedBy checks if a part of one of the two tries is contained by the given prefix block subnet or individual address.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Returns the root node of the contained sub-trie, or nil if no sub-trie is contained.
// The node returned need not be an "added" node.
// The returned sub-trie is backed by the containing trie, so changes in this trie are reflected in those nodes and vice-versa.
func (tries *DualIPv4v6Tries) ElementsContainedBy(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.elementsContainedBy(addr))
}

// RemoveElementsContainedBy removes any single address or prefix block subnet from ones of the two tries that is contained in the given individual address or prefix block subnet.
//
// This goes further than Remove, not requiring a match to an inserted node, and also removing all the sub-nodes of any removed node or sub-node.
//
// For example, after inserting 1.2.3.0 and 1.2.3.1, passing 1.2.3.0/31 to RemoveElementsContainedBy will remove them both,
// while the Remove method will remove nothing.
// After inserting 1.2.3.0/31, then Remove will remove 1.2.3.0/31, but will leave 1.2.3.0 and 1.2.3.1 in the trie.
//
// It cannot partially delete a node, such as deleting a single address from a prefix block represented by a node.
// It can only delete the whole node if the whole address or block represented by that node is contained in the given address or block.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Returns the root node of the sub-trie that was removed, or nil if nothing was removed.
func (tries *DualIPv4v6Tries) RemoveElementsContainedBy(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.removeElementsContainedBy(addr))
}

// GetAddedNode gets the trie node corresponding to the added address key.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Use Contains to check for the existence of a given address in the tries.
func (tries *DualIPv4v6Tries) GetAddedNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.getAddedNode(addr))
}

// LongestPrefixMatchNode returns the node of the address of the same version with the longest matching prefix compared to the provided address, or nil if no matching address.
func (tries *DualIPv4v6Tries) LongestPrefixMatchNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.longestPrefixMatchNode(addr))
}

// ShortestPrefixMatch returns the node of the address added to the trie of the same version with the shortest matching prefix compared to the provided address, or nil if no matching address.
func (tries *DualIPv4v6Tries) ShortestPrefixMatchNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.shortestPrefixMatchNode(addr))
}

// AddNode adds the address to this trie.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// The new or existing node for the address is returned.
func (tries *DualIPv4v6Tries) AddNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.addNode(addr))
}

// AddTrie adds nodes for the address keys from the trie with the argument trie root.
// AddTrie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6Tries) AddTrie(trie *TrieNode[*IPAddress]) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.addTrie(trie.toBase()))
}

// AddIPv6Trie adds nodes for the IPv6 address keys from the trie with the argument trie root.
// AddTrie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6Tries) AddIPv6Trie(trie *TrieNode[*IPv6Address]) *TrieNode[*IPAddress] {
	return addTrieToDual(tries, trie)
}

// AddIPv4Trie adds nodes for the IPv4 address keys from the trie with the argument trie root.
// AddTrie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6Tries) AddIPv4Trie(trie *TrieNode[*IPv4Address]) *TrieNode[*IPAddress] {
	return addTrieToDual(tries, trie)
}

func addTrieToDual[R interface {
	TrieKeyConstraint[R]
	ToIP() *IPAddress
}](tries *DualIPv4v6Tries, trie *TrieNode[R]) *TrieNode[*IPAddress] {
	return toAddressTrieNode(addAssociativeTrieToDual(&tries.baseDualIPv4v6Tries, trie.toBase(), func(e emptyValue) emptyValue { return e }))
}

// FloorAddedNode returns the added node whose address is the highest address of the same address version less than or equal to the given address.
func (tries *DualIPv4v6Tries) FloorAddedNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.floorAddedNode(addr))
}

// LowerAddedNode returns the added node whose address is the highest address of the same address version strictly less than the given address.
func (tries *DualIPv4v6Tries) LowerAddedNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.lowerAddedNode(addr))
}

// CeilingAddedNode returns the added node whose address is the lowest address of the same address version greater than or equal to the given address.
func (tries *DualIPv4v6Tries) CeilingAddedNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.ceilingAddedNode(addr))
}

// HigherAddedNode returns the added node whose address is the lowest address of the same address version strictly greater than the given address.
func (tries *DualIPv4v6Tries) HigherAddedNode(addr *IPAddress) *TrieNode[*IPAddress] {
	return toAddressTrieNode(tries.higherAddedNode(addr))
}

// NodeIterator returns an iterator that iterates through all the added nodes in the two tries in forward or reverse tree order.
// IPv4 comes first in forward order, IPv6 first in reverse order.
func (tries *DualIPv4v6Tries) NodeIterator(forward bool) IteratorWithRemove[*TrieNode[*IPAddress]] {
	return addrTrieNodeIteratorRem[*IPAddress, emptyValue]{tries.nodeIterator(forward)}
}

// ContainingFirstIterator returns an iterator that does a pre-order binary trie traversal of the added nodes.
// All added nodes will be visited before their added sub-nodes.
// For an address trie this means added containing subnet blocks will be visited before their added contained addresses and subnet blocks.
func (tries *DualIPv4v6Tries) ContainingFirstIterator(forwardSubNodeOrder bool) IteratorWithRemove[*TrieNode[*IPAddress]] {
	return addrTrieNodeIteratorRem[*IPAddress, emptyValue]{tries.containingFirstIterator(forwardSubNodeOrder)}
}

// ContainedFirstIterator returns an iterator that does a post-order binary trie traversal of the added nodes.
// All added sub-nodes will be visited before their parent nodes.
// For an address trie this means contained addresses and subnets will be visited before their containing subnet blocks.
func (tries *DualIPv4v6Tries) ContainedFirstIterator(forwardSubNodeOrder bool) IteratorWithRemove[*TrieNode[*IPAddress]] {
	return addrTrieNodeIteratorRem[*IPAddress, emptyValue]{tries.containedFirstIterator(forwardSubNodeOrder)}
}

// BlockSizeNodeIterator returns an iterator that iterates the added nodes in the two tries,
// ordered by keys from largest prefix blocks to smallest, and then to individual addresses.
//
// If lowerSubNodeFirst is true, for blocks of equal size the lower is first, otherwise the reverse order
func (tries *DualIPv4v6Tries) BlockSizeNodeIterator(lowerSubNodeFirst bool) IteratorWithRemove[*TrieNode[*IPAddress]] {
	return addrTrieNodeIteratorRem[*IPAddress, emptyValue]{tries.blockSizeNodeIterator(lowerSubNodeFirst)}
}

//
//
//
//
//

// DualIPv4v6AssociativeTries maintains a pair of associative tries to map both IPv4 and IPv6 addresses and subnets to values of the value type V.
type DualIPv4v6AssociativeTries[V any] struct {
	baseDualIPv4v6Tries[V]
}

// GetIPv4Trie provides direct access to the IPv4 associative trie.
func (tries *DualIPv4v6AssociativeTries[V]) GetIPv4Trie() *AssociativeTrie[*IPAddress, V] {
	tries.ensureRoots() // Since we are making a copy of ipv4Trie, we need to ensure the root exists, otherwise the returned trie will not share the same root
	return &AssociativeTrie[*IPAddress, V]{tries.ipv4Trie}
}

// GetIPv6Trie provides direct access to the IPv6 associative trie.
func (tries *DualIPv4v6AssociativeTries[V]) GetIPv6Trie() *AssociativeTrie[*IPAddress, V] {
	tries.ensureRoots() // Since we are making a copy of ipv6Trie, we need to ensure the root exists, otherwise the returned trie will not share the same root
	return &AssociativeTrie[*IPAddress, V]{tries.ipv6Trie}
}

// Equal returns whether the two given pair of tries is equal to this pair of tries
func (tries *DualIPv4v6AssociativeTries[V]) Equal(other *DualIPv4v6AssociativeTries[V]) bool {
	return tries.equal(&other.baseDualIPv4v6Tries)
}

// DeepEqual returns whether the given argument is a trie with a set of nodes with the same keys as in this trie according to the Compare method,
// and the same values according to the reflect.DeepEqual method
func (tries *DualIPv4v6AssociativeTries[V]) DeepEqual(other *DualIPv4v6AssociativeTries[V]) bool {
	return tries.ipv4Trie.deepEqual(&other.ipv4Trie) && tries.ipv6Trie.deepEqual(&other.ipv6Trie)
}

// String returns a string representation of the pair of tries.
func (tries *DualIPv4v6AssociativeTries[V]) String() string {
	return tries.TreeString(true)
}

// TreeString merges the trie strings of the two tries into a single merged trie string.
func (tries *DualIPv4v6AssociativeTries[V]) TreeString(withNonAddedKeys bool) string {
	if tries == nil {
		return nilString()
	}
	return tree.TreesString(withNonAddedKeys, tries.ipv4Trie.toTrie(), tries.ipv6Trie.toTrie())
}

// AddedNodesTreeString provides a string showing a flattened version of the two tries showing only the added nodes and their containment structure, which is non-binary.
func (tries *DualIPv4v6AssociativeTries[V]) AddedNodesTreeString() string {
	ipv4AddedTrie := tries.ipv4Trie.constructAddedNodesTree()
	ipv6AddedTrie := tries.ipv6Trie.constructAddedNodesTree()
	return tree.AddedNodesTreesString[trieKey[*IPAddress], V](&ipv4AddedTrie.trie, &ipv6AddedTrie.trie)
}

// Clone clones this pair of tries.
func (tries *DualIPv4v6AssociativeTries[V]) Clone() *DualIPv4v6AssociativeTries[V] {
	return &DualIPv4v6AssociativeTries[V]{tries.clone()}
}

// ElementsContaining finds the trie nodes in one of the two tries containing the given key and returns them as a linked list.
// Only added nodes are added to the linked list.
//
// If the argument is not a single address nor prefix block, this method will panic.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
func (tries *DualIPv4v6AssociativeTries[V]) ElementsContaining(addr *IPAddress) *ContainmentValuesPath[*IPAddress, V] {
	return &ContainmentValuesPath[*IPAddress, V]{*tries.elementsContaining(addr)}
}

// ElementsContainedBy checks if a part of one of the two tries is contained by the given prefix block subnet or individual address.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Returns the root node of the contained sub-trie, or nil if no sub-trie is contained.
// The node returned need not be an "added" node.
// The returned sub-trie is backed by the containing trie, so changes in this trie are reflected in those nodes and vice-versa.
func (tries *DualIPv4v6AssociativeTries[V]) ElementsContainedBy(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.elementsContainedBy(addr))
}

// RemoveElementsContainedBy removes any single address or prefix block subnet from ones of the two tries that is contained in the given individual address or prefix block subnet.
//
// This goes further than Remove, not requiring a match to an inserted node, and also removing all the sub-nodes of any removed node or sub-node.
//
// For example, after inserting 1.2.3.0 and 1.2.3.1, passing 1.2.3.0/31 to RemoveElementsContainedBy will remove them both,
// while the Remove method will remove nothing.
// After inserting 1.2.3.0/31, then Remove will remove 1.2.3.0/31, but will leave 1.2.3.0 and 1.2.3.1 in the trie.
//
// It cannot partially delete a node, such as deleting a single address from a prefix block represented by a node.
// It can only delete the whole node if the whole address or block represented by that node is contained in the given address or block.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Returns the root node of the sub-trie that was removed, or nil if nothing was removed.
func (tries *DualIPv4v6AssociativeTries[V]) RemoveElementsContainedBy(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.removeElementsContainedBy(addr))
}

// GetAddedNode gets the associative trie node corresponding to the added address key.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Use Contains to check for the existence of a given address in the tries.
func (tries *DualIPv4v6AssociativeTries[V]) GetAddedNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.getAddedNode(addr))
}

// LongestPrefixMatchNode returns the node of the address of the same version with the longest matching prefix compared to the provided address.
func (tries *DualIPv4v6AssociativeTries[V]) LongestPrefixMatchNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.longestPrefixMatchNode(addr))
}

// ShortestPrefixMatch returns the node of the address added to the trie of the same version with the shortest matching prefix compared to the provided address, or nil if no matching address.
func (tries *DualIPv4v6AssociativeTries[V]) ShortestPrefixMatchNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.shortestPrefixMatchNode(addr))
}

// AddNode adds the address to this trie.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// The new or existing node for the address is returned.
func (tries *DualIPv4v6AssociativeTries[V]) AddNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.addNode(addr))
}

// AddTrie adds nodes for the address keys from the trie with the argument trie root.
// All the new nodes will have values that are zero values.  To add both keys and values, use PutTrie.
// AddTrie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6AssociativeTries[V]) AddTrie(trie *AssociativeTrieNode[*IPAddress, V]) *AssociativeTrieNode[*IPAddress, V] {
	return AddAssociativeTrieToDual(tries, trie, nil)
}

// AddIPv6Trie adds nodes for the IPv6 address keys from the trie with the argument trie root.
// All the new nodes will have values that are zero values.  To add both keys and values, use PutIPv6Trie.
// AddIPv6Trie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6AssociativeTries[V]) AddIPv6Trie(trie *AssociativeTrieNode[*IPv6Address, V]) *AssociativeTrieNode[*IPAddress, V] {
	return AddAssociativeTrieToDual(tries, trie, nil)
}

// AddIPv4Trie adds nodes for the IPv4 address keys from the trie with the argument trie root.
// All the new nodes will have values that are zero values.  To add both keys and values, use PutIPv4Trie.
// AddIPv4Trie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6AssociativeTries[V]) AddIPv4Trie(trie *AssociativeTrieNode[*IPv4Address, V]) *AssociativeTrieNode[*IPAddress, V] {
	return AddAssociativeTrieToDual(tries, trie, nil)
}

// AddAssociativeTrie adds the given trie's entries to this trie.  The given trie's keys must have a ToIP() method to be convertible to *IPAddress, like *IPV4Address or *IPv6Address.
// If withValues is true, the values will be mapped with the given valueMap mapping.  If valueMap is nil, then all values will be napped to the V2 zero value.
// If withValues is false, then valueMap is ignored and can be nil.
// The given trie can map to different value types.  You must supply a function to map from the given trie's values to this trie's values.
// If you are using the same value type, then you can use DualIPv4v6AssociativeTries[V].AddIPv4Trie or DualIPv4v6AssociativeTries[V].AddIPv6Trie instead.
func AddAssociativeTrieToDual[R interface {
	TrieKeyConstraint[R]
	ToIP() *IPAddress
}, V, V2 any](tries *DualIPv4v6AssociativeTries[V], trie *AssociativeTrieNode[R, V2], valueMap func(v V2) V) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(addAssociativeTrieToDual(&tries.baseDualIPv4v6Tries, trie.toBase(), valueMap))
}

func addAssociativeTrieToDual[R interface {
	TrieKeyConstraint[R]
	ToIP() *IPAddress
}, V, V2 any](tries *baseDualIPv4v6Tries[V], trie *trieNode[R, V2], valueMap func(v V2) V) *tree.BinTrieNode[trieKey[*IPAddress], V] {
	if trie == nil {
		return nil
	}
	var targetTrie *trieBase[*IPAddress, V]
	rootKey := trie.getKey().ToIP()
	if rootKey.IsIPv4() {
		targetTrie = &tries.ipv4Trie
	} else if rootKey.IsIPv6() {
		targetTrie = &tries.ipv6Trie
	}
	if targetTrie != nil {
		return tree.AddConvertibleTrie(
			&targetTrie.trie,
			trie.toBinTrieNode(),
			func(r trieKey[R]) trieKey[*IPAddress] { return trieKey[*IPAddress]{r.address.ToIP()} },
			valueMap)
	}
	return nil
}

// FloorAddedNode returns the added node whose address is the highest address of the same address version less than or equal to the given address.
func (tries *DualIPv4v6AssociativeTries[V]) FloorAddedNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.floorAddedNode(addr))
}

// LowerAddedNode returns the added node whose address is the highest address of the same address version strictly less than the given address.
func (tries *DualIPv4v6AssociativeTries[V]) LowerAddedNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.lowerAddedNode(addr))
}

// CeilingAddedNode returns the added node whose address is the lowest address of the same address version greater than or equal to the given address.
func (tries *DualIPv4v6AssociativeTries[V]) CeilingAddedNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.ceilingAddedNode(addr))
}

// HigherAddedNode returns the added node whose address is the lowest address of the same address version strictly greater than the given address.
func (tries *DualIPv4v6AssociativeTries[V]) HigherAddedNode(addr *IPAddress) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(tries.higherAddedNode(addr))
}

// NodeIterator returns an iterator that iterates through all the added nodes in the two tries in forward or reverse tree order.
// IPv4 comes first in forward order, IPv6 first in reverse order.
func (tries *DualIPv4v6AssociativeTries[V]) NodeIterator(forward bool) IteratorWithRemove[*AssociativeTrieNode[*IPAddress, V]] {
	return associativeAddressTrieNodeIteratorRem[*IPAddress, V]{tries.nodeIterator(forward)}
}

// ContainingFirstIterator returns an iterator that does a pre-order binary trie traversal of the added nodes.
// All added nodes will be visited before their added sub-nodes.
// For an address trie this means added containing subnet blocks will be visited before their added contained addresses and subnet blocks.
func (tries *DualIPv4v6AssociativeTries[V]) ContainingFirstIterator(forwardSubNodeOrder bool) IteratorWithRemove[*AssociativeTrieNode[*IPAddress, V]] {
	return associativeAddressTrieNodeIteratorRem[*IPAddress, V]{tries.containingFirstIterator(forwardSubNodeOrder)}
}

// ContainedFirstIterator returns an iterator that does a post-order binary trie traversal of the added nodes.
// All added sub-nodes will be visited before their parent nodes.
// For an address trie this means contained addresses and subnets will be visited before their containing subnet blocks.
func (tries *DualIPv4v6AssociativeTries[V]) ContainedFirstIterator(forwardSubNodeOrder bool) IteratorWithRemove[*AssociativeTrieNode[*IPAddress, V]] {
	return associativeAddressTrieNodeIteratorRem[*IPAddress, V]{tries.containedFirstIterator(forwardSubNodeOrder)}
}

// BlockSizeNodeIterator returns an iterator that iterates the added nodes in the two tries,
// ordered by keys from largest prefix blocks to smallest, and then to individual addresses.
//
// If lowerSubNodeFirst is true, for blocks of equal size the lower is first, otherwise the reverse order
func (tries *DualIPv4v6AssociativeTries[V]) BlockSizeNodeIterator(lowerSubNodeFirst bool) IteratorWithRemove[*AssociativeTrieNode[*IPAddress, V]] {
	return associativeAddressTrieNodeIteratorRem[*IPAddress, V]{tries.blockSizeNodeIterator(lowerSubNodeFirst)}
}

// Get gets the value for the specified address key in the associative trie matching the given argument address version.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Returns the value for the given key.
// Returns nil if neither tries contains a mapping for that address key or if the mapped value is nil.
func (tries *DualIPv4v6AssociativeTries[V]) Get(addr *IPAddress) (V, bool) {
	return addressFuncDoubRetOp(addr, tries.ipv4Trie.get, tries.ipv6Trie.get)
}

// Put associates the specified value with the specified address key in the corresponding trie.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// If the corresponding trie previously contained a mapping for the key,
// the old value is replaced by the specified value, and false is returned along with the old value.
// If the corresponding trie did not previously contain a mapping for the key, true is returned along with a zero value.
// The boolean return value allows you to distinguish whether the address was previously mapped to the zero value or not mapped at all.
func (tries *DualIPv4v6AssociativeTries[V]) Put(addr *IPAddress, value V) (V, bool) {
	return addressFuncDoubArgDoubleRetOp(addr, value, tries.ipv4Trie.put, tries.ipv6Trie.put)
}

// PutNode associates the specified value with the specified key in this map.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
//
// Returns the node for the added address, whether it was already in the trie or not.
//
// If you wish to know whether the node was already there when adding, use PutNew, or before adding you can use GetAddedNode.
func (tries *DualIPv4v6AssociativeTries[V]) PutNode(addr *IPAddress, value V) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(addressFuncDoubArgOp(addr, value, tries.ipv4Trie.putNode, tries.ipv6Trie.putNode))
}

// PutTrie adds nodes with the address keys and values from the trie with the argument trie root.  To add only the keys, use AddTrie.
//
// For each added node from the given trie that does not exist in the associated trie of the same address version, a copy will be made,
// the copy including the node value, and the copy will be inserted into the associated trie.
//
// To add nodes without the associated values, use AddTrie.
//
// The address type/version of the keys must match.
//
// When adding one trie to another, this method is more efficient than adding each node of the first trie individually.
// When using this method, searching for the location to add sub-nodes starts from the inserted parent node.
//
// Returns the node corresponding to the given sub-root node, whether it was already in the trie or not.
func (tries *DualIPv4v6AssociativeTries[V]) PutTrie(trie *AssociativeTrieNode[*IPAddress, V]) *AssociativeTrieNode[*IPAddress, V] {
	return AddAssociativeTrieToDual(tries, trie, func(v V) V { return v })
}

// PutIPv4Trie adds nodes with the IPv4 address keys and values from the trie with the argument trie root.
//
// For each added node from the given trie that does not exist in the associated trie of the same address version, a copy will be made,
// the copy including the node value, and the copy will be inserted into the associated trie.
//
// To add nodes without the associated values, use AddIPv4Trie.
//
// PutIPv4Trie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6AssociativeTries[V]) PutIPv4Trie(trie *AssociativeTrieNode[*IPv4Address, V]) *AssociativeTrieNode[*IPAddress, V] {
	return AddAssociativeTrieToDual(tries, trie, func(v V) V { return v })
}

// PutIPv6Trie adds nodes with the IPv6 address keys and values from the trie with the argument trie root.
//
// For each added node from the given trie that does not exist in the associated trie of the same address version, a copy will be made,
// the copy including the node value, and the copy will be inserted into the associated trie.
//
// To add nodes without the associated values, use AddIPv6Trie.
//
// PutIPv6Trie returns the sub-node in the trie where the added trie begins, where the first node of the added trie is located.
func (tries *DualIPv4v6AssociativeTries[V]) PutIPv6Trie(trie *AssociativeTrieNode[*IPv6Address, V]) *AssociativeTrieNode[*IPAddress, V] {
	return AddAssociativeTrieToDual(tries, trie, func(v V) V { return v })
}

// Remap remaps node values in the two tries.
//
// This will look up the node corresponding to the given key.
// It will call the remapping function, regardless of whether the node is found or not.
//
// If the node is not found, or the node is not an "added" node, the existingValue argument will be the zero value.
// If the node is found, the existingValue argument will be the node's value, which can also be the zero value.
// The boolean "found" argument will be true if the node was found and it is an "added" node.
// If the node was not found or was not an "added" node, then the boolean "found" argument will be false.
//
// If the remapping function returns false as the "mapIt" argument, then the matched node will be removed or converted to a "non-added" node, if any.
// If it returns true, then either the existing node will be set to an "added" node with the "mapped" value given as the first argument,
// or if there was no matched node, it will create a new added node with the "mapped" value.
//
// The method will return the node involved, which is either the matched node, or the newly created node,
// or nil if there was no matched node nor newly created node.
//
// If the remapping function modifies the trie during its computation,
// and the returned values from the remapper requires changes to be made to the trie,
// then the trie will not be changed as required by the remapper, and Remap will panic.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
func (tries *DualIPv4v6AssociativeTries[V]) Remap(addr *IPAddress, remapper func(existingValue V, found bool) (mapped V, mapIt bool)) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(addressFuncDoubArgOp(addr, remapper, tries.ipv4Trie.remap, tries.ipv6Trie.remap))
}

// RemapIfAbsent remaps node values in the two tries, but only for nodes that do not exist or are not "added".
//
// This will look up the node corresponding to the given key.
// If the node is not found or not "added", then RemapIfAbsent will call the supplier function.
// It will create a new node with the value returned from the supplier function.
// If the node is found and "added", then RemapIfAbsent will not call the supplier function.
//
// The method will return the node involved, which is either the matched node, the newly created node, or nil if there was no matched node nor newly created node.
//
// If the supplier function modifies the trie during its computation,
// then the trie will not be changed and RemapIfAbsent will panic.
//
// If the argument is not a single address nor prefix block, this method will panic.
// The [Partition] type can be used to convert the argument to single addresses and prefix blocks before calling this method.
func (tries *DualIPv4v6AssociativeTries[V]) RemapIfAbsent(addr *IPAddress, supplier func() V) *AssociativeTrieNode[*IPAddress, V] {
	return toAssociativeTrieNode(addressFuncDoubArgOp(addr, supplier, tries.ipv4Trie.remapIfAbsent, tries.ipv6Trie.remapIfAbsent))
}

func addressFuncOp[T any](addr *IPAddress, ipv4Op, ipv6Op func(*IPAddress) T) T {
	return switchOp(addr, addr, ipv4Op, ipv6Op)
}

func switchOp[R interface {
	IsIPv4() bool
	IsIPv6() bool
}, S, T any](key R, arg S, ipv4Op, ipv6Op func(S) T) T {
	if key.IsIPv4() {
		return ipv4Op(arg)
	} else if key.IsIPv6() {
		return ipv6Op(arg)
	}
	var t T
	return t
}

// Get
func addressFuncDoubRetOp[T1, T2 any](addr *IPAddress, ipv4Op, ipv6Op func(*IPAddress) (T1, T2)) (T1, T2) {
	return switchDoubleRetOp(addr, addr, ipv4Op, ipv6Op)
}

func switchDoubleRetOp[R interface {
	IsIPv4() bool
	IsIPv6() bool
}, S, T1, T2 any](key R, arg S, ipv4Op, ipv6Op func(S) (T1, T2)) (T1, T2) {
	if key.IsIPv4() {
		return ipv4Op(arg)
	} else if key.IsIPv6() {
		return ipv6Op(arg)
	}
	var t1 T1
	var t2 T2
	return t1, t2
}

// PutNode, Remap, RemapIfAbsent
func addressFuncDoubArgOp[S, T any](addr *IPAddress, arg S, ipv4Op, ipv6Op func(*IPAddress, S) T) T {
	return switchDoubleArgOp(addr, addr, arg, ipv4Op, ipv6Op)
}

func switchDoubleArgOp[R interface {
	IsIPv4() bool
	IsIPv6() bool
}, S1, S2, T any](key R, arg1 S1, arg2 S2, ipv4Op, ipv6Op func(S1, S2) T) T {
	if key.IsIPv4() {
		return ipv4Op(arg1, arg2)
	} else if key.IsIPv6() {
		return ipv6Op(arg1, arg2)
	}
	var t T
	return t
}

// Put
func addressFuncDoubArgDoubleRetOp[S, T1, T2 any](addr *IPAddress, arg S, ipv4Op, ipv6Op func(*IPAddress, S) (T1, T2)) (T1, T2) {
	return switchDoubleArgDoubleRetOp(addr, addr, arg, ipv4Op, ipv6Op)
}

func switchDoubleArgDoubleRetOp[R interface {
	IsIPv4() bool
	IsIPv6() bool
}, S1, S2, T1, T2 any](key R, arg1 S1, arg2 S2, ipv4Op, ipv6Op func(S1, S2) (T1, T2)) (T1, T2) {
	if key.IsIPv4() {
		return ipv4Op(arg1, arg2)
	} else if key.IsIPv6() {
		return ipv6Op(arg1, arg2)
	}
	var t1 T1
	var t2 T2
	return t1, t2
}
