package gosdk_tester

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

func (suite *Snake) TestStopClient() {
	keyPath, err := repo.KeyPath()
	suite.Require().Nil(err)

	node0 := &rpcx.NodeInfo{Addr: "localhost:60011"}

	pk, err := asym.RestorePrivateKey(keyPath, repo.KeyPassword)

	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)

	err = client.Stop()
	suite.Require().Nil(err)
}

func (suite *Snake) TestSetPrivateKey() {
	suite.client.SetPrivateKey(nil)

	keyPath, err := repo.KeyPath()
	suite.Require().Nil(err)

	pk, err := asym.RestorePrivateKey(keyPath, "bitxhub")

	suite.client.SetPrivateKey(pk)
	suite.Require().Equal(pk, suite.pk)
}

func (suite *Snake) TestSendViewIsTrue() {
	BoltContractAddress := "0x000000000000000000000000000000000000000b"

	rand.Seed(time.Now().UnixNano())
	randKey := make([]byte, 20)
	_, err := rand.Read(randKey)
	suite.Require().Nil(err)

	tx1, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		types.NewAddressByStr(BoltContractAddress), "Set", pb.String(string(randKey)), pb.String("value"))
	suite.Require().Nil(err)
	tx1.Nonce = 1

	err = tx1.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt1, err := suite.client.SendView(tx1)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt1.Status)

	tx2, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		types.NewAddressByStr(BoltContractAddress), "Get", pb.String(string(randKey)))
	suite.Require().Nil(err)
	tx2.Nonce = 2

	err = tx2.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt2, err := suite.client.SendView(tx2)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, receipt2.Status)
}

func (suite *Snake) TestSendViewIsFalse() {
	BoltContractAddress := "0x000000000000000000000000000000000000000b"

	rand.Seed(time.Now().UnixNano())
	randKey := make([]byte, 20)
	_, err := rand.Read(randKey)
	suite.Require().Nil(err)

	tx1, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		types.NewAddressByStr(BoltContractAddress), "Set", pb.String(string(randKey)), pb.String("value"))
	suite.Require().Nil(err)
	tx1.Nonce = 1
	tx1.Payload = nil

	err = tx1.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendView(tx1)
	suite.Require().NotNil(err)

	tx2, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		types.NewAddressByStr(BoltContractAddress), "set", pb.String(string(randKey)), pb.String("value"))
	suite.Require().Nil(err)
	tx2.Nonce = 1

	err = tx2.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt, err := suite.client.SendView(tx2)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, receipt.Status)
}

func (suite Snake) TestSendTransactionIsTrue() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	receipt, err := suite.client.GetReceipt(hash)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
}

func (suite Snake) TestSendTransactionIsFalse() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 0,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
	suite.Contains(err.Error(), "tx payload and ibtp can't both be nil")
}

func (suite Snake) TestSendTransactionWithReceiptIsTrue() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
}

func (suite *Snake) TestSendTransactionWithReceiptWhenToIsNull() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From: suite.from,
		//To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}

	_, err = suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "tx to address is nil")
}

func (suite *Snake) TestSendTransactionWithReceiptWhenFromIsNull() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		//From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "from address can't be empty")
}

//TestGetReceiptByHashIsTrue same with TestSendTransactionIsTrue
func (suite *Snake) TestGetReceiptByHashIsFalse() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	receipt, err := suite.client.GetReceipt(hash + "1")
	suite.Require().Nil(err)
	suite.Require().Nil(receipt.Ret)
}

func (suite *Snake) TestGetTransactionIsTrue() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	receipt, err := suite.client.GetReceipt(hash)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	_, err = suite.client.GetTransaction(receipt.TxHash.String())
	suite.Require().Nil(err)
}

func (suite *Snake) TestGetTransactionIsFalse() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	receipt, err := suite.client.GetReceipt(hash)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	wronghash := receipt.TxHash.String() + "a123"
	transaction, err := suite.client.GetTransaction(wronghash)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "invalid format of tx hash for querying transaction")
	suite.Require().Nil(transaction)
}

func (suite *Snake) TestGetChainMeta() {
	meta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)
	suite.Require().NotNil(meta)
}

