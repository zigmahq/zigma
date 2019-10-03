package types

import (
	"net"
	"net/url"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

// Addrs type
type Addrs []*Addr

// UnmarshalYAML hook
func (a *Addrs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var addrs []*Addr
	err := unmarshal(&addrs)
	if err == nil {
		*a = addrs
		return nil
	}
	var addr *Addr
	err = unmarshal(&addr)
	if err == nil {
		*a = []*Addr{addr}
		return nil
	}
	return err
}

// PeerInfos converts list of address objects to peerstore.PeerInfo objects
func (a Addrs) PeerInfos() []peerstore.PeerInfo {
	var addrs []multiaddr.Multiaddr
	for _, addr := range a {
		m, err := addr.Multiaddr()
		if err != nil {
			continue
		}
		addrs = append(addrs, m)
	}
	i, err := peer.AddrInfosFromP2pAddrs(addrs...)
	if err != nil {
		return nil
	}
	return i
}

// Addr type
type Addr struct {
	*url.URL
}

// UnmarshalYAML hook
func (a *Addr) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	a.URL = u
	return nil
}

// MarshalYAML hook
func (a *Addr) MarshalYAML() (interface{}, error) {
	return a.URL.String(), nil
}

// Protocol returns address protocol name
func (a *Addr) Protocol() string {
	hn := a.URL.Hostname()
	ip := net.ParseIP(hn)
	switch {
	case hn == "0", ip != nil && strings.Contains(hn, "."):
		return "ip4"
	case ip != nil && strings.Contains(hn, ":"):
		return "ip6"
	case len(strings.Split(hn, ".")) > 1:
		return "dnsaddr"
	}
	panic("unknown protocol")
}

// Transport extracts the transport layer name
func (a *Addr) Transport() []string {
	switch a.URL.Scheme {
	case "tcp", "http", "https":
		return []string{"tcp"}
	default:
		return []string{"tcp", a.URL.Scheme}
	}
}

// Hostname returns the hostname of address
func (a *Addr) Hostname() string {
	hn := a.URL.Hostname()
	switch hn {
	case "0":
		return "0.0.0.0"
	default:
		return hn
	}
}

// Port returns the port number of address
func (a *Addr) Port() string {
	if port := a.URL.Port(); len(port) != 0 {
		return port
	}
	return "0"
}

// Multiaddr returns address in multiaddr format
func (a *Addr) Multiaddr() (multiaddr.Multiaddr, error) {
	var tmp []string
	tps := a.Transport()

	tmp = append(tmp, a.Protocol(), a.Hostname(), tps[0], a.Port())
	if len(tps) > 1 {
		tmp = append(tmp, tps[1])
	}
	if len(a.URL.Path) > 1 {
		tmp = append(tmp, "p2p")
		tmp = append(tmp, a.URL.Path[1:])
	}

	txt := "/" + strings.Join(tmp, "/")
	return multiaddr.NewMultiaddr(txt)
}

// MustMultiaddr returns address in multiaddr format, program should go panic
// if the address format is invalid
func (a *Addr) MustMultiaddr() multiaddr.Multiaddr {
	m, err := a.Multiaddr()
	if err != nil {
		panic(err)
	}
	return m
}

// NewAddr initializes an address
func NewAddr(s string) *Addr {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return &Addr{u}
}

func init() {
	multiaddr.SwapToP2pMultiaddrs()
}
