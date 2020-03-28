package models

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
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
}
