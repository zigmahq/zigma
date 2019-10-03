package config

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/zigmahq/zigma/config/types"
)

// RelayType for p2p connections
type RelayType string

// Define relay settings
const (
	RelayActive  RelayType = "active"
	RelayNat               = "nat"
	RelayDisable           = "disable"
)

// P2P encapsulates configuration options for zigma peer-to-peer communication
type P2P struct {
	Name            string        `yaml:"name"`              // the custom node identifier
	PrivKey         string        `yaml:"priv_key"`          // the path of the node private key
	Gossip          bool          `yaml:"gossip"`            // enable or disable gossip
	Relay           RelayType     `yaml:"relay"`             // enable or disable relay
	MinNumConns     int           `yaml:"min_num_conns"`     // the min number of connections allowed
	MaxNumConns     int           `yaml:"max_num_conns"`     // the max number of connections allowed
	ConnGracePeriod time.Duration `yaml:"conn_grace_period"` // the connect timeout settings
	Address         types.Addrs   `yaml:"laddr"`             // address to listen for incoming connections
	Seeds           types.Addrs   `yaml:"seeds"`             // list of seed nodes to connect to
	PersistentPeers types.Addrs   `yaml:"persistent_peers"`  // list of nodes to keep persistent connections
	QUIC            bool          `yaml:"quic"`              // experimental quic support
	MDNS            *MDNS         `yaml:"mdns"`
	RateLimit       *RateLimit    `yaml:"rate_limt"`
}

// MDNS encapsulates configuration options for mdns
type MDNS struct {
	Enable         bool          `yaml:"enable"`
	RescanInterval time.Duration `yaml:"rescan_interval"`
}

// RateLimit encapsulates configuration options for ratelimiting
type RateLimit struct {
	Enable      bool `yaml:"enable"`
	GlobalAvg   int  `yaml:"global_avg"`
	GlobalBurst int  `yaml:"global_burst"`
	PeerAvg     int  `yaml:"peer_avg"`
	PeerBurst   int  `yaml:"peer_burst"`
}

// DefaultP2P generates the default configuration for p2p server
func DefaultP2P() *P2P {
	p2p := &P2P{
		Address:         []*types.Addr{types.NewAddr("tcp://0:0")},
		ConnGracePeriod: 0,
		Gossip:          true,
		Relay:           RelayActive,
		RateLimit:       DefaultRateLimit(),
		QUIC:            true,
		MDNS:            DefaultMDNS(),
	}
	p2p.GenerateEd25519Key()
	return p2p
}

// DefaultMDNS generates the default configuration for mdns discovery
func DefaultMDNS() *MDNS {
	return &MDNS{
		Enable:         true,
		RescanInterval: time.Second * 10,
	}
}

// DefaultRateLimit generates the default configuration for p2p rate limiting
func DefaultRateLimit() *RateLimit {
	return &RateLimit{
		Enable:      false,
		GlobalAvg:   300,
		GlobalBurst: 500,
		PeerAvg:     300,
		PeerBurst:   500,
	}
}

// GenerateEd25519Key to generate public and private keys for node
func (p *P2P) GenerateEd25519Key() (crypto.PrivKey, error) {
	if k, err := p.DecodePrivateKey(); k != nil && err == nil {
		return k, err
	}
	f, err := ioutil.TempFile(os.TempDir(), "ed25519")
	if err != nil {
		return nil, err
	}
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	b, err := priv.Bytes()
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(f.Name(), b, 0644); err != nil {
		return nil, err
	}
	p.PrivKey = f.Name()
	return priv, nil
}

// DecodePrivateKey decodes the private key from private key path
func (p *P2P) DecodePrivateKey() (crypto.PrivKey, error) {
	b, err := ioutil.ReadFile(p.PrivKey)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(b)
}

// DecodePublicKey decodes and returns the public key from private key
func (p *P2P) DecodePublicKey() (crypto.PubKey, error) {
	k, err := p.DecodePrivateKey()
	if err != nil {
		return nil, err
	}
	return k.GetPublic(), nil
}

// DecodePeerID takes a public key and decode it to a peer.ID
func (p *P2P) DecodePeerID() (peer.ID, error) {
	k, err := p.DecodePublicKey()
	if err != nil {
		return "", err
	}
	return peer.IDFromPublicKey(k)
}

// DecodeListenAddrs decodes multiaddr addresses from listening addresses
func (p *P2P) DecodeListenAddrs() ([]multiaddr.Multiaddr, error) {
	var addrs = make([]multiaddr.Multiaddr, len(p.Address))
	for i, addr := range p.Address {
		m, err := addr.Multiaddr()
		if err != nil {
			return nil, err
		}
		addrs[i] = m
	}
	return addrs, nil
}

// DecodeListenAddrStrings decodes mulitaddr addresses from listen address urls
func (p *P2P) DecodeListenAddrStrings() ([]string, error) {
	var addrs = make([]string, len(p.Address))
	for i, addr := range p.Address {
		m, err := addr.Multiaddr()
		if err != nil {
			return nil, err
		}
		addrs[i] = m.String()
	}
	return addrs, nil
}
