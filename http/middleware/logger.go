package middleware

import (
	"time"

	"github.com/Muskchen/toolkits/logger"
	"github.com/gin-gonic/gin"
)

func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logger.Errorf("%s\n", c.Errors.String())
		errors := c.Errors.ByType(typ)
		if len(errors) > 0 {
			c.JSON(-1, errors)
		}
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latencyTime := time.Now().Sub(start)
		method := c.Request.Method
		path := c.Request.RequestURI
		code := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.Infof("code: %3d latencyTime: %d clientIP: %s method: %s requestUrl: %s\n",
			code, latencyTime, clientIP, method, path,
		)
	}
}
