package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func loginHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `login`)

	cmd := c.DefaultPostForm("cmd", "")

	checkMap := make(map[string]interface{})
	checkMap[`cmd`] = map[string]string{
		`type`: `:ss`,
	}
	if cmd == `login` {
		checkMap[`adm_id`] = map[string]string{
			`type`: `ss`,
			`name`: `用户名`,
		}
		checkMap[`adm_pwd`] = map[string]string{
			`type`: `ss`,
			`name`: `密码`,
		}
	} else {
		checkMap[`adm_id`] = map[string]string{
			`type`: `:ss`,
			`name`: `用户名`,
		}
		checkMap[`adm_pwd`] = map[string]string{
			`type`: `:s`,
			`name`: `密码`,
		}
	}

	e, arrData := ValidForm(c, checkMap)
	if !CheckValidResult(e) {
		fmt.Println(`%+v\n`, e)
		c.HTML(
			http.StatusOK,
			"login.html",
			MakeTemplateMap(c, e, arrData, nil),
		)
		return
	}

	if cmd == `login` {
		if CheckAdminLogin(c, arrData[`adm_id`].(string), arrData[`adm_pwd`].(string)) {
			c.Redirect(http.StatusFound, "/index")
			return
		} else {
			e[`adm_id`] = `用户名或密码不正确！`
		}
	}

	fmt.Println(`enter login`)

	c.HTML(
		http.StatusOK,
		"login.html",
		MakeTemplateMap(c, e, arrData, nil),
	)
}
