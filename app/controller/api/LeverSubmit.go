package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"okc/app/model"
	"okc/app/service"
	"okc/utils"
	"strconv"
	"time"
)

// LeverSubmit
// 合约提交订单
func LeverSubmit(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	share := c.PostForm("share")
	multiple := c.PostForm("multiple")
	types := c.PostForm("type")
	legalId := c.PostForm("legal_id")
	currencyId := c.PostForm("currency_id")
	status := c.DefaultPostForm("status", "1")
	targetPrice := c.PostForm("target_price")

	currencyIdParseInt, _ := strconv.ParseInt(currencyId, 10, 64)
	legalIdParseInt, _ := strconv.ParseInt(legalId, 10, 64)

	currencyMatches, err := new(model.CurrencyMatches).QueryCurrencyMatchesByCurrencyIdAndLegalId(currencyId, legalIdParseInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.CurrencyMatches).QueryCurrencyMatchesByCurrencyIdAndLegalId [ERROR] : %s ", err))
		return
	}

	if currencyMatches.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "指定交易对不存在",
		})
		return
	}

	if currencyMatches.OpenLever != 1 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "您未开通本交易对的交易功能",
		})
		return
	}

	shareParseInt, _ := strconv.ParseInt(share, 10, 64)

	if shareParseInt <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "手数必须是大于0的整数",
		})
		return
	}
	if currencyMatches.LeverMinShare > int(shareParseInt) {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": fmt.Sprintf("手数不能低于 : %d", currencyMatches.LeverMinShare),
		})
		return
	}

	if currencyMatches.LeverMaxShare < int(shareParseInt) {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": fmt.Sprintf("手数不能高于 : %d", currencyMatches.LeverMaxShare),
		})
		return
	}

	leverMultiple, err := new(model.LeverMultiple).QueryLeverMultipleByValue(multiple)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverMultiple).QueryLeverMultipleByValue [ERROR] : %s", err))
		return
	}
	if leverMultiple.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "选择倍数不在系统范围",
		})
		return
	}

	existsCloseTradCount, err := new(model.LeverTransaction).QueryLeverTransactionCountByUserIdAndStatus(userInfo.Id, status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryLeverTransactionCountByUserIdAndStatus [ERROR] : %s", err))
		return
	}

	//log.Println("existsCloseTradCount : ", existsCloseTradCount)

	if existsCloseTradCount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "您有正在平仓中的交易,暂不能进行买卖",
		})
		return
	}

	statusParseInt, _ := strconv.ParseInt(status, 10, 64)
	if !utils.InCheckIntOrIntSlice(int(statusParseInt), []int{0, 1}) {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "交易类型错误",
		})
		return
	}
	if statusParseInt == 0 {
		openLeverEntrust, err := new(model.Setting).QueryValueByKey("open_lever_entrust", "0")
		if err != nil {
			log.Println(fmt.Sprintf("new(model.Setting).QueryValueByKey [ERROR] : %s", err))
			return
		}
		openLeverEntrustParseInt, _ := strconv.ParseInt(openLeverEntrust, 10, 64)
		if openLeverEntrustParseInt <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "该功能暂未开放",
			})
			return
		}
	}

	targetPriceParseInt, _ := strconv.ParseInt(targetPrice, 10, 64)
	if statusParseInt == 0 && targetPriceParseInt <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "限价交易价格必须大于0",
		})
		return
	}

	// TODO: 未知变量
	var overnight float64
	var originPrice string
	var leverShareNum float64
	var spreadPrice float64
	var factPrice float64
	var allMoney float64

	if currencyMatches.Overnight != 0 {
		overnight = currencyMatches.Overnight
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

	lastPrice, err := service.QueryLastQuotationPriceByCurrencyName(currency.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("service.QueryLastQuotationPriceByCurrencyName [ERROR] : %s", err))
		return
	}

	typesParseInt, _ := strconv.ParseInt(types, 10, 64)
	if statusParseInt == 0 {

		if typesParseInt == 2 && targetPrice >= lastPrice {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "限价交易卖出不能低于当前价",
			})
			return
		} else if typesParseInt == 1 && targetPrice >= lastPrice {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "限价交易买入价格不能高于当前价",
			})
			return
		}
		originPrice = targetPrice
	} else {
		originPrice = lastPrice
	}

	if currencyMatches.LeverShareNum == 0 {
		leverShareNum = 1
	} else {
		leverShareNum = currencyMatches.LeverShareNum
	}
	leverShareNumParseStr := strconv.FormatFloat(leverShareNum, 'f', 4, 64)
	num, err := utils.BcMul(share, leverShareNumParseStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcMul(share, leverShareNumParseStr) [ERROR] : %s", err))
		return
	}

	spreadPrice = currencyMatches.Spread
	spreadPriceParseStr := strconv.FormatFloat(spreadPrice, 'f', 4, 64)

	if typesParseInt == 2 {
		spreadPriceFloat, err := utils.BcMul("-1", spreadPriceParseStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcMul(\"-1\", spreadPriceParseStr) : %s", err))
			return
		}
		spreadPriceFloatParse := strconv.FormatFloat(spreadPriceFloat, 'f', 4, 64)
		spreadPriceParseStr = spreadPriceFloatParse
	}

	factPrice, err = utils.BcAdd(originPrice, spreadPriceParseStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcAdd(originPrice, spreadPriceParseStr) [ERROR] : %s", err))
		return
	}
	factPriceParseStr := strconv.FormatFloat(factPrice, 'f', 4, 64)
	numParseStr := strconv.FormatFloat(num, 'f', 4, 64)
	allMoney, err = utils.BcMul(factPriceParseStr, numParseStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcMul(factPriceParseStr, numParseStr) [ERROR] : %s", err))
		return
	}

	var leverTradeFeeRate float64
	var tardFee float64
	leverTradeFeeRateParseStr := strconv.FormatFloat(currencyMatches.LeverTradeFee, 'f', 4, 64)

	if currencyMatches.LeverTradeFee > 0 {
		leverTradeFeeRate, err = utils.BcDiv(leverTradeFeeRateParseStr, "100")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcDiv(leverTradeFeeRateParseStr, \"100\") [ERROR] : %s", err))
			return
		}
	}

	allMoneyParseStr := strconv.FormatFloat(allMoney, 'f', 4, 64)
	leverTradeFeeRateParseStr = strconv.FormatFloat(leverTradeFeeRate, 'f', 4, 64)

	tardFee, err = utils.BcMul(allMoneyParseStr, leverTradeFeeRateParseStr)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcMul(allMoneyParseStr, leverTradeFeeRateParseStr) [ERROR] : %s", err))
		return
	}

	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userInfo.Id, legalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf(" new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] : %s", err))
		return
	}
	if userWallet.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "钱包未找到,请先添加钱包",
		})
		return
	}

	cautionMoney, err := utils.BcDiv(allMoneyParseStr, multiple)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf(" utils.BcDiv(allMoneyParseStr, multiple) [ERROR] : %s", err))
		return
	}

	cautionMoneyParseStr := strconv.FormatFloat(cautionMoney, 'f', 4, 64)
	tardFeeParseStr := strconv.FormatFloat(tardFee, 'f', 4, 64)

	//log.Println("cautionMoneyParseStr : ", cautionMoneyParseStr)
	//log.Println("tardFeeParseStr : ", tardFeeParseStr)

	shoudDeduct, err := utils.BcAdd(cautionMoneyParseStr, tardFeeParseStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcAdd(cautionMoneyParseStr, tardFeeParseStr) [ERROR] : %s", err))
		return
	}

	//log.Println("tardFeeParseStr : ", tardFeeParseStr)
	//log.Println("tardFee : ", tardFee)
	//log.Println("cautionMoneyParseStr : ", cautionMoneyParseStr)
	//log.Println("LegalBalance : ", userWallet.LegalBalance)
	//log.Println("shoudDeduct : ", shoudDeduct)

	if userWallet.LegalBalance < shoudDeduct {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "余额不足",
		})
		return
	}
	originPriceParseFloat, _ := strconv.ParseFloat(originPrice, 64)
	lastPriceParseFloat, _ := strconv.ParseFloat(lastPrice, 64)
	multipleParseInt, _ := strconv.ParseInt(multiple, 10, 64)

	leverTransaction := &model.LeverTransaction{
		Type:               int(typesParseInt),
		UserId:             userInfo.Id,
		Currency:           int(currencyIdParseInt),
		Legal:              int(legalIdParseInt),
		OriginPrice:        originPriceParseFloat,
		Price:              factPrice,
		UpdatePrice:        lastPriceParseFloat,
		Share:              int(shareParseInt),
		Number:             num,
		Multiple:           int(multipleParseInt),
		OriginCautionMoney: cautionMoney,
		CautionMoney:       cautionMoney,
		TradeFee:           tardFee,
		Overnight:          overnight,
		Status:             int(statusParseInt),
		CreateTime:         int(time.Now().Unix()),
		TransactionTime:    "0.00",
		UpdateTime:         "0.00",
		HandleTime:         "0.00",
		CompleteTime:       "0.00",
	}
	id, err := leverTransaction.AddLeverTransaction()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": fmt.Sprintf("提交失败 : %s", err),
		})
		return
	}

	extraData := map[string]interface{}{
		"trade_id":  id,
		"all_money": allMoney,
		"multiple":  multiple,
	}
	extraDataJson, _ := json.Marshal(extraData)

	res := service.ChangeUserWalletBalance(userWallet, 3, cautionMoney*-1, 30, "提交杠杆交易,扣除保证金", false, 0, 0, string(extraDataJson), 0, 0)
	switch res.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(error),
		})
		return
		break
	case string:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(string),
		})
		return
		break
	case bool:
		if res.(bool) == false {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "扣除保证金失败",
			})
			return
		}
		break
	}

	extraData = map[string]interface{}{
		"trade_id":             id,
		"all_money":            allMoney,
		"lever_trade_fee_rate": leverTradeFeeRate,
	}
	extraDataJson, _ = json.Marshal(extraData)
	res = service.ChangeUserWalletBalance(userWallet, 3, tardFee*-1, 35, "提交杠杆交易,扣除保证金", false, 0, 0, string(extraDataJson), 0, 0)
	switch res.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(error),
		})
		return
		break
	case string:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(string),
		})
		return
		break
	case bool:
		if res.(bool) == false {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "扣除手续费失败",
			})
			return
		}
		break
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "提成成功",
	})
	return
}
