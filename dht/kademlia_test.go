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
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/dht"
	"github.com/zigmahq/zigma/log"
	"github.com/zigmahq/zigma/store"
)

var (
	n         = 80
	kadList   = make([]*dht.Kademlia, n)
	rpcList   = make([]dht.KademliaRPC, n)
	nodeList  = make([]*dht.Node, n)
	storeList = make([]store.Store, n)
)

func init() {
	log.SetLevel(log.LogWarn)
}

func done() {
	for i := 0; i < n; i++ {
		db := storeList[i]
		kad := kadList[i]
		db.Close()
		kad.Stop()
	}
}

func TestNewKademlia(t *testing.T) {
	for i := 0; i < n; i++ {
		var (
			db   = store.TempBadgerStore()
			node = dht.MockNode(i)
			rpc  = dht.MockRPC(node, i == 0)
			kad  = dht.NewKademlia(node, db, rpc)
		)
		storeList[i] = db
		nodeList[i] = node
		rpcList[i] = rpc
		kadList[i] = kad

		assert.NotNil(t, db)
		assert.NotNil(t, node)
		assert.NotNil(t, rpc)
		assert.NotNil(t, kad)
	}
	for i := 0; i < n; i++ {
		kad := kadList[i]
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			node := nodeList[j]
			kad.Bootstrap(node)
		}
	}
}

func TestKademliaPing(t *testing.T) {
	ok := kadList[0].Ping(nodeList[1])
	assert.True(t, ok)

	ok = kadList[1].Ping(nodeList[2])
	assert.True(t, ok)
}

func TestKademliaStore(t *testing.T) {
	for i := 0; i < n; i++ {
		r := rand.Intn(n - 1)

		hs := dht.String(fmt.Sprintf("hello world %v", i))
		key, success := kadList[r].Store(hs)
		assert.NotEmpty(t, key)
		assert.True(t, success > 0)
	}
}

func TestKademliaFindValue(t *testing.T) {
	for i := 0; i < n; i++ {
		r := rand.Intn(n - 1)

		hs := dht.String(fmt.Sprintf("hello world %v", i))
		b, _ := kadList[r].FindValue(hs.Hash())
		assert.NotNil(t, b)
	}
}

func TestKademliaFindNode(t *testing.T) {
	defer done()
	hs := dht.String("hello world")

	node := kadList[0].FindNode(hs.Hash())
	assert.NotNil(t, node)
}
