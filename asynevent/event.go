package asynevent

import (
	"bufio"
	"gin-file/base"
	"gin-file/config"
	"gin-file/service"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/cihub/seelog"
	"os"
)

//删除文件
func delFile(file *base.FileStruct) {
	i, err := service.GetFileCountByHash(file.FileHash)
	if err != nil {
		log.Error(err.Error())
		return
	}
	//先删除oss里的文件
	err = config.Bucket.DeleteObject(file.OssPath)
	if err != nil {
		log.Error("oss文件删除出错:", err.Error())
	}
	//删除文件
	if i == 0 {
		err = os.Remove(file.FilePath)
		if err != nil {
			log.Error("文件删除出错:", err.Error())
			return
		}
		return
	}
	return
}

//将文件同步到oss
func saveFileToOSS(file *base.FileStruct) {
	fd, err := os.Open(file.FilePath)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer fd.Close()
	//storageType := oss.ObjectStorageClass(oss.StorageStandard)

	contentType := oss.ContentType("application/x-msdownload")
	// 上传文件
	err = config.Bucket.PutObject(file.OssPath, bufio.NewReader(fd), contentType)
	if err != nil {
		log.Error("Error:", err)
		return
	}
	log.Info("同步文件", file.OssPath, "成功")
	return
}
