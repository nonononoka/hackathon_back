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
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	env := os.Getenv("ENV")
	if env == "production" {
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPwd := os.Getenv("MYSQL_PWD")
		mysqlHost := os.Getenv("MYSQL_HOST")
		mysqlDatabase := os.Getenv("MYSQL_DATABASE")
		connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
		_db, err := sql.Open("mysql", connStr)
		//_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlUserPwd, mysqlDatabase))
		if err != nil {
			log.Fatal(err)
		}
		if err := _db.Ping(); err != nil {
			log.Fatal("fail: _db.Ping, %v\n", err)
		}
		db = _db
	} else {
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPwd := os.Getenv("MYSQL_PWD")
		mysqlDatabase := os.Getenv("MYSQL_DATABASE")
		_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3307)/%s", mysqlUser, mysqlPwd, mysqlDatabase))
		if err != nil {
			log.Fatalf("fail: sql.Open, %v\n", err)
		}
		if err := _db.Ping(); err != nil {
			log.Fatalf("fail: _db.Ping, %v\n", err)
		}
		db = _db
	}
}
