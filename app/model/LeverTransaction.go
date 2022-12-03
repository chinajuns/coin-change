package model

import "okc/utils"

// LeverTransaction
// 合约订单表
type LeverTransaction struct {
	Id                 int     `json:"id"`
	Type               int     `json:"type"`
	UserId             int     `json:"user_id"`
	Currency           int     `json:"currency"`
	Legal              int     `json:"legal"`
	OriginPrice        float64 `json:"origin_price"`
	Price              float64 `json:"price"`
	UpdatePrice        float64 `json:"update_price"`
	TargetProfitPrice  float64 `json:"target_profit_price"`
	StopLossPrice      float64 `json:"stop_loss_price"`
	Share              int     `json:"share"`
	Number             float64 `json:"number"`
	Multiple           int     `json:"multiple"`
	OriginCautionMoney float64 `json:"origin_caution_money"`
	CautionMoney       float64 `json:"caution_money"`
	FactProfits        float64 `json:"fact_profits"`
	TradeFee           float64 `json:"trade_fee"`
	Overnight          float64 `json:"overnight"`
	OvernightMoney     float64 `json:"overnight_money"`
	Status             int     `json:"status"`
	Settled            int     `json:"settled"`
	CreateTime         int     `json:"create_time"`
	TransactionTime    string  `json:"transaction_time"`
	UpdateTime         string  `json:"update_time"`
	HandleTime         string  `json:"handle_time"`
	CompleteTime       string  `json:"complete_time"`
	AgentPath          string  `json:"agent_path"`
}

// AddLeverTransaction
// 添加合约订单
func (l *LeverTransaction) AddLeverTransaction() (int, error) {
	t, err := utils.Db.Begin()
	if err != nil {
		return 0, err
	}
	sqls := "INSERT INTO lever_transaction (type, user_id, currency, legal, origin_price, price, update_price," +
		"target_profit_price, stop_loss_price, share, number, multiple, origin_caution_money, caution_money," +
		"fact_profits, trade_fee, overnight, overnight_money, status, settled, create_time, transaction_time," +
		"update_time, handle_time, complete_time, agent_path ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(l.Type, l.UserId, l.Currency, l.Legal, l.OriginPrice, l.Price, l.UpdatePrice,
		l.TargetProfitPrice, l.StopLossPrice, l.Share, l.Number, l.Multiple, l.OriginCautionMoney,
		l.CautionMoney, l.FactProfits, l.TradeFee, l.Overnight, l.OvernightMoney, l.Status,
		l.Settled, l.CreateTime, l.TransactionTime, l.UpdateTime, l.HandleTime, l.CompleteTime,
		l.AgentPath)
	if err != nil {
		t.Rollback()
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Rollback()
		return 0, err
	}
	t.Commit()
	return int(id), nil
}

// QueryLeverTransactionCountByUserIdAndStatus
// 根据用户id和订单状态获取订单总条数
func (l *LeverTransaction) QueryLeverTransactionCountByUserIdAndStatus(userId, status interface{}) (int, error) {
	db := utils.Db
	sqls := "SELECT count(*) FROM lever_transaction WHERE user_id = ? AND status = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var count int
	row := stmt.QueryRow(userId, status)
	row.Scan(&count)
	return count, nil
}

// QueryLeverTransactionById
// 根据合约订单id获取合约订单
func (l *LeverTransaction) QueryLeverTransactionById(id interface{}) (*LeverTransaction, error) {

	t, err := utils.Db.Begin()
	if err != nil {
		return nil, err
	}
	sqls := "SELECT * FROM lever_transaction WHERE id = ? FOR UPDATE "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return nil, err
	}
	defer stmt.Close()
	leverTransaction := new(LeverTransaction)
	leverTransactionPtr, err := utils.RefStructGetFieldPtr(leverTransaction, "*")
	if err != nil {
		t.Rollback()
		return nil, err
	}
	row := stmt.QueryRow(id)
	row.Scan(leverTransactionPtr...)
	t.Commit()
	return leverTransaction, nil
}

