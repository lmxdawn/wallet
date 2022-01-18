package engine

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lmxdawn/wallet/client"
	"github.com/lmxdawn/wallet/types"
	"math/big"
	"strings"
)

type EthWorker struct {
	confirms                  uint64 // 需要的确认数
	http                      *ethclient.Client
	contractTransferEventHash common.Hash
	contractTransferHash      common.Hash
	contract                  string  // 代币合约地址，为空表示主币
	contractAbi               abi.ABI // 合约的abi
}

func NewEthWorker(confirms uint64, contract string, url string) (*EthWorker, error) {
	http := client.NewEthClient(url)

	contractTransferHashSig := []byte("transfer(address,uint256)")
	contractTransferHash := crypto.Keccak256Hash(contractTransferHashSig)
	contractTransferEventHashSig := []byte("Transfer(address,address,uint256)")
	contractTransferEventHash := crypto.Keccak256Hash(contractTransferEventHashSig)
	tokenABI := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
	contractAbi, err := abi.JSON(strings.NewReader(tokenABI))
	if err != nil {
		return nil, err
	}

	return &EthWorker{
		confirms:                  confirms,
		contract:                  contract,
		http:                      http,
		contractTransferHash:      contractTransferHash,
		contractTransferEventHash: contractTransferEventHash,
		contractAbi:               contractAbi,
	}, nil
}

type LogTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
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

func (e *EthWorker) getTransaction(num uint64) ([]types.Transaction, uint64, error) {
	if e.contract == "" {
		return e.getBlockTransaction(num)
	} else {
		return e.getTokenTransaction(num)
	}

}

func (e *EthWorker) getBlockTransaction(num uint64) ([]types.Transaction, uint64, error) {

	block, err := e.http.BlockByNumber(context.Background(), big.NewInt(int64(num)))
	if err != nil {
		return nil, num, err
	}

	chainID, err := e.http.NetworkID(context.Background())
	if err != nil {
		return nil, num, err
	}

	var transactions []types.Transaction
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
			From:        strings.ToLower(msg.From().Hex()),
			To:          strings.ToLower(tx.To().Hex()),
			Value:       tx.Value(),
		})
	}
	return transactions, num + 1, nil
}

func (e *EthWorker) getTokenTransaction(num uint64) ([]types.Transaction, uint64, error) {
	contractAddress := common.HexToAddress(e.contract)
	toBlock := num + 100
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(num)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{
			contractAddress,
		},
	}
	logs, err := e.http.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, num, err
	}

	var transactions []types.Transaction
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case e.contractTransferEventHash.Hex():

			var transferEvent LogTransfer

			err = e.contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				continue
			}

			transactions = append(transactions, types.Transaction{
				BlockNumber: big.NewInt(int64(num)),
				BlockHash:   vLog.BlockHash.Hex(),
				Hash:        vLog.TxHash.Hex(),
				From:        strings.ToLower(common.HexToAddress(vLog.Topics[1].Hex()).Hex()),
				To:          strings.ToLower(common.HexToAddress(vLog.Topics[2].Hex()).Hex()),
				Value:       transferEvent.Value,
			})
		}
	}

	return transactions, toBlock, nil
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
func (e *EthWorker) sendTransaction(privateKeyStr string, toAddress string, amount int64) (string, error) {

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

	value := big.NewInt(amount) // in wei (1 eth)
	var gasLimit uint64
	gasLimit = uint64(80000) // in units
	gasPrice, err := e.http.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	var data []byte
	toAddressHex := common.HexToAddress(toAddress)
	if e.contract != "" {
		data, err = makeERC20TransferData(e.contractTransferHash, toAddressHex, value)
		if err != nil {
			return "", err
		}
		if err != nil {
			return "", err
		}
		value = big.NewInt(0)
		toAddressHex = common.HexToAddress(e.contract)
		// 代币转账把 gasPrice 乘以 2
		//gasPrice = gasPrice.Mul(gasPrice,big.NewInt(2))
	}

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

	err = e.http.SendTransaction(context.Background(), signTx)
	if err != nil {
		return "", err
	}

	return signTx.Hash().Hex(), nil
}

func makeERC20TransferData(contractTransferHash common.Hash, toAddress common.Address, amount *big.Int) ([]byte, error) {
	var data []byte
	data = append(data, contractTransferHash[:4]...)
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	data = append(data, paddedAddress...)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data = append(data, paddedAmount...)
	return data, nil
}
