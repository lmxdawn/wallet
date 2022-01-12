package scheduler

import (
	"github.com/lmxdawn/wallet/types"
)

// QueueScheduler 队列的调度器
type QueueScheduler struct {
	blockNum       chan int64      // 区块的通道
	blockNumWorker chan chan int64 // 区块worker的通道
	receipt        chan types.Transaction
	receiptWorker  chan chan types.Transaction
}

func NewQueueScheduler() *QueueScheduler {
	return &QueueScheduler{}
}

func (q *QueueScheduler) BlockWorkerChan() chan int64 {
	return make(chan int64)
}

func (q *QueueScheduler) BlockWorkerReady(w chan int64) {
	q.blockNumWorker <- w
}

func (q *QueueScheduler) BlockSubmit(blockNum int64) {
	q.blockNum <- blockNum
}

func (q *QueueScheduler) BlockRun() {
	q.blockNum = make(chan int64)
	q.blockNumWorker = make(chan chan int64)
	go func() {
		var nQ []int64
		var nWQ []chan int64
		for {
			var activateN int64
			var activateNW chan int64
			if len(nQ) > 0 && len(nWQ) > 0 {
				activateN = nQ[0]
				activateNW = nWQ[0]
			}
			select {
			case n := <-q.blockNum:
				nQ = append(nQ, n)
			case nw := <-q.blockNumWorker:
				nWQ = append(nWQ, nw)
			case activateNW <- activateN:
				nQ = nQ[1:]
				nWQ = nWQ[1:]
			}
		}

	}()
}

func (q *QueueScheduler) ReceiptWorkerChan() chan types.Transaction {
	return make(chan types.Transaction)
}

func (q *QueueScheduler) ReceiptWorkerReady(t chan types.Transaction) {
	q.receiptWorker <- t
}

func (q *QueueScheduler) ReceiptSubmit(transaction types.Transaction) {
	q.receipt <- transaction
}

func (q *QueueScheduler) ReceiptRun() {
	q.receipt = make(chan types.Transaction)
	q.receiptWorker = make(chan chan types.Transaction)
	go func() {
		var rQ []types.Transaction
		var rWQ []chan types.Transaction
		for {
			var activateR types.Transaction
			var activateRW chan types.Transaction
			if len(rQ) > 0 && len(rWQ) > 0 {
				activateR = rQ[0]
				activateRW = rWQ[0]
			}
			select {
			case r := <-q.receipt:
				rQ = append(rQ, r)
			case rw := <-q.receiptWorker:
				rWQ = append(rWQ, rw)
			case activateRW <- activateR:
				rQ = rQ[1:]
				rWQ = rWQ[1:]
			}
		}

	}()
}
