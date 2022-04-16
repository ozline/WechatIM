package model

import (
	"strings"
	"wechat/conf"
	"wechat/global"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection

func RabbitMQInit() bool {
	var err error
	path := strings.Join([]string{
		"amqp://",
		conf.Config.RabbitMQ.Username, ":",
		conf.Config.RabbitMQ.Password, "@",
		conf.Config.RabbitMQ.Address, "/",
	}, "")
	Conn, err = amqp.Dial(path)
	return global.UnifiedErrorHandle(err, "RabbitMQ连接")
}
