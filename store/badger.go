/* Copyright 2019 zigma authors
 * This file is part of the zigma library.
 *
 * The zigma library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The zigma library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with the zigma library. If not, see <http://www.gnu.org/licenses/>.
 */

package store

import (
	"time"

	"github.com/dgraph-io/badger"
)

// BadgerStore implements dgraph-io/badger storage
type BadgerStore struct {
	db   *badger.DB
	quit chan struct{}
}

// Init initializes the garbage collection task
func (b *BadgerStore) Init() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				b.db.RunValueLogGC(0.7)
			case <-b.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Set sets value to a badger key/value storage, passing in negative or 0 as
// expiration number would not set an expiration time to key
func (b *BadgerStore) Set(key, val []byte, expiration time.Duration) {
	txn := b.db.NewTransaction(true)
	defer txn.Commit()

	entry := badger.NewEntry(key, val)
	if expiration > 0 {
		entry.WithTTL(expiration)
	}
	_ = txn.SetEntry(entry)
}

// Get retrieves value from a badger storage
func (b *BadgerStore) Get(key []byte) ([]byte, bool) {
	txn := b.db.NewTransaction(false)
	item, err := txn.Get(key)
	if err != nil {
		return nil, false
	}
	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, false
	}
	return data, true
}

// Delete removes an existing key from badger storage
func (b *BadgerStore) Delete(key []byte) {
	txn := b.db.NewTransaction(true)
	defer txn.Commit()

	_ = txn.Delete(key)
}

// Iterate implements the iterate interface for badger storage
func (b *BadgerStore) Iterate(key []byte) Iterator {
	txn := b.db.NewTransaction(false)

	opts := badger.IteratorOptions{}
	opts.PrefetchSize = 10
	opts.Prefix = key

	iter := txn.NewIterator(opts)
	iter.Rewind()

	curr := -1
	return &BadgerIterator{iter, curr}
}

// Close closes a database
func (b *BadgerStore) Close() {
	close(b.quit)
	b.db.Close()
}

// NewBadgerStore opens a badger database from the provided path
func NewBadgerStore(path string) (Store, error) {
	opt := badger.DefaultOptions(path)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	quit := make(chan struct{}, 1)

	store := &BadgerStore{db, quit}
	store.Init()
	return store, nil
}
