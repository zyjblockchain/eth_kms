package models

import (
	"github.com/jinzhu/gorm"
)

// 交易记录表
type Transaction struct {
	gorm.Model
	TxHash string
	State  int // 交易状态，0：未发送，1：pending，2：发送失败，3：发送成功，4：交易超时
	TxType int // 交易类型，0：以太坊转账交易，1：以太坊上的usdt转账交易
	From   string
	To     string
	Amount string
	TxInfo string // 交易详细信息
	ErrMsg string // 如果交易发送失败，失败的msg
}

// Add
func (t *Transaction) Add() error {
	return DB.Create(t).Error
}

// Get
func (t *Transaction) Get() (*Transaction, error) {
	var tt *Transaction
	err := DB.Where(t).First(tt).Error
	return tt, err
}

// Update
func (t *Transaction) Update(tt Transaction) error {
	return DB.Model(t).Updates(tt).Error
}
