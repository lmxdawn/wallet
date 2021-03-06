package engine

import (
	"github.com/lmxdawn/wallet/client"
	"github.com/lmxdawn/wallet/config"
	"github.com/lmxdawn/wallet/db"
	"github.com/lmxdawn/wallet/scheduler"
	"github.com/lmxdawn/wallet/types"
	"github.com/rs/zerolog/log"
	"math/big"
	"strconv"
	"time"
)

type Worker interface {
	GetNowBlockNum() (uint64, error)
	GetTransaction(num uint64) ([]types.Transaction, uint64, error)
	GetTransactionReceipt(*types.Transaction) error
	GetBalance(address string) (*big.Int, error)
	CreateWallet() (*types.Wallet, error)
	Transfer(privateKeyStr string, toAddress string, value *big.Int, nonce uint64) (string, string, uint64, error)
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
	CollectionSendWorkerChan() chan db.WalletItem
	CollectionSendWorkerReady(c chan db.WalletItem)
	CollectionSendSubmit(c db.WalletItem)
	CollectionSendRun()
}

type ConCurrentEngine struct {
	scheduler Scheduler
	Worker    Worker
	config    config.EngineConfig
	Protocol  string
	CoinName  string
	db        db.Database
	http      *client.HttpClient
}

// Run 启动
func (c *ConCurrentEngine) Run() {
	// 关闭连接
	defer c.db.Close()

	go c.blockLoop()
	go c.collectionLoop()

	select {}
}

// blockLoop 区块循环监听
func (c *ConCurrentEngine) blockLoop() {
	// 读取当前区块
	blockNumber := c.config.BlockInit
	blockNumberStr, err := c.db.Get("block_number")
	if err == nil && blockNumberStr != "" {
		blockNumber, _ = strconv.ParseUint(blockNumberStr, 10, 64)
	}

	// 区块信息
	c.scheduler.BlockRun()
	// 交易信息
	c.scheduler.ReceiptRun()
	// 归集信息
	c.scheduler.CollectionSendRun()

	// 批量创建区块worker
	blockWorkerOut := make(chan types.Transaction)
	c.createBlockWorker(blockWorkerOut)

	// 批量创建交易worker
	for i := uint64(0); i < c.config.ReceiptCount; i++ {
		c.createReceiptWorker()
	}

	c.scheduler.BlockSubmit(blockNumber)

	go func() {
		for {
			transaction := <-blockWorkerOut
			//log.Info().Msgf("交易：%v", transaction)
			c.scheduler.ReceiptSubmit(transaction)
		}
	}()
}

// collectionLoop 归集循环监听
func (c *ConCurrentEngine) collectionLoop() {
	n := new(big.Int)
	collectionMax, ok := n.SetString(c.config.CollectionMax, 10)
	if !ok {
		panic("setString: error")
	}
	// 配置大于0才去自动归集
	if collectionMax.Cmp(big.NewInt(0)) > 0 {
		// 启动归集
		collectionWorkerOut := make(chan db.WalletItem)
		c.createCollectionWorker(collectionWorkerOut)

		// 启动归集发送worker
		for i := uint64(0); i < c.config.CollectionCount; i++ {
			c.createCollectionSendWorker(collectionMax)
		}

		go func() {
			for {
				collectionSend := <-collectionWorkerOut
				c.scheduler.CollectionSendSubmit(collectionSend)
			}
		}()
	}
}

// createBlockWorker 创建获取区块信息的工作
func (c *ConCurrentEngine) createBlockWorker(out chan types.Transaction) {
	in := c.scheduler.BlockWorkerChan()
	go func() {
		for {
			c.scheduler.BlockWorkerReady(in)
			num := <-in
			log.Info().Msgf("%v，监听区块：%d", c.config.CoinName, num)
			transactions, blockNum, err := c.Worker.GetTransaction(num)
			if err != nil || blockNum == num {
				log.Info().Msgf("等待%d秒，当前已是最新区块", c.config.BlockAfterTime)
				<-time.After(time.Duration(c.config.BlockAfterTime) * time.Second)
				c.scheduler.BlockSubmit(num)
				continue
			}
			err = c.db.Put("block_number", strconv.FormatUint(blockNum, 10))
			if err != nil {
				c.scheduler.BlockSubmit(num)
			} else {
				c.scheduler.BlockSubmit(blockNum)
			}
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
			err := c.Worker.GetTransactionReceipt(&transaction)
			if err != nil {
				log.Info().Msgf("等待%d秒，收据信息无效, err: %v", c.config.ReceiptAfterTime, err)
				<-time.After(time.Duration(c.config.ReceiptAfterTime) * time.Second)
				c.scheduler.ReceiptSubmit(transaction)
				continue
			}
			if transaction.Status != 1 {
				log.Error().Msgf("交易失败：%v", transaction.Hash)
				continue
			}
			//log.Info().Msgf("交易完成：%v", transaction.Hash)

			// 判断是否存在
			if ok, err := c.db.Has(c.config.HashPrefix + transaction.Hash); err == nil && ok {
				log.Info().Msgf("监听到提现事件，提现地址：%v，提现哈希：%v", transaction.To, transaction.Hash)
				orderId, err := c.db.Get(c.config.HashPrefix + transaction.Hash)
				if err != nil {
					log.Error().Msgf("未查询到订单：%v, %v", transaction.Hash, err)
					// 重新提交
					c.scheduler.ReceiptSubmit(transaction)
					continue
				}
				err = c.http.WithdrawSuccess(transaction.Hash, transaction.Status, orderId, transaction.To, transaction.Value.Int64())
				if err != nil {
					log.Error().Msgf("提现回调通知失败：%v, %v", transaction.Hash, err)
					// 重新提交
					c.scheduler.ReceiptSubmit(transaction)
					continue
				}
				_ = c.db.Delete(transaction.Hash)
			} else if ok, err := c.db.Has(c.config.WalletPrefix + transaction.To); err == nil && ok {
				log.Info().Msgf("监听到充值事件，充值地址：%v，充值哈希：%v", transaction.To, transaction.Hash)
				err = c.http.RechargeSuccess(transaction.Hash, transaction.Status, transaction.To, transaction.Value.Int64())
				if err != nil {
					log.Error().Msgf("充值回调通知失败：%v, %v", transaction.Hash, err)
					// 重新提交
					c.scheduler.ReceiptSubmit(transaction)
					continue
				}
			}
		}
	}()
}

