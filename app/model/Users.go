package model

import (
	"fmt"
	"okc/utils"
	"strconv"
	"strings"
)

// Users
// 用户表模型
type Users struct {
	Id                           int     `json:"id,omitempty"`                              // 用户id
	AreaCodeId                   int     `json:"area_code_id,omitempty"`                    // 国家区号 1默认大陆
	AreaCode                     int     `json:"area_code,omitempty"`                       // 区号
	AccountNumber                string  `json:"account_number,omitempty"`                  // 账号
	Type                         int     `json:"type,omitempty"`                            // 类型
	Phone                        string  `json:"phone,omitempty"`                           // 手机
	AgentId                      int     `json:"agent_id,omitempty"`                        // 0表示不是代理商，1以上表示该代理商id
	AgentNoteId                  int     `json:"agent_note_id,omitempty"`                   // 代理商节点id。当该用户是代理商时该值等于上级代理商Id，当该用户不是代理商时该值等于节点代理商id
	ParentId                     int     `json:"parent_id,omitempty"`                       // 推荐人id
	Email                        string  `json:"email,omitempty"`                           // 邮箱
	Password                     string  `json:"password,omitempty"`                        // 密码
	PayPassword                  string  `json:"pay_password,omitempty"`                    // 支付密码
	Time                         int     `json:"time,omitempty"`                            // 时间戳
	HeadPortrait                 string  `json:"head_portrait,omitempty"`                   //
	ExtensionCode                string  `json:"extension_code,omitempty"`                  //
	Status                       int     `json:"status,omitempty"`                          // 1，已锁定
	GesturePassword              string  `json:"gesture_password,omitempty"`                //
	IsAuth                       string  `json:"is_auth,omitempty"`                         //
	NikeName                     string  `json:"nike_name,omitempty"`                       // 昵称
	WalletAddress                string  `json:"wallet_address,omitempty"`                  // 钱包地址
	IsBlackList                  string  `json:"is_black_list,omitempty"`                   //
	ParentPath                   string  `json:"parent_path,omitempty"`                     // 上级推荐人节点
	PushStatus                   int     `json:"push_status,omitempty"`                     // 0:未实名认证 1:实名认证  2:直推3人 3:直推5人 4:直推10人  5:直推30人  6:直推50人
	CandyNumber                  float64 `json:"candy_number,omitempty"`                    // 糖果数量
	ZhituiRealNumber             int     `json:"zhitui_real_number,omitempty"`              // 实名认证过的直推人数
	RealTeamnumber               int     `json:"real_teamnumber,omitempty"`                 // 实名认证通过的团队人数
	TopUpnumber                  float64 `json:"top_upnumber,omitempty"`                    // 团队业绩充值金额
	IsRealname                   int     `json:"is_realname,omitempty"`                     // 1:未实名认证过  2：实名认证过
	IsAtelier                    int     `json:"is_atelier,omitempty"`                      // 是否工作室
	NewIsrealTime                int     `json:"new_isreal_time,omitempty"`                 // 最新通过的下级实名认证时间
	TodyRealTeamnumber           int     `json:"tody_real_teamnumber,omitempty"`            // 今日新增团队实名认证人数
	TodayLegalDealCancelNum      int     `json:"today_LegalDealCancel_num,omitempty"`       // 今天c2c订单已经取消次数
	LegalDealCancelNumUpdateTime int     `json:"LegalDealCancel_num_update_time,omitempty"` // c2c取消单子更新时间
	Risk                         int     `json:"risk,omitempty"`                            // -1.亏损,0.正常,1.盈利
	LockTime                     int     `json:"lock_time,omitempty"`                       // 锁定时间
	Level                        int     `json:"level,omitempty"`                           // 代数
	Fund                         float64 `json:"fund,omitempty"`                            // 秒合约资产
	IsService                    int     `json:"is_service,omitempty"`                      // 是否是客服
	AgentPath                    string  `json:"agent_path,omitempty"`                      // 代理商关系
	WalletPwd                    string  `json:"wallet_pwd,omitempty"`                      // 钱包密码
	CountryCode                  string  `json:"country_code,omitempty"`                    //
	Label                        string  `json:"label,omitempty"`                           //
	Nationality                  string  `json:"nationality,omitempty"`                     //
	LastLoginIp                  string  `json:"last_login_ip,omitempty"`                   //
	Score                        string  `json:"score,omitempty"`                           //
	Remark                       string  `json:"remark,omitempty"`                          //
}

// AddUserByEmail
// 根据邮箱注册用户
func (u *Users) AddUserByEmail() error {

	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return err
	}
	sql := "INSERT INTO users (account_number, area_code_id, area_code, email, password, time, extension_code, wallet_pwd) VALUES(?,?,?,?,?,?,?,?)"
	stmt, err := t.Prepare(sql)
	if err != nil {
		t.Rollback()
		return err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			_ = utils.WriteErrorLog(err)
		}
	}()
	_, err = stmt.Exec(u.AccountNumber, u.AreaCodeId, u.AreaCode, u.Email, u.Password, u.Time, u.ExtensionCode, u.WalletPwd)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}

