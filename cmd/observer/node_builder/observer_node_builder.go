package node_builder

//
//import (
//	"github.com/onflow/flow-go/cmd"
//	shared "github.com/onflow/flow-go/cmd/access/node_builder"
//	"github.com/onflow/flow-go/consensus/hotstuff"
//	"github.com/onflow/flow-go/consensus/hotstuff/notifications/pubsub"
//	"github.com/onflow/flow-go/engine/access/ingestion"
//	"github.com/onflow/flow-go/engine/access/rpc"
//	followereng "github.com/onflow/flow-go/engine/common/follower"
//	"github.com/onflow/flow-go/engine/common/requester"
//	synceng "github.com/onflow/flow-go/engine/common/synchronization"
//	"github.com/onflow/flow-go/model/flow"
//	"github.com/onflow/flow-go/module"
//	"github.com/onflow/flow-go/module/id"
//	"github.com/onflow/flow-go/module/mempool/stdmap"
//	"github.com/onflow/flow-go/module/synchronization"
//	"github.com/onflow/flow-go/network"
//	"github.com/onflow/flow-go/network/p2p"
//	"github.com/onflow/flow-go/state/protocol"
//	"github.com/onflow/flow/protobuf/go/flow/access"
//)
//
//func NewObserverNodeBuilder(builder *shared.FlowAccessNodeBuilder) *shared.UnstakedAccessNodeBuilder {
//	return shared.NewUnstakedAccessNodeBuilder(builder)
//}
//
//// ObserverNodeBuilder extends cmd.NodeBuilder and declares additional functions needed to bootstrap an Observer node.
//// The Staked network allows the staked nodes to communicate among themselves, while the observer network allows the
//// observer services and a staked Access node to communicate. Observer nodes can be attached to an observer or staked
//// access node to fetch read only block, and execution information. Observer nodes are scalable.
////
////                                 observer network                           staked network
////  +------------------------+
////  |        Observer Node 1 |
////  +------------------------+
////              |
////              v
////  +------------------------+
////  |        Observer Node 2 |<--------------------------|
////  +------------------------+                           |
////              |                                        |
////              v                                        v
////  +------------------------+                         +--------------------+                 +------------------------+
////  |        Observer Node 1 |<----------------------->| Staked Access Node |<--------------->| All other staked Nodes |
////  +------------------------+                         +--------------------+                 +------------------------+
////  +------------------------+                           ^
////  |        Observer Node 2 |<--------------------------|
////  +------------------------+
//
//type ObserverNodeBuilder interface {
//	cmd.NodeBuilder
//}
//
//type PublicNetworkConfig struct {
//	// NetworkKey crypto.PublicKey // TODO: do we need a different key for the public network?
//	BindAddress string
//	Network     network.Network
//	Metrics     module.NetworkMetrics
//}
//
//// ObserverNodeConfig defines all the user defined parameters required to bootstrap an observer node
//// For a node running as a standalone process, the config fields will be populated from the command line params,
//// while for a node running as a library, the config fields are expected to be initialized by the caller.
//type ObserverNodeConfig struct {
//	shared.SharedNodeConfig
//}
//
//func DefaultObserverNodeConfig() *ObserverNodeConfig {
//	return &ObserverNodeConfig{
//		SharedNodeConfig: *shared.DefaultSharedNodeConfig(),
//	}
//}
//
//// CustomAccessNodeConfig defines custom values of an unstaged access node based on observer parameters
//func NewObserverNodeConfig(config *ObserverNodeConfig, opts ...shared.Option) *shared.AccessNodeConfig {
//	accessConfig := shared.UnstakedAccessNodeConfig()
//	for _, opt := range opts {
//		opt(accessConfig)
//	}
//	return accessConfig
//}
//
//// FlowAccessNodeBuilder provides the common functionality needed to bootstrap a Flow staked and observer node
//// It is composed of the FlowNodeBuilder, the AccessNodeConfig and contains all the components and modules needed for the
//// staked access and observer nodes
//type FlowObserverNodeBuilder struct {
//	*cmd.FlowNodeBuilder
//	*ObserverNodeConfig
//
//	// components
//	LibP2PNode                 *p2p.Node
//	FollowerState              protocol.MutableState
//	SyncCore                   *synchronization.Core
//	RpcEng                     *rpc.Engine
//	FinalizationDistributor    *pubsub.FinalizationDistributor
//	FinalizedHeader            *synceng.FinalizedHeaderCache
//	CollectionRPC              access.AccessAPIClient
//	TransactionTimings         *stdmap.TransactionTimings
//	CollectionsToMarkFinalized *stdmap.Times
//	CollectionsToMarkExecuted  *stdmap.Times
//	BlocksToMarkExecuted       *stdmap.Times
//	TransactionMetrics         module.TransactionMetrics
//	PingMetrics                module.PingMetrics
//	Committee                  hotstuff.Committee
//	Finalized                  *flow.Header
//	Pending                    []*flow.Header
//	FollowerCore               module.HotStuffFollower
//	// for the observer node, the sync engine participants provider is the libp2p peer store which is not
//	// available until after the network has started. Hence, a factory function that needs to be called just before
//	// creating the sync engine
//	SyncEngineParticipantsProviderFactory func() id.IdentifierProvider
//
//	// engines
//	IngestEng   *ingestion.Engine
//	RequestEng  *requester.Engine
//	FollowerEng *followereng.Engine
//	SyncEng     *synceng.Engine
//}
//
//func FlowObserverNode() *shared.FlowAccessNodeBuilder {
//	config := DefaultObserverNodeConfig()
//	accessConfig := NewObserverNodeConfig(config)
//
//	return &shared.FlowAccessNodeBuilder{
//		AccessNodeConfig:        accessConfig,
//		FlowNodeBuilder:         cmd.FlowNode(flow.RoleObserverService.String(), config.BaseOptions...),
//		FinalizationDistributor: pubsub.NewFinalizationDistributor(),
//	}
//}
