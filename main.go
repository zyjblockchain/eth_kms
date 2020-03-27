package main

import (
	"eth_kms/conf"
	"eth_kms/routers"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 初始化数据库
	conf.Init()
	routers.NewRouter(":3000")
}
