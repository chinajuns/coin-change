package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// CoinTradeList
// 币币交易列表
func CoinTradeList(c *gin.Context) {
	lang := c.GetHeader("lang")
	limit := c.Query("limit")
	page := c.Query("page")
	currencyId := c.Query("currency_id")
	legalId := c.Query("legal_id")
	status := c.Query("status")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)

	coin, err := new(model.CoinTrade).QueryCoinTradePageByCurrencyIdAndLegalIdAndStatus(currencyId, legalId, status, int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.CoinTrade).QueryCoinTradePageByCurrencyIdAndLegalIdAndStatus [ERROR] : %s", err))
		return
	}
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyById [ERROR] : %s", err))
		return
	}
	for _, m := range coin {
		m["symbol"] = currency.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    coin,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
