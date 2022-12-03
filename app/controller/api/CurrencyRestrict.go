package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// CurrencyRestrict
// 币种限制
func CurrencyRestrict(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	currencyId := c.PostForm("currency")

	if currencyId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.CurrencyFindError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyById [ERROR] : %s", err))
		return
	}
	currencyIdParseInt, _ := strconv.ParseInt(currencyId, 10, 64)
	wallet, err := new(model.UsersWallet).QueryMapUserWalletByUserIdAndCurrencyId(userInfo.Id, int(currencyIdParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.UserWalletFindError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] : %s", err))
		return
	}

	data := make(map[string]interface{})
	data["rate"] = currency.Rate
	data["min_number"] = currency.MinNumber
	data["name"] = currency.Name
	data["legal_balance"] = wallet["legal_balance"]
	data["change_balance"] = wallet["change_balance"]

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return

}
