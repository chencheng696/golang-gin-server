package tbls

import (
	"strconv"
	"time"

	"Yinghao/klib"

	"github.com/jinzhu/gorm"
)

//案例分类表
type TCaseClass struct {
	CaseClassNo        int64     `gorm:"column:case_class_no;AUTO_INCREMENT;primary_key:true;"`
	CaseClassParentNo  int64     `gorm:"column:case_class_parent_no;not null;"` //0表示无父类
	CaseClassName      string    `gorm:"column:case_class_name;type:varchar(200);not null;"`
	CaseClassMemo      string    `gorm:"column:case_class_memo;type:text;"`
	CaseClassPicts     string    `gorm:"column:case_class_picts;type:text;"` //图片路径，多张以分号分割
	CaseClassInputdate time.Time `gorm:"column:case_class_inputdate"`
	CaseClassInputid   int64     `gorm:"column:case_class_inputid"`
	CaseClassUpdate    time.Time `gorm:"column:case_class_update"`
	CaseClassUpid      int64     `gorm:"column:case_class_upid"`
	CaseClassDelflg    string    `gorm:"column:case_class_delflg;type:char(1);"`

	Tbls
}

func (t *TCaseClass) TableName() string {
	return "m_case_class"
}

func (t *TCaseClass) GetList(db *gorm.DB) []TCaseClass {

	var res []TCaseClass

	db.Where(`case_class_delflg = '0'`).Order(`case_class_no`).Find(&res)

	return res
}

func (t *TCaseClass) GetData(db *gorm.DB, searchNo string) (bool, TCaseClass) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TCaseClass{}
	}

	var res []TCaseClass

	db.Where(`case_class_delflg = '0' AND case_class_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TCaseClass{}
	}
}

func (t *TCaseClass) GetDataNative(db *gorm.DB, searchNo string) []map[string]interface{} {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return make([]map[string]interface{}, 0)
	}

	rows, err := db.Raw(`SELECT A.*, B.case_class_name as case_class_parent_name 
		FROM m_case_class A 
			LEFT JOIN m_case_class B ON A.case_class_parent_no = B.case_class_no
		WHERE A.case_class_delflg = '0' AND A.case_class_no = ?`, value).Rows()
	if err != nil {
		return make([]map[string]interface{}, 0)
	}

	return klib.SqlRows2Array(rows)
}

func (t *TCaseClass) CheckNo(db *gorm.DB, searchNo string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TCaseClass{}).
		Where(`case_class_delflg = '0' AND case_class_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TCaseClass) CheckChildren(db *gorm.DB, searchNo string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TCaseClass{}).
		Where(`case_class_delflg = '0' AND case_class_parent_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TCaseClass) GetTreeData(db *gorm.DB) map[int64]interface{} {

	res := t.GetList(db)

	relation := make(map[int64]string)
	for _, item := range res {
		relation[item.CaseClassNo] = item.CaseClassName
	}

	return t.createTreeData(0, res, relation)
}

func (t *TCaseClass) createTreeData(parentNo int64, src []TCaseClass, relation map[int64]string) map[int64]interface{} {

	tree := make(map[int64]interface{})

	for _, item := range src {
		if item.CaseClassParentNo == parentNo {
			parentName, _ := relation[item.CaseClassParentNo]

			tree[item.CaseClassNo] = map[string]interface{}{
				`caseclassno`:         item.CaseClassNo,
				`caseclassparentno`:   item.CaseClassParentNo,
				`caseclassname`:       item.CaseClassName,
				`caseclassparentname`: parentName,
				`children`:            t.createTreeData(item.CaseClassNo, src, relation),
			}
		}
	}
	return tree
}

func (t *TCaseClass) GetTreeShow(parentNo int64, src []TCaseClass) []interface{} {

	arr := make([]interface{}, 0)

	for _, item := range src {
		if item.CaseClassParentNo == parentNo {
			tree := map[string]interface{}{
				`id`:   item.CaseClassNo,
				`text`: item.CaseClassName,
			}
			nodes := t.GetTreeShow(item.CaseClassNo, src)
			if len(nodes) > 0 {
				tree[`nodes`] = nodes
			}

			arr = append(arr, tree)
		}
	}
	return arr
}
