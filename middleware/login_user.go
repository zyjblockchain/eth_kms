package middleware

import (
	"eth_kms/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetLoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		// 如果用户登录了则uid不为空
		if uid != nil {
			fmt.Println("用户已经登录过一次了，uid = ", uid)
			// 通过用户uid在数据库中找到对应的user
			user, err := models.GetUserById(uid)
			if err == nil {
				// 返回的user保存进上下文中，用于后面验证用户是否登录过使用
				c.Set("user", user)
			}
		}
		// 保存此上下文状态传递到后面的handle中去
		c.Next()
	}
}
