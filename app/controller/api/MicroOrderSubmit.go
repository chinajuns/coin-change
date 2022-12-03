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
	"strconv"
	"time"
)

// MicroOrderSubmit
// 期权订单提交
func MicroOrderSubmit(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	matchId := c.PostForm("match_id")       // 交易对id
	currencyId := c.PostForm("currency_id") // 币种id
	types := c.PostForm("type")             // 类型
	second := c.PostForm("seconds")         // 时间
	number := c.PostForm("number")          // 数量

	// 转换
	numberParseFloat, _ := strconv.ParseFloat(number, 64)
	numberParseInt, _ := strconv.ParseInt(number, 10, 64)
	matchIdParseInt, _ := strconv.ParseInt(matchId, 10, 64)
	currencyIdParseInt, _ := strconv.ParseInt(currencyId, 10, 64)
	typesParseInt, _ := strconv.ParseInt(types, 10, 64)
	secondParseInt, _ := strconv.ParseInt(second, 10, 64)

	// 期权时间
	microSecond, err := new(model.MicroSecond).QueryMicroSecondBySecond(second)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.MicroSecond).QueryMicroSecondBySecond [ERROR] : %s", err))
		return
	}
	if microSecond.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "期权到期时间不存在",
		})
		_ = utils.WriteErrorLog("期权到期时间不存在")
		return
	}
	if int(numberParseInt) < microSecond.MinNum {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "下注金额不能小于最小金额",
		})
		_ = utils.WriteErrorLog("下注金额不能小于最小金额")
		return
	}

	// 保险开始时间
	setting := new(model.Setting)
	insuranceStart, err := setting.QueryValueByKey("insurance_start", "09:00")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("setting.QueryValueByKey [ERROR] :  : %s", err))
		return
	}

	// 保险结束时间
	insuranceEnd, err := setting.QueryValueByKey("insurance_end", "12:00")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("setting.QueryValueByKey [ERROR] :  : %s", err))
		return
	}
	now := time.Now()
	// 解析时间变成unix
	insuranceStartTime, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s:00", now.Format("2006-01-02"), insuranceStart))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("time.Parse [ERROR] :  : %s", err))
		return
	}
	insuranceEndTime, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s:00", now.Format("2006-01-02"), insuranceEnd))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("time.Parse [ERROR] :  : %s", err))
		return
	}

	// 当前时间戳
	unix := now.Unix()
	// 受保
	userInsuranceType := 0
	// 币种
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyById [ERROR] : %s", err))
		return
	}
	// 受保护时间
	if unix >= insuranceStartTime.Unix() && unix <= insuranceEndTime.Unix() {
		// 检查是否可以下单
		if currency.Insurancable == 1 {
			checkRes := service.CheckMicroInsurance(userInfo.Id, currencyId, int(numberParseInt))
			switch checkRes.(type) {
			case string:
				c.JSON(http.StatusOK, gin.H{
					"type":    500,
					"message": fmt.Sprintf("下单失败: %s", checkRes.(string)),
				})
				return
				break
			case error:
				c.JSON(http.StatusOK, gin.H{
					"type":    500,
					"message": fmt.Sprintf("下单失败: %s", checkRes.(error)),
				})
				return
				break
			case bool:
				break
			}

			userInsuranceSlice, err := new(model.UsersInsurance).QueryUsersInsuranceByUserIdAndCurrencyId(userInfo.Id, currencyId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"type":    500,
					"message": utils.GetLangMessage(lang, utils.ServerError),
				})
				_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersInsurance).QueryUsersInsuranceByUserIdAndCurrencyId [ERROR] : %s", err))
				return
			}
			userInsurance := userInsuranceSlice[0]
			// 受保
			userInsuranceType = userInsurance["type"].(int)
		}

	}

	// 挂单数量
	microOrderCount, err := new(model.MicroOrder).QueryMicroOrderCountByUserIdAndCurrencyId(userInfo.Id, currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.MicroOrder).QueryMicroOrderCountByUserIdAndCurrencyId [ERROR] : %s", err))
		return
	}

	// 如果当前不在受保时间段内或者所返币种不支持保险
	if currency.Insurancable != 1 && userInsuranceType == 0 && currency.MicroHoldtradeMax > 0 && microOrderCount >= currency.MicroHoldtradeMax {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "下单失败:超过最大持仓笔数限制",
		})
		_ = utils.WriteErrorLog("下单失败:超过最大持仓笔数限制")
		return
	}

	// 币种交易对
	currencyMatches, err := new(model.CurrencyMatches).QueryCurrencyMatchesByid(matchId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.CurrencyMatches).QueryCurrencyMatchesByid [ERROR] : %s", err))
		return
	}

	// 获取币种k线
	kline := new(service.CurrencyKlineStruct)
	mongo := utils.Mongo
	err = mongo.Collection(fmt.Sprintf("KLINE-%s", currency.Name)).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(&kline)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("mongo.Collection(currency.Name).FindOne [ERROR] : %s", err))
		return
	}

	// 用户钱包
	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userInfo.Id, currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] :  %s", err))
		return
	}
	if userWallet.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.UserWalletFindError),
		})
		return
	}
	// k线收盘价格
	price := kline.C
	priceParseFloat, _ := strconv.ParseFloat(price, 64)

	// 场控
	FluctuateMinStr := strconv.FormatFloat(currencyMatches.FluctuateMin, 'f', 15, 64)
	FluctuateMaxStr := strconv.FormatFloat(currencyMatches.FluctuateMax, 'f', 15, 64)

	decimalNum := utils.CheckDecimalLenByString(FluctuateMinStr)
	randFloat, err := utils.MatchRandFloatBySectionAndPec(FluctuateMinStr, FluctuateMaxStr, decimalNum)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("utils.MatchRandFloatBySectionAndPec [ERROR] : %s", err))
		return
	}
	if utils.MatchRandBoolSlicePop([]bool{true, false}) {
		priceParseFloat, err = utils.BcAdd(price, randFloat)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteInfoLog(fmt.Sprintf("场控增加: 原值(%s), 现值(%f), 精度(%s)", price, priceParseFloat, randFloat))
			_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcAdd [ERROR] : %s", err))
			return
		}
	} else {
		priceParseFloat, err = utils.BcSub(price, randFloat)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteInfoLog(fmt.Sprintf("场控减少: 原值(%s), 现值(%f), 精度(%s)", price, priceParseFloat, randFloat))
			_ = utils.WriteErrorLog(fmt.Sprintf("utils.BcAdd [ERROR] : %s", err))
			return
		}
	}

	// 类型
	var balanceType int
	if userInsuranceType != 0 {
		balanceType = 5
	} else {
		balanceType = 4
	}
	// 扣除本金
	changeRes := service.ChangeUserWalletBalance(userWallet, balanceType, numberParseFloat*-1, 502, "秒合约下单扣除本金", false, 0, 0, "", 0, 0)
	switch changeRes.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": changeRes.(error),
		})
		return
		break
	case string:
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": changeRes.(string),
		})
		return
		break
	case bool:
		if changeRes.(bool) == false {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": "钱包变更余额异常",
			})
			return
		}
		break
	}
	// 扣除手续费
	fee := currency.MicroTradeFee
	changeRes = service.ChangeUserWalletBalance(userWallet, balanceType, fee*-1, 502, fmt.Sprintf("秒合约下单扣除: %f 手续费", fee), false, 0, 0, "", 0, 0)
	switch changeRes.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": changeRes.(error),
		})
		return
		break
	case string:
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": changeRes.(string),
		})
		return
		break
	case bool:
		if changeRes.(bool) == false {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": "钱包变更余额异常",
			})
			return
		}
		break
	}

	// 下单
	microOder := &model.MicroOrder{
		UserId:          userInfo.Id,
		MatchId:         int(matchIdParseInt),
		CurrencyId:      int(currencyIdParseInt),
		Type:            int(typesParseInt),
		IsInsurance:     userInsuranceType,
		Seconds:         int(secondParseInt),
		Number:          numberParseFloat,
		OpenPrice:       priceParseFloat,
		EndPrice:        priceParseFloat,
		Fee:             fee,
		ProfitRatio:     microSecond.ProfitRatio,
		Status:          1,
		PreProfitResult: 0,
		CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:       time.Now().Format("2006-01-02 15:04:05"),
		HandledAt:       time.Now().Format("2006-01-02 15:04:05"),
		ReturnAt:        time.Now().Format("2006-01-02 15:04:05"),
	}
	err = microOder.AddMicroOder()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("microOder.AddMicroOder [ERROR] :  %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
