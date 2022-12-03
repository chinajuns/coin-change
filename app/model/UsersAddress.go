package model

import (
	"database/sql"
	"okc/utils"
)

// UsersAddress
// 用户地址表
type UsersAddress struct {
	Id       int    `json:"id,omitempty"`       //
	Currency int    `json:"currency,omitempty"` // 币种id
	UserId   int    `json:"user_id,omitempty"`  // 用户id
	NetType  string `json:"net_type,omitempty"` // 帐号所在网络类型
	Notes    string `json:"notes,omitempty"`    //
	Address  string `json:"address,omitempty"`  // 提币地址
}

// AddUserAddressByUserId
// 根据用户id添加用户地址
func (u *UsersAddress) AddUserAddressByUserId(userId int, address string) error {
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sqls := "INSERT INTO address (user_id, address, currency) VALUES (?,?,?)"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, address, 0)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// UpdateUserAddressByUserId
// 根据用户id更新用户地址
func (u *UsersAddress) UpdateUserAddressByUserId(userId int, address string) error {
	t, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sqls := "UPDATE address SET address = ? WHERE user_id = ?"
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(address, userId)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// QueryFindByUserId
// 根据用户id获取一条数据
func (u *UsersAddress) QueryFindByUserId(userId int) (*UsersAddress, error) {
	db := utils.Db
	sqls := "SELECT id, user_id, address FROM address WHERE user_id = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(userId)
	usersAddress := new(UsersAddress)
	err = row.Scan(&usersAddress.Id, &usersAddress.UserId, &usersAddress.Address)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return usersAddress, nil
}
