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

var mrpcs []*mockRPC

type mockRPC struct {
	self    *Node
	receive chan *Message
}

func (m *mockRPC) Write(msg *Message) {
	everyone := msg.Receiver == nil
	for _, r := range mrpcs {
		switch {
		case everyone && r.self.Equal(m.self):
			r.receive <- msg
		case r.self.Equal(msg.Receiver):
			r.receive <- msg
		}
	}
}

func (m *mockRPC) Read() <-chan *Message {
	return m.receive
}

// MockRPC returns a mock implementation of rpc for unit testing
func MockRPC(self *Node, reset ...bool) KademliaRPC {
	c := make(chan *Message)
	r := &mockRPC{self, c}
	if len(reset) > 0 && reset[0] {
		mrpcs = []*mockRPC{r}
	} else {
		mrpcs = append(mrpcs, r)
	}
	return r
}
