package mysql

import (
	"database/sql"
	"errors"
	"gin-file/base"
	"gin-file/base/code"
	"gin-file/config"
	log "github.com/cihub/seelog"
	uuid "github.com/satori/go.uuid"
)

//查看文件
func GetFileHash(hash string) (string, bool, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return "", false, err
	}
	str := "select filepath from file where filehash = ? limit 1"
	row := begin.QueryRow(str, hash)
	var filepa sql.NullString
	row.Scan(&filepa)
	if filepa.String == "" {
		begin.Rollback()
		return "", false, nil
	}
	begin.Commit()
	return filepa.String, true, nil
}

//上传文件
func AddFile(f *base.FileStruct) (*base.NormalResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	//生成uuid代表文件id
	fid := uuid.NewV4()

	str := "INSERT INTO file VALUES(?,?,?,?,?,?,?,?)"
	log.Debug(str)
	result, err := begin.Exec(str, fid.String(), f.FileName, f.FileSize, f.FilePath, f.UserId, f.FileTime, f.FileHash, f.OssPath)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}

	i, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	if i != 1 {
		log.Error("影响行数不为1")
		begin.Rollback()
		return nil, errors.New("影响行数不为1")
	}
	begin.Commit()
	return &base.NormalResponse{Code: 200, Message: "插入成功"}, nil
}

//查看文件
func GetFileList(id string) (*base.FilesResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	str := "select fileId,fileName,fileSize,userId,fileTime from file where userId = ?"
	log.Debug(str)
	rows, err := begin.Query(str, id)

	if err != nil {
		begin.Rollback()
		log.Error(err.Error())
		return nil, err
	}
	var files []*base.FileStruct
	for rows.Next() {
		var file base.FileStruct
		err := rows.Scan(&file.FileId, &file.FileName, &file.FileSize, &file.UserId, &file.FileTime)
		if err != nil {
			begin.Rollback()
			log.Error(err.Error())
			return nil, err
		}
		file.FileSize = file.FileSize + "kb"
		files = append(files, &file)
	}
	rows.Close()
	begin.Commit()
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

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	str := "select * from file where fileId = ? and userId = ?"
	log.Debug(str)
	rows, err := begin.Query(str, fid, uid)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	if !rows.Next() {
		begin.Rollback()
		return &base.NormalResponse{Code: code.UNAUTHORIZED_OPERATION, Message: "无权限"}, nil
	}

	rows.Close()
	begin.Commit()
	return &base.NormalResponse{Code: 200, Message: "success"}, nil
}

//查看文件的路径
func GetFileInfoById(fid string) (*base.FileStruct, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}

	str := "select fileName,filePath,fileHash,ossPath from file where fileId = ?"
	log.Debug(str)
	row := begin.QueryRow(str, fid)
	var file base.FileStruct
	err = row.Scan(&file.FileName, &file.FilePath, &file.FileHash, &file.OssPath)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	begin.Commit()
	return &file, nil

}

//删除文件
func DeleteFile(fid string) (*base.NormalResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	str := "delete from file where fileId = ?"
	log.Debug(str)
	result, err := begin.Exec(str, fid)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	i, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	if i != 1 {
		log.Error("affect row err")
		begin.Rollback()
		return nil, errors.New("affect row err")
	}
	begin.Commit()
	return &base.NormalResponse{Code: 200}, nil
}

//根据hash获取文件数量
func GetFileCountByHash(hash string) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return 0, err
	}
	str := "SELECT COUNT(fileHash) FROM file WHERE fileHash = ?"
	var cou int
	err = begin.QueryRow(str, hash).Scan(&cou)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return 0, err
	}
	begin.Commit()
	return cou, nil

}

//模糊查询文件
func GetFileByName(fname, uid string) (*base.FilesResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	str := `SELECT fileId,fileName,fileSize,userId,fileTime FROM file WHERE userId = ? AND fileName LIKE '%` + fname + `%' `
	log.Debug(str)

	rows, err := begin.Query(str, uid)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	if rows == nil {
		begin.Rollback()
		log.Error("结果为空")
		return &base.FilesResponse{Code: 5000}, nil
	}
	var files []*base.FileStruct
	for rows.Next() {
		var file base.FileStruct
		err := rows.Scan(&file.FileId, &file.FileName, &file.FileSize, &file.UserId, &file.FileTime)
		if err != nil {
			begin.Rollback()
			log.Error(err.Error())
			return nil, err
		}
		file.FileSize = file.FileSize + "kb"
		files = append(files, &file)
	}
	rows.Close()
	begin.Commit()
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

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	str := "SELECT SUM(fileSize) FROM file WHERE userId = ?"
	row := begin.QueryRow(str, uid)
	var used sql.NullFloat64
	err = row.Scan(&used)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	begin.Commit()

	//remian := (50*1024 - used)/1024

	userC := int(used.Float64)
	remainC := 50*1024 - userC
	capacity := &base.UserCapacity{
		TotalCapacity:  50 * 1024,
		UsedCapacity:   userC,
		RemainCapacity: remainC,
	}
	capacityResp := new(base.NormalResponse)
	capacityResp.Code = 200
	capacityResp.Data = capacity
	return capacityResp, err
}

//查询剩余空间
func GetUserRemain(uid string) (float64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	db := config.DB
	begin, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return 0, err
	}
	str := "SELECT SUM(fileSize) FROM file WHERE userId = ?"
	row := begin.QueryRow(str, uid)
	var used sql.NullFloat64
	err = row.Scan(&used)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return 0, err
	}
	begin.Commit()
	//计算剩余空间
	remain := 50*1024 - used.Float64
	return remain, nil
}
