package main

import (
	nodebuilder "github.com/onflow/flow-go/cmd/access/node_builder"
)

func main() {
	anb := nodebuilder.FlowAccessNode(nodebuilder.UnstakedService())

	anb.PrintBuildVersionDetails()

	// parse all the command line args
	if err := anb.ParseFlags(); err != nil {
		anb.Logger.Fatal().Err(err).Send()
	}

	// choose a staked or an unstaked node builder based on anb.staked
	var builder nodebuilder.AccessNodeBuilder
	builder = nodebuilder.NewUnstakedAccessNodeBuilder(anb)

	if err := builder.Initialize(); err != nil {
		anb.Logger.Fatal().Err(err).Send()
	}

	node, err := builder.Build()
	if err != nil {
		anb.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
