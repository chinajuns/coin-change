package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// MicroSecondList
// 秒合约（期权）时间列表
func MicroSecondList(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	_ = Cache.(*utils.UserClaims)

	data, err := new(model.MicroSecond).QueryMicroSecondList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.MicroSecond).QueryMicroSecondList [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
