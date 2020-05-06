package conf

import (
	"github.com/joho/godotenv"
	"github.com/zyjblockchain/eth_kms/models"
	"os"
)

// Init 初始化数据库
func Init() {
	// 从本地的配置文件中读取配置文件到环境变量中
	if err := godotenv.Load(".env_dev"); err != nil {
		panic(err)
	}
	// 链接数据库
	models.InitDB(os.Getenv("MYSQL_DSN"))
}
