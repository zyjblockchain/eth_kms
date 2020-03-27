package handles

import (
	"eth_kms/serializer"
	"eth_kms/services"
	"github.com/gin-gonic/gin"
)

//go:generate gencodec -type newRsp   -out new_rsp_json.go
type newKeyResult struct {
	Address string `json:"address"`
}

// NewEthKey 创建以太坊密钥对
func NewEthKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service services.NewKeysInfo
		if err := c.ShouldBind(&service); err != nil {
			address, err := service.NewKeys()
			if err != nil {
				serializer.ErrorResponse(c, 40001, "创建以太坊密钥对失败", err.Error())
			} else {
				serializer.SuccessResponse(c, newKeyResult{Address: address}, "创建以太坊秘钥对成功")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}
