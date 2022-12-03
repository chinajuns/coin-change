package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UserLogin
// 用户登录
func UserLogin(c *gin.Context) {
	// 语言
	lang := c.PostForm("lang")
	// 手机、邮箱、交易账号
	userString := c.PostForm("user_string")
	// 密码
	password := c.PostForm("password")

	if userString == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.UserNameInputError),
		})
		return
	}

	if password == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.PasswordInputError),
		})
		return
	}

	// 查询用户
	user, err := new(model.Users).QueryUserInfoByUserAccount(userString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "500",
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Users).QueryUserInfoByUserAccount [ERROR] : %s \n", err))
		return
	}
	// 用户不存在
	if user.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    "404",
			"message": utils.GetLangMessage(lang, utils.UserFindError),
		})
		return
	}
	// 密码错误
	if utils.GenerateSha256(password) != user.Password {
		c.JSON(http.StatusOK, gin.H{
			"type":    "404",
			"message": utils.GetLangMessage(lang, utils.PasswordError),
		})
		return
	}
	// 检查是否被冻结
	if user.Status == 1 {
		c.JSON(http.StatusOK, gin.H{
			"type":    "404",
			"message": utils.GetLangMessage(lang, utils.UserFindError),
		})
		return
	}
	// 生成token
	token, err := utils.GenerateToken(user.Id, user.AccountNumber, user.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "500",
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("/api/user/login [ERROR] : %s \n", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
		"data": gin.H{
			"token": token,
		},
	})
	return
}
