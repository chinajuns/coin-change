package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// CurrencyMatchesFind
// 获取交易对详情
func CurrencyMatchesFind(c *gin.Context) {
	lang := c.GetHeader("lang")
	legalId := c.Query("legal_id")
	currencyId := c.Query("currency_id")

	if legalId == "" || currencyId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	currencyMatches, err := new(model.CurrencyMatches).QueryCurrencyMatchesByCurrencyIdAndLegalIdAndIsDisplay(currencyId, legalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.CurrencyMatches).QueryCurrencyMatchesByCurrencyIdAndLegalIdAndIsDisplay [ERROR] : %s \n", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    currencyMatches,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
