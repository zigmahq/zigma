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
	h = multihash.SHA3_512

	// the maximum number of contacts stored in a bucket; this is normally 20
	k = 20

	// the number representing the degree of parallelism in network calls
	a = 3

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
	rpc   KademliaRPC
}

// KademliaRPC represents the rpc interface for kademlia dht server
type KademliaRPC interface {
	Write(*Message)
	Read() <-chan *Message
}

// Store stores data on the network. A sha-256 encoded identifier will be returned
// if the store operation is successful
func (kad *Kademlia) Store(data Hashable) error {
	msg := compose(kad.table.Self).store(data)
	kad.rpc.Write(msg)
	return nil
}

// FindValue retrieves data from the network with a key
func (kad *Kademlia) FindValue(key []byte) ([]byte, error) {
	msg := compose(kad.table.Self).findValue(key)
	kad.rpc.Write(msg)
	return nil, nil
}

// FindNode returns a node from the networking using key
func (kad *Kademlia) FindNode(key []byte) (*Node, error) {
	msg := compose(kad.table.Self).findNode(key)
	kad.rpc.Write(msg)
	return nil, nil
}

func (kad *Kademlia) listen() {
	for {
		select {
		case msg := <-kad.rpc.Read():
			switch msg.Type {
			// PING RPC involves one node sending a PING message to another, which presumably replies with a PONG.
			case MessageType_PING:
				kad.rpc.Write(msg.pong())

			// The sender of the STORE RPC provides a key and a block of data and requires that the recipient
			// store the data and make it available for later retrieval by that key.
			case MessageType_STORE:
				payload := msg.GetStore().Payload
				kad.table.Update(msg.Sender)
				kad.store.Set(payload.Key, payload.Data, 0)

			// FIND_VALUE returns the associated data if corresponding value is present. Otherwise the RPC
			// is equivalent to a FIND_NODE and a set of k triples is returned.
			case MessageType_FIND_VALUE:
				kad.table.Update(msg.Sender)
				if b, ok := kad.store.Get(msg.GetFind().Key); ok {
					kad.rpc.Write(msg.returnValue(b))
				} else {
					nodes := kad.table.Kclosest(k, &Node{Hash: msg.GetFind().Key}, msg.Sender)
					kad.rpc.Write(msg.returnClosest(nodes))
				}

			// FIND_NODE returns up to k triples for the contacts that it knows to be closest to the key
			case MessageType_FIND_NODE:
				kad.table.Update(msg.Sender)
				nodes := kad.table.Kclosest(k, &Node{Hash: msg.GetFind().Key}, msg.Sender)
				kad.rpc.Write(msg.returnClosest(nodes))
			}
		}
	}
}

// NewKademlia initializes a DHT kademlia service
func NewKademlia(self *Node, store store.Store, rpc KademliaRPC) *Kademlia {
	r := NewRoutingTable(self)
	k := &Kademlia{
		table: r,
		store: store,
		rpc:   rpc,
	}
	go k.listen()
	return k
}
