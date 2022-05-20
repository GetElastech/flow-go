package node_builder

import (
	"context"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/onflow/flow-go/engine"
	pingeng "github.com/onflow/flow-go/engine/access/ping"
	"github.com/onflow/flow-go/module/metrics/unstaked"
	"github.com/onflow/flow-go/network/p2p/unicast"
	relaynet "github.com/onflow/flow-go/network/relay"
	"github.com/onflow/flow-go/network/topology"
	"github.com/onflow/flow-go/utils/grpcutils"
	"google.golang.org/grpc"
	"strings"
	"time"

	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/onflow/flow-go/module/compliance"

	"github.com/onflow/flow-go/cmd"
	"github.com/onflow/flow-go/consensus"
	"github.com/onflow/flow-go/consensus/hotstuff"
	"github.com/onflow/flow-go/consensus/hotstuff/committees"
	"github.com/onflow/flow-go/consensus/hotstuff/notifications/pubsub"
	hotsignature "github.com/onflow/flow-go/consensus/hotstuff/signature"
	"github.com/onflow/flow-go/consensus/hotstuff/verification"
	recovery "github.com/onflow/flow-go/consensus/recovery/protocol"
	"github.com/onflow/flow-go/engine/access/ingestion"
	"github.com/onflow/flow-go/engine/access/rpc"
	"github.com/onflow/flow-go/engine/access/rpc/backend"
	"github.com/onflow/flow-go/engine/common/follower"
	followereng "github.com/onflow/flow-go/engine/common/follower"
	"github.com/onflow/flow-go/engine/common/requester"
	synceng "github.com/onflow/flow-go/engine/common/synchronization"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/buffer"
	finalizer "github.com/onflow/flow-go/module/finalizer/consensus"
	"github.com/onflow/flow-go/module/id"
	"github.com/onflow/flow-go/module/mempool/stdmap"
	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/module/synchronization"
	"github.com/onflow/flow-go/network"
	netcache "github.com/onflow/flow-go/network/cache"
	cborcodec "github.com/onflow/flow-go/network/codec/cbor"
	"github.com/onflow/flow-go/network/p2p"
	"github.com/onflow/flow-go/network/validator"
	"github.com/onflow/flow-go/state/protocol"
	badgerState "github.com/onflow/flow-go/state/protocol/badger"
	"github.com/onflow/flow-go/state/protocol/blocktimer"
	storage "github.com/onflow/flow-go/storage/badger"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/grpc/credentials"

	"github.com/onflow/flow-go/crypto"
)

// AccessNodeBuilder extends cmd.NodeBuilder and declares additional functions needed to bootstrap an Access node.
// The private network allows the staked nodes to communicate among themselves, while the public network allows the
// Observers and an Access node to communicate.
//
//                                 public network                           private network
//  +------------------------+
//  | Observer             1 |<--------------------------|
//  +------------------------+                           v
//  +------------------------+                         +--------------------+                 +------------------------+
//  | Observer             2 |<----------------------->| Staked Access Node |<--------------->| All other staked Nodes |
//  +------------------------+                         +--------------------+                 +------------------------+
//  +------------------------+                           ^
//  | Observer             3 |<--------------------------|
//  +------------------------+

type AccessNodeBuilder interface {
	cmd.NodeBuilder
}

// AccessNodeConfig defines all the user defined parameters required to bootstrap an access node
// For a node running as a standalone process, the config fields will be populated from the command line params,
// while for a node running as a library, the config fields are expected to be initialized by the caller.
type AccessNodeConfig struct {
	supportsObserverAndFollower  bool // True if this is an Access node that supports observers and consensus follower engines
	collectionGRPCPort           uint
	executionGRPCPort            uint
	pingEnabled                  bool
	nodeInfoFile                 string
	apiRatelimits                map[string]int
	apiBurstlimits               map[string]int
	rpcConf                      rpc.Config
	HistoricalAccessRPCs         []access.AccessAPIClient
	logTxTimeToFinalized         bool
	logTxTimeToExecuted          bool
	logTxTimeToFinalizedExecuted bool
	retryEnabled                 bool
	rpcMetricsEnabled            bool
	baseOptions                  []cmd.Option

	PublicNetworkConfig PublicNetworkConfig
}

