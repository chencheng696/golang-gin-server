package tbls

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

//语言表
type TLanguage struct {
	LngCode       string    `gorm:"column:lng_code;AUTO_INCREMENT;primary_key:true;"` //语言code 例如cn
	LngName       string    `gorm:"column:lng_name;type:varchar(50);not null;"`       //语言名字
	LngShowName   string    `gorm:"column:lng_show_name;type:varchar(50);not null;"`  //显示语言名字
	LngFontFamily string    `gorm:"column:lng_fontfamily;type:varchar(1000);"`        //语言字体(暂留)
	LngSeq        int64     `gorm:"column:lng_seq;not null;"`                         //显示顺序
	LngInputdate  time.Time `gorm:"column:lng_inputdate;not null;"`
	LngDelflg     string    `gorm:"column:lng_delflg;type:char(1);"`
}

func (t *TLanguage) TableName() string {
	return "m_language"
}

func (t *TLanguage) GetList(db *gorm.DB) []TLanguage {

	var res []TLanguage

	dbnew := t.GetListWhere(db)
	dbnew.Order(`lng_seq`).Find(&res)
	return res
}

func (t *TLanguage) GetListNative(db *gorm.DB) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db)
	return dbnew.Model(&TLanguage{}).Select("*").Rows()
}

func (t *TLanguage) GetListWhere(db *gorm.DB) *gorm.DB {
	return db.Where(`lng_delflg = '0'`).Order(`lng_seq`)
}

func (t *TLanguage) GetData(db *gorm.DB, code string) (bool, TLanguage) {
	var res []TLanguage

	db.Where(`lng_delflg = '0' AND lng_code = ?`, code).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TLanguage{}
	}
}
