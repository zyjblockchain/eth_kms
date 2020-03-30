package main

import (
	"github.com/zyjblockchain/eth_kms/conf"
	"github.com/zyjblockchain/eth_kms/routers"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 初始化数据库
	conf.Init()
	routers.NewRouter(":3000")
}
