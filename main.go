package main

import (
	"gin-file/asynevent"
	"gin-file/config"
	"gin-file/router"
	log "github.com/cihub/seelog"

	"net/http"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	//数据库服务
	err := config.DBInit()
	if err != nil {
		log.Error(err)
		return
	}
	//启动oss服务
	err = config.GetOSSBucket()
	if err != nil {
		log.Error(err)
		return
	}

	//开启异步处理协程
	asynevent.ChanInit()
	go asynevent.WaitEventMsg()

	//本地新建了一个分支为guozhao

	//设置路由
	r := router.Router()

	s := &http.Server{
		Addr:    ":5656",
		Handler: r,
	}
	//开启监听
	s.ListenAndServe()

}
