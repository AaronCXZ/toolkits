package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Message(c *gin.Context, v interface{}) {
	if v == nil {
		c.JSON(http.StatusOK, gin.H{"err": ""})
	}

	switch t := v.(type) {
	case string:
		c.JSON(http.StatusOK, gin.H{"err": t})
	case error:
		c.JSON(http.StatusOK, gin.H{"err": t.Error()})
	}
}

func Data(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"dat": data, "err": ""})
		return
	}
	Message(c, err.Error())
}
