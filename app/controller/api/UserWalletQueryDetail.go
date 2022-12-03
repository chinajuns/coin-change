package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// UserWalletQueryDetail
// 用户钱包详情
func UserWalletQueryDetail(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)

	// 币种id
	CurrencyId := c.PostForm("currency")
	// 类型
	//types := c.PostForm("type")

	if CurrencyId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	// TODO: 未知变量
	exRate, err := new(model.Setting).QueryValueByKey("USDTRate", "6.5")
	if err != nil {
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Setting).QueryValueByKey [ERROR] : %s", err))
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})

		return
	}
	currencyIdParseInt, _ := strconv.ParseInt(CurrencyId, 10, 64)

	userWallet, err := new(model.UsersWallet).QueryMapUserWalletByUserIdAndCurrencyId(userInfo.Id, int(currencyIdParseInt))
	if err != nil {
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] : %s", err))
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})

		return
	}

	coin, _ := new(model.Setting).QueryValueByKey("COIN_TRADE_FEE", "")
	userWallet["coin_trade_fee"] = coin

	userWallet["ExRate"] = exRate

	if utils.InCheckIntOrIntSlice(userWallet["currency"].(int), []int{1, 2, 3}) {
		userWallet["is_charge"] = true
	} else {
		userWallet["is_charge"] = false
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    userWallet,
		"message": utils.GetLangMessage(lang, utils.Success),
	})

	return
}
