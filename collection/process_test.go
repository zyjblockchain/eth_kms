package collection

import (
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zyjblockchain/eth_kms/models"
	"github.com/zyjblockchain/eth_kms/transaction"
	"testing"
)

// 初始化数据库
func init() {
	dsn := "kms_dev:kms_good123@tcp(rm-j6c49n23e4d07l8ijmo.mysql.rds.aliyuncs.com:3306)/eth_kms_dev?charset=utf8&parseTime=True&loc=Local"
	models.InitDB(dsn)
}

const contractAddress = "0x03332638A6b4F5442E85d6e6aDF929Cd678914f1" // lemo测试币的合约地址

func TestBatchGetCanCollectAddress(t *testing.T) {
	client := transaction.NewEthClient(transaction.RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()
	err := BatchGetCanCollectAddress(common.HexToAddress(contractAddress), 0, 10, client)
	t.Log(err)
}

func TestSendGasFeeForColAddr(t *testing.T) {
	client := transaction.NewEthClient(transaction.RINKEBYNET)
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()
	err := SendGasFeeForColAddr(0, 10, client)
	t.Log(err)
}
