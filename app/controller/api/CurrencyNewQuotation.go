package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// NewQuotation
// 获取币种和行情
func NewQuotation(c *gin.Context) {
	lang := c.GetHeader("lang")
	currency, err := new(model.Currency).QueryCurrencyJoinQuotation()
	if err != nil {
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyJoinQuotation() [ERROR] : %s", err))
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type": "ok",
		"message": gin.H{
			"quotation": currency,
		},
	})
	return
}
