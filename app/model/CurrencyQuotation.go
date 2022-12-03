package model

// CurrencyQuotation
// 行情表
type CurrencyQuotation struct {
	Id         int     `json:"id"`          //
	MatchId    int     `json:"match_id"`    // 交易对id
	LegalId    int     `json:"legal_id"`    //
	CurrencyId int     `json:"currency_id"` // 币种id
	Change     string  `json:"change"`      // 涨跌幅 带+ - 号
	Volume     float64 `json:"volume"`      // 成交量
	NowPrice   float64 `json:"now_price"`   // 当前价位
	AddTime    int     `json:"add_time"`    //
	Xm         int     `json:"xm"`          //
}
