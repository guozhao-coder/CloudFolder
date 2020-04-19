package mysql

import (
	"errors"
	"gin-file/base"
	"gin-file/base/code"
	"gin-file/config"
	log "github.com/cihub/seelog"
)

func UserLogin(u *base.UserStruct) (*base.NormalResponse, error) {
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
	str := "SELECT username FROM user WHERE userId = ? AND password = ? "
	rows, err := begin.Query(str, u.UserId, u.Password)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	if !rows.Next() {
		log.Error("用户名密码错误")
		begin.Rollback()
		return &base.NormalResponse{Code: code.LOGIN_PASSWORD_ACCOUNT_ERROR, Message: "用户名密码错误"}, nil
	}
	rows.Close()
	begin.Commit()
	return &base.NormalResponse{Code: code.SUCCESS, Message: "验证成功"}, nil

}

func UserRegister(u *base.UserStruct) (*base.NormalResponse, error) {
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
	str := "select * from user where userId = ?"
	rows, err := db.Query(str, u.UserId)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	if rows.Next() {
		log.Error("用户名已存在")
		begin.Rollback()
		return &base.NormalResponse{Code: code.DATA_EXIST, Message: "用户已存在"}, nil
	}
	rows.Close()
	str2 := "insert into user values(?,?,?)"
	result, err := begin.Exec(str2, u.UserId, u.Password, u.Username)
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
		log.Error("影响的行数不为1")
		begin.Rollback()
		return nil, errors.New("影响行数不为1")
	}
	begin.Commit()
	return &base.NormalResponse{Code: 200}, nil

}

func GetUserName(uid string) (string, error) {
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
		return "", err
	}
	str := "select username from user where userId = ?"
	row := begin.QueryRow(str, uid)
	var uname string
	err = row.Scan(&uname)
	if err != nil {
		log.Error(err.Error())
		begin.Rollback()
		return "", err
	}
	begin.Commit()
	return uname, nil
}
