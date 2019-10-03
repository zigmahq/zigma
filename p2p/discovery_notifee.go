package p2p

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zigmahq/zigma/log"
)

type discoveryNotifee struct {
	p2p *P2P
}

// interface to be called when a new peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	// do not establish connection to itself
	if n.p2p.id == pi.ID {
		return
	}

	// do not establish connection to peer if it has been connected or
	// connection failure was occurred
	switch n.p2p.host.Network().Connectedness(pi.ID) {
	case network.Connected, network.CannotConnect:
		return
	}

	// do not establish connection to peer if the peer address is found
	// inside the peerstore, address exists in peerstore means a connection
	// has been established or establishing
	if addrs := n.p2p.host.Peerstore().Addrs(pi.ID); len(addrs) > 0 {
		return
	}

	// there is no existing connection between local and remote peer node,
	// so establish a connection with remote node
	if err := n.p2p.host.Connect(n.p2p.ctx, pi); err != nil {
		n.p2p.logger.Error(
			"connection failure",
			log.String("peer-id", pi.ID.Pretty()),
			log.String("reason", err.Error()))
	}
}
