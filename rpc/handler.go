package rpc

import "github.com/gin-gonic/gin"

// CreateWallet ...
func CreateWallet(c *gin.Context) {

	var q CreateWalletReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if q.Uid == ""{
		APIResponse(c, ErrParam, nil)
	}

	// 创建钱包

	APIResponse(c, nil, nil)
}

// Withdraw ...
func Withdraw(c *gin.Context) {

	var q WithdrawReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 生成交易，发送裸签名

	APIResponse(c, nil, nil)
}
