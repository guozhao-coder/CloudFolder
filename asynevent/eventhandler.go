package asynevent

import (
	"gin-file/base"
	log "github.com/cihub/seelog"
)

var EventMsgChan chan EventMsg

//向控制并发的管道放的数据
const chanMSG = 1

type EventMsg struct {
	//需要处理的消息码
	Code int
	//消息体
	Data interface{}
}

//初始化管道
func ChanInit() {
	msgChan := make(chan EventMsg, 100)
	EventMsgChan = msgChan
}

//往管道放消息
func PushMsg(e EventMsg) {
	EventMsgChan <- e
	return
}

func WaitEventMsg() {
	//注册一个控制并发数的管道
	chanCtrl := make(chan int, 10)
	for {
		select {
		case E := <-EventMsgChan:
			switch E.Code {
			case 1: //删除实际文件
				chanCtrl <- chanMSG
				go delFile(E.Data.(*base.FileStruct), chanCtrl)
			case 2: //将文件转移到阿里云oss
				chanCtrl <- chanMSG
				go saveFileToOSS(E.Data.(*base.FileStruct), chanCtrl)
			default:
				log.Error("请求有误")
			}
		default:

		}
	}
}
