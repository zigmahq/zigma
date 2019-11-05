package dht

import (
	"time"

	"github.com/multiformats/go-multihash"
	"github.com/zigmahq/zigma/store"
)

var (
	prefixStoreReplication = []byte{0x72, 0x21}
	prefixStoreData        = []byte{0x64, 0x21}
)

// KademliaStore extends store.Store key-value storage
type KademliaStore struct {
	store.Store
	replicating bool
}

func (s *KademliaStore) dataKey(key []byte) []byte {
	return append(prefixStoreData, key...)
}

func (s *KademliaStore) replicationKey(key []byte) []byte {
	return append(prefixStoreReplication, key...)
}

// Get retrieves a key-value pair from storage
func (s *KademliaStore) Get(key []byte) (data []byte, found bool) {
	return s.Store.Get(s.dataKey(key))
}

// Set insert key value pair to storage
func (s *KademliaStore) Set(key, val []byte, ttl time.Duration) {
	s.Store.Set(s.dataKey(key), val, ttl)
	if b, err := time.Now().UTC().MarshalBinary(); err == nil {
		s.Store.Set(s.replicationKey(key), b, ttl)
	}
}

// Delete removes a key-value pair from storage
func (s *KademliaStore) Delete(key []byte) {
	s.Store.Delete(s.dataKey(key))
	s.Store.Delete(s.replicationKey(key))
}

// Iterate iterates key-value pairs existed in storage
func (s *KademliaStore) Iterate(key []byte) store.Iterator {
	return s.Store.Iterate(s.dataKey(key))
}

// PendingReplication returns pending replication items
func (s *KademliaStore) PendingReplication() <-chan Hashable {
	var ch = make(chan Hashable)
	go func() {
		defer close(ch)
		if s.replicating {
			return
		}
		s.replicating = true
		defer func() { s.replicating = false }()

		iter := s.Store.Iterate(prefixStoreReplication)
		defer iter.Done()

		for iter.Next() {
			var last time.Time
			item := iter.Item()
			last.UnmarshalBinary(item.Value())
			if !last.IsZero() && tReplicate > time.Since(last) {
				continue
			}

			dkey := item.Key()[len(prefixStoreReplication):]
			data, ok := s.Get(dkey)
			if !ok {
				continue
			}
			hash, err := multihash.Sum(data, h, -1)
			if err != nil {
				continue
			}
			ch <- &hashable{
				key:  dkey,
				data: data,
				hash: hash,
			}
		}
	}()
	return ch
}

// NewKademliaStore initializes kademlia store
func NewKademliaStore(store store.Store) *KademliaStore {
	return &KademliaStore{
		Store:       store,
		replicating: false,
	}
}
