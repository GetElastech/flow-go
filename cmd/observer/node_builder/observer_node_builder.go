package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/onflow/flow-go/follower"
	"github.com/onflow/flow-go/module/local"
	"github.com/onflow/flow-go/network/converter"
	"github.com/onflow/flow-go/network/p2p/keyutils"
	"github.com/onflow/flow-go/network/validator"
	"github.com/onflow/flow-go/state/protocol/events/gadgets"
	"github.com/onflow/flow-go/utils/io"
	"github.com/rs/zerolog"
	"strings"

	"github.com/onflow/flow-go/apiservice"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/onflow/flow/protobuf/go/flow/access"

	"github.com/onflow/flow-go/crypto"

	"github.com/onflow/flow-go/cmd"
	"github.com/onflow/flow-go/engine"
	pingeng "github.com/onflow/flow-go/engine/access/ping"
	"github.com/onflow/flow-go/engine/access/rpc"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/id"
	"github.com/onflow/flow-go/module/mempool/stdmap"
	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/network"
	netcache "github.com/onflow/flow-go/network/cache"
	"github.com/onflow/flow-go/network/p2p"
	"github.com/onflow/flow-go/network/p2p/unicast"
	relaynet "github.com/onflow/flow-go/network/relay"
	"github.com/onflow/flow-go/network/topology"
	"github.com/onflow/flow-go/utils/grpcutils"
)

// ObserverNodeBuilder builds a staked access node for simulation for now.
// Observer services can participate in the observer network
// publishing data for the observer node downstream.
type ObserverNodeBuilder struct {
	*FlowAccessNodeBuilder
	peerID peer.ID
}

func NewObserverNodeBuilder(builder *FlowAccessNodeBuilder) *ObserverNodeBuilder {
	return &ObserverNodeBuilder{
		FlowAccessNodeBuilder: builder,
	}
}

func (builder *ObserverNodeBuilder) InitIDProviders() {
	builder.Module("id providers", func(node *cmd.NodeConfig) error {
		idCache, err := p2p.NewProtocolStateIDCache(node.Logger, node.State, node.ProtocolEvents)
		if err != nil {
			return err
		}

		builder.IdentityProvider = idCache

		builder.SyncEngineParticipantsProviderFactory = func() id.IdentifierProvider {
			return id.NewIdentityFilterIdentifierProvider(
				filter.And(
					filter.HasRole(flow.RoleConsensus),
					filter.Not(filter.HasNodeID(node.Me.NodeID())),
					p2p.NotEjectedFilter,
				),
				idCache,
			)
		}

		builder.IDTranslator = p2p.NewHierarchicalIDTranslator(idCache, p2p.NewUnstakedNetworkIDTranslator())

		return nil
	})
}

func (builder *ObserverNodeBuilder) Initialize() error {
	builder.InitIDProviders()

	builder.EnqueueResolver()

	// enqueue the regular network
	builder.EnqueueNetworkInit()

	builder.EnqueuePingService()

	builder.EnqueueMetricsServerInit()

	if err := builder.RegisterBadgerMetrics(); err != nil {
		return err
	}

	builder.EnqueueTracer()
	builder.PreInit(cmd.DynamicStartPreInit)
	return nil
}

func (builder *ObserverNodeBuilder) enqueueRelayNetwork() {
	builder.Component("relay network", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		relayNet := relaynet.NewRelayNetwork(
			node.Network,
			builder.ObserverServiceConfig.PublicNetworkConfig.Network,
			node.Logger,
			[]network.Channel{engine.ReceiveBlocks},
		)
		node.Network = relayNet
		return relayNet, nil
	})
}

