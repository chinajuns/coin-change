package model

import "okc/utils"

// UserReal
// 用户认证表模型
type UserReal struct {
	Id           int    `json:"id,omitempty"`      // id
	UserId       int    `json:"user_id,omitempty"` // 用户id
	Name         string `json:"name,omitempty"`    // 名称
	CardId       int    `json:"card_id"`           // 卡id
	Type         int    `json:"type"`              // 1初级认证 2高级认证
	ReviewStatus int    `json:"review_status"`     // 1,未审核2,已审核
	FrontPic     string `json:"front_pic"`         //
	HandPic      string `json:"hand_pic"`          //
	CreateTime   int    `json:"create_time"`       // 创建时间
	UpdateTime   int    `json:"update_time"`       // 更新时间
}

// QueryFirstDataByUserId
// 根据用户ID获取一条数据
func (u *UserReal) QueryFirstDataByUserId(id interface{}) (*UserReal, error) {
	db := utils.Db
	sql := "SELECT * FROM user_real WHERE id = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	userReal := new(UserReal)
	row := stmt.QueryRow(id)
	userRealPtr, err := utils.RefStructGetFieldPtr(userReal, "*")
	if err != nil {
		return nil, err
	}
	row.Scan(userRealPtr...)
	return userReal, nil
}

// QueryFirstDataByUserIdAndTypes
// 根据用户ID和类型获取一条数据
func (u *UserReal) QueryFirstDataByUserIdAndTypes(id interface{}, types int) (*UserReal, error) {
	db := utils.Db
	sql := "SELECT * FROM user_real WHERE id = ? AND type = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	userReal := new(UserReal)
	row := stmt.QueryRow(id, types)
	userRealPtr, err := utils.RefStructGetFieldPtr(userReal, "*")
	if err != nil {
		return nil, err
	}
	row.Scan(userRealPtr...)
	return userReal, nil
}
