package service

import (
	"gin-file/asynevent"
	"gin-file/base"
	"gin-file/db/mongo"
	"gin-file/db/redis"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

func UserLogin(u *base.UserStruct) (*base.NormalResponse, error) {
	response, err := mongo.UserLogin(u)
	if err != nil {
		return nil, err
	}
	if response.Code != 200 {
		return response, nil
	}
	//验证成功，生成token
	tk := GenarateToken(u.UserId)
	//存入redis
	SaveToken(u.UserId, tk)
	return &base.NormalResponse{Code: 200, Message: "成功", Data: tk}, nil

}

func SaveToken(uid, utoken string) {
	redis.SaveToken(uid, utoken)
}

func UserLogout(id string) {
	redis.DeleteToken(id)
}
func GetUserName(uid string) (string, error) {
	userStruct, err := mongo.GetUserInfo(uid)
	return userStruct.Username, err
}

func GenarateToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,                                   //用户id
		"exp":    time.Now().Add(time.Hour * 2).Unix(), //过期时间
	})
	s, _ := token.SignedString([]byte("usersign"))

	return s
}

func UserRegister(u *base.UserStruct) (*base.NormalResponse, error) {
	return mongo.UserRegister(u)
}

func CheckMailCode(uid, ucode string) (bool, error) {
	return redis.CheckMailCode(uid, ucode)
}

func UpdateUserInfo(u *base.UserStruct) (bool, bool, error) {
	//首先比对用户名密码是否匹配
	response, err := mongo.UserLogin(u)
	if err != nil {
		return false, false, err
	}
	//用户名密码出错
	if response.Code != 200 {
		return false, false, nil
	}
	//执行修改
	b, err := mongo.UpdateUserInfo(u)
	if err != nil {
		return true, false, err
	}
	if b {
		return true, true, nil
	}
	return true, false, nil
}

//发送邮件
func SendMail(uid string) (*base.NormalResponse, error) {
	userStruct, err := mongo.GetUserInfo(uid)
	if err != nil {
		return &base.NormalResponse{Code: 1, Message: "查询此人出错"}, nil
	}
	//生成随机数，存到redis
	emailCode := (int(int32(1000) + rand.New(rand.NewSource(time.Now().Unix())).Int31n(10000-1000)))
	_, err = redis.SaveMailCode(uid, emailCode)
	if err != nil {
		return &base.NormalResponse{Code: 2, Message: "保存验证码失败"}, nil
	}
	//发送邮件（异步发送）
	//需要信息为：用户名，验证码，用户邮箱
	uMailInfo := &base.EmailStruct{
		Username: userStruct.Username,
		UserMail: userStruct.UserMail,
		MailCode: emailCode,
	}
	//将该信息送到管道(异步，响应时间由1.5319309s变为207.4343ms)
	asynevent.PushMsg(asynevent.EventMsg{Code: 3, Data: uMailInfo})
	//返回用户信息
	return &base.NormalResponse{Code: 200, Message: "已推送到管道"}, nil
}
