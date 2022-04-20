package api

import (
	"wechat/global"
	"wechat/middleware"
	"wechat/services"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	token := c.Request.Header.Get("AuthToken")
	if token != "" {
		claims, _ := middleware.JWTParse(token)
		global.UnifiedReturn(c, global.Success, global.MsgGeneral, claims, "")
	} else {
		global.UnifiedReturn(c, global.Success, global.MsgGeneral, nil, "")
	}
}

func GetUserID(c *gin.Context) string {
	return services.GetUserID(c.Request.Header.Get("AuthToken"))
}
