package structs

import (
	"github.com/golang-jwt/jwt"
	"github.com/streadway/amqp"
)

type User struct {
	Id       string `json:"id" example:"-1" form:"id"`
	Username string `json:"username" example:""`
	Password string `json:"password" example:"" form:"password"`
	Status   string `json:"status" example:"0"`
	CreateAt string `json:"createat" example:"0"`
	Nickname string `json:"nickname" example:"" form:"nickname"`
	Phone    string `json:"phone" example:"" form:"phone"`
	Email    string `json:"email" example:"" form:"email"`
	Avatar   string `json:"avatar" example:""`
	Gender   string `json:"gender" example:"0" form:"gender"`
	Profile  string `jsong:"profile" example:"" form:"profile"`
}

type Conf struct {
	HttpPort string `yaml:"httpPort"`
	Mysql    struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Address  string `yaml:"ip"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		Table    struct {
			Users string `yaml:"users"`
			Test  string `yaml:"test"`
			// Group   string `yaml:"messages_group"`
			// Private string `yaml:"messages_private"`
		}
	}
	RabbitMQ struct {
		Address  string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	Redis struct {
		Address  string `yaml:"address"`
		Port     string `yamp:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	Admin struct {
		Secret string `yaml:"secret"`
		Name   string `yaml:"name"`
	}
}

type JWTClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Status   int    `json:"status" example:"0"`
	jwt.StandardClaims
}

type RabbitMQ struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	Queue     amqp.Queue
	QueueName string `example:""` //队列名称
	Exchange  string `example:""` //交换机名称
	Key       string `example:""` //Key
	Mqurl     string `example:""` //链接信息
}

type Message struct {
	Sender   string `json:"sender"`   //发送者的USERID
	Receiver string `json:"receiver"` //0=群发 其他=用户ID
	// Type   int    `json:"type" example:"0"`           //消息类型 0=私聊 1=群发
	Msg string `json:"msg" example:"emptyMessage"` //消息正文
}
