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
	"unsafe"
)

// Contacts is used in order to sort a list of arbitrary nodes against a comparator
type Contacts struct {
	mutex      *sync.RWMutex
	uniq       map[string]int
	blacklist  map[string]struct{}
	Nodes      []*Node
	Comparator *Node
}

func (p *Contacts) nodeID(node *Node) string {
	return *(*string)(unsafe.Pointer(&node.Id))
}

// Append adds node to node list
func (p *Contacts) Append(node *Node) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	id := p.nodeID(node)
	if _, ok := p.blacklist[id]; ok {
		return false
	}
	if _, ok := p.uniq[id]; ok {
		return false
	}

	p.Nodes = append(p.Nodes, node)
	p.uniq[id] = len(p.Nodes) - 1
	return true
}

// Remove removes a node from the node list
func (p *Contacts) Remove(node *Node) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	id := p.nodeID(node)
	i, ok := p.uniq[id]
	if !ok {
		return false
	}

	p.Nodes = append(p.Nodes[:i], p.Nodes[i+1:]...)
	delete(p.uniq, id)
	for k, v := range p.uniq {
		if v > i {
			p.uniq[k]--
		}
	}
	return true
}

// IndexOf returns the index of node in node list
func (p *Contacts) IndexOf(node *Node) int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	id := p.nodeID(node)
	if i, ok := p.uniq[id]; ok {
		return i
	}
	return -1
}

// Sort sorts nodes based on distance in node list
func (p *Contacts) Sort() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.Len() > 1 {
		sort.Sort(p)
	}
}

// Len implements slice sorting interface
func (p *Contacts) Len() int {
	return len(p.Nodes)
}

// Swap implements slice sorting interface
func (p *Contacts) Swap(i, j int) {
	i1 := p.nodeID(p.Nodes[i])
	i2 := p.nodeID(p.Nodes[j])
	p.Nodes[i], p.Nodes[j] = p.Nodes[j], p.Nodes[i]
	p.uniq[i1], p.uniq[i2] = p.uniq[i2], p.uniq[i1]
}

// Less implements slice sorting interface
func (p *Contacts) Less(i, j int) bool {
	if p.Comparator != nil {
		d1 := p.Nodes[i].DistanceBetween(p.Comparator)
		d2 := p.Nodes[j].DistanceBetween(p.Comparator)
		return d1.Cmp(d2) == -1
	}
	return false
}

// NewContacts initializes a Contacts
func NewContacts(comparator *Node, blacklist ...*Node) *Contacts {
	pl := &Contacts{
		mutex:      new(sync.RWMutex),
		uniq:       make(map[string]int),
		blacklist:  make(map[string]struct{}),
		Nodes:      make([]*Node, 0),
		Comparator: comparator,
	}
	for _, node := range blacklist {
		id := pl.nodeID(node)
		pl.blacklist[id] = struct{}{}
	}
	return pl
}
