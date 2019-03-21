/*
注意：前台固定提交的key是小写不带下划线，例如checkMap[`itemno`]
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
)

func itemHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `item`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := itemValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `list_del` {
		itemDelete(c, e, arrData)
		return
	} else if cmd == `ajax_add` {
		itemAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		itemAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		itemAjaxDetail(c, e, arrData)
		return
	} else if cmd == `ajax_detail_lang` {
		itemAjaxDetailLang(c, e, arrData)
		return
	} else if cmd == `ajax_edit_lang` {
		itemAjaxEditLang(c, e, arrData)
		return
	} else if cmd == `list_download` {
		itemDownload(c, e, arrData)
		return
	}

	itemShowList(c, e, arrData)
}

func itemValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})
	checkMap[`searchName`] = map[string]string{
		`type`: `:s`,
		`name`: `名称`,
	}

	//根据不同场景单独设置
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`itemno`] = map[string]string{
				`type`: `:d`,
				`name`: `No`,
			}
		} else {
			checkMap[`itemno`] = map[string]string{
				`type`: `d`,
				`name`: `No`,
			}
		}
		checkMap[`itemname`] = map[string]string{
			`type`:   `s`,
			`name`:   `名称`,
			`maxlen`: `1000`,
		}
		checkMap[`itemtype`] = map[string]string{
			`type`:   `:s`,
			`name`:   `规格`,
			`maxlen`: `1000`,
		}
		checkMap[`itemclassno`] = map[string]string{
			`type`: `d`,
			`name`: `分类`,
		}
		checkMap[`itemstatus`] = map[string]string{
			`type`:   `ss`,
			`name`:   `状态`,
			`fixlen`: `1`,
		}
		checkMap[`iteminfo`] = map[string]string{
			`type`:   `:t`,
			`name`:   `详情`,
			`maxlen`: `10000`,
		}
	} else if cmd == `ajax_edit_lang` {
		checkMap[`langcode`] = map[string]string{
			`type`: `:s`,
		}
		checkMap[`itemno`] = map[string]string{
			`type`: `d`,
			`name`: `No`,
		}
		checkMap[`itemname`] = map[string]string{
			`type`:   `s`,
			`name`:   `名称`,
			`maxlen`: `1000`,
		}
		checkMap[`itemtype`] = map[string]string{
			`type`:   `:s`,
			`name`:   `规格`,
			`maxlen`: `1000`,
		}
		checkMap[`iteminfo`] = map[string]string{
			`type`:   `:t`,
			`name`:   `详情`,
			`maxlen`: `10000`,
		}
	}

	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`[item] valid error:`, e)

		if cmd == `ajax_detail` || cmd == `ajax_add` ||
			cmd == `ajax_edit` || cmd == `ajax_edit_lang` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			itemShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func itemShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	tData := new(tbls.TItem)
	arrData[`res`] = tData.GetListPage(db, arrData, appCfg.pagerow)
	arrData[`gMapHeadCount`] = gMapHeadCount
	arrData[`gMapTreatment`] = gMapTreatment

	itemClass := new(tbls.TItemClass)
	arrData[`res_class`] = itemClass.GetTreeShow(0, itemClass.GetList(db))

	fmt.Println(`[item]list`)

	SetGinGlobal(c, []string{`common`, `isHeaderSearch`}, true)
	SetGinGlobal(c, []string{`common`, `isHeaderListInfo`}, true)

	c.HTML(
		http.StatusOK,
		"item.html",
		MakeTemplateMap(c, e, arrData, tData),
	)
}

func itemAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]ajax detail`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TItem)
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

func itemAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]ajax add`)

	name := klib.MapForKey(arrData, `itemname`)
	itemtype := klib.MapForKey(arrData, `itemtype`)
	classno := klib.MapForKey(arrData, `itemclassno`)
	status := klib.MapForKey(arrData, `itemstatus`)
	info := klib.MapForKey(arrData, `iteminfo`)

	classData := new(tbls.TItemClass)
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

	tData := new(tbls.TItem)

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.ItemName = name
	tData.ItemType = itemtype
	tData.ItemClassNo, _ = strconv.ParseInt(classno, 10, 64)
	tData.ItemStatus = status
	tData.ItemInfo = info
	tData.ItemHit = 0
	tData.ItemPicts = ``
	tData.ItemInputdate = time.Now()
	tData.ItemInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.ItemInputid = int64(value.(float64))
	}
	tData.ItemDelflg = `0`
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func itemAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]ajax edit`)

	no := klib.MapForKey(arrData, `itemno`)
	name := klib.MapForKey(arrData, `itemname`)
	itemtype := klib.MapForKey(arrData, `itemtype`)
	classno := klib.MapForKey(arrData, `itemclassno`)
	status := klib.MapForKey(arrData, `itemstatus`)
	info := klib.MapForKey(arrData, `iteminfo`)

	tData := new(tbls.TItem)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"itemname": `数据不存在！`,
			},
		})
		return
	}

	classData := new(tbls.TItemClass)
	ok = classData.CheckNo(db, classno)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"msg": `分类不存在！`,
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.ItemNo, _ = strconv.ParseInt(no, 10, 64)
	tData.ItemName = name
	tData.ItemType = itemtype
	tData.ItemClassNo, _ = strconv.ParseInt(classno, 10, 64)
	tData.ItemStatus = status
	tData.ItemInfo = info
	tData.ItemPicts = ``
	tData.ItemUpdate = time.Now()
	tData.ItemUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.ItemUpid = int64(value.(float64))
	}
	tData.ItemDelflg = `0`

	db.Model(&tData).
		Where("item_delflg = '0' and item_no = ?", tData.ItemNo).
		Updates(tbls.TItem{
			ItemName:    tData.ItemName,
			ItemType:    tData.ItemType,
			ItemClassNo: tData.ItemClassNo,
			ItemStatus:  tData.ItemStatus,
			ItemInfo:    tData.ItemInfo,
			ItemPicts:   tData.ItemPicts,
			ItemUpdate:  tData.ItemUpdate,
			ItemUpid:    tData.ItemUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func itemAjaxDetailLang(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]ajax detail lang`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	langCode := klib.MapForKey(arrData, `langcode`)
	if len(searchNo) > 0 {
		tSrc := new(tbls.TItem)
		ok, sData := tSrc.GetData(db, searchNo)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1001,
				"msg": `数据不存在！`,
			})
			return
		}

		tData := new(tbls.TItemTrans)
		ok, data := tData.GetData(db, searchNo, langCode)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"ret":   0,
				"data":  data,
				"sData": sData,
			})
		} else {
			tData.ItemNo = sData.ItemNo
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

func itemAjaxEditLang(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]ajax edit lang`)

	langCode := klib.MapForKey(arrData, `langcode`)
	no := klib.MapForKey(arrData, `itemno`)
	name := klib.MapForKey(arrData, `itemname`)
	itemtype := klib.MapForKey(arrData, `itemtype`)
	info := klib.MapForKey(arrData, `iteminfo`)

	tData := new(tbls.TItem)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"itemname": `数据不存在！`,
			},
		})
		return
	}

	data := new(tbls.TItemTrans)
	data.ItemNo, _ = strconv.ParseInt(no, 10, 64)
	data.ItemName = name
	data.ItemType = itemtype
	data.ItemInfo = info
	data.ItemLanguage = langCode

	ok = data.CheckNo(db, no, langCode)
	if ok {
		db.Model(&data).
			Where("item_no = ? AND item_language = ?", data.ItemNo, langCode).
			Updates(tbls.TItem{
				ItemName: data.ItemName,
				ItemType: data.ItemType,
				ItemInfo: data.ItemInfo})
	} else {
		db.Create(&data)
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func itemDelete(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]delete`)

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
		tData := new(tbls.TItem)
		ok := tData.CheckNo(db, item)
		if !ok {
			continue
		}

		no, _ := strconv.ParseInt(item, 10, 64)

		tData.ItemNo = no
		tData.ItemUpdate = time.Now()
		tData.ItemUpid = 0
		if value, ok := sessionAdmin[`adm_no`]; ok {
			tData.ItemUpid = int64(value.(float64))
		}
		tData.ItemDelflg = `1`

		db.Model(&tData).
			Where("item_no = ?", tData.ItemNo).
			Updates(tbls.TItem{
				ItemUpdate: tData.ItemUpdate,
				ItemUpid:   tData.ItemUpid,
				ItemDelflg: tData.ItemDelflg})

		//删除多语言
		db.Where(`item_no = ?`, tData.ItemNo).
			Delete(tbls.TItemTrans{})

		//删除图片
		//....
	}

	itemShowList(c, e, arrData)
}

func itemDownload(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[item]download`)

	tData := new(tbls.TItem)
	config := tData.DownloadConfig()
	rows, err := tData.GetListNative(db, arrData)
	if err != nil {
		klib.WriteLog(`[item]GetListNative error`)
		itemShowList(c, e, arrData)
		return
	}
	defer rows.Close()

	filepath := SqlRows2Xlsx(config, rows)

	ok := CommonDownload(c, filepath, "商品.xlsx")
	if !ok {
		itemShowList(c, e, arrData)
	}
}
