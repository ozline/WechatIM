package services

import (
	"encoding/json"
	"time"
	"wechat/conf"
	"wechat/global"
	"wechat/middleware"
	"wechat/structs"

	"github.com/gorilla/websocket"
)

func RabbitMQHandler(ws *websocket.Conn, sender string, receiver string, mode int) {
	var err error
	var forever chan struct{}
	var rabbit structs.RabbitMQ

	//设定交换机模式信息
	if mode == 1 { //私聊
		rabbit.Kind.Sender = middleware.RabbitMQ_Direct
		rabbit.Kind.Receiver = middleware.RabbitMQ_Direct
	} else if mode == 2 { //群聊
		rabbit.Kind.Sender = middleware.RabbitMQ_Direct
		rabbit.Kind.Receiver = middleware.RabbitMQ_Fanout
	} else {
		return
	}

	//设定交换机ID信息
	if mode == 1 { //私聊
		rabbit.Exchange.Sender = middleware.GenerateTokenSHA256(sender)
		rabbit.Exchange.Receiver = middleware.GenerateTokenSHA256(receiver)
	} else if mode == 2 { //群聊
		rabbit.Exchange.Sender = middleware.GenerateTokenSHA256(sender)
		rabbit.Exchange.Receiver = receiver //数据库中已存SHA256结果
	} else {
		return
	}

	//设定交换机Key信息
	if mode == 1 { //私聊
		rabbit.Key.Sender = middleware.GenerateTokenSHA256(sender)
		rabbit.Key.Receiver = middleware.GenerateTokenSHA256(receiver)
	} else if mode == 2 { //群聊
		rabbit.Key.Sender = ""
		rabbit.Key.Receiver = ""
	} else {
		return
	}

	//在RedisDB中令用户上线
	if !RedisDB_UserOnline(sender) {
		return
	}

	//创建信道
	rabbit.Channel = middleware.RabbitMQCreateChannel()
	defer rabbit.Channel.Close()

	//创建交换机
	result := middleware.RabbitMQCreateExchange(rabbit.Channel, rabbit.Exchange.Sender, middleware.RabbitMQ_Direct) && middleware.RabbitMQCreateExchange(rabbit.Channel, rabbit.Exchange.Receiver, middleware.RabbitMQ_Direct)

	if !result {
		global.UnifiedPrintln("RabbitMQ创建交换机失败", nil)
		return
	}

	//创建随机Queue并绑定至Sender所属的交换机
	rabbit.Queue = middleware.RabbitMQCreateQueue(rabbit.Channel, "")
	result = middleware.RabbitMQQueueBind(rabbit.Channel, rabbit.Queue, rabbit.Exchange.Sender, rabbit.Key.Sender)

	//绑定对方交换机到我们交换机，获取对方传来的信息
	result = middleware.RabbitMQExchangeBind(rabbit.Channel, rabbit.Exchange.Sender, rabbit.Exchange.Receiver, rabbit.Key.Sender)
	if !result {
		global.UnifiedPrintln("RabbitMQ绑定交换机队列失败", nil)
	}

	//开始消费消息
	msgs := middleware.RabbitMQConsume(rabbit.Channel, rabbit.Queue)

	unSolvedMsg := middleware.RedisDSBHGetAll(receiver + "->" + sender)
	for d := range unSolvedMsg {
		ws.WriteMessage(websocket.TextMessage, []byte([]byte(unSolvedMsg[d])))
	}
	if len(unSolvedMsg) != 0 {
		middleware.RedisDBDel(receiver + "->" + sender)
	}

	//Websocket接收消息，推送至RabbitMQ
	go func() {
		for {
			//服务端心跳检测
			//设置链接限时，10分钟(600秒)没有操作则进入销毁
			ws.SetReadDeadline(time.Now().Add(time.Duration(600) * time.Second))
			messageType, message, err := ws.ReadMessage()

			if messageType == -1 { //退出
				RedisDB_UserExit(sender)
				ws.Close()
				return
			}
			if !global.UnifiedErrorHandle(err, "WebSocket ReadMessage") {
				return
			}

			//客户端心跳检测使用
			if string(message) == "ping" {
				ws.WriteMessage(websocket.TextMessage, []byte("pong"))
				continue
			}

			msgJson, err := json.Marshal(structs.Message{
				Sender:   sender,
				Receiver: receiver,
				Msg:      string(message),
				SendTime: middleware.GetTimestamp13(),
			})

			var res bool

			if mode == 1 && !middleware.RedisDBHexists(conf.Config.Redis.Key.Users, receiver) { //私聊模式，且对方不在线，则推送至Redis

				go RedisDB_SetMessage(sender+"->"+receiver, string(msgJson))

			} else { //群聊模式，或对方在线，直接通过RabbitMQ推送
				res = middleware.RabbitMQExchangePublish(rabbit.Channel, rabbit.Exchange.Receiver, msgJson, rabbit.Key.Receiver)
			}

			// res := global.UnifiedErrorHandle(err, "SolveJSON")

			if !res {
				return
			}
		}
	}()

	//接收RabbitMQ消息，推送至Websocket
	go func() {
		for d := range msgs {
			ws.SetWriteDeadline(time.Now().Add(time.Duration(600) * time.Second))
			//构建推送Json

			var msg structs.Message
			err = json.Unmarshal(d.Body, &msg)
			if !global.UnifiedErrorHandle(err, "SolveJSON") {
				return
			}

			if msg.Sender == sender { //发送者是我自己，不推给前端
				continue
			}

			//推送消息
			err = ws.WriteMessage(websocket.TextMessage, []byte(d.Body))
			if !global.UnifiedErrorHandle(err, "Websocket读取消息") {
				return
			}
		}
	}()

	<-forever
}
