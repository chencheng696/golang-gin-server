package tbls

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//新闻表
type TNews struct {
	NewsNo          int64     `gorm:"column:news_no;AUTO_INCREMENT;primary_key:true;"`
	NewsName        string    `gorm:"column:news_name;type:varchar(1000);not null;"` //标题
	NewsClassNo     int64     `gorm:"column:news_class_no;"`                         //分类
	NewsStatus      string    `gorm:"column:news_status;type:char(1);not null;"`     //1:发布 0:草稿
	NewsInfo        string    `gorm:"column:news_info;type:text;"`                   //详情
	NewsPicts       string    `gorm:"column:news_picts;type:text;"`                  //图片路径，多张以分号分割
	NewsHit         int64     `gorm:"column:news_hit;not null;"`                     //点击次数
	NewsAuthor      string    `gorm:"column:news_author;type:varchar(32);"`          //作者
	NewsShowDate    time.Time `gorm:"column:news_showdate"`                          //文章做成日期
	NewsPublishDate time.Time `gorm:"column:news_publishdate"`                       //文章前端开始显示日期
	NewsInputdate   time.Time `gorm:"column:news_inputdate"`
	NewsInputid     int64     `gorm:"column:news_inputid"`
	NewsUpdate      time.Time `gorm:"column:news_update"`
	NewsUpid        int64     `gorm:"column:news_upid"`
	NewsDelflg      string    `gorm:"column:news_delflg;type:char(1);"`

	Tbls
}

func (t *TNews) TableName() string {
	return "m_news"
}

func (t *TNews) GetList(db *gorm.DB, arrData map[string]interface{}) []TNews {

	var res []TNews

	dbnew := t.GetListWhere(db, arrData)
	dbnew.Order(`news_no`).Find(&res)
	return res
}

func (t *TNews) GetListNative(db *gorm.DB, arrData map[string]interface{}) (*sql.Rows, error) {
	dbnew := t.GetListWhere(db, arrData)
	return dbnew.Model(&TNews{}).Select("*").Rows()
}

func (t *TNews) GetListPage(db *gorm.DB, arrData map[string]interface{}, pageRow int) []TNews {

	var res []TNews

	dbnew := t.GetListWhere(db, arrData)
	//获取总行数
	dbnew.Model(&TNews{}).Count(&t.RowCount)

	pageNo := `1`
	if value, ok := arrData[`pageNo`]; ok {
		pageNo = value.(string)
	}
	t.Calclulate(pageNo, pageRow)

	dbnew.Offset((t.PageNo - 1) * pageRow).Limit(pageRow).Find(&res)

	return res
}

func (t *TNews) GetListWhere(db *gorm.DB, arrData map[string]interface{}) *gorm.DB {
	dbnew := db.Where(`news_delflg = '0'`).Order(`news_no`)

	if value, ok := arrData[`searchName`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`news_name LIKE ?`, "%"+seach+"%")
		}
	}
	if value, ok := arrData[`searchClassNo`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`news_class_no = ?`, seach)
		}
	}
	if value, ok := arrData[`searchStatus`]; ok {
		seach := value.(string)
		if len(seach) > 0 {
			dbnew = dbnew.Where(`news_status = ?`, seach)
		}
	}
	return dbnew
}

func (t *TNews) GetData(db *gorm.DB, searchNo string) (bool, TNews) {

	value, err := strconv.ParseInt(searchNo, 10, 64)
	if err != nil {
		return false, TNews{}
	}

	var res []TNews

	db.Where(`news_delflg = '0' AND news_no = ?`, value).Find(&res)

	if len(res) > 0 {
		return true, res[0]
	} else {
		return false, TNews{}
	}
}

func (t *TNews) CheckNo(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TNews{}).
		Where(`news_delflg = '0' AND news_no = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (t *TNews) CheckName(db *gorm.DB, v string) bool {

	count := 0

	db.Model(&TNews{}).
		Where(`news_delflg = '0' AND news_name = ?`, v).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}