// createCollectionWorker 创建归集Worker
func (c *ConCurrentEngine) createCollectionWorker(out chan db.WalletItem) {
	go func() {
		for {
			<-time.After(time.Duration(c.config.CollectionAfterTime) * time.Second)
			list, err := c.db.ListWallet(c.config.WalletPrefix)
			if err != nil {
				continue
			}
			for _, item := range list {
				out <- item
			}
		}
	}()
}

// collectionSendWorker 创建归集发送交易的worker
func (c *ConCurrentEngine) createCollectionSendWorker(max *big.Int) {
	in := c.scheduler.CollectionSendWorkerChan()
	go func() {
		for {
			c.scheduler.CollectionSendWorkerReady(in)
			collectionSend := <-in
			_, err := c.collection(collectionSend.Address, collectionSend.PrivateKey, max)
			if err != nil {
				// 归集失败，重新加入归集队列
				c.scheduler.CollectionSendSubmit(collectionSend)
				continue
			}
		}
	}()
}

// 归集
func (c ConCurrentEngine) collection(address, privateKey string, max *big.Int) (*big.Int, error) {
	balance, err := c.Worker.GetBalance(address)
	if err != nil {
		return nil, err
	}
	if balance.Cmp(max) < 0 {
		return big.NewInt(0), nil
	}

	// 开始归集
	_, _, _, err = c.Worker.Transfer(privateKey, c.config.CollectionAddress, balance, 0)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// Collection 归集某个地址
func (c *ConCurrentEngine) Collection(address string, max *big.Int) (*big.Int, error) {

	// 查询地址是否存在
	privateKey, err := c.db.Get(c.config.WalletPrefix + address)
	if err != nil {
		return nil, err
	}

	balance, err := c.collection(address, privateKey, max)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// CreateWallet 创建钱包
func (c *ConCurrentEngine) CreateWallet() (string, error) {
	wallet, err := c.Worker.CreateWallet()
	if err != nil {
		return "", err
	}
	_ = c.db.Put(c.config.WalletPrefix+wallet.Address, wallet.PrivateKey)
	log.Info().Msgf("创建钱包成功，地址：%v，私钥：%v", wallet.Address, wallet.PrivateKey)
	return wallet.Address, nil
}

// DeleteWallet 删除钱包
func (c *ConCurrentEngine) DeleteWallet(address string) error {
	err := c.db.Delete(c.config.WalletPrefix + address)
	if err != nil {
		return err
	}
	return nil
}

// Withdraw 提现
func (c *ConCurrentEngine) Withdraw(orderId string, toAddress string, value int64) (string, error) {

	_, hash, _, err := c.Worker.Transfer(c.config.WithdrawPrivateKey, toAddress, big.NewInt(value), 0)
	if err != nil {
		return "", err
	}
	_ = c.db.Put(c.config.HashPrefix+hash, orderId)
	return hash, nil
}

// GetTransactionReceipt 获取交易状态
func (c *ConCurrentEngine) GetTransactionReceipt(hash string) (int, error) {

	t := &types.Transaction{
		Hash:   hash,
		Status: 0,
	}

	err := c.Worker.GetTransactionReceipt(t)
	if err != nil {
		return 0, err
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
		worker, err = NewEthWorker(config.Confirms, config.Contract, config.Rpc)
		if err != nil {
			return nil, err
		}
	case "trx":
		worker, err = NewTronWorker(config.Confirms, config.Contract, config.ContractType, config.Rpc)
		if err != nil {
			return nil, err
		}
	case "btc":
		worker, err = NewBtcWorker(config.Confirms, config.Rpc, config.User, config.Pass, config.Rpc)
		if err != nil {
			return nil, err
		}
	}

	http := client.NewHttpClient(config.Protocol, config.CoinName, config.RechargeNotifyUrl, config.WithdrawNotifyUrl)

	return &ConCurrentEngine{
		//scheduler: scheduler.NewSimpleScheduler(), // 简单的任务调度器
		scheduler: scheduler.NewQueueScheduler(), // 队列的任务调度器
		Worker:    worker,
		config:    config,
		Protocol:  config.Protocol,
		CoinName:  config.CoinName,
		db:        keyDB,
		http:      http,
	}, nil
}
