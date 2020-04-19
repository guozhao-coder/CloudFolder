package controller

import (
	"fmt"
	"gin-file/asynevent"
	"gin-file/base"
	"gin-file/base/code"
	"gin-file/config"
	"gin-file/service"
	"gin-file/utils"
	log "github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/url"
	"os"
	"strconv"
	"time"
)

//秒传
func fastUploadFile(c *gin.Context) bool {
	fhash := c.Request.Header.Get("filehash")
	//该文件不支持秒传
	dst, bool, err := service.GetFileHash(fhash)
	if err != nil || !bool || dst == "" {
		return false
	}
	//支持秒传
	//直接存入数据库
	uid := c.MustGet("userId").(string)
	log.Info(uid, "秒传文件")
	//绑定（前端需要把文件名放到请求头）
	fname := c.Request.Header.Get("filename")
	//转码
	finame, _ := url.QueryUnescape(fname)
	//获取文件大小
	fSize := utils.GetFileSize(dst)
	//判断文件大小
	sz, _ := strconv.ParseFloat(utils.GetFileSize(dst), 32)
	//查看用户剩余空间()确保用户不会超量上传
	remain, err := service.GetUserRemain(uid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.GET_ERR, "message": "查询失败"})
		return true
	}
	if remain < sz {
		c.JSON(200, gin.H{"code": code.SPACE_NOT_ENOUGH, "message": "空间不足"})
		return true
	}
	//上传成功之后，写入数据库
	//获取当前时间
	savetime := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04")

	f := &base.FileStruct{
		FileName: finame,
		FileSize: fSize,
		FilePath: dst,
		UserId:   uid,
		FileTime: savetime,
		FileHash: fhash,
	}
	//增加信息到数据库
	addFile(f, c, uid)
	return true
}

//上传
func UpLoadFile(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			c.JSON(200, gin.H{"code": code.FAILURE, "message": "panic"})
		}
	}()
	//尝试秒传
	fastBool := fastUploadFile(c)
	if fastBool {
		log.Info("支持秒传")
		return
	}
	fhash := c.Request.Header.Get("filehash")
	uid := c.MustGet("userId").(string)
	log.Info(uid, "上传文件")
	//绑定（前端需要把文件名放到请求头）
	fname := c.Request.Header.Get("filename")
	//转码
	finame, _ := url.QueryUnescape(fname)

	//接收文件
	formFile, fileHeader, err := c.Request.FormFile(fname)
	if err != nil {
		log.Error(err.Error())
		c.JSON(200, gin.H{"code": code.FILE_ACCEPT_ERROR, "message": "服务器接受文件失败"})
		return
	}
	defer formFile.Close()

	log.Debug(fileHeader.Filename)
	time1 := strconv.Itoa(int(time.Now().Unix()))
	// 上传文件至指定目录
	//服务器文件名组成 用户id+当前时间+文件名
	dst := fmt.Sprintf("filepath/" + uid + "_" + time1 + "_" + finame)

	err = c.SaveUploadedFile(fileHeader, dst)
	if err != nil {
		log.Error(err.Error())
		c.JSON(200, gin.H{"code": code.FILE_ACCEPT_ERROR, "message": "服务器接受文件失败"})
		return
	}
	//获取文件大小
	fSize := utils.GetFileSize(dst)
	//判断文件大小：过大自动删除(应该是读取文件时候判断)
	sz, _ := strconv.ParseFloat(utils.GetFileSize(dst), 32)
	if sz > 5120 {
		os.Remove(dst)
		log.Error("文件过大")
		c.JSON(200, gin.H{"code": code.FILE_TOO_BIG, "message": "文件过大"})
		return
	}
	//查看用户剩余空间()确保用户不会超量上传
	remain, err := service.GetUserRemain(uid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.GET_ERR, "message": "查询失败"})
		os.Remove(dst)
		return
	}
	if remain < sz {
		c.JSON(200, gin.H{"code": code.SPACE_NOT_ENOUGH, "message": "空间不足"})
		os.Remove(dst)
		return
	}
	//获取当前时间
	savetime := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04")

	f := &base.FileStruct{
		FileName: finame,
		FileSize: fSize,
		FilePath: dst,
		UserId:   uid,
		FileTime: savetime,
		FileHash: fhash,
	}
	//增加信息到数据库
	addFile(f, c, uid)
}

