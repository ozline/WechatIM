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
		OutAuth.GET("/ping", api.Ping)                   //连通性测试
		OutAuth.POST("/user/login", api.UserLogin)       //登录
		OutAuth.POST("/user/register", api.UserRegister) //注册
	}
	Auth := router.Group("/api/auth/")
	{
		Auth.Use(middleware.JWTAuth())

		Auth.GET("/ping", api.Ping)                        //连通性测试
		Auth.GET("/chat/private/:userid", api.ChatPrivate) //私聊
		Auth.GET("/chat/rooms/:roomid", api.ChatRooms)     //群聊

		Auth.POST("/room", api.RoomCreate)                     //创建房间
		Auth.DELETE("/room/:roomid", api.RoomdDelete)          //删除房间
		Auth.GET("/subscribe", api.RoomGetSubscribe)           //获取订阅
		Auth.POST("/subscribe/:roomid", api.RoomSubscribe)     //订阅房间
		Auth.DELETE("/subscribe/:roomid", api.RoomUnSubscribe) //取消订阅
	}

	return router
}
