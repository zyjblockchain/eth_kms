package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// 秘钥存储模型
type KeysMgr struct {
	gorm.Model
	Address string
	PriKey  string
}

// AddRecord 添加记录
func AddRecord(newKeys *KeysMgr) (*KeysMgr, error) {
	err := DB.Create(newKeys).Error
	return newKeys, err
}

// 通过address查询记录，如果存在多条，只返回一条记录
func GetKeysByAddr(address string) (*KeysMgr, error) {
	var keysMgr KeysMgr
	err := DB.Where("address = ?", address).First(&keysMgr).Error
	return &keysMgr, err
}

// 通过主键id批量拉取记录,传入起始位置的上一个位置和拉取的数量
func GetBatchById(startId, limit uint) ([]*KeysMgr, uint, error) {
	var keysMgrs []*KeysMgr
	var total uint
	// 查询表中的记录总数
	err := DB.Model(&KeysMgr{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 开始的index已经大于了总数
	if startId > total {
		return nil, total, errors.New("start 超过了总数量")
	}

	err = DB.Limit(limit).Offset(startId).Find(&keysMgrs).Error
	return keysMgrs, total, err
}
