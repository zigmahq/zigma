package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/config"
)

var p2pcfg = []byte(`
p2p:
  peer_id: 12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
  listen_addrs:
    - /ip4/0.0.0.0/tcp
    - /ip6/::/tcp
  bootstrap_addrs:
    - "/ip4/127.0.0.1/tcp/1000/ipfs/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM"
`)

func TestP2PUnmarshal(t *testing.T) {
	c, err := config.Read(p2pcfg)

	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.P2P)
	assert.Equal(t, "12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM", c.P2P.PeerID)
	assert.Len(t, c.P2P.ListenAddrs, 2)
	assert.Equal(t, "/ip4/0.0.0.0/tcp", c.P2P.ListenAddrs[0])
	assert.Equal(t, "/ip6/::/tcp", c.P2P.ListenAddrs[1])
	assert.Len(t, c.P2P.BootstrapAddrs, 1)
	assert.Equal(t, "/ip4/127.0.0.1/tcp/1000/ipfs/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM", c.P2P.BootstrapAddrs[0].MustMultiaddr().String())
}
