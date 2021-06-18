package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Muskchen/toolkits/godis/interface/tcp"
	"github.com/Muskchen/toolkits/logs"
)

type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

// ListenAndServe 监听并提供服务,在closeChan接受到关闭信号后关闭服务
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan os.Signal) {
	// 正常关闭
	go func() {
		<-closeChan
		logs.Logger().Info("shutting down....")
		_ = listener.Close()
		_ = handler.Close()
	}()
	// 异常退出时释放资源
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		// 监听端口,阻塞直到新连接或者异常退出
		conn, err := listener.Accept()
		if err != nil {
			// 异常时break到waitDone.Wait,确保已建立的连接正常处理完成
			break
		}
		logs.Logger().Info("accept link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}

// ListenAndServeWithSignal 监听中断信号并通知关闭服务
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan os.Signal)
	signal.Notify(closeChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logs.Logger().Info(fmt.Sprintf("bind: %s, start listening....", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}
