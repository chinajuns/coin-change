package model

import "okc/utils"

// Admin
// 后台管理员账号表
type Admin struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	RoleId   int    `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
}

// QueryAdminByUsername
// 根据账号查询管理员
func (a *Admin) QueryAdminByUsername(username string) (*Admin, error) {
	db := utils.Db
	sqls := "SELECT `id`, `username`, `password`, `role_id` FROM admin WHERE `username` = ?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	admin := new(Admin)
	row := stmt.QueryRow(username)
	row.Scan(&admin.Id, &admin.Username, &admin.Password, &admin.RoleId)
	return admin, nil
}

// AddAdmin
// 添加管理员
func (a *Admin) AddAdmin() error {
	db := utils.Db
	sqls := "INSERT INTO admin (`username`, `password`, `role_id`) VALUES (?,?,?)"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Username, a.Password, a.RoleId)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAdminById
// 根据id更新管理员
func (a *Admin) UpdateAdminById() error {
	db := utils.Db
	sqls := "UPDATE  admin SET username=?, password=?, role_id=? WHERE id=?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Username, a.Password, a.RoleId, a.Id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAdminById
// 根据id删除管理员
func (a *Admin) DeleteAdminById() error {
	db := utils.Db
	sqls := "DELETE FROM admin WHERE id=?"
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

// QueryAdminListPage
// 获取管理员列表
func (a *Admin) QueryAdminListPage(page, limit int) ([]*Admin, int, error) {
	offset := (page - 1) * limit
	db := utils.Db
	sqls := "SELECT count(*) FROM admin ORDER BY `id` DESC LIMIT ?,?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	var total int
	row := stmt.QueryRow(offset, limit)
	row.Scan(&total)

	sqls = "SELECT a.`id`, a.`username`, a.`role_id`, r.name as role_name FROM admin as a LEFT JOIN admin_role as r ON a.role_id = r.id ORDER BY a.`id` DESC LIMIT ?,?"
	stmt, err = db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	data := make([]*Admin, 0)
	for rows.Next() {
		admin := new(Admin)
		rows.Scan(&admin.Id, &admin.Username, &admin.RoleId, &admin.RoleName)
		data = append(data, admin)
	}

	return data, total, nil
}
