package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"okc/app/model"
	"okc/app/service"
	"okc/utils"
)

// LeverDeal
// 合约交易信息
func LeverDeal(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	legalId := c.PostForm("legal_id")
	currencyId := c.PostForm("currency_id")

	if currencyId == "" || legalId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

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
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.Error),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("币种id:%s, 法币id:%s 没有交易对", currencyId, legalId))
		return
	}

	leverShareLimit := make(map[string]interface{})
	leverShareLimit["min"] = currencyMatches.LeverMinShare
	leverShareLimit["max"] = currencyMatches.LeverMaxShare

	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyById [ERROR] : %s", err))
		return
	}

	leverTransaction, err := new(model.LeverTransaction).QueryLeverTransactionByUserIdAndCurrencyIdAndLegalId(userInfo.Id, currencyId, legalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryLeverTransactionByUserIdAndCurrencyIdAndLegalId [ERROR] : %s", err))
		return
	}

	quotation := new(service.CurrencyQuotationStruct)

	mongo := utils.Mongo
	err = mongo.Collection(fmt.Sprintf("QHOTATION-%s", currency.Name)).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(quotation)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("mongo.Collection [ERROR] : %s", err))
		return
	}

	lastPrice := quotation.Last

	userLeverBalance := 0.0

	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userInfo.Id, currencyId)
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

	allLeverMoneyStr, err := new(model.LeverTransaction).QueryLeverTransactionAllMoneyByUserIdAndCurrencyIdAndLegalId(userInfo.Id, currencyId, legalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryLeverTransactionAllMoneyByUserIdAndCurrencyIdAndLegalId [ERROR] : %s", err))
		return
	}

	lastLerverT, err := new(model.LeverTransaction).QueryLeverTransactionByCurrencyIdAndLegalIdAndLimit(currencyId, legalId, 5)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryLeverTransactionByCurrencyIdAndLegalIdAndLimit [ERROR] : %s", err))
		return
	}

	ustdPrice := 1
	exRAte, err := new(model.Setting).QueryValueByKey("USDTRate", "6.5")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Setting).QueryValueByKey [ERROR] : %s", err))
		return
	}

	multiple, err := new(model.LeverMultiple).QueryLeverMultipleByCurrencyId(currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverMultiple).QueryLeverMultipleByCurrencyId [ERROR] : %s", err))
		return
	}

	data := make(map[string]interface{})
	data["lever_transaction"] = lastLerverT     // 最近几条
	data["my_transaction"] = leverTransaction   // 包含用户
	data["last_price"] = lastPrice              // 最新价格
	data["user_lever"] = userLeverBalance       // 用户杠杠余额
	data["all_levers"] = allLeverMoneyStr       // 全部杠杠
	data["ustd_price"] = ustdPrice              //
	data["ExRAte"] = exRAte                     //
	data["multiple"] = multiple                 // 合约倍数
	data["lever_share_limit"] = leverShareLimit //

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
