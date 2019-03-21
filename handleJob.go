/*
注意：前台固定提交的key是小写不带下划线，例如checkMap[`jobno`]
*/

package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"Yinghao/klib"
	"Yinghao/tbls"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func jobHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `job`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := jobValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `list_del` {
		jobDelete(c, e, arrData)
		return
	} else if cmd == `ajax_add` {
		jobAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		jobAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		jobAjaxDetail(c, e, arrData)
		return
	} else if cmd == `ajax_detail_lang` {
		jobAjaxDetailLang(c, e, arrData)
		return
	} else if cmd == `ajax_edit_lang` {
		jobAjaxEditLang(c, e, arrData)
		return
	} else if cmd == `list_download` {
		jobDownload(c, e, arrData)
		return
	}

	jobShowList(c, e, arrData)
}

func jobValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})
	checkMap[`searchName`] = map[string]string{
		`type`: `:s`,
		`name`: `标题`,
	}

	//根据不同场景单独设置
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`jobno`] = map[string]string{
				`type`: `:d`,
				`name`: `No`,
			}
		} else {
			checkMap[`jobno`] = map[string]string{
				`type`: `d`,
				`name`: `No`,
			}
		}
		checkMap[`jobname`] = map[string]string{
			`type`:   `ss`,
			`name`:   `标题`,
			`maxlen`: `200`,
		}
		checkMap[`jobheadcount`] = map[string]string{
			`type`:   `d`,
			`name`:   `人数`,
			`maxlen`: `200`,
		}
		checkMap[`jobaddress`] = map[string]string{
			`type`:   `ss`,
			`name`:   `工作地点`,
			`maxlen`: `200`,
		}
		checkMap[`jobtreatment`] = map[string]string{
			`type`:   `d`,
			`name`:   `待遇`,
			`maxlen`: `10`,
		}
		checkMap[`jobshowdate`] = map[string]string{
			`type`: `date`,
			`name`: `发布日期`,
		}
		checkMap[`jobperiod`] = map[string]string{
			`type`:   `d`,
			`name`:   `有效期限`,
			`maxlen`: `10`,
		}
		checkMap[`jobdescription`] = map[string]string{
			`type`:   `ss`,
			`name`:   `描述`,
			`maxlen`: `200`,
		}
	} else if cmd == `ajax_edit_lang` {
		checkMap[`langcode`] = map[string]string{
			`type`: `:s`,
		}
		checkMap[`jobno`] = map[string]string{
			`type`: `d`,
			`name`: `No`,
		}
		checkMap[`jobname`] = map[string]string{
			`type`:   `ss`,
			`name`:   `标题`,
			`maxlen`: `200`,
		}
		checkMap[`jobaddress`] = map[string]string{
			`type`:   `ss`,
			`name`:   `工作地点`,
			`maxlen`: `200`,
		}
		checkMap[`jobdescription`] = map[string]string{
			`type`:   `ss`,
			`name`:   `描述`,
			`maxlen`: `200`,
		}
	}

	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`[job] valid error:`, e)

		if cmd == `ajax_detail` || cmd == `ajax_add` ||
			cmd == `ajax_edit` || cmd == `ajax_edit_lang` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			jobShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func jobShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	tData := new(tbls.TJob)
	arrData[`res`] = tData.GetListPage(db, arrData, appCfg.pagerow)
	arrData[`gMapHeadCount`] = gMapHeadCount
	arrData[`gMapTreatment`] = gMapTreatment

	fmt.Println(`[job]list`)

	SetGinGlobal(c, []string{`common`, `isHeaderSearch`}, false)
	SetGinGlobal(c, []string{`common`, `isHeaderListInfo`}, true)

	c.HTML(
		http.StatusOK,
		"job.html",
		MakeTemplateMap(c, e, arrData, tData),
	)
}

func jobAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[job]ajax detail`)
	log.Println("----------------------")
	fmt.Println("------------------------")

	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TJob)
		ok, data := tData.GetData(db, searchNo)
		fmt.Println(data)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"ret":  0,
				"data": data,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 1001,
		"msg": "数据不存在",
	})
}

func jobAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[job]ajax add`)

	name := klib.MapForKey(arrData, `jobname`)
	headcount := klib.MapForKey(arrData, `jobheadcount`)
	address := klib.MapForKey(arrData, `jobaddress`)
	treatment := klib.MapForKey(arrData, `jobtreatment`)
	showdate := klib.MapForKey(arrData, `jobshowdate`)
	period := klib.MapForKey(arrData, `jobperiod`)
	description := klib.MapForKey(arrData, `jobdescription`)

	tData := new(tbls.TJob)

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.JobName = name
	tData.JobHeadcount, _ = strconv.ParseInt(headcount, 10, 64)
	tData.JobAddress = address
	tData.JobTreatment, _ = strconv.ParseInt(treatment, 10, 64)
	tData.JobShowDate, _ = time.Parse(`2006-01-02`, showdate)
	tData.JobPeriod, _ = strconv.ParseInt(period, 10, 64)
	tData.JobDescription = description
	tData.JobInputdate = time.Now()
	tData.JobInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.JobInputid = int64(value.(float64))
	}
	tData.JobDelflg = `0`
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func jobAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[job]ajax edit`)

	no := klib.MapForKey(arrData, `jobno`)
	name := klib.MapForKey(arrData, `jobname`)
	headcount := klib.MapForKey(arrData, `jobheadcount`)
	address := klib.MapForKey(arrData, `jobaddress`)
	treatment := klib.MapForKey(arrData, `jobtreatment`)
	showdate := klib.MapForKey(arrData, `jobshowdate`)
	period := klib.MapForKey(arrData, `jobperiod`)
	description := klib.MapForKey(arrData, `jobdescription`)

	tData := new(tbls.TJob)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"jobname": `数据不存在！`,
			},
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.JobNo, _ = strconv.ParseInt(no, 10, 64)
	tData.JobName = name
	tData.JobHeadcount, _ = strconv.ParseInt(headcount, 10, 64)
	tData.JobAddress = address
	tData.JobTreatment, _ = strconv.ParseInt(treatment, 10, 64)
	tData.JobShowDate, _ = time.Parse(`2006-01-02`, showdate)
	tData.JobPeriod, _ = strconv.ParseInt(period, 10, 64)
	tData.JobDescription = description
	tData.JobUpdate = time.Now()
	tData.JobUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.JobUpid = int64(value.(float64))
	}
	tData.JobDelflg = `0`

	db.Model(&tData).
		Where("job_delflg = '0' and job_no = ?", tData.JobNo).
		Updates(tbls.TJob{
			JobName:        tData.JobName,
			JobHeadcount:   tData.JobHeadcount,
			JobAddress:     tData.JobAddress,
			JobTreatment:   tData.JobTreatment,
			JobShowDate:    tData.JobShowDate,
			JobPeriod:      tData.JobPeriod,
			JobDescription: tData.JobDescription,
			JobUpdate:      tData.JobUpdate,
			JobUpid:        tData.JobUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func jobAjaxDetailLang(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[job]ajax detail lang`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	langCode := klib.MapForKey(arrData, `langcode`)
	if len(searchNo) > 0 {
		tSrc := new(tbls.TJob)
		ok, sData := tSrc.GetData(db, searchNo)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1001,
				"msg": `数据不存在！`,
			})
			return
		}

		tData := new(tbls.TJobTrans)
		ok, data := tData.GetData(db, searchNo, langCode)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"ret":   0,
				"data":  data,
				"sData": sData,
			})
		} else {
			tData.JobNo = sData.JobNo
			c.JSON(http.StatusOK, gin.H{
				"ret":   0,
				"data":  tData,
				"sData": sData,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 1001,
		"msg": "数据不存在",
	})
}

func jobAjaxEditLang(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[job]ajax edit lang`)

	langCode := klib.MapForKey(arrData, `langcode`)
	no := klib.MapForKey(arrData, `jobno`)
	name := klib.MapForKey(arrData, `jobname`)
	address := klib.MapForKey(arrData, `jobaddress`)
	description := klib.MapForKey(arrData, `jobdescription`)

	tData := new(tbls.TJob)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"jobname": `数据不存在！`,
			},
		})
		return
	}

	data := new(tbls.TJobTrans)
	data.JobNo, _ = strconv.ParseInt(no, 10, 64)
	data.JobName = name
	data.JobAddress = address
	data.JobDescription = description
	data.JobLanguage = langCode

	ok = data.CheckNo(db, no, langCode)
	if ok {
		db.Model(&data).
			Where("job_no = ? AND job_language = ?", data.JobNo, langCode).
			Updates(tbls.TJob{
				JobName:        data.JobName,
				JobAddress:     data.JobAddress,
				JobDescription: data.JobDescription})
	} else {
		db.Create(&data)
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func jobDelete(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[job]delete`)

	arrNo := make([]string, 0)

	searchNo := klib.MapForKey(arrData, `searchNo`)

	selectNo := make([]string, 0)
	if value, ok := arrData[`selectNo`]; ok {
		if reflect.TypeOf(value).Name() != "string" {
			selectNo = value.([]string)
		} else {
			selectNo = append(selectNo, value.(string))
		}
	}

	if searchNo != `` {
		arrNo = append(arrNo, searchNo)
	} else {
		arrNo = append(arrNo, selectNo...)
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	for _, item := range arrNo {
		tData := new(tbls.TJob)
		ok := tData.CheckNo(db, item)
		if !ok {
			continue
		}

		no, _ := strconv.ParseInt(item, 10, 64)

		tData.JobNo = no
		tData.JobUpdate = time.Now()
		tData.JobUpid = 0
		if value, ok := sessionAdmin[`adm_no`]; ok {
			tData.JobUpid = int64(value.(float64))
		}
		tData.JobDelflg = `1`

		db.Model(&tData).
			Where("job_no = ?", tData.JobNo).
			Updates(tbls.TJob{
				JobUpdate: tData.JobUpdate,
				JobUpid:   tData.JobUpid,
				JobDelflg: tData.JobDelflg})

		//删除多语言
		db.Where(`job_no = ?`, tData.JobNo).
			Delete(tbls.TJobTrans{})
	}

	jobShowList(c, e, arrData)
}

func jobDownload(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[job]download`)

	tData := new(tbls.TJob)
	config := tData.DownloadConfig()
	rows, err := tData.GetListNative(db, arrData)
	if err != nil {
		klib.WriteLog(`[job]GetListNative error`)
		jobShowList(c, e, arrData)
		return
	}
	defer rows.Close()

	filepath := SqlRows2Xlsx(config, rows)

	ok := CommonDownload(c, filepath, "招聘.xlsx")
	if !ok {
		jobShowList(c, e, arrData)
	}
}
