package utils

import (
	"github.com/gin-gonic/gin"
)

// gin response json
func Response(c *gin.Context, errcode int, errmessage string, data interface{}) {
	news := gin.H{
		"errcode":    errcode,
		"errmessage": errmessage,
		"data":       data,
	}

	c.JSON(200, news)
}

// gin response page json
func ResposePage(c *gin.Context, errcode int, errmessage string, data interface{}, total int) {
	news := gin.H{
		"errcode":    errcode,
		"errmessage": errmessage,
		"data":       data,
		"total":      total,
	}

	c.JSON(200, news)
}
