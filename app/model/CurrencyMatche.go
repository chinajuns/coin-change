package model

import "okc/utils"

// CurrencyMatche
// 币种交易对表

type CurrencyMatches struct {
	Id              int     `json:"id"`
	LegalId         int     `json:"legal_id"`
	CurrencyId      int     `json:"currency_id"`
	IsDisplay       int     `json:"is_display"`
	MarketFrom      int     `json:"market_from"`
	OpenTransaction int     `json:"open_transaction"`
	OpenLever       int     `json:"open_lever"`
	OpenMicrotrade  int     `json:"open_microtrade"`
	OpenCoinTrade   int     `json:"open_coin_trade"`
	Sort            int     `json:"sort"`
	MicroTradeFee   float64 `json:"micro_trade_fee"`
	LeverShareNum   float64 `json:"lever_share_num"`
	Spread          float64 `json:"spread"`
	Overnight       float64 `json:"overnight"`
	LeverTradeFee   float64 `json:"lever_trade_fee"`
	LeverMinShare   int     `json:"lever_min_share"`
	LeverMaxShare   int     `json:"lever_max_share"`
	FluctuateMin    float64 `json:"fluctuate_min"`
	FluctuateMax    float64 `json:"fluctuate_max"`
	RiskGroupResult int     `json:"risk_group_result"`
	CreateTime      int     `json:"create_time"`
}

// QueryCurrencyMatchesByid
// 根据id获取币种交易对
func (c *CurrencyMatches) QueryCurrencyMatchesByid(id interface{}) (*CurrencyMatches, error) {

	db := utils.Db
	sqls := "SELECT * FROM currency_matches WHERE id = ? "
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	currencyMatches := new(CurrencyMatches)
	currencyMatchesPtr, err := utils.RefStructGetFieldPtr(currencyMatches, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	row.Scan(currencyMatchesPtr...)
	return currencyMatches, nil

}

// QueryCurrencyMatchesByCurrencyIdAndLegalId
// 根据币种id和法币id获取币种交易对
func (c *CurrencyMatches) QueryCurrencyMatchesByCurrencyIdAndLegalId(currencyId, legalId interface{}) (*CurrencyMatches, error) {
	db := utils.Db
	sqls := "SELECT * FROM currency_matches WHERE currency_id = ? AND legal_id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	currencyMatches := new(CurrencyMatches)
	currencyMatchesPtr, err := utils.RefStructGetFieldPtr(currencyMatches, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(currencyId, legalId)
	row.Scan(currencyMatchesPtr...)
	return currencyMatches, nil
}

// QueryCurrencyMatchesByCurrencyIdAndLegalIdAndIsDisplay
// 根据币种id和法币id获取币种交易对
func (c *CurrencyMatches) QueryCurrencyMatchesByCurrencyIdAndLegalIdAndIsDisplay(currencyId, legalId interface{}) (*CurrencyMatches, error) {
	db := utils.Db
	sqls := "SELECT * FROM currency_matches WHERE currency_id = ? AND legal_id = ? AND is_display = 1"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	currencyMatches := new(CurrencyMatches)
	currencyMatchesPtr, err := utils.RefStructGetFieldPtr(currencyMatches, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(currencyId, legalId)
	row.Scan(currencyMatchesPtr...)
	return currencyMatches, nil
}
