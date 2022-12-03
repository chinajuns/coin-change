package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// PayableCurrencies
// 取允许支付的币种
func PayableCurrencies(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)

	currency, err := new(model.Currency).QueryCurrencyJoinMicroNumber()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyJoinMicroNumber [ERROR] : %s", err))
		return
	}
	for i, curr := range currency {
		// 险种
		insuranceType, err := new(model.InsuranceType).QueryInsuranceTypeByCurrencyId(curr["id"])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("new(model.InsuranceType).QueryInsuranceTypeByCurrencyId [ERROR] : %s", err))
			return
		}
		currency[i]["insurance_types"] = insuranceType

		// 用户钱包
		userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userInfo.Id, curr["id"])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] : %s", err))
			return
		}
		currency[i]["user_wallet"] = userWallet

		// 用户保险
		usersInsurance, err := new(model.UsersInsurance).QueryUsersInsuranceByUserIdAndCurrencyId(userInfo.Id, curr["id"])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersInsurance).QueryUsersInsuranceByUserId [ERROR] : %s", err))
			return
		}
		currency[i]["user_insurance"] = usersInsurance
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    currency,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
