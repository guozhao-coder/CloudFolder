package router

import (
	"gin-file/base/code"
	"gin-file/db/redis"
	"gin-file/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//权限过滤器
func auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.Request.Header.Get("web-token")
		if tokenStr == "" {
			ctx.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN, "message": "没有token"})
			ctx.Abort()
			return
		}
		claims, err := utils.ParseToken(tokenStr, []byte("usersign"))
		if err != nil {
			ctx.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN, "message": "验证token出错"})
			ctx.Abort()
			return
		}
		//获取到uid
		uid := claims.(jwt.MapClaims)["userId"]
		//log.Info("用户id：", uid.(string))
		//把解析出来的id传给后文
		ctx.Set("userId", uid.(string))
		//ctx.Abort()停止中间件用abort
		//此时需要到redis验证
		err, bool := redis.CheckUserToken(uid.(string), tokenStr)
		if err != nil {
			ctx.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN_TIMEOUT, "message": "该token过期"})
			ctx.Abort()
			return
		}
		if !bool {
			ctx.JSON(200, gin.H{"code": code.LOGIN_USER_NO_PERMISSION, "message": "该token错误"})
			ctx.Abort()
			return
		}
	}
}
