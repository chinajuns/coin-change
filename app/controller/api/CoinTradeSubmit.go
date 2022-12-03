package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/app/service"
	"okc/utils"
	"strconv"
)

// CoinTradeSubmit
// 币币交易订单提交
func CoinTradeSubmit(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	legalId := c.PostForm("legal_id")
	currencyId := c.PostForm("currency_id")
	targetPrice := c.PostForm("target_price")
	types := c.PostForm("type")
	amount := c.PostForm("amount")

	currencyMatches, err := new(model.CurrencyMatches).QueryCurrencyMatchesByCurrencyIdAndLegalId(currencyId, legalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.CurrencyMatches).QueryCurrencyMatchesByCurrencyIdAndLegalId [ERROR] : %s", err))
		return
	}
	if currencyMatches.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "找不到交易对",
		})
		return
	}
	if currencyMatches.OpenCoinTrade == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "您未开通本交易对的交易功能",
		})
		return
	}

	currencyIdParseInt, _ := strconv.ParseInt(currencyId, 10, 64)
	legalIdIdParseInt, _ := strconv.ParseInt(legalId, 10, 64)

	switch types {
	case "1":
		err := service.UserBuyCoint(userInfo.Id, int(currencyIdParseInt), int(legalIdIdParseInt), amount, targetPrice)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": err,
			})
			return
		}
		break
	case "2":
		err := service.UserSellCoin(userInfo.Id, int(currencyIdParseInt), int(legalIdIdParseInt), amount, targetPrice)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": err,
			})
			return
		}
		break
	default:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
		break
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
