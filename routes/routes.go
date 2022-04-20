package routes

import (
	"wechat/api"
	"wechat/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	OutAuth := router.Group("/api/")
	{
		OutAuth.GET("/ping", api.Ping) //测试通信
		OutAuth.GET("/chat/private", api.ChatPrivate)
		OutAuth.GET("/chat/rooms", api.ChatRooms)
		OutAuth.POST("/user/login", api.UserLogin)       //登录
		OutAuth.POST("/user/register", api.UserRegister) //注册
	}
	Auth := router.Group("/api/auth/")
	{
		Auth.Use(middleware.JWTAuth())

		Auth.GET("/ping", api.Ping)
		Auth.GET("/chat", api.ChatPrivate)
	}

	return router
}
