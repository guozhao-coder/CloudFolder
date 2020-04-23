package main

import (
	"testing"
)

// BenchmarkSprintf 对 fmt.Sprintf 函数进行基准测试
//func BenchmarkSprintf(b *testing.B)  {
//	//msgChan := make(chan asynevent.EventMsg, 100)
//	ch := make(chan asynevent.EventMsg,10)
//	msg := &asynevent.EventMsg{Code:1,Data:base.FileStruct{
//		FileId:   "10548e8c-76ad-4b1e-984a-b308a424a157",
//		FileName: "fiesql.sql",
//		FileSize: "2.8",
//		FilePath: "filepath/guozhao_1587603696_fiesql.sql",
//		UserId:   "guozhao",
//		FileTime: "2020-04-23 09:01",
//		FileHash: "01a6e27aa44c3ddd1a07257877f45870",
//		OssPath:  "guozhao1587603696fiesql.sql",
//	}}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		ch <- asynevent.EventMsg{Code:1,Data:msg}
//		asynevent.WaitEventMsg(ch)
//	}
//
//}

func BenchmarkSprintf(b *testing.B) {
	//msgChan := make(chan asynevent.EventMsg, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

	}

}
