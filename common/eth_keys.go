package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// aes加解密的盐值
const CIPHER = "~C·H!I@P#U$T%A^O&B*(T)U-P+.T/A:"

// GenerateEthKey 生成以太坊密钥对
func GenerateEthKey() (*AccountKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	// 获取hex类型的私钥
	privateToBytes := crypto.FromECDSA(privateKey)
	private := common.ToHex(privateToBytes)

	// 通过privateKey获取地址
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &AccountKey{
		Address: address,
		Private: private,
	}, nil
}

// 对私钥进行对称加密
func Encrypt(private string) (string, error) {
	data := common.FromHex(private)
	result, err := AesEncrypt(data, []byte(CIPHER))
	if err != nil {
		return "", err
	}
	return common.ToHex(result), nil
}

// 解密
func Decrypt(val string) (string, error) {
	data := common.FromHex(val)
	result, err := AesDecrypt(data, []byte(CIPHER))
	if err != nil {
		return "", err
	}
	return common.ToHex(result), nil
}
