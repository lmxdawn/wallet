package engine

import (
	"github.com/lmxdawn/wallet/db"
	"github.com/lmxdawn/wallet/scheduler"
	"github.com/lmxdawn/wallet/types"
	"github.com/rs/zerolog/log"
	"time"
)

type Worker interface {
	getTransaction(int64) ([]types.Transaction, error)
	getTransactionReceipt(*types.Transaction) error
}

type Scheduler interface {
	BlockWorkerChan() chan int64
	BlockWorkerReady(chan int64)
	BlockSubmit(int64)
	BlockRun()
	ReceiptWorkerChan() chan types.Transaction
	ReceiptWorkerReady(chan types.Transaction)
	ReceiptSubmit(types.Transaction)
	ReceiptRun()
}

type ConCurrentEngine struct {
	scheduler          Scheduler
	worker             Worker
	db                 db.Database
	blockWorkerCount   int // 区块worker数量
	blockAfterTime     int // 获取最新块的等待时间
	receiptWorkerCount int // 交易凭证worker数量
	receiptAfterTime   int // 获取交易信息的等待时间
}

// Run 启动
func (c *ConCurrentEngine) Run(blockNumber int64) {
	// 关闭连接
	defer c.db.Close()

	// 区块信息
	blockWorkerOut := make(chan types.Transaction)
	c.scheduler.BlockRun()
	// 交易信息
	c.scheduler.ReceiptRun()

	// 批量创建区块worker
	for i := 0; i < c.blockWorkerCount; i++ {
		c.createBlockWorker(blockWorkerOut)
	}

	// 批量创建交易worker
	for i := 0; i < c.receiptWorkerCount; i++ {
		c.createReceiptWorker()
	}

	c.scheduler.BlockSubmit(blockNumber)

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
			transactions, err := c.worker.getTransaction(num)
			if err != nil {
				log.Info().Msgf("wait %d seconds, the latest block is not obtained", c.blockAfterTime)
				<-time.After(time.Duration(c.blockAfterTime) * time.Second)
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
			err := c.worker.getTransactionReceipt(&transaction)
			if err != nil {
				log.Info().Msgf("wait %d seconds, the receipt information is invalid, err: %v", c.receiptAfterTime, err)
				<-time.After(time.Duration(c.receiptAfterTime) * time.Second)
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

// NewEthEngine 创建ETH
func NewEthEngine(blockCount int, receiptCount int, confirms uint64, url string, keyDB *db.KeyDB) *ConCurrentEngine {
	return &ConCurrentEngine{
		//scheduler: scheduler.NewSimpleScheduler(), // 简单的任务调度器
		scheduler:          scheduler.NewQueueScheduler(), // 队列的任务调度器
		worker:             NewEthWorker(confirms, url),
		db:                 keyDB,
		blockWorkerCount:   blockCount,
		blockAfterTime:     3,
		receiptWorkerCount: receiptCount,
		receiptAfterTime:   2,
	}
}
