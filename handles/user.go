package handles

import (
	"eth_kms/models"
	"eth_kms/serializer"
	"eth_kms/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Register 用户注册接口
func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service services.UserRegisterInfo
		// 解析传入的参数
		if err := c.ShouldBind(&service); err == nil {
			// 解析成功则进行注册
			if user, err := service.RegisterUser(); err != nil {
				serializer.ErrorResponse(c, 40001, "注册失败", err.Error())
			} else {
				// 返回注册的user存储信息
				serializer.SuccessResponse(c, user, "注册成功")
			}
		} else {
			serializer.ErrorResponse(c, 40001, "参数解析失败", err.Error())
		}
	}
}

// Login 用户登录接口
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginInfo services.UserLoginInfo
		if err := c.ShouldBind(&loginInfo); err == nil {
			user, err := loginInfo.UserLogin()
			if err != nil {
				// 用户登录失败
				serializer.ErrorResponse(c, 40001, "登录不成功", err.Error())
			} else {
				// 登录校验成功
				// 设置session
				s := sessions.Default(c)
				// 删除之前设置的kv
				s.Clear()
				// 重新设置kv
				s.Set("user_id", user.ID)
				// 保存当前请求的session
				_ = s.Save()

				// 返回登录结果
				serializer.SuccessResponse(c, *user, "登录成功")
			}
		} else {
			// 解析请求参数出错
			serializer.ErrorResponse(c, 40001, "参数解析失败", err.Error())
		}
	}
}

// 拉取自己的用户信息
func GetMine() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, exist := c.Get("user"); exist {
			if u, ok := user.(*models.User); ok {
				serializer.SuccessResponse(c, *u, "得到用户信息")
			}
		} else {
			serializer.ErrorResponse(c, 40001, "获取当前用户信息失败，可能没有登录", "not found")
		}
	}
}

// 登出
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 删除session记录
		s := sessions.Default(c)
		s.Clear()
		_ = s.Save()
		serializer.SuccessResponse(c, nil, "登出成功")
	}
}
