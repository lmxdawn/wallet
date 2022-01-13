package server

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// 定义一个全局翻译器T
var trans ut.Translator

// HandleValidatorError 处理字段校验异常
func HandleValidatorError(c *gin.Context, err error) {
	//如何返回错误信息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		APIResponse(c, InternalServerError, nil)
	}
	APIResponse(c, ErrParam, firstErr(errs.Translate(trans)))
	return
}

//   firstErr 返回第一个错误
func firstErr(filedMap map[string]string) string {
	for _, err := range filedMap {
		return err
	}
	return ""
}
