package controller

import (
	"gin-file/base"
	"gin-file/base/code"
	"gin-file/service"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func CheckIdAndPwd(c *gin.Context) {
	var user base.UserStruct
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(200, gin.H{"code": code.JSON_UNMARSHAL_ERROR, "message": "出错"})
		return
	}
	log.Info(user.UserId, "登陆")
	//调用service
	response, err := service.UserLogin(&user)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "服务繁忙", "errlog": err.Error()})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.LOGIN_PASSWORD_ACCOUNT_ERROR, "message": "账号密码错误"})
		return
	}
	c.JSON(200, gin.H{"code": code.SUCCESS, "message": "验证成功", "token": response.Data})
	return
}

func Logout(c *gin.Context) {
	uid := c.MustGet("userId").(string)
	log.Info(uid, "退出登陆")
	service.UserLogout(uid)
}

func GetUserName(c *gin.Context) {
	uid := c.MustGet("userId").(string)
	log.Info(uid, "查看用户")
	name, err := service.GetUserName(uid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "服务繁忙"})
		return
	}
	c.JSON(200, gin.H{"code": code.SUCCESS, "message": "成功", "data": name})
	return
}

func UserRegister(c *gin.Context) {
	log.Info("用户注册")
	var user base.UserStruct
	err := c.ShouldBindJSON(&user)
	if user.UserId == "" || user.Password == "" || user.Username == "" {
		c.JSON(200, gin.H{"code": code.PAR_PARAMETER_IS_NULL, "message": "参数不能为空"})
		return
	}
	if err != nil {
		c.JSON(200, gin.H{"code": code.JSON_MARSHAL_ERROR, "message": "json解析错误"})
		return
	}
	response, err := service.UserRegister(&user)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "服务繁忙"})
		return
	}
	if response.Code == 200 {
		c.JSON(200, gin.H{"code": code.SUCCESS, "message": "注册成功"})
		return
	}
	c.JSON(200, gin.H{"code": code.DATA_EXIST, "message": "用户已存在"})
	return
}
