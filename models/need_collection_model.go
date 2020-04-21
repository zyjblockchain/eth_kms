package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"math/big"
)

// 需要归集的地址表
type NeedCollectionAddress struct {
	gorm.Model
	Address     string  `gorm:"not null;unique"`
	State       int     `gorm:"index:n_idx"` // 0:未开始，1:进行中，2:已完成，3:失败
	EthBalance  big.Int `gorm:"index:n_idx"` // eth余额，单位为wei
	UsdtBalance big.Int `gorm:"index:n_idx"` // usdt余额，小数位6位的最小单位余额
}

//  Save 保存记录
func (n *NeedCollectionAddress) Save() error {
	return DB.Create(n).Error
}

// UpdateInfo 更新
func UpdateInfo(nca NeedCollectionAddress) error {
	return DB.Model(&NeedCollectionAddress{}).Where("address = ?", nca.Address).Updates(nca).Error
}

// 批量拉取 GetBatchAddrById
func GetBatchAddrById(startId, limit uint) ([]*NeedCollectionAddress, uint, error) {
	var records []*NeedCollectionAddress
	var total uint
	// 查询表中的记录总数
	err := DB.Model(&NeedCollectionAddress{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 开始的index已经大于了总数
	if startId > total {
		return nil, total, errors.New("startId 超过了数据库中的地址总数量")
	}
	err = DB.Limit(limit).Offset(startId).Find(&records).Error
	return records, total, err
}
