package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//商品表
type TItem struct {
	ItemNo        int64     `gorm:"column:item_no;AUTO_INCREMENT;primary_key:true;"`
	ItemName      string    `gorm:"column:item_name;type:varchar(1000);not null;"` //名称
	ItemType      string    `gorm:"column:item_type;type:varchar(1000);not null;"` //规格
	ItemClassNo   int64     `gorm:"column:item_class_no;"`                         //对于商品分类表No
	ItemStatus    string    `gorm:"column:item_status;type:char(1);not null;"`     //1:正常 0:下架
	ItemInfo      string    `gorm:"column:item_info;type:text;"`                   //详情
	ItemHit       int64     `gorm:"column:item_hit;not null;"`                     //点击次数
	ItemPicts     string    `gorm:"column:item_picts;type:text;"`                  //图片路径，多张以分号分割
	ItemInputdate time.Time `gorm:"column:item_inputdate"`
	ItemInputid   int64     `gorm:"column:item_inputid"`
	ItemUpdate    time.Time `gorm:"column:item_update"`
	ItemUpid      int64     `gorm:"column:item_upid"`
	ItemDelflg    string    `gorm:"column:item_delflg;type:char(1);"`

	Tbls
}

func (t *TItem) TableName() string {
	return "m_item"
}

//下载配置
func (t *TItem) DownloadConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:  `item_name`,
			`name`: `商品名称`,
		},
		map[string]string{
			`key`:  `item_type`,
			`name`: `商品规格`,
		},
		map[string]string{
			`key`:  `item_class_no`,
			`name`: `商品分类`,
		},
		map[string]string{
			`key`:     `item_status`,
			`name`:    `状态`,
			`convert`: `ConvertItemStatus`, //函数
		},
		map[string]string{
			`key`:  `item_info`,
			`name`: `商品详情`,
		},
		map[string]string{
			`key`:  `item_hit`,
			`name`: `点击次数`,
		},
	}

	return data
}

//导入配置
func (t *TItem) ImportConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:    `item_name`,
			`name`:   `商品名称`,
			`type`:   `ss`,
			`maxlen`: `1000`,
		},
		map[string]string{
			`key`:    `item_type`,
			`name`:   `商品规格`,
			`type`:   `:ss`,
			`maxlen`: `1000`,
		},
		map[string]string{
			`key`:    `item_class_no`,
			`name`:   `商品类别`,
			`type`:   `d`,
			`maxlen`: `10`,
		},
		map[string]string{
			`key`:    `item_status`,
			`name`:   `商品状态`,
			`type`:   `ss`,
			`fixlen`: `1`,
		},
		map[string]string{
			`key`:    `item_info`,
			`name`:   `商品详情`,
			`type`:   `:t`,
			`maxlen`: `200`,
		},
	}

	return data
}

func (t *TItem) GetList(db *gorm.DB, arrData map[string]interface{}) []TItem {

	var res []TItem

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`item_no`).Find(&res)
	return res
}

func (t *TItem) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TItem{}).Select("*").Rows()
}

func (t *TItem) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TItem {

	var res []TItem

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TItem{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TItem) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`item_delflg = '0'`).Order(`item_no`)

	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`item_name LIKE ?`, "%"+seach+"%")
		}
	}
	if value, ok := arrData[`searchClassNo`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`item_class_no = ?`, seach)
		}
	}
	if value, ok := arrData[`searchStatus`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`item_status = ?`, seach)
		}
	}
	return dbnew
}

func (t *TItem) GetData(db *gorm.DB, searchNo string) (bool, TItem) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TItem{}
	}

	var res []TItem

	db.Where(`item_delflg = '0' AND item_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TItem{}
	}
}

func (t *TItem) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TItem{}).
		Where(`item_delflg = '0' AND item_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TItem) CheckName(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TItem{}).
		Where(`item_delflg = '0' AND item_name = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
