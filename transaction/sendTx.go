package transaction

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/labstack/gommon/log"
	"math/big"
	"os"
)

func NewEthClient(rawurl string) *ethclient.Client {
	// 连接网络
	rpcDial, err := rpc.Dial(rawurl)
	if err != nil {
		return nil
	}
	return ethclient.NewClient(rpcDial)
}

// GetEthBalance 获取地址的eth余额
func GetEthBalance(address common.Address, client *ethclient.Client) (*big.Int, error) {
	return client.BalanceAt(context.Background(), address, nil)
}

// 发送以太坊交易
func SendEthTx(private string, nonce, gasLimit uint64, gasPrice *big.Int, to common.Address, amount *big.Int, client *ethclient.Client) (*types.Transaction, error) {
	// 构建原生交易
	rawTx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	// 对原生交易进行签名
	prv, err := crypto.ToECDSA(common.FromHex(private))
	if err != nil {
		panic(err)
	}
	// RINKEBYNET 的chainID为4,主网的chainID为1
	chainId, b := new(big.Int).SetString(os.Getenv("ETH_CHAIN_ID"), 10)
	if !b {
		panic("获取chainId失败")
	}
	log.Infof("chainId: %d", chainId)
	signedTx, err := SignRawTx(rawTx, chainId, prv)
	if err != nil {
		panic(err)
	}
	// 把签好名的交易发送到网络
	err = client.SendTransaction(context.Background(), signedTx)
	return signedTx, err
}

// SendTokenTx 发送token交易
func SendTokenTx(private string, nonce, gasLimit uint64, gasPrice *big.Int, tokenReceiver, contractAddr common.Address, tokenAmount *big.Int, client *ethclient.Client) (*types.Transaction, error) {
	rawTx := NewERC20TokenTx(nonce, tokenReceiver, contractAddr, gasLimit, gasPrice, tokenAmount)
	// 对原生交易进行签名
	prv, err := crypto.ToECDSA(common.FromHex(private))
	if err != nil {
		panic(err)
	}
	// RINKEBYNET 的chainID为4,主网的chainID为1
	chainId, b := new(big.Int).SetString(os.Getenv("ETH_CHAIN_ID"), 10)
	if !b {
		panic("获取chainId失败")
	}
	signedTx, err := SignRawTx(rawTx, chainId, prv)
	if err != nil {
		panic(err)
	}
	// 把签好名的交易发送到网络
	err = client.SendTransaction(context.Background(), signedTx)
	return signedTx, err
}
