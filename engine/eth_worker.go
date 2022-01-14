package engine

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/lmxdawn/wallet/client"
	"github.com/lmxdawn/wallet/types"
	"math"
	"math/big"
)

type EthWorker struct {
	confirms uint64 // 需要的确认数
	contract string // 代币合约地址，为空表示主币
	http     *ethclient.Client
}

func NewEthWorker(confirms uint64, contract string, url string) *EthWorker {
	http := client.NewEthClient(url)
	return &EthWorker{
		confirms: confirms,
		contract: contract,
		http:     http,
	}
}

func (e *EthWorker) getTransactionReceipt(transaction *types.Transaction) error {

	hash := common.HexToHash(transaction.Hash)

	receipt, err := e.http.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return err
	}

	// 获取最新区块
	latest, err := e.http.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// 判断确认数
	confirms := latest - receipt.BlockNumber.Uint64() + 1
	if confirms < e.confirms {
		return errors.New("the number of confirmations is not satisfied")
	}

	status := receipt.Status
	transaction.Status = uint(status)

	return nil
}

func (e *EthWorker) getTransaction(num uint64) ([]types.Transaction, error) {
	block, err := e.http.BlockByNumber(context.Background(), big.NewInt(int64(num)))
	if err != nil {
		return nil, err
	}
	var transactions []types.Transaction

	chainID, err := e.http.NetworkID(context.Background())
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
			BlockNumber: big.NewInt(int64(num)),
			BlockHash:   block.Hash().Hex(),
			Hash:        tx.Hash().Hex(),
			From:        msg.From().Hex(),
			To:          tx.To().Hex(),
			Value:       tx.Value(),
		})
	}
	return transactions, nil
}

func (e *EthWorker) createWallet() (*types.Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	privateKeyString := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyString := hexutil.Encode(publicKeyBytes)[4:]

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &types.Wallet{
		Address:    address,
		PublicKey:  publicKeyString,
		PrivateKey: privateKeyString,
	}, err
}

// sendTransaction 创建并发送裸交易
func (e *EthWorker) sendTransaction(privateKeyStr string, toAddress string, amount int64, decimals int) (string, error) {

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.http.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	value := big.NewInt(amount * int64(math.Pow10(decimals))) // in wei (1 eth)
	gasLimit := uint64(21000)                                 // in units
	gasPrice, err := e.http.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	toAddressHex := common.HexToAddress(toAddress)
	var data []byte
	txData := &ethTypes.LegacyTx{
		Nonce:    nonce,
		To:       &toAddressHex,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}
	tx := ethTypes.NewTx(txData)

	chainID, err := e.http.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signTx, err := ethTypes.SignTx(tx, ethTypes.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	rawTxBytes, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", err
	}
	rawTxHex := hex.EncodeToString(rawTxBytes)

	txSend := new(ethTypes.Transaction)
	err = rlp.DecodeBytes(rawTxBytes, &txSend)
	if err != nil {
		return "", err
	}
	err = e.http.SendTransaction(context.Background(), txSend)
	if err != nil {
		return "", err
	}

	return rawTxHex, nil
}
