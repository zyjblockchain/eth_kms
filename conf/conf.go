package conf

import (
	"github.com/zyjblockchain/eth_kms/models"
	"github.com/joho/godotenv"
	"os"
)

// Init 初始化数据库
func Init() {
	// 从本地的.env文件中读取配置文件到环境变量中
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	// 链接数据库
	models.InitDB(os.Getenv("MYSQL_DSN"))
}
