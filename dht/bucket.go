// Copyright 2019 zigma authors
// This file is part of the zigma library.
//
// The zigma library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The zigma library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the zigma library. If not, see <http://www.gnu.org/licenses/>.

package dht

import (
	"strings"
	"sync"
)

// Bucket implements the hashtable bucket
type Bucket struct {
	mutex *sync.RWMutex
	nodes [k]*Node
}

// Update adds a node to bucket
// Nodes within buckets are sorted by least recently seen e.g.
// [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  ^                                                           ^
//  └ Least recently seen                    Most recently seen ┘
func (b *Bucket) Update(node *Node) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !IsValidNode(node) {
		return
	}
	if idx := b.indexOf(node); idx > -1 {
		b.markSeen(idx)
		return
	}
	for i := 0; i < k-1; i++ {
		b.nodes[i] = b.nodes[i+1]
	}
	b.nodes[k-1] = node
}

// Remove removes a node from the bucket
// [a][b][c][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
// [ ][a][c][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//     ^
//     └ Remove node, then right pad the nodes on the left
func (b *Bucket) Remove(node *Node) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !IsValidNode(node) {
		return
	}
	if idx := b.indexOf(node); idx > -1 {
		var l = idx
		for i := idx; i > 0; i-- {
			if b.nodes[i] != nil {
				l = i
			}
			b.nodes[i] = b.nodes[i-1]
		}
		b.nodes[l] = nil
	}
}

// RemoveAll removes all nodes from bucket
func (b *Bucket) RemoveAll() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i := 0; i < k; i++ {
		b.nodes[i] = nil
	}
}

// Iterator iterate over active nodes in the bucket
func (b *Bucket) Iterator() <-chan *Node {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	ch := make(chan *Node, k)
	go func() {
		defer close(ch)
		for i := k - 1; i >= 0; i-- {
			if b.nodes[i] == nil {
				return
			}
			ch <- b.nodes[i]
		}
	}()
	return ch
}

// At returns the node at bucket index
func (b *Bucket) At(idx int) *Node {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.nodes[idx]
}

// Len calculates the node size of the bucket
func (b *Bucket) Len() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var total int
	for i := k - 1; i >= 0; i-- {
		if b.nodes[i] == nil {
			break
		}
		total++
	}
	return total
}

// Cap returns the maximum number of nodes that it could store
func (b *Bucket) Cap() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return len(b.nodes)
}

func (b *Bucket) String() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var sb strings.Builder
	var i int
	sb.WriteByte(0x5b)
	for node := range b.Iterator() {
		if i > 0 {
			sb.WriteByte(0x2c)
			sb.WriteByte(0x20)
		}
		s := node.HexString()
		if l := len(s); l > 8 {
			sb.Write(s[:8])
		} else {
			sb.Write(s[:])
		}
		i++
	}
	sb.WriteByte(0x5d)
	return sb.String()
}

func (b *Bucket) markSeen(idx int) {
	t := b.nodes[idx]
	for i := idx; i < k-1; i++ {
		b.nodes[i] = b.nodes[i+1]
	}
	b.nodes[k-1] = t
}

func (b *Bucket) indexOf(node *Node) int {
	for i := k - 1; i >= 0; i-- {
		if b.nodes[i] == nil {
			break
		}
		if b.nodes[i].Equal(node) {
			return i
		}
	}
	return -1
}

// NewBucket initializes a bucket instance
func NewBucket() *Bucket {
	return &Bucket{
		mutex: new(sync.RWMutex),
		nodes: [k]*Node{},
	}
}