type PublicNetworkConfig struct {
	// NetworkKey crypto.PublicKey // TODO: do we need a different key for the public network?
	BindAddress string
	Network     network.Network
	Metrics     module.NetworkMetrics
}

// DefaultAccessNodeConfig defines all the default values for the AccessNodeConfig
func DefaultAccessNodeConfig() *AccessNodeConfig {
	return &AccessNodeConfig{
		collectionGRPCPort: 9000,
		executionGRPCPort:  9000,
		rpcConf: rpc.Config{
			UnsecureGRPCListenAddr:    "0.0.0.0:9000",
			SecureGRPCListenAddr:      "0.0.0.0:9001",
			HTTPListenAddr:            "0.0.0.0:8000",
			RESTListenAddr:            "",
			CollectionAddr:            "",
			HistoricalAccessAddrs:     "",
			CollectionClientTimeout:   3 * time.Second,
			ExecutionClientTimeout:    3 * time.Second,
			MaxHeightRange:            backend.DefaultMaxHeightRange,
			PreferredExecutionNodeIDs: nil,
			FixedExecutionNodeIDs:     nil,
		},
		logTxTimeToFinalized:         false,
		logTxTimeToExecuted:          false,
		logTxTimeToFinalizedExecuted: false,
		pingEnabled:                  false,
		retryEnabled:                 false,
		rpcMetricsEnabled:            false,
		nodeInfoFile:                 "",
		apiRatelimits:                nil,
		apiBurstlimits:               nil,
		supportsObserverAndFollower:  false,
		PublicNetworkConfig: PublicNetworkConfig{
			BindAddress: cmd.NotSet,
			Metrics:     metrics.NewNoopCollector(),
		},
	}
}

// FlowAccessNodeBuilder provides the common functionality needed to bootstrap a Flow access node
// It is composed of the FlowNodeBuilder, the AccessNodeConfig and contains all the components and modules needed for the
// access nodes
type FlowAccessNodeBuilder struct {
	*cmd.FlowNodeBuilder
	*AccessNodeConfig

	// components
	LibP2PNode                 *p2p.Node
	FollowerState              protocol.MutableState
	SyncCore                   *synchronization.Core
	RpcEng                     *rpc.Engine
	FinalizationDistributor    *pubsub.FinalizationDistributor
	FinalizedHeader            *synceng.FinalizedHeaderCache
	CollectionRPC              access.AccessAPIClient
	TransactionTimings         *stdmap.TransactionTimings
	CollectionsToMarkFinalized *stdmap.Times
	CollectionsToMarkExecuted  *stdmap.Times
	BlocksToMarkExecuted       *stdmap.Times
	TransactionMetrics         module.TransactionMetrics
	PingMetrics                module.PingMetrics
	Committee                  hotstuff.Committee
	Finalized                  *flow.Header
	Pending                    []*flow.Header
	FollowerCore               module.HotStuffFollower
	// The sync engine participants provider is the libp2p peer store for the observer
	// which is not available until after the network has started.
	// Hence, a factory function that needs to be called just before creating the sync engine
	SyncEngineParticipantsProviderFactory func() id.IdentifierProvider

	// engines
	IngestEng   *ingestion.Engine
	RequestEng  *requester.Engine
	FollowerEng *followereng.Engine
	SyncEng     *synceng.Engine
}

