package tbls

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TSystem struct {
	SysKey    string    `gorm:"column:sys_key;type:varchar(32);not null;primary_key:true;"`
	SysVal    string    `gorm:"column:sys_val;type:text;"`
	SysUpdate time.Time `gorm:"column:sys_update;"`
	SysUpid   int64     `gorm:"column:sys_upid;"`
}

func (t *TSystem) TableName() string {
	return "t_system"
}

func (t *TSystem) ReadData(db *gorm.DB, k, def string) string {

	var res []TSystem

	db.Where(`sys_key = ?`, k).Find(&res)

	if len(res) > 0 {
		return res[0].SysVal
	} else {
		return def
	}
}

func (t *TSystem) WriteData(db *gorm.DB, k, v string) bool {

	t.SysKey = k
	t.SysVal = v
	t.SysUpdate = time.Now()
	t.SysUpid = 0

	if !t.CheckKey(db, k) {
		db.Create(&t)
	} else {
		db.Model(&t).
			Where("sys_key = ?", t.SysKey).
			Updates(TSystem{
				SysVal:    t.SysVal,
				SysUpdate: t.SysUpdate,
				SysUpid:   t.SysUpid})
	}
	return true
}

func (t *TSystem) CheckKey(db *gorm.DB, k string) bool {

	count := 0

	db.Model(&TSystem{}).
		Where(`sys_key = ?`, k).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
