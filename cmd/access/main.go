package main

import (
	nodebuilder "github.com/onflow/flow-go/cmd/access/node_builder"
)

func main() {
	builder := nodebuilder.FlowAccessNode() // use the generic Access Node builder till it is determined if this is a staked AN or an unstaked AN

	builder.PrintBuildVersionDetails()

	// parse all the command line args
	if err := builder.ParseFlags(); err != nil {
		builder.Logger.Fatal().Err(err).Send()
	}

	builder2 := builder
	if err := builder2.Initialize(); err != nil {
		builder2.Logger.Fatal().Err(err).Send()
	}

	node, err := builder2.Build()
	if err != nil {
		builder2.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
