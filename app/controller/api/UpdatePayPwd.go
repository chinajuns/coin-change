package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UpdatePayPwd
// 修改交易密码
func UpdatePayPwd(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	// 旧密码
	oldPassword := c.PostForm("oldpassword")
	// 新密码
	password := c.PostForm("password")
	// 再次输入密码
	rePassword := c.PostForm("re_password")

	// 判断密码长度
	if len(password) < 6 || len(password) > 16 {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": utils.GetLangMessage(lang, utils.PasswordLenError),
		})
		return
	}
	// 判断密码是否一致
	if password != rePassword {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": utils.GetLangMessage(lang, utils.PasswordHshError),
		})
		return
	}
	// 查询用户
	user, err := new(model.Users).QueryUserInfoById(userInfo.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Users).QueryUserInfoById(userInfo.Id) [ERROR] : %s", err))
		return
	}
	// 检查用户是否存在
	if user.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.UserFindError),
		})
		return
	}
	// 判断旧密码是否正确
	if utils.GenerateSha256(oldPassword) != user.PayPassword {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.PasswordError),
		})
		return
	}
	// 修改交易密码
	err = new(model.Users).UpdatePayPwdById(user.Id, utils.GenerateSha256(password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.UpdateError),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.UpdateSuccess),
	})
	return
}
