package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"Yinghao/klib"
	"Yinghao/tbls"
	"net/http"
)

func ArchtureHandle(c *gin.Context){
	SetGinGlobal(c, []string{`common`, `router`}, `architure`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := aboutatchitureValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `edit_save` {
		aboutatchitectureEdit(c, e, arrData)
		return
	}

	aboutarchitureShowIndex(c, e, arrData)
}

func aboutatchitureValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

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
		fmt.Println(`[architure] valid error:`, e)
		aboutarchitureShowIndex(c, e, arrData)
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func aboutatchitectureEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[organization]edit`)

	langCode := klib.MapForKey(arrData, `langcode`)
	content  := klib.MapForKey(arrData, `content`)

	var ret bool
	if langCode != `` {
		tData := new(tbls.TSystemTrans)
		ret = tData.WriteData(db, `organization`, langCode, content)
	} else {
		tData := new(tbls.TSystem)
		ret = tData.WriteData(db, `organization`, content)
	}

	if ret {
		e[`commomMsg`] = `保存成功`
	} else {
		e[`commomMsg`] = `保存失败`
	}
	aboutarchitureShowIndex(c, e, arrData)
}

func aboutarchitureShowIndex(c *gin.Context, e map[string]string, arrData map[string]interface{}) {
	langCode := klib.MapForKey(arrData, `langcode`)

	if langCode != `` {
		tData := new(tbls.TSystemTrans)
		arrData[`content`] = tData.ReadData(db, `organization`, langCode, ``)
	} else {
		tData := new(tbls.TSystem)
		arrData[`content`] = tData.ReadData(db, `organization`, ``)
	}

	fmt.Println(`[organization]index`)

	c.HTML(
		http.StatusOK,
		"organizational.html",
		MakeTemplateMap(c, e, arrData, nil),
	)
}
