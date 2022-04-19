package server

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/lmxdawn/wallet/engine"
	"math/big"
)

// CreateWallet ...
// @Tags 钱包
// @Summary 创建钱包地址
// @Produce json
// @Security ApiKeyAuth
// @Param login body CreateWalletReq true "参数"
// @Success 200 {object} Response{data=server.CreateWalletRes}
// @Router /api/createWallet [post]
func CreateWallet(c *gin.Context) {

	var q CreateWalletReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	v, ok := c.Get(q.Protocol + q.CoinName)
	if !ok {
		APIResponse(c, ErrEngine, nil)
		return
	}

	currentEngine := v.(*engine.ConCurrentEngine)

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
// @Tags 钱包
// @Summary 删除钱包地址
// @Produce json
// @Security ApiKeyAuth
// @Param login body DelWalletReq true "参数"
// @Success 200 {object} Response
// @Router /api/delWallet [post]
func DelWallet(c *gin.Context) {

	var q DelWalletReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	v, ok := c.Get(q.Protocol + q.CoinName)
	if !ok {
		APIResponse(c, ErrEngine, nil)
		return
	}

	currentEngine := v.(*engine.ConCurrentEngine)

	q.Address = common.HexToAddress(q.Address).Hex()

	err := currentEngine.DeleteWallet(q.Address)
	if err != nil {
		APIResponse(c, ErrCreateWallet, nil)
		return
	}

	APIResponse(c, nil, nil)
}

// Withdraw ...
// @Tags 钱包
// @Summary 提现
// @Produce json
// @Security ApiKeyAuth
// @Param login body WithdrawReq true "参数"
// @Success 200 {object} Response{data=server.WithdrawRes}
// @Router /api/withdraw [post]
func Withdraw(c *gin.Context) {

	var q WithdrawReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	v, ok := c.Get(q.Protocol + q.CoinName)
	if !ok {
		APIResponse(c, ErrEngine, nil)
		return
	}

	currentEngine := v.(*engine.ConCurrentEngine)

	q.Address = common.HexToAddress(q.Address).Hex()

	hash, err := currentEngine.Withdraw(q.OrderId, q.Address, q.Value)
	if err != nil {
		APIResponse(c, err, nil)
		return
	}

	res := WithdrawRes{Hash: hash}

	APIResponse(c, nil, res)
}

// Collection ...
// @Tags 归集某个地址
// @Summary 归集
// @Produce json
// @Security ApiKeyAuth
// @Param login body CollectionReq true "参数"
// @Success 200 {object} Response{data=server.CollectionRes}
// @Router /api/collection [post]
func Collection(c *gin.Context) {

	var q CollectionReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	v, ok := c.Get(q.Protocol + q.CoinName)
	if !ok {
		APIResponse(c, ErrEngine, nil)
		return
	}

	currentEngine := v.(*engine.ConCurrentEngine)

	n := new(big.Int)
	max, ok := n.SetString(q.Max, 10)
	if !ok {
		APIResponse(c, InternalServerError, nil)
		return
	}

	q.Address = common.HexToAddress(q.Address).Hex()

	balance, err := currentEngine.Collection(q.Address, max)
	if err != nil {
		APIResponse(c, err, nil)
		return
	}

	res := CollectionRes{Balance: balance.String()}

	APIResponse(c, nil, res)
}

// GetTransactionReceipt ...
// @Tags 钱包
// @Summary 获取交易结果
// @Produce json
// @Security ApiKeyAuth
// @Param login body TransactionReceiptReq true "参数"
// @Success 200 {object} Response{data=server.TransactionReceiptRes}
// @Router /api/getTransactionReceipt [get]
func GetTransactionReceipt(c *gin.Context) {

	var q TransactionReceiptReq

	if err := c.ShouldBindJSON(&q); err != nil {
		HandleValidatorError(c, err)
		return
	}

	v, ok := c.Get(q.Protocol + q.CoinName)
	if !ok {
		APIResponse(c, ErrEngine, nil)
		return
	}

	currentEngine := v.(*engine.ConCurrentEngine)

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
