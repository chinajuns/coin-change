package model

import "okc/utils"

// MicroSecond
// 期货到期时间表
type MicroSecond struct {
	Id          int     `json:"id"`
	Seconds     int     `json:"seconds"`
	MinNum      int     `json:"min_num"`
	Status      int     `json:"status"`
	ProfitRatio float64 `json:"profit_ratio"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// QueryMicroSecondList
// 获取到期时间列表
func (m *MicroSecond) QueryMicroSecondList() ([]map[string]interface{}, error) {
	db := utils.Db
	sqls := "SELECT * FROM micro_seconds WHERE status = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(1)
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

// QueryMicroSecondBySecond
// 根据时间获取到期时间
func (m *MicroSecond) QueryMicroSecondBySecond(second interface{}) (*MicroSecond, error) {
	db := utils.Db
	sqls := "SELECT * FROM micro_seconds WHERE seconds = ? AND status = 1"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	microSecond := new(MicroSecond)
	microSecondPtr, err := utils.RefStructGetFieldPtr(microSecond, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(second)
	row.Scan(microSecondPtr...)
	return microSecond, nil
}
