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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/dht"
)

func TestMockRPCReadWrite(t *testing.T) {
	n1 := dht.MockNode(0)
	r1 := dht.MockRPC(n1, true)
	n2 := dht.MockNode(1)
	r2 := dht.MockRPC(n2)

	assert.NotNil(t, n1)
	assert.NotNil(t, r1)
	assert.NotNil(t, n2)
	assert.NotNil(t, r2)

	var i int
	l := 10
	c := make(chan struct{})
	r := make(chan struct{})

	go func() {
		for {
			select {
			case <-r1.Read():
				r <- struct{}{}
			case <-r2.Read():
				r <- struct{}{}
			}
		}
	}()

	go func() {
		for i := 0; i < l; i++ {
			if i%2 == 0 {
				r1.Write(&dht.Message{Id: []byte{byte(i)}, Receiver: n2})
			} else {
				r2.Write(&dht.Message{Id: []byte{byte(i)}, Receiver: n1})
			}
		}
	}()

	go func() {
		for range r {
			i++
		}
	}()

	go func() {
		time.Sleep(time.Second)
		assert.Equal(t, l, i)
		c <- struct{}{}
	}()

	<-c
}

func TestMockRPCWriteResponse(t *testing.T) {
	n1 := dht.MockNode(0)
	r1 := dht.MockRPC(n1, true)
	n2 := dht.MockNode(1)
	r2 := dht.MockRPC(n2)

	go func() {
		for {
			select {
			case msg := <-r2.Read():
				msg.Sender, msg.Receiver = msg.Receiver, msg.Sender
				msg.IsResponse = true
				r2.Write(msg)
			}
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := []byte{byte(i)}
			reply := <-r1.Write(&dht.Message{
				Id:       id,
				Type:     dht.MessageType_NOOP,
				Sender:   n1,
				Receiver: n2,
			})(0)

			assert.NotNil(t, reply)
			assert.Equal(t, n1, reply.Receiver)
			assert.Equal(t, n2, reply.Sender)
			assert.Equal(t, id, reply.Id)
			assert.True(t, reply.IsResponse)
		}(i)
	}
	wg.Wait()
}

func TestMockRPCWriteTimeoutResponse(t *testing.T) {
	n1 := dht.MockNode(0)
	r1 := dht.MockRPC(n1, true)
	n2 := dht.MockNode(1)
	r2 := dht.MockRPC(n2)

	go func() {
		msg := <-r2.Read()
		msg.Sender, msg.Receiver = msg.Receiver, msg.Sender
		msg.IsResponse = true
		time.Sleep(time.Second * 10)
		r2.Write(msg)
	}()

	reply := <-r1.Write(&dht.Message{
		Id:       []byte("1"),
		Type:     dht.MessageType_NOOP,
		Sender:   n1,
		Receiver: n2,
	})(0)

	assert.Nil(t, reply)
}
