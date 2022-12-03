package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"okc/config"
	"okc/utils"
	"os"
	"strings"
	"time"
)

// UploadOne
// 单文件上传
func UploadOne(c *gin.Context) {

	lang := c.GetHeader("lang")
	file, _ := c.FormFile("file")
	fileName := file.Filename
	if strings.Index(fileName, ".") == -1 {
		log.Println("文件名称错误")
		return
	}

	fileSuffix := strings.Split(fileName, ".")[1]
	filePrefix := utils.GenerateMd5(time.Now().Format("2006-01-02 15:04:05"))
	fileName = fmt.Sprintf("%s.%s", filePrefix, fileSuffix)

	_, err := os.Stat("./storage/images")
	if err != nil {
		os.MkdirAll("./storage/images", 0755)
	}
	m := config.Config().APP.(map[interface{}]interface{})
	path := fmt.Sprintf("%s/storage/images/%s", m["APP_URL"], fileName)
	dst := "./storage/images/" + fileName
	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("c.SaveUploadedFile [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type": "ok",
		"data": gin.H{
			"src": path,
		},
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
