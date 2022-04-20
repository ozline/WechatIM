package api

import (
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
			ExpiresAt: int64(time.Now().Unix() + 3600), //签名过期时间 一小时后
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

func ChatPrivate(c *gin.Context) {
	var senderName = c.Request.Header.Get("username")
	var receiverName = c.Request.Header.Get("sendto")
	if senderName == "" || receiverName == "" {
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
	var senderName = c.Request.Header.Get("username")
	var receiverName = c.Request.Header.Get("sendto")
	if senderName == "" || receiverName == "" {
		global.UnifiedReturn(c, global.ErrorParams, "连接鉴权失败", nil, "")
		return
	}

	//升级为WebSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if !global.UnifiedErrorHandle(err, "升级WebSocket协议") {
		return
	}
	defer ws.Close()

	services.RabbitMQHandler(ws, senderName, receiverName, 2)
}
