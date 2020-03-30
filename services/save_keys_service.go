package services

import (
	"github.com/zyjblockchain/eth_kms/common"
	"github.com/zyjblockchain/eth_kms/models"
)

type SaveKeysInfo struct {
	Address string `form:"address" json:"address"`
	Private string `form:"private" json:"private"`
}

// SaveKeys
func (s *SaveKeysInfo) SaveKeys() error {
	// 秘钥加密保存数据库
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
	return err
}
