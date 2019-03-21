package main

import (
	"io"
	"net/http"
	"os"

	"Yinghao/klib"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//图片下载
func pictHandle(c *gin.Context) {

	flag := c.Param(`flag`)
	yyyymm := c.Param(`yyyymm`)
	name := c.Param(`name`)

	if yyyymm == `` || flag == `` || name == `` {
		klib.WriteLog(`[pictHandle]参数不正确`)
		c.String(http.StatusNotFound, ``)
		return
	}

	l := len(name)

	if l > 5 && name[(l-5):] == `-s220` {

		filefull := appCfg.uploadPict + flag + `/` + yyyymm + `/` + name[:(l-5)]
		if !klib.CheckFileIsExist(filefull) {
			klib.WriteLog(`[pictHandle]原图片不存在:` + filefull)
			c.String(http.StatusNotFound, ``)
			return
		}

		filename := appCfg.uploadPict + flag + `/` + yyyymm + `/` + name[:(l-5)]
		if klib.CheckFileIsExist(filename) {
			c.File(filename)
			return
		}

		localPath, format, _ := klib.IsPictureFormat(filefull)
		if localPath == `` || format == `` {
			klib.WriteLog(`[pictHandle]非图片文件:` + filefull)
			c.String(http.StatusNotFound, ``)
			return
		}

		if !klib.ImageCompress(
			func() (io.Reader, error) {
				return os.Open(localPath)
			},
			func() (*os.File, error) {
				return os.Open(localPath)
			},
			filename,
			100,
			150,
			format) {
			klib.WriteLog(`[pictHandle]生成缩略图失败:` + filefull)
		} else {
			klib.WriteLog(`[pictHandle]生成缩略图成功:` + filename)
			c.File(filename)
		}
	} else {
		if !klib.CheckFileIsExist(appCfg.uploadPict + flag + `/` + yyyymm + `/` + name) {
			c.String(http.StatusNotFound, ``)
			return
		}
		c.File(appCfg.uploadPict + flag + `/` + yyyymm + `/` + name)
	}
}
