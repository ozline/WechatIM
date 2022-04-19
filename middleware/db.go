package middleware

import (
	"database/sql"
	"strings"
	"wechat/conf"
	"wechat/global"

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
	err := DB.Ping()

	return global.UnifiedErrorHandle(err, "DB ping")
}

func DBDestory() {
	global.UnifiedErrorHandle(DB.Close(), "DB Destory")
}

//获取数量，where表示限定条件，不需要再写WHERE，但是多条件要加AND
func DBGetCount(table string, where string) (int, error) {
	cmd := "SELECT count(*) from `" + conf.Config.Mysql.DBName + "`.`" + table + "`"
	if where != "" {
		cmd += " WHERE " + where
	}

	var count = 0
	err := DB.QueryRow(cmd).Scan(&count)
	if global.UnifiedErrorHandle(err, "DB 查询行数") {
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
	_, err = tmp.Exec(args...)
	if err != nil {
		return false, err
	}
	tx.Commit()
	return true, nil
}
