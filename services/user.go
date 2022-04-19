package services

import (
	"errors"
	"fmt"
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

	cmd = "INSERT INTO `" + conf.Config.Mysql.DBName + "`.`" + conf.Config.Mysql.Table.Users + "` (`username`,`password`,`create_at`,`nickname`) VALUES (?, ?, ?, ?)"

	return middleware.DBCommit(cmd, info.Username, middleware.GenerateTokenSHA256(info.Password), middleware.GetTimestamp13(), info.Username)
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

func RedisDB_UserOnline(userid string) bool {
	tmp := make(map[string]interface{})
	tmp[userid] = fmt.Sprint(middleware.GetTimestamp13())
	return middleware.RedisDBHMSet(conf.Config.Redis.Key.Users, tmp)
}

func RedisDB_UserExit(userid string) bool {
	return middleware.RedisDBHDel(conf.Config.Redis.Key.Users, userid)
}

func RedisDB_CheckUser(userid string) (bool, string) {
	result := middleware.RedisDBHexists(conf.Config.Redis.Key.Users, userid)
	if result {
		return true, middleware.RedisDBHGet(conf.Config.Redis.Key.Users, userid)
	} else {
		return false, ""
	}
}
