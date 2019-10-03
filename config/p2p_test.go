package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/config"
)

var p2pcfg = []byte(`
p2p:
  priv_key: config/node_priv.key
  gossip: true
  relay: true
  min_num_conns: 200
  max_num_conns: 1000
  laddr: tcp://0.0.0.0:0
  persistent_peers:
    - tcp://0:0/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
  seeds:
    - tcp://0:0
    - tcp://0:0/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
    - tcp://0.0.0.0:0/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
    - tcp://127.0.0.1:1000/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
    - ws://0:0/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
    - ws://:::0/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
    - tcp://:::8000/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
    - tcp://::2:0/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
    - tcp://::2:0/12D
    - tcp://0:0:0/12D3KooWHvpXAvv68Pjysw2jYNSY4feWRCYRBUw9wJSeg2PRa7BM
`)

func TestP2PUnmarshal(t *testing.T) {
	c, err := config.Read(p2pcfg)

	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.P2P)
	assert.NotPanics(t, func() {
		_ = c.P2P.Seeds[0].MustMultiaddr()
		_ = c.P2P.Seeds[1].MustMultiaddr()
		_ = c.P2P.Seeds[2].MustMultiaddr()
		_ = c.P2P.Seeds[3].MustMultiaddr()
		_ = c.P2P.Seeds[4].MustMultiaddr()
		_ = c.P2P.Seeds[5].MustMultiaddr()
		_ = c.P2P.Seeds[6].MustMultiaddr()
	})
	assert.Panics(t, func() {
		_ = c.P2P.Seeds[7].MustMultiaddr()
		_ = c.P2P.Seeds[8].MustMultiaddr()
	})
}
