package common

import (
	"fmt"
	"sync/atomic"
	"time"

	crypto2 "github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

var AdminNonce uint64

func TransferFromAdmin(client *rpcx.ChainClient, adminPrivKey crypto2.PrivateKey, adminFrom *types.Address, address *types.Address, amount string) error {
	data := &pb.TransactionData{
		Amount: amount + "000000000000000000",
	}
	payload, err := data.Marshal()
	if err != nil {
		return err
	}

	tx := &pb.BxhTransaction{
		From:      adminFrom,
		To:        address,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:    adminFrom.String(),
		Nonce:   atomic.AddUint64(&AdminNonce, 1) - 1,
		PrivKey: adminPrivKey,
	})
	if err != nil {
		return err
	}
	if ret.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(ret.Ret))
	}
	return nil
}
