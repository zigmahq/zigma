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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/store"
)

func newBadgerDB(t *testing.T) (store.Store, string) {
	dir, err := ioutil.TempDir("", "badger")
	assert.Nil(t, err)

	tmp := filepath.Join(dir, "db")
	d, err := store.NewBadgerStore(tmp)
	assert.Nil(t, err)

	return d, dir
}

func TestNewBadgerStore(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	assert.NotEmpty(t, dir)
	assert.NotNil(t, db)
}

func TestBadgerSetEx(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	o := []byte{0x62}

	db.Set(o, o, time.Second*2)

	b, ok := db.Get(o)
	assert.True(t, ok)
	assert.Equal(t, o, b)

	time.Sleep(time.Second*2 + time.Millisecond)

	b, ok = db.Get(o)
	assert.False(t, ok)
	assert.Nil(t, b)
}

func TestBadgerStoreGetAfterSet(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	arr := make([][]byte, 500)

	for i := 0; i < len(arr); i++ {
		s := strconv.Itoa(i)
		arr[i] = []byte(s)
		db.Set(arr[i], arr[i], 0)
	}
	for i := 0; i < len(arr); i++ {
		o, found := db.Get(arr[i])
		assert.True(t, found)
		assert.Equal(t, arr[i], o)
	}
}

func TestBadgerStoreDeleteAfterSet(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	o := []byte{0x61}

	b, ok := db.Get(o)
	assert.False(t, ok)
	assert.Nil(t, b)

	db.Set(o, o, 0)

	b, ok = db.Get(o)
	assert.True(t, ok)
	assert.Equal(t, o, b)

	db.Delete(o)

	b, ok = db.Get(o)
	assert.False(t, ok)
	assert.Nil(t, b)
}

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
