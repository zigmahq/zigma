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

// BadgerIterator implements the iterator interface
type BadgerIterator struct {
	iter *badger.Iterator
	curr int
}

// BadgerIteratorItem implements the iterator item interface
type BadgerIteratorItem struct {
	item *badger.Item
}

// Seek would seek to the provided key if present, or rewind the
// iterator cursor all the way to zero-th position if the key is nil
func (b *BadgerIterator) Seek(key []byte) {
	b.curr = -1
	b.iter.Seek(key)
}

// Next returns the true if next item is available, returns
// false when iteration is done
func (b *BadgerIterator) Next() bool {
	if b.curr > -1 {
		b.iter.Next()
	}
	b.curr++
	return b.iter.Valid()
}

// Item returns the current key-value item
func (b *BadgerIterator) Item() Item {
	return &BadgerIteratorItem{b.iter.Item()}
}

// Done closes the badger iterator instance
func (b *BadgerIterator) Done() {
	b.iter.Close()
}

// Key returns the badger item key name
func (b *BadgerIteratorItem) Key() []byte {
	return b.item.Key()
}

// Value returns the badger item key value
func (b *BadgerIteratorItem) Value() []byte {
	t, err := b.item.ValueCopy(nil)
	if err != nil {
		return nil
	}
	return t
}

// TTL returns the badger item key expiration time
func (b *BadgerIteratorItem) TTL() time.Time {
	ex := b.item.ExpiresAt()
	if ex > 0 {
		sec := int64(b.item.ExpiresAt())
		return time.Unix(sec, 0)
	}
	return time.Time{}
}
