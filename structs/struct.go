package structs

import(
	"github.com/golang-jwt/jwt"
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
			Group string `yaml:"messages_group"`
			Private string `yaml:"messages_private"`
		}
	}
	Admin struct{
		Secret string `yaml:"secret"`
		Name string `yaml:"name"`
	}
}

type JWTClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Status   int    `json:"status" example:"0"`
	jwt.StandardClaims
}