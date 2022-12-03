package model

import (
	"fmt"
	"log"
	"okc/utils"
	"time"
)

// UsersWalletOut
// 用户钱包提币记录表
type UsersWalletOut struct {
	Id               int     `json:"id,omitempty"`               //
	UserId           int     `json:"user_id,omitempty"`          // 用户id
	Currency         int     `json:"currency,omitempty"`         // 币种id
	Address          string  `json:"address,omitempty"`          // 地址
	Nettype          string  `json:"nettype,omitempty"`          // 通道
	Number           float64 `json:"number,omitempty"`           //
	CreateTime       int     `json:"create_time,omitempty"`      // 创建时间
	Rate             float64 `json:"rate,omitempty"`             //
	Status           int     `json:"status,omitempty"`           //
	Notes            string  `json:"notes,omitempty"`            //
	RealNumber       float64 `json:"real_number,omitempty"`      //
	Txid             string  `json:"txid,omitempty"`             // 链上哈希
	Verificationcode string  `json:"verificationcode,omitempty"` //
	UpdateTime       int     `json:"update_time,omitempty"`      // 更新时间
	IsBank           int     `json:"is_bank,omitempty"`          // 1银行卡  0普通
	TibiRmb          float64 `json:"tibi_rmb,omitempty"`         // 提币金额 rmb
}

// AddUserWalletOut
// 添加用户提币记录
func (w *UsersWalletOut) AddUserWalletOut() error {

	if w.CreateTime == 0 {
		w.CreateTime = int(time.Now().Unix())
	}
	if w.UpdateTime == 0 {
		w.UpdateTime = int(time.Now().Unix())
	}

	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return err
	}
	sqls := "INSERT INTO users_wallet_out (user_id, currency, address, nettype, number, create_time, rate, status," +
		"notes, real_number, txid, verificationcode, update_time, is_bank, tibi_rmb) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.UserId, w.Currency, w.Address, w.Nettype, w.Number, w.CreateTime, w.Rate, w.Status, w.Notes,
		w.RealNumber, w.Txid, w.Verificationcode, w.UpdateTime, w.IsBank, w.TibiRmb)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// UpdateUserWalletOut
// 更新用户提币记录
func (w *UsersWalletOut) UpdateUserWalletOut() error {

	if w.UpdateTime == 0 {
		w.UpdateTime = int(time.Now().Unix())
	}

	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return err
	}
	sqls := "UPDATE users_wallet_out SET `user_id` = ?, `currency` = ?, `address` = ?, `nettype` = ?, `number` = ?, `create_time` = ?, `rate` = ?, `status` = ?," +
		"`notes` = ?, `real_number` = ?, `txid` = ?, `verificationcode` = ?, `update_time` = ?, `is_bank` = ?, `tibi_rmb` = ? WHERE id = ?"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.UserId, w.Currency, w.Address, w.Nettype, w.Number, w.CreateTime, w.Rate, w.Status, w.Notes,
		w.RealNumber, w.Txid, w.Verificationcode, w.UpdateTime, w.IsBank, w.TibiRmb, w.Id)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// QueryUserWalletOutById
// 根据id获取用户提币信息
func (w *UsersWalletOut) QueryUserWalletOutById(id interface{}) (*UsersWalletOut, error) {
	db := utils.Db
	sqls := "SELECT * FROM users_wallet_out WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	userWalletOut := new(UsersWalletOut)
	userWalletOutPtr, err := utils.RefStructGetFieldPtr(userWalletOut, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	row.Scan(userWalletOutPtr...)
	return userWalletOut, nil
}

// QueryUserWalletOutByIdAndLock
// 根据id获取用户提币信息(加锁)
func (w *UsersWalletOut) QueryUserWalletOutByIdAndLock(id interface{}) (*UsersWalletOut, error) {
	t, err := utils.Db.Begin()
	if err != nil {
		return nil, err
	}
	sqls := "SELECT * FROM users_wallet_out WHERE id = ? AND status <= 1 FOR UPDATE "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return nil, err
	}
	defer stmt.Close()
	userWalletOut := new(UsersWalletOut)
	userWalletOutPtr, err := utils.RefStructGetFieldPtr(userWalletOut, "*")
	if err != nil {
		t.Rollback()
		return nil, err
	}
	row := stmt.QueryRow(id)
	row.Scan(userWalletOutPtr...)
	t.Commit()
	return userWalletOut, nil
}

// QueryPagesByUserId
// 根据用户id获取分页记录
func (w *UsersWalletOut) QueryPagesByUserId(userId int, page, limit int) ([]map[string]string, int, error) {
	// 偏移量
	offset := (page - 1) * limit
	db := utils.Db
	sql := "SELECT count(*) FROM users_wallet_out AS w JOIN currency AS c ON w.currency = c.id WHERE w.user_id = ? ORDER BY w.id DESC LIMIT ?,?"
	var total int
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(userId, offset, limit)
	row.Scan(&total)

	sql = "SELECT w.*,c.name FROM users_wallet_out AS w JOIN currency AS c ON w.currency = c.id WHERE w.user_id = ? ORDER BY w.id DESC LIMIT ?,?"
	stmt, err = db.Prepare(sql)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// 获取列
	columns, _ := rows.Columns()
	// keys
	scanArgs := make([]interface{}, len(columns))
	// values
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// 获取数据
	data := make([]map[string]string, 0)
	for rows.Next() {
		rows.Scan(scanArgs...)
		maps := make(map[string]string)
		// 每行数据
		for i, col := range values {
			if col != nil {
				maps[columns[i]] = string(col)

			}
		}
		data = append(data, maps)
	}
	return data, total, nil
}

// QueryPagesByAccount
// 根据用户账号获取分页记录
func (w *UsersWalletOut) QueryPagesByAccount(account string, page, limit int) ([]map[string]string, int, error) {
	// 偏移量
	offset := (page - 1) * limit
	db := utils.Db
	sql := "SELECT count(*) FROM users_wallet_out AS w JOIN users AS u ON w.user_id = u.id "

	if account != "" {
		sql += fmt.Sprintf("WHERE u.phone = '%s' OR u.account_number = '%s' OR u.email = '%s' ", account, account, account)
	}

	sql += "ORDER BY w.id DESC LIMIT ?,?"
	log.Println(sql)
	var total int
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(offset, limit)
	row.Scan(&total)

	sql = "SELECT w.*, u.* FROM users_wallet_out AS w JOIN users AS u ON w.user_id = u.id "

	if account != "" {
		sql += fmt.Sprintf("WHERE u.phone = '%s' OR u.account_number = '%s' OR u.email = '%s' ", account, account, account)
	}
	sql += "ORDER BY w.id DESC LIMIT ?,?"
	stmt, err = db.Prepare(sql)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// 获取列
	columns, _ := rows.Columns()
	// keys
	scanArgs := make([]interface{}, len(columns))
	// values
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// 获取数据
	data := make([]map[string]string, 0)
	for rows.Next() {
		rows.Scan(scanArgs...)
		maps := make(map[string]string)
		// 每行数据
		for i, col := range values {
			if col != nil {
				maps[columns[i]] = string(col)

			}
		}
		data = append(data, maps)
	}
	return data, total, nil
}
