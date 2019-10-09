package node

import (
	"github.com/libp2p/go-libp2p-net"
	"github.com/zigmahq/zigma/log"
)

// PeerConnected implements the handler for P2P.Implementer.PeerConnected
func (n *Node) PeerConnected(net net.Network, conn net.Conn) {
	logger.Debug("connected", log.String("peer-id", conn.RemotePeer().Pretty()))
}

// PeerDisconnected implements the handler for P2P.Implementer.PeerDisconnected
func (n *Node) PeerDisconnected(net net.Network, conn net.Conn) {
	logger.Debug("cancelled", log.String("peer-id", conn.RemotePeer().Pretty()))
}

// PeerStreamOpened implements the handler for P2P.Implementer.PeerStreamOpened
func (n *Node) PeerStreamOpened(net net.Network, stream net.Stream) {}

// PeerStreamClosed implements the handler for P2P.Implementer.PeerStreamClosed
func (n *Node) PeerStreamClosed(net net.Network, stream net.Stream) {}
