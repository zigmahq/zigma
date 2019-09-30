package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
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

	// Allowed transports for the node
	Transports []string `yaml:"transports"`

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
	pri, pub, _ := crypto.GenerateEd25519Key(rand.Reader)
	pid, _ := peer.IDFromPublicKey(pub)
	bpri, _ := pri.Bytes()
	bpub, _ := pub.Bytes()

	p2p := &P2P{
		PeerID:      pid.String(),
		PubKey:      hex.EncodeToString(bpub),
		PriKey:      hex.EncodeToString(bpri),
		Transports:  []string{"tcp", "ws"},
		ListenAddrs: []string{"/ip4/0.0.0.0"},
	}
	return p2p
}

// DecodePeerID takes a string and decode it to a peer.ID
func (p *P2P) DecodePeerID() (peer.ID, error) {
	return peer.IDB58Decode(p.PeerID)
}

// DecodePrivateKey generates a private key through a hex string
func (p *P2P) DecodePrivateKey() (crypto.PrivKey, error) {
	b, err := hex.DecodeString(p.PriKey)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(b)
}

// DecodePublicKey generates a public key through a hex string
func (p *P2P) DecodePublicKey() (crypto.PubKey, error) {
	b, err := hex.DecodeString(p.PubKey)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPublicKey(b)
}

// Addrs returns the addresses with listening port
func (p *P2P) Addrs() []string {
	var addrs []string
	var transports = map[string]struct{}{}
	for _, tran := range p.Transports {
		transports[tran] = struct{}{}
	}
	for _, addr := range p.ListenAddrs {
		if _, ok := transports["tcp"]; ok {
			addrs = append(addrs, fmt.Sprintf("%s/tcp/%d", addr, p.Port))
		}
		if _, ok := transports["ws"]; ok {
			addrs = append(addrs, fmt.Sprintf("%s/tcp/%d/ws", addr, p.Port))
		}
	}
	return addrs
}
