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

package dht

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"math/bits"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multihash"
)

// MockNode initializes node in unit tests
func MockNode(i int) *Node {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(i))
	h, err := multihash.Sum(bs, h, -1)
	if err != nil {
		return nil
	}
	return NodeFromHash(h)
}

// NodeFromHash initializes a dht node from multihash
func NodeFromHash(mh multihash.Multihash) *Node {
	d, err := multihash.Decode(mh)
	if err != nil {
		return nil
	}
	n := &Node{
		Id:     mh,
		Hash:   d.Digest,
		PeerId: mh,
	}
	return n
}

// NodeFromPeerID initializes a dht node from go-libp2p peer.ID
func NodeFromPeerID(pid peer.ID) *Node {
	p, err := pid.ExtractPublicKey()
	if err != nil {
		return nil
	}
	r, err := p.Raw()
	if err != nil {
		return nil
	}
	h, err := multihash.Sum(r, h, -1)
	if err != nil {
		return nil
	}
	d, err := multihash.Decode(h)
	if err != nil {
		return nil
	}
	n := &Node{
		Id:     r,
		Hash:   d.Digest,
		PeerId: []byte(pid.Pretty()),
	}
	return n
}

// IsValidNode checks if the node is valid
func IsValidNode(node *Node) bool {
	if node == nil {
		return false
	}
	return len(node.Id) > 0 && len(node.Id) <= 66 && len(node.Hash) > 0 && len(node.Hash) <= 64
}

// HexString returns id in the hex format
func (n *Node) HexString() []byte {
	dst := make([]byte, hex.EncodedLen(len(n.Hash)))
	hex.Encode(dst, n.Hash)
	return dst
}

// Equal checks if two nodes are equal
func (n *Node) Equal(comparator *Node) bool {
	return bytes.Equal(n.Id, comparator.Id)
}

// XOR calculates instance between two nodes
func (n *Node) XOR(comparator *Node) []byte {
	l := len(n.Hash)
	if t := len(comparator.Hash); l > t {
		l = t
	}
	o := make([]byte, l)
	for i := 0; i < l; i++ {
		o[i] = n.Hash[i] ^ comparator.Hash[i]
	}
	return o
}

// ZeroPrefixLen returns the number of consecutive zeroes in a byte slice
// d(1, 4) = d(001₂, 100₂)
// = 001₂ XOR 100₂
// = 101₂
// = 5
func (n *Node) ZeroPrefixLen(comparator *Node) int {
	distance := n.XOR(comparator)
	for i, b := range distance {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(distance)*8 - 1
}

// DistanceBetween between two nodes
func (n *Node) DistanceBetween(comparator *Node) *big.Int {
	bf1 := new(big.Int).SetBytes(n.Hash)
	bf2 := new(big.Int).SetBytes(comparator.Hash)
	xor := new(big.Int).Xor(bf1, bf2)
	return xor
}
