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

import (
	"sync"
	"time"
)

// RoutingTable implements the routing table state
type RoutingTable struct {
	mutex   *sync.RWMutex
	b       int
	refresh []time.Time
	Self    *Node
	Buckets []*Bucket
}

func (r *RoutingTable) bucketFromNode(node *Node) *Bucket {
	idx := r.Self.ZeroPrefixLen(node)
	return r.Buckets[idx]
}

// Kclosest searches the routing table, and returns N number of closest node with contact id
func (r *RoutingTable) Kclosest(num int, contact *Node, ignoredNodes ...*Node) []*Node {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

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
	for i := d + 1; i < r.b && l.Len() < n; i++ {
		for node := range r.Buckets[i].Iterator() {
			if l.Append(node) && l.Len() >= n {
				break
			}
		}
	}
	l.Sort()
	return l.Nodes()
}

// Update insert a node to routing table
func (r *RoutingTable) Update(node *Node) {
	if !IsValidNode(node) {
		return
	}
	r.shouldUpdateBucketCap(node)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	bucket := r.bucketFromNode(node)
	bucket.Update(node)
}

// Remove removes a node from routing table
func (r *RoutingTable) Remove(node *Node) {
	if !IsValidNode(node) {
		return
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	bucket := r.bucketFromNode(node)
	bucket.Remove(node)
}

// Size returns the total number of nodes in routing table
func (r *RoutingTable) Size() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var n int
	for i := 0; i < len(r.Buckets); i++ {
		n += r.Buckets[i].Len()
	}
	return n
}

// BucketsNeededForRefresh returns a list of bucket index that needed for refresh
func (r *RoutingTable) BucketsNeededForRefresh() <-chan int {
	ch := make(chan int, len(r.refresh))
	go func() {
		r.mutex.RLock()
		defer r.mutex.RUnlock()
		defer close(ch)

		for i := 0; i < len(r.refresh); i++ {
			if r.shouldBucketRefresh(i) {
				ch <- i
			}
		}
	}()
	return ch
}

// MarkBucketRefreshed marks a bucket refreshed
func (r *RoutingTable) MarkBucketRefreshed(idx int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if idx > len(r.refresh)-1 {
		return
	}
	r.refresh[idx] = time.Now()
}

// RandomNodeFromBucket returns a random node picked from specified bucket
func (r *RoutingTable) RandomNodeFromBucket(idx int) *Node {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if idx > len(r.Buckets)-1 {
		return nil
	}
	bucket := r.Buckets[idx]
	for node := range bucket.Iterator() {
		return node
	}
	return nil
}

func (r *RoutingTable) shouldBucketRefresh(idx int) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if idx > len(r.refresh)-1 {
		return false
	}
	at := r.refresh[idx]
	return at.IsZero() || time.Since(at) > tRefresh
}

func (r *RoutingTable) shouldUpdateBucketCap(node *Node) {
	if b := len(node.Hash) * 8; b > r.b {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		t := make([]*Bucket, b-r.b)
		o := make([]time.Time, b-r.b)
		n := time.Now()
		for i := 0; i < len(t); i++ {
			t[i] = NewBucket()
			o[i] = n
		}
		r.Buckets = append(r.Buckets, t...)
		r.refresh = append(r.refresh, o...)
		r.b = b
	}
}

// NewRoutingTable initializes a new hashtable instance
func NewRoutingTable(self *Node) *RoutingTable {
	b := len(self.Hash) * 8
	n := time.Now()
	r := &RoutingTable{
		mutex:   new(sync.RWMutex),
		b:       b,
		refresh: make([]time.Time, b),
		Self:    self,
		Buckets: make([]*Bucket, b),
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := 0; i < len(r.Buckets); i++ {
		r.refresh[i] = n
		r.Buckets[i] = NewBucket()
	}
	return r
}
