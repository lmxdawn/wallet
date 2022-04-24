package engine

import (
	"github.com/lmxdawn/wallet/btc"
	"github.com/lmxdawn/wallet/types"
	"math/big"
)

type BtcWorker struct {
	confirms uint64 // 需要的确认数
	http     *btc.Web
	network  string // 网络{MainNet：主网，TestNet：测试网，TestNet3：测试网3，SimNet：测试网}
}

func NewBtcWorker(confirms uint64, host, user, pass, network string) (*BtcWorker, error) {
	http, err := btc.NewWeb(host, user, pass, network)
	if err != nil {
		return nil, err
	}

	return &BtcWorker{
		confirms: confirms,
		http:     http,
	}, nil
}

// GetNowBlockNum 获取最新块
func (e *BtcWorker) GetNowBlockNum() (uint64, error) {
	blockNumber, err := e.http.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return uint64(blockNumber), nil
}

// GetTransactionReceipt 获取交易的票据
func (e *BtcWorker) GetTransactionReceipt(transaction *types.Transaction) error {

	// TODO 待实现

	return nil

}

// GetTransaction 获取交易信息
func (e *BtcWorker) GetTransaction(num uint64) ([]types.Transaction, uint64, error) {
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

	//numInt := int64(num)
	//toBlockInt := int64(toBlock)

	// TODO 待实现

	var transactions []types.Transaction

	return transactions, toBlock, nil
}

// GetBalance 获取余额
func (e *BtcWorker) GetBalance(address string) (*big.Int, error) {
	balance, err := e.http.GetBalance(address)
	if err != nil {
		return nil, err
	}
	return big.NewInt(balance), nil
}

// CreateWallet 创建钱包
func (e *BtcWorker) CreateWallet() (*types.Wallet, error) {

	btcW, err := e.http.CreateWallet()
	if err != nil {
		return nil, err
	}

	return &types.Wallet{
		Address:    btcW.NestedSegWitAddress, // 使用3开头的地址
		PublicKey:  btcW.PublicKey,
		PrivateKey: btcW.PrivateKey,
	}, err
}

// Transfer 转账
func (e *BtcWorker) Transfer(privateKeyStr string, toAddress string, value *big.Int, nonce uint64) (string, string, uint64, error) {

	from, err := e.GetAddressByPrivateKey(privateKeyStr)
	if err != nil {
		return "", "", 0, err
	}
	hash, err := e.http.Transfer(privateKeyStr, from, toAddress, value.Int64(), 0.001*1e8)
	if err != nil {
		return "", "", 0, err
	}

	return from, hash, nonce, nil
}

// GetAddressByPrivateKey 根据私钥获取地址
func (e *BtcWorker) GetAddressByPrivateKey(privateKeyStr string) (string, error) {
	return e.http.GetWalletByPrivateKey(privateKeyStr)
}
