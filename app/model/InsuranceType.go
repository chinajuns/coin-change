package model

import "okc/utils"

// InsuranceType
// 币保险类型
type InsuranceType struct {
	Id                         int     `json:"id"`
	Name                       string  `json:"name"`
	CurrencyId                 int     `json:"currency_id"`
	Type                       int     `json:"type"`
	MinAmount                  float64 `json:"min_amount"`
	MaxAmount                  float64 `json:"max_amount"`
	InsuranceAssets            float64 `json:"insurance_assets"`
	ProfitTerminationCondition float64 `json:"profit_termination_condition"`
	DefectiveClaimsCondition   float64 `json:"defective_claims_condition"`
	DefectiveClaimsCondition2  float64 `json:"defective_claims_condition2"`
	ClaimsTimesDaily           int     `json:"claims_times_daily"`
	AutoClaim                  int     `json:"auto_claim"`
	ClaimRate                  float64 `json:"claim_rate"`
	ClaimDirection             int     `json:"claim_direction"`
	Status                     int     `json:"status"`
	IsTAdd1                    int     `json:"is_t_add_1"`
}

// QueryInsuranceTypeByCurrencyId
// 根据币种id获取币种险种
func (i InsuranceType) QueryInsuranceTypeByCurrencyId(currencyId interface{}) (*InsuranceType, error) {
	db := utils.Db
	sqls := "SELECT * FROM insurance_types WHERE currency_id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	insuranceType := new(InsuranceType)
	insuranceTypePtr, err := utils.RefStructGetFieldPtr(insuranceType, "*")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(currencyId)
	row.Scan(insuranceTypePtr...)
	return insuranceType, nil

}
