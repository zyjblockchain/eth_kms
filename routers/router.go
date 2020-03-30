package routers

import (
	"github.com/zyjblockchain/eth_kms/handles"
	"github.com/zyjblockchain/eth_kms/middleware"
	"github.com/gin-gonic/gin"
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

	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
