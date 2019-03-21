package main

import (
	"fmt"
	"net/http"

	"Yinghao/tbls"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"Yinghao/klib"
	"strconv"
	"time"
)

func newsClassHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `newsclass`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := newsclassValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `ajax_add` {
		newsclassAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		newsclassAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		newsclassAjaxDetail(c, e, arrData)
		return
	} else if cmd == `ajax_del` {
		newsclassAjaxDel(c, e, arrData)
		return
	} else if cmd == `ajax_detail_lang` {
		itemclassAjaxDetailLang(c, e, arrData)
		return
	} else if cmd == `ajax_edit_lang` {
		itemclassAjaxEditLang(c, e, arrData)
		return
	}

	newsclassShowList(c, e, arrData)
}

func newsclassValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})

	//根据不同场景单独设置
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`itemclassno`] = map[string]string{
				`type`: `:d`,
				`name`: `分类No`,
			}
		} else {
			checkMap[`itemclassno`] = map[string]string{
				`type`: `d`,
				`name`: `分类No`,
			}
		}
		checkMap[`itemclassname`] = map[string]string{
			`type`:   `ss`,
			`name`:   `名称`,
			`maxlen`: `200`,
		}
		checkMap[`itemclassparentno`] = map[string]string{
			`type`: `:d`,
		}
	} else if cmd == `ajax_edit_lang` {
		checkMap[`langcode`] = map[string]string{
			`type`: `:s`,
		}
		checkMap[`itemclassno`] = map[string]string{
			`type`: `d`,
			`name`: `No`,
		}
		checkMap[`itemclassname`] = map[string]string{
			`type`:   `ss`,
			`name`:   `标题`,
			`maxlen`: `200`,
		}
	}
	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`[newsclass] valid error`)

		if cmd == `ajax_detail` || cmd == `ajax_add` ||
			cmd == `ajax_edit` || cmd == `ajax_del` ||
			cmd == `ajax_edit_lang` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			newsclassShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}


func newsclassShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	tData := new(tbls.TNewsClass)
	res := tData.GetList(db)

	relation := make(map[int64]string)
	for _, item := range res {
		relation[item.NewsClassNo] = item.NewsClassName
	}

	arrData[`res`] = CreateTreeData(0, res, relation)

	fmt.Println(`[newsclass]list`)

	c.HTML(
		http.StatusOK,
		"newsclass.html",
		MakeTemplateMap(c, e, arrData, tData),
	)
}

func newsclassAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[newsclass]ajax add`)

	name := klib.MapForKey(arrData, `newsclassname`)
	parentno := klib.MapForKey(arrData, `newsclassparentno`)

	tData := new(tbls.TNewsClass)
	v, _ := strconv.ParseInt(parentno, 10, 64)
	if v > 0 {
		ok := tData.CheckNo(db, parentno)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1001,
				"error": gin.H{
					"newsclassname": `所属分类不存在！`,
				},
			})
			return
		}
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.NewsClassParentNo = v
	tData.NewsClassName = name
	tData.NewsClassInputdate = time.Now()
	tData.NewsClassInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.NewsClassInputid = int64(value.(float64))
	}
	tData.NewsClassDelflg = `0`
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func newsclassAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[newsclass]ajax edit`)

	no := klib.MapForKey(arrData, `newsclassno`)
	parentno := klib.MapForKey(arrData, `newsclassparentno`)
	name := klib.MapForKey(arrData, `newsclassname`)

	tData := new(tbls.TNewsClass)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"newsclassname": `分类不存在！`,
			},
		})
		return
	}

	v, _ := strconv.ParseInt(parentno, 10, 64)
	if v > 0 {
		ok := tData.CheckNo(db, parentno)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1001,
				"error": gin.H{
					"newsclassname": `所属分类不存在！`,
				},
			})
			return
		}
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.NewsClassNo, _ = strconv.ParseInt(no, 10, 64)
	tData.NewsClassParentNo = v
	tData.NewsClassName = name
	tData.NewsClassInputdate = time.Now()
	tData.NewsClassUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.NewsClassUpid = int64(value.(float64))
	}
	tData.NewsClassDelflg = `0`

	db.Model(&tData).
		Where("item_class_delflg = '0' and item_class_no = ?", tData.NewsClassNo).
		Updates(tbls.TNewsClass{
		NewsClassName:   tData.NewsClassName,
		NewsClassUpdate: tData.NewsClassUpdate,
		NewsClassUpid:   tData.NewsClassUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func newsclassAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[newsclass]ajax detail`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TNewsClass)
		res := klib.FormatCamelCase(tData.GetDataNative(db, searchNo))
		if len(res) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"ret":  0,
				"data": res[0],
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 1001,
		"msg": "数据不存在",
	})
}

func newsclassAjaxDel(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[itemclass]ajax delete`)

	searchNo := klib.MapForKey(arrData, `searchNo`)

	tData := new(tbls.TNewsClass)
	ok := tData.CheckNo(db, searchNo)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"msg": `分类不存在！`,
		})
		return
	}

	ok = tData.CheckChildren(db, searchNo)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"msg": `请清空子分类！`,
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	no, _ := strconv.ParseInt(searchNo, 10, 64)

	tData.NewsClassNo = no
	tData.NewsClassUpdate = time.Now()
	tData.NewsClassUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.NewsClassUpid = int64(value.(float64))
	}
	tData.NewsClassDelflg = `1`

	db.Model(&tData).
		Where("item_class_no = ?", tData.NewsClassNo).
		Updates(tbls.TItemClass{
		ItemClassUpdate: tData.NewsClassUpdate,
		ItemClassUpid:   tData.NewsClassUpid,
		ItemClassDelflg: tData.NewsClassDelflg})

	//删除多语言
	db.Where(`item_class_no = ?`, tData.NewsClassNo).
		Delete(tbls.TNewsClassTrans{})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}


