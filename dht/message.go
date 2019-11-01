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

import "github.com/google/uuid"

func compose(sender *Node) *Message {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	id, _ := uid.MarshalText()
	return &Message{Id: id, Sender: sender}
}

func (m *Message) ping() *Message {
	m.Type = MessageType_PING
	return m
}

func (m *Message) pong() *Message {
	var n = new(Message)
	*n = *m

	n.IsResponse = true
	n.Sender, n.Receiver = n.Receiver, n.Sender
	n.Request = nil
	return n
}

func (m *Message) store(payload Hashable) *Message {
	if payload == nil {
		return m
	}
	m.Type = MessageType_STORE
	m.Request = &Message_Store{
		Store: &StoreRequest{
			Payload: &Payload{
				Key:  payload.Hash(),
				Data: payload.Data(),
				Hash: payload.Hash(),
			},
		},
	}
	return m
}

func (m *Message) success(success bool) *Message {
	var n = new(Message)
	*n = *m

	n.IsResponse = true
	n.Sender, n.Receiver = n.Receiver, n.Sender
	n.Request = nil
	n.Response = &Message_Success{
		Success: success,
	}
	return n
}

func (m *Message) findNode(id []byte) *Message {
	m.Type = MessageType_FIND_NODE
	m.Request = &Message_Find{
		Find: &FindRequest{
			Key: id,
		},
	}
	return m
}

func (m *Message) returnClosest(nodes []*Node) *Message {
	var n = new(Message)
	*n = *m

	n.IsResponse = true
	n.Sender, n.Receiver = n.Receiver, n.Sender
	n.Request = nil
	n.Response = &Message_Closest{
		Closest: &Closest{
			Nodes: nodes,
		},
	}
	return n
}

func (m *Message) findValue(id []byte) *Message {
	m.Type = MessageType_FIND_VALUE
	m.Request = &Message_Find{
		Find: &FindRequest{
			Key: id,
		},
	}
	return m
}

func (m *Message) returnValue(data []byte) *Message {
	var n = new(Message)
	*n = *m

	n.IsResponse = true
	n.Sender, n.Receiver = n.Receiver, n.Sender
	n.Request = nil
	n.Response = &Message_Payload{
		Payload: &Payload{
			Key:  m.GetFind().Key,
			Data: data,
			Hash: nil,
			Sig:  nil,
		},
	}
	return n
}

func (m *Message) to(receiver *Node) *Message {
	m.Receiver = receiver
	return m
}

func (m *Message) isValid() bool {
	switch {
	case m.Sender == nil || m.Receiver == nil:
		return false

	case m.Type == MessageType_PING:
		return true

	case m.Type == MessageType_FIND_VALUE:
		return m.GetFind() != nil && len(m.GetFind().Key) > 0

	case m.Type == MessageType_FIND_NODE:
		return m.GetFind() != nil && len(m.GetFind().Key) > 0

	case m.Type == MessageType_STORE:
		if m.GetStore() != nil || m.GetStore().Payload != nil {
			payload := m.GetStore().Payload
			return len(payload.Key) > 0 && len(payload.Data) > 0
		}
	}
	return false
}
