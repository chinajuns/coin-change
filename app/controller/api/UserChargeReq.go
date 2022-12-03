package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// UserChargeReq
// 用户充值
func UserChargeReq(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	currencyId := c.PostForm("currency")
	number := c.DefaultPostForm("account", "1")
	amount := c.PostForm("amount")
	img := c.PostForm("img")

	if amount == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}
	currencyIdParseInt, _ := strconv.ParseInt(currencyId, 10, 64)
	amountParseFloat, _ := strconv.ParseFloat(amount, 64)
	chargeReq := &model.ChargeReq{
		Uid:         userInfo.Id,
		Amount:      amountParseFloat,
		UserAccount: number,
		Status:      1,
		CurrencyId:  int(currencyIdParseInt),
		IsBank:      1,
		Img:         img,
	}
	err := chargeReq.AddChargeReq()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("chargeReq.AddChargeReq() [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
