package services

import "github.com/zyjblockchain/eth_kms/common"

type GetPrivate struct {
	Address string `form:"address" json:"address"`
}

// GetPriv 返回私钥
func (g *GetPrivate) GetPriv() (string, error) {
	accKey, err := common.GetAccountKeyByAddr(g.Address)
	if err != nil {
		return "", err
	} else {
		return accKey.Private, err
	}
}
