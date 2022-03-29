package main

import (
	"context"
	"flag"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	accessAddr := ":9000"
	flag.StringVar(&accessAddr, "access-node-addr", ":9000", "Access insecure gRPC server address")
	flag.Parse()
	log.Info().Msgf("Access node: %v", accessAddr)
	flowClient, err := client.New(accessAddr, grpc.WithInsecure()) //nolint:staticcheck
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize gRPC client")
		return
	}
	pctx, pcancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pcancel()
	if err := flowClient.Ping(pctx); err != nil {
		log.Error().Err(err).Msg("failed to ping access node")
		return
	}
	sealedBlock, err := getLatestBlock(flowClient, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to get latest sealed block")
		return
	}
	log.Info().Msgf("Sealed block: %+v", sealedBlock)
	unsealedBlock, err := getLatestBlock(flowClient, false)
	if err != nil {
		log.Error().Err(err).Msg("failed to get latest unsealed block")
		return
	}
	log.Info().Msgf("Unsealed block: %+v", unsealedBlock)
}

func getLatestBlock(flowClient *client.Client, isSealed bool) (*flow.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return flowClient.GetLatestBlock(ctx, true)
}

