package transaction

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
