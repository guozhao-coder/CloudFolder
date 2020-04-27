package service

import (
	"gin-file/base"
	"gin-file/db/mongo"
	"gin-file/db/mysql"
	"gin-file/db/redis"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func AddFile(f *base.FileStruct) (*base.NormalResponse, error) {
	mysql.AddFile(f)
	return mongo.AddFile(f)
}

func GetFileHash(hash string) (string, bool, error) {
	//mysql.GetFileHash(hash)
	return mongo.GetFileHash(hash)
}

func GetFilesList(id string) (*base.FilesResponse, error) {
	//return mysql.GetFileList(id)
	return mongo.GetFileList(id)
}

func CheckFileMaster(uid, fid string) (*base.NormalResponse, error) {
	//return mysql.CheckFileMaster(uid, fid)
	return mongo.CheckFileMaster(uid, fid)
}

func GetFileInfoById(fid string) (*base.FileStruct, error) {
	//return mysql.GetFileInfoById(fid)
	return mongo.GetFileInfoById(fid)
}

func DeleteFile(fid string) (*base.NormalResponse, error) {
	//return mysql.DeleteFile(fid)
	return mongo.DeleteFile(fid)
}

func GetFileCountByHash(hash string) (int, error) {
	//return mysql.GetFileCountByHash(hash)
	return mongo.GetFileCountByHash(hash)
}

func GetFileByName(fname, uid string) (*base.FilesResponse, error) {
	//return mysql.GetFileByName(fname, uid)
	return mongo.GetFileByName(fname, uid)
}

func GetUserCapacity(uid string) (*base.NormalResponse, error) {
	//return mysql.GetUserCapacity(uid)
	return mongo.GetUserCapacity(uid)
}

func GetUserRemain(uid string) (float64, error) {
	//return mysql.GetUserRemain(uid)
	return mongo.GetUserRemain(uid)
}

func GenerateFileToken(uid, fid string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": uid,                                  //用户id
		"fileId": fid,                                  //文件id
		"exp":    time.Now().Add(time.Hour * 1).Unix(), //过期时间
	})
	s, _ := token.SignedString([]byte("filesign"))

	return s
}

func SaveFileToken(fid, ftoken string) {
	redis.SaveFileToken(fid, ftoken)
}

func CheckFileToken(fid, ftoken string) (bool, error) {
	return redis.CheckFileToken(fid, ftoken)
}
