package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UserCheckHasPayPassword
// 检查用户是否设置了交易密码
func UserCheckHasPayPassword(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)

	user, err := new(model.Users).QueryUserInfoById(userInfo.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.PayPasswordNilError),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": 1,
	})
	return
}
