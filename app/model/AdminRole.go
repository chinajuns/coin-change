package model

import "okc/utils"

// AdminRole
// 管理员角色表
type AdminRole struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	IsSuper int    `json:"is_super"`
	Right   string `json:"right"`
}

// AddAdminRole
// 添加管理员角色
func (a *AdminRole) AddAdminRole(name, right, isSuper interface{}) error {
	db := utils.Db
	sqls := "INSERT INTO `admin_role` (`name`, `right`, `is_super`) VALUES (?,?,?)"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, right, isSuper)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAdminRoleById
// 根据id更新管理员
func (a *AdminRole) UpdateAdminRoleById() error {
	db := utils.Db
	sqls := "UPDATE  `admin_role` SET `name` = ?, `right` = ?, `is_super` = ? WHERE `id` = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Name, a.Right, a.IsSuper, a.Id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAdminRoleById
// 根据id删除管理员角色
func (a *AdminRole) DeleteAdminRoleById() error {
	db := utils.Db
	sqls := "DELETE FROM `admin_role` WHERE `id`=?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Id)
	if err != nil {
		return err
	}
	return nil
}

// QueryAdminRoleByName
// 根据name获取管理员角色
func (a *AdminRole) QueryAdminRoleByName(name interface{}) (*AdminRole, error) {
	db := utils.Db
	sqls := "SELECT `id`, `name`,`right`, `is_super` FROM `admin_role` WHERE name=?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	adminRole := new(AdminRole)
	row := stmt.QueryRow(name)
	row.Scan(&adminRole.Id, &adminRole.Name, &adminRole.Right, &adminRole.IsSuper)
	return adminRole, nil
}

// QueryAdminRoleById
// 根据id获取管理员角色
func (a *AdminRole) QueryAdminRoleById(id interface{}) (*AdminRole, error) {
	db := utils.Db
	sqls := "SELECT `id`, `name`,`right`, `is_super` FROM `admin_role` WHERE id=?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	adminRole := new(AdminRole)
	row := stmt.QueryRow(id)
	row.Scan(&adminRole.Id, &adminRole.Name, &adminRole.Right, &adminRole.IsSuper)
	return adminRole, nil
}

// QueryAdminRoleListPage
// 获取管理员角色列表
func (a *AdminRole) QueryAdminRoleListPage(page, limit int) ([]*AdminRole, int, error) {
	offset := (page - 1) * limit
	db := utils.Db
	sqls := "SELECT count(*) FROM admin_role ORDER BY `id` DESC LIMIT ?,?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	var total int
	row := stmt.QueryRow(offset, limit)
	row.Scan(&total)

	sqls = "SELECT `id`, `name`, `right`, 'is_super' FROM admin_role ORDER BY `id` DESC LIMIT ?,?"
	stmt, err = db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	data := make([]*AdminRole, 0)
	for rows.Next() {
		adminRole := new(AdminRole)
		rows.Scan(&adminRole.Id, &adminRole.Name, &adminRole.Right, &adminRole.IsSuper)
		data = append(data, adminRole)
	}

	return data, total, nil
}

// QueryAdminRoleAll
// 获取管理员列表
func (a *AdminRole) QueryAdminRoleAll() ([]*AdminRole, error) {

	db := utils.Db

	sqls := "SELECT `id`, `name`, `right`, 'is_super' FROM admin_role ORDER BY `id`"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	data := make([]*AdminRole, 0)
	for rows.Next() {
		adminRole := new(AdminRole)
		rows.Scan(&adminRole.Id, &adminRole.Name, &adminRole.Right, &adminRole.IsSuper)
		data = append(data, adminRole)
	}

	return data, nil
}
