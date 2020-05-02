package email

import (
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

type SMTPMessage struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

type Message struct {
	Tos     []string `yaml:"tos"`
	Subject string   `yaml:"subject"`
	Content string   `yaml:"content"`
}

var (
	Smtp SMTPMessage
	conn *gomail.Dialer
)

func Init(cfg SMTPMessage) {
	Smtp = cfg
	conn = gomail.NewDialer(Smtp.Host, Smtp.Port, Smtp.User, Smtp.Password)
}

// 邮件内容校验
func (m *Message) validate() error {
	if m.Tos == nil || len(m.Tos) == 0 {
		return fmt.Errorf("%s：收件人为空。", m)
	}

	if len(m.Subject) == 0 {
		return fmt.Errorf("%s：邮件标题为空。", m)
	}

	if len(m.Content) == 0 {
		return fmt.Errorf("%s：邮件正文为空。", m)
	}

	return nil
}

// 打印邮件内容信息
func (m *Message) String() string {
	to := strings.Join(m.Tos, ",")
	return fmt.Sprintf("to:\t%s\tSub:\t%s\tBody:\t%s",
		to, m.Subject, m.Content)
}

// 初始化邮件内容
func NewMessage(to, subject, body string) (mail *Message, err error) {
	tos := strings.Split(to, ",")
	m := &Message{
		Tos:     tos,
		Subject: subject,
		Content: body,
	}
	if err := m.validate(); err == nil {
		return m, nil
	} else {
		return nil, err
	}
}

// 发送邮件
func Send(msg *Message) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", Smtp.User)      // 发件人
	m.SetHeader("To", msg.Tos...)       // 收件人
	m.SetHeader("Subject", msg.Subject) // 主题
	m.SetBody("text/html", msg.Content) // 正文
	if err := conn.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
