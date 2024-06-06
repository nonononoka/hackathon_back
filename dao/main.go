package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB

func init() {
	err := godotenv.Load()
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_PWD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@(localhost:3307)/%s", mysqlUser, mysqlUserPwd, mysqlDatabase)
	_db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatal("fail: _db.Ping, %v\n", err)
	}
	db = _db
}
