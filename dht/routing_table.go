/* Copyright 2019 zigma authors
 * This file is part of the zigma library.
 *
 * The zigma library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The zigma library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with the zigma library. If not, see <http://www.gnu.org/licenses/>.
 */

package dht

import "sync"

const (
	// Alpha is the number representing the degree of parallelism in network calls
	Alpha = 3
	// B is the size in bits of the keys used to identify nodes and store and retrieve data;
	// this is normally 160 bit with the use of sha1 hash type, but we are using ed25519
	// public key and sha-256 hash so we're passing in 256 here
	B = 256
	// K is the maximum number of contacts stored in a bucket; this is normally 20
	K = 20
)

// RoutingTable implements the routing table state
type RoutingTable struct {
	mutex   *sync.RWMutex
	Self    *Node
	Buckets [B]*Bucket
}

func (h *RoutingTable) bucketFromNode(node *Node) *Bucket {
	idx := h.Self.ZeroPrefixLen(node)
	return h.Buckets[idx]
}

func (h *RoutingTable) kclosest(num int, contact *Node, ignoredNodes ...*Node) []*Node {
	var (
		l = NewContacts(h.Self, ignoredNodes...)
		d = h.Self.ZeroPrefixLen(contact)
		m = K
	)
	if num > 0 {
		m = num
	}
	for i := d; i >= 0 && l.Len() < m; i-- {
		for node := range h.Buckets[i].Iterator() {
			if l.Append(node) && l.Len() >= m {
				break
			}
		}
	}
	for i := d + 1; i < B && l.Len() < m; i++ {
		for node := range h.Buckets[i].Iterator() {
			if l.Append(node) && l.Len() >= m {
				break
			}
		}
	}
	l.Sort()
	return l.Nodes
}

func (h *RoutingTable) addNode(node *Node) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if node == nil || len(node.Id) == 0 {
		return
	}

	bucket := h.bucketFromNode(node)
	bucket.AddNode(node)
}

func (h *RoutingTable) removeNode(node *Node) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if node == nil || len(node.Id) == 0 {
		return
	}

	bucket := h.bucketFromNode(node)
	bucket.RemoveNode(node)
}

func (h *RoutingTable) size() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var n int
	for i := 0; i < len(h.Buckets); i++ {
		n += h.Buckets[i].Len()
	}
	return n
}

// NewRoutingTable initializes a new hashtable instance
func NewRoutingTable(self *Node) *RoutingTable {
	h := &RoutingTable{
		mutex:   new(sync.RWMutex),
		Self:    self,
		Buckets: [B]*Bucket{},
	}
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < len(h.Buckets); i++ {
		h.Buckets[i] = NewBucket()
	}
	return h
}
