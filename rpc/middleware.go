package rpc

import (
	"github.com/gin-gonic/gin"
	"github.com/lmxdawn/wallet/db"
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
func SetDB(dbs ...*db.KeyDB) gin.HandlerFunc {

	return func(c *gin.Context) {
		for _, keyDB := range dbs {
			c.Set(keyDB.Protocol, keyDB)
		}
	}

}
