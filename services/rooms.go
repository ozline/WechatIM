package services

import (
	"wechat/conf"
	"wechat/global"
	"wechat/middleware"
	"wechat/structs"
)

func RoomCreate(room structs.Room) (bool, error) {
	cmd := "INSERT INTO `" + conf.Config.Mysql.DBName + "`.`" + conf.Config.Mysql.Table.Rooms + "` (`name`,`owner`,`create_at`,`exchange_id`,`status`) VALUES (?, ?, ?, ?, ?)"

	return middleware.DBCommit(cmd, room.Name, room.Owner, middleware.GetTimestamp13(), middleware.GenerateTokenSHA256("Room="+room.Name+room.Owner), 0)
}

func GetRoomExchangeID(roomid string) string {
	var room string

	cmd := "SELECT exchange_id FROM `" + conf.Config.Mysql.Table.Rooms + "`"
	cmd += " WHERE `roomid`=" + roomid
	err := middleware.DB.QueryRow(cmd).Scan(&room)
	if !global.UnifiedErrorHandle(err, "获取房间ExchangeID出错") {
		room = "-1"
	}
	return room
}

func GetRoomOwner(roomid string) string {
	var room string

	cmd := "SELECT owner FROM `" + conf.Config.Mysql.Table.Rooms + "`"
	cmd += " WHERE `roomid`=" + roomid
	err := middleware.DB.QueryRow(cmd).Scan(&room)
	if !global.UnifiedErrorHandle(err, "获取房间主人信息出错") {
		room = "-1"
	}
	return room
}

func RoomDelete(roomid string) (bool, error) {

	cmd := "DELETE FROM `" + conf.Config.Mysql.DBName + "`.`" + conf.Config.Mysql.Table.Rooms + "` WHERE `roomid`=" + roomid

	return middleware.DBCommit(cmd)
}

func RoomSubscribe(userid string, roomid string) (bool, error) {
	cmd := "INSERT INTO `" + conf.Config.Mysql.DBName + "`.`" + conf.Config.Mysql.Table.Subscribes + "` (`userid`,`roomid`,`create_at`) VALUES (?, ?, ?)"

	return middleware.DBCommit(cmd, userid, roomid, middleware.GetTimestamp13())
}

func RooomUnSubscribe(userid string, roomid string) (bool, error) {
	cmd := "DELETE FROM `" + conf.Config.Mysql.DBName + "`.`" + conf.Config.Mysql.Table.Subscribes + "` WHERE `roomid`=" + roomid + " AND `userid`=" + userid

	global.UnifiedPrintln(cmd, nil)
	return middleware.DBCommit(cmd)
}

func RoomGetSubscribe(userid string) {

	//TODO:照搬一下之前的代码，属于Ctrl CV事件了
}
