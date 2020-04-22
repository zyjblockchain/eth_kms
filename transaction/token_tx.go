package transaction

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

// 以太坊token交易
// getContractFunctionCode 计算合约函数code
func getContractFunctionCode(funcName string) []byte {
	h := crypto.Keccak256Hash([]byte(funcName))
	return h.Bytes()[:4]
}

// formatArgs 把参数转换成[32]byte的数组类型
func formatArgs(args string) []byte {
	b := common.FromHex(args)
	var h [32]byte
	if len(b) > len(h) {
		b = b[len(b)-32:]
	}
	copy(h[32-len(b):], b)
	return h[:]
}

// NewERC20TokenTx 返回的是rawTransaction
func NewERC20TokenTx(senderNonce uint64, receiver common.Address, contractAddr common.Address, gasLimit uint64, gasPrice *big.Int, tokenAmount *big.Int) *types.Transaction {
	/**
	transferFun := "0xa9059cbb"
	receiverAddrCode := 000000000000000000000000b1e15fdbe88b7e7c47552e2d33cd5a9b2e0fd478 // eg: 代币接收地址code
	tokenAmountCode := "0000000000000000000000000000000000000000000000000000000000000064" // eg: 转币数量100
	*/
	funcName := "transfer(address,uint256)"
	funcCode := getContractFunctionCode(funcName)
	receiverAddrCode := formatArgs(receiver.Hex())
	AmountCode := formatArgs(tokenAmount.Text(16))

	// 组合生成执行合约的input
	inputData := make([]byte, 0)
	inputData = append(append(funcCode, receiverAddrCode...), AmountCode...) // 顺序千万不能乱，可以在etherscan上找个合约交易查看input data

	// 组装以太坊交易
	return types.NewTransaction(senderNonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, inputData)
}

// SignRawTx对交易进行签名,主网的chainID 为 1
func SignRawTx(rawTx *types.Transaction, chainID *big.Int, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(rawTx, signer, prv)
	return signedTx, err
}

// GetTokenBalance
func GetTokenBalance(address, contractAddress common.Address, client *ethclient.Client) (*big.Int, error) {
	funcName := "balanceOf(address)"
	funcCode := getContractFunctionCode(funcName)

	// 组合生成执行合约的input
	inputData := make([]byte, 0)
	inputData = append(funcCode, formatArgs(address.Hex())...)

	callMsg := ethereum.CallMsg{
		From: address,          // 钱包地址
		To:   &contractAddress, // 代币合约地址
		Data: inputData,
	}
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}
	res := formatHex(hexutil.Encode(result))
	// res == "0x"
	if len(res) == 2 {
		return big.NewInt(0), nil
	} else {
		return hexutil.DecodeBig(res)
	}
}

// formatHex 去除前置的0
func formatHex(s string) string {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	// 去除前置的所有0
	ss := strings.TrimLeft(s, "0")
	return "0x" + ss
}

// EstimateTokenTxGas 预估代币转账交易gas used使用量
func EstimateTokenTxGas(client *ethclient.Client, tokenAmount *big.Int, from, contractAddress, receiver common.Address) (uint64, error) {
	funcName := "transfer(address,uint256)"
	funcCode := getContractFunctionCode(funcName)
	receiverAddrCode := formatArgs(receiver.Hex())
	AmountCode := formatArgs(tokenAmount.Text(16))
	// 组合生成执行合约的input
	inputData := make([]byte, 0)
	inputData = append(append(funcCode, receiverAddrCode...), AmountCode...)

	callMsg := ethereum.CallMsg{
		From:     from,
		To:       &contractAddress,
		GasPrice: nil,
		Data:     inputData,
	}
	return client.EstimateGas(context.Background(), callMsg)
}
