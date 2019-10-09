package node

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zigmahq/zigma/p2p"
	"github.com/zigmahq/zigma/upgrader"
	"github.com/zigmahq/zigma/version"
)

func (n *Node) handshake(ws *p2p.WrappedStream, msg p2p.Message) error {
	pi, err := peer.IDB58Decode(msg.Originator)
	if err != nil {
		return err
	}

	if phase, ok := msg.Headers["phase"]; !ok || phase != "request" {
		return ws.Dispatch(msg.WithBody(nil))
	}

	var ver version.Version
	if err := ver.Unmarshal(msg.Body); err != nil {
		return err
	}

	if !ver.IsValid() {
		n.Host().Network().ClosePeer(pi)
		return nil
	}

	b, err := version.Current.Marshal()
	if err != nil {
		return err
	}

	if ver.NewerThan(version.Current) {
		upgrader.FromPeer(n.Host(), pi, version.Current, ver)
	}

	n.Host().ConnManager().TagPeer(pi, "zigma-support", 1)
	return ws.Dispatch(msg.WithBody(b))
}
