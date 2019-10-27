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

import "strings"

// Bucket implements the hashtable bucket
type Bucket [K]*Node

// AddNode adds a node to bucket
// Nodes within buckets are sorted by least recently seen e.g.
// [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  ^                                                           ^
//  └ Least recently seen                    Most recently seen ┘
func (b *Bucket) AddNode(node *Node) {
	if node == nil || len(node.Id) == 0 {
		return
	}
	if idx := b.indexOf(node); idx > -1 {
		b.markSeen(idx)
		return
	}
	for i := 0; i < K-1; i++ {
		b[i] = b[i+1]
	}
	b[K-1] = node
}

// RemoveNode removes a node from the bucket
// [a][b][c][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
// [ ][a][c][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//     ^
//     └ Remove node, then right pad the nodes on the left
func (b *Bucket) RemoveNode(node *Node) {
	if node == nil || len(node.Id) == 0 {
		return
	}
	if idx := b.indexOf(node); idx > -1 {
		var l = idx
		for i := idx; i > 0; i-- {
			if b[i] != nil {
				l = i
			}
			b[i] = b[i-1]
		}
		b[l] = nil
	}
}

// RemoveAllNodes removes all nodes from bucket
func (b *Bucket) RemoveAllNodes() {
	for i := 0; i < K; i++ {
		b[i] = nil
	}
}

// Iterator iterate over active nodes in the bucket
func (b *Bucket) Iterator() <-chan *Node {
	ch := make(chan *Node, K)
	go func() {
		defer close(ch)
		for i := K - 1; i >= 0; i-- {
			if b[i] == nil {
				return
			}
			ch <- b[i]
		}
	}()
	return ch
}

// Len calculates the node size of the bucket
func (b *Bucket) Len() int {
	var total int
	for i := K - 1; i >= 0; i-- {
		if b[i] == nil {
			break
		}
		total++
	}
	return total
}

func (b *Bucket) String() string {
	var sb strings.Builder
	var i int
	sb.WriteByte(0x5b)
	for node := range b.Iterator() {
		if i > 0 {
			sb.WriteByte(0x2c)
			sb.WriteByte(0x20)
		}
		sb.Write(node.HexString()[:B/8/4])
		i++
	}
	sb.WriteByte(0x5d)
	return sb.String()
}

func (b *Bucket) markSeen(idx int) {
	t := b[idx]
	for i := idx; i < K-1; i++ {
		b[i] = b[i+1]
	}
	b[K-1] = t
}

func (b *Bucket) indexOf(node *Node) int {
	for i := K - 1; i >= 0; i-- {
		if b[i] == nil {
			break
		}
		if b[i].Equal(node) {
			return i
		}
	}
	return -1
}

// NewBucket initializes a bucket instance
func NewBucket() *Bucket {
	return new(Bucket)
}
