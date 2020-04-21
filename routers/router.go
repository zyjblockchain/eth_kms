package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zyjblockchain/eth_kms/handles"
	"github.com/zyjblockchain/eth_kms/middleware"
)

func NewRouter(addr string) {
	r := gin.Default()
	// 跨域访问过滤
	r.Use(middleware.Cors())

	// 生成秘钥对
	r.GET("/kms/new_key", handles.NewEthKeyHandle())
	// 传入密钥对进行存储
	r.POST("/kms/save", handles.SaveKeysHandle())
	// 分页拉取地址,传入参数{startId: 99, limit: 50}
	r.POST("/kms/batch_get_address", handles.BatchGetAddrHandle())
	// 签名接口 {address: ""0x123ssdd..., data: "0xkkkkkkkkkkkkkk..."}
	r.POST("/kms/sign", handles.SignDataHandle())
	// 拉取通过地址获取还原私钥 todo 此接口以后会删除，现在临时使用
	r.POST("kms/get_private", handles.GetPrivateHandle())

	// 设置归集资产接收地址 {address: "0x123ssdd..."}
	r.POST("/kms/set_collection_address", handles.SetCollectionAddrHandle())
	// get归集资产接收地址
	r.GET("/kms/get_collection_address", handles.GetCollectionAddrhandle())
	// 保存gas from address {address: "0x123ssdd...", private: "0xkkkkkkkkkkkkkk..."}
	r.POST("/kms/save_gas_from_address", handles.SaveGasFromAddrHandle())
	// get gas from address
	r.GET("/kms/get_gas_from_address", handles.GetGasFromAddrHandle())

	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
