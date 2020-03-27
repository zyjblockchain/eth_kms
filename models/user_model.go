package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
	// Active 激活用户
	Active string = "active"
	// Inactive 未激活用户
	Inactive string = "inactive"
	// Suspend 被封禁用户
	Suspend string = "suspend"
)

// user存储模型
type User struct {
	gorm.Model
	UserName       string
	PasswordDigest string
	NickName       string
	Status         string
}

// SetPassword 设置或者修改user在数据库中保存的密码
func (user *User) SetPassword(rawPassword string) error {
	// 对原始密码进行加密
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验用户密码
func (user *User) CheckPassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(rawPassword))
	return err == nil
}

// GetUserById 通过主键id查询user
func GetUserById(id interface{}) (*User, error) {
	var user User
	result := DB.First(&user, id)
	return &user, result.Error
}

// GetUserByUserName
func GetUserByUserName(userName string) (*User, error) {
	var user User
	result := DB.Where("user_name = ?", userName).First(&user)
	return &user, result.Error
}

// AddUser 增加user并返回结果
func AddUser(newUser *User) (error, *User) {
	err := DB.Create(newUser).Error
	return err, newUser
}

// 判断userName或者nickName是否被占用
func IsExistFromUserNameOrNickName(userName, nickName string) bool {
	count := 0
	if err := DB.Where("user_name = ?", userName).Or("nick_name = ?", nickName).Find(&User{}).Count(&count).Error; err != nil {
		// panic(err) // todo
		fmt.Println(err)
	}

	return count > 0
}
