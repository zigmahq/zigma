package node

import (
	"fmt"

	"github.com/zigmahq/zigma/p2p"
)

// MessageReceived implements the handler for P2P.Implementer.MessageReceived
func (n *Node) MessageReceived(ws *p2p.WrappedStream, m p2p.Message) error {
	switch m.Type {
	default:
		return n.handshake(ws, m)
	}
	return nil
}

// MessageSent implements the handler for P2P.Implementer.MessageSent
func (n *Node) MessageSent(ws *p2p.WrappedStream, m p2p.Message) error {
	fmt.Println("message sent", m)
	return nil
}
