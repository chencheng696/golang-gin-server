package tbls

import (
	"strconv"
	"time"

	"Yinghao/klib"

	"github.com/jinzhu/gorm"
)

//新闻分类表
type TNewsClass struct {
	NewsClassNo        int64     `gorm:"column:news_class_no;AUTO_INCREMENT;primary_key:true;"`
	NewsClassParentNo  int64     `gorm:"column:news_class_parent_no;not null;"` //0表示无父类
	NewsClassName      string    `gorm:"column:news_class_name;type:varchar(200);not null;"`
	NewsClassMemo      string    `gorm:"column:news_class_memo;type:text;"`
	NewsClassPicts     string    `gorm:"column:news_class_picts;type:text;"` //图片路径，多张以分号分割
	NewsClassInputdate time.Time `gorm:"column:news_class_inputdate"`
	NewsClassInputid   int64     `gorm:"column:news_class_inputid"`
	NewsClassUpdate    time.Time `gorm:"column:news_class_update"`
	NewsClassUpid      int64     `gorm:"column:news_class_upid"`
	NewsClassDelflg    string    `gorm:"column:news_class_delflg;type:char(1);"`

	Tbls
}

func (t *TNewsClass) TableName() string {
	return "m_news_class"
}

func (t *TNewsClass) GetList(db *gorm.DB) []TNewsClass {

	var res []TNewsClass

	db.Where(`news_class_delflg = '0'`).Order(`news_class_no`).Find(&res)

	return res
}

func (t *TNewsClass) GetData(db *gorm.DB, searchNo string) (bool, TNewsClass) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TNewsClass{}
	}

	var res []TNewsClass

	db.Where(`news_class_delflg = '0' AND news_class_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TNewsClass{}
	}
}

func (t *TNewsClass) GetDataNative(db *gorm.DB, searchNo string) []map[string]interface{} {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return make([]map[string]interface{}, 0)
	}

	rows, err := db.Raw(`SELECT A.*, B.news_class_name as news_class_parent_name 
		FROM m_news_class A 
			LEFT JOIN m_news_class B ON A.news_class_parent_no = B.news_class_no
		WHERE A.news_class_delflg = '0' AND A.news_class_no = ?`, value).Rows()
	if err != nil {
		return make([]map[string]interface{}, 0)
	}

	return klib.SqlRows2Array(rows)
}

func (t *TNewsClass) CheckNo(db *gorm.DB, searchNo string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TNewsClass{}).
		Where(`news_class_delflg = '0' AND news_class_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TNewsClass) CheckChildren(db *gorm.DB, searchNo string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TNewsClass{}).
		Where(`news_class_delflg = '0' AND news_class_parent_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TNewsClass) GetTreeData(db *gorm.DB) map[int64]interface{} {

	res := t.GetList(db)

	relation := make(map[int64]string)
	for _, item := range res {
		relation[item.NewsClassNo] = item.NewsClassName
	}

	return t.createTreeData(0, res, relation)
}

func (t *TNewsClass) createTreeData(parentNo int64, src []TNewsClass, relation map[int64]string) map[int64]interface{} {

	tree := make(map[int64]interface{})

	for _, item := range src {
		if item.NewsClassParentNo == parentNo {
			parentName, _ := relation[item.NewsClassParentNo]

			tree[item.NewsClassNo] = map[string]interface{}{
				`newsclassno`:         item.NewsClassNo,
				`newsclassparentno`:   item.NewsClassParentNo,
				`newsclassname`:       item.NewsClassName,
				`newsclassparentname`: parentName,
				`children`:            t.createTreeData(item.NewsClassNo, src, relation),
			}
		}
	}
	return tree
}

func (t *TNewsClass) GetTreeShow(parentNo int64, src []TNewsClass) []interface{} {

	arr := make([]interface{}, 0)

	for _, item := range src {
		if item.NewsClassParentNo == parentNo {
			tree := map[string]interface{}{
				`id`:   item.NewsClassNo,
				`text`: item.NewsClassName,
			}
			nodes := t.GetTreeShow(item.NewsClassNo, src)
			if len(nodes) > 0 {
				tree[`nodes`] = nodes
			}

			arr = append(arr, tree)
		}
	}
	return arr
}
