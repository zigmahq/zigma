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
	"encoding/hex"
	"math/big"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multihash"
)

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
	return node != nil && len(node.Id) > 0 && len(node.Hash) > 0
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
func (n *Node) ZeroPrefixLen(comparator *Node) int {
	distance := n.XOR(comparator)
	for i, b := range distance {
		for j := 0; j < 8; j++ {
			if (b>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
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
