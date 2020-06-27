package drivers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zbcheng/filestore/conf"
)

var db *sql.DB

func init() {
	var err error

	config := conf.GetConfig()
	username := config.MysqlConf.User
	password := config.MysqlConf.Password
	dbName := config.MysqlConf.DBName

	mysqlInfo := fmt.Sprintf("%s:%s@/%s", username, password, dbName)
	db, err = sql.Open("mysql", mysqlInfo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}
