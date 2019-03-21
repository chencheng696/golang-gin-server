package main

import (
	"fmt"
	"net/http"
	"Yinghao/tbls"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"Yinghao/klib"
	"time"
	"strconv"
	"reflect"
)

func PromotionHandle(c *gin.Context) {
	SetGinGlobal(c, []string{`common`, `router`}, `promotion`)
	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := promotionValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `list_del` {
		promotionDelete(c, e, arrData)
		return
	} else if cmd == `ajax_add` {
		promotionAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		promotionAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		promotionAjaxDetail(c, e, arrData)
		return
	} else if cmd == `ajax_detail_lang` {
		//jobAjaxDetailLang(c, e, arrData)
		return
	} else if cmd == `ajax_edit_lang` {
		//jobAjaxEditLang(c, e, arrData)
		return
	} else if cmd == `list_download` {
		//jobDownload(c, e, arrData)
		return
	}

	promotionShowList(c, e, arrData)
}

func promotionAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println("[promotion]ajax add------------------------")
	//name := klib.MapForKey(arrData,    `promotionrno`)

	name := klib.MapForKey(arrData,    `promotionname`)
	path := klib.MapForKey(arrData,    `promotionpath`)
	state := klib.MapForKey(arrData,   `promotionstate`)
	picpath := klib.MapForKey(arrData, `promotionpicpath`)

	tData := new(tbls.TPromotion)

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.PromName = name
	tData.PromUrl = path
	tData.PromStatus = state
	tData.PromPicts = picpath
	tData.PromUpdate = time.Now()
	tData.PromHit = 0
	tData.PromDescription = "just do it"
	tData.PromInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.PromInputid = int64(value.(float64))
	}
	tData.PromDelflg = `0`
	fmt.Println(tData)
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func promotionAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {
	fmt.Println(`[promotion]ajax edit`)
	no      := klib.MapForKey(arrData, `promotionno`)
	name    := klib.MapForKey(arrData, `promotionname`)
	url    := klib.MapForKey(arrData, `promotionpath`)
	state   := klib.MapForKey(arrData, `promotionstate`)
	picPath := klib.MapForKey(arrData, `promotionpicpath`)

	tData := new(tbls.TPromotion)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"promotionname": `数据不存在！`,
			},
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)
	tData.PromNo, _  = strconv.ParseInt(no, 10, 64)
	tData.PromName = name
	tData.PromUrl = url
	tData.PromStatus = state
	tData.PromPicts = picPath

	tData.PromUpid   = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.PromUpid = int64(value.(float64))
	}
	tData.PromDelflg = `0`

	db.Model(&tData).
		Where("pn_delflg = '0' and pn_no = ?", tData.PromNo).
		Updates(tbls.TPromotion{
		PromNo:        tData.PromNo,
		PromName:      tData.PromName,
		PromUrl:       tData.PromUrl,
		PromStatus:    tData.PromStatus,
		PromPicts:     tData.PromPicts,
		PromUpid:        tData.PromUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func promotionDelete(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[promotionHandle]delete`)

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
		tData := new(tbls.TPromotion)
		ok := tData.CheckNo(db, item)
		if !ok {
			continue
		}
		no, _ := strconv.ParseInt(item, 10, 64)

		tData.PromNo = no
		tData.PromUpdate = time.Now()
		tData.PromUpid = 0
		if value, ok := sessionAdmin[`adm_no`]; ok {
			tData.PromUpid = int64(value.(float64))
		}
		tData.PromDelflg = `1`

		db.Model(&tData).
			Where("prom_no = ?", tData.PromNo).
			Updates(tbls.TPromotion{
			PromUpdate: tData.PromUpdate,
			PromUpid:   tData.PromUpid,
			PromDelflg: tData.PromDelflg})
	}

	promotionShowList(c, e, arrData)
}

func promotionAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {
	fmt.Println(`[promotion]ajax detail`)
	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TPromotion)
		ok, data := tData.GetData(db, searchNo)
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


/*
func partnersAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[partners]ajax add`)

	name := klib.MapForKey(arrData, `partnername`)
	path := klib.MapForKey(arrData, `partnerpath`)
	state := klib.MapForKey(arrData, `partnertate`)
	picpath := klib.MapForKey(arrData, `partnerpicpath`)

	tData := new(tbls.TPartners)

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.PnName = name
	tData.PnUrl = path
	tData.PnStatus = state
	tData.PnPicts = picpath
	tData.PnUpdate = time.Now()
	tData.PnInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.PnInputid = int64(value.(float64))
	}
	tData.PnDelflg = `0`
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}
*/
func promotionValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})
	checkMap[`searchName`] = map[string]string{
		`type`: `:s`,
		`name`: `标题`,
	}

	//根据不同场景单独设置
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`promotionrno`] = map[string]string{
				`type`: `:d`,
				`name`: `No`,
			}
		} else {
			checkMap[`promotionrno`] = map[string]string{
				`type`: `d`,
				`name`: `No`,
			}
		}
		checkMap[`promotionname`] = map[string]string{
			`type`:   `ss`,
			`name`:   `标题`,
			`maxlen`: `200`,
		}
		checkMap[`promotionpath`] = map[string]string{
			`type`:   `ss`,
			`name`:   `路径`,
			`maxlen`: `200`,
		}
		checkMap[`promotionstate`] = map[string]string{
			`type`:   `ss`,
			`name`:   `状态`,
			`maxlen`: `200`,
		}
		checkMap[`promotionpicpath`] = map[string]string{
			`type`:   `ss`,
			`name`:   `图片路径`,
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
		fmt.Println(`[promotion] valid error:`, e)

		if cmd == `ajax_detail` || cmd == `ajax_add` ||
			cmd == `ajax_edit` || cmd == `ajax_edit_lang` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			promotionShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func promotionShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	tData := new(tbls.TPromotion)
	arrData[`res`] = tData.GetListPage(db, arrData, appCfg.pagerow)
	arrData[`gMapHeadCount`] = gMapHeadCount
	arrData[`gMapTreatment`] = gMapTreatment

	fmt.Println(`[promotion]list`)

	SetGinGlobal(c, []string{`common`, `isHeaderSearch`}, false)
	SetGinGlobal(c, []string{`common`, `isHeaderListInfo`}, true)

	c.HTML(
		http.StatusOK,
		"promotion.html",
		MakeTemplateMap(c, e, arrData, tData),
	)
}