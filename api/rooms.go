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
		return
	}

	if room.Name == "" {
		global.UnifiedReturn(c, global.ErrorParams, nil, nil, "")
		return
	}

	room.Owner = GetUserID(c)

	if room.Owner == "-1" {
		global.UnifiedReturn(c, global.ErrorParams, global.ErrorJWTChcek, nil, "")
		return
	}

	result, err := services.RoomCreate(room)
	if result && global.UnifiedErrorHandle(err, "创建房间") {
		global.UnifiedReturn(c, global.Success, nil, nil, services.UpdateAuthToken(c))
	} else {
		global.UnifiedReturn(c, global.ErrorDatabase, err.Error(), nil, services.UpdateAuthToken(c))
	}
}

func RoomdDelete(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := GetUserID(c)
	if services.GetRoomOwner(roomid) != userid {
		global.UnifiedReturn(c, global.ErrorPermission, "非房间主人", nil, services.UpdateAuthToken(c))
		return
	}

	result, err := services.RoomDelete(roomid)
	global.UnifiedErrorHandle(err, "DB删除Room")
	if result {
		global.UnifiedReturn(c, global.Success, nil, nil, services.UpdateAuthToken(c))
	} else {
		global.UnifiedReturn(c, global.ErrorDatabase, global.ErrorDefault, nil, services.UpdateAuthToken(c))
	}
}

func RoomSubscribe(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := GetUserID(c)

	//做到房间必须是已经存在的，否则不加入
	if services.GetRoomExchangeID(roomid) == "-1" || userid == "-1" {
		global.UnifiedReturn(c, global.ErrorPermission, nil, nil, services.UpdateAuthToken(c))
		return
	}

	result, err := services.RoomSubscribe(userid, roomid)
	global.UnifiedErrorHandle(err, "DB添加订阅")
	if result {
		global.UnifiedReturn(c, global.Success, nil, nil, services.UpdateAuthToken(c))
	} else {
		global.UnifiedReturn(c, global.ErrorDatabase, global.ErrorDefault, nil, services.UpdateAuthToken(c))
	}

}

func RoomUnSubscribe(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := GetUserID(c)

	if services.GetRoomExchangeID(roomid) == "-1" || userid == "-1" {
		global.UnifiedReturn(c, global.ErrorPermission, nil, nil, services.UpdateAuthToken(c))
		return
	}

	result, err := services.RooomUnSubscribe(userid, roomid)
	global.UnifiedErrorHandle(err, "DB删除订阅")
	if result {
		global.UnifiedReturn(c, global.Success, nil, nil, services.UpdateAuthToken(c))
	} else {
		global.UnifiedReturn(c, global.ErrorDatabase, global.ErrorDefault, nil, services.UpdateAuthToken(c))
	}
}

func RoomGetSubscribe(c *gin.Context) {

}
