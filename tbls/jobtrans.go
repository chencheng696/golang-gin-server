package tbls

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//招聘表-多语言
type TJobTrans struct {
	JobNo          int64  `gorm:"column:job_no;type:bigint(20);not null;primary_key:true;"`
	JobName        string `gorm:"column:job_name;type:varchar(200);not null;"`                    //招聘标题或对象
	JobAddress     string `gorm:"column:job_address;type:varchar(200);not null;"`                 //工作地点
	JobDescription string `gorm:"column:job_description;type:text;not null;"`                     //招聘要求描述
	JobLanguage    string `gorm:"column:job_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TJobTrans) TableName() string {
	return "m_job_trans"
}

func (t *TJobTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TJobTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TJobTrans{}
	}

	var res []TJobTrans

	db.Where(`job_no = ? AND job_language = ?`,
		value, langCode).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TJobTrans{}
	}
}

func (t *TJobTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TJobTrans{}).
		Where(`job_no = ? AND job_language = ?`,
			value, langCode).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
