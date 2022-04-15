package middleware

import (
	"wechat/conf"
	"wechat/structs"
	"wechat/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("AuthToken")
		if token == "" {
			global.UnifiedReturn(c,global.ErrorJWTCheck,"AuthToken验证失败",nil,"")
			c.Abort()
			return
		}
		claims, err := JWTParse(token)
		if err != nil || claims.ExpiresAt < time.Now().Unix() {
			global.UnifiedReturn(c,global.ErrorJWTCheck,"AuthToken验证失败",nil,"")
			c.Abort()
			return
		}

		c.Next()
	}
}

func JWTParse(token string) (*structs.JWTClaims, error) { //解析JWT
	tmp, err := jwt.ParseWithClaims(token, &structs.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Config.Admin.Secret), nil
	})
	if tmp != nil {
		claims, result := tmp.Claims.(*structs.JWTClaims)
		if result && tmp.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func JWTGenerate(claims structs.JWTClaims) (string, error) { //生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.Config.Admin.Secret))
}