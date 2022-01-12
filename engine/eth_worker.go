package engine

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lmxdawn/wallet/types"
	"math/big"
)

type EthWorker struct {
	Confirms uint64 // 需要的确认数
	client   *ethclient.Client
}

func NewEthWorker(confirms uint64, url string) *EthWorker {
	client, _ := ethclient.Dial(url)
	return &EthWorker{
		Confirms: confirms,
		client:   client,
	}
}

func (e *EthWorker) getTransactionReceipt(transaction *types.Transaction) error {

	hash := common.HexToHash(transaction.Hash)

	receipt, err := e.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return err
	}

	// 获取最新区块
	latest, err := e.client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// 判断确认数
	confirms := latest - receipt.BlockNumber.Uint64() + 1
	if confirms < e.Confirms {
		return errors.New("the number of confirmations is not satisfied")
	}

	status := receipt.Status
	transaction.Status = uint(status)

	return nil
}

func (e *EthWorker) getTransaction(num int64) ([]types.Transaction, error) {
	block, err := e.client.BlockByNumber(context.Background(), big.NewInt(num))
	if err != nil {
		return nil, err
	}
	var transactions []types.Transaction

	chainID, err := e.client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	for _, tx := range block.Transactions() {
		// 如果接收方地址为空，则是创建合约的交易，忽略过去
		if tx.To() == nil {
			continue
		}
		msg, err := tx.AsMessage(ethTypes.LatestSignerForChainID(chainID), tx.GasPrice())
		if err != nil {
			continue
		}
		transactions = append(transactions, types.Transaction{
			BlockNumber: big.NewInt(num),
			BlockHash:   block.Hash().Hex(),
			Hash:        tx.Hash().Hex(),
			From:        msg.From().Hex(),
			To:          tx.To().Hex(),
			Value:       tx.Value(),
		})
	}
	return transactions, nil
}
