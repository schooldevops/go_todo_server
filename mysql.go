package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// var connectionString = "boarduser:boarduser123@tcp(127.0.0.1:3306)/simpleboard?charset=utf8&parseTime=True&loc=Local"
var connectionString = "boarduser:boarduser123@tcp(go-todo-kido-db.cawucurfyhe9.ap-northeast-2.rds.amazonaws.com:3306)/simpleboard?charset=utf8&parseTime=True&loc=Local"

const maxOpenConn = 25
const maxIdleConn = 25
const maxConnLifetime = 10 * time.Minute

// get connect from db
func Connect() *sql.DB {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Fail to connect to db.", err)
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxConnLifetime)

	fmt.Println("InUsed: ", db.Stats().InUse)
	fmt.Println("Idle: ", db.Stats().Idle)
	fmt.Println("MaxOpenConnections: ", db.Stats().MaxOpenConnections)

	return db
}
