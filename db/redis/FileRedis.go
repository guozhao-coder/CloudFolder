package redis

import (
	"gin-file/config"
	log "github.com/cihub/seelog"
	"time"
)

func SaveFileToken(fId, fToken string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	conn := config.Redisdb
	s, err := conn.Set(fId+"_fileShareToken", fToken, time.Hour*1).Result()
	log.Error("结果：", s, ",错误：", err)
}

func CheckFileToken(fid, ftoken string) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	conn := config.Redisdb
	result, err := conn.Get(fid + "_fileShareToken").Result()
	if err != nil {
		log.Error("该token已过期")
		return false, err
	}
	if result != ftoken {
		log.Error("该token错误")
		return false, nil
	}
	return true, nil
}
