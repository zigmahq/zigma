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

package dht_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/dht"
)

func TestNodeXOR(t *testing.T) {
	expected := []byte{0x01, 0x00, 0x01}
	n1 := &dht.Node{Hash: []byte{0x00, 0x00, 0x01}}
	n2 := &dht.Node{Hash: []byte{0x01, 0x00, 0x00}}
	xor := n1.XOR(n2)
	assert.Equal(t, expected, xor)
}

func TestNodeZeroPrefixLen(t *testing.T) {
	n1 := &dht.Node{Hash: []byte{0x01}}
	n2 := &dht.Node{Hash: []byte{0x04}}
	zl := n1.ZeroPrefixLen(n2)
	assert.Equal(t, 5, zl)
}

func TestNodeDistanceBetween(t *testing.T) {
	n1 := &dht.Node{Hash: []byte{0x01}}
	n2 := &dht.Node{Hash: []byte{0x04}}
	distance := n1.DistanceBetween(n2)
	assert.Equal(t, 5, int(distance.Int64()))
}
