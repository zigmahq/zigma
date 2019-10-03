package p2p

import (
	"bufio"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/zigmahq/zigma/log"
)

func (n *P2P) handleStream(stream network.Stream) {
	// drop connection due to high traffic volume
	if n.cfg.RateLimit.Enable && !n.limiter.Allow() {
		stream.Conn().Close()
		return
	}

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	s, err := rw.ReadString('\n')
	if err != nil || len(s) == 0 {
		return
	}
	n.logger.Info("received", log.String("msg", s))
}