func (builder *ObserverNodeBuilder) Build() (cmd.Node, error) {
	builder.
		//BuildConsensusFollower().                                           //Remove
		//Module("collection node client", func(node *cmd.NodeConfig) error { //Remove
		//	// collection node address is optional (if not specified, collection nodes will be chosen at random)
		//	if strings.TrimSpace(builder.rpcConf.CollectionAddr) == "" {
		//		node.Logger.Info().Msg("using a dynamic collection node address")
		//		return nil
		//	}
		//
		//	node.Logger.Info().
		//		Str("collection_node", builder.rpcConf.CollectionAddr).
		//		Msg("using the static collection node address")
		//
		//	collectionRPCConn, err := grpc.Dial(
		//		builder.rpcConf.CollectionAddr,
		//		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcutils.DefaultMaxMsgSize)),
		//		grpc.WithInsecure(), //nolint:staticcheck
		//		backend.WithClientUnaryInterceptor(builder.rpcConf.CollectionClientTimeout))
		//	if err != nil {
		//		return err
		//	}
		//	builder.CollectionRPC = access.NewAccessAPIClient(collectionRPCConn)
		//	return nil
		//}).
		Module("observer public p2p collection", func(node *cmd.NodeConfig) error {
			if err := builder.deriveBootstrapPeerIdentities(); err != nil {
				return err
			}

			if builder.BaseConfig.BindAddr == cmd.NotSet {
				node.Logger.Info().Msg(fmt.Sprintf("using public libp2p network bind address %s", builder.PublicNetworkConfig.BindAddress))
				builder.BaseConfig.BindAddr = builder.PublicNetworkConfig.BindAddress
			}
			// TODO
			//if err := builder.validateParams(); err != nil {
			//	return err
			//}

			if err := builder.initPublicNodeInfo(); err != nil {
				return err
			}

			builder.InitIDProviders()

			builder.enqueueMiddleware()

			// TODO middleware conflict
			builder.enqueueUnstakedNetworkInit()

			if err := builder.enqueueConnectWithStakedAN(); err != nil {
				return err
			}

			//if builder.BaseConfig.MetricsEnabled {
			//	builder.EnqueueMetricsServerInit()
			//	if err := builder.RegisterBadgerMetrics(); err != nil {
			//		return err
			//	}
			//}

			builder.PreInit(builder.initUnstakedLocal())

			return nil
		}).
		Module("observer proxy", func(node *cmd.NodeConfig) error {
			// collection node address is optional (if not specified, collection nodes will be chosen at random)
			if len(builder.ObserverServiceConfig.bootstrapNodeAddresses) == 0 {
				node.Logger.Info().Msg("using a dynamic collection node address")
				return nil
			}

			err := builder.deriveBootstrapPeerIdentities()
			if err != nil {
				node.Logger.Info().Msg(fmt.Sprintf("cannot collect observer bootstrap peer identities %s", err.Error()))
				return nil
			}

			for _, identity := range builder.bootstrapIdentities {
				node.Logger.Info().
					Str("observer_node", builder.rpcConf.UnsecureGRPCListenAddr).
					Msg(fmt.Sprintf("proxying %s to %s with key %s ...", builder.rpcConf.UnsecureGRPCListenAddr, identity.Address, identity.NetworkPubKey))
			}
			proxy, err := apiservice.NewFlowAPIService(builder.bootstrapIdentities, builder.rpcConf.CollectionClientTimeout)
			if err != nil {
				node.Logger.Error().
					Str("observer_node", err.Error()).
					Msg(fmt.Sprintf("proxying %s failed", builder.rpcConf.UnsecureGRPCListenAddr))
				return err
			}

			if proxy == nil {
				node.Logger.Info().Msg("cannot create flow proxy")
				return nil
			}
			builder.Proxy = proxy

			return nil
		}).
		Module("historical access node clients", func(node *cmd.NodeConfig) error {
			addrs := strings.Split(builder.rpcConf.HistoricalAccessAddrs, ",")
			for _, addr := range addrs {
				if strings.TrimSpace(addr) == "" {
					continue
				}
				node.Logger.Info().Str("access_nodes", addr).Msg("historical access node addresses")

				historicalAccessRPCConn, err := grpc.Dial(
					addr,
					grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcutils.DefaultMaxMsgSize)),
					grpc.WithInsecure()) //nolint:staticcheck
				if err != nil {
					return err
				}
				builder.HistoricalAccessRPCs = append(builder.HistoricalAccessRPCs, access.NewAccessAPIClient(historicalAccessRPCConn))
			}
			return nil
		}).
		Module("transaction timing mempools", func(node *cmd.NodeConfig) error {
			var err error
			builder.TransactionTimings, err = stdmap.NewTransactionTimings(1500 * 300) // assume 1500 TPS * 300 seconds
			if err != nil {
				return err
			}

			builder.CollectionsToMarkFinalized, err = stdmap.NewTimes(50 * 300) // assume 50 collection nodes * 300 seconds
			if err != nil {
				return err
			}

			builder.CollectionsToMarkExecuted, err = stdmap.NewTimes(50 * 300) // assume 50 collection nodes * 300 seconds
			if err != nil {
				return err
			}

			builder.BlocksToMarkExecuted, err = stdmap.NewTimes(1 * 300) // assume 1 block per second * 300 seconds
			return err
		}).
		Module("transaction metrics", func(node *cmd.NodeConfig) error {
			builder.TransactionMetrics = metrics.NewTransactionCollector(builder.TransactionTimings, node.Logger, builder.logTxTimeToFinalized,
				builder.logTxTimeToExecuted, builder.logTxTimeToFinalizedExecuted)
			return nil
		}).
		Module("ping metrics", func(node *cmd.NodeConfig) error {
			builder.PingMetrics = metrics.NewPingCollector()
			return nil
		}).
		Module("server certificate", func(node *cmd.NodeConfig) error {
			node.Logger.Info().Msg(fmt.Sprintf("grpc public key %s", hex.EncodeToString(node.NetworkKey.PublicKey().Encode())))
			// generate the server certificate that will be served by the GRPC server
			x509Certificate, err := grpcutils.X509Certificate(node.NetworkKey)
			if err != nil {
				return err
			}
			tlsConfig := grpcutils.DefaultServerTLSConfig(x509Certificate)
			builder.rpcConf.TransportCredentials = credentials.NewTLS(tlsConfig)
			return nil
		}).
		//Module("public network", func(node *cmd.NodeConfig) error {
		//	builder.enqueuePublicNetworkInit()
		//	return nil
		//}).
		//Module("public network", func(node *cmd.NodeConfig) error { //Remove
		//	builder.enqueueRelayNetwork()
		//	return nil
		//}).
		Component("RPC engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
			builder.RpcEng = rpc.New(
				node.Logger,
				node.State,
				builder.rpcConf,
				builder.CollectionRPC,
				builder.HistoricalAccessRPCs,
				node.Storage.Blocks,
				node.Storage.Headers,
				node.Storage.Collections,
				node.Storage.Transactions,
				node.Storage.Receipts,
				node.Storage.Results,
				node.RootChainID,
				builder.TransactionMetrics,
				builder.collectionGRPCPort,
				builder.executionGRPCPort,
				builder.retryEnabled,
				builder.rpcMetricsEnabled,
				builder.apiRatelimits,
				builder.apiBurstlimits,
				builder.Proxy,
			)
			return builder.RpcEng, nil
		})
	//.Component("ingestion engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
	//	var err error
	//
	//	builder.RequestEng, err = requester.New(
	//		node.Logger,
	//		node.Metrics.Engine,
	//		node.Network,
	//		node.Me,
	//		node.State,
	//		engine.RequestCollections,
	//		filter.HasRole(flow.RoleCollection),
	//		func() flow.Entity { return &flow.Collection{} },
	//	)
	//	if err != nil {
	//		return nil, fmt.Errorf("could not create requester engine: %w", err)
	//	}
	//
	//	builder.IngestEng, err = ingestion.New(node.Logger, node.Network, node.State, node.Me, builder.RequestEng, node.Storage.Blocks, node.Storage.Headers, node.Storage.Collections, node.Storage.Transactions, node.Storage.Results, node.Storage.Receipts, builder.TransactionMetrics,
	//		builder.CollectionsToMarkFinalized, builder.CollectionsToMarkExecuted, builder.BlocksToMarkExecuted, builder.RpcEng)
	//	if err != nil {
	//		return nil, err
	//	}
	//	builder.RequestEng.WithHandle(builder.IngestEng.OnCollection)
	//	builder.FinalizationDistributor.AddConsumer(builder.IngestEng)
	//
	//	return builder.IngestEng, nil
	//}).
	//Component("requester engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
	//	// We initialize the requester engine inside the ingestion engine due to the mutual dependency. However, in
	//	// order for it to properly start and shut down, we should still return it as its own engine here, so it can
	//	// be handled by the scaffold.
	//	return builder.RequestEng, nil
	//})

	builder.Component("ping engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		ping, err := pingeng.New(
			node.Logger,
			node.IdentityProvider,
			node.IDTranslator,
			node.Me,
			builder.PingMetrics,
			builder.pingEnabled,
			builder.nodeInfoFile,
			node.PingService,
		)

		if err != nil {
			return nil, fmt.Errorf("could not create ping engine: %w", err)
		}

		return ping, nil
	})

	return builder.FlowAccessNodeBuilder.Build()
}

