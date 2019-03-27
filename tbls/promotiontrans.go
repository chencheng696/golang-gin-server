package tbls

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//网络推广表-多语言
type TPromotionTrans struct {
	PromNo          int64  `gorm:"column:prom_no;type:bigint(20);not null;primary_key:true;"`
	PromName        string `gorm:"column:prom_name;type:varchar(200);not null;"`                   //名称
	PromUrl         string `gorm:"column:prom_url;type:varchar(200);"`                             //网址
	PromDescription string `gorm:"column:prom_description;type:text;not null;"`                    //描述
	PromLanguage    string `gorm:"column:prm_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TPromotionTrans) TableName() string {
	return "m_promotion_trans"
}

func (t *TPromotionTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TPromotionTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TPromotionTrans{}
	}

	var res []TPromotionTrans

	db.Where(`prom_no = ? AND prom_language = ?`,
		value, langCode).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TPromotionTrans{}
	}
}

func (t *TPromotionTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TPromotionTrans{}).
		Where(`prom_no = ? AND prom_language = ?`,
			value, langCode).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
