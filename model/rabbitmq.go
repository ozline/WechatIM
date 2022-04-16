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

func RabbitMQCreateChannel() *amqp.Channel {
	ch, err := Conn.Channel()
	global.UnifiedErrorHandle(err, "RabbitMQ创建通道")
	return ch
}

func RabbitMQCreateQueue(ch *amqp.Channel, name string) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		true,
		false,
		nil,
	)
	global.UnifiedErrorHandle(err, "RabbitMQ声明队列")
	return q
}

func RabbitMQCreateExchange(ch *amqp.Channel, name string) bool {
	err := ch.ExchangeDeclare(
		name,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ声明交换机")
}

func RabbitMQQueueBind(ch *amqp.Channel, q amqp.Queue, exchange string) bool {
	err := ch.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ队列绑定交换机")
}

func RabbitMQQueuePublish(ch *amqp.Channel, q amqp.Queue, body []byte) bool {
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
	return global.UnifiedErrorHandle(err, "RabbitMQ发送消息_Queue")
}

func RabbitMQExchangePublish(ch *amqp.Channel, exchange string, body []byte) bool {
	err := ch.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ发送消息_Exchange")
}

func RabbitMQConsume(ch *amqp.Channel) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		"",
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
