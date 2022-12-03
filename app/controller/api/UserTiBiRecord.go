package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// UserTiBiRecord
// 用户提币记录
func UserTiBiRecord(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	limitInt, _ := strconv.ParseInt(limit, 10, 64)

	data, total, err := new(model.UsersWalletOut).QueryPagesByUserId(userInfo.Id, int(pageInt), int(limitInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWalletOut).QueryPagesByUserId [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"total":   total,
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return

}
