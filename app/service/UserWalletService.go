package service

import (
	"fmt"
	"okc/app/model"
	"okc/utils"
	"reflect"
	"strconv"
	"strings"
)

// ChangeUserWalletBalance
// @*model.UsersWallet 用户钱包模型
// @balanceType        1.法币,2.币币交易,3.杠杆交易,4.秒合约,5.保险
// @change             添加传正数，减少传负数
// @accountLogType     日志状态
// @remark             备注
// @isClock            是否是冻结或解冻资金
// @formUserId         用户id
// @extraSign          子场景标识
// @extraData          附加数据
// @zeroContinue       改变为0时继续执行,默认为假不执行
// @overflow           余额不足时允许继续处理,默认为假不允许
// 改变钱包余额
func ChangeUserWalletBalance(wallet *model.UsersWallet, balanceType int,
	change float64, accountLogType int, remark string, isClock bool, formUserId, extraSign int,
	extraData string, zeroContinue, overflow int) interface{} {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteErrorLog(fmt.Sprintf("ChangeUserWalletBalance [ERROR] : %s", err))
		}
	}()
	if zeroContinue == 0 && change == 0.0 {
		_ = utils.WriteInfoLog(fmt.Sprintf("改变金额为0,不处理"))
		return true
	}
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	if !utils.InCheckIntOrIntSlice(balanceType, []int{1, 2, 3, 4, 5}) {
		t.Rollback()
		return "货币类型不正确"
	}

	// 锁字段
	lockFields := []string{
		"",
		"legal_balance",
		"change_balance",
		"lever_balance",
		"micro_balance",
		"insurance_balance",
	}

	var lockField string
	var before float64
	var lockType int
	// 如果是冻结或解冻
	if isClock {
		lockField = fmt.Sprintf("lock_%s", lockFields[balanceType])
		lockType = 1
	} else {
		lockField = lockFields[balanceType]
		lockType = 0
	}

	// 刷新钱包
	err = wallet.Refresh()
	if err != nil {
		t.Rollback()
		return err
	}
	walletTypeOf := reflect.TypeOf(wallet)
	walletValueOf := reflect.ValueOf(wallet)
	walletFieldNum := walletTypeOf.Elem().NumField()
	//log.Println("lockField : ", lockField)
	for i := 0; i < walletFieldNum; i++ {

		fieldName := walletTypeOf.Elem().Field(i).Name
		fieldTag := walletTypeOf.Elem().Field(i).Tag.Get("json")

		if strings.Index(fieldTag, ",") != -1 {
			fieldTagSlice := strings.Split(fieldTag, ",")
			fieldTag = fieldTagSlice[0]
		}

		if lockField == fieldName {
			before = walletValueOf.Elem().FieldByName(fieldName).Float()
		} else if lockField == fieldTag {
			before = walletValueOf.Elem().FieldByName(fieldName).Float()
		}
	}
	beforeStr := strconv.FormatFloat(before, 'f', 4, 64)
	changeStr := strconv.FormatFloat(change, 'f', 4, 64)

	after, err := utils.BcAdd(beforeStr, changeStr)

	//log.Println("beforeStr : ", beforeStr)
	//log.Println("changeStr : ", changeStr)
	//log.Println("after : ", after)
	//log.Println("overflow : ", overflow)

	if err != nil {
		t.Rollback()
		return err
	}
	afterStr := strconv.FormatFloat(after, 'f', 4, 64)

	// 判断余额是否充足
	if after < 0 && overflow == 0 {
		t.Rollback()
		_ = utils.WriteErrorLog(fmt.Sprintf("币种钱包余额名称: %s, after=%f, overflow=%f,钱包余额不足", lockField, after, overflow))
		return "钱包余额不足"
	}

	// 创建日志
	accountLogModel := &model.AccountLog{
		UserId:   wallet.UserId,
		Value:    changeStr,
		Info:     remark,
		Type:     accountLogType,
		Currency: wallet.Currency,
	}
	accountLogId, err := accountLogModel.AddAccountLog()
	if err != nil {
		t.Rollback()
		return err
	}

	walletLogModel := &model.WalletLog{
		UserId:       wallet.UserId,
		FromUserId:   formUserId,
		AccountLogId: accountLogId,
		WalletId:     wallet.Id,
		BalanceType:  balanceType,
		LockType:     lockType,
		Before:       beforeStr,
		Change:       changeStr,
		After:        afterStr,
		Memo:         remark,
		ExtraSign:    extraSign,
		ExtraData:    extraData,
	}
	err = walletLogModel.AddWalletLog()
	if err != nil {
		t.Rollback()
		return err
	}
	//log.Println("wallet1 : ", wallet)
	for i := 0; i < walletFieldNum; i++ {

		fieldName := walletTypeOf.Elem().Field(i).Name
		fieldTag := walletTypeOf.Elem().Field(i).Tag.Get("json")

		if strings.Index(fieldTag, ",") != -1 {
			fieldTagSlice := strings.Split(fieldTag, ",")
			fieldTag = fieldTagSlice[0]
		}
		afterValueOf := reflect.ValueOf(after)
		if lockField == fieldName {

			walletValueOf.Elem().FieldByName(fieldName).Set(afterValueOf)
		} else if lockField == fieldTag {
			walletValueOf.Elem().FieldByName(fieldName).Set(afterValueOf)
		}
	}
	//walletJson, _ := json.Marshal(wallet)
	//log.Println("wallet : ", string(walletJson))
	err = wallet.UpdateUserWalletById()
	if err != nil {
		return err
	}

	return true
}
