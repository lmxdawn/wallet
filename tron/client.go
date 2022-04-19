package tron

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"strings"
	"time"
)

type Client struct {
	node string
	rpc  *client.GrpcClient
}

func NewClient(node string, withTLS bool) (*Client, error) {
	opts := make([]grpc.DialOption, 0)
	if withTLS {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	c := new(Client)
	c.node = node
	c.rpc = client.NewGrpcClient(node)
	err := c.rpc.Start(opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc client start error: %v", err)
	}
	return c, nil
}

// SetTimeout 设置超时时间
func (c *Client) SetTimeout(timeout time.Duration) error {
	if c == nil {
		return errors.New("client is nil ptr")
	}
	c.rpc = client.NewGrpcClientWithTimeout(c.node, timeout)
	err := c.rpc.Start()
	if err != nil {
		return fmt.Errorf("grpc start error: %v", err)
	}
	return nil
}

// keepConnect
func (c *Client) keepConnect() error {
	_, err := c.rpc.GetNodeInfo()
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return c.rpc.Reconnect(c.node)
		}
		return fmt.Errorf("node connect error: %v", err)
	}
	return nil
}

// Transfer 转账
func (c *Client) Transfer(from, to string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.rpc.Transfer(from, to, amount)
}

// TransferTrc10 trc10 转账
func (c *Client) TransferTrc10(from, to, assetName string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.rpc.TransferAsset(from, to, assetName, amount)
}

// TriggerContract 执行智能合约的方法
func (c *Client) TriggerContract(from, contractAddress, method, jsonString string, constant bool, feeLimit int64) (*api.TransactionExtention, error) {

	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	result := &api.TransactionExtention{}
	if constant {
		result, err = c.rpc.TriggerConstantContract(from, contractAddress, method, jsonString)
		if err != nil {
			return nil, err
		}
	} else {
		result, err = c.rpc.TriggerContract(from, contractAddress, method, jsonString, feeLimit, 0, "", 0)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// BroadcastTransaction 广播交易签名
func (c *Client) BroadcastTransaction(transaction *core.Transaction) error {
	err := c.keepConnect()
	if err != nil {
		return err
	}
	result, err := c.rpc.Broadcast(transaction)
	if err != nil {
		return fmt.Errorf("broadcast transaction error: %v", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("bad transaction: %v", string(result.GetMessage()))
	}
	if result.Result == true {
		return nil
	}
	d, _ := json.Marshal(result)
	return fmt.Errorf("tx send fail: %s", string(d))
}

// GetTransactionInfoByID 查询交易是否成功
func (c *Client) GetTransactionInfoByID(txid string, isRes bool) (*core.TransactionInfo, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}

	if !isRes {
		return c.rpc.GetTransactionInfoByID(txid)
	}

	count := 0
	for {
		if count >= 100 {
			return nil, errors.New("获取超时")
		}
		r, err := c.rpc.GetTransactionInfoByID(txid)
		if err != nil {
			count++
			<-time.After(3 * time.Second)
		} else {
			//dd, _ := json.Marshal(r)
			return r, nil
		}
	}
}

// GetBlock 获取区块信息
func (c *Client) GetBlock(start, end int64) (*api.BlockListExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}

	blocks, err := c.rpc.GetBlockByLimitNext(start, end)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}

// GetBlockByNum 获取区块信息
func (c *Client) GetBlockByNum(num int64) (*api.BlockExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}

	block, err := c.rpc.GetBlockByNum(num)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetNowBlock 获取最新区块
func (c *Client) GetNowBlock() (*api.BlockExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}

	block, err := c.rpc.GetNowBlock()
	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetAccount 获取账户余额
func (c *Client) GetAccount(address string) (*core.Account, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}

	return c.rpc.GetAccount(address)
}
