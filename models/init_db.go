package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	// // 配置MySQL连接参数
	// username := "kms_dev"  // 账号
	// password := "kms_good123" // 密码
	// host := "rm-j6c49n23e4d07l8ijmo.mysql.rds.aliyuncs.com" // 数据库地址，可以是Ip或者域名
	// port := 3306 // 数据库端口
	// Dbname := "eth_kms" // 数据库名
	//
	// // 通过前面的数据库参数，拼接MYSQL DSN， 其实就是数据库连接串（数据源名称）
	// // MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	// // 类似{username}使用花括号包着的名字都是需要替换的参数
	// dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	fmt.Println("dsn: ", dsn)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// 设置数据库日志级别
	if gin.Mode() == gin.ReleaseMode {
		db.LogMode(false)
	} else {
		db.LogMode(true)
	}

	DB = db
	autoCreateTable()
}

// 自动建表
func autoCreateTable() {
	DB.AutoMigrate(&KeysMgr{})
	DB.AutoMigrate(&Kv{})
}
