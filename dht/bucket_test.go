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
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/dht"
)

var bucket *dht.Bucket
var nodeset []*dht.Node

func init() {
	bucket = dht.NewBucket()
	nodeset = make([]*dht.Node, 1000)
	for i := range nodeset {
		nodeset[i] = newBucketNode(i)
	}
}

func newBucketNode(i int) *dht.Node {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(i))
	h := sha256.New()
	h.Write(bs)
	o := []byte(hex.EncodeToString(h.Sum(nil)))
	return &dht.Node{Id: o, Hash: o}
}

func BenchmarkNodesetInit(b *testing.B) {
	b.ReportAllocs()
	if b.N > len(nodeset) {
		nodeset = make([]*dht.Node, b.N)
	}
	for i := 0; i < b.N; i++ {
		nodeset[i] = newBucketNode(i)
	}
}

func BenchmarkNodeBucket(b *testing.B) {
	b.ReportAllocs()
	bucket.RemoveAllNodes()
	t := len(nodeset)
	for i := 0; i < b.N; i++ {
		bucket.AddNode(nodeset[i%t])
	}
	for i := 0; i < b.N; i++ {
		bucket.RemoveNode(nodeset[i%t])
	}
}

func TestBucketAddNode(t *testing.T) {
	var k = len(bucket)
	bucket.RemoveAllNodes()
	assert.Zero(t, bucket.Len())
	for i := 0; i < k; i++ {
		bucket.AddNode(nodeset[i])
		assert.Equal(t, i+1, bucket.Len())
		assert.Equal(t, nodeset[i], bucket[len(bucket)-1])
	}
	for i := 0; i < k; i++ {
		bucket.AddNode(nodeset[i])
		assert.Equal(t, k, bucket.Len())
		assert.Equal(t, nodeset[i], bucket[len(bucket)-1])
		if i < k-1 {
			assert.Equal(t, nodeset[i+1], bucket[0])
		}
	}
	bucket.AddNode(nodeset[1])
	assert.Equal(t, nodeset[1], bucket[len(bucket)-1])
	assert.Equal(t, nodeset[2], bucket[1])
}

func TestBucketRemoveNode(t *testing.T) {
	var k = len(bucket)
	bucket.RemoveAllNodes()
	assert.Zero(t, bucket.Len())
	bucket.AddNode(nodeset[0])
	bucket.AddNode(nodeset[1])
	bucket.AddNode(nodeset[2])
	bucket.AddNode(nodeset[3])
	assert.Equal(t, 4, bucket.Len())
	bucket.RemoveNode(nodeset[2])
	assert.Equal(t, bucket[k-2], nodeset[1])
	assert.Equal(t, bucket[k-3], nodeset[0])
	assert.Equal(t, 3, bucket.Len())
	bucket.RemoveNode(nodeset[3])
	assert.Equal(t, bucket[k-1], nodeset[1])
	assert.Equal(t, bucket[k-2], nodeset[0])
	assert.Equal(t, 2, bucket.Len())
	bucket.RemoveNode(nodeset[0])
	assert.Equal(t, bucket[k-1], nodeset[1])
	assert.Equal(t, 1, bucket.Len())
	bucket.RemoveNode(nodeset[1])
	assert.Equal(t, 0, bucket.Len())
}

func TestBucketIterator(t *testing.T) {
	var n = 5
	var count int
	bucket.RemoveAllNodes()
	assert.Zero(t, bucket.Len())
	for i := 0; i < n; i++ {
		bucket.AddNode(nodeset[i])
	}
	assert.Equal(t, n, bucket.Len())
	for node := range bucket.Iterator() {
		assert.NotNil(t, node)
		count++
	}
	assert.Equal(t, n, count)
}
