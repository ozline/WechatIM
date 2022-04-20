package services

import (
	"time"
	"wechat/conf"
	"wechat/middleware"
	"wechat/structs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetAuthToken(id string, username string, status int) (string, error) {
	claims := structs.JWTClaims{
		Id:       id,
		Username: username,
		Status:   status,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 3600), //签名生效时间 一小时前
			ExpiresAt: int64(time.Now().Unix() + 3600), //签名过期时间 一小时后
			Issuer:    conf.Config.Admin.Name,
		},
	}
	token, err := middleware.JWTGenerate(claims)
	return token, err
}

func UpdateAuthToken(c *gin.Context) string {
	before := c.Request.Header.Get("AuthToken")
	if before == "" {
		return ""
	}
	claims, err := middleware.JWTParse(before)
	if err != nil {
		return ""
	}
	token, err := GetAuthToken(claims.Id, claims.Username, claims.Status)
	if err != nil {
		return ""
	} else {
		return token
	}
}
