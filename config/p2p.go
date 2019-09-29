package config

import (
	"github.com/multiformats/go-multiaddr"
)

// P2PMultiAddr type
type P2PMultiAddr string

// Multiaddr parses and returns multi address from addr string
func (a P2PMultiAddr) Multiaddr() (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(string(a))
}

// MustMultiaddr parses and returns multi address from string, panic if failure
func (a P2PMultiAddr) MustMultiaddr() multiaddr.Multiaddr {
	m, err := a.Multiaddr()
	if err != nil {
		panic(err)
	}
	return m
}

// P2P encapsulates configuration options for zigma peer-to-peer communication
type P2P struct {
	// The node identifier
	PeerID string `yaml:"peer_id"`

	// The public and private key or the node
	PubKey string `yaml:"pub_key"`
	PriKey string `yaml:"pri_key"`

	// Port for p2p listener
	Port int `yaml:"listen_port"`

	// List of multi addresses to listen on
	ListenAddrs []string `yaml:"listen_addrs"`

	// List of seed addresses for bootstrapping
	BootstrapAddrs []P2PMultiAddr `yaml:"bootstrap_addrs"`
}

func maddr(s string) multiaddr.Multiaddr {
	ma, err := multiaddr.NewMultiaddr(s)
	if err != nil {
		panic(err)
	}
	return ma
}

// DefaultP2P generates the default configuration for p2p server
func DefaultP2P() *P2P {
	p2p := &P2P{
		ListenAddrs: []string{"/ip6/::/tcp", "/ip4/0.0.0.0/tcp"},
	}
	return p2p
}
