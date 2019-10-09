package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-pubsub"
)

// Broadcast sends a message to every single peer in the network
func (n *P2P) Broadcast(topic string, b []byte) error {
	return n.pubsub.Publish(n.cfg.Queue, b)
}

// Receive returns the next pending message. it blocks if no messages are available.
func (n *P2P) Receive(ctx context.Context) (*pubsub.Message, error) {
	msg, err := n.sub.Next(ctx)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
