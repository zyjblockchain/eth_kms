package transaction

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

const MAINNET = "https://mainnet.infura.io/KoLQDsHeWLs20urjat1X"
const RINKEBYNET = "https://rinkeby.infura.io/v3/36b98a13557c4b8583d57934ede2f74d"

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
	// TODO  RINKEBYNET 的chainID为4,主网的chainID为1
	signedTx, err := SignRawTx(rawTx, big.NewInt(4), prv)
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
	// TODO  RINKEBYNET 的chainID为4,主网的chainID为1
	signedTx, err := SignRawTx(rawTx, big.NewInt(4), prv)
	if err != nil {
		panic(err)
	}
	// 把签好名的交易发送到网络
	err = client.SendTransaction(context.Background(), signedTx)
	return signedTx, err
}

// // SendTx
// func SendTx(signedTx *types.Transaction, rawurl string) (common.Hash, error) {
// 	// 连接网络
// 	rpcDial, err := rpc.Dial(rawurl)
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
// 	client := ethclient.NewClient(rpcDial)
//
// 	// 发送交易
// 	err = client.SendTransaction(context.Background(),signedTx)
// 	// 返回交易hash
// 	return signedTx.Hash(), err
// }
