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
	log.Info("结果：", s, ",错误", err)
}

func SaveMailCode(uid string, code int) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	conn := config.Redisdb
	s, err := conn.Set(uid+"_fileSystemCode", code, time.Minute*5).Result()
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	log.Info("结果：", s)
	return true, nil
}

func CheckMailCode(uid, code string) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	conn := config.Redisdb
	result, err := conn.Get(uid + "_fileSystemCode").Result()
	if err != nil {
		log.Error("未查询到该键值对")
		return false, err
	}
	if result != code {
		log.Error("该code错误")
		return false, nil
	}
	return true, err
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
	log.Info("结果：", i, ",错误", err)
}
