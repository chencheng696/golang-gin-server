package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//招聘表
type TJob struct {
	JobNo          int64     `gorm:"column:job_no;AUTO_INCREMENT;primary_key:true;"`
	JobName        string    `gorm:"column:job_name;type:varchar(200);not null;"`    //招聘标题或对象
	JobHeadcount   int64     `gorm:"column:job_headcount;not null;"`                 //招聘人数
	JobAddress     string    `gorm:"column:job_address;type:varchar(200);not null;"` //工作地点
	JobTreatment   int64     `gorm:"column:job_treatment;not null;"`                 //待遇 0表示面议
	JobShowDate    time.Time `gorm:"column:job_showdate;"`                           //发布日期
	JobPeriod      int64     `gorm:"column:job_period;not null;"`                    //有限期限 单位年
	JobDescription string    `gorm:"column:job_description;type:text;not null;"`     //招聘要求描述
	JobInputdate   time.Time `gorm:"column:job_inputdate"`
	JobInputid     int64     `gorm:"column:job_inputid"`
	JobUpdate      time.Time `gorm:"column:job_update"`
	JobUpid        int64     `gorm:"column:job_upid"`
	JobDelflg      string    `gorm:"column:job_delflg;type:char(1);default:'0'"`

	Tbls
}

func (t *TJob) TableName() string {
	return "m_job"
}

//下载配置
func (t *TJob) DownloadConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:  `job_name`,
			`name`: `标题`,
		},
		map[string]string{
			`key`:  `job_headcount`,
			`name`: `人数`,
		},
		map[string]string{
			`key`:  `job_address`,
			`name`: `地点`,
		},
		map[string]string{
			`key`:  `job_treatment`,
			`name`: `待遇`,
		},
		map[string]string{
			`key`:        `job_showdate`,
			`name`:       `发布日期`,
			`formatdate`: `2006-01-02`,
		},
		map[string]string{
			`key`:  `job_period`,
			`name`: `有效期限`,
		},
	}

	return data
}

func (t *TJob) GetList(db *gorm.DB, arrData map[string]interface{}) []TJob {

	var res []TJob

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`job_no`).Find(&res)
	return res
}

func (t *TJob) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TJob{}).Select("*").Rows()
}

func (t *TJob) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TJob {

	var res []TJob

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TJob{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TJob) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`job_delflg = '0'`).Order(`job_no`)

	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`job_name LIKE ?`, "%"+seach+"%")
		}
	}
	return dbnew
}

func (t *TJob) GetData(db *gorm.DB, searchNo string) (bool, TJob) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TJob{}
	}

	var res []TJob

	db.Where(`job_delflg = '0' AND job_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TJob{}
	}
}

func (t *TJob) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TJob{}).
		Where(`job_delflg = '0' AND job_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
