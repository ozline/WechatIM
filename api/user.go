package api

import(
	"wechat/structs"
	"wechat/global"
	"wechat/services"

	"github.com/gin-gonic/gin"
)

//注册
func UserRegister(c *gin.Context) {
	var info structs.User
	info.Username = c.PostForm("username")
	info.Password = c.PostForm("password")
	info.Phone = c.PostForm("phone")

	result, err := services.AddUser(info)
	if err != nil{
		global.UnifiedReturn(c, global.ErrorUsers,err.Error(),nil,"")
	}else{
		global.UnifiedReturn(c, global.Success,result,nil,"")
	}
}

//登录
func UserLogin(c *gin.Context) {
	var info structs.User
	// c.ShouldBind(&info)
	info.Username = c.PostForm("username")
	info.Password = c.PostForm("password")

	if len(info.Username) == 0 || len(info.Password) == 0 {
		global.UnifiedReturn(c, global.ErrorUsers,"用户名或密码错误",nil,"")
		return
	}
	id, isadmin, err := services.CheckUser(info)
	if err != nil {
		global.UnifiedReturn(c, global.ErrorUsers,err.Error(),nil,"")
	} else {
		token, _ := GetAuthToken(id, info.Username, isadmin)
		global.UnifiedReturn(c, global.Success,nil,nil,token)
	}
}

// //gin绑定结构体
// //绑定了：id/password/nickname/phone/email/avatar/gender/profile
// //其中id只作为识别，不会进行修改
// func UpdateUser(c *gin.Context) {
// 	var tmp structs.User
// 	c.Bind(&tmp)
// 	success, err := services.UpdateUser(tmp)
// 	normalReturn(c, nil, err, success, updateAuthToken(c))
// }

// func UpdateUserAvatar(c *gin.Context) {
// 	dirAvatar := "avatar"
// 	userid := c.PostForm("id")
// 	result, fileName := model.OSSUPloadViaFormData(c, "file", "avatar")
// 	if !result {
// 		errReturn(c, "头像上传失败")
// 		return
// 	}
// 	filePath := "https://" + conf.Config.OSS.Endpoint + "/" + conf.Config.OSS.MainDir + "/" + dirAvatar + "/" + fileName
// 	success, err := services.UpdateUserAvatar(userid, filePath)
// 	normalReturn(c, filePath, err, success, updateAuthToken(c))
// }

// //TODO:可能存在漏洞
// //上传文件
// func UserUpload(c *gin.Context) {
// 	var upload structs.Video
// 	dirFile := "video"
// 	dirCover := "cover"

// 	//上传文件
// 	result, fileName := model.OSSUPloadViaFormData(c, "file", dirFile)
// 	if !result {
// 		errReturn(c, "视频文件传输失败")
// 		return
// 	}
// 	upload.FilePath = "https://" + conf.Config.OSS.Endpoint + "/" + conf.Config.OSS.MainDir + "/" + dirFile + "/" + fileName

// 	//上传封面
// 	result, fileName = model.OSSUPloadViaFormData(c, "cover", dirCover)
// 	if !result {
// 		errReturn(c, "封面文件传输失败")
// 		return
// 	}
// 	upload.CoverPath = "https://" + conf.Config.OSS.Endpoint + "/" + conf.Config.OSS.MainDir + "/" + dirCover + "/" + fileName

// 	//其他处理
// 	upload.Title = c.PostForm("title")
// 	upload.Category = c.PostForm("category")
// 	upload.Profile = c.PostForm("profile")
// 	upload.Status = 0
// 	upload.Author = c.PostForm("author")
// 	upload.User = c.PostForm("user")
// 	upload.Update = c.PostForm("time")

// 	res, err := services.AddVideo(upload)

// 	path := map[string]string{
// 		"filePath":  upload.FilePath,
// 		"coverPath": upload.CoverPath,
// 	}
// 	normalReturn(c, path, err, res, updateAuthToken(c))
// }

// //获取用户ID
// func GetUserID(c *gin.Context) string {
// 	token := c.Request.Header.Get("AuthToken")
// 	claims, _ := middleware.JWTParse(token)
// 	return claims.Id
// }