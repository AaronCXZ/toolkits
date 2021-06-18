package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Muskchen/toolkits/godis/lib/sync/atomic"
	"github.com/Muskchen/toolkits/godis/lib/sync/wait"
	"github.com/Muskchen/toolkits/logs"
)

// EchoHandler 将收到的数据回传给客户端,用于测试
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

// MakeEchoHandler 创建EchoHandler
func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// EchoClient EchoHandler的客户端,用于测试
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close 关闭连接
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// Handle 回传数据到客户端
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// 关闭处理程序拒绝的连接
		_ = conn.Close()
	}

	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				logs.Logger().Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logs.Logger().Warn(err.Error())
			}
			return
		}
		client.Waiting.Add(1)
		_, _ = conn.Write(msg)
		client.Waiting.Done()
	}
}

// Close 关闭
func (h *EchoHandler) Close() error {
	logs.Logger().Info("handler shutting down....")
	h.closing.Set(true)
	h.activeConn.Range(func(key, val interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
