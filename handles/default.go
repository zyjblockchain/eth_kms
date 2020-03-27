package handles

import (
	"eth_kms/serializer"
	"github.com/gin-gonic/gin"
)

// Ping 用于心跳检测
func Ping(c *gin.Context) {
	serializer.SuccessResponse(c, nil, "pong")
}
