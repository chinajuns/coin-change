package service

import (
	"errors"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// UserBuyCoint
// 买入
func UserBuyCoint(userId, currencyId, legalId int, amount, targetPrice string) error {
	// 用户钱包
	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userId, legalId)
	if err != nil {
		return err
	}
	if userWallet.Id == 0 {
		return errors.New("钱包不存在")
	}
	// 币种
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		return err
	}
	if currency.Id == 0 {
		return errors.New("币种不存在")
	}
	// 币种最新价格
	lastPrice, err := QueryLastQuotationPriceByCurrencyName(currency.Name)
	if err != nil {
		return err
	}
	constPrice, err := utils.BcMul(amount, targetPrice)
	if err != nil {
		return err
	}
	// 冻结资金
	res := ChangeUserWalletBalance(userWallet, 2, constPrice*-1, 50, "币币交易下单，资金冻结", false, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("币币交易下单，资金冻结失败")
		}
		break
	}
	res = ChangeUserWalletBalance(userWallet, 2, constPrice, 50, "币币交易下单，冻结资金增加", true, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("币币交易下单，冻结资金增加失败")
		}
		break
	}
	lastPriceParseFloat, err := strconv.ParseFloat(lastPrice, 64)
	if err != nil {
		return err
	}
	amountParseFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}
	targetPriceParseFloat, err := strconv.ParseFloat(targetPrice, 64)
	if err != nil {
		return err
	}
	fee, err := new(model.Setting).QueryValueByKey("COIN_TRADE_FEE", "0.0")
	if err != nil {
		return err
	}
	feeParseFloat, err := strconv.ParseFloat(fee, 64)
	if err != nil {
		return err
	}
	// 生成订单
	coin := &model.CoinTrade{
		UId:         userId,
		CurrencyId:  currencyId,
		LegalId:     legalId,
		Type:        1,
		TargetPrice: targetPriceParseFloat,
		TradePrice:  lastPriceParseFloat,
		TradeAmount: amountParseFloat,
		ChargeFee:   feeParseFloat,
		Status:      1,
	}
	err = coin.AddCoinTrade()
	if err != nil {
		return err
	}
	return nil
}

// UserSellCoin
// 卖出
func UserSellCoin(userId, currencyId, legalId int, amount, targetPrice string) error {
	// 用户钱包
	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userId, currencyId)
	if err != nil {
		return err
	}
	if userWallet.Id == 0 {
		return errors.New("钱包不存在")
	}
	// 币种
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		return err
	}
	if currency.Id == 0 {
		return errors.New("币种不存在")
	}
	// 币种最新价格
	lastPrice, err := QueryLastQuotationPriceByCurrencyName(currency.Name)
	if err != nil {
		return err
	}

	lastPriceParseFloat, err := strconv.ParseFloat(lastPrice, 64)
	if err != nil {
		return err
	}
	amountParseFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}
	targetPriceParseFloat, err := strconv.ParseFloat(targetPrice, 64)
	if err != nil {
		return err
	}
	fee, err := new(model.Setting).QueryValueByKey("COIN_TRADE_FEE", "0.0")
	if err != nil {
		return err
	}
	feeParseFloat, err := strconv.ParseFloat(fee, 64)
	if err != nil {
		return err
	}

	// 冻结资金
	res := ChangeUserWalletBalance(userWallet, 2, amountParseFloat*-1, 50, "币币交易下单，资金冻结", false, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("币币交易下单，资金冻结失败")
		}
		break
	}
	res = ChangeUserWalletBalance(userWallet, 2, amountParseFloat, 50, "币币交易下单，冻结资金增加", true, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("币币交易下单，冻结资金增加失败")
		}
		break
	}

	// 生成订单
	coin := &model.CoinTrade{
		UId:         userId,
		CurrencyId:  currencyId,
		LegalId:     legalId,
		Type:        2,
		TargetPrice: targetPriceParseFloat,
		TradePrice:  lastPriceParseFloat,
		TradeAmount: amountParseFloat,
		ChargeFee:   feeParseFloat,
		Status:      1,
	}
	err = coin.AddCoinTrade()
	if err != nil {
		return err
	}
	return nil
}

// UnUserBuyCoint
// 买入类型解除冻结还钱
func UnUserBuyCoint(userId, currencyId, legalId int, targetPrice, targetAmount string) error {
	// 用户钱包
	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userId, legalId)
	if err != nil {
		return err
	}
	if userWallet.Id == 0 {
		return errors.New("钱包不存在")
	}
	// 币种
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		return err
	}
	if currency.Id == 0 {
		return errors.New("币种不存在")
	}

	constPrice, err := utils.BcMul(targetPrice, targetAmount)
	if err != nil {
		return err
	}
	// 解除冻结资金
	res := ChangeUserWalletBalance(userWallet, 2, constPrice, 50, "取消币币交易，资金返还", false, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("取消币币交易，资金返还失败")
		}
		break
	}
	res = ChangeUserWalletBalance(userWallet, 2, constPrice*-1, 50, "取消币币交易，退还冻结资金", true, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("取消币币交易，退还冻结资金失败")
		}
		break
	}

	return nil
}

// UnUserSellCoin
// 买入类型解除冻结还钱
func UnUserSellCoin(userId, currencyId int, tradeAmount string) error {
	// 用户钱包
	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userId, currencyId)
	if err != nil {
		return err
	}
	if userWallet.Id == 0 {
		return errors.New("钱包不存在")
	}
	// 币种
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		return err
	}
	if currency.Id == 0 {
		return errors.New("币种不存在")
	}

	tradeAmountParseFloat, err := strconv.ParseFloat(tradeAmount, 64)
	if err != nil {
		return err
	}

	// 取消币币交易
	res := ChangeUserWalletBalance(userWallet, 2, tradeAmountParseFloat, 50, "取消币币交易", false, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("取消币币交易失败")
		}
		break
	}
	res = ChangeUserWalletBalance(userWallet, 2, tradeAmountParseFloat*-1, 50, "取消币币交易,退还冻结资金", true, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		return res.(error)
		break
	case string:
		return errors.New(res.(string))
		break
	case bool:
		if res.(bool) == false {
			return errors.New("取消币币交易,退还冻结资金")
		}
		break
	}

	return nil
}
