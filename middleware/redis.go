package middleware

import (
	"fmt"
	"wechat/conf"
	"wechat/global"

	"github.com/go-redis/redis"
)

var RedisDB *redis.Client

func RedisDBInit() bool {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     conf.Config.Redis.Address + ":" + conf.Config.Redis.Port,
		Password: conf.Config.Redis.Username + ":" + conf.Config.Redis.Password,
		DB:       0,
	})
	_, err := RedisDB.Ping().Result()

	tmp, _ := RedisDB.HGetAll("test123456").Result()

	for d := range tmp {
		global.UnifiedPrintln(d+" = "+tmp[d], nil)
	}

	return global.UnifiedErrorHandle(err, "RedisDB连接")
}

func RedisDBHMSet(key string, fields map[string]interface{}) bool {
	result, err := RedisDB.HMSet(key, fields).Result()
	return global.UnifiedErrorHandle(err, "RedisDB 哈希插入 return:"+result) && result == "OK"
}

func RedisDBHMSet_Message(key string, msg string) bool {
	tmp2 := make(map[string]interface{})
	tmp2[fmt.Sprint(GetTimestamp13())] = msg
	return RedisDBHMSet(key, tmp2)
}
