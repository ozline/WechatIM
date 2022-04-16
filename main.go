package main

import (
	"fmt"
	"os"
	"wechat/conf"
	"wechat/model"
	"wechat/routes"
)

func main() {
	//读取配置
	conf.Init()
	var HTTP_PORT string
	HTTP_PORT = os.Getenv("HTTP_PORT")
	if HTTP_PORT == "" {
		HTTP_PORT = conf.Config.HttpPort
	}

	//连接数据库和RabbitMQ
	if model.DBInit() && model.RabbitMQInit() {
		r := routes.NewRouter()
		_ = r.Run("0.0.0.0:" + HTTP_PORT)
	} else {
		fmt.Println("数据库或RabbitMQ连接失败")
	}
}
