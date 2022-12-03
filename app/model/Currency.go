package model

import (
	"okc/utils"
)

// Currency
// 币种表
type Currency struct {
	Id                int     `json:"id"`                  // ID
	Name              string  `json:"name"`                // 币种名称
	GetAddress        string  `json:"get_address"`         //
	Sort              int     `json:"sort"`                // 排序
	Logo              string  `json:"logo"`                // 币种logo
	IsDisplay         int     `json:"is_display"`          // 是否显示 0 否 1是
	MinNumber         float64 `json:"min_number"`          // 最小提币数量
	MaxNumber         float64 `json:"max_number"`          // 最大提币数量
	Rate              float64 `json:"rate"`                // 费率
	IsLever           int     `json:"is_lever"`            // 是否杠杆币 0否 1是
	IsLegal           int     `json:"is_legal"`            // 是否法币 0否 1是
	IsMatch           int     `json:"is_match"`            // 是否撮合交易 0否 1是
	IsMicro           int     `json:"is_micro"`            // 是否微交易 0否1是
	Insurancable      int     `json:"insurancable"`        // 是否可买保险
	ShowLegal         int     `json:"show_legal"`          // 是否显示法币商家 0否 1是
	Type              string  `json:"type"`                // 基于哪个区块链
	BlackLimt         int     `json:"black_limt"`          // 币种黑名单限制数量
	Key               string  `json:"key"`                 //
	ContractAddress   string  `json:"contract_address"`    //
	TotalAccount      string  `json:"total_account"`       //
	CollectAccount    string  `json:"collect_account"`     // 归拢地址
	CurrencyDecimals  float64 `json:"currency_decimals"`   //
	RmbRelation       float64 `json:"rmb_relation"`        // 折合人民币比例
	DecimalScale      int     `json:"decimal_scale"`       // 发布小数点
	ChainFee          float64 `json:"chain_fee"`           //
	Price             float64 `json:"price"`               // 价值(美元)
	MicroTradeFee     float64 `json:"micro_trade_fee"`     // 微交易手续费%
	MicroMin          float64 `json:"micro_min"`           // 最小下单数量
	MicroMax          float64 `json:"micro_max"`           // 最大下单数量
	MicroHoldtradeMax int     `json:"micro_holdtrade_max"` // 最大持仓笔数
	CreateTime        int     `json:"create_time"`         //
}

// QueryCurrencyById
// 根据id获取币种信息
func (c *Currency) QueryCurrencyById(id interface{}) (*Currency, error) {
	db := utils.Db
	sqls := "SELECT * FROM currency WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	currency := new(Currency)
	currencyPtr, err := utils.RefStructGetFieldPtr(currency, "*")
	if err != nil {
		return nil, err
	}
	row.Scan(currencyPtr...)
	return currency, nil

}

// QueryCurrencyJoinMicroNumber
// 获取币种和期权数量
func (c *Currency) QueryCurrencyJoinMicroNumber() ([]map[string]interface{}, error) {
	db := utils.Db
	sqls := "SELECT c.*,n.currency_id, n.number FROM currency as c LEFT JOIN micro_numbers as n ON c.id = n.currency_id WHERE c.is_micro = 1"
	rows, err := db.Query(sqls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

	return data, err
}

// QueryCurrencyJoinQuotation
// 获取币种和行情
func (c *Currency) QueryCurrencyJoinQuotation() ([]map[string]string, error) {
	db := utils.Db
	sqls := "SELECT c.*,q.match_id, q.legal_id, q.currency_id, q.`change`, q.volume, q.now_price, q.add_time, q.xm FROM currency as c LEFT JOIN currency_quotation as q ON c.id = q.id WHERE c.is_display = 1"
	rows, err := db.Query(sqls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
		maps := make(map[string]string)
		rows.Scan(scanArgs...)
		// 每行数据
		for i, col := range values {
			if col != nil {
				maps[columns[i]] = string(col)
			}
		}
		data = append(data, maps)
	}

	return data, err
}

// QueryCurrencyPage
// 获取币种列表
func (c *Currency) QueryCurrencyPage(page, limit int) (data []map[string]interface{}, total int, err error) {
	offset := (page - 1) * limit
	db := utils.Db
	sqls := "SELECT count(*) FROM currency LIMIT ?,? "
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(offset, limit)
	row.Scan(&total)

	sqls = "SELECT * FROM currency LIMIT ?,? "
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
	data = make([]map[string]interface{}, 0)
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

// QueryCurrencyByName
// 根据name获取币种
func (c *Currency) QueryCurrencyByName(name interface{}) (*Currency, error) {
	db := utils.Db
	sqls := "SELECT * FROM currency WHERE name = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	currency := new(Currency)
	currencyPtr, err := utils.RefStructGetFieldPtr(currency, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(name)
	row.Scan(currencyPtr...)
	return currency, nil
}

// AddCurrency
// 添加币种
func (c *Currency) AddCurrency() error {
	db := utils.Db
	sqls := "INSERT INTO currency (name, get_address, sort, logo, is_display, min_number, max_number, rate," +
		"is_lever, is_legal, is_match, is_micro, insurancable, show_legal, type, black_limt, `key`, " +
		"contract_address, total_account, collect_account, currency_decimals, rmb_relation, decimal_scale," +
		"chain_fee, price, micro_trade_fee, micro_min, micro_max, micro_holdtrade_max, create_time) VALUES (" +
		" ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Name, c.GetAddress, c.Sort, c.Logo, c.IsDisplay, c.MinNumber, c.MaxNumber, c.Rate,
		c.IsLever, c.IsLegal, c.IsMatch, c.IsMicro, c.Insurancable, c.ShowLegal, c.Type, c.BlackLimt, c.Key,
		c.ContractAddress, c.TotalAccount, c.CollectAccount, c.CurrencyDecimals, c.RmbRelation, c.DecimalScale,
		c.ChainFee, c.Price, c.MicroTradeFee, c.MicroMin, c.MicroMax, c.MicroHoldtradeMax, c.CreateTime)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCurrency
// 更新币种
func (c *Currency) UpdateCurrency() error {
	db := utils.Db
	sqls := "UPDATE currency SET name = ?, get_address = ?, sort = ?, logo = ? , is_display = ?, min_number = ?, max_number = ?, rate = ?," +
		"is_lever = ?, is_legal = ?, is_match = ?, is_micro = ?, insurancable = ?, show_legal = ?, type = ?, black_limt = ?, `key` = ?, " +
		"contract_address = ?, total_account = ?, collect_account = ?, currency_decimals = ?, rmb_relation = ?, decimal_scale = ?," +
		"chain_fee = ?, price = ?, micro_trade_fee = ?, micro_min = ?, micro_max = ? , micro_holdtrade_max = ?  WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Name, c.GetAddress, c.Sort, c.Logo, c.IsDisplay, c.MinNumber, c.MaxNumber, c.Rate,
		c.IsLever, c.IsLegal, c.IsMatch, c.IsMicro, c.Insurancable, c.ShowLegal, c.Type, c.BlackLimt, c.Key,
		c.ContractAddress, c.TotalAccount, c.CollectAccount, c.CurrencyDecimals, c.RmbRelation, c.DecimalScale,
		c.ChainFee, c.Price, c.MicroTradeFee, c.MicroMin, c.MicroMax, c.MicroHoldtradeMax, c.Id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCurrencyById
// 根据id删除币种
func (c *Currency) DeleteCurrencyById() error {
	db := utils.Db
	sqls := "DELETE FROM currency WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id)
	if err != nil {
		return err
	}
	return nil
}