func (builder *ObserverNodeBuilder) initPublicNodeInfo() error {
	// use the networking key that has been passed in the config, or load from the configured file
	builder.Logger.Info().Msg(fmt.Sprintf("loading public networking key from: %s\n", builder.observerNetworkingKeyPath))
	var err error
	networkingKey, err := loadNetworkingKey(builder.observerNetworkingKeyPath)
	if err != nil {
		return fmt.Errorf("could not load networking private key: %w", err)
	}

	pubKey, err := keyutils.LibP2PPublicKeyFromFlow(networkingKey.PublicKey())
	if err != nil {
		return fmt.Errorf("could not load networking public key: %w", err)
	}

	builder.peerID, err = peer.IDFromPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("could not get peer ID from public key: %w", err)
	}

	builder.NodeID, err = p2p.NewUnstakedNetworkIDTranslator().GetFlowID(builder.peerID)
	if err != nil {
		return fmt.Errorf("could not get flow node ID: %w", err)
	}

	builder.NodeConfig.NetworkKey = networkingKey // copy the key to NodeConfig
	builder.NodeConfig.StakingKey = nil           // no staking key for the unstaked node

	return nil
}

// enqueuePublicNetworkInit enqueues the public network component initialized for the staked node
func (builder *ObserverNodeBuilder) enqueuePublicNetworkInit() {
	// We still refer to "public network" here to reflect what the access node refers to
	// Observer implementations later should use a different name "public network"
	builder.Component("public network", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		builder.PublicNetworkConfig.Metrics = metrics.NewNetworkCollector(metrics.WithNetworkPrefix("unstaked"))

		libP2PFactory := builder.initLibP2PFactory(builder.NodeConfig.NetworkKey)

		msgValidators := observerNetworkMsgValidators(node.Logger.With().Bool("unstaked", true).Logger(), node.IdentityProvider, builder.NodeID)

		middleware := builder.initMiddleware(builder.NodeID, builder.PublicNetworkConfig.Metrics, libP2PFactory, msgValidators...)

		// topology returns empty list since peers are not known upfront
		top := topology.EmptyListTopology{}

		var heroCacheCollector module.HeroCacheMetrics = metrics.NewNoopCollector()
		if builder.HeroCacheMetricsEnable {
			heroCacheCollector = metrics.PublicNetworkReceiveCacheMetricsFactory(builder.MetricsRegisterer)
		}
		receiveCache := netcache.NewHeroReceiveCache(builder.NetworkReceivedMessageCacheSize,
			builder.Logger,
			heroCacheCollector)

		err := node.Metrics.Mempool.Register(metrics.ResourcePublicNetworkingReceiveCache, receiveCache.Size)
		if err != nil {
			return nil, fmt.Errorf("could not register networking receive cache metric: %w", err)
		}

		net, err := builder.initNetwork(builder.Me, builder.PublicNetworkConfig.Metrics, middleware, top, receiveCache)
		if err != nil {
			return nil, err
		}

		builder.ObserverServiceConfig.PublicNetworkConfig.Network = net

		node.Logger.Info().Msgf("public network will run on address: %s", builder.PublicNetworkConfig.BindAddress)
		return net, nil
	})
}

