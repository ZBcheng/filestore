package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var db *sql.DB
var mutex sync.Mutex

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "0000"
	dbname   = "filestore"
)

func init() {
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
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
