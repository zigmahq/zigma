package p2p

import (
	"context"
	"math"
	"math/rand"

	"github.com/zigmahq/zigma/config/types"
)

func randomPeerAddresses(in []*types.Addr, max int) []*types.Addr {
	n := int(math.Min(float64(max), float64(len(in))))
	var out []*types.Addr
	for _, val := range rand.Perm(len(in)) {
		out = append(out, in[val])
		if len(out) >= n {
			break
		}
	}
	return out
}

// ConnectSeedNodes to connect seed nodes
func (n *P2P) ConnectSeedNodes(ctx context.Context) {
	peers := randomPeerAddresses(n.cfg.Seeds, 4)
	n.ConnectWithAddr(peers...)
}
