package repo

import (
	"fmt"
	"path/filepath"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/fileutil"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/mitchellh/go-homedir"
)

const (
	// DefaultPathName is the default config dir name
	DefaultPathName = ".premo"
	// DefaultPathRoot is the path to the default config dir location.
	DefaultPathRoot = "~/" + DefaultPathName
	// KeyPassword is the password to generate privkey
	KeyPassword = "bitxhub"
	// APIName is api's Name
	APIName = "api"
)

// PathRoot returns root path (default .pier)
func PathRoot() (string, error) {
	return homedir.Expand(DefaultPathRoot)
}

// PathRootWithDefault gets current config path with default value
func PathRootWithDefault() (string, error) {
	path, err := PathRoot()
	if err != nil {
		return "", err
	}
	if !fileutil.Exist(path) {
		return "", fmt.Errorf("please run `premo init` first")
	}
	return path, nil
}

// filePath return filepath by filename with default path
func filePath(filename string) (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", nil
	}
	return filepath.Join(repoRoot, filename), nil
}

// Node1Path return node1.json path
func Node1Path() (string, error) {
	return filePath("node1.json")
}

// Node2Path return node2.json path
func Node2Path() (string, error) {
	return filePath("node2.json")
}

// Node3Path return node3.json path
func Node3Path() (string, error) {
	return filePath("node3.json")
}

// Node4Path return node4.json path
func Node4Path() (string, error) {
	return filePath("node4.json")
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

// KeyPriv return node4's privateKey and address
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
