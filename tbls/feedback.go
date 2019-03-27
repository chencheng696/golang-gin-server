package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//访客留言表
type TFeedback struct {
	FbNo        int64     `gorm:"column:fb_no;AUTO_INCREMENT;primary_key:true;"`
	FbName      string    `gorm:"column:fb_name;type:varchar(200);not null;"` //姓名
	FbCompany   string    `gorm:"column:fb_company;type:varchar(200);"`       //公司
	FbAddress   string    `gorm:"column:fb_address;type:varchar(200);"`       //地址
	FbSex       string    `gorm:"column:fb_sex;type:char(1);"`                //性别 0男 1女
	FbTel       string    `gorm:"column:fb_tel;type:varchar(20);"`            //电话
	FbPhone     string    `gorm:"column:fb_phone;varchar(20);"`               //手机
	FbFax       string    `gorm:"column:fb_fax;varchar(20);"`                 //传真
	FbEmail     string    `gorm:"column:fb_email;type:varchar(200);"`         //邮箱
	FbUrl       string    `gorm:"column:fb_url;type:varchar(200);"`           //网址
	FbTitle     string    `gorm:"column:fb_title;type:varchar(200);"`         //标题
	FbContent   string    `gorm:"column:fb_content;type:text;"`               //内容
	FbInputdate time.Time `gorm:"column:fb_inputdate"`
	FbInputid   int64     `gorm:"column:fb_inputid"`
	FbUpdate    time.Time `gorm:"column:fb_update"`
	FbUpid      int64     `gorm:"column:fb_upid"`
	FbDelflg    string    `gorm:"column:fb_delflg;type:char(1);default:'0'"`

	Tbls
}

func (t *TFeedback) TableName() string {
	return "t_feedback"
}

//下载配置
func (t *TFeedback) DownloadConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:  `fb_name`,
			`name`: `姓名`,
		},
		map[string]string{
			`key`:  `fb_email`,
			`name`: `邮箱`,
		},
		map[string]string{
			`key`:  `fb_title`,
			`name`: `标题`,
		},
		map[string]string{
			`key`:  `fb_content`,
			`name`: `内容`,
		},
		map[string]string{
			`key`:        `fb_inputdate`,
			`name`:       `日期`,
			`formatdate`: `2006-01-02`,
		},
	}

	return data
}

func (t *TFeedback) GetList(db *gorm.DB, arrData map[string]interface{}) []TFeedback {

	var res []TFeedback

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`fb_no`).Find(&res)
	return res
}

func (t *TFeedback) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TFeedback{}).Select("*").Rows()
}

func (t *TFeedback) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TFeedback {

	var res []TFeedback

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TFeedback{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TFeedback) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`fb_delflg = '0'`).Order(`fb_no`)

	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`fb_name LIKE ?`, "%"+seach+"%")
		}
	}
	return dbnew
}

func (t *TFeedback) GetData(db *gorm.DB, searchNo string) (bool, TFeedback) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TFeedback{}
	}

	var res []TFeedback

	db.Where(`fb_delflg = '0' AND fb_no = ?`, value).Find(&res)
	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TFeedback{}
	}
}

func (t *TFeedback) CheckNo(db *gorm.DB, searchNo string) bool {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false
	}

	count := 0

	db.Model(&TFeedback{}).
		Where(`fb_delflg = '0' AND fb_no = ?`, value).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
