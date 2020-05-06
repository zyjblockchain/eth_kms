package handles

import (
	"github.com/gin-gonic/gin"
	"github.com/zyjblockchain/eth_kms/serializer"
	"github.com/zyjblockchain/eth_kms/services"
)

// SetCollectionAddrHandle
func SetCollectionAddrHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service services.SetCollectionAddr
		if err := c.ShouldBind(&service); err == nil {
			if err := service.SetCollectionAddr(); err != nil {
				serializer.ErrorResponse(c, 40008, "设置资产归集接收地址调用失败", err.Error())
			} else {
				serializer.SuccessResponse(c, nil, "设置归集接收地址成功")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}

// GetCollectionAddrhandle
func GetCollectionAddrhandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var svr services.GetCollectionAddr
		if err := c.ShouldBind(&svr); err == nil {
			addr, err := svr.GetCollectionAddr()
			if err != nil {
				serializer.ErrorResponse(c, 40009, "获取归集接收地址失败", err.Error())
			} else {
				serializer.SuccessResponse(c, addr, "success")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}

// SaveGasFromAddrHandle
func SaveGasFromAddrHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var svr services.SetGasFrom
		if err := c.ShouldBind(&svr); err == nil {
			if err := svr.SaveGasFromAddr(); err != nil {
				serializer.ErrorResponse(c, 40010, "保存gas from address failed", err.Error())
			} else {
				serializer.SuccessResponse(c, nil, "success")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}

// GetGasFromAddrHandle
func GetGasFromAddrHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var svr services.GetGasFrom
		if err := c.ShouldBind(&svr); err == nil {
			addr, err := svr.GetGasFromAddr()
			if err != nil {
				serializer.ErrorResponse(c, 40011, "获取gas from address failed", err.Error())
			} else {
				serializer.SuccessResponse(c, addr, "success")
			}
		} else {
			serializer.ErrorResponse(c, 5001, "参数错误", err.Error())
		}
	}
}

// ProcessCollectUSDTHandle
func ProcessCollectUSDTHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		svr := &services.ProcessCollection{}
		if err := svr.ProcessCollectUSDT(); err != nil {
			serializer.ErrorResponse(c, 40011, "归集失败", err.Error())
		} else {
			serializer.SuccessResponse(c, nil, "success")
		}

	}
}
