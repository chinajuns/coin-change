package model

import (
	"okc/utils"
	"time"
)

// WalletLog
// 钱包日志
type WalletLog struct {
	Id           int         `json:"id"`
	UserId       int         `json:"user_id"`
	FromUserId   int         `json:"from_user_id"`
	AccountLogId int         `json:"account_log_id"`
	WalletId     int         `json:"wallet_id"`
	BalanceType  int         `json:"balance_type"`
	LockType     int         `json:"lock_type"`
	Before       string      `json:"before"`
	Change       string      `json:"change"`
	After        string      `json:"after"`
	Memo         string      `json:"memo"`
	ExtraSign    int         `json:"extra_sign"`
	ExtraData    interface{} `json:"extra_data"`
	CreateTime   int         `json:"create_time"`
}

// AddWalletLog
// 添加钱包日志
func (w *WalletLog) AddWalletLog() error {
	if w.CreateTime == 0 {
		w.CreateTime = int(time.Now().Unix())
	}
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sqls := "INSERT INTO wallet_log (`user_id`, `from_user_id`, `account_log_id`, `wallet_id`, `balance_type`, `lock_type`, `before`, `change`, `after`, `memo`, `extra_sign`, `extra_data`, `create_time`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	_, err = stmt.Exec(w.UserId, w.FromUserId, w.AccountLogId, w.WalletId, w.BalanceType,
		w.LockType, w.Before, w.Change, w.After, w.Memo, w.ExtraSign, w.ExtraData, w.CreateTime)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}
