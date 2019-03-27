package tbls

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TSystemTrans struct {
	SysKey      string    `gorm:"column:sys_key;type:varchar(32);not null;primary_key:true;"`
	SysVal      string    `gorm:"column:sys_val;type:text;"`
	SysUpdate   time.Time `gorm:"column:sys_update;"`
	SysUpid     int64     `gorm:"column:sys_upid;"`
	SysLanguage string    `gorm:"column:sys_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TSystemTrans) TableName() string {
	return "t_system_trasn"
}

func (t *TSystemTrans) ReadData(db *gorm.DB, k, l, def string) string {

	var res []TSystemTrans

	db.Where(`sys_key = ? AND sys_language = ?`, k, l).Find(&res)

	if len(res) > 0 {
		return res[0].SysVal
	} else {
		return def
	}
}

func (t *TSystemTrans) WriteData(db *gorm.DB, k, l, v string) bool {

	t.SysKey = k
	t.SysVal = v
	t.SysUpdate = time.Now()
	t.SysUpid = 0
	t.SysLanguage = l

	if !t.CheckKey(db, k, l) {
		db.Create(&t)
	} else {
		db.Model(&t).
			Where(`sys_key = ? AND sys_language = ?`, k, l).
			Updates(TSystemTrans{
				SysVal:    t.SysVal,
				SysUpdate: t.SysUpdate,
				SysUpid:   t.SysUpid})
	}
	return true
}

func (t *TSystemTrans) CheckKey(db *gorm.DB, k, l string) bool {

	count := 0

	db.Model(&TSystemTrans{}).
		Where(`sys_key = ? AND sys_language = ?`, k, l).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
