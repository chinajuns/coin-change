package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UserWalletQueryAddress
// 获取用户钱包地址
func UserWalletQueryAddress(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)

	usersAddress, err := new(model.UsersAddress).QueryFindByUserId(userInfo.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersAddress).QueryFindByUserIdAndAddress [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    usersAddress,
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
