package p2p

import (
	"time"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/peer"
)

// NewMessage initializes a new message type that contains a message body
func NewMessage(originator peer.ID, t uint32, body []byte) Message {
	ts := time.Now().UTC()
	return Message{
		Id:           NewMessageID(),
		Originator:   originator.Pretty(),
		CreatedSec:   ts.Unix(),
		CreatedNsec:  int64(ts.Nanosecond()),
		DeadlineSec:  ts.Unix(),
		DeadlineNsec: int64(ts.Nanosecond()),
		Headers:      make(map[string]string),
		Body:         body,
	}
}

// NewMessageID generates an unique identifier for message
func NewMessageID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Error(err.Error())
	}
	return uuid.String()
}

// WithReceiver sets the receiver of this message
func (m Message) WithReceiver(destination peer.ID) Message {
	return Message{
		Id:           m.Id,
		Originator:   m.Originator,
		Destination:  destination.Pretty(),
		CreatedSec:   m.CreatedSec,
		CreatedNsec:  m.CreatedNsec,
		DeadlineSec:  m.DeadlineSec,
		DeadlineNsec: m.DeadlineNsec,
		Headers:      m.Headers,
		Body:         m.Body,
	}
}

// WithHeaders attach header metadata to Message
func (m Message) WithHeaders(kv ...string) Message {
	headers := map[string]string{}
	for i := 0; i < len(kv)-1; i = i + 2 {
		headers[kv[i]] = kv[i+1]
	}
	return Message{
		Id:           m.Id,
		Originator:   m.Originator,
		Destination:  m.Destination,
		CreatedSec:   m.CreatedSec,
		CreatedNsec:  m.CreatedNsec,
		DeadlineSec:  m.DeadlineSec,
		DeadlineNsec: m.DeadlineNsec,
		Headers:      headers,
		Body:         m.Body,
	}
}

// WithBody returns a new message with a new body
func (m Message) WithBody(b []byte) Message {
	return Message{
		Id:           m.Id,
		Originator:   m.Originator,
		Destination:  m.Destination,
		CreatedSec:   m.CreatedSec,
		CreatedNsec:  m.CreatedNsec,
		DeadlineSec:  m.DeadlineSec,
		DeadlineNsec: m.DeadlineNsec,
		Headers:      m.Headers,
		Body:         b,
	}
}
