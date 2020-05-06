package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zyjblockchain/eth_kms/conf"
	"github.com/zyjblockchain/eth_kms/routers"
	"github.com/zyjblockchain/sandy_log/log"
)

func init() {
	// 初始化日志级别、格式、是否保存到文件
	log.Setup(log.LevelDebug, true, true)
}

func main() {
	// 初始化数据库
	conf.Init()
	routers.NewRouter(":3000")
}
