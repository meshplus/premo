package utils

import (
	"encoding/hex"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym/ecdsa"
	"github.com/meshplus/bitxhub-kit/crypto/sym"
)

const padding = "abcdefghijklmnopqrstuvwxyz"

func PrivKeyFromKey(key string) (crypto.PrivateKey, error) {
	bytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	des, err := sym.GenerateKey(sym.ThirdDES, []byte("bitxhub"+padding))
	if err != nil {
		return nil, err
	}
	bytes, err = des.Decrypt(bytes)
	if err != nil {
		return nil, err
	}

	var privKey crypto.PrivateKey
	privKey, err = ecdsa.UnmarshalPrivateKey(bytes, ecdsa.Secp256r1)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}
