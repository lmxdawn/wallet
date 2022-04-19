package tron

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"net/url"
)

type TronWeb struct {
	privateKey  string
	formAddress string
	client      *Client
	eventServer *EventServer
}

func NewTronWeb(node, eventServerUrl, privateKey string, withTLS bool) (*TronWeb, error) {
	var formAddress string
	var err error
	if privateKey != "" {
		formAddress, err = Private2TronAddress(privateKey)
		if err != nil {
			return nil, err
		}
	}
	var eventServer *EventServer
	if eventServerUrl != "" {
		eventServer = NewEventServer(eventServerUrl)
	}
	client, err := NewClient(node, withTLS)
	if err != nil {
		return nil, err
	}
	return &TronWeb{
		privateKey:  privateKey,
		formAddress: formAddress,
		client:      client,
		eventServer: eventServer,
	}, nil

}

func (t *TronWeb) CallContract(contractAddress, method, jsonString string) (string, error) {

	result, err := t.client.TriggerContract("", contractAddress, method, jsonString, true, 0)
	if err != nil {
		return "", err
	}

	data := common.BytesToHexString(result.GetConstantResult()[0])

	return data, nil

}

func (t *TronWeb) SendContract(privateKey, contractAddress, method, jsonString string) (string, error) {

	fmt.Println(t.formAddress, contractAddress, method, jsonString)
	tx, err := t.client.TriggerContract(t.formAddress, contractAddress, method, jsonString, false, 10000000000)
	if err != nil {
		return "", err
	}

	txid := common.BytesToHexString(tx.Txid)

	signTx, err := SignTransaction(tx.Transaction, privateKey)
	if err != nil {
		return "", err
	}

	err = t.client.BroadcastTransaction(signTx)
	if err != nil {
		return "", err
	}

	// 查询交易
	r, err := t.client.GetTransactionInfoByID(txid, true)
	if err != nil {
		return "", err
	}

	if r.Receipt.Result != 1 {
		return "", errors.New(string(r.ResMessage))
	}

	return common.BytesToHexString(r.ContractResult[0]), nil

}

// Transfer 转账/TRX/TRX10
func (t *TronWeb) Transfer(privateKey, from, to, assetName string, amount int64) (string, error) {
	var tx *api.TransactionExtention
	var err error
	if assetName != "" {
		tx, err = t.client.TransferTrc10(from, to, assetName, amount)

	} else {
		tx, err = t.client.Transfer(from, to, amount)
	}
	if err != nil {
		return "", err
	}

	txid := common.BytesToHexString(tx.Txid)

	signTx, err := SignTransaction(tx.Transaction, privateKey)
	if err != nil {
		return "", err
	}

	err = t.client.BroadcastTransaction(signTx)
	if err != nil {
		return "", err
	}

	return txid, nil

}

// GetTransactionInfoByID 查询交易信息
func (t *TronWeb) GetTransactionInfoByID(txid string, isRes bool) (*core.TransactionInfo, error) {

	return t.client.GetTransactionInfoByID(txid, isRes)

}

// GetBlockByNum 根据区块号获取区块信息
func (t *TronWeb) GetBlockByNum(num int64) (*api.BlockExtention, error) {

	return t.client.GetBlockByNum(num)

}

// GetNowBlockNum 获取最新区块
func (t *TronWeb) GetNowBlockNum() (int64, error) {

	block, err := t.client.GetNowBlock()
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(block.Blockid)), nil

}

// GetBalance 获取账户TRX余额
func (t *TronWeb) GetBalance(address string, token string) (int64, error) {

	account, err := t.client.GetAccount(address)

	if err != nil {
		return 0, err
	}

	if token != "" {
		if account.AssetV2 == nil {
			return 0, nil
		}
		return account.AssetV2[token], nil
	}

	return account.Balance, nil
}

// GetEventResult 获取事件的历史
func (t *TronWeb) GetEventResult(contractAddress string, blockNumber int64, res interface{}) error {

	params := url.Values{}
	if blockNumber > 0 {
		params.Set("block_number", fmt.Sprintf("%d", blockNumber))
	}
	err := t.eventServer.GetContractsEventsParams(contractAddress, params, &res)
	if err != nil {
		return err
	}
	return nil
}

// GetEventResultParams 获取事件的历史
func (t *TronWeb) GetEventResultParams(contractAddress string, params url.Values, res interface{}) error {

	err := t.eventServer.GetContractsEventsParams(contractAddress, params, &res)
	if err != nil {
		return err
	}
	return nil
}

func (t *TronWeb) GetTransaction(start, end int64) (*api.BlockListExtention, error) {

	return t.client.GetBlock(start, end)

}
