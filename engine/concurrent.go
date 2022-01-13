package engine

import (
	"github.com/lmxdawn/wallet/config"
	"github.com/lmxdawn/wallet/db"
	"github.com/lmxdawn/wallet/scheduler"
	"github.com/lmxdawn/wallet/types"
	"github.com/rs/zerolog/log"
	"time"
)

type Worker interface {
	getTransaction(uint64) ([]types.Transaction, error)
	getTransactionReceipt(*types.Transaction) error
	createWallet() (*types.Wallet, error)
	sendTransaction(string, string, int64, int) (string, error)
}

type Scheduler interface {
	BlockWorkerChan() chan uint64
	BlockWorkerReady(chan uint64)
	BlockSubmit(uint64)
	BlockRun()
	ReceiptWorkerChan() chan types.Transaction
	ReceiptWorkerReady(chan types.Transaction)
	ReceiptSubmit(types.Transaction)
	ReceiptRun()
}

type ConCurrentEngine struct {
	scheduler Scheduler
	Worker    Worker
	config    config.EngineConfig
	Protocol  string
	db        db.Database
}

// Run 启动
func (c *ConCurrentEngine) Run() {
	// 关闭连接
	defer c.db.Close()

	// 区块信息
	blockWorkerOut := make(chan types.Transaction)
	c.scheduler.BlockRun()
	// 交易信息
	c.scheduler.ReceiptRun()

	// 批量创建区块worker
	for i := uint64(0); i < c.config.BlockCount; i++ {
		c.createBlockWorker(blockWorkerOut)
	}

	// 批量创建交易worker
	for i := uint64(0); i < c.config.ReceiptCount; i++ {
		c.createReceiptWorker()
	}

	c.scheduler.BlockSubmit(c.config.BlockInit)

	for {
		transaction := <-blockWorkerOut
		//log.Info().Msgf("交易：%v", transaction)
		c.scheduler.ReceiptSubmit(transaction)
	}

}

// createBlockWorker 创建获取区块信息的工作
func (c *ConCurrentEngine) createBlockWorker(out chan types.Transaction) {
	in := c.scheduler.BlockWorkerChan()
	go func() {
		for {
			c.scheduler.BlockWorkerReady(in)
			num := <-in
			log.Info().Msgf("获取区块：%d", num)
			transactions, err := c.Worker.getTransaction(num)
			if err != nil {
				log.Info().Msgf("wait %d seconds, the latest block is not obtained", c.config.BlockAfterTime)
				<-time.After(time.Duration(c.config.BlockAfterTime) * time.Second)
				c.scheduler.BlockSubmit(num)
				continue
			}
			c.scheduler.BlockSubmit(num + 1)
			for _, transaction := range transactions {
				out <- transaction
			}
		}
	}()
}

// createReceiptWorker 创建获取区块信息的工作
func (c *ConCurrentEngine) createReceiptWorker() {
	in := c.scheduler.ReceiptWorkerChan()
	go func() {
		for {
			c.scheduler.ReceiptWorkerReady(in)
			transaction := <-in
			err := c.Worker.getTransactionReceipt(&transaction)
			if err != nil {
				log.Info().Msgf("wait %d seconds, the receipt information is invalid, err: %v", c.config.ReceiptAfterTime, err)
				<-time.After(time.Duration(c.config.ReceiptAfterTime) * time.Second)
				c.scheduler.ReceiptSubmit(transaction)
				continue
			}
			if transaction.Status == 1 {
				log.Info().Msgf("交易完成：%v", transaction.Hash)
			} else {
				log.Info().Msgf("交易失败：%v", transaction.Hash)
			}
			err = c.db.Put(transaction.To, "123")
			if err != nil {
				log.Info().Msgf("存入失败")
			}
			has, err := c.db.Has(transaction.To)
			if has && err == nil {
				log.Info().Msgf("查询到值")
			}
		}
	}()
}

// CreateWallet 创建钱包
func (c *ConCurrentEngine) CreateWallet() (string, error) {
	wallet, err := c.Worker.createWallet()
	if err != nil {
		return "", err
	}
	_ = c.db.Put(wallet.Address, wallet.PrivateKey)
	return wallet.Address, nil
}

// DeleteWallet 删除钱包
func (c *ConCurrentEngine) DeleteWallet(address string) error {
	err := c.db.Delete(address)
	if err != nil {
		return err
	}
	return nil
}

// SendTransaction 发送交易
func (c *ConCurrentEngine) SendTransaction(orderId string, toAddress string, value int64) (string, error) {
	hash, err := c.Worker.sendTransaction(c.config.WithdrawPrivateKey, toAddress, value, c.config.Decimals)
	if err != nil {
		return "", err
	}
	_ = c.db.Put(hash, orderId)
	return hash, nil
}

// GetTransactionReceipt 获取交易状态
func (c *ConCurrentEngine) GetTransactionReceipt(hash string) (int, error) {

	t := &types.Transaction{
		Hash:   hash,
		Status: 0,
	}

	err := c.Worker.getTransactionReceipt(t)
	if err != nil {
		return 0, err
	}
	if t.Status == 1 {

	}

	return int(t.Status), nil
}

// NewEngine 创建ETH
func NewEngine(config config.EngineConfig) (*ConCurrentEngine, error) {
	keyDB, err := db.NewKeyDB(config.File)
	if err != nil {
		return nil, err
	}

	var worker Worker
	switch config.Protocol {
	case "eth":
		worker = NewEthWorker(config.Confirms, config.Rpc)
	}

	return &ConCurrentEngine{
		//scheduler: scheduler.NewSimpleScheduler(), // 简单的任务调度器
		scheduler: scheduler.NewQueueScheduler(), // 队列的任务调度器
		Worker:    worker,
		config:    config,
		Protocol:  config.Protocol,
		db:        keyDB,
	}, nil
}
