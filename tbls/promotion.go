package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//网络推广表
type TPromotion struct {
	PromNo          int64     `gorm:"column:prom_no;AUTO_INCREMENT;primary_key:true;"`
	PromName        string    `gorm:"column:prom_name;type:varchar(200);not null;"` //名称
	PromUrl         string    `gorm:"column:prom_url;type:varchar(200);"`           //网址
	PromStatus      string    `gorm:"column:prom_status;type:char(1);not null;"`    //1:正常 0:停用
	PromHit         int64     `gorm:"column:prom_hit;not null;"`                    //点击次数
	PromPicts       string    `gorm:"column:prom_picts;type:text;"`                 //图片路径，多张以分号分割
	PromDescription string    `gorm:"column:prom_description;type:text;not null;"`  //描述
	PromInputdate   time.Time `gorm:"column:prom_inputdate"`
	PromInputid     int64     `gorm:"column:prom_inputid"`
	PromUpdate      time.Time `gorm:"column:prom_update"`
	PromUpid        int64     `gorm:"column:prom_upid"`
	PromDelflg      string    `gorm:"column:prom_delflg;type:char(1);default:'0'"`

	Tbls
}

func (t *TPromotion) TableName() string {
	return "m_promotion"
}

//下载配置
func (t *TPromotion) DownloadConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:  `prom_name`,
			`name`: `标题`,
		},
		map[string]string{
			`key`:  `prom_headcount`,
			`name`: `人数`,
		},
		map[string]string{
			`key`:  `prom_address`,
			`name`: `地点`,
		},
		map[string]string{
			`key`:  `prom_treatment`,
			`name`: `待遇`,
		},
		map[string]string{
			`key`:        `prom_showdate`,
			`name`:       `发布日期`,
			`formatdate`: `2006-01-02`,
		},
		map[string]string{
			`key`:  `prom_period`,
			`name`: `有效期限`,
		},
	}

	return data
}

func (t *TPromotion) GetList(db *gorm.DB, arrData map[string]interface{}) []TPromotion {

	var res []TPromotion

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`prom_no`).Find(&res)
	return res
}

func (t *TPromotion) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TPromotion{}).Select("*").Rows()
}

func (t *TPromotion) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TPromotion {

	var res []TPromotion

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TPromotion{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TPromotion) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`prom_delflg = '0'`).Order(`prom_no`)

	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`prom_name LIKE ?`, "%"+seach+"%")
		}
	}
	return dbnew
}

func (t *TPromotion) GetData(db *gorm.DB, searchNo string) (bool, TPromotion) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TPromotion{}
	}

	var res []TPromotion

	db.Where(`prom_delflg = '0' AND prom_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TPromotion{}
	}
}

func (t *TPromotion) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TPromotion{}).
		Where(`prom_delflg = '0' AND prom_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
