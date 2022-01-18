package server

type CreateWalletReq struct {
	Protocol string `json:"protocol" binding:"required"` // 协议
	CoinName string `json:"coinName" binding:"required"` // 币种名称
}

type DelWalletReq struct {
	Protocol string `json:"protocol" binding:"required"` // 协议
	CoinName string `json:"coinName" binding:"required"` // 币种名称
	Address  string `json:"address" binding:"required"`  // 地址
}

type WithdrawReq struct {
	Protocol string `json:"protocol" binding:"required"` // 协议
	CoinName string `json:"coinName" binding:"required"` // 币种名称
	OrderId  string `json:"orderId" binding:"required"`  // 订单号
	Address  string `json:"address" binding:"required"`  // 提现地址
	Value    int64  `json:"value" binding:"required"`    // 金额
}

type TransactionReceiptReq struct {
	Protocol string `json:"protocol" binding:"required"` // 协议
	CoinName string `json:"coinName" binding:"required"` // 币种名称
	Hash     string `json:"hash" binding:"required"`     // 交易哈希
}
