package p2p

import (
	"io"

	"github.com/libp2p/go-libp2p-core/network"
)

// stream errors
const (
	ErrStreamClosed = "stream closed"
	ErrStreamReset  = "stream reset"
	ErrUnknownCodec = "multicodec did not match"
)

func (n *P2P) streamHandlerWrapper(stream network.Stream) {
	n.streamHandler(WrapStream(stream), nil)
}

func (n *P2P) streamHandler(ws *WrappedStream, replies chan Message) {
	for {
		m, err := ws.receive()
		if err != nil {
			switch err.Error() {
			case io.EOF.Error(), ErrStreamClosed, ErrStreamReset, ErrUnknownCodec:
			default:
				logger.Error("stream error: " + err.Error())
			}
			break
		}

		if replies != nil {
			go func() { replies <- m }()
		}

		if n.itf != nil {
			if err := n.itf.MessageReceived(ws, m); err != nil {
				logger.Error("handler error: " + err.Error())
			}
		}
	}
}
