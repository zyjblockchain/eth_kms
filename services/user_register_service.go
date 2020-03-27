package services

import (
	"errors"
	"eth_kms/models"
)

// UserRegisterInfo
type UserRegisterInfo struct {
	UserName        string `json:"username" form:"username" binding:"required,min=5,max=30"`
	NickName        string `json:"nickname" form:"nickname" binding:"required,min=2,max=30"`
	Password        string `json:"password" form:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm" binding:"required,min=8,max=40"`
}

// Register 注册用户
func (s *UserRegisterInfo) RegisterUser() (*models.User, error) {
	// 1. 验证注册info
	if err := s.validate(); err != nil {
		return nil, err
	}

	// 2. 实例化user对象
	user := &models.User{
		UserName:       s.UserName,
		PasswordDigest: "",
		NickName:       s.NickName,
		Status:         models.Active,
	}

	// 3. 储存加密之后的password
	if err := user.SetPassword(s.Password); err != nil {
		return user, err
	}

	// 4. 持久化用户
	err, newUser := models.AddUser(user)
	if err != nil {
		// 保存数据失败
		return user, err
	} else {
		return newUser, nil
	}
}

func (s *UserRegisterInfo) validate() error {
	// 验证两次输入的密码是否相等
	if s.Password != s.PasswordConfirm {
		return errors.New("输入的两次密码不相等")
	}

	// 验证name是否被占用
	if models.IsExistFromUserNameOrNickName(s.UserName, s.NickName) {
		return errors.New("user name 和 nick name已经被占用")
	}
	return nil
}
