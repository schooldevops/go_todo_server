package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var connectionString = "boarduser:boarduser123@tcp(127.0.0.1:3306)/simpleboard?charset=utf8&parseTime=True&loc=Local"

const maxOpenConn = 25
const maxIdleConn = 25
const maxConnLifetime = 10 * time.Minute

// Connect get connection to db
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

// ConnectGormDB get connection to db using grom
func ConnectGormDB() *gorm.DB {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Fail to connect to db.", err)
		panic(err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		fmt.Println("Fail to open from GORM.", err)
		panic(err)
	}

	return gormDB
}
