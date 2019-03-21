package main

import (
	"fmt"
	"net/http"
	"strings"

	"Yinghao/klib"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//editor图片上传
func editoruploadHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `editorupload`)

	fmt.Println(`editorupload`)

	checkMap := make(map[string]interface{})
	checkMap[`file`] = map[string]string{
		`type`: `t`,
	}

	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`%+v\n`, e)
		c.JSON(http.StatusOK, gin.H{
			"ret":      1001,
			"msg":      `非法数据！`,
			"filepath": ``,
		})
		return
	}
	err, filepath := MoveUploadImage(arrData[`file`], `editor/`)
	if err != nil {
		klib.WriteLog(`[editorupload] error` + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"ret":      1001,
			"msg":      err.Error(),
			"filepath": ``,
		})
	} else {
		filepath = strings.Replace(filepath, klib.GetAppPath()+"/upload", "", -1)
		klib.WriteLog(`[editorupload] success:` + filepath)
		c.JSON(http.StatusOK, gin.H{
			"ret":      0,
			"msg":      ``,
			"filepath": filepath,
		})
	}
}
