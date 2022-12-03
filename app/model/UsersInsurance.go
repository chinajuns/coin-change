package model

import "okc/utils"

// UsersInsurance
// 用户保险
type UsersInsurance struct {
	Id              int     `json:"id"`
	UserId          int     `json:"user_id"`
	InsuranceTypeId int     `json:"insurance_type_id"`
	Amount          float64 `json:"amount"`
	InsuranceAmount float64 `json:"insurance_amount"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	YieldedAt       string  `json:"yielded_at"`
	RescindedAt     string  `json:"rescinded_at"`
	RescindedType   int     `json:"rescinded_type"`
	Status          int     `json:"status"`
	ClaimStatus     int     `json:"claim_status"`
}

// QueryUsersInsuranceByUserId
// 根据用户id获取用户保险
func (i *UsersInsurance) QueryUsersInsuranceByUserId(userId interface{}) (*UsersInsurance, error) {
	db := utils.Db
	sqls := "SELECT * FROM users_insurances WHERE user_id = ? AND status = 1"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	usersInsurance := new(UsersInsurance)
	usersInsurancePtr, err := utils.RefStructGetFieldPtr(usersInsurance, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(userId)
	row.Scan(usersInsurancePtr...)
	return usersInsurance, nil
}

// QueryUsersInsuranceByUserIdAndCurrencyId
// 根据用户id和币种id获取用户保险
func (i *UsersInsurance) QueryUsersInsuranceByUserIdAndCurrencyId(userId, currencyId interface{}) ([]map[string]interface{}, error) {
	db := utils.Db
	sqls := "SELECT ui.*, t.* FROM users_insurances as ui LEFT JOIN insurance_types as t on ui.insurance_type_id = t.id WHERE ui.user_id = ?  AND ui.status = 1 AND t.currency_id =?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, currencyId)
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
