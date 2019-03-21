package main

import (
	"os"
	"strconv"
	"time"

	"Yinghao/klib"
	"Yinghao/tbls"

	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//检查配置文件
func initConfig() bool {

	cfg, _ := ini.Load(klib.GetAppPath() + "/app.ini")

	appCfg.title = cfg.Section("main").Key("title").String()
	if appCfg.title == "" {
		klib.WriteLog("[ERROR] 配置文件title是空")
		return false
	}
	klib.WriteLog(`title:` + appCfg.title)

	appCfg.version = cfg.Section("main").Key("version").String()
	if appCfg.version == "" {
		klib.WriteLog("[ERROR] 配置文件version是空")
		return false
	}
	klib.WriteLog(`version:` + appCfg.version)

	appCfg.webport = cfg.Section("main").Key("webport").String()
	if appCfg.webport == "" {
		appCfg.webport = "80"
	}
	klib.WriteLog(`webport:` + appCfg.webport)

	appCfg.dbhost = cfg.Section("main").Key("dbhost").String()
	if appCfg.dbhost == "" {
		klib.WriteLog("[ERROR] 配置文件dbhost是空")
		return false
	}
	klib.WriteLog(`dbhost:` + appCfg.dbhost)

	appCfg.dbuser = cfg.Section("main").Key("dbuser").String()
	if appCfg.dbuser == "" {
		klib.WriteLog("[ERROR] 配置文件dbuser是空")
		return false
	}
	klib.WriteLog(`dbuser:` + appCfg.dbuser)

	appCfg.dbpassword = cfg.Section("main").Key("dbpassword").String()
	klib.WriteLog(`dbpassword:` + appCfg.dbpassword)

	appCfg.dbname = cfg.Section("main").Key("dbname").String()
	if appCfg.dbname == "" {
		klib.WriteLog("[ERROR] 配置文件dbname是空")
		return false
	}
	klib.WriteLog(`dbname:` + appCfg.dbname)

	var err error
	appCfg.pagerow, err = cfg.Section("main").Key("pagerow").Int()
	if err != nil {
		klib.WriteLog("[WARNING] 配置文件pagerow数据不正确")
		appCfg.pagerow = 10
	}
	klib.WriteLog(`pagerow:` + strconv.Itoa(appCfg.pagerow))

	value := cfg.Section("main").Key("debug").String()
	if value == "1" {
		appCfg.debug = true
		klib.WriteLog(`debug: true`)
	} else {
		appCfg.debug = false
		klib.WriteLog(`debug: false`)
	}

	value = cfg.Section("upload").Key("maxsize").String()
	appCfg.uploadMaxSize, err = strconv.Atoi(value)
	if err != nil {
		appCfg.uploadMaxSize = 10 //10M
		klib.WriteLog(`upload-maxsize: error, default 10`)
	} else {
		klib.WriteLog(`upload-maxsize: ` + strconv.Itoa(appCfg.uploadMaxSize) + `M`)
	}

	if ok, _ := klib.PathExists(klib.GetAppPath() + `/upload`); ok == false {
		os.Mkdir(klib.GetAppPath()+`/upload`, os.ModePerm)
	}
	appCfg.uploadTemp = cfg.Section("upload").Key("temp").String()
	if appCfg.uploadTemp == "" {
		appCfg.uploadTemp = klib.GetAppPath() + `/upload/temp/`
	} else {
		appCfg.uploadTemp = klib.GetAppPath() + `/` + appCfg.uploadTemp
	}
	klib.WriteLog(`upload-temp:` + appCfg.uploadTemp)
	if ok, _ := klib.PathExists(appCfg.uploadTemp); ok == false {
		os.Mkdir(appCfg.uploadTemp, os.ModePerm)
	}
	appCfg.uploadPict = cfg.Section("upload").Key("pict").String()
	if appCfg.uploadPict == "" {
		appCfg.uploadPict = klib.GetAppPath() + `/upload/pict/`
	} else {
		appCfg.uploadPict = klib.GetAppPath() + `/` + appCfg.uploadPict
	}
	klib.WriteLog(`upload-pict:` + appCfg.uploadPict)
	if ok, _ := klib.PathExists(appCfg.uploadPict); ok == false {
		os.Mkdir(appCfg.uploadPict, os.ModePerm)
	}

	if ok, _ := klib.PathExists(klib.GetAppPath() + `/download/`); ok == false {
		os.Mkdir(`download/`, os.ModePerm)
	}
	appCfg.downloadExcel = cfg.Section("download").Key("excel").String()
	if appCfg.downloadExcel == "" {
		appCfg.downloadExcel = klib.GetAppPath() + `/download/excel/`
	} else {
		appCfg.downloadExcel = klib.GetAppPath() + `/` + appCfg.downloadExcel
	}
	klib.WriteLog(`download-excel:` + appCfg.downloadExcel)
	if ok, _ := klib.PathExists(appCfg.downloadExcel); ok == false {
		os.Mkdir(appCfg.downloadExcel, os.ModePerm)
	}

	return true
}

func initDatabase() bool {

	var err error

	//db, err = gorm.Open("postgres", "host="+appCfg.dbhost+" user="+appCfg.dbuser+" password="+appCfg.dbpassword+" dbname="+appCfg.dbname+" sslmode=disable port=5432")
	db, err = gorm.Open("mysql", appCfg.dbuser+":"+appCfg.dbpassword+"@/"+appCfg.dbname+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		klib.WriteLog("[ERROR] 数据库连接" + err.Error())
		return false
	}

	db.AutoMigrate(&tbls.TAdmin{})
	db.AutoMigrate(&tbls.TItem{})
	db.AutoMigrate(&tbls.TItemTrans{})
	db.AutoMigrate(&tbls.TItemClass{})
	db.AutoMigrate(&tbls.TItemClassTrans{})
	db.AutoMigrate(&tbls.TNewsClass{})
	db.AutoMigrate(&tbls.TNewsClassTrans{})
	db.AutoMigrate(&tbls.TCaseClass{})
	db.AutoMigrate(&tbls.TCaseClassTrans{})
	db.AutoMigrate(&tbls.TJob{})
	db.AutoMigrate(&tbls.TJobTrans{})
	db.AutoMigrate(&tbls.TFeedback{})
	db.AutoMigrate(&tbls.TFeedbackTrans{})
	db.AutoMigrate(&tbls.TPartners{})
	db.AutoMigrate(&tbls.TPartnersTrans{})
	db.AutoMigrate(&tbls.TPromotion{})
	db.AutoMigrate(&tbls.TPromotionTrans{})
	db.AutoMigrate(&tbls.TPermission{})
	db.AutoMigrate(&tbls.TLanguage{})
	db.AutoMigrate(&tbls.TSystem{})
	db.AutoMigrate(&tbls.TSystemTrans{})

	//插入默认值
	//管理员
	admin := new(tbls.TAdmin)
	if !admin.CheckId(db, `admin`) {
		admin.AdmId = `admin`
		admin.AdmPass = klib.MD5ForStr(`admin`)
		admin.AdmName = `超级管理员`
		admin.AdmPerm = `0`
		admin.AdmInputdate = time.Now()
		admin.AdmInputid = 0
		admin.AdmDelflg = `0`
		db.Create(&admin)
	}

	//语言master
	lang := new(tbls.TLanguage)
	res := lang.GetList(db)
	if len(res) == 0 {
		var arr []tbls.TLanguage
		arr = append(arr, tbls.TLanguage{
			LngCode:     `cn`,
			LngName:     `简体中文`,
			LngShowName: `简体中文`,
		})
		arr = append(arr, tbls.TLanguage{
			LngCode:     `tw`,
			LngName:     `繁体中文`,
			LngShowName: `繁体中文`,
		})
		arr = append(arr, tbls.TLanguage{
			LngCode:     `en`,
			LngName:     `English`,
			LngShowName: `English`,
		})
		arr = append(arr, tbls.TLanguage{
			LngCode:     `jp`,
			LngName:     `日语`,
			LngShowName: `日本語`,
		})
		arr = append(arr, tbls.TLanguage{
			LngCode:     `kr`,
			LngName:     `韩语`,
			LngShowName: `한국어`,
		})

		for i := 0; i < len(arr); i++ {
			arr[i].LngSeq = int64(i)
			arr[i].LngInputdate = time.Now()
			arr[i].LngDelflg = `0`
			db.Create(&arr[i])
		}
	}

	return true
}
