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

package dht_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/dht"
)

var contacts *dht.Contacts

func init() {
	self := dht.MockNode(-1)
	contacts = dht.NewContacts(self)
}

func TestContactsAppend(t *testing.T) {
	var tmp = make([]*dht.Node, 1000)
	for i := 0; i < len(tmp); i++ {
		tmp[i] = dht.MockNode(i)
	}
	for _, node := range tmp {
		pl := contacts.Len()
		ok := contacts.Append(node)
		al := contacts.Len()
		assert.True(t, ok)
		assert.Equal(t, pl+1, al)
	}
	pl := contacts.Len()
	for _, node := range tmp {
		ok := contacts.Append(node)
		al := contacts.Len()
		assert.False(t, ok)
		assert.Equal(t, pl, al)
	}
}

func TestContactsRemove(t *testing.T) {
	nodes := contacts.Nodes()
	assert.Len(t, nodes, contacts.Len())
	assert.True(t, len(nodes) > 0)
	for _, node := range nodes {
		ps := contacts.IndexOf(node)
		assert.Equal(t, 0, ps)

		pl := contacts.Len()
		ok := contacts.Remove(node)
		al := contacts.Len()
		assert.True(t, ok)
		assert.Equal(t, pl-1, al)
	}
	pl := contacts.Len()
	for _, node := range nodes {
		ps := contacts.IndexOf(node)
		assert.Equal(t, -1, ps)

		ok := contacts.Remove(node)
		al := contacts.Len()
		assert.False(t, ok)
		assert.Equal(t, pl, al)
	}
}

func TestContactsSort(t *testing.T) {
	for i := 0; i < 1000; i++ {
		node := dht.MockNode(i)
		contacts.Append(node)
	}
	prev := contacts.Nodes()
	contacts.Sort()
	next := contacts.Nodes()
	assert.Equal(t, len(prev), len(next))

	var changed bool
	for i := 0; i < len(prev); i++ {
		if !prev[i].Equal(next[i]) {
			changed = true
			break
		}
	}
	assert.True(t, changed)
}
