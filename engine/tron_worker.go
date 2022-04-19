package engine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
	"github.com/lmxdawn/wallet/tron"
	"github.com/lmxdawn/wallet/types"
	"math/big"
)

type TronWorker struct {
	confirms  uint64 // 需要的确认数
	http      *tron.TronWeb
	token     string // 代币合约地址，为空表示主币
	tokenType string // 合约类型（trc10、trc20）区分
}

func NewTronWorker(confirms uint64, token string, tokenType string, url string) (*TronWorker, error) {
	http, err := tron.NewTronWeb(url, "", "", false)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &TronWorker{
		confirms:  confirms,
		http:      http,
		token:     token,
		tokenType: tokenType,
	}, nil
}

// GetNowBlockNum 获取最新块
func (e *TronWorker) GetNowBlockNum() (uint64, error) {
	blockNumber, err := e.http.GetNowBlockNum()
	if err != nil {
		return 0, err
	}
	return uint64(blockNumber), nil
}

// GetTransactionReceipt 获取交易的票据
func (e *TronWorker) GetTransactionReceipt(transaction *types.Transaction) error {
	receipt, err := e.http.GetTransactionInfoByID(transaction.Hash, false)
	if err != nil {
		return err
	}

	// 获取最新区块
	latest, err := e.http.GetNowBlockNum()
	if err != nil {
		return err
	}

	// 判断确认数
	confirms := latest - receipt.BlockNumber + 1
	if uint64(confirms) < e.confirms {
		return errors.New(fmt.Sprintf("哈希：%v，当前确认块：%d，需要确认块：%d", transaction.Hash, confirms, e.confirms))
	}

	status := 0
	if receipt.Result == core.TransactionInfo_SUCESS {
		status = 1
	}
	transaction.Status = uint(status)

	return nil

}

// GetTransaction 获取交易信息
func (e *TronWorker) GetTransaction(num uint64) ([]types.Transaction, uint64, error) {
	nowBlockNumber, err := e.GetNowBlockNum()
	if err != nil {
		return nil, num, err
	}
	toBlock := num + 100
	// 传入的num为0，表示最新块
	if num == 0 {
		toBlock = nowBlockNumber
	} else if toBlock > nowBlockNumber {
		toBlock = nowBlockNumber
	}

	numInt := int64(num)
	toBlockInt := int64(toBlock)

	blocks, err := e.http.GetTransaction(numInt, toBlockInt)
	if err != nil {
		return nil, num, err
	}

	var transactions []types.Transaction
	for _, block := range blocks.Block {
		for _, v := range block.Transactions {
			if v.Result == nil || !v.Result.Result {
				continue
			}

			rets := v.Transaction.Ret
			if len(rets) < 1 || rets[0].ContractRet != core.Transaction_Result_SUCCESS {
				continue
			}

			txid := common.Bytes2Hex(v.Txid)
			// https://tronscan.org/#/transaction/fede1aa9e5c5d7bd179fd62e23bdd11e3c1edd0ca51e41070e34a026d6a42569

			for _, v1 := range v.Transaction.RawData.Contract {
				from := ""
				to := ""
				amount := int64(0)
				if v1.Type == core.Transaction_Contract_TransferContract && e.token == "" {
					// trx 转账
					unObj := &core.TransferContract{}
					err := proto.Unmarshal(v1.Parameter.GetValue(), unObj)
					if err != nil {
						continue
					}
					from = common.EncodeCheck(unObj.GetOwnerAddress())
					to = common.EncodeCheck(unObj.GetToAddress())
					amount = unObj.GetAmount()
					//fmt.Println(form, to, unObj.GetAmount())
				} else if v1.Type == core.Transaction_Contract_TriggerSmartContract && e.token != "" && e.tokenType == "trc20" {
					// 调用合约
					// trc20 转账
					//fmt.Println(v1.Parameter.GetValue())
					// 调用智能合约
					unObj := &core.TriggerSmartContract{}
					err := proto.Unmarshal(v1.Parameter.GetValue(), unObj)
					if err != nil {
						continue
					}
					contract := common.EncodeCheck(unObj.GetContractAddress())
					if contract != e.token {
						continue
					}
					data := unObj.GetData()
					// unObj.Data  https://goethereumbook.org/en/transfer-tokens/ 参考eth 操作
					flag := false
					to, amount, flag = processData(data)
					// 只有调用了 transfer(address,uint256) 才是转账
					if !flag {
						continue
					}
					from = common.EncodeCheck(unObj.GetOwnerAddress())
					//fmt.Println(contract, txid, from, to, amount)
				} else if v1.Type == core.Transaction_Contract_TransferAssetContract && e.token != "" && e.tokenType == "trc10" {
					// 通证转账合约
					// trc10 转账
					unObj := &core.TransferAssetContract{}
					err := proto.Unmarshal(v1.Parameter.GetValue(), unObj)
					if err != nil {
						continue
					}
					// contract := common.EncodeCheck(unObj.GetAssetName())
					from = common.EncodeCheck(unObj.GetOwnerAddress())
					to = common.EncodeCheck(unObj.GetToAddress())
				}

				transactions = append(transactions, types.Transaction{
					BlockNumber: big.NewInt(int64(num)),
					BlockHash:   common.EncodeCheck(block.Blockid),
					Hash:        txid,
					From:        from,
					To:          to,
					Value:       big.NewInt(amount),
				})
			}

		}
	}

	return transactions, toBlock, nil
}

// GetBalance 获取余额
func (e *TronWorker) GetBalance(address string) (*big.Int, error) {
	if e.token != "" && e.tokenType == "TRC20" {
		// 获取代币合约余额
		jsonString := "[{\"address\":\"" + address + "\"}]"
		data, err := e.http.CallContract(e.token, "balanceOf", jsonString)
		if err != nil {
			return nil, err
		}
		balance, err := tron.ToNumber(data)
		if err != nil {
			return nil, err
		}
		return balance, nil
	}
	// 获取主币、trc10
	balance, err := e.http.GetBalance(address, e.token)
	if err != nil {
		return nil, err
	}
	return big.NewInt(balance), nil
}

// CreateWallet 创建钱包
func (e *TronWorker) CreateWallet() (*types.Wallet, error) {
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

// Transfer 转账
func (e *TronWorker) Transfer(privateKeyStr string, toAddress string, value *big.Int, nonce uint64) (string, string, uint64, error) {

	from, err := e.GetAddressByPrivateKey(privateKeyStr)
	if err != nil {
		return "", "", 0, err
	}
	hash, err := e.http.Transfer(privateKeyStr, from, toAddress, e.token, value.Int64())
	if err != nil {
		return "", "", 0, err
	}

	return from, hash, nonce, nil
}

// GetAddressByPrivateKey 根据私钥获取地址
func (e *TronWorker) GetAddressByPrivateKey(privateKeyStr string) (string, error) {
	return tron.Private2TronAddress(privateKeyStr)
}

// processData 处理合约的调用参数
func processData(data []byte) (to string, amount int64, flag bool) {

	if len(data) >= 68 {
		if common.Bytes2Hex(data[:4]) != "a9059cbb" {
			return
		}
		// 多1位41
		data[15] = 65
		to = common.EncodeCheck(data[15:36])
		amount = new(big.Int).SetBytes(common.TrimLeftZeroes(data[36:68])).Int64()
		flag = false
	}
	return
}