// CheckUserExistsByPhoneAndEmail
// 根据手机号或者邮箱查询用户是否存在
// true 存在 false 不存在
func (u *Users) CheckUserExistsByPhoneAndEmail(param string) (bool, error) {

	db := utils.Db
	sql := "SELECT id FROM users WHERE phone = ? OR email = ? OR account_number = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			_ = utils.WriteErrorLog(err)
		}
	}()
	row := stmt.QueryRow(param, param, param)
	users := Users{}
	_ = row.Scan(&users.Id)
	if users.Id != 0 {
		return true, nil
	}
	return false, err
}

// QueryUserIdByExtensionCode
// 根据推广码获取用户id
func (u *Users) QueryUserIdByExtensionCode(extensionCode string) (int, error) {

	db := utils.Db
	sql := "SELECT id FROM users WHERE extension_code = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			_ = utils.WriteErrorLog(err)
		}
	}()
	row := stmt.QueryRow(extensionCode)
	_ = row.Scan(&u.Id)
	return u.Id, nil
}

// QueryUserInfoById
// 根据用户id获取用户信息
func (u *Users) QueryUserInfoById(id int) (*Users, error) {

	user := new(Users)
	db := utils.Db
	sql := "SELECT * FROM users WHERE id = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	row := stmt.QueryRow(strconv.Itoa(id))

	fieldPtrSlice, err := utils.RefStructGetFieldPtr(user, "*")
	if err != nil {
		return nil, err
	}
	_ = row.Scan(fieldPtrSlice...)
	//_ = row.Scan(&user.Id, &user.Phone, &user.Password)

	return user, nil
}

// QueryUserInfoByUserAccount
// 根据用户账号获取用户信息
func (u *Users) QueryUserInfoByUserAccount(account string) (*Users, error) {

	user := new(Users)
	db := utils.Db
	sql := "SELECT * FROM users WHERE account_number = ?"

	stmt, err := db.Prepare(sql)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(account)
	fieldPtrSlice, err := utils.RefStructGetFieldPtr(user, "*")
	if err != nil {
		return nil, err
	}
	_ = row.Scan(fieldPtrSlice...)
	return user, nil
}

// UpdatePwdById
// 根据用户id修改密码
func (u *Users) UpdatePwdById(id interface{}, password string) error {
	db, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sql := "UPDATE users SET password = ? WHERE id = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		db.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(password, id)
	if err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

// UpdatePayPwdById
// 根据用户id修改交易密码
func (u *Users) UpdatePayPwdById(id interface{}, password string) error {
	db, err := utils.Db.Begin()
	if err != nil {
		return err
	}
	sql := "UPDATE users SET pay_password = ? WHERE id = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		db.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(password, id)
	if err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

// QueryUsersListPageByAccountAndNameAndRisk
// 根据account和name和risk获取用户分页列表
func (u *Users) QueryUsersListPageByAccountAndNameAndRisk(account, name string, risk, page, limit int) ([]map[string]interface{}, int, error) {
	offset := (page - 1) * limit
	db := utils.Db
	sqls := "SELECT count(*) FROM users as u LEFT JOIN user_real as e ON u.id = e.user_id "

	if account != "" {

		if strings.Index(sqls, "WHERE") == -1 {

			sqls += fmt.Sprintf("WHERE u.phone = '%v' OR u.email = '%v' OR u.account_number = '%v' ", account, account, account)
		} else {
			sqls += fmt.Sprintf("AND u.phone = '%v' OR u.email = '%v' OR u.account_number = '%v' ", account, account, account)
		}
	}

	if name != "" {

		if strings.Index(sqls, "WHERE") == -1 {

			sqls += fmt.Sprintf("WHERE e.name = '%v' ", name)
		} else {
			sqls += fmt.Sprintf("AND e.name = '%v' ", name)
		}
	}

	if risk != -2 {

		if strings.Index(sqls, "WHERE") == -1 {

			sqls += fmt.Sprintf("WHERE u.risk = '%v' ", risk)
		} else {
			sqls += fmt.Sprintf("AND u.risk = '%v' ", risk)
		}
	}

	if strings.Index(sqls[len(sqls)-5:], "AND") != -1 {
		sqls = sqls[:len(sqls)-5]
	}
	sqls += "ORDER BY u.id DESC LIMIT ?,?"
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	var total int
	row := stmt.QueryRow(offset, limit)
	row.Scan(&total)

	sqls = "SELECT u.*, e.* FROM users as u LEFT JOIN user_real as e ON u.id = e.user_id "
	if account != "" {

		if strings.Index(sqls, "WHERE") == -1 {

			sqls += fmt.Sprintf("WHERE u.phone = '%v' OR u.email = '%v' OR u.account_number = '%v' ", account, account, account)
		} else {
			sqls += fmt.Sprintf("AND u.phone = '%v' OR u.email = '%v' OR u.account_number = '%v' ", account, account, account)
		}
	}

	if name != "" {

		if strings.Index(sqls, "WHERE") == -1 {

			sqls += fmt.Sprintf("WHERE e.name = '%v' ", name)
		} else {
			sqls += fmt.Sprintf("AND e.name = '%v' ", name)
		}
	}

	if risk != -2 {

		if strings.Index(sqls, "WHERE") == -1 {

			sqls += fmt.Sprintf("WHERE u.risk = '%v' ", risk)
		} else {
			sqls += fmt.Sprintf("AND u.risk = '%v' ", risk)
		}
	}

	if strings.Index(sqls[len(sqls)-5:], "AND") != -1 {
		sqls = sqls[:len(sqls)-5]
	}
	sqls += "ORDER BY u.id DESC LIMIT ?,?"
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
