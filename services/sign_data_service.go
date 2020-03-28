package services

import (
	common2 "eth_kms/common"
	"github.com/ethereum/go-ethereum/common"
)

type SignInfo struct {
	Address string `form:"address" json:"address"`
	Data    string `form:"data" json:"data"`
}

// KeySign 数据签名
func (s *SignInfo) KeySign() (string, error) {
	byteData := common.FromHex(s.Data)
	accKey, err := common2.GetAccountKeyByAddr(s.Address)
	if err != nil {
		return "", err
	}
	sig, err := accKey.SignData(byteData)
	// bytes to hex
	return common.ToHex(sig), err
}
