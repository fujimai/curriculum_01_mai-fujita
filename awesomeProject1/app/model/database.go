package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var Db *sql.DB

func init() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PWD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))

	// データベースと接続。
	Db, err = sql.Open("mysql", dsn)
	/*
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPwd := os.Getenv("MYSQL_PWD")
		mysqlHost := os.Getenv("MYSQL_HOST")
		mysqlDatabase := os.Getenv("MYSQL_DATABASE")

		connStr := fmt.Sprintf("%s:%s@(localhost:8080)%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
		db, err := sql.Open("mysql", connStr)

		// ①-2 接続時のパラメータを指定
		_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlPwd, mysqlDatabase))
		if err != nil {
			log.Fatalf("fail: sql.Open, %v\n", err)
		}
		// ①-3　接続できているか確認、接続情報に誤りがあればエラーになる
		if err := _db.Ping(); err != nil {
			log.Fatalf("fail: _db.Ping, %v\n", err)
		}
		_ = _db
	*/
	// エラーハンドリング省略
	if err != nil {
		fmt.Println(err)
		return
	}
	// 疎通確認を行う。
	err = Db.Ping()

	// エラーハンドリング省略
	if err != nil {
		fmt.Println(err)
		return
	}
	/*// database.goがimportされたらinit関数が走り、このSQLが実行される。
	sql := `CREATE TABLE IF NOT EXISTS todos(
			id varchar(26) not null,
			name varchar(100) not null,
			status varchar(100) not null
		)`

	_, err = Db.Exec(sql)
	*/
	// エラーハンドリング省略
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connection has been established!")
}
