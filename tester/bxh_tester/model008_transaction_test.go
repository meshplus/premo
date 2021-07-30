package bxh_tester

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

type Model8 struct {
	*Snake
}

func (suite *Model8) SetupTest() {
	suite.T().Parallel()
}

func (suite *Model8) Test0801_TXEmptyFrom() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

func (suite *Model8) Test0802_TXEmptyTo() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

/*增加form和to都为空*/
func (suite *Model8) Test0803_TXEmptyFromAndTo() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

/*增加from和to相同*/
func (suite *Model8) Test0804_TXSameFromAndTo() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        from,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

func (suite *Model8) Test0805_TXEmptySig() {
	//node2, err := repo.Node2Path()
	//suite.Require().Nil(err)
	//pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	//suite.Require().Nil(err)
	//from, err := pk.PublicKey().Address()
	//suite.Require().Nil(err)
	//client := suite.NewClient(pk)
	//data := &pb.TransactionData{
	//	Amount: 1,
	//}
	//payload, err := data.Marshal()
	//suite.Require().Nil(err)
	//
	//tx := &pb.BxhTransaction{
	//	From:      from,
	//	To:        suite.to,
	//	Timestamp: time.Now().UnixNano(),
	//	Payload:   payload,
	//}
	//nonce := atomic.AddUint64(&nonce2, 1)
	//_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
	//	Nonce: nonce,
	//})
	//suite.Require().Nil(err)
}

func (suite *Model8) Test0806_TXWrongSigPrivateKey() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	client1 := suite.NewClient(pk1)

	hash, err := client1.SendTransaction(tx, nil)
	suite.Require().NotNil(err)

	_, err = client.GetReceipt(hash)
	suite.Require().NotNil(err)
}

func (suite *Model8) Test0807_TXWrongSigAlgorithm() {
	// K1
}

func (suite *Model8) Test0808_TXExtra10MB() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	MB10 := make([]byte, 10*1024*1024) // 10MB
	for i := 0; i < len(MB10); i++ {
		MB10[i] = uint8(rand.Intn(255))
	}

	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Extra:     MB10,
		Payload:   payload,
	}

	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "larger than max")
}

func (suite *Model8) Test0809_GetTxByHash() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.AddUint64(&nonce2, 1)
	hash, err := client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().Nil(err)

	var ret *pb.GetTransactionResponse
	err1 := retry.Retry(func(attempt uint) error {
		pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
		suite.Require().Nil(err)
		client1 := suite.NewClient(pk)
		ret, err = client1.GetTransaction(hash)
		if err != nil {
			return err
		}
		return nil
	},
		strategy.Limit(10),
		strategy.Backoff(backoff.Fibonacci(500*time.Millisecond)),
	)
	suite.Require().Nil(err1)
	suite.Require().NotNil(ret)
}

func (suite *Model8) Test0810_GetReceiptByHash() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	nonce := atomic.AddUint64(&nonce2, 1)
	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, ret.Status)
}

/*通过错误的hash值进行查询*/
func (suite *Model8) Test0811_GetReceiptByWrongHash() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	nonce := atomic.AddUint64(&nonce2, 1)
	hash, err := client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})

	hash = hash[0:len(hash)-5] + "12345"
	ret, err := client.GetReceipt(hash)
	suite.Require().Nil(ret)
}
