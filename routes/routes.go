package routes

import (
	// "wechat/api"
	// "wechat/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	// OutAuth := router.Group("/api/")
	{
		// OutAuth.GET("/ping", api.UserTest)                //测试通信
		// OutAuth.POST("/user/login", api.UserLogin)           //登录
		// OutAuth.POST("/user/register", api.UserRegister)     //注册
	}

	return router
}
