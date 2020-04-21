package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/mssql"
)

// 交易记录表
type Transaction struct {
	gorm.Model
	TxHash    string `gorm:"not null;unique"`
	State     int    `gorm:"index:t_idx"` // 交易状态
	TxType    int    `gorm:"index:t_idx"` // 交易类型
	TxContent mssql.JSON
}

// add
func (t *Transaction) Add() error {
	return DB.Create(t).Error
}

// Get
func (t *Transaction) Get() (*Transaction, error) {
	var tt *Transaction
	err := DB.Where(t).First(tt).Error
	return tt, err
}
