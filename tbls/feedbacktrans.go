package tbls

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//访客留言表-多语言
type TFeedbackTrans struct {
	FbNo       int64  `gorm:"column:fb_no;type:bigint(20);not null;primary_key:true;"`
	FbName     string `gorm:"column:fb_name;type:varchar(200);not null;"`                    //姓名
	FbCompany  string `gorm:"column:fb_company;type:varchar(200);"`                          //公司
	FbAddress  string `gorm:"column:fb_address;type:varchar(200);"`                          //地址
	FbUrl      string `gorm:"column:fb_url;type:varchar(200);"`                              //网址
	FbTitle    string `gorm:"column:fb_title;type:varchar(200);"`                            //标题
	FbContent  string `gorm:"column:fb_content;type:text;"`                                  //内容
	FbLanguage string `gorm:"column:fb_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TFeedbackTrans) TableName() string {
	return "t_feedback_trans"
}

func (t *TFeedbackTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TFeedbackTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TFeedbackTrans{}
	}

	var res []TFeedbackTrans

	db.Where(`fb_no = ? AND fb_language = ?`,
		value, langCode).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TFeedbackTrans{}
	}
}

func (t *TFeedbackTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TFeedbackTrans{}).
		Where(`fb_no = ? AND fb_language = ?`,
			value, langCode).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
