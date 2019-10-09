package upgrader

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zigmahq/zigma/log"
	"github.com/zigmahq/zigma/version"
)

var logger = log.DefaultLogger

// FromPeer handler
func FromPeer(h host.Host, pi peer.ID, v1, v2 version.Version) {
	logger.Debug("upgrade node software",
		log.String("current", v1.Number),
		log.String("target", v2.Number),
	)
}
