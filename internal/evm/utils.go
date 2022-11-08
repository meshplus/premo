package evm

import "math/big"

type signStruct struct {
	HashedMessage [32]byte
	V             uint8
	R             [32]byte
	S             [32]byte
}
type orgStruct struct {
	OrgId *big.Int
	BxmId string
	Extra string
	Sign  signStruct
}

type userStruct struct {
	UserAddr string
	OrgId    *big.Int
	Credit   *big.Int
	Extra    string
}

type userInput struct {
	User userStruct
	Sign signStruct
}

type creditPackage struct {
	Credit   *big.Int
	Quantity uint8
	Duration *big.Int
}
