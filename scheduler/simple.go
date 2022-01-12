package scheduler

import (
	"github.com/lmxdawn/wallet/types"
)

// SimpleScheduler 简单的调度器
type SimpleScheduler struct {
	blockNum chan int64             // 区块的通道
	receipt  chan types.Transaction // 交易的通道
}

func NewSimpleScheduler() *SimpleScheduler {
	return &SimpleScheduler{}
}

func (s *SimpleScheduler) BlockWorkerChan() chan int64 {
	return s.blockNum
}

func (s *SimpleScheduler) BlockWorkerReady(chan int64) {
}

func (s *SimpleScheduler) BlockSubmit(n int64) {
	go func() {
		s.blockNum <- n
	}()
}

func (s *SimpleScheduler) BlockRun() {
	s.blockNum = make(chan int64)
}

func (s *SimpleScheduler) ReceiptWorkerChan() chan types.Transaction {
	return s.receipt
}

func (s *SimpleScheduler) ReceiptWorkerReady(chan types.Transaction) {
}

func (s *SimpleScheduler) ReceiptSubmit(t types.Transaction) {
	go func() {
		s.receipt <- t
	}()
}

func (s *SimpleScheduler) ReceiptRun() {
	s.receipt = make(chan types.Transaction)
}
