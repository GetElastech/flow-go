package main

import (
	"context"
	"flag"
	"time"

	"github.com/onflow/flow-go-sdk/client"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	accessAddr := ":9000"
	flag.StringVar(&accessAddr, "access-address", ":9000", "Access insecure gRPC server address")
	flag.Parse()
	log.Info().Msgf("Access node: %v", accessAddr)
	flowClient, err := client.New(accessAddr, grpc.WithInsecure()) //nolint:staticcheck
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize gRPC client")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := flowClient.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("failed to ping access node")
		return
	}
}

