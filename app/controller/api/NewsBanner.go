package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// QueryNewsBanner
// 获取新闻轮播图
func QueryNewsBanner(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"type": "ok",
		"data": []string{
			"http://lucasie.info:8084/storage/images/705c6322007c09fff25258dbb9cd0e9f.jpeg",
			"http://lucasie.info:8084/storage/images/b844cd154d050df6ba25068607be4bbc.jpeg",
			"http://lucasie.info:8084/storage/images/fce9c533dd300267a5ef155ba423691f.jpeg",
		},
		"message": "success",
	})

	return
}
