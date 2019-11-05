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
	tExpire = time.Hour * 48

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
	rpc   KademliaRPC
	stop  chan struct{}
	store *KademliaStore
	table *RoutingTable
}

// KademliaReplyFn represents the wait-for-response function for KademliaRPC, passing
// a time.Duration as waiting timeout; default reply timeout is half a second
type KademliaReplyFn func(time.Duration) <-chan *Message

// KademliaRPC represents the network rpc interface for kademlia dht server
type KademliaRPC interface {
	Write(*Message) KademliaReplyFn
	Read() <-chan *Message
}

// Bootstrap adds seed nodes to the network
func (kad *Kademlia) Bootstrap(seeds ...*Node) {
	for _, seed := range seeds {
		go func(node *Node) {
			if kad.Ping(node) {
				kad.table.Update(node)
			}
		}(seed)
	}
}

// Stop stops the kademlia server
func (kad *Kademlia) Stop() {
	kad.stop <- struct{}{}
}

// Table returns the dht network routing table
func (kad *Kademlia) Table() *RoutingTable {
	return kad.table
}

// Ping the specified contact node; returns true if pong is returned from receiver
func (kad *Kademlia) Ping(node *Node) bool {
	msg := compose(kad.table.Self).to(node).ping()
	rec := kad.rpc.Write(msg)
	switch out := <-rec(0); {
	case out != nil:
		kad.table.Update(node)
		return true
	default:
		kad.table.Remove(node)
		return false
	}
}

// Store stores data on the network. A sha-256 encoded identifier will be returned
// if the store operation is successful
func (kad *Kademlia) Store(data Hashable) ([]byte, int) {
	if data != nil {
		if writes := kad.iterativeStore(data); writes > 0 {
			return data.Key(), writes
		}
	}
	return nil, 0
}

// FindNode returns a node from the networking using key
// The recipient of a the RPC returns k nodes it knows about closest to the target
// ID. These triples can come from a single k-bucket, or they may come from multiple
// k-buckets if the closest k-bucket is not full. In any case, the RPC recipient
// must return k items (unless there are fewer than k nodes in all its k-buckets
// combined, in which case it returns every node it knows about).
func (kad *Kademlia) FindNode(key []byte) *Node {
	if contacts := kad.iterativeFindNode(key); contacts.Len() > 0 {
		return contacts.Nodes()[0]
	}
	return nil
}

// FindValue retrieves data from the network with a key
func (kad *Kademlia) FindValue(key []byte) ([]byte, bool) {
	return kad.iterativeFindValue(key)
}

func (kad *Kademlia) iterativeStore(data Hashable) int {
	contacts := kad.iterativeFindNode(data.Hash())
	if contacts.Len() == 0 {
		return 0
	}
	var (
		c, i int
		ch   = make(chan bool)
	)
	for _, node := range contacts.Nodes() {
		go func(node *Node) {
			msg := compose(kad.table.Self).to(node).store(data)
			rec := kad.rpc.Write(msg)
			out := <-rec(0)
			ch <- out != nil && out.GetSuccess()
		}(node)
	}
	for success := range ch {
		if i == contacts.Len()-1 {
			close(ch)
			break
		}
		if success {
			c++
		}
		i++
	}
	return c
}

func (kad *Kademlia) iterativeFindNode(key []byte) *Contacts {
	contacts := NewContacts(kad.table.Self)
	nodes := kad.table.Kclosest(a, &Node{Hash: key})
	if len(nodes) == 0 {
		return nil
	}
	for _, node := range nodes {
		msg := compose(kad.table.Self).to(node).findNode(key)
		rec := kad.rpc.Write(msg)
		switch out := <-rec(0); {
		case out != nil:
			for _, node := range out.GetClosest().GetNodes() {
				contacts.Append(node)
			}
			kad.table.Update(node)
		default:
			kad.table.Remove(node)
		}
	}
	contacts.Sort()
	return contacts
}

func (kad *Kademlia) iterativeFindValue(key []byte) ([]byte, bool) {
	contacts := kad.iterativeFindNode(key)
	var i int
	for {
		if i == contacts.Len()-1 {
			break
		}
		node := contacts.Nodes()[i]
		msg := compose(kad.table.Self).to(node).findValue(key)
		rec := kad.rpc.Write(msg)
		switch out := <-rec(0); {
		case out != nil && out.GetPayload() != nil:
			return out.GetPayload().GetData(), true
		case out != nil && out.GetClosest() != nil:
			for _, node := range out.GetClosest().GetNodes() {
				contacts.Append(node)
			}
		}
		i++
	}
	return nil, false
}

func (kad *Kademlia) listen() {
	for {
		select {
		case <-kad.stop:
			return
		case msg := <-kad.rpc.Read():
			if msg == nil || msg.IsResponse || !msg.isValid() {
				continue
			}
			switch msg.Type {

			// PING RPC involves one node sending a PING message to another,
			// which presumably replies with a PONG.
			case MessageType_PING:
				kad.rpc.Write(msg.pong())

			// STORE RPC provides a key and a block of data and requires that the
			// recipient store the data and make it available for later retrieval
			// by that key.
			case MessageType_STORE:
				payload := msg.GetStore().Payload
				kad.table.Update(msg.Sender)
				kad.store.Set(payload.Key, payload.Data, tExpire)
				kad.rpc.Write(msg.success(true))

			// FIND_VALUE returns the associated data if corresponding value is
			// present. Otherwise the RPC is equivalent to a FIND_NODE and a set
			// of k triples is returned.
			case MessageType_FIND_VALUE:
				kad.table.Update(msg.Sender)
				if b, ok := kad.store.Get(msg.GetFind().Key); ok {
					kad.rpc.Write(msg.returnValue(b))
				} else {
					nodes := kad.table.Kclosest(k, &Node{Hash: msg.GetFind().Key}, msg.Sender)
					kad.rpc.Write(msg.returnClosest(nodes))
				}

			// FIND_NODE returns up to k triples for the contacts that it knows
			// to be closest to the key
			case MessageType_FIND_NODE:
				kad.table.Update(msg.Sender)
				nodes := kad.table.Kclosest(k, &Node{Hash: msg.GetFind().Key}, msg.Sender)
				kad.rpc.Write(msg.returnClosest(nodes))
			}
		}
	}
}

func (kad *Kademlia) refreshBuckets() {
	for idx := range kad.table.BucketsNeededForRefresh() {
		if node := kad.table.RandomNodeFromBucket(idx); node != nil {
			if kad.iterativeFindNode(node.Id).Len() > 0 {
				kad.table.MarkBucketRefreshed(idx)
			}
		}
	}
}

func (kad *Kademlia) replicaDatabase() {
	for hashable := range kad.store.PendingReplication() {
		kad.Store(hashable)
	}
}

func (kad *Kademlia) scheduleTasks() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			go kad.replicaDatabase()
			go kad.refreshBuckets()
		case <-kad.stop:
			ticker.Stop()
			return
		}
	}
}

// NewKademlia initializes a DHT kademlia service
func NewKademlia(self *Node, store store.Store, rpc KademliaRPC) *Kademlia {
	s := make(chan struct{})
	t := NewRoutingTable(self)
	r := NewKademliaStore(store)
	k := &Kademlia{
		rpc:   rpc,
		stop:  s,
		store: r,
		table: t,
	}
	go k.listen()
	go k.scheduleTasks()
	return k
}
