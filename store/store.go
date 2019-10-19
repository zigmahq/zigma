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

import "time"

// Store is the interface for implementing the basic storage mechanism
type Store interface {
	Set(key, val []byte, ttl time.Duration)
	Get(key []byte) (data []byte, found bool)
	Delete(key []byte)
	Iterate(key []byte) Iterator
	Close()
}

// Iterator is the interface for implementing iterator for key-value storage
type Iterator interface {
	Seek(key []byte)
	Next() bool
	Item() Item
	Done()
}

// Item is returned during iteration
type Item interface {
	Key() []byte
	Value() []byte
	TTL() time.Time
}
