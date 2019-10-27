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

// RoutingTable implements the routing table state
type RoutingTable struct {
	mutex   *sync.RWMutex
	Self    *Node
	Buckets [b]*Bucket
}

func (r *RoutingTable) bucketFromNode(node *Node) *Bucket {
	idx := r.Self.ZeroPrefixLen(node)
	return r.Buckets[idx]
}

func (r *RoutingTable) kclosest(num int, contact *Node, ignoredNodes ...*Node) []*Node {
	var (
		l = NewContacts(r.Self, ignoredNodes...)
		d = r.Self.ZeroPrefixLen(contact)
		n = k
	)
	if num > 0 {
		n = num
	}
	for i := d; i >= 0 && l.Len() < n; i-- {
		for node := range r.Buckets[i].Iterator() {
			if l.Append(node) && l.Len() >= n {
				break
			}
		}
	}
	for i := d + 1; i < b && l.Len() < n; i++ {
		for node := range r.Buckets[i].Iterator() {
			if l.Append(node) && l.Len() >= n {
				break
			}
		}
	}
	l.Sort()
	return l.Nodes
}

func (r *RoutingTable) addNode(node *Node) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if node == nil || len(node.Id) == 0 {
		return
	}

	bucket := r.bucketFromNode(node)
	bucket.AddNode(node)
}

func (r *RoutingTable) removeNode(node *Node) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if node == nil || len(node.Id) == 0 {
		return
	}

	bucket := r.bucketFromNode(node)
	bucket.RemoveNode(node)
}

func (r *RoutingTable) size() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var n int
	for i := 0; i < len(r.Buckets); i++ {
		n += r.Buckets[i].Len()
	}
	return n
}

// NewRoutingTable initializes a new hashtable instance
func NewRoutingTable(self *Node) *RoutingTable {
	r := &RoutingTable{
		mutex:   new(sync.RWMutex),
		Self:    self,
		Buckets: [b]*Bucket{},
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := 0; i < len(r.Buckets); i++ {
		r.Buckets[i] = NewBucket()
	}
	return r
}