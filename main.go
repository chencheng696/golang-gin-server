package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"Yinghao/klib"

	"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gorilla/websocket"
)

type appConfig struct {
	title      string
	version    string
	webport    string
	dbhost     string
	dbuser     string
	dbpassword string
	dbname     string
	pagerow    int
	debug      bool

	uploadMaxSize int //单位M
	uploadTemp    string
	uploadPict    string
	downloadExcel string
}

var (
	appCfg appConfig
	db     *gorm.DB //数据库连接对象
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true // 取消ws跨域校验
	},
}

func main() {
	//必须要先声明defer，否则不能捕获到panic异常
	defer func() {
		if r := recover(); r != nil {
			err := errors.New(``)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknow panic")
			}
			klib.WriteLog(err.Error())
		}
	}()

	if initConfig() == false {
		klib.WriteLog(`读取配置文件失败!`)
		os.Exit(1)
	}

	if initDatabase() == false {
		klib.WriteLog(`初始化数据库失败!`)
		os.Exit(1)
	}
	defer db.Close()

	// Logging to a file.
	f, _ := os.Create(klib.GetAppPath() + "/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	//gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(gin.Logger())

	//store := cookie.NewStore([]byte("secret"))
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	store.Options(sessions.Options{
		//Domain:   "localhost",
		MaxAge: 3600, //60*60=3600s
		//Path:     "/",
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("mysession", store))

	//添加中间件
	r.Use(CommonMiddleWare())

	//自定义模板函数
	r.SetFuncMap(template.FuncMap{
		"formatAsDate":   tmplFormatAsDate,
		"Int64ToString":  tmplInt64ToString,
		"unescaped":      tmplUnescaped,
		"unescapedjs":    tmplUnescapedjs,
		"unescapedcss":   tmplUnescapedcss,
		"formatOverflow": tmplFormatOverflow,
		"formatPictUrl":  tmplFormatPictUrl,
		"mapValue":       tmplMapValue,
		"multAdd":        tmplMultAdd,     //加
		"multMinus":      tmplMultMinus,   //减
		"multTimes":      tmplMultTimes,   //乘
		"multDivided":    tmplMultDivided, //除
	})

	r.Static("/assets", klib.GetAppPath()+"/asset")

	r.LoadHTMLGlob(klib.GetAppPath() + "/tmpl/*")

	//home page - login
	r.GET("/", loginHandle)
	r.POST("/", loginHandle)

	r.GET("/logout", logoutHandle)
	r.POST("/logout", logoutHandle)

	r.GET("/index", AuthMiddleWare(), indexHandle)
	r.POST("/index", AuthMiddleWare(), indexHandle)

	r.GET("/admin", AuthMiddleWare(), adminHandle)
	r.POST("/admin", AuthMiddleWare(), adminHandle)

	r.GET("/item", AuthMiddleWare(), itemHandle)
	r.POST("/item", AuthMiddleWare(), itemHandle)

	r.GET("/news", AuthMiddleWare(), newsHandle)
	r.POST("/news", AuthMiddleWare(), newsHandle)

	r.GET("/itemclass", AuthMiddleWare(), itemclassHandle)
	r.POST("/itemclass", AuthMiddleWare(), itemclassHandle)

	r.GET("/newsclass", AuthMiddleWare(), newsClassHandle)
	r.POST("/newsclass", AuthMiddleWare(), newsClassHandle)

	//公司简介
	r.GET("/aboutus", AuthMiddleWare(), aboutusHandle)
	r.POST("/aboutus", AuthMiddleWare(), aboutusHandle)

	r.GET("/organizational", AuthMiddleWare(), ArchtureHandle)
	r.POST("/organizational", AuthMiddleWare(), ArchtureHandle)

	r.GET("/glories", AuthMiddleWare(), GloriesHandle)
	r.POST("/glories", AuthMiddleWare(), GloriesHandle)

	r.GET("/culture", AuthMiddleWare(), CultureHandle)
	r.POST("/culture", AuthMiddleWare(), CultureHandle)

	//人才招聘
	r.GET("/job", AuthMiddleWare(), jobHandle)
	r.POST("/job", AuthMiddleWare(), jobHandle)
	//访客记录
	r.GET("/visator", AuthMiddleWare(), visatorHandle)
	r.POST("/visator", AuthMiddleWare(), visatorHandle)
	//合作伙伴
	r.GET("/partners", AuthMiddleWare(), PartnersHandle)
	r.POST("/partners", AuthMiddleWare(), PartnersHandle)
	//网络推广
	r.GET("/promotion", AuthMiddleWare(), PromotionHandle)
	r.POST("/promotion", AuthMiddleWare(), PromotionHandle)

	/*
		r.GET("/business", AuthMiddleWare(), businessHandle)
		r.POST("/business", AuthMiddleWare(), businessHandle)

		r.GET("/permission", AuthMiddleWare(), permissionHandle)
		r.POST("/permission", AuthMiddleWare(), permissionHandle)

		r.GET("/version", AuthMiddleWare(), versionHandle)
		r.POST("/version", AuthMiddleWare(), versionHandle)

			r.POST("/api/version", apiVerHandle)
	*/

	r.POST("/editorupload", AuthMiddleWare(), editoruploadHandle)

	r.GET("/pict/:flag/:yyyymm/:name", pictHandle)

	r.GET("/ws", func(c *gin.Context) {
		WsHandler(c.Writer, c.Request)
	})
	r.POST("/ws", func(c *gin.Context) {
		WsHandler(c.Writer, c.Request)
	})

	r.Run(":" + appCfg.webport)
}
