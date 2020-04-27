package mongo

import (
	"gin-file/base"
	"gin-file/base/code"
	"gin-file/config"
	log "github.com/cihub/seelog"
	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

//通过文件hash获取文件路径
func GetFileHash(hash string) (string, bool, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var file *base.FileStruct
	err := c.Find(bson.M{"file_hash": hash}).One(&file)
	if err != nil {
		log.Error(err)
		//如果错误类型为空，说明不存在该hash
		if err.Error() == "not found" {
			return "", false, nil
		}
		return "", false, err
	}
	return file.FilePath, true, nil
}

//增加文件
func AddFile(f *base.FileStruct) (*base.NormalResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	//生成uuid代表文件id
	fid := uuid.NewV4().String()
	f.FileId = fid
	err := c.Insert(f)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &base.NormalResponse{Code: 200, Message: "插入成功"}, nil
}

func GetFileList(id string) (*base.FilesResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var files []*base.FileStruct
	err := c.Find(bson.M{"user_id": id}).All(&files)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//加上单位
	for i := 0; i < len(files); i++ {
		files[i].FileSize += "kb"
	}
	fileRes := new(base.FilesResponse)
	fileRes.Code = 200
	fileRes.Files = files
	return fileRes, nil
}

//查看文件权限
func CheckFileMaster(uid, fid string) (*base.NormalResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var file []*base.FileStruct
	err := c.Find(bson.M{"file_id": fid, "user_id": uid}).All(&file)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if len(file) == 0 {
		return &base.NormalResponse{Code: code.UNAUTHORIZED_OPERATION, Message: "无权限"}, nil
	}
	return &base.NormalResponse{Code: 200, Message: "success"}, nil
}

//查看文件的路径
func GetFileInfoById(fid string) (*base.FileStruct, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var file *base.FileStruct
	err := c.Find(bson.M{"file_id": fid}).One(&file)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return file, nil
}

func DeleteFile(fid string) (*base.NormalResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	err := c.Remove(bson.M{"file_id": fid})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &base.NormalResponse{Code: 200}, nil
}

func GetFileCountByHash(hash string) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var files []*base.FileStruct
	err := c.Find(bson.M{"file_hash": hash}).All(&files)
	if err != nil {
		log.Error(err.Error())
		if err.Error() == "not found" {
			return 0, nil
		}
		return 0, err
	}
	return len(files), nil
}

//查询文件
func GetFileByName(fname, uid string) (*base.FilesResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var files []*base.FileStruct
	err := c.Find(bson.M{"file_name": fname, "user_id": uid}).All(&files)
	if err != nil {
		log.Error(err.Error())
		if err.Error() == "not found" {
			return &base.FilesResponse{Code: 200, Files: nil}, nil
		}
		return nil, err
	}
	fileRes := new(base.FilesResponse)
	fileRes.Code = 200
	fileRes.Files = files
	return fileRes, nil
}

//查询用户空间
func GetUserCapacity(uid string) (*base.NormalResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var files []*base.FileStruct
	err := c.Find(bson.M{"user_id": uid}).All(&files)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var usedCap float64
	for i := 0; i < len(files); i++ {
		fsize, _ := strconv.ParseFloat(files[i].FileSize, 64)
		usedCap += fsize
	}
	capacity := &base.UserCapacity{
		TotalCapacity:  50 * 1024,
		UsedCapacity:   int(usedCap),
		RemainCapacity: 50*1024 - int(usedCap),
	}
	capacityResp := new(base.NormalResponse)
	capacityResp.Code = 200
	capacityResp.Data = capacity
	return capacityResp, nil
}

//用户剩余控件
func GetUserRemain(uid string) (float64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	c := config.MgoClient.Copy().DB("").C("file")
	defer c.Database.Session.Close()
	var files []*base.FileStruct
	err := c.Find(bson.M{"user_id": uid}).All(&files)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	//已用空间
	var usedCap float64
	for i := 0; i < len(files); i++ {
		fsize, _ := strconv.ParseFloat(files[i].FileSize, 64)
		usedCap += fsize
	}
	return 50*1024 - usedCap, nil
}
