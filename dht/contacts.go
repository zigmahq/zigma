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
	"sort"
	"sync"
)

// the key length cap is set to 64+2 bytes. at the time of writing, this length is
// more than enough as the default hash function is set to sha3-256, which the output
// length is  32+2 bytes. but just a note that the length might change in the future
// on hashing algorithm upgrade. capacity: 2 bytes + (number of digest bits / 8)
// <varint hash function code><varint digest size in bytes><hash function output>
const cc = 66

// Contacts is used in order to sort a list of arbitrary nodes against a comparator
type Contacts struct {
	mutex      *sync.RWMutex
	uniq       map[[cc]byte]int
	blacklist  map[[cc]byte]struct{}
	nodes      []*Node
	Comparator *Node
}

// nodeID reads a {cc} bytes identifier key from node id
// the key is composed with the first {cc} bytes from node id only
func (c *Contacts) nodeID(node *Node) [cc]byte {
	var out [cc]byte
	l := len(node.Id)
	if n := len(out); l > n {
		l = n
	}
	for i := 0; i < l; i++ {
		out[i] = node.Id[i]
	}
	return out
}

// Append adds node to node list
func (c *Contacts) Append(node *Node) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	id := c.nodeID(node)
	if _, ok := c.blacklist[id]; ok {
		return false
	}
	if _, ok := c.uniq[id]; ok {
		return false
	}

	idx := len(c.nodes)
	c.nodes = append(c.nodes, node)
	c.uniq[id] = idx
	return true
}

// Remove removes a node from the node list
func (c *Contacts) Remove(node *Node) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	id := c.nodeID(node)
	i, ok := c.uniq[id]
	if !ok {
		return false
	}

	c.nodes = append(c.nodes[:i], c.nodes[i+1:]...)
	delete(c.uniq, id)

	for k, v := range c.uniq {
		if v > i {
			c.uniq[k]--
		}
	}
	return true
}

// IndexOf returns the index of node in node list
func (c *Contacts) IndexOf(node *Node) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	id := c.nodeID(node)
	if i, ok := c.uniq[id]; ok {
		return i
	}
	return -1
}

// Sort sorts nodes based on distance in node list
func (c *Contacts) Sort() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Len() > 1 {
		sort.Sort(c)
	}
}

// Nodes returns all nodes in contacts
func (c *Contacts) Nodes() []*Node {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	tmp := make([]*Node, len(c.nodes))
	for i := 0; i < len(c.nodes); i++ {
		tmp[i] = c.nodes[i]
	}
	return tmp
}

// Len implements slice sorting interface
func (c *Contacts) Len() int {
	return len(c.nodes)
}

// Swap implements slice sorting interface
func (c *Contacts) Swap(i, j int) {
	i1 := c.nodeID(c.nodes[i])
	i2 := c.nodeID(c.nodes[j])
	c.nodes[i], c.nodes[j] = c.nodes[j], c.nodes[i]
	c.uniq[i1], c.uniq[i2] = c.uniq[i2], c.uniq[i1]
}

// Less implements slice sorting interface
func (c *Contacts) Less(i, j int) bool {
	if c.Comparator != nil {
		d1 := c.nodes[i].DistanceBetween(c.Comparator)
		d2 := c.nodes[j].DistanceBetween(c.Comparator)
		return d1.Cmp(d2) == -1
	}
	return false
}

// NewContacts initializes a Contacts
func NewContacts(comparator *Node, blacklist ...*Node) *Contacts {
	c := &Contacts{
		mutex:      new(sync.RWMutex),
		uniq:       make(map[[cc]byte]int),
		blacklist:  make(map[[cc]byte]struct{}),
		nodes:      make([]*Node, 0),
		Comparator: comparator,
	}
	for _, node := range blacklist {
		id := c.nodeID(node)
		c.blacklist[id] = struct{}{}
	}
	return c
}
