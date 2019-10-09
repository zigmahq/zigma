package p2p

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-net"
)

// Implementer declare the communication interface between services and p2p server
type Implementer interface {
	PeerFound(peer.AddrInfo) error
	PeerConnected(net.Network, net.Conn)
	PeerDisconnected(net.Network, net.Conn)
	PeerStreamOpened(network.Network, network.Stream)
	PeerStreamClosed(network.Network, network.Stream)
	MessageReceived(*WrappedStream, Message) error
	MessageSent(*WrappedStream, Message) error
}
