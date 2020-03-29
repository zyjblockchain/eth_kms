package transaction

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
)

// TestSendTx 这里使用lemo测试代币进行测试
func TestSendTx(t *testing.T) {
	const contractAddress = "0x03332638A6b4F5442E85d6e6aDF929Cd678914f1"           // lemo测试币的合约地址
	fromPriv := "69F657EAF364969CCFB2531F45D9C9EFAC0A63E359CEA51E5F7D8340784168D2" // 发送者私钥
	fromAddr := "0x59375A522876aB96B0ed2953D0D3b92674701Cc2"                       // 发送者地址
	toAddr := "0x415979DC0266fd94A3CB90a04EC28853FCeB1A34"                         // 接收者地址

	client := NewEthClient(RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	// 查询发送者账户nonce
	nonce, err := client.NonceAt(context.Background(), common.HexToAddress(fromAddr), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("nonce: ", nonce)
	gasLimit := uint64(60000)
	gasPrice := big.NewInt(5000000000)
	tokenAmount := uint64(1000000000000000000)
	// 生成原生地址
	rawTx := NewERC20TokenTx(nonce, common.HexToAddress(toAddr), common.HexToAddress(contractAddress), gasLimit, gasPrice, tokenAmount)
	// 对原生交易进行签名
	prv, err := crypto.ToECDSA(common.FromHex(fromPriv))
	if err != nil {
		panic(err)
	}
	// TODO  RINKEBYNET 的chainID为4,主网的chainID为1
	signedTx, err := SignRawTx(rawTx, big.NewInt(4), prv)
	if err != nil {
		panic(err)
	}
	// 把签名好的交易发送到测试网络中
	err = client.SendTransaction(context.Background(), signedTx)
	fmt.Printf("txHash: %s\n err: %v", signedTx.Hash().Hex(), err)
}