//增加文件到数据库
func addFile(f *base.FileStruct, c *gin.Context, uid string) {
	log.Info(f.UserId, "同步文件", f.FileName)
	//文件名为用户id，时间戳，文件名。
	unixTime := strconv.Itoa(int(time.Now().Unix()))
	dst := uid + unixTime + f.FileName
	f.OssPath = dst
	//将文件同步到阿里云oss
	asynevent.PushMsg(asynevent.EventMsg{Code: 2, Data: f})
	//调用service的方法
	response, err := service.AddFile(f)
	if err != nil {
		c.JSON(200, gin.H{"code": code.ADD_ERR, "message": "数据库插入失败"})
		//os.Remove(dst)
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": response.Message})
		//os.Remove(dst)
		return
	}

	c.JSON(200, gin.H{"code": code.SUCCESS, "message": "ok"})
	return
}

//下载文件
func DownLoadFile(c *gin.Context) {
	log.Info("下载文件")
	//首先会去数据库查询该文件是不是属于某用户
	fid := c.Param("fid")
	token := c.Param("token")
	claims, err := utils.ParseToken(token, []byte("usersign"))
	if err != nil {
		c.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN, "message": "验证token出错"})
		return
	}
	//获取到uid
	uid := claims.(jwt.MapClaims)["userId"].(string)

	response, err := service.CheckFileMaster(uid, fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "查询出错"})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.UNAUTHORIZED_OPERATION, "message": "无操作权限"})
		return
	}
	//根据fid查询出文件名以及路径
	fileStruct, err := service.GetFileInfoById(fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "查询出错"})
		return
	}
	//c.JSON(200,gin.H{"code":code.SUCCESS,"message":"downloading file"})
	//开始下载
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileStruct.FileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fileStruct.FilePath)

}

//分享文件
func ShareFile(c *gin.Context) {
	log.Info("分享文件")
	//首先会去数据库查询该文件是不是属于某用户
	fid := c.Param("fid")
	uid := c.MustGet("userId").(string)
	response, err := service.CheckFileMaster(uid, fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "查询出错"})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.UNAUTHORIZED_OPERATION, "message": "无操作权限"})
		return
	}
	//有操作权限，生成token
	fToken := service.GenerateFileToken(uid, fid)
	//把token存到redis
	service.SaveFileToken(fid, fToken)
	//返回给前端一个下载文件的token
	c.JSON(200, gin.H{"code": code.SUCCESS, "message": "操作成功", "token": fToken})
}

//生成文件链接
func GenerateUrl(c *gin.Context) {
	//首先会去数据库查询该文件是不是属于某用户
	fid := c.Param("fid")
	uid := c.MustGet("userId").(string)
	log.Info(uid, "生成文件链接,文件id：", fid)
	response, err := service.CheckFileMaster(uid, fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "查询出错"})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.UNAUTHORIZED_OPERATION, "message": "无操作权限"})
		return
	}
	//有操作权限，查看文件路径
	fileStruct, err := service.GetFileInfoById(fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "查询出错"})
		return
	}
	//判断oss中是否存在该文件，存在就返回http的链接，否则生成服务器自己的链接
	bool := config.IfObjExist(fileStruct.OssPath)
	if bool {
		url := config.DownloadURL(fileStruct.OssPath)
		c.JSON(200, gin.H{"code": code.SUCCESS, "message": "操作成功", "token": (url)})
		return
	}
	//生成服务器自己的链接
	fileToken := service.GenerateFileToken(uid, fid)
	fileUrl := "http://120.26.78.161:5656/pan/file/downloadbyurl/"
	c.JSON(200, gin.H{"code": code.SUCCESS, "message": "操作成功", "token": (fileUrl + fileToken)})
	return
}

