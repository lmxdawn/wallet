package types

import "math/big"

type Wallet struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

type Transaction struct {
	BlockNumber *big.Int // 区块号
	BlockHash   string   // 区块哈希
	Hash        string   // 交易hash
	From        string   // 交易者
	To          string   // 接收者
	Value       *big.Int // 交易数量
	Status      uint     // 状态（0：未完成，1：已完成）
}
