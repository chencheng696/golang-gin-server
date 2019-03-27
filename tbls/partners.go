package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//合作伙伴表
type TPartners struct {
	PnNo          int64     `gorm:"column:pn_no;AUTO_INCREMENT;primary_key:true;"`
	PnName        string    `gorm:"column:pn_name;type:varchar(200);not null;"` //名称
	PnUrl         string    `gorm:"column:pn_url;type:varchar(200);"`           //网址
	PnStatus      string    `gorm:"column:pn_status;type:char(1);not null;"`    //1:正常 0:停用
	PnHit         int64     `gorm:"column:pn_hit;not null;"`                    //点击次数
	PnPicts       string    `gorm:"column:pn_picts;type:text;"`                 //图片路径，多张以分号分割
	PnDescription string    `gorm:"column:pn_description;type:text;not null;"`  //描述
	PnInputdate   time.Time `gorm:"column:pn_inputdate"`
	PnInputid     int64     `gorm:"column:pn_inputid"`
	PnUpdate      time.Time `gorm:"column:pn_update"`
	PnUpid        int64     `gorm:"column:pn_upid"`
	PnDelflg      string    `gorm:"column:pn_delflg;type:char(1);default:'0'"`

	Tbls
}

func (t *TPartners) TableName() string {
	return "m_partners"
}

//下载配置
func (t *TPartners) DownloadConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:  `pn_name`,
			`name`: `名称`,
		},
		map[string]string{
			`key`:  `pn_url`,
			`name`: `网址`,
		},
		map[string]string{
			`key`:  `pn_status`,
			`name`: `状态`,
		},
		map[string]string{
			`key`:  `pn_hit`,
			`name`: `点击次数`,
		},
		map[string]string{
			`key`:  `pn_description`,
			`name`: `描述`,
		},
	}

	return data
}

func (t *TPartners) GetList(db *gorm.DB, arrData map[string]interface{}) []TPartners {

	var res []TPartners

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`pn_no`).Find(&res)
	return res
}

func (t *TPartners) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TPartners{}).Select("*").Rows()
}

func (t *TPartners) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TPartners {

	var res []TPartners

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TPartners{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TPartners) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`pn_delflg = '0'`).Order(`pn_no`)

	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`pn_name LIKE ?`, "%"+seach+"%")
		}
	}
	return dbnew
}

func (t *TPartners) GetData(db *gorm.DB, searchNo string) (bool, TPartners) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TPartners{}
	}

	var res []TPartners

	db.Where(`pn_delflg = '0' AND pn_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TPartners{}
	}
}

func (t *TPartners) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TPartners{}).
		Where(`pn_delflg = '0' AND pn_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
