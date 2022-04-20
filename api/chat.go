package api

import (
	"net/http"
	"wechat/global"
	"wechat/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatPrivate(c *gin.Context) {
	var senderName = GetUserID(c)
	var receiverName = c.Param("userid")

	if senderName == "-1" || receiverName == "" || !services.CheckUserExist(receiverName) {
		global.UnifiedReturn(c, global.ErrorParams, "连接鉴权失败", nil, "")
		return
	}

	//升级为WebSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if !global.UnifiedErrorHandle(err, "升级WebSocket协议") {
		return
	}
	defer ws.Close()

	services.RabbitMQHandler(ws, senderName, receiverName, 1)
}

func ChatRooms(c *gin.Context) {
	var senderName = GetUserID(c)
	var receiverName = c.Param("roomid")

	receiver := services.GetRoomExchangeID(receiverName)

	if senderName == "" || receiverName == "" || receiver == "-1" {
		global.UnifiedReturn(c, global.ErrorParams, "连接鉴权失败", nil, "")
		return
	}

	//升级为WebSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if !global.UnifiedErrorHandle(err, "升级WebSocket协议") {
		return
	}
	defer ws.Close()

	services.RabbitMQHandler(ws, senderName, receiver, 2)
}
