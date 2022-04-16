package api

import (
	"net/http"
	"time"
	"wechat/conf"
	"wechat/global"
	"wechat/middleware"
	"wechat/structs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Ping(c *gin.Context) {
	token := c.Request.Header.Get("AuthToken")
	if token != "" {
		claims, _ := middleware.JWTParse(token)
		global.UnifiedReturn(c, global.Success, "ok", claims, "")
	} else {
		global.UnifiedReturn(c, global.Success, "ok", nil, "")
	}
}

func Chat(c *gin.Context) {
	//升级get请求为webSocket协议
	// go ChatProgress(c)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")

			//TODO:这部分似乎就是curd了
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func ChatProgress(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")

			//TODO:这部分似乎就是curd了
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
