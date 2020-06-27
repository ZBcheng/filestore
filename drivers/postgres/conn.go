package drivers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/zbcheng/filestore/conf"
)

var db *sql.DB

// var mutex sync.Mutex

func init() {

	var pgInfo string

	config := conf.GetConfig()

	host := config.PgConf.Host
	port := config.PgConf.Port
	user := config.PgConf.User
	dbname := config.PgConf.DBName
	password := config.PgConf.Password

	if password != "" {
		pgInfo = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host, port, user, dbname, password)
	} else {
		pgInfo = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			host, port, user, dbname)
	}

	db, _ = sql.Open("postgres", pgInfo)
	db.SetMaxOpenConns(1000)

	if err := db.Ping(); err != nil {
		fmt.Println("Failed to connect to postgres, err: " + err.Error())
		os.Exit(1)
	}
}

// DBConn : 返回postgres连接
func DBConn() *sql.DB {
	return db
}
