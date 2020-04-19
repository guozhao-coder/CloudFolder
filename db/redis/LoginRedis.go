package redis

import (
	"gin-file/config"
	log "github.com/cihub/seelog"
	"time"
)

func SaveToken(idKey, token string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	conn := config.Redisdb
	s, err := conn.Set(idKey+"_fileSystemLogin", token, time.Hour*2).Result()
	log.Error("结果：", s, ",错误", err)
}

func CheckUserToken(uid, token string) (error, bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	conn := config.Redisdb
	result, err := conn.Get(uid + "_fileSystemLogin").Result()
	if err != nil {
		log.Error("该token已过期")
		return err, false
	}
	if result != token {
		log.Error("该token错误")
		return nil, false
	}
	return nil, true
}

func DeleteToken(id string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	conn := config.Redisdb
	i, err := conn.Del(id + "_fileSystemLogin").Result()
	log.Error("结果：", i, ",错误", err)
}