/*
func itemclassAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[itemclass]ajax detail`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TItemClass)
		res := klib.FormatCamelCase(tData.GetDataNative(db, searchNo))
		if len(res) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"ret":  0,
				"data": res[0],
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 1001,
		"msg": "数据不存在",
	})
}

func itemclassAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[itemclass]ajax add`)

	name := klib.MapForKey(arrData, `itemclassname`)
	parentno := klib.MapForKey(arrData, `itemclassparentno`)

	tData := new(tbls.TItemClass)
	v, _ := strconv.ParseInt(parentno, 10, 64)
	if v > 0 {
		ok := tData.CheckNo(db, parentno)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1001,
				"error": gin.H{
					"itemclassname": `所属分类不存在！`,
				},
			})
			return
		}
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.ItemClassParentNo = v
	tData.ItemClassName = name
	tData.ItemClassInputdate = time.Now()
	tData.ItemClassInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.ItemClassInputid = int64(value.(float64))
	}
	tData.ItemClassDelflg = `0`
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func itemclassAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[itemclass]ajax edit`)

	no := klib.MapForKey(arrData, `itemclassno`)
	parentno := klib.MapForKey(arrData, `itemclassparentno`)
	name := klib.MapForKey(arrData, `itemclassname`)

	tData := new(tbls.TItemClass)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"itemclassname": `分类不存在！`,
			},
		})
		return
	}

	v, _ := strconv.ParseInt(parentno, 10, 64)
	if v > 0 {
		ok := tData.CheckNo(db, parentno)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1001,
				"error": gin.H{
					"itemclassname": `所属分类不存在！`,
				},
			})
			return
		}
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.ItemClassNo, _ = strconv.ParseInt(no, 10, 64)
	tData.ItemClassParentNo = v
	tData.ItemClassName = name
	tData.ItemClassUpdate = time.Now()
	tData.ItemClassUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.ItemClassUpid = int64(value.(float64))
	}
	tData.ItemClassDelflg = `0`

	db.Model(&tData).
		Where("item_class_delflg = '0' and item_class_no = ?", tData.ItemClassNo).
		Updates(tbls.TItemClass{
		ItemClassName:   tData.ItemClassName,
		ItemClassUpdate: tData.ItemClassUpdate,
		ItemClassUpid:   tData.ItemClassUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func itemclassAjaxDetailLang(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[itemclass]ajax detail lang`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	langCode := klib.MapForKey(arrData, `langcode`)
	if len(searchNo) > 0 {
		tSrc := new(tbls.TItemClass)
		ok, sData := tSrc.GetData(db, searchNo)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1001,
				"msg": `数据不存在！`,
			})
			return
		}

		tData := new(tbls.TItemClassTrans)
		ok, data := tData.GetData(db, searchNo, langCode)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"ret":   0,
				"data":  data,
				"sData": sData,
			})
		} else {
			tData.ItemClassNo = sData.ItemClassNo
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

func itemclassAjaxEditLang(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[itemclass]ajax edit lang`)

	langCode := klib.MapForKey(arrData, `langcode`)
	no := klib.MapForKey(arrData, `itemclassno`)
	name := klib.MapForKey(arrData, `itemclassname`)

	tData := new(tbls.TItemClass)
	ok := tData.CheckNo(db, no)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"itemclassname": `数据不存在！`,
			},
		})
		return
	}

	data := new(tbls.TItemClassTrans)
	data.ItemClassNo, _ = strconv.ParseInt(no, 10, 64)
	data.ItemClassName = name
	data.ItemClassLanguage = langCode

	ok = data.CheckNo(db, no, langCode)
	if ok {
		db.Model(&data).
			Where("item_class_no = ? AND item_class_language = ?", data.ItemClassNo, langCode).
			Updates(tbls.TItemClass{
			ItemClassName: data.ItemClassName})
	} else {
		db.Create(&data)
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func itemclassAjaxDel(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[itemclass]ajax delete`)

	searchNo := klib.MapForKey(arrData, `searchNo`)

	tData := new(tbls.TItemClass)
	ok := tData.CheckNo(db, searchNo)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"msg": `分类不存在！`,
		})
		return
	}

	ok = tData.CheckChildren(db, searchNo)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"msg": `请清空子分类！`,
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	no, _ := strconv.ParseInt(searchNo, 10, 64)

	tData.ItemClassNo = no
	tData.ItemClassUpdate = time.Now()
	tData.ItemClassUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.ItemClassUpid = int64(value.(float64))
	}
	tData.ItemClassDelflg = `1`

	db.Model(&tData).
		Where("item_class_no = ?", tData.ItemClassNo).
		Updates(tbls.TItemClass{
		ItemClassUpdate: tData.ItemClassUpdate,
		ItemClassUpid:   tData.ItemClassUpid,
		ItemClassDelflg: tData.ItemClassDelflg})

	//删除多语言
	db.Where(`item_class_no = ?`, tData.ItemClassNo).
		Delete(tbls.TItemClassTrans{})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}
*/
func CreateTreeData(parentNo int64, src []tbls.TNewsClass, relation map[int64]string) map[int64]interface{} {

	tree := make(map[int64]interface{})

	for _, item := range src {
		if item.NewsClassParentNo == parentNo {
			parentName, _ := relation[item.NewsClassParentNo]

			tree[item.NewsClassNo] = map[string]interface{}{
				`newsclassno`:         item.NewsClassNo,
				`newsclassparentno`:   item.NewsClassParentNo,
				`newsclassname`:       item.NewsClassName,
				`newsclassparentname`: parentName,
				`children`:            CreateTreeData(item.NewsClassNo, src, relation),
			}
		}
	}
	return tree
}