// UpdateLeverTransactionByStatusAndFactProfitsAndCreateTime
// 更新数据 [Status, FactProfits, CreateTime]
func (l *LeverTransaction) UpdateLeverTransactionByStatusAndFactProfitsAndCreateTime() error {
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sqls := "UPDATE lever_transaction SET status = ?, fact_profits = ?, create_time = ? WHERE id = ?"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.Status, l.FactProfits, l.CreateTime, l.Id)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// UpdateLeverTransactionByUpdateTimeAndUpdatePriceAndStatus
// 更新数据 [UpdateTime, UpdatePrice, Status]
func (l *LeverTransaction) UpdateLeverTransactionByUpdateTimeAndUpdatePriceAndStatus() error {
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sqls := "UPDATE lever_transaction SET update_price = ?, status = ?, update_time = ?, handle_time = ? WHERE id = ?"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.UpdatePrice, l.Status, l.UpdateTime, l.HandleTime, l.Id)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// Refresh
// 刷新数据
func (l *LeverTransaction) Refresh() error {

	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sqls := "SELECT * FROM lever_transaction WHERE id = ? FOR UPDATE "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	leverTransactionPtr, err := utils.RefStructGetFieldPtr(l, "*")
	if err != nil {
		t.Rollback()
		return err
	}
	row := stmt.QueryRow(l.Id)
	row.Scan(leverTransactionPtr...)
	t.Commit()
	return nil
}

// QueryLeverTransactionByUserIdAndCurrencyIdAndLegalId
// 根据用户id和币种id和法币id获取合约订单
func (l *LeverTransaction) QueryLeverTransactionByUserIdAndCurrencyIdAndLegalId(userId, currencyId, LegalId interface{}) ([]map[string]interface{}, error) {
	db := utils.Db
	sqls := "select * from `lever_transaction` where `user_id` = ? and `status` = 1 and `currency` = ? and `legal` = ? order by `id` desc limit 10"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, currencyId, LegalId)
	if err != nil {
		return nil, err
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
		rows.Scan(scanArgs...)
		maps := make(map[string]interface{})
		// 每行数据
		for i, col := range values {
			if col != nil {
				maps[columns[i]] = string(col)
			}
		}
		// 存入user
		sqls = "SELECT * FROM users WHERE id = ?"
		stmt, err = db.Prepare(sqls)
		if err != nil {
			return nil, err
		}
		user := new(Users)
		userPtr, err := utils.RefStructGetFieldPtr(user, "*")
		if err != nil {
			return nil, err
		}
		row := stmt.QueryRow(userId)
		row.Scan(userPtr...)
		maps["user"] = *user
		data = append(data, maps)
	}
	return data, nil
}

// QueryLeverTransactionAllMoneyByUserIdAndCurrencyIdAndLegalId
// 根据用户id和币种id和法币id获取所有杠杠余额
func (l *LeverTransaction) QueryLeverTransactionAllMoneyByUserIdAndCurrencyIdAndLegalId(userId, currencyId, legalId interface{}) (string, error) {
	db := utils.Db
	sqls := "SELECT sum(`number` * `price`) as `all_levers` FROM `lever_transaction` WHERE `legal` = ? AND `currency` = ? AND `user_id` = ? and `status` = 1"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var allLevers string
	row := stmt.QueryRow(legalId, currencyId, userId)
	row.Scan(&allLevers)
	return allLevers, nil
}

// QueryLeverTransactionByCurrencyIdAndLegalIdAndLimit
// 根据币种id和法币id和limit条数获取获取合约订单
func (l *LeverTransaction) QueryLeverTransactionByCurrencyIdAndLegalIdAndLimit(currencyId, legalId, limit interface{}) (map[string]interface{}, error) {
	db := utils.Db
	sqls := "SELECT * FROM lever_transaction WHERE currency = ? AND legal = ? AND type = 1 ORDER BY price DESC LIMIT ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(currencyId, legalId, limit)
	if err != nil {
		return nil, err
	}
	// 买入
	in := make([]LeverTransaction, 0)
	for rows.Next() {
		leverT := new(LeverTransaction)
		leverTPtr, err := utils.RefStructGetFieldPtr(leverT, "*")
		if err != nil {
			return nil, err
		}
		rows.Scan(leverTPtr...)
		in = append(in, *leverT)
	}

	sqls = "SELECT * FROM lever_transaction WHERE currency = ? AND legal = ? AND type = 2 ORDER BY price DESC LIMIT ?"
	stmt, err = db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err = stmt.Query(currencyId, legalId, limit)
	if err != nil {
		return nil, err
	}
	// 卖出
	out := make([]LeverTransaction, 0)
	for rows.Next() {
		leverT := new(LeverTransaction)
		leverTPtr, err := utils.RefStructGetFieldPtr(leverT, "*")
		if err != nil {
			return nil, err
		}
		rows.Scan(leverTPtr...)
		out = append(out, *leverT)
	}
	data := make(map[string]interface{})
	data["in"] = in
	data["out"] = out
	return data, nil
}

