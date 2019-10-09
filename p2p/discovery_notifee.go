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

	pid := pi.ID.Pretty()

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

	if n.p2p.itf == nil {
		return
	}

	if err := n.p2p.itf.PeerFound(pi); err != nil {
		logger.Error("PeerFound error: "+err.Error(), log.String("peer-id", pid))
	}
}
