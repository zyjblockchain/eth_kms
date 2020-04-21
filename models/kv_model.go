package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	GasFromAddrKey    = "gasFromAddr"    // 转出eth作为归集账户的交易gas的地址tag
	CollectionAddrKey = "CollectionAddr" // 接收归集的所有USDT地址tag
)

type Kv struct {
	gorm.Model
	Key string
	Val string // 地址
}

// newKey 生成对应的key
func newKey(tag string) string {
	return tag
}

// Set 创建或者更新
func Set(tag, addr string) error {
	// 查询是否存在，存在则更新，不存在则创建
	var count uint
	err := DB.Model(&Kv{}).Where("key = ?", newKey(tag)).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		// 创建
		kv := Kv{
			Key: newKey(tag),
			Val: addr,
		}
		return DB.Create(&kv).Error
	} else if count == 1 {
		// 更新
		return DB.Model(&Kv{}).Where("key = ?", tag).Update(newKey(tag), addr).Error
	} else {
		return errors.New(fmt.Sprintf("数据库中存在两个以上的tag: count = %d, tag = %s", count, tag))
	}
}

// Get
func Get(tag string) (string, error) {
	var kv Kv
	err := DB.Where("key = ?", newKey(tag)).First(&kv).Error
	if err != nil {
		return "", err
	}
	return kv.Val, nil
}
