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

	return global.UnifiedErrorHandle(err, "RedisDB连接")
}

func RedisDBHMSet(key string, fields map[string]interface{}) bool {
	result, err := RedisDB.HMSet(key, fields).Result()
	return global.UnifiedErrorHandle(err, "RedisDB 哈希插入 return:"+result) && result == "OK"
}

func RedisDBHDel(key string, fields string) bool {
	result, err := RedisDB.HDel(key, fields).Result()
	return global.UnifiedErrorHandle(err, "RedisDB 删除 Return:"+fmt.Sprint(result))
}

func RedisDBHexists(key string, fields string) bool {
	result, err := RedisDB.HExists(key, fields).Result()
	if global.UnifiedErrorHandle(err, "RedisDB 查询fields Return:"+fmt.Sprint(result)) {
		return result
	} else {
		return false
	}
}

func RedisDBHGet(key string, fields string) string {
	result, err := RedisDB.HGet(key, fields).Result()
	res := global.UnifiedErrorHandle(err, "RedisDB HGET Return:"+fmt.Sprint(result))
	if res {
		return result
	} else {
		return ""
	}
}

func RedisDSBHGetAll(key string) map[string]string {
	result, err := RedisDB.HGetAll(key).Result()
	global.UnifiedErrorHandle(err, "RedisDB HGetAll Return:"+fmt.Sprint(result))
	return result
}

func RedisDBDel(key string) int64 {
	result, err := RedisDB.Del(key).Result()
	global.UnifiedErrorHandle(err, "RedisDB DelKey Return:"+fmt.Sprint(result))
	return result
}
