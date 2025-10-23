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
	"strconv"
	"strings"
)

// Path is a list of nodes derived from following a path in a tree.
// Each node in the list corresponds to a node in the tree.
// Each node in the list corresponds to a tree node that is a direct or indirect sub-node of the tree node corresponding to the previous node in the list.
// Not all nodes in the pathway through the tree need to be included in the linked list.
//
// In other words, a path follows a pathway through a tree from root to leaf, but not necessarily including all nodes encountered along the way.
type Path[E Key, V any] struct {
	root, leaf *PathNode[E, V]
}

// GetRoot returns the beginning of the Path, which may or may not match the tree root of the originating tree.
func (path *Path[E, V]) GetRoot() *PathNode[E, V] {
	return path.root
}

// GetLeaf returns the end of the Path, which may or may not match a leaf in the originating tree.
func (path *Path[E, V]) GetLeaf() *PathNode[E, V] {
	return path.leaf
}

// String returns a visual representation of the Path with one node per line.
func (path *Path[E, V]) String() string {
	return path.ListString(true)
}

// ListString returns a visual representation of the tree with one node per line, with or without the non-added keys.
func (path *Path[E, V]) ListString(withNonAddedKeys bool) string {
	if path.Size() == 0 {
		builder := strings.Builder{}
		builder.WriteByte('\n')
		if path == nil || !withNonAddedKeys {
			builder.WriteString(nilString())
		} else {
			builder.WriteString(nonAddedNodeCircle)
		}
		builder.WriteByte('\n')
		return builder.String()
	}
	return path.GetRoot().ListString(withNonAddedKeys, true)
}

// Size returns the count of nodes in the list
// This is a constant-time operation since the size is maintained in each node.
func (path *Path[E, V]) Size() (storedSize int) {
	if path == nil || path.root == nil {
		return 0
	}
	return path.root.Size()
}

// PathNode is an element in the list of a Path
type PathNode[E Key, V any] struct {
	previous, next *PathNode[E, V]

	// the key for the node
	item E

	// only for associative trie nodes
	value V

	// the number of added nodes below this one, including this one if added
	storedSize int

	added bool
}

// Next returns the next node in the path
func (node *PathNode[E, V]) Next() *PathNode[E, V] {
	return node.next
}

// Previous returns the previous node in the path
func (node *PathNode[E, V]) Previous() *PathNode[E, V] {
	return node.previous
}

// GetKey gets the key used for placing the node in the tree.
func (node *PathNode[E, V]) GetKey() (key E) {
	if node != nil {
		return node.item
	}
	return
}

// GetValue returns the value assigned to the node
func (node *PathNode[E, V]) GetValue() (val V) {
	if node != nil {
		val = node.value
	}
	return
}

// IsAdded returns whether the node was "added".
// Some binary tree nodes are considered "added" and others are not.
// Those nodes created for key elements added to the tree are "added" nodes.
// Those that are not added are those nodes created to serve as junctions for the added nodes.
// Only added elements contribute to the size of a tree.
// When removing nodes, non-added nodes are removed automatically whenever they are no longer needed,
// which is when an added node has less than two added sub-nodes.
func (node *PathNode[E, V]) IsAdded() bool {
	return node != nil && node.added
}

// Size returns the count of nodes added to the sub-tree starting from this node as root and moving downwards to sub-nodes.
// This is a constant-time operation since the size is maintained in each node.
func (node *PathNode[E, V]) Size() (storedSize int) {
	if node != nil {
		storedSize = node.storedSize
		if storedSize == sizeUnknown {
			prev, next := node, node.next
			for ; next != nil && next.storedSize == sizeUnknown; prev, next = next, next.next {
			}
			var nodeSize int
			if next != nil {
				nodeSize = next.storedSize
			}
			for {
				if prev.IsAdded() {
					nodeSize++
				}
				prev.storedSize = nodeSize
				if prev == node {
					break
				}
				prev = prev.previous
			}
			storedSize = node.storedSize
		}
	}
	return
}

// Returns a visual representation of this node including the key, with an open circle indicating this node is not an added node,
// a closed circle indicating this node is an added node.
func (node *PathNode[E, V]) String() string {
	if node == nil {
		return NodeString[E, V](nil)
	}
	return NodeString[E, V](node)
}

// ListString returns a visual representation of the sub-list with this node as root, with one node per line.
//
// withNonAddedKeys: whether to show nodes that are not added nodes
// withSizes: whether to include the counts of added nodes in each sub-list
func (node *PathNode[E, V]) ListString(withNonAddedKeys, withSizes bool) string {
	builder := strings.Builder{}
	builder.WriteByte('\n')
	node.printList(&builder, indents{}, withNonAddedKeys, withSizes)
	return builder.String()
}

func (node *PathNode[E, V]) printList(builder *strings.Builder,
	indents indents,
	withNonAdded,
	withSizes bool) {
	if node == nil {
		builder.WriteString(indents.nodeIndent)
		builder.WriteString(nilString())
		builder.WriteByte('\n')
		return
	}
	next := node
	for {
		if withNonAdded || next.IsAdded() {
			builder.WriteString(indents.nodeIndent)
			builder.WriteString(next.String())
			if withSizes {
				builder.WriteString(" (")
				builder.WriteString(strconv.Itoa(next.Size()))
				builder.WriteByte(')')
			}
			builder.WriteByte('\n')
		} else {
			builder.WriteString(indents.nodeIndent)
			builder.WriteString(nonAddedNodeCircle)
			builder.WriteByte('\n')
		}
		if next = next.next; next == nil {
			break
		}
		indents.nodeIndent = indents.subNodeInd + rightElbow
		indents.subNodeInd += belowElbows
	}
}
