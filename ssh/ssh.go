package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHMessage struct {
	User, Password, Host, Cmd string
	Port                      int
	Client                    *ssh.Client
}

func NewSSH(user, password, host string, port int) *SSHMessage {
	return &SSHMessage{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (s *SSHMessage) Start() error {
	return s.start()
}

// 登录服务器
func (s *SSHMessage) start() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clinetConfig *ssh.ClientConfig
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(s.Password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clinetConfig = &ssh.ClientConfig{
		User:            s.User,
		Auth:            auth,
		HostKeyCallback: hostKeyCallbk,
	}

	addr = fmt.Sprintf("%s:%d", s.Host, s.Port)

	if s.Client, err = ssh.Dial("tcp", addr, clinetConfig); err != nil {
		return err
	}
	return nil

}

// 执行命令
func (s *SSHMessage) Run(cmd string) (out string, err error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var stdOut, stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	if err = session.Run(cmd); err != nil {
		return "", err
	}
	if stdErr.Len() != 0 {
		return "", errors.New(stdErr.String())
	}
	out = stdOut.String()
	return out, nil
}

// 退出登录
func (s *SSHMessage) Stop() error {
	return s.stop()
}

func (s *SSHMessage) stop() error {
	if err := s.Client.Close(); err != nil {
		return err
	}
	return nil
}
