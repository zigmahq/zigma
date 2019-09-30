package p2p

import "github.com/libp2p/go-libp2p-core/protocol"

const (
	// ZProtocolID defines the zigma protocol identifer
	ZProtocolID = protocol.ID("/zigma")
	// ZProtocolVersion defines the version of zigma protocol
	ZProtocolVersion = "/0.0.1"
	// ZMDNSServiceName defines the service name for mdns discovery
	ZMDNSServiceName = "_zigma-discovery._udp"
)
