package config

import (
	"database/sql"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/globalsign/mgo"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB
var Redisdb *redis.Client
var MgoClient *mgo.Session

var (
	//mysql连接参数
	driverName     = ""
	dataSourceName = ""
	//redis连接参数
	redisAddr = ""
	redisPwd  = ""
	//mongo连接的参数
	mongoAddr   = ""
	mongoUser   = ""
	mongoPwd    = ""
	mongoDBName = ""
)

//连接数据库
func DBInit() error {

	//连接redis
	if err := redisClient(); err != nil {
		return err
	}
	//连接mysql
	if err := mysqlClient(); err != nil {
		return err
	}
	//连接mongo
	if err := mongoClient(); err != nil {
		return err
	}

	return nil

}

func redisClient() error {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPwd,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Error("redis连接异常", err)
		return err
	}
	if pong == "PONG" {
		//_ = client.Options().Addr
		log.Info("redis连接成功")
		Redisdb = client
		return nil
	}

	return errors.New("redis conn  error......")

}

func mysqlClient() error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Error("err:", err.Error())
		return errors.New("db connect error...")
	}
	if err := db.Ping(); err != nil {
		log.Error("数据库连接异常", err.Error())
		return errors.New("db password error...")
	} else {
		log.Info("mysql连接成功")
		DB = db
		DB.SetMaxIdleConns(50)
		DB.SetMaxOpenConns(512)
		return nil
	}
	return errors.New("dbconn connect err....")
}

func mongoClient() error {
	var err error
	MgoClient, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:       []string{mongoAddr},
		Database:    mongoDBName,
		Username:    mongoUser,
		Password:    mongoPwd,
		MinPoolSize: 2048,
		PoolLimit:   2048,
		Timeout:     10 * time.Second,
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if err := MgoClient.Ping(); err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("mongo 连接成功")
	return nil
}
