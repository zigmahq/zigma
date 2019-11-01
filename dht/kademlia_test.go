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
	"github.com/zigmahq/zigma/log"
	"github.com/zigmahq/zigma/store"
)

var (
	kad1, kad2 *dht.Kademlia
	n1, n2     *dht.Node
	db1, db2   store.Store
)

func init() {
	log.SetLevel(log.LogWarn)
	db1 = store.TempBadgerStore()
	db2 = store.TempBadgerStore()
}

func TestNewKademlia(t *testing.T) {
	n1 = dht.MockNode(0)
	n2 = dht.MockNode(1)
	r1 := dht.MockRPC(n1, true)
	r2 := dht.MockRPC(n2)

	kad1 = dht.NewKademlia(n1, db1, r1)
	kad2 = dht.NewKademlia(n2, db2, r2)
	assert.NotNil(t, kad1)
	assert.NotNil(t, kad2)
}

func TestKademliaPing(t *testing.T) {
	ok := kad1.Ping(n2)
	assert.True(t, ok)

	ok = kad2.Ping(n1)
	assert.True(t, ok)
}

func TestKademliaStore(t *testing.T) {
	kad1.Store(dht.String("hello world 1"))
	kad1.Store(dht.Bytes([]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}))
}

func TestKademliaFindValue(t *testing.T) {
}

func TestKademliaFindNode(t *testing.T) {
	defer db1.Close()
	defer db2.Close()
}
