package api

import (
	"encoding/json"
	"net/http"
	"time"
	"wechat/conf"
	"wechat/global"
	"wechat/middleware"
	"wechat/services"
	"wechat/structs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetAuthToken(id string, username string, status int) (string, error) {
	claims := structs.JWTClaims{
		Id:       id,
		Username: username,
		Status:   status,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 3600), //签名生效时间 一小时前
			ExpiresAt: int64(time.Now().Unix() + 3600), //签名过期时间 按一小时算吧
			Issuer:    conf.Config.Admin.Name,
		},
	}
	token, err := middleware.JWTGenerate(claims)
	return token, err
}

func UpdateAuthToken(c *gin.Context) string {
	before := c.Request.Header.Get("AuthToken")
	if before == "" {
		return ""
	}
	claims, err := middleware.JWTParse(before)
	if err != nil {
		return ""
	}
	token, err := GetAuthToken(claims.Id, claims.Username, claims.Status)
	if err != nil {
		return ""
	} else {
		return token
	}
}

func Ping(c *gin.Context) {
	token := c.Request.Header.Get("AuthToken")
	if token != "" {
		claims, _ := middleware.JWTParse(token)
		global.UnifiedReturn(c, global.Success, global.MsgGeneral, claims, "")
	} else {
		global.UnifiedReturn(c, global.Success, global.MsgGeneral, nil, "")
	}
}

func Chat(c *gin.Context) {
	var senderName = c.Request.Header.Get("username")
	var receiverName = c.Request.Header.Get("sendto")
	if senderName == "" || receiverName == "" {
		global.UnifiedReturn(c, global.ErrorParams, "连接鉴权失败", nil, "")
		return
	}
	//升级为webSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if !global.UnifiedErrorHandle(err, "升级WebSocket协议") {
		return
	}
	defer ws.Close()

	RabbitMQHandler(ws, senderName, receiverName)
}

func RabbitMQHandler(ws *websocket.Conn, senderName string, receiverName string) {
	var err error
	var forever chan struct{}
	var rabbit structs.RabbitMQ

	//在RedisDB中令用户上线
	if !services.RedisDB_UserOnline(senderName) {
		return
	}

	//测试用
	rabbit.Exchange.Sender = middleware.GenerateTokenSHA256(senderName)
	rabbit.Exchange.Receiver = middleware.GenerateTokenSHA256(receiverName)

	//创建信道
	rabbit.Channel = middleware.RabbitMQCreateChannel()
	defer rabbit.Channel.Close()

	//创建交换机
	result := middleware.RabbitMQCreateExchange(rabbit.Channel, rabbit.Exchange.Sender, middleware.RabbitMQ_Direct) && middleware.RabbitMQCreateExchange(rabbit.Channel, rabbit.Exchange.Receiver, middleware.RabbitMQ_Direct)

	if !result {
		global.UnifiedPrintln("RabbitMQ创建交换机失败", nil)
		return
	}

	//创建随机Queue并绑定至本地交换机 - 用于消费消息
	rabbit.Queue = middleware.RabbitMQCreateQueue(rabbit.Channel, "")
	result = middleware.RabbitMQQueueBind(rabbit.Channel, rabbit.Queue, rabbit.Exchange.Sender, rabbit.Exchange.Sender)

	//绑定对方交换机到我们交换机，获取对方传来的信息
	result = middleware.RabbitMQExchangeBind(rabbit.Channel, rabbit.Exchange.Sender, rabbit.Exchange.Receiver, rabbit.Exchange.Sender)
	if !result {
		global.UnifiedPrintln("RabbitMQ绑定交换机队列失败", nil)
	}

	//创建消费者
	msgs := middleware.RabbitMQConsume(rabbit.Channel, rabbit.Queue)

	//Websocket接收消息，推送至RabbitMQ
	go func() {
		for {
			//服务端心跳检测使用
			//设置链接限时，10分钟(600秒)没有操作则进入销毁
			ws.SetReadDeadline(time.Now().Add(time.Duration(600) * time.Second))
			messageType, message, err := ws.ReadMessage()

			if messageType == -1 { //侦测到退出
				services.RedisDB_UserExit(senderName)
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

			tmp, err := json.Marshal(structs.Message{
				Sender:   senderName,
				Receiver: receiverName,
				Msg:      string(message),
			})

			res := global.UnifiedErrorHandle(err, "Websocket发送消息") && middleware.RabbitMQExchangePublish(rabbit.Channel, rabbit.Exchange.Receiver, tmp, rabbit.Exchange.Receiver)
			if !res {
				return
			}
		}
	}()

	//接收RabbitMQ消息，推送至Websocket
	go func() {
		for d := range msgs {
			//构建推送Json

			var msg structs.Message
			err = json.Unmarshal(d.Body, &msg)
			if !global.UnifiedErrorHandle(err, "生成发送JSON") {
				return
			}

			if msg.Sender == senderName { //发送者是我自己，不接收
				continue
			}
			//推送消息
			err = ws.WriteMessage(websocket.TextMessage, []byte(msg.Msg))
			if !global.UnifiedErrorHandle(err, "Websocket读取消息") {
				return
			}
		}
	}()

	<-forever
}
