package node

import (
	"errors"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zigmahq/zigma/log"
	"github.com/zigmahq/zigma/p2p"
	"github.com/zigmahq/zigma/upgrader"
	"github.com/zigmahq/zigma/version"
)

// Errors for PeerFound
var (
	ErrUnsupportedProtocol        = errors.New("unsupported protocol")
	ErrUnsupportedProtocolVersion = errors.New("unsupported protocol version")
	ErrHandshakeFailure           = errors.New("handshake failure")
	ErrUndeliveredMessage         = errors.New("undelivered message")
)

// PeerFound implements the handler for P2P.Implementer.PeerFound
func (n *Node) PeerFound(pi peer.AddrInfo) error {
	pid := pi.ID.Pretty()

	// attempt to upgrade the connection and open a zigma protocol channel
	if err := n.ConnectWithPeerInfo(pi); err != nil {
		return ErrUnsupportedProtocol
	}

	b, err := version.Current.Marshal()
	if err != nil {
		return ErrUnsupportedProtocolVersion
	}

	// compose and send handshake message
	replies := make(chan p2p.Message)
	req := p2p.NewMessage(n.Host().ID(), MsgHandshake, b)
	req = req.WithHeaders("phase", "request")

	if err := n.SendMessage(n.ctx, req, replies, pi.ID); err != nil {
		close(replies)
		return ErrUndeliveredMessage
	}

	// waiting for handshake response
	res := <-replies

	var ver version.Version
	if err := ver.Unmarshal(res.Body); err != nil {
		return ErrHandshakeFailure
	}

	if !ver.IsValid() {
		logger.Error("handshake failure", log.String("peer-id", pid))

		n.Host().Network().ClosePeer(pi.ID)
		return ErrHandshakeFailure
	}

	logger.Debug("handshake successful", log.String("peer-id", pid))
	n.Host().ConnManager().TagPeer(pi.ID, "zigma-support", 1)

	if ver.NewerThan(version.Current) {
		upgrader.FromPeer(n.Host(), pi.ID, version.Current, ver)
	}

	return nil
}
