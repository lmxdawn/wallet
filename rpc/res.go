package rpc

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
