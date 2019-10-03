package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-autonat-svc"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zigmahq/zigma/config"
	"github.com/zigmahq/zigma/log"
	"golang.org/x/time/rate"

	circuit "github.com/libp2p/go-libp2p-circuit"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	rdisc "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	pstoremem "github.com/libp2p/go-libp2p-peerstore/pstoremem"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	secio "github.com/libp2p/go-libp2p-secio"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery"
)

// PubSubFn function type for pubsub
type PubSubFn func(ctx context.Context, h host.Host, opts ...pubsub.Option) (*pubsub.PubSub, error)

// P2P encapsulates a peer-to-peer server
type P2P struct {
	id               peer.ID
	privateKey       crypto.PrivKey
	cfg              *config.P2P
	ctx              context.Context
	host             host.Host
	close            chan struct{}
	peers            int
	online           bool
	logger           log.Logger
	autonat          *autonat.AutoNATService
	psfn             PubSubFn
	pubs             *pubsub.PubSub
	subs             *pubsub.Subscription
	discovery        mdns.Service
	dht              *dht.IpfsDHT
	networkNotifee   *networkNotifee
	discoveryNotifee *discoveryNotifee
	limiter          *rate.Limiter
}

// ID returns the server peer id
func (n *P2P) ID() string {
	return n.id.String()
}

// Host returns node basic host
func (n *P2P) Host() host.Host {
	return n.host
}

// Start runs the znode server
func (n *P2P) Start() error {
	n.logger.Info("Node id: " + n.ID())
	n.logger.Info("Using default protocol \"" + string(ZProtocolID) + "\"")
	if addrs := n.cfg.Address; len(addrs) > 0 {
		n.logger.Info("Listening for new connections on:")
	}
	for _, addr := range n.cfg.Address {
		n.logger.Info("> " + addr.String())
	}
	n.online = true
	return nil
}

// Stop shuts down the host, its Network, and services.
func (n *P2P) Stop() error {
	close(n.close)
	n.online = false
	fs := []log.Field{log.Tag("kademlia")}
	if n.discovery != nil {
		fs = append(fs, log.Tag("mdns"))
	}
	n.logger.Info("Stopping peer discovery", fs...)
	if n.discovery != nil {
		if err := n.discovery.Close(); err != nil {
			return err
		}
	}
	if err := n.dht.Close(); err != nil {
		return err
	}
	n.logger.Info("Stopping p2p server")
	if err := n.host.Close(); err != nil {
		return err
	}
	return nil
}

// JoinOverlay triggers the host to join the DHT overlay
func (n *P2P) JoinOverlay(ctx context.Context) {
	rd := rdisc.NewRoutingDiscovery(n.dht)
	rdisc.Advertise(ctx, rd, string(ZProtocolID))
}

// NewServer initializes and returns a zigma node
func NewServer(ctx context.Context, logger log.Logger, p2pconf *config.P2P) (*P2P, error) {
	priv, err := p2pconf.DecodePrivateKey()
	if err != nil {
		return nil, err
	}

	id, err := p2pconf.DecodePeerID()
	if err != nil {
		return nil, err
	}

	addrs, err := p2pconf.DecodeListenAddrs()
	if err != nil {
		return nil, err
	}

	ps := pstoremem.NewPeerstore()
	ps.AddPrivKey(id, priv)
	ps.AddPubKey(id, priv.GetPublic())

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.UserAgent(p2pconf.Name),
		libp2p.Security(secio.ID, secio.New),
		libp2p.Peerstore(ps),
		libp2p.ListenAddrs(addrs...),
		libp2p.ConnectionManager(connmgr.NewConnManager(p2pconf.MinNumConns, p2pconf.MaxNumConns, p2pconf.ConnGracePeriod)),
		libp2p.DefaultMuxers,
		libp2p.DefaultTransports,
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

	dht, err := dht.New(ctx, host, dhtopts.Protocols(ZProtocolID))
	if err != nil {
		return nil, err
	}
	if err := dht.Bootstrap(ctx); err != nil {
		return nil, err
	}

	autonat, err := autonat.NewAutoNATService(ctx, host)
	if err != nil {
		return nil, err
	}

	p2p := &P2P{
		id:         id,
		privateKey: priv,
		cfg:        p2pconf,
		ctx:        ctx,
		host:       host,
		close:      make(chan struct{}),
		peers:      0,
		online:     false,
		logger:     logger,
		autonat:    autonat,
		pubs:       new(pubsub.PubSub),
		subs:       new(pubsub.Subscription),
		dht:        dht,
		limiter:    rate.NewLimiter(rate.Limit(p2pconf.RateLimit.GlobalAvg), p2pconf.RateLimit.GlobalBurst),
	}
	p2p.networkNotifee = &networkNotifee{p2p}
	p2p.discoveryNotifee = &discoveryNotifee{p2p}

	if p2pconf.Gossip {
		p2p.psfn = pubsub.NewFloodSub
	} else {
		p2p.psfn = pubsub.NewGossipSub
	}

	// attach the stream handler
	host.SetStreamHandler(ZProtocolID, p2p.handleStream)

	// implemenet and attach network notifee interface to receive
	// notifications from a network.
	host.Network().Notify(p2p.networkNotifee)

	// mdns discovery service, if enabled, attach the discovery notifee interface
	// to receive notifications from mdns service
	logger.Info("mDNS discovery service", log.Bool("active", p2pconf.MDNS.Enable))
	if p2pconf.MDNS.Enable {
		discovery, err := mdns.NewMdnsService(ctx, host, p2pconf.MDNS.RescanInterval, ZMDNSServiceName)
		if err != nil {
			return nil, err
		}
		discovery.RegisterNotifee(p2p.discoveryNotifee)
		p2p.discovery = discovery
	}

	return p2p, nil
}
