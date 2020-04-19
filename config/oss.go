package config

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/cihub/seelog"
)

var Bucket *oss.Bucket

const (
	OSSBucket          = ""
	OSSEndpoint        = ""
	OSSAccesskeyID     = ""
	OSSAccessKeySecret = ""
)

//获取oss的bucket实例
func GetOSSBucket() error {
	var err error
	client, err := oss.New(OSSEndpoint, OSSAccesskeyID, OSSAccessKeySecret)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	Bucket, err = client.Bucket(OSSBucket)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("oss实例创建成功")
	return nil
}

// 临时授权下载url
func DownloadURL(objName string) string {
	signedURL, err := Bucket.SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return signedURL
}

//判断oss是否存在该文件
func IfObjExist(s string) bool {
	b, err := Bucket.IsObjectExist(s)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return b
}
