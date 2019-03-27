package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//案例表
type TCase struct {
	CaseNo          int64     `gorm:"column:case_no;AUTO_INCREMENT;primary_key:true;"`
	CaseName        string    `gorm:"column:case_name;type:varchar(1000);not null;"` //标题
	CaseClassNo     int64     `gorm:"column:case_class_no;"`                         //分类
	CaseStatus      string    `gorm:"column:case_status;type:char(1);not null;"`     //1:发布 0:草稿
	CaseInfo        string    `gorm:"column:case_info;type:text;"`                   //详情
	CaseHit         int64     `gorm:"column:case_hit;not null;"`                     //点击次数
	CasePicts       string    `gorm:"column:case_picts;type:text;"`                  //图片路径，多张以分号分割
	CaseAuthor      string    `gorm:"column:case_author;type:varchar(32);"`          //作者
	CaseShowDate    time.Time `gorm:"column:case_showdate"`                          //文章做成日期
	CasePublishDate time.Time `gorm:"column:case_publishdate"`                       //文章前端开始显示日期
	CaseInputdate   time.Time `gorm:"column:case_inputdate"`
	CaseInputid     int64     `gorm:"column:case_inputid"`
	CaseUpdate      time.Time `gorm:"column:case_update"`
	CaseUpid        int64     `gorm:"column:case_upid"`
	CaseDelflg      string    `gorm:"column:case_delflg;type:char(1);"`

	Tbls
}

func (t *TCase) TableName() string {
	return "m_case"
}

func (t *TCase) GetList(db *gorm.DB, arrData map[string]interface{}) []TCase {

	var res []TCase

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`case_no`).Find(&res)
	return res
}

func (t *TCase) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TCase{}).Select("*").Rows()
}

func (t *TCase) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TCase {

	var res []TCase

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TCase{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TCase) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`case_delflg = '0'`).Order(`case_no`)

	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`case_name LIKE ?`, "%"+seach+"%")
		}
	}
	if value, ok := arrData[`searchClassNo`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`case_class_no = ?`, seach)
		}
	}
	if value, ok := arrData[`searchStatus`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`case_status = ?`, seach)
		}
	}
	return dbnew
}

func (t *TCase) GetData(db *gorm.DB, searchNo string) (bool, TCase) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TCase{}
	}

	var res []TCase

	db.Where(`case_delflg = '0' AND case_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TCase{}
	}
}

func (t *TCase) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TCase{}).
		Where(`case_delflg = '0' AND case_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TCase) CheckName(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TCase{}).
		Where(`case_delflg = '0' AND case_name = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
