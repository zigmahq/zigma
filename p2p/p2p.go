package p2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-autonat-svc"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/zigmahq/zigma/config"
	"github.com/zigmahq/zigma/log"
	"golang.org/x/time/rate"

	circuit "github.com/libp2p/go-libp2p-circuit"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	pstoremem "github.com/libp2p/go-libp2p-peerstore/pstoremem"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	secio "github.com/libp2p/go-libp2p-secio"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery"
)

var logger = log.DefaultLogger

// P2P encapsulates a peer-to-peer server
type P2P struct {
	id               peer.ID
	cfg              *config.P2P
	ctx              context.Context
	host             host.Host
	connManager      *connmgr.BasicConnMgr
	privateKey       crypto.PrivKey
	pubsub           *pubsub.PubSub
	sub              *pubsub.Subscription
	kad              *dht.IpfsDHT
	networkNotifee   *networkNotifee
	discoveryNotifee *discoveryNotifee
	discovery        mdns.Service
	routingDiscovery discovery.Discovery
	limiter          *rate.Limiter
	itf              Implementer
}

// ID returns the server peer id
func (n *P2P) ID() string {
	return n.id.String()
}

// Host returns node basic host
func (n *P2P) Host() host.Host {
	return n.host
}

// Config returns P2P server configurations
func (n *P2P) Config() *config.P2P {
	return n.cfg
}

// Start runs the znode server
func (n *P2P) Start() error {
	logger.Info("Node id: " + n.ID())
	logger.Info(fmt.Sprintf("Using default protocol \"%s/%s\"", n.cfg.Protocol.ID, n.cfg.Protocol.Version))
	if addrs := n.cfg.Address; len(addrs) > 0 {
		logger.Info("Listening for new connections on:")
	}
	for _, addr := range n.cfg.Address {
		logger.Info("> " + addr.String())
	}

	// connect to seed nodes
	n.ConnectSeedNodes(n.ctx)

	if n.cfg.DHT.Enable {
		// starts background process of searching for peers
		if err := n.kad.Bootstrap(n.ctx); err != nil {
			return err
		}
		// advertise ourselves for the purposes of peer discovery.
		discovery.Advertise(n.ctx, n.routingDiscovery, n.cfg.DHT.Rendezvous, discovery.TTL(n.cfg.DHT.TTL))
	}

	return nil
}

// Stop shuts down the host, its Network, and services.
func (n *P2P) Stop() error {
	var fs []log.Field
	if n.cfg.DHT.Enable {
		fs = append(fs, log.Tag("kademlia"))
	}
	if n.cfg.MDNS.Enable {
		fs = append(fs, log.Tag("mdns"))
	}
	logger.Info("Stopping peer discovery", fs...)
	if n.discovery != nil {
		if err := n.discovery.Close(); err != nil {
			return err
		}
	}
	if err := n.kad.Close(); err != nil {
		return err
	}
	logger.Info("Stopping p2p server")
	if err := n.host.Close(); err != nil {
		return err
	}
	return nil
}

// NewServer initializes and returns a zigma node
func NewServer(ctx context.Context, p2pconf *config.P2P, itf Implementer) (*P2P, error) {
	priv, err := p2pconf.DecodePrivateKey()
	if err != nil {
		return nil, err
	}

	id, err := p2pconf.DecodePeerID()
	if err != nil {
		return nil, err
	}

	addrs, err := p2pconf.ListenAddrs()
	if err != nil {
		return nil, err
	}

	ps := pstoremem.NewPeerstore()
	ps.AddPrivKey(id, priv)
	ps.AddPubKey(id, priv.GetPublic())

	bwc := metrics.NewBandwidthCounter()
	cnm := connmgr.NewConnManager(p2pconf.MinNumConns, p2pconf.MaxNumConns, p2pconf.ConnGracePeriod)

	var kad *dht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		kad, err = dht.New(ctx, h, dhtopts.Protocols(p2pconf.ProtocolID()))
		if err != nil {
			logger.Error("unable to initialize dht: " + err.Error())
		}
		return kad, err
	}

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.UserAgent(p2pconf.Name),
		libp2p.Peerstore(ps),
		libp2p.ListenAddrs(addrs...),
		libp2p.ConnectionManager(cnm),
		libp2p.BandwidthReporter(bwc),
		libp2p.DefaultMuxers,
		libp2p.Routing(newDHT),
		libp2p.DefaultTransports,
	}

	if p2pconf.Insecure {
		opts = append(opts, libp2p.NoSecurity)
	} else {
		opts = append(opts, libp2p.Security(secio.ID, secio.New))
	}

	switch p2pconf.Relay {
	case config.RelayActive:
		opts = append(opts, libp2p.EnableRelay(circuit.OptActive, circuit.OptHop), libp2p.NATPortMap())
	case config.RelayNat:
		opts = append(opts, libp2p.EnableRelay(), libp2p.NATPortMap())
	default:
		opts = append(opts, libp2p.DisableRelay())
	}

	logger.Info("QUIC transport protocol", log.Bool("active", p2pconf.QUIC))
	if p2pconf.QUIC {
		opts = append(opts, libp2p.Transport(quic.NewTransport))
	}

	host, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		_ = host.Close()
	}()

	if _, err := autonat.NewAutoNATService(ctx, host); err != nil {
		return nil, err
	}

	p2p := &P2P{
		id:               id,
		privateKey:       priv,
		cfg:              p2pconf,
		ctx:              ctx,
		host:             host,
		connManager:      cnm,
		routingDiscovery: discovery.NewRoutingDiscovery(kad),
		kad:              kad,
		limiter:          rate.NewLimiter(rate.Limit(p2pconf.RateLimit.GlobalAvg), p2pconf.RateLimit.GlobalBurst),
		itf:              itf,
	}
	p2p.networkNotifee = &networkNotifee{p2p}
	p2p.discoveryNotifee = &discoveryNotifee{p2p}

	// initializes pubsub server
	p2p.pubsub, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	// listen to a pubsub topic
	p2p.sub, err = p2p.pubsub.Subscribe(p2pconf.Queue)
	if err != nil {
		return nil, err
	}

	// attach the stream handler
	host.SetStreamHandler(p2pconf.ProtocolID(), p2p.streamHandlerWrapper)

	// implemenet and attach network notifee interface to receive
	// notifications from a network.
	host.Network().Notify(p2p.networkNotifee)

	// mdns discovery service, if enabled, attach the discovery notifee interface
	// to receive notifications from mdns service
	logger.Info("mDNS discovery service", log.Bool("active", p2pconf.MDNS.Enable))
	if p2pconf.MDNS.Enable {
		p2p.discovery, err = mdns.NewMdnsService(ctx, host, p2pconf.MDNS.RescanInterval, p2pconf.MDNS.ServiceName)
		if err != nil {
			return nil, err
		}
		p2p.discovery.RegisterNotifee(p2p.discoveryNotifee)
	}

	return p2p, nil
}