// initLibP2PFactory creates the LibP2P factory function for the given node ID and network key.
// The factory function is later passed into the initMiddleware function to eventually instantiate the p2p.LibP2PNode instance
// The LibP2P host is created with the following options:
// 		DHT as server
// 		The address from the node config or the specified bind address as the listen address
// 		The passed in private key as the libp2p key
//		No connection gater
// 		Default Flow libp2p pubsub options
func (builder *ObserverNodeBuilder) initLibP2PFactory(networkKey crypto.PrivateKey) p2p.LibP2PFactoryFunc {
	return func(ctx context.Context) (*p2p.Node, error) {
		connManager := p2p.NewConnManager(builder.Logger, builder.PublicNetworkConfig.Metrics)

		libp2pNode, err := p2p.NewNodeBuilder(builder.Logger, builder.PublicNetworkConfig.BindAddress, networkKey, builder.SporkID).
			SetBasicResolver(builder.Resolver).
			SetSubscriptionFilter(
				p2p.NewRoleBasedFilter(
					flow.RoleAccess, builder.IdentityProvider,
				),
			).
			SetConnectionManager(connManager).
			SetRoutingSystem(func(ctx context.Context, h host.Host) (routing.Routing, error) {
				return p2p.NewDHT(ctx, h, unicast.FlowPublicDHTProtocolID(builder.SporkID), p2p.AsServer())
			}).
			SetPubSub(pubsub.NewGossipSub).
			Build(ctx)

		if err != nil {
			return nil, fmt.Errorf("could not build libp2p node for staked access node: %w", err)
		}

		builder.LibP2PNode = libp2pNode

		return builder.LibP2PNode, nil
	}
}

