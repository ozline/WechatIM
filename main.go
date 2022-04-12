package main

import(
	"wechat/conf"
	"wechat/model"
	"wechat/routes"
	"fmt"
	"os"
)



func main(){
	conf.Init()
	var HTTP_PORT string
	HTTP_PORT = os.Getenv("HTTP_PORT")
	if HTTP_PORT == "" {
		HTTP_PORT = conf.Config.HttpPort
	}
	if model.DBInit() {
		r := routes.NewRouter()
		_ = r.Run("0.0.0.0:" + HTTP_PORT)
	} else {
		fmt.Println("启动失败")
	}
}