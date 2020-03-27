package middleware

import (
	"eth_kms/models"
	"eth_kms/serializer"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, exist := c.Get("user"); exist {
			if _, ok := user.(*models.User); ok {
				c.Next()
				return
			}
		}

		// 需要先登录
		serializer.ErrorResponse(c, 40001, "需要先登录", "")
		// 终止后面的handle执行
		c.Abort()
	}
}
