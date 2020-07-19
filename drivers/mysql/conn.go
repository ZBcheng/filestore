package drivers

import (
	"os"

	"github.com/arstd/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zbcheng/filestore/conf"
)

var db *gorm.DB

func init() {
	var err error

	config := conf.Load()
	username := config.MysqlConf.User
	password := config.MysqlConf.Password
	host := config.MysqlConf.Host
	port := config.MysqlConf.Port
	dbName := config.MysqlConf.DBName

	connInfo := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8" + "&parseTime=true"
	db, err = gorm.Open("mysql", connInfo)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func DBConn() *gorm.DB {
	return db
}
