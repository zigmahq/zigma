package config

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/zigmahq/zigma/config/types"
	"github.com/zigmahq/zigma/version"
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
	Queue           string        `yaml:"-"`                 // the pubsub queue name
	Relay           RelayType     `yaml:"relay"`             // enable or disable relay
	MinNumConns     int           `yaml:"min_num_conns"`     // the min number of connections allowed
	MaxNumConns     int           `yaml:"max_num_conns"`     // the max number of connections allowed
	ConnGracePeriod time.Duration `yaml:"conn_grace_period"` // the connect timeout settings
	Address         types.Addrs   `yaml:"laddr"`             // address to listen for incoming connections
	Seeds           types.Addrs   `yaml:"seeds"`             // list of seed nodes to connect to
	PersistentPeers types.Addrs   `yaml:"persistent_peers"`  // list of nodes to keep persistent connections
	QUIC            bool          `yaml:"quic"`              // experimental quic support
	Insecure        bool          `yaml:"insecure"`          // allow insecure connection
	MDNS            *MDNS         `yaml:"mdns"`              // settings for mdns discovery
	DHT             *DHT          `yaml:"dht"`               // settings for dht kademlia discovery
	Protocol        *Protocol     `yaml:"protocol"`          // protocol options
	RateLimit       *RateLimit    `yaml:"rate_limt"`         // ratelimiting for p2p
}

// Protocol encapsulates configuration options for p2p protocol
type Protocol struct {
	ID        string `yaml:"-"`
	Name      string `yaml:"-"`
	Version   string `yaml:"-"`
	Signature []byte `yaml:"-"`
}

// MDNS encapsulates configuration options for mdns
type MDNS struct {
	Enable         bool          `yaml:"enable"`
	ServiceName    string        `yaml:"service_name"`
	RescanInterval time.Duration `yaml:"rescan_interval"`
}

// DHT encapsulates configuration options for dht kademlia
type DHT struct {
	Enable     bool          `yaml:"enable"`
	Rendezvous string        `yaml:"rendezvous"`
	TTL        time.Duration `yaml:"ttl"`
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
		Queue:           "zigma",
		ConnGracePeriod: 0,
		Relay:           RelayActive,
		RateLimit:       DefaultRateLimit(),
		QUIC:            true,
		Insecure:        false,
		MDNS:            DefaultMDNS(),
		DHT:             DefaultDHT(),
		Protocol:        DefaultProtocol(),
	}
	p2p.GenerateEd25519Key()
	return p2p
}

// DefaultProtocol generates the default configuration for p2p protocol
func DefaultProtocol() *Protocol {
	return &Protocol{
		ID:        "zigma",
		Name:      version.Current.Name,
		Version:   version.Current.Number,
		Signature: version.Current.Signature,
	}
}

// DefaultMDNS generates the default configuration for mdns discovery
func DefaultMDNS() *MDNS {
	return &MDNS{
		Enable:         true,
		ServiceName:    "_zigma-discovery._udp",
		RescanInterval: time.Second * 30,
	}
}

// DefaultDHT generates the default configuration for dht discovery
func DefaultDHT() *DHT {
	return &DHT{
		Enable:     true,
		Rendezvous: "/zigma/dht",
		TTL:        time.Second * 10,
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

// ListenAddrs decodes multiaddr addresses from listening addresses
func (p *P2P) ListenAddrs() ([]multiaddr.Multiaddr, error) {
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

// ListenAddrStrings decodes mulitaddr addresses from listen address urls
func (p *P2P) ListenAddrStrings() ([]string, error) {
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

// ProtocolID generates p2p protocol-id with the name and version information
func (p *P2P) ProtocolID() protocol.ID {
	cfg := p.Protocol
	// tmp := fmt.Sprintf("/%s/%s", cfg.ID, cfg.Version)
	tmp := fmt.Sprintf("/%s", cfg.ID)
	return protocol.ID(tmp)
}
