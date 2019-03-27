package tbls

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//商品表-多语言
type TItemTrans struct {
	ItemNo       int64  `gorm:"column:item_no;type:bigint(20);not null;primary_key:true;"`
	ItemName     string `gorm:"column:item_name;type:varchar(1000);not null;"`                   //名称
	ItemType     string `gorm:"column:item_type;type:varchar(1000);not null;"`                   //规格
	ItemInfo     string `gorm:"column:item_info;type:text;"`                                     //详情
	ItemPicts    string `gorm:"column:item_picts;type:text;"`                                    //图片路径，多张以分号分割
	ItemLanguage string `gorm:"column:item_language;type:varchar(2);not null;primary_key:true;"` //参照m_language
}

func (t *TItemTrans) TableName() string {
	return "m_item_trans"
}

func (t *TItemTrans) GetData(db *gorm.DB, searchNo, langCode string) (bool, TItemTrans) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TItemTrans{}
	}

	var res []TItemTrans

	db.Where(`item_delflg = '0' AND item_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TItemTrans{}
	}
}

func (t *TItemTrans) CheckNo(db *gorm.DB, searchNo, langCode string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TItemTrans{}).
		Where(`item_delflg = '0' AND item_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
