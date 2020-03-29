package handles

import (
	"eth_kms/serializer"
	"eth_kms/services"
	"github.com/gin-gonic/gin"
)

type newKeyResult struct {
	Address string `json:"address"`
}

// NewEthKeyHandle 创建以太坊密钥对
func NewEthKeyHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service services.NewKeysInfo
		if err := c.ShouldBind(&service); err == nil {
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

// SaveKeysHandle 保存上传的秘钥对
func SaveKeysHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service services.SaveKeysInfo
		if err := c.ShouldBind(&service); err == nil {
			err = service.SaveKeys()
			if err != nil {
				serializer.ErrorResponse(c, 40001, "保存秘钥失败", err.Error())
			} else {
				serializer.SuccessResponse(c, nil, "秘钥保存成功")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}

type batchAddrResult struct {
	Addresses []string `json:"addresses"`
	Total     uint     `json:"total"`
}

// BatchGetAddrHandle 批量拉取地址
func BatchGetAddrHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service services.BatchGetAddrInfo
		if err := c.ShouldBind(&service); err == nil {
			addresses, total, err := service.BatchGetAddr()
			if err != nil {
				serializer.ErrorResponse(c, 40002, "批量拉取地址失败", err.Error())
			} else {
				serializer.SuccessResponse(c, batchAddrResult{Addresses: addresses, Total: total}, "批量拉取地址成功")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}

type signResult struct {
	Result string `json:"result"`
}

// SignDataHandle
func SignDataHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service services.SignInfo
		if err := c.ShouldBind(&service); err == nil {
			sig, err := service.KeySign()
			if err != nil {
				serializer.ErrorResponse(c, 40003, "签名失败", err.Error())
			} else {
				serializer.SuccessResponse(c, signResult{Result: sig}, "签名数据成功")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}
