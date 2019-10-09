package main

import (
	"context"

	"github.com/zigmahq/zigma/config"
	"github.com/zigmahq/zigma/log"
	"github.com/zigmahq/zigma/node"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewLogger()
	logger.Info("Zigma Node",
		log.String("version", "1.0.0-3687"),
		log.String("by", "Nakama Hiroyasu"),
	)
	logger.Info("Chain",
		log.String("specification", "Local Testnet"),
		log.String("node-name", "pretty-dress-0716"),
	)
	cfg := config.DefaultConfig()

	node, err := node.NewNode(ctx, cfg)
	if err != nil {
		panic(err)
	}

	if err := node.Start(); err != nil {
		panic(err)
	}
	if err := node.Stop(); err != nil {
		panic(err)
	}
}
