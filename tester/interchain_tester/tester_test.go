package interchain_tester

import (
	"path/filepath"
	"testing"

	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/appchain/ethereum"
	"github.com/meshplus/premo/pkg/appchain/fabric"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/etcd/pkg/fileutil"
)

func TestTester(t *testing.T) {
	repoRoot, err := repo.PathRoot()
	require.Nil(t, err)

	transferContractAddr := "0x668a209Dc6562707469374B8235e37b8eC25db08"
	etherConfigPath := filepath.Join(repoRoot, ".pier_ethereum", "ethereum")
	require.True(t, fileutil.Exist(etherConfigPath))

	ethClient, err := ethereum.New(etherConfigPath)
	require.Nil(t, err)

	transferAbi := filepath.Join(repoRoot, "transfer.abi")
	require.True(t, fileutil.Exist(transferAbi))

	ethLoadKey, err := asym.RestorePrivateKey(filepath.Join(repoRoot, ".pier_ethereum", "key.json"), repo.KeyPassword)
	require.Nil(t, err)
	address, err := ethLoadKey.PublicKey().Address()
	require.Nil(t, err)
	ethClientHelper := &EthClientHelper{
		EthClient:    ethClient,
		abiPath:      transferAbi,
		contractAddr: transferContractAddr,
		appchainId:   address.String(),
	}

	fabricLoadKey, err := asym.RestorePrivateKey(filepath.Join(repoRoot, ".pier_fabric", "key.json"), repo.KeyPassword)
	require.Nil(t, err)
	fabricClient, err := fabric.New(filepath.Join(repoRoot, ".pier_fabric", "fabric"))
	require.Nil(t, err)
	address, err = fabricLoadKey.PublicKey().Address()
	require.Nil(t, err)
	fabricClientHelper := &FabricClientHelper{
		FabricClient: fabricClient,
		appchainId:   address.String(),
	}

	suite.Run(t, &Interchain{
		repoRoot:     repoRoot,
		ethClient:    ethClientHelper,
		fabricClient: fabricClientHelper,
	})
}