func (builder *FlowAccessNodeBuilder) buildFollowerState() *FlowAccessNodeBuilder {
	builder.Module("mutable follower state", func(node *cmd.NodeConfig) error {
		// For now, we only support state implementations from package badger.
		// If we ever support different implementations, the following can be replaced by a type-aware factory
		state, ok := node.State.(*badgerState.State)
		if !ok {
			return fmt.Errorf("only implementations of type badger.State are currently supported but read-only state has type %T", node.State)
		}

		followerState, err := badgerState.NewFollowerState(
			state,
			node.Storage.Index,
			node.Storage.Payloads,
			node.Tracer,
			node.ProtocolEvents,
			blocktimer.DefaultBlockTimer,
		)
		builder.FollowerState = followerState

		return err
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) buildSyncCore() *FlowAccessNodeBuilder {
	builder.Module("sync core", func(node *cmd.NodeConfig) error {
		syncCore, err := synchronization.New(node.Logger, node.SyncCoreConfig)
		builder.SyncCore = syncCore

		return err
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) buildCommittee() *FlowAccessNodeBuilder {
	builder.Module("committee", func(node *cmd.NodeConfig) error {
		// initialize consensus committee's membership state
		// This committee state is for the HotStuff follower, which follows the MAIN CONSENSUS Committee
		// Note: node.Me.NodeID() is not part of the consensus committee
		committee, err := committees.NewConsensusCommittee(node.State, node.Me.NodeID())
		builder.Committee = committee

		return err
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) buildLatestHeader() *FlowAccessNodeBuilder {
	builder.Module("latest header", func(node *cmd.NodeConfig) error {
		finalized, pending, err := recovery.FindLatest(node.State, node.Storage.Headers)
		builder.Finalized, builder.Pending = finalized, pending

		return err
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) buildFollowerCore() *FlowAccessNodeBuilder {
	builder.Component("follower core", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		// create a finalizer that will handle updating the protocol
		// state when the follower detects newly finalized blocks
		final := finalizer.NewFinalizer(node.DB, node.Storage.Headers, builder.FollowerState, node.Tracer)

		packer := hotsignature.NewConsensusSigDataPacker(builder.Committee)
		// initialize the verifier for the protocol consensus
		verifier := verification.NewCombinedVerifier(builder.Committee, packer)

		followerCore, err := consensus.NewFollower(node.Logger, builder.Committee, node.Storage.Headers, final, verifier,
			builder.FinalizationDistributor, node.RootBlock.Header, node.RootQC, builder.Finalized, builder.Pending)
		if err != nil {
			return nil, fmt.Errorf("could not initialize follower core: %w", err)
		}
		builder.FollowerCore = followerCore

		return builder.FollowerCore, nil
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) buildFollowerEngine() *FlowAccessNodeBuilder {
	builder.Component("follower engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		// initialize cleaner for DB
		cleaner := storage.NewCleaner(node.Logger, node.DB, builder.Metrics.CleanCollector, flow.DefaultValueLogGCFrequency)
		conCache := buffer.NewPendingBlocks()

		followerEng, err := follower.New(
			node.Logger,
			node.Network,
			node.Me,
			node.Metrics.Engine,
			node.Metrics.Mempool,
			cleaner,
			node.Storage.Headers,
			node.Storage.Payloads,
			builder.FollowerState,
			conCache,
			builder.FollowerCore,
			builder.SyncCore,
			node.Tracer,
			compliance.WithSkipNewProposalsThreshold(builder.ComplianceConfig.SkipNewProposalsThreshold),
		)
		if err != nil {
			return nil, fmt.Errorf("could not create follower engine: %w", err)
		}
		builder.FollowerEng = followerEng

		return builder.FollowerEng, nil
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) buildFinalizedHeader() *FlowAccessNodeBuilder {
	builder.Component("finalized snapshot", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		finalizedHeader, err := synceng.NewFinalizedHeaderCache(node.Logger, node.State, builder.FinalizationDistributor)
		if err != nil {
			return nil, fmt.Errorf("could not create finalized snapshot cache: %w", err)
		}
		builder.FinalizedHeader = finalizedHeader

		return builder.FinalizedHeader, nil
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) buildSyncEngine() *FlowAccessNodeBuilder {
	builder.Component("sync engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		sync, err := synceng.New(
			node.Logger,
			node.Metrics.Engine,
			node.Network,
			node.Me,
			node.Storage.Blocks,
			builder.FollowerEng,
			builder.SyncCore,
			builder.FinalizedHeader,
			builder.SyncEngineParticipantsProviderFactory(),
		)
		if err != nil {
			return nil, fmt.Errorf("could not create synchronization engine: %w", err)
		}
		builder.SyncEng = sync

		return builder.SyncEng, nil
	})

	return builder
}

func (builder *FlowAccessNodeBuilder) BuildConsensusFollower() AccessNodeBuilder {
	builder.
		buildFollowerState().
		buildSyncCore().
		buildCommittee().
		buildLatestHeader().
		buildFollowerCore().
		buildFollowerEngine().
		buildFinalizedHeader().
		buildSyncEngine()

	return builder
}

type Option func(*AccessNodeConfig)

func FlowAccessNode(opts ...Option) *FlowAccessNodeBuilder {
	config := DefaultAccessNodeConfig()
	for _, opt := range opts {
		opt(config)
	}

	return &FlowAccessNodeBuilder{
		AccessNodeConfig:        config,
		FlowNodeBuilder:         cmd.FlowNode(flow.RoleAccess.String(), config.baseOptions...),
		FinalizationDistributor: pubsub.NewFinalizationDistributor(),
	}
}

func (builder *FlowAccessNodeBuilder) ParseFlags() error {

	builder.BaseFlags()

	builder.extraFlags()

	return builder.ParseAndPrintFlags()
}

func (builder *FlowAccessNodeBuilder) extraFlags() {
	builder.ExtraFlags(func(flags *pflag.FlagSet) {
		defaultConfig := DefaultAccessNodeConfig()
		var dummyBool bool
		var dummyString string
		var dummyStringSlice []string

		flags.UintVar(&builder.collectionGRPCPort, "collection-ingress-port", defaultConfig.collectionGRPCPort, "the grpc ingress port for all collection nodes")
		flags.UintVar(&builder.executionGRPCPort, "execution-ingress-port", defaultConfig.executionGRPCPort, "the grpc ingress port for all execution nodes")
		flags.StringVarP(&builder.rpcConf.UnsecureGRPCListenAddr, "rpc-addr", "r", defaultConfig.rpcConf.UnsecureGRPCListenAddr, "the address the unsecured gRPC server listens on")
		flags.StringVar(&builder.rpcConf.SecureGRPCListenAddr, "secure-rpc-addr", defaultConfig.rpcConf.SecureGRPCListenAddr, "the address the secure gRPC server listens on")
		flags.StringVarP(&builder.rpcConf.HTTPListenAddr, "http-addr", "h", defaultConfig.rpcConf.HTTPListenAddr, "the address the http proxy server listens on")
		flags.StringVar(&builder.rpcConf.RESTListenAddr, "rest-addr", defaultConfig.rpcConf.RESTListenAddr, "the address the REST server listens on (if empty the REST server will not be started)")
		flags.StringVarP(&builder.rpcConf.CollectionAddr, "static-collection-ingress-addr", "", defaultConfig.rpcConf.CollectionAddr, "the address (of the collection node) to send transactions to")
		flags.StringVarP(&dummyString, "script-addr", "s", "", "deprecated - the address (of the execution node) forward the script to")
		flags.StringVarP(&builder.rpcConf.HistoricalAccessAddrs, "historical-access-addr", "", defaultConfig.rpcConf.HistoricalAccessAddrs, "comma separated rpc addresses for historical access nodes")
		flags.DurationVar(&builder.rpcConf.CollectionClientTimeout, "collection-client-timeout", defaultConfig.rpcConf.CollectionClientTimeout, "grpc client timeout for a collection node")
		flags.DurationVar(&builder.rpcConf.ExecutionClientTimeout, "execution-client-timeout", defaultConfig.rpcConf.ExecutionClientTimeout, "grpc client timeout for an execution node")
		flags.UintVar(&builder.rpcConf.MaxHeightRange, "rpc-max-height-range", defaultConfig.rpcConf.MaxHeightRange, "maximum size for height range requests")
		flags.StringSliceVar(&builder.rpcConf.PreferredExecutionNodeIDs, "preferred-execution-node-ids", defaultConfig.rpcConf.PreferredExecutionNodeIDs, "comma separated list of execution nodes ids to choose from when making an upstream call e.g. b4a4dbdcd443d...,fb386a6a... etc.")
		flags.StringSliceVar(&builder.rpcConf.FixedExecutionNodeIDs, "fixed-execution-node-ids", defaultConfig.rpcConf.FixedExecutionNodeIDs, "comma separated list of execution nodes ids to choose from when making an upstream call if no matching preferred execution id is found e.g. b4a4dbdcd443d...,fb386a6a... etc.")
		flags.BoolVar(&builder.logTxTimeToFinalized, "log-tx-time-to-finalized", defaultConfig.logTxTimeToFinalized, "log transaction time to finalized")
		flags.BoolVar(&builder.logTxTimeToExecuted, "log-tx-time-to-executed", defaultConfig.logTxTimeToExecuted, "log transaction time to executed")
		flags.BoolVar(&builder.logTxTimeToFinalizedExecuted, "log-tx-time-to-finalized-executed", defaultConfig.logTxTimeToFinalizedExecuted, "log transaction time to finalized and executed")
		flags.BoolVar(&builder.pingEnabled, "ping-enabled", defaultConfig.pingEnabled, "whether to enable the ping process that pings all other peers and report the connectivity to metrics")
		flags.BoolVar(&builder.retryEnabled, "retry-enabled", defaultConfig.retryEnabled, "whether to enable the retry mechanism at the access node level")
		flags.BoolVar(&builder.rpcMetricsEnabled, "rpc-metrics-enabled", defaultConfig.rpcMetricsEnabled, "whether to enable the rpc metrics")
		flags.StringVarP(&builder.nodeInfoFile, "node-info-file", "", defaultConfig.nodeInfoFile, "full path to a json file which provides more details about nodes when reporting its reachability metrics")
		flags.StringToIntVar(&builder.apiRatelimits, "api-rate-limits", defaultConfig.apiRatelimits, "per second rate limits for Access API methods e.g. Ping=300,GetTransaction=500 etc.")
		flags.StringToIntVar(&builder.apiBurstlimits, "api-burst-limits", defaultConfig.apiBurstlimits, "burst limits for Access API methods e.g. Ping=100,GetTransaction=100 etc.")
		flags.BoolVar(&dummyBool, "staked", true, "deprecated - whether this node is a staked access node or not")
		flags.StringVar(&dummyString, "observer-networking-key-path", "", "deprecated - path to the networking key for observer")
		flags.StringSliceVar(&dummyStringSlice, "bootstrap-node-addresses", []string{}, "deprecated - the network addresses of the bootstrap access node if this is an unstaked access node e.g. access-001.mainnet.flow.org:9653,access-002.mainnet.flow.org:9653")
		flags.StringSliceVar(&dummyStringSlice, "bootstrap-node-public-keys", []string{}, "deprecated - the networking public key of the bootstrap access node if this is an unstaked access node (in the same order as the bootstrap node addresses) e.g. \"d57a5e9c5.....\",\"44ded42d....\"")
		flags.BoolVar(&builder.supportsObserverAndFollower, "supports-unstaked-node", defaultConfig.supportsObserverAndFollower, "true if this staked access node supports observer or follower connections")
		flags.StringVar(&builder.PublicNetworkConfig.BindAddress, "public-network-address", defaultConfig.PublicNetworkConfig.BindAddress, "staked access node's public network bind address")
	}).ValidateFlags(func() error {
		if builder.supportsObserverAndFollower && (builder.PublicNetworkConfig.BindAddress == cmd.NotSet || builder.PublicNetworkConfig.BindAddress == "") {
			return errors.New("public-network-address must be set if supports-unstaked-node is true")
		}

		return nil
	})
}

// initNetwork creates the network.Network implementation with the given metrics, middleware, initial list of network
// participants and topology used to choose peers from the list of participants. The list of participants can later be
// updated by calling network.SetIDs.
func (builder *FlowAccessNodeBuilder) initNetwork(nodeID module.Local,
	networkMetrics module.NetworkMetrics,
	middleware network.Middleware,
	topology network.Topology,
	receiveCache *netcache.ReceiveCache,
) (*p2p.Network, error) {

	codec := cborcodec.NewCodec()

	// creates network instance
	net, err := p2p.NewNetwork(
		builder.Logger,
		codec,
		nodeID,
		func() (network.Middleware, error) { return builder.Middleware, nil },
		topology,
		p2p.NewChannelSubscriptionManager(middleware),
		networkMetrics,
		builder.IdentityProvider,
		receiveCache,
	)
	if err != nil {
		return nil, fmt.Errorf("could not initialize network: %w", err)
	}

	return net, nil
}

func publicNetworkMsgValidators(log zerolog.Logger, idProvider id.IdentityProvider, selfID flow.Identifier) []network.MessageValidator {
	return []network.MessageValidator{
		// filter out messages sent by this node itself
		validator.ValidateNotSender(selfID),
		validator.NewAnyValidator(
			// message should be from a valid staked node
			validator.NewOriginValidator(
				id.NewIdentityFilterIdentifierProvider(filter.IsValidCurrentEpochParticipant, idProvider),
			),
			// or the message should be specifically targeted for this node
			validator.ValidateTarget(log, selfID),
		),
	}
}

// accessNodeBuilder builds a staked access node.
// The access node can optionally participate in the public network publishing data for the observers downstream.
func (builder *FlowAccessNodeBuilder) InitIDProviders() {
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

func (builder *FlowAccessNodeBuilder) Initialize() error {
	builder.InitIDProviders()

	builder.EnqueueResolver()

	// enqueue the regular network
	builder.EnqueueNetworkInit()

	// if this is an access node that supports followers, enqueue the public network
	if builder.supportsObserverAndFollower {
		builder.enqueuePublicNetworkInit()
		builder.enqueueRelayNetwork()
	}

	builder.EnqueuePingService()

	builder.EnqueueMetricsServerInit()

	if err := builder.RegisterBadgerMetrics(); err != nil {
		return err
	}

	builder.EnqueueTracer()
	builder.PreInit(cmd.DynamicStartPreInit)
	return nil
}

func (builder *FlowAccessNodeBuilder) enqueueRelayNetwork() {
	builder.Component("relay network", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		relayNet := relaynet.NewRelayNetwork(
			node.Network,
			builder.AccessNodeConfig.PublicNetworkConfig.Network,
			node.Logger,
			[]network.Channel{engine.ReceiveBlocks},
		)
		node.Network = relayNet
		return relayNet, nil
	})
}

func (builder *FlowAccessNodeBuilder) Build() (cmd.Node, error) {
	builder.
		BuildConsensusFollower().
		Module("collection node client", func(node *cmd.NodeConfig) error {
			// collection node address is optional (if not specified, collection nodes will be chosen at random)
			if strings.TrimSpace(builder.rpcConf.CollectionAddr) == "" {
				node.Logger.Info().Msg("using a dynamic collection node address")
				return nil
			}

			node.Logger.Info().
				Str("collection_node", builder.rpcConf.CollectionAddr).
				Msg("using the static collection node address")

			collectionRPCConn, err := grpc.Dial(
				builder.rpcConf.CollectionAddr,
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcutils.DefaultMaxMsgSize)),
				grpc.WithInsecure(), //nolint:staticcheck
				backend.WithClientUnaryInterceptor(builder.rpcConf.CollectionClientTimeout))
			if err != nil {
				return err
			}
			builder.CollectionRPC = access.NewAccessAPIClient(collectionRPCConn)
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
			// generate the server certificate that will be served by the GRPC server
			x509Certificate, err := grpcutils.X509Certificate(node.NetworkKey)
			if err != nil {
				return err
			}
			tlsConfig := grpcutils.DefaultServerTLSConfig(x509Certificate)
			builder.rpcConf.TransportCredentials = credentials.NewTLS(tlsConfig)
			return nil
		}).
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
				nil,
			)
			return builder.RpcEng, nil
		}).
		Component("ingestion engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
			var err error

			builder.RequestEng, err = requester.New(
				node.Logger,
				node.Metrics.Engine,
				node.Network,
				node.Me,
				node.State,
				engine.RequestCollections,
				filter.HasRole(flow.RoleCollection),
				func() flow.Entity { return &flow.Collection{} },
			)
			if err != nil {
				return nil, fmt.Errorf("could not create requester engine: %w", err)
			}

			builder.IngestEng, err = ingestion.New(node.Logger, node.Network, node.State, node.Me, builder.RequestEng, node.Storage.Blocks, node.Storage.Headers, node.Storage.Collections, node.Storage.Transactions, node.Storage.Results, node.Storage.Receipts, builder.TransactionMetrics,
				builder.CollectionsToMarkFinalized, builder.CollectionsToMarkExecuted, builder.BlocksToMarkExecuted, builder.RpcEng)
			if err != nil {
				return nil, err
			}
			builder.RequestEng.WithHandle(builder.IngestEng.OnCollection)
			builder.FinalizationDistributor.AddConsumer(builder.IngestEng)

			return builder.IngestEng, nil
		}).
		Component("requester engine", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
			// We initialize the requester engine inside the ingestion engine due to the mutual dependency. However, in
			// order for it to properly start and shut down, we should still return it as its own engine here, so it can
			// be handled by the scaffold.
			return builder.RequestEng, nil
		})

	if builder.supportsObserverAndFollower {
		builder.Component("unstaked sync request handler", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
			syncRequestHandler, err := synceng.NewRequestHandlerEngine(
				node.Logger.With().Bool("unstaked", true).Logger(),
				unstaked.NewUnstakedEngineCollector(node.Metrics.Engine),
				builder.AccessNodeConfig.PublicNetworkConfig.Network,
				node.Me,
				node.Storage.Blocks,
				builder.SyncCore,
				builder.FinalizedHeader,
			)

			if err != nil {
				return nil, fmt.Errorf("could not create unstaked sync request handler: %w", err)
			}

			return syncRequestHandler, nil
		})
	}

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

	return builder.FlowNodeBuilder.Build()
}

// enqueuePublicNetworkInit enqueues the public network component initialized for the access node
// Observers can connect to the public port to access features without affecting the system load much.
func (builder *FlowAccessNodeBuilder) enqueuePublicNetworkInit() {
	builder.Component("unstaked network", func(node *cmd.NodeConfig) (module.ReadyDoneAware, error) {
		builder.PublicNetworkConfig.Metrics = metrics.NewNetworkCollector(metrics.WithNetworkPrefix("unstaked"))

		libP2PFactory := builder.initLibP2PFactory(builder.NodeConfig.NetworkKey)

		msgValidators := publicNetworkMsgValidators(node.Logger.With().Bool("unstaked", true).Logger(), node.IdentityProvider, builder.NodeID)

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

		builder.AccessNodeConfig.PublicNetworkConfig.Network = net

		node.Logger.Info().Msgf("network will run on address: %s", builder.PublicNetworkConfig.BindAddress)
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
func (builder *FlowAccessNodeBuilder) initLibP2PFactory(networkKey crypto.PrivateKey) p2p.LibP2PFactoryFunc {
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
			SetPubSub(p2ppubsub.NewGossipSub).
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
func (builder *FlowAccessNodeBuilder) initMiddleware(nodeID flow.Identifier,
	networkMetrics module.NetworkMetrics,
	factoryFunc p2p.LibP2PFactoryFunc,
	validators ...network.MessageValidator) network.Middleware {

	// disable connection pruning for the AN which supports the observer
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

	return builder.Middleware
}
