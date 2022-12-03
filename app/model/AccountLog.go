package model

import (
	"fmt"
	"log"
	"okc/utils"
	"strings"
	"time"
)

// AccountLog
// 账号日志
type AccountLog struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Value       string `json:"value"`
	CreatedTime int    `json:"created_time"`
	Info        string `json:"info"`
	Type        int    `json:"type"`
	Currency    int    `json:"currency"`
}

// AddAccountLog
// 添加账号日志
func (a *AccountLog) AddAccountLog() (int, error) {
	if a.CreatedTime == 0 {
		a.CreatedTime = int(time.Now().Unix())
	}
	db := utils.Db
	//t, err := utils.Db.Begin()
	//if err != nil {
	//	return 0, err
	//}
	sqls := "INSERT INTO account_log (user_id, value, created_time, info, type, currency) VALUES (?,?,?,?,?,?)"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		//t.Rollback()
		return 0, err
	}
	res, err := stmt.Exec(a.UserId, a.Value, a.CreatedTime, a.Info, a.Type, a.Currency)
	if err != nil {
		//t.Rollback()
		return 0, err
	}
	id, _ := res.LastInsertId()
	//t.Commit()
	return int(id), nil
}

// QueryAccountLogListPage
// 分页获取账号日志列表
func (a *AccountLog) QueryAccountLogListPage(account, startTime, endTime, currency, types, sign interface{}, page, limit int) ([]map[string]interface{}, int, error) {
	offset := (page - 1) * limit
	db := utils.Db
	sqls := "SELECT count(*) FROM account_log as a "

	if account != "" {
		sqls = "SELECT count(*) FROM account_log as a JOIN users as u ON a.user_id = u.id " +
			"WHERE u.phone LIKE '%" + account.(string) + "%' " +
			"OR u.email LIKE '%" + account.(string) + "%' " +
			"OR u.account_number LIKE '%" + account.(string) + "%' "
	}

	if currency != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.currency = '%v' ", currency)
		} else {
			sqls += fmt.Sprintf("AND a.currency = '%v' ", currency)
		}

	}

	if types != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.type = '%v' ", types)
		} else {
			sqls += fmt.Sprintf("AND a.type = '%v' ", types)
		}
	}

	if startTime != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.created_time >= '%v' ", startTime)
		} else {
			sqls += fmt.Sprintf("AND a.created_time >= '%v' ", startTime)
		}
	}

	if endTime != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.created_time <= '%v' ", endTime)
		} else {
			sqls += fmt.Sprintf("AND a.created_time <= '%v' ", endTime)
		}
	}

	if sign != "" {

		if sign.(int) > 0 {

			if strings.Index(sqls, "WHERE") == -1 {
				sqls += fmt.Sprintf("WHERE a.value > 0 ")
			} else {
				sqls += fmt.Sprintf("AND a.value > 0 ")
			}
		} else {

			if strings.Index(sqls, "WHERE") == -1 {
				sqls += fmt.Sprintf("WHERE a.value < 0 ")
			} else {
				sqls += fmt.Sprintf("AND a.value < 0 ")
			}
		}
	}

	if strings.Index(sqls[len(sqls)-4:], "AND") != -1 {
		sqls = sqls[:len(sqls)-4]
	}

	sqls += "ORDER BY a.id DESC LIMIT ?,?"

	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	var total int
	row := stmt.QueryRow(offset, limit)
	row.Scan(&total)

	sqls = "SELECT a.* FROM account_log as a "

	if account != "" {
		sqls = "SELECT a.*, u.* FROM account_log as a JOIN users as u ON a.user_id = u.id " +
			"WHERE u.phone LIKE '%" + account.(string) + "%' " +
			"OR u.email LIKE '%" + account.(string) + "%' " +
			"OR u.account_number LIKE '%" + account.(string) + "%' "
	}

	if currency != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.currency = '%v' ", currency)
		} else {
			sqls += fmt.Sprintf("AND a.currency = '%v' ", currency)
		}

	}

	if types != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.type = '%v' ", types)
		} else {
			sqls += fmt.Sprintf("AND a.type = '%v' ", types)
		}
	}

	if startTime != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.created_time >= '%v' ", startTime)
		} else {
			sqls += fmt.Sprintf("AND a.created_time >= '%v' ", startTime)
		}
	}

	if endTime != "" {
		if strings.Index(sqls, "WHERE") == -1 {
			sqls += fmt.Sprintf("WHERE a.created_time <= '%v' ", endTime)
		} else {
			sqls += fmt.Sprintf("AND a.created_time <= '%v' ", endTime)
		}
	}

	if sign != "" {

		if sign.(int) > 0 {

			if strings.Index(sqls, "WHERE") == -1 {
				sqls += fmt.Sprintf("WHERE a.value > 0 ")
			} else {
				sqls += fmt.Sprintf("AND a.value > 0 ")
			}
		} else {

			if strings.Index(sqls, "WHERE") == -1 {
				sqls += fmt.Sprintf("WHERE a.value < 0 ")
			} else {
				sqls += fmt.Sprintf("AND a.value < 0 ")
			}
		}
	}

	if strings.Index(sqls[len(sqls)-4:], "AND") != -1 {
		sqls = sqls[:len(sqls)-4]
	}

	sqls += "ORDER BY a.id DESC LIMIT ?,?"

	log.Println("sqls : ", sqls)

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
