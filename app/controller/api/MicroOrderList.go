package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// MicroOrderList
// 秒合约（期权）交易列表
func MicroOrderList(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	status := c.DefaultQuery("status", "-1")
	matchId := c.DefaultQuery("match_id", "-1")
	currencyId := c.DefaultQuery("currency", "-1")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)
	statusParseInt, _ := strconv.ParseInt(status, 10, 64)
	matchIdParseInt, _ := strconv.ParseInt(matchId, 10, 64)
	currencyIdParseInt, _ := strconv.ParseInt(currencyId, 10, 64)

	data, total, err := new(model.MicroOrder).PageByUserIdAndStatusAndMatchIdAndCurrencyId(userInfo.Id,
		int(statusParseInt), int(matchIdParseInt), int(currencyIdParseInt), int(pageParseInt), int(limitParseInt))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.MicroOrder).PageByUserIdAndStatusAndMatchIdAndCurrencyId [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"total":   total,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
