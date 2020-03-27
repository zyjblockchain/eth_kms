package services

import (
	"errors"
	"eth_kms/models"
)

// 用户登录信息
type UserLoginInfo struct {
	UserName string `json:"username" form:"username" binding:"required,min=5,max=30"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=40"`
}

// UserLogin 用户登录
func (u *UserLoginInfo) UserLogin() (*models.User, error) {
	// 通过userName在数据库中查找对应的userInfo
	user, err := models.GetUserByUserName(u.UserName)
	if err != nil {
		// 表示没有查到
		return nil, err
	}
	// 验证password是否正确
	if !user.CheckPassword(u.Password) {
		return nil, errors.New("密码不正确")
	}
	return user, nil
}
