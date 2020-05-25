package email

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cfg := SMTPMessage{
		Host:     "smtp.exmail.qq.com",
		User:     "yw@tuanche.com",
		Password: "Abcd1234%",
		Port:     465,
	}
	t.Run("Init", func(t *testing.T) {
		Init(cfg)
		assert.Equal(t, "yw@tuanche.com", smtp.User)
		assert.Equal(t, "Abcd1234%", smtp.Password)
		assert.Equal(t, "smtp.exmail.qq.com", smtp.Host)
		assert.Equal(t, 465, smtp.Port)
	})
	t.Run("New Message and send", func(t *testing.T) {
		msg, err := NewMessage("xizhong.chen@tuanche.com", "测试", "测试邮件")
		if err != nil {
			assert.Error(t, err)
		}
		ok := fmt.Sprintf("to:\t%s\tSub:\t%s\tBody:\t%s",
			"xizhong.chen@tuanche.com", "测试", "测试邮件")
		assert.Equal(t, ok, msg.String())

		if err := Send(msg); err != nil {
			assert.Error(t, err)
		}
	})
}
