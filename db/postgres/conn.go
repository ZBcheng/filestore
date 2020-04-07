package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/Unknwon/goconfig"
	_ "github.com/lib/pq"
)

var db *sql.DB
var mutex sync.Mutex

func init() {

	config, err := goconfig.LoadConfigFile("db.conf")
	if err != nil {
		fmt.Println("Failed to read db.conf, err: " + err.Error())
		os.Exit(1)
	}

	host, _ := config.GetValue("postgres", "host")
	port, _ := config.GetValue("postgres", "port")
	user, _ := config.GetValue("postgres", "user")
	dbname, _ := config.GetValue("postgres", "dbname")
	password, _ := config.GetValue("postgres", "password")

	pgInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
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
