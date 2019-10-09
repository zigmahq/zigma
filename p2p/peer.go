package p2p

import (
	"context"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/zigmahq/zigma/config/types"
	"github.com/zigmahq/zigma/version"
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

// SendMessage opens a stream and sends a message to the peer target
func (n *P2P) SendMessage(ctx context.Context, msg Message, replies chan Message, pids ...peer.ID) error {
	for _, pid := range pids {
		if pid.Pretty() == n.ID() {
			continue
		}
		s, err := n.host.NewStream(ctx, pid, n.cfg.ProtocolID())
		if err != nil {
			return err
		}
		defer s.Close()
		ws := WrapStream(s)
		go n.streamHandler(ws, replies)
		if err := ws.Dispatch(msg); err != nil {
			return err
		}
		if n.itf != nil {
			if err := n.itf.MessageSent(ws, msg); err != nil {
				return err
			}
		}
	}
	return nil
}

// ActivePeers returns the total number of active peers connected with this node
func (n *P2P) ActivePeers() int {
	return n.connManager.GetInfo().ConnCount
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
		pis[i] = n.kad.FindLocal(v)
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

// PeerProtocol returns the protocol id from connection
func (n *P2P) PeerProtocol(pi peer.ID) (string, bool) {
	var (
		delimiter = "/"
		prefix    = delimiter + n.cfg.Protocol.ID + delimiter
	)
	p, err := n.host.Peerstore().GetProtocols(pi)
	if err != nil {
		return "", false
	}
	for _, o := range p {
		if strings.HasPrefix(o, prefix) {
			return o, true
		}
	}
	return "", false
}

// VerifyProtocol verifies protocol with release pubic key
func (n *P2P) VerifyProtocol(pi peer.ID, sig []byte) bool {
	var (
		delimiter = "/"
		i         = 2
	)
	p, ok := n.PeerProtocol(pi)
	if !ok {
		return false
	}
	if d := strings.Split(p, delimiter); len(d) > i {
		if version.Verifier.Verify([]byte(d[i]), sig) {
			return true
		}
	}
	return false
}
