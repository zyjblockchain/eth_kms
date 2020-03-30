package common

import (
	"errors"
	"github.com/zyjblockchain/eth_kms/models"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type AccountKey struct {
	Address common.Address `json:"address"`
	Private string         `json:"private"`
}

// GetAccountKeyByAddr
func GetAccountKeyByAddr(address string) (*AccountKey, error) {
	// 从数据库中查询对应的加密秘钥 todo 测试查询不存在的数据的情况是报错还是返回空数据
	keysMgr, err := models.GetKeysByAddr(address)
	if err != nil {
		return nil, err
	}
	// 解密数据库中的加密之后的密码
	private, err := Decrypt(keysMgr.PriKey)
	return &AccountKey{
		Address: common.HexToAddress(address),
		Private: private,
	}, err
}

// SignData 数据签名
func (a *AccountKey) SignData(data []byte) ([]byte, error) {
	privateKey, err := crypto.HexToECDSA(a.Private[2:])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("crypto.HexToECDSA(a.Private[2:]) error: %v", err))
	}

	hash := common.BytesToHash(data)
	sig, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("签名失败： %v", err))
	}
	return sig, nil
}
