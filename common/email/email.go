package email

import (
	"GopherAI/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	CodeMsg     = "GopherAI验证码如下，请保留好，后续可以用验证码登录 (验证码仅限于2分钟有效): "
	UserNameMsg = "GopherAI的账号如下，请保留好，后续可以用账号/邮箱登录 "
)

func SendCaptcha(email, code, msg string) error {
	m := gomail.NewMessage()

	// 发件人
	m.SetHeader("From", config.GetConfig().EmailConfig.Email)
	// 收件人
	m.SetHeader("To", email)
	// 主题
	m.SetHeader("Subject", "来自GopherAI的验证码")
	// 正文内容（纯文本形式，也可以用 HTML 格式）
	m.SetBody("text/plain", msg+" "+code)

	// 配置 SMTP 服务器和授权信息
	// 这里使用 QQ 邮箱的 SMTP 服务器，端口为 587（STARTTLS 端口）
	d := gomail.NewDialer("smtp.qq.com", 587, config.GetConfig().EmailConfig.Email, config.GetConfig().EmailConfig.Authcode)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("DialAndSend err %v:\n", err)
		return err
	}
	fmt.Printf("send mail success\n")
	return nil
}
