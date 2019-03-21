package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"Yinghao/klib"
	"Yinghao/tbls"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CommonMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		data := gin.H{
			`app`: gin.H{
				"version":   appCfg.version,
				"title":     appCfg.title,
				"languages": gLanguages,
			},
			`cookie`: GetCookie(c),
			`session`: gin.H{
				`admin`: GetSessionMap(c, `admin`),
			},
			`common`: gin.H{
				`router`:           ``,
				`isHeaderSearch`:   false, //头部检索按钮是否显示
				`isHeaderListInfo`: false, //右上角的 共X页，XX条是否显示
				"rowCount":         1,
				"pageNo":           1,
				"pageCount":        1,
				"pageRow":          appCfg.pagerow,
				"pageArray":        klib.MakePageNoArray(1, 1),
			},
			`error`: gin.H{},
			`data`:  gin.H{}, //此处存放页面所有其他数据
		}
		c.Set(`global`, data)
		c.Next()
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		admin := GetSessionMap(c, `admin`)
		if admin == nil {

			fmt.Println(`session timeout`)

			//如果客户端记住密码，存在cookie中，则自动登录
			cookie := GetCookie(c)
			cookieAutoLogin := `0`
			if value, ok := cookie[`cookie_auto_login`]; ok {
				cookieAutoLogin = value
			}
			cookieAdmId := ``
			if value, ok := cookie[`cookie_adm_id`]; ok {
				cookieAdmId = value
			}
			cookieAdmPwd := ``
			if value, ok := cookie[`cookie_adm_pwd`]; ok {
				cookieAdmPwd = value
			}
			if cookieAutoLogin == `1` &&
				cookieAdmId != `` &&
				cookieAdmPwd != `` {
				if CheckAdminLogin(c, cookieAdmId, cookieAdmPwd) {
					c.Next()
					return
				}
			}

			//如果是ajax，则返回json
			cmd := c.DefaultPostForm("cmd", "")
			if strings.Index(cmd, `ajax_`) == 0 {
				c.JSON(http.StatusOK, gin.H{
					"ret": 9999,
					"msg": `session timeout`,
				})
			} else {
				c.Redirect(http.StatusFound, "/")
			}

			c.Abort()
			return
		} else {

			//此操作为了刷新时间，否则timeout
			session := sessions.Default(c)
			session.Set(`admin`, klib.MapToJson(admin))
			session.Save()

			//fmt.Println(`session ok`)

			c.Next()
			return
		}
	}
}

func CheckAdminLogin(c *gin.Context, id, pwd string) bool {
	admin := new(tbls.TAdmin)

	ok, data := admin.Login(db, id, pwd)
	if ok {
		m := make(map[string]interface{})
		m[`adm_no`] = data.AdmNo
		m[`adm_id`] = data.AdmId
		m[`adm_name`] = data.AdmName
		m[`adm_perm`] = data.AdmPerm
		m[`adm_update`] = data.AdmUpdate

		session := sessions.Default(c)
		session.Set(`admin`, klib.MapToJson(m)) //不能直接存储map
		session.Save()

		return true
	} else {
		return false
	}
}

func MakeTemplateMap(c *gin.Context, e map[string]string, arrData map[string]interface{}, data interface{}) gin.H {

	global := gin.H{}
	if value, ok := c.Get(`global`); ok {
		global = value.(gin.H)
	}

	global[`error`] = e
	global[`data`] = arrData

	if data != nil {
		common := global[`common`].(gin.H)

		v := reflect.ValueOf(data).Elem()
		field := v.FieldByName(`Tbls`)
		if field.IsValid() {
			tb := field.Interface().(tbls.Tbls)

			common[`rowCount`] = tb.RowCount
			common[`pageNo`] = tb.PageNo
			common[`pageCount`] = tb.PageCount
		}

		common["pageArray"] = klib.MakePageNoArray(common["pageNo"].(int), common["pageCount"].(int))
		global[`common`] = common
	}

	return global
}
