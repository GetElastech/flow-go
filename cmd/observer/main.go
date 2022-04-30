package main

import (
	observer "github.com/onflow/flow-go/cmd/observer/node_builder"
)

func main() {
	anb := observer.FlowAccessNode() // use the generic Access Node builder till it is determined if this is a staked AN or an unstaked AN

	anb.PrintBuildVersionDetails()

	// parse all the command line args
	if err := anb.ParseFlags(); err != nil {
		anb.Logger.Fatal().Err(err).Send()
	}

	// choose a staked or an unstaked node builder based on anb.staked
	var builder observer.AccessNodeBuilder
	if anb.IsStaked() {
		builder = observer.NewUnstakedAccessNodeBuilder(anb)
	} else {
		builder = observer.NewUnstakedAccessNodeBuilder(anb)
	}

	if err := builder.Initialize(); err != nil {
		anb.Logger.Fatal().Err(err).Send()
	}

	node, err := builder.Build()
	if err != nil {
		anb.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
