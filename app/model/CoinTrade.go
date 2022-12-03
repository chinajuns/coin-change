package model

import (
	"fmt"
	"okc/utils"
	"time"
)

// CoinTrade
// 币币交易订单表
type CoinTrade struct {
	Id          int     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UId         int     `gorm:"column:u_id" json:"u_id"`
	CurrencyId  int     `gorm:"column:currency_id" json:"currency_id"`
	LegalId     int     `gorm:"column:legal_id" json:"legal_id"`
	Type        int     `gorm:"column:type" json:"type"`
	TargetPrice float64 `gorm:"column:target_price" json:"target_price"`
	TradePrice  float64 `gorm:"column:trade_price" json:"trade_price"`
	TradeAmount float64 `gorm:"column:trade_amount" json:"trade_amount"`
	ChargeFee   float64 `gorm:"column:charge_fee" json:"charge_fee"`
	Status      int     `gorm:"column:status;default:1" json:"status"`
	CreatedAt   string  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   string  `gorm:"column:updated_at" json:"updated_at"`
}

// AddCoinTrade
// 添加币币交易订单
func (c *CoinTrade) AddCoinTrade() error {
	if c.CreatedAt == "" {
		c.CreatedAt = time.Now().Format("2006-01-02 15:04-05")
	}
	if c.UpdatedAt == "" {
		c.UpdatedAt = time.Now().Format("2006-01-02 15:04-05")
	}
	db := utils.Db
	sqls := "INSERT INTO coin_trade (u_id, currency_id, legal_id, type, target_price, trade_price," +
		"trade_amount, charge_fee, status, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.UId, c.CurrencyId, c.LegalId, c.Type, c.TargetPrice, c.TradePrice,
		c.TradeAmount, c.ChargeFee, c.Status, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCoinTradeById
// 根据id更新币币交易订单
func (c *CoinTrade) UpdateCoinTradeById() error {
	if c.UpdatedAt == "" {
		c.UpdatedAt = time.Now().Format("2006-01-02 15:04-05")
	}
	db := utils.Db
	sqls := "UPDATE coin_trade SET u_id = ? , currency_id = ?, legal_id = ?, type = ?, target_price = ?," +
		"trade_price = ? , trade_amount = ? , charge_fee = ?, status = ?, created_at = ?, updated_at = ? " +
		"WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.UId, c.CurrencyId, c.LegalId, c.Type, c.TargetPrice, c.TradePrice,
		c.TradeAmount, c.ChargeFee, c.Status, c.CreatedAt, c.UpdatedAt, c.Id)
	if err != nil {
		return err
	}
	return nil
}

// QueryCoinTradeById
// 根据id获取币币交易订单
func (c *CoinTrade) QueryCoinTradeById(id interface{}) (*CoinTrade, error) {
	db := utils.Db
	sqls := "SELECT `id`,`u_id`, `currency_id`, `legal_id`, `type`, `target_price`, `trade_price`, `trade_amount`, `charge_fee`, `status`, `created_at`, `updated_at` FROM coin_trade WHERE id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	coin := new(CoinTrade)
	row := stmt.QueryRow(id)
	row.Scan(&coin.Id, &coin.UId, &coin.CurrencyId, &coin.LegalId, &coin.Type, &coin.TargetPrice, &coin.TradePrice,
		&coin.TradeAmount, &coin.ChargeFee, &coin.Status, &coin.CreatedAt, &coin.UpdatedAt)
	return coin, nil
}

// QueryCoinTradePageByCurrencyIdAndLegalIdAndStatus
// 获取币币交易订单分页
func (c *CoinTrade) QueryCoinTradePageByCurrencyIdAndLegalIdAndStatus(currencyId, legalId, status string, page, limit int) ([]map[string]interface{}, error) {
	offset := (page - 1) * limit
	db := utils.Db
	str := "WHERE "
	if currencyId != "" {
		str += fmt.Sprintf("currency_id = %s AND ", currencyId)
	}
	if legalId != "" {
		str += fmt.Sprintf("legal_id = %s AND ", legalId)
	}
	if status != "" {
		str += fmt.Sprintf("status = %s AND ", status)
	}
	if len(str) == 6 {
		str = ""
	} else {
		str = str[:len(str)-4]
	}
	sqls := fmt.Sprintf("SELECT * FROM coin_trade %s ORDER BY id DESC LIMIT ?,? ", str)
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(offset, limit)
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
