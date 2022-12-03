package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UserCenter2
// 获取用户信息(高级认证)
func UserCenter2(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
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

	if user.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.UserFindError),
		})
		return
	}
	userReal, err := new(model.UserReal).QueryFirstDataByUserIdAndTypes(user.Id, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UserReal).QueryFirstDataByUserId(user.Id) [ERROR] : %s", err))
		return
	}
	data := make(map[string]interface{})
	data["id"] = user.Id
	data["phone"] = user.Phone
	data["email"] = user.Email
	data["account_number"] = user.AccountNumber
	if userReal.Id == 0 {
		data["review_status"] = 0
	} else {
		data["review_status"] = userReal.ReviewStatus
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": data,
		},
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
