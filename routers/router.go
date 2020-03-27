package routers

import (
	"eth_kms/handles"
	"eth_kms/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

func NewRouter(addr string) {
	r := gin.Default()
	// 执行中间件
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.SetLoginUser())

	// 注册路由
	v1 := r.Group("/api/v1")
	{
		// 心跳检测接口
		v1.POST("ping", handles.Ping)

		// 1. 用户注册接口
		v1.POST("user/register", handles.Register())
		// 2. 用户登录接口
		v1.POST("user/login", handles.Login())

		// 3. 需要登录保护
		authed := v1.Group("/")
		// 需要登录授权才能访问的接口
		authed.Use(middleware.AuthRequired())
		{
			// 拉取自己的用户信息
			authed.GET("user/me", handles.GetMine())
			// 登出
			authed.DELETE("user/logout", handles.Logout())
			// 执行交易数据签名接口

		}
		// 生成秘钥对
		v1.GET("/new_key", handles.NewEthKey())

	}

	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
