package tbls

import (
	"strconv"
	"time"

	"Yinghao/klib"

	"github.com/jinzhu/gorm"
)

//商品分类表
type TItemClass struct {
	ItemClassNo        int64     `gorm:"column:item_class_no;AUTO_INCREMENT;primary_key:true;"`
	ItemClassParentNo  int64     `gorm:"column:item_class_parent_no;not null;"` //0表示无父类
	ItemClassName      string    `gorm:"column:item_class_name;type:varchar(200);not null;"`
	ItemClassMemo      string    `gorm:"column:item_class_memo;type:text;"`
	ItemClassPicts     string    `gorm:"column:item_class_picts;type:text;"` //图片路径，多张以分号分割
	ItemClassInputdate time.Time `gorm:"column:item_class_inputdate"`
	ItemClassInputid   int64     `gorm:"column:item_class_inputid"`
	ItemClassUpdate    time.Time `gorm:"column:item_class_update"`
	ItemClassUpid      int64     `gorm:"column:item_class_upid"`
	ItemClassDelflg    string    `gorm:"column:item_class_delflg;type:char(1);default:'0'"`

	Tbls
}

func (t *TItemClass) TableName() string {
	return "m_item_class"
}

func (t *TItemClass) GetList(db *gorm.DB) []TItemClass {

	var res []TItemClass

	db.Where(`item_class_delflg = '0'`).Order(`item_class_no`).Find(&res)

	return res
}

func (t *TItemClass) GetData(db *gorm.DB, searchNo string) (bool, TItemClass) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TItemClass{}
	}

	var res []TItemClass

	db.Where(`item_class_delflg = '0' AND item_class_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TItemClass{}
	}
}

func (t *TItemClass) GetDataNative(db *gorm.DB, searchNo string) []map[string]interface{} {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return make([]map[string]interface{}, 0)
	}

	rows, err := db.Raw(`SELECT A.*, B.item_class_name as item_class_parent_name 
		FROM m_item_class A 
			LEFT JOIN m_item_class B ON A.item_class_parent_no = B.item_class_no
		WHERE A.item_class_delflg = '0' AND A.item_class_no = ?`, value).Rows()
	if err != nil {
		return make([]map[string]interface{}, 0)
	}

	return klib.SqlRows2Array(rows)
}

func (t *TItemClass) CheckNo(db *gorm.DB, searchNo string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TItemClass{}).
		Where(`item_class_delflg = '0' AND item_class_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TItemClass) CheckChildren(db *gorm.DB, searchNo string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TItemClass{}).
		Where(`item_class_delflg = '0' AND item_class_parent_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TItemClass) GetTreeData(db *gorm.DB) map[int64]interface{} {

	res := t.GetList(db)

	relation := make(map[int64]string)
	for _, item := range res {
		relation[item.ItemClassNo] = item.ItemClassName
	}

	return t.createTreeData(0, res, relation)
}

func (t *TItemClass) createTreeData(parentNo int64, src []TItemClass, relation map[int64]string) map[int64]interface{} {

	tree := make(map[int64]interface{})

	for _, item := range src {
		if item.ItemClassParentNo == parentNo {
			parentName, _ := relation[item.ItemClassParentNo]

			tree[item.ItemClassNo] = map[string]interface{}{
				`itemclassno`:         item.ItemClassNo,
				`itemclassparentno`:   item.ItemClassParentNo,
				`itemclassname`:       item.ItemClassName,
				`itemclassparentname`: parentName,
				`children`:            t.createTreeData(item.ItemClassNo, src, relation),
			}
		}
	}
	return tree
}

func (t *TItemClass) GetTreeShow(parentNo int64, src []TItemClass) []interface{} {

	arr := make([]interface{}, 0)

	for _, item := range src {
		if item.ItemClassParentNo == parentNo {
			tree := map[string]interface{}{
				`id`:   item.ItemClassNo,
				`text`: item.ItemClassName,
			}
			nodes := t.GetTreeShow(item.ItemClassNo, src)
			if len(nodes) > 0 {
				tree[`nodes`] = nodes
			}

			arr = append(arr, tree)
		}
	}
	return arr
}
