package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UpdatePwd
// 修改账户密码
func UpdatePwd(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	// 旧密码
	oldPassword := c.PostForm("old_password")
	// 新密码
	password := c.PostForm("password")
	// 类型 type:1登录密码，type:2支付密码
	types := c.DefaultPostForm("type", "1")

	if password == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": utils.GetLangMessage(lang, utils.PasswordInputError),
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
	// 判断类型
	if types == "1" {
		// 匹配密码是否正确
		if utils.GenerateSha256(oldPassword) != user.Password {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": utils.GetLangMessage(lang, utils.PasswordHshError),
			})
			return
		}

		// 修改密码
		err = new(model.Users).UpdatePwdById(user.Id, utils.GenerateSha256(password))
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
	} else {

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

}
