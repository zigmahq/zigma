package p2p

import (
	"log"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type discoveryNotifee struct {
	node *ZNode
}

// interface to be called when a new peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	// do not establish connection to itself
	if n.node.id == pi.ID {
		return
	}

	// do not establish connection to peer if it has been connected or
	// connection failure was occurred
	switch n.node.host.Network().Connectedness(pi.ID) {
	case network.Connected, network.CannotConnect:
		return
	}

	// do not establish connection to peer if the peer address is found
	// inside the peerstore, address exists in peerstore means a connection
	// has been established or establishing
	if addrs := n.node.host.Peerstore().Addrs(pi.ID); len(addrs) > 0 {
		return
	}

	// there is no existing connection between local and remote peer node,
	// so establish a connection with remote node
	if err := n.node.host.Connect(n.node.ctx, pi); err != nil {
		log.Printf("[%s] connection failure: %s", pi.ID, err)
	}
}
