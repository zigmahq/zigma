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

package store_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadgerStoreIterator(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	ks := make([][]byte, 30)
	vs := make([][]byte, len(ks))
	var i int

	prefix := []byte{0x61, 0x3a, 0x62, 0x3a}

	for i := 0; i < len(ks); i++ {
		s := strconv.Itoa(10 + i)
		k := append(prefix, s...)
		v := []byte(s)

		ks[i] = k
		vs[i] = v
		db.Set(k, v, 0)
	}

	iter := db.Iterate(prefix)
	defer iter.Done()

	for iter.Next() {
		item := iter.Item()
		assert.Equal(t, ks[i], item.Key())
		assert.Equal(t, vs[i], item.Value())
		assert.Zero(t, item.TTL())
		i++
	}
}
