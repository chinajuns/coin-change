package model

import "okc/utils"

// LeverMultiple
// 合约倍数表
type LeverMultiple struct {
	Id         int    `json:"id"`
	Type       int    `json:"type"`
	Value      string `json:"value"`
	CurrencyId int    `json:"currency_id"`
}

// QueryLeverMultipleByCurrencyId
// 根据币种id获取合约倍数
func (l *LeverMultiple) QueryLeverMultipleByCurrencyId(currencyId interface{}) ([]map[string]interface{}, error) {
	db := utils.Db
	sqls := "SELECT * FROM lever_multiple WHERE currency_id = ? AND type = 1"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(currencyId)
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

// QueryLeverMultipleByValue
// 根据value获取合约倍数
func (l *LeverMultiple) QueryLeverMultipleByValue(value interface{}) (*LeverMultiple, error) {
	db := utils.Db
	sqls := "SELECT * FROM lever_multiple WHERE value = ? AND type = 1"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	leverMultiple := new(LeverMultiple)
	leverMultiplePtr, err := utils.RefStructGetFieldPtr(leverMultiple, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(value)
	row.Scan(leverMultiplePtr...)
	return leverMultiple, nil
}
