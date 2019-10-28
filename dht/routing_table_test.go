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
	"encoding/binary"
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/dht"
)

var table *dht.RoutingTable

func TestNewRoutingTable(t *testing.T) {
	table = dht.NewRoutingTable(dht.MockNode(-1))

	assert.NotNil(t, table)
	assert.NotNil(t, table.Self)
	assert.Len(t, table.Buckets, len(table.Self.Hash)*8)
}

func TestRoutingTableAddNode(t *testing.T) {
	l := 10000
	for i := 0; i < l; i++ {
		node := dht.MockNode(i)
		table.AddNode(node)
	}
	assert.True(t, table.Size() > 0)
}

func TestRoutingTableRemoveNode(t *testing.T) {
	l := 10000
	for i := 0; i < l; i++ {
		node := dht.MockNode(i)
		table.RemoveNode(node)
	}
	assert.Zero(t, table.Size())
}

func TestRoutingTableKclosest(t *testing.T) {
	ct1 := dht.MockNode(-2)
	ct2 := dht.MockNode(-3)

	table.AddNode(dht.MockNode(0))
	nodes := table.Kclosest(0, ct1)
	assert.Len(t, nodes, 1)
	node1 := nodes[0]

	table.AddNode(dht.MockNode(1))
	nodes = table.Kclosest(0, ct1)
	assert.Len(t, nodes, 2)
	assert.Equal(t, node1, nodes[0])

	table.AddNode(dht.MockNode(2))
	nodes = table.Kclosest(0, ct1)
	assert.Len(t, nodes, 3)
	assert.NotEqual(t, node1, nodes[0])

	for i := 3; i < 40; i++ {
		table.AddNode(dht.MockNode(i))
	}
	nodes = table.Kclosest(10, ct1)
	assert.Len(t, nodes, 10)
	nodes = table.Kclosest(0, ct1)
	assert.Len(t, nodes, 20)

	other := table.Kclosest(0, ct2)
	assert.Len(t, other, 20)

	var diff bool
	for i := 0; i < len(other); i++ {
		if !nodes[i].Equal(other[i]) {
			diff = true
			break
		}
	}
	assert.True(t, diff)

	other = table.Kclosest(0, ct1, nodes...)
	assert.NotEqual(t, len(other), len(nodes))
	assert.NotZero(t, len(other))

	for i := 40; i < 1000; i++ {
		table.AddNode(dht.MockNode(i))
	}
	other = table.Kclosest(0, ct1, nodes...)
	assert.Len(t, other, 20)
}

func TestRoutingTableBucketCap(t *testing.T) {
	b1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(b1, uint32(0))
	h1, _ := multihash.Sum(b1, multihash.SHAKE_128, -1)

	b2 := make([]byte, 4)
	binary.LittleEndian.PutUint32(b2, uint32(0))
	h2, _ := multihash.Sum(b2, multihash.SHA3_512, -1)

	table = dht.NewRoutingTable(dht.NodeFromHash(h1))
	assert.Len(t, table.Buckets, 256)

	table.AddNode(dht.NodeFromHash(h2))
	assert.Len(t, table.Buckets, 512)
}
