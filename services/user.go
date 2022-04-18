package services

import (
	"errors"
	"wechat/conf"
	"wechat/middleware"
	"wechat/structs"

	_ "github.com/go-sql-driver/mysql"
)

//注册
func AddUser(info structs.User) (bool, error) {
	var tmplory string
	cmd := "SELECT password FROM `" + conf.Config.Mysql.Table.Users + "` where username = '" + info.Username + "'"
	err := middleware.DB.QueryRow(cmd).Scan(&tmplory)
	if err == nil {
		return false, errors.New("用户已经存在")
	}

	tx, err := middleware.DB.Begin()
	if err != nil {
		return false, err
	}
	cmd = "INSERT INTO `" + conf.Config.Mysql.DBName + "`.`" + conf.Config.Mysql.Table.Users + "` (`username`,`password`,`create_at`,`nickname`) VALUES (?, ?, ?, ?)"
	tmp, err := tx.Prepare(cmd)
	if err != nil {
		return false, err
	}
	_, err = tmp.Exec(info.Username, middleware.GenerateTokenSHA256(info.Password), middleware.GetTimestamp(), info.Username)
	if err != nil {
		return false, err
	}
	tx.Commit()
	return true, nil
}

//登录
func CheckUser(info structs.User) (string, int, error) {
	var password string
	var id string
	var status int
	cmd := "SELECT password,id,status FROM `" + conf.Config.Mysql.Table.Users + "` where username = '" + info.Username + "'"
	err := middleware.DB.QueryRow(cmd).Scan(&password, &id, &status)
	if err != nil {
		return "-1", 0, err
	}
	if middleware.GenerateTokenSHA256(info.Password) == password {
		return id, status, nil
	} else {
		return "-1", 0, errors.New("用户名或密码错误")
	}
}
