package tbls

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//新闻分类表-多语言
type TNewsClassTrans struct {
	NewsClassNo       int64  `gorm:"column:news_class_no;type:bigint(20);not null;primary_key:true;"`
	NewsClassName     string `gorm:"column:news_class_name;type:varchar(200);not null;"`
	NewsClassPicts    string `gorm:"column:news_class_picts;type:text;"`                                    //图片路径，多张以分号分割
	NewsClassLanguage string `gorm:"column:news_class_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TNewsClassTrans) TableName() string {
	return "m_news_class_trans"
}

func (t *TNewsClassTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TNewsClassTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TNewsClassTrans{}
	}

	var res []TNewsClassTrans

	db.Where(`news_class_no = ? AND news_class_language = ?`, value, langCode).
		Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TNewsClassTrans{}
	}
}

func (t *TNewsClassTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TNewsClassTrans{}).
		Where(`news_class_no = ? AND news_class_language = ?`, value, langCode).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
