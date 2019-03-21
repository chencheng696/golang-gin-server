package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func logoutHandle(c *gin.Context) {

	session := sessions.Default(c)
	session.Clear()

	c.Redirect(http.StatusFound, "/")
}