// initMiddleware creates the network.Middleware implementation with the libp2p factory function, metrics, peer update
// interval, and validators. The network.Middleware is then passed into the initNetwork function.
func (builder *ObserverNodeBuilder) initMiddleware(nodeID flow.Identifier,
	networkMetrics module.NetworkMetrics,
	factoryFunc p2p.LibP2PFactoryFunc,
	validators ...network.MessageValidator) network.Middleware {

	if builder.Middleware == nil {
		// disable connection pruning for the staked AN which supports the observer service
		peerManagerFactory := p2p.PeerManagerFactory([]p2p.Option{p2p.WithInterval(builder.PeerUpdateInterval)}, p2p.WithConnectionPruning(false))

		builder.Middleware = p2p.NewMiddleware(
			builder.Logger.With().Bool("staked", false).Logger(),
			factoryFunc,
			nodeID,
			networkMetrics,
			builder.SporkID,
			p2p.DefaultUnicastTimeout,
			builder.IDTranslator,
			p2p.WithMessageValidators(validators...),
			p2p.WithPeerManager(peerManagerFactory),
			// use default identifier provider
		)
	}

	return builder.Middleware
}

// initUnstakedLocal initializes the unstaked node ID, network key and network address
// Currently, it reads a node-info.priv.json like any other node.
// TODO: read the node ID from the special bootstrap files
func (builder *ObserverNodeBuilder) initUnstakedLocal() func(node *cmd.NodeConfig) error {
	return func(node *cmd.NodeConfig) error {
		// for an unstaked node, set the identity here explicitly since it will not be found in the protocol state
		self := &flow.Identity{
			NodeID:        node.NodeID,
			NetworkPubKey: node.NetworkKey.PublicKey(),
			StakingPubKey: nil,             // no staking key needed for the unstaked node
			Role:          flow.RoleAccess, // unstaked node can only run as an access node
			Address:       builder.BindAddr,
		}

		var err error
		node.Me, err = local.NewNoKey(self)
		if err != nil {
			return fmt.Errorf("could not initialize local: %w", err)
		}
		return nil
	}
}

