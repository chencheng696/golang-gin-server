package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"Yinghao/tbls"
	"strconv"
	"time"
	"Yinghao/klib"
	"reflect"
)


func visatorValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})
	checkMap[`searchName`] = map[string]string{
		`type`: `:s`,
		`name`: `标题`,
	}

	fmt.Println("1111111111111111111111")
	//根据不同场景单独设置
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`fbno`] = map[string]string{
				`type`: `:d`,
				`name`: `No`,
			}
		} else {
			checkMap[`fbno`] = map[string]string{
				`type`: `d`,
				`name`: `No`,
			}
		}
		checkMap[`fbcontent`] = map[string]string{
			`type`:   `ss`,
			`name`:   `内容`,
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

	//根据不同场景单独设置
	e, arrData := ValidForm(c, checkMap)
	log.Println(arrData)
	if !CheckValidResult(e) {
		fmt.Println(`[visator] valid error:`, e)
		if cmd == `ajax_detail` || cmd == `ajax_add` ||
			cmd == `ajax_edit` || cmd == `ajax_edit_lang` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			visatorShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func visatorShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {
	fmt.Println(arrData)
	tData := new(tbls.TFeedback)
	arrData[`res`] = tData.GetListPage(db, arrData, appCfg.pagerow)
	arrData[`gMapHeadCount`] = gMapHeadCount
	arrData[`gMapTreatment`] = gMapTreatment

	fmt.Println(`[visator]list`)

	SetGinGlobal(c, []string{`common`, `isHeaderSearch`}, false)
	SetGinGlobal(c, []string{`common`, `isHeaderListInfo`}, true)
	c.HTML(
		http.StatusOK,
		"visator.html",
		MakeTemplateMap(c, e, arrData, tData),
	)

}

func visatorHandle(c *gin.Context) {
	SetGinGlobal(c, []string{`common`, `router`}, `visator`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := visatorValid(c, cmd)
	if !ok {

		return
	}

	if cmd == `list_del` {
		visitorDelete(c, e, arrData)
		return
	} else if cmd == `ajax_add` {
		//visitorAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		VisatorAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		visitorAjaxDetail(c, e, arrData)
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

	visatorShowList(c, e, arrData)

}

func VisatorAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[feedback]ajax edit`)

	no := klib.MapForKey(arrData, `fbno`)
	content := klib.MapForKey(arrData, `fbcontent`)

	tData := new(tbls.TFeedback)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"content": `数据不存在！`,
			},
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)
	tData.FbNo, _ = strconv.ParseInt(no, 10, 64)
	tData.FbContent = content
	tData.FbUpdate = time.Now()
	tData.FbUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.FbUpid = int64(value.(float64))
	}
	tData.FbDelflg = `0`

	db.Model(&tData).
		Where("fb_delflg = '0' and fb_no = ?", tData.FbNo).
		Updates(tbls.TFeedback{
		FbContent:        tData.FbContent,
		FbUpdate:         tData.FbUpdate,
		FbUpid:           tData.FbUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}


func visitorAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[visitor]ajax detail`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	fmt.Println(searchNo)
	if len(searchNo) > 0 {
		tData := new(tbls.TFeedback)
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


func visitorDelete(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[visitorHandle]delete`)

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
		tData := new(tbls.TFeedback)
		ok := tData.CheckNo(db, item)
		if !ok {
			continue
		}

		no, _ := strconv.ParseInt(item, 10, 64)

		tData.FbNo = no
		tData.FbUpdate = time.Now()
		tData.FbUpid = 0
		if value, ok := sessionAdmin[`adm_no`]; ok {
			tData.FbUpid = int64(value.(float64))
		}
		tData.FbDelflg = `1`

		db.Model(&tData).
			Where("fb_no = ?", tData.FbNo).
			Updates(tbls.TFeedback{
			FbUpdate: tData.FbUpdate,
			FbUpid:   tData.FbUpid,
			FbDelflg: tData.FbDelflg})
	}

	visatorShowList(c, e, arrData)
}

