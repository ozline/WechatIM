package services

import (
	"wechat/conf"
	"wechat/middleware"
	"wechat/structs"
)

func RoomCreate(room structs.Room) (bool, error) {
	cmd := "INSERT INTO `" + conf.Config.Mysql.DBName + "`.`" + conf.Config.Mysql.Table.Rooms + "` (`name`,`owner`,`create_at`,`exchange_id`) VALUES (?, ?, ?)"

	return middleware.DBCommit(cmd, room.Name, room.Owner, middleware.GetTimestamp13(), middleware.GenerateTokenSHA256(room.Name+room.Owner))
}

// func RoomGetAll(user structs.User)

func RoomDelete(roomid string) (bool, error) {
	return true, nil
}

// func RoomAddMember(user structs.User)

// func RoomDeleteMember(user structs.User)
