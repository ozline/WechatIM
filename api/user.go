package api

import (
	"wechat/global"
	"wechat/services"
	"wechat/structs"

	"github.com/gin-gonic/gin"
)

//注册
func UserRegister(c *gin.Context) {
	var info structs.User
	err := c.ShouldBind(&info)
	if !global.UnifiedErrorHandle(err, "Gin绑定表单") {
		global.UnifiedReturn(c, global.ErrorParams, err.Error(), nil, "")
	}

	result, err := services.AddUser(info)
	if err != nil {
		global.UnifiedReturn(c, global.ErrorUsers, err.Error(), nil, "")
	} else {
		global.UnifiedReturn(c, global.Success, result, nil, "")
	}
}

//登录
func UserLogin(c *gin.Context) {
	var info structs.User
	err := c.ShouldBind(&info)
	if !global.UnifiedErrorHandle(err, "Gin绑定表单") {
		global.UnifiedReturn(c, global.ErrorParams, err.Error(), nil, "")
	}

	if len(info.Username) == 0 || len(info.Password) == 0 {
		global.UnifiedReturn(c, global.ErrorUsers, "用户名或密码错误", nil, "")
		return
	}
	id, isadmin, err := services.CheckUser(info)
	if err != nil {
		global.UnifiedReturn(c, global.ErrorUsers, err.Error(), nil, "")
	} else {
		token, _ := services.GetAuthToken(id, info.Username, isadmin)
		global.UnifiedReturn(c, global.Success, nil, nil, token)
	}
}
