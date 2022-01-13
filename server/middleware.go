package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lmxdawn/wallet/engine"
)

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("x-token")
		if token != "" {
			c.Abort()
			APIResponse(c, ErrToken, nil)
		}

	}

}

// SetDB 设置db数据库
func SetDB(engines ...*engine.ConCurrentEngine) gin.HandlerFunc {

	return func(c *gin.Context) {
		for _, currentEngine := range engines {
			c.Set(currentEngine.Protocol, currentEngine)
		}
	}

}
