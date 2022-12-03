package model

import (
	"errors"
	"okc/utils"
)

// Setting
// 系统设置表
type Setting struct {
	Id    int    `json:"id,omitempty"`    //
	Key   string `json:"key,omitempty"`   //
	Value string `json:"value,omitempty"` //
	Notes string `json:"notes,omitempty"` //
}

// QueryValueByKey
// 根据key获取value
func (s *Setting) QueryValueByKey(key, def string) (string, error) {
	if key == "" {
		return def, nil
	}
	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return "", err
	}
	sqls := "SELECT `id`, `key`, `value`, notes FROM `settings` WHERE `key` = ? FOR UPDATE "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return "", err
	}
	defer stmt.Close()
	setting := new(Setting)
	row := stmt.QueryRow(key)
	row.Scan(&setting.Id, &setting.Key, &setting.Value, &setting.Notes)
	t.Commit()

	if setting.Value == "" {
		return def, nil
	}

	return setting.Value, nil
}

// UpdateValueByKey
// 根据key更新value
func (s *Setting) UpdateValueByKey(key, value string) error {
	if key == "" {
		return errors.New("UpdateValueByKey key parameter not nil")
	}
	t, err := utils.Db.Begin()
	if err != nil {
		t.Rollback()
		return err
	}
	sqls := "UPDATE `settings` SET `value` = ? WHERE `key` = ? "
	stmt, err := t.Prepare(sqls)
	if err != nil {
		t.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, key)
	if err != nil {
		t.Rollback()
		return err
	}
	t.Commit()
	return nil
}
