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
)

func TestHashableBytes(t *testing.T) {
	b := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}
	h := dht.Bytes(b)
	assert.NotNil(t, h)
	assert.Equal(t, b, h.Data())
	assert.Len(t, h.Hash(), 66)
}

func TestHashableString(t *testing.T) {
	s := "hello"
	h := dht.String(s)
	assert.NotNil(t, h)
	assert.Equal(t, []byte(s), h.Data())
	assert.Len(t, h.Hash(), 66)
}