//他人下载
func DownLoadByOther(c *gin.Context) {
	log.Info("下载文件")
	ftoken := c.Param("ftoken")
	//解析出fid
	claims, err := utils.ParseToken(ftoken, []byte("filesign"))
	if err != nil {
		c.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN, "message": "验证token出错"})
		return
	}
	fid := claims.(jwt.MapClaims)["fileId"].(string)
	//查看redis
	b, err := service.CheckFileToken(fid, ftoken)
	if err != nil {
		c.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN_TIMEOUT, "message": "验证token出错"})
		return
	}
	if !b {
		c.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN, "message": "验证token出错"})
		return
	}
	//根据fid查询出文件名以及路径
	fileStruct, err := service.GetFileInfoById(fid)
	if err != nil {
		//c.JSON(200,gin.H{"code":code.FAILURE,"message":"查询出错"})
		return
	}
	//c.JSON(200,gin.H{"code":code.SUCCESS,"message":"downloading file"})
	//开始下载
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileStruct.FileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fileStruct.FilePath)
}

//查看文件
func GetFilesList(c *gin.Context) {
	uid := c.MustGet("userId").(string)
	log.Info(uid, "查看文件列表")
	response, err := service.GetFilesList(uid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "获取文件列表失败"})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "获取文件列表失败"})
		return
	}
	c.JSON(200, gin.H{"code": code.SUCCESS, "result": response.Files})
	return
}

//下载文件
func DownLoadByUrl(c *gin.Context) {
	ftoken := c.Param("ftoken")
	//解析出fid
	claims, err := utils.ParseToken(ftoken, []byte("filesign"))
	if err != nil {
		c.JSON(200, gin.H{"code": code.LOGIN_NO_TOKEN, "message": "验证token出错"})
		return
	}
	fid := claims.(jwt.MapClaims)["fileId"].(string)
	log.Info("下载文件", fid)
	//根据fid查询出文件名以及路径
	fileStruct, err := service.GetFileInfoById(fid)
	if err != nil {
		//c.JSON(200,gin.H{"code":code.FAILURE,"message":"查询出错"})
		return
	}
	//c.JSON(200,gin.H{"code":code.SUCCESS,"message":"downloading file"})
	//开始下载
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileStruct.FileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fileStruct.FilePath)
}

//模糊查询文件
func GetFileByName(c *gin.Context) {
	uid := c.MustGet("userId").(string)
	log.Info(uid, "模糊查询文件")
	fname := c.Param("filename")
	response, err := service.GetFileByName(fname, uid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "获取文件列表失败"})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.RESPONSE_NIL, "message": "结果为空"})
		return
	}
	c.JSON(200, gin.H{"code": code.SUCCESS, "result": response.Files})
	return

}

//获取个人的容量信息
func GetUserCapacity(c *gin.Context) {
	uid := c.MustGet("userId").(string)
	log.Info(uid, "查看容量信息")
	response, err := service.GetUserCapacity(uid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "获取文件列表失败"})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "获取文件列表失败"})
		return
	}
	c.JSON(200, gin.H{"code": code.SUCCESS, "message": "查询成功", "result": response.Data})
}

//删除文件
func RemoveFile(c *gin.Context) {
	//首先会去数据库查询该文件是不是属于某用户
	//只是删除数据库的记录
	fid := c.Param("fid")
	uid := c.MustGet("userId").(string)
	log.Info(uid, "删除文件", fid)
	response, err := service.CheckFileMaster(uid, fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "查询出错"})
		return
	}
	if response.Code != 200 {
		c.JSON(200, gin.H{"code": code.UNAUTHORIZED_OPERATION, "message": "无操作权限"})
		return
	}

	//根据fid查询出文件名以及路径
	fileStruct, err := service.GetFileInfoById(fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.FAILURE, "message": "查询出错"})
		return
	}
	//删除数据库的记录
	res, err := service.DeleteFile(fid)
	if err != nil {
		c.JSON(200, gin.H{"code": code.DEL_ERR, "message": "删除出错"})
		return
	}
	if res.Code != 200 {
		c.JSON(200, gin.H{"code": code.DEL_ERR, "message": "删除出错"})
		return
	}
	//把删除文件消息推送到管道
	asynevent.PushMsg(asynevent.EventMsg{Code: 1, Data: fileStruct})
	c.JSON(200, gin.H{"code": code.SUCCESS, "message": "success!"})
	return

}
