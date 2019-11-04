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

import "github.com/multiformats/go-multihash"

// Hashable type
type Hashable interface {
	Key() []byte
	Data() []byte
	Hash() []byte
}

type hashable struct {
	customKey []byte
	data      []byte
	hash      []byte
}

func (h *hashable) Key() []byte {
	if len(h.customKey) > 0 {
		return h.customKey
	}
	return h.Hash()
}

func (h *hashable) Data() []byte {
	return h.data
}

func (h *hashable) Hash() []byte {
	return h.hash
}

// Bytes returns data in Hashable type
func Bytes(data []byte) Hashable {
	d := &hashable{
		data: data,
	}
	h, err := multihash.Sum(d.data, h, -1)
	if err != nil {
		panic(err)
	}
	d.hash = h
	return d
}

// String returns data in Hashable type
func String(data string) Hashable {
	return Bytes([]byte(data))
}
