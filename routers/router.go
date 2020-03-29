package routers

import (
	"eth_kms/handles"
	"eth_kms/middleware"
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

	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
