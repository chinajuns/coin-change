package model

import (
	"fmt"
	"log"
	"okc/utils"
	"time"
)

// ChargeReq
// 充值记录表
type ChargeReq struct {
	Id          int     `json:"id"`
	Uid         int     `json:"uid"`
	Amount      float64 `json:"amount"`
	UserAccount string  `json:"user_account"`
	Status      int     `json:"status"`
	CurrencyId  int     `json:"currency_id"`
	Remark      string  `json:"remark"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	IsBank      int     `json:"is_bank"`
	DaozhangNum float64 `json:"daozhang_num"`
	Img         string  `json:"img"`
}

// QueryChargeReqById
// 根据id获取充值记录
func (c ChargeReq) QueryChargeReqById(id interface{}) (*ChargeReq, error) {
	db := utils.Db
	sqls := "SELECT * FROM charge_req WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	chargeReq := new(ChargeReq)
	chargeReqPtr, err := utils.RefStructGetFieldPtr(chargeReq, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	row.Scan(chargeReqPtr...)
	return chargeReq, nil
}

// UpdateChargeReqById
// 根据id更新充值记录
func (c *ChargeReq) UpdateChargeReqById() error {
	db := utils.Db
	sqls := "UPDATE charge_req  SET `uid` = ?, `amount` = ?, `user_account` = ?, `status` = ?, `currency_id` = ?, `remark` = ?, `created_at` = ?, `updated_at` = ?, `is_bank` = ?,`daozhang_num` = ?, `img` = ? WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Uid, c.Amount, c.UserAccount, c.Status, c.CurrencyId, c.Remark, c.CreatedAt,
		c.UpdatedAt, c.IsBank, c.DaozhangNum, c.Img, c.Id)
	if err != nil {
		return err
	}
	return nil
}

// AddChargeReq
// 添加充值记录
func (c *ChargeReq) AddChargeReq() error {

	if c.CreatedAt == "" {
		c.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	}
	if c.UpdatedAt == "" {
		c.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	db := utils.Db
	sqls := "INSERT INTO charge_req (uid, amount, user_account, status, currency_id, remark, created_at, updated_at, is_bank," +
		"daozhang_num, img) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Uid, c.Amount, c.UserAccount, c.Status, c.CurrencyId, c.Remark, c.CreatedAt,
		c.UpdatedAt, c.IsBank, c.DaozhangNum, c.Img)
	if err != nil {
		return err
	}
	return nil
}

// QueryChargeReqListPageByTimeAndTypeAndAccount
// 根据endTime和startTime和account获取充值列表
func (c *ChargeReq) QueryChargeReqListPageByTimeAndTypeAndAccount(startTime, endTime, types, account interface{}, page, limit int) ([]map[string]interface{}, int, error) {
	offset := (page - 1) * limit
	db := utils.Db

	sqls := "SELECT count(*) FROM charge_req as r " +
		"JOIN users as u on u.id = r.uid " +
		"JOIN currency as c on c.id = r.currency_id "

	if types != "" {
		if types == "1" {
			sqls += "WHERE r.status = 1 "
		} else if types == "2" {
			sqls += "WHERE r.status > 1 "
		}
	}
	if startTime != "" {
		sqls += fmt.Sprintf("WHERE r.created_at >= %s ", startTime)
	}
	if endTime != "" {
		sqls += fmt.Sprintf("WHERE r.created_at <= %s ", endTime)
	}
	if account != "" {
		sqls += fmt.Sprintf("WHERE u.account_number = %s ", account)
	}
	sqls += "ORDER BY r.id DESC LIMIT ?,?"
	log.Println("sqls : ", sqls)
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	var total int
	row := stmt.QueryRow(offset, limit)
	row.Scan(&total)

	sqls = "SELECT r.*, u.account_number, c.name FROM charge_req as r " +
		"JOIN users as u on u.id = r.uid " +
		"JOIN currency as c on c.id = r.currency_id "

	if types != "" {
		if types == "1" {
			sqls += "WHERE r.status = 1 "
		} else if types == "2" {
			sqls += "WHERE r.status > 1 "
		}
	}
	if startTime != "" {
		sqls += fmt.Sprintf("WHERE r.created_at >= %s ", startTime)
	}
	if endTime != "" {
		sqls += fmt.Sprintf("WHERE r.created_at <= %s ", endTime)
	}
	if account != "" {
		sqls += fmt.Sprintf("WHERE u.account_number = %s ", account)
	}

	sqls += "ORDER BY r.id DESC LIMIT ?,?"
	stmt, err = db.Prepare(sqls)
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
