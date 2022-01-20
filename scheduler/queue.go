package scheduler

import (
	"github.com/lmxdawn/wallet/db"
	"github.com/lmxdawn/wallet/types"
)

// QueueScheduler 队列的调度器
type QueueScheduler struct {
	blockNum             chan uint64      // 区块的通道
	blockNumWorker       chan chan uint64 // 区块worker的通道
	receipt              chan types.Transaction
	receiptWorker        chan chan types.Transaction
	collectionSend       chan db.WalletItem
	collectionSendWorker chan chan db.WalletItem
}

func NewQueueScheduler() *QueueScheduler {
	return &QueueScheduler{}
}

func (q *QueueScheduler) BlockWorkerChan() chan uint64 {
	return make(chan uint64)
}

func (q *QueueScheduler) BlockWorkerReady(w chan uint64) {
	q.blockNumWorker <- w
}

func (q *QueueScheduler) BlockSubmit(blockNum uint64) {
	q.blockNum <- blockNum
}

func (q *QueueScheduler) BlockRun() {
	q.blockNum = make(chan uint64)
	q.blockNumWorker = make(chan chan uint64)
	go func() {
		var nQ []uint64
		var nWQ []chan uint64
		for {
			var activateN uint64
			var activateNW chan uint64
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

func (q *QueueScheduler) CollectionSendWorkerChan() chan db.WalletItem {
	return make(chan db.WalletItem)
}

func (q *QueueScheduler) CollectionSendWorkerReady(c chan db.WalletItem) {
	q.collectionSendWorker <- c
}

func (q *QueueScheduler) CollectionSendSubmit(c db.WalletItem) {
	q.collectionSend <- c
}

func (q *QueueScheduler) CollectionSendRun() {
	q.collectionSend = make(chan db.WalletItem)
	q.collectionSendWorker = make(chan chan db.WalletItem)
	go func() {
		var cQ []db.WalletItem
		var cWQ []chan db.WalletItem
		for {
			var activateR db.WalletItem
			var activateRW chan db.WalletItem
			if len(cQ) > 0 && len(cWQ) > 0 {
				activateR = cQ[0]
				activateRW = cWQ[0]
			}
			select {
			case c := <-q.collectionSend:
				cQ = append(cQ, c)
			case cw := <-q.collectionSendWorker:
				cWQ = append(cWQ, cw)
			case activateRW <- activateR:
				cQ = cQ[1:]
				cWQ = cWQ[1:]
			}
		}

	}()
}
