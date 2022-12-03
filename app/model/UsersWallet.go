package model

import (
	"log"
	"okc/utils"
)

// UsersWallet
// 用户钱包表
type UsersWallet struct {
	Id                    int     `json:"id,omitempty"`
	UserId                int     `json:"user_id,omitempty"`
	Currency              int     `json:"currency,omitempty"`
	Address               string  `json:"address,omitempty"`
	Address2              string  `json:"address_2,omitempty""`
	LegalBalance          float64 `json:"legal_balance,omitempty"`
	LockLegalBalance      float64 `json:"lock_legal_balance,omitempty"`
	ChangeBalance         float64 `json:"change_balance,omitempty"`
	LockChangeBalance     float64 `json:"lock_change_balance,omitempty"`
	LeverBalance          float64 `json:"lever_balance,omitempty"`
	LeverBalanceAddAllnum float64 `json:"lever_balance_add_allnum,omitempty"`
	LockLeverBalance      float64 `json:"lock_lever_balance,omitempty"`
	MicroBalance          float64 `json:"micro_balance,omitempty"`
	LockMicroBalance      float64 `json:"lock_micro_balance,omitempty"`
	InsuranceBalance      float64 `json:"insurance_balance,omitempty"`
	LockInsuranceBalance  float64 `json:"lock_insurance_balance,omitempty"`
	Status                int     `json:"status,omitempty"`
	CreateTime            int     `json:"create_time,omitempty"`
	OldBalance            string  `json:"old_balance,omitempty"`
	Private               string  `json:"private,omitempty"`
	Cost                  string  `json:"cost,omitempty"`
	GlTime                string  `json:"gl_time,omitempty"`
	Txid                  string  `json:"txid,omitempty"`
}

// QueryUserWalletByUserIdAndCurrencyId
// 根据用户id和币种id获取钱包
func (w *UsersWallet) QueryUserWalletByUserIdAndCurrencyId(userId, currencyId interface{}) (*UsersWallet, error) {
	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return nil, err
	}
	sqls := "SELECT * FROM users_wallet WHERE user_id = ? AND currency = ? FOR UPDATE "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return nil, err
	}
	defer stmt.Close()
	usersWallet := new(UsersWallet)
	row := stmt.QueryRow(userId, currencyId)
	row.Scan(&usersWallet.Id, &usersWallet.UserId, &usersWallet.Currency, &usersWallet.Address, &usersWallet.Address2, &usersWallet.LegalBalance,
		&usersWallet.LockLegalBalance, &usersWallet.ChangeBalance, &usersWallet.LockChangeBalance, &usersWallet.LeverBalance,
		&usersWallet.LeverBalanceAddAllnum, &usersWallet.LockLeverBalance, &usersWallet.MicroBalance, &usersWallet.LockMicroBalance,
		&usersWallet.InsuranceBalance, &usersWallet.LockInsuranceBalance, &usersWallet.Status, &usersWallet.CreateTime,
		&usersWallet.OldBalance, &usersWallet.Private, &usersWallet.Cost, &usersWallet.GlTime, &usersWallet.Txid)
	t.Commit()
	return usersWallet, err
}

// QueryMapUserWalletByUserIdAndCurrencyId
// 根据用户id和币种id获取钱包集合
func (w *UsersWallet) QueryMapUserWalletByUserIdAndCurrencyId(userId, currencyId interface{}) (map[string]interface{}, error) {
	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return nil, err
	}
	sqls := "SELECT * FROM users_wallet WHERE user_id = ? AND currency = ? FOR UPDATE "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return nil, err
	}
	defer stmt.Close()
	usersWallet := new(UsersWallet)
	row := stmt.QueryRow(userId, currencyId)
	row.Scan(&usersWallet.Id, &usersWallet.UserId, &usersWallet.Currency, &usersWallet.Address, &usersWallet.Address2, &usersWallet.LegalBalance,
		&usersWallet.LockLegalBalance, &usersWallet.ChangeBalance, &usersWallet.LockChangeBalance, &usersWallet.LeverBalance,
		&usersWallet.LeverBalanceAddAllnum, &usersWallet.LockLeverBalance, &usersWallet.MicroBalance, &usersWallet.LockMicroBalance,
		&usersWallet.InsuranceBalance, &usersWallet.LockInsuranceBalance, &usersWallet.Status, &usersWallet.CreateTime,
		&usersWallet.OldBalance, &usersWallet.Private, &usersWallet.Cost, &usersWallet.GlTime, &usersWallet.Txid)
	log.Println("usersWallet: ", *usersWallet)
	t.Commit()
	usersWalletMap, err := utils.RefStructChangeMap(usersWallet)
	return usersWalletMap, err
}

// QueryUserWalletListPageByUserId
// 根据用户id获取钱包列表
func (w *UsersWallet) QueryUserWalletListPageByUserId(userId interface{}, page, limit int) ([]map[string]interface{}, int, error) {
	offset := (page - 1) * limit
	db := utils.Db
	sqls := "SELECT count(*) FROM users_wallet WHERE user_id = ? LIMIT ?,?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	var total int
	row := stmt.QueryRow(userId, offset, limit)
	row.Scan(&total)

	sqls = "SELECT * FROM users_wallet WHERE user_id = ? LIMIT ?,?"
	stmt, err = db.Prepare(sqls)
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
	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		maps := make(map[string]interface{})
		rows.Scan(scanArgs...)
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

// UpdateUserWalletById
// 根据钱包id更新用户钱包
func (w *UsersWallet) UpdateUserWalletById() error {
	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return err
	}
	sqls := "UPDATE users_wallet SET address = ?, address_2 = ? , legal_balance = ? , lock_legal_balance = ? , change_balance = ? ," +
		"lock_change_balance = ? , lever_balance = ? , lever_balance_add_allnum = ? , lock_lever_balance = ? ," +
		"micro_balance = ? , lock_micro_balance = ? , insurance_balance = ? , lock_insurance_balance = ? ," +
		"status = ? , create_time = ? , old_balance = ? , private = ? , cost = ? , gl_time = ? , txid = ?" +
		" WHERE id = ?"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.Address, w.Address2, w.LegalBalance, w.LockLegalBalance, w.ChangeBalance,
		w.LockChangeBalance, w.LeverBalance, w.LeverBalanceAddAllnum, w.LockLeverBalance,
		w.MicroBalance, w.LockMicroBalance, w.InsuranceBalance, w.LockInsuranceBalance,
		w.Status, w.CreateTime, w.OldBalance, w.Private, w.Cost, w.GlTime, w.Txid, w.Id)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil

}

// Refresh
// 刷新数据
func (w *UsersWallet) Refresh() error {
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sqls := "SELECT * FROM users_wallet WHERE id = ? FOR UPDATE "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	userWalletPtr, err := utils.RefStructGetFieldPtr(w, "*")
	if err != nil {
		t.Rollback()
		return err
	}
	row := stmt.QueryRow(w.Id)
	row.Scan(userWalletPtr...)
	t.Commit()
	return nil
}
