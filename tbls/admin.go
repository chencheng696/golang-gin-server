package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"Yinghao/klib"

	"github.com/jinzhu/gorm"
)

//管理员表
type TAdmin struct {
	AdmNo        int64     `gorm:"column:adm_no;AUTO_INCREMENT;primary_key:true;"` //自增长 key
	AdmId        string    `gorm:"column:adm_id;varchar(200);not null;"`           //需要判断唯一
	AdmPass      string    `gorm:"column:adm_pass;type:varchar(32);not null;"`     //密码以md5形式存储
	AdmName      string    `gorm:"column:adm_name;type:varchar(16);"`
	AdmPerm      string    `gorm:"column:adm_perm;type:char(1);not null;"` //0:超级管理员 1:普通
	AdmMemo      string    `gorm:"column:adm_memo;type:varchar(200);"`
	AdmInputdate time.Time `gorm:"column:adm_inputdate"`
	AdmInputid   int64     `gorm:"column:adm_inputid"`
	AdmUpdate    time.Time `gorm:"column:adm_update"`
	AdmUpid      int64     `gorm:"column:adm_upid"`
	AdmDelflg    string    `gorm:"column:adm_delflg;type:char(1);"`

	Tbls
}

func (t *TAdmin) TableName() string {
	return "m_admin"
}

//下载配置
func (t *TAdmin) DownloadConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:  `adm_id`,
			`name`: `管理员ID`,
		},
		map[string]string{
			`key`:  `adm_name`,
			`name`: `管理员名字`,
		},
		map[string]string{
			`key`:     `adm_perm`,
			`name`:    `权限`,
			`convert`: `ConvertPerm`, //函数
		},
	}

	return data
}

//导入配置
func (t *TAdmin) ImportConfig() []map[string]string {

	data := []map[string]string{
		map[string]string{
			`key`:    `adm_id`,
			`name`:   `管理员ID`,
			`type`:   `ss`,
			`maxlen`: `20`,
		},
		map[string]string{
			`key`:    `adm_name`,
			`name`:   `管理员名字`,
			`type`:   `ss`,
			`maxlen`: `16`,
		},
		map[string]string{
			`key`:    `adm_pwd`,
			`name`:   `密码`,
			`type`:   `pass`,
			`minlen`: `6`,
			`maxlen`: `16`,
		},
		map[string]string{
			`key`:    `adm_perm`,
			`name`:   `权限`,
			`type`:   `ss`,
			`fixlen`: `1`,
		},
	}

	return data
}

func (t *TAdmin) Login(db *gorm.DB, id string, pwd string) (bool, TAdmin) {

	var res []TAdmin
	db.Where(`adm_delflg = '0' AND `+
		`adm_id = ? AND `+
		`adm_pass = ?`, id, klib.MD5ForStr(pwd)).Find(&res)

	if len(res) != 1 {
		return false, TAdmin{}
	} else {
		return true, res[0]
	}
}

func (t *TAdmin) GetList(db *gorm.DB, arrData map[string]interface{}) []TAdmin {

	var res []TAdmin

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`adm_no`).Find(&res)
	return res
}

func (t *TAdmin) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TAdmin{}).Select("*").Rows()
}

func (t *TAdmin) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TAdmin {

	var res []TAdmin

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TAdmin{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TAdmin) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`adm_delflg = '0'`).Order(`adm_no`)

	if value, ok := arrData[`searchId`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`adm_id LIKE ?`, "%"+seach+"%")
		}
	}
	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`adm_name LIKE ?`, "%"+seach+"%")
		}
	}
	return dbnew
}

func (t *TAdmin) GetData(db *gorm.DB, searchNo string) (bool, TAdmin) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TAdmin{}
	}

	var res []TAdmin

	db.Where(`adm_delflg = '0' AND adm_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TAdmin{}
	}
}

func (t *TAdmin) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TAdmin{}).
		Where(`adm_delflg = '0' AND adm_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TAdmin) CheckId(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TAdmin{}).
		Where(`adm_delflg = '0' AND adm_id = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
