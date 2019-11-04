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

import (
	"sync"
	"time"
	"unsafe"
)

var mrpcs []*mockRPC

type mockRPC struct {
	self    *Node
	receive chan *Message
	replies *sync.Map
}

func (m *mockRPC) Write(msg *Message) func(time.Duration) <-chan *Message {
	var id = *(*string)(unsafe.Pointer(&msg.Id))
	var wc = make(chan *Message, 1)

	if len(id) > 0 && !msg.IsResponse {
		m.replies.Store(id, make(chan *Message, 1))
	} else {
		wc <- nil
	}

	for _, r := range mrpcs {
		if r.self.Equal(msg.Receiver) {
			go func() { r.receive <- msg }()
			if len(id) > 0 && msg.IsResponse {
				go func() {
					if d, ok := r.replies.Load(id); ok {
						d.(chan *Message) <- msg
					}
				}()
			}
			break
		}
	}
	return func(timeout time.Duration) <-chan *Message {
		if len(id) > 0 && !msg.IsResponse {
			go func() {
				var (
					c chan *Message
					t time.Duration
				)
				if timeout > 0 {
					t = timeout
				} else {
					t = time.Second
				}
				if d, ok := m.replies.Load(id); ok {
					c = d.(chan *Message)
				}
				select {
				case msg := <-c:
					wc <- msg
				case <-time.After(t):
					wc <- nil
					close(c)
				}
				m.replies.Delete(id)
			}()
		}
		return wc
	}
}

func (m *mockRPC) Read() <-chan *Message {
	return m.receive
}

// MockRPC returns a mock implementation of rpc for unit testing
func MockRPC(self *Node, reset ...bool) KademliaRPC {
	c := make(chan *Message)
	r := new(sync.Map)
	m := &mockRPC{self, c, r}
	if len(reset) > 0 && reset[0] {
		mrpcs = []*mockRPC{m}
	} else {
		mrpcs = append(mrpcs, m)
	}
	return m
}
