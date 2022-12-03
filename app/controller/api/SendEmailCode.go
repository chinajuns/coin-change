package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/utils"
)

// SendEmailCode
// 发送邮箱验证码
func SendEmailCode(c *gin.Context) {
	lang := c.PostForm("lang")
	email := c.PostForm("user_string")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}
	code := utils.GenerateCode()

	// 缓存验证码60秒
	redis := utils.RedisConnect()
	defer redis.Close()
	key := "EMAIL_CODE_" + email
	_, err := redis.Do("SetEx", key, 60, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "500",
			"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
		})
		utils.WriteErrorLog(err)
		return
	}
	err = utils.SendEmailCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "500",
			"message": utils.GetLangMessage(lang, utils.SendCodeError) + err.Error(),
		})
		utils.WriteErrorLog(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
