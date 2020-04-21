package services

import (
	"github.com/zyjblockchain/eth_kms/common"
	"github.com/zyjblockchain/eth_kms/models"
)

// 设置归集gas分发地址
type SetGasFrom struct {
	Address string `form:"address" json:"address"`
	Private string `form:"private" json:"private"`
}

// SaveGasFromAddr
func (s *SetGasFrom) SaveGasFromAddr() error {
	// 密钥加密保存数据库
	result, err := common.Encrypt(s.Private)
	if err != nil {
		return err
	}
	// 保存加密之后的密钥对在数据库
	keys := &models.KeysMgr{
		Address: s.Address,
		PriKey:  result,
	}
	_, err = models.AddRecord(keys)
	if err != nil {
		return err
	}

	// 地址保存到kv表中
	return models.Set(models.GasFromAddrKey, s.Address)
}

type GetGasFrom struct {
}

// GetGasFromAddr
func (s *GetGasFrom) GetGasFromAddr() (string, error) {
	return models.Get(models.GasFromAddrKey)
}
