package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

// 需要归集的地址表
type CollectionAddress struct {
	gorm.Model
	Address      string
	State        int    // 0:待处理，1:可以进行归集操作，2:进行中，3:已完成，4:失败
	EthBalance   string // eth余额，单位为wei
	TokenBalance string // token余额，小数位6位的最小单位余额
	TokenAddress string // token所在的智能合约地址
	GasPrice     string // 发送交易的gas price
	GasLimit     uint64 // gasLimit
}

func (n *CollectionAddress) Println() string {
	return fmt.Sprintf("address: %s, state: %d, ethBalance: %s, "+
		"tokenBalance: %s, tokenAddress: %s",
		n.Address, n.State, n.EthBalance, n.TokenBalance, n.TokenAddress)
}

//  Save 保存记录
func (n *CollectionAddress) Save() error {
	return DB.Create(n).Error
}

// UpdateInfo 更新
func UpdateInfo(nca CollectionAddress) error {
	return DB.Model(&CollectionAddress{}).Where("address = ?", nca.Address).Updates(nca).Error
}

// 批量拉取 GetBatchAddrById
func GetBatchAddrById(startId, limit uint) ([]*CollectionAddress, uint, error) {
	var records []*CollectionAddress
	var total uint
	// 查询表中的记录总数
	err := DB.Model(&CollectionAddress{}).Count(&total).Error
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

//  InsertBatch 批量插入数据
func InsertBatch(records []*CollectionAddress) error {
	sql := "INSERT INTO `collection_addresses` (`address`,`state`,`token_balance`, `token_address`) VALUES "
	// 循环data数组,组合sql语句
	for key, value := range records {
		if len(records)-1 == key {
			// 最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s','%d','%s','%s');", value.Address, value.State, value.TokenBalance, value.TokenAddress)
		} else {
			sql += fmt.Sprintf("('%s','%d','%s','%s'),", value.Address, value.State, value.TokenBalance, value.TokenAddress)
		}
	}
	err := DB.Exec(sql).Error
	return err
}
