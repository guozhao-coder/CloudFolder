package utils

import (
	"fmt"
	"regexp"
)

const (

	//邮箱匹配
	EMAIL_MATCH = "^([a-z0-9_.-]+)@([da-z.-]+).([a-z.]{2,6})$"
	//姓名匹配
	NAME_MATCH = "^[\u4E00-\u9FA5]+$"
)

//正则匹配
func Match(data string, sig string) bool {
	r, err := regexp.Compile(sig)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return r.MatchString(data)
}
