package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lmxdawn/wallet/engine"
)

// CreateWallet ...
func CreateWallet(c *gin.Context) {

	var q CreateWalletReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	currentEngine := c.MustGet(q.Protocol).(engine.ConCurrentEngine)

	// 创建钱包
	address, err := currentEngine.CreateWallet()
	if err != nil {
		APIResponse(c, ErrCreateWallet, nil)
		return
	}

	res := CreateWalletRes{Address: address}

	APIResponse(c, nil, res)
}

// DelWallet ...
func DelWallet(c *gin.Context) {

	var q DelWalletReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	currentEngine := c.MustGet(q.Protocol).(engine.ConCurrentEngine)

	err := currentEngine.DeleteWallet(q.Address)
	if err != nil {
		APIResponse(c, ErrCreateWallet, nil)
		return
	}

	APIResponse(c, nil, nil)
}

// Withdraw ...
func Withdraw(c *gin.Context) {

	var q WithdrawReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	currentEngine := c.MustGet(q.Protocol).(engine.ConCurrentEngine)

	hash, err := currentEngine.SendTransaction(q.OrderId, q.Address, q.Value)
	if err != nil {
		APIResponse(c, nil, nil)
		return
	}

	res := WithdrawRes{Hash: hash}

	APIResponse(c, nil, res)
}

// GetTransactionReceipt ...
func GetTransactionReceipt(c *gin.Context) {

	var q TransactionReceiptReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	currentEngine := c.MustGet(q.Protocol).(engine.ConCurrentEngine)

	status, err := currentEngine.GetTransactionReceipt(q.Hash)
	if err != nil {
		APIResponse(c, InternalServerError, nil)
		return
	}

	res := TransactionReceiptRes{
		Status: status,
	}

	APIResponse(c, nil, res)
}
