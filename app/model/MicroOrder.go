package model

import (
	"fmt"
	"okc/utils"
	"time"
)

// MicroOrder
// 期权订单表
type MicroOrder struct {
	Id              int     `json:"id"`
	UserId          int     `json:"user_id"`
	MatchId         int     `json:"match_id"`
	CurrencyId      int     `json:"currency_id"`
	Type            int     `json:"type"`
	IsInsurance     int     `json:"is_insurance"`
	Seconds         int     `json:"seconds"`
	Number          float64 `json:"number"`
	OpenPrice       float64 `json:"open_price"`
	EndPrice        float64 `json:"end_price"`
	Fee             float64 `json:"fee"`
	ProfitRatio     float64 `json:"profit_ratio"`
	FactProfits     float64 `json:"fact_profits"`
	Status          int     `json:"status"`
	PreProfitResult int     `json:"pre_profit_result"`
	ProfitResult    int     `json:"profit_result"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	HandledAt       string  `json:"handled_at"`
	CompleteAt      string  `json:"complete_at"`
	ReturnAt        string  `json:"return_at"`
	AgentPath       string  `json:"agent_path"`
}

// AddMicroOder
// 添加期权订单
func (m *MicroOrder) AddMicroOder() error {
	if m.CompleteAt == "" {
		m.CompleteAt = time.Now().Format("2006-01-02 15:04:05")
	}
	if m.UpdatedAt == "" {
		m.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return err
	}
	sqls := "INSERT INTO micro_orders (user_id, match_id, currency_id, type, is_insurance, seconds," +
		"number, open_price, end_price, fee, profit_ratio, fact_profits, status, pre_profit_result," +
		"profit_result, created_at, updated_at, handled_at, complete_at, return_at, agent_path) VALUES " +
		"(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(m.UserId, m.MatchId, m.CurrencyId, m.Type, m.IsInsurance, m.Seconds,
		m.Number, m.OpenPrice, m.EndPrice, m.Fee, m.ProfitRatio, m.FactProfits, m.Status,
		m.PreProfitResult, m.ProfitResult, m.CreatedAt, m.UpdatedAt, m.HandledAt, m.CompleteAt,
		m.ReturnAt, m.AgentPath)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// PageByUserIdAndStatusAndMatchIdAndCurrencyId
// 根据userId和status和matchId和currencyId获取分页
func (m *MicroOrder) PageByUserIdAndStatusAndMatchIdAndCurrencyId(userId, status, matchId, currencyId, page, limit int) (data []map[string]interface{}, total int, err error) {
	// 偏移量
	offset := (page - 1) * limit

	sqls := "SELECT count(*) FROM micro_orders WHERE user_id = ? AND "
	if status != -1 {
		sqls += fmt.Sprintf("status = %d AND ", status)
	}

	if matchId != -1 {
		sqls += fmt.Sprintf("status = %d AND ", matchId)
	}

	if currencyId != -1 {
		sqls += fmt.Sprintf("status = %d AND ", currencyId)
	}
	sqls = sqls[:len(sqls)-4]
	sqls += "ORDER BY id DESC LIMIT ?,?"

	db := utils.Db
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(userId, offset, limit)
	row.Scan(&total)

	sqls = "SELECT * FROM micro_orders WHERE user_id = ? AND "
	if status != -1 {
		sqls += fmt.Sprintf("status = %d AND ", status)
	}

	if matchId != -1 {
		sqls += fmt.Sprintf("status = %d AND ", matchId)
	}

	if currencyId != -1 {
		sqls += fmt.Sprintf("status = %d AND ", currencyId)
	}
	sqls = sqls[:len(sqls)-4]
	sqls += "ORDER BY id DESC LIMIT ?,?"

	db = utils.Db
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
	data = make([]map[string]interface{}, 0)
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
	return data, total, nil
}

// QueryMicroOrderCountByUserIdAndCurrencyId
// 根据用户id和币种id获取期权订单数量
func (m *MicroOrder) QueryMicroOrderCountByUserIdAndCurrencyId(userId, currencyId interface{}) (int, error) {
	db := utils.Db
	sqls := "SELECT count(*) FROM micro_orders WHERE user_id = ? AND currency_id = ? AND status = 1"
	var count int
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(userId, currencyId)
	row.Scan(&count)
	return count, nil
}
