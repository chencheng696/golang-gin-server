/*
注意：前台固定提交的key是小写不带下划线，例如checkMap[`jobno`]
*/

package main

import (
	"fmt"
	"net/http"

	"Yinghao/klib"
	"Yinghao/tbls"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func aboutusHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `aboutus`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := aboutusValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `edit_save` {
		aboutusEdit(c, e, arrData)
		return
	}

	aboutusShowIndex(c, e, arrData)
}

func aboutusValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

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
		fmt.Println(`[aboutus] valid error:`, e)

		aboutusShowIndex(c, e, arrData)

		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func aboutusShowIndex(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	langCode := klib.MapForKey(arrData, `langcode`)

	if langCode != `` {
		tData := new(tbls.TSystemTrans)
		arrData[`content`] = tData.ReadData(db, `aboutus`, langCode, ``)
	} else {
		tData := new(tbls.TSystem)
		arrData[`content`] = tData.ReadData(db, `aboutus`, ``)
	}

	fmt.Println(`[aboutus]index`)

	c.HTML(
		http.StatusOK,
		"aboutus.html",
		MakeTemplateMap(c, e, arrData, nil),
	)
}

func aboutusEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[aboutus]edit`)

	langCode := klib.MapForKey(arrData, `langcode`)
	content := klib.MapForKey(arrData, `content`)

	var ret bool
	if langCode != `` {
		tData := new(tbls.TSystemTrans)
		ret = tData.WriteData(db, `aboutus`, langCode, content)
	} else {
		tData := new(tbls.TSystem)
		ret = tData.WriteData(db, `aboutus`, content)
	}

	if ret {
		e[`commomMsg`] = `保存成功`
	} else {
		e[`commomMsg`] = `保存失败`
	}
	aboutusShowIndex(c, e, arrData)
}
