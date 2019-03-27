package tbls

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type TPermission struct {
	PermNo        int64     `gorm:"column:perm_no;AUTO_INCREMENT;primary_key:true;"`
	PermParentNo  int64     `gorm:"column:perm_parent_no;not null;"`
	PermName      string    `gorm:"column:perm_name;size:32;not null;"`
	PermKey       string    `gorm:"column:perm_key;size:32;not null;"`
	PermMemo      string    `gorm:"column:perm_memo;type:varchar(200);"`
	PermInputdate time.Time `gorm:"column:perm_inputdate"`
	PermInputid   int64     `gorm:"column:perm_inputid"`
	PermUpdate    time.Time `gorm:"column:perm_update"`
	PermUpid      int64     `gorm:"column:perm_upid"`
	PermDelflg    string    `gorm:"column:perm_delflg;type:char(1);"`

	Tbls
}

func (t *TPermission) TableName() string {
	return "m_permission"
}

func (t *TPermission) GetList(db *gorm.DB) []TPermission {

	var res []TPermission

	db.Where(`perm_delflg = '0'`).Order(`perm_no`).Find(&res)

	return res
}

func (t *TPermission) GetData(db *gorm.DB, searchNo string) (bool, TPermission) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TPermission{}
	}

	var res []TPermission

	db.Where(`perm_delflg = '0' AND perm_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TPermission{}
	}
}

func (t *TPermission) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TPermission{}).
		Where(`perm_delflg = '0' AND perm_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TPermission) CheckKey(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TPermission{}).
		Where(`perm_delflg = '0' AND perm_key = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
