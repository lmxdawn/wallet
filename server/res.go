package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response ...
type Response struct {
	Code    int         `json:"code"`    // 错误code码
	Message string      `json:"message"` // 错误信息
	Data    interface{} `json:"data"`    // 成功时返回的对象
}

// APIResponse ....
func APIResponse(Ctx *gin.Context, err error, data interface{}) {
	if err == nil {
		err = OK
	}
	codeNum, message := DecodeErr(err)
	Ctx.JSON(http.StatusOK, Response{
		Code:    codeNum,
		Message: message,
		Data:    data,
	})
}

type CreateWalletRes struct {
	Address string `json:"address"` // 生成的钱包地址
}

type WithdrawRes struct {
	Hash string `json:"hash"` // 生成的交易hash
}

type TransactionReceiptRes struct {
	Status int `json:"status"` // 交易状态（0：未成功，1：已成功）
}
