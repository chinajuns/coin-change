package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// AccountLog
// 用户日志
type AccountLog struct {
}

// QueryAccountLogList
// 获取用户日志列表
func (a *AccountLog) QueryAccountLogList(c *gin.Context) {
	lang := c.GetHeader("lang")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	account := c.Query("account")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	currency := c.Query("currency_type")
	types := c.DefaultQuery("type", "0")
	sign := c.DefaultQuery("sign", "0")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)
	signParseInt, _ := strconv.ParseInt(sign, 10, 64)

	data, total, err := new(model.AccountLog).QueryAccountLogListPage(account, startTime, endTime, currency, types, int(signParseInt), int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AccountLog).QueryAccountLogListPage [ERROR] : %s", err))
		return
	}

	var sum interface{}

	if types != "" {
		var num float64
		for _, v := range data {
			valueParseFloat, _ := strconv.ParseFloat(v["value"].(string), 64)
			num += valueParseFloat
		}
		sum = num
	} else {
		sum = "选择日志类型进行统计"
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"total":   total,
		"sum":     sum,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
