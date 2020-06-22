package tester

import (
	"path/filepath"
	"testing"

	"github.com/meshplus/bitxhub-kit/key"
	"github.com/meshplus/premo/pkg/appchain/ethereum"
	"github.com/meshplus/premo/pkg/appchain/fabric"
	"go.etcd.io/etcd/pkg/fileutil"

	"github.com/meshplus/premo/internal/repo"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
)

func TestTester(t *testing.T) {
	repoRoot, err := repo.PathRoot()
	require.Nil(t, err)

	transferContractAddr := "0x668a209Dc6562707469374B8235e37b8eC25db08"
	ethAccountKeyPath := filepath.Join(repoRoot, ".pier_ethereum", "eth", "account.key")
	require.True(t, fileutil.Exist(ethAccountKeyPath))
	ethClient, err := ethereum.New("http://localhost:8545", ethAccountKeyPath)
	require.Nil(t, err)
	ethLoadKey, err := key.LoadKey(filepath.Join(repoRoot, ".pier_ethereum", "key.json"))
	require.Nil(t, err)
	ethClientHelper := &EthClientHelper{
		EthClient:    ethClient,
		abiPath:      "test_data/ethereum/transfer.abi",
		contractAddr: transferContractAddr,
		appchainId:   ethLoadKey.Address.Hex(),
	}

	fabricLoadKey, err := key.LoadKey(filepath.Join(repoRoot, ".pier_fabric", "key.json"))
	require.Nil(t, err)
	fabricClient, err := fabric.New(filepath.Join(repoRoot, ".pier_fabric", "fabric"))
	require.Nil(t, err)
	fabricClientHelper := &FabricClientHelper{
		FabricClient: fabricClient,
		appchainId:   fabricLoadKey.Address.Hex(),
	}

	suite.Run(t, &Interchain{
		repoRoot:     repoRoot,
		ethClient:    ethClientHelper,
		fabricClient: fabricClientHelper,
	})
}
