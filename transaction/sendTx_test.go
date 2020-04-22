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
	toAddr := "0x7AC954Ed6c2d96d48BBad405aa1579C828409f59"                         // 接收者地址

	client := NewEthClient(RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()
	// 查询发送者账户nonce
	nonce, err := client.NonceAt(context.Background(), common.HexToAddress(fromAddr), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("nonce: ", nonce)
	gasLimit := uint64(60000)
	gasPrice := big.NewInt(5000000000)
	tokenAmount, _ := new(big.Int).SetString("1000000000000000000", 10)
	// 生成原生交易
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

func TestGetTokenBalance(t *testing.T) {
	client := NewEthClient(RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()
	const contractAddress = "0x03332638A6b4F5442E85d6e6aDF929Cd678914f1" // lemo测试币的合约地址
	fromAddr := "0x59375A522876aB96B0ed2953D0D3b92674701Cc2"
	res, err := GetTokenBalance(common.HexToAddress(fromAddr), common.HexToAddress(contractAddress), client)
	t.Log("ddd: ", res)
	t.Log(res.String(), err)
}

func TestGetEthBalance(t *testing.T) {
	client := NewEthClient(RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()

	addr := "0x59375A522876aB96B0ed2953D0D3b92674701Cc2"
	balance, err := GetEthBalance(common.HexToAddress(addr), client)
	t.Log(balance.String(), err) // balance: 842893818200000000 wei
}

func TestSendEthTx(t *testing.T) {
	client := NewEthClient(RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()

	fromAddr := "0x59375A522876aB96B0ed2953D0D3b92674701Cc2"
	fromPriv := "69F657EAF364969CCFB2531F45D9C9EFAC0A63E359CEA51E5F7D8340784168D2" // 发送者私钥
	toAddr := "0x7AC954Ed6c2d96d48BBad405aa1579C828409f59"                         // 以太坊接收账户
	// 查询发送者账户nonce
	nonce, err := client.NonceAt(context.Background(), common.HexToAddress(fromAddr), nil)
	t.Log("111: ", err)
	suggestGasPrice, err := client.SuggestGasPrice(context.Background())
	t.Log("222: ", err)
	amount := big.NewInt(8200000000)
	tx, err := SendEthTx(fromPriv, nonce, uint64(22000), suggestGasPrice, common.HexToAddress(toAddr), amount, client)
	t.Log("333: ", err)
	t.Log("tx: ", tx.Hash().Hex())
}

func TestEstimateTokenTxGas(t *testing.T) {
	client := NewEthClient(RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()
	from := common.HexToAddress("0x59375a522876ab96b0ed2953d0d3b92674701cc2")
	receiver := common.HexToAddress("0xCef4DBEfd5E85D1f500B7a568a29208feeD2fb79")
	contractAddress := common.HexToAddress("0x03332638a6b4f5442e85d6e6adf929cd678914f1")
	amount, _ := new(big.Int).SetString("1000000000000000000", 10)
	gasUsed, err := EstimateTokenTxGas(client, amount, from, contractAddress, receiver)

	t.Log(gasUsed, err)
}