// enqueueMiddleware enqueues the creation of the network middleware
// this needs to be done before sync engine participants module
func (builder *ObserverNodeBuilder) enqueueMiddleware() {
	builder.
		Module("network middleware", func(node *cmd.NodeConfig) error {

			// NodeID for the unstaked node on the unstaked network
			unstakedNodeID := node.NodeID

			// Networking key
			unstakedNetworkKey := node.NetworkKey

			libP2PFactory := builder.initLibP2PFactory(unstakedNetworkKey)

			msgValidators := unstakedNetworkMsgValidators(node.Logger, node.IdentityProvider, unstakedNodeID)

			builder.initMiddleware(unstakedNodeID, node.Metrics.Network, libP2PFactory, msgValidators...)

			return nil
		})
}

// enqueueUnstakedNetworkInit enqueues the unstaked network component initialized for the unstaked node
func (builder *ObserverNodeBuilder) enqueueUnstakedNetworkInit() {

	builder.Component("unstaked network", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		var heroCacheCollector module.HeroCacheMetrics = metrics.NewNoopCollector()
		if builder.HeroCacheMetricsEnable {
			heroCacheCollector = metrics.NetworkReceiveCacheMetricsFactory(builder.MetricsRegisterer)
		}
		receiveCache := netcache.NewHeroReceiveCache(builder.NetworkReceivedMessageCacheSize,
			builder.Logger,
			heroCacheCollector)

		//err := node.Metrics.Mempool.Register(metrics.ResourceNetworkingReceiveCache, receiveCache.Size)
		//if err != nil {
		//	return nil, fmt.Errorf("could not register networking receive cache metric: %w", err)
		//}

		// topology is nil since it is automatically managed by libp2p
		net, err := builder.initNetwork(builder.Me, builder.Metrics.Network, builder.Middleware, nil, receiveCache)
		if err != nil {
			return nil, err
		}

		builder.Network = converter.NewNetwork(net, engine.SyncCommittee, engine.PublicSyncCommittee)

		builder.Logger.Info().Msgf("unstaked network will run on address: %s", builder.BindAddr)

		idEvents := gadgets.NewIdentityDeltas(builder.Middleware.UpdateNodeAddresses)
		builder.ProtocolEvents.AddConsumer(idEvents)

		return builder.Network, nil
	})
}

// enqueueConnectWithStakedAN enqueues the upstream connector component which connects the libp2p host of the unstaked
// AN with the staked AN.
// Currently, there is an issue with LibP2P stopping advertisements of subscribed topics if no peers are connected
// (https://github.com/libp2p/go-libp2p-pubsub/issues/442). This means that an unstaked AN could end up not being
// discovered by other unstaked ANs if it subscribes to a topic before connecting to the staked AN. Hence, the need
// of an explicit connect to the staked AN before the node attempts to subscribe to topics.
func (builder *ObserverNodeBuilder) enqueueConnectWithStakedAN() error {
	if len(builder.bootstrapIdentities) == 0 {
		return fmt.Errorf("empty boostrap identities")
	}

	builder.Component("upstream connector", func(_ *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		return follower.NewUpstreamConnector(builder.bootstrapIdentities, builder.LibP2PNode, builder.Logger), nil
	})

	return nil
}

func loadNetworkingKey(path string) (crypto.PrivateKey, error) {
	data, err := io.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read networking key (path=%s): %w", path, err)
	}

	keyBytes, err := hex.DecodeString(strings.Trim(string(data), "\n "))
	if err != nil {
		return nil, fmt.Errorf("could not hex decode networking key (path=%s): %w", path, err)
	}

	networkingKey, err := crypto.DecodePrivateKey(crypto.ECDSASecp256k1, keyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not decode networking key (path=%s): %w", path, err)
	}

	return networkingKey, nil
}

func unstakedNetworkMsgValidators(log zerolog.Logger, idProvider id.IdentityProvider, selfID flow.Identifier) []network.MessageValidator {
	return []network.MessageValidator{
		// filter out messages sent by this node itself
		validator.ValidateNotSender(selfID),
		validator.NewAnyValidator(
			// message should be either from a valid staked node
			validator.NewOriginValidator(
				id.NewIdentityFilterIdentifierProvider(filter.IsValidCurrentEpochParticipant, idProvider),
			),
			// or the message should be specifically targeted for this node
			validator.ValidateTarget(log, selfID),
		),
	}
}
