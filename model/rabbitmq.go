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

func RabbitMQCreateQueue(name string) (*amqp.Channel, amqp.Queue) {
	ch, err := Conn.Channel()
	global.UnifiedErrorHandle(err, "RabbitMQ创建通道")
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	global.UnifiedErrorHandle(err, "RabbitMQ声明队列")
	return ch, q
}

func RabbitMQPublish(ch *amqp.Channel, q amqp.Queue, body []byte) bool {
	err := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ发送消息")
}

func RabbitMQConsume(ch *amqp.Channel, q amqp.Queue) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	global.UnifiedErrorHandle(err, "RabbitMQ 注册Consume")
	return msgs
}
