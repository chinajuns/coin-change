package admin

import (
	"fmt"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminCurrency struct {
}

// CurrencyLists
// 币种管理列表
func (m *AdminCurrency) CurrencyLists(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)

	data, total, err := new(model.Currency).QueryCurrencyPage(int(pageParseInt), int(limitParseInt))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("语法错误", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(Currency).QueryCurrencyPage [ERROR] : %s", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":    http.StatusOK,
		"data":    data,
		"total":   total,
		"message": "success",
	})
}

// AddCurrency
// 添加币种
func (m *AdminCurrency) AddCurrency(c *gin.Context) {
	lang := c.GetHeader("lang")
	name := c.PostForm("name")
	sorts := c.DefaultPostForm("sort", "0")
	logo := c.PostForm("log")
	types := c.PostForm("type")
	isLegal := c.DefaultPostForm("is_legal", "0")
	isLever := c.DefaultPostForm("is_lever", "0")
	isMatch := c.DefaultPostForm("is_match", "0")
	isMicro := c.DefaultPostForm("is_micro", "0")
	microMin := c.DefaultPostForm("micro_min", "0")
	microMax := c.DefaultPostForm("micro_max", "0")
	microHoldtraeMax := c.DefaultPostForm("micro_holdtrade_max", "0")
	price := c.DefaultPostForm("price", "0")
	minNumber := c.DefaultPostForm("min_number", "0")
	maxNumber := c.DefaultPostForm("max_number", "0")
	rate := c.DefaultPostForm("rate", "0")
	rmbRelation := c.DefaultPostForm("rmb_relation", "0")
	microTradeFee := c.DefaultPostForm("micro_trade_fee", "0")
	contractAddress := c.DefaultPostForm("contract_address", "0")
	decimalScale := c.DefaultPostForm("decimal_scale", "18")
	chainFee := c.DefaultPostForm("chain_fee", "18")
	insurancable := c.DefaultPostForm("insurancable", "18")

	name = strings.ToUpper(name)

	if name == "" || sorts == "" || types == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	currency, err := new(model.Currency).QueryCurrencyByName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyByName [ERROR] : %s", err))
		return
	}

	if currency.Id > 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": fmt.Sprintf("币种%s已存在", name),
		})
		return
	}
	sortsParseInt, _ := strconv.ParseInt(sorts, 10, 64)
	isLegalParseInt, _ := strconv.ParseInt(isLegal, 10, 64)
	isLeverParseInt, _ := strconv.ParseInt(isLever, 10, 64)
	isMatchParseInt, _ := strconv.ParseInt(isMatch, 10, 64)
	isMicroParseInt, _ := strconv.ParseInt(isMicro, 10, 64)
	minNumberParseFloat, _ := strconv.ParseFloat(minNumber, 64)
	maxNumberParseFloat, _ := strconv.ParseFloat(maxNumber, 64)
	microHoldtraeMaxParseInt, _ := strconv.ParseInt(microHoldtraeMax, 10, 64)
	rateParseFloat, _ := strconv.ParseFloat(rate, 64)
	priceParseFloat, _ := strconv.ParseFloat(price, 64)
	microMinParseFloat, _ := strconv.ParseFloat(microMin, 64)
	microMaxParseFloat, _ := strconv.ParseFloat(microMax, 64)
	rmbRelationParseFloat, _ := strconv.ParseFloat(rmbRelation, 64)
	decimalScaleParseInt, _ := strconv.ParseInt(decimalScale, 10, 64)
	insurancableParseInt, _ := strconv.ParseInt(insurancable, 10, 64)
	chainFeeParseFloat, _ := strconv.ParseFloat(chainFee, 64)
	microTradeFeeParseFloat, _ := strconv.ParseFloat(microTradeFee, 64)

	currency.Name = name
	currency.Sort = int(sortsParseInt)
	currency.Logo = logo
	currency.IsLegal = int(isLegalParseInt)
	currency.IsLever = int(isLeverParseInt)
	currency.IsMatch = int(isMatchParseInt)
	currency.IsMicro = int(isMicroParseInt)
	currency.MinNumber = minNumberParseFloat
	currency.MaxNumber = maxNumberParseFloat
	currency.MicroHoldtradeMax = int(microHoldtraeMaxParseInt)
	currency.Rate = rateParseFloat
	currency.Price = priceParseFloat
	currency.MicroMin = microMinParseFloat
	currency.MicroMax = microMaxParseFloat
	currency.RmbRelation = rmbRelationParseFloat
	currency.DecimalScale = int(decimalScaleParseInt)
	currency.Insurancable = int(insurancableParseInt)
	currency.ChainFee = chainFeeParseFloat
	currency.ContractAddress = contractAddress
	currency.MicroTradeFee = microTradeFeeParseFloat
	currency.Type = types
	currency.TotalAccount = ""
	currency.CollectAccount = ""
	currency.Key = ""
	currency.CurrencyDecimals = 0.0
	currency.IsDisplay = 1

	err = currency.AddCurrency()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("currency.AddCurrency() [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}

// EditCurrency
// 编辑币种
func (m *AdminCurrency) EditCurrency(c *gin.Context) {
	lang := c.GetHeader("lang")
	id := c.PostForm("id")
	name := c.PostForm("name")
	sorts := c.DefaultPostForm("sort", "0")
	logo := c.PostForm("log")
	types := c.PostForm("type")
	isLegal := c.DefaultPostForm("is_legal", "0")
	isLever := c.DefaultPostForm("is_lever", "0")
	isMatch := c.DefaultPostForm("is_match", "0")
	isMicro := c.DefaultPostForm("is_micro", "0")
	microMin := c.DefaultPostForm("micro_min", "0")
	microMax := c.DefaultPostForm("micro_max", "0")
	microHoldtraeMax := c.DefaultPostForm("micro_holdtrade_max", "0")
	price := c.DefaultPostForm("price", "0")
	minNumber := c.DefaultPostForm("min_number", "0")
	maxNumber := c.DefaultPostForm("max_number", "0")
	rate := c.DefaultPostForm("rate", "0")
	rmbRelation := c.DefaultPostForm("rmb_relation", "0")
	microTradeFee := c.DefaultPostForm("micro_trade_fee", "0")
	contractAddress := c.DefaultPostForm("contract_address", "0")
	decimalScale := c.DefaultPostForm("decimal_scale", "18")
	chainFee := c.DefaultPostForm("chain_fee", "18")
	insurancable := c.DefaultPostForm("insurancable", "18")

	name = strings.ToUpper(name)

	if name == "" || sorts == "" || types == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	currency, err := new(model.Currency).QueryCurrencyByName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyByName [ERROR] : %s", err))
		return
	}

	if currency.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": fmt.Sprintf("币种%s不存在", name),
		})
		return
	}

	idParseInt, _ := strconv.ParseInt(id, 10, 64)
	sortsParseInt, _ := strconv.ParseInt(sorts, 10, 64)
	isLegalParseInt, _ := strconv.ParseInt(isLegal, 10, 64)
	isLeverParseInt, _ := strconv.ParseInt(isLever, 10, 64)
	isMatchParseInt, _ := strconv.ParseInt(isMatch, 10, 64)
	isMicroParseInt, _ := strconv.ParseInt(isMicro, 10, 64)
	minNumberParseFloat, _ := strconv.ParseFloat(minNumber, 64)
	maxNumberParseFloat, _ := strconv.ParseFloat(maxNumber, 64)
	microHoldtraeMaxParseInt, _ := strconv.ParseInt(microHoldtraeMax, 10, 64)
	rateParseFloat, _ := strconv.ParseFloat(rate, 64)
	priceParseFloat, _ := strconv.ParseFloat(price, 64)
	microMinParseFloat, _ := strconv.ParseFloat(microMin, 64)
	microMaxParseFloat, _ := strconv.ParseFloat(microMax, 64)
	rmbRelationParseFloat, _ := strconv.ParseFloat(rmbRelation, 64)
	decimalScaleParseInt, _ := strconv.ParseInt(decimalScale, 10, 64)
	insurancableParseInt, _ := strconv.ParseInt(insurancable, 10, 64)
	chainFeeParseFloat, _ := strconv.ParseFloat(chainFee, 64)
	microTradeFeeParseFloat, _ := strconv.ParseFloat(microTradeFee, 64)

	currency.Id = int(idParseInt)
	currency.Name = name
	currency.Sort = int(sortsParseInt)
	currency.Logo = logo
	currency.IsLegal = int(isLegalParseInt)
	currency.IsLever = int(isLeverParseInt)
	currency.IsMatch = int(isMatchParseInt)
	currency.IsMicro = int(isMicroParseInt)
	currency.MinNumber = minNumberParseFloat
	currency.MaxNumber = maxNumberParseFloat
	currency.MicroHoldtradeMax = int(microHoldtraeMaxParseInt)
	currency.Rate = rateParseFloat
	currency.Price = priceParseFloat
	currency.MicroMin = microMinParseFloat
	currency.MicroMax = microMaxParseFloat
	currency.RmbRelation = rmbRelationParseFloat
	currency.DecimalScale = int(decimalScaleParseInt)
	currency.Insurancable = int(insurancableParseInt)
	currency.ChainFee = chainFeeParseFloat
	currency.ContractAddress = contractAddress
	currency.MicroTradeFee = microTradeFeeParseFloat
	currency.Type = types
	currency.TotalAccount = ""
	currency.CollectAccount = ""
	currency.Key = ""
	currency.CreateTime = int(time.Now().Unix())
	currency.CurrencyDecimals = 0.0
	currency.IsDisplay = 1

	err = currency.UpdateCurrency()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("currency.AddCurrency() [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}

// DeleteCurrency
// 删除币种
func (m *AdminCurrency) DeleteCurrency(c *gin.Context) {
	lang := c.GetHeader("lang")
	id := c.PostForm("id")

	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	currency, err := new(model.Currency).QueryCurrencyById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyById [ERROR]  : %s", err))
		return
	}

	if currency.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "没有找到该币种",
		})
		return
	}

	idParseInt, _ := strconv.ParseInt(id, 10, 64)
	currency.Id = int(idParseInt)
	err = currency.DeleteCurrencyById()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "删除成功",
	})
	return

}
