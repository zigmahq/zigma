package p2p

import (
	"bufio"
	"io"

	"github.com/libp2p/go-libp2p-core/network"
	multicodec "github.com/multiformats/go-multicodec"
	protobuf "github.com/multiformats/go-multicodec/protobuf"
)

// WrappedStream wraps a libp2p stream
type WrappedStream struct {
	stream network.Stream
	enc    multicodec.Encoder
	dec    multicodec.Decoder
	w      *bufio.Writer
	r      *bufio.Reader
}

// WrapStream takes a stream
func WrapStream(s network.Stream) *WrappedStream {
	reader := bufio.NewReader(s)
	writer := bufio.NewWriter(s)
	enc := protobuf.Multicodec(new(Message)).Encoder(writer)
	dec := protobuf.Multicodec(new(Message)).Decoder(reader)
	return &WrappedStream{
		stream: s,
		enc:    enc,
		dec:    dec,
		r:      reader,
		w:      writer,
	}
}

func (ws *WrappedStream) receive() (Message, error) {
	var msg Message
	err := ws.dec.Decode(&msg)
	return msg, err
}

// Dispatch sends a message to peer
func (ws *WrappedStream) Dispatch(m Message) error {
	if err := ws.enc.Encode(&m); err != nil {
		return err
	}
	if err := ws.w.Flush(); err != nil {
		if err != io.EOF && err.Error() != "stream closed" {
			return err
		}
	}
	return nil
}
