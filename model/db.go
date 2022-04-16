package model

import (
	"database/sql"
	"fmt"
	"strings"
	"wechat/conf"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DBInit() bool { //连接RDS
	path := strings.Join([]string{
		conf.Config.Mysql.Username, ":",
		conf.Config.Mysql.Password, "@tcp(",
		conf.Config.Mysql.Address, ":",
		conf.Config.Mysql.Port, ")/",
		conf.Config.Mysql.DBName, "?charset=utf8&parseTime=True"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("Error: DBInit ", err)
	}
	fmt.Println("数据库连接成功")
	return true
}

func DBDestory() {
	if err := DB.Close(); err != nil {
		fmt.Println("Error: DBDestory ", err)
	}
}

//获取数量，where表示限定条件，不需要再写WHERE，但是多条件要加AND
func DBGetCount(table string, where string) (int, error) {
	cmd := "SELECT count(*) from `" + conf.Config.Mysql.DBName + "`.`" + table + "`"
	cmd += " WHERE " + where
	var count = 0
	err := DB.QueryRow(cmd).Scan(&count)
	if err != nil {
		return -1, err
	} else {
		return count, nil
	}
}

//学习一下可变参数
func DBCommit(statement string, args ...interface{}) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	tmp, err := tx.Prepare(statement)
	if err != nil {
		return false, err
	}
	_, err = tmp.Exec(args)
	if err != nil {
		return false, err
	}
	tx.Commit()
	return true, nil
}
