package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
)

func init() {
	dsn := "kms_dev:kms_good123@tcp(rm-j6c49n23e4d07l8ijmo.mysql.rds.aliyuncs.com:3306)/eth_kms_dev?charset=utf8&parseTime=True&loc=Local"
	InitDB(dsn)
}
func TestGet(t *testing.T) {
	str, err := Get("aa")
	t.Log(err, str)
}
