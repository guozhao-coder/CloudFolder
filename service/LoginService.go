package service

import (
	"gin-file/base"
	"gin-file/db/mongo"
	"gin-file/db/redis"
	"github.com/dgrijalva/jwt-go"
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
	tk := genarateToken(u.UserId)
	//存入redis
	redis.SaveToken(u.UserId, tk)
	return &base.NormalResponse{Code: 200, Message: "成功", Data: tk}, nil

}

func UserLogout(id string) {
	redis.DeleteToken(id)
}
func GetUserName(uid string) (string, error) {
	return mongo.GetUserName(uid)
}

func genarateToken(id string) string {
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
