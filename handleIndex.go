package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func indexHandle(c *gin.Context) {

	SetGinGlobal(c, []string{`common`, `router`}, `index`)

	fmt.Println(`[index]index`)

	c.HTML(
		http.StatusOK,
		"index.html",
		MakeTemplateMap(c, nil, nil, nil),
	)
}
