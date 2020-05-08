package router

import (
	"gin-file/controller"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	pan := r.Group("/pan")
	{
		//登陆部分路由
		user := pan.Group("/user")
		{
			//用户注册
			user.POST("/register", controller.UserRegister)
			//用户登陆
			user.POST("/login", controller.CheckIdAndPwd)
			//用户登出
			user.GET("/logout", auth(), controller.Logout)
			//用户信息
			user.GET("/info", auth(), controller.GetUserName)
			//修改个人信息
			user.POST("/info/update", auth(), controller.UpdateUserInfo)
		}

		//文件操作部门路由
		file := pan.Group("/file")
		file.Use(auth())
		{
			//上传文件
			file.POST("/upload", controller.UpLoadFile)
			//查看文件列表
			file.GET("/getlist", controller.GetFilesList)
			//删除文件
			file.GET("/remove/:fid", controller.RemoveFile)
			//模糊查询文件
			file.GET("/getbyname/:filename", controller.GetFileByName)
			//查看容量信息
			file.GET("/capacity", controller.GetUserCapacity)
			//生成文件链接
			file.GET("/url/genarate/:fid", controller.GenerateUrl)
			//分享文件()
			//file.GET("/share/:fid", controller.ShareFile)
		}
	}

	//第一版下载文件
	//用户自己下载文件
	//r.GET("/pan/file/download/:fid/:token", controller.DownLoadFile)
	//他人下载文件
	//r.GET("/pan/file/downloadbyother/:ftoken", controller.DownLoadByOther)

	//第二版下载文件
	//服务器下载文件(不区分用户自己或者他人)
	r.GET("/pan/file/downloadbyurl/:ftoken", controller.DownLoadByUrl)

	r.StaticFile("/", "front/index.html")
	r.Static("/js/", "front/js/")
	r.Static("/css/", "front/css/")
	r.Static("/img/", "front/img/")
	r.Static("/fonts/", "front/fonts/")
	return r
}
