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

func Start(r *gin.Engine, addr string) {
	srv.Addr = addr
	srv.Handler = r
	go func() {
		logs.SInfof("starting http server, listening on: %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logs.SErrorf("listening %s occur error: %s", srv.Addr, err)
		}
	}()
}

func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.SErrorf("cannot shutdown http server: %s", err)
	}
	select {
	case <-ctx.Done():
		logs.SInfo("shutdown http server timeout of 5 seconds.")
	default:
		logs.SInfo("http server stopped")
	}
}
