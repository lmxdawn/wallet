package btc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Client struct {
	rpc *rpcclient.Client
}

func NewClient(host, user, pass string) (*Client, error) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         host,
		User:         user,
		Pass:         pass,
	}, nil)

	if err != nil {
		return nil, err
	}

	return &Client{
		rpc: client,
	}, nil
}

// GetBlockCount 获取区块数量，调用返回本地最优链中的区块数量。
func (c *Client) GetBlockCount() (int64, error) {
	return c.rpc.GetBlockCount()
}

// GetBlockHash 获取指定高度区块的哈希。
func (c *Client) GetBlockHash(blockNumber int64) (*chainhash.Hash, error) {
	return c.rpc.GetBlockHash(blockNumber)
}

// ImportAddress 导入地址
func (c *Client) ImportAddress(address string) error {
	return c.rpc.ImportAddress(address)
}

// ListUnspent 获取未交易的UTX
func (c *Client) ListUnspent(address btcutil.Address) (listUnSpent []btcjson.ListUnspentResult, err error) {
	adds := [1]btcutil.Address{address}
	listUnSpent, err = c.rpc.ListUnspentMinMaxAddresses(1, 999999, adds[:])
	if err != nil {
		return
	}
	return
}

// SendRawTransaction 发送裸交易
func (c *Client) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {

	return c.rpc.SendRawTransaction(tx, false)

}

// ListSinceBlock 查询指定区块后发生的钱包交易
func (c *Client) ListSinceBlock(blockHash *chainhash.Hash) (*btcjson.ListSinceBlockResult, error) {
	return c.rpc.ListSinceBlock(blockHash)
}

// ListTransactions 查询最近发生的钱包交易
func (c *Client) ListTransactions(account string) ([]btcjson.ListTransactionsResult, error) {
	return c.rpc.ListTransactions(account)
}
