package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/zigmahq/zigma/config/types"
)

// ConnectWithPeerInfo connects a peer with peerstore.PeerInfo
func (n *P2P) ConnectWithPeerInfo(pis ...peerstore.PeerInfo) error {
	for _, pi := range pis {
		if err := n.Host().Connect(n.ctx, pi); err != nil {
			return err
		}
	}
	return nil
}

// ConnectWithMultiaddr to connect peers with multiaddrs(es)
func (n *P2P) ConnectWithMultiaddr(addrs ...multiaddr.Multiaddr) error {
	a, err := peer.AddrInfosFromP2pAddrs(addrs...)
	if err != nil {
		return err
	}
	return n.ConnectWithPeerInfo(a...)
}

// ConnectWithAddr to connect peers with address(es)
func (n *P2P) ConnectWithAddr(addrs ...*types.Addr) error {
	in := types.Addrs(addrs)
	pis := in.PeerInfos()
	return n.ConnectWithPeerInfo(pis...)
}

// Unicast to send another address a message
func (n *P2P) Unicast(ctx context.Context, target peerstore.PeerInfo, data []byte) error {
	if err := n.ConnectWithPeerInfo(target); err != nil {
		return err
	}
	stream, err := n.host.NewStream(ctx, target.ID, ZProtocolID)
	if err != nil {
		return err
	}
	defer stream.Close()
	_, err = stream.Write(data)
	return err
}

// Broadcast sends a message to every single peer in the network
func (n *P2P) Broadcast(topic string, b []byte) error {
	return n.pubs.Publish(topic, b)
}

// ActivePeers returns the total number of active peers connected with this node
func (n *P2P) ActivePeers() int {
	return n.peers
}

// Peers returns the active peers that are connected with this node
func (n *P2P) Peers(ctx context.Context) []peerstore.PeerInfo {
	peers := n.host.Peerstore().Peers()
	dedup := make(map[string]peer.ID)
	for _, peer := range peers {
		if len(peer) == 0 || n.host.ID() == peer {
			continue
		}
		dedup[peer.Pretty()] = peer
	}
	var i int
	pis := make([]peerstore.PeerInfo, len(dedup))
	for _, v := range dedup {
		pis[i] = n.dht.FindLocal(v)
		i++
	}
	return pis
}

// PeerInfo returns node peer information
func (n *P2P) PeerInfo() peerstore.PeerInfo {
	return peerstore.PeerInfo{
		ID:    n.host.ID(),
		Addrs: n.host.Addrs(),
	}
}
