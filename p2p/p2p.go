package p2p

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-autonat-svc"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/zigmahq/zigma/config"

	circuit "github.com/libp2p/go-libp2p-circuit"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	routing "github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pstoremem "github.com/libp2p/go-libp2p-peerstore/pstoremem"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	psrouter "github.com/libp2p/go-libp2p-pubsub-router"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

// ZNode encapsulates a zigma peer-to-peer node
type ZNode struct {
	// the node identifier
	id peer.ID
	// private key for encrypted communication and verifying identity
	privateKey crypto.PrivKey
	// p2p node configurations
	cfg *config.P2P
	// the base context for the p2p node
	ctx context.Context
	// host for p2p connections
	host host.Host
	// online state
	online           bool
	routing          routing.Routing
	routedhost       *rhost.RoutedHost
	autonat          *autonat.AutoNATService
	discovery        discovery.Service
	pubsub           *pubsub.PubSub
	psrouter         *psrouter.PubsubValueStore
	dht              *dht.IpfsDHT
	networkNotifee   *networkNotifee
	discoveryNotifee *discoveryNotifee
}

// ID returns the node peer id
func (n *ZNode) ID() string {
	return n.id.String()
}

// Close shuts down the host, its Network, and services.
func (n *ZNode) Close() error {
	n.online = false
	return n.host.Close()
}

func (n *ZNode) handleStream(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go func() {
		for {
			s, err := rw.ReadString('\n')
			if err != nil {
				continue
			}
			if len(s) == 0 {
				continue
			}
			if s != "\n" {
				fmt.Printf("\x1b[32m%s\x1b[0m>", s)
			}
		}
	}()
}

func (n *ZNode) bootstrapPeerDiscovery() error {
	if len(n.cfg.BootstrapAddrs) == 0 {
		return nil
	}
	var errs []error
	var wg sync.WaitGroup
	for _, addr := range n.cfg.BootstrapAddrs {
		wg.Add(1)
		go func(addr config.P2PMultiAddr) {
			defer wg.Done()
			m, err := addr.Multiaddr()
			if err != nil {
				return
			}
			p, err := peer.AddrInfoFromP2pAddr(m)
			if err != nil {
				return
			}
			if err := n.host.Connect(n.ctx, *p); err != nil {
				errs = append(errs, err)
			}
		}(addr)
	}
	wg.Wait()
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (n *ZNode) bootstrapDHT() error {
	return n.dht.Bootstrap(n.ctx)
}

func (n *ZNode) boostrapMDNS() error {
	discovery, err := discovery.NewMdnsService(n.ctx, n.host, time.Second*5, ZMDNSServiceName)
	if err != nil {
		return err
	}
	n.discovery = discovery

	// Implement and attach discovery notifee interface to receive
	// notifications from mdns service
	n.discovery.RegisterNotifee(n.discoveryNotifee)
	return nil
}

// NewZNode initializes and returns a zigma node
func NewZNode(ctx context.Context, p2pconf *config.P2P) (*ZNode, error) {
	node := &ZNode{
		ctx: ctx,
		cfg: p2pconf,
	}
	node.networkNotifee = &networkNotifee{node}
	node.discoveryNotifee = &discoveryNotifee{node}

	id, err := p2pconf.DecodePeerID()
	if err != nil {
		return nil, err
	}
	node.id = id

	priv, err := p2pconf.DecodePrivateKey()
	if err != nil {
		return nil, err
	}
	node.privateKey = priv

	ps := pstoremem.NewPeerstore()
	ps.AddPrivKey(node.id, priv)
	ps.AddPubKey(node.id, priv.GetPublic())

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.Security(secio.ID, secio.New),
		libp2p.Peerstore(ps),
		libp2p.ListenAddrStrings(p2pconf.Addrs()...),
		libp2p.EnableRelay(circuit.OptHop, circuit.OptDiscovery),
		libp2p.ConnectionManager(connmgr.NewConnManager(1000, 4000, time.Minute)),
		libp2p.ChainOptions(
			libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
			libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
		),
	}

	node.host, err = libp2p.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// Implemenet and attach network notifee interface to receive
	// notifications from a Network.
	node.host.Network().Notify(node.networkNotifee)
	log.Printf("# %s/p2p/%s\n", node.host.Addrs()[0], node.ID())

	node.autonat, err = autonat.NewAutoNATService(ctx, node.host)
	if err != nil {
		return nil, err
	}

	node.host.SetStreamHandler(ZProtocolID, node.handleStream)
	node.online = true

	node.dht, err = dht.New(ctx, node.host)
	if err != nil {
		return nil, err
	}

	node.routedhost = rhost.Wrap(node.host, node.dht)

	node.bootstrapPeerDiscovery()
	node.bootstrapDHT()
	node.boostrapMDNS()

	return node, nil
}
