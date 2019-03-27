package tbls

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//合作伙伴表-多语言
type TPartnersTrans struct {
	PnNo          int64  `gorm:"column:pn_no;type:bigint(20);not null;primary_key:true;"`
	PnName        string `gorm:"column:pn_name;type:varchar(200);not null;"`                    //名称
	PnUrl         string `gorm:"column:pn_url;type:varchar(200);"`                              //网址
	PnDescription string `gorm:"column:pn_description;type:text;not null;"`                     //描述
	PnLanguage    string `gorm:"column:pn_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TPartnersTrans) TableName() string {
	return "m_partners_trans"
}

func (t *TPartnersTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TPartnersTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TPartnersTrans{}
	}

	var res []TPartnersTrans

	db.Where(`pn_no = ? AND pn_language = ?`,
		value, langCode).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TPartnersTrans{}
	}
}

func (t *TPartnersTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TPartnersTrans{}).
		Where(`pn_no = ? AND pn_language = ?`,
			value, langCode).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