func (suite Snake) TestGetBlocksIsTrue() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt1, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt1.Status)

	transaction1, err := suite.client.GetTransaction(receipt1.TxHash.String())
	suite.Require().Nil(err)
	height1 := transaction1.GetTxMeta().BlockHeight

	receipt2, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt2.Status)

	transaction2, err := suite.client.GetTransaction(receipt2.TxHash.String())
	suite.Require().Nil(err)
	height2 := transaction2.GetTxMeta().BlockHeight

	blocks, err := suite.client.GetBlocks(height1, height2)
	suite.Require().Nil(err)
	l := len(blocks.Blocks)
	suite.Require().Equal(2, l)
}

func (suite Snake) TestGetBlocksIsFalse() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt1, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt1.Status)

	transaction1, err := suite.client.GetTransaction(receipt1.TxHash.String())
	suite.Require().Nil(err)
	height1 := transaction1.GetTxMeta().BlockHeight

	receipt2, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt2.Status)

	transaction2, err := suite.client.GetTransaction(receipt2.TxHash.String())
	suite.Require().Nil(err)
	height2 := transaction2.GetTxMeta().BlockHeight

	//height2>height1
	blocks, err := suite.client.GetBlocks(height2, height1)
	suite.Require().Nil(err)
	l := len(blocks.Blocks)
	suite.Require().Equal(0, l)

	//get more blocks
	blocks, err = suite.client.GetBlocks(height1, height2+100)
	suite.Require().Nil(err)
	l = len(blocks.Blocks)
	suite.Require().Equal(2, l)

	//get does not exist blocks
	blocks, err = suite.client.GetBlocks(height1+100, height2+100)
	suite.Require().Nil(err)
	l = len(blocks.Blocks)
	suite.Require().Equal(0, l)

	//get Illegal height blocks
	blocks, err = suite.client.GetBlocks(0, 0)
	suite.Require().Nil(err)
	l = len(blocks.Blocks)
	suite.Require().Equal(0, l)
}

func (suite *Snake) TestGetBlockIsTrue() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	transaction, err := suite.client.GetTransaction(receipt.TxHash.String())
	suite.Require().Nil(err)

	block1, err := suite.client.GetBlock(fmt.Sprint(transaction.GetTxMeta().GetBlockHeight()), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)

	block2, err := suite.client.GetBlock(fmt.Sprint(types.NewHash(transaction.GetTxMeta().GetBlockHash())), pb.GetBlockRequest_HASH)
	suite.Require().Nil(err)

	suite.Require().Equal(block1, block2)
}

func (suite *Snake) TestGetBlockIsFalse() {
	td := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     1,
		Payload:   payload,
	}
	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	transaction, err := suite.client.GetTransaction(receipt.TxHash.String())
	suite.Require().Nil(err)

	block1, err := suite.client.GetBlock(fmt.Sprint(transaction.GetTxMeta().GetBlockHeight()+1), pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Nil(block1)

	block2, err := suite.client.GetBlock(fmt.Sprint(types.NewHash(transaction.GetTxMeta().GetBlockHash()))+"123", pb.GetBlockRequest_HASH)
	suite.Require().NotNil(err)
	suite.Require().Nil(block2)
}

func (suite *Snake) TestGetChainStatus() {
	status, err := suite.client.GetChainStatus()
	suite.Require().Nil(err)
	suite.Equal("normal", string(status.Data))
}

func (suite *Snake) TestGetValidators() {
	Validators, err := suite.client.GetValidators()
	suite.Require().Nil(err)
	suite.Require().NotNil(Validators)
}

func (suite *Snake) TestGetNetworkMeta() {
	meta, err := suite.client.GetNetworkMeta()
	suite.Require().Nil(err)
	suite.Require().NotNil(meta)
}

func (suite Snake) TestGetAccountBalanceIsTrue() {
	address, err := suite.pk.PublicKey().Address()
	suite.Require().Nil(err)

	balance, err := suite.client.GetAccountBalance(address.String())
	suite.Require().Nil(err)
	suite.Require().NotNil(balance)
}

func (suite Snake) TestGetAccountBalanceIsFalse() {
	address, err := suite.pk.PublicKey().Address()
	suite.Require().Nil(err)

	balance, err := suite.client.GetAccountBalance(address.String() + "123")
	suite.Require().NotNil(err)
	suite.Contains(err.Error(), "invalid account address")
	suite.Require().Nil(balance)
}

