package utils

import (
	"fmt"

	logs "github.com/cihub/seelog"
)

func init() {
	log, err := logs.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logs.ReplaceLogger(log)
}
