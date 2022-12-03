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

// LeverMyTrade
// 我的合约订单
func LeverMyTrade(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	currencyId := c.DefaultPostForm("currency_id", "0")
	legalId := c.DefaultPostForm("legal_id", "0")
	page := c.DefaultPostForm("page", "1")
	limit := c.DefaultPostForm("limit", "10")

	userLeverBalance := 0.0

	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userInfo.Id, 3)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] : %s", err))
		return
	}

	if userWallet.Id != 0 {
		userLeverBalance = userWallet.LeverBalance
	}

	var profitsAll string            // 交易对盈亏总额
	var cautionMoneyAll string       // 交易对可用本金总额
	var originCautionMoneyAll string // 交易对原始保证金

	profitsAll, cautionMoneyAll, originCautionMoneyAll, err = new(model.LeverTransaction).QueryUserProfitByUserIdAndLegalId(userInfo.Id, legalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryUserProfitByUserIdAndLegalId [ERROR] : %s", err))
		return
	}
	profitsAll = utils.TrimRightZeroByFloatStr(profitsAll)
	cautionMoneyAll = utils.TrimRightZeroByFloatStr(cautionMoneyAll)
	originCautionMoneyAll = utils.TrimRightZeroByFloatStr(originCautionMoneyAll)

	var profits string            // 交易对盈亏总额
	var cautionMoney string       // 交易对可用本金总额
	var originCautionMoney string // 交易对原始保证金
	profits, cautionMoney, originCautionMoney, err = new(model.LeverTransaction).QueryUserProfitByUserIdAndLegalIdAndCurrencyId(userInfo.Id, legalId, currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryUserProfitByUserIdAndLegalId [ERROR] : %s", err))
		return
	}

	profits = utils.TrimRightZeroByFloatStr(profits)
	cautionMoney = utils.TrimRightZeroByFloatStr(cautionMoney)
	originCautionMoney = utils.TrimRightZeroByFloatStr(originCautionMoney)

	hazardRate, err := service.QueryWalletHazardRate(userWallet)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("service.QueryWalletHazardRate [ERROR] : %s", err))
		return
	}
	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)
	leverList, err := new(model.LeverTransaction).QueryLeverTransactionPageByCurrencyIdAndLegalId(currencyId, legalId, int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryLeverTransactionPageByCurrencyIdAndLegalId [ERROR] : %s", err))
		return
	}

	rateProfitsTotal := make(map[string]interface{})
	rateProfitsTotal["hazard_rate"] = hazardRate
	rateProfitsTotal["profits_total"] = profitsAll

	data := make(map[string]interface{})
	data["rate_profits_total"] = rateProfitsTotal
	data["lever_balance"] = userLeverBalance
	data["message"] = leverList

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
