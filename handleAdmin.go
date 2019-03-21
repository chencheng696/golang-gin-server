/*
注意：前台固定提交的key是小写不带下划线，例如checkMap[`jobno`]
*/

package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"Yinghao/klib"
	"Yinghao/tbls"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func adminHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `admin`)

	cmd := c.DefaultPostForm("cmd", "")

	ok, e, arrData := adminValid(c, cmd)
	if !ok {
		return
	}

	if cmd == `list_del` {
		adminDelete(c, e, arrData)
		return
	} else if cmd == `ajax_add` {
		adminAjaxAdd(c, e, arrData)
		return
	} else if cmd == `ajax_edit` {
		adminAjaxEdit(c, e, arrData)
		return
	} else if cmd == `ajax_detail` {
		adminAjaxDetail(c, e, arrData)
		return
	} else if cmd == `list_download` {
		adminDownload(c, e, arrData)
		return
	} else if cmd == `list_upload` {
		adminUpload(c, e, arrData)
		return
	}

	adminShowList(c, e, arrData)
}

func adminValid(c *gin.Context, cmd string) (bool, map[string]string, map[string]interface{}) {

	checkMap := make(map[string]interface{})
	checkMap[`searchId`] = map[string]string{
		`type`: `:ss`,
		`name`: `管理员ID`,
	}
	checkMap[`searchName`] = map[string]string{
		`type`: `:t`,
		`name`: `管理员名字`,
	}
	checkMap[`uploadFile`] = map[string]string{
		`type`: `:t`,
		`name`: `文件`,
	}
	if cmd == `ajax_add` || cmd == `ajax_edit` {
		if cmd == `ajax_add` {
			checkMap[`admno`] = map[string]string{
				`type`: `:d`,
				`name`: `管理员no`,
			}
		} else {
			checkMap[`admno`] = map[string]string{
				`type`: `d`,
				`name`: `管理员no`,
			}
		}
		checkMap[`admid`] = map[string]string{
			`type`:   `ss`,
			`name`:   `管理员ID`,
			`maxlen`: `20`,
		}
		checkMap[`admname`] = map[string]string{
			`type`:   `ss`,
			`name`:   `管理员名字`,
			`maxlen`: `16`,
		}
		if cmd == `ajax_add` {
			checkMap[`admpwd`] = map[string]string{
				`type`:   `pass`,
				`name`:   `密码`,
				`minlen`: `6`,
				`maxlen`: `16`,
			}
		} else {
			checkMap[`admpwd`] = map[string]string{
				`type`: `:ss`,
				`name`: `密码`,
			}
		}
		checkMap[`admpwd2`] = map[string]string{
			`type`: `:ss`,
			`name`: `密码`,
		}

		checkMap[`admperm`] = map[string]string{
			`type`:   `ss`,
			`name`:   `权限`,
			`fixlen`: `1`,
		}
	}
	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`[admin] valid error`)

		if cmd == `ajax_detail` || cmd == `ajax_add` || cmd == `ajax_edit` {
			c.JSON(http.StatusOK, gin.H{
				"ret":   1000,
				"error": e,
			})
		} else {
			adminShowList(c, e, arrData)
		}
		return false, e, arrData
	} else {
		return true, e, arrData
	}
}

func adminShowList(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	tData := new(tbls.TAdmin)
	arrData[`res`] = tData.GetListPage(db, arrData, appCfg.pagerow)

	fmt.Println(`[admin]list`)

	//插入测试数据
	/*
		sessionAdmin := GetSessionMap(c, `admin`)
		for i := admin.RowCount; i < admin.RowCount+100; i++ {

			data := new(tbls.TAdmin)

			data.AdmId = `test` + klib.PadLeft(strconv.Itoa(i), 10, `0`)
			data.AdmPass = klib.MD5ForStr(`aaaa1111`)
			data.AdmName = `test` + klib.PadLeft(strconv.Itoa(i), 10, `0`)
			data.AdmPerm = `1`
			data.AdmInputdate = time.Now()
			data.AdmInputid = 0
			if value, ok := sessionAdmin[`adm_no`]; ok {
				data.AdmInputid = int64(value.(float64))
			}
			data.AdmDelflg = `0`
			db.Create(&data)
		}
	*/

	SetGinGlobal(c, []string{`common`, `isHeaderSearch`}, true)
	SetGinGlobal(c, []string{`common`, `isHeaderListInfo`}, true)

	c.HTML(
		http.StatusOK,
		"admin.html",
		MakeTemplateMap(c, e, arrData, tData),
	)
}

