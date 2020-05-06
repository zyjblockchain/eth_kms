package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	GasFromAddrKey       = "gasFromAddr"          // 转出eth作为归集账户的交易gas的地址tag
	CollectionAddrKey    = "CollectionAddr"       // 接收归集的所有USDT地址tag
	CollectionAddrOffset = "CollectionAddrOffset" // 记录上次归集地址的偏移位置
	SendGasFeeOffset     = "SendGasFeeOffset"     // 记录为归集地址转gas费用的偏移位置

)

type Kv struct {
	gorm.Model
	K string
	V string // 地址
}

// newKey 生成对应的key
func newKey(tag string) string {
	return tag
}

// Set 创建或者更新
func Set(tag, val string) error {
	// 查询是否存在，存在则更新，不存在则创建
	var count uint
	err := DB.Model(&Kv{}).Where("k = ?", newKey(tag)).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		// 创建
		kv := Kv{
			K: newKey(tag),
			V: val,
		}
		return DB.Create(&kv).Error
	} else if count == 1 {
		// 更新
		return DB.Model(&Kv{}).Where("k = ?", tag).Update(newKey(tag), val).Error
	} else {
		return errors.New(fmt.Sprintf("数据库中存在两个以上的tag: count = %d, tag = %s", count, tag))
	}
}

// Get
func Get(tag string) (string, error) {
	var kv Kv
	err := DB.Where("k = ?", newKey(tag)).First(&kv).Error
	if err != nil {
		return "", err
	}
	return kv.V, nil
}
