package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"Yinghao/klib"

	"github.com/tealeg/xlsx"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func GetCookie(c *gin.Context) map[string]string {
	m := make(map[string]string)
	for _, v := range c.Request.Cookies() {
		m[v.Name] = v.Value
	}
	return m
}

func GetSession(c *gin.Context, args ...string) interface{} {
	session := sessions.Default(c)

	var v interface{}
	if len(args) > 0 {
		v = session.Get(args[0])
	}

	var ok bool

	for n := 1; n < len(args); n++ {
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.String {
			var m map[string]interface{}

			err := json.Unmarshal([]byte(v.(string)), &m)
			if err != nil {
				return ``
			}

			v, ok = m[args[n]]
			if !ok {
				return ``
			}
		} else {
			return ``
		}
	}
	return v
}

func GetSessionMap(c *gin.Context, k string) map[string]interface{} {
	session := sessions.Default(c)
	v := session.Get(k)
	if v == nil {
		return nil
	}

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.String {
		var m map[string]interface{}

		err := json.Unmarshal([]byte(v.(string)), &m)
		if err != nil {
			return nil
		}

		return m
	} else {
		return nil
	}
}

//m,k...k,v
func SetGinGlobal(c *gin.Context, keys []string, value interface{}) {
	if len(keys) < 1 {
		return
	}

	m := gin.H{}
	if value, ok := c.Get(`global`); ok {
		m = value.(gin.H)
	} else {
		return
	}

	n := m
	for i := 0; i < len(keys)-1; i++ {
		if v, ok := n[keys[i]]; ok {
			n = v.(gin.H)
		}
	}

	n[keys[len(keys)-1]] = value
}

func MoveUploadFile(data interface{}) (error, string) {

	t := reflect.TypeOf(data)
	if t.String() != `*multipart.FileHeader` {
		klib.WriteLog(`没有上传文件`)
		return errors.New(`请选择文件！`), ``
	} else {
		header := data.(*multipart.FileHeader)

		n := float64(header.Size) / (1024.0 * 1024.0)
		if n > float64(appCfg.uploadMaxSize) {
			klib.WriteLog(`上传文件大小超出限制(` + strconv.Itoa(appCfg.uploadMaxSize) + `M)`)
			return errors.New(`上传文件大小超出限制(` + strconv.Itoa(appCfg.uploadMaxSize) + `M)`), ``
		}

		file, err := header.Open()
		if err != nil {
			klib.WriteLog(`打开上传文件失败` + err.Error())
			return errors.New(`上传文件失败！`), ``
		}

		i := 1
		tempfile := ``
		for tempfile == `` || klib.CheckFileIsExist(tempfile) {
			tempfile = appCfg.uploadTemp +
				strconv.FormatInt(time.Now().UTC().UnixNano(), 10) +
				strconv.Itoa(i)
			i++
		}

		out, err := os.Create(tempfile)
		if err != nil {
			klib.WriteLog(`创建临时文件失败` + err.Error())
			return errors.New(`上传文件失败！`), ``
		}
		_, err = io.Copy(out, file)
		if err != nil {
			klib.WriteLog(`复制上传文件到临时文件失败` + err.Error())
			return errors.New(`上传文件失败！`), ``
		}
		file.Close()
		out.Close()

		return nil, tempfile
	}
}

//path是uploadPict下一级目录名
func MoveUploadImage(data interface{}, path string) (error, string) {

	t := reflect.TypeOf(data)
	if t.String() != `*multipart.FileHeader` {
		return errors.New(`请选择文件！`), ``
	} else {
		header := data.(*multipart.FileHeader)
		contentType := header.Header.Get(`Content-Type`)
		if contentType != `image/png` &&
			contentType != `image/gif` &&
			contentType != `image/jpeg` {
			return errors.New(`上传图片格式不正确`), ``
		}
		contentType = `.` + strings.Replace(contentType, `image/`, ``, -1)

		n := float64(header.Size) / (1024.0 * 1024.0)
		if n > float64(appCfg.uploadMaxSize) {
			return errors.New(`上传文件大小超出限制(` + strconv.Itoa(appCfg.uploadMaxSize) + `M)`), ``
		}

		file, err := header.Open()
		if err != nil {
			return errors.New(`上传文件失败！`), ``
		}

		if path != `` {
			if ok, _ := klib.PathExists(appCfg.uploadPict + path); ok == false {
				os.Mkdir(appCfg.uploadPict+path, os.ModePerm)
			}
		}
		//upload/pict/yyyymm001/xxxxx.jpg
		monthDir := time.Now().Format("200601") + `/`
		if ok, _ := klib.PathExists(appCfg.uploadPict + path + monthDir); ok == false {
			os.Mkdir(appCfg.uploadPict+path+monthDir, os.ModePerm)
		}

		i := 1
		tempfile := ``
		for tempfile == `` || klib.CheckFileIsExist(tempfile) {
			tempfile = appCfg.uploadPict + path + monthDir +
				strconv.FormatInt(time.Now().UTC().UnixNano(), 10) +
				strconv.Itoa(i) + contentType
			i++
		}

		out, err := os.Create(tempfile)
		if err != nil {
			return errors.New(`上传文件失败！`), ``
		}
		_, err = io.Copy(out, file)
		if err != nil {
			return errors.New(`上传文件失败！`), ``
		}
		file.Close()
		out.Close()

		return nil, tempfile
	}
}

func SqlRows2Xlsx(config []map[string]string, rows *sql.Rows) string {

	var xlsRowTitle, xlsRow *xlsx.Row
	var cell *xlsx.Cell

	xlsxFile := xlsx.NewFile()
	sheet, _ := xlsxFile.AddSheet("sheet1")

	//add title
	xlsRowTitle = sheet.AddRow()
	for _, v := range config {
		cell = xlsRowTitle.AddCell()
		cell.Value, _ = v[`name`]
	}

	//add row
	columns, _ := rows.Columns()
	for rows.Next() {
		m := klib.SqlRow2Map(columns, rows)

		xlsRow = sheet.AddRow()
		for _, v := range config {
			k, ok := v[`key`]
			if !ok {
				continue
			}

			cell = xlsRow.AddCell()
			if value, ok := v[`convert`]; ok {
				f, ok := gMapFunc[value]
				if ok {
					cell.Value = f(m[k].(string))
				} else {
					cell.Value = m[k].(string)
				}
			} else {
				if reflect.TypeOf(m[k]).Name() == "Time" {
					t := m[k].(time.Time)
					if value, ok := v[`formatdate`]; ok {
						cell.Value = t.Format(value)
					} else {
						cell.Value = t.Format("2006-01-02 15:04:05")
					}
				} else {
					cell.Value = m[k].(string)
				}
			}
		}
	}

	i := 1
	filepath := ``
	for filepath == `` || klib.CheckFileIsExist(filepath) {
		filepath = appCfg.downloadExcel +
			strconv.FormatInt(time.Now().UTC().UnixNano(), 10) +
			strconv.Itoa(i)
		i++
	}
	xlsxFile.Save(filepath)
	return filepath
}

func CommonDownload(c *gin.Context, filepath, outname string) bool {

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}

	if fileInfo.Size() <= 0 {
		return false
	} else {
		file, err := os.Open(filepath)
		if err != nil {
			klib.WriteLog(err.Error())
			return false
		}
		defer file.Close()

		extraHeaders := map[string]string{
			`Content-Disposition`: `attachment; filename="` + outname + `"`,
		}

		c.DataFromReader(http.StatusOK, fileInfo.Size(), `application/octet-stream`, file, extraHeaders)
		return true
	}
}
