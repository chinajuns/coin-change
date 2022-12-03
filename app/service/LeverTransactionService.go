package service

import (
	"encoding/json"
	"errors"
	"okc/app/model"
	"okc/utils"
	"strconv"
	"time"
)

// QueryWalletHazardRate
// 取钱包的风险率
func QueryWalletHazardRate(wallet *model.UsersWallet) (float64, error) {

	hazardRate := 0.0
	totalMoney := 0.0

	if wallet.Id == 0 {
		return hazardRate, errors.New("*model.CurrencyMatches not found")
	}

	// 取盈亏总额
	var profitsAll string            // 交易对盈亏总额
	var cautionMoneyAll string       // 交易对可用本金总额
	var originCautionMoneyAll string // 交易对原始保证金

	profitsAll, cautionMoneyAll, originCautionMoneyAll, err := new(model.LeverTransaction).QueryUserProfitByUserIdAndLegalId(wallet.UserId, wallet.Currency)
	if err != nil {
		return hazardRate, err
	}
	profitsAll = utils.TrimRightZeroByFloatStr(profitsAll)
	cautionMoneyAll = utils.TrimRightZeroByFloatStr(cautionMoneyAll)
	originCautionMoneyAll = utils.TrimRightZeroByFloatStr(originCautionMoneyAll)

	balance := wallet.LeverBalance
	balanceParseStr := strconv.FormatFloat(balance, 'f', 15, 64)
	totalMoney, err = utils.BcAdd(balanceParseStr, originCautionMoneyAll)
	if err != nil {
		return hazardRate, err
	}

	totalMoneyStr := strconv.FormatFloat(totalMoney, 'f', 15, 64)
	add, err := utils.BcAdd(totalMoneyStr, profitsAll)
	if err != nil {
		return hazardRate, err
	}
	addStr := strconv.FormatFloat(add, 'f', 15, 64)
	div, err := utils.BcDiv(addStr, originCautionMoneyAll)
	if err != nil {
		return hazardRate, err
	}
	divStr := strconv.FormatFloat(div, 'f', 4, 64)
	hazardRate, err = utils.BcMul(divStr, "100")
	if err != nil {
		return hazardRate, err
	}

	return hazardRate, nil
}

// CheckLeverClose
// 检查合约订单进行平仓
func CheckLeverClose(l *model.LeverTransaction) error {
	if l.Id == 0 {
		return errors.New("*model.LeverTransaction not nil")
	}

	// 事务处理
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}

	currency, err := new(model.Currency).QueryCurrencyById(l.Currency)
	if err != nil {
		t.Rollback()
		return err
	}
	// 获取行情最新价格
	lastPrice, err := QueryLastQuotationPriceByCurrencyName(currency.Name)
	lastPriceParseFloat, _ := strconv.ParseFloat(lastPrice, 64)
	if err != nil {
		t.Rollback()
		return err
	}
	// 数据刷新
	err = l.Refresh()
	if err != nil {
		t.Rollback()
		return err
	}
	if l.Status != 1 {
		t.Rollback()
		return errors.New("该笔交易状态异常,不能平仓")
	}
	unix := time.Now().Unix()
	unixStr := strconv.Itoa(int(unix))
	// 更新状态
	l.UpdatePrice = lastPriceParseFloat
	l.UpdateTime = unixStr
	l.Status = 2
	l.HandleTime = unixStr
	err = l.UpdateLeverTransactionByUpdateTimeAndUpdatePriceAndStatus()
	if err != nil {
		t.Rollback()
		return err
	}

	// 用户钱包
	userWatll, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(l.UserId, l.Currency)
	if err != nil {
		t.Rollback()
		return err
	}

	// 盈亏计算
	profitParStr := strconv.FormatFloat(l.FactProfits, 'f', 4, 64)
	CautionMoneyParStr := strconv.FormatFloat(l.CautionMoney, 'f', 4, 64)
	chage, err := utils.BcAdd(CautionMoneyParStr, profitParStr)
	diff := 0
	if err != nil {
		t.Rollback()
		return err
	}
	// 序列化数据
	extraData := map[string]interface{}{
		"trade_id":      l.Id,
		"caution_money": l.CautionMoney,
		"profit":        profitParStr,
		"diff":          diff,
	}
	extraDataJson, _ := json.Marshal(extraData)
	// 改变钱包余额
	//log.Println("extarDataJson : ", string(extraDataJson))

	res := ChangeUserWalletBalance(userWatll, 3, chage, 31, "平仓资金处理", false, 0, diff, string(extraDataJson), 1, 1)
	switch res.(type) {
	case error:
		t.Rollback()
		return res.(error)
		break
	case string:
		t.Rollback()
		return errors.New(res.(string))
	case bool:
		if res.(bool) == false {
			t.Rollback()
			return errors.New("平仓失败:更新处理状态失败")
		}
		break
	}
	// 再次更新订单状态
	err = l.Refresh()
	if err != nil {
		t.Rollback()
		return err
	}
	profitParFloat, _ := strconv.ParseFloat(profitParStr, 64)
	l.Status = 3
	l.FactProfits = profitParFloat
	l.CreateTime = int(time.Now().Unix())
	err = l.UpdateLeverTransactionByStatusAndFactProfitsAndCreateTime()
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()

	return nil
}
