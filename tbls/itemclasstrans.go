package tbls

import (
	"database/sql"
	"strconv"

	"github.com/jinzhu/gorm"
)

//商品分类表-多语言
type TItemClassTrans struct {
	ItemClassNo       int64  `gorm:"column:item_class_no;type:bigint(20);not null;primary_key:true;"`
	ItemClassName     string `gorm:"column:item_class_name;type:varchar(200);not null;"`
	ItemClassPicts    string `gorm:"column:item_class_picts;type:text;"`                                    //图片路径，多张以分号分割
	ItemClassLanguage string `gorm:"column:item_class_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TItemClassTrans) TableName() string {
	return "m_item_class_trans"
}

func (t *TItemClassTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TItemClassTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TItemClassTrans{}
	}

	var res []TItemClassTrans

	db.Where(`item_class_no = ? AND item_class_language = ?`, value, langCode).
		Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TItemClassTrans{}
	}
}

func (t *TItemClassTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TItemClassTrans{}).
		Where(`item_class_no = ? AND item_class_language = ?`, value, langCode).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