func (suite *Snake) TestGetBlockHeaderIsTrue() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	headers := make(chan *pb.BlockHeader)

	err := suite.client.GetBlockHeader(ctx, 1, 2, headers)
	suite.Require().Nil(err)
	for {
		select {
		case header, ok := <-headers:
			suite.Require().Equal(true, ok)
			suite.Require().Equal(uint64(1), header.Number)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (suite *Snake) TestGetBlockHeaderIsFalse() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	headers := make(chan *pb.BlockHeader)
	err := suite.client.GetBlockHeader(ctx, 2, 1, headers)
	suite.Require().Nil(err)
	for {
		select {
		case header, ok := <-headers:
			suite.Require().Equal(false, ok)
			suite.Require().Nil(header)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (suite *Snake) TestGetInterchainTxWrappersIsTrue() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	//sendInterchain
	_, _, _, to, receipt, err := suite.sendInterchain()
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
	//get
	meta, err := suite.client.GetChainMeta()
	ch := make(chan *pb.InterchainTxWrappers, 10)
	err = suite.client.GetInterchainTxWrappers(ctx, to.String(), meta.Height, meta.Height+100, ch)
	suite.Require().Nil(err)

	for {
		select {
		case wrappers, ok := <-ch:
			suite.Require().Equal(true, ok)
			suite.Require().NotNil(wrappers.InterchainTxWrappers[0])
			wrapper := wrappers.InterchainTxWrappers[0]
			suite.Require().GreaterOrEqual(meta.Height, wrapper.Height)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (suite *Snake) TestGetInterchainTxWrappersIsFalse() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	//sendInterchain
	_, _, _, to, receipt, err := suite.sendInterchain()
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	//get
	meta, err := suite.client.GetChainMeta()
	ch := make(chan *pb.InterchainTxWrappers, 10)
	err = suite.client.GetInterchainTxWrappers(ctx, to.String(), meta.Height, meta.Height-1, ch)
	suite.Require().Nil(err)

	for {
		select {
		case wrappers, ok := <-ch:
			suite.Require().Equal(false, ok)
			suite.Require().Nil(wrappers)
			goto label1
		case <-ctx.Done():
			return
		}
	}
label1:
	ch = make(chan *pb.InterchainTxWrappers, 10)
	err = suite.client.GetInterchainTxWrappers(ctx, to.String()+"123", meta.Height, meta.Height+100, ch)
	suite.Require().Nil(err)

	for {
		select {
		case wrappers, ok := <-ch:
			suite.Require().Equal(true, ok)
			suite.Require().Nil(wrappers.InterchainTxWrappers[0].Transactions)
			goto label2
		case <-ctx.Done():
			return
		}
	}
label2:
	ch = make(chan *pb.InterchainTxWrappers, 10)
	err = suite.client.GetInterchainTxWrappers(ctx, to.String()+"123", meta.Height, meta.Height-1, ch)
	suite.Require().Nil(err)

	for {
		select {
		case wrappers, ok := <-ch:
			suite.Require().Equal(false, ok)
			suite.Require().Nil(wrappers)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (suite *Snake) TestSubscribe_BLOCK() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c, err := suite.client.Subscribe(ctx, pb.SubscriptionRequest_BLOCK, nil)
	suite.Require().Nil(err)

	td := &pb.TransactionData{
		Amount: 10,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	for {
		select {
		case block, ok := <-c:
			suite.Require().Equal(true, ok)
			suite.Require().NotNil(block)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (suite *Snake) TestSubscribe_BLOCK_HEADER() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c, err := suite.client.Subscribe(ctx, pb.SubscriptionRequest_BLOCK_HEADER, nil)
	suite.Require().Nil(err)

	td := &pb.TransactionData{
		Amount: 10,
	}
	payload, err := td.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	err = tx.Sign(suite.pk)
	suite.Require().Nil(err)

	receipt, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	for {
		select {
		case block_header, ok := <-c:
			suite.Require().Equal(true, ok)
			suite.Require().NotNil(block_header)
			return
		case <-ctx.Done():
			return
		}
	}
}

/*func (suite *Snake) TestSubscribe_INTERCHAIN_TX() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	c, err := suite.client.Subscribe(ctx, pb.SubscriptionRequest_INTERCHAIN_TX, nil)
	suite.Require().Nil(err)
	//sendInterchain
	_, _, _, _, receipt, err := suite.sendInterchain()
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS,receipt.Status)
	for {
		select {
		case interchainTx, ok := <-c:
			suite.Require().Equal(true,ok)
			suite.Require().NotNil(interchainTx)
			fmt.Println(interchainTx)
			return
		case <-ctx.Done():
			return
		}
	}
}*/

func (suite *Snake) TestSubscribe_INTERCHAIN_TX_WRAPPER() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c, err := suite.client.Subscribe(ctx, pb.SubscriptionRequest_INTERCHAIN_TX_WRAPPER, []byte(suite.from.String()))
	suite.Require().Nil(err)

	//sendInterchain
	_, _, _, _, receipt, err := suite.sendInterchain()
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	for {
		select {
		case interchainTxWrapper, ok := <-c:
			suite.Require().Equal(true, ok)
			suite.Require().NotNil(interchainTxWrapper)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (suite *Snake) TestSubscribe_UNION_INTERCHAIN_TX_WRAPPER() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c, err := suite.client.Subscribe(ctx, pb.SubscriptionRequest_UNION_INTERCHAIN_TX_WRAPPER, []byte(suite.from.String()))
	suite.Require().Nil(err)

	//sendInterchain
	_, _, _, _, receipt, err := suite.sendInterchain()
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	for {
		select {
		case interchainTxWrapper, ok := <-c:
			suite.Require().Equal(true, ok)
			suite.Require().NotNil(interchainTxWrapper)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (suite *Snake) TestDeployContract() {
	contract, err := ioutil.ReadFile("../bxh_tester/testdata/example.wasm")
	suite.Require().Nil(err)

	address, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
}

/*
func (suite *Snake) TestGenerateContractTx() {
	BoltContractAddress := "0x000000000000000000000000000000000000000b"
	rand.Seed(time.Now().UnixNano())
	randKey := make([]byte, 20)
	_, err := rand.Read(randKey)
	suite.Require().Nil(err)
	tx1, err := suite.client.GenerateContractTx(pb.TransactionData_BVM, types.NewAddressByStr(BoltContractAddress),
		"Set", pb.String(string(randKey)), pb.String("value"))
	suite.Require().Nil(err)
	tx1.Nonce = 1
	err = tx1.Sign(suite.pk)
	suite.Require().Nil(err)
	receipt1, err := suite.client.SendView(tx1)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS,receipt1.Status)
	tx2, err := suite.client.GenerateContractTx(pb.TransactionData_BVM,
		types.NewAddressByStr(BoltContractAddress), "Get", pb.String(string(randKey)))
	suite.Require().Nil(err)
	tx2.Nonce = 2
	err = tx2.Sign(suite.pk)
	suite.Require().Nil(err)
	receipt2, err := suite.client.SendView(tx2)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS,receipt2.Status)
}
*/

func (suite *Snake) TestInvokeXVMContractIsTrue() {
	contract, err := ioutil.ReadFile("../bxh_tester/testdata/example.wasm")
	suite.Require().Nil(err)

	address, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)

	receipt, err := suite.client.InvokeXVMContract(address, "a", nil, rpcx.Int32(1), rpcx.Int32(2))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
	suite.Require().Equal("336", string(receipt.Ret))
}

func (suite *Snake) TestInvokeXVMContractIsFalse() {
	contract, err := ioutil.ReadFile("../bxh_tester/testdata/example.wasm")
	suite.Require().Nil(err)

	address, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)

	receipt, err := suite.client.InvokeXVMContract(address, "abc", nil, rpcx.Int32(1), rpcx.Int32(2))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, receipt.Status)
}

func (suite Snake) TestInvokeBVMContractIsTrue() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(),
		"Set", nil, rpcx.String("a"), rpcx.String("10"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	receipt, err = suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(),
		"Get", nil, rpcx.String("a"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
	suite.Require().Equal("10", string(receipt.Ret))
}

func (suite Snake) TestInvokeBVMContractIsFalse() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(),
		"set", nil, rpcx.String("a"), rpcx.String("10"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, receipt.Status)

	receipt, err = suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(),
		"get", nil, rpcx.String("a"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, receipt.Status)
}

func (suite *Snake) TestGetTPSIsTrue() {
	BoltContractAddress := "0x000000000000000000000000000000000000000b"

	rand.Seed(time.Now().UnixNano())
	randKey := make([]byte, 20)
	_, err := rand.Read(randKey)
	suite.Require().Nil(err)

	tx1, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		types.NewAddressByStr(BoltContractAddress), "Set", pb.String(string(randKey)), pb.String("value"))
	suite.Require().Nil(err)

	err = tx1.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendTransactionWithReceipt(tx1, nil)
	suite.Require().Nil(err)

	meta0, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	for i := 0; i < 10; i++ {
		tx2, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
			types.NewAddressByStr(BoltContractAddress), "Set", pb.String(string(randKey)), pb.String("value"))
		suite.Require().Nil(err)

		err = tx2.Sign(suite.pk)
		suite.Require().Nil(err)

		_, err = suite.client.SendTransaction(tx2, nil)
		suite.Require().Nil(err)
	}

	time.Sleep(time.Second)
	meta1, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	tps, err := suite.client.GetTPS(meta0.Height, meta1.Height)
	suite.Require().Nil(err)
	suite.Require().NotNil(tps)
	suite.Require().True(tps > 0)
}

func (suite *Snake) TestGetTPSIsFalse() {
	BoltContractAddress := "0x000000000000000000000000000000000000000b"

	rand.Seed(time.Now().UnixNano())
	randKey := make([]byte, 20)
	_, err := rand.Read(randKey)
	suite.Require().Nil(err)

	meta0, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	for i := 0; i < 10; i++ {
		tx, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
			types.NewAddressByStr(BoltContractAddress), "Set", pb.String(string("a")), pb.String("1"))
		suite.Require().Nil(err)

		err = tx.Sign(suite.pk)
		suite.Require().Nil(err)

		_, err = suite.client.SendTransaction(tx, nil)
		suite.Require().Nil(err)
	}
	time.Sleep(time.Second)

	meta1, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	_, err = suite.client.GetTPS(meta1.Height, meta0.Height)
	suite.Require().NotNil(err)
}

func (suite Snake) TestGetPendingNonceByAccountIsTrue() {
	nextNonce, err := suite.client.GetPendingNonceByAccount(suite.from.String())
	suite.Require().Nil(err)
	suite.Require().NotNil(nextNonce)
	suite.Require().True(nextNonce > 0)
}

func (suite Snake) TestGetPendingNonceByAccountIsFalse() {
	nextNonce, err := suite.client.GetPendingNonceByAccount(suite.from.String() + "123456")
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(1), nextNonce)
}

func genContractTransaction(
	vmType pb.TransactionData_VMType, privateKey crypto.PrivateKey,
	address *types.Address, method string, args ...*pb.Arg) (*pb.Transaction, error) {
	from, err := privateKey.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	pl := &pb.InvokePayload{
		Method: method,
		Args:   args[:],
	}

	data, err := pl.Marshal()
	if err != nil {
		return nil, err
	}

	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  vmType,
		Payload: data,
	}

	payload, err := td.Marshal()
	if err != nil {
		return nil, err
	}

	tx := &pb.Transaction{
		From:      from,
		To:        address,
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
	}

	if err := tx.Sign(privateKey); err != nil {
		return nil, fmt.Errorf("tx sign: %w", err)
	}
	tx.TransactionHash = tx.Hash()

	return tx, nil
}

func (suite *Snake) prepare() (crypto.PrivateKey, crypto.PrivateKey, *types.Address, *types.Address) {
	kA := suite.pk
	//suite.Require().Nil(err)
	kB, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	from, err := kA.PublicKey().Address()
	suite.Require().Nil(err)
	to, err := kB.PublicKey().Address()
	suite.Require().Nil(err)

	return kA, kB, from, to
}

func (suite *Snake) RegisterAppchain(pk crypto.PrivateKey, chainType string) {
	pubBytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	suite.client.SetPrivateKey(pk)
	var pubKeyStr = hex.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String(chainType),          //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	appChain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appChain)
	suite.Require().Nil(err)
	suite.Require().NotNil(appChain.ID)
}

func (suite *Snake) RegisterRule(pk crypto.PrivateKey, ruleFile string) {
	suite.client.SetPrivateKey(pk)

	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	// deploy rule
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := suite.client.DeployContract(bytes, nil)
	suite.Require().Nil(err)

	// register rule
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, pb.String(from.String()), pb.String(addr.String()))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
}

func (suite Snake) sendInterchain() (crypto.PrivateKey, crypto.PrivateKey, *types.Address, *types.Address, *pb.Receipt, error) {
	//sendInterchain
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "../../config/rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return kA, kB, from, to, res, nil
}
