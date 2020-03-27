package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 跨域中间件
func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	// 重新设置默认设置中的跨域字段
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie"}
	config.AllowOrigins = []string{"http://localhost:8080", "https://www.gourouting.com"}
	config.AllowCredentials = true
	return cors.New(config)
}
