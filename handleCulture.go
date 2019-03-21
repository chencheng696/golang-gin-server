package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"Yinghao/klib"
	"Yinghao/tbls"
	"net/http"
)

func CultureHandle(c *gin.Context){
	SetGinGlobal(c, []string{`common`, `router`}, `culture`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := cultureValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `edit_save` {
		cultureEdit(c, e, arrData)
		return
	}

	gloriesShowIndex(c, e, arrData)
}

func cultureValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	//默认从结构体tag中读取 校验信息
	checkMap := make(map[string]interface{})
	checkMap[`searchName`] = map[string]string{
		`type`: `:s`,
		`name`: `标题`,
	}
	checkMap[`langcode`] = map[string]string{
		`type`: `:ss`,
	}

	//根据不同场景单独设置
	if cmd == `edit_save` {
		checkMap[`content`] = map[string]string{
			`type`:   `:t`,
			`name`:   `内容`,
			`maxlen`: `10000`,
		}
	}

	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`[culture] valid error:`, e)
		cultureShowIndex(c, e, arrData)
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func cultureEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[culture]edit`)

	langCode := klib.MapForKey(arrData, `langcode`)
	content  := klib.MapForKey(arrData, `content`)

	var ret bool
	if langCode != `` {
		tData := new(tbls.TSystemTrans)
		ret = tData.WriteData(db, `culture`, langCode, content)
	} else {
		tData := new(tbls.TSystem)
		ret = tData.WriteData(db, `culture`, content)
	}

	if ret {
		e[`commomMsg`] = `保存成功`
	} else {
		e[`commomMsg`] = `保存失败`
	}
	cultureShowIndex(c, e, arrData)
}

func cultureShowIndex(c *gin.Context, e map[string]string, arrData map[string]interface{}) {
	langCode := klib.MapForKey(arrData, `langcode`)

	if langCode != `` {
		tData := new(tbls.TSystemTrans)
		arrData[`content`] = tData.ReadData(db, `culture`, langCode, ``)
	} else {
		tData := new(tbls.TSystem)
		arrData[`content`] = tData.ReadData(db, `culture`, ``)
	}

	fmt.Println(`[culture]index`)

	c.HTML(
		http.StatusOK,
		"culture.html",
		MakeTemplateMap(c, e, arrData, nil),
	)
}
