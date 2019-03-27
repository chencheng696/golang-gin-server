package tbls

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//案例分类表-多语言
type TCaseClassTrans struct {
	CaseClassNo       int64  `gorm:"column:case_class_no;type:bigint(20);not null;primary_key:true;"`
	CaseClassName     string `gorm:"column:case_class_name;type:varchar(200);not null;"`
	CaseClassPicts    string `gorm:"column:case_class_picts;type:text;"`                                    //图片路径，多张以分号分割
	CaseClassLanguage string `gorm:"column:case_class_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TCaseClassTrans) TableName() string {
	return "m_case_class_trans"
}

func (t *TCaseClassTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TCaseClassTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TCaseClassTrans{}
	}

	var res []TCaseClassTrans

	db.Where(`case_class_no = ? AND case_class_language = ?`, value, langCode).
		Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TCaseClassTrans{}
	}
}

func (t *TCaseClassTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TCaseClassTrans{}).
		Where(`case_class_no = ? AND case_class_language = ?`, value, langCode).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
