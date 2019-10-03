package p2p

import "github.com/libp2p/go-libp2p-core/protocol"

const (
	// ZProtocolVersion defines the version of zigma protocol
	ZProtocolVersion = "0.0.1"
	// ZProtocolID defines the zigma protocol identifer
	ZProtocolID = protocol.ID("/zigma/" + ZProtocolVersion)
	// ZMDNSServiceName defines the service name for mdns discovery
	ZMDNSServiceName = "_zigma-discovery._udp"
)
