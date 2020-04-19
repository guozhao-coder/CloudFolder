package asynevent

import (
	"gin-file/base"
	log "github.com/cihub/seelog"
)

var EventMsgChan chan EventMsg

type EventMsg struct {
	//需要处理的消息码
	Code int
	//消息体
	Data interface{}
}

//初始化管道
func ChanInit() {
	msgChan := make(chan EventMsg, 10)
	EventMsgChan = msgChan
}

//往管道放消息
func PushMsg(e EventMsg) {
	EventMsgChan <- e
	return
}

func WaitEventMsg() {
	for {
		select {
		case E := <-EventMsgChan:
			switch E.Code {
			case 1: //删除实际文件
				go delFile(E.Data.(*base.FileStruct))
			case 2: //将文件转移到阿里云oss
				go saveFileToOSS(E.Data.(*base.FileStruct))
			default:
				log.Error("请求有误")
			}
		default:

		}
	}
}
