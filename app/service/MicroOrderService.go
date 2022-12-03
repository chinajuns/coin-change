package service

import (
	"okc/app/model"
	"time"
)

// CheckMicroInsurance
// 检查期权保险
func CheckMicroInsurance(userId, currencyId interface{}, number int) interface{} {

	// 查询用户是否买了保险
	userInsuranceSlice, err := new(model.UsersInsurance).QueryUsersInsuranceByUserIdAndCurrencyId(userId, currencyId)
	if err != nil {
		return err
	}
	userInsurance := userInsuranceSlice[0]
	insuranceId := userInsurance["id"]
	if insuranceId.(int) == 0 {
		return "尚未申购或理赔保险"
	}
	// 检查类型
	isTAdd1 := userInsurance["is_t_add_1"]
	if isTAdd1.(int) == 1 {
		userInsuranceCreatedAtTime, _ := time.Parse("2006-01-02 15:04:05", userInsurance["created_at"].(string))
		if time.Now().Unix() < userInsuranceCreatedAtTime.Unix() {
			return "申购的保险T+1生效"
		}
	}

	// 用户钱包
	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userId, currencyId)
	if err != nil {
		return err
	}
	// 受保资产为0不允许下单
	if userWallet.InsuranceBalance == 0 {
		return "受保资产为零"
	}

	switch userInsurance["type"].(int) {
	case 1:
		// 正向险种，受保资产小于等于【条件1额度】，不允许下单
		if userWallet.InsuranceBalance <= userInsurance["defective_claims_condition2"].(float64) {
			return "受保资产小于等于可下单条件"
		}
		break
	case 2:
		// 反向险种，受保资产小于等于【条件2额度】，不允许下单
		if userWallet.InsuranceBalance <= userInsurance["defective_claims_condition2"].(float64) {
			return "您已超过持仓限制，暂停下单"
		}
		break
	default:
		return "未知的险种类型"
		break
	}

	// 判断持仓数量
	if number >= 500 {
		return "超过最大持仓数量限制"
	}

	// 判断挂单数量
	microOrderCount, err := new(model.MicroOrder).QueryMicroOrderCountByUserIdAndCurrencyId(userId, currencyId)
	if err != nil {
		return err
	}
	if microOrderCount >= 5 {
		return "交易中的订单大于最大挂单数量"
	}
	return true
}
