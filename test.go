package main

import (
	db "filestore/db/postgres"
)

func main() {
	conn := db.DBConn()
	conn.Ping()
}
