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

func PartnersHandle(c *gin.Context) {
	SetGinGlobal(c, []string{`common`, `router`}, `partners`)
	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := partnersValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `list_del` {
		partnersDelete(c, e, arrData)
		return
	} else if cmd == `ajax_add` {
		partnersAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		partnerAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		partnerAjaxDetail(c, e, arrData)
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

	partnersShowList(c, e, arrData)
}

func partnersDelete(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[partnerHandle]delete`)

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
		tData := new(tbls.TPartners)
		ok := tData.CheckNo(db, item)
		if !ok {
			continue
		}

		no, _ := strconv.ParseInt(item, 10, 64)

		tData.PnNo = no
		tData.PnUpdate = time.Now()
		tData.PnUpid = 0
		if value, ok := sessionAdmin[`adm_no`]; ok {
			tData.PnUpid = int64(value.(float64))
		}
		tData.PnDelflg = `1`

		db.Model(&tData).
			Where("pn_no = ?", tData.PnNo).
			Updates(tbls.TPartners{
			PnUpdate: tData.PnUpdate,
			PnUpid:   tData.PnUpid,
			PnDelflg: tData.PnDelflg})
	}

	partnersShowList(c, e, arrData)
}



func partnerAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[partner]ajax detail`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TPartners)
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

func partnersAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println("[partners]ajax add---------------------------------")

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

func partnerAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[partner]ajax edit`)

	no      := klib.MapForKey(arrData, `pnno`)
	name    := klib.MapForKey(arrData, `pnname`)
	url    := klib.MapForKey(arrData, `pnpath`)
	state   := klib.MapForKey(arrData, `pnstatus`)
	picPath := klib.MapForKey(arrData, `pnpicts`)

	tData := new(tbls.TPartners)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"partnername": `数据不存在！`,
			},
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.PnNo, _  = strconv.ParseInt(no, 10, 64)
	tData.PnName   = name
	tData.PnUrl    = url
	tData.PnStatus = state
    tData.PnPicts  = picPath
	tData.PnUpid   = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.PnUpid = int64(value.(float64))
	}
	tData.PnDelflg = `0`

	db.Model(&tData).
		Where("pn_delflg = '0' and pn_no = ?", tData.PnNo).
		Updates(tbls.TPartners{
		PnNo:        tData.PnNo,
		PnName:      tData.PnName,
		PnUrl:       tData.PnUrl,
		PnStatus:    tData.PnStatus,
		PnPicts:     tData.PnPicts,
		PnUpid:        tData.PnUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func partnersValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})
	checkMap[`searchName`] = map[string]string{
		`type`: `:s`,
		`name`: `标题`,
	}

	//根据不同场景单独设置
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`pnno`] = map[string]string{
				`type`: `:d`,
				`name`: `No`,
			}
		} else {
			checkMap[`pnno`] = map[string]string{
				`type`: `d`,
				`name`: `No`,
			}
		}
		checkMap[`pnname`] = map[string]string{
			`type`:   `ss`,
			`name`:   `标题`,
			`maxlen`: `200`,
		}
		checkMap[`pnurl`] = map[string]string{
			`type`:   `ss`,
			`name`:   `路径`,
			`maxlen`: `200`,
		}
		checkMap[`pnstatus`] = map[string]string{
			`type`:   `ss`,
			`name`:   `状态`,
			`maxlen`: `200`,
		}
		checkMap[`pnpicts`] = map[string]string{
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
		fmt.Println(`[partners] valid error:`, e)

		if cmd == `ajax_detail` || cmd == `ajax_add` ||
			cmd == `ajax_edit` || cmd == `ajax_edit_lang` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			partnersShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func partnersShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	tData := new(tbls.TPartners)
	arrData[`res`] = tData.GetListPage(db, arrData, appCfg.pagerow)
	arrData[`gMapHeadCount`] = gMapHeadCount
	arrData[`gMapTreatment`] = gMapTreatment

	fmt.Println(`[partners]list`)

	SetGinGlobal(c, []string{`common`, `isHeaderSearch`}, false)
	SetGinGlobal(c, []string{`common`, `isHeaderListInfo`}, true)

	c.HTML(
		http.StatusOK,
		"partners.html",
		MakeTemplateMap(c, e, arrData, tData),
	)
}