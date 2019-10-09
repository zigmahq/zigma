package node

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/zigmahq/zigma/config"
	"github.com/zigmahq/zigma/log"
	"github.com/zigmahq/zigma/p2p"
)

var logger = log.DefaultLogger

// Node encapsulates a znode server
type Node struct {
	*p2p.P2P
	ctx  context.Context
	stop chan os.Signal
}

// Start starts the node and p2p services
func (n *Node) Start() error {
	if err := n.P2P.Start(); err != nil {
		return err
	}
	signal.Notify(n.stop, syscall.SIGINT)
	<-n.stop
	logger.NL()
	return nil
}

// Stop to stop all the active services
func (n *Node) Stop() error {
	return n.P2P.Stop()
}

// NewNode initializes and returns a zigma node
func NewNode(ctx context.Context, conf *config.Config) (*Node, error) {
	n := &Node{
		ctx:  ctx,
		stop: make(chan os.Signal, 1),
	}

	p, err := p2p.NewServer(ctx, conf.P2P, n)
	if err != nil {
		return nil, err
	}
	n.P2P = p

	return n, nil
}
