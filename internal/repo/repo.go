package repo

import (
	"fmt"
	"path/filepath"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"

	"github.com/meshplus/bitxhub-kit/fileutil"
	"github.com/mitchellh/go-homedir"
)

const (
	// DefaultPathName is the default config dir name
	DefaultPathName = ".premo"

	// DefaultPathRoot is the path to the default config dir location.
	DefaultPathRoot = "~/" + DefaultPathName

	KeyPassword = "bitxhub"

	// API name
	APIName = "api"
)

var RootPath string

// PathRoot returns root path (default .pier)
func PathRoot() (string, error) {
	if RootPath != "" {
		return RootPath, nil
	}

	return homedir.Expand(DefaultPathRoot)
}

// PathRootWithDefault gets current config path with default value
func PathRootWithDefault(path string) (string, error) {
	var err error
	if len(path) == 0 {
		path, err = PathRoot()
		if err != nil {
			return "", err
		}
	}

	if !fileutil.Exist(path) {
		return "", fmt.Errorf("please run `premo init` first")
	}

	return path, nil
}
func KeyPath() (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", nil
	}

	return filepath.Join(repoRoot, "key.json"), nil
}

func Node1Path() (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", nil
	}

	return filepath.Join(repoRoot, "node1.json"), nil
}

func Node2Path() (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", nil
	}

	return filepath.Join(repoRoot, "node2.json"), nil
}

func Node3Path() (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", nil
	}

	return filepath.Join(repoRoot, "node3.json"), nil
}

func Node4Path() (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", nil
	}

	return filepath.Join(repoRoot, "node4.json"), nil
}

// getPrivByPath return privateKey and address by path
func getPrivByPath(path string) (crypto.PrivateKey, *types.Address, error) {
	pk, err := asym.RestorePrivateKey(path, KeyPassword)
	if err != nil {
		return nil, nil, err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return nil, nil, err
	}
	return pk, from, nil
}

// Node1Priv return node1's privateKey and address
func Node1Priv() (crypto.PrivateKey, *types.Address, error) {
	path, err := Node1Path()
	if err != nil {
		return nil, nil, err
	}
	return getPrivByPath(path)
}

// Node2Priv return node2's privateKey and address
func Node2Priv() (crypto.PrivateKey, *types.Address, error) {
	path, err := Node2Path()
	if err != nil {
		return nil, nil, err
	}
	return getPrivByPath(path)
}

// Node3Priv return node3's privateKey and address
func Node3Priv() (crypto.PrivateKey, *types.Address, error) {
	path, err := Node3Path()
	if err != nil {
		return nil, nil, err
	}
	return getPrivByPath(path)
}

// Node4Priv return node4's privateKey and address
func Node4Priv() (crypto.PrivateKey, *types.Address, error) {
	path, err := Node4Path()
	if err != nil {
		return nil, nil, err
	}
	return getPrivByPath(path)
}

// KeyPriv return privateKey and address
func KeyPriv() (crypto.PrivateKey, *types.Address, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, nil, err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return nil, nil, err
	}
	return pk, from, nil
}

func RulePath() (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", nil
	}

	return filepath.Join(repoRoot, "simple_rule.wasm"), nil
}
