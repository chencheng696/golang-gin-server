/*
注意：前台固定提交的key是小写不带下划线，例如checkMap[`itemno`]
*/

package main

import (
	"fmt"
	"net/http"

	"Yinghao/tbls"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"time"
	"Yinghao/klib"
	"reflect"
)

func newsHandle(c *gin.Context) {
	SetGinGlobal(c, []string{`common`, `router`}, `news`)

	cmd := c.DefaultPostForm("cmd", "")

	fmt.Println(cmd)
	ok, e, arrData := newsValid(c, cmd)
	if !ok {
		return
	}

	fmt.Println(cmd)
	if cmd == `list_del` {
		newsDelete(c, e, arrData)
		return
	} else if cmd == `ajax_add` {
		newsAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		newsAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		newsAjaxDetail(c, e, arrData)
		return
	} else if cmd == `ajax_detail_lang` {
		//itemAjaxDetailLang(c, e, arrData)
		return
	} else if cmd == `ajax_edit_lang` {
		//itemAjaxEditLang(c, e, arrData)
		return
	} else if cmd == `list_download` {
		//itemDownload(c, e, arrData)
		return
	}

	newsShowList(c, e, arrData)
}

func newsValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})
	checkMap[`searchName`] = map[string]string{
		`type`: `:s`,
		`name`: `名称`,
	}

	//根据不同场景单独设置
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`newsno`] = map[string]string{
				`type`: `:d`,
				`name`: `No`,
			}
		} else {
			checkMap[`newsno`] = map[string]string{
				`type`: `d`,
				`name`: `No`,
			}
		}
		checkMap[`newsname`] = map[string]string{
			`type`:   `s`,
			`name`:   `名称`,
			`maxlen`: `1000`,
		}
		//checkMap[`newstype`] = map[string]string{
		//	`type`:   `:s`,
		//	`name`:   `规格`,
		//	`maxlen`: `1000`,
		//}
		checkMap[`newsclassno`] = map[string]string{
			`type`: `d`,
			`name`: `分类`,
		}
		checkMap[`newsstatus`] = map[string]string{
			`type`:   `ss`,
			`name`:   `状态`,
			`fixlen`: `1`,
		}
		checkMap[`newsinfo`] = map[string]string{
			`type`:   `:t`,
			`name`:   `详情`,
			`maxlen`: `10000`,
		}
	} else if cmd == `ajax_edit_lang` {
		checkMap[`langcode`] = map[string]string{
			`type`: `:s`,
		}
		checkMap[`newsno`] = map[string]string{
			`type`: `d`,
			`name`: `No`,
		}
		checkMap[`newsname`] = map[string]string{
			`type`:   `s`,
			`name`:   `名称`,
			`maxlen`: `1000`,
		}
		checkMap[`newstype`] = map[string]string{
			`type`:   `:s`,
			`name`:   `规格`,
			`maxlen`: `1000`,
		}
		checkMap[`newsinfo`] = map[string]string{
			`type`:   `:t`,
			`name`:   `详情`,
			`maxlen`: `10000`,
		}
	}

	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`[news] valid error:`, e)

		if cmd == `ajax_detail` || cmd == `ajax_add` ||
			cmd == `ajax_edit` || cmd == `ajax_edit_lang` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			newsShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func newsShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {
	tData := new(tbls.TNews)
	arrData[`res`] = tData.GetListPage(db, arrData, appCfg.pagerow)
	arrData[`gMapHeadCount`] = gMapHeadCount
	arrData[`gMapTreatment`] = gMapTreatment

	fmt.Println(`[news]list`)

	SetGinGlobal(c, []string{`common`, `isHeaderSearch`}, true)
	SetGinGlobal(c, []string{`common`, `isHeaderListInfo`}, true)

	c.HTML(
		http.StatusOK,
		"news.html",
		MakeTemplateMap(c, e, arrData, tData),
	)
}

func newsAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[news]ajax add`)

	name := klib.MapForKey(arrData, `newsname`)
	//itemtype := klib.MapForKey(arrData, `newstype`)
	classno := klib.MapForKey(arrData, `newsclassno`)
	status := klib.MapForKey(arrData, `newsstatus`)
	info := klib.MapForKey(arrData, `newsinfo`)

	classData := new(tbls.TNewsClass)
	ok := classData.CheckNo(db, classno)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"msg": `分类不存在！`,
		})
		return
	}

	//处理图片
	//...

	tData := new(tbls.TNews)

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.NewsName = name
	//tData.ItemType = itemtype
	tData.NewsClassNo, _ = strconv.ParseInt(classno, 10, 64)
	tData.NewsStatus = status
	tData.NewsInfo = info
	tData.NewsHit = 0
	tData.NewsPicts = ``
	tData.NewsInputdate = time.Now()
	tData.NewsInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.NewsInputid = int64(value.(float64))
	}
	tData.NewsDelflg = `0`
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func newsDelete(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[news]delete`)

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
		tData := new(tbls.TNews)
		ok := tData.CheckNo(db, item)
		if !ok {
			continue
		}

		no, _ := strconv.ParseInt(item, 10, 64)

		tData.NewsNo = no
		tData.NewsUpdate = time.Now()
		tData.NewsUpid = 0
		if value, ok := sessionAdmin[`adm_no`]; ok {
			tData.NewsUpid = int64(value.(float64))
		}
		tData.NewsDelflg = `1`

		db.Model(&tData).
			Where("item_no = ?", tData.NewsNo).
			Updates(tbls.TNews{
			NewsUpdate: tData.NewsUpdate,
			NewsUpid:   tData.NewsUpid,
			NewsDelflg: tData.NewsDelflg})

		//删除多语言
		db.Where(`item_no = ?`, tData.NewsNo).
			Delete(tbls.TNewsClassTrans{})

		//删除图片
		//....
	}

	newsShowList(c, e, arrData)
}

func newsAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[news]ajax edit`)

	no := klib.MapForKey(arrData, `newsno`)
	name := klib.MapForKey(arrData, `newsname`)
	//itemtype := klib.MapForKey(arrData, `newstype`)
	classno := klib.MapForKey(arrData, `newsclassno`)
	status := klib.MapForKey(arrData, `newsstatus`)
	info := klib.MapForKey(arrData, `newsinfo`)

	tData := new(tbls.TNews)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"newsname": `数据不存在！`,
			},
		})
		return
	}

	classData := new(tbls.TNewsClass)
	ok = classData.CheckNo(db, classno)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"msg": `分类不存在！`,
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.NewsNo, _ = strconv.ParseInt(no, 10, 64)
	tData.NewsName = name
	//tData.ItemType = itemtype
	tData.NewsClassNo, _ = strconv.ParseInt(classno, 10, 64)
	tData.NewsStatus = status
	tData.NewsInfo = info
	tData.NewsPicts = ``
	tData.NewsUpdate = time.Now()
	tData.NewsUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.NewsUpid = int64(value.(float64))
	}
	tData.NewsDelflg = `0`

	db.Model(&tData).
		Where("item_delflg = '0' and item_no = ?", tData.NewsNo).
		Updates(tbls.TNews{
		NewsName:    tData.NewsName,
		//ItemType:    tData.ItemType,
		NewsClassNo: tData.NewsClassNo,
		NewsStatus:  tData.NewsStatus,
		NewsInfo:    tData.NewsInfo,
		NewsPicts:   tData.NewsPicts,
		NewsUpdate:  tData.NewsUpdate,
		NewsUpid:    tData.NewsUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func newsAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]ajax detail`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TNews)
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
