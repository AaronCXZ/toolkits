package http

import (
	"context"
	"net/http"
	"time"

	"github.com/Muskchen/toolkits/logs"
	"github.com/gin-gonic/gin"
)

var srv = &http.Server{
	ReadTimeout:    time.Duration(10) * time.Second,
	WriteTimeout:   time.Duration(10) * time.Second,
	MaxHeaderBytes: 1 << 20,
}
var logger = logs.GetSLogger()

func Start(r *gin.Engine, addr string) {
	srv.Addr = addr
	srv.Handler = r
	go func() {
		logger.Infof("starting http server, listening on: %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("listening %s occur error: %s", srv.Addr, err)
		}
	}()
}

func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("cannot shutdown http server: %s", err)
	}
	select {
	case <-ctx.Done():
		logger.Info("shutdown http server timeout of 5 seconds.")
	default:
		logger.Info("http server stopped")
	}
}
