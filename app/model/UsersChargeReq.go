package model

import "okc/utils"

// UsersChargeReq
// 用户充币记录表
type UsersChargeReq struct {
	Id          int     `json:"id,omitempty"`           //
	Uid         int     `json:"uid,omitempty"`          //
	Amount      float64 `json:"amount,omitempty"`       //
	UserAccount string  `json:"user_account,omitempty"` //
	Status      int     `json:"status,omitempty"`       //
	Remark      string  `json:"remark,omitempty"`       //
	CreateAt    string  `json:"create_at,omitempty"`    //
	UpdateAt    string  `json:"update_at,omitempty"`    //
	IsBank      int     `json:"is_bank,omitempty"`      //
	DaozhangNum float64 `json:"daozhang_num,omitempty"` //
	Img         string  `json:"img,omitempty"`          //
}

// QueryPagesByUserId
// 根据用户id获取分页记录
func (w *UsersChargeReq) QueryPagesByUserId(userId int, page, limit int) ([]map[string]string, int, error) {
	// 偏移量
	offset := (page - 1) * limit
	db := utils.Db
	sql := "SELECT count(*) FROM charge_req AS r JOIN currency AS c ON r.currency_id = c.id WHERE r.uid = ? ORDER BY r.id DESC LIMIT ?,?"
	var total int
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(userId, offset, limit)
	row.Scan(&total)

	sql = "SELECT r.*,c.name FROM charge_req AS r JOIN currency AS c ON r.currency_id = c.id WHERE r.uid = ? ORDER BY r.id DESC LIMIT ?,?"
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
