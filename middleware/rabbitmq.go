package middleware

import (
	"strings"
	"wechat/conf"
	"wechat/global"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection

const (
	RabbitMQ_Fanout = "fanout"
	RabbitMQ_Direct = "direct"
)

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

func RabbitMQCreateExchange(ch *amqp.Channel, name string, kind string) bool {
	//kind 有 fanout:群发 和 direct: 定向
	err := ch.ExchangeDeclare(
		name,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ声明交换机")
}

func RabbitMQQueueBind(ch *amqp.Channel, q amqp.Queue, exchange string, key string) bool {
	err := ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil,
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ队列绑定交换机")
}

func RabbitMQExchangeBind(ch *amqp.Channel, destination string, source string, key string) bool {
	err := ch.ExchangeBind(
		destination,
		key,
		source,
		true,
		nil,
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ交换机绑定交换机")
}

func RabbitMQExchangePublish(ch *amqp.Channel, exchange string, body []byte, key string) bool {
	err := ch.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	return global.UnifiedErrorHandle(err, "RabbitMQ发送消息_Exchange")
}

func RabbitMQConsume(ch *amqp.Channel, queue amqp.Queue) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		queue.Name,
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
