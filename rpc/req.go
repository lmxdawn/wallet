package rpc

type CreateWalletReq struct {
	Protocol string `json:"protocol" binding:"required"` // 协议
	Uid      string `json:"uid" binding:"required"`      // 用户ID
}

type WithdrawReq struct {
	Protocol string `json:"protocol" binding:"required"` // 协议
	OrderId  string `json:"orderId" binding:"required"`  // 订单号
	Address  string `json:"address" binding:"required"`  // 提现地址
	Value    string `json:"value" binding:"required"`    // 金额
}
