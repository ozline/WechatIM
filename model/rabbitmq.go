package model

import (
	"fmt"
	"strings"
	"wechat/conf"

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
	fmt.Println(path)
	Conn, err = amqp.Dial(path)
	if err != nil {
		return false
	}
	return true
}