func adminAjaxDetail(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[admin]ajax detail`)

	searchNo := klib.MapForKey(arrData, `searchNo`)
	if len(searchNo) > 0 {
		tData := new(tbls.TAdmin)
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

func adminAjaxAdd(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[admin]ajax add`)

	admId := klib.MapForKey(arrData, `admid`)
	admPwd := klib.MapForKey(arrData, `admpwd`)
	admName := klib.MapForKey(arrData, `admname`)
	admPerm := klib.MapForKey(arrData, `admperm`)

	tData := new(tbls.TAdmin)
	ok := tData.CheckId(db, admId)
	if ok {
		c.JSON(
			http.StatusOK,
			gin.H{
				"ret": 1001,
				"error": gin.H{
					"admid": `管理员ID已经存在！`,
				},
			},
		)
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	tData.AdmId = admId
	tData.AdmPass = klib.MD5ForStr(admPwd)
	tData.AdmName = admName
	tData.AdmPerm = admPerm
	tData.AdmInputdate = time.Now()
	tData.AdmInputid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.AdmInputid = int64(value.(float64))
	}
	tData.AdmDelflg = `0`
	db.Create(&tData)

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func adminAjaxEdit(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[admin]ajax edit`)

	admNo := klib.MapForKey(arrData, `admno`)
	admId := klib.MapForKey(arrData, `admid`)
	admPwd := klib.MapForKey(arrData, `admpwd`)
	admName := klib.MapForKey(arrData, `admname`)
	admPerm := klib.MapForKey(arrData, `admperm`)

	tData := new(tbls.TAdmin)
	ok := tData.CheckNo(db, admNo)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1001,
			"error": gin.H{
				"admid": `管理员不存在！`,
			},
		})
		return
	}

	sessionAdmin := GetSessionMap(c, `admin`)

	no, _ := strconv.ParseInt(admNo, 10, 64)

	tData.AdmNo = no
	tData.AdmId = admId
	tData.AdmPass = klib.MD5ForStr(admPwd)
	tData.AdmName = admName
	tData.AdmPerm = admPerm
	tData.AdmUpdate = time.Now()
	tData.AdmUpid = 0
	if value, ok := sessionAdmin[`adm_no`]; ok {
		tData.AdmUpid = int64(value.(float64))
	}
	tData.AdmDelflg = `0`

	db.Model(&tData).
		Where("adm_delflg = '0' and adm_no = ?", tData.AdmNo).
		Updates(tbls.TAdmin{
			AdmName:   tData.AdmName,
			AdmPerm:   tData.AdmPerm,
			AdmUpdate: tData.AdmUpdate,
			AdmUpid:   tData.AdmUpid})

	c.JSON(http.StatusOK, gin.H{"ret": 0})
}

func adminDelete(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[admin]delete`)

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
		tData := new(tbls.TAdmin)
		ok := tData.CheckNo(db, item)
		if !ok {
			continue
		}

		no, _ := strconv.ParseInt(item, 10, 64)
		if no == 1 {
			continue
		}

		tData.AdmNo = no
		tData.AdmUpdate = time.Now()
		tData.AdmUpid = 0
		if value, ok := sessionAdmin[`adm_no`]; ok {
			tData.AdmUpid = int64(value.(float64))
		}
		tData.AdmDelflg = `1`

		db.Model(&tData).
			Where("adm_no = ?", tData.AdmNo).
			Updates(tbls.TAdmin{
				AdmUpdate: tData.AdmUpdate,
				AdmUpid:   tData.AdmUpid,
				AdmDelflg: tData.AdmDelflg})
	}

	adminShowList(c, e, arrData)
}

func adminDownload(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[admin]download`)

	tData := new(tbls.TAdmin)
	config := tData.DownloadConfig()
	rows, err := tData.GetListNative(db, arrData)
	if err != nil {
		klib.WriteLog(`[admin]GetListNative error`)
		adminShowList(c, e, arrData)
		return
	}
	defer rows.Close()

	filepath := SqlRows2Xlsx(config, rows)

	ok := CommonDownload(c, filepath, "管理员.xlsx")
	if !ok {
		adminShowList(c, e, arrData)
	}
}

func adminUpload(c *gin.Context, e map[string]string, arrData map[string]interface{}) {

	fmt.Println(`[admin]upload`)

	err, filepath := MoveUploadFile(arrData[`uploadFile`])
	if err != nil {
		e[`commomMsg`] = err.Error()
	} else {

		defer os.Remove(filepath)

		ok, list := klib.ReadXlsxMap(filepath, true)
		if !ok {
			e[`commomMsg`] = `文件格式不正确，请上传Excel文件！`
			adminShowList(c, e, arrData)
			return
		}

		//处理导入数据
		tData := new(tbls.TAdmin)
		config := tData.ImportConfig()

		err = ValidUpload(config, list)
		if err != nil {
			e[`commomMsg`] = err.Error()
			adminShowList(c, e, arrData)
			return
		}
		fmt.Println(list)

		e[`commomMsg`] = `导入成功！`
	}

	adminShowList(c, e, arrData)
}
