package structs

import (
	"github.com/golang-jwt/jwt"
	"github.com/streadway/amqp"
)

type User struct {
	Id       string `json:"id" example:"-1" form:"id"`
	Username string `json:"username" example:"" form:"username"`
	Password string `json:"password" example:"" form:"password"`
	Status   string `json:"status" example:"0"`
	CreateAt string `json:"createat" example:"0"`
	Nickname string `json:"nickname" example:"" form:"nickname"`
	Phone    string `json:"phone" example:"" form:"phone" form:"phone"`
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
			Rooms string `yaml:"rooms"`
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
		Key      struct {
			Users string `yaml:"users"`
		}
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
	Conn     *amqp.Connection
	Channel  *amqp.Channel
	Queue    amqp.Queue
	Exchange struct { //交换机ID
		Sender   string
		Receiver string
	}
	Kind struct { //交换机种类
		Sender   string
		Receiver string
	}
	Key struct { //交换机Key
		Sender   string `example:""`
		Receiver string `exmple:""`
	}
}

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Msg      string `json:"msg" example:"emptyMessage"` //消息正文
	SendTime int64  `json:"sendtime"`
}

type Room struct {
	Roomid       string
	Name         string `form:"name"`
	Owner        string `form:"owner"`
	Status       string `form:"status"`
	ExchangeName string //交换机标识
	CreateAt     string //创建时间,13位时间戳
}

type Test struct {
	Count string
}