// QueryUserProfitByUserIdAndLegalId
// 根据用户id和法币id获取用户盈利和保证金
func (l *LeverTransaction) QueryUserProfitByUserIdAndLegalId(userId, LegalId interface{}) (string, string, string, error) {
	profitsTotal := "0.0"            // 交易对盈亏总额
	cautionMoneyTotal := "0.0"       // 交易对可用本金总额
	originCautionMoneyTotal := "0.0" // 交易对原始保证金

	db := utils.Db
	sqls := "SELECT SUM((CASE `type` WHEN 1 THEN `update_price` - `price` WHEN 2 THEN `price` - `update_price` END) * `number`) AS `profits_total`, SUM(`caution_money`) AS `caution_money_total`, SUM(`origin_caution_money`) AS `origin_caution_money_total` FROM lever_transaction WHERE user_id = ? AND legal = ? AND status = 1 GROUP BY user_id"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return "", "", "", err
	}
	defer stmt.Close()
	row := stmt.QueryRow(userId, LegalId)
	row.Scan(&profitsTotal, &cautionMoneyTotal, &originCautionMoneyTotal)

	return profitsTotal, cautionMoneyTotal, originCautionMoneyTotal, nil
}

// QueryUserProfitByUserIdAndLegalIdAndCurrencyId
// 根据用户id和法币id和币种id获取用户盈利和保证金
func (l *LeverTransaction) QueryUserProfitByUserIdAndLegalIdAndCurrencyId(userId, LegalId, currencyId interface{}) (string, string, string, error) {
	profitsTotal := "0.0"            // 交易对盈亏总额
	cautionMoneyTotal := "0.0"       // 交易对可用本金总额
	originCautionMoneyTotal := "0.0" // 交易对原始保证金

	db := utils.Db
	sqls := "SELECT SUM((CASE `type` WHEN 1 THEN `update_price` - `price` WHEN 2 THEN `price` - `update_price` END) * `number`) AS `profits_total`, SUM(`caution_money`) AS `caution_money_total`, SUM(`origin_caution_money`) AS `origin_caution_money_total` FROM lever_transaction WHERE user_id = ? AND legal = ? AND currency = ? AND status = 1 GROUP BY user_id"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return "", "", "", err
	}
	defer stmt.Close()
	row := stmt.QueryRow(userId, LegalId, currencyId)
	row.Scan(&profitsTotal, &cautionMoneyTotal, &originCautionMoneyTotal)

	return profitsTotal, cautionMoneyTotal, originCautionMoneyTotal, nil
}

// QueryLeverTransactionPageByCurrencyIdAndLegalId
// 根据币种id和法币id分页获取合约订单
func (l *LeverTransaction) QueryLeverTransactionPageByCurrencyIdAndLegalId(currencyId, legalId interface{}, page, limit int) ([]map[string]interface{}, error) {

	offset := (page - 1) * limit
	db := utils.Db
	sqls := "SELECT * FROM lever_transaction WHERE currency = ? AND legal = ? AND status = 1 ORDER BY id DESC LIMIT ?,?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(currencyId, legalId, offset, limit)
	if err != nil {
		return nil, err
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
		rows.Scan(scanArgs...)
		maps := make(map[string]interface{})
		// 每行数据
		for i, col := range values {
			if col != nil {
				maps[columns[i]] = string(col)
			}
		}

		data = append(data, maps)
	}
	return data, nil
}

// QueryLeverTransactionCountByUserIdAndInStatus
// 根据用户id和订单状态获取订单数量
func (l *LeverTransaction) QueryLeverTransactionCountByUserIdAndInStatus(userId interface{}, status string) (int, error) {
	db := utils.Db
	sqls := "SELECT count(*) FROM lever_transaction WHERE user_id = ? AND status IN (?)"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var count int
	row := stmt.QueryRow(userId, status)
	row.Scan(&count)
	return count, nil
}
