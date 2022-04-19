package services

import (
	"fmt"
	"wechat/middleware"
)

func RedisDB_SetMessage(key string, msg string) bool {
	tmp := make(map[string]interface{})
	tmp[fmt.Sprint(middleware.GetTimestamp13())] = msg
	return middleware.RedisDBHMSet(key, tmp)
}
