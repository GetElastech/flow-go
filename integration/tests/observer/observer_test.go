package access

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/onflow/flow-go/integration/testnet"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/utils/unittest"
	"github.com/onflow/flow/protobuf/go/flow/access"
)

func TestObserver(t *testing.T) {
	suite.Run(t, new(ObserverSuite))
}

type ObserverSuite struct {
	suite.Suite

	log zerolog.Logger

	// root context for the current test
	ctx    context.Context
	cancel context.CancelFunc

	net *testnet.FlowNetwork
}

func (s *ObserverSuite) TearDownTest() {
	s.log.Info().Msg("================> Start TearDownTest")
	s.net.Remove()
	s.cancel()
	s.log.Info().Msg("================> Finish TearDownTest")
}

func (suite *ObserverSuite) SetupTest() {
	logger := unittest.LoggerWithLevel(zerolog.InfoLevel).With().
		Str("testfile", "observer_test.go").
		Str("testcase", suite.T().Name()).
		Logger()
	suite.log = logger
	suite.log.Info().Msg("================> SetupTest")
	defer func() {
		suite.log.Info().Msg("================> Finish SetupTest")
	}()

	nodeConfigs := []testnet.NodeConfig{
		testnet.NewNodeConfig(flow.RoleObserverService, testnet.WithLogLevel(zerolog.InfoLevel)),
	}

	// need one access node
	accessConfig := testnet.NewNodeConfig(flow.RoleAccess, testnet.WithLogLevel(zerolog.FatalLevel))
	nodeConfigs = append(nodeConfigs, accessConfig)

	// need one execution node
	exeConfig := testnet.NewNodeConfig(flow.RoleExecution, testnet.WithLogLevel(zerolog.FatalLevel))
	nodeConfigs = append(nodeConfigs, exeConfig)

	// need one dummy verification node (unused ghost)
	verConfig := testnet.NewNodeConfig(flow.RoleVerification, testnet.WithLogLevel(zerolog.FatalLevel), testnet.AsGhost())
	nodeConfigs = append(nodeConfigs, verConfig)

	// need one controllable collection node (unused ghost)
	collConfig := testnet.NewNodeConfig(flow.RoleCollection, testnet.WithLogLevel(zerolog.FatalLevel), testnet.AsGhost())
	nodeConfigs = append(nodeConfigs, collConfig)

	// need three consensus nodes (unused ghost)
	for n := 0; n < 3; n++ {
		conID := unittest.IdentifierFixture()
		nodeConfig := testnet.NewNodeConfig(flow.RoleConsensus,
			testnet.WithLogLevel(zerolog.FatalLevel),
			testnet.WithID(conID),
			testnet.AsGhost())
		nodeConfigs = append(nodeConfigs, nodeConfig)
	}

	conf := testnet.NewNetworkConfig("access_api_test", nodeConfigs)
	suite.net = testnet.PrepareFlowNetwork(suite.T(), conf)

	// start the network
	suite.T().Logf("starting flow network with docker containers")
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
	suite.net.Start(suite.ctx)
}

func (suite *ObserverSuite) TestObserver() {
	gRPCAddress := fmt.Sprintf(":%s", suite.net.AccessPorts[testnet.AccessNodeAPIPort])
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, gRPCAddress, grpc.WithInsecure())
	require.NoError(suite.T(), err, "http proxy port not open on the observer node")
	defer func() {
		if err := conn.Close(); err != nil {
			require.NoError(suite.T(), err, "failed to close gRPC connection")
		}
	}()
	grpcClient := access.NewAccessAPIClient(conn)
	_, err = grpcClient.Ping(context.Background(), &access.PingRequest{})
	require.NoError(suite.T(), err, "insecure grpc port not open on the observer node")

	accid := "0xf8d6e0586b0a20c7"

	accresp, err := grpcClient.GetAccountAtLatestBlock(context.Background(), &access.GetAccountAtLatestBlockRequest{
		Address: flow.HexToAddress(accid).Bytes(),
	})
	require.NoErrorf(suite.T(), err, "account %s does not exist", accid)
	assert.NotNil(suite.T(), accresp.Account, "account not found")
}
