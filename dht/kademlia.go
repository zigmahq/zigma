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
	"time"

	"github.com/multiformats/go-multihash"
	"github.com/zigmahq/zigma/store"
)

const (
	// the default hashing function
	h = multihash.SHA3_256

	// the maximum number of contacts stored in a bucket; this is normally 20
	k = 20

	// the number representing the degree of parallelism in network calls
	alpha = 3

	// the time after which a key/value pair expires; this is a time-to-live (TTL)
	// from the original publication date; this is normally 86400s
	tExpire = time.Hour * 24

	// the time after which an otherwise unaccessed bucket must be refreshed
	tRefresh = time.Hour

	// the interval between kademlia replication events, when a node is required to
	// publish its entire database
	tReplicate = time.Hour

	// the time after which the original publisher must republish a key/value pair
	tRepublish = time.Hour * 24
)

// Kademlia represents the state of the local node in the distributed hash table
type Kademlia struct {
	table *RoutingTable
	store store.Store
}

// Store stores data on the network. A sha-256 encoded identifier will be returned
// if the store operation is successful
func (kad *Kademlia) Store(data Data) ([]byte, error) {
	key := data.Key()
	if len(key) == 0 {
	}
	return nil, nil
}

// FindValue retrieves data from the network with a key
func (kad *Kademlia) FindValue(key []byte) ([]byte, error) {
	return nil, nil
}

// FindNode returns a node from the networking using key
func (kad *Kademlia) FindNode(key []byte) (*Node, error) {
	return nil, nil
}

// NewKademlia initializes a DHT kademlia service
func NewKademlia(self *Node, store store.Store) *Kademlia {
	r := NewRoutingTable(self)
	k := &Kademlia{
		table: r,
		store: store,
	}
	return k
}
