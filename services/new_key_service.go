package services

import (
	"eth_kms/common"
	"eth_kms/models"
)

// 生成密钥对
type NewKeysInfo struct {
}

// NewKeys 返回以太坊地址
func (k *NewKeysInfo) NewKeys() (string, error) {
	// 生成以太坊的密钥对
	accountKey, err := common.GenerateEthKey()
	if err != nil {
		return "", err
	}
	// 地址
	strAddress := accountKey.Address.String()
	// 为私钥进行aes对称加密
	result, err := common.Encrypt(accountKey.Private)
	if err != nil {
		return "", err
	}
	// 保存加密之后的密钥对在数据库
	keys := &models.KeysMgr{
		Address: strAddress,
		PriKey:  result,
	}
	_, err = models.AddRecord(keys)
	if err != nil {
		return "", err
	}
	return strAddress, nil
}
