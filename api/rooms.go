package api

import (
	"wechat/global"
	"wechat/services"
	"wechat/structs"

	"github.com/gin-gonic/gin"
)

func RoomCreate(c *gin.Context) {
	var room structs.Room
	err := c.ShouldBind(&room)
	if !global.UnifiedErrorHandle(err, "Gin绑定表单") {
		global.UnifiedReturn(c, global.ErrorParams, err.Error(), nil, "")
	}

	result, err := services.RoomCreate(room)
	if result && global.UnifiedErrorHandle(err, "创建房间") {
		global.UnifiedReturn(c, global.Success, nil, nil, UpdateAuthToken(c))
	} else {
		global.UnifiedReturn(c, global.ErrorDatabase, nil, nil, UpdateAuthToken(c))
	}
}

func RoomdDelete(c *gin.Context) {
	// var room structs.Room
}
