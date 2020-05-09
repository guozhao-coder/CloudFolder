package utils

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/jordan-wright/email"
	"mime"
	"net/smtp"
)

const (
	hostPwd  = ""
	hostUser = "1733755545@qq.com"
	hostComm = "smtp.qq.com"
)

func SendMail(code int, receAddress, receName string) error {
	e := email.NewEmail()
	e.From = mime.QEncoding.Encode("UTF-8", "CloudFolder") + "<" + hostUser + ">"
	e.To = []string{receAddress}
	e.Subject = "CloudFolder找回密码操作"
	mailContext := fmt.Sprintf("<h1>%v，您好，您正在找回密码，验证码为：%d，请不要泄露！</h1>", receName, code)
	e.HTML = []byte(mailContext)

	auth := smtp.PlainAuth("", hostUser, hostPwd, hostComm)
	err := e.Send("smtp.qq.com:25", auth)
	if err != nil {
		seelog.Error("Send Mail to", receName, " Error:", err.Error())
		return err
	}
	seelog.Info("Send Mail to", receName, "Successfully")
	return nil
}
